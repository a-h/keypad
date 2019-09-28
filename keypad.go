package keypad

import (
	"time"

	"github.com/stianeikeland/go-rpio"
)

// New creates a key pad, where the pins 1-7 of the keypad are connected to the input
// GPIOs values.
func New(col1, col2, col3, col4, row1, row2, row3, row4 rpio.Pin) Pad {
	col1.Output()
	col1.Low()
	col2.Output()
	col2.Low()
	col3.Output()
	col3.Low()
	col4.Output()
	col4.Low()
	row1.Input()
	row2.Input()
	row3.Input()
	row4.Input()

	return Pad{
		col1:                col1,
		col2:                col2,
		col3:                col3,
		col4:                col4,
		row1:                row1,
		row2:                row2,
		row3:                row3,
		row4:                row4,
		previousPressedKeys: make([]bool, 16, 16),
		pressedKeys:         make([]bool, 16, 16),
		lastEvent:           make([]time.Time, 16, 16),
	}
}

// Pad is a 16 key pad.
// e.g. https://www.sparkfun.com/products/14881
type Pad struct {
	col1, col2, col3, col4, row1, row2, row3, row4 rpio.Pin
	previousPressedKeys                            []bool
	pressedKeys                                    []bool
	lastEvent                                      []time.Time
}

func (p *Pad) setPressedKeys() {
	for i := 0; i < len(p.pressedKeys); i++ {
		p.pressedKeys[i] = false
	}
	p.readColumn(p.col1, 0)
	p.readColumn(p.col2, 4)
	p.readColumn(p.col3, 8)
	p.readColumn(p.col4, 12)
}

func (p *Pad) readColumn(col rpio.Pin, offset int) {
	col.High()
	defer col.Low()
	if p.row1.Read() == rpio.High {
		p.pressedKeys[offset+0] = true
	}
	if p.row2.Read() == rpio.High {
		p.pressedKeys[offset+1] = true
	}
	if p.row3.Read() == rpio.High {
		p.pressedKeys[offset+2] = true
	}
	if p.row4.Read() == rpio.High {
		p.pressedKeys[offset+3] = true
	}
}

var allKeys = []string{"1", "4", "7", "*", "2", "5", "8", "0", "3", "6", "9", "#", "A", "B", "C", "D"}

func (p *Pad) Read() (keys []string, ok bool) {
	p.setPressedKeys()

	// Work out which keys have been pressed and released.
	for i, pressed := range p.pressedKeys {
		if pressed {
			// Continue to check.
			if !p.previousPressedKeys[i] {
				// Limit the keypress rate to 1 every 200 milliseconds.
				if time.Now().Before(p.lastEvent[i].Add(time.Millisecond * 200)) {
					continue
				}
				p.lastEvent[i] = time.Now()
			}
		} else {
			if p.previousPressedKeys[i] {
				// It was pressed before, but it isn't now.
				keys = append(keys, allKeys[i])
			}
		}
		p.previousPressedKeys[i] = pressed
	}

	return keys, len(keys) > 0
}
