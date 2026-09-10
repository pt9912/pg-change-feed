# Welle welle-2 — Reale PostgreSQL-Integration — Closure-Notiz

> **Zitier-Form** *(bleibt stehen — Norm, kein Ausfüll-Hinweis).* Dieses
> Artefakt friert ein; was es zitiert, bewegt sich weiter. Deshalb: **Kennung,
> nicht Adresse** — `slice-NNN` statt seines Lifecycle-Pfads, `make <target>`
> statt eines Links auf die Sensor-Datei, eine Baseline-Stelle als
> `v<X.Y.Z>` · `regelwerk/<datei>.md` §<Abschnitt> statt als Link
> (Baseline-Regelwerk `grundlagen-harness-dateien.md` §harness/README.md als
> Einstiegspunkt — beim Ausfüllen mit dem adoptierten Tag schreiben).

**Welle:** welle-2
**Abschluss:** 2026-09-10
**Verantwortlich:** pt9912

## Was wurde geliefert?

<!-- BEDIENHINWEIS: Ergebnis, nicht Taetigkeit. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — *was gelernt wurde*: geliefert · was
funktionierte · was anders lief. Mit ID-Bezug, wo es einen gibt.

- PostgreSQL-ChangeStore-Adapter real (Persist mit
  Idempotenz-Deduplizierung über die Transaktions-ID, deterministisches
  Lesen — Teil-Beleg [`LH-QA-REL-001`](../../../../spec/lastenheft.md)/[`LH-QA-REL-002`](../../../../spec/lastenheft.md),
  10/10 Store-Tests gegen reale gepinnte PostgreSQL).
- Replication-Stream-Adapter real: `pgoutput` dekodiert (pglogrepl,
  [`ADR-0008`](../../../plan/adr/0008-pgoutput-standard.md)-Infrastruktur), Stream-Vertragsverstöße sichtbar, realer ACK —
  Beleg [`LH-FA-CAP-001`](../../../../spec/lastenheft.md)…003 am realen Pfad.
- Compose-Umgebung und MVP-Integrationstest: CDC aktivieren →
  INSERT/UPDATE/DELETE → lesen → Reihenfolge/Inhalt — Beleg
  `make test-integration` Exit 0 (MVP-Schnitt, Abschnitt 1).
- Bootstrap-Verdrahtung (wellenlos): das Binary trägt die reale
  Pipeline ([`ADR-0026`](../../../plan/adr/0026-composition-root.md)) — der Feed-Container fährt
  CDC-Runtime statt `--version`-Smoke (Digest `ee795a88…` →
  `899f3954…`, Beleg am HEAD).

## Was hat funktioniert?

<!-- BEDIENHINWEIS: was du im naechsten Zyklus bewusst wieder so machen wuerdest. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3.

- Die Rollen-Sequenz mit Übergabe-Artefakten lief über fünf Slices
  (slice-004…007) — Review und Verifier in frischem Kontext, alle
  Sensor-Belege selbst gefahren.
- Das Standing-Gate `commit-traceability` ([`ADR-0045`](../../../plan/adr/0045-commit-traceability-standing-gate.md)) färbte beim
  ersten echten Verstoß selbst rot (F-3-Commit des slice-007-Fix-Zugs)
  — die Maschine greift.
- Die a-check-Edges tragen die Schichten-Constraints am vollen Baum
  (alle Layer-Globs matchen echten Content).

## Steering-Loop-Einträge

<!-- BEDIENHINWEIS — keine Norm; faellt beim Kopieren weg (README.md
§Verwendung, Schritt 5) und darf deshalb nichts Tragendes halten. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 (hier stehen **nur** Beobachtungen, die im
Register 3× erreicht haben; jeder Eintrag nennt seine `BEO-<NNN>`) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (Feld
und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der Backticks; die
**Spec-Lücke** trägt statt `liegt in` ihre `LH-*`-ID — das ist kein Versehen).

- **Standing-Gate `commit-traceability`** (neuer Sensor, [`ADR-0045`](../../../plan/adr/0045-commit-traceability-standing-gate.md)):
  je Commit ≥ 1 `LH-*`/`ADR-*` **und** keine Struktur-ID im Betreff —
  färbte beim ersten echten Verstoß selbst rot (F-3-Commit des
  slice-007-Fix-Zugs).
  — liegt in `.d-check.yml` (commits-Abschnitt) und
  `tools/harness/commit-traceability.sh` ·
  Auslöser: Commit-Kennung-Klasse (slice-005 F-5 = 4. Auftreten der
  Struktur-ID-Variante, review-slice-006 F-2 = 3. Auftreten der
  Kennungslos-Variante — Sequenz lief bei slice-003, Verkörperung
  [`ADR-0045`](../../../plan/adr/0045-commit-traceability-standing-gate.md)).
- **Datei-Abschluss-Zeilenumbrüche** (Klasse 3×: slice-001 F-6,
  slice-004 F-5, slice-005 F-6, bestätigt durch verify-slice-005): die
  Fix-Runden schlossen alle berührten Dateien; der Ausgang läuft über
  die Rotation — kein BEO-Eintrag nötig, da die Klasse maschinell nicht
  belegt ist (kein d-check-Modul).
- **Zwei gemischte Commits** (`13abe88`, `e3e36f2`): parallele
  Rollen-Läufe mischten Implementer-Content mit Review-/Plan-Artefakten,
  eine Review-Basis (`193716e`) wurde verdrängt — disclosed in den
  slice-004/005-Closure-Notizen; Konsequenz: Rollen-Läufe strikt
  sequenzieren (Implementer-Ende vor Review-Artefakt-Commits).
- **Die Image-Beleg-Formel** (`96c47af`) wurde von review-slice-005
  F-7 falsifiziert und in [`ADR-0044`](../../../plan/adr/0044-image-beleg-semantik.md)
  korrigiert (Digest lauf-gebunden; Inhalts-Streits über Binary-sha256).
- **Die Commit-Kennung-Klasse** zählte 4× (Struktur-IDs in Messagen,
  Commits ohne Kennung) — verkörpert als Standing-Gate ([`ADR-0045`](../../../plan/adr/0045-commit-traceability-standing-gate.md));
  die Rotation aus dem 5-Commit-Fenster ist der Ausgang.
- **Der d-migrate-1.2.0-Blocker** am CHECK-Ausdruck erzwang die
  Nacharbeit-Ausweichform; die gemeldete Fix-Release retiriert sie
  (Pin-Hebung + `chk_change_operation` direkt ins YAML).
