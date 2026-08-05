package middlewares

import (
	"net/http"
	"slices"
	"strings"
)

type HppOptions struct {
	CheckQueryParams        bool
	CheckBodyParams         bool
	CheckBodyForContentType string
	Whitelist               []string
}

func filterQueryParams(r *http.Request, whitelist []string) {
	queryParams := r.URL.Query()

	for k, v := range queryParams {
		if isWhitelisted(k, whitelist) {
			continue
		}
		if len(v) > 1 {
			queryParams.Set(k, v[len(v)-1])
		}
	}

	r.URL.RawQuery = queryParams.Encode()
}

func filterBodyParams(w http.ResponseWriter, r *http.Request, whitelist []string) {
	// r.ParseForm only contains body
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Cannot parse form", http.StatusUnprocessableEntity)
		return
	}

	for k, v := range r.PostForm {
		if isWhitelisted(k, whitelist) {
			continue
		}
		if len(v) > 1 {
			r.Form.Set(k, v[len(v)-1])
		}
	}
}

func isAcceptedContentType(r *http.Request, contentType string) bool {
	return strings.Contains(r.Header.Get("Content-Type"), contentType)
}

func isWhitelisted(param string, whitelist []string) bool {
	if len(whitelist) == 0 {
		return true
	}
	return slices.Contains(whitelist, param)
}

func Hpp(opts HppOptions) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if opts.CheckQueryParams {
				filterQueryParams(r, opts.Whitelist)
			}
			if opts.CheckBodyParams &&
				(r.Method == http.MethodPost ||
					r.Method == http.MethodPut ||
					r.Method == http.MethodPatch) &&
				isAcceptedContentType(r, opts.CheckBodyForContentType) {
				filterBodyParams(w, r, opts.Whitelist)
			}

			next.ServeHTTP(w, r)
		})
	}
}
