package database

import (
	"github.com/TipsyPixie/custom-go-webserver/config"
	"testing"
)

func TestOpen(t *testing.T) {
	testConfig := config.Config{
		Env: "test",
		Database: struct {
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
		}{
			Driver:             "sqlite3",
			Hostname:           "file::memory:?cache=shared",
			MaxOpenConnections: 10,
			MaxIdleConnections: 20,
			ConnectionLifetime: 30,
			EngineOptions: map[string]string{
				"sslmode": "disable",
			},
		},
	}

	db, err := Open(testConfig)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}
