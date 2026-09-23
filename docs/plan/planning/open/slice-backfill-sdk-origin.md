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

**Verantwortlich:** — (noch nicht priorisiert).

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

- [ ] Die drei HTTP-Lesemodelle tragen `origin` als optionales Feld: eine
      Antwort mit `origin` liefert den Wert, eine Antwort ohne das Feld liefert
      `wal`; je Sprache ein Unit-Test mit beiden Fällen und einer Fixture, deren
      Ursprung (reale Antwort von `GET /changes`) im Test genannt ist. *Zu
      belegen durch:* `make sdk-pack-csharp`, `make sdk-pack-kotlin` und
      `make sdk-pack-python` — der Bau führt die Tests aus und bricht bei rotem
      Test ab; jeder Aufruf trägt `--build-context proto=proto` (er ist für den
      **ganzen** Bau jedes SDK zwingend, nicht nur für gRPC:
      `BEO-PGC/zusatzkontext-kopplung-breiter-als-dod-wortlaut`, offen, 3×).
- [ ] Die Versionen stehen in den drei Metadaten-Quellen (`.csproj` `<Version>`,
      `build.gradle.kts` `version`, `pyproject.toml` `[project] version`) auf der
      nächsten Minor-Version nach dem Muster der bisherigen Hebungen (Wahl und
      Grund im Bericht); die Träger sind mitgezogen (§3.13-Suchlauf), die
      real erzeugten Artefaktnamen in den Beschreibungen sind aus **realen**
      Pack-Läufen geschrieben.
- [ ] Die Sprachreinheit der Kopien: die drei Sprach-Änderungen sind Kopien
      derselben Form; jede Kopie ist auf verbliebene deutsche Wortfragmente in
      Code, Kommentaren, Docstrings und README geprüft, nicht nur auf ihre
      Übereinstimmung mit dem Vorbild
      (`BEO-PGC/formvorbild-kopie-traegt-deutsches-wortfragment-weiter`,
      `BEO-PGC/deutsches-fachwort-im-englischen-sdk-readme`, je verkörpert, 3×).
- [ ] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein Self-Review (Modul 8).
- [ ] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff)
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [ ] Doku-Update: SDK-READMEs (englisch) nennen das Feld und die Version; Benutzerhandbuch (`**SDK:**`-Absätze der HTTP-Beschreibung) und Änderungshistorie, soweit sie die Version oder die Feldmenge tragen.
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

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaften: „die Version der drei Packages (0.2.0 → neue Version)", „die Feldmenge des HTTP-Lesemodells"; beide Stände gemessen):**

| Träger | Suchbefehl | Befund | Behandlung |
|---|---|---|---|
| Versionsangaben | `grep -rn '0\.2\.0' spec harness sdks docs/user` (Treffer nach Bezug auf die drei Packages lesen; andere Versionen sind keine Treffer) | *(Implementer trägt ein)* | Träger nachziehen; `docs/reviews/**` und `done/**` sind Records und bleiben |
| Artefaktnamen (`PgChangeFeed.Client.0.2.0.nupkg`, `pgchangefeed-kotlin-0.2.0.jar`, `pgchangefeed-0.2.0-…`) | `grep -rn 'nupkg\|\.jar\|\.whl' harness docs/user spec` | *(Implementer trägt ein)* | aus realen `make sdk-pack-*`-Läufen neu schreiben |
| Feldbeschreibungen des HTTP-Lesemodells in READMEs und Handbuch | `grep -rn 'old_image\|oldImage\|OldImage' sdks docs/user` (Fundstellen nach HTTP-Bezug lesen) | *(Implementer trägt ein)* | nachziehen |
| Beispiel-Clients, die `GET /changes` dekodieren | Lesen von `examples/**` | *(Implementer trägt ein)* | anpassen oder Punkt als entfallen im Bericht |
| Publish-Workflows (Tag-Abgleich gegen die Metadaten-Quelle) | Lesen der drei `sdk-*-release.yml` (nur lesen) | *(Implementer trägt ein)* | unverändert; der Tag-Abgleich liest die Quellen, die dieser Slice hebt |

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
