package core

import(
	"log"
)

type ToolsCore struct {
}

func (t *ToolsCore) Test() string{
	log.Println("func ToolsCore - Test")

	return "test-01"
}