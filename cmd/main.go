package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	appelli "github.com/VaiTon/uniboappelli"
	"github.com/gorilla/feeds"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	debug := flag.Bool("debug", false, "sets log level to debug")
	flag.Parse()

	corsi := flag.Args()

	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	for _, corso := range corsi {
		rssForCorso(corso)
	}
}

func rssForCorso(corso string) {
	log.Info().Str("corso", corso).Msg("generating rss")

	prove, err := appelli.GetProve(corso)

	log.Info().Str("corso", corso).Int("nProve", len(prove)).Msg("got prove")

	if err != nil {
		log.Err(err).Str("corso", corso).Msg("Could not get appelli url")
		return
	}

	file, err := os.ReadFile("data.json")
	oldProve := make([]appelli.Prova, 0, 10)
	if err != nil && !os.IsNotExist(err) {
		log.Error().Msg("Could not read file")
		os.Exit(1)
	} else if err == nil {
		err = json.Unmarshal(file, &oldProve)
		if err != nil {
			log.Error().Msg("Could not parse file")
			os.Exit(1)
		}
	}

	appelliUrl := appelli.GetAppeliUrl(corso)

	feed := &feeds.Feed{
		Title:       "Notifiche Appelli",
		Description: "Notifiche Appelli",
		Author:      &feeds.Author{Name: "Notifiche Appelli", Email: ""},
		Link:        &feeds.Link{Href: appelliUrl},
		Created:     time.Now(),
	}

	feed.Items = make([]*feeds.Item, 0, len(prove))

	log.Debug().Msg("Sorting prove")
	var proveSort appelli.Prove = prove
	sort.Sort(sort.Reverse(proveSort))

	for _, p := range proveSort {
		builder := strings.Builder{}

		builder.WriteString(fmt.Sprintf("Data: %s\n", p.DataEOra))
		builder.WriteString(fmt.Sprintf("Tipo: %s\n", p.Tipo))
		builder.WriteString(fmt.Sprintf("Docente: %s\n", p.Materia.Docente))
		builder.WriteString(fmt.Sprintf("Materia: %s\n", p.Materia.Titolo))

		feed.Items = append(feed.Items, &feeds.Item{
			Title:       p.Materia.Titolo,
			Link:        &feeds.Link{Href: appelliUrl},
			Updated:     time.Now(),
			Description: builder.String(),
		})
	}

	if os.MkdirAll("rss", 0755) != nil {
		log.Error().Msg("Could not create RSS folder")
		os.Exit(1)
	}

	fileName := fmt.Sprintf("rss/%s.rss", corso)
	atomFile, err := os.Create(fileName)
	if err != nil {
		log.Error().Msg("Could not create RSS")
		os.Exit(1)
	}
	defer atomFile.Close()

	err = feed.WriteAtom(atomFile)
	if err != nil {
		log.Error().Msg("Could not create RSS")
		os.Exit(1)
	}
}

func diffProve(new []appelli.Prova, old []appelli.Prova) []appelli.Prova {
	diff := make([]appelli.Prova, 0, 1)

	log.Debug().Int("new", len(new)).Int("old", len(old)).Msg("diffing prove")

	oldMap := appelli.Prove.ToHashMap(old)

	for _, newProva := range new {
		if _, ok := oldMap[newProva.String()]; !ok {
			diff = append(diff, newProva)
		}
	}

	return diff
}
