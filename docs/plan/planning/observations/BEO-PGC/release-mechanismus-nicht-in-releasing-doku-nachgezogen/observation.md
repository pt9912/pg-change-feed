# BEO-PGC/release-mechanismus-nicht-in-releasing-doku-nachgezogen

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die
Betreiber-Dokumentation des Release-Prozesses, keine eigene Sub-Area im
Sinn der Modus-Deklaration).

Die Beobachtung: `docs/user/releasing.md` ist laut eigenem §1-Zweck der
Träger, der beschreibt „wie ein Release ausgelöst wird" und trägt eine
Secret-Tabelle (§4) sowie eine Liste begleitender Workflows (§5). Ein neuer,
eigenständiger Release-Trigger — ein zweiter Tag-Namensraum mit eigenem
Workflow und eigenem Secret, unabhängig vom bestehenden `v*`-Server-Release —
entsteht in einem Slice, ohne dass dieser Doku-Träger im selben Zug mitzieht:
`docs-check` prüft Referenzen, nicht Vollständigkeit, kein Sensor fängt eine
fehlende inhaltliche Ergänzung.

**Abgrenzung zu `BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`:**
Derselbe zugrunde liegende Mechanismus (eine Meta-Pflicht an einer
Doku-Datei fällt im schreibenden Kontext aus dem Blick, weil kein Sensor sie
erzwingt) — aber jener Eintrag ist über Pfad und Text eng auf
`docs/user/benutzerhandbuch.md` (Betreiber-Oberfläche des laufenden
Prozesses: Umgebungsvariablen, administrative Aufgaben) gezogen und bereits
`verkörpert` (kombinierter Mechanismus: Selbstprüf-Instruktion +
Reviewer-HIGH-Punkt, beide **namentlich auf `docs/user/benutzerhandbuch.md`**
verengt — `.harness/skills/reviewer.md` nennt `releasing.md` nirgends). Dieser
Eintrag führt die Geschwister-Klasse für einen **anderen** Träger
(`docs/user/releasing.md`, Release-Auslösung statt Betriebs-Oberfläche) und
bleibt bewusst getrennt, bis ein zweiter Beleg zeigt, ob beide Klassen
zusammenfallen oder eigenständig bleiben sollten.

Deklaration: Planner-Agent (Closure-Rolle), 2026-09-19, auf Grundlage des
Review-Reports zu `slice-sdk-csharp-publish-workflow`, F-1 (erstes
Auftreten, dort ausdrücklich als „kein Steering-Loop-Zähler-Eintrag nötig
unter 3x" markiert — dieser Eintrag legt den Zähler an, statt ihn erst beim
dritten Auftreten nachzutragen, analog `BEO-PGC/dockerignore-default-deny-blockiert-neuen-pfad`).
