package tools

import(
	"time"
	"github.com/rs/zerolog/log"
)

var childLogger = log.With().Str("go-core", "tools").Logger()

type ToolsCore struct {
}

func (t *ToolsCore) Test() string{
	childLogger.Debug().Msg("func ToolsCore - Test")
	return "test-01"
}

func (t *ToolsCore) ConvertToDate(date_str string) (*time.Time, error){
	childLogger.Debug().Str("func","ConvertToDate").Send()

	layout := "2006-01-02"

	date, err := time.Parse(layout, date_str)
	if err != nil {
		childLogger.Info().Err(err).Send()
		return nil, err
	}

	return &date, nil
}