package raterw

import (
	"context"
	"io"

	"golang.org/x/time/rate"
)

func NewRateReader(r io.Reader, l rate.Limit, burst int) io.Reader {
	return &RateReader{
		R: r,
		l: rate.NewLimiter(l, burst),
	}
}

type RateReader struct {
	R io.Reader
	l *rate.Limiter
}

func (r *RateReader) Read(p []byte) (n int, err error) {
	lenp := len(p)
	ctx := context.Background()
	burst := r.l.Burst()
	b := make([]byte, burst)
	for {
		size := lenp - n
		if size < burst {
			b = b[:size]
		} else {
			size = burst
		}
		err = r.l.WaitN(ctx, size)
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
