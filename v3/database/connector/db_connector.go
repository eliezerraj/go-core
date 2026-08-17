package connector

import(
	"context"
	"time"

	"github.com/eliezerraj/go-core/v3/logger"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// IQueryExecutor defines the interface for executing queries against a database.
type IQueryExecutor interface {
	Exec(context.Context, string, ...interface{}) (commandTag pgconn.CommandTag, err error)
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) (pgx.Row)
	BeginTx(context.Context, pgx.TxOptions) (pgx.Tx, error)
	Stat() *pgxpool.Stat
	Close()
}

// Connector configuration struct for database connections.
type ConnectorConfig struct {
	DSN                               string
	MaxConnIdleTime, MaxConnLifeTime, DBConnTimeout, HealthCheckPeriod time.Duration
	MaxConns, MinConns                int32
}

// DatabaseConnector struct manages database connections for reading and writing.
type DatabaseConnector struct {
	application                  string
	readerConnectorConfig, writerConnectorConfig ConnectorConfig
	reader, writer               IQueryExecutor
}

// IDatabaseConnector defines the interface for a database connector that provides separate readers and writers.
type IDatabaseConnector interface {
	//AcquireTransaction(ctx context.Context) IQueryExecutor
	Reader() IQueryExecutor
	Writer() IQueryExecutor
}

// Just a placeholder for the actual implementation of the DatabaseConnector methods.
func (dbc *DatabaseConnector) Reader() IQueryExecutor { return dbc.reader }
func (dbc *DatabaseConnector) Writer() IQueryExecutor { return dbc.writer }

// NewDatabaseConnector creates a new DatabaseConnector configuration with separate reader and writer connections.
func NewDatabaseConnector(application string, readerConnectorConfig, writerConnectorConfig ConnectorConfig) (IDatabaseConnector, error) {
	logger.InfoOutCtx("initializing database connector SUCCESSFULLY")

	var writer IQueryExecutor
	var reader IQueryExecutor

	readerPgConnection := &PgConnection{}
	reader, err := readerPgConnection.NewPool(context.Background(), readerConnectorConfig)
	if err != nil {
		logger.ErrorOutCtx("failed to create reader pool")
		return nil, err
	}

	writerPgConnection := &PgConnection{}
	writer, err = writerPgConnection.NewPool(context.Background(), writerConnectorConfig)
	if err != nil {
		logger.ErrorOutCtx("failed to create writer pool")
		return nil, err
	}

	return &DatabaseConnector{
		application:    application,
		readerConnectorConfig:  readerConnectorConfig,
		writerConnectorConfig:  writerConnectorConfig,
		reader:         reader,
		writer:         writer,
	}, nil
}