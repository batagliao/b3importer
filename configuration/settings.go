package configuration

import (
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/joho/godotenv"
)

var instance *settings

const (
	MYSQL_ADDRESS = "MYSQL_ADDRESS"
	YEAR_START    = "YEAR_START"
)

type settings struct {
	MysqlAddress string
	FullImport   bool
	YearStart    int
}

func loadSettings() *settings {
	err := godotenv.Load()
	if err != nil {
		log.Debug().Err(err)
	}
	value := &settings{}

	env_mysql_address, exists := os.LookupEnv(MYSQL_ADDRESS)
	if exists {
		value.MysqlAddress = env_mysql_address
	}

	env_year_start, exists := os.LookupEnv(YEAR_START)
	t := time.Now()
	if exists && env_year_start != "" {
		year, err := strconv.Atoi(env_year_start)
		if err != nil {
			log.Warn().Msgf("Invalid value for %s variable: %s", env_year_start, env_year_start)
			value.YearStart = t.Year()
		}
		value.YearStart = year
	} else {
		value.YearStart = t.Year()
	}

	return value
}

func Get() *settings {
	if instance == nil {
		instance = loadSettings()
	}
	return instance
}
