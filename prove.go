package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/goodsign/monday"
	"github.com/rs/zerolog/log"
)

func getAppeliUrl(corso string) string {
	return fmt.Sprintf("https://corsi.unibo.it/laurea/%s/appelli", corso)
}

func getProve(corso string) ([]Prova, error) {
	prove := make([]Prova, 0, 1)

	bStart := 0
	for {
		url := getAppeliUrl(corso)
		if bStart > 0 {
			url = fmt.Sprintf("%s?b_start=%d", url, bStart)
		}

		res, err := http.Get(url)
		if err != nil {
			log.Error().Msg("Could not get page")
			os.Exit(1)
		}

		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			return nil, err
		}

		listaCorsi := doc.Find(".dropdown-component").First()
		if listaCorsi.Length() == 0 {
			log.Warn().Str("corso", corso).Msg("lista corsi non trovata")
			break
		}

		nCorsi := listaCorsi.Children().Length() / 2
		log.Debug().Msgf("corsi trovati: %d", nCorsi)
		if nCorsi == 0 {
			log.Debug().Str("corso", corso).Int("bStart", bStart).Msg("lista corsi vuota")
			break
		}

		for i := 0; i < nCorsi; i++ {

			tab := listaCorsi.Find(fmt.Sprintf("#tab%d", i)).First()
			panel := listaCorsi.Find(fmt.Sprintf("#panel%d", i)).First()

			code := tab.Find(".code").Text()
			docente := tab.Find(".docente").Text()

			titolo := tab.Text()
			titolo = strings.ReplaceAll(titolo, code, "")
			titolo = strings.ReplaceAll(titolo, docente, "")
			titolo = strings.TrimSpace(titolo)
			materia := Materia{code, titolo, docente}

			appelli := panel.Find(".single-item")

			for i := range appelli.Nodes {
				s := appelli.Eq(i)

				tds := s.Find("td")

				tipoProva := strings.TrimSpace(tds.Eq(2).Text())

				dataEOra := strings.TrimSpace(tds.Eq(0).Text())
				dataEOra = strings.ReplaceAll(dataEOra, " ore ", " ")

				locale, err := time.LoadLocation("Europe/Rome")
				if err != nil {
					return nil, err
				}
				dataParsed, err := monday.ParseInLocation("02 January 2006 15:04", dataEOra, locale, monday.LocaleItIT)
				if err != nil {
					return nil, err
				}

				prove = append(prove, Prova{dataParsed, tipoProva, &materia})
			}

		}

		bStart += nCorsi

	}

	return prove, nil
}
