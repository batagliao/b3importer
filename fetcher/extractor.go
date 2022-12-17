package fetcher

import (
	"archive/zip"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
)

type Extractor struct {
}

func NewExtractor() *Extractor {
	return &Extractor{}
}

func (me *Extractor) UnzipFile(zipfile string) ([]string, error) {
	log.Debug().Msgf("Extracting file %s", zipfile)
	zipReader, err := zip.OpenReader(zipfile)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := zipReader.Close(); err != nil {
			panic(err)
		}
	}()

	// build destination dir
	dest := strings.TrimSuffix(zipfile, filepath.Ext(zipfile))

	// create destination dir
	err = os.MkdirAll(dest, 0755)
	if err != nil {
		return nil, err
	}

	extractedFiles := []string{}
	for _, f := range zipReader.File {
		extracted, err := me.extractItem(dest, f)
		if err != nil {
			return nil, err
		}

		if extracted != "" {
			extractedFiles = append(extractedFiles, extracted)
		}
	}

	return extractedFiles, nil
}

func (me *Extractor) extractItem(destinationDir string, file *zip.File) (string, error) {
	rc, err := file.Open()
	if err != nil {
		return "", nil
	}

	defer func() {
		if err := rc.Close(); err != nil {
			panic(err)
		}
	}()

	targetPath := path.Join(destinationDir, file.Name)
	if file.FileInfo().IsDir() {
		err := os.MkdirAll(targetPath, file.Mode())
		if err != nil {
			return "", nil
		}
		return "", nil
	}

	err = os.MkdirAll(filepath.Dir(targetPath), file.Mode())
	if err != nil {
		return "", nil
	}

	targetFile, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
	if err != nil {
		return "", err
	}

	defer func() {
		if err := targetFile.Close(); err != nil {
			panic(err)
		}
	}()

	_, err = io.Copy(targetFile, rc)
	if err != nil {
		return "", err
	}

	return targetPath, nil

}
