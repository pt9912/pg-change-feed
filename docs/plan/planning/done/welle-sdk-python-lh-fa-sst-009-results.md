# Welle welle-sdk-python-lh-fa-sst-009 — Closure-Notiz

**Welle:** welle-sdk-python-lh-fa-sst-009
**Abschluss:** 2026-09-19
**Verantwortlich:** pt9912

## Was wurde geliefert?

[`ADR-0107`](../../adr/0107-python-pypi-zweites-sdk-package.md) (`Accepted`)
und ihre partielle Supersession
[`ADR-0108`](../../adr/0108-python-sdk-uv-statt-build-twine.md) (`Accepted`)
vollständig umgesetzt — das zweite offizielle Client-Bibliothek-Package für
[`LH-FA-SST-009`](../../../../spec/lastenheft.md), in vier Slices, rein
sequenziell (keine Parallelisierung — anders als bei der C#-Welle gibt es
nur eine Draht-Fläche, HTTP):

- **`slice-sdk-python-projektgeruest`**: `sdks/python/pgchangefeed/` neu
  angelegt — `pyproject.toml` mit PyPA-Standard-Metadaten (`name =
  "pgchangefeed"`, `version = "0.1.0"`), eigenständiges, digest-gepinntes
  Docker-Bau-Setup, englischsprachiges `README.md`, ein minimales
  Konfigurations-Skelett (`ClientOptions`). Ein realer, nutzergetriebener
  Nachzug hob `requires-python`/die Docker-Basis von `>=3.11`/
  `python:3.13-slim` auf `>=3.14`/`python:3.14-slim` an.
- **`slice-sdk-python-http-client-flaeche`**: öffentliche Python-API-Fläche
  für alle neun Port-gedeckten Fähigkeiten von `SPEC-018` plus das
  Changes-Lesen (`SPEC-022`) — Bearer-Token-Auth, eine zentrale
  Exception-Hierarchie über einen gemeinsamen `_handle`-Helfer für alle
  zehn Methoden, 23 eigene, netzlos laufende `pytest`-Tests.
- **`slice-sdk-python-pack-werkzeug`**: `make sdk-pack-python` (Docker-only,
  kein Gate, `uv build --no-sources` gemäß `ADR-0108` statt `build`+`twine`)
  baut, testet und paketiert die HTTP-Fläche zu einem realen
  `pgchangefeed-0.1.0-py3-none-any.whl` + `pgchangefeed-0.1.0.tar.gz` —
  Träger-Nachzug in `spec/pflichtenheft.md` §1 (`LH-FA-SST-009.a`) und §6
  (neue `SPEC-027`-Zeile) sowie `harness/README.md` §Werkzeuge.
- **`slice-sdk-python-publish-workflow`**: `.github/workflows/sdk-python-release.yml`
  (Trigger `push: tags: ['sdk-python-v*']`, eigener Tag-Namensraum getrennt
  von `release.yml`/`sdk-csharp-release.yml`), Tag-vs-`pyproject.toml`-
  Versionsabgleich (eigenständige PEP-440-Einfachfall-Regex, nicht geteilt
  mit `tools/harness/semver-regex.sh`), `uv publish` gegen PyPI (Token über
  `UV_PUBLISH_TOKEN`), referenziert `PYPI_API_TOKEN` (real geprüft: **noch
  nicht** angelegt — anders als bei der C#-Welle), `docs/user/releasing.md`-
  Nachzug im selben Slice (vorab geplant, kein Reviewer-Finding nötig).

`make gates` grün auf dem Endstand (839 Dateien, 0 `docs-check`-Befunde,
Coverage 82,80 % ≥ 80 %, `a-check`/`generated-sync`/`commit-traceability` je
ohne Befund, `baseline-verify` OK). Reale `.whl`/`.tar.gz`-Artefakte
existieren (`slice-sdk-python-pack-werkzeug` §7, mehrfach unabhängig
bestätigt). Kein realer `sdk-python-v*`-Tag wurde gepusht (`git tag -l
"sdk-python-v*"` leer) — diese Welle liefert bewusst nur den **Mechanismus**;
`PYPI_API_TOKEN` existiert real **nicht** als Repository-Secret (`gh secret
list`: nur `DOCKERHUB_TOKEN`, `DOCKERHUB_USERNAME`, `NUGET_API_KEY`), ein
realer Tag-Push würde am fehlenden Secret scheitern, bevor er PyPI überhaupt
erreicht — anders als bei der C#-Welle, wo `NUGET_API_KEY` bereits existierte.

## Was hat funktioniert?

Über den gesamten Wellen-Zyklus hinweg trug eine durchgehende, mehrfach
unabhängige Verifikation, oft vier- bis sechsfach:

- `slice-sdk-python-projektgeruest`: fünf unabhängige Docker-Builds
  (Implementer, Reviewer-Erstlauf, Reviewer-Fixrunden-Nachprüfung, Verifier,
  diese Planner-Closure), 3/3 Tests grün, keine Divergenz — auch über den
  nachträglichen Python-3.14-Bump hinweg stabil.
- `slice-sdk-python-http-client-flaeche`: der sauberste Zyklus **nach** der
  Fixrunde — Reviewer und Verifikation führten je eine vollständige,
  unabhängige Feld-für-Feld-Gegenprobe aller zehn Fähigkeiten plus zwölf
  `Change`-Felder gegen `SPEC-018`/`SPEC-022` durch, beide ohne Abweichung.
- `slice-sdk-python-pack-werkzeug`: 0 HIGH/MEDIUM/LOW auf Anhieb (1 LOW, 1
  INFO), der `uv`-Digest wurde real gegen die Registry reverifiziert und war
  identisch zum `ADR-0108`-Wert; der `tar`-Stream-Export trug beim ersten
  Versuch für beide Artefakte in einem Stream.
- `slice-sdk-python-publish-workflow`: 0 Findings auf Anhieb, vier
  unabhängige Reverifikationen (inkl. SHA-Pins, PEP-440-Isolation), keine
  Fixrunde.

Die proaktive Vermeidung der aus der C#-Welle bekannten Fixrunden-Klassen hat
**teilweise** funktioniert: Der Implementer von `slice-sdk-python-http-client-flaeche`
wandte den `AGENTS.md` §3.13-Träger-Nachzug-Suchlauf bereits vor dem
Erstreview proaktiv an (zwei Stellen korrigiert, `README.md`/`__init__.py`)
und vermied beide im C#-Review gefundenen Design-Findings-Klassen
(methodenspezifischer Zweitpfad, malformter 2xx-Erfolgskörper) von
vornherein durch einen zentralen `_handle`-Helfer. Das etablierte
3-Commit-Move-Muster (`git mv` · Inhalt · `git mv`) hielt `AGENTS.md` §3.3
über alle vier Slices sauber.

## Was ging anders als geplant?

Zwei der vier Slices durchliefen eine Reviewer-Fixrunde, zwei blieben beim
ersten Durchlauf ohne HIGH/MEDIUM:

- `slice-sdk-python-projektgeruest`: 1 HIGH (F-1, Slice-Chronik im
  `__init__.py`-Docstring) — Fixrunde behoben (Commit `70aa64ce`),
  Fixrunden-Nachprüfung ohne neuen Fund. Zusätzlich ein eigener kleiner
  Fixrunden-Zyklus für einen **durch den realen Python-Versions-Bump
  (3.11→3.14, Nutzerentscheidung nach Verifikation) entstandenen
  Plan-Widerspruch**: der neue Nachzug-Absatz in §3 stand unverbunden neben
  dem dadurch überholten alten Absatz (F-2, MEDIUM, Coordinator-Fix
  `e7670f6f`, ohne erneuten Reviewer-Durchlauf — Planner-Closure prüfte die
  Widerspruchsfreiheit selbst nach).
- `slice-sdk-python-http-client-flaeche`: 1 HIGH (F-1) + 1 MEDIUM (F-2).
  **Trotz** der expliziten Anweisung, die C#-Fixrunden-Fehlerklassen vorab
  zu vermeiden, trat eine verwandte Klasse trotzdem auf, nur an einer
  anderen Stelle: `pyproject.toml:22-24` behauptete weiterhin „folgt erst"
  für genau die jetzt gelieferte Fläche — derselbe Träger-Nachzug-Suchlauf,
  der `README.md`/`__init__.py` bereits proaktiv korrigiert hatte, deckte
  weder die Datei (Glob-Lücke) noch die deutsche Formulierung
  (Wortmuster-Lücke „folgt erst" statt nur „follow-up"/„added by") ab. Das
  ist selbst ein lehrreicher Fund: Ein aus einem vorigen Vorkommen
  übernommener Suchlauf garantiert keine Vollständigkeit am nächsten,
  strukturell ähnlichen Vorgang — beide Achsen (welche Dateien, welche
  Formulierungen) müssen unabhängig breit genug gewählt werden. F-2
  (Commit-Message-Zahlen-Ungenauigkeit „23" vs. Kategorisierung) wurde per
  dokumentierter Klarstellung statt `--amend` behoben (`AGENTS.md` §3.12).
  Fixrunde `d0da688b`, Fixrunden-Nachprüfung ohne neuen Fund.
- `slice-sdk-python-pack-werkzeug` und `slice-sdk-python-publish-workflow`:
  glatte Durchläufe, 0 HIGH/MEDIUM.

**Der reale Python-Versions-Bump (3.11→3.14):** eine Nutzerentscheidung nach
Verifikation, dieselbe Analogie zur `net10.0`-Wahl bei C# — jetzt real
vollzogen statt nur als Option offengehalten, mit dem oben genannten eigenen
kleinen Fixrunden-Zyklus für den dabei entstandenen Plan-Widerspruch.

**Die uv-statt-build+twine-Supersession (`ADR-0108`):** während der
laufenden Welle entschieden (Nutzerrückfrage per `AskUserQuestion`, echte
Alternativen-Lücke in `ADR-0107`s Tabelle E) und rückwirkend in die noch
offenen Slice-Pläne (`slice-sdk-python-pack-werkzeug`,
`slice-sdk-python-publish-workflow`, die Welle-Datei selbst) nachgezogen,
bevor sie umgesetzt wurden — kein Code musste migriert werden, weil zum
Entscheidungszeitpunkt noch nichts von `build`/`twine` real existierte.

**Der real aufgedeckte Gate-Maskierungsfund:** Der Reviewer von
`slice-sdk-python-pack-werkzeug` bemerkte im eigenen Report-Entwurf einen
unclosed-backtick-Fehler, korrigierte ihn vor dem eigenen Commit, und
vermutete (ohne selbst zu prüfen), dasselbe Muster könnte im bereits
gemergten `docs/reviews/review-slice-sdk-csharp-pack-werkzeug.md` vorliegen.
Die Planner-Closure dieses Slice maß die Vermutung real nach: ein
Backtick-Gesamtzahl-Zähllauf (`grep -o` + `wc -l`) ergab 403 (ungerade) — ein
realer unclosed-backtick-Defekt an Zeile 122f. eines bereits gemergten
Nachbar-Reports, unabhängig vom Diff dieser Welle. Die Korrektur (403 → 404)
legte einen **zweiten**, bis dahin verdeckten Fund offen: `make docs-check`
meldete danach einen `id-unlinked`-Befund für eine nackte
`ADR-0106`-Erwähnung an Zeile 132 desselben Reports — real unverlinkt seit
dem ursprünglichen Merge, solange durch den Paritäts-Defekt fälschlich als
„innerhalb eines Codespans" maskiert. `make gates`/`make docs-check` liefen
über den gesamten Zeitraum grün — nicht weil kein echter Fund vorlag,
sondern weil der Backtick-Defekt ihn verbarg. Beide Funde in eigenen,
separat committeten Fixes behoben.

**Zwei reale, externe Doku-Drifts in `docs/user/releasing.md`, verursacht
durch Ereignisse außerhalb dieser Welle:** Zwischen der C#-Welle-Closure und
dieser Python-Welle fand real ein `sdk-csharp-v0.1.0`-Tag-Push und
NuGet-Publish statt (eine externe Betreiber-Handlung, keine dieser Welle
zurechenbare Arbeit). `docs/user/releasing.md` §1 sowie zwei weitere Absätze
behaupteten danach fälschlich weiterhin „kein realer Tag gesetzt". Beide
Drifts wurden während dieser Welle gefunden und behoben (Commit `38a137f7`,
außerhalb der Slice-DoD, aber vom Reviewer/Verifier von
`slice-sdk-python-publish-workflow` mitgeprüft und gegen NuGet.org/`git tag
-l` real bestätigt) — dies ist eine **neue Fehlerklasse**: das überholende
Ereignis lag zum ersten Mal vollständig **außerhalb** jedes Commits und
jeder Versionskontrolle dieses Repos (ein realer, externer Publish-Vorgang),
nicht in einem Diff. Kein Sensor und kein Diff hätte das zeigen können;
gefunden hat es der Coordinator beim Gegenlesen des Nachbarabschnitts, nicht
der auf den eigenen Slice-Gegenstand begrenzte `AGENTS.md` §3.13-Suchlauf
des Implementers.

## Trigger-Audit

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 2 — drei Artefaktklassen, je eine belegte
Feststellung.

- **Carveouts (Modul 7):** 0 offen — `docs/plan/carveouts/` enthält
  ausschließlich `.gitkeep`, kein Carveout referenziert `ADR-0107`,
  `ADR-0108` oder eine dieser vier Slices.
- **Bootstrap-aware Gates (Modul 13):** 0 betroffen — `coverage-gate` steht
  bereits auf der Endstufe (80 %, ausgeschöpft, `harness/mk/coverage.mk`)
  und wurde von keinem der vier Slices berührt (reiner Python-Baum,
  außerhalb der netzlos prüfbaren Go-Fläche `./internal/...`+`./cmd/...`+
  `./gen/...`); die gemessene Coverage auf dem Endstand dieser Welle
  (82,80 %) liegt unverändert über der Schwelle. `.a-check.yml`
  `languages: go` liest `sdks/python/**` strukturell nicht (real bestätigt:
  `a-check` meldet 0 Befunde auf dem Endstand).
- **Entscheidung/ADR — ZWEI ADRs zu prüfen:**
  - [`ADR-0107`](../../adr/0107-python-pypi-zweites-sdk-package.md)
    §Re-Evaluierungs-Trigger, durch diese Welle **nicht** ausgelöst:
    - Trigger 1 (dritte Sprache/dritter Vertriebsweg verlangt) — nicht
      eingetreten, out-of-scope §6 der Welle-Datei.
    - Trigger 2 (real auf PyPI veröffentlicht + Nutzungsdaten) — **nicht**
      eingetreten. Das Package ist real paketierbar (`.whl`/`.tar.gz`
      existieren, `slice-sdk-python-pack-werkzeug`), aber **nicht** real
      veröffentlicht: `git tag -l "sdk-python-v*"` ist leer, kein Tag-Push
      erfolgte, `uv publish` lief nie real gegen PyPI, und
      `PYPI_API_TOKEN` fehlt real als Repository-Secret (`gh secret list`).
      Ohne reale Veröffentlichung gibt es keine Download-Zahlen und keine
      Issue-Nachfrage. **Beobachtbare Auslösebedingung** (analog dem
      `ADR-0106`-Präzedenzfall in
      `welle-sdk-csharp-lh-fa-sst-009-results.md`): zunächst die externe
      Anlage von `PYPI_API_TOKEN` als Repository-Secret (Voraussetzung, die
      bei der C#-Welle bereits erfüllt war, hier noch nicht), danach der
      erste reale `git tag sdk-python-v<PEP 440> && git push --tags`, der
      einen grünen `sdk-python-release.yml`-Lauf erzeugt, und danach
      messbare Nutzungssignale (PyPI-Downloadzahlen, GitHub-Issues mit
      SDK-Bezug) — prüfbar über `git tag -l "sdk-python-v*"` (nicht mehr
      leer), einen grünen Actions-Lauf und einen PyPI-Paketseiten-Blick.
      Bis dahin bleibt `ADR-0107` Festlegung 1 (Umfang: nur HTTP-API, kein
      gRPC/SSE/NATS) unverändert in Kraft.
    - Trigger 3 (Python-Referenz-Client entsteht nachträglich) — nicht
      eingetreten, `examples/python/` existiert weiterhin nicht.
    - Trigger 4 (PyPI Trusted Publishing für ein bereits real existierendes
      Package geprüft) — nicht eingetreten, kein Package real veröffentlicht.
  - [`ADR-0108`](../../adr/0108-python-sdk-uv-statt-build-twine.md)
    §Re-Evaluierungs-Trigger, durch diese Welle **nicht** ausgelöst:
    - Trigger 1 (`uv publish` erweist sich im ersten realen
      `sdk-python-release.yml`-Lauf als instabil/inkompatibel) —
      **noch nicht beobachtbar**, da `uv publish` real noch nie gegen PyPI
      gelaufen ist (kein Tag-Push, kein Secret). Die ADR benennt diesen
      Reifegrad selbst ausdrücklich als offenen Punkt (§Kontext), nicht als
      bewiesenen Erfolg — diese Welle ändert daran nichts, weder in die eine
      noch in die andere Richtung.
    - Trigger 2 (PyPA/packaging.python.org erklärt `uv` offiziell zum
      empfohlenen Minimalweg) — nicht eingetreten, keine externe
      Ereignis-Prüfung Teil dieser Welle.
    - Trigger 3 (`uv` erreicht Version 1.0 mit expliziter
      Stabilitätszusage für `uv publish`) — nicht eingetreten; der real
      gemessene Stand bleibt `0.12.17` (`slice-sdk-python-pack-werkzeug`
      §7, gegen die Registry reverifiziert).
    - Trigger 4 (`ADR-0107` Trigger 4 tritt ein) — nicht eingetreten,
      dieselbe Abhängigkeit wie oben.

## Steering-Loop-Einträge

Kein neuer Eintrag erreicht in dieser Welle erstmals die 3×-Schwelle mit
Verkörperungs-Bedarf. Bereits vor Wellen-Beginn verkörpert und unverändert
gültig, mit neuen Belegen in dieser Welle:

- `BEO-PGC/arbeit-ueberholt-stehenden-traeger` (verkörpert seit `AGENTS.md`
  §3.13): zwei neue Belege in dieser Welle
  (`slice-sdk-python-http-client-flaeche`, 18. Beleg — die
  `pyproject.toml`-Fundstelle; `slice-sdk-python-publish-workflow`, 19.
  Beleg — der `docs/user/releasing.md`-Fund, dessen überholendes Ereignis
  erstmals vollständig außerhalb jeder Versionskontrolle lag), Zähler jetzt
  real 19× (`ls .../evidence | wc -l`). Bestätigt: bereits verkörpert, keine
  erneute Verkörperung nötig.
- `BEO-PGC/report-nackte-id-ohne-link` (verkörpert seit `slice-063` in
  `AGENTS.md` §3.9): geprüft, kein neuer Beleg in dieser Welle — Zähler
  bleibt real bei 8×.

**Neu angelegt, unter der Schwelle:**

- `BEO-PGC/unclosed-backtick-taeuscht-nackte-id-vor` (Erstauftreten,
  `slice-sdk-python-pack-werkzeug`, 1×) — ein einzelnes fehlendes
  schließendes Backtick verschiebt die Codespan-Parität für den Rest eines
  Dokuments und kann dadurch eine an ihrer eigenen Stelle korrekt
  gebacktickte Kennung an anderer Stelle desselben Dokuments vor `d-check`
  verbergen (Fernwirkung, real feuernd nachgewiesen, siehe „Was ging anders
  als geplant"). Abgrenzung zu `BEO-PGC/report-nackte-id-ohne-link`: dort
  fehlt die Verlinkung an der eigenen Stelle, hier maskiert ein entfernter
  Syntaxfehler eine an ihrer eigenen Stelle korrekte Kennung. Kein neuer
  Sensor — Kandidat für eine künftige Reviewer-Skill-Schärfung, aber unter
  der 3×-Schwelle.
- `BEO-PGC/nachzug-laesst-ueberholten-text-stehen` (Erstauftreten,
  `slice-sdk-python-projektgeruest`, 1×) — ein Plan-Nachzug-Absatz wird
  korrekt und rechtzeitig ergänzt, ohne den dadurch überholten Absatz an
  anderer Stelle desselben Dokuments zu bereinigen oder zu markieren.
  Andere Klasse als `BEO-PGC/plan-nachzug` (Nachzug fehlt ganz) und
  `BEO-PGC/arbeit-ueberholt-stehenden-traeger` (überholter Träger liegt
  außerhalb des Diffs) — hier liegt der überholte Text in derselben Datei
  und demselben Bearbeitungsanlass. Kein neuer Sensor (Prosa-Kohärenz ist
  kein Gate-fähiges Muster); Wächter bleibt Review/Planner-Closure-Lektüre.

Die aus der Eröffnungs-Sichtung der Welle-Datei §6 bereits bekannten
Einträge blieben ohne dritten Treffer:

- `BEO-PGC/workaround-uebersieht-etabliertes-muster-im-bestand` (Zähler
  weiterhin 2×, `slice-104`/`slice-sdk-csharp-projektgeruest`) — kein
  drittes Auftreten in dieser Welle, unter der Schwelle. Real getragen hat
  der Vorbild-Suchlauf (`sdks/csharp/**`) in
  `slice-sdk-python-projektgeruest` — kein neuer Beleg dieser Klasse, das
  übernommene README-COPY-Workaround-Muster trug real, auch für ein
  strukturell anderes Problem als beim C#-Vorbild.
- `BEO-PGC/release-mechanismus-nicht-in-releasing-doku-nachgezogen` (Zähler
  weiterhin 1×, `slice-sdk-csharp-publish-workflow`) — kein zweites
  Auftreten: `slice-sdk-python-publish-workflow` nahm den
  `docs/user/releasing.md`-Nachzug bewusst vorab als eigenen DoD-Punkt auf,
  genau um ein zweites Auftreten zu vermeiden. Gelungen — der real
  aufgetretene `docs/user/releasing.md`-Fund dieser Welle (siehe „Was ging
  anders als geplant") ist eine **andere** Fehlerklasse (externe
  Ereignis-Drift, nicht „Mechanismus nicht nachgezogen").
- `BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`
  (Zähler weiterhin 3×, bereits verkörpert) — kein neuer Beleg, der
  DoD-Punkt in `slice-sdk-python-http-client-flaeche` trug ihn bereits vorab.
- `BEO-PGC/github-actions-unverifizierbar-lokal` (verkörpert als `AGENTS.md`
  §3.10, Zähler weiterhin 7×) — das Post-Push-Risiko dieser Welle bleibt
  strukturell offen (siehe Trigger-Audit oben), kein neuer Beleg, da kein
  realer Tag-Push stattfand.

## Beobachtungs-Register (Zeiger)

Der Zähler steht in [`../observations/`](../observations/). Kein Eintrag
dieser Welle hat die 3×-Schwelle **erstmals** erreicht; `BEO-PGC/arbeit-ueberholt-stehenden-traeger`
steht bereits bei 19× und ist als `AGENTS.md` §3.13 bereits verkörpert —
diese Closure bestätigt das, keine erneute Verkörperung nötig. Die beiden
neuen Einträge (`unclosed-backtick-taeuscht-nackte-id-vor`,
`nachzug-laesst-ueberholten-text-stehen`) stehen bei 1×, siehe oben.

## Nebenbefunde (außerhalb des Scopes dieser Welle, hier nur gemeldet)

- **Zwei reale, externe `docs/user/releasing.md`-Doku-Drifts, verursacht
  durch den `sdk-csharp-v0.1.0`-Tag-Push/NuGet-Publish zwischen den beiden
  Wellen** — bereits während dieser Welle gefunden und im Commit `38a137f7`
  behoben (siehe „Was ging anders als geplant"). Gemeldet hier als eigene,
  neue Fehlerklasse: eine Doku-Drift, deren Ursache vollständig außerhalb
  jeder Versionskontrolle dieses Repos liegt.
- **`PYPI_API_TOKEN` existiert real nicht als Repository-Secret** (`gh
  secret list`, 2026-09-19: nur `DOCKERHUB_TOKEN`, `DOCKERHUB_USERNAME`,
  `NUGET_API_KEY`) — anders als bei der C#-Welle (`NUGET_API_KEY` existierte
  dort bereits). Ein realer `sdk-python-v*`-Tag-Push würde am
  Publish-Schritt scheitern, bevor er PyPI überhaupt erreicht. Diese Welle
  legt das Secret bewusst **nicht** an — eine externe, außerhalb dieses
  Repos liegende Betreiber-Handlung.
- **`pgchangefeed` ist real frei auf PyPI** (`curl -s -o /dev/null -w
  '%{http_code}' https://pypi.org/pypi/pgchangefeed/json` → `404`, real
  geprüft während `slice-sdk-python-publish-workflow`) — kein
  Namensraum-Konflikt für einen künftigen ersten Publish-Versuch.

## Folge-Slices

Keine. `ADR-0107`/`ADR-0108` sind mit dieser Welle für Python/PyPI
vollständig umgesetzt. Was ansteht, sind keine weiteren Slices dieser Welle,
sondern eigenständige, künftige Entscheidungen:

- Ein realer `sdk-python-v0.1.0`-Tag-Push — eine irreversible, extern
  sichtbare Aktion (öffentlicher PyPI-Push), die zusätzlich die externe
  Anlage von `PYPI_API_TOKEN` voraussetzt und nur nach expliziter,
  gesonderter Rückfrage beim Auftraggeber läuft (Welle-Datei §3 vierter
  Punkt, `AGENTS.md` §3.10).
- Eine dritte SDK-Sprache oder ein zweiter Python-Vertriebsweg — bleibt laut
  `ADR-0107` Re-Evaluierungs-Trigger 1 offen für eine künftige, separate
  ADR; kein Slice dieser Welle nimmt das vorweg.
- Ein Folge-Package für gRPC/SSE/NATS in Python — bleibt laut `ADR-0107`
  §Entscheidung Festlegung 1/§Re-Evaluierungs-Trigger 2/3 offen für ein
  künftiges Release, nicht Teil dieser Welle.

## Verifikation

- `docs/reviews/review-slice-sdk-python-projektgeruest.md` (1 HIGH,
  Fixrunde `70aa64ce` real behoben, Fixrunden-Nachprüfung ohne neuen Fund;
  zusätzlicher Nachtrag zum Versions-Bump mit 1 MEDIUM, Coordinator-Fix
  `e7670f6f`) + `docs/reviews/verifikation-slice-sdk-python-projektgeruest.md`.
- `docs/reviews/review-slice-sdk-python-http-client-flaeche.md` (1 HIGH + 1
  MEDIUM, Fixrunde `d0da688b` real behoben, Fixrunden-Nachprüfung ohne neuen
  Fund) + `docs/reviews/verifikation-slice-sdk-python-http-client-flaeche.md`.
- `docs/reviews/review-slice-sdk-python-pack-werkzeug.md` (0 HIGH/MEDIUM, 1
  LOW, 1 INFO, keine Fixrunde) +
  `docs/reviews/verifikation-slice-sdk-python-pack-werkzeug.md`.
- `docs/reviews/review-slice-sdk-python-publish-workflow.md` (0
  HIGH/MEDIUM/LOW, keine Fixrunde) +
  `docs/reviews/verifikation-slice-sdk-python-publish-workflow.md`.
- `make gates`: grün auf dem Endstand (839 Dateien, 0 Befunde; Coverage
  82,80 % ≥ 80 %; `a-check`/`generated-sync`/`commit-traceability`/
  `baseline-verify` je ohne Befund), ungepiped geprüft (`AGENTS.md` §3.9)
  nach jedem Commit dieser Welle sowie erneut vor dieser Closure.
- Reale `.whl`/`.tar.gz`-Artefakte (`pgchangefeed-0.1.0-py3-none-any.whl`,
  `pgchangefeed-0.1.0.tar.gz`) — mehrfach unabhängig real erzeugt
  (`slice-sdk-python-pack-werkzeug` §7).
- `git tag -l "sdk-python-v*"`: leer — bestätigt, dass kein realer
  PyPI-Publish-Versuch stattfand (Welle-Datei §3 vierter Punkt eingehalten,
  Post-Push-Risiko bleibt strukturell offen nach `AGENTS.md` §3.10).
- `gh secret list`: `PYPI_API_TOKEN` fehlt real — bestätigt, dass ein realer
  Tag-Push am Publish-Schritt scheitern würde.
- `docs/plan/carveouts/`: nur `.gitkeep` — 0 offene Carveouts dieser Welle.

## Archivierung

Dieses Repo führt kein Archivierungs-Werkzeug für Wellen-Zeitdokumente (kein
`archiv`-Ziel in `Makefile`/`harness/mk/*.mk`, real geprüft per `grep`) — die
Bedingung für Schritt 4 der Closure-Prozedur ist nicht eingetreten, keine
Handarbeit als Ersatz.
