package bwgtest

import "time"

type Quotation struct {
	UpdateID   int       `json:"update_id" db:"update_id"`
	CodeFrom   string    `json:"code_from" db:"code_from" binding:"required"`
	CodeTo     string    `json:"code_to" db:"code_to" binding:"required"`
	Rate       float32   `json:"rate" db:"rate"`
	UpdateTime time.Time `json:"update_time" db:"update_time"`
	Success    bool      `json:"success" db:"success"`
}
