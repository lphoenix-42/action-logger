package model

import "time"

const (
	ActionTypeUnspecified = iota
	ActionTypeBuy
	ActionTypeRefund
)

type ActionType int16

type Action struct {
	ID   int64
	Info *ActionInfo
}

type ActionInfo struct {
	UserID     int64
	ActionType ActionType
	Timestamp  time.Time
	Details    map[string]any
}
