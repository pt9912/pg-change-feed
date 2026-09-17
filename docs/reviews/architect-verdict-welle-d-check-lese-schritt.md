# Architect-Verdikt — Lese-Schritt der `welle-d-check`-Closure (zwei fällige Beobachtungs-Einträge)

**Datum:** 2026-09-17 · **Stand:** `7e598bb` · **Rolle:** Architect (eigener
Kontext; alle Zitate und Zeilenbezüge dieses Zugs sind selbst am
referenzierten Original nachgeprüft, nicht aus den Evidence-Dateien
übernommen) · **Anlass:** Modul-6/8-Lese-Schritt der `welle-d-check`-Closure
(Schritt 3b) — zwei Einträge im Beobachtungs-Register haben die 3×-Schwelle
erreicht: `BEO-PGC/regel-weiter-als-ihr-sensor` und
`BEO-PGC/zitat-nennt-die-falsche-stelle`.

---

## Eintrag 1 — `BEO-PGC/regel-weiter-als-ihr-sensor`

### 1.1 Eigene Prüfung der drei Belege

Gelesen: `observation.md`, `state.md`, alle drei `evidence/*.md`,
`AGENTS.md` §3.11/§3.13, `harness/sensors/docs-check.md` §Grenze Punkt 8,
`harness/sensors/coverage-gate.md` §Grenze Punkt 4/5,
`docs/plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md`.

Die drei Manifestationen bestätigen sich, und sie liegen in drei
unterschiedlichen technischen Domänen:

| Beleg | Regel | Sensor-Lücke | Ursache der Lücke |
|---|---|---|---|
| `slice-078` | `AGENTS.md` §3.11 (hostpaths) deckt die **ganze** Markdown-Fläche inkl. Fenced-Blöcke | Modul `hostpaths` lässt Fences **per Design** frei | Werkzeug-Design (d-check), zusätzlich ein **strukturelles** Bedürfnis: `AGENTS.md` §3.11 selbst braucht Fenced-`**Falsch:**`-Beispiele, um eine verbotene Pfadform überhaupt zu illustrieren |
| `slice-079` | `ADR-0071`s Fitness Function verlangt Freiheit von drei dienstgebundenen Paketen | Die Prozent-Schwelle **nähert** die Eigenschaft nur an — `postgresack`s Rücknahme allein bliebe grün | Eine Eigenschaft („braucht externen Dienst") ist ein Urteil, keine mechanisch prüfbare Größe; die Schwelle ist ein Proxy |
| `slice-096` | `AGENTS.md` §3.13 (Träger-Nachzug) selbst | (1) Suchform fängt Symbolnamen, keine Zahlen-Lokatoren; (2) Suchergebnis hat keinen committeten Träger | Die Regel benennt ihre eigene Suchform nicht, und sie benennt keinen Ablageort für ihr Ergebnis |

### 1.2 Verdikt: **keine gemeinsame Verkörperung über alle drei — die Domänen sind zu verschieden.** Konkret verkörpert wird nur die self-referentielle dritte Lücke (§3.13); die ersten beiden bleiben bewusst benannt, kein neuer Sensor.

**Gegen eine gemeinsame Prüf-Disziplin.** Eine allgemeine Regel „bei jeder
Hard-Rule-Formulierung die Sensor-Deckungslücke explizit benennen" würde
keinen der drei Fälle *zusätzlich* schließen — für `slice-078` und
`slice-079` steht die Lücke bereits explizit in `AGENTS.md` §3.11 bzw.
`harness/sensors/coverage-gate.md` §Grenze; die Praxis existiert längst,
ohne dass eine Meta-Regel sie erzwungen hätte. Eine formale Pflicht würde
das nur nachträglich zertifizieren, was schon geschieht — dieselbe
Begründungsstruktur, mit der `AGENTS.md` §3.12 §Die benannte Grenze eine
Formpflicht-auf-Prosa ablehnt („eine Formpflicht auf Prosa erzeugte
Pflichterfüllung"). Drei derart unterschiedliche technische Ursachen
(Markdown-Werkzeug-Design mit einem strukturellen Illustrations-Bedürfnis,
eine inhärent approximative Prozent-Metrik, eine Suchform mit zwei
unabhängigen Lücken) teilen keine gemeinsame Abhilfe — nur denselben
Oberbegriff. Eine gemeinsame Verkörperung wäre eine Etikettierung ohne
neuen Effekt.

**Für `slice-078` (hostpaths/Fences): kein Grenz-Hälften-Sensor jetzt.**
`state.md` benennt diese Entscheidung ausdrücklich als bei „Auftraggeber"
offen — dieser Lese-Schritt ist der Ort, an dem sie fällt: **nein**, aktuell
nicht commissioniert. Ein `make commit-traceability`-artiger
Grenz-Hälften-Sensor (positive Hälfte über das bestehende Modul, negative
Hälfte über einen neuen Fenced-Block-Scan) müsste zusätzlich eine
Opt-out-Markierung einführen — sonst würde er sein eigenes
Regel-Dokument (`AGENTS.md` §3.11, das ein `**Falsch:**`-Beispiel in einem
Fence führt, wenn auch mit Platzhalter `<Host-Wurzel>`) und jede künftige
`**Falsch:**`-Illustration eines realen Pfads rot färben. Das ist ein
eigenständiger Konventions-Entwurf (Opt-out-Marker-Semantik, analog
`d-check:ignore`), kein Ein-Zeilen-Fix, und der reale Befund-Stand
rechtfertigt ihn nicht: kein Vorgang hat bislang einen echten
Host-Pfad-in-Fence-Verstoß gemeldet — nur die strukturelle Lücke selbst
(dreimal). Der Wächter bleibt das Review, wie in `AGENTS.md` §3.11 bereits
benannt.

**Für `slice-079` (coverage-gate/Paketausschluss): kein neuer Sensor
jetzt.** `harness/sensors/coverage-gate.md` §Grenze Punkt 4/5 quantifiziert
das Residualrisiko bereits präzise (je Paket der reale Abstand zur
Schwelle, in Prozentpunkten und Statements) — reifer als eine neue
Struktur-Prüfung liefern könnte. Ein Sensor, der die deklarierte
Ausschlussliste (`ADR-0071`) gegen die tatsächliche `-coverpkg`-Aufrufform
hält, wäre technisch baubar und würde Grenzpunkt 5 („Disziplin, kein
Sensor") schließen — er schlösse aber **nicht** Grenzpunkt 4 (dass die
Prozent-Schwelle eine Rücknahme nicht zuverlässig färbt), weil das
Kern-Problem eine Eigenschaft ist, kein Listen-Abgleich. Ein Sensor, der
nur die Hälfte einer bereits benannten Lücke schließt, bei bereits
exzellenter Dokumentation der verbleibenden Hälfte, rechtfertigt die
zusätzliche Maschinerie nicht. **Verdikt:** bleibt benannt, kein Sensor.

**Für `slice-096` (§3.13 selbst): konkrete Textänderung — Verkörperung.**
Anders als die ersten beiden trifft dieser Beleg eine Hard Rule, die noch
gar keine eigene §Grenze-Notiz trägt (§3.11 hat sie in sich selbst,
`ADR-0071` in `coverage-gate.md`) — hier ist die Lücke nicht nur
unmechanisierbar, sie ist bislang **nirgends benannt**. Zwei Fixes:

1. **Gap 1 (Suchform).** `AGENTS.md` §3.13 nennt ihre Suchform nicht — sie
   sagt „`grep` über die Träger nach der bewegten Eigenschaft", ohne zu
   sagen, dass das nur **Symbolnamen** zuverlässig trifft und **Zahlen**
   (Zeilen-Lokatoren, Abschnittsnummern) durchfallen. Zielort: ein neuer
   Absatz in `AGENTS.md` §3.13, unmittelbar nach dem Absatz „Warum diese
   Regel einen Träger braucht und keinen Vorsatz" — er benennt die Grenze
   explizit (Vorbild: `harness/sensors/coverage-gate.md` §Grenze), nennt
   die Reviewer-Rolle als Schließer des Rests (real bereits so gelaufen,
   siehe `docs/reviews/review-slice-096.md` F-2, vierter Anker), und
   verlangt **keine** Erweiterung der Suchform selbst — eine `grep`-Form,
   die zuverlässig „dieselbe Klasse" auf Zahlen anwendet, bräuchte eine
   Semantik-Entscheidung, welche Zahl zu welcher Eigenschaft gehört
   (dieselbe Art Sensor-Unmöglichkeit wie bei `slice-079`).
2. **Gap 2 (kein committeter Träger) — dies braucht einen committeten
   Träger, ja.** §3.13 sagt „sein Ergebnis steht im Bericht des Slice" —
   ein Implementer-Bericht ist ein Handoff, über Läufe hinweg nicht
   nachlesbar (bestätigt: kein Träger in `slice-096`s Closure-Notiz nennt
   den rohen Suchlauf-Befund, nur die *nachträglich gefundene* Lücke
   selbst, `docs/plan/planning/done/slice-096-konfigurationsdatei-nachzug.md`
   §7 Punkt „Steering-Loop-Eintrag"). Zielort: `AGENTS.md` §3.13, Absatz
   „Wer sie liest" — ergänzt um die Pflicht, dass das Suchergebnis
   (Gefundenes **und** Nichtgefundenes) in einem strukturell vorgesehenen,
   committeten Feld des Slice-Plans landet, nicht nur im Lauf-Bericht:
   naheliegend eine DoD-Zeile analog dem bestehenden Muster
   (`docs/plan/planning/*.template.md` §2 Definition of Done), die den
   Implementer zwingt, den Suchlauf und sein Ergebnis **im Plan selbst**
   zu protokollieren, bevor die Closure-Notiz (§7) geschrieben wird — die
   Closure-Notiz allein reicht nicht, weil sie ex post und narrativ ist
   und den rohen Befund nur dann trägt, wenn ihn jemand später noch für
   erwähnenswert hält.

**Ausgang für diesen Eintrag: verkörpert** (Zielort `AGENTS.md` §3.13,
Herkunfts-Anker `seit welle-d-check`) — mit explizit benanntem Rest: die
Grenzen aus `slice-078`/`slice-079` bleiben unverändert benannt, kein
Sensor wird commissioniert; das ist eine getroffene Entscheidung, keine
offene Frage mehr.

---

## Eintrag 2 — `BEO-PGC/zitat-nennt-die-falsche-stelle`

### 2.1 Eigene Prüfung der drei Belege und ihrer Klassifikation

Gelesen: `observation.md`, `state.md`, alle drei `evidence/*.md`,
`docs/reviews/review-slice-090-delta.md` D-1, `docs/reviews/review-slice-102.md`
F-3, `docs/reviews/review-slice-d-check-tracked-modul.md` F-2,
`AGENTS.md` §3.12, `.harness/skills/reviewer.md` (voll gelesen).

Die drei Belege bestätigen sich inhaltlich. Auffällig — und für die
Verdikt-Frage entscheidend — ist, **wie unterschiedlich sie klassifiziert
wurden**, obwohl es strukturell derselbe Fehler ist (ein Verweis auf eine
Stelle eines anderen Dokuments trägt die Aussage nicht):

| Beleg | Gefunden von | `kategorie` | `quelle`/`klasse` |
|---|---|---|---|
| `slice-090` (D-1) | Delta-Review (**nicht** die erste Review-/Verify-Runde) | MEDIUM | `AGENTS.md` §3.12 Instanz A/B |
| `slice-102` (F-3) | Reviewer, unabhängig vom Verifier reproduziert | INFO | „Maintainability (`AGENTS.md` §3.13...)" |
| `slice-d-check-tracked-modul` (F-2) | Reviewer | **HIGH** | „Beleg trägt seinen Satz nicht" (bestehender Punkt) |

Das bestätigt die im Auftrag formulierte These direkt: Der **existierende**
HIGH-Punkt „Beleg trägt seinen Satz nicht" (`.harness/skills/reviewer.md`
Zeilen 114–133) wurde nur beim **dritten** Vorkommen tatsächlich als
einschlägig erkannt und angewendet. Beim ersten Vorkommen hat ihn **keine**
der beiden reguläten Prüf-Rollen (Review erste Runde, Verifier) überhaupt
gefunden — nur eine spätere, aus anderem Anlass angesetzte Delta-Review;
beim zweiten Vorkommen wurde er zwar gefunden, aber als INFO unter
`AGENTS.md` §3.13 statt als HIGH unter der einschlägigen Klasse eingeordnet.

### 2.2 Eigene Prüfung: reicht der bestehende HIGH-Punkt aus?

Der bestehende Punkt zählt als Beleg-Formen „einen Befehl, eine Abfrage,
einen Pfad, eine Adresse, eine Mutationsangabe oder eine Assertion" —
**keine dieser sechs Formen ist wörtlich ein Verweis auf eine Stelle eines
anderen Dokuments** (Abschnittsnummer, ADR-Festlegung, Slice-/Welle-
Kennung als Zitat-Ziel). Die Beispiel-Herkunft des Punkts
(„`BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht`", fünf Fundstellen: `git
diff` ohne Pathspec, `go list` ohne Testdatei-Feld, ein Testkommentar, eine
Adresse auf ein nicht existierendes Artefakt, eine Assertion, die auch aus
einem anderen Pfad hält) ist durchweg **ausführbar/mechanisch nachfahrbar**
(„den genannten Befehl ausführen, die genannte Adresse auflösen"). Ein
Abschnitts- oder Festlegungs-Zitat ist dagegen **nicht ausführbar** — es
wird **aufgeschlagen**, nicht gefahren. Diese Verweis-Form ist strukturell
eine Nachbar-Form, aber die Wortwahl des Punkts deckt sie nicht **explizit**
— und die Klassifikations-Drift über die drei realen Belege zeigt, dass
diese Uneindeutigkeit real zu unterschiedlichen Kategorien und
unterschiedlichem Schweregrad geführt hat (MEDIUM/INFO/HIGH für dieselbe
Fehlerklasse ist selbst ein Konsistenz-Problem der Skill-Anwendung).

`AGENTS.md` §3.12 Instanz B (Tatsachenbehauptung trägt ihren Beleg-Anker)
trägt das Prinzip auf der **Schreiber**-Seite bereits korrekt — eine
Aussage, die einen Beleg-Anker nennt, ist damit prüfbar. Was fehlt, ist die
**Leser**-seitige Pflicht beim Reviewer: den genannten Anker tatsächlich
aufzuschlagen, statt eine im Kopf mitgeführte oder aus einer Zusammenfassung
(Berichts-Kopfzeile) übernommene Vorstellung seines Inhalts ungeprüft zu
akzeptieren. Das ist exakt der in `observation.md` benannte Mechanismus
(„Ein solcher Kopf sieht aus wie eine Quelle und ist eine Wiedergabe").

### 2.3 Verdikt: **verkörpert** — der bestehende HIGH-Punkt wird um die Verweis-Form und eine explizite Gegenprobe-Pflicht ergänzt, keine neue Kategorie

Kein neuer, eigenständiger HIGH-Punkt — das würde die Klasse künstlich
spalten, wo sie strukturell dieselbe ist (der genannte Beleg trägt die
Aussage nicht; nur *wie* man ihn prüft, unterscheidet sich: ausführen vs.
aufschlagen). Stattdessen eine **präzisierende Erweiterung** des
bestehenden Punkts „Beleg trägt seinen Satz nicht" in
`.harness/skills/reviewer.md` (Zeilen 114–133):

1. Die Aufzählung der Beleg-Formen („einen Befehl, eine Abfrage, einen
   Pfad, eine Adresse, eine Mutationsangabe oder eine Assertion") wird um
   „einen Verweis auf eine Stelle eines anderen Dokuments — eine
   Abschnittsnummer, eine ADR-Festlegung, eine Slice-/Welle-Kennung"
   ergänzt.
2. Der Imperativ-Satz („Wer einen Beleg nennt, fährt ihn: den genannten
   Befehl ausführen, …") bekommt einen zweiten Halbsatz für die
   nicht-ausführbare Form: „— oder schlägt die genannte Stelle im
   **Original** auf, nicht in einer Zusammenfassung (die Kopfzeile eines
   Berichts, die eigene Erinnerung)." Das benennt den in `observation.md`
   beschriebenen Mechanismus direkt in der Handlungsanweisung, nicht nur
   in der Begründung.
3. Die Herkunfts-Zeile wird um `BEO-PGC/zitat-nennt-die-falsche-stelle`
   (3×, `slice-090`/`slice-102`/`slice-d-check-tracked-modul`) ergänzt,
   analog zur bestehenden `beleg-befehl-traegt-seinen-satz-nicht`-Zeile,
   mit Anker `· seit welle-d-check`.

**Warum keine neue Kategorie und kein neuer Sensor:** Ein Sensor ist
strukturell ausgeschlossen — ob eine zitierte Stelle die Aussage trägt, ist
eine Lese-Handlung am Original, dieselbe Grenze, die `AGENTS.md` §3.12
§Die benannte Grenze für Instanz A/B bereits zieht. Der einzig verfügbare
Wächter ist der Reviewer als Leser; die Erweiterung macht seine bereits
einmal (beim dritten Vorkommen) richtig geübte Praxis zur **expliziten**
Regel, statt sie dem Zufall zu überlassen, ob der Prüfende die
Nachbar-Form von sich aus unter denselben Punkt fasst.

**Ausgang für diesen Eintrag: verkörpert** (Zielort
`.harness/skills/reviewer.md`, bestehender HIGH-Punkt „Beleg trägt seinen
Satz nicht", Herkunfts-Anker `seit welle-d-check`).

---

## 3. Folgearbeit — nicht Teil dieses Zugs (Implementer-Fixrunde)

Dieses Verdikt ändert keinen lebenden Träger. Zwei Textänderungen sind
nachzuziehen:

- **`AGENTS.md` §3.13** — neuer Grenz-Absatz (Symbolname-vs-Zahl-Lücke,
  Reviewer als Schließer) **und** eine Trägerpflicht für das
  Suchlauf-Ergebnis (committeter Ort im Slice-Plan statt reinem
  Lauf-Bericht) — siehe §1.2 Punkt 1/2 oben für den genauen Zielort und
  Wortlaut-Vorschlag.
- **`.harness/skills/reviewer.md`** — Erweiterung des HIGH-Punkts „Beleg
  trägt seinen Satz nicht" um die Verweis-Form, die Gegenprobe-Pflicht
  und die Herkunfts-Zeile — siehe §2.3 oben.

Beide Änderungen sind reine `AGENTS.md`-/Skill-Textpflege, keine
`Accepted`-ADR-Korrektur (§3.5 unberührt) und keine neue ADR (§3.6
unberührt — keine Schwellen-Senkung).

## 4. Was nicht getan wurde

Keine Änderung an `AGENTS.md`, `.harness/skills/reviewer.md`,
`harness/sensors/*.md`, `.d-check.yml` oder einer ADR. Keine neue ADR
geschrieben. Kein `state.md`/`observation.md` im Beobachtungs-Register
verändert — die Ausgangs-Zuweisung im Register selbst ist der nachfolgende
Planner-Zug (Sequenz Planner → Architect → Planner, Vorbild
`BEO-PGC/commit-traceability-kein-vorab-hook`s Lese-Schritt bei
`welle-16`). Kein Produktionscode berührt.

## 5. Gates

`make gates` am Stand nach dem Commit dieses Verdikts: Exit-Code direkt und
ungepiped geprüft (`AGENTS.md` §3.9) — siehe Commit-Historie für den
Lauf-Beleg dieses Zugs.
