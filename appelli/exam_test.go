package appelli

import (
	"fmt"
	"testing"
	"time"
)

func TestDiff(t *testing.T) {

	course := Course{Code: "01", Title: "Course 1", Teacher: "Teacher 1"}
	exam1 := Exam{
		Time:   time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC),
		Type:   "Scritto",
		Course: Course{Code: "01", Title: "Course 1", Teacher: "Teacher 1"},
	}
	exam2 := Exam{
		Time:   time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC),
		Type:   "Scritto",
		Course: course,
	}

	diff := Diff(Exams{exam1}, Exams{exam2})
	if len(diff) != 0 {
		t.Error("Diff should return empty", diff)
	}

}

func TestDiff2(t *testing.T) {
	exam1 := Exam{
		Time:   time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC),
		Type:   "Scritto",
		Course: Course{Code: "01", Title: "Course 1", Teacher: "Teacher 1"},
	}
	exam3 := Exam{
		Time:   time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC),
		Type:   "Scritto",
		Course: Course{Code: "02", Title: "Course 2", Teacher: "Teacher 2"},
	}
	diff := Diff(Exams{exam1}, Exams{exam3})
	if len(diff) != 1 {
		t.Error("Diff should return 1 element", diff)
	}
}

func BenchmarkDiff(b *testing.B) {
	exams1, exams2 := createTestExams()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Diff(exams1, exams2)
	}
}

func createTestExams() ([]Exam, []Exam) {
	exams1 := make([]Exam, 0, 100)
	exams2 := make([]Exam, 0, 100)

	for i := 0; i < 100; i++ {
		exams1 = append(exams1, Exam{
			Time:   time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC),
			Type:   fmt.Sprintf("Scritto %d", i),
			Course: Course{Code: fmt.Sprintf("01%d", i), Title: fmt.Sprintf("Course %d", i), Teacher: fmt.Sprintf("Teacher %d", i)},
		})

		exams2 = append(exams2, Exam{
			Time:   time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC),
			Type:   fmt.Sprintf("Scritto %d", i),
			Course: Course{Code: fmt.Sprintf("01%d", i), Title: fmt.Sprintf("Course %d", i), Teacher: fmt.Sprintf("Teacher %d", i)},
		})

	}
	return exams1, exams2
}
