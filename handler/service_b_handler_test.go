package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestGetWeatherHandler(t *testing.T) {
	type args struct {
		path   string
		status int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Valid CEP returns 200",
			args: args{
				path:   "/weather/06233903",
				status: 200,
			},
		},
		{
			name: "Invalid CEP returns 422",
			args: args{
				path:   "/weather/invalid",
				status: 422,
			},
		},
		{
			name: "Empty CEP returns 404",
			args: args{
				path:   "/weather/12345678",
				status: 404,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := mux.NewRouter()
			router.HandleFunc("/weather/{cep}", GetWeatherHandler).Methods("GET")

			req, _ := http.NewRequest(http.MethodGet, tt.args.path, nil)
			req.RemoteAddr = "0.0.0.1:8000"
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			assert.Equal(t, tt.args.status, rr.Code)
		})
	}
}
