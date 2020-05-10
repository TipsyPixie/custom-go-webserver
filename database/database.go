package database

import (
	"bytes"
	"custom-go-webserver/config"
	"database/sql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"strings"
	"time"
)

func getUrlString(connectionConfig config.Config) string {
	urlBuffer := bytes.Buffer{}
	if connectionConfig.Database.Username != "" || connectionConfig.Database.Password != "" {
		urlBuffer.WriteString(connectionConfig.Database.Username + ":" + connectionConfig.Database.Password + "@")
	}
	urlBuffer.WriteString(connectionConfig.Database.Hostname)
	if connectionConfig.Database.Port != "" {
		urlBuffer.WriteString(":" + connectionConfig.Database.Port)
	}
	if connectionConfig.Database.Driver != "sqlite3" {
		urlBuffer.WriteString("/" + connectionConfig.Database.DatabaseName)
		if len(connectionConfig.Database.EngineOptions) > 0 {
			optionPairs := make([]string, 0, len(connectionConfig.Database.EngineOptions))
			for key, value := range connectionConfig.Database.EngineOptions {
				optionPairs = append(optionPairs, key+"="+value)
			}
			urlBuffer.WriteString("?" + strings.Join(optionPairs, "&"))
		}
	}

	return urlBuffer.String()
}

func setConnectionPool(db *sql.DB, connectionConfig config.Config) {
	if connectionConfig.Database.MaxIdleConnections > 0 {
		db.SetMaxOpenConns(connectionConfig.Database.MaxOpenConnections)
	}
	if connectionConfig.Database.MaxOpenConnections > 0 {
		db.SetMaxIdleConns(connectionConfig.Database.MaxIdleConnections)
	}
	if connectionConfig.Database.ConnectionLifetime > 0 {
		db.SetConnMaxLifetime(time.Second * time.Duration(connectionConfig.Database.ConnectionLifetime))
	}
}

func Open(connectionConfig config.Config) (*sql.DB, error) {
	dataSourceUrl := getUrlString(connectionConfig)
	db, err := sql.Open(connectionConfig.Database.Driver, dataSourceUrl)
	if err != nil {
		return nil, err
	}
	setConnectionPool(db, connectionConfig)
	return db, nil
}
