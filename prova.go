package main

import (
	"fmt"
	"time"
)

type Materia struct {
	Codice  string
	Titolo  string
	Docente string
}

type Prova struct {
	DataEOra time.Time
	Tipo     string
	Materia  *Materia
}

// Creates a map from a slice of Prova.
// The key is the string representation of the Prova.
func proveHashMap(prove []Prova) map[string]Prova {
	hashMap := make(map[string]Prova, len(prove))
	for _, prova := range prove {
		hashMap[prova.String()] = prova
	}
	return hashMap
}

func (p *Prova) String() string {
	return fmt.Sprintf("%s %s %s %s", p.DataEOra, p.Tipo, p.Materia.Codice, p.Materia.Titolo)
}

// Implements sort.Interface for []Prova based on the DataEOra field.
type Prove []Prova

func (p Prove) Len() int { return len(p) }
func (p Prove) Less(i, j int) bool {
	return p[i].DataEOra.Before(p[j].DataEOra)
}
func (p Prove) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
