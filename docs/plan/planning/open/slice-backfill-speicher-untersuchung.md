# Slice backfill-speicher-untersuchung: Speicher-Untersuchung des Backfills — Spitze des Feed-Containers über dem Blockbedarf, Ursache und gemessene Grenze

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle — der Slice trägt keine Closure-Bedingung, die von
seiner DoD verschieden wäre. Er hat keine Kante zu einer offenen Welle und ist
von den Slices der Welle [welle-transformationen](../welle-transformationen.md)
unabhängig.

**Bezug:** [`LH-FA-CAP-009`](../../../../spec/lastenheft.md) (Backfill des
Bestands),
[`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md)
(Bestand-Backfill; sein Re-Evaluierungs-Trigger zur Ausbaustufe),
[`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md)
(Warnkriterium mit der Richtgröße von 4.000.000 geschätzten Zeilen),
[`ADR-0054`](../../adr/0054-coverage-gate-und-benchmark-infrastruktur.md) §(b)
(Benchmark-Infrastruktur), Architect-Verdikt
[`architect-verdict-welle-backfill-bestand-lese-schritt`](../../../reviews/architect-verdict-welle-backfill-bestand-lese-schritt.md)
§5 (h).

**Berührte Spec-Stellen:** [`SPEC-005`](../../../../spec/pflichtenheft.md)
(TransactionBufferPort; der Satz in `spec/pflichtenheft.md`, Zeile 72, gelesen
am 2026-09-25: „Große Transaktionen dürfen nicht unbegrenzt im RAM gehalten
werden“) — gelesen, nicht geändert.

**Verantwortlich:** — (noch nicht priorisiert).

**Autor:** Planner-Agent, Closure der Welle
[welle-backfill-bestand](../done/welle-backfill-bestand.md). **Datum:** 2026-09-25.

---

## 1. Ziel und Abgrenzung

**Ziel:** Die Ursache der Speicher-Spitze des Feed-Containers im Backfill-Run
ist benannt, und ihre Abhängigkeit von der Tabellengröße ist gemessen. Das
Benutzerhandbuch nennt unter „Grenzwerte“ (Speicher des Feed-Containers) den
Stand *übernommen* aus einem Lauf, der im Repository nicht auflösbar ist:
Stufe mit 200.000 Zeilen, Spitze im Run 401,7 bis 467 MiB bei 176,3 bis 196,0
MiB in Ruhe davor, höchster gemessener Wert 641,7 MiB; zwei Runs über je
1.000.000 Zeilen mit 435 MiB und **1.544 MiB**; der Bedarf eines Blocks liegt
bei etwa 74 KB, „die Ursache ist nicht untersucht, ein Zusammenhang mit der
Tabellengröße ist nicht belegt“ (Handbuch §9, Zeilen 1653 bis 1668, übernommen).
Die Warn-Richtgröße steht bei 4.000.000 geschätzten Zeilen
(`internal/application/usecase/backfill/warn.go`, gemessen am 2026-09-25); eine
Richtgröße, die auf einer ungeklärten Speicherlage steht, ist ein Risiko für den
ersten Betreiber.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Eine Änderung am Backfill-Code** (Bytelimit, Zeilenbreiten-Wache, andere
  Blockgröße). Der Slice ist eine Untersuchung; ergibt sie ein Wachstum mit der
  Tabellengröße, folgen ein Befund und ein eigener Änderungs-Slice mit
  Entscheidung des Architects ([`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md)
  Re-Evaluierung: Ausbaustufe für Durchsatz).
- **Ein Gate oder eine Pass/Fail-Schwelle für den Speicher.** Die Messung ist
  ohne Schwelle geführt (`harness/targets/bench-backfill.md`); eine Schwelle
  brauchte eine ADR ([`AGENTS.md`](../../../../AGENTS.md) §3.6).
- **Die Kopierdauer und der Durchsatz** — Gegenstand der bestehenden Messung;
  der Slice liest sie mit, ändert sie nicht.
- **Breite Zeilen als eigene Messreihe**, sofern die Ursache nicht an der
  Zeilenbreite hängt: die Zeilenbreite (etwa 74 Bytes) ist im Ursprung jeder
  Kopier-Zahl genannt, breitere Zeilen sind ungemessen und bleiben so, bis die
  Ursache einen Zusammenhang zeigt.

## 2. Definition of Done

- [ ] Die Messreihe steht mit gedruckten Zeilen im Repository: `tools/bench-backfill.sh`
      (vorhanden, Vertrag `harness/targets/bench-backfill.md`) liefert je Stufe
      den Speicher des Feed-Containers **und** eine Grundlinie ohne Run
      (Ruhe-Speicher vor und nach), über mindestens die Stufen bis 1.000.000
      Zeilen (`--full`); jede Zahl trägt Host, Lauf und die gedruckte Zeile
      (`AGENTS.md` §3.12 Instanz A), der Bericht trennt Messung von Übernahme.
      *Zu belegen durch:* die gedruckten Zeilen im Bericht des Slice unter
      `docs/reviews/` (auflösbar, anders als der übernommene Lauf
      `20260924T233628Z`); der Lauf, aus dem die 1.544 MiB stammen, wird
      wiederholt oder bleibt als *übernommen* gekennzeichnet.
- [ ] Die Ursache ist benannt und belegt: ein Lauf, der sie an- und abschaltet
      (eine Änderung genau einer Größe — Blockgröße, Zeilenbreite,
      Garbage-Collector-Einstellung des Feed-Containers oder der Zeitpunkt der
      Probe —, dieselbe Messung vor und nach), oder die Ursache steht als
      *ungeklärt* mit der gemessenen Grenze im Bericht. Der Blockbedarf von etwa
      74 KB (Handbuch, übernommen) erklärt Spitzen von einigen hundert MiB nicht;
      welche Größe sie trägt, ist offen und im Slice zu erproben.
- [ ] Der Ausgang ist gesetzt: wächst der Speicher mit der Tabellengröße, steht
      ein Befund im Bericht und ein Änderungs-Slice existiert als Datei in `open/`
      (Entscheidung des Architects, ob eine ADR nötig ist); sonst trägt das
      Benutzerhandbuch unter „Grenzwerte“ die gemessene Grenze samt Ursprung,
      und die Warn-Richtgröße von 4.000.000 geschätzten Zeilen ist gegen die
      Messung bewertet (bleibt, oder eine Nachschärfung nach dem Trigger von
      [`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md)
      „Eine Messung liegt vor“ ist beantragt); die Handbuch-Version und die
      Änderungshistorie tragen eine Zeile. *Zu belegen durch:* Lesen des
      Handbuch-Abschnitts und `make docs-check`.
- [ ] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8).
- [ ] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff)
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [ ] Doku-Update: siehe dritter Liefer-Punkt (Handbuch §Grenzwerte,
      `harness/targets/bench-backfill.md`, falls der Vertrag der Messung eine
      Grundlinie führt).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag (geschärfte Regel · neuer
      Sensor · benannte Spec-Lücke).
- [ ] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem
      Repo (Greenfield).
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder weitere `evidence/`-Datei; kein Anfall ist ebenfalls
      eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter
      offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von
      der Closure der nächsten Welle (die Roadmap führt
      [welle-transformationen](../welle-transformationen.md) unter *Offene
      Wellen*, das Ereignis kann eintreten; ein Slice ohne Welle wird von ihr
      mitgeprüft).

**Umfang:** M — Schätzung, nicht gemessen: die Messreihe fährt die 1.000.000-Stufe
(Laufzeit und Speicherbedarf des Messhosts sind ungemessen), die Ursachensuche
hat kein bekanntes Ende.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/bench-backfill.sh` (falls die Grundlinie fehlt) | update | Speicher ohne Run vor und nach der Stufe, im Ursprung jeder Zahl genannt. |
| `harness/targets/bench-backfill.md` | update | Vertrag der Messung nennt die Grundlinie. |
| `docs/user/benutzerhandbuch.md` (§Grenzwerte, Version, Änderungshistorie) | update | gemessene Grenze samt Ursprung statt der übernommenen Zahlen. |
| Bericht des Slice unter `docs/reviews/` | neu | gedruckte Zeilen der Messreihe, Ursache, Ausgang. |

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „die Aussage über den
Speicher des Feed-Containers im Backfill und die Warn-Richtgröße“; beide Stände
gemessen; die Befehle stehen im Codeblock, der Implementer trägt Stand und
Trefferzahl ein):**

```text
git grep -n -i -E 'Speicher|MiB|1\.544|4\.000\.000|Richtgröße|nicht untersucht' -- docs/user harness spec internal tools
git grep -n -E 'estimatedRowsGuideline|DefaultBlockSize' -- internal
```

| Träger | Befund | Behandlung |
|---|---|---|
| Handbuch §Grenzwerte, `harness/targets/bench-backfill.md`, Doc-Kommentare an `DefaultBlockSize` und `estimatedRowsGuideline` | *(Implementer trägt ein)* | jede Zahl trägt Lauf und Ursprung; Aussagen zur Ursache folgen dem Ausgang |

## 4. Trigger

**Start** (`next` → `in-progress`): kein weiterer Slice in `in-progress/`
(WIP-Limit 1); sonst unabhängig von jedem anderen Slice. Der Slice muss `done`
sein, **bevor** ein Server-Release veröffentlicht wird, dessen Commit den
Backfill trägt — beobachtbar: `git tag -l 'v*'` nennt ein `v*`-Tag über dem
höchsten Tag `v0.1.2`, dessen Commit die Backfill-Änderungen enthält (gemessen am
2026-09-25: kein Server-Tag trägt sie). Ein mechanischer Wächter für diese
Bedingung existiert nicht; der Träger ist die Roadmap-Zeile dieses Slice
(Benannte Grenze, §6).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): falls Messreihe und
  Ursachensuche nicht in einem Review tragen — der abtrennbare Teil ist die
  Messreihe (erster Liefer-Punkt) als eigener Slice, die Ursachensuche folgt.
- `in-progress` → `open` (blockiert): falls die Messreihe an der Kapazität des
  Messhosts scheitert (1.000.000 Zeilen brauchen Speicher und Zeit; kein
  Ersatz-Host verfügbar) — dann steht die Grenze der Messung im Bericht statt
  einer Ursache.

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + die Messreihe mit gedruckten Zeilen im
Bericht + Ausgang gesetzt (Befund und Änderungs-Slice als Datei in `open/`,
oder Handbuch-Grenze) + Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- **Die Messung hängt am Host** (Speicher und Laufzeit des Messhosts). *Erwartet,
  zu belegen durch:* jede
  Zahl nennt Host und Lauf; das DoD-Kriterium „Ursache benannt“ ist hostunabhängig
  gefasst (`BEO-PGC/dod-kriterium-haengt-am-messhost`, 1×, offen).
  **Ausgang:** *(bei Closure)*
- **Die Ursache bleibt ungeklärt.** *Erwartet, zu belegen durch:* der Ausgang
  „ungeklärt mit gemessener Grenze“ ist zulässig und steht im Bericht und im
  Handbuch; er ist kein stiller Abschluss. **Ausgang:** *(bei Closure)*
- **Der Slice läuft nicht vor dem ersten Server-Release mit Backfill** (kein
  Wächter). *Erwartet, zu belegen durch:* die Roadmap-Zeile nennt die
  Bedingung; die Prüfung liegt beim Planner der Release-Vorbereitung. **Ausgang:**
  *(bei Closure)*
- **Die Warn-Richtgröße bleibt unbegründet, weil kein Wachstum gemessen wird**
  (`BEO-PGC/blockgroesse-zaehlt-zeilen-nicht-bytes`, 1×, Ausgang: dieser Slice).
  *Erwartet, zu belegen durch:* die Bewertung im dritten Liefer-Punkt.
  **Ausgang:** *(bei Closure)*
- **Speicher des Messhosts während der Messreihe.** Die Stufe mit 1.000.000
  Zeilen belastet Docker-Volumes und den Feed-Container; kein `docker volume
  prune` und kein `docker system prune` im Lauf. *Erwartet, zu belegen durch:*
  der Bericht nennt die Aufräum-Schritte des Runners. **Ausgang:** *(bei
  Closure)*

## 7. Closure-Notiz

- **Was hat funktioniert:** *(zu tragen bei Closure)*
- **Was ging anders als geplant:** *(zu tragen bei Closure)*
- **Steering-Loop-Eintrag (Lerneintrag):** *(zu tragen bei Closure —
  geschärfte Regel · neuer Sensor · benannte Spec-Lücke; ohne ihn kein
  `done/`-Übergang)*
- **Beobachtungs-Register (`../observations/`):** *(je Anfall Beleg oder
  „keine Beobachtung angefallen“ als notierte Antwort)*
- **Folge-Slices:** *(zu tragen bei Closure)*
- **Risiken aus §6:** *(je ein Ausgang)*
- **Drei Paarungen:** dieser Slice hat keine Welle; die Prüfung läuft
  regelkonform bei der Closure der nächsten Welle.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area
`*`/`PGC` (Greenfield); Bench-Skripte, Handbuch und Use Case des Runs sind
keine eigenen Sub-Areas — kein Anlass zur Ausdifferenzierung.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen
(Zähler gemessen am 2026-09-25 mit `ls evidence | wc -l` je Eintrag) —
`BEO-PGC/blockgroesse-zaehlt-zeilen-nicht-bytes` (1×, Ausgang: dieser Slice),
`BEO-PGC/backfill-adapter-startwerte-ohne-messung` (offen, 1×, verwandt: dieselbe
Messung, andere Größe), `BEO-PGC/dod-kriterium-haengt-am-messhost` (offen, 1×,
Risiko §6), `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (verkörpert, 21×,
die übernommenen Zahlen des Handbuchs).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
