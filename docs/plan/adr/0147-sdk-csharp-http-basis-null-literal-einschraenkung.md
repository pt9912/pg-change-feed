# ADR-0147: C#-SDK 0.6.0 — „quellseitig erhalten“ gilt nicht für ein `null`-Literal an der dritten Stelle der protected HTTP-Basis (Schärft ADR-0145 Festlegung 4)

**Status:** Accepted — **kein** Supersedes.

**Datum:** 2026-10-04

**Autor:** pt9912 (Architect-Rolle, Modul 8)

**Bezug:** [`LH-FA-SST-009`](../../../spec/lastenheft.md) (SDK-Packages),
[ADR-0145](0145-sdk-meldungscodes-eigenschaft-der-fehlertypen.md) (Festlegung 4 ist der
Gegenstand), [ADR-0106](0106-csharp-nuget-erstes-sdk-package.md) (C#-Package),
[ADR-0083](0083-herkunft-von-aussagen-in-traegern.md), `AGENTS.md` §3.5, §3.12.

**Schärft:** [ADR-0145](0145-sdk-meldungscodes-eigenschaft-der-fehlertypen.md) Festlegung 4
(Satz „Damit bleibt jede 0.5.x-Signatur binär und quellseitig erhalten“ und Satz „Ein Eintrag
‚Upgrading 0.6.0‘ in den READMEs entfällt“, nur für C#). ADR-0145 bleibt unverändert
(`AGENTS.md` §3.5). Spec-Stelle: keine (`SPEC-026` bis `SPEC-028` nennen Fehlertypen nicht,
*übernommen* aus ADR-0145).

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

**(1) Gemessener Befund.** `make test-sdk-kompat` (C#-null-Matrix, 14 Fälle je Version;
Slice `sdk-0-6-kompatibilitaet-messen` §7, Review F-1): 13 Fälle `0.5.0 ok, 0.6.0 ok`, der Fall
`http-basis-3-argumente-null-literal` `0.5.0 ok, 0.6.0 CS0121`. Die abstrakte Basis
`PgChangeFeedException` (`sdks/csharp/PgChangeFeed.Client/Http/PgChangeFeedException.cs`) hat
seit 0.6.0 neben `(int, string, Exception)` den Konstruktor `(int, string, string?)`; ein
`base(500, "x", null)` einer fremden Unterklasse passt auf beide. Gemessen sind ebenso: Cast
(`(Exception)null`) und benanntes Argument (`innerException: null`) übersetzen unter 0.5.0 und
0.6.0; die öffentlichen Blatt-Typen und die gRPC-Basis zeigen den Fehler nicht; die
Binärkompatibilität trägt (A2, drei Sprachen).

**(2) Wer leitet die Basis ab.** Die Basis ist `abstract` mit `protected` Konstruktoren; die
Blatt-Typen sind `sealed` mit öffentlichen Konstruktoren. Gemessen am Baum (`git grep -E
':\s*PgChangeFeedException\b' -- '*.cs'`, 2026-10-04): Ableitungen außerhalb des Packages
stehen nur in den Mess-Gästen unter `tools/harness/sdk-kompat/csharp/`; `examples/` und `test/`
leiten sie nicht ab (0 Treffer). Fremde Nutzung über NuGet ist nicht messbar; die Erwartung
„selten“ ist *hergeleitet* (Ableiten einer Fehler-Basis ist nicht der vorgesehene Gebrauch;
der Anwender fängt die Typen).

**(3) Stand der Veröffentlichung.** 0.5.0 und 0.6.0 sind auf NuGet. Ein Code-Fix wäre eine
neue Version 0.6.1.

## Entscheidung

Wir wählen **die Einschränkung annehmen, kein Code-Fix, mit Hinweis im C#-README**.

1. **Geltung.** Die Aussage „quellseitig erhalten“ aus ADR-0145 Festlegung 4 gilt für das
   C#-SDK **mit einer benannten Ausnahme**: ein `null`-Literal an der dritten Stelle eines
   Aufrufs der protected Basis-Konstruktoren von `PgChangeFeedException` (HTTP) ist unter 0.6.0
   mehrdeutig (`CS0121`). Die öffentlichen Konstruktoren der Blatt-Typen, die gRPC-Basis,
   Kotlin und Python sind nicht betroffen (C#/Kotlin/Python: gemessen im genannten Slice für
   die dort geprobten Formen; die Menge sind die 14 Fälle der Matrix, 15 Python-Typen und die
   Gast-Quelltexte; jede Verallgemeinerung darauf hinaus ist *hergeleitet*).
2. **Umgehung** (gemessen: Cast und benanntes Argument übersetzen unter beiden Versionen):
   `base(status, text, (Exception?)null)` oder `base(status, text, innerException: null)`.
3. **Hinweis.** Das C#-README bekommt unter „Upgrading“ einen Eintrag 0.6.0 für Ableitungen der
   Basis (Wortlaut ohne interne Kennung, `ADR-0134`). Damit entfällt für C# der Satz
   „Upgrading 0.6.0 entfällt“ aus ADR-0145; für Kotlin und Python bleibt er.
4. **Kein 0.6.1 aus diesem Grund.** Der Hinweis erscheint mit dem nächsten ohnehin fälligen
   C#-Release (das README reist mit dem Package; das schon veröffentlichte 0.6.0 trägt ihn
   nicht — *hergeleitet*, am Package nicht geprüft).

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun, nur die ADR-Aussage berichtigen | kein Aufwand | der Anwender erfährt die Umgehung nicht, ein Treffer ist ein Rätsel an einer Compiler-Meldung |
| **B — Einschränkung annehmen + README-Hinweis (gewählt)** | kleinster Aufwand mit Wirkung beim Betroffenen; kein Release; Umgehung gemessen | die Ausnahme bleibt bestehen; Betroffene sehen den Hinweis erst im nächsten Release-README oder im Repo |
| C — Code-Fix, 0.6.1: Konstruktor `(int, string, string?)` entfernen, nur die Vierer-Form `(int, string, Exception?, string?)` anbieten | beseitigt `CS0121` | entfernt eine in 0.6.0 veröffentlichte protected Signatur (bricht gegen 0.6.0 gebaute Ableitungen, Binärkompatibilität zu 0.6.0), neue Version und Slice (Code, Tests, `sdk-kompat`-Matrix, Release) für eine selten genutzte Form; ob es ohne den Konstruktor ginge, ist nicht gebaut und nur *hergeleitet* (zwei Überladungen mit Referenztyp an gleicher Stelle sind bei `null` immer mehrdeutig) |
| D — Konstruktor `(int, string, string?)` zusätzlich als `[Obsolete]`/Kennzeichnung | — | ändert die Mehrdeutigkeit nicht |

## Konsequenzen

- Positiv: keine neue Version, keine Binärbrüche; die Aussage der Spec-Kette ist ehrlich.
- Negativ: ein Anwender, der die Basis fremd ableitet und ein `null`-Literal übergibt, bekommt
  beim Upgrade `CS0121` und muss eine Zeile ändern.
- Folgepflicht: (1) README-Eintrag in `sdks/csharp/README.md` §Upgrading (kleiner Nachzug, mit
  dem nächsten C#-Release oder einem Doku-Zug; kein eigener Slice nötig). (2) Der Plan
  `slice-sdk-0-6-kompatibilitaet-messen` nennt den Befund bereits; sein Träger
  `harness/targets/sdk-kompat.md` zitiert die Aussage und die Ausnahme für den Fall
  `http-basis-3-argumente-null-literal` bereits — beim Nachzug auf diese ADR verweisen.
  (3) Die Release-Notiz des nächsten C#-Releases nennt den Hinweis.
- Akzeptiertes Negativ: 0.6.0 auf NuGet bleibt, wie es ist; sein README trägt den Hinweis nicht.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| null-Matrix (erprobt, eine Instanz) | `NullMatrix.cs` (14 Fälle): Fall `http-basis-3-argumente-null-literal` `0.5.0 ok, 0.6.0 CS0121`, alle anderen `ok, ok`; gesehen im Lauf des Slice `sdk-0-6-kompatibilitaet-messen` (*übernommen* aus dem Slice, hier nicht nachgefahren) | `make test-sdk-kompat` |
| README-Hinweis (hergeleitet; nicht erprobt) | der Eintrag steht nach dem Nachzug in `sdks/csharp/README.md`; keine interne Kennung | `make sdk-public-doc-check` |

## Re-Evaluierungs-Trigger

Meldet ein Anwender eine fremde Ableitung der Basis, oder erscheint ohnehin ein C#-Release
mit Konstruktor-Änderungen an der Basis (Option C wird dann billig mitgenommen).

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-04 | Accepted — Architect-Entscheidung zu Review-Befund F-1 | `slice-sdk-0-6-kompatibilitaet-messen` |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0147` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
