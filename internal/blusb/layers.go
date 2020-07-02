// Copyright (c) 2020 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package blusb

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const (
	matrixRows = 8
	matrixCols = 20
)

// Layer represents one layer.
type Layer struct {
	// Each layer has a matrixRows*matrixCols matrix, representing up to
	// 160 keys.
	//
	// Each key is 2-bytes with the higher byte representing any modifiers
	// and the lower byte being a key code.
	Matrix [matrixRows][matrixCols]uint16
}

func (l Layer) String() string {
	buf := &bytes.Buffer{}
	buf.WriteString("    ")
	for c := range l.Matrix[0] {
		fmt.Fprintf(buf, "C%-3d  ", c)
	}
	buf.WriteString("\n\n")

	for r := range l.Matrix {
		fmt.Fprintf(buf, "R%-1d  ", r)
		for c := range l.Matrix[r] {
			fmt.Fprintf(buf, "%04X  ", l.Matrix[r][c])
		}
		buf.WriteByte('\n')
	}

	return buf.String()
}

// Layers represents all configured layers.
type Layers []Layer

func (ls Layers) String() string {
	buf := &bytes.Buffer{}
	for i, l := range ls {
		fmt.Fprintf(buf, "Layer %d/%d\n\n%s\n", i+1, len(ls), l)
	}

	return buf.String()
}

// MarshalBinary encodes a layers data packet.  It consists of 1-byte to
// indicate the number of layers and 160-bytes for each layer.
func (ls Layers) MarshalBinary() (data []byte, err error) {
	//		          1                   2                   3
	//    0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//   |    # Layers   | Layer 1 Row 0 Col 0 Key Code  | Row 0, Col..  |
	//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//   | 1 Key Code    |     Row 0, Col 2 Key Code     | ...           |
	//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//   | ...           | Layer 2 Row 0 Col 0 Key Code  | Row 0, Col..  |
	//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//   | 1 Key Code    |     Row 0, Col 2 Key Code     | ...           |
	//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	data = make([]byte, 0, 1+len(ls)*matrixRows*matrixCols)
	data = append(data, byte(len(ls)))
	for _, l := range ls {
		for r := range l.Matrix {
			for _, kc := range l.Matrix[r] {
				data = append(data, byte(kc), byte(kc>>8))
			}
		}
	}

	return
}

// UnmarshalBinary decodes a layers data packet.  It consists of 1-byte to
// indicate the number of layers and 160-bytes for each layer.
func (ls *Layers) UnmarshalBinary(data []byte) error {
	var l Layer
	var r, c int
	var numLayers = int(data[0])
	for i := 1; i < len(data); i += 2 {
		l.Matrix[r][c] = uint16(data[i+1])<<8 | uint16(data[i])
		if c < matrixCols-1 {
			c++
		} else if r < matrixRows-1 {
			r++
			c = 0
		} else {
			*ls = append(*ls, l)
			if len(*ls) >= numLayers {
				break
			}
			r, c = 0, 0
		}
	}

	return nil
}

// MarshalText composes CSV formatted layers consisting of one line for each
// layer with each layer consisting of 160 hexadecimal combined modifier and key
// codes separated by commas.
func (ls Layers) MarshalText() ([]byte, error) {
	buf := &bytes.Buffer{}
	for _, l := range ls {
		for r := range l.Matrix {
			for c := range l.Matrix[r] {
				if r > 0 || c > 0 {
					buf.WriteString(", ")
				}
				fmt.Fprintf(buf, "%X", l.Matrix[r][c])
			}
		}
		buf.WriteByte('\n')
	}

	return buf.Bytes(), nil
}

// UnmarshalText parses CSV formatted layers consisting of one line for each
// layer with each layer consisting of 160 hexadecimal combined modifier and
// key codes separated by commas.
func (ls *Layers) UnmarshalText(text []byte) error {
	var l Layer
	var r, c int
	s := bufio.NewScanner(bytes.NewReader(text))
	s.Split(bufio.ScanWords)
	for i := 0; s.Scan(); i++ {
		u, err := strconv.ParseUint(strings.Trim(s.Text(), ", "), 16, 16)
		if err != nil {
			return err
		}

		l.Matrix[r][c] = uint16(u)
		if c < matrixCols-1 {
			c++
		} else if r < matrixRows-1 {
			r++
			c = 0
		} else {
			*ls = append(*ls, l)
			r, c = 0, 0
		}
	}

	return nil
}

const (
	layersPageHeadSize = 0x3
	layersPageDataSize = 0x100
	layersPageSize     = layersPageHeadSize + layersPageDataSize
)

type layersPager struct {
	*bytes.Buffer

	// Layers page header
	//
	//		          1                   2                   3
	//    0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//   |       ID     |  Total pages   | Current Page  | Data          |
	//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	prevPageRead  int
	prevPageWrite int
}

func (p *layersPager) Reset() {
	p.Buffer.Reset()

	p.prevPageRead = 0
	p.prevPageWrite = 0
}

func (p *layersPager) Read(b []byte) (int, error) {
	if p.Len() < 1 {
		return 0, io.EOF
	}

	totalPages := p.Cap()/(len(b)-layersPageHeadSize) + 1
	p.prevPageRead++
	h := []byte{firmLayers, byte(totalPages), byte(p.prevPageRead)}
	if n := copy(b[:len(h)], h); n != layersPageHeadSize {
		return n, io.ErrShortBuffer
	}

	n, err := p.Buffer.Read(b[layersPageHeadSize:])
	n += layersPageHeadSize
	if err != nil {
		return n, err
	}

	// Pad the last page with 0xff's.
	if n < len(b) {
		n += copy(b[n:], bytes.Repeat([]byte{0xff}, len(b)-n))
	}

	return n, nil
}

func (p *layersPager) Write(b []byte) (int, error) {
	const (
		idOff          = 0
		totalPagesOff  = 1
		currentPageOff = 2
	)

	if len(b) < layersPageHeadSize {
		return 0, io.ErrShortWrite
	}
	// XXX ID should be firmLayers but it comes through as zero
	if (b[idOff] != 0 && b[idOff] != firmLayers) ||
		int(b[currentPageOff]) != p.prevPageWrite+1 {
		return 0, io.ErrNoProgress
	}
	p.prevPageWrite++

	n, err := p.Buffer.Write(b[layersPageHeadSize:])
	n += layersPageHeadSize
	if err != nil {
		return n, err
	}

	if b[totalPagesOff] == b[currentPageOff] {
		return n, io.EOF
	}

	return n, nil
}

// GetLayers returns the layers stored in the controller.
func (c Controller) GetLayers() (ls Layers, err error) {
	// Read all layer pages into a buffer.
	p := layersPager{Buffer: &bytes.Buffer{}}
	for {
		page := make([]byte, layersPageSize)
		_, err = c.getControlReport(firmLayers, page)
		if err != nil {
			return
		}

		_, err = p.Write(page)
		if err != nil {
			if err == io.EOF {
				break
			}
			return
		}
	}

	// Unmarshal the buffer into layers.
	err = ls.UnmarshalBinary(p.Bytes())
	return
}

// SetLayers sets the controller layers.
func (c Controller) SetLayers(ls Layers) error {
	data, err := ls.MarshalBinary()
	if err != nil {
		return err
	}

	// Create layer pages from data and write to controller.
	p := layersPager{Buffer: bytes.NewBuffer(data)}
	for {
		page := make([]byte, layersPageSize)
		_, err := p.Read(page)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		if err := c.setControlReport(page); err != nil {
			return err
		}
	}

	return nil
}
