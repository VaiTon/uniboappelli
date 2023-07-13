package appelli

import (
	"time"
)

type Exam struct {
	Time   time.Time
	Type   string
	Course Course
}
