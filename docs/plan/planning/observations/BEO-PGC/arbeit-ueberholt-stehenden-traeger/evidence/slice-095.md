# Beleg: slice-095

Vorgang: `slice-095` — die drei fehlenden Beispiel-Clients (2. Durchlauf,
Liefer-Punkt 2/3: `examples/grpc-client/` + Rest-Träger-Nachzug).

Fund: §3 der eigenen Plan-Datei begründete die geplante Änderung an
`.a-check.yml` mit „keine Kante — `ADR-0079` Festlegung 1". Diese Begründung
war bei ihrer ursprünglichen Niederschrift (erster `in-progress`-Durchlauf,
vor der Rückführung) für alle vier Beispiel-Clients wahr. Der zweite
Durchlauf lieferte real eine neue Kante `{from: examples, to: contract}`,
seit `slice-097`s Umzug der Vertragsfläche auf den öffentlichen Pfad
`gen/cdc/stream/v1` nötig — die Lieferung selbst ist ADR-konform und korrekt
begründet (Commit `1812e36`), aber sie fasste die eigene §3-Begründungszeile
des Plans nicht an. Die Zeile blieb stehen und behauptete weiterhin „keine
Kante" für alle vier Clients, obwohl der gRPC-Client jetzt eine hat.

Gefunden hat es der Reviewer (F-1, LOW, Review zu `slice-095`) — nicht der
Implementer-Suchlauf für §3.13 (der war auf `grep`-Treffer für interne
Pfad-Referenzen im Code beschränkt und fand hier nichts, weil dieser Fund
kein Code-Pfad, sondern eine Plan-Prosa-Zeile ist).

**Selbstreferentielle Variante:** anders als die bisherigen fünf Belege
(externe Träger wie `harness/sensors/coverage-gate.md` oder — im Fall
`slice-097` — eine **andere** Slice-Plan-Datei) ist der überholte Träger hier
dieselbe Plan-Datei, die auch die überholende Arbeit trägt. Der Mechanismus
bleibt derselbe: eine Arbeit bewegt eine beschriebene Eigenschaft (die
`.a-check.yml`-Kanten-Struktur) und macht damit einen Satz falsch, der diese
Eigenschaft beschreibt, ohne dass die Arbeit selbst diesen Satz anfasst.

In dieser Closure direkt korrigiert (Planner-Trägerpflege, Commit 1) — kein
Carveout, kein Folge-Slice nötig, weil sofort behebbar (Formvorbild
`slice-097` §6, dort ebenfalls „eingetreten, direkt korrigiert").

Quelle: Review zu `slice-095`, F-1 ·
`docs/plan/planning/in-progress/slice-095-beispiel-clients-drei.md` §3
(berichtigt bei dieser Closure) ·
`docs/plan/planning/done/altbestand/slice-097-umzug-vertragsflaeche.md` §6/§7
(Formvorbild) · `AGENTS.md` §3.13.
