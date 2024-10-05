package controllers

import (
    "merchant-bank-go/models"
    "merchant-bank-go/repository"
	"merchant-bank-go/utils"
    "encoding/json"
    "net/http"
	// "errors"
    "time"
)

type PaymentController struct {
    CustomerRepo *repository.CustomerRepository
    HistoryRepo  *repository.HistoryRepository
}

func (pc *PaymentController) Pay(w http.ResponseWriter, r *http.Request) {
    var payment struct {
        ToUsername string  `json:"to_username"`
        Amount     float64 `json:"amount"`
    }
    err := json.NewDecoder(r.Body).Decode(&payment)
    if err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    if payment.Amount <= 0 {
        http.Error(w, "Amount must be greater than zero", http.StatusBadRequest)
        return
    }

    fromCustomerID, err := extractCustomerID(r)
    if err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    fromCustomer, err := pc.CustomerRepo.FindByID(fromCustomerID)
    if err != nil || fromCustomer == nil {
        http.Error(w, "Customer not found", http.StatusNotFound)
        return
    }

    toCustomer, err := pc.CustomerRepo.FindByUsername(payment.ToUsername)
    if err != nil || toCustomer == nil {
        http.Error(w, "Recipient not found", http.StatusNotFound)
        return
    }

    if fromCustomer.Balance < payment.Amount {
        http.Error(w, "Insufficient balance", http.StatusBadRequest)
        return
    }

    fromCustomer.Balance -= payment.Amount
    toCustomer.Balance += payment.Amount

    err = pc.CustomerRepo.Update(fromCustomer)
    if err != nil {
        http.Error(w, "Failed to update sender balance", http.StatusInternalServerError)
        return
    }

    err = pc.CustomerRepo.Update(toCustomer)
    if err != nil {
        http.Error(w, "Failed to update recipient balance", http.StatusInternalServerError)
        return
    }

    history := &models.History{
        ID:         generateID(),
        CustomerID: fromCustomer.ID,
        Action:     "payment",
        Amount:     payment.Amount,
        ToCustomer: toCustomer.ID,
        Timestamp:  time.Now(),
    }
    err = pc.HistoryRepo.Add(history)
    if err != nil {
        utils.ErrorLogger.Println("Failed to log history:", err)
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Payment successful"))
}
