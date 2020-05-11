package migration

import (
	"database/sql"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func formatFileName(revisionName string) (string, string) {
	now := strconv.FormatInt(time.Now().UnixNano(), 10)
	base := now + "_" + revisionName
	return base + ".up.sql", base + ".down.sql"
}

//func formatRevisionDirectory(revisionDir string) (string, error) {
//	colonSlashIndex := strings.Index(revisionDir, "://")
//	if colonSlashIndex < 0 {
//		revisionDir = "file://" + revisionDir
//	} else if schema := revisionDir[:colonSlashIndex]; schema != "file" {
//		return "", errors.New(revisionDir + "is not a filesystem directory")
//	}
//	return revisionDir, nil
//}
func NewRevision(revisionDir string, revisionName string) (string, string, error) {
	upFileName, downFileName := formatFileName(revisionName)
	upRevisionPath := filepath.Join(revisionDir, upFileName)
	upFile, err := os.Create(upRevisionPath)
	if err != nil {
		return "", "", err
	}
	defer upFile.Close()

	downRevisionPath := filepath.Join(revisionDir, downFileName)
	downFile, err := os.Create(downRevisionPath)
	if err != nil {
		return "", "", err
	}
	defer downFile.Close()

	return upFileName, downFileName, nil
}

func getDriver(driverName string, db *sql.DB) (database.Driver, error) {
	switch driverName {
	case "postgres":
		return postgres.WithInstance(db, &postgres.Config{})
	case "sqlite3":
		return sqlite3.WithInstance(db, &sqlite3.Config{})
	default:
		return nil, errors.New("does not support " + driverName)
	}
}

func New(revisionUrl string, driverName string, databaseName string, db *sql.DB) (*migrate.Migrate, error) {
	driver, err := getDriver(driverName, db)
	if err != nil {
		return nil, err
	}

	migration, err := migrate.NewWithDatabaseInstance(revisionUrl, databaseName, driver)
	if err != nil {
		return nil, err
	}

	return migration, nil
}
