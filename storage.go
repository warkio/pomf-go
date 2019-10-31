package pomf

import (
	"encoding/json"
	"io"
	"net/http"
)

type Storage interface {
	SelfChecker

	Upload(originalName string, size int64, r io.Reader) (*UploadResult, error)
	SetURLPrefix(prefix string) error
}

type UploadResultFile struct {
	Hash string `json:"hash"`
	Name string `json:"name"`
	URL  string `json:"url"`
	Size int64  `json:"size"`
}

type UploadResult struct {
	Success bool `json:"success"`

	Files []UploadResultFile `json:"files,omitempty"`

	ErrorCode   int    `json:"errorcode,omitempty"`
	Description string `json:"description,omitempty"`
}

func (ur *UploadResult) WriteToResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	enc := json.NewEncoder(w)
	err := enc.Encode(ur)

	return err
}

func combineResults(results []UploadResult) UploadResult {
	var files []UploadResultFile
	for _, result := range results {
		if !result.Success {
			return result
		}

		files = append(files, result.Files...)
	}

	ur := UploadResult{
		Success: true,
		Files:   files,
	}

	return ur
}
