# ADR-0093: Digest-Korrektur — `ADR-0087`s Kotlin-Basis-Image-Zeile

**Status:** Accepted — Supersedes [`ADR-0087`](0087-beispiel-clients-csharp-kotlin.md)
in **einer** Tabellenzelle: der `eclipse-temurin:21-jdk`-Zeile ihrer
Digest-Pinning-Tabelle (§Entscheidung, Festlegung 3). Alles Übrige der
[`ADR-0087`](0087-beispiel-clients-csharp-kotlin.md) bleibt **hiermit
bestätigt** und wird nicht wiederholt — auch nicht die bereits von
[`ADR-0090`](0090-beispiel-clients-volle-matrix.md) superseded Umfangs-Klauseln,
die von dieser Korrektur unberührt bleiben.

**Datum:** 2026-09-17

**Autor:** pt9912 (Architect-Rolle, Modul 8 §Konflikt-Pfad, Fall „ADR wird
per Folge-ADR supersedet"; Übergabe-Artefakt für das Reviewer-Finding F-1 im
Review zu `slice-099`, Commit `a65bf29`) <!-- d-check:status-provenance -->

**Bezug:** [`ADR-0087`](0087-beispiel-clients-csharp-kotlin.md) (die
korrigierte Zeile), [`ADR-0073`](0073-zitat-korrektur-an-immutablen-dokumenten.md)
(warum diese Korrektur **keine** Zitat-Korrektur ist), `AGENTS.md` §3.5
(ADR-Immutabilität), `AGENTS.md` §3.12 (Zahl im Träger trägt ihren Ursprung)

**Schärft:** — *(Prozess-Korrektur ohne eigenes Spec-Stratum; die berührte
Aussage ist ein Digest-Wert innerhalb einer ADR, keine Spec-Stelle)*

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Der Reviewer von `slice-099` hat beim Prüfen von `examples/kotlin/Dockerfile` <!-- d-check:status-provenance -->
gegen `ADR-0087`s Digest-Pinning-Tabelle (§Entscheidung, Festlegung 3,
Zeile 254) eine Diskrepanz gefunden
(Review zu `slice-099`, Finding F-1, Commit `a65bf29`): der dort <!-- d-check:status-provenance -->
für `eclipse-temurin:21-jdk` genannte Wert

```
sha256:085eb93e049c7397f725bd8be31c4f52ba4777a75168e851428508234aa224e
```

hat **63** Hex-Zeichen statt der für SHA-256 zwingenden 64 — ein
strukturell ungültiger, nie real auflösbarer Digest. Der `Dockerfile`
verwendet stattdessen bereits

```
sha256:085eb93e049c7397f725bd8be31c4fd52ba4777a75168e851428508234aa224e
```

(`examples/kotlin/Dockerfile:16`) — 64 Hex-Zeichen, ein fehlendes „d"
gegenüber dem Tabellenwert.

**Eigene Verifikation dieses Zugs** (Architect-Rolle, unabhängig
nachgemessen, nicht aus dem Reviewer-Finding übernommen):

- `echo -n "085eb93e049c7397f725bd8be31c4f52ba4777a75168e851428508234aa224e" | wc -c` → `63`
- `echo -n "085eb93e049c7397f725bd8be31c4fd52ba4777a75168e851428508234aa224e" | wc -c` → `64`
- `docker manifest inspect eclipse-temurin:21-jdk` (2026-09-17), Eintrag
  `platform: {architecture: amd64, os: linux}` →
  `digest: sha256:085eb93e049c7397f725bd8be31c4fd52ba4777a75168e851428508234aa224e`
  — deckt sich exakt mit dem `Dockerfile`-Wert, nicht mit dem `ADR-0087`-Wert.

Der Fund ist damit **bestätigt**: `ADR-0087` Zeile 254 enthält einen
Transkriptionsfehler; der real gebaute Code war davon nie betroffen.

**Warum keine Zitat-Korrektur (`ADR-0073`).** Die Tabelle steht innerhalb von
`ADR-0087` §Entscheidung, Festlegung 3 — laut
[`ADR-0073`](0073-zitat-korrektur-an-immutablen-dokumenten.md) Punkt 1 vom
Zitat-Korrektur-Kanal ausgeschlossen, sobald sie §Entscheidung berührt. Ein
Digest-Wert ist zudem keine Instanz des Zitat-/Verweisgerüsts (host-lokaler
Pfad, Linkziel, Zeilen-Lokator) bei unverändertem Referenten, sondern Teil
der Entscheidungssubstanz selbst — der Referent (welcher Digest gilt) ändert
sich mit der Korrektur. Nach `AGENTS.md` §3.5 ist eine In-place-Reparatur
der `Accepted`-ADR deshalb unzulässig; die Korrektur läuft als eigene,
engräumige Folge-ADR mit `Supersedes ADR-0087`.

## Entscheidung

Wir korrigieren **ausschließlich** den `eclipse-temurin:21-jdk`-Digest in
`ADR-0087`s Digest-Pinning-Tabelle (Zeile 254). Kein anderer Teil von
`ADR-0087`s Entscheidung wird berührt oder neu bewertet.

**Korrigierte Tabellenzeile** (ersetzt in der Lesart von `ADR-0087` §Entscheidung
Festlegung 3 die dortige Zeile „Kotlin | `eclipse-temurin:21-jdk` | …"):

| Sprache | Kandidat (Tag) | korrekter gemessener amd64-Digest |
|---|---|---|
| Kotlin | `eclipse-temurin:21-jdk` | `sha256:085eb93e049c7397f725bd8be31c4fd52ba4777a75168e851428508234aa224e` |

Dieser Wert ist identisch mit dem in `examples/kotlin/Dockerfile:16`
tatsächlich gepinnten und gebauten Digest.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun, Finding nur im Review-Report belassen | kein Aufwand | ein künftiger Zug, der `ADR-0087`s Tabellenwert statt eines eigenen Re-Measurements kopiert, pinnt einen strukturell ungültigen, nie auflösbaren Digest; die Tabelle bleibt als Beleg unzuverlässig |
| B — `ADR-0087` in-place korrigieren (Zitat-Korrektur) | einfachste Form, ein Wert an einer Stelle | nach `AGENTS.md` §3.5/`ADR-0073` unzulässig: der Fehler sitzt in §Entscheidung, keine Verweisgerüst-Instanz, sondern Entscheidungssubstanz — eine `Accepted`-ADR wird dadurch **inhaltlich** überschrieben |
| **C — engräumige Folge-ADR mit `Supersedes ADR-0087` für genau diese Zelle (gewählt)** | erfüllt `AGENTS.md` §3.5 vollständig; die Korrektur ist selbst auditierbar (eigene ADR-Nummer, eigene Geschichte-Zeile); `ADR-0087` bleibt unangetastet bis auf die zulässige §Geschichte-Ergänzung; künftige Züge, die den Digest brauchen, finden den korrekten Wert an einer eindeutigen, `Accepted`-Adresse | eine weitere ADR-Datei für einen einzelnen Tabellenwert — bewusst in Kauf genommen, weil die Immutabilitäts-Regel keine Ausnahme für „nur ein Wert" kennt |

## Konsequenzen

- Positiv: `ADR-0087`s Digest-Pinning-Tabelle trägt jetzt einen strukturell
  gültigen, real auflösbaren Wert für `eclipse-temurin:21-jdk` — künftige
  Kopien aus der Tabelle sind sicher.
- Positiv: **Keine Auswirkung auf gebauten Code.** `examples/kotlin/Dockerfile`
  verwendete bereits den korrekten Wert; nichts an ihm wird geändert.
- Negativ: keine — die Korrektur ist rein deklarativ und betrifft nur
  `ADR-0087`s Text.
- Folgepflicht: `ADR-0087` erhält gemäß `AGENTS.md` §3.5 eine Zeile in ihrer
  eigenen §Geschichte-Tabelle (Datum, Ereignis, Commit-Kennung) — die einzige
  zulässige Änderung an ihrer Datei; §Entscheidung/§Konsequenzen/
  §Verglichene Alternativen/§Status/`Supersedes`-Kette bleiben unberührt.
- Folgepflicht: ADR-Index (`docs/plan/adr/README.md`) bekommt die neue
  Zeile für `ADR-0093`.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| — | Kein Sensor prüft Digest-Korrektheit einer ADR-Tabellenzelle gegen die Registry; die Verifikation ist eine einmalige, im Kontext dieser ADR dokumentierte Messung (`docker manifest inspect`), kein Standing-Gate. | — |
| `make docs-check` | Referenzen/Links dieser ADR und der `ADR-0087`-Geschichte-Zeile bleiben auflösbar. | `make docs-check` |

## Re-Evaluierungs-Trigger

Permanent — die Korrektur ist einmalig und abgeschlossen. Ein künftiger
Basis-Image-Wechsel (`eclipse-temurin:21-jdk` wird abgekündigt oder neu
gemessen) läuft über `ADR-0087`s eigenen Trigger 4, nicht über diese ADR.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-17 | Accepted — Anlass: Reviewer-Finding F-1 im Review zu `slice-099` (Commit `a65bf29`); Architect-Rolle verifiziert den Fehler eigenständig (`docker manifest inspect eclipse-temurin:21-jdk`, Zeichenzählung) und korrigiert `ADR-0087`s Digest-Pinning-Tabelle per Folge-ADR (Modul 8 §Konflikt-Pfad, Fall „ADR wird per Folge-ADR supersedet") | Review zu `slice-099` F-1; `docker manifest inspect eclipse-temurin:21-jdk` (amd64/linux, 2026-09-17) <!-- d-check:status-provenance --> |
| 2026-09-18 | Zitat-Korrektur — `docs/reviews/**`-Pfade durch Kennung ersetzt (`ADR-0073`) | PENDING_COMMIT |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0093` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
