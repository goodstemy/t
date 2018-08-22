package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"encoding/json"
	"fmt"
	"strconv"
)

type Store struct {
	Accounts map[int]*Account
}

type Account struct {
	Name string `json:"name,omitempty"`
	Surname string `json:"surname,omitempty"`
	Balance float32 `json:"balance,omitempty"`
	ID int `json:"id,omitempty"`
}

type Transaction struct {
	From int `json:"from,omitempty"`
	To int `json:"to,omitempty"`
	Amount float32 `json:"amount,omitempty"`
}

func main() {
	router := mux.NewRouter()
	store := Store{
		Accounts: map[int]*Account{},
	}

	router.HandleFunc("/account", func(w http.ResponseWriter, r *http.Request) {
		store.CreateAccount(w, r)
	}).Methods("POST")
	router.HandleFunc("/sendMoney", func(w http.ResponseWriter, r *http.Request) {
		store.SendMoney(w, r)
	}).Methods("POST")
	router.HandleFunc("/account/{id}", func(w http.ResponseWriter, r *http.Request) {
		store.GetAccount(w, r)
	}).Methods("GET")

	log.Println("Listening on :8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func (s *Store) CreateAccount(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var account Account

	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user data")
		return
	}

	if account.Name == "" || account.Surname == "" || account.Balance < 0 || account.ID <= 0 {
		respondWithError(w, http.StatusBadRequest, "Invalid user data")
		return
	}

	if _, ok := s.Accounts[account.ID]; ok {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("User with ID %d already created", account.ID))
		return
	}

	s.Accounts[account.ID] = &account

	w.WriteHeader(http.StatusOK)
}

func (s *Store) GetAccount(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)

	if _, ok := vars["id"]; !ok {
		respondWithError(w, http.StatusBadRequest, "Invalid id")
		return
	}

	var id int
	var err error

	if id, err = strconv.Atoi(vars["id"]); err != nil {
		respondWithError(w, http.StatusBadRequest, "Id must be a string")
		return
	}

	var account *Account
	var ok bool

	if account, ok = s.Accounts[id]; !ok {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Can not find any user with id %d", id))
		return
	}

	resp, err := json.Marshal(&account)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJson(w, http.StatusOK, resp)
}

func (s *Store) SendMoney(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var tx Transaction

	if err := json.NewDecoder(r.Body).Decode(&tx); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid transaction data")
		return
	}

	if _, ok := s.Accounts[tx.From]; !ok {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Can not find user with id %d", tx.From))
		return
	}

	if _, ok := s.Accounts[tx.To]; !ok {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Can not find user with id %d", tx.To))
		return
	}

	if s.Accounts[tx.From].Balance < tx.Amount {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Not enough of money on balance with id %d", tx.From))
		return
	}

	s.Accounts[tx.From].Balance -= tx.Amount
	s.Accounts[tx.To].Balance += tx.Amount

	w.WriteHeader(http.StatusOK)
}

func respondWithJson(w http.ResponseWriter, status int, json []byte) {
	w.WriteHeader(status)
	fmt.Fprintf(w, string(json))
}

func respondWithError(w http.ResponseWriter, status int, errMessage string) {
	w.WriteHeader(status)
	fmt.Fprintf(w, errMessage)
}