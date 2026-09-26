Zustand: offen — Ausgang: **weiter offen**, adressiert. Zähler (abgeleitet): 1× (evidence/slice-transformationen-start-reihenfolge.md).

Die Frage ist eine Entscheidung, keine Umsetzungsarbeit: ob der Vorlauf eine Zeitgrenze
trägt (Frist je Antrag oder je Durchlauf), was bei Ablauf geschieht (der Antrag endet
`failed`, oder der Stream startet und die Goroutine wiederholt den Antrag) und ob
`diagnose` und der `--healthcheck` das Warten zeigen (`LH-FA-ADM-003`). Sie berührt
`ADR-0112` Folgepflicht 5 (die Ordnung „Antrag vor erster Transaktion“ und ihre
Kosten) und `LH-QA-REL-001.a` (die Erfassung startet später, kein ACK entfällt).

Adresse: die Frage (e) in `welle-transformationen` §5 („Fragen für den nächsten
Architect-Zug“) und das Closure-Kriterium dort in §3: ein Architect-Verdikt vor der
Closure der Welle, im Umfang von `slice-transformationen-e2e-abhilfe` (dessen Beleg
fährt den Vorlauf mit wartendem Antrag am System) und
`slice-transformationen-betriebsdoku` (die Betreiber-Aussage: ein langer Antrag hält
den Stream-Start an, in Plan §6 von `slice-transformationen-start-reihenfolge` als
Aufschub mit Adresse geführt). Entscheidet das Verdikt eine Zeitgrenze, ist ein
Umsetzungs-Slice im Startpfad fällig; entscheidet es ein akzeptiertes Negativ, trägt es
die Aussage mit einer Messung am Prozess (der Wert „gesund“ während des Wartens ist
bisher hergeleitet) und einem Trigger.
