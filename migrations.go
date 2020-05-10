package main

import (
    "fmt"
    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database"
    "github.com/golang-migrate/migrate/v4/database/postgres"
    "github.com/golang-migrate/migrate/v4/database/sqlite3"
    _ "github.com/golang-migrate/migrate/v4/source/file"
    _ "github.com/lib/pq"
    "os"
    "time"
)

type UnknownDriver struct{}

func (receiver UnknownDriver) Error() string {
    return "unknown database driver"
}

func getMigration(sqlDir string) (*migrate.Migrate, error) {
    config, err := getConfig()
    if err != nil {
        return nil, err
    }
    dbDriverName := config.Database.Driver

    db, err := getDb()
    if err != nil {
        return nil, err
    }

    var driver database.Driver
    switch dbDriverName {
    case "postgres":
        driver, err = postgres.WithInstance(db, &postgres.Config{})
    case "sqlite3":
        driver, err = sqlite3.WithInstance(db, &sqlite3.Config{})
    default:
        return nil, UnknownDriver{}
    }

    migration, err := migrate.NewWithDatabaseInstance("file://"+sqlDir, dbDriverName, driver)
    if err != nil {
        return nil, err
    }
    return migration, nil
}

func createFile(filename string) error {
    sqlFile, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer sqlFile.Close()

    return nil
}

func generateRevision(name string) ([]string, error) {
    config, err := getConfig()
    if err != nil {
        return nil, err
    }
    sqlDir := config.Application.SqlDirectory

    currentTime := time.Now().UTC().Format("20060102150405")
    var createdFilePaths []string

    upSqlPath := fmt.Sprintf("%s/%s_%s.up.sql", sqlDir, currentTime, name)
    err = createFile(upSqlPath)
    if err != nil {
        return createdFilePaths, err
    }
    createdFilePaths = append(createdFilePaths, upSqlPath)

    downSqlPath := fmt.Sprintf("%s/%s_%s.down.sql", sqlDir, currentTime, name)
    err = createFile(downSqlPath)
    if err != nil {
        return createdFilePaths, err
    }
    createdFilePaths = append(createdFilePaths, downSqlPath)

    return createdFilePaths, nil
}

func migrateUp() error {
    config, err := getConfig()
    if err != nil {
        return err
    }

    migration, err := getMigration(config.Application.SqlDirectory)
    if err != nil {
        return err
    }

    return migration.Up()
}

func migrateDown() error {
    config, err := getConfig()
    if err != nil {
        return err
    }

    migration, err := getMigration(config.Application.SqlDirectory)
    if err != nil {
        return err
    }

    return migration.Down()
}

func migrateBySteps(steps int) error {
    config, err := getConfig()
    if err != nil {
        return err
    }

    migration, err := getMigration(config.Application.SqlDirectory)
    if err != nil {
        return err
    }

    return migration.Steps(steps)
}

func getVersion() (uint, error) {
    config, err := getConfig()
    if err != nil {
        return 0, err
    }

    migration, err := getMigration(config.Application.SqlDirectory)
    if err != nil {
        return 0, err
    }

    version, _, err := migration.Version()
    if err != nil {
        return 0, err
    }
    return version, nil
}
