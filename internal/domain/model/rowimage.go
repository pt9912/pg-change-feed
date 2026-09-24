package model

import (
	"bytes"
	"encoding/json"
)

// BuildRowImage trägt das JSON-Row-Image (`SPEC-002`, `ADR-0016`) einer
// Zeile: ein JSON-Objekt über die gesendeten Spalten-Werte in der
// Reihenfolge von `columns`. Werte sind JSON-Strings — der Text-Stand der
// Quelle geht unverändert in das Bild, ohne Typ-Interpretation; Namen und
// Werte tragen die Maskierung von `encoding/json`. `values[i]` gehört zu
// `columns[i]`; ein nil-Wert trägt NULL oder unverändertes TOAST und ist
// Abwesenheit (`LH-FA-CAP-008` Boundary), ebenso eine Spalte ohne Wert
// (`i >= len(values)`); überzählige Werte ohne Spalte bleiben unbeachtet.
// Ein in `excluded` geführter Spaltenname wird ebenso übersprungen
// (`LH-FA-CFG-005`, `LH-QA-SEC-004`): derselbe Abwesenheits-Vertrag wie beim
// nil-Wert, kein eigener Platzhalter (`LH-FA-DAT-005` Boundary), sein Wert
// wird nie serialisiert. Eine nil-Werteliste liefert kein Bild (nil, kein
// Fehler); eine Werteliste ohne tragenden Wert liefert `{}`.
//
// Die Funktion ist rein: sie hält keinen Zustand, liest ihre Eingaben nur
// und teilt keinen Speicher mit ihnen. Sie ist die eine Konstruktionsstelle
// für Row Images (`ADR-0111` Teilfrage 2): jeder Pfad, der ein Row Image
// erzeugt, ruft sie; gerufen wird sie vom Replication-Mapper.
func BuildRowImage(columns []string, values []*string, excluded []string) ([]byte, error) {
	if values == nil {
		return nil, nil
	}
	var image bytes.Buffer
	image.WriteByte('{')
	first := true
	for i, column := range columns {
		if i >= len(values) || values[i] == nil || containsName(excluded, column) {
			continue
		}
		name, err := json.Marshal(column)
		if err != nil {
			return nil, err
		}
		value, err := json.Marshal(*values[i])
		if err != nil {
			return nil, err
		}
		if !first {
			image.WriteByte(',')
		}
		first = false
		image.Write(name)
		image.WriteByte(':')
		image.Write(value)
	}
	image.WriteByte('}')
	return image.Bytes(), nil
}

// containsName meldet, ob ein Spaltenname in einer Ausschluss-Liste steht;
// die Liste trägt die Ausschlüsse einer Tabelle und bleibt klein, die
// lineare Suche damit ohne eigenen Index.
func containsName(names []string, name string) bool {
	for _, candidate := range names {
		if candidate == name {
			return true
		}
	}
	return false
}
