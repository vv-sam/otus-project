package main

import (
	"encoding/json"
	"flag"
	"io"
	"os"
	"path"

	"github.com/vv-sam/otus-project/server/internal/configuration"
)

var (
	dataPath = flag.String("data_path", "/data", "path to service data folder")
	confPath = path.Join(*dataPath, "config.json")
)

func main() {
	confFile, err := os.Open(confPath)
	if err != nil {
		panic(err)
	}

	rawConfig, err := io.ReadAll(confFile)
	if err != nil {
		panic(err)
	}

	err = confFile.Close()
	if err != nil {
		panic(err)
	}

	var config configuration.Config
	err = json.Unmarshal(rawConfig, &config)
	if err != nil {
		panic(err)
	}

	// хочу сохранять конфиг при изменении в рантайме
	config.HandleConfigChanged(func(conf *configuration.Config) error {
		data, err := json.Marshal(conf)
		if err != nil {
			return err
		}

		err = os.WriteFile(confPath, data, os.ModeAppend)
		return err
	})

}
