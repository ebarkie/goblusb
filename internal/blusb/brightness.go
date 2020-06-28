// Copyright (c) 2020 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package blusb

// GetBrightness returns the Num Lock, Caps Lock, and Scroll Lock LED
// brightness values stored in the controller for USB and Bluetooth modes.
// The value range is 0-255.
func (c Controller) GetBrightness() (uint, uint, error) {
	data := make([]byte, 8)
	_, err := c.getControlReport(featBrightness, data)
	if err != nil {
		return 0, 0, err
	}

	return uint(data[0]), uint(data[1]), nil
}

// SetBrightness sets the controller Num Lock, Caps Lock, and Scroll Lock LED
// brightness values for USB and Bluetooth modes.  The value range is 0-255.
func (c Controller) SetBrightness(usb, bt uint) error {
	if usb > 255 || bt > 255 {
		return ErrInvalidBrightness
	}

	data := make([]byte, 8)
	copy(data, []byte{featBrightness, byte(usb), byte(bt)})
	return c.setControlReport(data)
}
