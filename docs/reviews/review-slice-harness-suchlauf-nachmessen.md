# Review-Report: slice-harness-suchlauf-nachmessen — 2026-09-26

**Review-Art:** Code — der Diff liefert das Nachmess-Werkzeug `make suchlauf-nachmessen`
(Skript, Vertrag, Test, Make-Ziele), den Wrapper `tools/schema/rollout-restore.sh` samt Test und
die Umstellung aller sieben Aufrufer von `make schema-rollout`, die Träger-Zeilen (Harness-README,
`AGENTS.md` §3.13, `implement-slice` Schritt 18, Reviewer-Skill, Schema-Rollout-Vertrag) und ersetzt
vier Register-Links durch Kennungs-Zitate; geprüft gegen Plan, ADRs und `AGENTS.md` Hard Rules
(Modul 10 §Drei Review-Arten). Kein DoD-Abgleich — das ist Verifier-Aufgabe (Modul 11).

**Gegenstand:** Slice `slice-harness-suchlauf-nachmessen` (ohne Welle), Diff-Range `71ff3e44..fc0f5629`
(6 Commits, 23 Dateien, +608/−30; Baum sauber, gepusht).

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09“ (seither um weitere
HIGH-Klassen ergänzt, u. a. Zahl-im-Träger mit `suchlauf`-Probe, Beleg-Satz,
Zusage-ohne-Eingabeseite, Nachzug-Nachbar).
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-26.

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis)*. Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link
> (`v<X.Y.Z>` · `regelwerk/<datei>.md` §<Abschnitt>). Der vendored Baum trägt
> genau einen Tag; der Sprung löscht den alten, und ein Link darauf färbt beim
> nächsten Bump ein Artefakt rot, das niemand mehr anfassen darf. Ein `pfad`-Feld
> auf den **geprüften Gegenstand** ist davon nicht betroffen — es zitiert den
> Stand des Laufs und darf ihn festhalten (`v<X.Y.Z>` ·
> `regelwerk/grundlagen-harness-dateien.md` §harness/README.md als
> Einstiegspunkt — diese Zeile ist selbst ein Beispiel der Form).

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde — ohne
diese Liste ist der Lauf nicht reproduzierbar):

- Slice-Plan `slice-harness-suchlauf-nachmessen` (§1 Ziel und Abgrenzung, §2 DoD-Wortlaut als
  Bezug, §3 Plan mit Abweichungs-Zeilen und Suchlauf-Feld, §6 Risiken); Architect-Verdikt
  `architect-verdict-welle-backfill-bestand-lese-schritt` (Folge-Slice A)
- [`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md) (Herkunft von Aussagen, Grenze
  „kein Sensor auf Prosa“), [`ADR-0043`](../plan/adr/0043-schemamigrationen-mit-d-migrate.md)
  (Schema-Rollout), [`ADR-0071`](../plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
  (Messgegenstand der netzlos prüfbaren Fläche; Makefile-Konventionen)
- `AGENTS.md` (Hard Rules §3.1 Docker-only, §3.3, §3.7, §3.9, §3.12, §3.13), `harness/conventions.md`
  (`MR-000`/`MR-001`)
- Report-Gerüst: `docs/reviews/review-report.template.md`, Formvorbild
  `docs/reviews/review-slice-sdk-readme-nutzerdoku.md`

**Eigenständig durchgeführte Prüfungen** (gemessen, nicht aus dem Implementer-Bericht
übernommen; Exit-Codes ungepiped in Log-Dateien gesichert, gedruckte Zeilen zitiert):

- **Tests am Stand `fc0f5629`:** `make test-suchlauf-nachmessen` Exit 0, gedruckt
  „run-suchlauf-nachmessen-tests: alle Fälle bestanden“; `make test-rollout-restore` Exit 0, gedruckt
  „run-rollout-restore-tests: alle Fälle bestanden“; `make docs-check` Exit 0, gedruckt
  „d-check: 1198 Datei(en) geprüft, 0 Befund(e)“; `make gates` (einmal) Exit 0, gedruckt
  „gesamt: 0 Befund(e)“, `git status --short` danach leer.
- **Werkzeug am Plan selbst:** `make suchlauf-nachmessen PLAN=<Plan-Datei>` Exit 0, acht Zeilen
  `OK`, gedruckt „suchlauf-nachmessen: 8 Zeilen stimmen“ (soll=ist: 14/74/9/12/59/74/35/45).
  **Von Hand** mit `git grep` nachgefahren, ohne das Werkzeug: Zeile 1 (Parent `8717c4fb`, `-i -E
  'suchlauf|nachmess|Suchform'` über die Träger-Pfade) 14 Trefferzeilen (AGENTS.md 9,
  Benutzerhandbuch 3, Reviewer-Skill 2); Zeile 2 (Diff, gleicher Befehl, Plan-Datei per
  `:(exclude,glob)` ausgenommen) 74, je Datei 13/5/3/3/10/4/7/11/18; Zeile 4 (Diff, Formmuster,
  ganzer Baum) 12 in AGENTS.md 2, `implement-slice.md` 1, Sensor-Doku 2, Register 7; Zeile 6 (Diff,
  `plan\.yaml|down\.sql`) 74; Zeile 7 (Parent, ohne Plan-Datei) 35; Zeile 8 (Diff) 45 in 17 Dateien.
  Auch die Datei-Aufteilungen der Befund-Zellen (`AGENTS.md` 9 zu 13, Skill 2 zu 4, „15 zu 17
  Dateien“) stimmen. „Nicht gefunden: weitere Beschreibung in `.claude/commands/*`“ am Parent:
  0 Zeilen mit `suchlauf` in `.claude` und `harness/sensors`.
- **Eingabeseiten-Mutationen am Werkzeug** (je eine Kopie des Skripts per `sed`, Test mit `TOOL=<Kopie>`
  gefahren, Exit ungepiped; das Original blieb unverändert):

  | Nr. | Mutation an `suchlauf-nachmessen.sh` | Ergebnis |
  |---|---|---|
  | M1 | Ausschluss der Plan-Datei entfernt | Exit 1, Fall „stimmt“ rot (Ist 6 statt 3) |
  | M2 | Ausschluss als Pfad statt als Dateiname-Glob | Exit 1, rot |
  | M3 | Zahlenvergleich umgekehrt (`-ne`) | Exit 1, rot |
  | M4 | `HEAD` als Stand erlaubt | Exit 1, Fall „HEAD abgelehnt“ rot |
  | M5 | leerer Block gilt als bestanden | Exit 1, Fall „kein Block (leer)“ rot |
  | M6 | `git grep`-Exit über 1 ignoriert | Exit 0, „alle Fälle bestanden“ — nicht rot |
  | M7 | nicht geschlossener Block ignoriert | Exit 0 — nicht rot |
  | M8 | Soll ohne Zahl ohne Prüfung | Exit 0 — nicht rot |
  | M9 | Anführungszeichen als Gruppierung entfernt | Exit 1, Fall „stimmt“ rot |
  | M10 | `--`-Trennung von Optionen und Pathspec entfernt | Exit 0 — nicht rot; am echten Plan-Feld dieselbe Kopie: `ABWEICHUNG  soll=14 ist=74` |
  | M11 | Stand-Kennung ohne Existenzprüfung | Exit 0 — nicht rot |
  | M12 | Zeile ohne Suchmuster zugelassen | Exit 0 — nicht rot |

- **Eingabeseiten-Mutationen am Wrapper** (Kopien, `TOOL=<Kopie>`): Sicherung ohne `trap`, `.fehlt`-Marke
  ignoriert, Sicherung entfernt, Exit-Code verschluckt, Wiederherstellung per `git checkout` statt aus
  der Sicherung, nur `plan.yaml` gesichert, Exit 0 ohne Kommando: **sieben von sieben** färben
  `make test-rollout-restore` rot (Exit 1). Aufrufer-Prüfung: `tools/bench-lib.sh` auf den alten Aufruf
  zurückgesetzt → „FEHLER: Aufrufer ohne rollout-restore.sh: tools/bench-lib.sh:76“, Exit 1; danach
  `git checkout`.
- **Injektion über Plan-Inhalt:** eine Plan-Zeile mit `-O'touch <Datei>;true'` bzw.
  `--open-files-in-pager='touch <Datei>;true'` als `git grep`-Argument legte die Datei an, das Werkzeug
  druckte `OK  soll=0 ist=0` und endete mit Exit 0 (siehe F-1). Weitere Eingaben: Soll `007` (dezimal
  gelesen), `:(top)`-Pathspec (durchgereicht), `--no-index` auf einen Pfad außerhalb des Repos
  (Exit 2, `git grep` Exit 128), `--and` ohne Muster (Exit 2), unbekannte Kennung `abcdef0` (Exit 2),
  `git` nicht im `PATH` (Exit 2 mit Shell-Meldung), `make suchlauf-nachmessen` ohne `PLAN` (Exit 2,
  `$(error …)` vor jedem Befehl), Abweichung über `make` (Skript Exit 1, Make-Exit 2).
- **Aufrufer von `make schema-rollout` von Hand:** `git grep -n 'schema-rollout' -- ':!docs'
  ':!*.md'` gelesen: sieben Aufruf-Zeilen (`apply-rollout.sh`, `bench-lib.sh`, `run-integration-tests.sh`,
  drei `run-sdk-*-integration-tests.sh`, `examples/bootstrap.sh`), alle über `rollout-restore.sh`;
  der Guard-Test trägt seine eigene Sicherung; die übrigen Treffer sind Kommentare, Fehlertexte,
  `Makefile`-Kopfzeile und `compose.yaml`-Kommentar. Kein Aufrufer außerhalb des Wrappers.
- **Reale Wirkung:** `make test-store` einmal, Exit 0, gedruckt „db-coverage: OK — DB-Adapter-Coverage
  82.61% erfuellt Schwelle 80%“; danach `git status --short tools/schema` und `git status --short`
  leer. Mutation: der Wrapper aus `apply-rollout.sh` entfernt, derselbe Lauf Exit 0, danach
  `M tools/schema/plan.yaml` (Arbeitsbaum verändert); beide Änderungen per `git checkout`
  zurückgenommen, Baum sauber. `make test-replication` habe ich nicht gefahren (nicht beauftragt).
- **Umgebung:** `free -m` vor den Läufen 17,5 GB verfügbar, danach 17,7 GB; ein schwerer Docker-Lauf
  zugleich; dangling Volumes (`docker volume ls -q -f dangling=true | wc -l`) 34 vor und 34 nach allen
  Läufen; kein `prune`; keine eigenen Container oder Images zurückgeblieben (ein fremder Container
  `bats/bats` lief unabhängig von diesem Review).
- **Commits:** alle sechs Betreffs nennen `ADR-0083` (Bau-Commit zusätzlich `ADR-0043`), keiner
  trägt `SPEC-*`/`ARC-*`; die beiden Lifecycle-Moves (`42e5ebb8`, `8717c4fb`) sind reine Renames
  (0 Zeilen), der Inhalt steht in eigenen Commits (`AGENTS.md` §3.3).
- **Nicht gefahren (Grenze):** `make test-replication`, `make test-integration`,
  `make test-sdk-*-integration`, `make bench`, `make example-demo-up` — die sechs direkten Aufrufer
  außerhalb von `apply-rollout.sh` sind gelesen und über die Aufrufer-Prüfung gebunden, nicht real
  gefahren.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | HIGH | Der Zerleger reicht jedes Wort einer `suchlauf`-Zeile unverändert an `git grep` weiter; `-O<Kommando>` und `--open-files-in-pager=<Kommando>` lassen `git grep` das Kommando über die Shell starten. Gemessen: eine Plan-Zeile `diff 0 -O'touch <Datei>;true' -e … -- …` legt die Datei an, das Werkzeug meldet `OK` und Exit 0. Der Plan sagt in §3 „eigener Zerleger, kein `eval`: eine Plan-Zeile führt keinen Shell-Code aus“, der Vertrag nennt keine Optionsgrenze; der Reviewer und der Verifier führen das Werkzeug auf Plan-Inhalt aus, der im Diff steht. | Sicherheits-Anti-Pattern (Reviewer-Skill, HIGH-Liste); Plan §3 „Form der Zeile“ | `tools/harness/suchlauf-nachmessen.sh:106-115` (Optionen ungeprüft), `:137`/`:139` (Aufruf) | ja — Plan-Zeile mit `-O` in `make suchlauf-nachmessen`, Datei entsteht; der Tabellentest hat keinen Fall dafür | Argument-Injektion über Werkzeug-Option; Beleg trägt seinen Satz nicht („führt keinen Shell-Code aus“) |
| F-2 | MEDIUM | Der Tabellentest bindet sechs Zusagen des Vertrags nicht an ihre Eingabeseite: `git grep`-Fehler (Exit über 1 → Exit 2), nicht geschlossener Block, Soll ohne Zahl, unbekannte Stand-Kennung, Zeile ohne Suchmuster (M6, M7, M8, M11, M12 grün) und die Kombination Commit-Stand mit `--`-Pathspec (M10 grün; dieselbe Kopie zählt am echten Plan-Feld 74 statt 14). Die Fälle mit Pathspec fahren nur den Stand `diff`. | Maintainability; fehlende Negativtests bei neuem öffentlichem Vertrag (Reviewer-Skill, MEDIUM-Liste) | `tools/harness/run-suchlauf-nachmessen-tests.sh:66-73` (Fall 1) und Fälle 1–5 gesamt; Vertrag `harness/sensors/suchlauf-nachmessen.md` §Ausgabe und Ausgänge | ja — Mutationen M6–M8, M10–M12 an einer Kopie, `make test-suchlauf-nachmessen` bleibt Exit 0 | Vertrags-Exit-Zweige ohne Testbindung |
| F-3 | MEDIUM | Der Plan-Wortlaut und seine Abweichungs-Zeile stehen nebeneinander, ohne dass einer auf den anderen verweist: §1 („`tools/schema/apply-rollout.sh` stellt … wieder her“) und der bereits abgehakte dritte DoD-Punkt („`apply-rollout.sh` sichert … vor `make schema-rollout` in ein `mktemp`-Verzeichnis“) nennen das Skript als Ort der Rücknahme; §3 sagt, die Rücknahme liege im Wrapper, „nicht im Skript selbst“, und begründet das damit, der Plan lasse „`apply-rollout.sh` oder das Target offen“ — §1 und DoD lassen es nicht offen. | Maintainability; Nachzug widerspricht dem Nachbarn im selben Träger (Reviewer-Skill, MEDIUM-Liste) | `docs/plan/planning/in-progress/slice-harness-suchlauf-nachmessen.md` §1 (Ziel), §2 dritter Punkt, §3 Zeile `tools/schema/rollout-restore.sh` | ja — Lesen der drei Stellen | Nachzug-Nachbar im Plan; Begründung nennt eine Plan-Aussage, die der Plan nicht trägt |
| F-4 | LOW | Ein leeres Argument (`''`) in einer Plan-Zeile verschwindet: das Array wird über `read -a` mit `US`-getrenntem Text zurückgelesen, ein abschließendes leeres Feld fällt weg. Gemessen: `diff 1 -e '' -- AGENTS.md` misst mit verändertem Argumentsatz (`ist=2`) statt die Zeile abzulehnen. Daneben zählt `2>&1` Stderr-Zeilen eines `git grep`-Laufs mit Exit 0 oder 1 als Trefferzeilen. | Maintainability | `tools/harness/suchlauf-nachmessen.sh:122-123`, `:132-134`, `:137-147` | ja — Zeile mit leerem Argument | Zerleger verliert leeres Argument |
| F-5 | INFO | Der Wrapper sichert beim Start und stellt beim Ende her; zwei gleichzeitig laufende Wrapper-Aufrufe (zwei DB-Tier-Läufe parallel) können den Zwischenstand des anderen als „Zustand vor dem Lauf“ sichern und nach dem Ende wiederherstellen. Kein Ziel und kein Gate startet zwei Rollouts zugleich (beide schreiben ohnehin dieselben Dateien); die Grenze steht nicht im Skriptkopf. Ein SIGKILL des Wrappers lässt die Erzeugnisse verändert (kein `trap` fängt ihn). | Maintainability | `tools/schema/rollout-restore.sh:19-47` | nein | Wrapper nicht nebenläufig |
| F-6 | INFO | Vereinbarkeit mit `AGENTS.md` §3.1: das Werkzeug braucht auf dem Host `bash`, `git` und `realpath`; der Vertrag nennt „nur `bash` und `git`“. Präzedenz gleicher Form: `apply-rollout.sh`, `run-release-tag-info-tests.sh` (`git rev-parse` auf dem Host); `git` ist Voraussetzung jeder Arbeit an diesem Repo, es wird nichts installiert. Das §6-Risiko des Plans führt den Ausgang noch als offen; die Vereinbarkeit ist eine Architect-Aussage, kein Befund. | Docker-only (`AGENTS.md` §3.1) | `tools/harness/suchlauf-nachmessen.sh:24-25`; `harness/sensors/suchlauf-nachmessen.md` §Vertrag | nein | Host-Werkzeug-Grenze von §3.1 |
| F-7 | INFO | Die Aufrufer-Prüfung von `make test-rollout-restore` liest eine Zeile mit `make` und `schema-rollout` unter `tools/*.sh` und `examples/*.sh`. Formen außerhalb (Fortsetzungszeile mit `schema-rollout` auf der Folgezeile, `$MAKE`, ein Makefile-Rezept, ein Compose-Kommando) sieht sie nicht; von Hand nachgefahren gibt es solche Aufrufer heute nicht (siehe Prüfungen). Die Zeilen-Bindung prüft nur, dass `rollout-restore.sh` in derselben Zeile steht. | Maintainability | `tools/harness/run-rollout-restore-tests.sh:19-31` | ja — Mutation eines Aufrufers auf `make` gemessen (rot); die genannten Formen nicht | Aufrufer-Prüfung auf Einzelzeilen-Form begrenzt |
| F-8 | INFO | Vier Register-`state.md` tragen den Slice jetzt als Kennung in Inline-Code statt als Link (der Link auf `open/` wäre mit dem Lifecycle-Move gebrochen); der Zustandstext dort führt den Slice weiter als „geplant“, obwohl er in `in-progress/` liegt. Der Plan meldet das als Closure-Arbeit des Planners (Träger in fremder Datei, §3.13). | `AGENTS.md` §3.13 (Träger-Nachzug) | `docs/plan/planning/observations/BEO-PGC/*/state.md` (vier Einträge) | nein | Träger-Nachzug fremder Datei (gemeldet) |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `tools/harness/suchlauf-nachmessen.sh` (Zerleger, Anführungszeichen, Leerzeichen in Mustern, `--`-Trennung, `:(…)`-Pathspec, Zahlenvergleich, `HEAD`/Namen abgelehnt, Exit-Codes, Fehlerpfade) | geprüft, ohne Befund über F-1, F-2, F-4 hinaus: kein `eval`, jedes Wort als eigenes Argument; Befehl immer `git grep -n …`; Quoting mit Leerzeichen (M9 rot), Soll `007` dezimal, Stand `diff` und Commit-Kennung getrennt gemessen, `HEAD`/`HEAD~1` vor jedem Lauf abgelehnt (M4 rot), Exit 0/1/2 nach Vertrag, Exit über `make` als Make-Exit 2 wie dokumentiert, `git` nicht im `PATH` und außerhalb eines Repos Exit 2; Selbstverweis-Ausschluss über Dateinamen zieht dieselbe Datei in anderen Verzeichnissen mit ab — im Vertrag genannt („in jedem Verzeichnis“) |
| `tools/harness/run-suchlauf-nachmessen-tests.sh` (fünf Fälle, Wegwerf-Repo) | geprüft, ohne Befund über F-2 hinaus: die fünf benannten Fälle sind an ihre Eingabeseite gebunden (M1–M5, M9 rot); das Wegwerf-Repo setzt `user.*` und `commit.gpgsign` selbst, der Aufräum-`trap` läuft, `TOOL` übersteuerbar; keine Abhängigkeit vom Host-Repo-Stand |
| `tools/schema/rollout-restore.sh` und `tools/harness/run-rollout-restore-tests.sh` | geprüft, ohne Befund über F-5, F-7 hinaus: Sicherung, Wiederherstellung, lokale Änderung bleibt, fehlende Datei bleibt fehlend, Exit-Code des Kommandos, `trap` auf EXIT/INT/TERM; sieben von sieben Wrapper-Mutationen rot, Aufrufer-Mutation rot; die Wahl „Wrapper statt Rücknahme im Skript“ ist im Plan §3 begründet (sechs weitere Aufrufer) und trägt: ein Ort für sieben Aufrufer, der Betrieb ruft `make schema-rollout` weiter direkt und behält Report und Rollback-Artefakt |
| Die sieben Aufrufer (`apply-rollout.sh`, `bench-lib.sh`, `run-integration-tests.sh`, drei `run-sdk-*-integration-tests.sh`, `examples/bootstrap.sh`) | geprüft, ohne Befund: je eine Zeile auf `bash tools/schema/rollout-restore.sh make …` umgestellt, Argumente unverändert; `bench-lib.sh` ruft den Wrapper über den absoluten Repo-Pfad, alle übrigen laufen bereits im Repo-Wurzel-Verzeichnis; Kommentare in `apply-rollout.sh` und `run-integration-tests.sh` beschreiben den Ist-Zustand; reale Wirkung nur für `make test-store` gemessen (sauberer Baum, Mutation rot) |
| `Makefile` (drei Ziele) | geprüft, ohne Befund: `.PHONY`, `$(error …)` bei fehlendem `PLAN` vor jedem Befehl, Zeilen in `make help`, keines der Ziele in einem Gate-Bündel |
| `harness/sensors/suchlauf-nachmessen.md`, `harness/README.md`-Zeilen, `AGENTS.md` §3.13-Absatz, `implement-slice` Schritt 18, Reviewer-Skill | geprüft, ohne Befund über F-6 hinaus: Ist-Zustand in Indikativ, kein Gate behauptet, Grenze („Zahlen und Stände, nicht Vollständigkeit von Suchraum und Muster“) ehrlich und an allen fünf Stellen gleich; der §3.13-Absatz nennt Werkzeug und Aufruf, ohne den Regelcharakter zu ändern (Regel bleibt der Absatz „Suchform“, der Suchraum-Ausschluss von `docs/reviews/**`, `done/` und Baseline steht weiter in der Zeile, das Werkzeug erzwingt ihn nicht — im Vertrag als Grenze genannt); der Satz „Bestehende Pläne tragen keinen solchen Block“ stimmt (`git grep -l '^ *```suchlauf'` nennt nur den eigenen Plan) |
| `harness/targets/schema-rollout.md` (Abschnitt „Erzeugnisse in Test-, Bench- und Beispiel-Läufen“) | geprüft, ohne Befund: die genannten sechs Läufe entsprechen den Aufrufern; Betrieb und Guard-Test korrekt abgegrenzt |
| Suchlauf-Feld des Plans §3 (acht Zeilen, beide Stände) | geprüft, ohne Befund: mit dem Werkzeug Exit 0 und von Hand nachgefahren (siehe Prüfungen), Befund-Zellen (Datei-Aufteilungen, „Nicht gefunden“-Satz) stimmen; Nichtgefundenes je Träger steht im Feld |
| Vier Register-`state.md` (Link durch Kennung ersetzt) | geprüft, ohne Befund über F-8 hinaus: nur die Verweis-Form geändert, `make docs-check` Exit 0 |
| Kommentare nach `AGENTS.md` §3.7 in allen neuen und geänderten Skripten | geprüft, ohne Befund: Kopf- und Zeilenkommentare tragen Zusage, Kopplung oder Abgrenzung im Indikativ; keine verworfene Alternative, kein Vorher/Nachher, keine Slice-Chronik in Produktionspfaden; der Tabellentest nennt seine Fälle als Test-Subjekt |
| Hard Rules und Commit-Struktur | geprüft, ohne Befund: keine Suppression, kein host-lokaler absoluter Pfad, keine Toolchain-Installation, `git mv` rein und Inhalt getrennt, Betreffs mit `ADR-0083`, ohne `SPEC-`/`ARC-`; Exit-Disziplin `AGENTS.md` §3.9 im Wrapper (`"$@"` als letzter Befehl, Exit-Code des Kommandos) |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 2 |
| LOW | 1 |
| INFO | 4 |

**Finding-Klassen dieses Laufs:** Argument-Injektion über Werkzeug-Option · Vertrags-Exit-Zweige ohne Testbindung · Nachzug-Nachbar im Plan · Zerleger verliert leeres Argument · Wrapper nicht nebenläufig · Host-Werkzeug-Grenze von §3.1 · Aufrufer-Prüfung auf Einzelzeilen-Form begrenzt · Träger-Nachzug fremder Datei (gemeldet)

## Verdikt

**Merge-blockierend:** ja — F-1 (HIGH) ist eine Ausführung beliebiger Kommandos aus Plan-Inhalt; F-2 und F-3
(MEDIUM) gehören in dieselbe Fixrunde. Die Wrapper-Seite (Rücknahme von `plan.yaml`/`down.sql`, Umstellung der
sieben Aufrufer, Träger-Zeilen, Register-Kennungen) trägt ohne Befund über INFO hinaus.

**Übergabe:** Findings gehen an den Implementer (Fixrunde nötig, deshalb bleibt die DoD-Zeile „Review
durchgeführt“ im Plan offen und wird bei Schritt 21 des Implementer-Workflows nachgezogen); die
**Finding-Klassen** gehen zusätzlich in die Slice-Closure §7 und von dort in den Zähler. F-8 geht als
Meldung an den Planner (Frist: Closure dieses Slice). Dieser Report selbst ist ein **Lauf-Beleg**
(Audit: dieser Diff, dieser Skill, dieses Modell, dieses Verdikt) — er wird über Läufe hinweg nicht wieder
gelesen, und muss es nicht. Der Report ersetzt keine Verifikation — DoD-/Spec-Konformität prüft der
Verifier separat (Modul 11; anderes Prüf-Artefakt, anderer Eingabe-Kontext).
