package main

import (
	"fmt"
	appelli "github.com/VaiTon/uniboappelli"
	"github.com/gorilla/feeds"
	"github.com/samber/lo"
	"strings"
	"time"
)

type TimedExam struct {
	Updated time.Time
	Exam    appelli.Exam
}

type TimedExams []TimedExam

func NewTimedExams(exams appelli.Exams, time time.Time) TimedExams {
	return lo.Map(exams, func(exam appelli.Exam, _ int) TimedExam {
		return NewTimedExam(exam, time)
	})
}

func NewTimedExam(exam appelli.Exam, time time.Time) TimedExam {
	return TimedExam{Updated: time, Exam: exam}
}

func (e TimedExams) ToExams() appelli.Exams {
	return lo.Map(e, func(p TimedExam, _ int) appelli.Exam {
		return p.Exam
	})
}

func (e TimedExam) ToFeedItem(url string) *feeds.Item {
	b := strings.Builder{}

	b.WriteString(fmt.Sprintf("Data: %s\n", e.Updated))
	b.WriteString(fmt.Sprintf("Type: %s\n", e.Exam.Type))
	b.WriteString(fmt.Sprintf("Teacher: %s\n", e.Exam.Course.Teacher))
	b.WriteString(fmt.Sprintf("Course: %s\n", e.Exam.Course.Title))

	return &feeds.Item{
		Title:       e.Exam.Course.Title,
		Link:        &feeds.Link{Href: url},
		Updated:     e.Updated,
		Description: b.String(),
	}
}
