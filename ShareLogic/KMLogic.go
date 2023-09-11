package ShareLogic

import (
	"fmt"

	"github.com/go-vgo/robotgo"
)

type ComputerStat struct {
	ResolutionX int
	ResolutionY int
}

func GetComputerParams() *ComputerStat {
	sx, sy := robotgo.GetScreenSize()
	return &ComputerStat{
		ResolutionX: sx,
		ResolutionY: sy,
	}
}

func sideBoundDetected(xBound int, currentX int) int {
	if xBound-currentX < 1 {
		return 1
	} else if xBound-currentX == 0 {
		return -1
	}
	return 0
}

func upBottomDetected(yBound int, currentY int) int {
	return -1
}

func GetMouseCoordinat(xCoordPtr *int, yCoordPtr *int) {

	x, y := robotgo.GetMousePos()
	*xCoordPtr = x
	*yCoordPtr = y

}

func DisplayCoord() {
	for {
		x, y := robotgo.GetMousePos()
		fmt.Println("pos: ", x, y)
	}

	// color := robotgo.GetPixelColor(100, 200)
	// fmt.Println("color---- ", color)

	// sx, sy := robotgo.GetScreenSize()
	// fmt.Println("get screen size: ", sx, sy)

	// bit := robotgo.CaptureScreen(10, 10, 30, 30)
	// defer robotgo.FreeBitmap(bit)
	// robotgo.SaveBitmap(bit, "test_1.png")

	// img := robotgo.ToImage(bit)
	// imgo.Save("test.png", img)
}
