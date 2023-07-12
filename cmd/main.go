package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/samber/lo"
	"os"
	"path"
	"time"

	appelli "github.com/VaiTon/uniboappelli"
	"github.com/gorilla/feeds"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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
	degreeLog := log.With().Str("degree", degree).Logger()

	newExams, err := appelli.GetExams(degree)
	if err != nil {
		degreeLog.Err(err).Msg("could not get exams")
		return
	}
	degreeLog.Info().Int("exams", len(newExams)).Msg("exams found")

	// Create data directory if it does not exist

	dataPath := "data"
	err = os.MkdirAll(dataPath, 0755)
	if err != nil {
		degreeLog.Fatal().Err(err).Msg("could not create data directory")
	}

	dataFileName := fmt.Sprintf("%s.json", degree)
	dataFileName = path.Join(dataPath, dataFileName)

	oldProve := TimedExams{}
	data, err := readDataFile(dataFileName)
	if err != nil {
		degreeLog.Fatal().Err(err).Msg("could not read data file")
	} else if data != nil {
		oldProve = data
	}

	diff := appelli.Diff(newExams, oldProve.ToExams())
	if len(diff) > 0 {
		degreeLog.Info().Int("exams", len(diff)).Msg("new exams found")
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

	// Add data exams to feed
	feed.Items = append(feed.Items, examsToFeedItems(oldProve, pageUrl)...)

	// Add new exams to feed
	timedDiff := NewTimedExams(diff, time.Now())
	feed.Items = append(feed.Items, examsToFeedItems(timedDiff, pageUrl)...)

	if err = os.MkdirAll("rss", 0755); err != nil {
		degreeLog.Fatal().Err(err).Msg("could not create RSS folder")
	}

	fileName := fmt.Sprintf("rss/%s.rss", degree)
	atomFile, err := os.Create(fileName)
	if err != nil {
		degreeLog.Fatal().Err(err).Msg("could not create RSS")
	}

	err = feed.WriteAtom(atomFile)
	if err != nil {
		degreeLog.Fatal().Err(err).Msg("could not create RSS")
	}
	err = atomFile.Close()
	if err != nil {
		degreeLog.Fatal().Err(err).Msg("could not close dataFile")
	}

	oldProve = append(oldProve, timedDiff...)
	proveJson, err := json.Marshal(oldProve)
	if err != nil {
		degreeLog.Fatal().Err(err).Msg("Could not marshal newExams")
	}

	err = os.WriteFile(dataFileName, proveJson, 0644)
	if err != nil {
		degreeLog.Fatal().Err(err).Msg("Could not write dataFile")
	}
}

func examsToFeedItems(exams TimedExams, url string) []*feeds.Item {
	return lo.Map(exams, func(p TimedExam, _ int) *feeds.Item {
		return p.ToFeedItem(url)
	})
}

func readDataFile(path string) (exams TimedExams, err error) {
	// If dataFile does not exist, ignore error
	dataFile, err := os.ReadFile(path)

	if err != nil && !os.IsNotExist(err) {
		// If dataFile exists but there is an error, return it
		return nil, err
	} else if err != nil {
		// If dataFile does not exist, return empty array
		return nil, nil
	} else {
		err = json.Unmarshal(dataFile, &exams)
		if err != nil {
			return nil, err
		}
	}

	return exams, nil
}
