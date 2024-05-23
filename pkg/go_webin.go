package pkg

import (
	"bytes"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"

	"github.com/disintegration/imaging"
	"github.com/nickalie/go-webpbin"
)

func Decode(img io.Reader) (io.Reader, error) {
	i, _, err := image.Decode(img)
	if err != nil {
		return nil, err
	}
	iFit := imaging.Fit(i, 800, 400, imaging.Lanczos)

	var buff bytes.Buffer

	if os.Getenv("ENVIRONMENT") == "PRODUCTION" {
		cwebp := webpbin.NewCWebP(webpbin.SetVendorPath("/usr/local/bin/"), webpbin.SetSkipDownload(true))
		err := cwebp.Quality(40).InputImage(iFit).Output(&buff).Run()
		if err != nil {
			return nil, err
		}
	} else {
		enc := webpbin.Encoder{
			Quality: 40,
		}
		if err := enc.Encode(&buff, i); err != nil {
			return nil, err
		}
	}	

	return &buff, nil
}