# ADR-0085: Build-Kontext-Ausnahme der `coverage`-Stufe; `test-only` auf den Zweck geschärft

**Status:** Accepted — Supersedes [`ADR-0082`](0082-coverage-schnittmass-composition-root-nicht-netzlos.md)
in genau **zwei Klauseln**: der Folgepflicht „test-only" (Implementer-Zug, je
Test-Slice) in deren §Konsequenzen und der Test-only-Zeile ihrer
§Fitness Function. Alles Übrige aus `ADR-0082` — §Entscheidung mit den fünf
Festlegungen (Schnittmaß, ungezogene Naht am Composition Root, Endstufe 80 %,
Messgegenstand, nicht baubarer Bootstrap-Smoke-Test), §Kontext, §Verglichene
Alternativen, die übrigen Folgepflichten, die übrigen Fitness-Function-Zeilen
und die Re-Evaluierungs-Trigger — bleibt unverändert bestehen und wird hier
nicht wiederholt.

**Datum:** 2026-09-16

**Autor:** pt9912 (Architect-Rolle, unabhängiger Architect-Zug — anderer
Kontext als der Implementer-Lauf von `slice-093` <!-- d-check:status-provenance -->,
der die Abweichung selbst offen nannte, und als der Reviewer-Lauf, der sie als
F-1 (HIGH) gemessen hat; Modul 8 §Konflikt-Pfad als Rollen-Sequenz, drittes
Verdikt: *Lockerung legitim, aber undokumentiert* — der Zug gehört der
Architect-Rolle, nicht dem Implementer)

**Bezug:** [`ADR-0082`](0082-coverage-schnittmass-composition-root-nicht-netzlos.md)
(in zwei Klauseln abgelöst, sonst bestätigt) ·
[`ADR-0071`](0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
(Messgegenstand — unberührt) · [`ADR-0054`](0054-coverage-gate-und-benchmark-infrastruktur.md)
(die `coverage`-Stufe und ihr Gate) ·
[`ADR-0044`](0044-image-beleg-semantik.md) (Image-Beleg-Semantik; Punkt 3 —
die bestehende Regel, auf der Merkmal (iv) dieser ADR aufsitzt) ·
[`ADR-0073`](0073-zitat-korrektur-an-immutablen-dokumenten.md) (die
Zitat-Korrektur-Klasse — sie deckt §Konsequenzen und Fitness-Function-Regeln
**nicht**) · `AGENTS.md` §3.5 (Accepted-ADRs sind immutable), §3.6 (Schwellen
nur per ADR), §3.7 (Ist-Zustand), §3.12 (Herkunft von Aussagen) · `Dockerfile`
(Stufen `deps`, `proto`, `coverage`, `build`, `runtime`) · `.dockerignore` ·
`harness/mk/coverage.mk` (`THRESHOLD`) · `tools/coverage-gate.sh` ·
`harness/sensors/coverage-gate.md` · `harness/README.md` §Sensors/§Werkzeuge
(`make coverage-gate`, `make image`) · `harness/image-hash.txt` ·
`docs/plan/planning/welle-20.md` · Review zu `slice-093` <!-- d-check:status-provenance -->
(F-1) · `docs/plan/planning/in-progress/slice-093-coverage-cluster-d2.md` <!-- d-check:status-provenance -->
(§1, §3) · `docs/plan/planning/observations/BEO-PGC/coverage-stage-dockerignore-blockiert-tooling/observation.md` ·
`docs/plan/planning/observations/BEO-PGC/arbeit-ueberholt-stehenden-traeger/observation.md` ·
`internal/bootstrap/roles_rollout_file_internal_test.go` (der Leser der
Ausnahme)

**Schärft:** — (Prozess-/Tooling-ADR ohne Spec-Stratum, wie `ADR-0054`,
`ADR-0071`, `ADR-0082`)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

**(1) Der Anlass: eine Folgepflicht, die im Wortlaut bricht, während ihr Zweck
hält.** `ADR-0082` §Konsequenzen trägt für jeden Test-Slice dieser Welle:

> Folgepflicht (Implementer-Zug, je Test-Slice): **test-only**. Der Diff eines
> Slice dieser Welle enthält ausschließlich `*_test.go` (und, für Cluster A,
> keinen Produktionscode außerhalb der Testdatei); keine Zeile `internal/**`
> oder `cmd/**` außerhalb von Tests ändert sich, um Coverage zu gewinnen.

Der Diff, den der Review prüft (`014f29c` — ein Commit, acht Pfade, <!-- d-check:status-provenance -->
Gegenstand des Reports) trägt **zwei** Dateien ohne Test-Endung: die Sensor-Doku
(`harness/sensors/coverage-gate.md`, in §3 des Slice-Plans vorgesehen) und
`.dockerignore` — die zweite als Abweichung, die der Implementer selbst
offenlegte. Der Review hat sie gemessen und als F-1 (HIGH) geführt: sie ist
**notwendig** (ohne die Ausnahme fällt der Test im Gate real um — eigener Bau
der `coverage`-Stufe mit der Parent-`.dockerignore`, `Rollout-Datei … nicht
lesbar: no such file or directory`, Bau-EC **1**; übernommene Messung aus dem
Review-Lauf zu `slice-093` <!-- d-check:status-provenance -->),
**minimal** (eine Zeile, keine Muster-/Verzeichnis-Ausnahme),
**image-neutral** (Digest und Binär-sha256 unverändert) und **mit Präzedenz**
(`BEO-PGC/coverage-stage-dockerignore-blockiert-tooling` beschreibt denselben
Fall für `tools/coverage-gate.sh`, 1×). Sein Verdikt: die Lockerung ist
legitim, aber undokumentiert — das Artefakt gehört der Architect-Rolle.

**(2) Warum die Ausnahme notwendig ist — der Build-Kontext ist die
Eingabe-Grenze der netzlosen Tests.** Die `coverage`-Stufe ist `FROM deps AS
coverage` mit `COPY . .` (`Dockerfile:51-63`); der Build-Kontext ist eine
geschlossene Liste (`.dockerignore`): `*` mit Negationen für `cmd/`,
`internal/`, `go.mod`, `go.sum`, `tools/coverage-gate.sh` und
`tools/schema/nacharbeit-roles.sql`. Der Test
`internal/bootstrap/roles_rollout_file_internal_test.go` liest die **reale**
Rollout-Datei über `filepath.Join(filepath.Dir(runtime.Caller(0)), …)` — ein
Pfad, der relativ zum Quelltext aufgelöst wird und damit an den **Kontext**
gebunden ist, nicht an den Arbeitsbaum. Steht die Datei nicht im Kontext, ist
sie in `/src` nicht vorhanden, der Leseversuch schlägt fehl und **der Bau der
Stufe** endet rot. Seit der Test-Arbeit dieser Welle ist der Build-Kontext
damit nicht mehr nur ein Bau-Detail: er ist die Grenze, welche realen Artefakte
ein netzloser Test überhaupt binden kann. Der Testkopf dieses Tests nennt die
Grenze selbst — die View-Klassifikation in `tools/schema/schema.yaml` liegt
„**nicht** im Build-Kontext der `coverage`-Stufe" und ist ihm deshalb
verschlossen.

**(3) Die Größe der Ausnahme — eine Datei, nicht ein Baum.** *Abgeleitet* aus
dem `.dockerignore`-Muster, Methode `git ls-tree -r 014f29c -- cmd internal`:
**166** Dateien unter `cmd/` und `internal/`, dazu `go.mod`, `go.sum` und die
zwei negierten Dateien unter `tools/` ergibt **170** Kontext-Einträge; dieselbe
Rechnung mit der Parent-`.dockerignore` (eine Negation weniger) ergibt **169**.
Der Review-Lauf hat bau-seitig **169 → 170** gemessen (Lauf `slice-093` <!-- d-check:status-provenance -->);
Ableitung und Messung stimmen überein.

**(4) Warum das Image unberührt ist — strukturell, und mit Beleg.** `make image`
baut ohne `--target`, also die **letzte** Stufe (`Makefile:27`); `runtime` ist
`FROM gcr.io/distroless/static-debian12:nonroot@sha256:afa5…` plus
`COPY --from=build /out/pg-change-feed` (`Dockerfile:82-85`). Die
`coverage`-Stufe liegt **nicht** im Abhängigkeits-Graphen dieser Stufe, und
keine Kontext-Datei erreicht das Image: die Negationen wirken auf den Kontext,
aus dem `build`/`coverage` kopieren. Gemessen: der Rebuild dieses Zugs ergab
den Digest `sha256:84bdca56…` — identisch mit dem in `harness/image-hash.txt`
geführten Beleg (Wert selbst gelesen; die Datei trägt ihn seit Commit `1abc7c5`)
— und denselben Binär-sha256 `bbbf7135…` (übernommene Messung, Review-Lauf zu
`slice-093` <!-- d-check:status-provenance -->). Die Pflicht dazu besteht
bereits: „Ein Zug, der Build-Kontext-Dateien ändert, läuft `make image` vor
seiner Closure" (`harness/README.md` §Werkzeuge,
[`ADR-0044`](0044-image-beleg-semantik.md) Punkt 3).

**(5) Was die Klausel schützt — und wo ihr Mittel über das Ziel hinausschießt.**
Geschützt ist: **keine Zeile Produktionscode ändert sich, um Coverage zu
gewinnen.** Das Mittel, mit dem die Klausel das prüft — „jede geänderte Datei
endet auf `_test.go`" — ist **strenger** als dieser Zweck und trifft mit der
`.dockerignore`-Zeile genau den Fall, den der Zweck nicht meint: eine
Bau-Kontext-Plumbing-Zeile, die keine Produktionszeile berührt und die Zahl
nicht bewegt. Das Mittel hat dort zusätzlich **Kosten**: die Welle will Tests,
die an ein reales Artefakt gebunden sind (LP2 ist ihr teuerster und wertvollster
Lieferpunkt), und die Datei-Endung verbietet genau die Änderung, die solche
Tests überhaupt lauffähig macht.

**(6) Die Träger der Klausel — geprüft, nicht erinnert.**
`grep -rn "test-only|Test-only|ausschließlich \`\*_test.go\`" --include=*.md .`
ohne `.harness/` über den Stand dieses Zugs: der normative Träger sind **zwei
Zeilen** — `ADR-0082:336-337` (§Konsequenzen) und `ADR-0082:349`
(§Fitness Function); diese ADR selbst ist als die ablösende Stelle nicht
mitgezählt. Alle weiteren Treffer sind Records (`docs/reviews/**`) und
ein Register-Eintrag — beide werden nicht umgeschrieben. Es gibt **keinen**
dritten stehenden Träger; diese ADR löst genau die zwei Zeilen ab.

## Entscheidung

Wir wählen: **die Folgepflicht wird auf ihren Zweck geschärft; die legitime
Ausnahme wird eine benannte Klasse mit vier Merkmalen; ihr Träger bleibt die
`.dockerignore`-Negation.** Fünf Festlegungen:

**1. Abgelöst wird eine Folgepflicht — nicht das Schnittmaß.** `ADR-0082`
§Entscheidung wird ausdrücklich **bestätigt**: das Schnittmaß ist der Tail, der
Composition Root bleibt im Gegenstand ohne Naht, die Endstufe 80 % und der
Messgegenstand bleiben unverändert. Die geltende Stufe steht unverändert auf
`THRESHOLD ?= 70` (`harness/mk/coverage.mk:15`, Stand dieses Zugs); die Rampe
auf 80 % hebt die Wellen-Closure. Abgelöst sind allein die zwei Zeilen aus
Kontext (6).

**2. Die geschärfte Folgepflicht — Zweck statt Datei-Endung.** Sie lautet
künftig:

> Folgepflicht (Implementer-Zug, je Test-Slice dieser Welle): **kein Griff in
> den Messgegenstand.** Keine Datei unter `internal/**` oder `cmd/**`, die
> nicht auf `_test.go` endet, ändert sich — keine Zeile Produktionscode ändert
> sich, um Coverage zu gewinnen. Unverändert bleiben außerdem die **Messfläche**
> (die Paketliste des `coverage`-Filters, `-coverpkg`/`-covermode`, die Zählung
> in `tools/coverage-gate.sh`) und die **Schwelle** (`THRESHOLD`,
> `AGENTS.md` §3.6). Änderungen **außerhalb** dieser zwei Bäume und dieser drei
> Größen sind zulässig, wenn sie die vier Merkmale der Build-Kontext-Ausnahme
> (Festlegung 3) tragen oder eine durch die Arbeit **überholte Aussage**
> nachziehen (`BEO-PGC/arbeit-ueberholt-stehenden-traeger`) — und sie nennen
> ihre Begründung im Diff.

**3. Die Build-Kontext-Ausnahme — vier Merkmale.** Eine Änderung an
`.dockerignore`, die dem Build-Kontext der `coverage`-Stufe genau eine Datei
hinzufügt, trägt, wenn sie alle vier erfüllt:

- **(i) Gelesen, nicht gebaut.** Die aufgenommene Datei ist **Eingabe** eines
  Tests (oder eines Skripts der Stufe) und wird als reales Artefakt gelesen; sie
  wird nicht ins Produkt kompiliert, und für sie ändert sich keine
  Produktionszeile.
- **(ii) Genau eine Datei.** Die Ausnahme benennt **eine** Datei — kein
  Verzeichnis, kein Muster, keine Wildcard. Die Kontext-Größe wächst um genau
  einen Eintrag.
- **(iii) Der Leser steht dabei.** Der Kommentar über der Ausnahme nennt den
  Test oder das Skript, das die Datei liest (Kopplung, `AGENTS.md` §3.7) — eine
  veraltete Ausnahme ist damit ohne Bau sichtbar. Beide Einträge des Stands
  erfüllen das: `tools/coverage-gate.sh` ← die Stufe selbst,
  `tools/schema/nacharbeit-roles.sql` ←
  `roles_rollout_file_internal_test.go`.
- **(iv) Das Image bleibt unberührt — mit Beleg.** Der Zug, der eine Ausnahme
  setzt oder ändert, läuft `make image` vor seiner Closure und führt den
  Vergleich in der Form aus, die
  [`ADR-0044`](0044-image-beleg-semantik.md) vorschreibt: Digest innerhalb
  desselben Builders, und bei einem Inhalts-Streit der sha256 des extrahierten
  Binaries; der Digest-Commit entfällt bei unverändertem Digest.

Die Klasse ist damit **enger als ihr Name nahelegt**: eine Ausnahme ist nicht
„`.dockerignore` darf wachsen", sondern ein Eintrag, der seinen Leser nennt und
seine Image-Neutralität belegt.

**4. Der Träger bleibt die `.dockerignore`-Negation.** Der Weg über den
Kontext ist gewählt: er ist die kleinste Änderung, die den Test real lauffähig
macht, er berührt **keinen** Messmechanismus, und die Image-Neutralität ist
strukturell (Kontext (4)) **und** gemessen. Der Bind-Mount des Arbeitsbaums in
die `coverage`-Stufe (Muster `make proto-generate`) ist die **größere** Änderung
und wird mit Trigger (a)/(d) vertagt — die Abwägung steht vollständig in
§Verglichene Alternativen Option D, samt der **einen** Gegenstärke: er
beseitigt die Klasse ganz, weil er den Kontext-Filter als Eingabe-Grenze
auflöst.

**5. Was diese ADR nicht tut — und warum.**

- **Kein Sensor.** Die geschärfte Regel ist **wellen-skopiert** („je Test-Slice
  dieser Welle"), die verfügbaren Gate-Ranges sind es nicht: ein Standing-Gate
  über die letzten Commits enthielte Commits fremder Züge, und ein Gate über den
  Wellen-Diff braucht ein Range-Wissen, das kein bestehendes Target hat. Ein
  Sensor an dieser Stelle wäre eine neue Gate-Mechanik mit eigener Beweislast —
  und ihr Gegenstand („ist die Änderung ein Griff in den Messgegenstand?") ist
  ein **Zweck**-Urteil, genau wie in `ADR-0082` („Disziplin, Review-Gegenstand;
  kein Sensor"). Der Wächter bleibt der Review.
- **Keine Messfläche, kein Messmechanismus, keine Schwelle.** Der Paketfilter
  der Stufe, `-coverpkg`/`-covermode`, die Zählung in
  `tools/coverage-gate.sh`, die `--no-cache-filter`-Sicherung der Stufe und
  `THRESHOLD` bleiben unverändert; [`ADR-0071`](0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
  und `AGENTS.md` §3.6 sind unberührt.
- **Kein Produktionscode, kein Image-Eingriff.** Diese ADR ist reiner Text an
  der Entscheidungslage; der ausgelieferte Bau-Kontext folgt ihr bereits.

### Warum die Ausnahme keine zweite Quelle für die Zusage ist

Der Einwand liegt nahe: eine Ausnahme neben einer Regel ist eine zweite Stelle,
an der „wo gilt die Zusage nicht" steht ([`ADR-0073`](0073-zitat-korrektur-an-immutablen-dokumenten.md)
§Verglichene Alternativen B verwirft genau das für die Doku). Hier trifft er
**nicht**: die Zusage, die das Repo ausschließt, ist „keine Zeile
Produktionscode um Coverage zu gewinnen" — und die Ausnahme berührt diese Zusage
nicht, weil sie kein Produktionscode ist. Was sie berührt, ist das **Mittel**,
mit dem die Zusage geprüft wurde; ersetzt wird also das Mittel, nicht die
Zusage. Die vier Merkmale stehen darum **nicht** als Ausnahmeliste neben der
Regel, sondern sind die Form, in der die **Regel selbst** ihre Grenze nennt.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun; die Abweichung bleibt die Notiz des Implementers, die zwei Zeilen bleiben stehen | kein Eingriff an einer `Accepted`-ADR | die zwei Zeilen bleiben im Wortlaut falsch; jeder folgende Test-Slice der Welle legt die Ausnahme neu gegen sie aus, und der nächste Review findet dieselbe HIGH-Klasse — genau der Zustand, den [`ADR-0067`](0067-capture-publish-einbindung-fitness-function-korrektur.md) Option A für dieselbe Klasse („zitiert unrevidiert") verworfen hat |
| B — die Folgepflicht **in-place** schärfen, ohne Folge-ADR | kleinster Text-Eingriff an einer Stelle | `AGENTS.md` §3.5 verbietet es: §Konsequenzen und die Fitness-Function-Regeln sind ausdrücklich **keine** Zitat-Korrektur ([`ADR-0073`](0073-zitat-korrektur-an-immutablen-dokumenten.md) §Entscheidung Punkt 1) — ein in-place-Eingriff wäre der stille Kanal, den dieselbe ADR eng fasst |
| C — die Lockerung als **Zusatzregel** führen, die `ADR-0082` nur `Schärft` (das Feld auf eine ADR gerichtet) | die alte Klausel bleibt unangetastet, kein Supersede | **das Feld trägt diese Richtung nicht**: `Schärft:` weist auf die Spec-Straten (`SPEC-*`, `ARC-*`, `LH-*.a`) — im ganzen ADR-Bestand zeigt keine ADR damit auf eine andere ADR. Und solange die Klausel stehenbleibt, behauptet sie etwas, das die Welle real gebrochen hat; die Auslegung bliebe bei jedem Lauf |
| D — **den Messmechanismus ändern:** den Arbeitsbaum in die `coverage`-Stufe bind-mounten (Muster `make proto-generate`), statt den Kontext zu erweitern | der Kontext-Filter verliert seine Rolle als Test-Sichtgrenze: **jede** Datei des Repos wird für Tests lesbar, die Ausnahmen dieser Klasse entfallen ganz; die Stufe verliert zugleich ihren `--no-cache-filter`-Wächter, weil ein Lauf keinen Layer-Cache kennt | verschiebt den **Messmechanismus** von Bau-Zeit auf Lauf-Zeit und damit ein Gate (`make coverage-gate` würde ein `docker run` über den Arbeitsbaum statt ein `docker build` über den gefilterten Kontext): die Zählbasis-Aussage („das Profil der `coverage`-Stufe des gepinnten Images", `harness/sensors/coverage-gate.md` §Zählbasis; `ADR-0082` §Fitness Function Zeile 3) wäre neu zu belegen, der Lauf hinge an der Laufzeit-Umgebung des Aufrufers (unversionierter Baum samt `.git`), und die heute **benannte** Sichtgrenze fiele ersatzlos — eine eigene Entscheidung mit eigener Beweislast, nicht die Antwort auf eine Zeile |
| **E — die Folgepflicht auf ihren Zweck umformulieren und die Ausnahme als Klasse mit vier Merkmalen fassen (Festlegungen 2–4), `ADR-0082`s Schnittmaß bestätigt (gewählt)** | die Klausel wird wieder wahr und trägt dieselbe Wirkung („keine Zeile Produktionscode um Coverage zu gewinnen" — als „keine Änderung unter `internal/**`/`cmd/**` außer `_test.go`" prüfbar); die legitime Ausnahme ist **benannt** statt still geduldet; Messfläche, Messmechanismus und Schwelle bleiben unberührt, das Image bleibt unberührt (strukturell **und** gemessen); kein Produktcode, kein neuer Sensor | zwei ADRs müssen zusammengelesen werden (`0082` + `0085`); die Klasse ist eine **Disziplin** mit vier Merkmalen — ihr Wächter ist der Review, kein Sensor —, und sie ist nicht beliebig erweiterbar (der dritte Eintrag löst Trigger (a) aus) |

**Fazit:** E. B und C scheitern an der **Form** (§3.5; `Schärft:`-Richtung), A
lässt falschen Text stehen, D ist an derselben Stelle die größere Änderung und
für den Anlass nicht nötig — sie bleibt als Trigger (a)/(d) ausdrücklich
vorgemerkt, weil sie die Klasse **ganz** beseitigen würde.

## Konsequenzen

- Positiv: Die zwei abgelösten Zeilen sind wieder wahr, und die Wirkung, die sie
  tragen sollten, ist unverändert: die zwei Bäume des Messgegenstands und die
  drei Messgrößen sind benannt und unverändert.
- Positiv: Die Ausnahme ist **benannt** statt geduldet. Merkmal (iii) macht eine
  veraltete Ausnahme ohne Bau sichtbar; Merkmal (iv) hängt sie an eine
  **bestehende** Regel (`harness/README.md` §Werkzeuge,
  [`ADR-0044`](0044-image-beleg-semantik.md) Punkt 3) statt an eine neue Pflicht.
- Positiv: Kein Produktionscode, kein Messmechanismus, keine Schwelle, kein
  Image-Eingriff; §3.6 und die Out-of-Scope-Liste der Welle bleiben unberührt.
  Die Sichtgrenze des Kontexts bleibt bestehen — und mit ihr die benannte Grenze
  im Testkopf, die ohne sie verschwände.
- Negativ mit Grenze: Die geschärfte Folgepflicht ist im Wortlaut **weiter** als
  die alte (Dateien außerhalb `internal/**`/`cmd/**` sind nicht mehr generell
  verboten) und an der tragenden Stelle **enger** (die zwei Bäume und die drei
  Messgrößen sind namentlich gesetzt). Das ist die Setzung: die alte Fassung war
  ein **Mittel**-Verbot, die neue ein **Zweck**-Verbot mit benannter Ausnahme.
  Ihr Wächter ist der Review.
- Negativ mit Grenze: Die Ausnahme-Klasse ist **nicht mechanisch gewächtert**.
  Eine fehlende Ausnahme zeigt sich erst am realen roten Bau (so gemessen,
  `BEO-PGC/coverage-stage-dockerignore-blockiert-tooling`), und eine zu weite
  (Verzeichnis statt Datei) ist am Muster erkennbar, aber ungeprüft. Merkmal
  (ii) ist damit eine Form-Zusage an den Review, kein Sensor. Nach diesem Zug
  steht die Klasse bei **2×**; die Schwelle 3× löst Trigger (a) aus.
- Folgepflicht (Implementer-Zug): `harness/sensors/coverage-gate.md` §Grenze
  bekommt einen Punkt, der die Klasse und ihre vier Merkmale benennt — dort
  liest sie der Implementer eines Coverage-Slice. Der Kopf-Kommentar in
  `.dockerignore` beschreibt die zwei Ausnahmen als „Stage"-gebunden; sie sind
  **kontextweit** (sie gelten jeder Stufe, die den Baum kopiert). Die tragende
  Aussage („das Runtime-Image bleibt unberührt") ist wahr und bleibt; die
  Formulierung wird auf den Ist-Zustand gezogen und nennt je Eintrag seinen
  Leser (§3.7).
- Folgepflicht (Planner-Zug): Die geschärfte Folgepflicht gehört in die
  Slice-Pläne der Welle — dieselbe Pflicht, die `ADR-0082` für die Welle bereits
  formuliert; sie ist die **Ersetzung** des dort zitierten Satzes, keine
  Zusatzbedingung. Der **Erinnerungs-Träger** für Trigger (a) ist der
  Register-Eintrag `BEO-PGC/coverage-stage-dockerignore-blockiert-tooling`
  (nach diesem Zug 2×): der Gegenstand ist eine Beobachtung, ihr Zähler steht
  dort — ein Erinnerungs-Slice wäre eine **zweite Quelle** für denselben Zustand
  (Modul 6).

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| `git diff --name-only` des Slice-Commits | **Kein Griff in den Messgegenstand:** keine geänderte Datei unter `internal/**` oder `cmd/**`, die nicht auf `_test.go` endet. Die zulässigen Änderungen **außerhalb** dieser zwei Bäume sind die der Build-Kontext-Ausnahmen (vier Merkmale) und der überholten Aussagen — sie sind **nicht** durch eine Endung charakterisiert und bleiben Urteil | — (Disziplin, Review-Gegenstand; kein Sensor — Begründung: Festlegung 5) |
| `make image` + `harness/image-hash.txt` | **Das Image bleibt unberührt:** ein Zug, der eine Build-Kontext-Ausnahme setzt oder ändert, läuft `make image` vor seiner Closure; Digest innerhalb desselben Builders und sha256 des extrahierten Binaries sind der Beleg ([`ADR-0044`](0044-image-beleg-semantik.md)) | `make image` (kein Gate) |
| `.dockerignore` (eine Zeile je Ausnahme) | **Genau eine Datei je Ausnahme** (kein Verzeichnis, kein Muster) und **der Leser steht dabei** — beides am Diff prüfbar | — (Form, Review-Gegenstand) |

Die **erste** und die **dritte** Fitness-Function-Zeile aus `ADR-0082`
(Gesamt-Coverage über der netzlos prüfbaren Fläche gegen `THRESHOLD`;
Zählbasis-Nachweis des Laufs) bleiben unverändert in Kraft; ihre **zweite**
Zeile (Test-only) ist durch die erste Zeile dieser Tabelle ersetzt.

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

Beobachtbare Trigger:

**(a)** Die Ausnahme-Klasse erreicht **3×** — der Register-Eintrag
`BEO-PGC/coverage-stage-dockerignore-blockiert-tooling` steht nach diesem Zug
bei 2× (der Beleg zu diesem Auftreten entsteht bei der Slice-Closure des
Vorgangs) —, dann ist der **Mechanismus** neu zu entscheiden: der Lese-Schritt
des Registers und dieser Trigger fallen zusammen, und die Antwort ist Option D
(Bind-Mount) oder ein eigener Kontext für die Messung. Ein Erinnerungs-Slice
tritt an ihre Stelle nicht (Konsequenzen).

**(b)** Ein Slice dieser Welle ändert eine Datei unter `internal/**` oder
`cmd/**` außerhalb von Tests **oder** eine der drei Messgrößen — dann trägt die
geschärfte Folgepflicht nicht, und die Antwort ist ein Sensor oder eine Naht,
nicht eine weitere Formulierung.

**(c)** `make image` ergibt für eine Ausnahme einen anderen **Binär-sha256**
(nicht nur einen anderen Digest) — dann ist die Ausnahme nicht image-neutral,
Merkmal (iv) ist verletzt, und der Träger ist neu zu bewerten.

**(d)** Die `coverage`-Stufe verliert ihre Bau-Zeit-Form aus einem anderen Grund
(etwa weil sie eine Eingabe braucht, die kein Einzeldatei-Eintrag tragen kann —
ein Verzeichnis, eine zweite Wurzel) — dann ist Option D nicht mehr die größere,
sondern die naheliegende Änderung, und die Zählbasis-Aussagen sind mit ihr
nachzuziehen. Das ist derselbe Ausgang wie (a) mit einem anderen Anlass.

Sonst permanent: die zwei abgelösten Klauseln gelten in der Fassung dieser ADR
weiter, unabhängig davon, wie viele Test-Slices die Welle noch schneidet.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-16 | Accepted — Anlass: Review zu `slice-093` <!-- d-check:status-provenance --> (F-1, HIGH) belegt, dass eine Zeile in `.dockerignore` die Folgepflicht „test-only" aus `ADR-0082` im Wortlaut bricht, und gibt das Verdikt „Lockerung legitim, aber undokumentiert" an die Architect-Rolle (Modul 8 §Konflikt-Pfad, drittes Verdikt). Entscheidet: die Folgepflicht wird auf ihren Zweck geschärft (keine Änderung unter `internal/**`/`cmd/**` außer `_test.go`; Messfläche, Messmechanismus und Schwelle unverändert), die Ausnahme wird eine benannte Klasse mit vier Merkmalen, ihr Träger bleibt die `.dockerignore`-Negation; der Bind-Mount des Arbeitsbaums (Option D) ist mit Trigger (a)/(d) vertagt. `ADR-0082`s §Entscheidung und alle übrigen Klauseln bestätigt | Review zu `slice-093` F-1 <!-- d-check:status-provenance --> (Bau-EC 1 ohne die Ausnahme; Kontext 169 → 170; Digest und Binär-sha256 unverändert) · eigene Ableitung des Kontext-Umfangs aus dem Muster (166 + 2 + 2 = 170; `git ls-tree -r 014f29c`) · `grep` über die Träger der Klausel (Kontext (6)) | <!-- d-check:status-provenance -->
| 2026-09-18 | Zitat-Korrektur — `docs/reviews/**`-Pfade durch Kennung ersetzt (`ADR-0073`) | PENDING_COMMIT |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0085` — eine **Zitat-Korrektur** ausgenommen
([`ADR-0073`](0073-zitat-korrektur-an-immutablen-dokumenten.md))
(Baseline-Regelwerk `modul-04-adrs.md` §Hard Rule für Accepted-ADRs).
