package raterw

import (
	"context"
	"io"

	"golang.org/x/time/rate"
)

func NewRateReadCloser(r io.ReadCloser, l rate.Limit, burst int) io.ReadCloser {
	return &RateReadCloser{
		R: r,
		l: rate.NewLimiter(l, burst),
	}
}

func AddRateReadCloser(r io.ReadCloser, l *rate.Limiter) io.ReadCloser {
	return &RateReadCloser{
		R: r,
		l: l,
	}
}

type RateReadCloser struct {
	R io.ReadCloser
	l *rate.Limiter
}

func (r *RateReadCloser) Close() error {
	return r.R.Close()
}

func (r *RateReadCloser) Read(p []byte) (n int, err error) {
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
