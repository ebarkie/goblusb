// Copyright (c) 2020 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package blusb

import "time"

// GetDebounce returns the debounce duration stored in the controller.
func (c Controller) GetDebounce() (time.Duration, error) {
	data := make([]byte, 8)
	_, err := c.getControlReport(firmDebounce, data)
	if err != nil {
		return 0, err
	}

	return time.Duration(data[0]) * time.Millisecond, nil
}

// SetDebounce sets the controller debounce duration.
func (c Controller) SetDebounce(dur time.Duration) error {
	if dur < 1*time.Millisecond || dur > 255*time.Millisecond {
		return ErrInvalidDebounceDur
	}

	data := make([]byte, 8)
	copy(data, []byte{firmDebounce, byte(dur.Milliseconds())})
	return c.setControlReport(data)
}
