package configkey

const (
	AppName = "app.name"

	WebPort    = "web.port"
	WebPrefork = "web.prefork"

	LogLevel = "log.level"

	DatabaseMigrations   = "database.migrations"
	DatabaseUsername     = "database.username"
	DatabasePassword     = "database.password"
	DatabaseHost         = "database.host"
	DatabasePort         = "database.port"
	DatabaseName         = "database.name"
	DatabasePoolIdle     = "database.pool.idle"
	DatabasePoolMax      = "database.pool.max"
	DatabasePoolLifetime = "database.pool.lifetime"

	KafkaBootstrapServers = "kafka.bootstrap.servers"
	KafkaGroupId          = "kafka.group.id"
	KafkaAutoOffsetReset  = "kafka.auto.offset.reset"
	KafkaProducerEnabled  = "kafka.producer.enabled"

	TelemetryOTLPEndpoint = "telemetry.otlp.endpoint"

	AuthJWTSecret        = "auth.jwt.secret"
	AuthJWTIssuer        = "auth.jwt.issuer"
	AuthJWTExpireSeconds = "auth.jwt.expire_seconds"
)
