package config

import (
	"encoding/json"
	"github.com/heindrichpaul/cleanupUtil/cleanupTasks"
	"log"
	"os"
	"time"
)

//Config is a struct that holds all the information needed for the cleanup job.
type Config struct {
	ModulesToClean cleanupTasks.CleanupTasks `json:"cleanup"`
	LogLocation    string                    `json:"logfile"`
	Location       string                    `json:"location"`
	loc            *time.Location
}

//NewConfig creates a new config structure from the given filename.
func NewConfig(configFile string) (conf *Config) {
	c, err := os.Open(configFile)

	if err != nil {
		log.Fatalf("Could not read config file (%s)\n", configFile)

	}

	json.NewDecoder(c).Decode(&conf)

	conf.loc, err = time.LoadLocation(conf.Location)
	if err != nil {
		log.Fatalf("Could not load (%s) location", conf.Location)
	}

	return
}

func (z *Config) GetLocation() *time.Location {
	return z.loc
}
