

func (m *Middleware) Headers(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'")
		w.Header().Set("Content-Security-Policy", "require-sri-for style script")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Add("X-Frame-Options", "SAMEORIGIN")
		w.Header().Set("Strict-Transport-Security", "max-age=315360000; includeSubdomains; preload")
		w.Header().Set("Cache-Control", "no-cache, no-store")
		w.Header().Set("Public-Key-Pins", "pin-sha256=base64==; max-age=315360000")
		w.Header().Set("X-Powered-By", "")
		w.Header().Set("Server", "")
		next.ServeHTTP(w, r)
	})
}