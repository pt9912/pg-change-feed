# Welle 18: Spaltenauswahl — Antrags-Queue-Erweiterung, Assembler-Filterung, E2E-Beleg (`ADR-0059`)

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach
`done/` (neben ihre `welle-18-results.md`). Der Zustand ist die
Verzeichnis-Position — kein Status-Feld.

**Zielmeilenstein:** kein Meilenstein-Bezug.

**Verantwortlich:** —. **Datum:** 2026-09-14.

---

## 1. Welle-Ziel

`LH-FA-CFG-005` (Spaltenauswahl) bekommt den in
[ADR-0059](../../adr/0059-spaltenauswahl-mechanismus.md) entschiedenen,
vollständigen Umsetzungspfad: zwei neue Antragsarten
(`exclude_column`/`include_column`) auf der bestehenden Antrags-Queue aus
`ADR-0050`, Filterung im `Assembler` bei der Row-Image-Konstruktion (vor
jeder Serialisierung), Live-Reload ohne Prozess-Neustart über den bereits
etablierten Mechanismus, und ein realer E2E-Beleg gegen den laufenden
Feed-Container. `ADR-0059` (Accepted) trifft bereits alle fünf Teilfragen
inklusive Alternativen-Vergleich; diese Welle setzt sie in drei Slices um,
exakt nach der in `ADR-0059` §Konsequenzen/Folgepflicht empfohlenen
Zerlegung.

## 2. Trigger (Welle startet)

- `ADR-0059` liegt Accepted vor (bereits erfüllt — Architect-Entscheidung
  vom 2026-09-14).

## 3. Closure-Trigger (Welle schließt)

- `slice-066`, `slice-067`, `slice-068` liegen in `done/`.
- `make gates` grün.
- Der reale E2E-Beleg aus `slice-068` läuft grün — das *Mehr* gegenüber den
  einzelnen Slice-DoDs: `slice-068` kann ohne `slice-066` (Antrag existiert
  nicht) und ohne `slice-067` (kein Filter-Wirkort) nicht laufen, und erst
  der volle Rundlauf (SQL-Antrag → Live-Reload → gefiltertes Row Image)
  belegt das Zusammenspiel aller drei Schichten, das keine einzelne
  Slice-DoD für sich beansprucht.
- Closure-Notiz in `welle-18-results.md`.

## 4. Slices in dieser Welle

| Slice | Titel | Bezug |
|---|---|---|
| slice-066 | SQL-Funktionen + Antrags-Verarbeitung — `cdc.exclude_column`/`cdc.include_column` | [LH-FA-CFG-005](../../../../spec/lastenheft.md), [ADR-0059](../../adr/0059-spaltenauswahl-mechanismus.md) |
| slice-067 | Assembler-Filterung + Live-Reload-Verdrahtung — `TableBinding.ExcludedColumns` | [LH-FA-CFG-005](../../../../spec/lastenheft.md), [LH-FA-SCH-003](../../../../spec/lastenheft.md), [LH-FA-DAT-005](../../../../spec/lastenheft.md) |
| slice-068 | E2E-Beleg — Spaltenausschluss am laufenden Feed-Container | [LH-FA-CFG-005](../../../../spec/lastenheft.md) |

## 5. Abhängigkeiten

- `slice-067` baut auf `slice-066` auf (die Live-Reload-Methode wird von den
  neuen `applyAdministrationRequest`-Fällen aus `slice-066` aufgerufen).
- `slice-068` baut auf `slice-066` **und** `slice-067` auf (der E2E-Rundlauf
  braucht Antrag, Verarbeitung und Filterung zusammen) — kein
  Implementierungs-Blocker vor Beginn, aber eine harte Reihenfolge.
- Blockiert: keine andere offene Welle (`Nächste Wellen` ist derzeit leer;
  `welle-16`/`welle-17` sind unabhängige, bereits offene Wellen).
- Wird blockiert von: keiner (`ADR-0059` liegt bereits Accepted vor).

## 6. Out-of-Scope für diese Welle

- **Quellen-/musterweiter Spaltenausschluss über mehrere Tabellen hinweg**
  (`ADR-0059` Teilfrage 2, Option C) — die ADR verschiebt das ausdrücklich
  auf einen eigenen Re-Evaluierungs-Trigger (3×-Beobachtung oder geschärfte
  `LH-QA-SEC-004`); diese Welle liefert ausschließlich die pro-Tabelle-
  Granularität (Option B).
- **Boot-Zeit-Ausschlussfeld in `CDC_TABLES`/der optionalen YAML-Datei** —
  beide bleiben unverändert Erstaktivierungs-Seed für eine leere Datenbank
  (`ADR-0059` Teilfrage 1); ein initial ausgeschlossenes Feld ist durch
  diese Welle nicht ausgeschlossen, aber nicht ihr Gegenstand.
- **Echte technische Brücke für synchrone SQL→Go-Aufrufe** (FDW, `dblink`)
  — eigener, von `ADR-0050` geerbter Re-Evaluierungs-Trigger, unverändert
  unberührt von dieser Welle.
- **Dauerhaftigkeit des Ausschlussstandes über Prozess-Neustart und
  Bindungs-Zyklus** — diese Welle liefert den Ausschluss für den *laufenden*
  Prozess: die Antrags-Seite (`slice-066`), die Filterung samt Live-Reload
  (`slice-067`) und den realen E2E-Beleg des laufenden Pfads (`slice-068`).
  Dass der Stand keinen dauerhaften Träger hat — er geht bei jedem Neustart
  **und** bei einem `cdc.disable_table`/`enable_table`-Zyklus verloren —, ist
  im Review von `slice-067` real gefunden und vom Architect-Zug als **Lücke**
  entschieden worden, nicht als zulässige Grenze; er löst `ADR-0065`
  (`Supersedes ADR-0059`, nur die Dauerhaftigkeits-Aussage) und den
  wellenlosen Folge-Slice `slice-075` aus. `welle-18`s Closure-Trigger
  berührt das nicht (er verlangt den Beleg des laufenden Pfads, nicht
  Dauerhaftigkeit). Verdikt:
  [`docs/reviews/architect-verdict-spaltenausschluss-dauerhaftigkeit.md`](../../../reviews/architect-verdict-spaltenausschluss-dauerhaftigkeit.md).

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-traceability.md`
§Herkunfts-Anker für Steering-Loop-Regeln — dort die **Ruheort-Regel**: Die
beiden Zeiger unten sind so zu schreiben, wie sie vom Ruheort `done/` auflösen,
nicht vom Schreibort.

Ergebnis: [welle-18-results.md](welle-18-results.md), Geschwister im Ruheort `done/`
Zähler: [../observations/](../observations/)`BEO-PGC/`, eine Ebene über dem Ruheort
