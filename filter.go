package pomf

import (
	"errors"
	"io"
	"path"
)

var ErrNotAllowed = errors.New("not allowed")

const DefaultMaxFileSize = 50 * 1024 * 1024

type FilterFile interface {
	io.Reader
	io.ReaderAt
	io.Seeker
}

type Filter interface {
	IsAllowed(name string, size int64, f FilterFile) (bool, error)
}

var DefaultFilter Filter = &defaultFilter{}

type defaultFilter struct{}

func (df *defaultFilter) IsAllowed(name string, size int64, f FilterFile) (bool, error) {
	if size <= 0 || size > DefaultMaxFileSize {
		return false, nil
	}

	allowedExtensions := []string{
		".gif",
		".jpeg",
		".jpg",
		".png",

		".txt",
	}

	ext := path.Ext(name)
	for _, allowedExtension := range allowedExtensions {
		if ext == allowedExtension {
			return true, nil
		}
	}

	return false, nil
}
