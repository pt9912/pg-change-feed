# Slice bench-source-impact-absolut: das PER-001-Bench prüft die Zusatzlatenz je Commit gegen die gemessene Festschreib-Latenz

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

**Bezug:** [`LH-QA-PER-001`](../../../../spec/lastenheft.md),
[`ADR-0155`](../../adr/0155-per-001-zusatzlatenz-relativ-zur-festschreib-latenz.md)
(Entscheidung; ersetzt `ADR-0151` und die PER-001-Schwelle von `ADR-0104`),
[`ADR-0054`](../../adr/0054-coverage-gate-und-benchmark-infrastruktur.md)
(Bench-Infrastruktur, unverändert).

**Berührte Spec-Stellen:** `SPEC-025` · `SPEC-036` (beide bereits nachgezogen,
Voraussetzung des Starts, siehe §4); der Slice ändert die Spec nicht.

**Verantwortlich:** — bis zur Priorisierung.

**Autor:** pt9912 (Architect-Zuschnitt). **Datum:** 2026-10-05.

---

## 1. Ziel und Abgrenzung

**Ziel:** `tools/bench-source-impact.sh` misst im selben Lauf die
Festschreib-Latenz `t_sync` und endet rot, wenn die zusätzliche Latenz je
Quelltransaktion Δ die Grenze max(0,10 ms; 1,5 × `t_sync`) überschreitet —
auf jedem Host, ohne die alte 35-%-Schwelle.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Die Latenz selbst senken (Commits bündeln, Store trennen) — berührt
  `ADR-0011`/`ADR-0010`, eigene Entscheidung (`ADR-0155` §Nicht Teil).
- `LH-QA-PER-002`/`LH-QA-PER-003` und deren Skripte — ihre Schwellen bleiben
  nach `ADR-0104`.
- Eine Aufnahme von `make bench` in `make gates` — das Ziel braucht Docker-Netz,
  Image und DB-Zugang (Bestand, `harness/README.md` §Sensors).
- Eine Messung auf einem zweiten Host zur Absicherung des Faktors 1,5 — das ist
  der Re-Evaluierungs-Trigger von `ADR-0155`, kein Liefer-Punkt: ein
  Einzellauf auf dem Entwicklungshost genügt der DoD.

## 2. Definition of Done

- [x] `tools/bench-source-impact.sh` misst `t_sync` (siehe §3), druckt in jedem
      Lauf Δ, `t_sync`, das Verhältnis Δ / `t_sync` und die Grenze, endet bei
      Δ über der Grenze mit Exit 1 und bei nicht messbarem `t_sync` mit Exit 2;
      `THRESHOLD_PCT` und der Prozent-Vergleich sind entfernt. `LH-QA-PER-001`
      erfüllt, Messung dokumentiert: gedruckte Zeile eines `make bench`-Laufs
      auf dem Entwicklungshost im Bericht (erwartet: Exit 0 bei Δ ≈ 3 ms und
      `t_sync` ≈ 2,9 ms — **zu belegen durch den Lauf**, nicht vorab behauptet).
- [x] Mutationsprobe an einer **Scratchpad-Kopie** des Skripts (`AGENTS.md` §3.1,
      nie `sed -i`): Faktor 1,5 → 0,5 macht den Lauf auf dem Entwicklungshost rot
      (Exit 1); `pg_test_fsync` durch einen nicht vorhandenen Befehl ersetzt gibt
      Exit 2. Der Bericht nennt Stellen, Instanz und gesehene Farbe (§3.12);
      nicht Gefahrenes steht als hergeleitet.
- [x] Die Träger der bewegten Eigenschaft sind nachgezogen (Liste §3; Suchlauf im
      Bericht und in §6).
- [x] `make gates` grün (Exit 0, ungefiltert gelaufen, Exit separat gesichert).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`), kein Self-Review.
- [x] Doku-Update: `docs/user/bench-abdeckung.md` (vom Generator geschrieben),
      kein weiterer öffentlicher Vertrag berührt.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register fortgeschrieben (`BEO-PGC/dod-kriterium-haengt-am-messhost`
      um einen Beleg ergänzt, wenn der Lauf ihn liefert) oder „keine
      Beobachtung angefallen“ in §7.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

(Liefer-Punkte: zwei — Skript samt Probe, Träger-Nachzug.)

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/bench-source-impact.sh` | update | `t_sync` vor Phase 1 messen: `docker exec <pg> pg_test_fsync -s 2 -f <Datei im Datenverzeichnis>`, Zeile der Methode aus `SHOW wal_sync_method` (Regelfall `fdatasync`), Wert `usecs/op` in ms; Grenze und Vergleich per `awk` (Gleitkomma); Exit 1/2 wie `ADR-0155` Entscheidung 5/6; `bench::record_row`-Schwelle-Text auf „Zusatzlatenz je Commit ≤ max(0,10 ms; 1,5 × Festschreib-Latenz) (`SPEC-025`)“ — [`LH-QA-PER-001`](../../../../spec/lastenheft.md) |
| `tools/bench-lib.sh` | update | Kommentar Zeile 160 und Generator-Kopf (Zeile 190, `docs/user/bench-abdeckung.md`): Quelle der PER-001-Schwelle nennt `ADR-0155`, die Quellen der PER-002/003-Schwellen bleiben `ADR-0104`. **Kein Helfer** `bench::sync_latency_ms`: die Messfunktion `measure_t_sync` steht im Skript selbst (Reduktion, damit die Mutationsprobe an einer Scratchpad-Kopie des einen Skripts greift) |
| `Makefile` | update | Hilfetext des Ziels `bench` (Zeile 246) nennt die neue Messgröße statt „mit Schwelle“ allein, falls er die 35 % trägt |
| `docs/user/bench-abdeckung.md` | neu erzeugt | vom Lauf geschrieben (Zeile PER-001, Kopf); nicht von Hand ändern; die Datei ist Nutzerdokument — Kennungsfreiheit prüft `make handbuch-public-doc-check`? Nein: die Datei ist ausgenommen (Bestand), die Kennungen in der Generator-Zeile bleiben wie bisher |
| `harness/README.md` | update | §Sensors Zeile `make bench`: Beschreibung des PER-001-Skripts (35-%-Schwelle, 25,0 %/28,2 %, 87,5 % bis 95,8 %, „endet dort rot“) auf die neue Messgröße; zuvor an beiden Ständen nachmessen (§6); Messwerte des Verdikts bleiben als **übernommen** gekennzeichnet |
| `harness/sensors/*` / `harness/targets/*` | prüfen | Suchlauf zeigt: nur `harness/targets/bench-backfill.md` nennt Bench; falls es PER-001 beschreibt, nachziehen, sonst Negativbefund im Bericht. **Negativbefund:** `git grep -n -i -E 'PER-001\|source-impact\|35' harness/targets/bench-backfill.md harness/sensors` fand 0 Treffer, nichts nachzuziehen |
| `Makefile` (Ist) | update | Hilfetext des Ziels `bench` nennt PER-001 mit der neuen Messgröße (die 35 % standen dort nicht) |

Ansatz: `pg_test_fsync` liegt im PostgreSQL-Image (Verdikt: gemessen im
`postgres:18-alpine`-Container; für den **Digest-Pin** `PG_TEST_IMAGE` in
`tools/bench-lib.sh` *hergeleitet*, der Implementer prüft es mit `docker exec`
als ersten Schritt — fehlt es, ist das ein Befund an `ADR-0155`, kein stiller
Ausweg).

```suchlauf
71c6a886 5 -E 'THRESHOLD_PCT|35 ?-?%|SCHWELLE ÜBERSCHRITTEN' -- tools/bench-source-impact.sh
71c6a886 4 -E 'ADR-0104|ADR-0151' -- tools Makefile harness docs/user
71c6a886 2 -E '35[ -]?%' -- . :!docs/reviews :!docs/plan/planning/done :!docs/plan/planning/observations :!.harness/baseline :!docs/plan/adr
71c6a886 1 -i -E 'fdatasync|pg_test_fsync' -- tools Makefile harness docs/user
diff 0 -E 'THRESHOLD_PCT|35 ?-?%|SCHWELLE ÜBERSCHRITTEN' -- tools/bench-source-impact.sh
diff 5 -E 'ADR-0104|ADR-0151' -- tools Makefile harness docs/user
diff 1 -E '35[ -]?%' -- . :!docs/reviews :!docs/plan/planning/done :!docs/plan/planning/observations :!.harness/baseline :!docs/plan/adr
diff 6 -i -E 'fdatasync|pg_test_fsync' -- tools Makefile harness docs/user
```

Die `diff`-Zeilen sind der Stand nach dem Nachzug: das Skript trägt keine
Prozent-Schwelle mehr (0); die fünf `ADR-0104`-Zeilen sind Quellen der
PER-002/003-Schwellen bzw. benannte Herkunft (Generator-Kopf, README, drei
Kommentar-/Kopfzeilen der Bibliothek, davon zwei mit `ADR-0155` daneben); die
eine `35 %`-Zeile ist die als **übernommen** gekennzeichnete Herkunftsangabe des
früheren Verdikts in `harness/README.md`; die sechs `fdatasync`-/`pg_test_fsync`-Zeilen
sind das neue Verdikt (fünf im Skript, eine in `harness/README.md`).

**Messung (Implementer-Lauf, 2026-10-05, Host Linux 6.8.0-139-generic,
`postgres:18-alpine` aus dem Pin `PG_TEST_IMAGE`, `wal_sync_method` = `fdatasync`,
`pg_test_fsync` im Pin vorhanden — geprüft mit `docker run --rm --entrypoint sh <Pin> -c 'which pg_test_fsync'`,
Ausgabe `/usr/local/bin/pg_test_fsync`):** `bash tools/bench-source-impact.sh` direkt
(nicht über `make bench`; die übrigen Bench-Skripte sind nicht Gegenstand), gedruckte Zeilen:

```text
Lauf 1 (Exit 0):
bench-source-impact: t_sync 1.782 ms (pg_test_fsync, Methode aus wal_sync_method)
bench-source-impact: Ergebnis (LH-QA-PER-001) — ohne CDC 9480 ms (Median von 5 Läufen), mit CDC 18258 ms (Median von 5 Läufen), Differenz 8778 ms für je 5000 Schreibtransaktionen auf public.bench_source_impact
bench-source-impact: Zusatzlatenz Δ 1.756 ms je Transaktion, t_sync 1.782 ms, Verhältnis Δ/t_sync 0.985, Grenze max(0.10 ms; 1.5 × t_sync) = 2.673 ms
Lauf 2 (Exit 0):
bench-source-impact: Ergebnis (LH-QA-PER-001) — ohne CDC 10400 ms (Median von 5 Läufen), mit CDC 19446 ms (Median von 5 Läufen), Differenz 9046 ms für je 5000 Schreibtransaktionen auf public.bench_source_impact
bench-source-impact: Zusatzlatenz Δ 1.809 ms je Transaktion, t_sync 2.795 ms, Verhältnis Δ/t_sync 0.647, Grenze max(0.10 ms; 1.5 × t_sync) = 4.192 ms
```

Abweichung gegen die Erwartung der DoD (Δ ≈ 3 ms, `t_sync` ≈ 2,9 ms, übernommen aus
`ADR-0155`): auf diesem Lauf ist Δ ≈ 1,8 ms und `t_sync` 1,78 ms bzw. 2,80 ms
(gemessen); das Verhältnis liegt zwischen 0,65 und 0,99, also unter 1,5. `t_sync`
streut zwischen zwei Läufen desselben Hosts um den Faktor 1,6; die Grenze folgt ihm.

**Mutationsprobe** (Scratchpad-Kopien des Skripts, erzeugt mit `sed … > Kopie`, Aufruf
`bash <Kopie>`; je Zusage die mutierte Eingabe und die gesehene Farbe):

| Zusage | mutierte Eingabe | gesehenes Ergebnis |
|---|---|---|
| Δ über Grenze → Exit 1 | Zeile `FACTOR=1.5` → `FACTOR=0.5` (Grenze 0.882 ms bei Δ 1.837 ms, t_sync 1.765 ms) | Exit 1, Meldung `GRENZE ÜBERSCHRITTEN …` |
| `t_sync` nicht messbar → Exit 2, kein Rückfall | `pg_test_fsync` → `pg_test_fsync_fehlt` im `docker exec`-Aufruf | Exit 2, Meldung `t_sync nicht messbar …`, vor Phase 1 |
| Untergrenze 0,10 ms (`MIN_MS`) | nicht gefahren: greift nur bei `t_sync` < 0,067 ms; auf diesem Host nicht erreichbar — **hergeleitet** | — |
| Zeilenwahl der Methode (`SHOW wal_sync_method`) | nicht gefahren — **hergeleitet** (Lauf gegen `fdatasync`, die Zeile ist der erste Treffer des Abschnitts „one 8kB write“) | — |

Die Mutation von `FACTOR` mutiert die Grenze, nicht die Messung von Δ; die Eingabeseite
des Exit-1-Zweigs (Δ selbst, etwa durch einen künstlich verzögernden Feed) ist
**nicht gefahren** (hergeleitet).

Suchraum: Baum ohne `docs/reviews/**`, `done/`, `observations/` (Records),
`.harness/baseline/**`, ADRs (immutabel). Die Zahlen sind am Parent `71c6a886`
zu messen.
Stand der Zeilen: **gemessen** am Parent `71c6a886` mit
`bash tools/harness/suchlauf-nachmessen.sh <Plan-Datei>` (vier Zeilen stimmen, Exit 0);
die Vollständigkeit von Suchraum und Muster ist Lese-Handlung des Reviewers.

## 4. Trigger

**Start** (`next` → `in-progress`): `ADR-0155` ist `Accepted`, das Lastenheft
trägt die neu gefasste `LH-QA-PER-001` (0.16.0) und das Pflichtenheft
`SPEC-025`/`SPEC-036` in der Fassung von `ADR-0155` (zuvor committet), und der
Implementer hat ein geladenes `:dev`-Image (`make image`).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): Die Messung von `t_sync` verlangt mehr als
  einen `docker exec` und eine Auswertung (z. B. mehrere Datenträger) — dann
  Zerlegung in „Messung“ und „Träger“.
- `in-progress` → `open` (blockiert): `pg_test_fsync` fehlt im gepinnten Image —
  dann Folge-ADR zur Messmethode.

## 5. Closure-Trigger

DoD vollständig, `make gates` grün, Review- und Verifier-Bericht liegen vor,
Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- Das Skript hängt am Host: der Lauf dauert Minuten und braucht Docker-Netz; ein
  roter Lauf auf dem Entwicklungshost nach der Umstellung wäre ein Befund gegen
  den Faktor 1,5 — **Ausgang:** weiter offen bis zum Lauf; tritt er ein, Folge-ADR
  mit `Supersedes ADR-0155` (Re-Evaluierungs-Trigger dort).
- Der Suchlauf der Träger fängt Symbolnamen, nicht verschobene Zahlen
  (`AGENTS.md` §3.13 Grenze) — **Ausgang:** Lese-Handlung des Reviewers.
- Die Zahlen des Verdikts in `harness/README.md` (87,5 % bis 95,8 %, 94,6 %)
  bleiben nur als **übernommen** stehen — **Ausgang:** entfallen als Risiko, wenn
  die Zeile die Herkunft nennt (`AGENTS.md` §3.12).

## 7. Closure-Notiz

- **Was hat funktioniert:** —
- **Was ging anders als geplant:** —
- **Steering-Loop-Eintrag:** —
- **Beobachtungs-Register:** —
- **Folge-Slices:** —
- **Risiken aus §6:** —
- **Drei Paarungen:** —

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist `*` (Default, `tools/bench-*`,
`harness/`, `docs/user/`); eine Ausdifferenzierung ist nicht nötig.

**Vorgelagert — offene Beobachtungen sichten:** `BEO-PGC/dod-kriterium-haengt-am-messhost`
trifft den Gegenstand (Zählerstand beim Start aus den Dateien des Eintrags lesen).

Alle berührten Sub-Areas GF.
