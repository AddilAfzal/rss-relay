package parser

import (
	"github.com/addilafzal/rss-relay/internal/rss"
	"gopkg.in/yaml.v3"
)

type Data struct {
	Source []rss.Source `yaml:"source"`
}

func ParseConfigFile(raw []byte) (*Data, error) {
	var dataObj Data
	if err := yaml.Unmarshal(raw, &dataObj); err != nil {
		return nil, err
	}
	return &dataObj, nil
}
