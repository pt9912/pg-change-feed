Stand: **offen** (1×, unter der Schwelle) — die Beobachtung trägt keinen der
drei Ausgänge; sie ist der Lerneintrag dieses Slice in der Form *benannte
Spec-Lücke*.

Die Lücke: Der Abschnitt §4 des Slice-Plans verlangt, die
Rückführungs-Bedingung **vorab** zu benennen (Baseline-Regelwerk
`modul-05-planning-harness.md` §Lifecycle als State Machine), und der Übergang
selbst trägt den *Grund* nach. Keine Stelle verlangt, die Bedingung **beim
Schreiben der Umsetzung** auszuwerten — gelesen wird sie erst bei der Closure.
Ist sie dann eingetreten, ist der benannte Weg (`in-progress` → `next`)
versperrt, weil das Artefakt bereits existiert, und der tatsächliche Weg ist
der Konflikt-Pfad.

Gelesen wird der Eintrag im Sichtungs-Schritt der Slice-Planung
(`docs/plan/planning/observations/README.md`). Zähler (abgeleitet): **2×**
(evidence/slice-073.md, evidence/slice-101.md) — weiterhin unter der
3×-Schwelle, Stand bleibt `offen`. Der zweite Beleg (`slice-101`) betrifft
denselben Mechanismus an einer anderen Fundstelle: Statt einer
Kongruenz-Prüfung (Hook vs. Prüfungen, `slice-073`) ist es hier eine
transitive Abhängigkeit samt Lizenz-/Sicherheits-Bewertung
(`io.nats:jnats` → `bcprov-lts8on`); anders als bei `slice-073` deckte die
nachträgliche Prüfung diesmal **keine** Diskrepanz auf, sondern bestätigte
das bereits (unvollständig belegt) behauptete Ergebnis — siehe
`evidence/slice-101.md`.
