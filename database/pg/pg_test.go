package pg

import (
	"testing"
	"context"
	"encoding/json"
	"time"
)

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

func TestGoCore_DatabasePGServer(t *testing.T){
	var databasePGServer DatabasePGServer

	databaseConfig := DatabaseConfig{
		Host: "127.0.0.1", 				
		Port: "5432", 				
		Schema:	"public",			
		DatabaseName: "postgres",		
		User: "postgres",				
		Password: "postgres",			
		Postgres_Driver: "postgres",		
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

	t.Logf("stats: %v", string(jsonBytes) )
}