# Architect-Verdikt: Form-Vorbild-Kopie und die Sprachreinheit je Form-Teil

**Rolle:** Architect (Modul 8)
**Anlass:** die offene Übergabe aus dem Lese-Schritt der
`welle-sdk-reale2e`-Closure
([`../plan/planning/done/welle-sdk-reale2e-results.md`](../plan/planning/done/welle-sdk-reale2e-results.md)
§Steering-Loop-Einträge, zweiter Eintrag):
`BEO-PGC/formvorbild-kopie-traegt-deutsches-wortfragment-weiter` erreichte
mit dem dritten Beleg die 3×-Schwelle, und die Closure ließ die Frage
offen — „Reviewer-Skill-Lese-Pflicht je Form-Teil oder Fitness
Function?". Das `state.md` des Eintrags trägt dieselbe Übergabe (Zustand
`offen`, Lesungs-Vermerk). Mitgeschlossen wird die Schwesterklasse
`BEO-PGC/deutsches-fachwort-im-englischen-sdk-readme` — 3×, derselbe
Ausgang (offen, Architect-Entscheidung, offen geblieben seit der
`welle-sdk-kotlin-lh-fa-sst-009`-Closure), dieselbe Mechanik: sie teilt
sich mit dieser Entscheidung Träger und Wortlaut (siehe §Verdikt).
**Rolleninhaber:** pt9912 (Architect-Rolle, dieser Lauf)
**Datum:** 2026-09-23
**Bezug:** die Beleg-Dateien beider Klassen
(`docs/plan/planning/observations/BEO-PGC/formvorbild-kopie-traegt-deutsches-wortfragment-weiter/evidence/slice-sdk-csharp-projektgeruest.md`,
`…/slice-sdk-kotlin-projektgeruest.md`,
`…/slice-sdk-python-http-reale2e.md`;
`docs/plan/planning/observations/BEO-PGC/deutsches-fachwort-im-englischen-sdk-readme/evidence/`
mit den drei `projektgeruest`-Belegen), Review zu
`slice-sdk-csharp-projektgeruest` (F-1), Review zu
`slice-sdk-kotlin-projektgeruest` (F-1, F-2), Review zu
`slice-sdk-kotlin-reale2e` (F-2, F-3 — Runner-Familie und
Quell-Fragment), Review zu `slice-sdk-python-http-reale2e` (F-4),
`docs/plan/planning/done/slice-sdk-csharp-reale2e.md` §6 (die je-teilige
Sichtung als Zusage, real gegangen), `.harness/skills/reviewer.md` (Träger
der Entscheidung), `AGENTS.md` §3.6 (Gate = ADR), §3.12, §3.13 (die
benachbarten Anker), [`../plan/adr/0106-csharp-nuget-erstes-sdk-package.md`](../plan/adr/0106-csharp-nuget-erstes-sdk-package.md),
[`../plan/adr/0107-python-pypi-zweites-sdk-package.md`](../plan/adr/0107-python-pypi-zweites-sdk-package.md),
[`../plan/adr/0109-kotlin-github-packages-drittes-sdk-package.md`](../plan/adr/0109-kotlin-github-packages-drittes-sdk-package.md)
(je Festlegung 3 — die Formvorbild-Praxis, die die Klasse trägt), Modul 4
§Verglichene Alternativen, Modul 8 §Rollen-Regeln, Modul 10 §Pflege,
Modul 13 §Kein behauptetes Gate ohne Deckung

---

## Frage

Trägt die Klasse „Form-Vorbild-Kopie trägt ein sprachgebrochenes
Wortfragment weiter" eine verkörperte Regel — und wenn ja, an welchem
Träger und in welcher Kategorie? Zu prüfen ist dabei die Familien-Grenze,
die die Klasse bisher trug (deutsche Runner-Kommentar-Familie in
`tools/harness/` ≠ englische SDK-Doc-Familie): der dritte Beleg saß in
Plan-Prosa — trägt die Zwei-Familien-Form ihn noch?

## Empirischer Befund

**Sechs Vorkommen in zwei Klassen; im laufenden Vorgang gefunden hat die
Klasse durchgehend der Reviewer** — die zwei rückwirkend erfassten
README-Instanzen wurden bei derselben Kotlin-Planner-Closure nachgemessen,
nachdem deren Review die Klasse formal benannt hatte. Die Zähler beider
Einträge stehen je bei 3×; die Mechanik teilen sie vollständig.

**Klasse `formvorbild-kopie` (3×):**

1. **C#-Erstauftreten** (`slice-sdk-csharp-projektgeruest`, Review F-1,
   LOW): das Fragment „unstrittige" gelangte aus der deutschen §1-Prosa
   des Plans („ein echter, unstrittiger gemeinsamer Nenner") frisch
   formuliert in den englischen XML-Doc-Kommentar von
   `PgChangeFeedClientOptions.cs` — kein Vorbild war beteiligt, der Slip
   sitzt am Übersetzungsübergang. Keine Fixrunde.
2. **Kotlin-Kopie** (`slice-sdk-kotlin-projektgeruest`, Review F-1, LOW):
   derselbe Kommentartext wortgleich aus dem C#-Vorbild übernommen
   (`PgChangeFeedClientOptions.kt:9`) — der Kopier-Vorgang trug das
   Fragment weiter. Die Korrektur blieb „dem nächsten Slice überlassen,
   der dieselbe Datei ohnehin berührt"; keine Fixrunde.
3. **Python-Plan-Nachzug** (`slice-sdk-python-http-reale2e`, Review F-4,
   LOW): „degenerater Pfad" wanderte wortgleich aus der Kotlin-Plan-§7
   (`slice-sdk-kotlin-reale2e.md:359`) in den Plan-Nachzug §3 des
   Python-Plans — die §8-Deklaration „dieser Slice kopiert kein
   Form-Vorbild" prüfte ihre eigene Prosa nicht mit. Gezogen in der
   Fixrunde (`d668b9cc`), vom Verifier gegen beide Stände bestätigt.

**Klasse `deutsches-fachwort` (3×):** derselbe Mechanismus an den
README-Eröffnungssätzen aller drei Sprachpakete — das interne deutsche
Fachwort „Vollinhalt" ([`../plan/adr/0100-nats-dritter-vollinhalts-zustellweg.md`](../plan/adr/0100-nats-dritter-vollinhalts-zustellweg.md))
stand unflektiert im englischen Satz („NATS-vollinhalt delivery remain
out of scope"), je Paket in eigener Formulierung aus derselben internen
Terminologie; nur die Kotlin-Instanz war
formal befundet, die zwei übrigen rückwirkend erfasst. **Heute gezogen:**
`grep -rni vollinhalt sdks/*/README.md sdks/*/pgchangefeed-kotlin/README.md`
→ 0 Treffer (gemessen in diesem Zug).

**Bestandsmessung in diesem Zug** (`grep -rni unstrittige sdks/`): die
beiden Code-Kommentar-Fragmente aus 1 und 2 stehen **bis heute im
Bestand** (`PgChangeFeedClientOptions.cs:7`,
`PgChangeFeedClientOptions.kt:9`) — die als LOW an den „nächsten Slice,
der dieselbe Datei berührt" weitergereichte Korrektur hat sie nie
gezogen, obwohl die benannten Folgeslices real gelaufen sind. Das ist die
gemessene Konsequenz der LOW-Weiterreichung für diese Klasse.

**Die Grenzfälle um die Klasse herum bestätigen die Abgrenzung.** Die
Runner-Familie („die vier Flaeche", „FlaecheN" — byte-identische Kopien
zwischen den `tools/harness/`-Köpfen) trägt dieselbe
Kopier-Dynamik an orthografischen Formen; beide Reviews führten sie
korrekt als eigene Nachbar-Form (Kotlin Review F-2: „Form-Fragment in der
Vorbild-Kopie"), **außerhalb** des Zählers dieser Klasse — und bemerkten
selbst, dass die zugesagte separate Sprachreinheits-Prüfung diesen Treffer
nicht ausgesiebt hat. Die deutschen „Vollinhalt"-Nennungen in den
Build-Datei-Kommentaren der SDK-Bäume (`.csproj`, `build.gradle.kts`) sind
an ihrem Ort — deutschsprachiger Kommentar, kein sprachfremdes Fragment
(gemessen: 5 Treffer, alle in deutschen Kommentaren). **Positiv-
Kontrolle:** `slice-sdk-csharp-reale2e` §6 deklarierte die je-teilige
Sichtung als Zusage, sie lief real und trug — kein unübersetztes Fragment
in Runner, Dockerfile und Make-Target-Kommentar der C#-Kette.

## Diagnose

**Die tragende Fundstruktur über alle sechs Vorkommen:** ein Form-Teil,
der aus der deutschen Plan-Prosa stammt (frisch formuliert oder aus einem
Vorgänger-Plan übernommen) oder aus einem Formvorbild kopiert wurde, trug
ein Fragment, das die Sprache seines Trägers bricht — und die
Sprachreinheits-Sichtung des Form-Teils selbst wurde nie ausgeführt.
Stattdessen standen zwei Ersatzformen ein: die **wortgleiche
Übereinstimmung mit dem Vorbild** (Fall 2; Runner-„FlaecheN") und die
**§8-Deklaration „nicht einschlägig"** (Fall 3). Der Reviewer war die
Fundinstanz jedes laufenden Vorgangs; die zwei rückwirkend erfassten
README-Instanzen wurden erst nachgefasst, nachdem der Kotlin-Review die
Klasse benannt hatte — kein Implementer- oder Verifier-Lauf hat die
Klasse je selbst gefangen.

**Die Familien-Grenze trägt den dritten Beleg nicht.** Die Zwei-Familien-
Form der Welle-Closure (deutsche Runner-Kommentar-Familie ≠ englische
SDK-Doc-Familie) benennt die Plan-Prosa nirgends — die Closure musste die
Lehre deshalb als Zusatzsatz führen („auch in Plan-Prosa"). Die Grenze
wird hier auf die Mechanismusebene neu gefasst, unter der alle drei
Belege und die README-Klasse fallen: **die Sprache des Trägers
entscheidet, nicht die Datei-Familie.** Ein Fragment, das die Sprache
seines Trägers bricht — das unübersetzte deutsche Fachwort im englischen
Satz, das Hybrid-Fragment in der deutschen Prosa — ist Verstoßfläche, wo
auch immer der Träger steht (SDK-Doku, Code-Kommentar, README,
Plan-Prosa); ein deutsches Fachwort in einem deutschsprachigen Kommentar
ist an seinem Ort (Runner-Köpfe, Build-Datei-Kommentare). Die orthografische
Form-Familie der mono-lingualen Kopien („FlaecheN") bleibt die
Nachbar-Form der jeweiligen Reviews — sie ist kein zweiter Zähler dieser
Klasse geworden, und diese Entscheidung ändert das nicht.

**Die Kategorie war Teil des Fehlers.** Alle drei Vorkommen der
formvorbild-Klasse liefen als LOW mit deferral — „keine Fixrunde, der
nächste Slice, der dieselbe Datei berührt" — und genau dieser Pfad hat
die Korrektur nie gezogen (Bestandsmessung oben). Die Klasse ist kein
„einmaliger Tippfehler": die Kopie macht das Fragment zum Muster, das in
jedem künftigen Geschwister-Paket wieder ankommt (die README-Familie
stand bei 3/3 Paketen, bevor sie gezogen wurde). Der benannte Anker
verhindert die Subsumtion unter LOW — dieselbe Begründung wie beim
Architect-Verdikt zur Zusage ohne Bindung an ihre Eingabeseite, wo die
Subsumtion unter die Nachbar-Kategorie den Ausgang der Fixrunde-Frage
stellte.

## Verglichene Alternativen

**(a) Reviewer-Skill-Zeile (benannter HIGH-Punkt mit Lese-Pflicht je
Form-Teil) — gewählt.** Der Reviewer ist die Fundinstanz jedes
laufenden Vorgangs; die drei Präzedenz-Punkte derselben Liste („Zahl im Träger
ohne Ursprung", „Beleg trägt seinen Satz nicht", „Zusage ohne Bindung an
ihre Eingabeseite") sind dieselbe Hausform für Leser-Disziplin ohne
Sensor. Pro: kleinste Form mit realer Wirkung — der Probe-Satz ist direkt
anwendbar (Vorbild in die Hand, je Form-Teil sichten). Contra: eine
Lese-Pflicht ist nicht mechanisch erzwingbar — sie trägt so weit, wie der
Reviewer sie ausführt.

**(b) Fitness Function/Gate — verworfen.** Ein mechanischer Prüf-Begriff
existiert nicht: Prosa-Orthografie und Sprachfremdheit sind durch keines
der drei Sensor-Prinzipien dieses Repos entscheidbar (Referenz,
Coverage-Zahl, Byte-Vergleich); ein Wörterbuch- oder
Sprachdetektor-Gate bräuchte eine Semantik-Entscheidung je Träger, welches
Wort sprachfremd ist, und müsste Prosaflächen lesen, die kein Sensor des
Repos liest. Modul 13 (kein behauptetes Gate ohne Deckung) und
`AGENTS.md` §3.6 (Gate = ADR) sperren den Weg; der Kotlin-Review hat die
Grenze bereits so notiert („verifizierbar: nein — Text lesen; kein Sensor
für Prosa-Orthografie"). Zusätzlich die Verhältnismäßigkeit: ein
blockierendes Gate für einen Befund ohne semantische Auswirkung wäre
überstark — die
Falsifikation dieser Klasse ist das Lesen, nicht der Lauf.

**(c) Keine Schärfung (Regel trägt sich bereits) — verworfen.** Der
Lese-Schritt der Welle-Closure hat die benachbarten Anker explizit
geprüft und verneint: `AGENTS.md` §3.13 trägt den Suchlauf über Träger
einer bewegten Eigenschaft, §3.12 die Zahl-Herkunft — keine der beiden
Formen nennt die Sprachreinheit eines kopierten Form-Teils. Eine echte
Lücke bleibt, und die Bestandsmessung trägt ihren Preis: zwei Fragmente
stehen seit ihren Slices unverändert im Bestand, weil kein stehender
Träger die Sichtung verlangt hat — die je-teilige Sichtung existierte
bisher nur als freiwillige §6-Zusage einzelner Pläne, und der eine Plan,
der sie brauchte (Python-HTTP), hatte sich per §8 selbst ausgeschlossen.

## Verdikt

**Verkörpert — ein Träger, ein benannter HIGH-Punkt.**
`.harness/skills/reviewer.md` trägt in der HIGH-Liste (nach „Zusage ohne
Bindung an ihre Eingabeseite") den Punkt „**Form-Vorbild-Kopie trägt ein
sprachgebrochenes Wortfragment weiter**" in der Hausform der übrigen
repo-spezifischen HIGH-Punkte: benannte Klasse, Probe in Kursivschrift
(jeder Form-Teil einzeln gesichtet, mit dem Formvorbild in der Hand, die
wortgleiche Übereinstimmung als nicht-Beleg), die neu gefasste
Familien-Grenze (Sprache des Trägers entscheidet, nicht die
Datei-Familie), die Kein-Gate-Begründung und die Herkunft beider Klassen
mit dem `seit`-Anker.

**Warum ein Träger genügt.** Beim Zusage-ohne-Bindung-Präzedenzfall
wurden zwei Träger verkörpert, weil dort eine Implementer-Disziplin
(Schritt 19, Mutations-Pflicht) bereits lief und nur die Richtung
fehlte. Hier existiert **kein** Workflow-Schritt, der die Sichtung
bereits ausübt — die je-teilige Sichtung war bisher freiwillige
Plan-Deklaration —, und eine neue Implementer-Pflicht einzuführen wäre
ein neuer Prozesskörper ohne Beleg: in allen Vorkommen war der Fix eine
Ein-Zeilen-Änderung, und der Reviewer hat die Klasse in jedem Vorgang
gefunden, bevor sie sich weiter ausbreitete (der Python-Plan-Nachzug ist
das einzige Vorkommen, das im selben Vorgang gezogen wurde — Fixrunde
`d668b9cc`). Die Lese-Pflicht beim Reviewer sitzt an der einzigen Stelle,
die den
Form-Teil **und** sein Vorbild im selben Kontext sieht — genau der
Kontext, der in Fall 3 fehlte, als die §8-Deklaration ihre eigene Prosa
nicht mit las.

**Die Kategorie: HIGH, nicht LOW.** Die LOW-Weiterreichung an „den
nächsten Slice" ist für diese Klasse real gescheitert (Bestandsmessung
oben); die Kategorie entscheidet über die Fixrunde, und die Fixrunde im
eigenen Slice ist die einzige Behandlungsform, die in den Vorkommen
nachweislich zog (Fall 3, `d668b9cc`). Der Punkt
bekommt damit einen eigenen Anker statt der Subsumtion unter LOW
(„einmalige Tippfehler") — dieselbe Begründung wie bei den
Schwester-Präzedenzen derselben Liste.

## Was dieses Verdikt NICHT tut

- Es schreibt **keine** ADR: kein Hard Rule in `AGENTS.md`, kein Gate,
  keine Schwelle — die Verkörperung ist ein Rollen-Skill, die Hausform
  dafür ist ein Verdikt (Vergleich: der Verdikt zur Zusage ohne Bindung
  an ihre Eingabeseite, 4×).
- Es bumpt **keinen** Zähler und legt **keine** Beleg-Datei an; beide
  Klassen bleiben bei real 3×, die Ausgänge stehen in ihren `state.md`.
- Es fasst **keinen** Produktionscode an — auch nicht die zwei im Bestand
  verbliebenen „unstrittige"-Kommentare; sie sind als realer Befund im
  Herkunfts-Feld des neuen HIGH-Punkts verankert (gemessen, auflösbar)
  und ziehen, wenn der nächste Slice die Dateien berührt — mit dem
  HIGH-Anker hinter der Sichtung.
- Es ändert die **Familien-Abgrenzung der Runner-Familie nicht**: die
  orthografische Nachbar-Form („Flaeche"/„FlaecheN") bleibt außerhalb
  dieses Zählers; erreicht sie unter eigener Buchführung 3×, entscheidet
  ein eigener Zug.
- Es führt **keine** neue Rolle und **keinen** neuen Workflow-Schritt ein.

## Konsequenzen für die Rollen

**Reviewer:** ein Diff, der Form-Teile aus einem Formvorbild oder aus der
deutschen Plan-Prosa trägt, liest er je Form-Teil mit dem Vorbild in der
Hand; ein sprachgebrochenes Fragment ist ein HIGH-Finding mit eigenem
Anker, nicht mehr „einmaliger Tippfehler". Die Negativbefund-Zeile des
Reports trägt das Ergebnis der Sichtung je betrachtetem Form-Teil.

**Planner:** die §8-Deklaration ist Review-Prosa — „dieser Slice kopiert
kein Form-Vorbild" muss ihren eigenen Text überstehen. Kein neuer
Plan-Abschnitt, keine neue DoD-Zeile; die Zusage-Form aus
`slice-sdk-csharp-reale2e` §6 bleibt das freiwillige Vorbild, nicht
Pflicht.

**Implementer:** ein HIGH-Befund auf einem kopierten Fragment ist eine
Fixrunde im selben Slice — der Weiterreich-Pfad („der nächste Slice, der
dieselbe Datei berührt") ist für diese Klasse geschlossen.

**Verifier:** nichts Neues — die gezogene Korrektur wird wie jede
andere gegen beide Stände geprüft.

## Re-Evaluierungs-Trigger

- **Ein mechanischer Prüf-Begriff entsteht** — z. B. ein Doku-Lint mit
  realem deutschem/englischem Lexikon über die SDK-Doku-Flächen: dann
  öffnet sich die Gate-Frage neu (`AGENTS.md` §3.6, Modul 13 — der
  Kein-Gate-Satz des HIGH-Punkts trägt seinen Grund, und der Grund wäre
  entfallen).
- **Ein viertes Auftreten nach der Verkörperung** — ein kopiertes
  Fragment, das die benannte Probe passiert hat: dann schärft der
  nächste Zug Probe oder Familien-Grenze nach (Modul 10 §Pflege, neuer
  Beleg, neuer `evidence/`-Eintrag).
- **Die Nachbar-Form erreicht unter eigener Buchführung 3×** (die
  orthografische Runner-Kopie-Familie): eigene Entscheidung, kein
  Vorgriff durch dieses Verdikt.

## Was dieser Zug geändert hat — und was nicht

**Geändert:**

- `.harness/skills/reviewer.md` (ein neuer HIGH-Punkt in der HIGH-Liste)
- `state.md` beider Klassen (Ausgang `verkörpert` mit Anker auf den
  Punkt)
- diese Verdikt-Datei

**Nicht geändert:** `AGENTS.md` (keine neue Hard Rule), jede ADR (auch
keine `Supersedes`-Kette), kein Gate/Sensor/Makefile, kein
Slice-/Welle-Plan, `spec/**`, `internal/**`, `sdks/**` (die zwei
„unstrittige"-Kommentare stehen unverändert im Bestand — gemessen,
verankert im Herkunfts-Feld), `tools/**`. Kein Commit in diesem Zug —
der schreibende Commit folgt außerhalb des Verdikts.