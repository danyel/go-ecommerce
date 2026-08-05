package initializer

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/danyel/ecommerce/cmd/config"
	"github.com/danyel/ecommerce/cmd/database"
	_ "github.com/lib/pq" // ← REQUIRED for Goose + sql.Open("postgres")
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/modules/rabbitmq"
	"gorm.io/gorm"
)

type BackendInitializer struct {
	context               *context.Context
	db                    *gorm.DB
	pg                    *tcpostgres.PostgresContainer
	BrokerConfiguration   *config.BrokerConfiguration
	DatabaseConfiguration *config.DatabaseConfiguration
}

func (b *BackendInitializer) initializeDatabaseConfiguration() {
	envFile := os.Getenv("ENV")
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
	envFile := os.Getenv("ENV")
	b.BrokerConfiguration.Username = "developer"
	b.BrokerConfiguration.Password = "developer"
	b.BrokerConfiguration.Protocol = "amqp"
	if envFile == "dev" {
		b.BrokerConfiguration.Addr = "172.17.0.1"
	} else {
		b.BrokerConfiguration.Addr = "localhost"
	}
}

func (b *BackendInitializer) initializeRabbitMqTestContainer() testcontainers.Container {
	log.Printf("Initializing RabbitMQ container %s %s %s %s", b.BrokerConfiguration.Addr, b.BrokerConfiguration.Username, b.BrokerConfiguration.Password, b.BrokerConfiguration.Protocol)
	rabbitmqContainer, err := rabbitmq.Run(context.Background(), "rabbitmq:3-management", rabbitmq.WithAdminUsername(b.BrokerConfiguration.Username), rabbitmq.WithAdminPassword(b.BrokerConfiguration.Password))
	if err != nil || rabbitmqContainer == nil {
		log.Printf("failed to start container: %s", err)
		return nil
	}

	host, err := rabbitmqContainer.Host(*b.context)
	if err != nil {
		log.Printf("failed to get container host: %s", err)
		return nil
	}
	b.BrokerConfiguration.Addr = host
	amqpURI, err := rabbitmqContainer.MappedPort(*b.context, "5672/tcp")
	if err != nil {
		log.Printf("failed to start container: %s", err)
		return nil
	}

	b.BrokerConfiguration.Port = amqpURI.Port()
	log.Printf("container started successfully: %s://%s:%s@%s:%s", b.BrokerConfiguration.Protocol, b.BrokerConfiguration.Username, b.BrokerConfiguration.Password, b.BrokerConfiguration.Addr, b.BrokerConfiguration.Port)

	return rabbitmqContainer
}

func (b *BackendInitializer) initializePostgresTestContainer() error {
	pg, err := tcpostgres.Run(*b.context,
		"postgres:18-alpine",
		tcpostgres.WithDatabase(b.DatabaseConfiguration.Database),
		tcpostgres.WithUsername(b.DatabaseConfiguration.Username),
		tcpostgres.WithPassword(b.DatabaseConfiguration.Password),
		tcpostgres.BasicWaitStrategies(),
	)

	if pg == nil {
		log.Fatalf("failed to start container")
	}
	port, err := pg.MappedPort(*b.context, "5432/tcp")

	if err != nil {
		log.Fatalf("failed to fetch port: %v", err)
	}

	b.DatabaseConfiguration.Port = port.Port()
	log.Printf("DatabaseConfiguration: postgres://%s:%s@%s:%s/%s", b.DatabaseConfiguration.Username, b.DatabaseConfiguration.Password, b.DatabaseConfiguration.Host, b.DatabaseConfiguration.Port, b.DatabaseConfiguration.Schema)
	b.pg = pg
	return err
}

func (b *BackendInitializer) initializeMigrationScripts(pg *tcpostgres.PostgresContainer) {
	dsn, err := pg.ConnectionString(*b.context)
	if err != nil {
		log.Fatalf("failed to get DSN: %v", err)
	}

	dsn = dsn + "sslmode=disable"

	log.Printf("Migration url %s", dsn)

	sqlDB, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("sql open failed: %v", err)
	}
	defer func(sqlDB *sql.DB) {
		err := sqlDB.Close()
		if err != nil {
			log.Fatalf("failed to close connection: %v", err)
		}
	}(sqlDB)

	goose.SetBaseFS(nil)
	err = goose.SetDialect("postgres")
	if err != nil {
		log.Fatalf("failed to set goose dialect: %v", err)
	}

	if err := goose.Up(sqlDB, "../../migrations"); err != nil {
		log.Fatalf("goose migration failed: %v", err)
	}
}

func (b *BackendInitializer) connect(c *config.DatabaseConfiguration) (*gorm.DB, error) {
	return database.Connect(c)
}

func (b *BackendInitializer) Terminate() {
	err := b.pg.Terminate(*b.context)
	if err != nil {
		log.Fatalf("failed to terminate postgres container: %v", err)
	}
}

func (b *BackendInitializer) Db() *gorm.DB {
	return b.db
}

func (b *BackendInitializer) TestContainers(t *testing.T) {
	b.initializeBrokerConfiguration()
	rabbitMqTestContainer := b.initializeRabbitMqTestContainer()

	if rabbitMqTestContainer != nil {
		t.Cleanup(func() {
			ctx := context.Background()
			if err := rabbitMqTestContainer.Terminate(ctx); err != nil {
				log.Printf("failed to terminate container: %s", err)
			}
		})
	}
	b.initializeDatabaseConfiguration()
	err := b.initializePostgresTestContainer()
	if err != nil {
		log.Fatalf("failed to initialize test pg: %v", err)
	}
}

func (b *BackendInitializer) Run() {
	b.initializeDatabaseConfiguration()
	var err error
	b.initializeMigrationScripts(b.pg)
	b.db, err = b.connect(b.DatabaseConfiguration)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
}

func NewBackendInitializer() *BackendInitializer {
	c := context.Background()
	configuration := config.NewBrokerConfiguration()
	databaseConfiguration := config.NewDatabaseConfiguration()
	return &BackendInitializer{
		context:               &c,
		BrokerConfiguration:   &configuration,
		DatabaseConfiguration: &databaseConfiguration,
	}
}
