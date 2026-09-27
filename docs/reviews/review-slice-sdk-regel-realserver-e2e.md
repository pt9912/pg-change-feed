# Review-Report: slice-sdk-regel-realserver-e2e — 2026-09-27

**Review-Art:** Code — Diff gegen Plan (`docs/plan/planning/in-progress/slice-sdk-regel-realserver-e2e.md`)
+ Konventionen (`AGENTS.md`, `harness/conventions.md`).

**Gegenstand:** `git diff 300637cf..47efceee` — Feature-Commit `ad3af754`
(„Row-Image-Umbenennung real an allen drei SDK-Tiers belegt“, [LH-FA-CFG-007](../../spec/lastenheft.md),
[ADR-0112](../plan/adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)) + DoD-Nachtrag-Commit `47efceee` (§3-Suchlauf-Feld).

**Skill:** `.harness/skills/reviewer.md` @ Accepted (Stand 2026-09-09,
zuletzt inhaltlich erweitert für `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht`)
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-27

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-sdk-regel-realserver-e2e.md` (Plan, §2 DoD, §3 Suchlauf, §6 Risiken)
- `ADR-0112` (Transformationsform, Folgepflicht 6/Teilfrage 8), `ADR-0110` §Entscheidung Festlegung 2, `ADR-0125`, `ADR-0126`
- `LH-FA-CFG-007`, `LH-FA-SST-009`, `LH-FA-SST-008`, `LH-FA-SST-006`
- `AGENTS.md` §3.1 (Docker-only), §3.7 (Kommentarform), §3.12 Instanz A/B (Herkunft von Aussagen), §3.13 (Suchlauf)
- `harness/sensors/kommentar-kennungen.md`, `harness/sensors/suchlauf-nachmessen.md`

---

## Findings

### F-1 — DoD-Checkboxen Liefer-Punkt 2/3 als erfüllt markiert, obwohl der literal geforderte Beleg-Umfang nur teilweise erbracht und die Lücke nirgends im Plan benannt ist

- `kategorie`: HIGH
- `quelle`: `AGENTS.md` §3.12 Instanz B („Ein Kriterium, das erst nach der
  Arbeit belegt werden kann, ist eine Zusage — es ist als solche formuliert
  und nicht als Feststellung“) · Reviewer-Skill HIGH „Beleg trägt seinen Satz
  nicht“
- `pfad`: `docs/plan/planning/in-progress/slice-sdk-regel-realserver-e2e.md:153-171`
  (Liefer-Punkt 2/3, „Zu belegen durch: … dieselben drei Mutationen der
  Eingabeseite, rot gesehen“)
- `befund`: Der DoD-Text für Liefer-Punkt 2 (Kotlin) und 3 (Python) verlangt
  wörtlich „dieselben drei Mutationen der Eingabeseite, rot gesehen“ (die drei
  Mutationen aus Liefer-Punkt 1: Regel-Antrag entfällt, Regel benennt falsche
  Spalte, Gegenlesen prüft den falschen Schlüssel). Laut Commit-Message
  `ad3af754` wurde für Kotlin und Python real nur **eine** der drei Mutationen
  unabhängig gefahren (der entfallende `set_transformation`-Aufruf); die
  beiden anderen sind ausschließlich strukturell begründet („Code-Identität
  der Hilfsdatei/Runner-Form, kein eigener Lauf“). Beide Checkboxen sind trotzdem
  auf `[x]` gesetzt — im Feature-Commit selbst, ohne Fixrunde. Die Lücke steht
  **nirgends im Plan-Dokument** (weder im DoD-Text, noch in §6 „Risiken“, noch
  im §3-Suchlauf-Nachtrag): `grep -n "Verifikationsspalt\|nicht unabhängig\|Code-Identität"
  docs/plan/planning/in-progress/slice-sdk-regel-realserver-e2e.md` liefert
  keinen Treffer. Nur die Commit-Message trägt die Einschränkung — die
  Plan-Datei selbst behauptet mit dem gesetzten Häkchen einen vollständig
  erbrachten Beleg, den die eigene "Zu belegen durch"-Formulierung nicht
  hergibt.
  **Eigene Nachprüfung:** Ich habe die im Plan benannte Lücke selbst
  geschlossen für einen Fall — Mutation (b) (Regel benennt `id` statt `name`)
  auf dem Python-Tier, gRPC-Fläche real gefahren (Kopie der gemeinsamen
  Hilfsdatei `tools/harness/lib-sdk-rule-fixture.sh` mutiert, `make
  test-sdk-python-integration` real ausgeführt, Container-seitig rot gesehen,
  Datei per `git checkout` zurückgesetzt). Ergebnis: der Test scheiterte
  korrekt — `AssertionError: assert '501' == 'PythonGrpcRuleSdkE2ESentinel'`
  (der `id`-Wert `501` landet unter `display_name` statt des Sentinels) —,
  d. h. **kein Funktionsfehler**, die Python-Assertion ist an dieser einen
  Fläche tatsächlich gebunden. Die übrigen drei Python-Flächen und alle vier
  Kotlin-Flächen für Mutation (b) sowie (c) an beiden Sprachen bleiben
  weiterhin nur strukturell begründet, nicht erprobt.
  Der Befund ist also ein **Prozess-/Dokumentationsfehler** (Checkbox behauptet
  mehr, als der Diff selbst belegt), kein bestätigter Funktionsfehler — die
  Einstufung bleibt trotzdem HIGH, weil die DoD-Zeile als Feststellung
  formuliert ist, obwohl sie zum Zeitpunkt des Commits eine Zusage war
  (`AGENTS.md` §3.12 Instanz B), und weil dieselbe Lücke bereits vom
  Implementer selbst erkannt, aber nicht in den Träger geschrieben wurde, der
  dafür vorgesehen ist (der Plan, nicht die Commit-Message).
- `verifizierbar`: ja — Plan-Text vs. Commit-Message-Text (gelesen),
  zusätzlich die eigene Mutation heute real gefahren (Log im Scratchpad
  dieser Review-Session, nicht committet).
- `klasse`: „Beleg trägt seinen Satz nicht — DoD-Checkbox ohne vollständige
  Zu-belegen-durch-Erfüllung“

**Empfehlung an den Implementer (keine Lösungsvorschlags-Pflicht, nur
Handlungsraum):** entweder die beiden Checkboxen wieder auf `[ ]` setzen und
den „Zu belegen durch“-Text um den tatsächlich erbrachten Umfang ergänzen
(„Mutation (a) unabhängig gefahren; (b)/(c) strukturell begründet, nicht
erprobt“ als benannte Grenze statt stiller Vollständigkeitsbehauptung), oder
die restlichen Mutationen real nachfahren, bevor der Haken gesetzt bleibt.

---

### F-2 — Kostenklasse-Messwerte aus §4 („Laufzeit je Tier“, `free -m`, dangling Volumes) fehlen im Bericht

- `kategorie`: LOW
- `quelle`: Plan §4 „Kostenklasse — schwere Realserver-Läufe“ (kein Gate,
  keine Hard Rule — Maintainability/Plan-Treue)
- `pfad`: `docs/plan/planning/in-progress/slice-sdk-regel-realserver-e2e.md:361-376`
- `befund`: §4 verlangt explizit, dass der Implementer die Laufzeit je Tier
  vor/nach der Erweiterung sowie `free -m` und die Zahl dangling Volumes vor
  dem ersten und nach dem letzten Lauf „im Bericht“ einträgt. Die
  Commit-Message `ad3af754` nennt keinen dieser Werte (nur Exit-Codes und
  Coverage-Prozent). Da dieser Wert nicht an eine Closure-Bedingung
  (`(bei Closure)`) delegiert ist, sondern gegenwärtig gefordert wird, fehlt
  er im vorliegenden Übergabe-Artefakt.
- `verifizierbar`: nein — hängt davon ab, ob eine Bericht-Instanz außerhalb der
  Commit-Message existiert, die dieser Review keinen Zugriff hat
  (`verifizierbar: nein`, wie Referenz-Fall F-3 des Reviewer-Skills).
- `klasse`: „Plan-Vorgabe zur Lauf-Beobachtung nicht getragen“

---

## Negativbefunde

- geprüft, ohne Befund: `tools/harness/lib-sdk-rule-fixture.sh` (Kommentarform
  §3.7, Docker-only §3.1, Wiederverwendung durch alle drei Runner identisch)
- geprüft, ohne Befund: `tools/harness/run-sdk-csharp-integration-tests.sh`,
  `run-sdk-kotlin-integration-tests.sh`, `run-sdk-python-integration-tests.sh`
  (Diff-Struktur der drei Runner ist deckungsgleich — kein
  `BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall`-Fund; keine
  Host-Toolchain-Aufrufe, kein `sed -i`/Umleitung auf Repo-Dateien)
- geprüft, ohne Befund: die zwölf neuen Testdateien/-klassen
  (`GrpcRuleRealserverTests`/`Test`/`.py`,
  `HttpRuleRealserverTests`/`Test`/`.py`,
  `NatsRuleRealserverTests`/`Test`/`.py`,
  `SseRuleRealserverTests`/`Test`/`.py`) — Assertion-Form (Zielschlüssel
  trägt Sentinel, Quellschlüssel fehlt) ist über alle drei Sprachen und alle
  vier Wege identisch aufgebaut; kein Chronik-Ton, kein „ff.“, Docstrings
  ausschließlich in Englisch ohne deutsches Wortfragment
- geprüft, ohne Befund: `PhaseEnvironment.cs`/`.kt` (Ergänzung
  `RuleSourceKey`/`RuleTargetKey` bzw. `ruleSourceKey`/`ruleTargetKey`,
  identisches Muster zu den bestehenden Feldern)
- geprüft, ohne Befund: `make sdk-public-doc-check` real ausgeführt gegen den
  Arbeitsbaum — Exit 0, keine interne Kennung (`SPEC-`, `ADR-`, `LH-FA-`/`LH-QA-`,
  `slice-`/`welle-`) unter `sdks/`
- geprüft, ohne Befund: `make suchlauf-nachmessen
  PLAN=docs/plan/planning/in-progress/slice-sdk-regel-realserver-e2e.md` real
  ausgeführt — 22/22 Zeilen stimmen (Exit 0); die behauptete Zahl „22/22“ ist
  damit selbst nachgemessen, nicht nur übernommen
- geprüft, ohne Befund: `git diff --name-only 71024045 47efceee -- sdks` —
  ausschließlich die drei Test-Verzeichnisse, keine `.csproj`/`pyproject.toml`/
  `build.gradle.kts`, keine Versionshebung
- geprüft, ohne Befund: `make kommentar-kennungen DIFF=300637cf` — 0
  Kandidaten (der Diff berührt keine Go-Datei, die Regel greift nicht)
- geprüft, ohne Befund: DoD-Rollentrennung — offene Checkboxen (Review,
  Closure-Notiz, Beobachtungs-Register, Risiken-Ausgänge, drei Paarungen)
  sind korrekt `[ ]` belassen; einzige Auffälligkeit ist F-1 oben
  (Liefer-Punkt 2/3 fälschlich `[x]`)
- geprüft, ohne Befund: `harness/README.md` §Sensors, `harness/mk/sdk.mk` —
  Nachzug der drei Sensor-Zeilen und der beiden Hilfetexte („acht Phasen“)
  deckungsgleich mit dem tatsächlichen Umfang; die Nichtänderung der
  Python-Hilfszeile ist konsistent (sie trug „vier Phasen“ nie)
- geprüft, ohne Befund: `docs/user/sdk-e2e-abdeckung.md` — genau drei neue
  Zeilen (eine je Sprache), marker-gegrenzte Abschnitte unverändert sonst
- geprüft, ohne Befund: Kombination von Slice-Name und ADR-Anker im selben
  Bash-Kommentarblock (`lib-sdk-rule-fixture.sh:3`,
  `run-sdk-*-integration-tests.sh` Kopfkommentare) — entspricht bereits
  bestehender Repo-Konvention (`tools/harness/sdk-pack-csharp.sh`,
  `run-notify-tests.sh`, `release-tag-info.sh`, `sdk-pack-kotlin.sh` tragen
  dasselbe Muster); die strikte „höchstens eine Kennung“-Regel und das
  Werkzeug `make kommentar-kennungen` sind laut eigenem Vertrag auf
  Go-Kommentare skopiert — kein Befund dieser Klasse

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** Beleg trägt seinen Satz nicht — DoD-Checkbox
ohne vollständige Zu-belegen-durch-Erfüllung · Plan-Vorgabe zur
Lauf-Beobachtung nicht getragen

## Verdikt

**Merge-blockierend:** ja — F-1 ist HIGH. Die Blockade ist eine
Dokumentations-/Prozessfrage, kein bestätigter Funktionsfehler: die eigene
Mutationsprobe dieser Review (Mutation (b), Python-Tier, gRPC-Fläche) fand
keinen Fehler, sie bestätigte im Gegenteil, dass die Assertion an dieser
einen Stelle korrekt gebunden ist. Der Implementer kann F-1 ohne neuen
Realserver-Lauf schließen — Checkbox zurücknehmen und die tatsächlich
erbrachte Beleg-Menge im Plan benennen (Alternative: die verbleibenden
Mutationen real nachfahren). F-2 (LOW) blockiert nicht.

**Übergabe:** F-1 geht an den Implementer zur Fixrunde (DoD-Text/Checkbox
korrigieren oder Mutationen nachfahren); F-2 kann im selben Zug mitgezogen
werden. Da eine Fixrunde erwartet wird, bleibt die Plan-Checkbox „Review
durchgeführt“ bewusst `[ ]` (kein Nachzug ohne Fixrunde, siehe
`.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug ohne Fixrunde — dieser
Fall greift hier nicht). Der Report ersetzt keine Verifikation — DoD-/
Spec-Konformität prüft der Verifier separat.
