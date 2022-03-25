package constants

const (
	EnvironmentAppRestPort = "APP_REST_PORT"
	EnvironmentDirectory   = "../../.env"

	EnvironmentJWTSecretKey = "JWT_SECRET_KEY"

	EnvironmentLogMode         = "LOG_MODE"
	EnvironmentLogPostgresMode = "LOG_POSTGRES_MODE"

	EnvironmentWhatsmeowDBHost     = "WHATSMEOW_DB_HOST"
	EnvironmentWhatsmeowDBUser     = "WHATSMEOW_DB_USER"
	EnvironmentWhatsmeowDBPassword = "WHATSMEOW_DB_PASSWORD"
	EnvironmentWhatsmeowDBName     = "WHATSMEOW_DB_NAME"
	EnvironmentWhatsmeowDBPort     = "WHATSMEOW_DB_PORT"

	EnvironmentPostgresMigrationDirectory = "POSTGRES_MIGRATION_DIRECTORY"
	EnvironmentPostgresMigrationDialect   = "POSTGRES_MIGRATION_DIALECT"
	EnvironmentPostgresDBHost             = "POSTGRES_DB_HOST"
	EnvironmentPostgresDBUser             = "POSTGRES_DB_USER"
	EnvironmentPostgresDBPassword         = "POSTGRES_DB_PASSWORD"
	EnvironmentPostgresDBName             = "POSTGRES_DB_NAME"
	EnvironmentPostgresDBPort             = "POSTGRES_DB_PORT"
	EnvironmentPostgresDBSSLMode          = "POSTGRES_DB_SSL_MODE"
)
