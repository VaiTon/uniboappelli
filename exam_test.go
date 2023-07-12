package appelli

import (
	"testing"
	"time"
)

func TestDiff(t *testing.T) {

	course1 := Course{Code: "01", Title: "Course 1", Teacher: "Teacher 1"}
	exam1 := Exam{
		Time:   time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC),
		Type:   "Scritto",
		Course: &course1,
	}

	exam2 := Exam{
		Time:   time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC),
		Type:   "Scritto",
		Course: &course1,
	}

	diff := Diff(Exams{exam1}, Exams{exam2})
	if len(diff) != 0 {
		t.Error("Diff should return nil", diff)
	}
}
