package media

import (
	"image"
	"image/color"
	"image/draw"
)

// RemoveBordersRobust removes near-uniform borders from an image.
// tol: tolerance for pixel deviation from border color (0-255).
// safetyCrop: optional extra crop in pixels to remove residual lines.
func RemoveBorders(img image.Image, tol uint8, safetyCrop int) image.Image {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	if width == 0 || height == 0 {
		return img
	}

	// Estimate border color from corners (top-left, top-right, bottom-left, bottom-right)
	corners := []color.Gray{}
	cornerSize := 5
	for y := bounds.Min.Y; y < bounds.Min.Y+cornerSize && y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Min.X+cornerSize && x < bounds.Max.X; x++ {
			corners = append(corners, rgbToGray(img.At(x, y)))
		}
		for x := bounds.Max.X - cornerSize; x < bounds.Max.X; x++ {
			if x >= bounds.Min.X {
				corners = append(corners, rgbToGray(img.At(x, y)))
			}
		}
	}
	for y := bounds.Max.Y - cornerSize; y < bounds.Max.Y; y++ {
		if y < bounds.Min.Y {
			continue
		}
		for x := bounds.Min.X; x < bounds.Min.X+cornerSize && x < bounds.Max.X; x++ {
			corners = append(corners, rgbToGray(img.At(x, y)))
		}
		for x := bounds.Max.X - cornerSize; x < bounds.Max.X; x++ {
			if x >= bounds.Min.X {
				corners = append(corners, rgbToGray(img.At(x, y)))
			}
		}
	}

	borderColor := medianGray(corners)

	// Find crop bounds
	minX, minY := bounds.Max.X, bounds.Max.Y
	maxX, maxY := bounds.Min.X, bounds.Min.Y

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			pix := rgbToGray(img.At(x, y))
			if absInt(int(pix.Y)-int(borderColor.Y)) > int(tol) {
				if x < minX {
					minX = x
				}
				if x > maxX {
					maxX = x
				}
				if y < minY {
					minY = y
				}
				if y > maxY {
					maxY = y
				}
			}
		}
	}

	// Nothing found, return original
	if minX > maxX || minY > maxY {
		return img
	}

	// Apply safety crop
	minX = max(bounds.Min.X, minX+safetyCrop)
	minY = max(bounds.Min.Y, minY+safetyCrop)
	maxX = min(bounds.Max.X, maxX-safetyCrop)
	maxY = min(bounds.Max.Y, maxY-safetyCrop)

	cropped := image.NewRGBA(image.Rect(0, 0, maxX-minX+1, maxY-minY+1))
	draw.Draw(cropped, cropped.Bounds(), img, image.Point{X: minX, Y: minY}, draw.Src)
	return cropped
}

// helpers

func rgbToGray(c color.Color) color.Gray {
	r, g, b, _ := c.RGBA()
	return color.Gray{Y: uint8((r>>8 + g>>8 + b>>8) / 3)}
}

func medianGray(vals []color.Gray) color.Gray {
	if len(vals) == 0 {
		return color.Gray{Y: 0}
	}
	tmp := make([]uint8, len(vals))
	for i, v := range vals {
		tmp[i] = v.Y
	}
	quickSelect(tmp, len(tmp)/2)
	return color.Gray{Y: tmp[len(tmp)/2]}
}

func absInt(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Quickselect for median
func quickSelect(a []uint8, k int) {
	if len(a) <= 1 {
		return
	}
	pivot := a[len(a)/2]
	left, right := 0, len(a)-1
	for i := 0; i <= right; {
		if a[i] < pivot {
			a[i], a[left] = a[left], a[i]
			i++
			left++
		} else if a[i] > pivot {
			a[i], a[right] = a[right], a[i]
			right--
		} else {
			i++
		}
	}
	if k < left {
		quickSelect(a[:left], k)
	} else if k > right {
		quickSelect(a[right+1:], k-right-1)
	}
}
