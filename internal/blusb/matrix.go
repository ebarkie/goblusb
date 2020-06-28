// Copyright (c) 2020 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package blusb

// GetMatrix is a gets the matrix row and column for the current key being
// pressed.  This is non-blocking and if no key is being pressed it returns
// 0, 0.
func (c Controller) GetMatrix() (int, int, error) {
	data := make([]byte, 8)
	_, err := c.getControlReport(featMatrix, data)
	if err != nil {
		return 0, 0, err
	}

	return int(data[0]), int(data[1]), nil
}
