package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type JSONAdapter struct{}

func (a JSONAdapter) Render(w http.ResponseWriter, payload any, status int) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(payload)
}

func (a JSONAdapter) Error(w http.ResponseWriter, err error, status int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	_err := json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
	if _err != nil {
		panic(_err)
	}
}

type HTMLAdapter struct{}

func (a HTMLAdapter) Render(w http.ResponseWriter, payload any, status int) error {
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(status)

	body, ok := payload.([]byte)
	if !ok {
		return fmt.Errorf("body must be []byte but got %T", payload)
	}
	_, err := w.Write(body)
	return err
}

func (a HTMLAdapter) Error(w http.ResponseWriter, err error, status int) {
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(status)

	body := []byte(fmt.Sprintf(`
  <html>
    <body>
    <h1>Error: %s</h1>
    <p>Status: %d</p>
    </body>
  </html>
  `, err, status))

	if _, err := w.Write(body); err != nil {
		panic(err)
	}
}

type TextAdapter struct{}

func (a TextAdapter) Render(w http.ResponseWriter, payload any, status int) error {
	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(status)

	body, ok := payload.([]byte)
	if !ok {
		return fmt.Errorf("body must be []byte but got %T", payload)
	}
	_, err := w.Write(body)
	return err
}

func (a TextAdapter) Error(w http.ResponseWriter, err error, status int) {
	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(status)

	body := []byte(fmt.Sprintf(`
    Error: %s \n
    Status: %d
  `, err, status))

	if _, err := w.Write(body); err != nil {
		panic(err)
	}
}
