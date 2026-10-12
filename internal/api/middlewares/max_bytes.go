// https://destel.dev/blog/max-bytes-middleware-in-go
package middlewares

import (
	"context"
	"io"
	"net/http"
)

var limiterCtxKey = new(int)

type limitReader struct {
	io.Reader
	io.Closer
	limit int64
	n     int64 // bytes read so far
	err   error // sticky error
}

func (r *limitReader) Read(p []byte) (n int, err error) {
	if r.err != nil {
		return 0, r.err
	}
	if len(p) == 0 {
		return 0, nil
	}

	// Ask the body for at most one byte past the limit.
	if remaining := r.limit - r.n + 1; int64(len(p)) > remaining {
		p = p[:remaining]
	}

	n, err = r.Reader.Read(p)
	r.n += int64(n)

	// Within the limit: hand out everything and remember the error.
	if r.n <= r.limit {
		r.err = err
		return n, err
	}

	// Past the limit: hand out only the bytes within the limit and fail.
	n -= int(r.n - r.limit)
	r.err = &http.MaxBytesError{Limit: r.limit}
	return n, r.err
}

func MaxBytes(size int64) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			limiter, _ := ctx.Value(limiterCtxKey).(*limitReader)

			if limiter != nil {
				limiter.limit = size
				next.ServeHTTP(w, r)
				return
			}

			limiter = &limitReader{
				Reader: r.Body,
				Closer: r.Body,
				limit:  size,
			}

			ctx = context.WithValue(ctx, limiterCtxKey, limiter)
			r = r.WithContext(ctx)
			r.Body = limiter

			next.ServeHTTP(w, r)
		})
	}
}
