package raterw

import (
	"io"

	"golang.org/x/time/rate"
)

type RateReadWriter struct {
	Limit *rate.Limiter
	io.Reader
	io.Writer
}

func NewRateReadWriter(rw io.ReadWriter, l rate.Limit, burst int) io.ReadWriter {
	limit := rate.NewLimiter(l, burst)
	r := AddRateReader(rw, limit)
	w := AddRateWriter(rw, limit)
	return &RateReadWriter{Limit: limit, Reader: r, Writer: w}
}

func AddRateReadWriter(rw io.ReadWriter, limit *rate.Limiter) io.ReadWriter {
	r := AddRateReader(rw, limit)
	w := AddRateWriter(rw, limit)
	return &RateReadWriter{Limit: limit, Reader: r, Writer: w}
}
