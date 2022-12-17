package main

import (
	"b3importer/configuration"
	"b3importer/fetcher"
	"b3importer/importer"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	log.Logger = log.With().Caller().Logger().Output(zerolog.ConsoleWriter{
		Out: os.Stderr,
	})
}

func main() {
	settings := configuration.Get()

	// Download the files
	downloader := fetcher.NewDownloader()
	files, err := downloader.DownloadAll(settings.YearStart)
	if err != nil {
		log.Fatal().Err(err)
	}

	// Extract zip
	extractor := fetcher.NewExtractor()
	extractedFiles := []string{}
	for _, f := range files {
		extracted, err := extractor.UnzipFile(f)
		if err != nil {
			log.Fatal().Err(err)
		}
		extractedFiles = append(extractedFiles, extracted...)
	}

	dataImporter := importer.NewDataImporter()
	for _, f := range extractedFiles {
		err := dataImporter.ImportFile(f)
		if err != nil {
			log.Fatal().Err(err)
		}
	}
}
