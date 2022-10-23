package id

import (
	"github.com/utkarsh-pro/use/pkg/config"
	"github.com/utkarsh-pro/use/pkg/id/snowflake"
)

type Gen interface {
	// Next returns the next ID.
	Next() uint64
}

// New returns a new ID generator.
func New() Gen {
	return snowflake.New(0, int64(config.WorkerID))
}
