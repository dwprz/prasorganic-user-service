package config

import (
	"github.com/sirupsen/logrus"
)

type currentApp struct {
	RestfulAddress string
	GrpcPort       string
}

type apiGateway struct {
	BaseUrl           string
	BasicAuth         string
	BasicAuthUsername string
	BasicAuthPassword string
}

type postgres struct {
	Url      string
	Dsn      string
	User     string
	Password string
}

type redis struct {
	AddrNode1 string
	AddrNode2 string
	AddrNode3 string
	AddrNode4 string
	AddrNode5 string
	AddrNode6 string
	Password  string
}

type Config struct {
	CurrentApp *currentApp
	Postgres   *postgres
	Redis      *redis
	ApiGateway *apiGateway
}

func New(appStatus string, logger *logrus.Logger) *Config {
	var config *Config

	if appStatus == "DEVELOPMENT" {

		config = setUpForDevelopment(logger)
		return config
	}

	config = setUpForNonDevelopment(appStatus, logger)
	return config
}
