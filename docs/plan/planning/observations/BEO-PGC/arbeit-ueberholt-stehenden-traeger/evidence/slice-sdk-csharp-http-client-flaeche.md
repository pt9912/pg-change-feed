# Beleg: slice-sdk-csharp-http-client-flaeche

Vorgang: `slice-sdk-csharp-http-client-flaeche` — liefert
`PgChangeFeedHttpClient` (zehn öffentliche Methoden, `SPEC-018`/`SPEC-022`)
als real existierende HTTP-API-Client-Fläche des .NET-SDK-Packages.

Fund: Der Vorgänger-Slice (`slice-sdk-csharp-projektgeruest`) hatte
`sdks/csharp/README.md` §Status bewusst geschrieben — mit Verweis auf genau
diesen Folge-Slice als den Zug, der die HTTP-Fläche real liefert
(`done/slice-sdk-csharp-projektgeruest.md` §6/§7, „Folge-Slices:
`slice-sdk-csharp-http-client-flaeche`, …"). Der Satz war bei Niederschrift
wahr („Client surfaces … are added by follow-up releases") und wurde durch
diese Arbeit — nicht durch eine Nachlässigkeit, sondern durch die korrekte
Auslieferung des angekündigten Umfangs — falsch: `PgChangeFeedHttpClient`
existiert jetzt real mit zehn Methoden, aber das README behauptete
weiterhin, die Fläche stünde noch aus. Der Implementer-eigene
§3.13-Suchlauf fand diesen Träger nicht (Plan-Nachzug-Abschnitt des Slice
nennt ihn nicht) — gefunden hat ihn erst der **Reviewer** (F-1, HIGH,
merge-blockierend), über eine gezielte Lektüre des Package-READMEs gegen
den realen Methodenbestand, nicht über einen Diff-Treffer (die Datei stand
nicht im Diff dieses Slice). Damit reiht sich dieser Fund in die bereits
belegte Unter-Klasse „gefunden vom Reviewer, nicht vom Implementer-eigenen
Suchlauf" ein (vgl. `slice-095`, `slice-097`) statt in die häufigere Form
„Implementer findet über den vorgeschriebenen Suchlauf selbst".

Behoben in der Fixrunde (`7bf7dece`): das README trägt jetzt den realen
Lieferstand („a full HTTP API client surface … the nine `SPEC-018`
capabilities plus `GET /changes`"), der gRPC-Stream bleibt korrekt als
„follow-up release" benannt (real noch offen). Fixrunden-Nachprüfung und
Verifikation bestätigen unabhängig voneinander (letztere über eine breitere
Wortliste als der Reviewer-Suchlauf), dass kein weiterer stehen gebliebener
Satz im README oder in benachbarten Dateien (`PgChangeFeedClientOptions.cs`,
`PgChangeFeedClientOptionsTests.cs`) verblieben ist — die dort verbliebenen
zwei Treffer sind Scope-Aussagen über eine andere Klasse, keine
Lieferstand-Behauptung, und bleiben zutreffend.

Quelle: `docs/reviews/review-slice-sdk-csharp-http-client-flaeche.md` <!-- d-check:status-provenance -->
F-1 (Klasse dort bereits selbst benannt: „`BEO-PGC/arbeit-ueberholt-
stehenden-traeger`-Familie … hier eine neue Datei-Instanz derselben
Fehlerklasse") und §Fixrunden-Nachprüfung F-1 · `docs/reviews/verifikation-
slice-sdk-csharp-http-client-flaeche.md` §3.2 · Commit `7bf7dece`.
