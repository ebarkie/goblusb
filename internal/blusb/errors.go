// Copyright (c) 2020 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package blusb

import "errors"

// Errors
var (
	ErrInvalidBrightness  = errors.New("brightness value must be between 0 and 255")
	ErrInvalidDebounceDur = errors.New("debounce duration must be between 1ms and 255ms")
	ErrControllerNotFound = errors.New("blusb controller not found")
)
