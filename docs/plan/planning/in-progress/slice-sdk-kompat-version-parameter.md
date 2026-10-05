# Slice sdk-kompat-version-parameter: `make test-sdk-kompat` liest die Version der Bibliothek aus den Version-Dateien der SDKs statt sie fest zu tragen

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle — es gibt keine Closure-Bedingung jenseits der DoD
dieses Slice, siehe Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht (Modul 6).

**Bezug:** [`ADR-0145`](../../adr/0145-sdk-meldungscodes-eigenschaft-der-fehlertypen.md)
(Festlegung 4: die 0.5.x-Signaturen bleiben erhalten — die Aussage, die das Ziel
misst). Keine Anforderung des Lastenhefts ist berührt: der Slice ändert ein
Harness-Werkzeug, kein Produkt.

**Berührte Spec-Stellen:** — (keine).

**Verantwortlich:** pt9912 (Implementer-Agent im Auftrag).

**Autor:** pt9912 (Planner, Meldung mit Frist aus der Closure von
`slice-sdk-tls-optionen`). **Datum:** 2026-10-05.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung**.

**Anlass (gemessen).** `make test-sdk-kompat` im Standardmodus
(`SDK_KOMPAT_NEU=dist`) endet rot an der Auflösung des C#-Pakets: das
Werkzeug erwartet die Version `0.6.0` in `sdks/csharp/dist`, dort liegt `0.6.1`
(`NU1603 … PgChangeFeed.Client 0.6.0 was not found. … 0.6.1 was resolved instead`,
gelaufen im Fixrunden-Lauf von `slice-sdk-tls-optionen`, Exit 2). Die
Package-Versionen stehen seit dem Stand `0.6.1` in `PgChangeFeed.Client.csproj`
(`<Version>`), `build.gradle.kts` (`version`) und `pyproject.toml` (`version`);
das Werkzeug trägt `0.6.0` an 39 Stellen (Suchlauf unten, gemessen am Parent).
Für Kotlin und Python ist der Rotbefund **nicht gemessen** (*hergeleitet*:
dort ersetzt der Lauf eine Datei nach Namen bzw. nutzt ein Wheel; ob der
Dateiname `…-0.6.0.jar` bei dist `0.6.1` bricht, misst der Implementer).

**Ziel:** Im Modus `dist` liest das Werkzeug die Version der zu messenden
Bibliothek aus den Version-Dateien der drei SDKs (eine Quelle je Sprache) und
trägt sie nirgends als Literal; im Modus `registry` bleibt die feste, veröffentlichte
Version eine benannte Konstante an genau einer Stelle. `make test-sdk-kompat`
endet im Standardmodus mit Exit 0 am Arbeitsstand.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Die Grundlinie 0.5.0.** Sie ist die Messbasis („altes Paket“), kein Stand,
  der mitläuft; ihr Literal bleibt, als eine Konstante je Sprache.
- **Neue Kompatibilitäts-Messungen** (etwa das C#-Mapping von `Internal` auf
  `UnexpectedStatus` aus `slice-sdk-tls-optionen`) — anderer Gegenstand: das Ziel
  misst die Fehlertypen aus `ADR-0145`; eine Erweiterung ist ein eigener Slice
  mit eigener Fitness-Function.
- **Ein Aufnehmen des Ziels in `make gates`.** Das Ziel braucht Netz (NuGet, PyPI,
  Cloudsmith, Maven Central) und bleibt Werkzeug, siehe
  `harness/targets/sdk-kompat.md`.
- **Anheben der Package-Versionen.** Gehört in den Release-Zug mit neuer Freigabe
  des Auftraggebers.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; Gate-Läufe und Closure-Pflichten
zählen nicht mit.

- [x] **Version aus der Quelle:** kein `0.6.0`-Literal im Modus `dist` mehr in
      `tools/harness/sdk-kompat` (Suchlauf: die Trefferzahl der Stellen ohne die
      benannte Konstante des Modus `registry` ist 0); die Version je Sprache
      kommt aus `.csproj`, `build.gradle.kts` bzw. `pyproject.toml`.
      Beleg: Suchlauf in §3 (`tools/harness/sdk-kompat` 0 Treffer, die Konstante
      `REGISTRY_VERSION=0.6.0` als einziger Treffer in
      `tools/harness/run-sdk-kompat-tests.sh`); Lauf in §3 *Messung*.
- [x] **Messung:** `make test-sdk-kompat` im Standardmodus Exit 0 am Arbeitsstand
      (gedruckte Schlusszeile, Exit direkt gelesen), `SDK_KOMPAT_NEU=registry`
      unverändert Exit 0; eine Mutation der Version (Eingabeseite: ein anderer
      Wert in einer Version-Datei einer Kopie im Scratchpad) färbt den Lauf rot
      mit lesbarer Meldung statt eines stillen Rückfalls auf eine andere Version.
      Beleg: §3 *Messung* (Läufe, Schlusszeilen, Exit) und *Mutationen*.
- [x] **Vertrag nachgezogen:** `harness/targets/sdk-kompat.md` und die Kommentare
      in `harness/mk/sdk.mk`, `tools/harness/run-sdk-kompat-tests.sh` nennen keine
      feste Version mehr als Gegenstand des Standardmodus (Lese-Handlung des
      Reviewers; der Suchlauf in §3 deckt nur `tools/harness/sdk-kompat`).
      Dazu die Index-Zeile in `harness/README.md` (§3, Nachzug).
- [x] `make gates` grün (Lauf nach dem letzten Commit dieses Laufs, Exit 0 direkt
      gelesen).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`), kein Self-Review.
- [x] Doku-Update: nur der Target-Vertrag (kein öffentlicher Vertrag berührt);
      dazu die Index-Zeile in `harness/README.md`, `docs/user/` ist nicht berührt.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register fortgeschrieben oder „keine Beobachtung“ in §7.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

(Liefer-Punkte: drei — Version aus der Quelle, Messung, Vertrag.)

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/harness/run-sdk-kompat-tests.sh` | update | liest die Version je Sprache aus der Version-Datei (Modus `dist`) und reicht sie als Build-Argument bzw. Umgebung an die Sprach-Läufer; die feste Version des Modus `registry` als eine benannte Konstante |
| `tools/harness/sdk-kompat/{csharp,kotlin,python}/Dockerfile`, `run.sh` | update | `-p:PgcfVersion=0.6.0`, Jar-Name `pgchangefeed-kotlin-0.6.0.jar`, Wheel-Auswahl: Literal wird Parameter; die 39 Stellen sind gezählt, nicht alle sind Version-Literale (Kommentare, Anzeigetexte) — der Implementer trennt sie |
| `harness/targets/sdk-kompat.md`, `harness/mk/sdk.mk` | update | Vertrag und Kommentare: „Bibliothek 0.6.0“ wird „Bibliothek des Arbeitsstands (dist) bzw. 0.6.0 (registry)“ |
| `harness/README.md` (Index-Zeile `make test-sdk-kompat`) | update — **Nachzug** | die Zeile nannte „SDK-Packages 0.6.0“ als Gegenstand; der Vertrag führt sie als Kurzfassung, ein Satz |
| `tools/harness/sdk-kompat/python/run.sh` (A2) | update — **Nachzug** | über den Plan hinaus: der Austausch auf die neue Bibliothek ist geprüft (pip-Exit und installierte Version); vor dem Fix lief A2 nach gescheitertem `pip install` still gegen 0.5.0 grün (Rotbefund unten) |
| `tools/harness/sdk-kompat/csharp/run.sh` (A2) | update — **Nachzug** | trägt `/neu` das Artefakt der neuen Version und fehlt die daraus gewonnene DLL, ist A2 rot statt auf die veröffentlichte DLL zu fallen |
| `harness/targets/sdk-kompat.md` §Versionen, §Grenze 6/7, §Test | update — **Nachzug** | neuer Abschnitt zur Quelle der Versionen; Grenzen „Dateiname, nicht Inhalt“ und „`REGISTRY_VERSION` zieht niemand nach“; drei Mutations-Zeilen; der Satz „A3 und A5 lesen immer das veröffentlichte Paket“ gilt gemessen nur für C# und Kotlin (Python: `A3-Grundlage` und zweite Signatur-Lesung laufen unter der Bibliothek aus A2) und ist entsprechend berichtigt |
| `harness/mk/sdk.mk` (Ziel `test-sdk-altserver`) | **nicht geändert — gemeldet** | Kommentar und Hilfetext nennen „0.6.0-SDKs“ bzw. „SDK 0.6.0“, das Ziel fährt aber die SDK-Quelle des Arbeitsstands (0.6.1); anderes Ziel, anderer Vertrag (`harness/targets/sdk-altserver.md`) — Frist: Closure dieses Slice, der Planner zieht nach oder benennt den Träger mit Adresse |

```suchlauf
3e028ee3 39 -n 0\.6\.0 -- tools/harness/sdk-kompat
diff 0 -n 0\.6\.0 -- tools/harness/sdk-kompat
diff 1 -n 0\.6\.0 -- tools/harness/run-sdk-kompat-tests.sh
3b173e28 5 -n 0\.6\.0 -- tools/harness/run-sdk-kompat-tests.sh
3b173e28 17 -n -e Bibliothek.0\.6\.0 -e Packages.0\.6\.0 -e 0\.6\.0-Pakete -- harness tools/harness
diff 0 -n -e Bibliothek.0\.6\.0 -e Packages.0\.6\.0 -e 0\.6\.0-Pakete -- harness tools/harness
```

Stand der ersten Zeile: **gemessen** mit `git grep -n` am Parent `3e028ee3` (neun Dateien,
je Datei 8, 1, 12, 4, 8, 1, 1, 1, 3 Treffer); der Implementer wiederholt sie mit
`make suchlauf-nachmessen PLAN=<Plan-Datei>`. Der Suchlauf trifft das Literal,
nicht jede Stelle, die eine Version implizit trägt (Dateinamen mit Versionsteil,
Umgebungsvariablen): das liest der Reviewer.

**Suchlauf des Implementers (gemessen, `make suchlauf-nachmessen`).** Zeile 2:
`tools/harness/sdk-kompat` trägt kein `0.6.0` mehr (Literal und Kommentare der
Gast-Quellen sind Parameter bzw. „die Bibliothek mit Meldungscode“). Zeile 3/4:
im Runner bleibt genau ein Treffer, die Konstante `REGISTRY_VERSION=0.6.0`
(Parent: fünf, alle im Kopfkommentar). Zeile 5/6: die
Beschreibung der bewegten Eigenschaft („Bibliothek 0.6.0“, „Packages 0.6.0“,
„0.6.0-Pakete“) in `harness/` und `tools/harness/` — am Parent `3b173e28` 17
Treffer in sechs Dateien, am Diff 0. **Gefunden, nicht geändert:** die
Altserver-Zeilen in `harness/mk/sdk.mk` (Tabelle oben, gemeldet); die
`0.6.0`-Treffer im Ganzen-Baum-Suchlauf außerhalb von `harness/` und `tools/`
(Lastenheft-/ADR-/Handbuch-Text über das Release 0.6.0) beschreiben die
veröffentlichte Version, nicht den Gegenstand dieses Ziels. **Nicht gefunden:**
ein weiterer Träger, der `make test-sdk-kompat` mit einer festen Version nennt
(`git grep -n sdk-kompat` am Diff, ohne `docs/reviews`, `done/`,
`.harness/baseline`: [`ADR-0147`](../../adr/0147-sdk-csharp-http-basis-null-literal-einschraenkung.md) nennt `0.6.0 CS0121` als Messwert eines Laufs am
Stand 0.6.0 — ein Record, kein Gegenstand).

**Messung (gemessen, Läufe dieses Implementer-Laufs, Arbeitsstand vor dem Plan-Commit,
Code-Stand `500bebb2`).** Netz: NuGet, PyPI, Cloudsmith (`…/pgchangefeed-kotlin/0.6.0/…pom`)
und Maven Central antworteten je HTTP 200 (Abruf per `curl` im Container).

- **Vor dem Fix** (Parent `3b173e28`, dist `0.6.1` je Sprache): C# Exit 2 am Bau
  (`NU1603 … PgChangeFeed.Client 0.6.0 was not found. … 0.6.1 was resolved instead`,
  danach `cp /kompat/pk-dist/pgchangefeed.client/0.6.0/…` gescheitert). Kotlin
  Exit 0, aber die Zeile lügt: `KOMPAT kotlin A2: Bibliothek 0.6.0 (Artefakt
  pgchangefeed-kotlin-0.6.1.jar) …` (das Jar wurde unter dem Namen 0.6.0 kopiert).
  Python Exit 2: `ERROR: Could not find a version that satisfies the requirement
  pgchangefeed==0.6.0 (from versions: 0.6.1)`, danach **still grün**
  `KOMPAT python A2: Bibliothek pgchangefeed 0.5.0 (Artefakt
  pgchangefeed-0.6.1-py3-none-any.whl) ersetzt 0.5.0` / `A2: 45 Aufrufe ok`
  (A2 maß 0.5.0 gegen 0.5.0), rot erst an `A3-Grundlage`.
- **Nach dem Fix, `make test-sdk-kompat`** (dist): Exit 0, Schlusszeile
  `run-sdk-kompat-tests: Kompatibilitätsmessung (dist) grün für: csharp kotlin python`;
  A2-Zeilen `KOMPAT csharp A2: Bibliothek 0.6.1 (Artefakt PgChangeFeed.Client.0.6.1.nupkg)
  e5fa52ed7181 …` / `53 Aufrufe ok`, `KOMPAT kotlin A2: Bibliothek 0.6.1 (Artefakt
  pgchangefeed-kotlin-0.6.1.jar) af19721beee7 …` / `45 Aufrufe ok`,
  `KOMPAT python A2: Bibliothek pgchangefeed 0.6.1 (Artefakt
  pgchangefeed-0.6.1-py3-none-any.whl) …` / `45 Aufrufe ok`; keine Zeile `ROT`.
- **`SDK_KOMPAT_NEU=registry make test-sdk-kompat`:** Exit 0, Schlusszeile
  `run-sdk-kompat-tests: Kompatibilitätsmessung (registry) grün für: csharp kotlin python`;
  A2 je Sprache `Bibliothek 0.6.0 (Registry …)` (C# DLL `f8e6415b23e5`, verschieden
  vom dist-Artefakt `e5fa52ed7181`); keine Zeile `ROT`.

**Mutationen (Eingabeseite, je einmal gefahren).** Zusage · mutierte Eingabe ·
gesehenes Rot; die Läufe an einem Klon des Repos im Scratchpad (Stand
`500bebb2`), Artefakte aus `sdks/<sprache>/dist` über `SDK_KOMPAT_DIST_<SPRACHE>`.

- Version aus der Version-Datei · `<Version>0.6.9</Version>` in der Klon-`.csproj` ·
  Exit 2 vor jedem Bau: `… trägt kein Artefakt PgChangeFeed.Client.0.6.9.nupkg
  (Version 0.6.9 aus sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.csproj);
  vorhanden: PgChangeFeed.Client.0.6.1.nupkg — make sdk-pack-csharp vorher …`.
- Keine Version ohne Versionszeile · Zeile `version = "0.6.1"` aus der
  Klon-`pyproject.toml` entfernt · Exit 2: `sdks/python/pgchangefeed/pyproject.toml
  trägt nicht genau eine Version der Form X.Y.Z (gelesen: keine)`.
- Eindeutige Version · zweite Zeile `version = "0.6.2"` in der Klon-`build.gradle.kts` ·
  Exit 2: `… build.gradle.kts trägt nicht genau eine Version der Form X.Y.Z
  (gelesen: 0.6.1 0.6.2)`. Gegenprobe: derselbe Klon ohne Mutation, Kotlin, Exit 0
  (`Kompatibilitätsmessung (dist) grün für: kotlin`).
- Version-Datei muss existieren · Klon-`.csproj` gelöscht · Exit 2:
  `Version-Datei sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.csproj der
  Sprache csharp fehlt`.
- Artefakt genau dieser Version · Verzeichnis mit
  `pgchangefeed-0.6.2-py3-none-any.whl` (`SDK_KOMPAT_DIST_PYTHON`) · Exit 2:
  `… trägt kein Artefakt pgchangefeed-0.6.1-py3-none-any.whl (Version 0.6.1 aus
  sdks/python/pgchangefeed/pyproject.toml); vorhanden: pgchangefeed-0.6.2-py3-none-any.whl …`.
- Python-Austausch ohne stillen Rückfall · das veröffentlichte Rad 0.6.0 unter dem
  Namen `pgchangefeed-0.6.1-py3-none-any.whl` (`SDK_KOMPAT_DIST_PYTHON`) · Exit 2
  über `make`: `KOMPAT python A2: ROT — Austausch auf 0.6.1 gescheitert (pip Exit 1,
  installiert 0.5.0, Artefakt pgchangefeed-0.6.1-py3-none-any.whl)`.
- **Nicht gefahren:** die bestehende Mutationsprobe A4 (Signatur aus einer
  SDK-Kopie entfernt). Der Pfad, über den das dist-Artefakt in A2 ankommt, ist durch
  die verschiedenen Prüfsummen belegt (dist `e5fa52ed7181`, Registry `f8e6415b23e5`);
  dass eine entfernte Signatur mit der neuen Versionsführung weiter rot wird, ist
  *hergeleitet*.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): der Slice steht vor dem Release-Zug der
SDK-Versionen oder vor der nächsten Änderung an einem Package, die `make
test-sdk-kompat` braucht; ein WIP-Limit-1-Platz in `in-progress/` ist frei; Netz
für NuGet, PyPI, Cloudsmith und Maven Central.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): nicht erwartet (drei Sprach-Läufer, eine Quelle).
- `in-progress` → `open` (blockiert): ein Registry-Pfad ist nicht erreichbar —
  dann Messung des Modus `dist` allein, der Modus `registry` als Folge.

## 5. Closure-Trigger

DoD vollständig, `make gates` grün, Review-Bericht liegt vor, Closure-Notiz mit
Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — jedes Risiko bekommt genau einen
Ausgang.

- Das Werkzeug misst gegen ein Paket, dessen Version sich zwischen Versionssprüngen
  ändert, ohne dass ein Sensor es meldet (kein Gate) — **Ausgang:** offen bis zur
  Closure; die Quelle der Version ist danach die Version-Datei, ein Sprung ist kein
  Bruch mehr.
- Kotlin und Python sind für den Sprung nicht gemessen (*hergeleitet*, siehe §1)
  — **Ausgang:** offen bis zur Closure; der Implementer misst vor dem Fix den
  Rotbefund je Sprache und nennt ihn im Bericht.

## 7. Closure-Notiz

- **Was hat funktioniert:** —
- **Was ging anders als geplant:** —
- **Steering-Loop-Eintrag:** —
- **Beobachtungs-Register:** —
- **Folge-Slices:** —
- **Risiken aus §6:** —
- **Drei Paarungen:** —

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist `*` (Default; `tools/harness`,
`harness/`); eine Ausdifferenzierung ist nicht nötig.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen; kein Treffer
zu `sdk-kompat`.

Alle berührten Sub-Areas GF.
