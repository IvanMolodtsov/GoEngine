package primitives

import (
	"image"
	"image/png"
	"os"
)

type Image struct {
	Data image.Image
}

func LoadImage(path string) *Image {
	var (
		i    Image
		file *os.File
		err  error
	)
	file, err = os.Open(path)
	if err != nil {
		// Handle error (e.g., log it, return)
		panic(err)
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		// Handle error (e.g., log it, return)
		panic(err)
	}
	i.Data = img
	return &i
}

func (img *Image) GetPixel(u, v, l float64) uint32 {
	w := float64(img.Data.Bounds().Dx())
	u = min(max(u*w, 0), w)
	h := float64(img.Data.Bounds().Dy())
	v = min(max(v*h, 0), h)
	color := img.Data.At(int(u), int(v))
	return ToHex(color, l)
}
