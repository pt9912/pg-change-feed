Zustand: **verkörpert** — Ausgang: **verkörpert** → `AGENTS.md` §3.12 **Instanz B**;
Leser sind Verifier und Planner · seit slice-089. Der **Lese-Schritt der
`welle-20`-Closure** hat zwei frühere Kandidaten **verworfen** (§3.7 und der
ADR-Index-Kopf) — beide hat `ADR-0083` §Entscheidung Festlegung 3 ausgeschlossen.

**Ein Träger ist inzwischen benannt, aber nicht gebaut:** der Architect-Zug zu
`ADR-0078` hat für diese Wurzel eine Regel formuliert — *jeder Zahlenwert in
einem `Accepted`-Dokument trägt seinen Ursprung* (gemessen / übernommen /
**abgeleitet**); ein abgeleiteter Wert darf nie als gemessen erscheinen. Als
Träger schlägt er den Kopf von `docs/plan/adr/README.md` vor oder `AGENTS.md`
§3.7. **Das Benennen ist nicht das Bauen** — die Verkörperung ist Planner-Arbeit
und steht aus; deshalb ist der Ausgang hier *nicht* vorwegnehmend als
`verkörpert` gesetzt.

Zähler (abgeleitet): **8×** (evidence/slice-036.md, evidence/slice-082.md,
evidence/slice-081.md, evidence/slice-083.md, evidence/slice-095.md,
evidence/slice-sdk-kotlin-sse-client-flaeche.md,
evidence/slice-backfill-snapshot-reader.md,
evidence/slice-sdk-kotlin-cloudsmith.md) — **Schwelle erreicht**,
Ausgang beim Lese-Schritt der `welle-20`-Closure. Der achte Beleg
(`slice-sdk-kotlin-cloudsmith`, Review F-1, MEDIUM) ist eine Ausprägung im
Anwender-/Betreiber-Träger: eine noch nicht belegte Aussage („anonym lesbar“) stand
als Tatsache in fünf Trägern, der Plan führte sie als Erwartung mit Adresse, der
Post-Push-Lauf löste sie ein; die Aussage war wahr. Der siebte Beleg
(`slice-backfill-snapshot-reader`) trifft eine Plan-Begründung und ein DoD-Kriterium
(„der Testcontainer bringt die Extension-Typen nicht mit"), die der
Nachprüfungs-Review durch Messung widerlegte; verwandt mit
`BEO-PGC/mechanismus-erklaerung-ohne-werkzeugbeleg`, gezählt hier (Instanz B);
Ausgang bleibt **verkörpert**. Der sechste Beleg
(`slice-sdk-kotlin-sse-client-flaeche`) ist ein weiterer Vorgang derselben
Klasse und kein neuer Handlungsbedarf, mit einer bemerkenswerten
Herkunfts-Variante: die ungeprüfte Übernahme betraf hier einen
Registerstand (`BEO-PGC/zusatzkontext-kopplung-breiter-als-dod-wortlaut`),
der durch eine **parallele**, nicht die eigene Arbeit bereits überholt war
(`slice-sdk-csharp-sse-client-flaeche` einer Geschwister-Welle, vor
Slice-Start) — ein eigener neuer Beobachtungs-Eintrag für diese Variante
wurde bei der Closure geprüft und verworfen, weil der Mechanismus
(ungeprüfte Übernahme statt frischer Prüfung zum Schreibzeitpunkt)
identisch mit dieser Klasse bleibt, unabhängig von der Herkunft des
überholenden Ereignisses; gefunden hat es der Verifier, nicht der
Reviewer, dessen Diff-Prüfung „keine Änderung an `observations/`" korrekt,
aber ohne Abgleich gegen den absoluten Registerstand blieb (Details:
`evidence/slice-sdk-kotlin-sse-client-flaeche.md`). Der fünfte Beleg
(`slice-095`, 2. Durchlauf) ist ein weiterer Vorgang derselben Klasse und
kein neuer Handlungsbedarf: eine §2-DoD-Zeile behauptete
ungeprüft „dieses Repo führt Wellen-Betrieb", ein unangepasster
Vorlagen-Standardtext, der dem eigenen Slice-Kopf widersprach — bei der
Slice-Closure gefunden und korrigiert. Der vierte Beleg ist ein weiterer
Vorgang derselben Klasse und kein neuer Handlungsbedarf: `slice-083` berief sich
auf **nicht existierende** Nachbar-Clients („wie die drei anderen"), und der
Verfasser war der Planner. Das Erstvorkommen (`slice-036`) wurde seinerzeit **ohne
Kennung** notiert: das Review nannte das Label, legte aber kein Verzeichnis an.
Der Eintrag entstand mit dem zweiten Auftreten und zitiert das erste über
**dasselbe Label** — die Zuordnung ist belegt, nicht abgeleitet. Zwei Funde **im
selben** Vorgang (`slice-082`: Review F-1 und Verifikation V-1; `slice-081`:
Review F-1 und die Berichtigung in §1/§2) sind je *eine* Gelegenheit — der
Zähler misst Wiederholung über Vorgänge, nicht die Zahl der Funde.
