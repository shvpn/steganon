package utils

import (
	"fmt"
	"image"
	"image/color"
)

// EncodeMessageInImage hides a message in an image using LSB steganography
func EncodeMessageInImage(img image.Image, message string) *image.RGBA {
	bounds := img.Bounds()
	encoded := image.NewRGBA(bounds)

	// Copy original image
	copyImage(encoded, img, bounds)

	// Prepare data to hide (length header + message)
	dataToHide := prepareDataWithHeader(message)

	// Encode data into image pixels
	encodeDataInPixels(encoded, dataToHide, bounds)

	return encoded
}

// DecodeMessageFromImage extracts a hidden message from an image
func DecodeMessageFromImage(img image.Image) string {
	bounds := img.Bounds()

	// Extract message length from header
	messageLen := extractMessageLength(img, bounds)
	if messageLen <= 0 || messageLen > 1000000 {
		return ""
	}

	// Extract the actual message
	messageBytes := extractMessageData(img, bounds, messageLen)

	return string(messageBytes)
}

// copyImage copies all pixels from source to destination
func copyImage(dst *image.RGBA, src image.Image, bounds image.Rectangle) {
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			dst.Set(x, y, src.At(x, y))
		}
	}
}

// prepareDataWithHeader adds length header to message data
func prepareDataWithHeader(message string) []byte {
	messageBytes := []byte(message)
	messageLen := len(messageBytes)
	// 8-byte length header + message data
	return append([]byte(fmt.Sprintf("%08d", messageLen)), messageBytes...)
}

// encodeDataInPixels encodes data bits into image pixels using LSB
func encodeDataInPixels(encoded *image.RGBA, dataToHide []byte, bounds image.Rectangle) {
	bitIndex := 0
	totalBits := len(dataToHide) * 8

	for y := bounds.Min.Y; y < bounds.Max.Y && bitIndex < totalBits; y++ {
		for x := bounds.Min.X; x < bounds.Max.X && bitIndex < totalBits; x++ {
			r, g, b, a := encoded.At(x, y).RGBA()

			// Convert to 8-bit values
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)
			a8 := uint8(a >> 8)

			// Encode bits in RGB channels (3 bits per pixel)
			r8 = encodeBitInChannel(r8, dataToHide, &bitIndex, totalBits)
			g8 = encodeBitInChannel(g8, dataToHide, &bitIndex, totalBits)
			b8 = encodeBitInChannel(b8, dataToHide, &bitIndex, totalBits)

			encoded.SetRGBA(x, y, color.RGBA{R: r8, G: g8, B: b8, A: a8})
		}
	}
}

// encodeBitInChannel encodes a single bit into a color channel
func encodeBitInChannel(channel uint8, data []byte, bitIndex *int, totalBits int) uint8 {
	if *bitIndex < totalBits {
		byteIndex := *bitIndex / 8
		bitOffset := *bitIndex % 8
		bit := (data[byteIndex] >> (7 - bitOffset)) & 1
		channel = (channel & 0xFE) | bit
		*bitIndex++
	}
	return channel
}

// extractMessageLength extracts the 8-byte length header from the image
func extractMessageLength(img image.Image, bounds image.Rectangle) int {
	lengthBytes := make([]byte, 8)
	bitIndex := 0

	for y := bounds.Min.Y; y < bounds.Max.Y && bitIndex < 64; y++ {
		for x := bounds.Min.X; x < bounds.Max.X && bitIndex < 64; x++ {
			r, g, b, _ := img.At(x, y).RGBA()

			// Extract LSB from each channel
			extractBitFromChannel(lengthBytes, &bitIndex, uint8(r>>8), 64)
			extractBitFromChannel(lengthBytes, &bitIndex, uint8(g>>8), 64)
			extractBitFromChannel(lengthBytes, &bitIndex, uint8(b>>8), 64)
		}
	}

	// Parse the length
	var messageLen int
	fmt.Sscanf(string(lengthBytes), "%08d", &messageLen)
	return messageLen
}

// extractMessageData extracts the actual message data from the image
func extractMessageData(img image.Image, bounds image.Rectangle, messageLen int) []byte {
	messageBytes := make([]byte, messageLen)
	totalBits := messageLen * 8
	bitIndex := 0

	for y := bounds.Min.Y; y < bounds.Max.Y && bitIndex < totalBits; y++ {
		for x := bounds.Min.X; x < bounds.Max.X && bitIndex < totalBits; x++ {
			r, g, b, _ := img.At(x, y).RGBA()

			// Skip the first 64 bits (length header)
			currentBit := bitIndex + 64

			// Calculate pixel position for current bit
			pixelBit := currentBit % 3
			pixelIndex := currentBit / 3

			if pixelIndex != (x-bounds.Min.X)+(y-bounds.Min.Y)*(bounds.Max.X-bounds.Min.X) {
				continue
			}

			if pixelBit == 0 && bitIndex < totalBits {
				extractBitFromChannel(messageBytes, &bitIndex, uint8(r>>8), totalBits)
			}

			if pixelBit <= 1 && bitIndex < totalBits {
				extractBitFromChannel(messageBytes, &bitIndex, uint8(g>>8), totalBits)
			}

			if pixelBit <= 2 && bitIndex < totalBits {
				extractBitFromChannel(messageBytes, &bitIndex, uint8(b>>8), totalBits)
			}
		}
	}

	return messageBytes
}

// extractBitFromChannel extracts a bit from a color channel
func extractBitFromChannel(data []byte, bitIndex *int, channel uint8, maxBits int) {
	if *bitIndex < maxBits {
		byteIndex := *bitIndex / 8
		bitOffset := *bitIndex % 8
		bit := channel & 1
		data[byteIndex] |= bit << (7 - bitOffset)
		*bitIndex++
	}
}
