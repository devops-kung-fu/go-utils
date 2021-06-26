package util

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/afero"
)

func sanitizeExtractPath(filePath string, destination string) error {
	destpath := filepath.Join(destination, filePath)
	if !strings.HasPrefix(destpath, destination) {
		return fmt.Errorf("%s: illegal file path", filePath)
	}
	return nil
}

// DownloadFile - Download a file from a URL
func DownloadFile(fs afero.Fs, url string) (string, error) {
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer checkClose(response.Body)

	if response.StatusCode != 200 {
		return "", errors.New("Received non 200 response code")
	}
	thisUUID := getUUID()
	fileName := fmt.Sprintf("/tmp/%s.zip", thisUUID)
	file, err := fs.Create(fileName)
	if err != nil {
		return "", err
	}
	defer checkClose(file)

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return "", err
	}
	log.Printf("Downloaded %s as %s\n", url, fileName)
	return thisUUID, nil
}

// UnZip - Extracts a zip archive
func UnZip(source string, destination string) error {
	archive, err := zip.OpenReader(source)
	if err != nil {
		return err
	}
	defer checkClose(archive)

	for _, file := range archive.Reader.File {
		reader, err := file.Open()
		if err != nil {
			return err
		}
		defer checkClose(reader)

		path := filepath.Join(destination, file.Name)
		// Remove file if it already exists; no problem if it doesn't; other cases can error out below
		_ = os.Remove(path)
		// Create a directory at path, including parents
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
		// If file is _supposed_ to be a directory, we're done
		if file.FileInfo().IsDir() {
			continue
		}
		// otherwise, remove that directory (_not_ including parents)
		err = os.Remove(path)
		if err != nil {
			return err
		}
		err = sanitizeExtractPath(file.Name, destination)
		if err != nil {
			return err
		}
		// and create the actual file.  This ensures that the parent directories exist!
		// An archive may have a single file with a nested path, rather than a file for each parent dir
		writer, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer checkClose(writer)

		_, err = io.Copy(writer, reader)
		if err != nil {
			return err
		}
	}
	log.Printf("Extracted %s into %s", source, destination)
	return nil
}

// FindFiles - Recursively search for files matching a pattern.
func FindFiles(fs afero.Fs, root string, re string) ([]string, error) {
	libRegEx, e := regexp.Compile(re)
	if e != nil {
		return nil, e
	}
	var files []string
	e = afero.Walk(fs, root, func(filePath string, info os.FileInfo, err error) error {
		if err == nil && libRegEx.MatchString(info.Name()) {
			files = append(files, filePath)
		}
		return nil
	})
	if e != nil {
		return nil, e
	}
	return files, nil
}

func checkClose(v interface{}) {
	if d, ok := v.(io.ReadCloser); ok {
		_ = d.Close()
	} else if d, ok := v.(io.Closer); ok {
		_ = d.Close()
	} else if d, ok := v.(zip.ReadCloser); ok {
		_ = d.Close()
	} else if d, ok := v.(os.File); ok {
		_ = d.Close()
	}
}
