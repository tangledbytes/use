package snowflake

import "time"

// Snowflake is a struct that holds the state of the snowflake generator.
type Snowflake struct {
	epochOffset int64
	workerID    int64
	sequence    int64
	lastts      int64
}

// New returns a new Snowflake generator.
//
// If epochOffset < 0 then it is set to 0.
// If workerID < 0 or workerID > 1023 then it is set to 0.
func New(epochOffset int64, workerID int64) *Snowflake {
	if epochOffset < 0 {
		epochOffset = 0
	}

	if workerID < 0 || workerID > 1023 {
		workerID = 0
	}

	return &Snowflake{
		epochOffset: epochOffset,
		workerID:    workerID,
		sequence:    0,
		lastts:      0,
	}
}

// Next returns the next snowflake ID.
func (s *Snowflake) Next() uint64 {
	ts := time.Now().UnixMilli()

	// Handle with time going backwards or staying the same.
	// If time goes backwards, we wait until it goes forward.
	// If time stays the same, we increment the sequence.
	// The consequence is that we can only generate 4096 IDs per millisecond.
	// and if the skew lasts for a very long time then this function will block.
	if ts <= s.lastts {
		s.sequence = (s.sequence + 1) & 0xfff
		if s.sequence == 0 {
			for ts <= s.lastts {
				ts = time.Now().UnixMilli()
			}
		}
	} else {
		s.sequence = 0
	}

	s.lastts = ts

	id := uint64(ts-s.epochOffset) << 22 // 41 bits
	id |= uint64(s.workerID) << 12       // 10 bits
	id |= uint64(s.sequence)             // 12 bits

	return id
}
