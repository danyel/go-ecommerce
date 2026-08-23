package initializer

import (
	Context "context"
	SQL "database/sql"
	Log "log"
	OS "os"
	Testing "testing"

	Configuration "github.com/danyel/ecommerce/cmd/config"
	DatabaseConnection "github.com/danyel/ecommerce/cmd/database"
	_ "github.com/lib/pq" // ← REQUIRED for Goose + sql.Open("postgres")
	Goose "github.com/pressly/goose/v3"
	TestContainers "github.com/testcontainers/testcontainers-go"
	Postgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	RabbitMQ "github.com/testcontainers/testcontainers-go/modules/rabbitmq"
	Database "gorm.io/gorm"
)

type BackendInitializer struct {
	context                    *Context.Context
	database                   *Database.DB
	postgresContainer          *Postgres.PostgresContainer
	MessageBrokerConfiguration *Configuration.MessageBrokerConfiguration
	DatabaseConfiguration      *Configuration.DatabaseConfiguration
}

func (backendInitializer *BackendInitializer) initializeDatabaseConfiguration() {
	backendInitializer.DatabaseConfiguration.Username = "test"
	backendInitializer.DatabaseConfiguration.Password = "test"
	backendInitializer.DatabaseConfiguration.Database = "ecommerce"
	backendInitializer.DatabaseConfiguration.Schema = "ecommerce"
	if isDevelopment() {
		backendInitializer.DatabaseConfiguration.Host = "172.17.0.1"
	} else {
		backendInitializer.DatabaseConfiguration.Host = "localhost"
	}
}

func (backendInitializer *BackendInitializer) initializeBrokerConfiguration() {
	backendInitializer.MessageBrokerConfiguration.Username = "developer"
	backendInitializer.MessageBrokerConfiguration.Password = "developer"
	backendInitializer.MessageBrokerConfiguration.Protocol = "amqp"
	if isDevelopment() {
		backendInitializer.MessageBrokerConfiguration.Addr = "172.17.0.1"
	} else {
		backendInitializer.MessageBrokerConfiguration.Addr = "localhost"
	}
}

func (backendInitializer *BackendInitializer) initializeRabbitMqTestContainer() TestContainers.Container {
	Log.Printf("Initializing RabbitMQ container %s %s %s %s", backendInitializer.MessageBrokerConfiguration.Addr, backendInitializer.MessageBrokerConfiguration.Username, backendInitializer.MessageBrokerConfiguration.Password, backendInitializer.MessageBrokerConfiguration.Protocol)
	rabbitmqContainer, err := RabbitMQ.Run(Context.Background(), "rabbitmq:3-management", RabbitMQ.WithAdminUsername(backendInitializer.MessageBrokerConfiguration.Username), RabbitMQ.WithAdminPassword(backendInitializer.MessageBrokerConfiguration.Password))
	if err != nil || rabbitmqContainer == nil {
		Log.Printf("failed to start container: %s", err)
		return nil
	}

	host, err := rabbitmqContainer.Host(*backendInitializer.context)
	if err != nil {
		Log.Printf("failed to get container host: %s", err)
		return nil
	}
	backendInitializer.MessageBrokerConfiguration.Addr = host
	amqpURI, err := rabbitmqContainer.MappedPort(*backendInitializer.context, "5672/tcp")
	if err != nil {
		Log.Printf("failed to start container: %s", err)
		return nil
	}

	backendInitializer.MessageBrokerConfiguration.Port = amqpURI.Port()
	Log.Printf("container started successfully: %s://%s:%s@%s:%s", backendInitializer.MessageBrokerConfiguration.Protocol, backendInitializer.MessageBrokerConfiguration.Username, backendInitializer.MessageBrokerConfiguration.Password, backendInitializer.MessageBrokerConfiguration.Addr, backendInitializer.MessageBrokerConfiguration.Port)

	return rabbitmqContainer
}

func (backendInitializer *BackendInitializer) initializePostgresTestContainer() error {
	postgresContainer, err := Postgres.Run(*backendInitializer.context,
		"postgres:18-alpine",
		Postgres.WithDatabase(backendInitializer.DatabaseConfiguration.Database),
		Postgres.WithUsername(backendInitializer.DatabaseConfiguration.Username),
		Postgres.WithPassword(backendInitializer.DatabaseConfiguration.Password),
		Postgres.BasicWaitStrategies(),
	)

	if postgresContainer == nil {
		Log.Fatalf("failed to start container")
	}
	port, err := postgresContainer.MappedPort(*backendInitializer.context, "5432/tcp")

	if err != nil {
		Log.Fatalf("failed to fetch port: %v", err)
	}

	backendInitializer.DatabaseConfiguration.Port = port.Port()
	Log.Printf("DatabaseConfiguration: postgres://%s:%s@%s:%s/%s", backendInitializer.DatabaseConfiguration.Username, backendInitializer.DatabaseConfiguration.Password, backendInitializer.DatabaseConfiguration.Host, backendInitializer.DatabaseConfiguration.Port, backendInitializer.DatabaseConfiguration.Schema)
	backendInitializer.postgresContainer = postgresContainer
	return err
}

func (backendInitializer *BackendInitializer) initializeMigrationScripts(postgresContainer *Postgres.PostgresContainer) {
	connectionString, err := postgresContainer.ConnectionString(*backendInitializer.context)
	if err != nil {
		Log.Fatalf("failed to get DSN: %v", err)
	}

	connectionString = connectionString + "sslmode=disable"

	Log.Printf("Migration url %s", connectionString)

	connection, err := SQL.Open("postgres", connectionString)
	if err != nil {
		Log.Fatalf("sql open failed: %v", err)
	}
	defer func(DB *SQL.DB) {
		err := DB.Close()
		if err != nil {
			Log.Fatalf("failed to close connection: %v", err)
		}
	}(connection)

	Goose.SetBaseFS(nil)
	err = Goose.SetDialect("postgres")
	if err != nil {
		Log.Fatalf("failed to set goose dialect: %v", err)
	}

	if err := Goose.Up(connection, "../../migrations"); err != nil {
		Log.Fatalf("goose migration failed: %v", err)
	}
}

func (backendInitializer *BackendInitializer) connect(databaseConfiguration *Configuration.DatabaseConfiguration) (*Database.DB, error) {
	return DatabaseConnection.Connect(databaseConfiguration)
}

func (backendInitializer *BackendInitializer) Terminate() {
	err := backendInitializer.postgresContainer.Terminate(*backendInitializer.context)
	if err != nil {
		Log.Fatalf("failed to terminate postgres container: %v", err)
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
				Log.Printf("failed to terminate container: %s", err)
			}
		})
	}
	backendInitializer.initializeDatabaseConfiguration()
	err := backendInitializer.initializePostgresTestContainer()
	if err != nil {
		Log.Fatalf("failed to initialize test pg: %v", err)
	}
}

func (backendInitializer *BackendInitializer) Run() {
	backendInitializer.initializeDatabaseConfiguration()
	var err error
	backendInitializer.initializeMigrationScripts(backendInitializer.postgresContainer)
	backendInitializer.database, err = backendInitializer.connect(backendInitializer.DatabaseConfiguration)
	if err != nil {
		Log.Fatalf("failed to connect to database: %v", err)
	}
}

func isDevelopment() bool {
	return OS.Getenv("ENV") == "dev"
}

func NewBackendInitializer() *BackendInitializer {
	backgroundContext := Context.Background()
	messageBrokerConfiguration := Configuration.NewMessageBrokerConfiguration()
	databaseConfiguration := Configuration.NewDatabaseConfiguration()
	return &BackendInitializer{
		context:                    &backgroundContext,
		MessageBrokerConfiguration: &messageBrokerConfiguration,
		DatabaseConfiguration:      &databaseConfiguration,
	}
}
