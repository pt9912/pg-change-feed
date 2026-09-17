# Beleg: slice-101

Vorgang: `slice-101` — NATS-Client in C# und Kotlin, konkret sein
Abschnitt §4 (Rückführungen, Trigger 1) gegen den tatsächlichen Verlauf.

Fund: §4 nannte die Bedingung wörtlich vorab — „wenn eine der beiden
gepinnten Client-Bibliotheken transitiv eine zweite, gepinnte Abhängigkeit
erzwingt, die eine eigene Bewertung braucht (Lizenz, Sicherheits-Historie)
— dann ist die Größenannahme… falsch." Die Bedingung trat ein: `io.nats:
jnats:2.26.3` erzwingt real, unconditional, `org.bouncycastle:bcprov-
lts8on:2.73.12.1` (`compile`-Scope). Der Implementer schrieb das Artefakt
(Pin-Kommentar in `examples/kotlin/nats-client/build.gradle.kts`, Commit
`7656c36`) samt eigener Bewertung, **bevor** die §4-Bedingung an
Architect/Planner zur Prüfung ging — die Bewertung selbst blieb zudem
unvollständig belegt (Lizenz mit Beleg, Sicherheits-Historie nur als
Reputationsaussage, `docs/reviews/review-slice-101.md` F-1/F-2). Der
benannte Weg (`in-progress` → `next`) war zum Zeitpunkt der
Reviewer-Prüfung nicht mehr gangbar, ohne einen bereits fertigen,
getesteten Kotlin-Client zu verwerfen; der tatsächliche Weg war der
Konflikt-Pfad (Reviewer → Architect, `docs/reviews/
architect-verdict-slice-101-jnats-bouncycastle.md`).

**Zweite, unabhängige Instanz derselben Klasse wie `slice-073`:** Dort
deckte die nachträgliche Prüfung eine echte Diskrepanz auf (der Hook
bildete zwei Prüfungen nicht deckungsgleich ab); hier deckt die
nachträgliche Prüfung (Architect-Verdikt §2/§3) **keine** Diskrepanz auf,
sondern bestätigt — mit zusätzlicher Tiefe (drei weitere Advisories über
NVD, alle vor oder exakt bei der gepinnten Version behoben) — genau das
Ergebnis, das der Implementer bereits (unvollständig belegt) behauptet
hatte. Das Muster selbst — Bedingung vorab benannt, Artefakt vor der
Klärung geschrieben, Klärung läuft über den Konflikt-Pfad statt über einen
Neuschnitt — ist in beiden Fällen identisch; nur der inhaltliche Ausgang
unterscheidet sich (dort Korrektur nötig, hier Fortsetzen bestätigt).

Zusätzlich benannt (noch nicht Gegenstand dieses Zählers): Der Architect
hat Modul 8 §Konflikt-Pfad — dem Wortlaut nach für ADR-Konflikte
geschrieben — hier erstmals explizit auf einen Plan-Trigger-Konflikt
übertragen und diese Übertragung begründet; der Verifier hat das als
eigenständige, plausible aber nicht die einzig mögliche
Auslegungsentscheidung benannt (`docs/reviews/verify-slice-101.md` §3.1).

Quelle: `docs/plan/planning/done/slice-101-nats-client-csharp-kotlin.md`
§4/§7 · `docs/reviews/review-slice-101.md` F-1/F-2 ·
`docs/reviews/architect-verdict-slice-101-jnats-bouncycastle.md` ·
`docs/reviews/verify-slice-101.md` §3.
