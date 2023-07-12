package appelli

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/goodsign/monday"
	"net/http"
	"strings"
)

const (
	localeItIT = monday.LocaleItIT
	examUrl    = "https://corsi.unibo.it/laurea/%s/appelli"
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
	teacher := tab.Find(".teacher").Text()

	title := tab.Text()
	title = strings.ReplaceAll(title, code, "")
	title = strings.ReplaceAll(title, teacher, "")
	title = strings.TrimSpace(title)

	return Course{code, title, teacher}
}

func parseExam(sel *goquery.Selection, course Course) (Exam, error) {
	tds := sel.Find("td")

	examType := strings.TrimSpace(tds.Eq(2).Text())

	examTime := strings.TrimSpace(tds.Eq(0).Text())
	examTime = strings.ReplaceAll(examTime, " ore ", " ")

	dataParsed, err := monday.Parse("02 January 2006 15:04", examTime, localeItIT)
	if err != nil {
		return Exam{}, err
	}

	return Exam{dataParsed, examType, &course}, nil
}

func Diff(new, old Exams) Exams {
	oldMap := old.ToHashMap()
	diff := make([]Exam, 0, 1)
	for _, newExam := range new {
		_, found := oldMap[newExam.String()]
		if !found {
			diff = append(diff, newExam)
		}
	}

	return diff
}