package pg

import (
	"testing"
	"context"
	"time"
)

func TestGoCore_DatabasePGServer(t *testing.T){
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

	databasePG, err := NewDatabasePGServer(ctx, databaseConfig)
	if err != nil {
		t.Errorf("failed to open database : %s", err)
	}

	t.Logf("databasePG: %v", databasePG)
}