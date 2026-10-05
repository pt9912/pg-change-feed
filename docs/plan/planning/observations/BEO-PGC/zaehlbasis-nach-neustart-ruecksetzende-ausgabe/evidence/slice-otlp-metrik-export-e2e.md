**Vorgang:** slice-otlp-metrik-export-e2e (Post-Push-Lauf, Planner-Closure).

**Fund:** e2e.yml Lauf `37225751324` (Commit `19e83251`), Leg PostgreSQL 18, Phase „Wiederaufnahme nach
docker start“ rot. Gemessen am gepinnten Collector 0.162.0: der Collector legt die Datei seines
file-Exporters bei jedem Start leer neu an; der Runner verglich gegen die Zeilenzahl vor dem Stopp
(Stand 10) und brauchte mindestens 11 Exporte zu je 5 s, mehr als das 40-s-Fenster. Ein Produktfehler
ist ausgeschlossen (Exporter mit IP-Wechsel gemessen). Lokal lief dieselbe Phase grün.

**Behebung:** `a19c28cc` — der Runner wartet auf „Everything is ready“ per `docker logs --since` und
nimmt die Basis an der neuen Datei; danach lokal `make test-integration` Exit 0 (3 Exporte in 16 s),
Post-Push-Lauf `37254000299` am Stand `759d2f2f` beide Legs `success`.
