**Vorgang:** slice-retention-lauf-speicher-begrenzung (Review F-3, MEDIUM; Verifikation V-1, MEDIUM)

**Fund:** Der Architect-Zug zu den zwei offenen Bewertungen (neuer Consumer im laufenden Lauf, Transaktion über zwei Seiten) lieferte Nachzüge, die im Plan und im Handbuch fehlten: die Risiko-Ausgänge in Plan §6 trugen „bei Closure“, die Plan-Zeile „nicht breiter geworden (Dauer des Laufs nicht länger)“ stand gegen die abgeleiteten Grenzen des Verdikts (Regelbetrieb 0,38 bis 2,0 s, Extremlauf etwa 26 s statt etwa 3,4 s), und der Handbuch-Satz „gilt seine Position erst im nächsten Durchlauf“ trug die Folge (die Löschung dieses Durchlaufs) und die Nicht-Atomarität über Transaktionen nicht. Der Verifier fand die drei Träger durch den Abgleich von Verdikt und Träger; die Closure zieht sie nach (Handbuch 1.64).

Quelle: `docs/reviews/review-slice-retention-lauf-speicher-begrenzung.md` (F-3) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-retention-lauf-speicher-begrenzung.md` (§9 V-1). <!-- d-check:status-provenance -->
