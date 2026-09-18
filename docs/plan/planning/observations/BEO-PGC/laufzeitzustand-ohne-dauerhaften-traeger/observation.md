# BEO-PGC/laufzeitzustand-ohne-dauerhaften-traeger

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die
Verdrahtung administrativer Zustände zwischen Antrags-Datensatz und
laufender Erfassung, keine eigene Sub-Area im Sinn der Modus-Deklaration).

Die Beobachtung: Ein administrativer Zustand, den ein Antrags-Datensatz
(`applied`) als verarbeitet ausweist, lebt ausschließlich im
Prozessspeicher und geht beim Prozess-Neustart **und** bei einem
`cdc.disable_table`/`cdc.enable_table`-Zyklus verloren — ohne dass ein
Statuswert, ein Log oder ein Gate den Verlust meldet. Belegt am
Spaltenausschluss ([`LH-FA-CFG-005`](../../../../../../spec/lastenheft.md)):
`TableBinding.ExcludedColumns` liegt in der `Assembler`-Bindung;
`activatedTableBindings` (`internal/bootstrap/wiring.go`) baut den
Bindungsstand aus den Aktivierungs-Zeilen und der Schema-Version auf und
liest die Spalten-Antragsarten nicht; `RemoveBinding` (Deaktivierung)
entfernt die Bindung samt Stand, `AddBinding` (Aktivierung) legt sie ohne
Vorgänger neu an. Nach beiden Auslösern liest die `applied`-Zeile in
`cdc.administration_request` weiter Erfolg, während die Spalte wieder
erfasst und persistiert wird — genau der Zustand, den
[`ADR-0059`](../../../../adr/0059-spaltenauswahl-mechanismus.md) Teilfrage 3
Option D als Wirkort ausschließt („der sensible Wert wird nie serialisiert
und nie an die Persistenzschicht übergeben"). Ein persistierter Wert ist
nicht zurückholbar; die Zusage von `LH-QA-SEC-004` hängt damit an der
Prozesslebensdauer.

Deklaration: `slice-067` (Review-Findings F-1 und F-2,
Review zu `slice-067`;
die Findings sind als Entscheidungen an Planner/Architect gereicht, Modul 8
§Konflikt-Pfad). Der dauerhafte, tabellen-scoped Träger ist entschieden
([`ADR-0065`](../../../../adr/0065-spaltenausschluss-dauerhafter-traeger.md)),
der umsetzende Slice ist
`slice-075`
(`docs/plan/planning/open/`). Der Beleg trägt den Dateinamen seines Reviews statt der reinen
Slice-Nummer;
`slice-067` und sein Review sind **eine** Gelegenheit, kein zweites
Auftreten (Modul 6 §Ein Vorgang zählt einmal).
Zustand: **geplant** — die Entscheidung ist gefallen, der Träger ist
benannt; die Beobachtung braucht die 3×-Schwelle nicht (dieselbe direkte
Auflösung unter der Schwelle wie
`BEO-PGC/schema-evolution-nicht-dynamisch` und
`BEO-PGC/retention-keine-loeschausfuehrung`). Zähler (abgeleitet): 1×
(evidence/review-slice-067.md).
