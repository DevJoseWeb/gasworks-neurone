/*
 * Copyright (c) Clinton Freeman 2013
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and
 * associated documentation files (the "Software"), to deal in the Software without restriction,
 * including without limitation the rights to use, copy, modify, merge, publish, distribute,
 * sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or
 * substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT
 * NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
 * NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
 * DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */
package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	// "github.com/huin/goserial"
	// "log"
	// "time"
	"io"
)

// updateArduinoEnergy transmits a new energy level over the nominated serial port to the arduino. Returns an error
// on failure, nil otherwise. Arduino code takes the energy level and turns it into a lighting sequence.
func updateArduinoEnergy(energy float32, serialPort io.ReadWriteCloser) error {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, energy)
	if err != nil {
		return err
	}

	_, err = serialPort.Write(buf.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func Axon(delta_e chan float32) {
	// Connect to the arduino over serial.
	// c := &goserial.Config{Name: "/dev/tty.usbserial-A1017HU2", Baud: 9600}
	// s, err := goserial.OpenPort(c)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// When connecting to an arduino, you need to wait a little while it resets.
	// time.Sleep(1 * time.Second)

	var energy float32
	energy = 0.0

	for i := 0; i < 500; i++ {
		de := <-delta_e
		energy += de

		// Neuron has reached threshold. Fire axon.
		if energy > 1.0 {

			// Feed into the dendrites of adjacent neurons.
			//req, err := http.NewRequest("GET", "http://neuronip:8080/?e=%d", nil)

			energy = 0.0
		}

		// updateArduinoEnergy(energy, s)
		fmt.Printf("Energy level %f %f\n", energy, de)
	}
}