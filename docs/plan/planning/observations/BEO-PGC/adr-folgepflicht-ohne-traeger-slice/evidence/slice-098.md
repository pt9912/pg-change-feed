# Beleg: slice-098

Vorgang: `slice-098` — C#-Sprachwurzel und HTTP-Client.

Fund (Review zu `slice-098`, F-1):
`spec/pflichtenheft.md` `SPEC-023` Zeile *Sprachen und Umfang* trägt weiterhin
den `ADR-0087`-Wortlaut ("C# und Kotlin die **HTTP-Familie** … weitere
Oberflächen **perspektivisch**"), obwohl `ADR-0090` (Accepted, 2026-09-17, vor
Beginn dieses Slice) genau diese Beschränkung als aufzulösen benennt und die
Aktualisierung als eigene "Folgepflicht (Spec-Zug) … ohne ADR-/Slice-Kennung
im Spec-Text" führt. Weder `ADR-0090`s eigener Annahme-Commit (`72d026b`)
noch dieser Slice-Diff ändern die Zeile (`git diff 9b6b4fe..HEAD --
spec/pflichtenheft.md` ist leer).

**Warum kein Ausgang „eingetreten" gegen `BEO-PGC/arbeit-ueberholt-stehenden-traeger`:**
Jene Klasse verlangt, dass eine Arbeit eine beschriebene Eigenschaft bewegt
und dadurch einen zuvor wahren Satz falsch macht (`AGENTS.md` §3.13). Die
`SPEC-023`-Zeile war bereits **vor** `slice-098` falsch — seit `ADR-0090`s
Annahme, die diesem Slice vorausging und die dieser Slice-Diff nicht
berührt. Es handelt sich um eine andere Fehlerklasse: eine ADR benennt ihre
eigene Folgepflicht, ohne ihr einen Träger zuzuweisen, und die Pflicht bleibt
seither adresslos liegen.

**Ausgang bei dieser Closure: weiter offen.** `slice-098`s Plan hat
`SPEC-023`/die HTTP-API-Vertragsfläche ausdrücklich als Out-of-Scope geführt
(§1: „Eine Änderung an der HTTP-API selbst … wäre eine Spec-Änderung, kein
Beispiel-Umbau") — die fehlende Aktualisierung ist deshalb kein Versäumnis
dieses Slice, sondern eine seit `ADR-0090`s Annahme offene, adresslose
Folgepflicht. Folge-Hinweis für die Planung: `slice-103` (letzter Slice der
Matrix) benennt `SPEC-023` „Sprachen und Umfang" bereits als bedingten
Prüfpunkt in seinem eigenen §6; einer der dazwischenliegenden Slices
(`099`–`102`) kann die Zeile ebenso früher nachziehen, sobald ein Planner sie
einem konkreten Liefer-Punkt zuweist.

Quelle: Review zu `slice-098`, F-1 ·
`docs/plan/planning/in-progress/slice-098-csharp-sprachwurzel-http-client.md`
§6 (fünfte Risiko-Zeile) · `docs/plan/adr/0090-beispiel-clients-volle-matrix.md`
§Konsequenzen (Folgepflicht Spec-Zug) · `spec/pflichtenheft.md` `SPEC-023`
§Sprachen und Umfang ·
`docs/plan/planning/open/slice-103-grpc-client-kotlin.md` §6.
