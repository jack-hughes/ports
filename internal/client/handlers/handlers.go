package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jack-hughes/ports/internal/client/service"
	"go.uber.org/zap"
	"net/http"
)

// Get returns a JSON body for a singular port
func Get(ctx context.Context, svc service.Service, log *zap.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log = log.With(zap.String("component", "handler"))
		w.Header().Set("Content-Type", "application/json")
		par := mux.Vars(r)
		log.Debug(fmt.Sprintf("making get request made for port ID: %v", par["port_id"]))
		port, err := svc.Get(ctx, par["port_id"])
		if err != nil {
			log.Error("failed to get port: %v", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if port.ID == "" {
			log.Error(fmt.Sprintf("port does not exist: %v", par["port_id"]))
			w.WriteHeader(http.StatusNotFound)
			return
		}

		j, err := json.Marshal(port)
		if err != nil {
			log.Error("failed to marshal json %v", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = w.Write(j)
		if err != nil {
			log.Error(fmt.Sprintf("unexpected error on get for port: %v, error: %v", par["port_id"], zap.Error(err)))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

// List returns a JSON body for all ports
func List(ctx context.Context, svc service.Service, log *zap.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log = log.With(zap.String("component", "handler"))
		log.Debug("making list request")
		w.Header().Set("Content-Type", "application/json")
		pSlice, err := svc.List(ctx)
		if err != nil {
			log.Error("failed to list ports: %v", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		j, err := json.Marshal(pSlice)
		if err != nil {
			log.Error("failed to marshal json %v", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = w.Write(j)
		if err != nil {
			log.Error("failed to return result %v", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
