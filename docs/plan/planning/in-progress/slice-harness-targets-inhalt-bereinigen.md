# Slice harness-targets-inhalt-bereinigen: die Sensor- und Target-Dateien tragen ihren Inhalt einmal und im Ist-Zustand

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle — es gibt keine Closure-Bedingung jenseits der DoD
dieses Slice (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle
braucht).

**Bezug:** [`ADR-0051`](../../adr/0051-cicd-pipeline-github-actions.md) (CI/CD-Pipeline über GitHub Actions, mit Pin-Inventar),
[`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md) (Herkunft von
Aussagen: eine berichtigte Zahl trägt ihren Ursprung). Keine `LH-*`-Anforderung
ist berührt: der Slice ändert Harness-Dokumente, nicht das Produkt.

**Berührte Spec-Stellen:** —

**Verantwortlich:** pt9912 (Implementer-Agent im Auftrag).

**Autor:** pt9912 (Planner, Closure von `slice-harness-readme-zellen-kuerzen`).
**Datum:** 2026-10-05.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung**.

**Anlass.** `slice-harness-readme-zellen-kuerzen` zog die ausführlichen Zellen
der `harness/README.md` wortgleich in Sensor- und Target-Dateien um und durfte
keinen Wortlaut ändern. Seine Closure-Notiz (§7, Befund-Liste Punkte 2 bis 11)
nennt, was dabei als doppelt, widersprüchlich oder veraltet auffiel; 15 Dateien
tragen am Ende einen Abschnitt `## Fassung im Gate-Index` neben dem eigenen
Vertrag (Review-Befund F-1, `BEO-PGC/nachzug-laesst-ueberholten-text-stehen`).

**Ziel:** Jede Sensor- und Target-Datei trägt ihren Vertrag einmal, ohne
`## Fassung im Gate-Index`, und jede Aussage aus der Befund-Liste ist nachgemessen,
als übernommen gekennzeichnet oder gestrichen.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Die Index-Zeilen der `harness/README.md`.** Sie sind kurz und von der
  `structure`-Regel gehalten; ändert sich ein Kurzsatz durch eine Berichtigung,
  wird er mitgezogen, sonst bleibt die README unberührt. Die Vorlagen-Platzhalter
  der README (§Sensors-Musterzeile, Werkzeug-Musterzeilen, §Safety, §Leseordnung)
  trägt `slice-abgeleitete-dokumente-vorlagen-nachzug`.
- **Neue Gates oder Änderungen an Skripten.** Die Berichtigung betrifft Prosa;
  stellt sich eine Aussage als Defekt eines Werkzeugs heraus, ist das ein anderer
  Vorgang mit eigenem Slice.
- **Die `Accepted`-ADRs, die eine Zeile der alten Tabelle zitieren.** Sie frieren
  ihren Stand ein (`AGENTS.md` §3.5).

**Keine Mindestzahl.** Die vier Klassen des Ausschlusses sind ein Suchraster,
keine Ausfüll-Liste.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**. Alle Beleg-Angaben sind **Zusagen**
(„zu belegen durch …“).

- [ ] **Fassungen zusammenführen (Liefer-Punkt 1).** In den 15 Dateien mit
      `## Fassung im Gate-Index` (`git grep -l '^## Fassung im Gate-Index' --
      harness/`) ist der Inhalt der Fassung in die Abschnitte davor eingearbeitet,
      der Abschnitt samt Vorrang-Satz entfernt; ein Widerspruch zwischen beiden
      (z. B. „die vier“ gegen „die fünf“ ausgenommenen Dateien in
      `harness/sensors/handbuch-public-doc-check.md`) ist am Gegenstand gemessen
      entschieden, und Verweise auf `§Fassung im Gate-Index` (README-Bindung
      `make doc-tracked`) zeigen auf den neuen Ort. Zu belegen durch: der
      `git grep` oben liefert 0 Treffer; je Datei nennt der Bericht, welche Sätze
      zusammengeführt und welche gestrichen wurden.
- [ ] **Veraltete Aussagen berichtigen (Liefer-Punkt 2).** Jeder Punkt 2 bis 8,
      10 und 11 der Befund-Liste von `slice-harness-readme-zellen-kuerzen` §7
      ist erledigt: nachgemessen (Lauf genannt), als übernommen gekennzeichnet
      oder gestrichen; Chronik-Sprache steht im Ist-Zustand (`AGENTS.md` §3.7).
      Zu belegen durch eine Tabelle Punkt → Ausgang → Beleg im Bericht.
- [ ] `make gates` grün (Exit direkt ausgewertet, am Endstand;
      [`AGENTS.md`](../../../../AGENTS.md) §3.9), `make docs-check` Exit 0.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update: Träger, die `§Fassung im Gate-Index` oder eine berichtigte
      Zahl zitieren, sind nachgezogen (Suchlauf, `AGENTS.md` §3.13).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von der Slice-Closure selbst, solange die Roadmap unter *Offene Wellen* keine Welle führt.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| die 15 Dateien unter `harness/sensors/` und `harness/targets/` mit `## Fassung im Gate-Index` | update | Liefer-Punkt 1 |
| `harness/targets/image.md`, `workflows.md`, `sdk-pack.md`, `pin-stale.md`, `sdk-integration.md`, `examples.md`, `test-integration.md`; `harness/sensors/baseline-verify.md`; Kopfkommentar `.github/workflows/ci.yml` | update | Liefer-Punkt 2 (Befunde 2 bis 6, 8, 10, 11) |
| `harness/README.md` | update, nur falls ein Kurzsatz sich ändert | Bindung `make doc-tracked`, §1 |

Der Suchlauf (§3.13) wird beim Start gesetzt; bewegte Eigenschaften: die Menge
der Dateien mit `## Fassung im Gate-Index` und die berichtigten Zahlen.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `in-progress/` trägt keinen anderen Slice
(WIP-Limit 1); läuft `slice-baseline-6-14-0-dokumente-nachziehen` zuerst, liest
dieser Slice dessen Änderung an `harness/targets/pin-stale.md` mit.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): die Zusammenführung
  trägt in einem Review nicht alle 15 Dateien — dann Teilung nach
  `harness/sensors/` und `harness/targets/`.
- `in-progress` → `open` (blockiert): eine Befund-Aussage ist ohne Netz oder
  Release nicht messbar und lässt sich nicht als übernommen kennzeichnen —
  Frage an den Auftraggeber.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

Die DoD aus §2 ist vollständig, `make gates` endet mit Exit 0 am Stand der
Closure, und die Closure-Notiz in §7 trägt einen Lerneintrag.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang.

- **Zusammenführen verliert eine Zusage:** ein Satz der Fassung steht nur dort und
  fällt beim Streichen weg. Gegenmaßnahme: Abgleich der Eigennamen (Make-Ziele,
  Pfade, Zahlen) der Fassung gegen den Endstand der Datei. — **Ausgang:** offen
  bis zur Closure.
- **Eine Aussage ist am Host nicht messbar** (Release-Stand, Registry): dann
  *übernommen* mit Quelle statt Messung. — **Ausgang:** offen bis zur Closure.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<NNN>` **zitieren** statt neu
formulieren) · `grundlagen-traceability.md` §Herkunfts-Anker für
Steering-Loop-Regeln. Ging der Gegenstand an einen anderen Slice oder entfiel er,
trägt diese Sektion die Zeile `Gegenstand:` mit Kennung oder Grund.

*Der Plan füllt diese Sektion nicht; sie wird bei der Closure vor dem
`git mv` nach `done/` geschrieben.*

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung.

**Der Abschnitt selbst entfällt nie.**

**Vorgelagert — Sub-Area-Wahl prüfen:** [`harness/conventions.md`](../../../../harness/conventions.md)
§Modus-Deklaration führt eine Sub-Area (`*`, Kürzel `PGC`, Greenfield); die
berührten Pfade (`harness/`, `.github/workflows/ci.yml`) liegen in ihr.

**Vorgelagert — offene Beobachtungen sichten:** beim Start durch den Implementer
(Register `docs/plan/planning/observations/BEO-PGC/`). Bekannt bei der Anlage:
`BEO-PGC/nachzug-laesst-ueberholten-text-stehen` (verkörpert, Anlass F-1) und
`BEO-PGC/bindung-spalte-uneinheitlich-tief` (offen, 1×; berührt, wenn die
Bindung `make doc-tracked` umgestellt wird).

**Modus:** alle berührten Sub-Areas GF (`*`/`PGC`, Greenfield).
