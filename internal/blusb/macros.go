// Copyright (c) 2020 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package blusb

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
)

const (
	macroSize = 8  // Size of each macro in bytes
	numMacros = 24 // Total number of macros
)

// Macro represents one macro.
type Macro struct {
	Mods     uint8    // Internal modifier key codes
	Reserved uint8    // Reserved for future use
	Key      [6]uint8 // Up to 6 key codes
}

// Macros represents the full macro table.
type Macros [numMacros]Macro

// MarshalBinary encodes a 193-byte macro table data packet.
func (ms Macros) MarshalBinary() (data []byte, err error) {
	data = make([]byte, 1+len(ms)*macroSize)
	data[0] = firmMacros
	for i := range ms {
		data[1+i*macroSize] = ms[i].Mods
		data[1+i*macroSize+1] = ms[i].Reserved
		copy(data[1+i*macroSize+2:i*macroSize+8], ms[i].Key[:])
	}

	return
}

// UnmarshalBinary decodes a 192-byte macro table data packet.
func (ms *Macros) UnmarshalBinary(data []byte) error {
	// XXX Why isn't the ID included but the data size has
	// space for it?
	for i := range ms {
		(*ms)[i].Mods = data[i*macroSize]
		(*ms)[i].Reserved = data[i*macroSize+1]
		copy((*ms)[i].Key[:], data[i*macroSize+2:i*macroSize+8])
	}

	return nil
}

// MarshalText composes a CSV formatted macro table consisting of one line for
// each macro with each macro consisting of it's parts encoded as hexadecimal
// codes and separated by commas.
func (ms Macros) MarshalText() ([]byte, error) {
	buf := &bytes.Buffer{}
	for _, m := range ms {
		fmt.Fprintf(buf, "%X, %X", m.Mods, m.Reserved)
		for _, k := range m.Key {
			fmt.Fprintf(buf, ", %X", k)
		}
		buf.WriteByte('\n')
	}

	return buf.Bytes(), nil
}

// UnmarshalText parses a CSV formatted macro table that consists of one line
// for each macro with each macro consisting of it's parts encoded as
// hexadecimal codes and separated by commas.
func (ms *Macros) UnmarshalText(text []byte) error {
	s := bufio.NewScanner(bytes.NewReader(text))
	for i := 0; s.Scan(); i++ {
		macroParts := strings.Split(s.Text(), ",")
		if len(macroParts) != macroSize {
			return csv.ErrFieldCount
		}

		for p := range macroParts {
			u, err := strconv.ParseUint(strings.TrimSpace(macroParts[p]), 16, 8)
			if err != nil {
				return err
			}

			switch p {
			case 0:
				(*ms)[i].Mods = uint8(u)
			case 1:
				(*ms)[i].Reserved = uint8(u)
			default:
				(*ms)[i].Key[p-2] = uint8(u)
			}
		}
	}

	return nil
}

func (ms Macros) String() string {
	buf := &bytes.Buffer{}
	buf.WriteString("     MODS  RSVD  KEY1  KEY2  KEY3  KEY4  KEY5  KEY6\n\n")

	for i, m := range ms {
		fmt.Fprintf(buf, "M%02d  %02X    %02X", i+1, m.Mods, m.Reserved)
		for _, k := range m.Key {
			fmt.Fprintf(buf, "    %02X", k)
		}
		buf.WriteByte('\n')
	}

	return buf.String()
}

// GetMacros returns the macro table stored in the controller.
func (c Controller) GetMacros() (ms Macros, err error) {
	data := make([]byte, len(ms)*macroSize)
	_, err = c.getControlReport(firmMacros, data)
	if err != nil {
		return
	}

	err = ms.UnmarshalBinary(data)
	return
}

// SetMacros sets the controller macro table.
func (c Controller) SetMacros(ms Macros) error {
	data, err := ms.MarshalBinary()
	if err != nil {
		return err
	}

	return c.setControlReport(data)
}
