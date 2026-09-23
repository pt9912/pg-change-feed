# Welle welle-sdk-reale2e — Closure-Notiz

**Welle:** welle-sdk-reale2e
**Abschluss:** 2026-09-23
**Verantwortlich:** pt9912

## Was wurde geliefert?

Die Zielmatrix dieser Welle ist geschlossen: **alle zwölf Zustellweg-Flächen
der drei SDK-Packages (3 Sprachen × 4 Wege) tragen je einen realen
Realserver-Beleg** gegen eine laufende Server-Instanz — eingelöst über
[`ADR-0110`](../../adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
§Entscheidung Festlegung 2/Folgepflicht 1 (die etablierte Mechanik-Klasse)
in drei strikt sequentiellen Slices:

- **`slice-sdk-csharp-reale2e`**: die vier C#-Realserver-Phasen (gRPC, SSE,
  NATS-Vollinhalt, HTTP) mit je SQL-Gegenprüfung und je Ablehnungs-Beleg;
  führt den C#-Runner, das Make-Target `test-sdk-csharp-integration`, den
  Abdeckungs-Träger [`sdk-e2e-abdeckung.md`](../../../../docs/user/sdk-e2e-abdeckung.md)
  (C#-Abschnitt, Erzeugnis des Runners) und den `trace.coverage`-Eintrag
  (Label `SDK-E2E`) in `.d-check.yml` ein.
- **`slice-sdk-kotlin-reale2e`**: dieselbe Mechanik, auf den Kotlin-Baum
  gespiegelt — additive `integration`-Stufe (kein neuer Pin, Wiederverwendung
  aus `build`), Runner `tools/harness/run-sdk-kotlin-integration-tests.sh`,
  Make-Target `test-sdk-kotlin-integration`, Kotlin-Abschnitt im Träger.
- **`slice-sdk-python-http-reale2e`**: die einzige bislang nur netzlos
  belegte Fläche (Python-HTTP) nachgezogen — vierte Phase im bestehenden
  Python-Runner (Parametrisierung `received_grep`/`sql_kind`, SQL-Variante
  `consumer` gegen `cdc.consumer`), Python-Abschnitt im Träger, Träger-
  Erweiterung und Matrix-Vollständigkeits-Prüfung als eigene DoD-Punkte.

Der Abdeckungs-Träger
[`docs/user/sdk-e2e-abdeckung.md`](../../../../docs/user/sdk-e2e-abdeckung.md)
deklariert real die volle Matrix: **12 Zeilen** (4 C# + 4 Kotlin + 4 Python;
vom Verifier des Schluss-Slices selbst gezählt, Verifikation §2 Zeile 3),
alle drei Abschnitte marker-gegrenzt und byte-identisch dem Erzeugnis ihres
je eigenen Runner-Generators (idempotent geschrieben, kein Lauf-Beleg — die
Läufe trägt die Verifikations-Reports); der `trace.coverage`-Eintrag
(Label `SDK-E2E`) ist real im `.d-check.yml`-Bestand (Zeilen 318–319) und
macht [`LH-FA-SST-009`](../../../../spec/lastenheft.md) über `make doc-trace`
sichtbar (79 Anforderungen, 2 Waisen — je Verifikations-Report §1,
`LH-FA-SST-009` über `SDK-E2E`, Status `ok`).

`make gates` grün auf dem Endstand (938 Datei(en) im Gate-Lauf, 0 `docs-check`-Befunde,
Coverage 82,60–82,80 % über die Gate-Läufe dieser Closure (run-to-run-Variation der Messung, jeder Lauf über der 80 %-Endstufe), `a-check` 0 Befunde,
`generated-sync`/`commit-traceability`/`baseline-verify` je ohne Befund),
ungepiped geprüft zu jedem Commit dieser Closure — der Ausgang je
Commit-Klasse steht in §Verifikation ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
Kein Tag-Push (`git tag -l "sdk-*-v*"` liefert ausschließlich
`sdk-csharp-v0.1.0` und `sdk-python-v0.1.0` — eigenständig gemessen zu
dieser Closure): diese Welle ändert keine öffentliche API-Fläche und kein
Package-Artefakt, wie im Closure-Trigger (§3 der Welle-Datei) vorab
festgelegt.

## Was hat funktioniert?

- Die strikte Sequenz trug sich: der Abdeckungs-Träger entstand **mit dem
  ersten Slice** als Erzeugnis des C#-Runners (kein retrospektives
  Nach-Generieren, Welle-Plan §4 Reihenfolge-Grund 1), die Kotlin-Mechanik
  spiegelte die durchlaufene Form (dieselbe `integration`-Stufen-/Runner-
  Form, keine zweite Infrastruktur), und die Vollständigkeits-Prüfung lag
  beim Schluss-Slice, der den Träger um seinen Abschnitt erweiterte.
- Drei Runner schreiben je einen marker-gegrenzten Abschnitt derselben
  Träger-Datei, ohne sich zu überbieten: der beidseitige Rest-Erhalt
  (C#-Review F-7-Kette) hielt über alle drei Slices — die Abschnitte sind
  byte-identisch dem Erzeugnis ihres je eigenen Generators, die Nachbar-
  Abschnitte in jedem Diff rein additions (je Verifikation §2).
- Die Review-Ketten fingen Träger-/Plan-Lücken, bevor sie in `done/`
  landeten: drei Fixrunden (`86892bdb`, `c6523009`, `d668b9cc`), jede vom
  Verifier am Endstand gegen beide Stände bestätigt; die Suchlauf-Felder
  wurden je vom Verifikations-Lauf am Baum nachgemessen (alle Feld-Zeilen
  bestätigt, Verifikation §3 je Slice).
- Das Move-/Reconciliation-Muster ([`AGENTS.md`](../../../../AGENTS.md) §3.3)
  hielt über die drei Slice-Moves und den Welle-Move sauber; die
  Kennungs-Form in Inline-Code (Muster `a7723fa7`) hält die Berichte
  stabil gegen weitere Moves.

## Was ging anders als geplant?

- Die Klasse `BEO-PGC/arbeit-ueberholt-stehenden-traeger` trat dreifach in
  dieser Welle auf (C#-F-1, Kotlin-F-1/F-6, Python-HTTP-F-2) — in drei
  Variationen derselben Lücken-Struktur: der §3.13-Suchlauf war zu schmal
  gebunden (deklarierter Mustersatz bzw. Eintäger-Liste), während die
  bewegte Eigenschaft an Nachbar-Trägern derselben Aussage weiterbeschrieben
  stand. Siehe Steering-Loop-Einträge unten.
- `BEO-PGC/formvorbild-kopie-traegt-deutsches-wortfragment-weiter`
  erreichte mit dem Python-HTTP-Slice die 3×-Schwelle — erstmals in
  Plan-Prosa statt im SDK-Code-Kommentar; die §8-Deklaration „dieser Slice
  kopiert kein Form-Vorbild" prüfte ihre eigene Prosa nicht mit.
- Sprachspezifische Formen, die keine Spiegelung mitträgt: die Gradle-
  Config-Resolvability (eine reine Vererbung endet real in der
  ProviderNotFoundException — `grpc-netty-shaded` fehlt im Lauf-Classpath)
  und die gson-Null-Semantik (`JsonNull.INSTANCE` statt Kotlin-null) im
  Kotlin-Slice; die Form-Residuen der Phasen-Parametrisierung im
  Python-HTTP-Slice (F-5, zwei Diagnose-Zeilen phasenneutral gezogen, der
  tote Insert-Anteil bleibt als deklariertes Residuum).

## Steering-Loop-Einträge

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — hier steht der Lese-Schritt
dieser Closure (Modul 6).

- **`BEO-PGC/arbeit-ueberholt-stehenden-traeger`** (verkörpert, Anker
  [`AGENTS.md`](../../../../AGENTS.md) §3.13): drei neue Belege in dieser
  Welle — der vierundzwanzigste
  (`evidence/slice-sdk-csharp-reale2e.md`: Suchlauf an den deklarierten
  Mustersatz gebunden, drei Schreib-Verbreiterungs-Fundstellen, Kette
  F-1/Fixrunde `86892bdb`), der fünfundzwanzigste
  (`evidence/slice-sdk-kotlin-reale2e.md`: die Prüf-Angaben des
  Suchlauf-Felds selbst widerlegt — Adresse ohne Artefakt F-1,
  Behandlungs-Spalte F-6, Fixrunde `c6523009`) und der sechsundzwanzigste
  (`evidence/slice-sdk-python-http-reale2e.md`: Endklassen-Kette der zwei
  Nachbar-Zeilen derselben Werkzeug-Tabelle, F-2, Fixrunde `d668b9cc`).
  **Lese-Schritt: keine `AGENTS.md`-Schärfung nötig.** §3.13s Wortlaut
  trägt die Eigenschafts-Form bereits („`grep` über die Träger nach der
  bewegten Eigenschaft, beide Stände gemessen"); alle drei Welle-Funde
  sind Anwendungs-Verfehlungen der Regel **als geschrieben**, gefangen von
  den bestehenden Lesern (Reviewer, Fixrunde, Verifikation §3 je am
  Suchlauf-Feld) — dieselbe Entscheidung wie in der
  `welle-sdk-python-vollabdeckung`-Closure; die im vierundzwanzigsten
  Beleg offengelassene Frage („ob die Schärfung einen eigenen Satz im
  Träger trägt") wird hier damit beantwortet: nein, Register-Evidence
  trägt sie, keine doppelte Verkörperung, keine echte Lücke, kein
  Folge-Zug. Zähler 26× (Datei-Anzahl, real ausgezählt), Ausgang bleibt
  **verkörpert**.
- **`BEO-PGC/formvorbild-kopie-traegt-deutsches-wortfragment-weiter`** —
  erreichte in dieser Welle mit dem dritten Beleg (F-4 des
  Python-HTTP-Slices: das Fragment „degenerater Pfad" wanderte wortgleich
  aus der Kotlin-Plan-§7 in die Plan-Prosa) die **3×-Schwelle**. Die
  geschärfte Regel aus den drei Belegen: **Form-Vorbild-Kopien tragen je
  Form-Teil eine eigene Sprachreinheits-Sichtung** — nicht die
  Übereinstimmung mit dem Vorbild als Beleg werten, jeden kopierten
  Form-Teil selbst sichten, auch in Plan-Prosa; die Familien-Grenze
  trägt die Klasse selbst (deutsche Runner-Kommentar-Familie in
  `tools/harness/`, kein unübersetztes Fragment im englischen SDK-Doc —
  der deutschsprachige ASCII-transliterierte Runner-Kommentar ist keine
  Verstoßfläche dieser Klasse). **Lese-Schritt: die bereits verkörperten
  Anker ([`AGENTS.md`](../../../../AGENTS.md) §3.13/§3.12) tragen die
  Schärfung nicht** — eine echte Lücke bleibt. Die fällige Regelschärfung
  (eine Reviewer-Skill-Lese-Pflicht je Form-Teil oder eine Fitness
  Function) ist eine **Architect-Entscheidung** (Modul 4/8, Präzedenz
  `BEO-PGC/deutsches-fachwort-im-englischen-sdk-readme`, dieselbe
  Behandlung in der `welle-sdk-kotlin-lh-fa-sst-009`-Closure) — dieser
  Lese-Schritt trägt den Fund als offene Übergabe weiter, statt ihn
  einseitig zu embodyen; `state.md` bleibt `offen` mit dem Vermerk, dass
  die Lesung stattgefunden hat. **Folge-Zug:** eine künftige
  Architect-Sichtung (Review-Skill-Ergänzung oder Fitness Function).
  Zähler 3× (Datei-Anzahl, real ausgezählt).

**Unter der Schwelle, unverändert, kein Verkörperungsbedarf:**

- `BEO-PGC/test-runner-stiller-ausschluss` — weiterhin 2×, kein Vorkommen:
  die Runner folgen dem Muster der expliziten Testdatei-Auswahl je Phase
  (Umgebungsvariable mit laut scheiterndem Guard, Review-Negativbefund je
  Slice).
- `BEO-PGC/nachzug-laesst-ueberholten-text-stehen` — weiterhin 2×, kein
  Vorkommen: der Python-Runner-Kopf-Nachzug zog die Dreier-Phasen-Form
  real (DoD des Schluss-Slices, Verifikation §2 Zeile 4).
- `BEO-PGC/deutsches-fachwort-im-englischen-sdk-readme` — weiterhin 3×
  (offen, Architect-Entscheidung), kein neuer Beleg: die Slices berührten
  READMEs höchstens im Suchlauf-Nachzug, nicht als Schreib-Ziel.
- `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (verkörpert) —
  kein neuer Beleg: die Träger-Zahlen (RTM 79/2, Coverage 82,80 %, 12
  Träger-Zeilen) wurden je Verifikations-Lauf nachgemessen.
- `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht` (verkörpert) — kein
  neuer Beleg (die schmale Erzeuger-Glob-Form des C#-Reviews, F-5, wurde
  in der Kette gelöst und vom C#-Closure-Lese-Schritt als nicht gezählte
  Kurzform notiert); Zähler bleibt 8×.
- `BEO-PGC/github-actions-unverifizierbar-lokal` (verkörpert als
  [`AGENTS.md`](../../../../AGENTS.md) §3.10) — kein neuer Beleg: kein
  Workflow dieses Repos wurde in dieser Welle geändert, kein realer
  Tag-Push stattfand.

## Beobachtungs-Register (Zeiger)

Der Zähler steht in [`../observations/`](../observations/README.md). Ein
Eintrag erreichte in dieser Welle **neu** die 3×-Schwelle —
`BEO-PGC/formvorbild-kopie-traegt-deutsches-wortfragment-weiter`, dessen
Lese-Schritt oben als offene Übergabe an eine künftige Architect-Sichtung
getragen wird. `BEO-PGC/arbeit-ueberholt-stehenden-traeger` erhielt drei
neue Belege bei bereits verkörpertem Ausgang ohne erneute Verkörperung.

## Trigger-Audit

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 2 — drei Artefaktklassen, je eine
belegte Feststellung, jede Zeile der Welle-§3 mit eigener Messung geprüft.

- **Carveouts (Modul 7):** 0 offen — `docs/plan/carveouts/` enthält
  ausschließlich `.gitkeep`, kein Carveout referenziert [`ADR-0110`](../../adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
  oder einen der drei Slices dieser Welle (gemessen).
- **Bootstrap-aware Gates (Modul 13):** 0 betroffen — `coverage-gate`
  steht unverändert bei der 80 %-Endstufe (real gemessen 82,60–82,80 % über die Gate-Läufe dieser Closure) und
  wurde von keinem der drei Slices berührt (SDK-/Werkzeug-Bäume, außerhalb
  der netzlos prüfbaren Go-Fläche `./internal/...`+`./cmd/...`+`./gen/...`);
  `.a-check.yml` `languages: go` liest `sdks/**` strukturell nicht (real
  bestätigt: `a-check` meldet 0 Befunde auf dem Endstand).
- **Entscheidung/ADR — [`ADR-0110`](../../adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
  §Re-Evaluierungs-Trigger, durch diese Welle **nicht** ausgelöst:**
  - Trigger 1 (ein Python-Referenz-Client entsteht nachträglich) — nicht
    eingetreten: `examples/` trägt `csharp`, `kotlin` und die flache
    Go-Wurzel, **kein** `examples/python` (Verzeichnis gemessen,
    Verifikation §2 des Schluss-Slices); [`ADR-0110`](../../adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
    Festlegung 2 mandatiert keinen Beispiel-Baum.
  - Trigger 2 (ein reales Protokoll-Missverständnis im Python-Folge-Release)
    — nicht eingetreten: der reale Rot-Beleg des Schluss-Slices ging auf
    eine korrekt gebundene Negativ-Assertion zurück (gültiger Token im
    401-Test → `DID NOT RAISE` → rot, Revert, grün — Verifikation §1),
    nicht auf eine falsch verstandene HTTP-Semantik; die `SPEC-018`-
    Konformität hielt am Wire (Verifikation §5).
  - Trigger 3 ([`ADR-0107`](../../adr/0107-python-pypi-zweites-sdk-package.md)
    Triggers 1/2/4: dritte Sprache/dritter Vertriebsweg verlangt; reale
    PyPI-Nutzungsdaten; Trusted-Publishing-Umstellung) — nicht eingetreten:
    kein realer Tag-Push (`git tag -l "sdk-*-v*"` liefert ausschließlich
    `sdk-csharp-v0.1.0` und `sdk-python-v0.1.0` — gemessen), also keine
    Nutzungsdaten; keine dritte Sprache; die drei `sdk-*-release.yml`-
    Workflows wurden von keinem der drei Slices verändert.

  [`ADR-0110`](../../adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
  bleibt `Accepted` und inhaltlich unberührt durch diese Welle.

## Folge-Slices

Keine. [`ADR-0110`](../../adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
§Entscheidung Festlegung 2/Folgepflicht 1 ist mit dieser Welle für die
volle Drei-Sprachen-Vier-Wege-Matrix umgesetzt. Was ansteht, sind keine
weiteren Slices dieser Welle, sondern eigenständige, künftige
Entscheidungen:

- Die Regelschärfung der
  `BEO-PGC/formvorbild-kopie-traegt-deutsches-wortfragment-weiter`-Klasse
  (Reviewer-Skill-Lese-Pflicht je Form-Teil oder Fitness Function) —
  Architect-Entscheidung, als offene Übergabe des Lese-Schritts oben
  getragen, kein Slice.
- Ein realer `sdk-*-v*`-Tag-Push — irreversible, extern sichtbare
  Betreiber-Handlung nach [`AGENTS.md`](../../../../AGENTS.md) §3.10 (wie
  bei den vorigen SDK-Wellen); diese Welle verlangt ihn ausdrücklich
  nicht (Closure-Trigger, Welle-Datei §3).

## Verifikation

- `docs/reviews/review-slice-sdk-csharp-reale2e.md` (F-1…F-8, 0 HIGH ·
  1 MEDIUM · 5 LOW · 2 INFO) + Fixrunde (`86892bdb`, F-1 bis F-6) +
  `docs/reviews/verifikation-slice-sdk-csharp-reale2e.md` (DoD-Verdikt:
  erfüllt, eigener Pflichtbeleg-Lauf EXIT=0 — gRPC `change_id=804-1`,
  SSE `807-1`, NATS `810-1`, `consumer_id=csharp-sdk-e2e-20260923082238`
  laufgebunden; der DoD-Ursprungs-Lauf trug
  `consumer_id=csharp-sdk-e2e-20260923075217`).
- `docs/reviews/review-slice-sdk-kotlin-reale2e.md` (F-1…F-6, 0 HIGH ·
  1 MEDIUM · 2 LOW · 3 INFO) + Fixrunde (`c6523009`, F-1, F-2, F-3, F-6) +
  `docs/reviews/verifikation-slice-sdk-kotlin-reale2e.md` (DoD-Verdikt:
  erfüllt mit drei nicht-blockierenden Beobachtungen, eigener Lauf EXIT=0
  — gRPC `805-1`, SSE `810-1`, NATS `817-1`, `consumer_id=kotlin-sdk-e2e-
  20260923105937` laufgebunden; der DoD-Ursprungs-Lauf trug
  `806-1`/`814-1`/`818-1`; zusätzlich `make sdk-pack-kotlin` EXIT=0 über
  die Integrations-Quellmenge).
- `docs/reviews/review-slice-sdk-python-http-reale2e.md` (F-1…F-6, 0 HIGH ·
  2 MEDIUM · 3 LOW · 1 INFO) + Fixrunde (`d668b9cc`, F-1 bis F-5) +
  `docs/reviews/verifikation-slice-sdk-python-http-reale2e.md`
  (DoD-Verdikt: erfüllt, 0 DoD-Abweichungen; eigener Pflichtbeleg-Lauf
  EXIT=0 — gRPC `804-1`, SSE `807-1`, NATS `810-1`,
  `consumer_id=python-sdk-e2e-570d402081d4` laufgebunden; der
  DoD-Ursprungs-Lauf trug `813-1`/`consumer_id=python-sdk-e2e-4ad61b9f8a0d`;
  der reale Rot-Beleg durch eigene Mutation, EXIT=2 → Revert → grün;
  Link-Berichtigung `a78aed81`). Die DoD-/Verifier-Werte sind aus den
  Reports übernommen und dort laufgebunden verankert (`AGENTS.md` §3.12
  Instanz A).
- `make gates`: grün nach jedem Inhalts-Commit dieser Closure; nach dem
  reinen Welle-Move (`bb3dbfda`) misste der Lauf real **46
  target-missing-Befunde** (Exit 2, docs-check only — Muster `848d2263`,
  der Move-Commit bleibt rein) und wurde vom Reconciliations-Commit
  (`66bade56`) auf Grün geführt (ungepiped, Exit-Code direkt geprüft,
  [`AGENTS.md`](../../../../AGENTS.md) §3.9) — 937–938 Datei(en) je
  Gate-Lauf dieser Closure, 0 Befunde; Coverage 82,60–82,80 % ≥ 80 %;
  `a-check`/`generated-sync`/`commit-traceability`/`baseline-verify` je
  ohne Befund; Stempel-Gleichheit je Lauf.
- Realer, grüner Pflichtbeleg je Sprache (Closure-Trigger, Welle-Datei §3):
  `make test-sdk-csharp-integration`, `make test-sdk-kotlin-integration`
  und `make test-sdk-python-integration` je EXIT=0 mit allen vier
  Flächen-Phasen gegen eine reale, laufende Server-Instanz — Implementer-/
  Verifier-Läufe, laufgebunden in den Reports (Zeiger oben), je mit
  SQL-Gegenprüfung (`cdc.changes` gRPC/SSE/NATS, `cdc.consumer` HTTP) und
  je Ablehnungs-Beleg (gRPC `Unauthenticated`, SSE/HTTP `401`,
  NATS-Verbindungsablehnung).
- `git tag -l "sdk-*-v*"`: ausschließlich `sdk-csharp-v0.1.0` und
  `sdk-python-v0.1.0` — bestätigt, dass kein realer Tag-Push stattfand
  (Welle-Datei §3 vorab festgelegt).
- `docs/plan/carveouts/`: nur `.gitkeep` — 0 offene Carveouts dieser Welle.

## Archivierung

Dieses Repo führt kein Archivierungs-Werkzeug für Wellen-Zeitdokumente
(kein `archiv`-Ziel in `Makefile`/`harness/mk/*.mk`, real geprüft per
`grep` — dieselbe Feststellung wie bei den vorigen SDK-Wellen) — die
Bedingung für Schritt 4 der Closure-Prozedur ist nicht eingetreten, keine
Handarbeit als Ersatz.