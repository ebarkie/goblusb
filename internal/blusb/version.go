// Copyright (c) 2020 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package blusb

// GetVersion returns the controller firmware version as major and minor
// integers.  These are usually written as "major.minor".
func (c Controller) GetVersion() (int, int, error) {
	data := make([]byte, 8)
	_, err := c.getControlReport(featVersion, data)
	if err != nil {
		return 0, 0, err
	}

	return int(data[0]), int(data[1]), nil
}
