package transfer

import (
	"mongo-transfer/database"
	"time"
)

type Transfer struct {
	Source *database.Server
	Target *database.Server

	StartedAt  time.Time
	FinishedAt time.Time
}
