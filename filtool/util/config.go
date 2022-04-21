package util

// TODO: add toString function

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Config struct {
	Fullnode     ApiInfo `json:"fullnode"`
	StorageMiner ApiInfo `json:"storageminer"`
	MarketMiner  ApiInfo `json:"marketminer"`
}

type ApiInfo struct {
	Token string `json:"token"`
	Node  string `json:"node"`
}

func GetConfig() (*Config, error) {
	cfgPath := os.Getenv("CONFIG_PATH")
	if cfgPath == "" {
		cfgPath = "config.json"
	}

	f, err := os.Open(cfgPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	byteValue, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var cfg Config
	json.Unmarshal(byteValue, &cfg)

	return &cfg, nil
}
