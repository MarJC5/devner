// mkicon rasterises an SVG into a macOS-menubar-friendly template
// PNG: opaque where the SVG was black, transparent where the SVG was
// white. That inversion matches how SetTemplateIcon expects the alpha
// channel (opaque = tinted by theme, transparent = lets menubar show
// through), producing a filled "badge" with the glyph cut out.
//
// Supersample at 2× then downscale with a simple box filter so
// anti-aliased edges convert to partial alpha cleanly.
//
// Usage: go run ./cmd/mkicon <input.svg> <output.png>
package main

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"

	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
)

func main() {
	in, out := os.Args[1], os.Args[2]

	data, err := os.ReadFile(in)
	if err != nil {
		panic(err)
	}
	icon, err := oksvg.ReadIconStream(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}

	// Target output 44×44 @2x. Supersample at 2× → render at 88×88,
	// then downscale for smooth alpha transitions at small size.
	const finalSize = 44
	const ss = 2
	const renderSize = finalSize * ss

	icon.SetTarget(0, 0, renderSize, renderSize)
	rgba := image.NewRGBA(image.Rect(0, 0, renderSize, renderSize))
	// White background so the rasteriser has something to draw over;
	// we'll flip white → transparent in the alpha conversion below.
	draw.Draw(rgba, rgba.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)
	icon.Draw(rasterx.NewDasher(
		renderSize, renderSize,
		rasterx.NewScannerGV(renderSize, renderSize, rgba, rgba.Bounds()),
	), 1)

	// Downscale + convert luminance → alpha.
	dst := image.NewRGBA(image.Rect(0, 0, finalSize, finalSize))
	for y := 0; y < finalSize; y++ {
		for x := 0; x < finalSize; x++ {
			var sum int
			for sy := 0; sy < ss; sy++ {
				for sx := 0; sx < ss; sx++ {
					r, g, b, _ := rgba.At(x*ss+sx, y*ss+sy).RGBA()
					// Luminance (8-bit) via average — R/G/B all
					// equal for greyscale input, so a simple avg is
					// fine.
					lum := (r + g + b) / (3 * 256)
					// Flip: 0 (black) → alpha 255, 255 (white) → alpha 0.
					sum += 255 - int(lum)
				}
			}
			a := uint8(sum / (ss * ss))
			dst.SetRGBA(x, y, color.RGBA{0, 0, 0, a})
		}
	}

	f, err := os.Create(out)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := png.Encode(f, dst); err != nil {
		panic(err)
	}
}
