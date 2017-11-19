package main

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Match struct {
	IP     string `json:"ip"`
	Regexp string `json:"regexp"`
	Proxy  string `json:"proxy"`
	Direct bool   `json:"direct"`
}

type Config struct {
	Proxies map[string]string `json:"proxies"`
	Matches []Match           `json:"matches"`
}

func ParseConfig(data []byte) (*Config, error) {
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, errors.Wrap(err, "failed to parse config")
	}
	return &config, nil
}
