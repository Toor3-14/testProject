package handler

import (
	"net/http"
	"encoding/json"
	"context"
	"fmt"
	"time"

	"github.com/Toor3-14/testProject/pkg/helper"
	"github.com/Toor3-14/testProject/internal/app"
)

type IncrJson struct {
	Key string `json:"key"`
	Value int `json:"value"`
}

func Incr(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var incr IncrJson
	err := decoder.Decode(&incr)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if err := helper.ValidateFields(incr.Key, incr.Value); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	incr.Value += 1

	// TO CHANGE: key value is stored for a minute
	err = app.Redis().Set(context.Background(), incr.Key, incr.Value, time.Minute).Err()
	if err != nil {
		app.ErrLog().Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintf(w, "{\n\"value\":%d\n}", incr.Value)
}