package utils

import (
    "merchant-bank-go/models"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"

    "golang.org/x/crypto/bcrypt"
)

func main() {
    data, err := ioutil.ReadFile("data/customers.json")
    if err != nil {
        log.Fatalf("Failed to read customers.json: %v", err)
    }

    var customers []models.Customer
    err = json.Unmarshal(data, &customers)
    if err != nil {
        log.Fatalf("Failed to unmarshal customers.json: %v", err)
    }

    for i, customer := range customers {
        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(customer.Password), bcrypt.DefaultCost)
        if err != nil {
            log.Fatalf("Failed to hash password for user %s: %v", customer.Username, err)
        }
        customers[i].Password = string(hashedPassword)
        fmt.Printf("Hashed password for user %s\n", customer.Username)
    }

    updatedData, err := json.MarshalIndent(customers, "", "  ")
    if err != nil {
        log.Fatalf("Failed to marshal updated customers: %v", err)
    }

    err = ioutil.WriteFile("data/customers.json", updatedData, 0644)
    if err != nil {
        log.Fatalf("Failed to write updated customers.json: %v", err)
    }

    fmt.Println("Passwords hashed successfully.")
}
