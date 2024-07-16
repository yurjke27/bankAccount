package main

import (
	"atm/account"
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
)

var (
	accounts = make(map[int]*account.Account)
	nextID   = 1
	mu       sync.Mutex
)

func createAccount(w http.ResponseWriter, _ *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	acc := &account.Account{ID: nextID}
	accounts[nextID] = acc
	nextID++
	w.WriteHeader(http.StatusCreated)
	err := json.NewEncoder(w).Encode(acc)
	if err != nil {
		return
	}
}

func deposit(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing account ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid account ID", http.StatusBadRequest)
		return
	}

	acc, exists := accounts[id]
	if !exists {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}

	var requestData struct {
		Amount float64 `json:"amount"`
	}
	err = json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = acc.Deposit(requestData.Amount)
	if err != nil {
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(map[string]float64{"balance": acc.GetBalance()})
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func withdraw(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing account ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid account ID", http.StatusBadRequest)
		return
	}

	acc, exists := accounts[id]
	if !exists {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}

	var requestData struct {
		Amount float64 `json:"amount"`
	}
	err = json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = acc.Withdraw(requestData.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(map[string]float64{"balance": acc.GetBalance()})
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func getBalance(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing account ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid account ID", http.StatusBadRequest)
		return
	}

	acc, exists := accounts[id]
	if !exists {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}

	balance := acc.GetBalance()
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(map[string]float64{"balance": balance})
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

}

func main() {
	http.HandleFunc("/accounts", createAccount)
	http.HandleFunc("/accounts/deposit", deposit)
	http.HandleFunc("/accounts/withdraw", withdraw)
	http.HandleFunc("/accounts/getbalance", getBalance)
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		return
	}
}
