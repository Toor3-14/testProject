package handler

import (
	"net/http"
	"fmt"

	"github.com/Toor3-14/testProject/pkg/helper"
	"github.com/Toor3-14/testProject/internal/app"
)

func Users(w http.ResponseWriter, r *http.Request) {
	badRequest := func() {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	fields, err := helper.DecodeJson(r.Body, "name", "age")
	if err != nil { badRequest(); return}

	name, ok := fields["name"].(string)
	if !ok { badRequest(); return}

	age, ok := fields["age"].(float64)
	if !ok { badRequest(); return}

	id, err := app.UserRepo().Insert(name, int(age))
	if err != nil {
		badRequest()
		fmt.Fprint(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintf(w, "{\n\"id\":%d\n}", id)
}