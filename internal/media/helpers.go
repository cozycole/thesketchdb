package media

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"math"
	"net/http"

	"github.com/google/uuid"
	xdraw "golang.org/x/image/draw"
)

func cover(src image.Image, tw, th int) (image.Image, error) {
	b := src.Bounds()
	sw, sh := b.Dx(), b.Dy()
	if sw == 0 || sh == 0 {
		return nil, errors.New("invalid source dimensions")
	}

	scale := math.Max(float64(tw)/float64(sw), float64(th)/float64(sh))
	rw := maxInt(1, int(math.Round(float64(sw)*scale)))
	rh := maxInt(1, int(math.Round(float64(sh)*scale)))

	resized := image.NewRGBA(image.Rect(0, 0, rw, rh))
	xdraw.CatmullRom.Scale(resized, resized.Bounds(), src, b, xdraw.Over, nil)

	// center-crop to exact target
	x0 := (rw - tw) / 2
	y0 := (rh - th) / 2
	cropRect := image.Rect(x0, y0, x0+tw, y0+th)

	dst := image.NewRGBA(image.Rect(0, 0, tw, th))
	draw.Draw(dst, dst.Bounds(), resized, cropRect.Min, draw.Src)
	return dst, nil
}

func contain(src image.Image, tw, th int, fmt Format) (image.Image, error) {
	b := src.Bounds()
	sw, sh := b.Dx(), b.Dy()
	if sw == 0 || sh == 0 {
		return nil, errors.New("invalid source dimensions")
	}

	scale := math.Min(float64(tw)/float64(sw), float64(th)/float64(sh))
	rw := maxInt(1, int(math.Round(float64(sw)*scale)))
	rh := maxInt(1, int(math.Round(float64(sh)*scale)))

	resized := image.NewRGBA(image.Rect(0, 0, rw, rh))
	xdraw.CatmullRom.Scale(resized, resized.Bounds(), src, b, xdraw.Over, nil)

	// letterbox/pad to exact target
	dst := image.NewRGBA(image.Rect(0, 0, tw, th))

	// For JPEG, transparency becomes black anyway; we explicitly paint black.
	if fmt == FormatJPEG {
		draw.Draw(dst, dst.Bounds(), image.Black, image.Point{}, draw.Src)
	}

	offX := (tw - rw) / 2
	offY := (th - rh) / 2
	draw.Draw(dst, image.Rect(offX, offY, offX+rw, offY+rh), resized, image.Point{}, draw.Over)
	return dst, nil
}

// ---- encoding ----

func encode(img image.Image, fmt Format, quality int) ([]byte, string, error) {
	var buf bytes.Buffer

	switch fmt {
	case FormatJPEG:
		err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: clampInt(quality, 1, 100)})
		return buf.Bytes(), "image/jpeg", err

	case FormatPNG:
		enc := png.Encoder{CompressionLevel: png.DefaultCompression}
		err := enc.Encode(&buf, img)
		return buf.Bytes(), "image/png", err

	case FormatWEBP:
		return nil, "", errors.New("webp encoding not implemented (add a webp encoder or use jpeg/png)")

	default:
		return nil, "", errors.New("unknown format")
	}
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func clampInt(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// files

var mimeToExt = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
}

func GenerateFileName(file []byte) (string, error) {
	thumbnailId := uuid.New().String()
	thumbnailExtension, ok := mimeToExt[http.DetectContentType(file)]
	if !ok {
		return "", fmt.Errorf("Mime does not exist in extension table, bad file")
	}
	return thumbnailId + thumbnailExtension, nil
}
