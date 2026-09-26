package model

import (
	"bytes"
	"encoding/json"
	"fmt"

	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
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
// `rules` sind die Transformationsregeln der Tabelle (`LH-FA-CFG-007`,
// `SPEC-030`). Sie werden in derselben Schleife und erst nach dem Ausschluss
// und der Abwesenheits-Prüfung ausgewertet: eine ausgeschlossene oder
// abwesende Spalte erreicht keine Regel, ihr Schlüssel steht weder unter dem
// Quell- noch unter einem Zielnamen. Ein umbenannter Schlüssel behält die
// Position seiner Quellspalte. Eine leere Regelmenge liefert dieselben Bytes
// wie ein Aufruf ohne Regeln. Die Anwendbarkeit der Regeln (`SPEC-030`)
// prüft der Aufrufer vor dem Aufruf mit `Transformation.CheckApplicable`,
// damit eine Änderung als nicht anwendbar endet, bevor ein Bild entsteht.
// Eine Regel, deren Spalte nicht in `columns` steht, trifft keine Spalte und
// wirkt nicht. Die Funktion schreibt nie zwei gleichnamige Schlüssel:
// gleicht der Zielname einer treffenden Regel einer Spalte aus `columns`
// (ausgeschlossene und wertlose eingeschlossen) oder dem Zielnamen einer
// anderen im Bild umbenannten Spalte, liefert sie
// `ErrTransformationTargetCollides` und kein Bild.
//
// Die Funktion ist rein: sie hält keinen Zustand, liest ihre Eingaben nur
// und teilt keinen Speicher mit ihnen. Sie ist die eine Konstruktionsstelle
// für Row Images (`ADR-0111` Teilfrage 2): jeder Pfad, der ein Row Image
// erzeugt, ruft sie — der Replication-Mapper mit dem Regelstand der Bindung,
// der Backfill-Lauf mit dem Regelstand des Blocks.
func BuildRowImage(columns []string, values []*string, excluded []string, rules []Transformation) ([]byte, error) {
	if values == nil {
		return nil, nil
	}
	var image bytes.Buffer
	image.WriteByte('{')
	first := true
	renamed := make([]string, 0, 4)
	for i, column := range columns {
		if i >= len(values) || values[i] == nil || containsName(excluded, column) {
			continue
		}
		key, text := applyTransformations(rules, column, *values[i])
		if key != column {
			if containsName(columns, key) || containsName(renamed, key) {
				return nil, fmt.Errorf("%w: %s", domainerrors.ErrTransformationTargetCollides, key)
			}
			renamed = append(renamed, key)
		}
		name, err := json.Marshal(key)
		if err != nil {
			return nil, err
		}
		value, err := json.Marshal(text)
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

// containsName meldet, ob ein Spaltenname in einer Namensliste steht (die
// Ausschlüsse einer Tabelle, die Spalten einer Relation); die Listen bleiben
// klein, die lineare Suche damit ohne eigenen Index.
func containsName(names []string, name string) bool {
	for _, candidate := range names {
		if candidate == name {
			return true
		}
	}
	return false
}
