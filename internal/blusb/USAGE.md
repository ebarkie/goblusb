# blusb
```go
import "gitlab.com/ebarkie/goblusb/internal/blusb"
```

Package blusb implements protocol and data structures for configuring the
various settings of the Blusb Universal BT-USB Model M Controller.

## Usage

```go
var (
	// Controller Vendor ID
	VID gousb.ID = 0x04b3

	// Controller Product ID
	PID gousb.ID = 0x301c

	// Debug logger
	Debug *log.Logger = log.New(ioutil.Discard, "[DBUG] ", 0)
)
```
Defaults

```go
var (
	ErrInvalidBrightness  = errors.New("brightness value must be between 0 and 255")
	ErrInvalidDebounceDur = errors.New("debounce duration must be between 1ms and 255ms")
	ErrControllerNotFound = errors.New("blusb controller not found")
)
```
Errors

#### type Controller

```go
type Controller struct {

	// Skip set operations so nothing is changed (enabled with check mode)
	SkipSets bool
}
```

Controller holds the Blusb Universal BT-USB Model M Controller context.

#### func  Open

```go
func Open() (Controller, error)
```
Open opens the controller and claims the default interface.

#### func (*Controller) Close

```go
func (c *Controller) Close()
```
Close releases the controller.

#### func (Controller) EnterBoot

```go
func (c Controller) EnterBoot() error
```
EnterBoot signals the firmware to enter the bootloader.

#### func (Controller) ExitBoot

```go
func (c Controller) ExitBoot() error
```
ExitBoot signals the bootloader to exit and boot the firmware.

#### func (Controller) GetBrightness

```go
func (c Controller) GetBrightness() (uint, uint, error)
```
GetBrightness returns the Num Lock, Caps Lock, and Scroll Lock LED brightness
values stored in the controller for USB and Bluetooth modes. The value range is
0-255.

#### func (Controller) GetDebounce

```go
func (c Controller) GetDebounce() (time.Duration, error)
```
GetDebounce returns the debounce duration stored in the controller.

#### func (Controller) GetLayers

```go
func (c Controller) GetLayers() (ls Layers, err error)
```
GetLayers returns the layers stored in the controller.

#### func (Controller) GetMacros

```go
func (c Controller) GetMacros() (ms Macros, err error)
```
GetMacros returns the macro table stored in the controller.

#### func (Controller) GetMatrix

```go
func (c Controller) GetMatrix() (pos MatrixPos, err error)
```
GetMatrix is a gets the matrix row and column for the current key being pressed.
This is non-blocking and if no keys being pressed it returns a zero value.

#### func (Controller) GetVersion

```go
func (c Controller) GetVersion() (int, int, error)
```
GetVersion returns the controller firmware version as major and minor integers.
These are usually written as "major.minor".

#### func (Controller) MonitorMatrix

```go
func (c Controller) MonitorMatrix(ctx context.Context) <-chan MatrixPos
```
MonitorMatrix returns a channel and sends keypress matrix positions. A key that
is held down will be sent as a single keypress.

#### func (Controller) SetBrightness

```go
func (c Controller) SetBrightness(usb, bt uint) error
```
SetBrightness sets the controller Num Lock, Caps Lock, and Scroll Lock LED
brightness values for USB and Bluetooth modes. The value range is 0-255.

#### func (Controller) SetDebounce

```go
func (c Controller) SetDebounce(dur time.Duration) error
```
SetDebounce sets the controller debounce duration.

#### func (Controller) SetLayers

```go
func (c Controller) SetLayers(ls Layers) error
```
SetLayers sets the controller layers.

#### func (Controller) SetMacros

```go
func (c Controller) SetMacros(ms Macros) error
```
SetMacros sets the controller macro table.

#### func (Controller) String

```go
func (c Controller) String() string
```

#### type Layer

```go
type Layer struct {
	// Each layer has a matrixRows*matrixCols matrix, representing up to
	// 160 keys.
	//
	// Each key is 2-bytes with the higher byte representing any modifiers
	// and the lower byte being a key code.
	Matrix [matrixRows][matrixCols]uint16
}
```

Layer represents one layer.

#### func (Layer) String

```go
func (l Layer) String() string
```

#### type Layers

```go
type Layers []Layer
```

Layers represents all configured layers.

#### func (Layers) MarshalBinary

```go
func (ls Layers) MarshalBinary() (data []byte, err error)
```
MarshalBinary encodes a layers data packet. It consists of 1-byte to indicate
the number of layers and 160-bytes for each layer.

#### func (Layers) MarshalText

```go
func (ls Layers) MarshalText() ([]byte, error)
```
MarshalText composes CSV formatted layers consisting of one line for each layer
with each layer consisting of 160 hexadecimal combined modifier and key codes
separated by commas.

#### func (Layers) String

```go
func (ls Layers) String() string
```

#### func (*Layers) UnmarshalBinary

```go
func (ls *Layers) UnmarshalBinary(data []byte) error
```
UnmarshalBinary decodes a layers data packet. It consists of 1-byte to indicate
the number of layers and 160-bytes for each layer.

#### func (*Layers) UnmarshalText

```go
func (ls *Layers) UnmarshalText(text []byte) error
```
UnmarshalText parses CSV formatted layers consisting of one line for each layer
with each layer consisting of 160 hexadecimal combined modifier and key codes
separated by commas.

#### type Macro

```go
type Macro struct {
	Mods     uint8    // Internal modifier key codes
	Reserved uint8    // Reserved for future use
	Key      [6]uint8 // Up to 6 key codes
}
```

Macro represents one macro.

#### type Macros

```go
type Macros [numMacros]Macro
```

Macros represents the full macro table.

#### func (Macros) MarshalBinary

```go
func (ms Macros) MarshalBinary() (data []byte, err error)
```
MarshalBinary encodes a 193-byte macro table data packet.

#### func (Macros) MarshalText

```go
func (ms Macros) MarshalText() ([]byte, error)
```
MarshalText composes a CSV formatted macro table consisting of one line for each
macro with each macro consisting of it's parts encoded as hexadecimal codes and
separated by commas.

#### func (Macros) String

```go
func (ms Macros) String() string
```

#### func (*Macros) UnmarshalBinary

```go
func (ms *Macros) UnmarshalBinary(data []byte) error
```
UnmarshalBinary decodes a 192-byte macro table data packet.

#### func (*Macros) UnmarshalText

```go
func (ms *Macros) UnmarshalText(text []byte) error
```
UnmarshalText parses a CSV formatted macro table that consists of one line for
each macro with each macro consisting of it's parts encoded as hexadecimal codes
and separated by commas.

#### type MatrixPos

```go
type MatrixPos struct {
	Row, Col int
}
```

MatrixPos represents a keyboard matrix position.

#### func (MatrixPos) IsZero

```go
func (p MatrixPos) IsZero() bool
```
IsZero indicates if the position is empty.

#### func (MatrixPos) String

```go
func (p MatrixPos) String() string
```

#### func (*MatrixPos) UnmarshalBinary

```go
func (p *MatrixPos) UnmarshalBinary(data []byte) error
```
UnmarshalBinary decodes an 8-byte matrix report data packet.
