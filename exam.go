package appelli

import (
	"fmt"
	"time"
)

type Exam struct {
	Time   time.Time
	Type   string
	Course *Course
}

func (p *Exam) String() string {
	return fmt.Sprintf("%s %s %s %s", p.Time, p.Type, p.Course.Code, p.Course.Title)
}

// Exams implements sort.Interface for []Exam based on the Updated field.
type Exams []Exam

// ToHashMap creates a map from a slice of Exam.
//
// The key is the string representation of the Exam.
func (e Exams) ToHashMap() map[string]Exam {
	hashMap := make(map[string]Exam)
	for _, exam := range e {
		hashMap[exam.String()] = exam
	}
	return hashMap
}

func (e Exams) Len() int { return len(e) }
func (e Exams) Less(i, j int) bool {
	return e[i].Time.Before(e[j].Time)
}
func (e Exams) Swap(i, j int) { e[i], e[j] = e[j], e[i] }
