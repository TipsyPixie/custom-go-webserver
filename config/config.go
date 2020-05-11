package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"path"
)

type Config struct {
	Env         string
	Application struct {
		Secret string
		Debug  bool
	}
	Database struct {
		Driver             string
		Username           string
		Password           string
		Hostname           string
		Port               string
		DatabaseName       string            `yaml:"databaseName"`
		MaxOpenConnections int               `yaml:"maxOpenConnections"`
		MaxIdleConnections int               `yaml:"maxIdleConnections"`
		ConnectionLifetime int               `yaml:"connectionLifeTime"`
		EngineOptions      map[string]string `yaml:"engineOptions,omitempty"`
	}
	Migration struct {
		RevisionUrl string `yaml:"revisionUrl"`
	}
}

func Load(directory string, env string) (*Config, error) {
	configFileContent, err := ioutil.ReadFile(path.Join(directory, env))
	if err != nil {
		return nil, err
	}

	var newConfig Config
	err = yaml.Unmarshal(configFileContent, &newConfig)
	if err != nil {
		return nil, err
	}
	newConfig.Env = env
	return &newConfig, nil
}
