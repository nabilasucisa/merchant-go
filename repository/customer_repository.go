package repository

import (
    "encoding/json"
    "io/ioutil"
    "sync"

    "merchant-bank-go/models"
)

type CustomerRepository struct {
    filePath  string
    mutex     sync.Mutex
    Customers []models.Customer
}

func NewCustomerRepository(filePath string) (*CustomerRepository, error) {
    repo := &CustomerRepository{filePath: filePath}
    data, err := ioutil.ReadFile(filePath)
    if err != nil {
        return nil, err
    }
    err = json.Unmarshal(data, &repo.Customers)
    if err != nil {
        return nil, err
    }
    return repo, nil
}

func (r *CustomerRepository) Save() error {
    r.mutex.Lock()
    defer r.mutex.Unlock()
    data, err := json.MarshalIndent(r.Customers, "", "  ")
    if err != nil {
        return err
    }
    return ioutil.WriteFile(r.filePath, data, 0644)
}

func (r *CustomerRepository) FindByUsername(username string) (*models.Customer, error) {
    for _, customer := range r.Customers {
        if customer.Username == username {
            return &customer, nil
        }
    }
    return nil, nil
}

func (r *CustomerRepository) FindByID(id string) (*models.Customer, error) {
    for _, customer := range r.Customers {
        if customer.ID == id {
            return &customer, nil
        }
    }
    return nil, nil
}

func (r *CustomerRepository) Update(updatedCustomer *models.Customer) error {
    for i, customer := range r.Customers {
        if customer.ID == updatedCustomer.ID {
            r.Customers[i] = *updatedCustomer
            return r.Save()
        }
    }
    return nil
}
