package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

type tile struct {
	X int
	Y int
}

type pt struct {
	x int
	y int
}

func main() {

	var tileWidth, tileHeight, multiplier int
	var xTiles, yTiles int
	var output string

	flag.IntVar(&tileWidth, "width", 128, "tile width in pixels")
	flag.IntVar(&tileHeight, "height", 128, "base tile height in pixels")
	flag.IntVar(&multiplier, "multiplier", 2, "tile height multiplier")
	flag.IntVar(&xTiles, "x", 8, "number of tiles across")
	flag.IntVar(&yTiles, "y", 8, "number of tiles down")
	flag.StringVar(&output, "output", "", "output filename")

	flag.Parse()

	if output == "" {
		flag.Usage()
		os.Exit(1)
	}

	t := tile{tileWidth, tileHeight}

	width := xTiles * t.X
	height := yTiles * t.Y

	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	for x := 0; x < xTiles; x++ {
		for y := 0; y < yTiles; y++ {
			drawDiamond(img, t, x, y, multiplier)
		}
	}

	// Encode as PNG.
	f, err := os.Create(output)
	if err != nil {
		fmt.Printf("could not output to %s - %s", output, err)
		os.Exit(1)
	}
	png.Encode(f, img)
	f.Close()
}

func drawDiamond(img *image.RGBA, t tile, x int, y int, multiplier int) {

	tileX := x * t.X
	tileY := y * t.Y

	yTop := t.Y - (t.Y / multiplier)
	yMid := t.Y - (t.Y / (multiplier * 2))

	pt1 := pt{tileX, tileY + yMid}           // left middle
	pt2 := pt{tileX + t.X/2, tileY + t.Y}    // bottom
	pt3 := pt{tileX + t.X - 1, tileY + yMid} // right
	pt4 := pt{tileX + t.X/2, tileY + yTop}   // top

	drawLine(img, pt1.x, pt1.y, pt2.x, pt2.y, color.Black)
	drawLine(img, pt2.x, pt2.y, pt3.x, pt3.y, color.Black)
	drawLine(img, pt3.x, pt3.y, pt4.x, pt4.y, color.Black)
	drawLine(img, pt4.x, pt4.y, pt1.x, pt1.y, color.Black)
	addLabel(img, tileX+t.X/2-10, tileY+30, fmt.Sprintf("%d,%d", x, y))

	// tl => bl
	drawLine(img, tileX, tileY, tileX, tileY+t.Y-1, color.White)
	// bl => br
	drawLine(img, tileX, tileY+t.Y-1, tileX+t.X-1, tileY+t.Y-1, color.White)
	// br => tr
	drawLine(img, tileX+t.X-1, tileY+t.Y-1, tileX+t.X-1, tileY, color.White)
	// tr => tl
	drawLine(img, tileX+t.X-1, tileY, tileX, tileY, color.White)

}

// thanks to https://github.com/StephaneBunel/bresenham
func drawLine(img *image.RGBA, x1, y1, x2, y2 int, col color.Color) {
	var dx, dy, e, slope int

	// Because drawing p1 -> p2 is equivalent to draw p2 -> p1,
	// I sort points in x-axis order to handle only half of possible cases.
	if x1 > x2 {
		x1, y1, x2, y2 = x2, y2, x1, y1
	}

	dx, dy = x2-x1, y2-y1
	// Because point is x-axis ordered, dx cannot be negative
	if dy < 0 {
		dy = -dy
	}

	switch {

	// Is line a point ?
	case x1 == x2 && y1 == y2:
		img.Set(x1, y1, col)

	// Is line an horizontal ?
	case y1 == y2:
		for ; dx != 0; dx-- {
			img.Set(x1, y1, col)
			x1++
		}
		img.Set(x1, y1, col)

	// Is line a vertical ?
	case x1 == x2:
		if y1 > y2 {
			y1, y2 = y2, y1
		}
		for ; dy != 0; dy-- {
			img.Set(x1, y1, col)
			y1++
		}
		img.Set(x1, y1, col)

	// Is line a diagonal ?
	case dx == dy:
		if y1 < y2 {
			for ; dx != 0; dx-- {
				img.Set(x1, y1, col)
				x1++
				y1++
			}
		} else {
			for ; dx != 0; dx-- {
				img.Set(x1, y1, col)
				x1++
				y1--
			}
		}
		img.Set(x1, y1, col)

	// wider than high ?
	case dx > dy:
		if y1 < y2 {
			// BresenhamDxXRYD(img, x1, y1, x2, y2, col)
			dy, e, slope = 2*dy, dx, 2*dx
			for ; dx != 0; dx-- {
				img.Set(x1, y1, col)
				x1++
				e -= dy
				if e < 0 {
					y1++
					e += slope
				}
			}
		} else {
			// BresenhamDxXRYU(img, x1, y1, x2, y2, col)
			dy, e, slope = 2*dy, dx, 2*dx
			for ; dx != 0; dx-- {
				img.Set(x1, y1, col)
				x1++
				e -= dy
				if e < 0 {
					y1--
					e += slope
				}
			}
		}
		img.Set(x2, y2, col)

	// higher than wide.
	default:
		if y1 < y2 {
			// BresenhamDyXRYD(img, x1, y1, x2, y2, col)
			dx, e, slope = 2*dx, dy, 2*dy
			for ; dy != 0; dy-- {
				img.Set(x1, y1, col)
				y1++
				e -= dx
				if e < 0 {
					x1++
					e += slope
				}
			}
		} else {
			// BresenhamDyXRYU(img, x1, y1, x2, y2, col)
			dx, e, slope = 2*dx, dy, 2*dy
			for ; dy != 0; dy-- {
				img.Set(x1, y1, col)
				y1--
				e -= dx
				if e < 0 {
					x1++
					e += slope
				}
			}
		}
		img.Set(x2, y2, col)
	}
}

func addLabel(img *image.RGBA, x, y int, label string) {
	col := color.RGBA{20, 20, 240, 255}
	point := fixed.Point26_6{X: fixed.I(x), Y: fixed.I(y)}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(label)
}

// // Colors are defined by Red, Green, Blue, Alpha uint8 values.
// cyan := color.RGBA{100, 200, 200, 0xff}

// // Set color for each pixel.
// for x := 0; x < width; x++ {
// 	for y := 0; y < height; y++ {
// 		switch {
// 		case x < width/2 && y < height/2: // upper left quadrant
// 			img.Set(x, y, cyan)
// 		case x >= width/2 && y >= height/2: // lower right quadrant
// 			img.Set(x, y, color.White)
// 		default:
// 			// Use zero value.
// 		}
// 	}
// }
