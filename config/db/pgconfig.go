package pgconfig

import (
	"fmt"
	"os"
	"strings"
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

func NewDbSettings(envPrefix string) *dbSettings {
	envPrefix = strings.ToUpper(envPrefix)
	envName := func(env string) string { return envPrefix + env }
	sslMode, ok := os.LookupEnv(envName("DB_SSLMODE"))
	if !ok {
		sslMode = "disable"
	}
	timeZone, ok := os.LookupEnv(envName("DB_TIMEZONE"))
	if !ok {
		timeZone = "Etc/GMT+8"
	}
	return &dbSettings{
		hostName: os.Getenv(envName("DB_HOST")),
		user:     os.Getenv(envName("DB_USER")),
		password: os.Getenv(envName("DB_PASSWORD")),
		dbName:   os.Getenv(envName("DB_NAME")),
		port:     os.Getenv(envName("DB_PORT")),
		sslMode:  sslMode,
		timeZone: timeZone}
}

func (settings *dbSettings) AsString() string {
	dbSettingsString := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		settings.hostName,
		settings.user,
		settings.password,
		settings.dbName,
		settings.port,
		settings.sslMode,
		settings.timeZone,
	)
	return dbSettingsString
}

func (settings *dbSettings) AsUrl() string {
	dbSettingsUrl := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s&timezone=%s",
		settings.user,
		settings.password,
		settings.hostName,
		settings.port,
		settings.dbName,
		settings.sslMode,
		settings.timeZone,
	)
	return dbSettingsUrl
}
