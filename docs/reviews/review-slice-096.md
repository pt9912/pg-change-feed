# Review-Report: slice-096 — Konfigurationsdatei-Nachzug (Durchleitung und Feldmenge) · 2026-09-17

**Review-Art:** Code — geprüft gegen **Plan und Entscheidungen**. DoD-/Spec-Konformität
(Verifier), §6-Ausgänge, Register und die drei Paarungen sind **nicht** Gegenstand.

**Gegenstand:** `b8fd926` (Parent `3c8b6ff`) — **ein** Commit, **4** Dateien, +359/−65:
`internal/bootstrap/config_file.go`, `internal/bootstrap/wiring.go`,
`internal/bootstrap/config_file_internal_test.go`, `docs/user/benutzerhandbuch.md`.
Kein Gate, keine Spec, kein `THRESHOLD`, keine Naht.

**Skill:** `.harness/skills/reviewer.md` · **Datum:** 2026-09-17.

**Eingangs-Kontext:** Slice-Plan `slice-096`, [ADR-0088](../plan/adr/0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md),
`ADR-0052`, `ADR-0073` (Zitat-Korrektur), `AGENTS.md` §3.7/§3.11/§3.12/§3.13,
`harness/conventions.md`, `SPEC-016`.

**Mess-Grundlage:** Arbeitskopien beider Commits in einem Scratch-Verzeichnis außerhalb des
Repos — der Working Tree trug zum Review-Zeitpunkt fremde, uncommittete Testdateien. Alle
Läufe netzlos im gepinnten Toolchain-Container.

---

## Findings

### F-1 — `ADR-0088`s §Fitness Function nennt `make docs-check` als Träger der Feldmengen-Paarung; der genannte Lauf trägt sie nicht

- `kategorie`: **HIGH**
- `quelle`: Skill HIGH „Beleg trägt seinen Satz nicht" (`BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht`,
  5×) · `AGENTS.md` §3.12 Instanz B
- `pfad`: `docs/plan/adr/0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md:333`
- `befund`: Die Zeile sagt, `make docs-check` prüfe, dass „die Feldmenge in `SPEC-016`, im
  Handbuch §5.2 und die Schlüssel-Prüfung im Code … dieselben Namen nennen". Eigene Messung:
  (a) `.d-check.yml:8` führt `modules: [links, anchors, ids, matrix, versions, structure,
  hostpaths]` — **kein** Modul liest Go-Quelltext, und **keine** `structure`-Regel stellt zwei
  Dokumente einander gegenüber. (b) Probe in einer Arbeitskopie: `grpc_addr` → `grpc_addr_x` in
  der `SPEC-016`-Feldtabelle, `nats_url` aus der Handbuch-Klassenliste, `http_addr`/`grpc_addr`
  aus dem §5.2-Beispiel gestrichen → `make docs-check` bleibt **grün** („788 Datei(en), 0
  Befunde", Exit 0); unmutiert ebenso. Der genannte Beleg liefert mit und ohne Paarung dasselbe
  Ergebnis. Die drei Träger nennen heute tatsächlich dieselben Namen — die **Aussage** ist wahr,
  ihre **Stütze** nicht.
- `verifizierbar`: **ja**
- `klasse`: „Beleg trägt seinen Satz nicht (genanntes Tool trägt die Aussage nicht)"
- **Aktor ist der Architect, nicht der Implementer.** Der Träger liegt außerhalb des Diffs;
  `ADR-0073` §Entscheidung 1 nimmt „die Fitness-Function-Regeln" **ausdrücklich** von der
  Zitat-Korrektur aus → inhaltliche Änderung, also Folge-ADR mit `Supersedes`.
  **Kein Reviewer→Implementer-Pfeil.** Für diesen Slice heißt es: der **Verifier** liest die
  Zeile als Beleg für LP2/LP3, und sie prüft nichts.

### F-2 — Drei Anker in `ADR-0088` lösen nach diesem Diff nicht mehr auf; der Suchlauf fand den Symbolnamen, die zwei Zeilen-Lokatoren nicht

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.13 · `ADR-0073` §Entscheidung 1
- `pfad`: `docs/plan/adr/0088-…md:40` (§Bezug), `:82` und `:85` (§Kontext (3))
- `befund`: Der Symbolname `forbiddenFileDSNKeys` (Z. 40) existiert nach der Umbenennung nicht
  mehr — das hat der Bericht gemeldet, richtig (§3.13: fremden Träger melden, nicht still
  mitändern). **Nicht** gemeldet: die zwei Zeilen-Lokatoren derselben ADR. Am Parent gemessen:
  `wiring.go:276-280` = genau die fünf `cfg.X = getenv(...)`-Zeilen, `config_file.go:148` =
  `func mergeConfig`, `:193` = ihr `}`. Am Review-Commit liegen die fünf Zeilen bei **282-286**,
  `mergeConfig` bei **167-222** — der zitierte Bereich beginnt jetzt in `overrideString` und
  endet mitten in `mergeConfig`. Der Suchlauf nach der bewegten Eigenschaft sucht
  **Symbolnamen**; Lokatoren sind Zahlen und fallen durch.
- `verifizierbar`: **ja**
- `klasse`: „Anker löst nach Umbenennung/Layout-Verschiebung nicht auf"
- **Stellungnahme (verlangt):** §Bezug Z. 40 und die zwei Lokatoren sind **Zitat-Korrektur**
  nach `ADR-0073` (gebrochene Referenzform bzw. Zeilen-Lokator bei unverändertem Referenten) —
  in-place zulässig; Belegform: die Commit-Message nennt `ADR-0073`, `ADR-0088` erhält **eine**
  §Geschichte-Zeile. **Nicht** korrigieren: die Ist-Stand-Zeile (Z. 71, „führt **3** Schlüssel")
  und die §Geschichte-Zeile (Z. 367) — sie zitieren den **Parent-Stand als Messung**; die neuen
  Namen dort einzutragen machte die Messung falsch. Die Klasse ist im Repo bereits getragen und
  deshalb **LOW**: `ADR-0082` führt heute zwei ebenso veraltete Lokatoren.

### F-3 — Der Testkommentar bindet „die Prädikate, die `Run` tatsächlich zieht" nur zu einem Drittel

- `kategorie`: LOW
- `quelle`: Skill „Beleg trägt seinen Satz nicht" · §3.12 Instanz B
- `pfad`: `internal/bootstrap/config_file_internal_test.go:458-467` gegen `wiring.go:704`, `:759`
- `befund`: `changeStreamEnabled` ist **real** gebunden (Mutation `return true` → Suite rot,
  Exit 1); die zwei `!= ""`-Vergleiche stehen im Test als **eigene** Ausdrücke über `cfg`.
  Mutiert man `Run`s eigene Grenze (`if cfg.HTTPAddr != ""` → `!= "disabled"`), bleibt die
  gesamte `bootstrap`-Suite **grün** (Exit 0) — die Grenzen, „die `Run` tatsächlich zieht", sind
  assertion-seitig nachgebaut, nicht gebunden. Die **benannte** Grenze des Kommentars (kein
  laufender Server) ist dagegen ehrlich: `Run` konstruiert den Store (`wiring.go:439`) vor dem
  HTTP-Server (`:704`). Der Kommentar nennt eine Grenze — nur nicht diese.
- `verifizierbar`: **ja** — eigene Mutation, Suite-Exit 0
- `klasse`: „Beleg re-statement statt Bindung an den Verzweigungspfad"

### F-4 — Der Bericht nennt für „kein Sensor fand es" das falsche Modul; die tragende Grenze ist eine andere — und sie ist benannt

- `kategorie`: INFO · `quelle`: §3.13 · `harness/sensors/docs-check.md` §Grenze 7
- `pfad`: Bericht des Laufs gegen `harness/sensors/docs-check.md:85-90`
- `befund`: „`docs-check` führt das `citations`-Modul nicht" trifft zu, ist aber nicht der Grund:
  das Modul **allein** gefahren meldet „788 Datei(en), 0 Befunde", Exit 0 — es hätte den
  veralteten Symbolnamen ebenso wenig gefunden. Die tragende Grenze steht benannt da: „**Kein
  Modul prüft einen Symbol- oder Funktionsnamen**". Für die zwei Zeilen-Lokatoren ist die
  Nachbar-Grenze `codepaths` (aus) — es prüft Pfade, keine Zeilenbereiche. **Benannte Grenze,
  keine Lücke** — der benannte Träger ist nur der falsche.
- `klasse`: „Grenze benannt, aber am falschen Träger verortet"

### F-5 — Plan §3 nennt eine Testdatei, die es nicht gibt

- `kategorie`: INFO · `quelle`: Maintainability (Plan-vs-Diff)
- `pfad`: `slice-096-…md:137`
- `befund`: §3 führt `internal/bootstrap/config_file_test.go`; im Baum liegen
  `config_file_internal_test.go` und `config_file_rest_internal_test.go`. §3 erklärt sich als
  Kandidatenliste, die „der Implementer-Agent in seinem ersten Lauf erweitert" — geschehen ist
  das nicht. Ohne Wirkung auf Code, Träger oder DoD.
- `klasse`: „Plan-Kandidatenliste nicht nachgezogen"

---

## Was der Diff richtig macht (nachgemessen, nicht übernommen)

**(1) Der Defekt ist wirklich behoben, und die Durchleitung wirkt.** Eigene Zählung: `Config`
trägt **15** Felder; `mergeConfig` erreichte am Parent **10** distinkte Felder, heute **15**;
`ConfigFromEnv` liest 13, die Differenz zum Parent-Pfad ist **genau** die fünf genannten. Die
Zahlen des Commits (15/10/5) halten punktgenau.

**(2) Der Diskriminator ist exakt umgesetzt.** `forbiddenFileCredentialKeys` = 6 Schlüssel (drei
DSN, zwei Token, `nats_url`); `fileConfig` deklariert `http_addr`/`grpc_addr`, keines der sechs
verbotenen; die Fehlerzeile nennt Klasse und Schlüssel; die Gegenprobe (Tippfehler-Schlüssel) ist
im Test. `SPEC-016` und Handbuch §5.2 stimmen mit dem Code überein — nachgezählt.

**(3) Die Mutations-Bindungen halten (Eingabeseite).** Fünf eigene Mutationen, je Exit 1:
`NatsURL`-Zuweisung gestrichen · beide Token-Zuweisungen gestrichen · `api_token_reader` aus der
Liste gestrichen · `overrideString` → `getenv` bei beiden Adressfeldern (2 Präzedenz-Fälle rot).

**(4) Das abgedruckte Handbuch-Beispiel ist kein Dekor.** Provokations-Test im Scratch-Klon: der
YAML-Block aus §5.2 durch `ConfigFromFile` + `mergeConfig` → Exit 0,
`HTTPAddr=":8090" GRPCAddr=":9090" Source="quelle-1"`.

**(5) Träger-Zug und Versionshistorie.** Handbuch `Version: 1.18` + Zeile im selben Diff; die
neue Betreiber-Wirkung (Oberflächen unter Datei **an**) ist in §5.2 beschrieben.

**(6) Records bleiben richtig unangetastet.** `review-slice-041*`, `verify-slice-041` und
`done/slice-041-*` nennen weiterhin `forbiddenFileDSNKeys` — nach `ADR-0073` §Entscheidung 4
korrekt (Record-Inhalt unantastbar, keine §Geschichte-Zeile). **Kein** Befund.

**(7) Kein neuer Satz mit Chronik-Ton** in Produktionskommentaren (`grep` über den Diff: 0 Treffer).

## Zahlen nachgemessen (§3.12)

| Zahl | Ergebnis | Ursprung |
|---|---|---|
| 15 Felder / Parent 10 / jetzt 15 | **15** · **10** · **15**; `ConfigFromEnv` liest 13, Differenz = die fünf | gemessen — hält |
| acht Mutationen | 5 + 2 + 1 = **8**; **5** selbst gefahren, jede rot | gemessen (Stichprobe) — hält |
| Coverage 83,40 % (vorher 83,10) | eigener Lauf: HEAD **83.4 %** (28× ok), Parent **83.1 %** | gemessen — hält punktgenau |
| „die neuen Tests decken `config_file.go` mit" | `ConfigFromFile`, `ConfigFromEnvAndFile`, `overrideString`, `mergeConfig`, `mergeTables` je **100,0 %** | gemessen — hält |
| „kein Lauf setzt `CDC_CONFIG_FILE`" | `git grep`: Treffer nur in `docs/**` und `internal/bootstrap/**` | gemessen — hält |

*Benannte Abweichung (§3.12):* die Commandfolge der `coverage`-Stage wurde in einem `docker run`
derselben gepinnten Toolchain nachgefahren, **nicht** als `docker build --target coverage` — die
Zahl ist gemessen, nicht übernommen.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | **1** |
| MEDIUM | 0 |
| LOW | 2 |
| INFO | 2 |

## Verdikt

**Der Diff trägt in der Sache.** Der stille Defekt ist real behoben (15/15), der Diskriminator
aus `ADR-0088` ist **genau** in den drei Klassen umgesetzt, die Bindungen halten an der
**Eingabeseite**, das Handbuch-Beispiel lädt real durch den Loader, Versionshistorie und Träger
sind nachgezogen, kein neuer Satz trägt Chronik-Ton. Die Zahlen halten punktgenau.

**Fixrunde am Implementer: nein.** F-1 (HIGH) und F-2 (LOW) gehören dem **Architect** — F-1 als
Folge-ADR (`ADR-0073` nimmt Fitness-Function-Regeln von der Zitat-Korrektur aus), F-2 als
Zitat-Korrektur nach `ADR-0073`; beide gehen als benanntes Übergabe-Artefakt diesen Report-Weg,
**kein** Reviewer→Implementer-Pfeil. F-3 ist eine **Kommentar-Zeile** im Test (Nachzug, kein
Blocker). F-4/F-5 sind INFO.

**Vor der Closure:** F-1 sollte durch einen Architect-Zug aufgelöst werden — sonst schließt
dieser Slice gegen ein Beleg-Versprechen, das nicht existiert.

**DoD-Häkchen „Review durchgeführt":** nicht nachgezogen — dieser Lauf legt keine Datei an; wer
den Report committet, zieht das Häkchen im selben Commit nach. Es ist **kein**
Implementer-Fixrunde-Pfad offen, der es offen halten müsste.

**Nicht gefahren:** `make a-check`, `make baseline-verify`, `make image`,
`make test-store`/`-replication`/`-notify`/`-integration` — keine Architektur-Kante, keine
Baseline, kein Build-Kontext, keine Naht und kein Container-Vertrag berührt.
