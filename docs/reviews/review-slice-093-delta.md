# Review-Report: slice-093 — **Delta** (Fixrunde `4298c4c`) · 2026-09-16

**Review-Art:** Delta-Review **eines** Fix-Commits — geprüft gegen **Plan und
Entscheidungen**, gegen die Findings der ersten Runde (`docs/reviews/review-slice-093.md`
F-1…F-6) und gegen die Zusage des Implementers im Commit-Text. Baseline-Regelwerk
`v6.5.0` · `regelwerk/modul-10-review-harness.md` §Drei Review-Arten. DoD-/Spec-Konformität
(Verifier), §6-Risiko-Ausgänge, Register-Zähler und die drei Paarungen (Planner-Closure)
sind **nicht** Gegenstand. Arbeitsbaum vor und nach **jedem** Lauf sauber.

**Gegenstand** — `4298c4c` (Parent `66eb313`):

| Pfad | Zeilen | im Gegenstand |
|---|---|---|
| `internal/bootstrap/roles_rollout_file_internal_test.go` | **+97/−14** | **ja** |
| `internal/adapters/driven/telemetry/slog_levels_internal_test.go` | **+13/−4** | **ja** |
| `docs/plan/adr/0085-build-kontext-ausnahme-test-only-zweck.md` | +128/−0 | **nein** (paralleler Architect-Zug zu F-1) |

*Zahl-Korrektur zum Auftrag (gemessen):* „+110/−18" und „+17" sind die **Summen beider
Testdateien**, nicht zwei Datei-Maße; die Einzelwerte stehen oben.

**Skill:** `.harness/skills/reviewer.md` @ `HEAD` (`03fc53a`; `git diff --stat 4298c4c HEAD
-- internal/` → **leer**, die zwei Gegenstands-Dateien sind unverändert) · **Datum:** 2026-09-16.

---

## Findings

### D-1 — Der neue Satz beschreibt den **abwesenden** Text der Vorgänger-Fassung statt der geltenden Zusage

- `kategorie`: **HIGH**
- `quelle`: `AGENTS.md` §3.7 („**Falsch:** *die frühere Fassung prüfte nur die Länge* —
  beschreibt abwesenden Text. **Richtig:** die geltende Zusage nennen; die vorige hält
  `git`.") · `.harness/skills/reviewer.md` HIGH „Kommentar trägt keine der Kommentar-Klassen"
- `pfad`: `internal/bootstrap/roles_rollout_file_internal_test.go:176-180`, Fundstelle `:179-180`
- `befund`: Der Satz lautet:

  ```text
  176 // Zwei Prüfungen sind **Regeln über den geparsten Grant-Bestand**, keine
  177 // Namenslisten: (6) kein Rollen-Grant über eine ganze Objektklasse und
  178 // (7) der Leser teilt kein Objekt mit einer schreibenden Rolle. Die zweite
  179 // ersetzt eine frühere Acht-Namen-Liste von Basistabellen — sie war
  180 // handverlesen und ließ jedes neunte Objekt durch.
  ```

  Sein Subjekt ist in `:179-180` **der entfernte Text**: „ersetzt eine **frühere**
  Acht-Namen-Liste … **sie war** handverlesen und **ließ** jedes neunte Objekt durch". Die
  Liste ist im Baum nicht mehr vorhanden; die Aussage ist am geltenden Stand **nicht
  prüfbar**, nur in `git`. Die Zahl „Acht" ist eine Zahl über einen abwesenden Gegenstand.
  Repo-weite Probe: **die einzige Stelle dieser Form**:

  ```text
  $ grep -rnE "^[[:space:]]*//.*(frühere Fassung|vorige Fassung|früher stand|ersetzt eine frühere|bislang stand|die alte Fassung)" --include=*.go . | grep -v "^./.git"
  internal/bootstrap/roles_rollout_file_internal_test.go:179:// ersetzt eine frühere Acht-Namen-Liste von Basistabellen — sie war
  ```

  *Abgrenzung:* Man kann den Satz als **Abgrenzung** lesen („warum eine Regel und keine
  Liste") — eine legitime Kommentar-Klasse. Sie trägt aber nicht, weil die Begründung im
  **Vergangenheits-Ton über die eigene Datei** steht; derselbe Inhalt ist als Ist-Zustand
  formulierbar.
- `verifizierbar`: **ja** — `grep`-Probe; `git show 4298c4c`; die Liste selbst ist aus dem
  Baum nicht auflösbar
- `klasse`: „Absent-Text im Kommentar (Vorher/Nachher)" — **verwandt**, nicht dasselbe
  Vorkommen wie `BEO-PGC/slice-chronik-in-code-kommentar`

### D-2 — Die Ausnahme der neuen Regel (7) ist **weiter** als ihr Satz behauptet

- `kategorie`: **MEDIUM**
- `quelle`: LP2 („an das **reale Artefakt** gebunden") · §3.12 Instanz B · `LH-QA-SEC-001`
- `pfad`: `internal/bootstrap/roles_rollout_file_internal_test.go:270-278`
- `befund`: Der Satz schließt die Ausnahme **einelementig**: „Ausgenommen ist allein der
  Schema-USAGE-Grant — er ist die Vorbedingung beider Seiten, kein Objektzugriff." Der Code
  (`:278`) nimmt per `strings.HasPrefix(objekt, "schema ")` **jedes** Schema-Objekt aus.
  Gemessen kostet das zwei reale Regressionen, die weder (6) noch (7) fängt und die **keine**
  der drei benannten Grenzen nennt:

  ```text
  M7  GRANT CREATE ON SCHEMA cdc TO cdc_reader;   → TEST-EXIT=0   ok …internal/bootstrap
  M13 GRANT ALL ON SCHEMA cdc TO cdc_reader;      → TEST-EXIT=0   ok …internal/bootstrap
  ```

  `GRANT ALL ON SCHEMA cdc TO cdc_reader` gibt dem Leser `CREATE` auf dem Schema — eine
  Least-Privilege-Regression, die der Test **still** durchlässt, während der Satz daneben
  eine engere Ausnahme zusichert, als der Code zieht.
- `verifizierbar`: **ja** — zwei eigene Mutationen, beide **EC=0**
- `klasse`: „benannte Ausnahme enger als die Ausnahme im Code" — **derselbe Vorgang** wie F-2

### D-3 — Der Beleg-Adressat „Lauf-Bericht" löst nirgends auf

- `kategorie`: LOW
- `quelle`: §3.12 Instanz B · Instanz A
- `pfad`: `internal/bootstrap/roles_rollout_file_internal_test.go:182`
- `befund`: `// Rot färbende Mutationen (real gefahren, Exit-Codes im Lauf-Bericht):` —
  „Lauf-Bericht" ist **kein Artefakt dieses Repos** (repo-weit genau ein Treffer: die Zeile
  selbst). Der Leser kann die sechs Exit-Codes nicht nachschlagen. Dazu: die Auflistung nennt
  **sechs** Mutationen, der Commit-Text derselben Runde sagt „**Fuenf** Mutationen, alle Exit
  1". Der Reviewer hat **alle sechs** gefahren und **alle sechs rot** gesehen — die Auflistung
  trägt; die Zahl des Commit-Textes ist die unstimmige.
- `verifizierbar`: **ja** — `grep -rn "Lauf-Bericht" .` → 1 Treffer; sechs Mutationsläufe
- `klasse`: „Beleg-Adressat ohne auflösbares Artefakt"

### D-4 — Die Begründung von Regel (6) ist weiter als die Messung

- `kategorie`: LOW
- `quelle`: §3.12 Instanz B
- `pfad`: `internal/bootstrap/roles_rollout_file_internal_test.go:258-261`, Fundstelle `:260`
- `befund`: „`ALL TABLES IN SCHEMA`-Zeile ist keine Aufzählung, sondern eine Rundum-Vergabe —
  **sie hebt jeden geprüften Objekt-Grant auf**". Gemessen hebt sie **nichts** auf: mit
  angehängtem `GRANT ALL ON ALL TABLES IN SCHEMA cdc TO cdc_reader` bleiben alle Grants der
  Datei stehen, und (0)–(5) laufen weiter durch — der **Fatal kommt allein aus (6)**.
  Aufgehoben ist die **Trennung**, nicht ein Grant einer anderen Rolle.
- `verifizierbar`: **ja** — M4: EC=1, die Fatal-Zeile ist die von (6)
- `klasse`: „Regel-Begründung weiter als die Messung"

### D-5 — „seit `LH-FA-RET-005`" ist die Herkunfts-Aussage in einer Nicht-Anker-Form

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.7 (Herkunft als **ein** auflösbares Feld in den genannten Formen)
- `pfad`: `internal/bootstrap/roles_rollout_file_internal_test.go:247`
- `befund`: Die Wendung hängt „seit" an eine `LH-*`-ID; die kanonische Anker-Form ist
  `· seit slice-<NNN>`. Die Aussage ist **wahr** (`LH-FA-RET-005` = „Erkennbarkeit
  blockierender Consumer", und der Lesezugriffsweg läuft real über `cdc_reader`). Sie
  spiegelt außerdem wörtlich eine Zeile **außerhalb** des Diffs
  (`tools/schema/nacharbeit-roles.sql:101`). Kein Rückgabe-Pfeil.
- `verifizierbar`: **ja**
- `klasse`: „Herkunfts-Anker in Nicht-Anker-Form"

---

## Negativbefunde

**F-2 ist im Kern erledigt — die zwei Regeln binden, was sie behaupten.** 18 LP2-Läufe
(1 Baseline + 17 Mutationen): **11 rot / 6 grün**, Baseline grün.

| Probe | Mutation | EC |
|---|---|---|
| M1 | `GRANT USAGE ON SCHEMA cdc …` entfernt | **1** (Prüfung (0)) |
| M2 | `cdc.retention_blockers` aus dem Reader-Grant | **1** |
| M4 | `GRANT ALL ON ALL TABLES IN SCHEMA cdc TO cdc_reader` | **1** (Regel (6)) |
| M6 | **neuntes** Objekt: `cdc.administration_request` an beide Rollen | **1** (Regel (7)) |
| M8 | `GRANT SELECT ON ALL SEQUENCES IN SCHEMA cdc TO cdc_reader` | **1** (Regel (6)) |
| M10 | `GRANT SELECT ON cdc.process_heartbeat TO cdc_reader` | **1** (Regel (7)) |
| M11 | `GRANT SELECT ON cdc.table_schema TO cdc_reader` | **1** (Regel (7)) |
| M14/M15/M16 | die drei alten Bindungen | **1** / **1** / **1** |
| M17 | Reader-Grant auf `cdc.consumer` | **1** (Regel (7)) |
| M3 | `DO $$…$$`-Grant entfernt | **0** (Grenze 1, korrekt benannt) |
| M5/M9 | Lesegrant ohne schreibende Rolle | **0** (Grenze 2, korrekt benannt) |
| M12 | `GRANT USAGE ON SCHEMA public TO cdc_reader` | **0** (keine falsche Auslösung) |
| **M7/M13** | `GRANT CREATE` / `GRANT ALL ON SCHEMA cdc` | **0** → **D-2** |

Die frühere Acht-Namen-Liste ist damit **ohne messbaren Verlust** ersetzt: alle acht Objekte
werden von mindestens einer schreibenden Rolle gehalten, also fängt (7) sie einzeln — und
**das neunte** fängt sie mit (M6).

**Regel (6) fängt den legitimen `GRANT USAGE ON SCHEMA cdc` nicht** — der grüne Lauf trägt
den Grant selbst; `istKlassenGrant` prüft `all `/` all `/` in schema `, nicht `schema`.

**Die drei benannten Grenzen sind wirklich Grenzen — und ihre Begründungen tragen.**
- **Grenze 1 (`DO`-Block):** M3 → EC=0; der Parser trifft sie bewusst nicht. Zweite Hälfte
  trägt: `roles_test.go` fährt `SET ROLE cdc_admin` und verlangt `CREATE PUBLICATION`
  **erfolgreich** — ohne den Grant wäre das rot.
- **Grenze 2:** M9 → EC=0; gemessen liegt `tools/schema/schema.yaml` **nicht** im
  Build-Kontext (eigener Kontext-Lauf: **170** Dateien, `schema.yaml` **0 Treffer**), und
  `postgresstorage/schema.sql` hat **fünf** `CREATE TABLE`, **keine** View, kein `cdc.metrics`.
- **Grenze 3:** `nacharbeit-observability.sql:71` grantet `cdc.metrics` real; der Test liest
  genau eine Datei.

**F-3 ist erledigt — beide neuen Mutations-Behauptungen sind real und treffen genau:**
`Warn → ErrorContext` → EC=1, `level-Feld: ERROR, wollen WARN`; `Warn → InfoContext` →
EC=1, `unexpected end of JSON input`. Der Kommentar zitiert beide Meldungen wörtlich richtig
und beschreibt die **verschiedene** Wirkung korrekt.

**Der Fix hat die von §1 ausgeschlossene Datei nicht angefasst** —
`git diff --name-only 4298c4c^..4298c4c -- tools/schema/` → **0 Pfade**.

**Die geänderten Tests laufen im Gate, die Zahl ist unverändert** — `make coverage-gate`
Exit 0, `OK — Coverage 80.00% erfüllt Schwelle 70%`.

**Hygiene, Traceability, Struktur** — kein `nolint`, kein host-lokaler Pfad; Betreff mit
`ADR-0082`, keine Struktur-ID; kein `git mv`, kein `THRESHOLD`, keine Accepted-ADR berührt.

## Eigene Messungen (Exit-Codes direkt, ungepiped)

| Lauf | Exit | Ergebnis |
|---|---|---|
| `git show --numstat --format="" 4298c4c` | 0 | 3 Pfade: **+97/−14**, **+13/−4**, +128/−0 |
| `git diff --stat 4298c4c HEAD -- internal/` | 0 | **leer** — Review-Stand = `HEAD` |
| `git diff --name-only 4298c4c^..4298c4c -- tools/schema/ \| wc -l` | 0 | **0** |
| **18 LP2-Läufe** (Toolchain-Container, `--network none`) | je | **11 rot**, **6 grün**, Baseline grün |
| 2 Telemetrie-Mutationen | je 1 | die zwei zitierten Meldungen |
| Kontext-Messung | 0 | **170** Dateien; `schema.yaml` **0 Treffer** |
| `make coverage-gate` | **0** | `OK — Coverage 80.00% erfüllt Schwelle 70%` |
| `make docs-check` (auf `HEAD`) | 0 | 766 Dateien, 0 Befunde |
| `grep -rn "Lauf-Bericht" .` | 0 | **1** Treffer (die Zeile selbst) → D-3 |
| `grep`-Probe Absent-Text (`*.go`, repo-weit) | 0 | **1** Treffer (`:179`) → D-1 |

## Antwort auf die Schwerpunkte

**(1) F-2: Regeln statt Namensliste.** (a) **Ja, beide binden.** (b) **Ja, das neunte Objekt
wird gefangen** (M6) — der alte Test hätte es durchgelassen. (c) **Nein, (6) fängt den
legitimen USAGE-Grant nicht.** **Aber die Umstellung hat eine vierte Lücke neben den drei
benannten** — die Ausnahme von (7) ist klassenweit, während `:274-275` sie als *den einen*
Grant zusichert → **D-2**.

**(2) Die drei benannten Grenzen — sind es welche?** **Ja, alle drei, und ihre Begründungen
tragen** (gemessen, nicht geglaubt). **Die Liste ist trotzdem nicht vollständig** — sie nennt
die vierte Stelle nicht, sondern behauptet sie weg. Eine benannte Lücke ist zulässig; eine zu
enge Benennung ist es nicht.

**(3) F-3 — beide Mutationen, zwei Wirkungen.** Beide real gefahren, wörtlich richtig
zitiert. **F-3 geschlossen.**

**(4) Hat die Fixrunde etwas Neues falsch gemacht? — Ja, drei Stellen, alle in Sätzen.** Der
**Mechanismus** ist sauber: keine Mutation fälscht ein Ergebnis, kein Grant-Test verliert
Bindung. Neu falsch sind **Aussagen**: D-1 (abwesender Text), D-3 (Adresse ohne Artefakt),
D-4 (Begründung weiter als die Messung). **Der Fund der Runde liegt wieder im neuen Satz,
nicht im neuen Mechanismus.**

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | **1** |
| MEDIUM | 1 |
| LOW | 2 |
| INFO | 1 |

**Register-Hinweis (Modul 6):** D-1 ist **kein** Vorkommen von
`BEO-PGC/slice-chronik-in-code-kommentar` (dort Slice-/Wellen-Nummern; hier der abwesende
Text) — die Zuordnung entscheidet der Lese-Schritt. D-2 ist eine Fortsetzung des offenen
F-2-Vorgangs, **kein** zweiter Beleg.

## Verdikt

**Trägt die Fixrunde? — Im Mechanismus ja, in den Sätzen nicht.** F-2 ist als **Klasse**
geschlossen; F-3 ist belegt. Der Fix hat aber **drei neue Aussage-Defekte** in dieselben zwei
Dateien geschrieben, darunter einen HIGH nach der geschlossenen Liste.

**Merge-blockierend:** **ja** — D-1 (HIGH) und D-2 (MEDIUM). Kein Code-Nachzug nach außen;
beide Fundstellen liegen in den zwei Testdateien dieses Commits.

**Rückgabe-Pfeil Reviewer → Implementer: nötig, kurz** — vier Adressen, alle in
`internal/bootstrap/roles_rollout_file_internal_test.go`:
1. `:179-180` — den Satz auf die **geltende** Zusage stellen (D-1).
2. `:274-278` — Ausnahme und Code auf dieselbe Weite bringen **oder** die klassenweite
   Ausnahme als vierte Grenze benennen (D-2).
3. `:182` — den Beleg-Adressaten benennen, den dieser Lauf wirklich führt (D-3).
4. `:260` — das Subjekt der Begründung auf das ziehen, was die Messung zeigt (D-4).

D-5 (INFO) liefert **keinen** Rückgabe-Pfeil. Der [ADR-0085](../plan/adr/0085-build-kontext-ausnahme-test-only-zweck.md)-Pfad bleibt **außerhalb** dieses
Reports — er ist in `583a2bf`/`03fc53a` vom Architect fortgeschrieben und damit F-1s Träger.

**DoD-Häkchen „Review durchgeführt":** bleibt **offen** — es kommt eine Fixrunde.

**Nicht gefahren:** `make test` (`-race`, ganzer Baum) · `make test-store`/`-replication`/
`-integration`/`-notify` · `make gates` als Ganzes (stattdessen die Coverage-Stufe,
`docs-check` und `commit-traceability` einzeln, je mit eigenem Exit) · DoD/LP1–LP3 und die
§6-Ausgänge (Verifier/Closure).
