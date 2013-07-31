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

// #cgo darwin CFLAGS: -I/opt/local/include -I/opt/local/include/opencv
// #cgo linux CFLAGS: -I/usr/include -I/usr/include/opencv
// #cgo CFLAGS: -Wno-error
// #cgo darwin LDFLAGS: -L/opt/local/lib
// #cgo linux LDFLAGS: -L/usr/lib
// #cgo LDFLAGS: -lopencv_highgui -lopencv_core -lopencv_video
// #include "cv.h"
// #include "highgui.h"
import "C"

 import (	
	"fmt"
	"math"
	"unsafe"
)

func calcDeltaEnergy(flow *C.IplImage) float64 {
	var i C.int
	var dx, dy float64

	// Accumulate the change in flow across all the pixels.
	totalPixels := flow.width * flow.height
	for i = 0; i < totalPixels; i++ {
		value := C.cvGet2D(unsafe.Pointer(flow), i/flow.width, i%flow.width)
		dx += math.Abs(float64(value.val[0]))
		dy += math.Abs(float64(value.val[1]))
	}

	// average out the magnitude of dx and dy across the whole image.
	dx = dx / float64(totalPixels)
	dy = dy / float64(totalPixels)

	// The magnitude of accumulated flow forms our change in energy for the frame.
	deltaE := math.Sqrt((dx * dx) + (dy * dy))
	fmt.Printf("flow %f \n", deltaE)

	// Clamp the energy to start at 0 for 'still' frames with little/no motion.
	deltaE = math.Max(0.0, (deltaE - 1.0))

	// Scale the flow to be less than 0.1 for 'active' frames with lots of motion.
	deltaE = deltaE / 1000.0

	return deltaE
}

func DendriteCam() {
	camera := C.cvCaptureFromCAM(-1)
	C.cvSetCaptureProperty(camera, C.CV_CAP_PROP_FRAME_WIDTH, 160)
	C.cvSetCaptureProperty(camera, C.CV_CAP_PROP_FRAME_HEIGHT, 120)

	// Capture original frame.
	prev := C.cvCloneImage(C.cvQueryFrame(camera))

	// file := C.CString("a.png")
	// C.cvSaveImage(file, unsafe.Pointer(prev), nil)
	// C.free(unsafe.Pointer(file))

	flow := C.cvCreateImage(C.cvSize(prev.width, prev.height), C.IPL_DEPTH_32F, 2)
	prevG := C.cvCreateImage(C.cvSize(prev.width, prev.height), C.IPL_DEPTH_8U, 1)
	nextG := C.cvCreateImage(C.cvSize(prev.width, prev.height), C.IPL_DEPTH_8U, 1)
	C.cvConvertImage(unsafe.Pointer(prev), unsafe.Pointer(prevG), 0)

	for i := 1; i < 500; i++ {
		C.cvGrabFrame(camera)

		// Capture the new frame and convert it to grayscale.
		next := C.cvCloneImage(C.cvQueryFrame(camera))
		C.cvConvertImage(unsafe.Pointer(prev), unsafe.Pointer(prevG), 0)
		C.cvConvertImage(unsafe.Pointer(next), unsafe.Pointer(nextG), 0)

		C.cvCalcOpticalFlowFarneback(unsafe.Pointer(prevG), unsafe.Pointer(nextG), unsafe.Pointer(flow), 0.5, 2, 5, 2, 5, 1.1, 0)
		delta := float32(calcDeltaEnergy(flow))
		//energy += delta
		fmt.Printf("energy: %f \n", delta)
		//		updateArduinoEnergy(energy, s)

		C.cvReleaseImage(&prev)
		prev = next
	}

	C.cvReleaseImage(&prev)

	C.cvReleaseImage(&nextG)
	C.cvReleaseImage(&prevG)
	C.cvReleaseImage(&flow)
	C.cvReleaseCapture(&camera)
}