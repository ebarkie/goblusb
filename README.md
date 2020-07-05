![Push](https://github.com/ebarkie/goblusb/workflows/Push/badge.svg)

# go-blusb

Go package and Command Line Interface for interacting with the Blusb Universal
BT-USB Model M Controller.

Things that it doesn't do but would be nice:

* Unit tests!
* Update firmware
* Enter/exit bootloader
* Break down modifier and keycodes into clearer structs and types instad of
  uint's.
* Marshal/unmarshal to alternate file formats

## Installation

```sh
$ go get github.com/ebarkie/goblusb
```

## Usage

```
Usage of ./goblusb:
  -check
    	don't actually set anything
  -debug
    	enable extra debug output
  -get-brightness
    	get usb and bt brightness
  -get-debounce
    	get debounce duration
  -get-layers
    	get layers
  -get-macros
    	get macro keys
  -monitor-matrix
    	monitor for key presses
  -set-brightness value
    	set usb,bt brightness
  -set-debounce duration
    	set debounce duration
  -set-layers string
    	set layers from file
  -set-macros string
    	set macro keys fom file
  -to string
    	write to file
  -version
    	firmware version
```

## License

Copyright (c) 2020 Eric Barkie. All rights reserved.  
Use of this source code is governed by the MIT license
that can be found in the [LICENSE](LICENSE) file.
