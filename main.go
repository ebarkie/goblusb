// Copyright (c) 2020 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package main

import (
	"context"
	"encoding"
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"gitlab.com/ebarkie/goblusb/internal/blusb"
)

const ok = "OK"

type uints struct {
	Want int
	S    []uint
}

func (u uints) String() string {
	s := make([]string, len(u.S))
	for i := range u.S {
		s[i] = strconv.FormatUint(uint64(u.S[i]), 10)
	}
	return strings.Join(s, ",")
}

func (u *uints) Set(value string) error {
	for _, s := range strings.Split(value, ",") {
		i, err := strconv.ParseUint(s, 10, 32)
		if err != nil {
			return err
		}
		u.S = append(u.S, uint(i))
	}

	if len(u.S) != u.Want {
		return fmt.Errorf("want %d values but got %d", u.Want, len(u.S))
	}

	return nil
}

func writeTextFile(v encoding.TextMarshaler, filename string) error {
	text, err := v.MarshalText()
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, text, 0644)
}

func main() {
	monitorMatrix := flag.Bool("monitor-matrix", false, "monitor for key presses")

	version := flag.Bool("version", false, "firmware version")

	getBright := flag.Bool("get-brightness", false, "get usb and bt brightness")
	getDebounce := flag.Bool("get-debounce", false, "get debounce duration")
	getLayers := flag.Bool("get-layers", false, "get layers")
	getMacros := flag.Bool("get-macros", false, "get macro keys")
	to := flag.String("to", "", "write to file")

	setBright := uints{Want: 2}
	flag.Var(&setBright, "set-brightness", "set usb,bt brightness")
	setDebounce := flag.Duration("set-debounce", 0, "set debounce duration")
	setLayers := flag.String("set-layers", "", "set layers from file")
	setMacros := flag.String("set-macros", "", "set macro keys fom file")

	flag.Parse()

	c, err := blusb.Open()
	if err != nil {
		fmt.Printf("Open device error: %s\n", err)
		return
	}
	defer c.Close()
	fmt.Printf("Blusb Controller - %s\n\n", c)

	// Exclusive operations.
	if *monitorMatrix {
		// Monitoring too quickly after pressing enter to start this
		// command bombards standard input with repeating carriage
		// returns so sleep a little bit.
		time.Sleep(500 * time.Millisecond)

		// XXX Safeguard in case something goes wrong and another
		// keyboard isn't handy.  Maybe remove this in the future.
		dur := 30 * time.Second
		fmt.Printf("Monitoring matrix for up to %s.  Press the same key twice in a row to exit sooner.\n\n", dur)

		ctx, cancel := context.WithTimeout(context.Background(), dur)
		defer cancel()
		var prevPos blusb.MatrixPos
		for pos := range c.MonitorMatrix(ctx) {
			fmt.Println(pos)

			if pos == prevPos {
				return
			}
			prevPos = pos
		}
	}

	// Combinable operations.
	if *version {
		maj, min, err := c.GetVersion()
		if err != nil {
			fmt.Printf("Get version error: %s\n", err)
			return
		}
		fmt.Printf("Version %d.%d\n", maj, min)
	}

	if *getBright {
		usb, bt, err := c.GetBrightness()
		if err != nil {
			fmt.Printf("Get brightness error: %s\n", err)
			return
		}
		fmt.Println("Brightness")
		fmt.Printf("\tUSB is %d/255\n", usb)
		fmt.Printf("\tBluetooth is %d/255\n", bt)
	}

	if *getDebounce {
		db, err := c.GetDebounce()
		if err != nil {
			fmt.Printf("Get debounce error: %s\n", err)
			return
		}
		fmt.Printf("Debounce time is %s\n", db)
	}

	if *getLayers {
		layers, err := c.GetLayers()
		if err != nil {
			fmt.Printf("Get layers error: %s\n", err)
			return
		}
		fmt.Printf("%s", layers)

		if *to != "" {
			if err := writeTextFile(layers, *to); err != nil {
				fmt.Printf("Save layers error: %s\n", err)
				return
			}
		}
	}

	if *getMacros {
		macros, err := c.GetMacros()
		if err != nil {
			fmt.Printf("Get macros error: %s\n", err)
			return
		}
		fmt.Printf("Macro key table:\n\n%s\n", macros)

		if *to != "" {
			if err := writeTextFile(macros, *to); err != nil {
				fmt.Printf("Save macros error: %s\n", err)
			}
		}
	}

	if len(setBright.S) > 0 {
		fmt.Println("Setting brightness")
		fmt.Printf("\tUSB to %d/255\n", setBright.S[0])
		fmt.Printf("\tBluetooth to %d/255\n", setBright.S[1])
		if err := c.SetBrightness(setBright.S[0], setBright.S[1]); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(ok)
		}
	}

	if *setDebounce > 0 {
		fmt.Printf("Setting debounce to %s\n", *setDebounce)
		if err := c.SetDebounce(*setDebounce); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(ok)
		}
	}

	if *setLayers != "" {
		text, err := ioutil.ReadFile(*setLayers)
		if err != nil {
			fmt.Printf("Set layers parse error: %s\n", err)
			return
		}

		var layers blusb.Layers
		if err := layers.UnmarshalText(text); err != nil {
			fmt.Printf("Set layers parse error: %s\n", err)
			return
		}
		fmt.Printf("Setting layers to:\n\n%s", layers)
		if err := c.SetLayers(layers); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(ok)
		}
	}

	if *setMacros != "" {
		text, err := ioutil.ReadFile(*setMacros)
		if err != nil {
			fmt.Printf("Set macros parse error: %s\n", err)
			return
		}

		var macros blusb.Macros
		if err := macros.UnmarshalText(text); err != nil {
			fmt.Printf("Set macros parse error: %s\n", err)
			return
		}
		fmt.Printf("Setting macros to:\n\n%s\n", macros)
		if err := c.SetMacros(macros); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(ok)
		}
	}
}
