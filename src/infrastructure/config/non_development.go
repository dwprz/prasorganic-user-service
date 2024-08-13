package config

import (
	"context"
	"encoding/base64"
	"os"
	"strings"

	"github.com/dwprz/prasorganic-user-service/src/common/log"
	vault "github.com/hashicorp/vault/api"
	"github.com/sirupsen/logrus"
)

func setUpForNonDevelopment(appStatus string) *Config {
	defaultConf := vault.DefaultConfig()
	defaultConf.Address = os.Getenv("PRASORGANIC_CONFIG_ADDRESS")

	client, err := vault.NewClient(defaultConf)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "config.setUpForNonDevelopment", "section": "vault.NewClient"}).Fatal(err)
	}

	client.SetToken(os.Getenv("PRASORGANIC_CONFIG_TOKEN"))

	mountPath := "prasorganic-secrets" + "-" + strings.ToLower(appStatus)

	userServiceSecrets, err := client.KVv2(mountPath).Get(context.Background(), "user-service")
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "config.setUpForNonDevelopment", "section": "KVv2.Get"}).Fatal(err)
	}

	apiGatewaySecrets, err := client.KVv2(mountPath).Get(context.Background(), "api-gateway")
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "config.setUpForNonDevelopment", "section": "KVv2.Get"}).Fatal(err)
	}

	jwtSecrets, err := client.KVv2(mountPath).Get(context.Background(), "jwt")
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "config.setUpForNonDevelopment", "section": "KVv2.Get"}).Fatal(err)
	}

	imageKitSecrets, err := client.KVv2(mountPath).Get(context.Background(), "imagekit")
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "config.setUpForNonDevelopment", "section": "KVv2.Get"}).Fatal(err)
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

	jwtPrivateKey := jwtSecrets.Data["PRIVATE_KEY"].(string)
	base64Byte, err := base64.StdEncoding.DecodeString(jwtPrivateKey)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "config.setUpForNonDevelopment", "section": "base64.StdEncoding.DecodeString"}).Fatal(err)
	}
	jwtPrivateKey = string(base64Byte)

	jwtPublicKey := jwtSecrets.Data["Public_KEY"].(string)
	base64Byte, err = base64.StdEncoding.DecodeString(jwtPublicKey)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "config.setUpForNonDevelopment", "section": "base64.StdEncoding.DecodeString"}).Fatal(err)
	}
	jwtPublicKey = string(base64Byte)

	jwtConf.PrivateKey = loadRSAPrivateKey(jwtPrivateKey)
	jwtConf.PublicKey = loadRSAPublicKey(jwtPublicKey)

	imageKitConf := new(imageKit)
	imageKitConf.Id = imageKitSecrets.Data["ID"].(string)
	imageKitConf.BaseUrl = imageKitSecrets.Data["BASE_URL"].(string)
	imageKitConf.PrivateKey = imageKitSecrets.Data["PRIVATE_KEY"].(string)
	imageKitConf.PublicKey = imageKitSecrets.Data["PUBLIC_KEY"].(string)

	return &Config{
		CurrentApp: currentAppConf,
		Postgres:   postgresConf,
		Redis:      redisConf,
		ApiGateway: apiGatewayConf,
		Jwt:        jwtConf,
		ImageKit:   imageKitConf,
	}
}
