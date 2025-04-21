package server

import (
	"net/http"
	"service/pkg"
)

func ResponseCheckerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := pkg.GetLogger()

		logger.Infof("w: %v", w.Header())
		logger.Infof("r: %v", r.Context())

		next.ServeHTTP(w, r)
	})
}
