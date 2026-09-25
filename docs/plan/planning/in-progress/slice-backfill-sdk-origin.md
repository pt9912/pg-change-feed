# Slice backfill-sdk-origin: SDK-HTTP-Lesemodelle — `origin` als optionales Feld in C#, Kotlin und Python; Package-Versionen

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [welle-backfill-bestand](../welle-backfill-bestand.md).

**Bezug:** [`LH-FA-SST-009`](../../../../spec/lastenheft.md) (Client-Bibliotheken — die HTTP-Fläche der drei
Packages), [`LH-FA-SST-006`](../../../../spec/lastenheft.md) (HTTP-API, `GET /changes`), [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 8
(Reichweite in SDKs: `origin` als optionales Feld, fehlt es, gilt `wal`; die
Package-Version hebt der SDK-Slice nach dem bestehenden Muster),
[`ADR-0106`](../../adr/0106-csharp-nuget-erstes-sdk-package.md), [`ADR-0107`](../../adr/0107-python-pypi-zweites-sdk-package.md), [`ADR-0109`](../../adr/0109-kotlin-github-packages-drittes-sdk-package.md) (die drei Package-Entscheidungen),
[`ADR-0110`](../../adr/0110-python-sdk-umfang-erweitert-vollmatrix.md) (Sprachmatrix, Belegklasse der Flächen).

**Berührte Spec-Stellen:** [`SPEC-022`](../../../../spec/pflichtenheft.md) (`GET /changes`, durch `spec-nachzug`
um `origin` ergänzt) — gelesen; [`SPEC-026`](../../../../spec/pflichtenheft.md), [`SPEC-027`](../../../../spec/pflichtenheft.md), [`SPEC-028`](../../../../spec/pflichtenheft.md)
(Package-Zeilen mit der aktuellen Version) — geändert.

**Verantwortlich:** Implementer-Agent, 2026-09-25.

**Autor:** Planner-Agent, Welle-Eröffnung [welle-backfill-bestand](../welle-backfill-bestand.md). **Datum:** 2026-09-23.

---

## 1. Ziel und Abgrenzung

**Ziel:** Die HTTP-Lesemodelle der drei SDK-Packages tragen `origin` als
**optionales** Feld (`wal`, wenn der Server es nicht sendet) und die Packages
heben ihre Version: C# `PgChangeFeed.Client`
(`sdks/csharp/PgChangeFeed.Client/Http/Models/Changes.cs`), Kotlin
`pgchangefeed-kotlin` (`sdks/kotlin/pgchangefeed-kotlin/src/main/kotlin/io/github/pt9912/pgchangefeed/http/model/Changes.kt`)
und Python `pgchangefeed` (`sdks/python/pgchangefeed/src/pgchangefeed/models.py` samt
`http_client.py`). Ohne diese Anpassung funktionieren die Packages unverändert
weiter (erwartet: die Decoder ignorieren das Feld) — der Slice macht das Feld
**nutzbar** und belegt beide Fälle (mit und ohne Feld).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Die Live-Flächen** (gRPC, SSE, NATS-Vollinhalt) — sie tragen kein `origin`
  ([`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 8); die zehn Felder bleiben.
- **Ein Release** — ein Tag-Push (`sdk-*-v*`) ist Betreiber-Handlung außerhalb der
  Welle ([`AGENTS.md`](../../../../AGENTS.md) §3.10); die drei `sdk-*-release.yml`-Workflows bleiben
  unverändert. Der Slice hebt nur die Version in den Metadaten-Quellen.
- **Realserver-Läufe der SDKs** (`make test-sdk-*-integration`) — der Bestand
  dieser Läufe (HTTP: Registrierung und Listen-Aufruf) ruft `GET /changes` nicht
  auf; der Slice erweitert sie nicht.
- **Beispiel-Clients** (`examples/**`) — [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) nennt einen
  Go-HTTP-Beispiel-Client; ob ein Beispiel `GET /changes` überhaupt dekodiert, ist
  ungeprüft (die HTTP-Beispiele rufen laut `.a-check.yml`-Kommentar `GET /tables`
  auf, die NATS-Beispiele bauen eine Abfrage-URL). *Zu belegen durch:* Lesen der
  Dekodierung je Beispiel; findet sich keine, **entfällt** der Punkt und der
  Bericht sagt es.

## 2. Definition of Done

- [x] Die drei HTTP-Lesemodelle tragen `origin` als optionales Feld: eine
      Antwort mit `origin` liefert den Wert, eine Antwort ohne das Feld liefert
      `wal`; je Sprache ein Unit-Test mit beiden Fällen (dazu unbekannter und
      leerer Wert durchgereicht, JSON-`null` liest als `wal`) und einer
      Fixture, deren Ursprung im Test genannt ist: die Antwortform des
      Server-Handlers (`TestReadChangesTraegtOriginAlsLetztesFeld`), **kein**
      Lauf gegen einen Server-Container. *Zu
      belegen durch:* `make sdk-pack-csharp`, `make sdk-pack-kotlin` und
      `make sdk-pack-python` — der Bau führt die Tests aus und bricht bei rotem
      Test ab; jeder Aufruf trägt `--build-context proto=proto` (er ist für den
      **ganzen** Bau jedes SDK zwingend, nicht nur für gRPC:
      `BEO-PGC/zusatzkontext-kopplung-breiter-als-dod-wortlaut`, offen, 3×).
- [x] Die Versionen in den drei Metadaten-Quellen (`.csproj` `<Version>`,
      `build.gradle.kts` `version`, `pyproject.toml` `[project] version`) stehen
      auf `0.2.0`: der nächsten Minor-Version nach den vorhandenen Tags
      (`0.1.0`), die noch kein Tag trägt — keine weitere Hebung, Wahl und
      Grund in §3 „Versionsentscheidung"; die Träger sind geprüft
      (§3.13-Suchlauf), die Artefaktnamen sind durch **reale** Pack-Läufe
      bestätigt.
- [x] Die Sprachreinheit der Kopien: die drei Sprach-Änderungen sind Kopien
      derselben Form; jede Kopie ist auf verbliebene deutsche Wortfragmente in
      Code, Kommentaren, Docstrings und README geprüft, nicht nur auf ihre
      Übereinstimmung mit dem Vorbild
      (`BEO-PGC/formvorbild-kopie-traegt-deutsches-wortfragment-weiter`,
      `BEO-PGC/deutsches-fachwort-im-englischen-sdk-readme`, je verkörpert, 3×).
- [x] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor, kein
      offenes HIGH/MEDIUM (`.harness/skills/reviewer.md`; die zwei HIGH
      F-1/F-2 im Plan-Text, die LOW/INFO F-3 bis F-8 in der Fixrunde
      aufgelöst, siehe §3) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein Self-Review (Modul 8).
- [x] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff)
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [x] Doku-Update: SDK-READMEs (englisch) nennen das Feld und die Version; Benutzerhandbuch (`**SDK:**`-Absätze der HTTP-Beschreibung) und Änderungshistorie, soweit sie die Version oder die Feldmenge tragen.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag (geschärfte Regel ·
      neuer Sensor · benannte Spec-Lücke).
- [ ] Reconciliation-Register — entfällt: keine Reconciliation-Datei in
      diesem Repo (Greenfield).
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder weitere `evidence/`-Datei; kein Anfall ist ebenfalls
      eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen /
      weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen —
      von der Closure der Welle [welle-backfill-bestand](../welle-backfill-bestand.md) (die Roadmap führt sie unter
      *Offene Wellen*, das Ereignis kann eintreten).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `sdks/csharp/PgChangeFeed.Client/Http/Models/Changes.cs` (+ Tests in `PgChangeFeed.Client.Tests`) | update | Feld `Origin`; Default `wal`. |
| `sdks/kotlin/pgchangefeed-kotlin/src/main/kotlin/io/github/pt9912/pgchangefeed/http/model/Changes.kt` (+ Test) | update | Feld; Gson setzt ein fehlendes Feld nicht auf den Kotlin-Default (erwartet, zu belegen) — die Abbildung trägt den Fall. |
| `sdks/python/pgchangefeed/src/pgchangefeed/models.py`, `http_client.py` (+ `tests/test_http_client.py`) | update | Feld und Dekodierung mit expliziten Feldern. |
| `sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.csproj`, `sdks/kotlin/pgchangefeed-kotlin/build.gradle.kts`, `sdks/python/pgchangefeed/pyproject.toml` | update | Version. |
| `sdks/*/README.md` | update | Feld und Version; englisch. |
| `spec/pflichtenheft.md` §6 (`SPEC-026`, `SPEC-027`, `SPEC-028`) und §7 Historie | update | die Zeilen tragen die aktuelle Version; eine Historie-Zeile je Package im Muster der bisherigen Nachzüge, ohne ADR-/Slice-Bezug. |
| `harness/README.md` §Sensors (`make sdk-pack-*`) | update | real erzeugte Artefaktnamen; aus **realen** Läufen geschrieben. |
| `harness/mk/sdk.mk` | prüfen | trägt der Kommentar die Version? |
| `docs/user/benutzerhandbuch.md` | prüfen/update | `**SDK:**`-Absätze; Änderungshistorie. |
| Kommentare mit Feldanzahl in den SDK-Bäumen (gemeldet von `slice-backfill-change-origin`, Fixrunde, Suchlauf über `sdks/` und `examples/`) | update | `Sse/Models/Change.cs` Z. 17 und `Nats/Models/Change.cs` Z. 21 (C#, „eleven fields" für `GET /changes`), `sse/model/Change.kt` Z. 19–20 (Kotlin, „twelve fields"), `models.py` Z. 264 (Python, „zwölf des HTTP-Lesezugriffs"): der HTTP-Lesezugriff trägt mit `origin` dreizehn Felder; die C#-Kommentare nennen schon am Parent eine Zahl, die nicht zu den gezählten zwölf Feldern passt. Ebenso `Http/Models/Changes.cs` Z. 8 und `http/model/Changes.kt` Z. 8 („the same ten fields as the domain type"): das HTTP-Modell trägt am Parent zwölf, jetzt dreizehn Felder. Die Kommentare sind Träger der bewegten Eigenschaft „Feldmenge des HTTP-Lesemodells" und werden mit dem Feld nachgezogen. |
| **Versionsentscheidung: keine Hebung** — `.csproj` `<Version>`, `build.gradle.kts` `version` (Zeile 108, die eingerückte Zeile 230 im `publishing`-Block trägt dieselbe Zahl), `pyproject.toml` `version` bleiben `0.2.0` | Reduktion (Plan-Punkt „Version heben" entfällt) | *Beleg (gemessen 2026-09-25, Fixrunde):* `git tag -l 'sdk-*'` nennt `sdk-csharp-v0.1.0` (Commit `1a20335a`) und `sdk-python-v0.1.0` (`824f11a3`), Kotlin: kein Tag; `git show <Tag>:<Quelle>` liefert in beiden Tags `0.1.0`. Die Quellen stehen bei `0.2.0` (Arbeitsbaum: `<Version>0.2.0</Version>`, `version = "0.2.0"` in `pyproject.toml` Z. 7 und `build.gradle.kts` Z. 108). `git log -S'<Version>0.2.0'` (C#) und `git log -S'version = "0.2.0"'` (Python, Kotlin) nennen genau `d86d1965` (C#), `0f8cc4f2` (Python) und `c8c9e3ae` (Kotlin); das sind die NATS-Commits der Vollabdeckungs-Welle, und ihre Parents tragen `0.1.0` (`git show <Commit>^:<Quelle>`) — die Hebung `0.1.0` → `0.2.0` liegt je in diesem einen Commit, kein späterer Commit hob erneut. Beide Tags liegen vor diesen Commits (`git merge-base --is-ancestor <Tag> <Commit>` Exit 0), also ist `0.2.0` in allen drei Quellen unveröffentlicht: eine additive Erweiterung faltet sich in diese anstehende Version. `origin` ist additiv und rückwärtskompatibel (SemVer-Minor); `0.2.0` ist damit bereits die nächste Minor-Version nach dem letzten Tag und trägt das Feld im ersten Release, der `0.2.0` nennt. Eine Hebung auf `0.3.0` übersprünge eine Version, die kein Konsument je sah. Die Release-Tag-Info-Skripte prüfen nur die Grammatik des Tags (`sdk-<sprache>-v<Version>`); `bash tools/harness/sdk-{csharp,kotlin,python}-release-tag-info.sh sdk-<sprache>-v0.2.0` liefert je `version=0.2.0`, und der Tag-Abgleich der Workflows liest die Quellen (Kotlin: die Zeile `^version = "…"`), die dieser Slice nicht ändert. Für den Release-Zug nach der Wellen-Closure: Tag = Quell-Version = `sdk-csharp-v0.2.0`, `sdk-python-v0.2.0`, `sdk-kotlin-v0.2.0`. Folge: die Träger mit `0.2.0` (Spec-§6-Zeilen, Artefaktnamen in Spec/`harness/mk/sdk.mk`/Dockerfiles, Kotlin-README) bleiben unverändert und wahr; die Artefaktnamen werden durch die realen Pack-Läufe dieses Slice bestätigt. |
| `sdks/python/pgchangefeed/src/pgchangefeed/http_client.py` | Reduktion (kein Diff) | Der Client liest `GET /changes` über `ReadChangesResponse.from_json` → `Change.from_json` in `models.py`; `http_client.py` trägt keine Feldliste und ändert sich nicht. |
| Kotlin `Change`: `wireOrigin: String? = null` plus abgeleitete Eigenschaft `origin` (`wireOrigin ?: "wal"`); C# `Origin = "wal"` als Parameter-Default plus `OriginConverter` (`HandleNull`, liest JSON-`null` als `wal`); Python `origin: str = "wal"` und `"wal" if data.get("origin") is None else data["origin"]` | Ausführungsform | Gson setzt bei einer Klasse ohne Konstruktor-Aufruf den Parameter-Default nicht (Risiko §6, der Test „ohne `origin`" belegt es); der Wert bleibt der Server-String (kein Enum). Die Regel gilt in allen drei Sprachen gleich: Feld fehlt oder JSON-`null` → `wal`; jeder andere Server-String (auch leer oder unbekannt) wird unverändert durchgereicht (Entscheidung des Auftraggebers zu Review-Finding F-3). |
| C# `Http/Models/Changes.cs` (`OriginConverter`, Test-Fälle `origin:null` und `origin:""`), Python `models.py` (`is None`-Prüfung, Test-Fall `empty-value`), Kotlin-Test-Fall „empty origin" | update (Fixrunde F-3) | Die drei Modelle stimmen an den Randwerten überein: JSON-`null` liest als `wal` (C#: `OriginConverter`), ein leerer String bleibt leer (Python: `is None`-Prüfung statt `or`); je Änderung eine Mutation der Eingabeseite rot gesehen (Bericht). |
| Kotlin `wireOrigin` als öffentlicher Konstruktorparameter | belassen (Fixrunde F-4), Risiko §6 | Gson ruft keinen Konstruktor auf und übernimmt keinen Kotlin-Default (Mutation M5 des Reviews: `@SerializedName("origin") val origin: String = "wal"` färbt „without origin" und „JSON-null" rot); eine Normalisierung im Konstruktor greift damit nicht, ein Lesepfad ohne `wireOrigin` verlangte einen eigenen Gson-`TypeAdapter` für `Change` statt des einfachen `Gson()` in `PgChangeFeedHttpClient.kt` — kein kleiner Eingriff. `equals`/`copy`/`component13()` arbeiten auf `wireOrigin`; die KDoc nennt es. |
| `spec/pflichtenheft.md` `LH-FA-SST-009.a` (Zeile 288) und §7 | update (Fixrunde F-5) | die Aufzählung der bestehenden Zustellwege nennt `SPEC-022` (die HTTP-Fläche der Packages deckt `GET /changes`); Historie-Zeile ohne ADR-/Slice-Bezug; die `origin`-Zeile trägt die Regel mit JSON-`null`. |
| `docs/user/benutzerhandbuch.md` (Python-Absatz Zeile 1030–1036, `origin`-Absatz, Version 1.59, Historie) | update (Fixrunde F-6) | der Python-Absatz führt SSE und den NATS-Vollinhalts-Stream als vom Package getragen (Python `0.2.0` trägt beide: README, `SPEC-027`); der `origin`-Absatz nennt die Regel; Suchlauf nach überholten Zukunftsaussagen im Suchlauf-Feld. |
| Kommentare „Feldzahl der HTTP-Lesemodelle": `Http/Models/Changes.cs`, `http/model/Changes.kt`, `models.py` (SSE-Abschnittskommentar) (Fixrunde F-7) | update | die Live-Flächen tragen zehn Felder, das HTTP-Lesemodell dreizehn (zehn plus `commit_position`, `committed_at`, `origin`); der Domänentyp `model.Change` trägt selbst `Origin` (elf Felder, `internal/domain/model/change.go`), deshalb nennen die Kommentare „the ten the live surfaces carry" statt „the ten of the domain type". |
| Docstring/KDoc/XML-Doc der `origin`-Regel (C#, Kotlin, Python) und die drei SDK-READMEs (Fixrunde F-3) | update | die Regel (JSON-`null` → `wal`, sonstiger Wert unverändert) steht an allen Trägern in derselben Fassung. |
| `spec/pflichtenheft.md` §6 (`SPEC-026`/`-027`/`-028`) und §7 | update (abweichend von „eine Historie-Zeile je Package") | Die Zeilen nennen `SPEC-022` unter den gedeckten Drahtverträgen (`origin` ist dessen Antwortfeld); **eine** Historie-Zeile für alle drei Packages, weil der Inhalt identisch ist; die Version wird nicht gehoben, die Zeilen tragen weiter `0.2.0`. |
| `docs/user/benutzerhandbuch.md` | update | Ein Absatz nach den drei `**SDK:**`-Absätzen der HTTP-Beschreibung (nennt `origin` für alle drei Packages, `wal` bei fehlendem Feld, Live-Wege ohne), Version 1.58 samt Historie-Zeile. |
| Beispiel-Clients `examples/**` | entfällt | Kein Beispiel dekodiert `GET /changes` (Suchlauf unten). |

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaften: „die Version der drei Packages (0.2.0 → neue Version)", „die Feldmenge des HTTP-Lesemodells"; beide Stände gemessen):**

| Träger | Suchbefehl | Befund | Behandlung |
|---|---|---|---|
| Versionsangaben | `git grep -n '0\.2\.0' <Stand> -- spec harness sdks docs/user tools .github Makefile` (Stand: Parent `fea14159`, Stand vor der Fixrunde `e80b4f64`; der Diff-Stand ist derselbe Befehl ohne `<Stand>` am Arbeitsbaum; `git grep` liest nur getrackte Dateien, `dist/`/`bin/`/`obj/`/`build/` sind ausgeschlossen; `grep -rn` am Arbeitsbaum mit `--exclude-dir` für diese vier zählt gleich) — Treffer nach Bezug auf die drei Packages lesen | **Gefunden, gemessen 2026-09-25:** 29 Zeilen an **allen drei Ständen** (Parent `fea14159` 29, `e80b4f64` 29, Diff-Stand 29; `grep -rn` am Arbeitsbaum 29), kein Treffer verändert. Je Datei (Parent-Zeilen): `spec/pflichtenheft.md` 9 (Z. 302/308/315 Fließtext, Z. 820–822 §6-Zeilen, Z. 864–866 Historie), `docs/user/benutzerhandbuch.md` 4 (Z. 1780 Historie 1.32 ohne Package-Bezug; Z. 1786/1787/1791 Historie 1.38/1.39/1.43, frühere Einträge, bleiben), `harness/README.md` 3 (Z. 157–159, die drei `make sdk-pack-*`-Zeilen mit den Artefaktnamen `PgChangeFeed.Client.0.2.0.nupkg`, `pgchangefeed-0.2.0.tar.gz`, `pgchangefeed-kotlin-0.2.0.jar`), `harness/mk/sdk.mk` 2 (Z. 27, 49), `sdks/csharp/Dockerfile` 1 (Z. 66), `sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.csproj` 2 (Z. 15/23), `sdks/kotlin/Dockerfile` 2 (Z. 77/104), `sdks/kotlin/pgchangefeed-kotlin/README.md` 1 (Z. 38), `sdks/kotlin/pgchangefeed-kotlin/build.gradle.kts` 3 (Z. 97/108/230), `sdks/python/pgchangefeed/pyproject.toml` 1 (Z. 7), `spec/lastenheft.md` 1 (Z. 1361, Dokumentversion ohne Package-Bezug); Summe 9 + 4 + 3 + 2 + 1 + 2 + 2 + 1 + 3 + 1 + 1 = 29. **Nicht gefunden:** kein Träger nennt eine andere Package-Version als `0.2.0` (außer den frühen Historie-Zeilen mit `0.1.0`, siehe Zeile Artefaktnamen); `docs/user/version.md` und `docs/user/releasing.md` tragen keine `0.2.0` (`git grep -c` liefert für beide keine Zeile, Parent und Diff-Stand); `tools/`, `.github/` und `Makefile` tragen keine Zeile. Der Suchraum enthält weder die Plan-Datei (`docs/plan/`) noch `docs/reviews/`, es gibt also keinen Selbstverweis dieses Feldes im Treffer-Satz. | kein Nachzug (Version unverändert); `docs/reviews/**` und `done/**` sind Records und bleiben |
| Artefaktnamen (`PgChangeFeed.Client.0.2.0.nupkg`, `pgchangefeed-kotlin-0.2.0.jar`, `pgchangefeed-0.2.0-…`) | `git grep -c 'nupkg\|\.jar\|\.whl' <Stand> -- harness docs/user spec` (Stand: `fea14159`, `e80b4f64`, Arbeitsbaum) | **Gefunden, gemessen 2026-09-25, identisch an allen drei Ständen:** `harness/mk/sdk.mk` 8, `spec/pflichtenheft.md` 9, `harness/README.md` 5, `docs/user/releasing.md` 2 Zeilen. Die Namen tragen `0.2.0` in `spec/pflichtenheft.md` (Z. 302/308/315, 864–866), `harness/mk/sdk.mk` (Z. 27, 49) und `harness/README.md` (Z. 157–159: `PgChangeFeed.Client.0.2.0.nupkg`, `pgchangefeed-0.2.0.tar.gz`, `pgchangefeed-kotlin-0.2.0.jar`); `harness/README.md` Z. 146/147 und `docs/user/releasing.md` Z. 148/196 nennen `.nupkg`/`.whl` ohne Version. **Real bestätigt** (Läufe der Fixrunde, Exit je 0, `ls` der `dist/`-Verzeichnisse): `make sdk-pack-csharp` erzeugt `sdks/csharp/dist/PgChangeFeed.Client.0.2.0.nupkg`, `make sdk-pack-python` `pgchangefeed-0.2.0-py3-none-any.whl` und `pgchangefeed-0.2.0.tar.gz`, `make sdk-pack-kotlin` `pgchangefeed-kotlin-0.2.0.jar`; die Lauf-Kennungen der ersten Implementer-Läufe sind im Repo nicht auflösbar, der Review-Lauf bestätigt dieselben Namen (`docs/reviews/review-slice-backfill-sdk-origin.md`, Abschnitt „Pack-Läufe"). **Nicht gefunden:** kein Artefaktname mit einer anderen Version als `0.2.0` außer den drei Historie-Zeilen `spec/pflichtenheft.md` Z. 859/861/863 (`0.1.0`, frühere Einträge, bleiben); keine Änderung nötig. | unverändert; durch reale `make sdk-pack-*`-Läufe bestätigt |
| Feldbeschreibungen des HTTP-Lesemodells in READMEs und Handbuch | `grep -rn 'old_image\|oldImage\|OldImage' sdks docs/user` (Fundstellen nach HTTP-Bezug lesen) | **Gefunden, Parent `fea14159`:** in `sdks/**` nur Quelltext und Tests (Symbolnamen), keine README-Feldliste (die READMEs nennen `GET /changes` nur als Fähigkeit); `docs/user/benutzerhandbuch.md` Z. 979 (Feldliste von `GET /changes`, trägt `origin` bereits), Z. 1080 (gRPC-Feldtabelle) und Z. 1313 (SSE-Feldliste) — die beiden letzten sind Live-Wege. Zusätzlich gesucht nach Zählwörtern der Feldmenge (`twelve`/`eleven`/`zwölf`/`ten fields`/`zehn`) an beiden Ständen: Parent trägt sie in den sechs Kommentaren (`Http/Models/Changes.cs` Z. 8 „ten", `Sse/Models/Change.cs` Z. 17 und `Nats/Models/Change.cs` Z. 21 „eleven", `http/model/Changes.kt` Z. 8 „ten", `sse/model/Change.kt` Z. 19 „twelve", `models.py` Z. 264 „zwölf"); Diff-Stand: alle sechs nennen dreizehn. **Nicht gefunden / unverändert:** die „zehn Felder" der Live-Flächen (gRPC/SSE/NATS in Handbuch, Spec §2 Z. 594/613/689, Kommentare, Tests) bleiben, ebenso „zehn Fähigkeiten" (Fähigkeiten der HTTP-Oberfläche, keine Felder); keine README-Feldliste zu ergänzen — die READMEs nennen `origin` in ihrer HTTP-Beschreibung, das Handbuch führt die Feldliste. | Handbuch-Absatz und README-Sätze ergänzt; die sechs Zählwort-Kommentare auf dreizehn gezogen |
| Beispiel-Clients, die `GET /changes` dekodieren | Lesen von `examples/**` | **Nicht gefunden, Parent und Diff-Stand:** `grep -rn -i -E 'commit_position|commitposition|"changes"|ReadChanges|readChanges|read_changes' examples` liefert an beiden Ständen keinen Treffer; die HTTP-Beispiele (Go/C#/Kotlin) rufen `GET /tables` auf, die NATS-Beispiele bauen nur die Abfrage-URL für `/changes` (`examples/nats-client/subject.go` Z. 20) und dekodieren keine Antwort, SSE/NATS-Vollinhalt sind Live-Wege. Kein Beispiel dekodiert `GET /changes`. | Punkt entfällt, keine Änderung an `examples/**`; kein `make examples-*`-Lauf nötig |
| Publish-Workflows (Tag-Abgleich gegen die Metadaten-Quelle) | Lesen der drei `sdk-*-release.yml` (nur lesen) | **Gelesen, unverändert:** der Tag-Abgleich extrahiert die Version aus `<Version>` der `.csproj` (`sdk-csharp-release.yml` Z. 68), aus `^version = "…"` der `pyproject.toml` (`sdk-python-release.yml` Z. 85) und der `build.gradle.kts` (`sdk-kotlin-release.yml` Z. 110; die eingerückte Zeile im `publishing`-Block trifft das `^`-Muster nicht). Gemessen mit denselben Befehlen an den Quellen: `cs=0.2.0`, `py=0.2.0`, `kt=0.2.0`; eine auf `0.3.0` verstellte Kopie der `pyproject.toml` liefert `py_mut=0.3.0` (Abweichung zum Tag `0.2.0` wäre sichtbar). | unverändert; die Quellen tragen `0.2.0`, Tag = `sdk-<sprache>-v0.2.0` |
| **Fixrunde** — bewegte Eigenschaft: „die Regel für `origin` in den SDK-Lesemodellen (fehlendes Feld oder JSON-`null` → `wal`, sonstiger Server-String unverändert)" | `git grep -n -i -E '(liest\|lesen\|read\|reads\|gilt\|defaults?) (das Package \|the response )?(als\|as\|to) (.\|<c>)?wal\|without the field\|ohne dieses Feld\|fehlender Wert\|fehlendes Feld\|Default .wal' <Stand> -- sdks docs/user spec harness tools examples` (Stand: `e80b4f64` und Arbeitsbaum) | **Gefunden, gemessen 2026-09-25:** an `e80b4f64` 21 Zeilen, am Diff-Stand 23 Zeilen (das Muster sucht Formulierungen; die Zahl ändert sich mit der Umformulierung der Träger — Handbuch 5 → 7 durch umbrochenen Absatz und Historie-Zeile 1.59, C# 1 → 2, Kotlin 2 → 1 —, nicht durch neue Träger). Träger der SDK-Regel, angepasst: `sdks/csharp/PgChangeFeed.Client/Http/Models/Changes.cs` (XML-Doc), `http/model/Changes.kt` (KDoc), `sdks/python/pgchangefeed/src/pgchangefeed/models.py` Z. 227 (Docstring), die drei `sdks/*/README.md` Z. 9, `docs/user/benutzerhandbuch.md` Z. 1057 (Absatz) samt Historie, `spec/pflichtenheft.md` Historie-Zeile. **Nicht gefunden / unverändert:** die Server-Regel für gespeicherte Änderungen (`spec/pflichtenheft.md` Z. 354/645, `docs/user/benutzerhandbuch.md` Z. 395/981, `tools/schema/schema.yaml` Z. 167 „NULL liest als `wal`") betrifft die SQL-/Speicher-Seite, nicht das SDK-Lesemodell, und bleibt; Treffer außerhalb (`fehlender Wert` in `examples/**`-CliTests und `harness/erfassung-feldliste.md`, `tools/harness/natssub/main.go`) betreffen andere Gegenstände. | die sechs SDK-Träger und die Handbuch-Fassung nachgezogen |
| **Fixrunde** — bewegte Eigenschaft: „welche Zustellwege das Python-Package trägt (SSE und NATS-Vollinhalt)" | `git grep -n -i 'folgen\|im selben Folge-Release\|noch nicht' <Stand> -- docs/user/benutzerhandbuch.md` (Stand: `e80b4f64`, Arbeitsbaum); zusätzlich `git grep -n -i -E 'noch nicht\|folgen\|Folge-Release\|vorerst\|demnächst\|geplant\|später' -- sdks/*/README.md` | **Gefunden:** Handbuch an beiden Ständen je 12 Zeilen. `e80b4f64`: 8 Betriebs-Zustände (Z. 206/237/270/506/549/726/804/821, kein Package-Bezug), Z. 1034 (Python-Absatz: „SSE und der NATS-Vollinhalts-Stream folgen im selben Folge-Release" — die einzige überholte Zukunftsaussage über ein Package) und drei Historie-Zeilen (1.40/1.42/1.56). Diff-Stand: dieselben 8 Betriebs-Zustände und vier Historie-Zeilen (1.40/1.42/1.56/1.59, Records, bleiben); der Python-Absatz führt SSE und NATS-Vollinhalt als getragen (Verweis auf die beiden Abschnitte). Der Kotlin-Absatz nennt alle drei Streams als getragen, der C#-Absatz führt nur die HTTP-Fläche und keine Zukunftsaussage. **Nicht gefunden:** in den drei `sdks/*/README.md` kein Treffer (Zählung 0). | Python-Absatz nachgezogen; Kotlin/C#: kein Nachzug |
| **Fixrunde** — bewegte Eigenschaft: „die Zuordnung der Feldmenge im HTTP-Lesemodell (zehn Live-Felder plus drei)" | `git grep -n -c -E 'of the (Go )?domain type\|same ten fields as the\|domain type' <Stand> -- sdks docs/user spec harness examples` (Stand: `e80b4f64`, Arbeitsbaum) | **Gefunden, gemessen 2026-09-25:** `e80b4f64` `Http/Models/Changes.cs` 1, `http/model/Changes.kt` 1, `models.py` 2 Zeilen (vier Zeilen in drei Dateien); Diff-Stand: kein Treffer in `sdks docs/user spec harness examples`. **Nicht gefunden:** `Sse/Models/Change.cs`, `Nats/Models/Change.cs`, `sse/model/Change.kt` tragen „thirteen fields" für das HTTP-Modell und „the ten fields" für die Live-Fläche, ohne den Domänentyp zu nennen; `internal/domain/model/change.go` trägt elf Felder (Go, nicht Träger). | die drei Kommentare nachgezogen |
| **Fixrunde** — bewegte Eigenschaft: „die Aufzählung der durch die Packages gedeckten Drahtverträge" | `git grep -n -E 'SPEC-018.{4}SPEC-020' e80b4f64 -- spec harness docs/user` und derselbe Aufruf ohne Stand-Argument am Arbeitsbaum | **Gefunden:** 6 Zeilen an beiden Ständen: `spec/pflichtenheft.md` Z. 288 (`LH-FA-SST-009.a` „bestehende Zustellwege", ohne `SPEC-022` am Parent), Z. 655 und 671 (Beispiel-Clients, tragen `SPEC-022`), Z. 820–822 (`SPEC-026`/`-027`/`-028`, tragen `SPEC-022`). Nur Z. 288 wird geändert (`SPEC-022` ergänzt); die Aufzählung in Z. 288 setzt sich in Z. 289 fort. **Nicht gefunden:** kein weiterer Träger der Aufzählung in `harness/` und `docs/user/`. | Z. 288 nachgezogen |

## 4. Trigger

**Start** (`next` → `in-progress`): wenn `change-origin` in `done/` liegt (das
Feld existiert am Draht) und kein anderer Slice in `in-progress/` liegt
(WIP-Limit 1). Von `snapshot-reader` bis `bench-richtgroesse` technisch
unabhängig; die Reihenfolge der Welle stellt ihn ans Ende.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): falls drei Sprachen mit
  Version, Spec-Zeilen und Handbuch nicht in einem Review tragen — der
  abtrennbare Teil ist je ein Slice je Sprache (dann `sdk-origin-csharp`,
  `-kotlin`, `-python`).
- `in-progress` → `open` (blockiert): falls ein Decoder ein unbekanntes Feld nicht
  ignoriert (dann ein Befund gegen die Annahme von [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md), mit
  Architect-Frage nach der Versionierung des Drahtvertrags).

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + die drei `make sdk-pack-*`-Läufe real
grün + Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- **Ein Decoder ignoriert das neue Feld nicht.** [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) hält fest, dass
  C# `System.Text.Json`, Kotlin Gson und Python (`json` + `from_json`) unbekannte
  Felder ignorieren (gelesen); ein realer Beleg gegen den Server fehlt. *Erwartet,
  zu belegen durch:* die Unit-Tests mit einer Fixture aus einer realen Antwort
  (Ursprung im Test genannt). **Ausgang:** *(bei Closure)*
- **Kotlin/Gson und der Default `wal`.** Gson setzt bei einem fehlenden Feld auf
  eine Kotlin-Klasse mit Default-Parameter nicht den Default (erwartet: `null`);
  das Modell muss `null` als `wal` lesen. *Erwartet, zu belegen durch:* der Test
  „Antwort ohne `origin`". **Ausgang:** *(bei Closure)*
- **Träger nennen die alte Version** (`BEO-PGC/arbeit-ueberholt-stehenden-traeger`,
  verkörpert, 26×): Spec-§6-Zeilen, `harness/README.md`, Kommentare in
  `harness/mk/sdk.mk`, SDK-READMEs, Handbuch. *Erwartet, zu belegen durch:* der
  Suchlauf mit beiden Ständen. **Ausgang:** *(bei Closure)*
- **Zahl im Träger ohne Ursprung** ([`AGENTS.md`](../../../../AGENTS.md) §3.12 Instanz A): die Artefaktnamen
  mit Version sind Läufe. *Erwartet, zu belegen durch:* die Zeilen entstehen aus
  realen `make sdk-pack-*`-Läufen dieses Slice. **Ausgang:** *(bei Closure)*
- **Sprachfragmente aus der Vorlage** (siehe DoD 3): eine Kopie der Form über
  drei Sprachen trägt ein deutsches Fragment leicht weiter. *Erwartet, zu belegen
  durch:* der Reviewer-Punkt zur Form-Vorbild-Kopie. **Ausgang:** *(bei Closure)*
- **Der Docker-Kontext `proto`** ist für jeden Bau der SDKs zwingend
  (`BEO-PGC/zusatzkontext-kopplung-breiter-als-dod-wortlaut`, offen, 3×); ein
  Bau ohne ihn bricht an der `COPY --from=proto`-Zeile ab — der Slice nennt ihn in
  jedem Aufruf und hält die Klasse damit sichtbar. **Ausgang:** *(bei Closure)*
- **Version ohne Tag** — die Metadaten-Quelle ist höher als jeder vorhandene Tag;
  der Tag-Abgleich der Publish-Workflows läuft erst beim Tag-Push. Kein
  Workflow-Zug, [`AGENTS.md`](../../../../AGENTS.md) §3.10 greift nicht. **Ausgang:** *(bei Closure)*
- **Kotlin `Change` trägt den Roh-Wert als öffentlichen Konstruktorparameter**
  (`wireOrigin: String?`): `equals`, `hashCode`, `toString`, `copy` und
  `component13()` arbeiten auf dem Roh-Wert, zwei Changes mit fehlendem Feld
  und mit `origin = "wal"` sind ungleich (Review-Finding F-4). Bewusst
  akzeptierte, in der KDoc genannte Eigenschaft. *Beleg:* Mutation M5 des
  Reviews — `@SerializedName("origin") val origin: String = "wal"` als
  Konstruktorparameter färbt „reads a response without origin as wal" und
  „reads a JSON-null origin as wal" rot, weil Gson keinen Konstruktor aufruft
  und keinen Kotlin-Default übernimmt; ein Lesepfad ohne `wireOrigin` verlangt
  einen eigenen `TypeAdapter` statt des einfachen `Gson()` in
  `PgChangeFeedHttpClient.kt`. **Ausgang:** *(bei Closure)*

## 7. Closure-Notiz

- **Was hat funktioniert:** *(zu tragen bei Closure)*
- **Was ging anders als geplant:** *(zu tragen bei Closure)*
- **Steering-Loop-Eintrag (Lerneintrag):** *(zu tragen bei Closure —
  geschärfte Regel · neuer Sensor · benannte Spec-Lücke; ohne ihn kein
  `done/`-Übergang)*
- **Beobachtungs-Register (`../observations/`):** *(je Anfall Beleg oder
  „keine Beobachtung angefallen" als notierte Antwort)*
- **Folge-Slices:** *(zu tragen bei Closure)*
- **Risiken aus §6:** *(je ein Ausgang)*
- **Drei Paarungen:** dieser Slice gehört zu [welle-backfill-bestand](../welle-backfill-bestand.md) (offen) — die
  Prüfung läuft regelkonform bei deren Closure.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt sind die SDK-Bäume `sdks/csharp/`,
`sdks/kotlin/` und `sdks/python/` — je bereits mit ihrem Projektgerüst eröffnet
(Greenfield); keine erneute Ausdifferenzierung nötig.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
`BEO-PGC/zusatzkontext-kopplung-breiter-als-dod-wortlaut` (offen, 3×, Ausgang
nicht zugewiesen, einschlägig — DoD 1 und Risiko §6),
`BEO-PGC/formvorbild-kopie-traegt-deutsches-wortfragment-weiter` (verkörpert, 3×)
und `BEO-PGC/deutsches-fachwort-im-englischen-sdk-readme` (verkörpert, 3×) —
DoD 3, `BEO-PGC/arbeit-ueberholt-stehenden-traeger` (verkörpert, 26×, Suchlauf
§3), `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (verkörpert, 13×,
Artefaktnamen), `BEO-PGC/release-mechanismus-nicht-in-releasing-doku-nachgezogen`
(offen, 1×, gesichtet — der Slice berührt keinen Release-Mechanismus),
`BEO-PGC/github-actions-unverifizierbar-lokal` (verkörpert, 7× — nicht
einschlägig: kein Workflow-Zug), `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht`
(verkörpert, 8×, die `make sdk-pack-*`-Zeilen beschreiben nur gefahrene Läufe).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
