# ADR-0094: `review`-Matrixklasse — ADRs zitieren Review-Reports über Kennung, nicht Adresse

**Status:** Accepted

**Datum:** 2026-09-18

**Autor:** Architect (pt9912; anderer Kontext als der Planner-/Implementer-Zug,
der die 25 betroffenen ADRs zitat-korrigiert und `.d-check.yml` anschließend
ändert — `modul-08-agentenrollen.md` §Rollen-Regeln)

**Bezug:** `AGENTS.md` §3.5 (ADR-Immutabilität, Zitat-Korrektur-Ausnahme) ·
[`ADR-0073`](0073-zitat-korrektur-an-immutablen-dokumenten.md) (die Klasse
„Zitat-Korrektur"; diese ADR beantwortet die dort offen gelassene Frage nach
der Mechanisierung) · [`ADR-0045`](0045-commit-traceability-standing-gate.md),
[`ADR-0072`](0072-hostpaths-modul-aktiviert-ohne-ausnahme.md),
[`ADR-0084`](0084-sync-gate-fuer-generierte-artefakte.md) (Präzedenz: jede
bisherige Aktivierung eines neuen `d-check`-Moduls bzw. einer neuen
`matrix`-Regel trägt ihre eigene ADR) · Baseline-Regelwerk
`modul-06-roadmap.md` §Wellen-Closure-Prozedur (Review-Reports sind
Lauf-Belege ohne eigene Identität, bekommen keinen Stub) · der vorausgehende
Architect-Verdikt zur ADR-Review-Zitat-Korrektur (2026-09-18) — Bestandsaufnahme
der 25 betroffenen `Accepted`-ADRs, technischer Befund zum Hänger-Mechanismus
und der hier übernommene Matrix-Vorschlag (§7 dort).

**Schärft:** — (Prozess-ADR ohne Spec-Stratum, wie
[`ADR-0073`](0073-zitat-korrektur-an-immutablen-dokumenten.md),
[`ADR-0045`](0045-commit-traceability-standing-gate.md); ergänzt zugleich
`ADR-0073`s eigene Fitness-Function-Zeile, die dort weiterhin „—" lautet — der
Rückverweis steht bewusst nur hier, nicht als In-place-Ergänzung an
`ADR-0073` selbst, siehe §Konsequenzen).

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

`ADR-0073` definiert die Klasse „Zitat-Korrektur" und stuft sie als **Urteil,
kein Sensor** ein — ihre eigene Fitness-Function-Zeile lautet „—", weil kein
Werkzeug entscheiden kann, ob eine Korrektur nur das Zitat-Gerüst ändert.
Unabhängig davon besteht ein **separates, mechanisch fassbares** Teilproblem:
25 `Accepted`-ADRs verlinken oder erwähnen live den Datei-Basisnamen eines
Review-Reports oder Architect-Verdikts unter `docs/reviews/`. Baseline-Regelwerk
`modul-06-roadmap.md` §Wellen-Closure-Prozedur behandelt Review-Reports als
Lauf-Belege **ohne eigene Identität jenseits ihres Slice** — sie wandern bei
Archivierung in ein Wellen-Archiv und bekommen **keinen** Stub, anders als
Slices und Wellen. Ein lebender Markdown-Link oder Bare-Pfad aus einer
permanent im Repo stehenden `Accepted`-ADR macht den zitierten Report faktisch
nie archivierbar, ohne die ADR zu brechen.

Das ist kein hypothetisches Risiko: `ai-harness-init archive-welle --vorschau
welle-d-check` bricht real mit einer `[haenger]`-Sperre ab, weil mehrere
`Accepted`-ADRs Review-Reports zitieren, die der Lauf verschieben würde. Der
zugrundeliegende Scan (`internal/archive/scan.go`) läuft über reinen
Basisnamen-Teilstring-Vergleich, nicht über Link-Auflösung — jede Erwähnung
zählt, Markdown-Link, Inline-Code-Bare-Pfad oder nackter Fließtext
gleichermaßen. Delinken allein löst den Hänger deshalb nicht; nötig ist eine
Zitierform, die den Datei-Basisnamen an keiner Stelle mehr reproduziert.

Der vorausgehende Architect-Verdikt zur ADR-Review-Zitat-Korrektur hat diesen
Bestand vollständig aufgenommen (25 ADRs, Fundort-Klassen, Ziel-Zitierform)
und geprüft, ob die Fehlerklasse zusätzlich zur einmaligen Bereinigung
**mechanisch** über `make docs-check` erkennbar gemacht werden soll und kann.
Ergebnis dort: ja, über eine neue `matrix`-Klasse `review` mit `token`, analog
zur bestehenden `slice`-Klasse. Diese ADR trifft die dafür nötige
Grundsatzentscheidung — `AGENTS.md` §3.6 verlangt zwar nur für
Schwellen-*Senkungen* zwingend eine ADR, aber der Bestand dieses Repos zeigt
durchgängig ein stärkeres, bereits geübtes Muster: jede bisherige Aktivierung
eines neuen `d-check`-Moduls oder einer neuen `matrix`-Regel trägt ihre
eigene ADR (`ADR-0045` für das Modul `commits`, `ADR-0072` für das Modul
`hostpaths`, `ADR-0084` für das Sync-Gate) — unabhängig davon, ob die
Änderung eine Lockerung oder eine Verschärfung war.

## Entscheidung

Wir aktivieren eine neue `matrix`-Klasse `review` und eine neue Regel
`{from: adr, to: review, allow: false}` in `.d-check.yml`. Der exakte Diff
(anzuwenden durch die nachfolgende Implementer-Fixrunde, **nach** der
Zitat-Korrektur der 25 betroffenen ADRs, nicht davor — sonst liefert
`make docs-check` sofort 25+ `matrix-forbidden`-Befunde, bevor der
Implementer zu schreiben beginnt):

```yaml
matrix:
  classes:
    - name: spec
      paths: [spec/lastenheft.md, spec/pflichtenheft.md, spec/architecture.md]
    - name: adr
      paths: ["docs/plan/adr/[0-9]*.md"]
      token: 'ADR-\d{4}'
    - name: slice
      paths: ["docs/plan/planning/**/slice-*.md"]
      token: 'slice-\d{3}'
    - name: review                                    # NEU
      paths: ["docs/reviews/*.md"]                     # NEU
      token: 'docs/reviews/[\w.-]+\.md'                 # NEU — fängt auch Bare-Pfad-Erwähnungen ohne Link
  rules:
    - {from: spec, to: adr, allow: false}
    - {from: spec, to: slice, allow: false}
    - {from: adr, to: slice, allow: false}
    - {from: adr, to: review, allow: false}            # NEU
```

Bewusst **keine** neue Regel für `slice → review` (ein Slice zitiert seinen
eigenen Review-Report routinemäßig und beide werden gemeinsam archiviert,
kein Hänger-Risiko) und für `review → review` (Delta-Reviews verweisen auf
ihren Ausgangs-Report, beide liegen in derselben Archiv-Einheit). `adr → adr`
bleibt unverändert erlaubt (Supersedes-Kette, Bezug-Verweise).

Bewusst **kein** Provenance-Marker-Escape für `adr → review`, anders als bei
der `slice`-Klasse: Ein `slice`-Klassen-Mitglied behält beim Archivieren
einen Stub und bleibt daher als reine Provenance-Kennung dauerhaft
auflösbar — der `<!-- d-check:status-provenance -->`-Marker ist dort ein
legitimer Escape-Hatch. Ein `review`-Klassen-Mitglied hat **keine** dauerhafte
Adresse, das ist genau der Grund für diese Entscheidung. Ein Escape-Hatch
würde die Klasse aushöhlen, die hier geschlossen wird: jede künftige
`adr → review`-Referenz, auch eine „nur als Provenance gemeinte", läuft über
die Kennung-Form (Themen-Prosa bzw. `slice-NNN` + Finding-ID, wie vom
vorausgehenden Verdikt für die 25 Bestandsfälle bereits festgelegt), nicht
über den Marker.

Die mechanisierte Regel prüft ausschließlich das **Vorhandensein einer
Adresse** (Datei-Basisname als Token). Sie prüft **nicht**, ob eine
Umformulierung die Aussage-Semantik einer Zitat-Korrektur wahrt — das bleibt
weiterhin Review-Prüfpflicht (`AGENTS.md` §3.12 Instanz B, wie in `ADR-0073`
§Konsequenzen bereits benannt: „der Wächter ist ein Urteil, kein Sensor").
Diese ADR mechanisiert also einen **Teil** des in `ADR-0073` beschriebenen
Urteils, nicht das ganze Urteil.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — kein Sensor; die 25 ADRs einmalig bereinigen, künftige `adr → review`-Referenzen bleiben allein Review-Prüfpflicht | kein Eingriff in `.d-check.yml`, kleinster Schnitt | genau diese Lücke hat zum realen Bestand von 25 unbemerkten Verstößen geführt, die erst über einen dedizierten Architect-Zug aufgedeckt wurden — Review-Prüfpflicht allein hat die Klasse in der Vergangenheit nicht gehalten |
| B — Escape-Marker-Ausnahme wie bei der `slice`-Klasse (`<!-- d-check:status-provenance -->` auch für `review`-Token zulassen) | konsistent mit dem bestehenden Escape-Muster, geringerer Umschreibungsdruck auf künftige ADR-Autoren | höhlt die Klasse aus, die diese Entscheidung gerade schließt: anders als ein `slice`-Mitglied behält ein `review`-Mitglied beim Archivieren keinen Stub und bleibt nicht dauerhaft auflösbar — ein Marker würde eine strukturell nicht vorhandene Dauerhaftigkeit vortäuschen |
| **C — neue `matrix`-Klasse `review` mit `token`, Regel `adr → review: false`, ohne Escape-Marker — gewählt** | schließt exakt die Fehlerklasse, die real 25× aufgetreten ist; nutzt dieselbe, bereits im Repo geübte Klassen-/Token-Form wie `slice`; lässt `slice → review` und `review → review` bewusst unberührt (kein Overreach) | mechanisiert nur das Vorhandensein einer Adresse, nicht die Klassenzugehörigkeit „ist es wirklich eine reine Zitat-Korrektur" — dieser Rest bleibt Review-Prüfpflicht |

**Fazit:** C. A hat sich am realen Bestand bereits als unzureichend erwiesen;
B würde die neue Klasse durch denselben Mechanismus aushöhlen, den sie
schließen soll.

## Konsequenzen

- Positiv: Die Fehlerklasse „ADR verlinkt live in `docs/reviews/`" wird
  künftig mechanisch über `make docs-check` erkannt, nicht mehr nur
  einmalig von Hand bereinigt.
- Positiv: Nach der vorausgehenden Zitat-Korrektur der 25 betroffenen ADRs
  entsperrt diese Aktivierung `ai-harness-init archive-welle` für genau
  diese 25 Fälle. Das ist **nicht notwendig repo-weit** — der zugrundeliegende
  Hänger-Scan läuft über alle getrackten Dateien; bare Erwähnungen desselben
  Review-Basisnamens außerhalb der 25 ADRs (`done/`-Records,
  `docs/reviews/**` selbst, Nicht-Markdown-Dateien) bleiben ein offener Rest
  und sind nicht Gegenstand dieser Entscheidung.
- Positiv: Konsistente Form mit der bestehenden `slice`-Klasse — dieselbe
  `paths`/`token`-Struktur, keine neue Konzept-Klasse in `.d-check.yml`.
- Negativ: Künftige ADRs dürfen `docs/reviews/`-Pfade an keiner Stelle mehr
  live zitieren (Markdown-Link, Inline-Code-Bare-Pfad, nackter Fließtext) —
  Kennung-Form (Themen-Prosa für Architect-Verdikte ohne Slice-Nummer,
  `slice-NNN` + Finding-ID für Reports mit Slice-Bezug) ist ab sofort Pflicht,
  nicht mehr nur empfohlene Praxis.
- Negativ mit Grenze: Die Regel prüft nur das Ergebnis (Adresse vorhanden
  oder nicht), nicht die Klassenzugehörigkeit einer Umformulierung als
  echte Zitat-Korrektur — dieser Teil bleibt Review-Prüfpflicht, wie in
  `ADR-0073` bereits benannt.
- Folgepflicht (Implementer-Zug, in dieser Reihenfolge): (1) alle 25
  betroffenen ADR-Texte auf Kennung-Form umschreiben, (2) `.d-check.yml` um
  die `review`-Klasse und -Regel ergänzen, (3) `make docs-check`/`make gates`
  grün, erst dann Commit/Closure.
- Zu `ADR-0073`: Diese ADR beantwortet dessen offen gelassene
  Mechanisierungsfrage, ändert aber nichts an `ADR-0073`s eigener,
  unberührbaren Aussage — weder an §Entscheidung noch an der
  Fitness-Function-Zeile selbst. Eine bloße Textkorrektur an einer
  Fitness-Function-Zeile ist in diesem Repo bereits als **nicht**
  In-place-fähig behandelt worden (`ADR-0067` brauchte für exakt diesen
  Korrekturtyp eine eigene Supersedes-ADR). Da `ADR-0073`s Klausel selbst
  unverändert bleibt, ist ein `Supersedes` hier nicht nötig — der
  Rückverweis lebt stattdessen ausschließlich in dieser ADR (§Schärft, §Bezug
  oben), `ADR-0073` bleibt textlich unangetastet.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| `d-check` Modul `matrix` (Klasse `review`, Regel `{from: adr, to: review, allow: false}`) | Nach Aktivierung liefert `make docs-check` **0** `matrix-forbidden`-Befunde für `adr → review` — **Zusage**, real zu messen erst nach der Implementer-Fixrunde (die `.d-check.yml`-Änderung existiert zum Zeitpunkt dieser ADR noch nicht), nicht als bereits erfüllte Tatsache behauptet (`AGENTS.md` §3.12 Instanz B) | `make docs-check` (in `make gates`) |

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

Beobachtbare Trigger: **(a)** die Implementer-Fixrunde liefert nach
Aktivierung weiterhin `matrix-forbidden`-Befunde außerhalb der 25 bereits
bekannten ADRs (der in §Konsequenzen benannte offene Rest wird real
angetroffen) — dann eigener Folge-Slice zur Bereinigung, keine Änderung
dieser ADR nötig; **(b)** ein begründeter Bedarf für einen
Provenance-Marker-Escape bei `adr → review` entsteht (z. B. eine Referenzform,
die dauerhaft auflösbar bleibt, obwohl der Report selbst keinen Stub bekommt)
— dann Folge-ADR mit `Supersedes ADR-0094`, die diese Entscheidung
klassenweise nachzieht. Andernfalls permanent.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-18 | Accepted — `review`-Matrixklasse und Regel `adr → review: false` beschlossen, Anlass: 25 `Accepted`-ADRs mit Live-Zitat auf `docs/reviews/**`, `archive-welle`-Hänger real reproduziert | der vorausgehende Architect-Verdikt zur ADR-Review-Zitat-Korrektur (2026-09-18) |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0094` (Baseline-Regelwerk `modul-04-adrs.md` §Hard Rule für
Accepted-ADRs).
