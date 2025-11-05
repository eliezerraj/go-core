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

type PoolStats struct {
	AcquireCount        int64 `json:"acquire_count"`
	AcquiredConns       int32 `json:"acquired_conns"`
	CanceledAcquireCount int64 `json:"canceled_acquire_count"`
	ConstructingConns   int32 `json:"constructing_conns"`
	EmptyAcquireCount   int64 `json:"empty_acquire_count"`
	IdleConns           int32 `json:"idle_conns"`
	MaxConns            int32 `json:"max_conns"`
	TotalConns          int32 `json:"total_conns"`
}

type DatabaseConfig struct {
    Host 				string `json:"host"`
    Port  				string `json:"port"`
	DatabaseName		string `json:"databaseName"`
	User				string `json:"user"`
	Password			string `json:"password"`
	DbMax_Connection	int	`json:"db_max_connection"`
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

func Config(database_url string, maxConns ...int) (*pgxpool.Config) {
	childLogger.Debug().Str("func","Config").Send()

	// Default max connection
	var_max_conns := 10

	if len(maxConns) > 0 {
		var_max_conns = maxConns[0]
	}

	defaultMaxConns := int32(var_max_conns)
	const defaultMinConns = int32(2)
	const defaultMaxConnLifetime = 30 * time.Minute
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

// About create a database service
func (d DatabasePGServer) NewDatabasePGServer(ctx context.Context, databaseConfig DatabaseConfig) (DatabasePGServer, error) {
	childLogger.Debug().Str("func","NewDatabasePGServer").Send()
	
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", 
							databaseConfig.User, 
							databaseConfig.Password, 
							databaseConfig.Host, 
							databaseConfig.Port, 
							databaseConfig.DatabaseName) 
							
	connPool, err := pgxpool.NewWithConfig(ctx, Config(connStr, databaseConfig.DbMax_Connection))
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

// About acquire connection from pool
func (d DatabasePGServer) Acquire(ctx context.Context) (*pgxpool.Conn, error) {
	childLogger.Debug().Str("func","NewDatabasePGServer").Send()
	
	connection, err := d.connPool.Acquire(ctx)
	if err != nil {
		childLogger.Error().Err(err).Msg("Error while acquiring connection from the database pool!!")
		return nil, err
	} 

	return connection, nil
}

// About release connection
func (d DatabasePGServer) Release(connection *pgxpool.Conn) {
	childLogger.Debug().Str("func","Release").Send()

	defer connection.Release()
}

// About close a get connection
func (d DatabasePGServer) GetConnection() (*pgxpool.Pool) {
	childLogger.Debug().Str("func","GetConnection").Send()

	return d.connPool
}

// About close a connection
func (d DatabasePGServer) CloseConnection() {
	childLogger.Debug().Str("func","CloseConnection").Send()

	defer d.connPool.Close()
}

// About start a transaction
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

// About release the connection to pool connection
func (d DatabasePGServer) ReleaseTx(connection *pgxpool.Conn) {
	childLogger.Debug().Str("func","ReleaseTx").Send()

	defer connection.Release()
}

// About get Stats from database
func (d DatabasePGServer) Stat() (*pgxpool.Stat){
	childLogger.Debug().Str("func","Stat").Send()

	return d.connPool.Stat()
}