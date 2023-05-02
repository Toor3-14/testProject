package handler

import (
	"net/http"
	"testing"
	"strings"
	"io"
)

func TestIncrHandler(t *testing.T) {

	cases := []struct{
		test string
		statusWant int
		messageWant string
	}{
		{
			`{
				"key": "Age",
				"value": 79
			}`,
			http.StatusOK,
			"{\n\"value\":80\n}",
		},
		{
			`{
				"value": 29,
				"key": "Number"
			}`,
			http.StatusOK,
			"{\n\"value\":30\n}",
		},
		{
			// wait - error, status 400 Bad Request
			`{
				"value": 26,
			}`,
			http.StatusBadRequest,
			"Bad Request",
		},
		{
			// test for incorrect data type of filed, wait error, status 400 Bad Request
			`{
				"key": "Number",
				"value": "90",
			}`,
			http.StatusBadRequest,
			"Bad Request",
		},
	}

	client := new(http.Client)
	for _, c := range cases {
		got, err := client.Post("http://localhost:8080/redis/incr", "application/json", strings.NewReader(c.test))
		if err != nil {
			t.Fatal(err)
		}
		if got.StatusCode != c.statusWant {
			body, err := io.ReadAll(got.Body)
			if err != nil {
				t.Fatal(err)
			}
			if string(body) != c.messageWant {
				t.Fatalf("Non-expected status code %d, want %d `/redis/incr`:\nHave body:%s\nWait body:%s", 
						got.StatusCode, c.statusWant, c.test, body)
			}
		}
	}

}