package raterw

import (
	"context"
	"io"

	"golang.org/x/time/rate"
)

func NewRateWriter(w io.Writer, l rate.Limit, burst int) io.Writer {
	return &RateWriter{
		W: w,
		L: rate.NewLimiter(l, burst),
	}
}

func AddRateWriter(w io.Writer, l *rate.Limiter) io.Writer {
	return &RateWriter{
		W: w,
		L: l,
	}
}

type RateWriter struct {
	W io.Writer
	L *rate.Limiter
}

func (r *RateWriter) Write(p []byte) (n int, err error) {
	lenp := len(p)
	ctx := context.Background()
	burst := r.L.Burst()
	for {
		size := lenp - n
		if size > burst {
			size = burst
		}
		err = r.L.WaitN(ctx, size)
		if err != nil {
			return
		}
		num, err := r.W.Write(p[n : n+size])
		n += num
		if n == lenp {
			return n, nil
		}
		if err != nil {
			return n, err
		}
	}
}
