package snowflakego

import (
	"errors"
	"sync"
	"time"
)

const (
	epoch         int64 = 1740304858525
	timeBits            = 41
	sequenceBits        = 12
	machineIDBits       = 10

	machineIDMax = -1 ^ (-1 << machineIDBits)
	sequenceMax  = -1 ^ (-1 << sequenceBits)
	timeMax      = -1 ^ (-1 << timeBits)
)

var (
	ErrInvalidMachineID   = errors.New("invalid machine id")
	ErrTimestampIsInvalid = errors.New("timestamp is invalid")
)

type Snowflake struct {
	mu            sync.Mutex
	lastTimestamp int64
	sequence      int64
	machineID     int64
}

func New(machineID int64) (*Snowflake, error) {

	if machineID > machineIDMax {
		return nil, ErrInvalidMachineID
	}

	return &Snowflake{
		machineID:     machineID,
		sequence:      0,
		lastTimestamp: 0,
	}, nil
}

func (s *Snowflake) getMilliSeconds() int64 {
	return time.Now().UnixMilli()
}

func (s *Snowflake) GenerateID() (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.nextID()
}

func (s *Snowflake) nextID() (int64, error) {

	timestamp := s.getMilliSeconds()

	if timestamp < s.lastTimestamp {
		return 0, ErrTimestampIsInvalid
	}

	if timestamp == s.lastTimestamp {
		s.sequence = (s.sequence + 1) & sequenceMax

		if s.sequence == 0 {
			// Sequence overflow, wait for the next millisecond
			for timestamp <= s.lastTimestamp {
				timestamp = s.getMilliSeconds()
			}
		}
	} else {
		s.sequence = 0
	}

	s.lastTimestamp = timestamp

	timestamp -= epoch

	id := (timestamp << (machineIDBits + sequenceBits)) |
		(s.machineID << sequenceBits) |
		s.sequence

	return id, nil
}
