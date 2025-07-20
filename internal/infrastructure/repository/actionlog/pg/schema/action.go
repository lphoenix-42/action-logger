package schema

import (
	"encoding/json"
	"time"
)

type ActionSchema struct {
	ID         int64           `db:"id"`
	UserID     int64           `db:"user_id"`
	ActionType int16           `db:"action_type"`
	Timestamp  time.Time       `db:"timestamp"`
	Details    json.RawMessage `db:"details"`
}
