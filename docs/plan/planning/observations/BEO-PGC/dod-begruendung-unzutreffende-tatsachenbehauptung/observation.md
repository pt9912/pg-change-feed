# BEO-PGC/dod-begruendung-unzutreffende-tatsachenbehauptung

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die Form, in der
ein DoD-Kriterium seine Erfüllung begründet, keine eigene Sub-Area im Sinn der
Modus-Deklaration).

Die Beobachtung: Ein DoD-Kriterium begründet seine Erfüllung mit einer
**Tatsachenbehauptung, die nicht geprüft wurde** — übernommen aus einem
Bericht, einem Nachbardokument oder der eigenen Erinnerung, statt am Gegenstand
gemessen. Im Kriterium ist die Behauptung nicht als Annahme erkennbar; sie
liest sich wie ein Beleg. Die Klasse trifft den **Plan**, nicht den Code: der
Diff war in beiden belegten Fällen korrekt.

Belegt an `slice-036`: eine Begründung war ungeprüft aus `harness/README.md`
übernommen und deckte sich nicht mit dem zitierten Text; ein Datei-Diff gegen
den zitierten Text deckte es auf. Belegt an `slice-082`, **zweimal im selben
Vorgang**: das Kriterium sagte, der Tier-Lauf beweise die reale Anwesenheit
**beider** Tabellen (er beweist nur die, deren Fehlen fatal ist — die andere
wird best-effort geschrieben), und es nannte „15 Tabellen, 6 Sichten,
4 Funktionen" — eine aus dem Review-Report übernommene, nie selbst gemessene
Zahl; der Rollout-Report führt 10 Tabellen und 4 Sichten.

**Warum das zählt:** Das Kriterium ist die Prüf-Form des Liefer-Punkts
(Baseline-Regelwerk `modul-05-planning-harness.md` §Ziel-Form: Slice). Trägt es
eine ungeprüfte Zahl, prüft die Verifikation gegen eine Behauptung statt gegen
den Gegenstand — die DoD sieht erfüllt aus und ist es vielleicht nicht. Der
Mechanismus ist beide Male derselbe: **nicht geprüfte Übernahme**.
