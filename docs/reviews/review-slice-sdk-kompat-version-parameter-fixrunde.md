# Review-Report: slice-sdk-kompat-version-parameter (Fixrunde) — 2026-10-05

**Review-Art:** Code — Re-Review der Fixrunde gegen Plan, Entscheidungen und Hard Rules (nicht gegen die DoD; die prüft der Verifier)

**Gegenstand:** Diff 9d1edf2a..6abac454, die Commits a8e9ddf0 (Runner `version_lesen`, Umbenennung in Gast-Dockerfiles und `run.sh`, Vertrag) und 6abac454 (Plan: Abschnitt „Fixrunde“ in §3, zwei Suchlauf-Zeilen; Vertrag: Mutations-Zeile zum Suffix). Bezug: Review `review-slice-sdk-kompat-version-parameter` (Commit 9d1edf2a), F-1 bis F-6.

**Skill:** `.harness/skills/reviewer.md` @ 7fff2cac
**Modell:** claude-opus-5-5 · **Datum:** 2026-10-05

> **Zitier-Form.** Dieser Report friert ein; er zitiert Kennungen, nicht
> Adressen (`slice-<Kennung>`, `make <target>`, Baseline-Stellen als Tag +
> Pfad in Inline-Code).

**Eingangs-Kontext:**

- Slice-Plan `slice-sdk-kompat-version-parameter` (§3, Abschnitt „Fixrunde“, Nachzug-Zeilen, `suchlauf`-Block)
- Review `review-slice-sdk-kompat-version-parameter` (F-1 bis F-6)
- [`ADR-0145`](../plan/adr/0145-sdk-meldungscodes-eigenschaft-der-fehlertypen.md) Festlegung 4
- Vertrag `harness/targets/sdk-kompat.md` (§Versionen, §Ausgänge, §Grenze 8/9, §Test)
- `AGENTS.md` §3.1, §3.7, §3.9, §3.12, §3.13; keine `LH-*`-Kennung berührt (Harness-Werkzeug)

**Eigene Proben dieses Laufs** (alle an einem frischen Klon des Repos, Stand 6abac454, im Scratchpad; `docker` ist ein Stub, der seine Argumente protokolliert und Exit 0 liefert — kein Bau, kein Netz; ein Lauf, der die Versionsprüfung passiert, endet mit dem Stub mit Exit 1 an `A1 — keine gedruckte Zeile`; je Fall ein Verzeichnis mit einer leeren Datei des Artefakt-Namens der eingesetzten Version über `SDK_KOMPAT_DIST_<SPRACHE>`, `SDK_KOMPAT_SPRACHEN` auf die Sprache gesetzt; gemessen):

- **Suffix, je Version-Datei:** `<Version>0.7.0-rc.1</Version>` und `<Version>0.7.0+b1</Version>` in der `.csproj`, `version = "0.7.0-rc.1"` in `build.gradle.kts` (oberste Zeile), `version = "0.7.0-rc.1"` und `version = "0.7.0rc1"` in `pyproject.toml`: je Exit 2, null `docker`-Aufrufe, Meldung `… trägt nicht genau eine Version der Form X.Y.Z, nur Ziffern ohne Suffix (gelesen: <Wert>)`.
- **Gegenprobe `0.6.1`** (unveränderte Datei) je Sprache: Versionsprüfung passiert, Kopfzeile `… Modus dist, neue Bibliothek 0.6.1, veröffentlicht 0.6.0`, drei `docker`-Aufrufe, Exit 1 erst am Stub.
- **Grenzfälle**, je Sprache: `0.6` und `0.6.1.2` → Exit 2 (`gelesen: 0.6` bzw. `0.6.1.2`); führendes und nachgestelltes Leerzeichen im Wert (`<Version> 0.6.1</Version>`, `"0.6.1 "` usw.) → Exit 2 (`gelesen:  0.6.1` bzw. `0.6.1 `); `v0.6.1` in der `.csproj` → Exit 2; Leerzeichen hinter dem schließenden Anführungszeichen in `build.gradle.kts` → Exit 2 (`gelesen: keine`); Leerzeichen hinter `</Version>` → passiert (das Muster lässt es zu, der Wert ist rein); `00.06.01` in `pyproject.toml` → passiert (F-1 unten).
- **Mutations-Zeile „ohne Ersatzwert“ des Vertrags** am Stand 6abac454 nachgefahren (die vier Fälle der Zeile: `.csproj` auf `0.6.9`, Versionszeile aus `pyproject.toml` entfernt, zweite oberste `version = "0.6.2"` in `build.gradle.kts`, `.csproj` gelöscht): je Exit 2, Meldungen wie im Vertrag zitiert, die Form-Meldung mit dem Zusatz „, nur Ziffern ohne Suffix“ — die Kürzungs-Notiz der Zeile stimmt.
- `make suchlauf-nachmessen PLAN=<Plan>`: 8 Zeilen stimmen, Exit 0.
- `make kommentar-kennungen DIFF=dad07666`: kein Kandidat, Exit 0 (Probe der Form, kein Beleg der Wahrheit).
- `git grep -n -E 'alt06|neu06|LIB06|\bd06\b|wheels/06|Gast-alt06|Gast-neu06' -- . ':!docs/reviews' ':!.harness'`: Treffer nur im Plan (Umbenennungs-Liste und Suchlauf-Zeile), keiner unter `tools/` oder `harness/`.
- Gelesen gegen die Plan-Aussagen (Abschnitt „Fixrunde“, *Läufe am Endstand*): die Logs des Implementer-Laufs im Scratchpad, `kv/dist2.log` (Exit 0 laut `kv/dist2.rc`, Schlusszeile `… (dist) grün für: csharp kotlin python`, null Zeilen `ROT`, A2 `Bibliothek 0.6.1 (Artefakt …)` mit den Prüfsummen `e5fa52ed7181`/`af19721beee7` wie in `kv/dist.log`, `KOMPAT kotlin javap Gast (neureg)` und `javap Bibliothek REG`) und `kv/reg2.log` (Exit 0 laut `kv/reg2.rc`, `… (registry) grün für: csharp kotlin python`, null Zeilen `ROT`). Die Läufe von `make test-sdk-kompat` selbst sind **nicht wiederholt, übernommen** aus diesen Logs.

---

## Prüfung der Findings des ersten Reviews

| Finding | Behauptung der Fixrunde | Ergebnis dieses Laufs |
|---|---|---|
| F-1 (HIGH) | `version_lesen` verlangt `^[0-9]+\.[0-9]+\.[0-9]+$`; Kommentar, Meldung, §Versionen und §Ausgänge sagen dasselbe | **gelöst.** Suffix in allen drei Version-Dateien endet mit Exit 2 vor jedem `docker`-Aufruf; Gegenprobe `0.6.1` passiert; `0.6`, `0.6.1.2`, Leerzeichen im Wert enden mit Exit 2. Kommentar (`kein Suffix wie -rc.1`), Meldung und Vertrag tragen dieselbe Form. Rest: F-1 unten (INFO) |
| F-2 (LOW) | §3-Tabelle meldet zusätzlich Titel von `harness/targets/sdk-altserver.md` und Kopf von `tools/harness/run-sdk-altserver-tests.sh`, Frist Closure; Suchlauf-Zeile 7 | **gelöst** (als Meldung): die Zeile nennt alle drei Träger und die Frist; Suchlauf-Zeile 7 misst 4 Treffer am Arbeitsbaum. Der Nachzug selbst bleibt beim Planner bis zur Closure |
| F-3 (INFO) | §Grenze 8 im Vertrag | **gelöst.** C# A2 kopiert `alt05` und ersetzt nur die DLL, Kotlin A2 kopiert `altreg/lib` — wie Grenze 8 sagt. Rest: F-3 unten (INFO) |
| F-4 (INFO) | §Grenze 9; §Versionen „die nicht eingerückte Zeile“ | **gelöst.** Das Muster `^version = "…"$` liest nur die nicht eingerückte Zeile; die Bindung von Jar-Version und POM-`<version>` steht in der `pack`-Probe von `sdks/kotlin/Dockerfile` (gelesen, nicht gefahren) |
| F-5 (INFO) | Zweig behalten, als über den Runner nicht erreichbar gekennzeichnet, Kopplung an die RUN-Schicht `lib-dist` | **gelöst**, Kennzeichnung steht; die genannte Schicht (`mkdir -p /kompat/lib-dist`, gleiche Bedingung, `set -eu`) existiert. Wortlaut der Kopplung: F-2 unten (INFO) |
| F-6 (INFO) | `alt06`→`altreg`, `neu06`→`neureg`, `LIB06`→`LIBREG`, `d06`→`dreg`, `/wheels/06`→`/wheels/neu`; gedruckte Zeile `Bibliothek REG` | **gelöst**, vollständig und konsistent über beide Gast-Dockerfiles (Verzeichnisse und `-o`-Ziele), die drei `run.sh`, Python-Dockerfile und `run.sh` (`/wheels/neu` an Schreib- und Lesestelle); kein toter Bezeichner, keine Gast-Quelle und kein Vertrag trug die alten Namen; die Logs zeigen die neuen Zeilen. Die Suchlauf-Zeile 8 deckt einen der fünf Namen nicht: F-4 unten (LOW) |

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | INFO | Die Formprüfung lässt führende Nullen durch: `version = "00.06.01"` in `pyproject.toml` passiert sie (gemessen, Kopfzeile `neue Bibliothek 00.06.01`). In Python normalisiert PEP 440 das zu `0.6.1`; das Rad hieße anders als der verlangte Name, und A2 vergliche `installiert` mit dem nicht normalisierten Text (*hergeleitet*, nicht gebaut). Kein Package trägt heute eine solche Version. | Maintainability | `tools/harness/run-sdk-kompat-tests.sh` · `grep -Eq '^[0-9]+\.[0-9]+\.[0-9]+$'` | ja — Lauf des Runners mit `00.06.01` | Form enger als Normalform nicht benannt |
| F-2 | INFO | Der neue Satz „Kopplung: ändert sich die Bedingung dort, gilt der Zweig hier weiter.“ nennt die Stelle (RUN-Schicht `lib-dist`), aber nicht, was mit ihr zusammen zu ändern ist; er liest sich als Zusage, dass der Zweig ohne Änderung weiter greift. Inhaltlich stimmt das (der Zweig färbt A2 rot, wenn die DLL fehlt). | `AGENTS.md` §3.7 (Kopplung) | `tools/harness/sdk-kompat/csharp/run.sh` · `Kopplung: ändert sich die Bedingung dort, gilt der Zweig hier weiter.` | nein — Lese-Handlung | Kopplung ohne mitzuändernde Handlung |
| F-3 | INFO | §Grenze 8 sagt ohne Sprach-Einschränkung „A5 vergleicht die Abhängigkeitsmengen von 0.5.0 und `REGISTRY_VERSION`“; ein solcher Vergleich steht nur in C# (`deps`) und Kotlin (`d05`/`dreg`), Python-A5 vergleicht keine Abhängigkeiten (`git grep -n abhaengigkeiten -- tools/harness/sdk-kompat`: Treffer nur in `csharp/run.sh` und `kotlin/run.sh`). | Maintainability | `harness/targets/sdk-kompat.md` · `A5 vergleicht die` | ja — der genannte `git grep` | Grenze ohne Geltungsbereich |
| F-4 | LOW | Der Plan stützt die Umbenennung der fünf Namen mit „Suchlauf-Zeile 8: 0 Treffer“; das Muster der Zeile (`-e alt06 -e neu06 -e LIB06 -e wheels/06`) trägt `d06` nicht. Die Aussage stimmt (eigene Messung `\bd06\b`: 0 Treffer unter `tools/`), die Zeile misst sie für diesen Namen nicht. Gewogen gegen die HIGH-Klasse „Beleg trägt seinen Satz nicht“: nicht genommen, weil die Zeile genau misst, was ihr Muster nennt, und die Lücke die Vollständigkeit des Musters ist, die `AGENTS.md` §3.13 dem Reviewer als Lese-Handlung zuweist. | `AGENTS.md` §3.13 (Suchform: Symbolnamen der bewegten Eigenschaft) | Plan `slice-sdk-kompat-version-parameter` · `diff 0 -n -e alt06 -e neu06 -e LIB06 -e wheels/06 -- tools/harness/sdk-kompat` | ja — `git grep -n -w d06 -- tools/harness/sdk-kompat` | Suchmuster deckt nicht jeden Symbolnamen |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `tools/harness/run-sdk-kompat-tests.sh` — `version_lesen` nach der Fixrunde | geprüft, ohne Befund außer F-1: Suffix (`-rc.1`, `+b1`, PEP-440-Form `rc1`), zwei und vier Ziffernfolgen, Leerzeichen im Wert und Präfix `v` enden in allen drei Sprachen mit Exit 2 vor jedem `docker`-Aufruf; `0.6.1` passiert; Modus `registry` liest weiter keine Version-Datei (Code unverändert) |
| Kommentar `version_lesen` gegen den Code | geprüft, ohne Befund: die Zusage „Exit 2, wenn … nicht genau eine Versionszeile in der Form X.Y.Z (drei Ziffernfolgen, kein Suffix wie -rc.1)“ trägt der Code (gemessen) |
| Umbenennung in `tools/harness/sdk-kompat/{csharp,kotlin,python}` | geprüft, ohne Befund: Schreib- und Lesestellen jedes Namens gepaart (`Gast-altreg`/`bin/altreg`, `neureg`, `LIBREG` über `LIB$v` mit `v in 05 REG`, `dreg`, `/wheels/neu`); kein Rest der alten Namen außerhalb des Plans |
| Vertrag `harness/targets/sdk-kompat.md` (§Versionen, §Ausgänge, §Grenze 8/9, §Test) | geprüft, Befund F-3; die PEP-440-Aussage ist als *hergeleitet* gekennzeichnet, die neue Mutations-Zeile zitiert die Meldungen wörtlich (selbst nachgefahren), die Kürzungs-Notiz der älteren Zeile stimmt am Stand 6abac454 (selbst nachgefahren) |
| Plan, Abschnitt „Fixrunde“, gegen Code und Logs | geprüft, Befund F-4; Exit-Werte, Schlusszeilen, Prüfsummen und die Zeile `Bibliothek REG` in `kv/dist2.log`/`kv/reg2.log` stimmen mit dem Plan überein |
| Kommentare (`AGENTS.md` §3.7) | geprüft, ohne Befund außer F-2: keine Vorher/Nachher-Sprache, kein Konjunktiv über eine verworfene Alternative in den geänderten Kommentaren, keine Kennungskette (`make kommentar-kennungen DIFF=dad07666`) |
| Neue Befunde durch die Fixrunde | geprüft: F-1 bis F-4 dieses Reports; kein neuer Fehlerpfad, keine geänderte Ausgangsklasse |
| Docker-only, Suppression, ADR-Immutabilität, Traceability | geprüft, ohne Befund: keine Host-Toolchain, kein `-i`, kein Suppression-Kommentar, keine ADR geändert, beide Commits nennen `ADR-0145` |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Form enger als Normalform nicht benannt · Kopplung ohne mitzuändernde Handlung · Grenze ohne Geltungsbereich · Suchmuster deckt nicht jeden Symbolnamen

## Verdikt

**Merge-blockierend:** nein — F-1 des ersten Reviews (HIGH) ist gelöst, kein HIGH offen. F-4 (LOW) und die INFOs brauchen keine weitere Fixrunde; der Implementer akzeptiert oder begründet (Modul 8, isolierte LOW/INFO). F-2 des ersten Reviews bleibt als Meldung an den Planner mit Frist Closure dieses Slice.

**DoD-Zeile „Review durchgeführt“:** Nach `.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug zöge der Reviewer sie bei diesem Verdikt im selben Commit nach. Dieser Lauf ändert auf Anweisung des Auftraggebers keine Repo-Datei außer dem Report; die Zeile bleibt offen und geht als offener Punkt an den Auftraggeber.

**Übergabe:** Findings an den Implementer; die Finding-Klassen in die Slice-Closure §7. Dieser Report ist Lauf-Beleg und ersetzt keine Verifikation.
