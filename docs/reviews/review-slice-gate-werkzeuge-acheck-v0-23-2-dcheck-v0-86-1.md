# Review-Report: slice-gate-werkzeuge-acheck-v0-23-2-dcheck-v0-86-1 — 2026-10-10

**Review-Art:** Code — gegen Plan, Entscheidungen und Hard Rules (nicht gegen die DoD; die prüft der Verifier)

**Gegenstand:** Diff `d496d0ee..606b7560` (5 Commits, nicht gepusht): `05c14aac` (Reichweiten-Abschnitt und `SPEC-040`-Tabelle im Plan), `0845d926` (Pin a-check), `17611a23` (Pin d-check), `bf8e1762` (Pin-Belege, Suchlauf am Stand `diff`, Liefer-Punkte), `606b7560` (Gate-Lauf belegt).

**Skill:** `.harness/skills/reviewer.md` @ 606b7560
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-10

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis; die `<Platzhalter>` darin sind Formbeispiele)*. Dieser
> Report friert ein; was er zitiert, bewegt sich
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

- Slice-Plan `slice-gate-werkzeuge-acheck-v0-23-2-dcheck-v0-86-1` (Stand `606b7560`)
- [`ADR-0051`](../plan/adr/0051-cicd-pipeline-github-actions.md) (Pin-Inventar P6/P7, Entscheidung 7), [`ADR-0041`](../plan/adr/0041-a-check-maschinenform-architekturpruefung.md), [`ADR-0045`](../plan/adr/0045-commit-traceability-standing-gate.md), [`ADR-0163`](../plan/adr/0163-schaerft-spec-040-provenance-marker-und-exempt-paths.md) (Entscheidung 1, Trigger (a) bis (c)), [`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md)
- Bump-Regel `harness/targets/pin-stale.md` §Bump eines Gate-Werkzeugs: Reichweite vor dem Pin-Commit
- [`SPEC-040`](../../spec/pflichtenheft.md#7-festlegungen-der-harness-werkzeuge) und `SPEC-039` (Kanten des Slice)
- [`AGENTS.md`](../../AGENTS.md) §3.1, §3.5, §3.9, §3.12, §3.13
- Changelogs (`CHANGELOG.md`) von a-check und d-check, selbst per `gh api` geholt; `gh release list` beider Repos
- Keine `LH-*`-ID berührt (Harness-Werkzeuge)

**Eigene Messungen dieses Laufs** (Klon von `d496d0ee` im Scratchpad, `git clone --no-hardlinks`, vor jeder Mutation `pwd` im Klon geprüft; Mutationen per Redirect in den Klon, Rücknahme per `git checkout`/`git reset --hard`; Arbeitsbaum des Repos unberührt):

- Digests: `docker buildx imagetools inspect ghcr.io/pt9912/a-check:v0.23.2` → `sha256:2368f7b3…f422`; `ghcr.io/pt9912/d-check:v0.86.1` → `sha256:3e0b9779…ce0e`; beide `application/vnd.oci.image.index.v1+json`. Gleich den Pins in `a-check.mk`/`d-check.mk` und dem Plan.
- `make pin-stale-acheck`: Exit 0, `OK … == sha256:2368f7b3…`; `make pin-stale-dcheck`: Exit 0, `OK DCHECK_DIGEST … == sha256:3e0b9779…` und `OK DCHECK_IMAGE Tag-Frische (v0.86.1) == neuester Release`.
- `git show --stat`: `0845d926` eine Datei (`a-check.mk`, 1 Zeile), `17611a23` eine Datei (`d-check.mk`, 2 Zeilen: Tag und Digest). `git log`: `05c14aac` (Messung) liegt vor beiden Pin-Commits; alle fünf Messages tragen eine `ADR-*`-Kennung.
- Alt gegen neu über `make` mit Kommandozeilen-Digest: `make docs-check` (beide Exit 0, `1848 Datei(en) geprüft, 0 Befund(e)`, Ausgabe byte-gleich); `make a-check` (beide Exit 0, byte-gleich); `make doc-immutable` mit `RANGE=d496d0ee..d496d0ee` (Exit 2, Leerfall), `dc04087e..d496d0ee` (Exit 0), `3f88e0c2~1..3f88e0c2` (Exit 2, 4 Befunde), `8e00e831~1..8e00e831` (Exit 2, 1 Befund) — je alt/neu Exit gleich, nach `sort` byte-gleich.
- Mutationen, die rot werden müssen: `internal/domain/model/zz_mut.go` mit Import von `…/adapters/driven/grpcstream` → `make a-check` alt und neu Exit 2, `wrong-direction: domain -> adapters`, `gesamt: 1 Befund(e)`, byte-gleich; leerer Commit `chore: ohne Kennung` → `make doc-commits RANGE=HEAD~1..HEAD` alt und neu Exit 2, 1 × `commit-untraceable`, gleich; `allow-supersede-lineage` aus `.d-check.yml` entfernt → alt und neu Exit 1, 16 × `matrix-inactive`; zwei Home-relative Pfade mit Schrägstrich in Inline-Code → neu Exit 1, 2 × `hostpath-forbidden`; drei fehlende Linkziele → 5 Läufe `make docs-check`, je Exit 2, letzte Zeile `make: *** [d-check.mk:20: docs-check] Fehler 1`, Summenzeile in drei Läufen vor, in zwei Läufen nach den Befundzeilen (Reihenfolge nicht festgelegt, wie im Plan).
- Modul `reviews` direkt: v0.82.0 und v0.84.0 → Exit 0, `0 Befund(e)`; v0.85.0 → Exit 1, `132 Befund(e)`, alle `review-missing`; v0.86.1 → Exit 1, 132 (wie im Plan).
- Nicht im Plan gelistete Einzel-Ziele, alt gegen neu: `make doc-trace` (Exit 0, 90 Zeilen), `make doc-complete` (Exit 0, 90 Zeilen), `make doc-doctor` (Exit 0, 2 Zeilen), `make doc-repair` (Exit 0, 1 Zeile) — jeweils byte-gleich.
- Changelog-Gegenlesung (a-check 0.21.0 bis 0.23.2, d-check 0.83.0 bis 0.86.1 samt 0.84.0-rc.1): keine Änderung an einem in diesem Repo aktiven Modul (`links`, `anchors`, `ids`, `matrix`, `versions`, `structure`, `hostpaths`, `tracked` in `.d-check.yml`; `commits`, `vcs`, `planning`, `targets` nur in Einzel-Zielen) und keine Erweiterung der a-check-Regeln (`shapes` in `.a-check.yml` nicht gesetzt). Die zwei „nicht rein additiven“ Einträge betreffen nur `reviews`. `git grep -e '--manual\|--suggest-config'` über `Makefile`, `*.mk`, `tools`, `.github`, `.claude` → kein Treffer.
- Suchlauf §3.13 am Baum (alte/neue Digests, Tags `v0.82.0`/`v0.20.0`, `0.8x.y`/`0.2x.y`, `pt9912/(a|d)-check`) über `harness`, `README.md`, `docs/user`, `docs/maintainer`, `AGENTS.md`, `.claude`, `.harness/skills`, `tools`, `Makefile`, `*.mk`, `.github`: kein weiterer Stand-Träger außer den Pins; die Nennungen in `Accepted`-ADRs sind im Plan als Bestand begründet.
- `make suchlauf-nachmessen PLAN=<Plan>`: `18 Zeilen stimmen`, Exit 0.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Die latente Folge der Erweiterung am Modul `reviews` (bei Aktivierung 132 × `review-missing`; gemessen: v0.84.0 0 Befunde, v0.85.0 132) hat weder einen Ausgang noch eine Adresse. §6 legt für das Risiko „nicht rein additive Änderung“ den erwarteten Ausgang *entfallen* fest; das deckt „läuft heute nicht mit“, nicht „meldet bei Aktivierung 132“. Die Zahl steht nur im Abschnitt „Reichweite …“ des Plans, der mit dem Slice nach `done/` und später ins Archiv geht; `harness/README.md` nennt `reviews` (Ziel `doc-reviews`) als Soll-Kante, das Ziel gibt es nicht. | Baseline `v6.16.0` · `regelwerk/modul-05-planning-harness.md` §Offene Risiken werden bei Closure aufgelöst; [`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md) | Slice-Plan · „entfallen, mit diesem Grund.“ (§6, erstes Risiko); „Das ist die angekündigte „nicht rein additive“ Änderung“ (§3) | ja — `docker run … d-check@sha256:3e0b9779… --enable reviews` gegen den Klon: Exit 1, 132 Befunde (Existenz belegt, Adresse nicht mechanisch prüfbar) | Aufschub ohne Adresse |
| F-2 | LOW | Die Spaltenüberschrift der `SPEC-040`-Tabelle lautet „v0.86.1 (gleich v0.82.0)“, der Vergleichstext behauptet für jeden Lauf gleiche Ausgabe beider Stände; die Zeile „Basislauf mit `--enable reviews`“ in derselben Tabelle zeigt Exit 1 und 132 Befunde, v0.82.0 liefert dort Exit 0. Die 67 gezählten Läufe schließen diese Zeile aus, die Tabelle sagt es nicht. | `AGENTS.md` §3.12 Instanz A; Reviewer-Skill „Nachzug widerspricht dem Nachbarn im selben Träger“ | Slice-Plan · „v0.86.1 (gleich v0.82.0)“ und „Basislauf mit `--enable reviews`“ (§3) | nein — Lese-Handlung | Nachzug widerspricht Nachbar im selben Träger |
| F-3 | INFO | Der Plan schreibt die „nicht rein additive“ Änderung mit den 132 Befunden `0.85.0/0.86.0` zu. Die Nachmessung zeigt sie bereits mit v0.85.0 (v0.84.0: 0, v0.85.0: 132, v0.86.1: 132). Die Zuordnung zu 0.86.0 ist weder gemessen noch als *hergeleitet* gekennzeichnet; das Urteil „keine Reichweiten-Änderung eines Ziels“ bleibt davon unberührt. | [`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md) | Slice-Plan · „von 0.85.0/0.86.0; sie erreicht kein Ziel dieses Repos“ (§3) | ja — `d-check:v0.85.0 --enable reviews` gegen den Klon | Ursprung einer Aussage nicht gekennzeichnet |
| F-4 | INFO | Liefer-Punkt 1 listet die Einzel-Ziele `doc-trace`, `doc-complete`, `doc-doctor` und `doc-repair` nicht. Alle vier laufen mit dem neuen Digest und sind im Review alt gegen neu byte-gleich (Exit 0); `doc-complete` ist laut `d-check.mk` ein „Vollständigkeits-Gate“. Die Aufzählung der Ziele im Plan ist nicht erschöpfend gegen `d-check.mk`. | `harness/targets/pin-stale.md` §Bump eines Gate-Werkzeugs | Slice-Plan · „dazu die Einzel-Ziele“ (§2, Liefer-Punkt 1) | ja — `make doc-complete DCHECK_DIGEST=…` | Messliste unvollständig gegen Zielliste |
| F-5 | INFO | Für `make doc-planning` und `make doc-targets` belegt der Plan nur Gleichheit am Bestand und an Mutationen, die in keinem Stand rot wurden; die Reichweite dieser zwei Ziele ist nicht durch gesehenes Rot gestützt. Der Plan benennt die Grenze selbst; kein Mangel, festgehalten, damit „keine Reichweiten-Änderung“ dort als Gleichheitsbefund und nicht als Reichweitenbeleg gelesen wird. | [`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md) | Slice-Plan · „damit nicht durch ein gesehenes Rot belegt“ (§3) | nein | Zusage ohne Rot-Beleg (benannt) |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Bump-Regel: Urteil je übersprungener Version (a-check 0.21.0, 0.22.0, 0.23.0, 0.23.1, 0.23.2; d-check 0.83.0, 0.84.0-rc.1, 0.84.0, 0.85.0, 0.86.0, 0.86.1) | geprüft, ohne Befund — jede Version trägt „keine Reichweiten-Änderung“ mit Grund und Artefakt; gegen die Changelogs keine übersehene Erweiterung eines aktiven Moduls |
| Reihenfolge der Commits | geprüft, ohne Befund — Messung `05c14aac` vor `0845d926` und `17611a23`; keine Erweiterung, daher kein Architect-Zug nötig |
| Pins `a-check.mk`, `d-check.mk` | geprüft, ohne Befund — Digests gleich Registry-Index, je Pin-Commit eine Datei, `make pin-stale-acheck`/`make pin-stale-dcheck` ohne `DRIFT` |
| Träger-Nachzug (Suchlauf §3.13) | geprüft, ohne Befund — kein weiterer Träger der alten Stände in `harness/`, Workflows, Doku, Skripten; `Accepted`-ADRs bleiben als Bestand (`AGENTS.md` §3.5) |
| `SPEC-040`/`SPEC-039` gegen v0.86.1 | geprüft, ohne Befund — Stichproben (Leerfall, Form-/Umzugs-Commit, `matrix-inactive`, `hostpath-forbidden`, `make`-Exit-Verhalten) reproduzieren die Tabelle; Trigger (a) bis (c) von [`ADR-0163`](../plan/adr/0163-schaerft-spec-040-provenance-marker-und-exempt-paths.md) zu Recht „nicht eingetreten“ (kein Changelog-Eintrag an `matrix`) |
| „hergeleitet“/„übernommen“-Kennzeichnung (`AGENTS.md` §3.12) | geprüft, ohne Befund außer F-3 — Digests als übernommen, Erwartungen und nicht gemessene Sätze als *hergeleitet* gekennzeichnet |
| Suchlauf-Block und Gate-Beleg | geprüft, ohne Befund — `make suchlauf-nachmessen`: 18 Zeilen stimmen; Gate-Lauf-Beleg nennt Stand `bf8e1762` und Exit 0 |
| Docker-only, Suppression, Kommentare | geprüft, ohne Befund — der Diff berührt zwei Pin-Zeilen und den Plan |
| `spec/`, `docs/plan/adr/`, `harness/conventions/` | nicht berührt (`git diff --stat`) — `make doc-immutable` über die ganze Range entspricht der Erwartung im Plan |
| M32c (57 × `matrix-inactive`) und weitere Tabellenzeilen außerhalb der Stichprobe | nicht nachgemessen (eine Mutation des Laufs traf die falsche Konfigurationszeile und zählt nicht); übernommen |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 1 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Aufschub ohne Adresse · Nachzug widerspricht Nachbar im selben Träger · Ursprung einer Aussage nicht gekennzeichnet · Messliste unvollständig gegen Zielliste · Zusage ohne Rot-Beleg (benannt)

## Verdikt

**Merge-blockierend:** nein — F-1 und F-2 sind Textstellen des Plans, die der Planner bei der Closure ohnehin schreibt (F-1: Ausgang des ersten Risikos in §6/§7 als *weiter offen* mit Register-Eintrag oder als *eingetreten* mit Adresse statt *entfallen*; F-2: Spaltenüberschrift). Es gibt keine Rückgabe an den Implementer und keine Code- oder Pin-Änderung; Pins, Reihenfolge und Reichweiten-Urteile tragen. F-1 bleibt bis zur Closure offen und ist dort zu schließen.

**Übergabe:** Findings gehen an den Planner (Closure-Notiz §7, Risiko-Ausgänge §6); die **Finding-Klassen** gehen zusätzlich in die Slice-Closure §7 und von dort in den Zähler. Dieser Report selbst ist ein **Lauf-Beleg**. DoD-/Spec-Konformität prüft der Verifier separat.
