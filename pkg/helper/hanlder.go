package helper

import (
	"net/http"
	"errors"
	"io"
	"encoding/json"
)

func JustMethod(handler http.HandlerFunc,necessaryMethod string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != necessaryMethod {
			w.Header().Set("Allow", http.MethodPost)
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		handler(w,r)
	}

}
func DecodeJson(body io.ReadCloser, correctFields ...string) (map[string]interface{}, error) {
	fields := make(map[string]interface{})

	decoder := json.NewDecoder(body)
	err := decoder.Decode(&fields)
	if err != nil {
		return nil, err
	}
	for field := range fields {
		fail := 0
		for _,correct := range correctFields {
			if field != correct {
				fail++
			}
		}
		if fail == len(correctFields) || fail < len(correctFields)-1 {
			return nil, errors.New("Incorrect json field " + field)
		}
	}
	return fields, nil
}
func ValidateFields(objParams ...any) error {
	for _, v := range objParams {
		if v == nil || v == "" || v == 0 || v == 0.0 {
			return errors.New("Incorrect json fields")
		}
	}
	return nil
}
