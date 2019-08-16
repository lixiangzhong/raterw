package raterw

import (
	"context"
	"io"

	"golang.org/x/time/rate"
)

func NewRateReader(r io.Reader, l rate.Limit, burst int) io.Reader {
	return &RateReader{
		R: r,
		L: rate.NewLimiter(l, burst),
	}
}

func AddRateReader(r io.Reader, l *rate.Limiter) io.Reader {
	return &RateReader{
		R: r,
		L: l,
	}
}

type RateReader struct {
	R io.Reader
	L *rate.Limiter
}

func (r *RateReader) Read(p []byte) (n int, err error) {
	lenp := len(p)
	ctx := context.Background()
	burst := r.L.Burst()
	b := make([]byte, burst)
	for {
		size := lenp - n
		if size < burst {
			b = b[:size]
		} else {
			size = burst
		}
		err = r.L.WaitN(ctx, size)
		if err != nil {
			return
		}
		var num int
		num, err = r.R.Read(b)
		n += copy(p[n:], b[:num])
		if n == lenp {
			return
		}
		if err != nil {
			return
		}
	}
}
