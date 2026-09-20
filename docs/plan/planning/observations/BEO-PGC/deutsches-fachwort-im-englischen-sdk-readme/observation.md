# BEO-PGC/deutsches-fachwort-im-englischen-sdk-readme

**Sub-Area:** `sdks/*/` (die SDK-Sprachpakete dieser Repo-Reihe —
`harness/conventions.md` §Modus-Deklaration, Default-Sub-Area `*`/`PGC`
Greenfield; dieselbe Formvorbild-Kopie-Dynamik wie
`BEO-PGC/formvorbild-kopie-traegt-deutsches-wortfragment-weiter`, hier an
einer anderen Textstelle).

Die Beobachtung: Der einleitende README-Satz zum SSE-/NATS-Umfang trägt in
**allen drei** bislang existierenden SDK-Sprachpaketen dasselbe deutsche
Fachwort unflektiert mitten im englischen Satz — „SSE and
NATS-vollinhalt(s) delivery remain out of scope …". Das Wort stammt aus der
internen Terminologie dieses Repos (`Vollinhalts-Zustellweg`, siehe z. B.
[`ADR-0100`](../../../../adr/0100-nats-dritter-vollinhalts-zustellweg.md)
„NATS dritter Vollinhalts-Zustellweg") und wurde beim Formulieren des
README-Satzes nicht ins Englische übersetzt:

- `sdks/csharp/README.md:9` — „SSE and NATS-vollinhalts delivery remain
  uncovered …"
- `sdks/python/README.md:9` — „gRPC, SSE and NATS-vollinhalt delivery
  remain out of scope …"
- `sdks/kotlin/pgchangefeed-kotlin/README.md:9` — „SSE and NATS-vollinhalt
  delivery remain out of scope …"

Alle drei real nachgemessen bei der Planner-Closure von
`slice-sdk-kotlin-projektgeruest` (2026-09-20):
`grep -n "vollinhalt" sdks/csharp/README.md sdks/python/README.md
sdks/kotlin/pgchangefeed-kotlin/README.md` — drei Treffer, einer je Datei.

**Nur die Kotlin-Instanz wurde formal befundet.** Weder
`review-slice-sdk-csharp-projektgeruest.md` noch das Python-Pendant
benannten diese Formulierung als eigene Finding-Klasse — der Fehler war
real bereits zweimal vorhanden, bevor ihn `review-slice-sdk-kotlin-projektgeruest.md`
(F-2) zum ersten Mal formal als Sprachbruch-Klasse benannte. Dieser Eintrag
schließt die Lücke rückwirkend, statt auf ein „viertes" Auftreten zu
warten, das den realen Bestand ignorieren würde.

**Warum das zählt:** Dieselbe Formvorbild-Kopie-Dynamik wie beim
KDoc-Wortfragment — wortgleiche Übernahme trägt unbemerkt Fehler weiter.
Zusätzlich zeigt dieser Fall, dass ein Review, das ausschließlich den
eigenen Diff prüft, eine bereits **zweifach** im Bestand vorhandene
Formulierung nicht als Muster erkennt, solange niemand gezielt gegen den
Bestand der Geschwister-Pakete sucht.

Deklaration: `slice-sdk-csharp-projektgeruest`, `slice-sdk-python-projektgeruest`
(beide rückwirkend erfasst, real vorhanden, nie formal befundet),
`slice-sdk-kotlin-projektgeruest` (dritte Instanz, erste formale Befundung
als F-2 in `review-slice-sdk-kotlin-projektgeruest.md`; Eintrag angelegt
bei dieser Planner-Closure, 2026-09-20).
