package utils

import (
	"net/http"
	"time"

	"chino/pkg/log"

	"github.com/go-chi/render"
)

type Envelope struct {
	Success   bool        `json:"success"`
	RequestID string      `json:"request_id"`
	Time      string      `json:"time"`
	Response  interface{} `json:"response"`
	Error     string      `json:"error,omitempty"`
	// Message        string      `json:"message,omitempty"`
	HTTPStatusCode int `json:"-"`
}

func NewEnvelope() *Envelope {
	return &Envelope{}
}

func (e *Envelope) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (e *Envelope) SetResponse(resp interface{}) *Envelope {
	e.Success = true
	e.Time = time.Now().Format(time.RFC822)
	e.Response = resp
	return e
}

func (e *Envelope) SetError(resp error) *Envelope {
	e.Success = false
	e.Time = time.Now().Format(time.RFC822)
	e.Response = resp.Error()
	return e
}

func Render(w http.ResponseWriter, r *http.Request, e *Envelope) {
	if err := render.Render(w, r, e); err != nil {
		log.Error(r.Context(), err)
	}
}
