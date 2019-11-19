package utils

import (
	"bytes"
	"encoding/base64"
	"image/png"

	"github.com/skip2/go-qrcode"
)

func GenerateQrcode(s *string) (string, bool) {
	q, _ := qrcode.New(*s, qrcode.Medium)
	img := q.Image(256)
	buffer := new(bytes.Buffer)
	if err := png.Encode(buffer, img); err != nil {
		return "", false
	}

	return base64.StdEncoding.EncodeToString(buffer.Bytes()), true
}
