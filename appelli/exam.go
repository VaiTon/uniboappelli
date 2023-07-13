package appelli

import (
	"time"
)

type Exam struct {
	Time   time.Time
	Type   string
	Course Course
}

// Exams implements sort.Interface for []Exam based on the Updated field.
type Exams []Exam

func (e Exams) Len() int { return len(e) }
func (e Exams) Less(i, j int) bool {
	return e[i].Time.Before(e[j].Time)
}
func (e Exams) Swap(i, j int) { e[i], e[j] = e[j], e[i] }
