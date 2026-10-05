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
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`), kein Self-Review.
      Beleg: `docs/reviews/review-slice-sdk-kompat-version-parameter.md`
      (1 HIGH, gelöst) und `docs/reviews/review-slice-sdk-kompat-version-parameter-fixrunde.md`
      (kein HIGH); Verifikation `docs/reviews/verify-slice-sdk-kompat-version-parameter.md`
      (DoD bestätigt, V-3).
- [x] Doku-Update: nur der Target-Vertrag (kein öffentlicher Vertrag berührt);
      dazu die Index-Zeile in `harness/README.md`, `docs/user/` ist nicht berührt.
- [x] Gemeldete Träger fremder Dateien mit der Closure nachgezogen (§3 Tabelle,
      Ziel `test-sdk-altserver`; `AGENTS.md` §3.13).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register fortgeschrieben oder „keine Beobachtung“ in §7.
- [x] Jedes Risiko aus §6 trägt einen Ausgang.
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

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
| `harness/mk/sdk.mk` (Ziel `test-sdk-altserver`), `harness/targets/sdk-altserver.md` (Titel), `tools/harness/run-sdk-altserver-tests.sh` (Kopfkommentar) | **gemeldet, in der Closure nachgezogen** (versionsneutral: „SDK des Arbeitsstands“ bzw. „Quelle des Arbeitsstands“, §7) | Kommentar und Hilfetext in `sdk.mk`, der Titel des Vertrags und der Kopf des Runners nennen „0.6.0-SDKs“ bzw. „SDK 0.6.0“, das Ziel fährt aber die SDK-Quelle des Arbeitsstands (0.6.1); anderes Ziel, anderer Vertrag — Frist: Closure dieses Slice, der Planner zieht nach oder benennt den Träger mit Adresse (die zwei letzten ergänzt in der Fixrunde, Review F-2) |
| `tools/harness/run-sdk-kompat-tests.sh` (`version_lesen`) | update — **Fixrunde** | strikt `X.Y.Z` ohne Suffix, Meldung „…, nur Ziffern ohne Suffix“ (Review F-1) |
| `tools/harness/sdk-kompat/{csharp,kotlin}/{Dockerfile,run.sh}`, `python/{Dockerfile,run.sh}` | update — **Fixrunde** | Bezeichner nach Rolle statt Version: `alt06`→`altreg`, `neu06`→`neureg`, `LIB06`→`LIBREG`, `d06`→`dreg`, `/wheels/06`→`/wheels/neu`; Kennzeichnung des über den Runner nicht erreichbaren C#-Zweigs (Review F-5, F-6) |
| `harness/targets/sdk-kompat.md` (§Versionen, §Ausgänge, §Grenze 8/9, §Test) | update — **Fixrunde** | Form `X.Y.Z` ohne Suffix samt PEP-440-Hinweis; Grenzen „A2 tauscht nicht die Abhängigkeiten“ und „Kotlin: nur die oberste `version`-Zeile“; Mutations-Zeile zum Suffix (Review F-1, F-3, F-4) |

```suchlauf
3e028ee3 39 -n 0\.6\.0 -- tools/harness/sdk-kompat
diff 0 -n 0\.6\.0 -- tools/harness/sdk-kompat
diff 1 -n 0\.6\.0 -- tools/harness/run-sdk-kompat-tests.sh
3b173e28 5 -n 0\.6\.0 -- tools/harness/run-sdk-kompat-tests.sh
3b173e28 17 -n -e Bibliothek.0\.6\.0 -e Packages.0\.6\.0 -e 0\.6\.0-Pakete -- harness tools/harness
diff 0 -n -e Bibliothek.0\.6\.0 -e Packages.0\.6\.0 -e 0\.6\.0-Pakete -- harness tools/harness
75d2dee8 4 -n -e 0\.6\.0-SDKs -e SDK.0\.6\.0 -- harness tools/harness
diff 0 -n -e 0\.6\.0-SDKs -e SDK.0\.6\.0 -- harness tools/harness
9d1edf2a 26 -n -e alt06 -e neu06 -e LIB06 -e wheels/06 -e d06 -- tools/harness/sdk-kompat
diff 0 -n -e alt06 -e neu06 -e LIB06 -e wheels/06 -e d06 -- tools/harness/sdk-kompat
```

**Closure-Nachzug des Suchlaufs (Planner, gemessen mit `git grep -n`).** Zeile 7
(in der Fixrunde `diff 4`) steht auf dem Stand vor dem Träger-Nachzug (`75d2dee8`,
4 Treffer in drei Dateien), Zeile 8 misst am Arbeitsbaum der Closure 0 Treffer: der Nachzug in
`harness/mk/sdk.mk`, `harness/targets/sdk-altserver.md` und
`tools/harness/run-sdk-altserver-tests.sh` (§7). Zeile 9/10 tragen den fünften
umbenannten Namen `d06` (Re-Review F-4): am Parent der Fixrunde `9d1edf2a`
26 Treffer, davon 3 für `d06` allein (`git grep -n -w d06 9d1edf2a -- tools/harness/sdk-kompat`),
am Arbeitsbaum 0.

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

### Fixrunde (Review `review-slice-sdk-kompat-version-parameter`, Code-Stand `a8e9ddf0`)

- **F-1 (HIGH) — Zusage „Form X.Y.Z“ nicht getragen.** Der Code wurde bewegt, nicht
  die Zusage: `version_lesen` verlangt `^[0-9]+\.[0-9]+\.[0-9]+$`. Alle drei Packages
  tragen heute diese Form, und für eine Version mit Suffix stimmte der
  Python-Dateiname nicht mit dem Artefakt überein, das der Pack schreibt (PEP-440-Normalform).
  Die Meldung lautet jetzt „… trägt nicht genau eine Version der Form X.Y.Z, nur
  Ziffern ohne Suffix (gelesen: …)“, Kommentar und Vertrag (§Versionen, §Ausgänge)
  sagen dasselbe. **Mutation (gemessen)** an einem Klon (Stand `a8e9ddf0`) im
  Scratchpad: `<Version>0.7.0-rc.1</Version>` in der `.csproj` mit einem Verzeichnis
  `PgChangeFeed.Client.0.7.0-rc.1.nupkg` über `SDK_KOMPAT_DIST_CSHARP` → Exit 2,
  `sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.csproj trägt nicht genau eine
  Version der Form X.Y.Z, nur Ziffern ohne Suffix (gelesen: 0.7.0-rc.1)`, **rot**;
  `version = "0.7.0-rc.1"` in `pyproject.toml` mit `pgchangefeed-0.7.0-rc.1-py3-none-any.whl`
  über `SDK_KOMPAT_DIST_PYTHON` → Exit 2, dieselbe Meldung für `pyproject.toml`,
  **rot**. **Gegenprobe** mit `0.6.1` (Arbeitsstand, `make test-sdk-kompat`, dist):
  Exit 0, **grün** (unten).
- **F-2 (LOW) — Träger-Meldung unvollständig.** Ergänzt in der §3-Tabelle:
  `harness/targets/sdk-altserver.md` (Titel) und `tools/harness/run-sdk-altserver-tests.sh`
  (Kopf), nur gemeldet, Frist Closure; Suchlauf-Zeile 7 (4 Treffer in drei Dateien).
- **F-3 (INFO) — Laufzeit-Abhängigkeiten in A2.** Als §Grenze 8 im Vertrag.
- **F-4 (INFO) — Kotlin-Versionszeile.** Als §Grenze 9 im Vertrag; §Versionen
  sagt „die nicht eingerückte Zeile“.
- **F-5 (INFO) — nicht erreichbarer C#-Zweig.** Behalten und im Kommentar als über
  den Runner nicht erreichbar gekennzeichnet, mit der Kopplung an die RUN-Schicht
  `lib-dist` im Dockerfile. Grund: ohne den Zweig liefe A2 bei einer künftigen
  Änderung dieser Schicht mit der 0.5.0-DLL aus `alt05` weiter — derselbe stille
  Rückfall, den der Slice in Python beseitigt. Keine Mutation (über den Runner keine
  erreichbare Eingabe); dass der Zweig rot färbt, ist *hergeleitet*.
- **F-6 (INFO) — Version im Bezeichner.** Umbenannt: `alt06`→`altreg`,
  `neu06`→`neureg`, `LIB06`→`LIBREG`, `d06`→`dreg`, `/wheels/06`→`/wheels/neu`;
  die gedruckte Zeile `KOMPAT kotlin javap Bibliothek 06` heißt `… Bibliothek REG`.
  Suchlauf-Zeile 8 (nach dem Closure-Nachzug Zeile 10): 0 Treffer (das Muster
  trug `d06` nicht; nachgezogen in der Closure, Re-Review F-4).

**Läufe am Endstand der Fixrunde (Code `a8e9ddf0`, gemessen).**
`make test-sdk-kompat` (dist): Exit 0, Schlusszeile
`run-sdk-kompat-tests: Kompatibilitätsmessung (dist) grün für: csharp kotlin python`,
keine Zeile `ROT`, A2 je Sprache `Bibliothek 0.6.1 (Artefakt …)` mit denselben
Prüfsummen wie oben. `SDK_KOMPAT_NEU=registry make test-sdk-kompat`: Exit 0,
Schlusszeile `… (registry) grün für: csharp kotlin python`, keine Zeile `ROT`
(gefahren, weil die Umbenennung die Pfade beider Modi berührt).
`make kommentar-kennungen DIFF=dad07666`: Exit 0, kein Kandidat; der
Konjunktiv-Kandidatenlauf (Schritt 20) über `*.sh` und die Gast-Dockerfiles seit
`dad07666`: kein Treffer.

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
  ändert, ohne dass ein Sensor es meldet (kein Gate) — **Ausgang:** entfallen —
  im Modus `dist` kommt die Version aus der Version-Datei je Sprache, kein Literal
  bleibt (Suchlauf Zeile 2, 0 Treffer); ein Sprung ist kein Bruch mehr: am Stand
  `0.6.1` Exit 0 in allen drei Sprachen (§3 *Messung*, vom Verifier reproduziert),
  eine Version-Datei ohne passendes Artefakt endet laut mit Exit 2 (§3
  *Mutationen*). Was bleibt, ist `REGISTRY_VERSION`, die niemand nachzieht — benannt
  in §Grenze 7 des Vertrags, kein Risiko dieses Slice.
- Kotlin und Python sind für den Sprung nicht gemessen (*hergeleitet*, siehe §1)
  — **Ausgang:** entfallen — vor dem Fix gemessen (§3 *Vor dem Fix*): Kotlin Exit 0
  mit falscher Zeile (Jar 0.6.1 unter dem Namen 0.6.0), Python A2 still grün gegen
  0.5.0 und rot erst an `A3-Grundlage`; beide Befunde trägt der Fix (A2-Zeilen mit
  `Bibliothek 0.6.1`, Python-Mutation „Austausch ohne stillen Rückfall“ rot).

## 7. Closure-Notiz

- **Geliefert:** `make test-sdk-kompat` liest im Modus `dist` die Version je Sprache
  aus `.csproj`, `build.gradle.kts` (oberste Zeile) bzw. `pyproject.toml`, verlangt
  die Form `X.Y.Z` ohne Suffix und das Artefakt genau dieser Version, ohne
  Ersatzwert; `REGISTRY_VERSION=0.6.0` ist die einzige feste Version (Modus
  `registry`, A3/A5 in C# und Kotlin). Python-A2 prüft den Austausch (pip-Exit und
  installierte Version), C#-A2 fällt nicht still auf die veröffentlichte DLL zurück;
  Gast-Bezeichner nach Rolle statt Version. Vertrag `harness/targets/sdk-kompat.md`
  (§Versionen, §Grenze 6–10, §Test), Index-Zeile in `harness/README.md`. Am Stand
  `0.6.1`: `dist` und `registry` je Exit 0 in allen drei Sprachen. Review
  `docs/reviews/review-slice-sdk-kompat-version-parameter.md` (1 HIGH, 1 LOW, 4 INFO;
  Fixrunde `a8e9ddf0`, `6abac454`), Re-Review
  `docs/reviews/review-slice-sdk-kompat-version-parameter-fixrunde.md` (kein HIGH;
  F-4 LOW, F-1–F-3 INFO), Verifikation
  `docs/reviews/verify-slice-sdk-kompat-version-parameter.md` (DoD bestätigt, V-1 LOW,
  V-2/V-3 INFO, `make gates` Exit 0).
- **Was hat funktioniert:** Der Rotbefund vor dem Fix wurde in allen drei Sprachen
  gemessen statt hergeleitet; er fand zwei Fehler, die das Rot in C# verdeckt hätte:
  Kotlin grün mit falscher A2-Zeile, Python-A2 still grün gegen 0.5.0. Beide sind
  jetzt Mutationen mit sichtbarem Rot. Die getrennten Prüfsummen von dist- und
  Registry-DLL (`e5fa52ed7181` gegen `f8e6415b23e5`) belegen den Weg des Artefakts in
  A2 ohne A4-Lauf. Reviewer und Verifier haben je eigene Proben mit einem
  `docker`-Stub gefahren und so Aussagen geprüft, die der Implementer nur an einer
  Sprache gemessen hatte.
- **Was ging anders als geplant:** (1) **HIGH, Zusage weiter als ihr Code:**
  Kommentar, Meldung und Vertrag sagten „Form X.Y.Z“, der Code ließ Suffixe durch
  (Review F-1); die Fixrunde zog den Code auf die Zusage (strikt `X.Y.Z`, Mutation
  `0.7.0-rc.1` rot in C# und Python). (2) Zwei Nachzüge über den Plan hinaus
  (Python- und C#-A2), §3 Tabelle. (3) Die Fremd-Träger-Meldung nannte zuerst nur
  `harness/mk/sdk.mk` (Review F-2, LOW), die Fixrunde ergänzte zwei Träger.
  (4) „Vor jedem Bau“ war an einer Sprache gemessen und für den Lauf über alle
  Sprachen gesagt (V-1, LOW); die Closure berichtigte den Wortlaut.
- **Re-Review:** nötig und gefahren, weil die Fixrunde Logik änderte
  (`version_lesen`, Umbenennung der Pfade beider Modi). Es fand keinen HIGH, aber
  vier Reste: F-1 führende Nullen passieren die Form, F-2 Kopplungssatz ohne
  mitzuändernde Handlung, F-3 Grenze 8 ohne Sprach-Einschränkung, F-4 Suchmuster
  ohne `d06`. Regel: `.claude/commands/implement-slice.md` Schritt 21.
- **Re-Review- und Verifikations-Pflichten in der Closure:** F-4 → Suchlauf-Zeilen
  9/10 mit `d06` (§3, Parent `9d1edf2a` 26, Arbeitsbaum 0). F-2 → der Kopplungssatz
  in `tools/harness/sdk-kompat/csharp/run.sh` nennt Bedingung und Zielpfad der
  RUN-Schicht `lib-dist` als Auslöser und `[ -f "$neu_nupkg" ]`/`neu_dll` als
  mitzuändernde Stelle (`AGENTS.md` §3.7). F-3 → Grenze 8 gilt für C# und Kotlin,
  Python-A5 vergleicht keine Abhängigkeiten. F-1 → Grenze 10 im Vertrag. V-1 →
  Vertrag §Versionen und Kopf des Runners sagen „vor dem Bau dieser Sprache“ und
  nennen die Schleife; der Code bleibt. V-2 → siehe *Befunde*. V-3 → DoD-Haken.
- **Träger-Nachzug bei der Closure** (`AGENTS.md` §3.13, Frist Closure):
  *nachgezogen*, versionsneutral auf „SDK des Arbeitsstands“ bzw. „Quelle des
  Arbeitsstands“: `harness/mk/sdk.mk` (Kommentar und Hilfetext von
  `test-sdk-altserver`), `harness/targets/sdk-altserver.md` (Titel; kein Verweis auf
  den Titel-Anker, `git grep 'sdk-altserver.md#'` 0 Treffer),
  `tools/harness/run-sdk-altserver-tests.sh` (Kopf). Suchlauf Zeile 7/8: `75d2dee8` 4,
  Arbeitsbaum 0.
- **Befunde und Folgearbeit:**
  - **V-1, Code-Seite:** eine Vorprüfung der Versionen aller Sprachen vor der
    Schleife, damit ein Eingabefehler den Lauf vor jedem Bau beendet. Kein
    Folge-Slice: der Ausgang ist schon heute Exit 2, der Vertrag beschreibt das
    Verhalten, der Preis ist Bauzeit der vorderen Sprachen. Adresse: die nächste
    Arbeit an `tools/harness/run-sdk-kompat-tests.sh`.
  - **V-2:** am Stand `0.6.1` ist die Quellseite (A5) und die Gegenrichtung (A3) nur
    in Python gegen die neue Version gemessen; C# und Kotlin übersetzen beide gegen
    `REGISTRY_VERSION` 0.6.0 (gedruckt `KOMPAT csharp A5 quelle 0.6.0`,
    `KOMPAT kotlin A5 quelle 0.6.0`). Die Binärseite (A2) ist in allen drei Sprachen
    für 0.6.1 gemessen. Der Vertrag nennt das (§Versionen, §Grenze 7/8).
  - **Re-Review F-1:** führende Nullen — Grenze 10 im Vertrag, keine Code-Änderung
    (kein Package trägt eine solche Version).
  - **Review F-3 bis F-6, Re-Review F-2 bis F-4:** in Fixrunde bzw. Closure
    getragen (oben).
- **Steering-Loop-Eintrag:** zwei Lernpunkte, keine neue Regel im Text, kein neuer
  Sensor. (a) *Eine Eingabeprüfung in einer Schleife wird an mehr als einem
  Durchlauf gemessen:* alle Mutationen liefen mit einer Sprache, deshalb sah keine,
  dass „vor jedem Bau“ für drei Sprachen nicht gilt; geschärfte Rückfrage an die
  Mutationszeilen eines Plans: deckt die Eingabe die Mehrzahl, für die die Aussage
  gilt? (b) *Wer einen Wert prüft, misst die Zusage an einem Wert außerhalb der
  Menge, die heute im Baum steht:* alle drei Packages trugen `X.Y.Z`, deshalb blieb
  die Suffix-Lücke bis zum Review unsichtbar. Beide gezählt unter
  `BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad` (Ausgang verkörpert im
  Reviewer-Skill, Klausel *Zusage*; die Probe „den zugesagten Pfad nachfahren“ hat
  F-1 gefunden).
- **Beobachtungs-Register (`../observations/`):**
  `evidence/slice-sdk-kompat-version-parameter.md` in
  `BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad/` angelegt (Review F-1 HIGH,
  Verifikation V-1 LOW; ein Vorgang, eine Datei) — Zähler 10×, Ausgang unverändert
  verkörpert, `state.md` nennt jetzt den Deckel. Unter dem Deckel ohne Datei:
  Review F-2 (LOW) in `BEO-PGC/arbeit-ueberholt-stehenden-traeger` (Fremd-Träger-Meldung
  unvollständig, 34×) und Re-Review F-4 (LOW) in
  `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht` (Suchlauf-Zeile stützt einen Satz
  über fünf Namen mit einem Muster über vier, 22×). Re-Review F-2 (INFO, Kopplung ohne
  mitzuändernde Handlung) ohne Register-Eintrag: getragen in der Closure, keine
  wiederkehrende Klasse im Register. Kein Eintrag erreicht mit diesem Slice die
  Schwelle 3× neu; kein Lese-Schritt fällig.
- **Folge-Slices:** keine (V-1-Code-Seite benannt mit Adresse, oben).
- **Risiken aus §6:** je ein Ausgang, siehe §6 — beide entfallen (Version aus der
  Quelle gemessen; Kotlin und Python vor dem Fix gemessen).
- **Drei Paarungen:** nach dem `git mv` gemessen — (a) *Anker:* kein Eintrag dieser
  Notiz trägt das Feld `liegt in` (nichts verkörpert), kein Gegenstand der Paarung;
  (b) *Folge-Slice:* keiner genannt; (c) *Register:* die drei genannten Verzeichnisse
  (`BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad`,
  `BEO-PGC/arbeit-ueberholt-stehenden-traeger`,
  `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht`) existieren, jedes mit nicht leerem
  `evidence/` (10, 34, 22 Dateien). Ergebnis: getragen. Durch den Move brach kein
  Verweis (`make docs-check` 0 Befunde; `git grep` nach dem `in-progress`-Pfad dieses
  Plans: 0 Treffer). Der Ruhe-Marker der Roadmap steht wieder, Wortlaut gleich
  `6b179119` (`diff` Exit 0); `in-progress/` trägt nur `roadmap.md`.
- **Gates der Closure:** `make docs-check` Exit 0 nach dem Move; `make gates` und
  `make suchlauf-nachmessen` nach dem letzten Commit stehen im Bericht der Sitzung,
  nicht in diesem Plan.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist `*` (Default; `tools/harness`,
`harness/`); eine Ausdifferenzierung ist nicht nötig.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen; kein Treffer
zu `sdk-kompat`.

Alle berührten Sub-Areas GF.
