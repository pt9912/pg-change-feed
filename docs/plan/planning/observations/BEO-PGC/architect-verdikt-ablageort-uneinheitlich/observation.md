# BEO-PGC/architect-verdikt-ablageort-uneinheitlich

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die
Harness-Struktur/Dokumentation, keine eigene Sub-Area im Sinn der
Modus-Deklaration).

Die Beobachtung: Architect-Verdikt-Reports aus Modul 8s Konflikt-Pfad-/
Trigger-Audit-Verfahren („kein neues ADR nötig, bestehende ADRs decken das
ab" o. ä.) landen an zwei verschiedenen Orten mit zwei verschiedenen
Namensschemata, statt an einem konsistenten Ort:

- 9 Dateien unter `docs/plan/adr/architect-review-slice-NNN.md` bzw.
  `architect-review-welle-NN.md` (`slice-011`, `slice-013`, `slice-014`,
  `slice-015`, `slice-016`, `slice-021`, `welle-1`, `welle-5`, `welle-6`).
- 1 Datei mit dem Namensschema `architect-verdict-` unter `docs/reviews/`
  (der Architect-Verdikt zu `slice-030`/`ADR-0015`).

Beide Gruppen sind derselbe Artefakt-Typ (Architect-Verdikt, Modul 8
§Konflikt-Pfad als Rollen-Sequenz) — kein ADR selbst, deshalb korrekt
außerhalb des ADR-Index (`docs/plan/adr/README.md`, „eine Zeile je
ADR-Datei"); das ist keine Index-Lücke. Die Uneinheitlichkeit liegt im
Ablageort (`docs/plan/adr/` vs. `docs/reviews/`) und im Namensschema
(`architect-review-` vs. `architect-verdict-`) selbst.

Deklaration: direkte Nutzerfrage zum ADR-Index, 2026-09-13 — kein Slice
oder Welle hat diese Beobachtung ausgelöst.
