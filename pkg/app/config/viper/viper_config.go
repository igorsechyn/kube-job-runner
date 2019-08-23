package viper

import (
	"strings"

	"kube-job-runner/pkg/app/config"

	"github.com/spf13/viper"
)

func LoadConfig() config.Config {
	viper.AutomaticEnv()
	appConfig := config.Config{}
	appConfig.PgDatabase = viper.GetString("PG_DB_SCHEMA")
	appConfig.PgHost = viper.GetString("PG_DB_HOST")
	appConfig.PgPassword = viper.GetString("PG_DB_PASSWORD")
	appConfig.PgUser = viper.GetString("PG_DB_ROLE")
	appConfig.PgPort = viper.GetInt("PG_DB_PORT")
	appConfig.SQSWaitTimeSeconds = viper.GetInt("SQS_WAIT_TIME_SECONDS")
	appConfig.Environment = viper.GetString("MICROS_ENV")
	appConfig.SQSQueueName = viper.GetString("SQS_Q_QUEUE_NAME")
	appConfig.SQSQueueRegion = viper.GetString("SQS_Q_QUEUE_REGION")
	appConfig.SQSQueueName = viper.GetString("SQS_Q_QUEUE_NAME")
	appConfig.SQSQueueUrl = viper.GetString("SQS_Q_QUEUE_URL")
	appConfig.Namespace = getEnvOrDefault("NAMESPACE", "default")
	appConfig.Components = getComponents("COMPONENTS")
	return appConfig
}

func getEnvOrDefault(key string, defaultValue string) string {
	if viper.IsSet(key) {
		return viper.GetString(key)
	}

	return defaultValue
}

func getComponents(key string) []string {
	componentsConfig := getEnvOrDefault(key, "")
	if componentsConfig == "" {
		return []string{}
	}

	return strings.Split(componentsConfig, ",")
}
