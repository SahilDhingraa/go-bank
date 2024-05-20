package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func WriteJSON(w http.ResponseWriter, status int, val any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json ")
	return json.NewEncoder(w).Encode(val)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type APIError struct {
	Error string
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
		}
	}
}

type APIServer struct {
	listenAddress string
}

func NewAPIServer(listenAddr string) *APIServer {
	sim := &APIServer{
		listenAddress: listenAddr,
	}
	return sim
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/account", makeHTTPHandleFunc(s.HandleAccount))

	router.HandleFunc("/account/{id}", makeHTTPHandleFunc(s.HandleGetAccount))

	log.Println("JSON API SERVER RUNNING ON PORT ", s.listenAddress)

	http.ListenAndServe(s.listenAddress, router)
	http.NewServeMux()
}

func (s *APIServer) HandleAccount(w http.ResponseWriter, r *http.Request) error {

	switch r.Method {
	case "GET":
		return s.HandleGetAccount(w, r)
	case "POST":
		return s.HandleCreateAccount(w, r)
	case "DELETE":
		return s.HandleDeleteAccount(w, r)
	default:
		return fmt.Errorf("method not allowed %s", r.Method)
	}
}

func (s *APIServer) HandleGetAccount(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)["id"]
	
	// account := NewAccount("Sahil", "Dhingra")
	return WriteJSON(w, http.StatusOK, vars)
}

func (s *APIServer) HandleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) HandleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) HandleTransfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}
