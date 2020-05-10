package database

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestOpen(t *testing.T) {
	dbDir, err := ioutil.TempDir("", "custom-go-webserver-database")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer os.RemoveAll(dbDir)

	dbFile, err := ioutil.TempFile(dbDir, "tmp")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	testConfig := Config{
		Driver:             "sqlite3",
		Hostname:           dbFile.Name(),
		MaxOpenConnections: 10,
		MaxIdleConnections: 20,
		ConnectionLifetime: 30,
		EngineOptions: map[string]string{
			"sslmode": "disable",
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
