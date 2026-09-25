# Slice sdk-readme-nutzerdoku: Die drei SDK-READMEs und Paket-Metadaten als Anwender-Dokumentation — Aufbau, belegte Beispiele, API-Übersicht; Version 0.2.1

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle — es gibt keine Closure-Bedingung, die von der DoD dieses
Slice verschieden ist (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine
Welle braucht).

**Bezug:** [`LH-FA-SST-009`](../../../../spec/lastenheft.md) (offizielle Client-Bibliotheken — die drei
Packages, deren Paketbeschreibung dieser Slice ist),
[`ADR-0106`](../../adr/0106-csharp-nuget-erstes-sdk-package.md), [`ADR-0107`](../../adr/0107-python-pypi-zweites-sdk-package.md), [`ADR-0109`](../../adr/0109-kotlin-github-packages-drittes-sdk-package.md)
(die drei Package-Entscheidungen: Vertriebsweg, Metadaten-Quellen, Version unabhängig
vom Server), [`ADR-0110`](../../adr/0110-python-sdk-umfang-erweitert-vollmatrix.md) (Sprachmatrix, Flächen der Packages).

**Berührte Spec-Stellen:** [`SPEC-026`](../../../../spec/pflichtenheft.md), [`SPEC-027`](../../../../spec/pflichtenheft.md), [`SPEC-028`](../../../../spec/pflichtenheft.md)
(Package-Zeilen mit der aktuellen Version) und §1 `LH-FA-SST-009.a` (real erzeugte
Artefaktnamen) — geändert; [`SPEC-018`](../../../../spec/pflichtenheft.md) bis [`SPEC-024`](../../../../spec/pflichtenheft.md) (Drahtverträge, aus denen die
Anwender-Beschreibung ihre Aussagen liest) — gelesen.

**Verantwortlich:** —

**Autor:** Implementer-Agent (Nutzer-Anweisung „ohne Nachfragen weiter“, Rückmeldung zur
PyPI-Paketseite). **Datum:** 2026-09-25.

---

## 1. Ziel und Abgrenzung

**Ziel:** Die README jedes der drei SDK-Packages beschreibt in Anwender-Sprache, was
das Package tut, wie man es installiert, einsetzt und welche Klassen und Methoden es
bietet — belegt durch Beispiele aus dem Quelltext und den Tests der SDKs, ohne
interne Kennungen und ohne Projekt-Chronik; die Metadaten-Felder, die Anwender auf den
Paket-Seiten sehen, sind ebenso frei von Kennungen; die Quell-Version aller drei
Packages steht auf `0.2.1`, damit die Beschreibung mit dem nächsten Release erscheint.

**Ausgangslage (Nutzer-Rückmeldung, wörtlich):** (1) Der PyPI-Status-Text sei „nicht
aktuell/sinnvoll“ (er nennt, es gebe noch keinen Referenz-Client, und zählt Kennungen
auf); (2) „Warum steht dort nicht die API Beschreibung“ — die README enthält keine
API-Beschreibung, nur Verweise auf die Spezifikation; (3) mit den `SPEC-`-Kennungen
könne kein Benutzer etwas anfangen. Die README ist die Paketbeschreibung auf PyPI.org
(`pyproject.toml` `readme`) und auf NuGet.org (`PackageReadmeFile`); die Kotlin-README
gelangt nicht in das Artefakt (§3).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Ein Release / ein Tag** — ein Tag-Push (`sdk-*-v*`) ist Betreiber-Handlung
  ([`AGENTS.md`](../../../../AGENTS.md) §3.10); die drei `sdk-*-release.yml`-Workflows bleiben unverändert. Der
  Slice hebt nur die Quell-Version, damit ein späterer Tag eine noch nicht
  veröffentlichte Version trägt.
- **Änderungen an Signaturen oder Verhalten der SDKs** — die Beschreibung folgt dem
  Code, nicht umgekehrt. Auffälligkeiten im Code, die eine README-Aussage
  erschweren (z. B. transitive Abhängigkeiten der Kotlin-Flächen), werden gemeldet,
  nicht mitgeändert.
- **Das Benutzerhandbuch** (`docs/user/benutzerhandbuch.md`) — ein anderes Dokument
  mit eigener Konvention für Betreiber und Integratoren; es verweist auf die READMEs,
  ohne Aussagen zu ihrem Inhalt oder Paketstand zu machen (Suchlauf §3). Ändert sich
  daran nichts, entfällt auch Version und Historie.
- **Dauerhafter Wächter für die C#- und Kotlin-Beispiele** — die Beispiele werden
  einmalig durch Übersetzen belegt (§3). Für Python läuft der Wächter als Test im
  Package-Bau, weil die README dort bereits im Bau-Kontext liegt; für C# und Kotlin
  läge die README außerhalb der Test-Bau-Kontexte, ein Wächter dort wäre eine
  Änderung der Bau-Kontexte (ein anderer Vorgang).
- **Kennungen in Quelltext-Kommentaren der Build-Dateien** (`.csproj`,
  `pyproject.toml`, `build.gradle.kts`) — sie stehen in Kommentaren, die kein
  Paket-Metadatenfeld erreichen; `AGENTS.md` §3.7 lässt Rang-Zeiger und Herkunftsanker
  dort zu.

## 2. Definition of Done

- [ ] Die drei READMEs tragen den Anwender-Aufbau (Was ist das · Installation · Schnellstart ·
      Live-Streams · API-Übersicht · Change-Objekt mit `origin` · Fehlerbehandlung ·
      Zustellwege im Vergleich · Links · Lizenz), jedes Code-Beispiel trägt einen
      Beleg-Anker in §3 und ist übersetzt (C#, Kotlin) bzw. durch den Wächter im Bau
      geprüft (Python); keine `SPEC-`/`ADR-`/`LH-`/`ARC-`-Kennung, keine Chronik.
- [ ] Die Paket-Metadaten-Felder (Beschreibung, Schlagwörter, URLs) tragen keine Kennung; die
      Quell-Version der drei Packages ist `0.2.1`, alle Träger, die sie nennen, stehen
      auf demselben Stand (Suchlauf §3, beide Stände gemessen).
- [ ] `make sdk-pack-csharp`, `make sdk-pack-python`, `make sdk-pack-kotlin` grün, die
      Artefakte tragen die neue README bzw. Metadaten (Listing im Bericht),
      `make docs-check` und `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben oder „keine Beobachtung
      angefallen“ in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **mit**
      Wellen von der nächsten Welle-Closure geprüft (auch für Slices ohne
      Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Ist-Zustand gemessen am Parent (`6839f738`): Die drei READMEs tragen einen
„Status“-Absatz mit Kennungen (`SPEC-`/`ADR-`/`LH-`) und den Satz „no `examples/python/`
reference client exists yet“ (Python), verweisen für die API auf `spec/pflichtenheft.md`
und haben keine API-Übersicht, keine Beispiele (Ausnahme: der Kotlin-Installationsweg).
Wohin die README gelangt: Python `pyproject.toml` `readme = "README.md"` (der Docker-Bau
kopiert `sdks/python/README.md` an `pgchangefeed/README.md`); C# `PackageReadmeFile` (das
`.csproj` packt `../README.md` an die Paket-Wurzel); Kotlin: das `jar` trägt keine README,
der `publishing`-Block der `build.gradle.kts` trägt keine `pom`-Beschreibung — Anwender
sehen dort nur Koordinate und Version.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `sdks/python/README.md` | update | Anwender-Aufbau; Beispiele aus `sdks/python/pgchangefeed/integration/test_*_realserver.py` und `src/pgchangefeed/*.py` (Beleg-Anker unten je Beispiel). |
| `sdks/csharp/README.md` | update | dito, Beispiele aus `sdks/csharp/PgChangeFeed.Client.Integration/*RealserverTests.cs`. |
| `sdks/kotlin/pgchangefeed-kotlin/README.md` | update | dito, Beispiele aus `src/integrationTest/kotlin/…/*RealserverTest.kt`; der Installationsweg (GitHub Packages, Zugangsdaten) bleibt in Anwender-Sprache. |
| `sdks/python/pgchangefeed/pyproject.toml` | update | `keywords`, `[project.urls]` (Dokumentation, Fehlerberichte), `version = "0.2.1"`. |
| `sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.csproj` | update | `<PackageTags>`, `<Version>0.2.1</Version>`, Versions-Kommentar. |
| `sdks/kotlin/pgchangefeed-kotlin/build.gradle.kts` | update | `pom { name, description, url, licenses }` im `publishing`-Block (das Feld, das GitHub Packages anzeigt), `version` `0.2.1` an beiden Stellen (Kopf und `publishing`-Block), Versions-Kommentar. |
| `sdks/python/pgchangefeed/tests/test_readme_examples.py` | neu | Wächter: jeder ` ```python `-Block der README übersetzt (`compile`), jeder `from pgchangefeed… import`-Name löst auf, jede Methode der API-Tabelle existiert an der genannten Klasse. Greift, wenn README und Code auseinanderlaufen. |
| `spec/pflichtenheft.md` §1 (`LH-FA-SST-009.a`, drei Artefakt-Sätze), §6 `SPEC-026`/`-027`/`-028`, §7 | update | die Zeilen tragen die aktuelle Version bzw. den real erzeugten Artefaktnamen; eine Historie-Zeile. |
| `harness/README.md` §Sensors (`make sdk-pack-*`), `harness/mk/sdk.mk`, `sdks/csharp/Dockerfile`, `sdks/kotlin/Dockerfile` | update | Träger, die den Artefaktnamen mit der Version nennen (Suchlauf unten). |
| `docs/user/benutzerhandbuch.md` | prüfen | Absätze `**SDK:**` verweisen auf die READMEs; nur ändern, wenn sie Aussagen zu README/Paketstand machen, die nach der Änderung nicht mehr stimmen. |

**Beleg-Anker der Code-Beispiele** (je Beispiel die Quelle im Repo; werden beim Schreiben
eingetragen und bei der Plan-Nachführung um die Übersetzungs-Belege ergänzt):

*(Abschnitt wird mit dem Bau gefüllt.)*

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaften: „die Version der drei
Packages (0.2.0 → 0.2.1)“, „Kennungen in Paket-Beschreibung und Metadaten-Feldern“,
„der Text der README als Paketbeschreibung“; beide Stände gemessen):**

*(Abschnitt wird mit dem Bau gefüllt.)*

## 4. Trigger

**Start** (`next` → `in-progress`): kein anderer Slice liegt in `in-progress/` (WIP-Limit 1);
Nutzer-Anweisung liegt vor.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): wenn die Beispiel-Belege
  Änderungen an den SDK-Signaturen verlangten (dann Code-Slice je Package).
- `in-progress` → `open` (blockiert — Carveout?): wenn ein Pack-Lauf am Netz oder Speicher
  des Hosts scheitert und nicht wiederholbar ist.

## 5. Closure-Trigger

DoD vollständig, `make gates` grün, Review-Report liegt vor und ist aufgelöst, Verifikation
bestätigt die DoD, Closure-Notiz mit Steering-Loop-Eintrag geschrieben.

## 6. Risiken und offene Punkte

- Eine README-Aussage stimmt nicht mit dem Verhalten des Servers überein, weil sie aus
  dem Quelltext des SDK statt aus dem Server gelesen wurde — **Ausgang:** offen bis zur
  Verifikation; Gegenmaßnahme: Aussagen zu Positionen, Tokens und `limit` sind gegen
  `spec/pflichtenheft.md` (`SPEC-018`, `SPEC-022`) und das Benutzerhandbuch gelesen.
- Die Beschreibung erscheint erst nach einem Release (Tag-Push, Betreiber-Handlung) auf den
  Paket-Seiten; ob sie dort wie erwartet gerendert wird, ist lokal nicht prüfbar
  ([`AGENTS.md`](../../../../AGENTS.md) §3.10 sinngemäß: externe Oberfläche) — **Ausgang:** weiter offen bis zum
  ersten Release mit `0.2.1`.
- Die Kotlin-Beschreibung liegt im `pom`; ob GitHub Packages sie anzeigt, ist lokal nicht
  prüfbar — **Ausgang:** weiter offen bis zum ersten Release mit `0.2.1`.

## 7. Closure-Notiz

*(wird bei der Closure durch den Planner gefüllt)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist `sdks/*` (die SDK-Sprachpakete) und
`spec/`/`harness/` als Träger von Versionsangaben; die Modus-Deklaration führt nur die
Default-Sub-Area `*` (`PGC`, Greenfield) — keine Ausdifferenzierung nötig, keine Schwelle
verfehlt.

**Vorgelagert — offene Beobachtungen sichten:** Register `BEO-PGC` durchgegangen; Treffer je
berührter Sub-Area: `deutsches-fachwort-im-englischen-sdk-readme` (Zähler 3×, verkörpert:
Sprachreinheit der englischen READMEs — die drei neuen READMEs sind je Form-Teil gesichtet,
Suchlauf auf deutsche Wortfragmente im Bericht), `formvorbild-kopie-traegt-deutsches-wortfragment-weiter`
(3×, verkörpert), `zahl-in-traeger-driftet-gegen-die-messung` (21×, verkörpert: die Trefferzahlen
im Suchlauf stehen wie gedruckt), `nachzug-laesst-ueberholten-text-stehen` (8×, verkörpert: die
Versionsträger werden am Ende gegen die Beschreibung gesucht), `arbeit-ueberholt-stehenden-traeger`
(31×, verkörpert: Suchlauf §3 als committetes Feld), `beleg-befehl-traegt-seinen-satz-nicht`
(13×, verkörpert: jeder Beleg nennt den Lauf und die gedruckte Zeile),
`drei-sprachen-kopie-divergiert-am-randfall` (1×, offen: die Aussagen zu `origin` und Fehlerklassen
je Sprache werden gegen den Quelltext der jeweiligen Sprache gelesen, nicht aus einer Schwester-README
kopiert).

**Modus:** alle berührten Sub-Areas GF.
