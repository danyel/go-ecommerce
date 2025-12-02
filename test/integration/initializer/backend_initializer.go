package initializer

import (
	"context"
	"database/sql"
	"log"
	"os"

	"github.com/danyel/ecommerce/cmd/config"
	"github.com/danyel/ecommerce/cmd/database"
	_ "github.com/lib/pq" // ← REQUIRED for Goose + sql.Open("postgres")
	"github.com/pressly/goose/v3"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"gorm.io/gorm"
)

type BackendInitializer struct {
	context *context.Context
	db      *gorm.DB
	pg      *tcpostgres.PostgresContainer
}

func (b *BackendInitializer) initializeDatabaseConfiguration() config.DatabaseConfiguration {
	envFile := os.Getenv("ENV")
	databaseConfiguration := config.NewDatabaseConfiguration()
	databaseConfiguration.Username = "test"
	databaseConfiguration.Password = "test"
	databaseConfiguration.Database = "ecommerce"
	databaseConfiguration.Schema = "ecommerce"
	if envFile == "dev" {
		databaseConfiguration.Host = "172.17.0.1"
	} else {
		databaseConfiguration.Host = "localhost"
	}
	return databaseConfiguration
}

func (b *BackendInitializer) updateDatabaseConfiguration(c *tcpostgres.PostgresContainer, d *config.DatabaseConfiguration) {
	ports, err := c.Ports(*b.context)
	if err != nil {
		log.Fatalf("failed to get ports: %v", err)
	}
	d.Port = ports["5432/tcp"][0].HostPort
}

func (b *BackendInitializer) initializeTestContainer(d *config.DatabaseConfiguration) (*tcpostgres.PostgresContainer, error) {
	pg, err := tcpostgres.Run(*b.context,
		"postgres:18-alpine",
		tcpostgres.WithDatabase(d.Database),
		tcpostgres.WithUsername(d.Username),
		tcpostgres.WithPassword(d.Password),
		tcpostgres.BasicWaitStrategies(),
	)

	if pg == nil {
		log.Fatalf("failed to start container")
	}

	return pg, err
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

func (b *BackendInitializer) Run() {
	databaseConfiguration := b.initializeDatabaseConfiguration()
	pg, err := b.initializeTestContainer(&databaseConfiguration)
	if err != nil {
		log.Fatalf("failed to initialize test pg: %v", err)
	}
	b.updateDatabaseConfiguration(pg, &databaseConfiguration)
	b.initializeMigrationScripts(pg)
	b.db, err = b.connect(&databaseConfiguration)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
}

func NewBackendInitializer() *BackendInitializer {
	c := context.Background()
	return &BackendInitializer{
		context: &c,
	}
}
