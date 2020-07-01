// Copyright (c) 2020 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

// Package blusb implements protocol and data structures for configuring the
// various settings of the Blusb Universal BT-USB Model M Controller.
package blusb

import (
	"fmt"
	"io"

	"github.com/google/gousb"
)

// Bitmap request types
const (
	reqGetReport uint8 = 0x01
	reqSetReport uint8 = 0x09
)

// Specific requests
const (
	repIn uint16 = (iota + 1) << 8
	repOut
	repFeature
)

// Bootloader features
const (
	bootPageData byte = iota + 0x01 // Send page data (firmware)
	bootExit                        // Exit bootloader and boot the firmware
)

// Firmware features
const (
	firmLayers     byte = iota + 0x01 // Get or set layers
	firmMacros                        // Get or set macros
	firmMatrix                        // Read matrix input
	firmBrightness                    // Get or set brightness
	firmGetVersion                    // Get firmware version
	firmDebounce                      // Get or set debounce duration
	firmEnterBoot                     // Enter bootloader
)

// Default Vendor and Product ID's.
var (
	VID gousb.ID = 0x04b3
	PID gousb.ID = 0x301c
)

// Controller holds the Blusb Universal BT-USB Model M Controller context.
type Controller struct {
	ctx  *gousb.Context // USB device handling context
	dev  *gousb.Device  // Opened USB device
	done func()         // Release claimed interface and config
}

// Open opens the controller and claims the default interface.
func Open() (Controller, error) {
	ctx := gousb.NewContext()

	dev, err := ctx.OpenDeviceWithVIDPID(VID, PID)
	if dev == nil {
		return Controller{}, ErrControllerNotFound
	}
	if err != nil {
		return Controller{}, err
	}

	// We only send control requests and don't care about the interface but
	// Linux doesn't like it when we don't claim it.
	_, done, err := dev.DefaultInterface()
	if err != nil {
		return Controller{}, err
	}

	return Controller{
		ctx:  ctx,
		dev:  dev,
		done: done,
	}, nil
}

// Close releases the controller.
func (c *Controller) Close() {
	c.done()
	c.dev.Close()
	c.ctx.Close()
}

func (c Controller) String() string {
	m, _ := c.dev.Manufacturer()
	p, _ := c.dev.Product()

	return fmt.Sprintf("Bus: %d Address: %d Manufacturer: %s (%s) Product: %s (%s)",
		c.dev.Desc.Bus, c.dev.Desc.Address, m, c.dev.Desc.Vendor, p, c.dev.Desc.Product)
}

func (c Controller) getControlReport(feat byte, b []byte) (int, error) {
	n, err := c.dev.Control(gousb.ControlInterface|gousb.ControlIn|gousb.ControlClass,
		reqGetReport, repFeature|uint16(feat), 0, b)
	if err != nil {
		return n, err
	}
	if n != len(b) {
		return n, io.EOF
	}

	return n, nil
}

func (c Controller) setControlReport(b []byte) error {
	_, err := c.dev.Control(gousb.ControlInterface|gousb.ControlOut|gousb.ControlClass,
		reqSetReport, repFeature|uint16(b[0]), 0, b)
	return err
}
