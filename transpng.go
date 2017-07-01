package transpng

import (
	"image"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"io"
	"math"
)

// A Decoder reads and decodes image from an input stream.
type Decoder struct {
	r   io.Reader
	img image.Image
}

// NewDecoder returns a new decoder that reads from r.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		r: r,
	}
}

// Decode reads the Transparency-PNG image from its input and stores it in
// the value pointed to by w.
func (dec *Decoder) Decode(w io.Writer) error {
	if dec.img == nil {
		img, _, err := image.Decode(dec.r)
		if err != nil {
			return err
		}
		dec.img = img
	}

	var dst *image.NRGBA
	switch img := dec.img.(type) {
	case *image.NRGBA:
		dst = img
	default:
		dst = image.NewNRGBA(dec.img.Bounds())
		draw.Draw(dst, dst.Bounds(), dec.img, image.ZP, draw.Src)
	}

	p := dst.Bounds().Max
	x, y := p.X-1, p.Y-1
	c := dst.NRGBAAt(x, y)
	if c.A < math.MaxUint8 {
		c.A += 1
		dst.SetNRGBA(x, y, c)
	}

	return png.Encode(w, dst)
}

// An Encoder writes image to an output stream.
type Encoder struct {
	w io.Writer
}

// NewEncoder returns a new encoder that writes to w.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		w: w,
	}
}

// Encode writes the Transparency PNG of r to the stream.
func (enc *Encoder) Encode(r io.Reader) error {
	img, _, err := image.Decode(r)
	if err != nil {
		return err
	}

	var dst *image.NRGBA
	switch img := img.(type) {
	case *image.NRGBA:
		dst = img
	default:
		dst = image.NewNRGBA(img.Bounds())
		draw.Draw(dst, dst.Bounds(), img, image.ZP, draw.Src)
	}

	p := dst.Bounds().Max
	x, y := p.X-1, p.Y-1
	c := dst.NRGBAAt(x, y)
	if c.A == math.MaxUint8 {
		c.A -= 1
		dst.SetNRGBA(x, y, c)
	}

	return png.Encode(enc.w, dst)
}
