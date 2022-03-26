package domain

import (
	"database/sql"
	"net/http"
	"os"
	"strconv"
	"svc-whatsapp/domain/constants"
	"svc-whatsapp/libraries"
	messagebroker "svc-whatsapp/libraries/messagesbroker"

	socketio "github.com/googollee/go-socket.io"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"go.mau.fi/whatsmeow/store/sqlstore"
)

type Config struct {
	App                *gin.Engine
	Socket             *socketio.Server
	Validator          *validator.Validate
	SecretKey          string
	Postgres           *gorm.DB
	PostgresConnection *sql.DB
	StoreContainer     *sqlstore.Container
	NsqProducer        messagebroker.Producer
	NsqConsumer        messagebroker.Consumer
	WhatsappWorker     *libraries.WorkerPool
}

func LoadConfiguration() (config Config, err error) {
	// load env
	if err = godotenv.Load(constants.EnvironmentDirectory); err != nil {
		return config, err
	}

	// validator
	config.Validator = validator.New()

	// gin-gonic
	config.App = gin.Default()

	// go-socketio
	config.Socket = socketio.NewServer(nil)

	// nsq-go
	config.NsqProducer = messagebroker.NewProducer("0.0.0.0:4150", "svc-whatsapp")
	config.NsqConsumer = messagebroker.NewConsumer("0.0.0.0:4150", "svc-whatsapp")

	// Whatsmeow
	whatsmeowLibrary := libraries.WhatsmeowLibrary{
		DBHost:     os.Getenv(constants.EnvironmentWhatsmeowDBHost),
		DBUser:     os.Getenv(constants.EnvironmentWhatsmeowDBUser),
		DBPassword: os.Getenv(constants.EnvironmentWhatsmeowDBPassword),
		DBPort:     os.Getenv(constants.EnvironmentWhatsmeowDBPort),
		DBName:     os.Getenv(constants.EnvironmentWhatsmeowDBName),
	}
	config.StoreContainer, err = whatsmeowLibrary.Connect()
	if err != nil {
		return config, err
	}

	// postgres
	dbLogMode, err := strconv.Atoi(os.Getenv(constants.EnvironmentLogPostgresMode))
	if err != nil {
		return config, err
	}

	postgresLibrary := libraries.PostgresLibrary{
		MigrationDirectory: os.Getenv(constants.EnvironmentPostgresMigrationDirectory),
		MigrationDialect:   os.Getenv(constants.EnvironmentPostgresMigrationDialect),
		DBHost:             os.Getenv(constants.EnvironmentPostgresDBHost),
		DBUser:             os.Getenv(constants.EnvironmentPostgresDBUser),
		DBPassword:         os.Getenv(constants.EnvironmentPostgresDBPassword),
		DBPort:             os.Getenv(constants.EnvironmentPostgresDBPort),
		DBName:             os.Getenv(constants.EnvironmentPostgresDBName),
		DBSSLMode:          os.Getenv(constants.EnvironmentPostgresDBSSLMode),
		LogMode:            dbLogMode,
	}

	config.Postgres, config.PostgresConnection, err = postgresLibrary.ConnectAndValidate()

	if err != nil {
		return config, err
	}

	//postgres migrate
	if err = postgresLibrary.Migrate(config.PostgresConnection); err != nil {
		return config, err
	}

	// load secret key JWT
	config.SecretKey = os.Getenv(constants.EnvironmentJWTSecretKey)

	// whatsapp-worker
	config.WhatsappWorker = libraries.NewWorkerPool(10, config.StoreContainer, config.Postgres)

	return config, err
}

func HttpListen(app *gin.Engine) *http.Server {
	return &http.Server{
		Addr:    os.Getenv(constants.EnvironmentAppRestPort),
		Handler: app,
	}
}
