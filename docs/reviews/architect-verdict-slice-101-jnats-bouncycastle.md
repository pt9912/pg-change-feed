# Architect-Verdikt — `io.nats:jnats`/`bouncycastle`-Trigger in slice-101

**Datum:** 2026-09-17 · **Stand:** `8ddb2e0` · **Rolle:** Architect (eigener Kontext; die
Sicherheits-Historie-Prüfung dieses Zugs ist selbst durchgeführt, nicht übernommen) ·
**Anlass:** `docs/reviews/review-slice-101.md` F-1 (HIGH) — Modul 8 §Konflikt-Pfad als
Rollen-Sequenz, Eskalation Reviewer → Architect.

## 1. Verdikt: **Fortsetzen ohne Rückführung** — der Trigger ist rückwirkend als nicht
   rückführungspflichtig geschlossen, der Prozessfehler bleibt separat benannt

Zwei Fragen, zwei getrennte Antworten (der Review trennt sie bereits korrekt in F-1/F-2):

1. **Trägt die Bewertung (Lizenz + Sicherheits-Historie) die Fortsetzung inhaltlich?** Ja —
   siehe §2/§3: beide Hälften jetzt vollständig und mit Beleg geprüft, Ergebnis unverändert
   gegenüber der Implementer-Bewertung.
2. **Überwiegt der Prozessfehler trotzdem?** Nein — siehe §4: eine formale Rückführung
   `in-progress → next` würde einen fertigen, getesteten Kotlin-Client verwerfen und exakt
   dieselbe Bewertung ein zweites Mal erzeugen, ohne dass sich am Ergebnis etwas änderte. Das
   ist der Fall, den Modul 8 §Konflikt-Pfad als *„Lockerung legitim, aber undokumentiert"*
   führt — nur ist der Gegenstand hier kein ADR-Widerspruch, sondern ein Plan-Trigger (§4 des
   Slice-Plans), und das nachziehende Artefakt ist dieses Verdikt-Dokument selbst, nicht eine
   Folge-ADR (siehe §4, Begründung der Abweichung von der ADR-Form).

## 2. Sicherheits-Historie — vollständig, mit Datum und Quelle nachgeprüft

**Geprüftes Objekt:** `org.bouncycastle:bcprov-lts8on:2.73.12.1` (transitiv über
`io.nats:jnats:2.26.3`, `compile`-Scope, unconditional — vom Reviewer bereits gegen
`jnats-2.26.3.pom`/`.module` bestätigt, hier nicht erneut gemessen).

**Geprüft am 2026-09-17, drei unabhängige Quellen, alle mit Netzzugriff real abgefragt:**

| Quelle | Abfrage | Ergebnis |
|---|---|---|
| OSV.dev, exakte Version | `POST /v1/query` mit `{"package":{"name":"org.bouncycastle:bcprov-lts8on","ecosystem":"Maven"},"version":"2.73.12.1"}` | `{}` — keine offene Vulnerability für genau diese Version |
| OSV.dev, ganzes Paket (ohne Version) | `POST /v1/query` mit nur `package` (keine Version) | 2 Treffer: `GHSA-4h8f-2wvx-gg5w`/`CVE-2024-34447` (bcprov-lts8on: `introduced 0`, `fixed 2.73.6`), `GHSA-mx76-r943-rf8g`/`CVE-2026-8149` (bcprov-lts8on: `introduced 2.73.0`, `fixed 2.73.11`) |
| GitHub Advisory Database | `GET /advisories?ecosystem=maven&affects=org.bouncycastle:bcprov-lts8on` | dieselben zwei GHSA-IDs, keine dritte |
| NVD, Keyword-Suche | `GET /rest/json/cves/2.0?keywordSearch=bcprov-lts8on` | 4 Treffer gesamt: `CVE-2024-34447` (= oben), `CVE-2026-8149` (= oben), `CVE-2025-9341` (LTS `2.73.0`–`2.73.7`, fixed danach), `CVE-2025-12194` (LTS `2.73.0`–`2.73.7`, fixed danach), `CVE-2026-15997` (LTS `< 2.73.12.1`, **fixed genau bei `2.73.12.1`**) |
| OSV.dev, Cross-Check der zwei zusätzlichen NVD-Funde | `GET /v1/vulns/CVE-2025-9341`, `CVE-2025-12194`, `CVE-2026-15997` | bestätigt dieselben Grenzen (Aliase `GHSA-jfcv-jv9g-2vx2`, `GHSA-jv6h-4262-q663`; `CVE-2026-15997` ohne GHSA-Alias, `fixed`-Commit `053e59f` bei genau `2.73.12.1`) |

**Fünf bekannte Advisories insgesamt** für `org.bouncycastle:bcprov-lts8on` (Stand
2026-09-17), keine davon offen für `2.73.12.1`:

| ID | Betroffener Bereich (bcprov-lts8on) | Fix-Version | `2.73.12.1` betroffen? |
|---|---|---|---|
| `CVE-2024-34447` / `GHSA-4h8f-2wvx-gg5w` | `< 2.73.6` | `2.73.6` | nein |
| `CVE-2025-9341` / `GHSA-jfcv-jv9g-2vx2` | `2.73.0`–`2.73.7` | `2.73.8`(+) | nein |
| `CVE-2025-12194` / `GHSA-jv6h-4262-q663` | `2.73.0`–`2.73.7` | `2.73.8`(+) | nein |
| `CVE-2026-8149` / `GHSA-mx76-r943-rf8g` | `2.73.0`–`2.73.10` | `2.73.11` | nein |
| `CVE-2026-15997` (kein GHSA-Alias) | `2.73.0`–`< 2.73.12.1` | **`2.73.12.1`** | nein — die gepinnte Version **ist** die Fix-Version |

`CVE-2026-15997` ist der schärfste Fund: Die gepinnte Version liegt nicht nur hinter dem
Fix, sie **ist** der Fix-Commit (`fixed: 2.73.12.1`, Out-of-Bounds-Write in der
ARM-SHA3/SHAKE-`restoreFullState`, nur relevant bei „memoable SHA3/SHAKE-States aus nicht
vertrauenswürdiger Quelle" — im NATS-Client-Anwendungsfall ohnehin nicht einschlägig).
Diese Zeile fehlt in der Reviewer-Prüfung (der Review nannte nur die zwei GHSA-Treffer aus
OSV.dev, nicht die NVD-Keyword-Suche) — deckt aber dieselbe Schlussfolgerung: **kein
offener Fund**, und zusätzlich der stärkere Beleg, dass der Pin exakt auf dem
Sicherheits-Fix sitzt, nicht nur zufällig danach.

**Lizenz-Hälfte** (vom Review bereits geprüft, hier gegengelesen): eigener Abruf von
`bouncycastle.org/licence.html` am 2026-09-17 bestätigt wörtlich den MIT-Text
(„Permission is hereby granted, free of charge…"). Keine Diskrepanz zum Review.

**Ergebnis:** Beide Hälften des §4-Triggers („Lizenz, Sicherheits-Historie") sind jetzt mit
Datum, Quelle und Gegenprobe über drei unabhängige Datenbanken belegt — nicht nur die vom
Implementer im Pin-Kommentar behauptete Pauschalaussage. Das löst F-2 des Reviews
(„Aussage ohne Ursprungsbeleg") auf: Der Beleg steht jetzt in diesem Dokument.

## 3. Trägt die Bewertung die Fortsetzung? Ja — Gegenprobe gegen den Trigger-Wortlaut

§4 des Slice-Plans: „wenn eine der beiden gepinnten Client-Bibliotheken transitiv eine
zweite, gepinnte Abhängigkeit erzwingt, die eine eigene Bewertung braucht (Lizenz,
Sicherheits-Historie) — dann ist die Größenannahme… falsch."

Der Trigger-Wortlaut sagt: **Erzwingt** eine Bewertung ⇒ Größenannahme falsch. Er sagt
nicht: **Jedes Ergebnis** einer erzwungenen Bewertung macht die Größenannahme falsch. Mit
der jetzt vollständigen Bewertung lässt sich das Gegenteil dessen zeigen, was der Trigger
befürchtet: Eine dritte gepinnte Abhängigkeit *ohne* offenen Sicherheits- oder
Lizenz-Befund fügt keine neue Bewertungs-*Kategorie* hinzu, die der Slice nicht schon trägt
(er pinnt bereits zwei Bibliotheken mit Lizenz-/Registry-Bewertung, `ADR-0090` Festlegung 4)
— sie verlängert nur die bestehende Prüfung um eine dritte Zeile derselben Form. Das ist
ein Argument für „die Prüftiefe war nötig und ist jetzt erbracht", nicht für „der Slice ist
zu groß".

## 4. Warum kein Rückgriff auf die Formstrenge — Prozessfehler vs. Substanz

**Der Prozessfehler ist real und wird nicht kleingeredet:** Der Implementer hat eine
vorab als Rückführungs-Grund deklarierte, real eingetretene Bedingung selbst für nicht
bindend erklärt, statt den Konflikt-Pfad auszulösen (Modul 8 §Rollen-Regeln — dieselbe
Grenze wie bei ADRs: „darf höchstens vorschlagen, niemals stillschweigend widersprechen").
Das ist kein Formfehler, den man mit einem nachgereichten Beleg wegwischt.

**Trotzdem keine formale Rückführung, aus einem Grund, der in der Sache liegt, nicht in der
Bequemlichkeit:** Eine Rückführung `in-progress → next` hieße hier: den bereits fertigen,
getesteten Kotlin-NATS-Client verwerfen, den Slice neu zerlegen, und in der Neuzerlegung
exakt dieselbe Frage stellen, die dieses Dokument gerade beantwortet hat — mit demselben
Ergebnis. Eine Rückführung, die mit Sicherheit zum selben Bau-Artefakt und derselben
Bibliotheks-Wahl zurückführt, korrigiert nichts Inhaltliches; sie wiederholt Aufwand, um
eine Prozess-Lücke zu schließen, die auch anders schließbar ist (siehe §5, §6). Das
unterscheidet diesen Fall von der `ADR-0069`/`-0070`-Historie
(`BEO-PGC/vorab-bedingung-nach-umsetzung-geprueft`, §6): Dort deckte die nachträgliche
Prüfung eine **echte Diskrepanz** auf (der Hook bildete die Prüfungen nicht
deckungsgleich ab) — hier deckt die nachträgliche Prüfung **keine** Diskrepanz auf, sondern
bestätigt exakt das, was der Implementer bereits (unvollständig belegt) behauptet hatte.

**Diese Architect-Bewertung trägt jetzt rückwirkend die im Plan §4 geforderte „eigene
Bewertung".** Die verlangte Prüftiefe — Lizenz **und** Sicherheits-Historie, mit Beleg —
ist mit §2/§3 real erreicht, auch wenn sie zeitlich nach der Implementierung kam. Der
Trigger 1 aus §4 gilt damit als **erfüllt und beantwortet**, nicht als offen oder
übersprungen: Er verlangte eine Bewertung, keine bestimmte Rollen-Reihenfolge für ihre
Erstellung. Der Trigger ist damit **rückwirkend als nicht rückführungspflichtig
geschlossen** — der Slice bleibt in `in-progress/`.

**Warum keine neue oder Folge-ADR:** `ADR-0090` trifft keine Aussage über eine konkrete
transitive Abhängigkeit einer konkreten Bibliotheksversion; es gibt nichts, dem dieses
Verdikt widerspricht oder das es ablöst. Modul 8 §Konflikt-Pfad-Tabelle nennt die ADR als
Übergabe-Artefakt für den Fall „Lockerung legitim, aber undokumentiert" — hier ist der
Gegenstand aber ein Plan-Trigger, nicht eine ADR-Festlegung; das nachziehende,
dokumentierende Artefakt *ist* dieses Verdikt-Dokument (Übergabe an die Planner-Closure,
§6).

## 5. Rollen-Klärung — Beobachtungs-Register geprüft, ein passender Eintrag existiert

`grep -ril "rueckfuehrung\|rückführung\|trigger" docs/plan/planning/observations/` und
händische Sichtung finden **einen** treffenden Bestandseintrag:

**`BEO-PGC/vorab-bedingung-nach-umsetzung-geprueft`** (Stand: offen, 1×, Beleg
`evidence/slice-073.md`). Die dort beschriebene Lücke deckt exakt dieses Muster:
„Eine §4-Bedingung, die die Klärung einer Frage vor die Umsetzung legt, wird erst bei der
Closure ausgewertet — geschrieben wird die Umsetzung vorher. Ist die Bedingung dann
eingetreten, steht der Slice zwischen zwei Ausgängen: der benannte Weg… ist nicht mehr
gangbar… der tatsächliche Weg ist der Konflikt-Pfad." `slice-101` ist eine zweite,
unabhängige Instanz derselben Klasse: Der Implementer schrieb das Artefakt, bevor die
§4-Bedingung ausgewertet war, die Bedingung trat ein, und der tatsächliche Weg war (nach
Reviewer-Eskalation) der Konflikt-Pfad — nicht der benannte `in-progress → next`.

**Empfehlung an die Planner-Closure (nicht Architect-Aufgabe, hier nur benannt):** Eine
weitere Beleg-Datei `evidence/slice-101.md` in diesem bestehenden Verzeichnis anlegen (kein
neuer Beobachtungs-Eintrag — die Bezeichnung ist bereits vergeben und passt). Das hebt den
abgeleiteten Zähler auf **2×** — weiterhin unter der 3×-Schwelle, der Stand bleibt „offen".
Diese Architect-Eskalation selbst ist **keine dritte, neue Beobachtung**: Sie ist derselbe
Befund wie `slice-073`, nur an einem anderen Slice gemessen.

## 6. Übergabe

An **Planner-Closure** von `slice-101`:

- Trigger 1 (§4) gilt als erfüllt **und** beantwortet — keine Rückführung nötig; dieses
  Dokument ist der Beleg, §7 des Slice-Plans zitiert es.
- F-2 des Reviews ist durch §2 dieses Dokuments aufgelöst — kein separater Implementer-Fix
  nötig, sofern die Closure-Notiz oder der Pin-Kommentar auf dieses Dokument verweist.
- `evidence/slice-101.md` in `BEO-PGC/vorab-bedingung-nach-umsetzung-geprueft/` anlegen
  (§5) — Zähler-Stand nach Anlage: 2×, Stand bleibt „offen".
- Risiken aus Slice-Plan §6 bleiben unberührt von diesem Verdikt; sie werden regulär bei
  Closure aufgelöst.

An **Reviewer**: F-1 ist mit diesem Verdikt geschlossen (Modul 8 §Konflikt-Pfad,
Fall „Lockerung legitim, aber undokumentiert", adaptiert für einen Plan-Trigger statt
einer ADR). Kein Re-Review nötig, sofern die Closure-Notiz auf dieses Dokument verweist.

## 7. Was nicht getan wurde

Kein Produktionscode, kein Slice-Plan, kein `git mv`, kein DoD-Häkchen geändert; keine ADR
geschrieben (§4 begründet, warum keine nötig ist). Alle Netzabfragen sind schreibfrei
(GET/POST-Query gegen öffentliche, lesende APIs); keine Artefakte außerhalb dieses
Dokuments angelegt.
