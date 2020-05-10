package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestGet(t *testing.T) {
	configDir, err := ioutil.TempDir("", "custom-go-webserver-test")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer os.RemoveAll(configDir)

	configFileContent := []byte(`
application:
  secret: secretValue
  migrationDir: Morikawa Yuki
database:
  engineOptions:
    option1: value1
`)
	configFile, err := ioutil.TempFile(configDir, "tmp")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if _, err := configFile.Write(configFileContent); err != nil {
		t.Error(err)
		t.FailNow()
	}

	config, err := Get(configDir, filepath.Base(configFile.Name()))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if env := config.Env; env != filepath.Base(configFile.Name()) {
		t.Error(env + "!=" + filepath.Base(configFile.Name()))
	}
	if migrationDir := config.Application.MigrationDir; migrationDir != "Morikawa Yuki" {
		t.Error(migrationDir + "!=" + "Morikawa Yuki")
	}
	if option1 := config.Database.EngineOptions["option1"]; option1 != "value1" {
		t.Error(option1 + "!+" + "value1")
	}
}
