package config

import (
	"context"
	vault "github.com/hashicorp/vault/api"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"strings"
)

func setUpForNonDevelopment(appStatus string, logger *logrus.Logger) *Config {
	defaultConf := vault.DefaultConfig()
	defaultConf.Address = os.Getenv("PRASORGANIC_CONFIG_ADDRESS")

	client, err := vault.NewClient(defaultConf)
	if err != nil {
		log.Fatalf("vault new client: %v", err)
	}

	client.SetToken(os.Getenv("PRASORGANIC_CONFIG_TOKEN"))

	mountPath := "prasorganic-secrets" + "-" + strings.ToLower(appStatus)

	userServiceSecrets, err := client.KVv2(mountPath).Get(context.Background(), "user-service")
	if err != nil {
		logger.WithFields(logrus.Fields{"location": "config.setUpForNonDevelopment", "section": "KVv2.Get"}).Fatal(err)
	}

	apiGatewaySecrets, err := client.KVv2(mountPath).Get(context.Background(), "api-gateway")
	if err != nil {
		logger.WithFields(logrus.Fields{"location": "config.setUpForNonDevelopment", "section": "KVv2.Get"}).Fatal(err)
	}

	jwtSecrets, err := client.KVv2(mountPath).Get(context.Background(), "jwt")
	if err != nil {
		logger.WithFields(logrus.Fields{"location": "config.setUpForNonDevelopment", "section": "KVv2.Get"}).Fatal(err)
	}

	currentAppConf := new(currentApp)
	currentAppConf.RestfulAddress = userServiceSecrets.Data["RESTFUL_ADDRESS"].(string)
	currentAppConf.GrpcPort = userServiceSecrets.Data["GRPC_PORT"].(string)

	postgresConf := new(postgres)
	postgresConf.Url = userServiceSecrets.Data["POSTGRES_URL"].(string)
	postgresConf.Dsn = userServiceSecrets.Data["POSTGRES_DSN"].(string)
	postgresConf.User = userServiceSecrets.Data["POSTGRES_USER"].(string)
	postgresConf.Password = userServiceSecrets.Data["POSTGRES_PASSWORD"].(string)

	redisConf := new(redis)
	redisConf.AddrNode1 = userServiceSecrets.Data["REDIS_ADDR_NODE_1"].(string)
	redisConf.AddrNode2 = userServiceSecrets.Data["REDIS_ADDR_NODE_2"].(string)
	redisConf.AddrNode3 = userServiceSecrets.Data["REDIS_ADDR_NODE_3"].(string)
	redisConf.AddrNode4 = userServiceSecrets.Data["REDIS_ADDR_NODE_4"].(string)
	redisConf.AddrNode5 = userServiceSecrets.Data["REDIS_ADDR_NODE_5"].(string)
	redisConf.AddrNode6 = userServiceSecrets.Data["REDIS_ADDR_NODE_6"].(string)
	redisConf.Password = userServiceSecrets.Data["REDIS_PASSWORD"].(string)

	apiGatewayConf := new(apiGateway)
	apiGatewayConf.BaseUrl = apiGatewaySecrets.Data["BASE_URL"].(string)
	apiGatewayConf.BasicAuth = apiGatewaySecrets.Data["BASIC_AUTH"].(string)
	apiGatewayConf.BasicAuthUsername = apiGatewaySecrets.Data["BASIC_AUTH_PASSWORD"].(string)
	apiGatewayConf.BasicAuthPassword = apiGatewaySecrets.Data["BASIC_AUTH_USERNAME"].(string)

	jwtConf := new(jwt)
	jwtConf.PrivateKey = loadRSAPrivateKey(jwtSecrets.Data["PRIVATE_KEY"].(string), logger)
	jwtConf.PublicKey = loadRSAPublicKey(jwtSecrets.Data["PUBLIC_KEY"].(string), logger)

	return &Config{
		CurrentApp: currentAppConf,
		Postgres:   postgresConf,
		Redis:      redisConf,
		ApiGateway: apiGatewayConf,
		Jwt:        jwtConf,
	}
}
