package httperr

import (
	"net/http"

	logger "github.com/lGreenLightl/fitness-club-service/internal/app/logs/logrus"

	"github.com/go-chi/render"
)

type ErrResponse struct {
	Slug       string `json:"slug"`
	httpStatus int
}

func (e ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(e.httpStatus)

	return nil
}

func httpResponseWithErr(r *http.Request, w http.ResponseWriter, slug string, err error, logMessage string, status int) {
	logger.LogEntry(r).WithError(err).WithField("error-slug", slug).Warn(logMessage)

	if err := render.Render(w, r, ErrResponse{Slug: slug, httpStatus: status}); err != nil {
		panic(err)
	}
}

func BadRequest(r *http.Request, w http.ResponseWriter, slug string, err error) {
	httpResponseWithErr(r, w, slug, err, "Bad request", http.StatusBadRequest)
}

func InternalErr(r *http.Request, w http.ResponseWriter, slug string, err error) {
	httpResponseWithErr(r, w, slug, err, "Internal server error", http.StatusInternalServerError)
}

func Unauthorized(r *http.Request, w http.ResponseWriter, slug string, err error) {
	httpResponseWithErr(r, w, slug, err, "Unauthorized", http.StatusUnauthorized)
}
