package main

import "net/http"

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	env := envelope{"error": message}

	err := app.writeJSON(w, status, env, nil)
	if err != nil {
		app.logError(r, err)
	}
}

// func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request) {

// }

func (app *application) logError(r *http.Request, err error) {
	app.logger.Error().Err(err).Str("request_method", r.Method).Str("request_url", r.URL.String())
}
