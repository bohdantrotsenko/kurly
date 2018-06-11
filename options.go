package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/alsm/ioprogress"
	"github.com/urfave/cli"
)

type Field struct {
	Type      string
	IsFile    bool
	Value     string
	Filealias string
}

type FormData map[string]Field

type Options struct {
	outputFilename string
	fileUpload     string
	remoteName     bool
	continueAt     string
	verbose        bool
	maxTime        uint
	remoteTime     bool
	cookie         string
	cookieJar      string
	followRedirect bool
	maxRedirects   uint
	redirectsTaken uint
	silent         bool
	method         string
	headers        []string
	agent          string
	user           string
	expectTimeout  uint
	data           []string
	dataAscii      []string
	dataRaw        []string
	dataBinary     []string
	dataURLEncode  []string
	form           []string
	head           bool
	insecure       bool
	fdata          FormData // fdata is the field for processed form data
}

func (o *Options) getOptions(app *cli.App) {
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "output, o",
			Usage:       "Filename to name url content to",
			Destination: &o.outputFilename,
		},
		cli.StringFlag{
			Name:        "upload-file, T",
			Usage:       "File to upload",
			Destination: &o.fileUpload,
		},
		cli.BoolFlag{
			Name:        "remote-name, O",
			Usage:       "Save output to file named with file part of URL",
			Destination: &o.remoteName,
		},
		cli.StringFlag{
			Name:        "continue-at, C",
			Usage:       "Resume transfer from offset",
			Destination: &o.continueAt,
		},
		cli.BoolFlag{
			Name:        "v",
			Usage:       "Verbose output",
			Destination: &o.verbose,
		},
		cli.UintFlag{
			Name:        "max-time, m",
			Usage:       "Maximum time to wait for an operation to complete in seconds",
			Destination: &o.maxTime,
		},
		cli.BoolFlag{
			Name:        "R",
			Usage:       "Set the timestamp of the local file to that of the remote file, if available",
			Destination: &o.remoteTime,
		},
		cli.StringFlag{
			Name:        "cookie, b",
			Usage:       "Set the cookies to be sent along with this request",
			Destination: &o.cookie,
		},
		cli.StringFlag{
			Name:        "cookie-jar, c",
			Usage:       "File to which the cookies have to be written (in cURL's cookie-jar file format)",
			Destination: &o.cookieJar,
		},
		cli.BoolFlag{
			Name:        "location, L",
			Usage:       "Follow 3xx redirects",
			Destination: &o.followRedirect,
		},
		cli.UintFlag{
			Name:        "max-redirs",
			Usage:       "Maximum number of 3xx redirects to follow",
			Destination: &o.maxRedirects,
			Value:       10,
		},
		cli.BoolFlag{
			Name:        "silent, s",
			Usage:       "Mute kurly entirely, operation without any output",
			Destination: &o.silent,
		},
		cli.StringFlag{
			Name:        "request, X",
			Usage:       "HTTP method to use",
			Destination: &o.method,
			Value:       "GET",
		},
		cli.StringFlag{
			Name:        "user-agent, A",
			Usage:       "User agent to set for this request",
			Destination: &o.agent,
			Value:       "Kurly/" + version,
		},
		cli.StringFlag{
			Name:        "user, u",
			Usage:       "User authentication data to set for this request",
			Destination: &o.user,
		},
		cli.StringSliceFlag{
			Name:  "header, H",
			Usage: "Extra headers to be sent with the request",
		},
		cli.UintFlag{
			Name:        "expect100-timeout",
			Usage:       "Timeout in seconds for Expect: 100-continue wait period",
			Destination: &o.expectTimeout,
			Value:       1,
		},
		cli.StringSliceFlag{
			Name:  "data, d",
			Usage: "Sends the specified data in a POST request to the server",
		},
		cli.StringSliceFlag{
			Name:  "data-ascii",
			Usage: "The same as --data, -d",
		},
		cli.StringSliceFlag{
			Name:  "data-raw",
			Usage: "Basically the same as --data-binary (no @ interpretation)",
		},
		cli.StringSliceFlag{
			Name:  "data-binary",
			Usage: "Sends the data as binary",
		},
		cli.StringSliceFlag{
			Name:  "data-urlencode",
			Usage: "Sends the data as urlencoded ascii",
		},
		cli.StringSliceFlag{
			Name:  "form, F",
			Usage: "Send HTTP multipart post data",
		},
		cli.BoolFlag{
			Name:        "head, I",
			Usage:       "Get HEAD from URL only",
			Destination: &o.head,
		},
		cli.BoolFlag{
			Name:        "insecure, k",
			Usage:       "Allow insecure server connections when using TLS",
			Destination: &o.insecure,
		},
	}
}

func (o *Options) checkRedirect(req *http.Request, via []*http.Request) error {
	o.redirectsTaken++

	resp := req.Response

	if !o.followRedirect || o.redirectsTaken >= o.maxRedirects {
		return http.ErrUseLastResponse
	}

	if resp != nil {
		fmt.Fprintf(Incoming, "%s %s\n", resp.Proto, resp.Status)

		for k, v := range resp.Header {
			fmt.Fprintln(Incoming, k, v)
		}
		fmt.Fprintln(Incoming)
	}

	if o.verbose {
		Status.Println(" Ignoring the response body")
		Status.Printf(" Issuing request to this URL : %s\n", req.URL)
	}
	return nil
}

// BuildCommonOptions function is used to build all the common options from the cli.Context
// which is common for all the URL args passed from commandline.
// This this process takes place only once.
func (opts *Options) BuildCommonOptions(c *cli.Context) error {
	opts.headers = c.StringSlice("header")
	opts.user = c.String("user")
	opts.dataAscii = c.StringSlice("data")
	opts.dataAscii = append(opts.dataAscii, c.StringSlice("data-ascii")...)
	opts.dataBinary = c.StringSlice("data-binary")
	opts.dataRaw = c.StringSlice("data-raw")
	opts.dataURLEncode = c.StringSlice("data-urlencode")
	opts.form = c.StringSlice("form")

	// If verbose set the logs writers
	if opts.verbose {
		Incoming.(*LogWriter).SetOutput(os.Stderr)
		Outgoing.(*LogWriter).SetOutput(os.Stderr)
	}

	// Process form data or url-encoded data
	opts.ProcessData()
	d, err := opts.ProcessFormData()
	if err != nil {
		return err
	}
	if d != nil && len(opts.data) > 0 {
		return fmt.Errorf("only one type of body can be accepted : either multipart form or url encoded values")
	}
	opts.fdata = d

	// Set the request method if Head option is specified
	if opts.head {
		opts.method = "HEAD"
		Incoming = io.MultiWriter(os.Stdout, Incoming.(*LogWriter))
	}

	// Start the timer
	if opts.maxTime > 0 {
		maxTime(opts.maxTime)
	}

	return nil
}

// BuildTargetSpecificOptions function is used to build the options specific to a given URL target.
// This has to run for every URL separately.
func (opts *Options) BuildTargetSpecificOptions(target string) (io.Reader, error) {
	var body io.Reader
	// Set the output filename from the remote URL, if -O is passed.
	if opts.remoteName {
		opts.outputFilename = path.Base(target)
	}

	// Initialize the file upload if specified.
	if opts.fileUpload != "" {
		body = opts.uploadFile()
	}

	// Process headers and post data
	if len(opts.data) > 0 || len(opts.fdata) > 0 {
		var data bytes.Buffer
		opts.method = "POST"

		header := ""

		if len(opts.data) > 0 {
			header = "Content-Type: application/x-www-form-urlencoded"
			for i, d := range opts.data {
				data.WriteString(d)
				if i < len(opts.data)-1 {
					data.WriteRune('&')
				}
			}
		}

		if len(opts.fdata) > 0 {
			w := multipart.NewWriter(&data)
			for key, field := range opts.fdata {
				err := writeToMultipart(w, key, field)
				if err != nil {
					return nil, fmt.Errorf("unable to create http request; %s\n", err)
				}
			}
			w.Close()
			header = "Content-Type: " + w.FormDataContentType()
		}

		opts.headers = append([]string{header}, opts.headers...)
		body = &data
	}

	return body, nil
}

func (o *Options) ProcessData() {
	var uriEncodes url.Values
	for _, d := range o.dataAscii {
		parts := strings.SplitN(d, "=", 2)
		if len(parts) == 1 {
			if strings.HasPrefix(parts[0], "@") {
				data, err := ioutil.ReadFile(strings.TrimPrefix(parts[0], "@"))
				if err != nil {
					Status.Fatalf("Unable to read file %s : %s", strings.TrimPrefix(parts[0], "@"), err)
				}
				o.data = append(o.data, string(data))
				continue
			}
			o.data = append(o.data, d)
			continue
		}
		if strings.HasPrefix(parts[1], "@") {
			data, err := ioutil.ReadFile(strings.TrimPrefix(parts[1], "@"))
			if err != nil {
				Status.Fatalf("Unable to read file %s for data element %s\n", strings.TrimPrefix(parts[1], "@"), parts[0])
			}
			data = []byte(strings.Replace(string(data), "\r", "", -1))
			data = []byte(strings.Replace(string(data), "\n", "", -1))
			o.data = append(o.data, fmt.Sprintf("%s=%s", parts[0], string(data)))
		} else {
			o.data = append(o.data, d)
		}
	}
	for _, d := range o.dataRaw {
		parts := strings.SplitN(d, "=", 2)
		o.data = append(o.data, fmt.Sprintf("%s=%s", parts[0], parts[1]))
	}
	for _, d := range o.dataBinary {
		parts := strings.SplitN(d, "=", 2)
		o.data = append(o.data, fmt.Sprintf("%s=%s", parts[0], parts[1]))
	}
	for _, d := range o.dataURLEncode {
		parts := strings.SplitN(d, "=", 2)
		uriEncodes.Add(parts[0], parts[1])
	}
	if len(uriEncodes) > 0 {
		o.data = append(o.data, uriEncodes.Encode())
	}
}

func (o *Options) openOutputFile() *os.File {
	var err error
	var outputFile *os.File
	if o.outputFilename != "" {
		if outputFile, err = os.OpenFile(o.outputFilename, os.O_CREATE|os.O_RDWR, 0666); err != nil {
			Status.Fatalf("Error: Unable to create/open file '%s' for output\n", o.outputFilename)
		}
	} else {
		outputFile = os.Stdout
	}
	return outputFile
}

func (o *Options) uploadFile() io.Reader {
	o.method = "PUT"

	tr := &http.Transport{
		ExpectContinueTimeout: time.Duration(o.expectTimeout) * time.Second,
	}
	client.Transport = tr
	o.headers = append(o.headers, "Expect: 100-continue")

	reader, err := os.Open(o.fileUpload)
	if err != nil {
		Status.Fatalf("Error opening %s\n", o.fileUpload)
	}

	if !o.silent {
		fi, err := reader.Stat()
		if err != nil {
			Status.Fatalf("Unable to get file stats for %v\n", o.fileUpload)
		}
		return &ioprogress.Reader{
			Reader: reader,
			Size:   fi.Size(),
			DrawFunc: ioprogress.DrawTerminalf(os.Stderr, func(progress, total int64) string {
				return fmt.Sprintf(
					"%s %s",
					(ioprogress.DrawTextFormatBarWithIndicator(40, '>'))(progress, total),
					ioprogress.DrawTextFormatBytes(progress, total))
			}),
		}
	}

	return reader
}

// ProcessFormData is used to parse the form data passed as commandline arguments.
func (o *Options) ProcessFormData() (FormData, error) {
	if len(o.form) == 0 {
		return nil, nil
	}

	fd := make(FormData, 0)

	for _, field := range o.form {
		fieldName, f, err := parseField(field)
		if err != nil {
			return nil, err
		}

		fd[fieldName] = f
	}
	return fd, nil
}

func parseField(raw string) (string, Field, error) {
	f := Field{}
	returnKey := ""
	parts := splitFormParams(raw)
	for i, part := range parts {
		key, value, err := getKeyVal(part)
		if err != nil {
			return "", f, err
		}
		switch key {
		case "type":
			f.Type = value
		case "filename":
			f.Filealias = value
		default:
			if returnKey != "" || i != 0 {
				return "", Field{}, errors.New("malformed form data")
			}
			returnKey = key
			if strings.HasPrefix(value, "@") && len(value) > 1 {
				f.IsFile = true
				f.Value = value[1:]
				continue
			}

			f.Value = value
			f.IsFile = false
		}
	}
	if f.IsFile && f.Filealias == "" {
		f.Filealias = filepath.Base(f.Value)
	}

	return returnKey, f, nil
}

func splitFormParams(raw string) []string {
	inQuotes := false
	temp := ""
	strList := make([]string, 0)
	for _, ch := range raw {
		t := string(ch)
		if t == "\"" {
			inQuotes = !inQuotes
			continue
		}

		if t == ";" && !inQuotes {
			if temp != "" {
				strList = append(strList, temp)
				temp = ""
			}
			continue
		}

		temp += t
	}

	if temp != "" {
		strList = append(strList, temp)
	}

	return strList
}

func getKeyVal(raw string) (string, string, error) {
	splits := strings.Split(raw, "=")
	if len(splits) < 2 {
		return "", "", errors.New("not a valid key-value pair")
	}
	key := splits[0]
	val := strings.Join(splits[1:], "=")
	return key, val, nil
}
