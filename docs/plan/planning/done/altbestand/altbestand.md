# Sammelpunkt `altbestand` — wellenloser Archiv-Bestand vor dem ersten Archivierungslauf

**Schlüssel:** `altbestand` — kein Welle-Plan im Sinne von Modul 6 (keine
eigene Zielsetzung, kein Slice-Schnitt), sondern der von
[`ADR-0096`](../../adr/0096-altbestand-schluessel-fuer-wellenlosen-archiv-bestand.md)
gewählte, dauerhafte Sammel-Schlüssel für **wellenlose** `done/`-Slices, die
vor dem ersten `archive-welle`-Lauf dieses Repos flach lagen. Dient
ausschließlich dazu, dem Werkzeug (`ai-harness-init archive-welle`) den
formal geforderten Plan-Kandidaten unter `docs/plan/planning/done/` zu
liefern (`internal/archive/collect.go` `planName`).

**Verantwortlich:** pt9912. **Datum:** 2026-09-18.

---

## 1. Ziel

Kein Liefer-Ziel im Sinne eines Slice-Plans. Der reale Inhalt dieses Laufs —
welche Slices eingesammelt wurden, warum, und mit welchem Ausgang — steht in
[`altbestand-results.md`](altbestand-results.md).

## 2. Bezug

[`ADR-0096`](../../adr/0096-altbestand-schluessel-fuer-wellenlosen-archiv-bestand.md)
(Entscheidung für diesen Schlüssel), `slice-archive-altbestand-vollzug`
(vollzieht den Lauf).
