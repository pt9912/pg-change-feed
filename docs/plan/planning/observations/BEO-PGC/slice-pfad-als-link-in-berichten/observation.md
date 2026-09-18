# BEO-PGC/slice-pfad-als-link-in-berichten

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die
Verweisform auf Slice-Pläne in Berichten und Entscheidungsdokumenten, keine
eigene Sub-Area im Sinn der Modus-Deklaration).

Die Beobachtung: Ein Slice-Plan wechselt im Lifecycle seine Ablage
(`open/` → `next/` → `in-progress/` → `done/`). Wird er aus einem Bericht
oder einer Entscheidung als **Markdown-Link mit festem Verzeichnis**
adressiert (`[…](../plan/planning/in-progress/slice-NNN-….md)`), bricht der
Link beim nächsten Übergang — `docs-check` meldet `target-missing` und der
Gate-Lauf wird rot, ohne dass sich am Inhalt etwas geändert hätte. Die
Repo-Konvention ist die **Kennungs-Zitierung** (`` `slice-NNN` ``) oder ein
Inline-Code-Pfad; nur echte, lage-stabile Ziele werden verlinkt.

Deklaration: Review und Verifikationsbericht zu `slice-068` sowie
der Architect-Verdikt zur Spaltenausschluss-Dauerhaftigkeit (beide
Fälle in derselben Sitzung real aufgetreten und vom Doku-Gate gefangen).
