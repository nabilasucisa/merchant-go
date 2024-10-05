package main

import (
    "log"
    "net/http"

    "merchant-bank-go/controllers"
    "merchant-bank-go/repository"
    "merchant-bank-go/utils"

    "github.com/gorilla/mux"
)

func main() {
    utils.InitLogger()
    utils.InfoLogger.Println("Logger initialized")

    customerRepo, err := repository.NewCustomerRepository("data/customers.json")
    if err != nil {
        utils.ErrorLogger.Fatalln("Failed to load customers:", err)
    }
    utils.InfoLogger.Println("Customer repository initialized")

    historyRepo, err := repository.NewHistoryRepository("data/history.json")
    if err != nil {
        utils.ErrorLogger.Fatalln("Failed to load history:", err)
    }
    utils.InfoLogger.Println("History repository initialized")

    authController := &controllers.AuthController{
        CustomerRepo: customerRepo,
        HistoryRepo:  historyRepo,
    }

    paymentController := &controllers.PaymentController{
        CustomerRepo: customerRepo,
        HistoryRepo:  historyRepo,
    }

    r := mux.NewRouter()

    r.HandleFunc("/login", authController.Login).Methods("POST")
    r.HandleFunc("/logout", authController.Logout).Methods("POST")
    r.HandleFunc("/payment", paymentController.Pay).Methods("POST")

    utils.InfoLogger.Println("Server starting on :8080")
    log.Println("Server starting on :8080")
    if err := http.ListenAndServe(":8080", r); err != nil {
        utils.ErrorLogger.Fatalln("Failed to start server:", err)
    }
}
