package fetcher

import (
	"b3importer/filesystem"
	"errors"
	"fmt"
	"net/url"
	"path"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/cavaliergopher/grab/v3"
)

const (
	BASE_ZIP_FILE     = "dfp_cia_aberta_%d.zip"
	BASE_DOWNLOAD_URL = "https://dados.cvm.gov.br/dados/CIA_ABERTA/DOC/DFP/DADOS/"
	DOWNLOAD_DIR      = "./downloads"
	FIRST_YEAR        = 2010
)

type Downloader struct {
	client           grab.Client
	currentDownloads []*grab.Response
}

func NewDownloader() *Downloader {
	return &Downloader{
		client:           *grab.NewClient(),
		currentDownloads: []*grab.Response{},
	}
}

func (me *Downloader) download_main_zip(year int) string {
	filename := fmt.Sprintf(BASE_ZIP_FILE, year)
	target_file := path.Join(DOWNLOAD_DIR, filename)
	source_url, err := url.JoinPath(BASE_DOWNLOAD_URL, filename)

	if err != nil {
		log.Error().Err(err)
		return ""
	}

	req, err := grab.NewRequest(target_file, source_url)
	if err != nil {
		log.Error().Err(err)
		return ""
	}

	req.BeforeCopy = func(r *grab.Response) error {
		log.Debug().Msgf("starting download of %s", filename)
		return nil
	}

	req.AfterCopy = func(r *grab.Response) error {
		log.Debug().Msgf("%s download successfull", filename)
		return nil
	}

	req.NoResume = true

	if filesystem.FileExists(target_file) {
		// if the file already exists we will delete it to download again
		// this way is safer as the file might change in the server
		log.Debug().Msgf("removing existing file %s", filename)
		err := filesystem.DeleteFile(target_file)
		if err != nil {
			log.Error().Err(err)
		}
	}

	response := me.client.Do(req)
	me.currentDownloads = append(me.currentDownloads, response)
	if err := response.Err(); err != nil {
		log.Error().Err(err)
	}

	return target_file
}

func (me *Downloader) DownloadAll(startYear int) ([]string, error) {
	if startYear < FIRST_YEAR {
		return nil, errors.New("invalid start year")
	}

	files := []string{}
	currentYear := time.Now().Year()

	for y := startYear; y <= currentYear; y++ {
		file := me.download_main_zip(y)
		files = append(files, file)
	}

	// wait for all downloads to complete
	log.Info().Msg("Waiting all downloads to complete")
	for _, r := range me.currentDownloads {
		r.Wait()
	}

	return files, nil
}
