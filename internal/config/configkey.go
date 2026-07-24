package config

const (
	AppName = "app.name"

	AuthJWTSecret        = "auth.jwt.secret"
	AuthJWTIssuer        = "auth.jwt.issuer"
	AuthJWTExpireSeconds = "auth.jwt.expire_seconds"

	AWSRegion       = "aws.region"
	AWSBaseEndpoint = "aws.base_endpoint"

	DatabaseMigrations   = "database.migrations"
	DatabaseUsername     = "database.username"
	DatabasePassword     = "database.password"
	DatabaseHost         = "database.host"
	DatabasePort         = "database.port"
	DatabaseName         = "database.name"
	DatabasePoolIdle     = "database.pool.idle"
	DatabasePoolMax      = "database.pool.max"
	DatabasePoolLifetime = "database.pool.lifetime"

	ElasticsearchAddress = "elasticsearch.address"

	KafkaBootstrapServers   = "kafka.bootstrap.servers"
	KafkaAutoOffsetReset    = "kafka.auto.offset.reset"
	KafkaConsumerMaxRetries = "kafka.consumer.max_retries"
	KafkaProducerEnabled    = "kafka.producer.enabled"

	OutboxPollIntervalSeconds = "outbox.poll_interval_seconds"
	OutboxBatchSize           = "outbox.batch_size"

	LogLevel = "log.level"

	IdempotencyCleanupIntervalSeconds = "idempotency.cleanup_interval_seconds"

	RedisHost     = "redis.host"
	RedisPort     = "redis.port"
	RedisDB       = "redis.db"
	RedisUsername = "redis.username"
	RedisPassword = "redis.password"

	TelemetryOTLPEndpoint = "telemetry.otlp.endpoint"

	WebPort    = "web.port"
	WebPrefork = "web.prefork"
)
