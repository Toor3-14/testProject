package handler

import (
	"net/http"
	"testing"
	"strings"
)
func TestUserHandler(t *testing.T) {

	cases := []struct{
		test string
		want int
	}{
		{
			`{
				"name": "Alex",
				"age": 25
			}`,
			http.StatusOK,
		},
		{
			`{
				"name": "Alex",
				"age": 47
			}`,
			http.StatusBadRequest,
		},
		{
			`{
				"name": "Jhon",
				"hello": "test"
			}`,
			http.StatusBadRequest,
		},
		{
			`{
				"name": "Philip"
			}`,
			http.StatusBadRequest,
		},
		{
			`{
				"name": "August",
				"age": "51"
			}`,
			http.StatusBadRequest,
		},
	}

	client := new(http.Client)
	for _, c := range cases {
		got, err := client.Post("http://localhost:8080/postgres/users", "application/json", strings.NewReader(c.test))
		if err != nil {
			t.Fatal(err)
		}
		if got.StatusCode != c.want {
			t.Fatalf("Non-expected status code %d, want %d `/postgres/users`:\nWith body:%s", 
					got.StatusCode, c.want, c.test)
		}
	}
}