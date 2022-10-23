package snowflake

import (
	"testing"
	"time"
)

func TestSnowflake_Next(t *testing.T) {
	t.Run("integrity test", func(t *testing.T) {
		const second = 1000
		const runForSeconds = 5

		s := New(0, 0)

		ids := make(map[uint64]struct{})
		now := time.Now().UnixMilli()

		for time.Now().UnixMilli() < now+runForSeconds*second {
			id := s.Next()
			if _, ok := ids[id]; ok {
				t.Errorf("duplicate ID: %d", id)
			}
			ids[id] = struct{}{}
		}
	})
}
