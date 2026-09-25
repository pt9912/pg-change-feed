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
      `wal`; je Sprache ein Unit-Test mit beiden Fällen (dazu unbekannter Wert
      durchgereicht; Kotlin/Python: JSON-`null` liest als `wal`) und einer
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
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
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
| **Versionsentscheidung: keine Hebung** — `.csproj` `<Version>`, `build.gradle.kts` `version` (Zeile 108, die eingerückte Zeile 230 im `publishing`-Block trägt dieselbe Zahl), `pyproject.toml` `version` bleiben `0.2.0` | Reduktion (Plan-Punkt „Version heben" entfällt) | *Beleg (gemessen 2026-09-25):* `git tag -l` nennt `sdk-csharp-v0.1.0` und `sdk-python-v0.1.0` (Kotlin: kein Tag); die Quellen stehen bei `0.2.0`, gesetzt in `d86d1965` (C#), `0f8cc4f2` (Python) und `c8c9e3ae` (Kotlin) — je im letzten Slice der Vollabdeckungs-Welle, und der jeweils zweite Slice (NATS nach SSE) hob **nicht** erneut: eine additive Erweiterung faltet sich in die anstehende, noch unveröffentlichte Version (Muster der Vorgänger). `origin` ist additiv und rückwärtskompatibel (SemVer-Minor); `0.2.0` ist damit bereits die nächste Minor-Version nach dem letzten Tag und trägt das Feld im ersten Release, der `0.2.0` nennt. Eine Hebung auf `0.3.0` übersprünge eine Version, die kein Konsument je sah. Die Release-Tag-Info-Skripte prüfen nur die Grammatik des Tags (`sdk-<sprache>-v<Version>`); `bash tools/harness/sdk-{csharp,kotlin,python}-release-tag-info.sh sdk-<sprache>-v0.2.0` liefert je `version=0.2.0`, und der Tag-Abgleich der Workflows liest die Quellen (Kotlin: die Zeile `^version = "…"`), die dieser Slice nicht ändert. Für den Release-Zug nach der Wellen-Closure: Tag = Quell-Version = `sdk-csharp-v0.2.0`, `sdk-python-v0.2.0`, `sdk-kotlin-v0.2.0`. Folge: die Träger mit `0.2.0` (Spec-§6-Zeilen, Artefaktnamen in Spec/`harness/mk/sdk.mk`/Dockerfiles, Kotlin-README) bleiben unverändert und wahr; die Artefaktnamen werden durch die realen Pack-Läufe dieses Slice bestätigt. |
| `sdks/python/pgchangefeed/src/pgchangefeed/http_client.py` | Reduktion (kein Diff) | Der Client liest `GET /changes` über `ReadChangesResponse.from_json` → `Change.from_json` in `models.py`; `http_client.py` trägt keine Feldliste und ändert sich nicht. |
| Kotlin `Change`: `wireOrigin: String? = null` plus abgeleitete Eigenschaft `origin` (`wireOrigin ?: "wal"`); C# `Origin = "wal"` als Parameter-Default; Python `origin: str = "wal"` und `data.get("origin") or "wal"` | Ausführungsform | Gson setzt bei einer Klasse ohne Konstruktor-Aufruf den Parameter-Default nicht (Risiko §6, der Test „ohne `origin`" belegt es); der Wert bleibt der Server-String (kein Enum), ein unbekannter Wert wird durchgereicht statt zu scheitern; ein JSON-`null` liest in Kotlin und Python als `wal`. |
| `spec/pflichtenheft.md` §6 (`SPEC-026`/`-027`/`-028`) und §7 | update (abweichend von „eine Historie-Zeile je Package") | Die Zeilen nennen `SPEC-022` unter den gedeckten Drahtverträgen (`origin` ist dessen Antwortfeld); **eine** Historie-Zeile für alle drei Packages, weil der Inhalt identisch ist; die Version wird nicht gehoben, die Zeilen tragen weiter `0.2.0`. |
| `docs/user/benutzerhandbuch.md` | update | Ein Absatz nach den drei `**SDK:**`-Absätzen der HTTP-Beschreibung (nennt `origin` für alle drei Packages, `wal` bei fehlendem Feld, Live-Wege ohne), Version 1.58 samt Historie-Zeile. |
| Beispiel-Clients `examples/**` | entfällt | Kein Beispiel dekodiert `GET /changes` (Suchlauf unten). |

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaften: „die Version der drei Packages (0.2.0 → neue Version)", „die Feldmenge des HTTP-Lesemodells"; beide Stände gemessen):**

| Träger | Suchbefehl | Befund | Behandlung |
|---|---|---|---|
| Versionsangaben | `grep -rn '0\.2\.0' spec harness sdks docs/user` (Treffer nach Bezug auf die drei Packages lesen; andere Versionen sind keine Treffer) | **Gefunden, Parent `fea14159`** (`git grep … fea14159`): Package-Bezug in `spec/pflichtenheft.md` (Z. 302/308/315 Fließtext, Z. 820–822 §6-Zeilen, Z. 864–866 Historie), `harness/mk/sdk.mk` (Z. 27, 49), `sdks/csharp/Dockerfile` Z. 66, `sdks/kotlin/Dockerfile` Z. 77/104, `sdks/kotlin/pgchangefeed-kotlin/README.md` Z. 38, `build.gradle.kts` Z. 97/108/230, `.csproj` Z. 15/23, `pyproject.toml` Z. 7, `docs/user/benutzerhandbuch.md` Historie 1.38/1.39/1.43 (Z. 1786/1787/1791, frühere Einträge, bleiben); ohne Package-Bezug (kein Treffer): `spec/lastenheft.md` Z. 1361 (Dokumentversion), `docs/user/benutzerhandbuch.md` Z. 1780 (Historie 1.32, Software-Version-Kopf). **Diff-Stand:** `grep -rn '0\.2\.0' spec harness sdks docs/user tools .github Makefile` (Ausschluss `dist/`/`bin/`/`obj/`/`build/`) liefert 26 Zeilen — dieselbe Zahl wie am Parent, kein Treffer verändert (Handbuch-Eintrag 1.58 nennt `0.2.0` nicht). **Nicht gefunden:** kein Träger nennt eine andere Package-Version als `0.2.0`; `docs/user/version.md` und `docs/user/releasing.md` tragen keine SDK-Version (`grep -c` je 0). | kein Nachzug (Version unverändert); `docs/reviews/**` und `done/**` sind Records und bleiben |
| Artefaktnamen (`PgChangeFeed.Client.0.2.0.nupkg`, `pgchangefeed-kotlin-0.2.0.jar`, `pgchangefeed-0.2.0-…`) | `grep -rn 'nupkg\|\.jar\|\.whl' harness docs/user spec` | **Gefunden, Parent und Diff-Stand** (`grep -c` je Datei, identisch an beiden Ständen): `harness/mk/sdk.mk` 8, `spec/pflichtenheft.md` 9, `harness/README.md` 5, `docs/user/releasing.md` 2 Zeilen. Die Namen tragen `0.2.0` (Spec, `sdk.mk`) bzw. keine Version (`harness/README.md`, `releasing.md`). **Real bestätigt:** `make sdk-pack-csharp` erzeugt `sdks/csharp/dist/PgChangeFeed.Client.0.2.0.nupkg`, `make sdk-pack-python` `pgchangefeed-0.2.0-py3-none-any.whl` und `pgchangefeed-0.2.0.tar.gz`, `make sdk-pack-kotlin` `pgchangefeed-kotlin-0.2.0.jar` (Läufe dieses Slice, Zeilen im Bericht). **Nicht gefunden:** kein Artefaktname mit einer anderen Version als `0.2.0` außer den drei Historie-Zeilen `spec/pflichtenheft.md` Z. 859/861/863 (`0.1.0`, frühere Einträge, bleiben); keine Änderung nötig. | unverändert; durch die realen `make sdk-pack-*`-Läufe bestätigt |
| Feldbeschreibungen des HTTP-Lesemodells in READMEs und Handbuch | `grep -rn 'old_image\|oldImage\|OldImage' sdks docs/user` (Fundstellen nach HTTP-Bezug lesen) | **Gefunden, Parent `fea14159`:** in `sdks/**` nur Quelltext und Tests (Symbolnamen), keine README-Feldliste (die READMEs nennen `GET /changes` nur als Fähigkeit); `docs/user/benutzerhandbuch.md` Z. 979 (Feldliste von `GET /changes`, trägt `origin` bereits), Z. 1080 (gRPC-Feldtabelle) und Z. 1313 (SSE-Feldliste) — die beiden letzten sind Live-Wege. Zusätzlich gesucht nach Zählwörtern der Feldmenge (`twelve`/`eleven`/`zwölf`/`ten fields`/`zehn`) an beiden Ständen: Parent trägt sie in den sechs Kommentaren (`Http/Models/Changes.cs` Z. 8 „ten", `Sse/Models/Change.cs` Z. 17 und `Nats/Models/Change.cs` Z. 21 „eleven", `http/model/Changes.kt` Z. 8 „ten", `sse/model/Change.kt` Z. 19 „twelve", `models.py` Z. 264 „zwölf"); Diff-Stand: alle sechs nennen dreizehn. **Nicht gefunden / unverändert:** die „zehn Felder" der Live-Flächen (gRPC/SSE/NATS in Handbuch, Spec §2 Z. 594/613/689, Kommentare, Tests) bleiben, ebenso „zehn Fähigkeiten" (Fähigkeiten der HTTP-Oberfläche, keine Felder); keine README-Feldliste zu ergänzen — die READMEs nennen `origin` in ihrer HTTP-Beschreibung, das Handbuch führt die Feldliste. | Handbuch-Absatz und README-Sätze ergänzt; die sechs Zählwort-Kommentare auf dreizehn gezogen |
| Beispiel-Clients, die `GET /changes` dekodieren | Lesen von `examples/**` | **Nicht gefunden, Parent und Diff-Stand:** `grep -rn -i -E 'commit_position|commitposition|"changes"|ReadChanges|readChanges|read_changes' examples` liefert an beiden Ständen keinen Treffer; die HTTP-Beispiele (Go/C#/Kotlin) rufen `GET /tables` auf, die NATS-Beispiele bauen nur die Abfrage-URL für `/changes` (`examples/nats-client/subject.go` Z. 20) und dekodieren keine Antwort, SSE/NATS-Vollinhalt sind Live-Wege. Kein Beispiel dekodiert `GET /changes`. | Punkt entfällt, keine Änderung an `examples/**`; kein `make examples-*`-Lauf nötig |
| Publish-Workflows (Tag-Abgleich gegen die Metadaten-Quelle) | Lesen der drei `sdk-*-release.yml` (nur lesen) | **Gelesen, unverändert:** der Tag-Abgleich extrahiert die Version aus `<Version>` der `.csproj` (`sdk-csharp-release.yml` Z. 68), aus `^version = "…"` der `pyproject.toml` (`sdk-python-release.yml` Z. 85) und der `build.gradle.kts` (`sdk-kotlin-release.yml` Z. 110; die eingerückte Zeile im `publishing`-Block trifft das `^`-Muster nicht). Gemessen mit denselben Befehlen an den Quellen: `cs=0.2.0`, `py=0.2.0`, `kt=0.2.0`; eine auf `0.3.0` verstellte Kopie der `pyproject.toml` liefert `py_mut=0.3.0` (Abweichung zum Tag `0.2.0` wäre sichtbar). | unverändert; die Quellen tragen `0.2.0`, Tag = `sdk-<sprache>-v0.2.0` |

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
