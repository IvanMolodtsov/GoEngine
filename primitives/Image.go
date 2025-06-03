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

func (img *Image) GetPixel(u, v float64) uint32 {
	// println(img.Data.Bounds().Dx(), img.Data.Bounds().Dy())
	U := int(u * float64(img.Data.Bounds().Dx()))
	V := int(v * float64(img.Data.Bounds().Dy()))
	// println(U, V)
	color := img.Data.At(U, V)
	// println(color.RGBA())
	return toHex(color)
}
