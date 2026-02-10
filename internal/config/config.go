package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	*viper.Viper
}

func NewConfig() *Config {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("json")
	v.AddConfigPath("./../../../")
	v.AddConfigPath("./../../")
	v.AddConfigPath("./../")
	v.AddConfigPath("./")

	config := &Config{Viper: v}

	err := config.ReadInConfig()
	if err != nil {
		panic(err)
	}

	return config
}

func (c *Config) GetAppName() string {
	return c.GetString(AppName)
}

func (c *Config) GetAuthJWTSecret() string {
	return c.GetString(AuthJWTSecret)
}

func (c *Config) SetAuthJWTSecret(value string) {
	c.Set(AuthJWTSecret, value)
}

func (c *Config) GetAuthJWTIssuer() string {
	return c.GetString(AuthJWTIssuer)
}

func (c *Config) SetAuthJWTIssuer(value string) {
	c.Set(AuthJWTIssuer, value)
}

func (c *Config) GetAuthJWTExpireSeconds() int {
	return c.GetInt(AuthJWTExpireSeconds)
}

func (c *Config) SetAuthJWTExpireSeconds(value int) {
	c.Set(AuthJWTExpireSeconds, value)
}

func (c *Config) GetAWSRegion() string {
	return c.GetString(AWSRegion)
}

func (c *Config) GetAWSBaseEndpoint() string {
	return c.GetString(AWSBaseEndpoint)
}

func (c *Config) GetDatabaseMigrations() string {
	return c.GetString(DatabaseMigrations)
}

func (c *Config) SetDatabaseMigrations(value string) {
	c.Set(DatabaseMigrations, value)
}

func (c *Config) GetDatabaseUsername() string {
	return c.GetString(DatabaseUsername)
}

func (c *Config) GetDatabasePassword() string {
	return c.GetString(DatabasePassword)
}

func (c *Config) GetDatabaseHost() string {
	return c.GetString(DatabaseHost)
}

func (c *Config) GetDatabasePort() int {
	return c.GetInt(DatabasePort)
}

func (c *Config) GetDatabaseName() string {
	return c.GetString(DatabaseName)
}

func (c *Config) GetDatabasePoolIdle() int {
	return c.GetInt(DatabasePoolIdle)
}

func (c *Config) GetDatabasePoolMax() int {
	return c.GetInt(DatabasePoolMax)
}

func (c *Config) GetDatabasePoolLifetime() int {
	return c.GetInt(DatabasePoolLifetime)
}

func (c *Config) GetKafkaBootstrapServers() string {
	return c.GetString(KafkaBootstrapServers)
}

func (c *Config) GetKafkaAutoOffsetReset() string {
	return c.GetString(KafkaAutoOffsetReset)
}

func (c *Config) GetKafkaProducerEnabled() bool {
	return c.GetBool(KafkaProducerEnabled)
}

func (c *Config) GetLogLevel() string {
	return c.GetString(LogLevel)
}

func (c *Config) GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", c.GetRedisHost(), c.GetRedisPort())
}

func (c *Config) GetRedisHost() string {
	return c.GetString(RedisHost)
}

func (c *Config) GetRedisPort() int {
	return c.GetInt(RedisPort)
}

func (c *Config) GetRedisDB() int {
	return c.GetInt(RedisDB)
}

func (c *Config) GetRedisUsername() string {
	return c.GetString(RedisUsername)
}

func (c *Config) GetRedisPassword() string {
	return c.GetString(RedisPassword)
}

func (c *Config) GetTelemetryOTLPEndpoint() string {
	return c.GetString(TelemetryOTLPEndpoint)
}

func (c *Config) GetWebPort() string {
	return c.GetString(WebPort)
}

func (c *Config) GetWebPrefork() bool {
	return c.GetBool(WebPrefork)
}
