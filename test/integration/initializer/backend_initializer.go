package initializer

import (
	Context "context"
	SQL "database/sql"
	Error "errors"
	Fmt "fmt"
	OS "os"
	Testing "testing"

	Configuration "github.com/danyel/ecommerce/cmd/config"
	DatabaseConnection "github.com/danyel/ecommerce/cmd/database"
	Logger "github.com/danyel/ecommerce/cmd/logger"
	TestUtils "github.com/danyel/ecommerce/test/testutils"
	_ "github.com/lib/pq" // ← REQUIRED for Goose + sql.Open("postgres")
	Goose "github.com/pressly/goose/v3"
	TestContainers "github.com/testcontainers/testcontainers-go"
	Postgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	RabbitMQ "github.com/testcontainers/testcontainers-go/modules/rabbitmq"
	Database "gorm.io/gorm"
)

type BackendInitializer struct {
	context           *Context.Context
	database          *Database.DB
	postgresContainer *Postgres.PostgresContainer
}

func (backendInitializer *BackendInitializer) initializeDatabaseConfiguration() {
	Configuration.Database().Username = "test"
	Configuration.Database().Password = "test"
	Configuration.Database().Database = "ecommerce"
	Configuration.Database().Schema = "ecommerce"
	if isDevelopment() {
		Configuration.Database().Host = "172.17.0.1"
	} else {
		Configuration.Database().Host = "localhost"
	}
}

func (backendInitializer *BackendInitializer) initializeBrokerConfiguration() {
	Configuration.MessageBroker().Username = "developer"
	Configuration.MessageBroker().Password = "developer"
	Configuration.MessageBroker().Protocol = "amqp"
	if isDevelopment() {
		Configuration.MessageBroker().Addr = "172.17.0.1"
	} else {
		Configuration.MessageBroker().Addr = "localhost"
	}
}

func (backendInitializer *BackendInitializer) initializeRabbitMqTestContainer() TestContainers.Container {
	Logger.Log.Info("Initializing RabbitMQ container %s %s %s %s", Configuration.MessageBroker().Addr, Configuration.MessageBroker().Username, Configuration.MessageBroker().Password, Configuration.MessageBroker().Protocol)
	rabbitmqContainer, err := RabbitMQ.Run(Context.Background(), "rabbitmq:3-management", RabbitMQ.WithAdminUsername(Configuration.MessageBroker().Username), RabbitMQ.WithAdminPassword(Configuration.MessageBroker().Password))
	if err != nil || rabbitmqContainer == nil {
		Logger.Log.Fatalf("failed to start container: %s", err)
		return nil
	}

	host, err := rabbitmqContainer.Host(*backendInitializer.context)
	if err != nil {
		Logger.Log.Fatalf("failed to get container host: %s", err)
		return nil
	}
	Configuration.MessageBroker().Addr = host
	amqpURI, err := rabbitmqContainer.MappedPort(*backendInitializer.context, "5672/tcp")
	if err != nil {
		Logger.Log.Fatalf("failed to start container: %s", err)
		return nil
	}

	Configuration.MessageBroker().Port = amqpURI.Port()
	Logger.Log.Info("container started successfully: %s://%s:%s@%s:%s", Configuration.MessageBroker().Protocol, Configuration.MessageBroker().Username, Configuration.MessageBroker().Password, Configuration.MessageBroker().Addr, Configuration.MessageBroker().Port)

	return rabbitmqContainer
}

func (backendInitializer *BackendInitializer) initializePostgresTestContainer() error {
	postgresContainer, err := Postgres.Run(
		*backendInitializer.context,
		"postgres:18-alpine",
		Postgres.WithDatabase(Configuration.Database().Database),
		Postgres.WithUsername(Configuration.Database().Username),
		Postgres.WithPassword(Configuration.Database().Password),
		Postgres.BasicWaitStrategies(),
	)
	if err != nil {
		return err
	}

	if postgresContainer == nil {
		return Error.New("could not create postgress container")
	}
	port, err := postgresContainer.MappedPort(*backendInitializer.context, "5432/tcp")
	if err != nil {
		return Fmt.Errorf("failed to fetch port: %v", err)
	}

	Configuration.Database().Port = port.Port()
	Logger.Log.Info("DatabaseConfiguration: postgres://%s:%s@%s:%s/%s", Configuration.Database().Username, Configuration.Database().Password, Configuration.Database().Host, Configuration.Database().Port, Configuration.Database().Schema)
	backendInitializer.postgresContainer = postgresContainer
	return err
}

func (backendInitializer *BackendInitializer) initializeMigrationScripts(postgresContainer *Postgres.PostgresContainer) {
	connectionString, err := postgresContainer.ConnectionString(*backendInitializer.context)
	if err != nil {
		Logger.Log.Fatalf("failed to get DSN: %v", err)
	}

	connectionString = connectionString + "sslmode=disable"

	Logger.Log.Info("Migration url %s", connectionString)

	connection, err := SQL.Open("postgres", connectionString)
	if err != nil {
		Logger.Log.Fatalf("sql open failed: %v", err)
	}
	defer func(DB *SQL.DB) {
		err := DB.Close()
		if err != nil {
			Logger.Log.Fatalf("failed to close connection: %v", err)
		}
	}(connection)

	Goose.SetBaseFS(nil)
	err = Goose.SetDialect("postgres")
	if err != nil {
		Logger.Log.Fatalf("failed to set goose dialect: %v", err)
	}

	if err := Goose.Up(connection, "../../migrations"); err != nil {
		Logger.Log.Fatalf("goose migration failed: %v", err)
	}
}

func (backendInitializer *BackendInitializer) connect() (*Database.DB, error) {
	return DatabaseConnection.Connect(Configuration.Database())
}

func (backendInitializer *BackendInitializer) Terminate() {
	err := backendInitializer.postgresContainer.Terminate(*backendInitializer.context)
	if err != nil {
		Logger.Log.Fatalf("failed to terminate postgres container: %v", err)
	}
}

func (backendInitializer *BackendInitializer) DatabaseConnection() *Database.DB {
	return backendInitializer.database
}

func (backendInitializer *BackendInitializer) TestContainers(unitTest *Testing.T) {
	backendInitializer.initializeBrokerConfiguration()
	rabbitMqTestContainer := backendInitializer.initializeRabbitMqTestContainer()

	if rabbitMqTestContainer != nil {
		unitTest.Cleanup(func() {
			backgroundContext := Context.Background()
			if err := rabbitMqTestContainer.Terminate(backgroundContext); err != nil {
				Logger.Log.Fatalf("failed to terminate container: %s", err)
				OS.Exit(0)
			}
		})
	}
	backendInitializer.initializeDatabaseConfiguration()
	err := backendInitializer.initializePostgresTestContainer()
	if err != nil {
		Logger.Log.Fatalf("failed to initialize test pg: %v", err)
	}
}

func (backendInitializer *BackendInitializer) Run() {
	backendInitializer.initializeDatabaseConfiguration()
	var err error
	backendInitializer.initializeMigrationScripts(backendInitializer.postgresContainer)
	backendInitializer.database, err = backendInitializer.connect()
	if err != nil {
		Logger.Log.Fatalf("failed to connect to database: %v", err)
	}
}

func isDevelopment() bool {
	return OS.Getenv("ENV") == "dev"
}

func NewBackendInitializer() *BackendInitializer {
	TestUtils.PreInitTest()
	backgroundContext := Context.Background()
	return &BackendInitializer{
		context: &backgroundContext,
	}
}
