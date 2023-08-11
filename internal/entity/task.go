package entity

import (
)

type Task struct {
	Status string      `json:"status,omitempty"`
	Title    string    `json:"title" binding:"required,max=200"`
	ActiveAt string    `json:"activeAt" binding:"required"`
}