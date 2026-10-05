# Review-Report: slice-sdk-kompat-version-parameter — 2026-10-05

**Review-Art:** Code — gegen Plan, Entscheidungen und Hard Rules (nicht gegen die DoD; die prüft der Verifier)

**Gegenstand:** Diff dad07666..f13a0e4a, die Commits 500bebb2 (Runner, Gast-Dockerfiles, `run.sh`, Vertrag, `sdk.mk`, Index-Zeile) und f13a0e4a (Plan: Messung, Mutationen, Suchlauf, Nachzug). Der Roadmap-Commit 3b173e28 ist Lifecycle-Arbeit und nicht Gegenstand.

**Skill:** `.harness/skills/reviewer.md` @ 7fff2cac
**Modell:** claude-opus-5-5 · **Datum:** 2026-10-05

> **Zitier-Form.** Dieser Report friert ein; er zitiert Kennungen, nicht
> Adressen (`slice-<Kennung>`, `make <target>`, Baseline-Stellen als Tag +
> Pfad in Inline-Code).

**Eingangs-Kontext:**

- Slice-Plan `slice-sdk-kompat-version-parameter` (§1 Ziel und Abgrenzung, §3 Plan samt Nachzug-Zeilen)
- [`ADR-0145`](../plan/adr/0145-sdk-meldungscodes-eigenschaft-der-fehlertypen.md) Festlegung 4
- Vertrag `harness/targets/sdk-kompat.md` (Stand des Diffs)
- `AGENTS.md` §3.1, §3.7, §3.9, §3.12, §3.13; keine `LH-*`-Kennung berührt (Harness-Werkzeug)

**Eigene Proben dieses Laufs:**

- `make suchlauf-nachmessen PLAN=<Plan>`: 6 Zeilen stimmen, Exit 0.
- `make kommentar-kennungen DIFF=dad07666`: kein Kandidat, Exit 0 (Probe der Form, kein Beleg der Wahrheit).
- Fünf Läufe des Runners an einem Klon des Repos (Stand f13a0e4a) im Scratchpad, mit einem Stub für `docker`, der seine Argumente protokolliert und Exit 0 liefert (kein Bau, kein Netz; Exit 1 der Läufe mit Stub kommt aus den fehlenden `KOMPAT`-Zeilen):
  - **R0** (Gegenprobe, Modus `dist`, Artefakt-Verzeichnisse mit leeren Dateien der Namen `…0.6.1…` über `SDK_KOMPAT_DIST_*`): je Sprache `docker build … --build-arg NEU_VERSION=0.6.1 --build-arg REGISTRY_VERSION=0.6.0 …`, Kopfzeile `neue Bibliothek 0.6.1, veröffentlicht 0.6.0`.
  - **R1** (nur die eingerückte Publikations-Zeile `version = "0.9.9"` in `build.gradle.kts`): der Runner liest `0.6.1`, keine Meldung.
  - **R3** (`<Version></Version>` in der `.csproj`): Exit 2, `… PgChangeFeed.Client.csproj trägt nicht genau eine Version der Form X.Y.Z (gelesen: keine)`.
  - **R4** (Modus `registry`, `.csproj` gelöscht): kein Lesen der Version-Datei, Bau mit `NEU_VERSION=0.6.0`.
  - **R5** (`<Version>0.7.0-rc.1</Version>` in der `.csproj`, Verzeichnis mit `PgChangeFeed.Client.0.7.0-rc.1.nupkg`): kein Exit 2, der Runner baut (`neue Bibliothek 0.7.0-rc.1`, zwei `docker build`-Aufrufe); dieselbe Version in `pyproject.toml`: Exit 2 erst an der Artefakt-Prüfung mit dem Hinweis `make sdk-pack-python vorher`.
- Gelesen gegen die Plan-Aussagen (§3 *Messung*, *Mutationen*): die Logs des Implementer-Laufs im Scratchpad — `kv/dist.log` und `kv/registry.log` (Exit 0, Schlusszeilen und A2-Zeilen wie im Plan, Prüfsummen `e5fa52ed7181`/`f8e6415b23e5`), `kv/m6.log` (`KOMPAT python A2: ROT — Austausch auf 0.6.1 gescheitert (pip Exit 1, installiert 0.5.0, …)`), `vorher-dist.log` (`NU1603 …`, Exit 2), `vorher-kotlin.log` (Exit 0, `Bibliothek 0.6.0 (Artefakt pgchangefeed-kotlin-0.6.1.jar)`), `vorher-python.log` (Exit 2, `A2: Bibliothek pgchangefeed 0.5.0 … 45 Aufrufe ok`). Die Läufe von `make test-sdk-kompat` selbst sind **nicht wiederholt, übernommen** aus diesen Logs.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | HIGH | Der Kommentar von `version_lesen`, der Vertrag (§Versionen, Exit-Code-Tabelle) und die Meldung sagen Exit 2 zu, wenn die Version-Datei nicht genau eine Version „der Form X.Y.Z“ trägt; der Code lässt jede Version mit Suffix durch (`([-+][0-9A-Za-z.-]+)?`). Gemessen (R5): `<Version>0.7.0-rc.1</Version>` endet nicht mit Exit 2, der Runner baut; in Python endet dieselbe Version erst an der Artefakt-Prüfung mit dem Hinweis `make sdk-pack-python vorher`, den ein neuer Pack nicht erfüllt, weil das Rad den normalisierten Namen trägt (*hergeleitet*, PEP-440-Normalisierung, nicht gebaut). | `AGENTS.md` §3.7 (Zusage); Reviewer-Skill HIGH „Kommentar trägt keine der Kommentar-Klassen“ (*Zusage*) | `tools/harness/run-sdk-kompat-tests.sh` · `nicht genau eine Versionszeile in der Form X.Y.Z trägt`; `harness/targets/sdk-kompat.md` · `trägt sie nicht genau eine Version der Form` | ja — Lauf des Runners mit einer Version mit Suffix (R5) | Zusage im Kommentar nicht getragen |
| F-2 | LOW | Die Meldung des überholten Trägers nennt nur `harness/mk/sdk.mk` (Ziel `test-sdk-altserver`); dieselbe Beschreibung „SDK 0.6.0“ bzw. „0.6.0-SDKs“ für ein Ziel, das die SDK-Quelle des Arbeitsstands (0.6.1) fährt, steht auch im Titel von `harness/targets/sdk-altserver.md` und im Kopf von `tools/harness/run-sdk-altserver-tests.sh`. Das Suchmuster in §3 (`Bibliothek.0\.6\.0`, `Packages.0\.6\.0`, `0\.6\.0-Pakete`) trifft keine der drei Formen. | `AGENTS.md` §3.13 (Meldung eines fremden Trägers, Frist Closure) | `harness/targets/sdk-altserver.md` · `SDK 0.6.0 gegen einen Server vor 0.6.0`; `tools/harness/run-sdk-altserver-tests.sh` · `die Fehlerfälle der 0.6.0-SDKs` | ja — `git grep -n -e '0\.6\.0-SDKs' -e 'SDK 0\.6\.0' -- harness tools/harness` (4 Treffer in 3 Dateien am Stand f13a0e4a) | Träger-Meldung unvollständig |
| F-3 | INFO | A2 tauscht nur die Bibliotheksdatei; die übrigen Laufzeit-Abhängigkeiten sind in Kotlin die der Version `REGISTRY_VERSION` (`alt06/lib`), in C# die des 0.5.0-Publish (`alt05`), und A5 vergleicht die Abhängigkeitsmengen von 0.5.0 und `REGISTRY_VERSION`, nicht die der neuen Version. Seit dieser Arbeit unterscheiden sich `NEU_VERSION` und `REGISTRY_VERSION`; die Grenze steht nicht in §Grenze. | Maintainability | `tools/harness/sdk-kompat/kotlin/run.sh` · `cp -r "$BIN/alt06/lib" /tmp/a2lib` | nein — eine Lese-Handlung | Grenze nicht benannt |
| F-4 | INFO | `build.gradle.kts` trägt zwei Zeilen `version = "…"` (oberste Ebene und Publikation, eingerückt); der Runner liest nur die erste (R1: Publikation auf `0.9.9`, keine Meldung), der Vertrag spricht von „die Zeile `version = "…"`“. Die Gleichheit beider hält die Probe in `make sdk-pack-kotlin` (dort beschrieben, hier nicht nachgefahren). | Maintainability | `harness/targets/sdk-kompat.md` · `die Zeile` | nein | Grenze nicht benannt |
| F-5 | INFO | Der Zweig `KOMPAT csharp A2: ROT — Bibliotheksdatei … fehlt` ist über den Runner nicht erreichbar: dieselbe Bedingung (`/neu/PgChangeFeed.Client.$NEU_VERSION.nupkg` vorhanden) lässt im Dockerfile `cp …/lib/net10.0/PgChangeFeed.Client.dll` unter `set -eu` laufen, ein Fehlen bricht den Bau ab. Der Plan (§3 Nachzug) führt den Zweig, aber keine Mutation; das Verhalten stimmt an der Stelle, ist aber ungemessen. | Maintainability | `tools/harness/sdk-kompat/csharp/run.sh` · `Bibliotheksdatei $neu_dll fehlt` | nein | Pfad ohne erreichbare Eingabe |
| F-6 | INFO | Bezeichner tragen die Version weiter implizit: `alt06`, `neu06`, `LIB06`, `d06` stehen für `REGISTRY_VERSION`, `/wheels/06` in Python für `NEU_VERSION` (0.6.1). Der Suchlauf trifft sie nicht (§3 nennt das als Lese-Handlung des Reviewers). | Maintainability | `tools/harness/sdk-kompat/python/Dockerfile` · `pip download --no-cache-dir -d /wheels/06 "$rad"` | nein | implizite Version im Bezeichner |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `tools/harness/run-sdk-kompat-tests.sh` — Lesen der Version, fail-closed | geprüft, ohne Befund außer F-1: `exit 2` in `version_lesen`/`neu_verzeichnis` wirkt über die Kommandosubstitution unter `set -euo pipefail` (R3: Exit 2); `grep -c .` bei null Treffern liefert `0` ohne Abbruch; kein Ersatzwert; Modus `registry` liest keine Version-Datei (R4) |
| Weiterreichen der Build-Argumente | geprüft, ohne Befund: `NEU_VERSION`/`REGISTRY_VERSION` kommen in jedem Gast-Bau an (R0); C# und Kotlin deklarieren beide `ARG`, Python nur `NEU_VERSION`; `RUN test -n …` bricht ohne Wert ab; `run.sh` prüft `${…:?}` |
| `tools/harness/sdk-kompat/{csharp,kotlin,python}` — Literal `0.6.0` | geprüft, ohne Befund: 0 Treffer am Diff, 39 am Stand 3e028ee3 (`make suchlauf-nachmessen`); einziger Treffer im Runner ist `REGISTRY_VERSION=0.6.0` |
| Python-A2-Austausch (über den Plan hinaus, §3 Nachzug) | geprüft, ohne Befund: pip-Exit und installierte Version gebunden, Rot gemessen (`kv/m6.log`) |
| Vertrag `harness/targets/sdk-kompat.md` (§Versionen, Grenze 6/7, §Test) | geprüft, Befunde F-1, F-3, F-4; Mutationszeilen nennen Stellen und Menge, die Übertragung ist als *hergeleitet* gekennzeichnet |
| Index-Zeile `harness/README.md` (`make test-sdk-kompat`) | geprüft, ohne Befund: Beschreibungszelle 89 Zeichen |
| Kommentare (`AGENTS.md` §3.7) | geprüft, ohne Befund außer F-1: keine Vorher/Nachher-Sprache, keine Kennungskette (`make kommentar-kennungen DIFF=dad07666`) |
| `harness/mk/sdk.mk` (`test-sdk-kompat`) | geprüft, ohne Befund |
| Plan-Aussagen gegen die Logs des Implementers | geprüft, ohne Befund: Exit-Werte und Zeilen in `kv/` und `vorher-*.log` stimmen mit §3 überein |
| Docker-only, Suppression, ADR-Immutabilität, Traceability | geprüft, ohne Befund: keine Host-Toolchain, kein `-i`, kein Suppression-Kommentar, keine ADR geändert, beide Commits nennen `ADR-0145` |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 4 |

**Finding-Klassen dieses Laufs:** Zusage im Kommentar nicht getragen · Träger-Meldung unvollständig · Grenze nicht benannt · Pfad ohne erreichbare Eingabe · implizite Version im Bezeichner

## Verdikt

**Merge-blockierend:** ja — F-1 (HIGH) braucht eine Fixrunde des Implementers; ob der Code oder die Zusage (Kommentar, Vertrag, Meldung) bewegt wird, entscheidet der Implementer. F-2 geht als Meldung an den Planner (Frist: Closure dieses Slice). Die DoD-Zeile „Review durchgeführt“ bleibt offen und wird nach der Fixrunde nachgezogen.

**Übergabe:** Findings an den Implementer; die Finding-Klassen in die Slice-Closure §7. Dieser Report ist Lauf-Beleg und ersetzt keine Verifikation.
