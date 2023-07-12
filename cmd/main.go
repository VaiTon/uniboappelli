package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	appelli "github.com/VaiTon/uniboappelli"
	"github.com/gorilla/feeds"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

func main() {

	debugFlag := flag.Bool("debugFlag", false, "sets log level to debugFlag")
	jsonFlag := flag.Bool("jsonFlag", false, "output jsonFlag to console instead of pretty printed text")
	flag.Parse()

	// Set global log level to debugFlag or info
	if *debugFlag {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	if !*jsonFlag {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	degrees := flag.Args()
	if len(degrees) == 0 {
		log.Fatal().Msg("No degrees specified")
	}

	for _, degree := range degrees {
		log.Debug().Str("degree", degree).Msg("starting analysis")
		doDegree(degree)
	}
}

func doDegree(degree string) {
	logger := log.With().Str("degree", degree).Logger()

	newExams, err := appelli.GetExams(degree)
	if err != nil {
		logger.Err(err).Msg("could not get exams")
		return
	}
	logger.Info().Int("exams", len(newExams)).Msg("exams found")

	// Create data directory if it does not exist

	dataPath := "data"
	err = os.MkdirAll(dataPath, 0755)
	if err != nil {
		logger.Fatal().Err(err).Msg("could not create data directory")
	}

	dataFileName := fmt.Sprintf("%s.json", degree)
	dataFileName = path.Join(dataPath, dataFileName)

	dataFile, err := os.ReadFile(dataFileName)

	oldProve := make(TimedExams, 0, 10)
	// If dataFile does not exist, ignore error
	if err != nil && !os.IsNotExist(err) {
		logger.Fatal().Err(err).Msg("Could not read dataFile")
	} else if err == nil {
		err = json.Unmarshal(dataFile, &oldProve)
		if err != nil {
			logger.Fatal().Msg("Could not parse dataFile")
		}
	}

	diff := appelli.Diff(newExams, oldProve.ToExams())
	if len(diff) > 0 {
		logger.Info().Int("exams", len(diff)).Msg("new exams found")
	}

	pageUrl := appelli.GetExamsUrl(degree)

	feed := &feeds.Feed{
		Title:       "Notifiche Appelli",
		Description: "Notifiche Appelli",
		Author:      &feeds.Author{Name: "Notifiche Appelli", Email: ""},
		Link:        &feeds.Link{Href: pageUrl},
		Created:     time.Now(),
	}

	feed.Items = make([]*feeds.Item, 0, len(diff)+len(oldProve))

	// Add old exams to feed
	feed.Items = append(feed.Items, examsToFeedItems(oldProve, pageUrl)...)

	// Add new exams to feed
	timedDiff := NewTimedExams(diff, time.Now())
	feed.Items = append(feed.Items, examsToFeedItems(timedDiff, pageUrl)...)

	if err = os.MkdirAll("rss", 0755); err != nil {
		logger.Fatal().Err(err).Msg("could not create RSS folder")
	}

	fileName := fmt.Sprintf("rss/%s.rss", degree)
	atomFile, err := os.Create(fileName)
	if err != nil {
		logger.Fatal().Err(err).Msg("could not create RSS")
	}

	err = feed.WriteAtom(atomFile)
	if err != nil {
		logger.Fatal().Err(err).Msg("could not create RSS")
	}
	err = atomFile.Close()
	if err != nil {
		logger.Fatal().Err(err).Msg("could not close dataFile")
	}

	oldProve = append(oldProve, timedDiff...)
	proveJson, err := json.Marshal(oldProve)
	if err != nil {
		logger.Fatal().Err(err).Msg("Could not marshal newExams")
	}

	err = os.WriteFile(dataFileName, proveJson, 0644)
	if err != nil {
		logger.Fatal().Err(err).Msg("Could not write dataFile")
	}
}

func examsToFeedItems(exams TimedExams, url string) []*feeds.Item {
	items := make([]*feeds.Item, 0, len(exams))
	for _, p := range exams {
		b := strings.Builder{}

		b.WriteString(fmt.Sprintf("Data: %s\n", p.Updated))
		b.WriteString(fmt.Sprintf("Type: %s\n", p.Exam.Type))
		b.WriteString(fmt.Sprintf("Teacher: %s\n", p.Exam.Course.Teacher))
		b.WriteString(fmt.Sprintf("Course: %s\n", p.Exam.Course.Title))

		items = append(items, &feeds.Item{
			Title:       p.Exam.Course.Title,
			Link:        &feeds.Link{Href: url},
			Updated:     p.Updated,
			Description: b.String(),
		})
	}

	return items
}

type TimedExam struct {
	Updated time.Time
	Exam    appelli.Exam
}

type TimedExams []TimedExam

func (t TimedExams) ToExams() appelli.Exams {
	return lo.Map(t, func(p TimedExam, _ int) appelli.Exam {
		return p.Exam
	})
}

func NewTimedExams(exams appelli.Exams, time time.Time) TimedExams {
	return lo.Map(exams, func(exam appelli.Exam, _ int) TimedExam {
		return TimedExam{Updated: time, Exam: exam}
	})
}
