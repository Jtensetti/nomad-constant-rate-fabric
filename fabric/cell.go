package fabric

import (
	"crypto/rand"
	"io"
)

const CellSize = 1200

// Cell is the fixed-size unit emitted by the research traffic shaper.
type Cell [CellSize]byte

func RandomCell() (Cell, error) {
	var c Cell
	_, err := io.ReadFull(rand.Reader, c[:])
	return c, err
}
