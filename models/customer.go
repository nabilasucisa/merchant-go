package models

type Customer struct {
    ID       string  `json:"id"`
    Username string  `json:"username"`
    Password string  `json:"password"` 
    Balance  float64 `json:"balance"`
}
