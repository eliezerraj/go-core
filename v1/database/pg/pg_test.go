package pg

import (
	"testing"
	"context"
	"encoding/json"
	"time"
)

func TestGoCore_DatabasePGServer(t *testing.T){
	var databasePGServer DatabasePGServer

	databaseConfig := DatabaseConfig{
		Host: "127.0.0.1", 				
		Port: "5432", 						
		DatabaseName: "postgres",		
		User: "postgres",				
		Password: "postgres",
		DbMax_Connection: 30,		
	}

	ctx, cancel := context.WithTimeout(	context.Background(), 
										time.Duration( 30 ) * time.Second)
	defer cancel()

	databasePG, err := databasePGServer.NewDatabasePGServer(ctx, databaseConfig)
	if err != nil {
		t.Errorf("failed to open database : %s", err)
	}

	_, conn, err := databasePG.StartTx(ctx)
	if err != nil {
		t.Errorf("failed to starttx : %s", err)
	}

	databasePG.ReleaseTx(conn)
	
	t.Logf("databasePG: %v", databasePG)

	stats := databasePG.Stat()

	statsJSON := PoolStats{
		AcquireCount:         stats.AcquireCount(),
		AcquiredConns:        stats.AcquiredConns(),
		CanceledAcquireCount: stats.CanceledAcquireCount(),
		ConstructingConns:    stats.ConstructingConns(),
		EmptyAcquireCount:    stats.EmptyAcquireCount(),
		IdleConns:            stats.IdleConns(),
		MaxConns:             stats.MaxConns(),
		TotalConns:           stats.TotalConns(),
	}

	// Convert to JSON
	jsonBytes, _ := json.MarshalIndent(statsJSON, "", "  ")

	t.Logf("stats: %v", string(jsonBytes))

	err = databasePG.Ping()
	if err != nil {
		t.Errorf("failed to ping : %s", err)
	}

	t.Logf("Ping Successful !!!")
}