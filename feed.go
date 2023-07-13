package main

import (
	"fmt"
	"github.com/gorilla/feeds"
	"github.com/rs/zerolog"
	"os"
	"time"
)

func createFeedFile(
	degree, pageUrl string,
	newExams TimedExams,
	oldProve TimedExams,
	degreeLog zerolog.Logger,
) {

	feed := &feeds.Feed{
		Title:       "Notifiche Appelli",
		Description: "Notifiche Appelli",
		Author:      &feeds.Author{Name: "Notifiche Appelli", Email: ""},
		Link:        &feeds.Link{Href: pageUrl},
		Created:     time.Now(),
	}

	// Create feed items array
	feed.Items = make([]*feeds.Item, 0, len(newExams)+len(oldProve))

	// Add data exams to feed
	feed.Items = append(feed.Items, examsToFeedItems(oldProve, pageUrl)...)

	// Add new exams to feed
	feed.Items = append(feed.Items, examsToFeedItems(newExams, pageUrl)...)

	if err := os.MkdirAll("rss", 0755); err != nil {
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

}

func examsToFeedItems(exams TimedExams, url string) (items []*feeds.Item) {
	items = make([]*feeds.Item, 0, len(exams))
	for _, exam := range exams {
		items = append(items, exam.ToFeedItem(url))
	}
	return
}
