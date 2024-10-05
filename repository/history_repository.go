package repository

import (
    "encoding/json"
    "io/ioutil"
    "sync"

    "merchant-bank-go/models"
)

type HistoryRepository struct {
    filePath string
    mutex    sync.Mutex
    Histories []models.History
}

func NewHistoryRepository(filePath string) (*HistoryRepository, error) {
    repo := &HistoryRepository{filePath: filePath}
    data, err := ioutil.ReadFile(filePath)
    if err != nil {
        return nil, err
    }
    err = json.Unmarshal(data, &repo.Histories)
    if err != nil {
        return nil, err
    }
    return repo, nil
}

func (r *HistoryRepository) Add(history *models.History) error {
    r.mutex.Lock()
    defer r.mutex.Unlock()
    r.Histories = append(r.Histories, *history)
    data, err := json.MarshalIndent(r.Histories, "", "  ")
    if err != nil {
        return err
    }
    return ioutil.WriteFile(r.filePath, data, 0644)
}
