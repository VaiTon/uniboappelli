package appelli

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	set "github.com/deckarep/golang-set/v2"
	"github.com/goodsign/monday"
)

const (
	localeItIT = monday.LocaleItIT
	examUrl    = "https://corsi.unibo.it/laurea/%s/appelli"
)

var (
	space = regexp.MustCompile(`\s+`)
)

func GetExamsUrl(degree string) string {
	return fmt.Sprintf(examUrl, degree)
}

func GetExams(degree string) ([]Exam, error) {
	prove := make([]Exam, 0, 20)

	bStart := 0
	for {
		url := GetExamsUrl(degree)
		if bStart > 0 {
			url = fmt.Sprintf("%s?b_start=%d", url, bStart)
		}

		res, err := http.Get(url)
		if err != nil {
			return nil, err
		}

		body := res.Body
		doc, err := goquery.NewDocumentFromReader(body)
		if err != nil {
			return nil, err
		}

		err = body.Close()
		if err != nil {
			return nil, err
		}

		coursesSel := doc.Find(".dropdown-component").First()
		if coursesSel.Length() == 0 {
			break
		}

		nCourses := coursesSel.Children().Length() / 2

		if nCourses == 0 {
			break
		}

		for i := 0; i < nCourses; i++ {

			tabId := fmt.Sprintf("#tab%d", i)
			tab := coursesSel.Find(tabId).First()
			course := parseCourse(tab)

			panelId := fmt.Sprintf("#panel%d", i)
			panel := coursesSel.Find(panelId).First()
			examsSel := panel.Find(".single-item")

			for i := range examsSel.Nodes {
				examSel := examsSel.Eq(i)

				exam, err := parseExam(examSel, course)
				if err != nil {
					return nil, err
				}
				prove = append(prove, exam)
			}

		}
		bStart += nCourses
	}

	return prove, nil
}

func parseCourse(tab *goquery.Selection) Course {
	code := tab.Find(".code").Text()
	code = strings.TrimSpace(code)

	teacher := tab.Find(".teacher").Text()
	teacher = strings.ReplaceAll(teacher, "\n", "")
	teacher = strings.TrimSpace(teacher)
	teacher = strings.TrimSpace(teacher)

	title := tab.Text()
	title = strings.ReplaceAll(title, code, "")
	title = strings.ReplaceAll(title, teacher, "")
	title = strings.ReplaceAll(title, "\n", "")
	title = space.ReplaceAllString(title, " ")
	title = strings.TrimSpace(title)

	return Course{code, title, teacher}
}

func parseExam(sel *goquery.Selection, course Course) (Exam, error) {
	tds := sel.Find("td")

	examType := tds.Eq(2).Text()
	examType = strings.ReplaceAll(examType, "\n", "")
	examType = strings.TrimSpace(examType)
	examType = space.ReplaceAllString(examType, " ")

	examTime := tds.Eq(0).Text()
	examTime = strings.ReplaceAll(examTime, " ore ", " ")
	examTime = space.ReplaceAllString(examTime, " ")
	examTime = strings.TrimSpace(examTime)

	dataParsed, err := monday.Parse("02 January 2006 15:04", examTime, localeItIT)
	if err != nil {
		return Exam{}, err
	}

	return Exam{dataParsed, examType, course}, nil
}

func Diff(new, old []Exam) (diff []Exam) {
	diff = make([]Exam, 0, 10)

	oldExamSet := set.NewSet(old...)
	for _, exam := range new {
		if !oldExamSet.Contains(exam) {
			diff = append(diff, exam)
		}
	}
	return
}
