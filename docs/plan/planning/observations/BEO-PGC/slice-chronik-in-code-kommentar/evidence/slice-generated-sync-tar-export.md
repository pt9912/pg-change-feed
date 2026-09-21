**Vorgang:** slice-generated-sync-tar-export
**Fund:** Zwei getrennte Instanzen, dieselbe Klasse, zwei unterschiedliche
Träger. (1) Der Implementer trug im Kopf-Kommentar von
`tools/harness/generated-sync.sh` implizite Chronik-Sprache ein („Frueher
lief hier `docker build --target proto` gefolgt von einem manuellen
`protoc`-Aufruf gegen einen … Bind-Mount. Der Mount war die Fehlerquelle
auf …") — ein Vorher/Nachher-Vergleich ohne Slice-Zahl, aber wörtlich das
`AGENTS.md` §3.7-Falsch-Beispiel („die frühere Fassung prüfte nur die
Länge"). Erstes Vorkommen dieser Klasse in einem **Shell-Skript**
(bislang: Go, SQL, `schema.yaml`) — Datei-Typ-Erweiterung der bestehenden
Verkörperung, keine neue Klasse. Der unabhängige Reviewer fing die Stelle
vor Merge (F-1, HIGH).

Quelle: `docs/reviews/review-slice-generated-sync-tar-export.md` <!-- d-check:status-provenance -->
Fixrunde real geprüft, 0 HIGH danach.

Quelle: `docs/reviews/review-slice-generated-sync-tar-export-fixrunde.md` <!-- d-check:status-provenance -->
Achter Beleg, dieselbe Diagnose: die Verkörperung (Reviewer-HIGH-Punkt +
Schritt-20-Enumerationslauf) trägt, kein Hard-Rule-Verstoß hat `main`
erreicht.

(2) **Zweite Instanz der bereits als offene Grenze notierten
Doku-Prosa-Variante** (siehe `evidence/slice-schema-rollout-zentrale-
idempotenz-wache.md` Teil 2, dort erstmals in `harness/README.md`
gefunden): derselbe Chronik-Stil stand in `harness/sensors/generated-sync.md`
(„Bis `slice-generated-sync-tar-export` lief hier stattdessen …") — vom
Reviewer selbst gefunden (F-2, INFO, nicht erst vom Verifier wie beim
ersten Vorkommen), bewusst INFO statt HIGH eingestuft, weil
`AGENTS.md` §3.7 wörtlich „Code, Konfiguration und Skripte" nennt und kein
Gate Doku-Prosa deckt. Blieb über die Fixrunde hinweg unverändert stehen
(bestätigt in der Fixrunde und im Verifikationsbericht) — bewusst kein
Fixrunde-Anlass, dieselbe Einordnung wie beim ersten Vorkommen. Zweiter
realer Beleg dafür, dass die Doku-Prosa-Variante kein Einzelfall war;
weiterhin unter der Zwei-Beleg-Schwelle für eine eigene Registrierung, aber
ein Signal mehr für den nächsten Lese-Schritt (Skopus-Erweiterung des
HIGH-Punkts auf Sensor-/README-Doku erwägen).
