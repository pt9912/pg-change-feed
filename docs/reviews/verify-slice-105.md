# Verifikationsbericht: slice-105 — 2026-09-17

**Rolle:** Verifier (Modul 8/11) — „Bauen wir es richtig?" gegen den
DoD-Vertrag (`slice-105` §2, LP1–LP3) und die dort referenzierten
Entscheidungen. Keine ADR ist für die Baseline-Materialisierungs-Mechanik
selbst einschlägig (eigener Grep über `docs/plan/adr/*.md`, unten
bestätigt); geprüft wird [`ADR-0073`](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md)
(Zitat-Korrektur) auf die drei Link-Fixes sowie `AGENTS.md` §3.9
(Exit-Code-Disziplin) und §3.5 (ADR-Immutabilität, für den unberührten
`ADR-0051`-Beleg). **Nicht** gegen den Diff als solchen (Reviewer-Aufgabe,
[`review-slice-105.md`](review-slice-105.md)) und **nicht** gegen realen
Bedarf (Validator, nicht ausgelöst).

**Frischer Kontext.** Diese Sitzung hat den Slice-Plan vollständig gelesen
(§1–§8), alle vier Implementer-Commits (`47d86b2`, `94de7bb`, `9b66dc1`,
`9edad9b`) samt vollem Diff, den Review-Report (`5fe76ae`), `MR-002`
isoliert (ohne Rückgriff auf den Slice-Plan) und `harness/conventions.md`
im aktuellen Stand gelesen. Jede Zahl/jeder Befund in diesem Bericht stammt
aus einem hier selbst gefahrenen Lauf oder einer hier selbst
gelesenen/abgefragten Quelle — Review-Zahlen und -Befunde waren Kontext,
nicht übernommen (eigener Download, eigener Hash, eigener `diff -r`, eigene
Gate-Läufe, eigene Reproduktion von F-1).

**Gegenstand.** `47d86b2` (LP1) → `94de7bb` (LP2 + Träger-Nachzug +
`.d-check.yml`-Ausnahme) → `9b66dc1` (LP3, `MR-002`) → `9edad9b`
(DoD-Checkboxen) → `5fe76ae` (Review-Report, 0 HIGH, 1 MEDIUM, 1 INFO). Der
Slice liegt in `in-progress/`; der `git mv` nach `done/` ist **nicht**
erfolgt — erwartungsgemäß (Closure-Notiz, Beobachtungs-Register-Eintrag,
Risiko-Ausgänge und die drei Paarungen stehen laut DoD noch aus). `git
status` war vor dieser Sitzung sauber; alle in dieser Sitzung erzeugten
Artefakte (Download, Entpackung, testweise Einfügungen für die F-1-
Reproduktion) liegen im Scratchpad bzw. wurden vollständig zurückgenommen
(`git checkout --` nach jeder Reproduktionsprobe, `git status --short`
danach leer bestätigt).

---

## 1. Eigene Messungen dieses Laufs

Jeder Lauf ungepiped, Exit-Code direkt aus einem eigenen, abgeschlossenen
Schritt gelesen (`AGENTS.md` §3.9).

| # | Lauf | Exit | Ergebnis |
|---|---|---|---|
| 1 | `curl` gegen das reale `v6.9.0`-Release-Asset (`lab-regelwerk.zip`), eigener `sha256sum` | — | `8a4e0aaf597a9c67404cb7a350a6fba992f0c98195e011e025c073660ee55cce` — **identisch** zur Implementer-/Reviewer-Behauptung, unabhängig ermittelt |
| 2 | `unzip` des selbst geladenen Assets, `diff -r` gegen den committeten `.harness/baseline/v6.9.0/{regelwerk,templates}/`-Baum | 0 (diff) | **byte-identisch**, keine Abweichung |
| 3 | `ls -d .harness/baseline/*/` | — | genau **ein** Verzeichnis (`v6.9.0/`) — `v6.5.0/` real entfernt |
| 4 | `make baseline-verify` | **0** | „v6.9.0 OK — 54 Dateien" |
| 5 | `wc -l < .harness/baseline/v6.9.0/SHA256SUMS` | — | 54 Zeilen, `find … -type f \| wc -l` (inkl. `SHA256SUMS` selbst) → 55 — konsistent |
| 6 | `grep -rn "v6\.5\.0" . --exclude-dir=.git` (repo-weit, eigener Lauf) | — | Treffer ausschließlich in: `docs/plan/planning/in-progress/slice-105-*.md` (eigene Plan-Datei, `d-check:ignore`-annotiert, historische Zielstand-Erwähnung), `docs/reviews/**` (jetzt exempt), `docs/plan/planning/done/**` (Records: `welle-18-results.md`, `welle-19-results.md`), `docs/plan/adr/0051-*.md` (P8-Zeile, s. u.) — **kein** lebender Treffer außerhalb dieser vier Klassen |
| 7 | `git show 94de7bb -- docs/reviews/architect-verdict-aufschub-adresse-verfaellt.md docs/reviews/review-slice-041-fixrunde-2.md` | — | beide Diffs ändern **ausschließlich** das Versions-Pfadsegment (`v6.5.0`→`v6.9.0`); Referent (dieselbe Regelwerk-/Templates-Datei bzw. derselbe Modul-6-Anker) unverändert |
| 8 | Ziel-Anker-Check: `grep -n "^### Das Beobachtungs-Register" .harness/baseline/v6.9.0/regelwerk/modul-06-roadmap.md` | — | Zeile 79, Anker `#das-beobachtungs-register-modul-6` löst auf (zusätzlich indirekt bestätigt: `docs-check` liefe sonst `anchor-missing`, ist aber grün) |
| 9 | Eigene F-1-Reproduktion (a): stale `v6.5.0`-Pfad in `docs/reviews/verify-slice-104.md` eingefügt, `make docs-check`, danach `git checkout --` | 0 | **grün trotz stale Pin** — `versions`-Modul schlägt nicht an (861 Dateien, 0 Befunde) |
| 10 | Eigene F-1-Reproduktion (b): derselbe stale Pfad in `harness/README.md` (außerhalb der Ausnahme) eingefügt, `make docs-check`, danach `git checkout --` | **2** | `harness/README.md:224 v6.5.0 version-stale` — bestätigt: außerhalb der Ausnahme färbt derselbe Pin korrekt rot |
| 11 | `RANGE=3c5e18a..HEAD make commit-traceability` | **0** | 5 Commits (inkl. Review-Commit), 0 Befunde, Betreffs ohne Struktur-ID |
| 12 | `make gates` (auf `HEAD` = `5fe76ae`), Log in Datei umgeleitet, Exit separat gelesen | **0** | alle sechs Gates grün (`baseline-verify`, `docs-check`, `a-check`, `commit-traceability`, `coverage-gate`, `generated-sync`) |
| 13 | `grep -lI "Baseline" docs/plan/adr/*.md` + gezielte Suche nach „Baseline-Mechanik"/„baseline-verify" in ADR-Texten | — | einziger Treffer: `ADR-0051` P8-Zeile (Pin-Inventur-Momentaufnahme), keine ADR zur Baseline-*Mechanik* selbst |
| 14 | `grep -rli "baseline\|v6\." docs/plan/planning/observations/` | — | Treffer betreffen andere Themen (z. B. `commit-traceability-kein-vorab-hook`); kein Register-Eintrag zu Baseline-Version/Namenskonvention — bestätigt Plan §8 eigenständig |

---

## 2. LP1 — Baseline `v6.9.0` netzlos materialisiert, `v6.5.0` entfernt, genau ein Verzeichnis

**Erfüllt, eigenständig gemessen — nicht nur übernommen.** Der SHA256 des
selbst geladenen Release-Assets (Lauf 1) ist identisch zur Implementer-/
Reviewer-Behauptung; der entpackte Inhalt ist byte-identisch zum
committeten Baum (Lauf 2, `diff -r` über beide Unterbäume). `make
baseline-verify` läuft eigenständig grün gegen genau ein `<tag>`-Verzeichnis
(Lauf 3, 4). Die im Plan dokumentierte Erkenntnis — dass der Sensor
`$base/*/` zählt und deshalb **keine** zwei parallelen Stände (auch nicht
unter `.harness/baseline/done/`) verträgt, weshalb `v6.5.0/` per `git rm -r`
statt `git mv` entfernt wurde — ist durch Lauf 3 bestätigt: `git log`/`git
show 47d86b2^:.harness/baseline/v6.5.0/...` trägt den alten Stand weiterhin
vollständig, kein Datenverlust.

## 3. LP2 — `harness/conventions.md` §Baseline aktuell, keine lebende `v6.5.0`-Erwähnung

**Erfüllt, eigenständig gemessen.** `harness/conventions.md` §Baseline zeigt
`Stand: v6.9.0`, `Datum der Adoption: 2026-09-17`. Der eigene repo-weite
`grep` (Lauf 6) findet **keine** lebenden Treffer außerhalb der vier
erwarteten Klassen (eigene Plan-Datei mit `d-check:ignore`, `docs/reviews/**`
jetzt exempt, `done/`-Records, `ADR-0051` P8-Pininventur). Die über den
DoD-Wortlaut hinausgehende, real gefundene Zusatzarbeit — drei per Entfernen
von `v6.5.0/` `target-missing` gewordene Links in `docs/reviews/**` per
Zitat-Korrektur (`ADR-0073`) korrigiert, alle übrigen reinen
Text-Erwähnungen dort über eine neue `.d-check.yml`-`versions.exempt-paths`-
Zeile ausgenommen — ist real vollzogen (Lauf 7) und die Ausnahme selbst
funktioniert wie behauptet (Lauf 9/10, siehe §5 unten für das eigene Urteil
zur *Breite* dieser Ausnahme). `ADR-0051` §Kontext (P8-Zeile) bleibt
unverändert — korrekt: reine historische Pin-Inventur-Momentaufnahme zum
`Accepted`-Zeitpunkt, kein Zitat-Fehler-Anlass (`AGENTS.md` §3.5).

## 4. LP3 — `MR-002` inhaltlich korrekt, in der Aktive-Adaptionen-Tabelle verlinkt

**Erfüllt.** `harness/conventions.md` §Aktive Adaptionen trägt die Zeile
`MR-002 <a id="mr-002"></a> | [Slice-/Welle-Kennungen sind Namen, nicht
Nummern …]`. Die Datei `harness/conventions/MR-002-slice-welle-kennungen-
sind-namen.md` trägt alle Pflichtfelder (Datum, Geltungsbereich,
Ersetzt-Baseline-Regel mit auflösendem Anker, Adaption, Begründung,
Auflösungs-Trigger). Der Anker `#vergabe-woher-die-nächste-kennung-kommt`
löst gegen die materialisierte `v6.9.0`-Datei auf (eigen gelesen, §Vergabe:
woher die nächste Kennung kommt, Zeile 331). Die beiden wörtlichen Zitate
(„Welle- und Slice-Kennungen sind Namen, nicht Nummern — unabhängig von der
Schreiberzahl" und „Schreiber ist, was committet: ein Mensch, ein Agent,
ein Automat") stimmen zeichengenau mit der materialisierten Quelle überein
(eigen gegengelesen). Zum inhaltlichen Selbsttrage-Urteil siehe §6 unten.

---

## 5. Eigenes Urteil zu F-1 — `.d-check.yml`-Ausnahme für `docs/reviews/**`

**Ich teile die MEDIUM-Einstufung und die Einschätzung „nicht
fix-pflichtig" — mit eigener, unabhängiger Begründung, nicht nur
Zustimmung.**

Der Befund selbst ist real reproduziert (Lauf 9/10): Die Ausnahme schützt
**ab dem ersten Commit**, nicht erst ab einem Abschluss-Ereignis wie beim
Vorbild `harness/conventions/done/**` (dort beginnt der Schutz erst mit dem
`git mv`, während eine aktive Adaptions-Datei in Bearbeitung weiterhin
geprüft wird). Das ist eine reale Asymmetrie und kein Wahrnehmungsfehler
des Reviewers.

**Warum ich es dennoch nicht auf HIGH hebe:** Der Maßstab für HIGH ist eine
Verletzung einer Hard Rule oder eines ADRs, oder ein Fund mit realem
Schadenspotential für eine *künftige Entscheidung*. Beides trifft hier
nicht zu, aus einem Grund, der über die Review-Begründung hinausgeht:
`docs/reviews/**`-Dateien sind laut Modul 10 explizit „Lauf-Beleg[e] …
über Läufe hinweg nicht wieder gelesen" — sie werden nicht als kanonische
Quelle konsultiert, aus der ein späterer Agentenlauf eine Handlungsanweisung
ableitet (anders als `harness/README.md`, `AGENTS.md` oder
`harness/conventions.md`, wo ein stiller stale Pin tatsächlich eine falsche
künftige Entscheidung stützen könnte). Ein unbemerkter Tippfehler-Pin in
einem Report bleibt damit ein **kosmetischer** Fehler in einem Dokument,
dessen einzige nachgelagerte Lese-Operation die Summary-Zeile für den
Beobachtungs-Register-Zähler ist (Modul 5 §Closure- und
Lerneintrag-Regeln) — nicht der Versions-Pin im Fließtext. Das
Schadenspotential ist damit strukturell begrenzt, nicht bloß unwahrscheinlich.

**Was ich zusätzlich zum Reviewer anmerke:** Die vom Reviewer offen
gelassene Frage — engerer Zuschnitt der Ausnahme vs. benannte Akzeptanz des
Trade-offs — sollte bei der Slice-Closure nicht nur „zu entscheiden" bleiben,
sondern tatsächlich entschieden und in §7 der Closure-Notiz als
Steering-Loop-Kandidat (geschärfte Regel: Ausnahmen mit Reifegrenze statt
Freibrief-ab-Anlage) benannt werden — sonst wiederholt sich dieselbe
Asymmetrie beim nächsten Baseline-Sprung unreflektiert. Das ist eine
Empfehlung an den Planner, kein DoD-Blocker dieses Slice.

**Fazit F-1:** MEDIUM korrekt, nicht fix-pflichtig korrekt — mit der
Präzisierung, dass die Begrenzung des Schadens aus der Lese-Semantik von
`docs/reviews/**` selbst folgt, nicht nur aus der Seltenheit eines
Tippfehlers.

**F-2 (INFO, vorbestehender kaputter Anker in `MR-001`):** eigenständig
nachvollzogen — der Anker
`grundlagen-source-precedence.md#spec-straten-mehr-als-ein-spec-dokument`
existiert real nur in `grundlagen-referenz-richtung.md`, gegen `v6.9.0`
**und** (per `git show 3c5e18a:…`) bereits gegen `v6.5.0` bestätigt
vorbestehend. Kein Handlungsbedarf in diesem Slice, korrekt außerhalb des
Diff-Umfangs eingeordnet.

---

## 6. `MR-002` isoliert gelesen — ist sie selbsttragend?

Gelesen ohne Rückgriff auf `slice-105` oder die Session-Historie, nur die
Datei selbst plus (wie von der Datei selbst verlinkt) die
Baseline-Quellstelle.

**Trägt weitgehend.** Ein künftiger Planner-Lauf, der `slice-106` (oder
später) anlegt, findet in `MR-002` alles Nötige: Geltungsbereich
(„alle nach `slice-105` neu angelegten Slice- und Welle-Plan-Dateien"),
die operative Regel („Name trägt das Präfix eines vorhandenen Ankers
(`LH-*`, `ADR-*`, `CO-*`), wenn einer existiert, sonst ein freier Slug"),
den Bestandsschutz-Cutoff (`slice-001`–`slice-105` bleiben nummeriert) und
einen auflösenden Link auf die vollständige Baseline-Regel für Details, die
`MR-002` selbst nicht wörtlich wiederholt (z. B. den
Welle-Namensraum-Präfix für einen während einer Welle entstehenden Slice,
`slice-<welle-name>-<aspekt-slug>` — dieser Fall ist nicht in `MR-002`s
eigenem Text, aber über den Link erreichbar, und das Repo läuft ohnehin
wellenlos).

**Eine reale, aber geringfügige Lücke:** Der operative Absatz („Adaption")
spricht ausschließlich von **Slice**-Bestandsschutz und -Namensform; eine
symmetrische Aussage für **Welle**-Kennungen fehlt im Fließtext, obwohl
Titel und Geltungsbereich beide Klassen nennen. Das ist bei genauem
Nachdenken auflösbar (es gibt keine vorbestehende `welle-NN`-Nummerierung in
diesem Repo, die geschützt werden müsste — der Bestandsschutz-Absatz für
Wellen wäre leer, weil es nichts zu schützen gibt), aber ein Planner-Lauf,
der zum ersten Mal überhaupt eine Welle anlegt, muss diese Schlussfolgerung
selbst ziehen statt sie explizit vorzufinden. Das ist kein Blocker — die
Konsequenz eines Fehllesens wäre höchstens eine nummerierte statt einer
benannten ersten Welle, korrigierbar per Review — aber es ist die einzige
Stelle, an der „isoliert lesbar" nicht ganz „isoliert vollständig" bedeutet.
Empfehlung: bei Gelegenheit (kein eigener Slice nötig) einen Satz ergänzen,
der den Welle-Fall explizit als „keine Altbestand-Ausnahme, volle
Baseline-Regel ab der ersten Welle" benennt.

**Insgesamt:** `MR-002` ist **selbsttragend genug** für den Regelfall
(nächster Slice), mit einer benannten, nicht blockierenden Lücke für den
noch nie eingetretenen Wellen-Fall.

---

## 7. §6-Risiken — Ausgänge (eigenes Urteil)

| Risiko | Ausgang | Beleg |
|---|---|---|
| 1. `SHA256SUMS`-Prüfung des Release-Assets schlägt fehl | **entfallen** — nicht eingetreten | Lauf 1 (eigener Download, identischer Hash) |
| 2. `make baseline-verify` verträgt strukturell keinen Versionswechsel | **entfallen** — nicht eingetreten | Lauf 3, 4 (generischer Glob, grün gegen `v6.9.0`) |
| 3. Freshness-Audit findet zusätzliche relevante Wellen, die die ursprüngliche Recherche übersah | **entfallen** — kein weiterer Kategorienfehler gefunden; §7 „Zur Kenntnis" listet 132/133/134/136/137 mit Begründung „kein Handlungsbedarf", eigener Gegen-Grep (Lauf 14) bestätigt keinen Register-Treffer, der dem widerspräche | Plan §7, Lauf 14 |

**Zusätzliches, nicht vorab antizipiertes Risiko (vom Implementer selbst
gefunden):** „Entfernen von `v6.5.0/` macht Links/Text-Erwähnungen in
Records stale." Dieses Risiko stand in keiner der drei §6-Zeilen — es ist
eine reale Lücke im ursprünglichen Risiko-Katalog, nicht nur eine
Unterkategorie von Risiko 3 (jenes betraf Kategorienfehler der
*Migrations-Recherche* selbst, nicht Referenz-Bruch durch die *eigene*
Löschung). **Es ist eingetreten** (drei reale `target-missing`-Links,
gefunden und per `ADR-0073`-Zitat-Korrektur behoben) und **korrekt
mitigiert** — kein Schaden verblieben (eigener repo-weiter Gegen-Grep,
Lauf 6, findet keinen verbliebenen kaputten Link; `make gates` grün
schließt `target-missing` mit ein). Die gewählte Mitigationsform (die
`.d-check.yml`-Ausnahme für die verbleibenden reinen Text-Erwähnungen) ist
der Gegenstand von F-1 oben — die Mitigation selbst war nötig und wirksam,
ihre *Form* ist die dokumentierte Abwägung.

**Empfehlung für die Closure-Notiz:** Dieses vierte, unerwartete Risiko
gehört als eigener Punkt in §7 (Closure-Notiz) benannt — nicht rückwirkend
in §6 eingefügt (§6 ist Vorab-Benennung, die Nachträglichkeit wäre
unehrlich), sondern als Fund dieses Umsetzungslaufs. Kandidat für einen
Beobachtungs-Register-Eintrag, falls ein künftiger Baseline-Sprung dasselbe
Muster zeigt (Erst-Auftreten jetzt, Zähler bei 1).

Keines der drei geplanten Risiken ist schädlich eingetreten; das
unerwartete vierte Risiko ist eingetreten, aber folgenlos mitigiert.

---

## 8. Verdikt

**DoD-konform für die drei geprüften Liefer-Punkte (LP1, LP2, LP3) —
Closure-Arbeit steht wie erwartet noch aus.** Alle drei eigenständig
gemessen, keine Diskrepanz zur Implementer-/Reviewer-Behauptung gefunden.
`make gates` läuft eigenständig grün (Lauf 12, ungepiped, Exit-Code direkt
gelesen). Keine ADR ist für die Baseline-Mechanik selbst einschlägig
(bestätigt, Lauf 13); `ADR-0073` ist auf die drei Link-Fixes korrekt
angewendet (Referent unverändert, reine Pfadsegment-Korrektur).

**F-1:** MEDIUM, nicht fix-pflichtig — eigenes Urteil teilt die
Einstufung, mit eigener Begründung über die Lese-Semantik von
`docs/reviews/**` (§5 oben). **F-2:** INFO, vorbestehend, kein
Handlungsbedarf — bestätigt.

**`MR-002`:** selbsttragend für den Regelfall (nächster Slice), mit einer
benannten, nicht blockierenden Lücke für den bisher nie eingetretenen
Wellen-Fall (§6 oben).

**§6-Risiken:** alle drei geplanten Risiken *entfallen*; ein viertes,
unerwartetes Risiko ist eingetreten und korrekt mitigiert — Empfehlung, es
explizit in die Closure-Notiz aufzunehmen (siehe §7 oben).

**Für die Closure noch offen (Planner-Arbeit, DoD-Zeilen bereits als offen
markiert):**

1. Closure-Notiz §7 füllen — inkl. des vierten, unerwarteten Risikos und
   des F-1-Steering-Loop-Kandidaten (engerer Ausnahme-Zuschnitt oder
   benannte Akzeptanz).
2. Beobachtungs-Register: keine Beobachtung mit 3× erreicht (eigener
   Vorab-Check, Lauf 14); neuer Eintrag für das vierte Risiko ist eine
   Planner-Ermessensfrage (Erstauftreten, Zähler 1), kein Zwang.
3. §6-Risiko-Ausgänge gemäß §7 dieses Berichts eintragen (alle drei
   *entfallen*).
4. Die drei Paarungen (Anker · Folge-Slice · Register) — Anker-Paarung
   betrifft hier `seit slice-105`, falls das F-1-Trade-off oder das vierte
   Risiko verkörpert wird; Folge-Slice-Paarung: die beiden in §1
   angekündigten Folge-Slices (Gate-Index-Konsolidierung, DoD-Rot-vor-
   Grün-Regel) existieren noch nicht als Datei — das ist zulässig, da sie
   „noch zu vergebende Namen" ausdrücklich offenlassen; kein Blocker.
5. `git mv` nach `done/` erst nach 1–4.

**Nicht Gegenstand dieser Verifikation:** Validierung gegen realen Bedarf
(Validator, hier nicht ausgelöst).
