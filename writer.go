package raterw

import (
	"context"
	"io"

	"golang.org/x/time/rate"
)

func NewRateWriter(w io.Writer, l rate.Limit, burst int) io.Writer {
	return &RateWriter{
		w: w,
		l: rate.NewLimiter(l, burst),
	}
}

func AddRateWriter(w io.Writer, l *rate.Limiter) io.Writer {
	return &RateWriter{
		w: w,
		l: l,
	}
}

type RateWriter struct {
	w io.Writer
	l *rate.Limiter
}

func (r *RateWriter) Write(p []byte) (n int, err error) {
	lenp := len(p)
	ctx := context.Background()
	burst := r.l.Burst()
	for {
		size := lenp - n
		if size > burst {
			size = burst
		}
		err = r.l.WaitN(ctx, size)
		if err != nil {
			return
		}
		num, err := r.w.Write(p[n : n+size])
		n += num
		if n == lenp {
			return n, nil
		}
		if err != nil {
			return n, err
		}
	}
}
