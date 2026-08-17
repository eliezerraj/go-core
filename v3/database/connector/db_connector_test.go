package connector

import(
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
)

func TestNewDatabaseConnector(t *testing.T) {
	application := "go-core-v3-test"

	host := "localhost"
	port := "5432"
	dbName := "postgres"

	user_read := "postgres"
	password_read := "postgres"
	
	readerConfig := ConnectorConfig{
		DSN:              "postgres://" + user_read + ":" + password_read + "@" + host + ":" + port + "/" + dbName,
		MaxConnIdleTime:  5 * time.Minute,
		MaxConnLifeTime:  30 * time.Minute,
		MaxConns:         10,
		MinConns:         1,
	}

	user_write := "postgres"
	password_write := "postgres"

	writerConfig := ConnectorConfig{
		DSN:              "postgres://" + user_write + ":" + password_write + "@" + host + ":" + port + "/" + dbName,
		MaxConnIdleTime:  5 * time.Minute,
		MaxConnLifeTime:  30 * time.Minute,
		MaxConns:         10,
		MinConns:         1,
	}
	
	ctx := context.Background()

	connector, err := NewDatabaseConnector(application, readerConfig, writerConfig)
	if err != nil {
		t.Fatalf("Failed to create database connector: %v", err)
	}

	t.Logf("connector: %+v", connector)

	pgConnection := &PgConnection{}
	pool, err := pgConnection.NewPool(ctx, readerConfig)
	if err != nil {
		t.Errorf("Failed to create database pool: %v", err)
	}
	if pool == nil {
		t.Errorf("Expected non-nil pool, got nil")
	}

	err = pgConnection.Ping(ctx)
	if err != nil {
		t.Errorf("Failed to ping database pool: %v", err)
	} else {
		t.Logf("Successfully pinged the database pool")
	}

	query := "SELECT 1+1"
	rows, err := pool.Query(ctx, query)
	if err != nil {
		t.Errorf("Failed to execute query: %v", err)
	}
	defer rows.Close()

	var result int
	if rows.Next() {
		err = rows.Scan(&result)
		if err != nil {
			t.Errorf("Failed to scan result: %v", err)
		} else {
			t.Logf("Query result: %d", result)
		}
	}

	connectorReader := connector.Reader()
	if connectorReader == nil {
		t.Errorf("Expected non-nil reader, got nil")
	}
	rows, err = connectorReader.Query(ctx, "SELECT 2+2")
	if err != nil {
		t.Errorf("Failed to execute query: %v", err)
	} else {
		defer rows.Close()
		if rows.Next() {
			err = rows.Scan(&result)
			if err != nil {
				t.Errorf("Failed to scan result: %v", err)
			} else {
				t.Logf("Query result: %d", result)
			}
		}
	}

	connectorWriter := connector.Writer()
	if connectorWriter == nil {
		t.Errorf("Expected non-nil writer, got nil")
	}
	rows, err = connectorWriter.Query(ctx, "SELECT 3+3")
	if err != nil {
		t.Errorf("Failed to execute query: %v", err)
	} else {
		defer rows.Close()
		if rows.Next() {
			err = rows.Scan(&result)
			if err != nil {
				t.Errorf("Failed to scan result: %v", err)
			} else {
				t.Logf("Query result: %d", result)
			}
		}
	}

	tx, err := pgConnection.BeginTx(ctx)
	if err != nil {
		t.Errorf("Failed to begin transaction: %v", err)
	} else {
		t.Logf("Successfully began a transaction")
		err = tx.Rollback(ctx)
		if err != nil {
			t.Errorf("Failed to rollback transaction: %v", err)
		} else {
			t.Logf("Successfully rolled back the transaction")
		}
	}

	tx, err = connectorWriter.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		t.Errorf("Failed to begin transaction: %v", err)
	} else {
		t.Logf("Successfully began a transaction")
		err = tx.Rollback(ctx)
		if err != nil {
			t.Errorf("Failed to rollback transaction: %v", err)
		} else {
			t.Logf("Successfully rolled back the transaction")
		}
	}

	tx.Commit(ctx)
}