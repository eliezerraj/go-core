package connector
import (
	"context"
	"errors"
	"go.uber.org/zap"

	"github.com/eliezerraj/go-core/v3/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgConnection struct {
	connString         string
	minConns           int32
	maxConns           int32
	conn               *pgxpool.Pool
}

// Pool returns the underlying pgxpool.Pool instance.
func (pgc *PgConnection) Pool() *pgxpool.Pool { return pgc.conn }

// NewPool initializes a new database connection pool using the provided configuration.
func (p *PgConnection) NewPool(ctx context.Context, connectorConfig ConnectorConfig) (*pgxpool.Pool, error) {
	logger.Info(ctx, "initializing a database pool SUCCESSFULLY")

	p.connString = connectorConfig.DSN
	pgxConfig, err := pgxpool.ParseConfig(p.connString)
	if err != nil {
		logger.Error(ctx, "failed to parse database connection string", zap.Error(err))
		return nil, err
	}

	pgxConfig.HealthCheckPeriod = connectorConfig.MaxConnIdleTime / 2
	pgxConfig.ConnConfig.ConnectTimeout = connectorConfig.DBConnTimeout
	pgxConfig.MaxConnLifetime = connectorConfig.MaxConnLifeTime
	pgxConfig.MaxConnIdleTime = connectorConfig.MaxConnIdleTime
	pgxConfig.MaxConns = connectorConfig.MaxConns
	pgxConfig.MinConns = connectorConfig.MinConns

	p.conn, err = pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		logger.Error(ctx, "failed to create database connection pool", zap.Error(err))
		return nil, err
	}
	return p.conn, nil
}

// Ping check.
func (p *PgConnection) Ping(ctx context.Context) error {
	logger.Info(ctx, "pinging the database pool")

	conn, err := p.Pool().Acquire(ctx)
	if err != nil {
		return err
	}

	err = conn.Conn().Ping(ctx)
	if err != nil {
		conn.Release()
		return err
	}

	conn.Release()
	return nil
}

// ClosePool closes the database connection pool.
func (p *PgConnection) Close() {
	logger.InfoOutCtx("closing the database pool")
	if p.Pool() != nil {
		p.Pool().Close()
	}
}

// Query Executor
func (p *PgConnection) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	logger.Info(ctx, "executing Query: ", zap.String("sql", sql), zap.Any("args", args))

	if p.conn == nil {
		logger.Error(ctx, "database connection pool is not initialized", zap.Error(errors.New("database connection pool is not initialized")))
		return nil, errors.New("database connection pool is not initialized")
	}

	f := p.conn.Query

	res, err := f(ctx, sql, args...)
	
	if err != nil {
		logger.Error(ctx, "failed to execute query", zap.Error(err))
		return nil, err	
	}

	if res.Err() != nil {
		logger.Error(ctx, "query execution error", zap.Error(res.Err()))
		return nil, res.Err()
	}

	return res, nil
}

// BeginTx starts a new transaction and returns the transaction object.
func (p *PgConnection) BeginTx(ctx context.Context) (pgx.Tx, error) {
	logger.Info(ctx, "request and starting a new transaction")

	if p.conn == nil {
		logger.Error(ctx, "database connection pool is not initialized", zap.Error(errors.New("database connection pool is not initialized")))
		return nil, errors.New("database connection pool is not initialized")
	}

	tx, err := p.conn.Begin(ctx)
	if err != nil {
		logger.Error(ctx, "failed to begin transaction", zap.Error(err))
		return nil, err
	}

	return tx, nil
}

// Query Executor
func (p *PgConnection) QueryRow(ctx context.Context, tx *pgx.Tx, sql string, args ...interface{}) (pgx.Rows, error) {
	logger.Info(ctx, "executing QueryRow: ", zap.String("sql", sql), zap.Any("args", args))

	if p.conn == nil {
		logger.Error(ctx, "database connection pool is not initialized", zap.Error(errors.New("database connection pool is not initialized")))
		return nil, errors.New("database connection pool is not initialized")
	}

	f := p.conn.Query

	if tx != nil {
		f = (*tx).Query
	}

	res, err := f(ctx, sql, args...)
	
	if err != nil {
		logger.Error(ctx, "failed to execute query", zap.Error(err))
		return nil, err	
	}

	if res.Err() != nil {
		logger.Error(ctx, "query execution error", zap.Error(res.Err()))
		return nil, res.Err()
	}

	return res, nil
}

// Release is intentionally a no-op because pgxpool.Pool does not expose a pool-level
// Release method. Individual acquired connections are released via pgxpool.Conn.Release(),
// while the pool lifecycle is managed by Close().
func (p *PgConnection) Release() error {
	return nil
}

func (p *PgConnection) Stat() *pgxpool.Stat{
	if p.conn != nil {
		stats := p.conn.Stat()
		return stats
	}
	return nil
}