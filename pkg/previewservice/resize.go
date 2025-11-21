package previewservice

import (
	"image"
	"image/jpeg"
	"log"
	"os"

	"golang.org/x/image/draw"
)

// Scale creates a resized image from the path and writes it to the resizedImagePath.
// The desired width and height of the scaled image is specified in pixels,
// and the resulting width and height will be calculated
// to preserve the aspect ratio.
// If the ratio does not the same (in and out images):
// in this case, one of the desired dimensions
// (outWidth or outHeight) will be retained.
// If the size of incoming image is smaller than the desired,
// incoming image will be returned as a response.
func Scale(path string, resizedImagePath string, outWidth int, outHeight int) error {
	log.Println("path to scale:", path)
	// Open file.
	inputFile, err := os.Open(path)
	if err != nil {
		return err
	}

	// -------------------------------------------------------------------------

	originalImage, err := jpeg.Decode(inputFile)
	if err != nil {
		return err
	}

	inHeight := originalImage.Bounds().Max.Y
	inWidth := originalImage.Bounds().Max.X

	ok := checkSize(outWidth, outHeight, inWidth, inHeight)
	if !ok {
		return &ResizeError{} // fmt.Errorf("nothing to resize")
	}

	checkHeight := inHeight * outWidth / inWidth
	checkWidth := inWidth * outHeight / inHeight

	if checkHeight <= outHeight {
		outHeight = checkHeight
	} else if checkWidth <= outWidth {
		outWidth = checkWidth
	}

	dst := image.NewRGBA(image.Rect(0, 0, outWidth, outHeight))

	draw.NearestNeighbor.Scale(dst, dst.Rect, originalImage, originalImage.Bounds(), draw.Over, nil)

	scaledImageFile, err := os.Create(resizedImagePath)
	if err != nil {
		return err
	}
	defer scaledImageFile.Close() //nolint

	err = jpeg.Encode(scaledImageFile, dst, nil)
	if err != nil {
		return err
	}

	return nil
}

type ResizeError struct {
}

func (e *ResizeError) Error() string {
	return "nothing to resize"
}

type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}

// Crop crops jpg image from the path and writes it to the resizedImagePath.
// The desired width and height of the cropped image is specified in pixels
// (outWidth and outHeight).
// The resulting size is cropped from the center of the incomming image.
func Crop(path string, resizedImagePath string, outWidth int, outHeight int) error {
	log.Println("path", path)
	// Open file.
	inputFile, err := os.Open(path)
	if err != nil {
		return err
	}

	// -------------------------------------------------------------------------

	originalImage, err := jpeg.Decode(inputFile)
	if err != nil {
		return err
	}

	bounds := originalImage.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	xMin := (width - outWidth) / 2
	xMax := xMin + outWidth
	yMin := (height - outHeight) / 2
	yMax := yMin + outHeight

	cropSize := image.Rect(xMin, yMin, xMax, yMax)

	croppedImage := originalImage.(SubImager).SubImage(cropSize)

	// -------------------------------------------------------------------------

	croppedImageFile, err := os.Create(resizedImagePath)
	if err != nil {
		return err
	}
	defer croppedImageFile.Close() //nolint

	err = jpeg.Encode(croppedImageFile, croppedImage, nil)
	if err != nil {
		return err
	}

	return nil
}

// checkSize checks if the desired size of result image less than destination image
func checkSize(outWidth, outHeight, inWidth, inHeight int) bool {
	if outWidth <= inWidth && outHeight <= inHeight {
		return true
	}
	return false
}
