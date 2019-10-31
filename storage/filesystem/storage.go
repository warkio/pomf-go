package filesystem

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/warkio/pomf-go"
)

var _ pomf.Storage = (*Storage)(nil)

const (
	MinNameLength     = 4
	DefaultNameLength = 8
)

type Storage struct {
	prefix string

	NameLength      int
	UploadDirectory string
}

func (fs *Storage) SelfCheck() error {
	if fs.NameLength == 0 {
		fs.NameLength = DefaultNameLength
	}

	if fs.NameLength < MinNameLength {
		return fmt.Errorf("NameLength too small")
	}

	var (
		err error
		f   *os.File
	)

	f, err = os.Open(fs.UploadDirectory)
	if err != nil {
		return err
	}
	defer f.Close()

	var fi os.FileInfo

	fi, err = f.Stat()
	if err != nil {
		return err
	}

	if !fi.IsDir() {
		err = fmt.Errorf("not a directory: %s", fs.UploadDirectory)

		return err
	}

	return nil
}

func (fs *Storage) Upload(originalName string, size int64, r io.Reader) (*pomf.UploadResult, error) {
	ext := path.Ext(originalName)
	ext = strings.ToLower(ext)

	var (
		err  error
		name string
	)

	name, err = RandomName(fs.NameLength)
	if err != nil {
		return nil, err
	}

	fileName := name + ext

	file := path.Join(fs.UploadDirectory, fileName)

	var f *os.File

	f, err = os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0640)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	h := sha256.New()

	tee := io.TeeReader(r, f)

	_, err = io.Copy(h, tee)
	if err != nil {
		return nil, err
	}

	sum := h.Sum(nil)

	fileURL := fs.prefix + fileName

	ur := &pomf.UploadResult{
		Success: true,
		Files: []pomf.UploadResultFile{
			pomf.UploadResultFile{
				Name: originalName,
				Size: size,
				Hash: fmt.Sprintf("%x", sum),
				URL:  fileURL,
			},
		},
	}

	return ur, nil
}

func (fs *Storage) SetURLPrefix(prefix string) error {
	fs.prefix = prefix

	return nil
}
