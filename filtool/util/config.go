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

func ParseConfig(cfg_path string) (*Config, error) {
	if cfg_path == "" {
		cfg_path = "config.json"
	}

	f, err := os.Open(cfg_path)
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
