**Vorgang:** slice-routing-antragsweg (Review F-1, HIGH; derselbe Vorgang wie die Datei im Eintrag `beleg-befehl-traegt-seinen-satz-nicht`)

**Fund:** Die Umsetzungsentscheidung „Schreibweise der Zahl“ im Plan begründete eine Einschränkung der Annahmemenge von `order` mit der Tatsachenbehauptung, `jsonb` bewahre die Schreibweise und die Exponent-Form sei darum am Parser erreichbar und abzulehnen. Die Behauptung über PostgreSQL war weder gemessen noch als erwartet gekennzeichnet; die Messung (`jsonb` normalisiert `1e1` zu `10`) widerlegte sie. Die Einschränkung selbst ging zudem über den Wortlaut der Spec hinaus (Review F-2); die Fixrunde nahm sie zurück und trägt die Messung mit Messort im Plan.

**Form (Ausprägung):** Instanz B von `AGENTS.md` §3.12 am Träger **Slice-Plan**: eine Aussage über das Verhalten des Fremdsystems als Begründung einer Umsetzungsentscheidung, ohne Beleg-Anker. Der Vorgang zählt hier mit dem Mechanismus „die Behauptung trägt nicht“, im Eintrag `beleg-befehl-traegt-seinen-satz-nicht` mit dem Mechanismus „der Beleg trägt den Satz nicht“ (der Testfall); derselbe Vorgang, zwei Mechanismen. Schwere HIGH, daher eine Datei trotz Deckel; vor dem Merge vom Reviewer gefunden, Ausgang unverändert **verkörpert**.

Quelle: `docs/reviews/review-slice-routing-antragsweg.md` (F-1, F-2) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-routing-antragsweg.md` (§5 Zeilen F-1 und F-2). <!-- d-check:status-provenance -->
