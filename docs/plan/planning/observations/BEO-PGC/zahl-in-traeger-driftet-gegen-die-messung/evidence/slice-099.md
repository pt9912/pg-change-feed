# Beleg: slice-099

Vorgang: `slice-099` — Kotlin-Sprachwurzel und HTTP-Client (`ADR-0090`,
`ADR-0087`).

Fund: `ADR-0087`s Digest-Pinning-Tabelle (§Entscheidung Festlegung 3, Zeile
254) nennt für `eclipse-temurin:21-jdk` einen sha256-Wert mit **63** statt
64 Hex-Zeichen (fehlendes „d“) — strukturell ungültig, nie real auflösbar.
Der Reviewer hat den Fund selbst unter genau der Klasse dieses Eintrags
eingeordnet (Review zu `slice-099`, F-1, `klasse`: „Zahl im
Träger ohne Ursprung — oder gegen die Messung driftend“), der Verifier hat
ihn per eigener, unabhängiger `docker manifest inspect`-Messung bestätigt
(Verifikationsbericht zu `slice-099`, #1).

**Abweichende Form gegenüber den bisherigen sieben Belegen:** Hier hat sich
der Gegenstand nicht *bewegt* (kein Drift durch fortschreitende Arbeit) —
der Tabellenwert war seit `ADR-0087`s Annahme **nie** korrekt (63 statt 64
Zeichen), unbemerkt bis zum ersten realen Bau dieser Zelle. Das fällt unter
die erste Hälfte des Klassennamens („ohne Ursprung“), nicht unter die
zweite („driftend“): eine `Accepted`-ADR behauptete einen scheinbar
gemessenen Digest, der nie eine reale, verifizierte 64-stellige Registry-
Antwort war. Der tatsächlich gebaute und gepinnte Wert
(`examples/kotlin/Dockerfile:16`) war davon nie betroffen.

**Auflösung:** engräumige Folge-ADR `ADR-0093` (`Supersedes ADR-0087` für
genau diese eine Tabellenzelle, Architect-Zug, parallel zum laufenden
Slice) — kein neuer Sensor: `ADR-0093`s eigene Fitness-Function-Tabelle
hält fest, dass kein Werkzeug Digest-Korrektheit einer ADR-Tabellenzelle
gegen die Registry prüft; die verfügbare Falsifikation bleibt die
Messung selbst (unverändert gegenüber dem bereits verkörperten Ausgang
dieses Eintrags).

**Ausgang bei dieser Closure:** kein neuer Schwellen-Übertritt — die Klasse
ist bereits verkörpert (`AGENTS.md` §3.12 Instanz A,
`.harness/skills/reviewer.md` HIGH-Punkt, seit `slice-089`). Dieser Beleg
erweitert ihre Reichweite auf ADR-Tabellenwerte, die eine ADR selbst nie
gegen ihre reale Quelle nachprüfte, statt nur auf Zahlen, die durch
spätere Arbeit stale wurden.

Quelle: Review zu `slice-099`, F-1 · Verifikationsbericht zu `slice-099`,
#1 · `docs/plan/adr/0093-digest-korrektur-adr-0087-kotlin-basis-image.md`.
