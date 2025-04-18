package pg

import (
	"context"
	"fmt"
	"time"
	"errors"
	"github.com/rs/zerolog/log"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var childLogger = log.With().Str("component","go-core").Str("package", "database.pg").Logger()

type DatabaseConfig struct {
    Host 				string `json:"host"`
    Port  				string `json:"port"`
	Schema				string `json:"schema"`
	DatabaseName		string `json:"databaseName"`
	User				string `json:"user"`
	Password			string `json:"password"`
	Db_timeout			int	`json:"db_timeout"`
	Postgres_Driver		string `json:"postgres_driver"`
}

type DatabasePG interface {
	GetConnection() (*pgxpool.Pool)
	Acquire(context.Context) (*pgxpool.Conn, error)
	Release(*pgxpool.Conn)
	CloseConnection()
}

type DatabasePGServer struct {
	connPool   	*pgxpool.Pool
}

func Config(database_url string) (*pgxpool.Config) {
	childLogger.Debug().Str("func","Config").Send()

	const defaultMaxConns = int32(10)
	const defaultMinConns = int32(2)
	const defaultMaxConnLifetime = time.Hour
	const defaultMaxConnIdleTime = time.Minute * 10
	const defaultHealthCheckPeriod = time.Minute
	const defaultConnectTimeout = time.Second * 5
   
	dbConfig, err := pgxpool.ParseConfig(database_url)
	if err!=nil {
		childLogger.Error().Err(err).Send()
	}
   
	dbConfig.MaxConns = defaultMaxConns
	dbConfig.MinConns = defaultMinConns
	dbConfig.MaxConnLifetime = defaultMaxConnLifetime
	dbConfig.MaxConnIdleTime = defaultMaxConnIdleTime
	dbConfig.HealthCheckPeriod = defaultHealthCheckPeriod
	dbConfig.ConnConfig.ConnectTimeout = defaultConnectTimeout
   
	dbConfig.BeforeAcquire = func(ctx context.Context, c *pgx.Conn) bool {
		childLogger.Debug().Msg("Before acquiring connection pool !")
	 	return true
	}
   
	dbConfig.AfterRelease = func(c *pgx.Conn) bool {
		childLogger.Debug().Msg("After releasing connection pool !")
	 	return true
	}
   
	dbConfig.BeforeClose = func(c *pgx.Conn) {
		childLogger.Debug().Msg("Closed connection pool !")
	}
   
	return dbConfig
}

func (d DatabasePGServer) NewDatabasePGServer(ctx context.Context, databaseConfig DatabaseConfig) (DatabasePGServer, error) {
	childLogger.Debug().Str("func","NewDatabasePGServer").Send()
	
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", 
							databaseConfig.User, 
							databaseConfig.Password, 
							databaseConfig.Host, 
							databaseConfig.Port, 
							databaseConfig.DatabaseName) 
							
	connPool, err := pgxpool.NewWithConfig(ctx, Config(connStr))
	if err != nil {
		return DatabasePGServer{}, err
	}
	
	err = connPool.Ping(ctx)
	if err != nil {
		return DatabasePGServer{}, err
	}

	return DatabasePGServer{
		connPool: connPool,
	}, nil
}

func (d DatabasePGServer) Acquire(ctx context.Context) (*pgxpool.Conn, error) {
	childLogger.Debug().Str("func","NewDatabasePGServer").Send()
	
	connection, err := d.connPool.Acquire(ctx)
	if err != nil {
		childLogger.Error().Err(err).Msg("Error while acquiring connection from the database pool!!")
		return nil, err
	} 

	return connection, nil
}

func (d DatabasePGServer) Release(connection *pgxpool.Conn) {
	childLogger.Debug().Str("func","Release").Send()

	defer connection.Release()
}

func (d DatabasePGServer) GetConnection() (*pgxpool.Pool) {
	childLogger.Debug().Str("func","GetConnection").Send()

	return d.connPool
}

func (d DatabasePGServer) CloseConnection() {
	childLogger.Debug().Str("func","CloseConnection").Send()

	defer d.connPool.Close()
}

func (d DatabasePGServer) StartTx(ctx context.Context) (pgx.Tx, *pgxpool.Conn, error) {
	childLogger.Debug().Str("func","StartTx").Send()

	conn, err := d.Acquire(ctx)
	if err != nil {
		childLogger.Error().Err(err).Msg("error acquire")
		return nil, nil, errors.New(err.Error())
	}
	
	tx, err := conn.Begin(ctx)
    if err != nil {
        return nil, nil, errors.New(err.Error())
    }

	return tx, conn, nil
}

func (d DatabasePGServer) ReleaseTx(connection *pgxpool.Conn) {
	childLogger.Debug().Str("func","ReleaseTx").Send()

	defer connection.Release()
}