package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGETGetHandler(t *testing.T) {
	t.Run("testing", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()

		GetHandler(response, request)
		fmt.Print(response)
		got := response.Body.String()
		want := "request- 1"

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}

	})

}

// func TestPOSTAddHandler(t *testing.T) {
// 	t.Run("Add a TODO item", func(t *testing.T) {
// 		request, _ := http.NewRequest(http.MethodGet, "/add", nil)
// 		response := httptest.NewRecorder()

// 		AddHandler(response, request)

// 	})

// }
