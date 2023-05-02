package handler

import (
	"net/http"
	"testing"
	"strings"
)
func TestHmacscha512Handler(t *testing.T) {

	cases := []struct{
		test string
		want int
	}{
		{
			`{
				"text": "Hello, world!",
				"key": "hello"
			}`,
			http.StatusOK,
		},
		{
			`{
				"text": "Sweet home",
				"key": 414241
			}`,
			http.StatusBadRequest,
		},
		{
			`{
				"key": "somekey",
				"text": "Lorem"
			}`,
			http.StatusOK,
		},
		{
			`{
				"key": "Secret"
			}`,
			http.StatusBadRequest,
		},
		{
			`{
				"text": "August"
			}`,
			http.StatusBadRequest,
		},
	}

	client := new(http.Client)
	for _, c := range cases {
		got, err := client.Post("http://localhost:8080/sign/hmacsha512", "application/json", strings.NewReader(c.test))
		if err != nil {
			t.Fatal(err)
		}
		if got.StatusCode != c.want {
			t.Fatalf("Non-expected status code %d, want %d `/sign/hmacsha512`:\nWith body:%s", 
					got.StatusCode, c.want, c.test)
		}
	}
}