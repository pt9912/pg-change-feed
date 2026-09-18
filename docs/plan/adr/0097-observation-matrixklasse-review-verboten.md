# ADR-0097: `observation`-Matrixklasse — Beobachtungs-Register zitiert Review-Reports über Kennung, nicht Adresse

**Status:** Accepted

**Datum:** 2026-09-18

**Autor:** Architect (pt9912)

**Bezug:** [`ADR-0094`](0094-review-matrixklasse-kennung-statt-adresse.md)
(dieselbe `review`-Matrixklasse, hier um eine zweite Quell-Klasse erweitert),
[`ADR-0073`](0073-zitat-korrektur-an-immutablen-dokumenten.md) (Zitat-Korrektur-Klasse),
Baseline-Regelwerk `modul-06-roadmap.md` §Das Beobachtungs-Register
(`evidence/<vorgangs-id>.md` ist „unveränderlich ab Merge, eine je Auftreten" —
dieselbe Dauerhaftigkeit wie eine `Accepted`-ADR, obwohl kein ADR-Dokument).

**Schärft:** — (Prozess-ADR ohne Spec-Stratum, wie `ADR-0094`).

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

`ADR-0094` hat die Fehlerklasse „ein permanentes Dokument zitiert live einen
Datei-Basisnamen unter `docs/reviews/`" für `Accepted`-ADRs geschlossen. Beim
Versuch, den realen `archive-welle`-Lauf durchzuführen, zeigte sich dieselbe
Klasse ein zweites Mal, in deutlich größerem Umfang: **83 Dateien** im
Beobachtungs-Register (`docs/plan/planning/observations/**`), überwiegend
`evidence/*.md`-Belege, zitieren Review-Reports per Basisname (Muster:
„Quelle: <Report-Dateiname unter der `review`-Klasse> (F-N)"). Diese Dateien sind laut
Baseline-Regelwerk genauso dauerhaft wie eine `Accepted`-ADR — „unveränderlich
ab Merge" —, obwohl sie kein ADR-Dokument sind und `AGENTS.md` §3.5 sie
deshalb nicht wörtlich erfasst. Derselbe Grund, der `ADR-0094` trug, trägt
hier ebenso: der zitierte Report wird bei Archivierung ohne Stub entfernt,
der Basisname bleibt in der dauerhaften Belegdatei stehen, und
`ai-harness-init archive-welle`s Hänger-Scan (reiner Basisnamen-Vergleich,
`ADR-0094` §Kontext) sperrt real jeden Lauf, der einen betroffenen Report
verschieben würde.

## Entscheidung

Wir ergänzen `.d-check.yml`s `matrix`-Sektion um eine Klasse `observation`
und eine Regel `{from: observation, to: review, allow: false}` — dieselbe
Form, dasselbe Ziel wie `ADR-0094`s `adr → review`-Regel, kein neues Konzept:

```yaml
matrix:
  classes:
    - name: observation                                          # NEU
      paths: ["docs/plan/planning/observations/**/*.md"]          # NEU
  rules:
    - {from: observation, to: review, allow: false}               # NEU
```

Die bestehende `review`-Klasse (Pfade, `token`) aus `ADR-0094` wird
unverändert wiederverwendet — nur eine zweite Quell-Klasse kommt hinzu. Kein
Provenance-Marker-Escape, aus demselben Grund wie in `ADR-0094` §Entscheidung:
ein `review`-Klassen-Mitglied hat keine dauerhafte Adresse.

**Bewusst keine Regel für** `{from: slice, to: review}` in dieser ADR: Ein
`done/`-Slice, der die REST-Kanten dieser Klasse betrifft (real gemessen: 4
Fälle, nicht 83), zitiert teils den eigenen, teils einen fremden Review — die
Unterscheidung braucht Einzelfallprüfung, kein blanket Verbot. Diese Lücke
ist benannt, nicht verschwiegen (§Konsequenzen).

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun, Register-Einträge bleiben Review-Prüfpflicht | kein Eingriff | dieselbe Erosion, die `ADR-0094` bereits für ADRs widerlegt hat — real 83 unbemerkte Fälle |
| **B — `observation`-Klasse, Regel `observation → review: false`, `review`-Klasse wiederverwendet — gewählt** | schließt die real größte Fehlerklasse (83 Dateien); keine neue `review`-Definition nötig | deckt nicht die 4 `done/`-Slice-Fälle (bewusst offen, siehe oben) |
| C — `observation`-Pfade in die bestehende `adr`-Klasse aufnehmen statt einer eigenen Klasse | ein Rule-Eintrag weniger | verwischt die Klassengrenze — ein Beobachtungs-Eintrag ist keine ADR und braucht keine `ADR-\d{4}`-Token-Erkennung; eine eigene Klasse ist die im Repo bereits geübte, klarere Form (`ADR-0094` §Entscheidung) |

**Fazit:** B.

## Konsequenzen

- Positiv: Die real größte offene Hänger-Quelle (83 von 127 betroffenen
  Dateien vor dieser Runde) wird künftig mechanisch erkannt.
- Negativ mit Grenze: `{from: slice, to: review}` bleibt unmechanisiert (4
  reale Fälle) — benannter, nicht behobener Rest, ggf. eigener Folge-Slice
  bei Bedarf.
- Folgepflicht (Implementer-Zug): (1) alle betroffenen Beobachtungs-Register-
  Dateien auf Kennung-Form umschreiben, (2) `.d-check.yml` um die Klasse/Regel
  ergänzen, (3) `make gates` grün, erst dann Commit.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| `d-check` Modul `matrix` (Klasse `observation`, Regel `{from: observation, to: review, allow: false}`) | Nach Aktivierung liefert `make docs-check` **0** `matrix-forbidden`-Befunde für `observation → review` — Zusage, real zu messen nach der Implementer-Fixrunde | `make docs-check` (in `make gates`) |

## Re-Evaluierungs-Trigger

Beobachtbarer Trigger: die verbleibenden `{from: slice, to: review}`-Fälle
erreichen wiederholtes Auftreten (Beobachtungs-Register, 3×) — dann Folge-ADR
mit `Supersedes ADR-0097` zur Frage, ob auch diese Kante mechanisiert wird.
Andernfalls permanent.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-18 | Accepted — `observation`-Matrixklasse und Regel `observation → review: false` beschlossen, Anlass: 83 Beobachtungs-Register-Dateien mit Live-Zitat auf `docs/reviews/**`, real beim Versuch des `archive-welle`-Laufs gefunden | Slice `slice-archive-altbestand-vollzug` |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0097`.
