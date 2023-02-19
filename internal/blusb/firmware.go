// Copyright (c) 2020 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package blusb

// EnterBoot signals the firmware to enter the bootloader.
func (c Controller) EnterBoot() error {
	data := make([]byte, 8)
	data[0] = firmEnterBoot
	return c.setControlReport(data)
}

// ExitBoot signals the bootloader to exit and boot the firmware.
func (c Controller) ExitBoot() error {
	data := make([]byte, 8)
	data[0] = bootExit
	return c.setControlReport(data)
}

func (c Controller) UpdateFirmware(filename string) error {
	return nil
}
