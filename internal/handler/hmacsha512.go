package handler

import (
	"net/http"
	"encoding/json"
	"crypto/hmac"
	"crypto/sha512"
	"fmt"

	"github.com/Toor3-14/testProject/pkg/helper"
	"github.com/Toor3-14/testProject/internal/app"
)

type Hmacsha512Json struct {
	Text string `json:"text"`
	Key string `json:"key"`
}

func Hmacsha512(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var obj Hmacsha512Json
	err := decoder.Decode(&obj)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if err := helper.ValidateFields(obj.Text, obj.Key); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	mac := hmac.New(sha512.New, []byte(obj.Key))
	if _, err := mac.Write([]byte(obj.Text)); err != nil {
		app.ErrLog().Println(err)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w ,"%x", mac.Sum(nil))
}