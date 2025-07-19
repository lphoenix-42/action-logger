package model

import "time"

type ActionSearch struct {
	UserID         *int64
	ActionTypes    []ActionType
	From           *time.Time
	To             *time.Time
	DetailsFilters []JsonFilter
}

type JsonFilter struct {
	Path  []string
	Value string
}
