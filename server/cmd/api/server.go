package main

import (
	"errors"
	"fmt"
	"net/http"
)

func (app *application) server() error {
	srv := &http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%d", app.config.port),
		Handler: app.routes(),
	}

	app.logger.Info().Str("addr", srv.Addr).Msg("Starting server")
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
