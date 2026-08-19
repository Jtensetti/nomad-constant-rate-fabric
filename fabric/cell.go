package fabric

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"io"
)

const CellSize = 1200

// Cell is the fixed-size unit visible to the transport simulation.
// The protocol deliberately exposes no variable-length application payload.
type Cell [CellSize]byte

// FillRandom fills a cell with cryptographically secure random bytes.
func FillRandom(c *Cell) error {
	_, err := io.ReadFull(rand.Reader, c[:])
	return err
}

// StampEpoch places a test-only epoch marker into the first eight bytes.
// It is intended for local simulation and MUST NOT be used as a wire identifier.
func StampEpoch(c *Cell, epoch uint64) {
	binary.BigEndian.PutUint64(c[:8], epoch)
}

// ReadEpoch reads the test-only marker set by StampEpoch.
func ReadEpoch(c Cell) (uint64, error) {
	if len(c) != CellSize {
		return 0, errors.New("invalid cell size")
	}
	return binary.BigEndian.Uint64(c[:8]), nil
}
