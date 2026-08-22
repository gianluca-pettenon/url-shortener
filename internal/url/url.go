package url

import "time"

type URL struct {
	ID          uint64
	OriginalURL string
	CreatedAt   time.Time
}
