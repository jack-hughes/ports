package handlers

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jack-hughes/ports/internal/client/service"
	"log"
	"net/http"
)

func Get(ctx context.Context, svc service.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		par := mux.Vars(r)
		port, err := svc.Get(ctx, par["port_id"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			// TODO: Log errors
			return
		}
		if port.ID == "" {
			w.WriteHeader(http.StatusNotFound)
			// TODO: Log errors
			return
		}

		j, err := json.Marshal(port)
		_, err = w.Write(j)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			// TODO: Log errors
			return
		}
	}
}

func List(ctx context.Context, svc service.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		pSlice, err := svc.List(ctx)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatal(err)
			return
		}

		j, err := json.Marshal(pSlice)
		_, err = w.Write(j)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatal(err)
			return
		}
	}
}
