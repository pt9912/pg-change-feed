# Slice release-binary-version-injektion: `--version` zeigt die reale Release-Version

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`.

**Welle:** ohne Welle — die Closure-Bedingung ist identisch mit der DoD
dieses einen Slice, kein *Mehr* gegenüber ihr.

**Bezug:** [`ADR-0051`](../../adr/0051-cicd-pipeline-github-actions.md)
(additive Erweiterung — `docs/user/version.md`/der Git-Tag ist bereits
die Quelle der Wahrheit für die Version, dieser Slice schließt nur die
letzte Lücke: das Binary selbst kannte sie bisher nicht).

**Berührte Spec-Stellen:** — (reine Build-/Betriebsdetail-Änderung ohne
Spec-Stratum-Bezug — `--version` selbst trägt keine eigene `LH-*`-ID,
siehe `LH-FA-SST-003`s Aufzählung in `spec/lastenheft.md`, die
`--version`/`--healthcheck` nur als bestehenden Mechanismus nennt, ohne
Aussage über den gelieferten Versionswert).

**Verantwortlich:** Implementer-Agent (priorisiert 2026-09-19, direkter
Auftrag: "Ja, Slice anlegen und umsetzen").

**Autor:** Planner-Agent, direkt beauftragt. **Datum:** 2026-09-19.

---

## 1. Ziel und Abgrenzung

**Ziel:** `pg-change-feed --version` zeigt die tatsächlich ausgelieferte
Release-Version (`0.1.1`, `0.1.2`, …) statt eines fest einkompilierten,
seit Ersteinführung nie aktualisierten Strings
(`cmd/pg-change-feed/main.go:19`, real beim ersten `v0.1.1`-Release
gefunden: zeigte `0.2.0-verdrahtung`, obwohl das ausgelieferte Image
`0.1.1` war). Die Version wird zur Build-Zeit per Linker-Flag aus
demselben `VERSION`-Parameter injiziert, den `make image`/`release.yml`
bereits für die Image-Tags verwenden — kein zweiter, manuell zu
pflegender Ort.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Ein `git describe`-basierter automatischer Versions-Fallback für den
  lokalen `:dev`-Build — der `:dev`-Pfad bekommt bewusst den festen
  Platzhalter `dev` (siehe §3), ein automatisch aus der Git-Historie
  abgeleiteter Wert wäre eine eigene, nicht angefragte Erweiterung.
- Rückwirkende Korrektur bereits veröffentlichter Images (`v0.1.0`,
  `v0.1.1` zeigen weiterhin den alten String, wenn sie erneut gezogen
  werden) — Images sind unveränderlich nach dem Push, ein Nachtrag wäre
  ein neuer Release, kein technischer Fix an einer bestehenden Version.
- Versionsangaben in `diagnose`/`--healthcheck` — beide melden
  Betriebszustand, keine Versionsinformation; kein bestehender Vertrag
  wird hier erweitert.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] `cmd/pg-change-feed/main.go`: `const version = "0.2.0-verdrahtung"`
      wird zu `var version = "dev"` (Linker-Flag `-X` kann nur Variablen,
      keine Konstanten, überschreiben). Kein Verhaltensunterschied für
      den `--version`-Aufruf selbst.
- [ ] `Dockerfile`s `build`-Stufe bekommt `ARG VERSION=dev` und injiziert
      sie über `-ldflags="-s -w -X main.version=$VERSION"` (ergänzt die
      bestehenden Strip-Flags, ersetzt sie nicht). `Makefile`s
      `image:`-Target übergibt im `VERSION=`-Zweig
      `--build-arg VERSION=$(VERSION)`; der `VERSION`-lose `:dev`-Zweig
      übergibt keinen Wert und nutzt damit automatisch den
      Dockerfile-Default `dev`. Real geprüft: ein Build mit
      `VERSION=0.1.1-test` liefert `pg-change-feed --version` →
      `pg-change-feed 0.1.1-test`; ein Build ohne `VERSION` (unverändert
      `:dev`) liefert `pg-change-feed dev`.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update für <Schnittstelle X> falls öffentlicher Vertrag berührt — entfällt: `--version` selbst bleibt unverändert als Aufrufform, nur sein gelieferter Wert wird korrekt; kein neuer öffentlicher Vertrag.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `cmd/pg-change-feed/main.go` | update | `const version` → `var version = "dev"` (Voraussetzung für `-ldflags -X`). |
| `Dockerfile` (`build`-Stufe) | update | `ARG VERSION=dev`, `-ldflags` injiziert `main.version`. |
| `Makefile` (`image:`-Target) | update | `VERSION=`-Zweig übergibt `--build-arg VERSION=$(VERSION)`. |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): sofort — keine Abhängigkeit von
anderen offenen Slices.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): entfällt aus
  heutiger Sicht — Umfang ist eine Variable, ein Dockerfile-Arg, eine
  Makefile-Zeile.
- `in-progress` → `open` (blockiert — Carveout?): keine bekannte
  Blockade-Bedingung.

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + Closure-Notiz mit Lerneintrag
geschrieben.

## 6. Risiken und offene Punkte

- Der lokale `:dev`-Build zeigt nach dieser Änderung `pg-change-feed dev`
  statt des bisherigen `0.2.0-verdrahtung` — ein sichtbarer, aber
  gewollter Verhaltensunterschied für jeden, der `--version` gegen ein
  `:dev`-Image aufruft. **Ausgang:** entfallen als Risiko — `dev` ist die
  ehrlichere Aussage für ein nie veröffentlichtes Image, kein
  Informationsverlust gegenüber einem ohnehin schon bedeutungslosen
  Alt-String.
- `AGENTS.md` §3.10 gilt für die geänderte `Dockerfile`/`Makefile`-
  Mechanik (mittelbar über `release.yml`) unverändert — ein realer,
  grüner Post-Push-Lauf mit tatsächlich korrekt gesetztem `--version`-
  Wert bleibt der einzige volle Beleg. **Ausgang:** weiter offen,
  strukturell (derselbe Fall wie
  `BEO-PGC/github-actions-unverifizierbar-lokal`, bereits verkörpert).

## 7. Closure-Notiz

*(wird bei Bearbeitung gefüllt.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Sub-Area „CI/CD-Pipeline"
(`Dockerfile`, `Makefile`) — bereits mehrfach berührt
(`welle-release-pipeline-adr-0051`, `release-multi-arch-image`),
Schwelle erfüllt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
`BEO-PGC/github-actions-unverifizierbar-lokal` (bereits verkörpert,
siehe §6) betrifft diese Sub-Area; kein weiterer Treffer.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
