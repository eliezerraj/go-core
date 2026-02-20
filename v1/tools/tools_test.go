package tools

import (
	"time"
	"testing"
)

func TestCore_Test(t *testing.T){
	var toolsCore ToolsCore
	toolsCore.Test()
}

func TestCore_ConvertDate(t *testing.T){
	var toolsCore ToolsCore

	date_str := "2021-01-01"
	var res *time.Time
	res, err := toolsCore.ConvertToDate(date_str)
	if err != nil {
		t.Errorf("failed to ConvertToDate : %s", err)
	}
	t.Logf("res: %v", res)
}