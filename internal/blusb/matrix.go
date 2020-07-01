// Copyright (c) 2020 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package blusb

import (
	"context"
	"fmt"
)

// MatrixPos represents a keyboard matrix position.
type MatrixPos struct {
	Row, Col int
}

// IsZero indicates if the position is empty.
func (p MatrixPos) IsZero() bool { return p == MatrixPos{} }

func (p MatrixPos) String() string {
	return fmt.Sprintf("Row %2d, Col %2d", p.Row, p.Col)
}

// UnmarshalBinary decodes an 8-byte matrix report data packet.
func (p *MatrixPos) UnmarshalBinary(data []byte) error {
	p.Row, p.Col = int(data[0]), int(data[1])
	return nil
}

// GetMatrix is a gets the matrix row and column for the current key being
// pressed.  This is non-blocking and if no keys being pressed it returns a
// zero value.
func (c Controller) GetMatrix() (pos MatrixPos, err error) {
	data := make([]byte, 8)
	_, err = c.getControlReport(firmMatrix, data)
	if err != nil {
		return
	}

	err = pos.UnmarshalBinary(data)
	return
}

// MonitorMatrix returns a channel and sends keypress matrix positions.
// A key that is held down will be sent as a single keypress.
func (c Controller) MonitorMatrix(ctx context.Context) <-chan MatrixPos {
	ch := make(chan MatrixPos)
	go func() {
		defer close(ch)

		var prevPos MatrixPos
		for {
			pos, err := c.GetMatrix()
			if err != nil {
				return
			}

			if !pos.IsZero() && pos != prevPos {
				select {
				case ch <- pos:
				case <-ctx.Done():
					return
				}
			} else {
				select {
				case <-ctx.Done():
					return
				default:
				}
			}

			prevPos = pos
		}
	}()

	return ch
}
