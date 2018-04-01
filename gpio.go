package gpio

import (
	"log"
	"os"
	"sync"
	"syscall"
	"unsafe"
)

// PinMode type
type PinMode uint8

// Pin mode (pin can be set in Input or Output)
const (
	ModeIN PinMode = iota
	ModeOUT
)

// PinValue type
type PinValue uint8

// Represent state of pin, Low / High
const (
	PinLOW PinValue = iota
	PinHIGH
)

// Pin interface GPIO pin
type Pin interface {
	// set pin mode
	SetPinMode(PinMode)
	// set pin high
	High()
	// set pin low
	Low()
	// gets the current pin state
	Read() int
}

type pin struct {
	number uint8
	mode   PinMode
}

var (
	gpioMap []*uint32
	mem     []uint8

	lock sync.Mutex
)

func init() {
	loadMem()
}

func loadMem() {
	memFile, err := os.OpenFile("/dev/gpiomem", os.O_RDWR|os.O_SYNC, 0)
	if err != nil {
		log.Fatalf("Unable to open /dev/gpiomem: %v", err)
	}
	defer memFile.Close()

	lock.Lock()
	defer lock.Unlock()

	memFd := int(memFile.Fd())
	mem, err = syscall.Mmap(memFd, BCM2835_GPIO_BASE, BCM2835_BLOCK_SIZE, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		log.Fatalf("Unable to mmap GPIO page: %v", err)
	}

	gpioMap = []*uint32{
		(*uint32)(unsafe.Pointer(&mem[BCM2835_GPFSEL0])),
		(*uint32)(unsafe.Pointer(&mem[BCM2835_GPFSEL1])),
		(*uint32)(unsafe.Pointer(&mem[BCM2835_GPFSEL2])),
		(*uint32)(unsafe.Pointer(&mem[BCM2835_GPFSEL3])),
		(*uint32)(unsafe.Pointer(&mem[BCM2835_GPFSEL4])),
		(*uint32)(unsafe.Pointer(&mem[BCM2835_GPFSEL5])),
		(*uint32)(unsafe.Pointer(&mem[BCM2835_GPSET0])),
		(*uint32)(unsafe.Pointer(&mem[BCM2835_GPSET1])),
		(*uint32)(unsafe.Pointer(&mem[BCM2835_GPCLR0])),
		(*uint32)(unsafe.Pointer(&mem[BCM2835_GPCLR1])),
		(*uint32)(unsafe.Pointer(&mem[BCM2835_GPLEV0])),
		(*uint32)(unsafe.Pointer(&mem[BCM2835_GPLEV1])),
	}
}

// Setup pin
func Setup(number int, mode PinMode) Pin {
	p := &pin{
		number: uint8(number),
		mode:   mode,
	}
	p.setWithMode()
	return p
}

// High - sets the pin level high
func (p *pin) High() {
	offset := p.number/32 + getRegisterShift(PinHIGH)
	shift := p.number % 32
	*gpioMap[offset] = (1 << shift)
}

// Low - sets the pin level low
func (p *pin) Low() {
	offset := p.number/32 + getRegisterShift(PinLOW)
	shift := p.number % 32
	*gpioMap[offset] = (1 << shift)
}

// Output - set pin as output
func Output(number int, value PinValue) {
	pin := uint8(number)
	offset := pin/32 + getRegisterShift(value)
	shift := pin % 32
	*gpioMap[offset] = (1 << shift)
}

// SetPinMode (ModeIN or ModeOUT)
func (p *pin) SetPinMode(mode PinMode) {
	p.mode = mode
	p.setWithMode()
}

// Read - read pin state
func (p *pin) Read() int {
	offset := p.number/32 + 10
	shift := p.number % 32
	return int(*gpioMap[offset] & (1 << shift) / 4)
}

// Cleanup ...
func Cleanup() error {
	lock.Lock()
	defer lock.Unlock()

	return syscall.Munmap(mem)
}

func (p *pin) setWithMode() {
	offset := p.number / 10
	shift := (p.number % 10) * 3
	value := *gpioMap[offset]
	mask := BCM2835_GPIO_FSEL_MASK << shift
	value &= ^uint32(mask)
	value |= uint32(p.mode) << shift
	*gpioMap[offset] = value & mask
}

func getRegisterShift(value PinValue) uint8 {
	switch value {
	case PinLOW:
		return 8
	case PinHIGH:
		return 6
	default:
		return 0
	}
}
