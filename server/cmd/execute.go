package main

import "net/http"

func (app *application) createExecTask(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Stdin    string `json:"stdin"`
		Solution string `json:"solution"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

}
