go_gpio Library
=======

Simple Go library for accessing [GPIO](https://www.raspberrypi.org/documentation/usage/gpio-plus-and-raspi2/)
on [Raspberry Pi](https://www.raspberrypi.org).

### Usage

```go
import (
    gpio "github.com/DoZator/go_gpio"
)
```

To setup GPIO pin as an output:

```go
pin := gpio.Setup(18, gpio.ModeOUT)
```

To configure GPIO pin as an input:

```go
pin := gpio.Setup(18, gpio.ModeIN)
```

To set the output state of a GPIO pin:

```go
gpio.Output(18, gpio.PinHIGH)
```

Set pin High

```go
pin.High()
```

Set pin Low

```go
pin.Low()
```

*Good luck in coding :)*


