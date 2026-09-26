**Vorgang:** slice-transformationen-antragsweg-schema (Review F-6, LOW; Mutation F1 des Reviewers; erste Fassung des Tests, in der Fixrunde vom Implementer an der Mutation rot gesehen)

**Fund:** Die Zusage „`SECURITY DEFINER`, gepinnter `search_path`“ der zwei Funktionen war an ihrer Eingabeseite nur zur Hälfte gebunden: `SECURITY INVOKER` färbte Lauf 5 des Guard-Skripts rot, das Entfernen von `SET search_path = cdc, pg_temp` färbte keinen Test — kein Test las `prosecdef` oder `proconfig` (`git grep` 0 Treffer); dieselbe Lücke bestand für die fünf Funktionen des Parents. Die Fixrunde band beide Eigenschaften für alle sieben Funktionen über den Katalog. Die erste Fassung des Tests (`proconfig = ARRAY[…]` in einer Verneinung) blieb bei fehlendem `search_path` grün, weil `NOT (… AND NULL)` in SQL NULL ist; der Implementer sah es an der Mutation und stellte auf `IS DISTINCT FROM` um.

**Form (Ausprägung):** neue Form (anderer Träger-Typ: **Katalog-Eigenschaft** einer Datenbankfunktion im Store-Test): die Eingabeseite einer Sicherheits-Härtung ist ihr Fehlen, und die Prüfung liest die Katalog-Spalte, deren NULL-Fall die Verneinung nicht trägt.

Quelle: `docs/reviews/review-slice-transformationen-antragsweg-schema.md` (F-6, Mutationen F1, F2) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-transformationen-antragsweg-schema.md` (§4 Zeile M3, §5 Zeile F-6). <!-- d-check:status-provenance -->
