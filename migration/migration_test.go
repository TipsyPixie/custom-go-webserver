package migration

import (
	"github.com/TipsyPixie/custom-go-webserver/config"
	"github.com/TipsyPixie/custom-go-webserver/database"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestNewRevision(t *testing.T) {
	revisionDir, err := ioutil.TempDir("", "custom-go-webserver-migration")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer os.RemoveAll(revisionDir)

	upRevisionFilename, downRevisionFilename, err := NewRevision(revisionDir, "test")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	_, err = os.Stat(filepath.Join(revisionDir, upRevisionFilename))
	if err != nil {
		t.Error(err)
	}

	_, err = os.Stat(filepath.Join(revisionDir, downRevisionFilename))
	if err != nil {
		t.Error(err)
	}
}

func TestNew(t *testing.T) {
	revisionDir, err := ioutil.TempDir("", "custom-go-webserver-migration")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer os.RemoveAll(revisionDir)

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
	_, _, err = NewRevision(revisionDir, "former")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	_, _, err = NewRevision(revisionDir, "latter")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	db, err := database.Open(testConfig)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer db.Close()

	migration, err := New("file://"+revisionDir, testConfig.Database.Driver, testConfig.Database.DatabaseName, db)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer migration.Close()

	err = migration.Up()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}
