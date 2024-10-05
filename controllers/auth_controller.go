package controllers

import (
    "merchant-bank-go/models"
    "merchant-bank-go/repository"
    "merchant-bank-go/utils"
    "encoding/json"
    "errors"
    "net/http"
    "strings"
    "time"
    "fmt"

    "golang.org/x/crypto/bcrypt"
)

type AuthController struct {
    CustomerRepo *repository.CustomerRepository
    HistoryRepo  *repository.HistoryRepository
}

var idCounter int = 1

func generateID() string {
    id := fmt.Sprintf("%d", idCounter)
    idCounter++
    return id
}

func (ac *AuthController) Login(w http.ResponseWriter, r *http.Request) {
    var creds struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }
    err := json.NewDecoder(r.Body).Decode(&creds)
    if err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    customer, err := ac.CustomerRepo.FindByUsername(creds.Username)
    if err != nil || customer == nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    err = bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(creds.Password))
    if err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    token, err := utils.GenerateJWT(customer.ID)
    if err != nil {
        http.Error(w, "Could not generate token", http.StatusInternalServerError)
        return
    }

    history := &models.History{
        ID:         generateID(),
        CustomerID: customer.ID,
        Action:     "login",
        Timestamp:  time.Now(),
    }
    err = ac.HistoryRepo.Add(history)
    if err != nil {
        utils.ErrorLogger.Println("Failed to log history:", err)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (ac *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
    customerID, err := extractCustomerID(r)
    if err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    history := &models.History{
        ID:         generateID(),
        CustomerID: customerID,
        Action:     "logout",
        Timestamp:  time.Now(),
    }
    err = ac.HistoryRepo.Add(history)
    if err != nil {
        utils.ErrorLogger.Println("Failed to log history:", err)
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Logged out successfully"))
}

func extractCustomerID(r *http.Request) (string, error) {
    authHeader := r.Header.Get("Authorization")
    if authHeader == "" {
        return "", errors.New("no token provided")
    }

    parts := strings.Split(authHeader, " ")
    if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
        return "", errors.New("invalid token format")
    }

    tokenStr := parts[1]
    customerID, err := utils.ValidateJWT(tokenStr)
    if err != nil {
        return "", err
    }
    return customerID, nil
}
