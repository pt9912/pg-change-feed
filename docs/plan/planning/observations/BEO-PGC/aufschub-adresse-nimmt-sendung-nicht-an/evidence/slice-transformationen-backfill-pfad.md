**Vorgang:** slice-transformationen-backfill-pfad (Verifikation V-4, LOW)

**Fund:** Der Slice trägt den Handbuch-Gegenstand „Regeln gelten auch für einen Backfill“ als Aufschub mit
Adresse (`slice-transformationen-betriebsdoku` §2, Stichwort „Backfill-Bezug“). Der Verifier fand, dass die
Adresse die Sendung nur dem Wort nach annimmt: die Folgepflicht 4 der Run-Klassen-Entscheidung verlangt einen
Betriebshinweis zur Abhilfe im Run (Regel entfernen, **neuer** `cdc.backfill_table`-Antrag, kein Prozessneustart),
und die Abhilfe-Zeile der Adresse nannte nur den Erfassungspfad („Prozess starten“); weder die Entscheidung noch der
neue Antrag standen dort (`git grep -n '0117\|neuer Antrag'` im Plan der Adresse: kein Treffer). Gefunden vor dem
Merge; die Closure schrieb die Run-Abhilfe als committeten Text in den Plan der Adresse.

**Form (Ausprägung):** die Adresse deckt den **Wortlaut** des Gegenstands („Backfill-Bezug“), nicht seinen
**Ablauf**; die zweite Hälfte des Ausgangs-Vorschlags (der Gegenstand steht als committeter Text im Plan der
Adresse) hätte den Fall am Absender gefangen.

Quelle: `docs/reviews/verifikation-slice-transformationen-backfill-pfad.md` (§8 V-4, §7 Zeile Handbuch-Kandidatenlauf). <!-- d-check:status-provenance -->
