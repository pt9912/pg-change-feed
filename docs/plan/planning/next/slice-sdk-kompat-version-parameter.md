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

**Verantwortlich:** — bis zur Priorisierung.

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

- [ ] **Version aus der Quelle:** kein `0.6.0`-Literal im Modus `dist` mehr in
      `tools/harness/sdk-kompat` (Suchlauf: die Trefferzahl der Stellen ohne die
      benannte Konstante des Modus `registry` ist 0); die Version je Sprache
      kommt aus `.csproj`, `build.gradle.kts` bzw. `pyproject.toml`.
- [ ] **Messung:** `make test-sdk-kompat` im Standardmodus Exit 0 am Arbeitsstand
      (gedruckte Schlusszeile, Exit direkt gelesen), `SDK_KOMPAT_NEU=registry`
      unverändert Exit 0; eine Mutation der Version (Eingabeseite: ein anderer
      Wert in einer Version-Datei einer Kopie im Scratchpad) färbt den Lauf rot
      mit lesbarer Meldung statt eines stillen Rückfalls auf eine andere Version.
- [ ] **Vertrag nachgezogen:** `harness/targets/sdk-kompat.md` und die Kommentare
      in `harness/mk/sdk.mk`, `tools/harness/run-sdk-kompat-tests.sh` nennen keine
      feste Version mehr als Gegenstand des Standardmodus (Lese-Handlung des
      Reviewers; der Suchlauf in §3 deckt nur `tools/harness/sdk-kompat`).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`), kein Self-Review.
- [ ] Doku-Update: nur der Target-Vertrag (kein öffentlicher Vertrag berührt).
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

```suchlauf
3e028ee3 39 -n 0\.6\.0 -- tools/harness/sdk-kompat
```

Stand der Zeile: **gemessen** mit `git grep -n` am Parent `3e028ee3` (neun Dateien,
je Datei 8, 1, 12, 4, 8, 1, 1, 1, 3 Treffer); der Implementer wiederholt sie mit
`make suchlauf-nachmessen PLAN=<Plan-Datei>`. Der Suchlauf trifft das Literal,
nicht jede Stelle, die eine Version implizit trägt (Dateinamen mit Versionsteil,
Umgebungsvariablen): das liest der Reviewer.

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
