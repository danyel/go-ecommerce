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
	context               *Context.Context
	db                    *Database.DB
	pg                    *Postgres.PostgresContainer
	BrokerConfiguration   *Configuration.BrokerConfiguration
	DatabaseConfiguration *Configuration.DatabaseConfiguration
}

func (b *BackendInitializer) initializeDatabaseConfiguration() {
	envFile := OS.Getenv("ENV")
	b.DatabaseConfiguration.Username = "test"
	b.DatabaseConfiguration.Password = "test"
	b.DatabaseConfiguration.Database = "ecommerce"
	b.DatabaseConfiguration.Schema = "ecommerce"
	if envFile == "dev" {
		b.DatabaseConfiguration.Host = "172.17.0.1"
	} else {
		b.DatabaseConfiguration.Host = "localhost"
	}
}

func (b *BackendInitializer) initializeBrokerConfiguration() {
	envFile := OS.Getenv("ENV")
	b.BrokerConfiguration.Username = "developer"
	b.BrokerConfiguration.Password = "developer"
	b.BrokerConfiguration.Protocol = "amqp"
	if envFile == "dev" {
		b.BrokerConfiguration.Addr = "172.17.0.1"
	} else {
		b.BrokerConfiguration.Addr = "localhost"
	}
}

func (b *BackendInitializer) initializeRabbitMqTestContainer() TestContainers.Container {
	Log.Printf("Initializing RabbitMQ container %s %s %s %s", b.BrokerConfiguration.Addr, b.BrokerConfiguration.Username, b.BrokerConfiguration.Password, b.BrokerConfiguration.Protocol)
	rabbitmqContainer, err := RabbitMQ.Run(Context.Background(), "rabbitmq:3-management", RabbitMQ.WithAdminUsername(b.BrokerConfiguration.Username), RabbitMQ.WithAdminPassword(b.BrokerConfiguration.Password))
	if err != nil || rabbitmqContainer == nil {
		Log.Printf("failed to start container: %s", err)
		return nil
	}

	host, err := rabbitmqContainer.Host(*b.context)
	if err != nil {
		Log.Printf("failed to get container host: %s", err)
		return nil
	}
	b.BrokerConfiguration.Addr = host
	amqpURI, err := rabbitmqContainer.MappedPort(*b.context, "5672/tcp")
	if err != nil {
		Log.Printf("failed to start container: %s", err)
		return nil
	}

	b.BrokerConfiguration.Port = amqpURI.Port()
	Log.Printf("container started successfully: %s://%s:%s@%s:%s", b.BrokerConfiguration.Protocol, b.BrokerConfiguration.Username, b.BrokerConfiguration.Password, b.BrokerConfiguration.Addr, b.BrokerConfiguration.Port)

	return rabbitmqContainer
}

func (b *BackendInitializer) initializePostgresTestContainer() error {
	pg, err := Postgres.Run(*b.context,
		"postgres:18-alpine",
		Postgres.WithDatabase(b.DatabaseConfiguration.Database),
		Postgres.WithUsername(b.DatabaseConfiguration.Username),
		Postgres.WithPassword(b.DatabaseConfiguration.Password),
		Postgres.BasicWaitStrategies(),
	)

	if pg == nil {
		Log.Fatalf("failed to start container")
	}
	port, err := pg.MappedPort(*b.context, "5432/tcp")

	if err != nil {
		Log.Fatalf("failed to fetch port: %v", err)
	}

	b.DatabaseConfiguration.Port = port.Port()
	Log.Printf("DatabaseConfiguration: postgres://%s:%s@%s:%s/%s", b.DatabaseConfiguration.Username, b.DatabaseConfiguration.Password, b.DatabaseConfiguration.Host, b.DatabaseConfiguration.Port, b.DatabaseConfiguration.Schema)
	b.pg = pg
	return err
}

func (b *BackendInitializer) initializeMigrationScripts(pg *Postgres.PostgresContainer) {
	dsn, err := pg.ConnectionString(*b.context)
	if err != nil {
		Log.Fatalf("failed to get DSN: %v", err)
	}

	dsn = dsn + "sslmode=disable"

	Log.Printf("Migration url %s", dsn)

	sqlDB, err := SQL.Open("postgres", dsn)
	if err != nil {
		Log.Fatalf("sql open failed: %v", err)
	}
	defer func(sqlDB *SQL.DB) {
		err := sqlDB.Close()
		if err != nil {
			Log.Fatalf("failed to close connection: %v", err)
		}
	}(sqlDB)

	Goose.SetBaseFS(nil)
	err = Goose.SetDialect("postgres")
	if err != nil {
		Log.Fatalf("failed to set goose dialect: %v", err)
	}

	if err := Goose.Up(sqlDB, "../../migrations"); err != nil {
		Log.Fatalf("goose migration failed: %v", err)
	}
}

func (b *BackendInitializer) connect(c *Configuration.DatabaseConfiguration) (*Database.DB, error) {
	return DatabaseConnection.Connect(c)
}

func (b *BackendInitializer) Terminate() {
	err := b.pg.Terminate(*b.context)
	if err != nil {
		Log.Fatalf("failed to terminate postgres container: %v", err)
	}
}

func (b *BackendInitializer) Db() *Database.DB {
	return b.db
}

func (b *BackendInitializer) TestContainers(t *Testing.T) {
	b.initializeBrokerConfiguration()
	rabbitMqTestContainer := b.initializeRabbitMqTestContainer()

	if rabbitMqTestContainer != nil {
		t.Cleanup(func() {
			ctx := Context.Background()
			if err := rabbitMqTestContainer.Terminate(ctx); err != nil {
				Log.Printf("failed to terminate container: %s", err)
			}
		})
	}
	b.initializeDatabaseConfiguration()
	err := b.initializePostgresTestContainer()
	if err != nil {
		Log.Fatalf("failed to initialize test pg: %v", err)
	}
}

func (b *BackendInitializer) Run() {
	b.initializeDatabaseConfiguration()
	var err error
	b.initializeMigrationScripts(b.pg)
	b.db, err = b.connect(b.DatabaseConfiguration)
	if err != nil {
		Log.Fatalf("failed to connect to database: %v", err)
	}
}

func NewBackendInitializer() *BackendInitializer {
	c := Context.Background()
	configuration := Configuration.NewBrokerConfiguration()
	databaseConfiguration := Configuration.NewDatabaseConfiguration()
	return &BackendInitializer{
		context:               &c,
		BrokerConfiguration:   &configuration,
		DatabaseConfiguration: &databaseConfiguration,
	}
}
