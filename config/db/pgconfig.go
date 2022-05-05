package pgconfig

import (
	"fmt"
	"os"
)

type dbSettings struct {
	hostName string
	user     string
	password string
	dbName   string
	port     string
	sslMode  string
	timeZone string
}

func newDbSettings() *dbSettings {
	sslMode, ok := os.LookupEnv("DB_SSLMODE")
	if !ok {
		sslMode = "disable"
	}
	timeZone, ok := os.LookupEnv("DB_TIMEZONE")
	if !ok {
		timeZone = "Etc/GMT+8"
	}
	return &dbSettings{
		hostName: os.Getenv("DB_HOST"),
		user:     os.Getenv("DB_USER"),
		password: os.Getenv("DB_PASSWORD"),
		dbName:   os.Getenv("DB_NAME"),
		port:     os.Getenv("DB_PORT"),
		sslMode:  sslMode,
		timeZone: timeZone}
}

func DBSettingsAsString() string {
	dbSettings := newDbSettings()
	dbSettingsString := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		dbSettings.hostName,
		dbSettings.user,
		dbSettings.password,
		dbSettings.dbName,
		dbSettings.port,
		dbSettings.sslMode,
		dbSettings.timeZone,
	)
	return dbSettingsString
}
