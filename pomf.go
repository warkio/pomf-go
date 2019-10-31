package pomf

import (
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/warkio/logger-go"
)

type ResponseFormat int

const (
	ResponseFormatJSON ResponseFormat = iota + 1
	ResponseFormatText
)

const multipartFormMaxMemory = 1 * 1024 * 1024

var _ SelfChecker = (*Pomf)(nil)

type Pomf struct {
	Logger  logger.Logger
	Storage Storage
	Filter  Filter
}

func (p *Pomf) SelfCheck() error {
	var err error

	if p.Storage == nil {
		return fmt.Errorf("nil Storage")
	}

	err = p.Storage.SelfCheck()
	if err != nil {
		return err
	}

	return nil
}

func (p *Pomf) logger() logger.Logger {
	if p.Logger == nil {
		return logger.Discard
	}

	return p.Logger
}

func (p *Pomf) filter() Filter {
	if p.Filter == nil {
		return DefaultFilter
	}

	return p.Filter
}

func (p *Pomf) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log := p.logger()
	reqPath := req.URL.Path

	var err error

	switch reqPath {
	case "/upload":
		fallthrough
	case "/upload.json":
		err = p.handleUpload(ResponseFormatJSON, w, req)
	case "/upload.txt":
		err = p.handleUpload(ResponseFormatText, w, req)
	default:
		err = writeNotFound(w)
	}

	if err != nil {
		log.Print(err)
	}
}

func (p *Pomf) handleUpload(responseFormat ResponseFormat, w http.ResponseWriter, req *http.Request) error {
	if req.Method != http.MethodPost {
		return writeMethodNotAllowed(w)
	}

	var err error

	err = req.ParseMultipartForm(multipartFormMaxMemory)
	if err != nil {
		return err
	}

	mf := req.MultipartForm

	var results []UploadResult

	filter := p.filter()

	for _, fileHeader := range mf.File["files[]"] {
		err = func(fh *multipart.FileHeader) error {
			var (
				err error

				f multipart.File
			)

			f, err = fh.Open()
			if err != nil {
				return err
			}
			defer f.Close()

			name := fh.Filename
			size := fh.Size

			var isAllowed bool

			isAllowed, err = filter.IsAllowed(name, size, f)
			if err != nil {
				return err
			}
			if !isAllowed {
				return ErrNotAllowed
			}

			var ur *UploadResult

			ur, err = p.Storage.Upload(name, size, f)
			if err != nil {
				return err
			}

			results = append(results, *ur)

			return nil
		}(fileHeader)
	}

	ur := combineResults(results)

	err = ur.WriteToResponse(w)
	if err != nil {
		return err
	}

	return nil
}
