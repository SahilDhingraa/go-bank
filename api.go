package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type apiFunc func(http.ResponseWriter, *http.Request) error
type APIServer struct {
	listenAddress string
	store         Storage
}
type APIError struct {
	Error string
}

func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddress: listenAddr,
		store:         store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/account", makeHTTPHandleFunc(s.HandleAccount))

	router.HandleFunc("/account/{id}", makeHTTPHandleFunc(s.HandleGetAccountByID))

	log.Println("JSON API SERVER RUNNING ON PORT ", s.listenAddress)

	http.ListenAndServe(s.listenAddress, router)
	// http.NewServeMux()
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
	accounts, err := s.store.GetAccounts()
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, accounts)
}
func (s *APIServer) HandleGetAccountByID(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)["id"]

	// account := NewAccount("Sahil", "Dhingra")
	return WriteJSON(w, http.StatusOK, vars)
}

func (s *APIServer) HandleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	createAccountReq := new(CreateAccountRequest)
	// createAccountRequest := CreateAccountRequest{}
	if err := json.NewDecoder(r.Body).Decode(createAccountReq); err != nil {
		// if err := json.NewDecoder(r.Body).Decode(&createAccountRequest); err != nil {
		return err
	}

	account := NewAccount(createAccountReq.FirstName, createAccountReq.LastName)
	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) HandleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) HandleTransfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func WriteJSON(w http.ResponseWriter, status int, val any) error {
	w.Header().Set("Content-Type", "application/json ")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(val)
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
		}
	}
}
