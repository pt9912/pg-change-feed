# Review-Report: slice-sdk-python-grpc-client-flaeche — Fixrunde — 2026-09-23

**Review-Art:** Code-Review — **Fixrunden-Re-Review** — geprüft gegen die
Befunde der Haupt-Review
([F-1…F-6](review-slice-sdk-python-grpc-client-flaeche.md)), die
Verifikations-Auflagen
([V-1…V-3](verifikation-slice-sdk-python-grpc-client-flaeche.md)),
[`ADR-0110`](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
(Accepted, Folgepflicht 2/3), [`SPEC-020`](../../spec/pflichtenheft.md),
[`LH-FA-SST-009`](../../spec/lastenheft.md), `AGENTS.md` §3
(§3.7/§3.11/§3.12/§3.13), `.harness/skills/reviewer.md`. **Nicht geprüft
gegen die DoD-Substanz** — das bleibt Verifier-Aufgabe (Modul 11); der
Verifier-Lauf liegt vor und ist hier nur als aufgelöste Auflagen-Kette
Eingang.

**Gegenstand:** `git diff b7a993fe~1..59e42d22` — zwei Commits:
`b7a993fe` (Fixrunde F-1…F-4/F-6), `59e42d22` (Verifikations-Report +
Auflagen V-1…V-3). **Während dieses Laufs gelandet und mitgeprüft:**
`5d5fa3eb` (docs(plan): id-unlinked-Befund im Fixrunden-Nachzug behoben) —
ändert dieselbe Nachzug-Tabelle (Version-Label 1.39 → 1.40 + [`ADR-0107`](../plan/adr/0107-python-pypi-zweites-sdk-package.md)-Link,
aus einem `d-check`-Lauf); sein Inhalt ist in die Befunde eingerechnet, der
Review-Range selbst bleibt der genannte.

**Skill:** `.harness/skills/reviewer.md` (Stand der Schärfung 2026-09-09,
inkl. §3.12-Instanz-A-HIGH, „Beleg trägt seinen Satz nicht",
„Zusage ohne Bindung an ihre Eingabeseite") · <!-- d-check:ignore (Adopter-spezifischer Skill-Pfad, existiert im Ziel-Repo ggf. nicht) -->
**Modell:** glm-5.3-flash (Claude-Agent-SDK, Reviewer-Rolle) · **Datum:**
2026-09-23

**Eingangs-Kontext** (Verträge, gegen die geprüft wurde):

- [Slice-Plan](../plan/planning/in-progress/slice-sdk-python-grpc-client-flaeche.md)
  (§2 DoD, §3 Plan + Plan-Nachzug)
- [Haupt-Review F-1…F-6](review-slice-sdk-python-grpc-client-flaeche.md),
  [Verifikation V-1…V-3](verifikation-slice-sdk-python-grpc-client-flaeche.md)
- [`ADR-0110`](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
  (Folgepflicht 2: `SPEC-027`/`LH-FA-SST-009.a` gebündelt; Folgepflicht 3:
  Handbuch-SDK-Hinweis je Oberfläche im selben Zug),
  [`ADR-0107`](../plan/adr/0107-python-pypi-zweites-sdk-package.md)
- [`proto/cdc/stream/v1/changestream.proto`](../../proto/cdc/stream/v1/changestream.proto)
  (`message Change` — Feldtypen der zehn Felder)
- `AGENTS.md` §3.7/§3.11/§3.12/§3.13, MR-000 (ID-Schema)

---

## Findings

<!-- Kein Fließtext, kein Lösungsvorschlag im Befund. -->

### F-1 — Fixrunden-Regression: drei String-Identitätsfelder der Feldvollständigkeits-Prüfung nur noch typgebunden

- `kategorie`: MEDIUM
- `quelle`: Beleg trägt seinen Satz nicht
  (`BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht` — dieselbe Klasse wie
  Haupt-Review F-2, schmaler); Verwandtschaft: „Zusage ohne Bindung an ihre
  Eingabeseite"
- `pfad`: `sdks/python/pgchangefeed/integration/test_grpc_realserver.py:59-77`
- `befund`: Der Testkommentar behauptet „Feldvollstaendigkeit am realen
  Wire-Image (`SPEC-020`): **alle zehn Felder getragen**"; für
  `transaction_id`, `source_table_id` und `schema_version` trägt die neue
  `isinstance`-Assertion (Zeilen 64, 65, 70) nur den **Typ** — ein Wire-Image
  ohne diese Felder (proto3 non-optional: ungesetzt = `""`, Feld nicht
  serialisiert) bleibt grün; die rot machende Mutation für die
  „getragen"-Aussage ist für diese drei Felder nicht benennbar, und der
  Runner bindet sie nicht (kein Treffer für
  `transaction_id|source_table_id|schema_version` im
  `run-sdk-python-integration-tests.sh`). Vor der Fixrunde war genau diese
  Bindung korrekt vorhanden (Haupt-Review F-2 nannte die sieben
  String-Identitäten ausdrücklich als korrekt gebunden); die typbewusste
  Neuschreibung hat sie für diese drei Felder **eingetauscht** statt
  ergänzt. Abgrenzung: `sequence` (int) ist typgebunden — ohne
  Domänenannahme die stärkste bindbare Form, kein Befund; Typ-Änderung und
  Feldentfall (`AttributeError`) färben alle zehn Felder rot, der MEDIUM
  betrifft nur die leere-Wert-Lücke der drei String-Identitäten
  (`change_id` bleibt über `!= ""` plus SQL-Gegenprüfung gebunden).
- `verifizierbar`: ja — Mutationsprobe: Server sendet ein `Change` ohne
  Feld 2 (oder mit leerem `transaction_id`); alle zehn Asserts des
  Integrationstests bleiben grün. Die Probe ist die Mutation selbst (kein
  Mutations-Harness im Repo)
- `klasse`: „Beleg trägt seinen Satz nicht" (Teilvakuum der
  Feldvollständigkeits-Assertion, 3 von 10 Feldern)

### F-2 — Handbuch: gRPC-Abschnitt trägt `**SDK:**`-Absätze für C#/Kotlin, keinen für Python; der 1.40-Historien-Satz nennt eine Klasse, die der Handbuch-Körper nirgends nennt

- `kategorie`: LOW
- `quelle`: Maintainability (Muster-Abweichung im eigenen Handbuch-Muster);
  Teil-Form von „Beleg trägt seinen Satz nicht"; [`ADR-0110`](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md) Folgepflicht 3 im
  Kern erfüllt — deshalb LOW
- `pfad`: `docs/user/benutzerhandbuch.md:682-685` (Python-Pointer),
  `docs/user/benutzerhandbuch.md:706-796` (gRPC-Abschnitt, SDK-Block nur
  C#/Kotlin), `docs/user/benutzerhandbuch.md:1243` (1.40-Historien-Zeile)
- `befund`: Die F-3-Korrektur trägt die Python-gRPC-Präsenz nur als
  Zeigerform aus dem HTTP-Absatz („derselbe Package trägt außerdem den
  gRPC-Change-Stream (siehe unten …)"); der Ziel-Absatz (Zeile 706 ff.)
  trägt `**SDK:**`-Absätze für C# (1.34) und Kotlin (1.37), keinen für
  Python, und nennt `PgChangeFeedGrpcClient` nur für C#/Kotlin (Zeilen
  772, 782). Die 1.40-Historien-Zeile behauptet „trägt jetzt die
  gRPC-Stream-Fläche `PgChangeFeedGrpcClient`" — die genannte Klasse steht
  im Handbuch-Körper an keiner Stelle (nur der Oberflächen-Pointer steht).
  `grep -n "PgChangeFeedGrpcClient" docs/user/benutzerhandbuch.md` liefert
  nur C#- (772), Kotlin- (782) und die 1.34/1.40-Historien-Zeilen.
- `verifizierbar`: ja — obiger `grep` (beide Stände: Parent `59e42d22` und
  HEAD identisch)
- `klasse`: „Beleg trägt seinen Satz nicht" (Historien-Satz über den
  Absatz-Inhalt) / Muster-Abweichung (in-section-Absatz vs. Pointer)

### F-3 — DoD-Zeile 6 des Plans durch die eigene F-3-Korrektur teilweise falsch geworden („Doku-Update bewusst nicht in diesem Slice")

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.13 („Arbeit überholt stehenden Träger", Planform)
- `pfad`: `docs/plan/planning/in-progress/slice-sdk-python-grpc-client-flaeche.md:121-123`
  gegen `docs/plan/planning/in-progress/slice-sdk-python-grpc-client-flaeche.md:170`
- `befund`: Die DoD-Zeile „Doku-Update … bewusst **nicht** in diesem Slice —
  gebündelt im letzten Flächen-Slice" ist durch die Fixrunde für den
  gRPC-Handbuch-Teil falsch geworden — die Nachzug-Zeile 170 desselben
  Dokuments deklariert die Handbuch-Korrektur (Version 1.40 + Historie)
  als im selben Slice erledigt; beide Stellen widersprechen sich, ohne dass
  die DoD-Zeile auf die Ausnahme verweist (die Bündelung bleibt nur für
  `spec/pflichtenheft.md` und den SSE-/NATS-Handbuch-Teil wahr).
- `verifizierbar`: nein — reine Plan-Konsistenzprobe (beide Zeilen im
  selben Dokument gelesen)
- `klasse`: „Arbeit überholt stehenden Träger" (DoD-Zeile vs. Nachzug-Zeile)

### F-4 — Handbuch-`Stand`-Feld hinter dem Historien-Stand

- `kategorie`: LOW
- `quelle`: Maintainability (Zustandsfeld, §3.7-Disziplin)
- `pfad`: `docs/user/benutzerhandbuch.md:5` gegen
  `docs/user/benutzerhandbuch.md:1243`
- `befund`: Der Kopf trägt weiter `Stand: 2026-09-22`, während die neueste
  Änderungshistorie-Zeile auf 2026-09-23 datiert; `d86d1965` etablierte die
  Konvention (Version **und** Stand gemeinsam gehoben, 1.38/09-22), die
  Fixrunde hob nur die Version. Kein `d-check`-Befund (kein gebrochener
  Verweis), aber ein Zustandsfeld, das seinen Ist-Stand um einen Tag
  hinterläuft.
- `verifizierbar`: ja — `sed -n '5p' docs/user/benutzerhandbuch.md` gegen
  `git log -S"Stand: 2026-09-22" --oneline -- docs/user/benutzerhandbuch.md`
- `klasse`: „Arbeit überholt stehenden Träger" (Zustandsfeld)

### F-5 — Zahl-/Stand-Drift-Klasse: drittes Auftreten in dieser Korrektur-Kette (am Range-Ende offen, in `5d5fa3eb` gelöst)

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.12 Instanz A — Zahl im Träger driftet gegen die
  Messung (Herkunft der Klasse:
  `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`)
- `pfad`: `docs/plan/planning/in-progress/slice-sdk-python-grpc-client-flaeche.md:170`
  (Stand am Range-Ende `59e42d22`)
- `befund`: Am Ende des geprüften Bereichs trug die Plan-Nachzug-Zeile
  „Version **1.39**" gegen den Ist-Stand des Handbuchs (1.40) — derselbe
  Klassentyp wie F-1/V-1, beim dritten Nachziehen wieder aufgetreten (5/29 →
  6/29 → 1.39/1.40). `5d5fa3eb` (während dieses Laufs gelandet, aus einem
  `d-check`-id-unlinked-Befund) hob die Zeile auf 1.40 und trug die
  V-2-Nachzug-Notiz — am Ist-Baum (HEAD `5d5fa3eb`) **gelöst**, kein offener
  Befund; die Wiederholung zählt für den Steering-Loop (drittes Auftreten
  derselben Klasse im selben Träger-Paar Plan↔Handbuch).
- `verifizierbar`: ja — `grep -n "Version 1\."
  docs/plan/planning/in-progress/slice-sdk-python-grpc-client-flaeche.md`
  (HEAD: 1.40) gegen `git show 59e42d22:<pfad>` (1.39)
- `klasse`: „Zahl im Träger ohne Ursprung — oder gegen die Messung driftend"

### F-6 — Verifikations-Report und Auflagen-Korrekturen in einem Commit; die Report-Auflagen-Liste liest sich als offen, ist im selben Commit angewendet

- `kategorie`: INFO
- `quelle`: Maintainability (Übergabe-Artefakt-Zuordnung, Modul 8)
- `pfad`: `59e42d22` (Commit), `docs/reviews/verifikation-slice-sdk-python-grpc-client-flaeche.md` §6
- `befund`: `59e42d22` trägt den Verifikations-Report **und** die V-1…V-3
  Korrekturen in einem Commit; der Report-§6 listet die drei Auflagen als
  „vor Closure" erforderlich, während derselbe Baum sie bereits anwendet —
  die Rollen-Zuordnung der Korrektur (Implementer vs. Verifier) ist aus dem
  Commit allein nicht ablesbar. Die Zuordnung tragen die Plan-Nachzug-Notizen
  („Verifikation V-3 …", „V-2-Nachzug …") — die Übergabe-Artefakte sind
  damit intakt; dieser Hinweis erwartet keine Aktion.
- `verifizierbar`: nein — Prozessbeobachtung
- `klasse`: Übergabe-Artefakt in Sammel-Commit

### F-7 — Fortführung Haupt-Review F-5: offene NATS-Slice-Planzeile weiter überholt (vom Fix-Auftrag nicht gedeckt)

- `kategorie`: INFO
- `quelle`: Maintainability („Arbeit überholt stehenden Träger", Planform —
  Übernahme aus [Haupt-Review
  F-5](review-slice-sdk-python-grpc-client-flaeche.md))
- `pfad`: `docs/plan/planning/open/slice-sdk-python-nats-stream-client-flaeche.md:127`
- `befund`: Die Planzeile des offenen NATS-Folge-Slices trägt weiter die
  inzwischen erledigte `options.py`-Docstring-Korrektur; die Fixrunde löste
  sie nicht (nicht in ihrem Auftrag F-1…F-4/F-6), und sie ist auch jetzt
  nicht Teil der Korrektur-Kette — bleibt beim nächsten Planner-Zug
  reduzierbar, damit der Folgelauf keinen bereits korrigierten Träger
  „korrigiert".
- `verifizierbar`: ja — `grep -n "options.py"
  docs/plan/planning/open/slice-sdk-python-nats-stream-client-flaeche.md`
  gegen den Ist-Stand von
  `sdks/python/pgchangefeed/src/pgchangefeed/options.py` (Satzform
  „no gRPC/SSE/NATS surface exists yet" dort nicht mehr vorhanden)
- `klasse`: „Arbeit überholt stehenden Träger" (Planform)

---

## Negativbefunde (geprüft, ohne Befund)

- **F-1 (Zahl) — gelöst und gemessen:** Plan §2 DoD-Zeile 78-83 trägt jetzt
  „31 Unit-Tests grün, darunter 8 neue gRPC-Tests (`grep -c "^def test_"` =
  8; 20 + 3 + 8 = 31 …)". Eigene Messung: `grep -c "^def test_"`
  `tests/test_grpc_client.py` = **8**, `test_http_client.py` = 20,
  `test_options.py` = 3 → **31** — deckungsgleich; Messform und Endstand
  („nach der Fixrunde, die F-4 … aufgelöst hat") sind am Ort genannt.
  Kein Drift. (Die Auflage V-1 ist damit real gelöst.)
- **F-4 (Timeout-Bindung) — gelöst und an der Eingabeseite gebunden:** der
  Fake zeichnet je Aufruf `(metadata, request, timeout)` (Zeilen 54, 71);
  `test_stream_changes_forwards_the_call_timeout` (30.0) und
  `..._leaves_the_call_deadline_unbounded_by_default` (None) halten beide
  Richtungen — die rot machende Mutation ist je Aussage benennbar:
  `timeout=timeout` in `grpc_client.py:70` streichen → Weiterleitungs-Test
  rot; dort fest verdrahten → Default-Test rot. Alle vier
  `invocations`-Consumer entpacken konsistent den 3-Tupel (Zeilen 168-190,
  kein veralteter 2-Tupel-Unpack). Der Fake-Kommentar trägt Zusage/Kopplung
  in Indikativ (keine Chronik, §3.7).
- **F-2 (typbewusste Prüfung) — für die drei Nicht-String-Felder gelöst:**
  alle zehn `isinstance`-Arten decken die Proto-Feldtypen exakt
  (`proto/cdc/stream/v1/changestream.proto` `message Change`: 7× string,
  1× int64, 2× bytes); der Nicht-`bool`-Guard schließt die
  bool-als-int-Lücke; `old_image == b""` am INSERT und der Sentinel im
  `new_image` binden die bytes-Felder inhaltlich; Feldentfall bricht real
  über `getattr` mit `AttributeError`. Rest-Lücke der drei
  String-Identitätsfelder: F-1 dieses Laufs.
- **F-3 (Handbuch + Versionshistorie) — die HIGH-Regel
  „Handbuch-Versionshistorie nicht fortgeschrieben" ist erfüllt:** der
  Diff ändert das Handbuch inhaltlich **und** hebt den `Version:`-Kopf
  (1.38 → 1.39 → 1.40) **und** ergänzt die Historie-Zeile im selben Diff;
  V-2 real gelöst: die neue Zeile trägt 1.40 (2026-09-23), `1.39` ist
  eindeutig (Kotlin, 2026-09-22), die Historie ist in Version/Datum
  konsistent geordnet — `grep -c "^| 1.39|^| 1.40"` = 2 (je eine Zeile).
  Die falsch werdende Satzform samt `ADR-0107`-Festlegung-1-Zitat ist aus
  dem Python-Absatz entfernt; „SSE und der NATS-Vollinhalts-Stream folgen
  im selben Folge-Release (`ADR-0110`)" ist eine wahre, verankerte
  Vorgriffs-Form (offene Welle-Slices). Die 1.40-Historien-Zeile zitiert
  die alte Satzform als Historie — exempt.
- **§3.13 Träger-Suchlauf (beide Stände gemessen: Range-Basis
  `b7a993fe~1` = `b6f82269` und HEAD `5d5fa3eb`)** — Suchlauf über
  `sdks/python/**`, `docs/**`, `README.md`, `spec/**` nach der bewegten
  Satzform-Klasse („gRPC/SSE/NATS bleiben außerhalb des Python-Packages"):
  `sdks/python/README.md` §Status trägt die gRPC-Fläche positiv mit
  [`ADR-0110`](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)-Anker und wahrer Out-of-Scope-Form für SSE/NATS;
  `__init__.py`-Docstring, `options.py`-Docstring und Attribut-Doku
  (V-3 deklariert, inhaltlich geprüft: `address`/`api_token` tragen die
  gRPC-Form), Handbuch-Python-Absatz — keine Rest-Satzform der Klasse.
  Root-`README.md` Zeile 27 („the HTTP API as an official Python client
  library") bleibt für das **veröffentlichte** 0.1.0-Package wahr, Nachzug
  beim Version-Bump gebündelt — dieselbe deklarierte Praxis wie bei C#
  (1.38: Bump mit Vollabdeckung) und Kotlin (1.39), kein Befund in diesem
  Slice. `SPEC-027`/`LH-FA-SST-009.a` gebündelt ([`ADR-0110`](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md) Folgepflicht 2,
  konsistent).
- **§3.11 (host-lokale Pfade)** — der ganze Fix-Diff (Plan, Handbuch,
  Verifikations-Report, beide Testdateien) enthält keinen host-lokalen
  absoluten Pfad; der Verifikations-Report nennt das Log-Verzeichnis nur
  als „Host-Temp-Verzeichnis" ohne Pfad.
- **§3.7 (Kommentar-Klassen)** — die neuen Code-Kommentare tragen Klassen:
  Fake-Kommentar (Zusage/Kopplung, Indikativ über den Ist-Zustand);
  Testkommentar `getypte Pruefung, weil ein \`!= ""\` … vakuum waere` ist
  Abgrenzung über ein Muster-Eigenschafts-Verhältnis (warum die Typ-Form),
  keine Vorher/Nachher-Sprache über diesen Datei-Zustand — kein
  Chronik-Befund. Die Plan-Nachzug-Zeilen tragen Review-/Verifikations-Anker
  als Herkunft (Plan-Artefakt, zulässige Form).
- **§3.12 Instanz A (neue Zahlen im Fix-Diff)** — 31/8/20+3+8 (gemessen,
  identisch), 1.40/1.39 (gemessen, identisch), `change_id=804-1` (trägt
  Lauf, unverändert), Verifikations-Report-Zahlen (82.80 %, „31 passed",
  55 Import-Zeilen, 270 Zeilen) tragen je Log-/Mess-Anker; keine Zahl im
  Fix-Diff ohne Ursprung (F-5 ist die eine, in `5d5fa3eb` gelöste
  Ausnahme am Range-Ende).
- **Commit-Traceability** — `b7a993fe`, `59e42d22` und `5d5fa3eb` tragen je
  `LH-FA-SST-009` + `ADR-0110` (Betreff geprüft).
- **Gate-Stempel-Frische** — `.harness/state/gates-passed.diffsha`
  (`09156789…`) ist byte-gleich dem frisch berechneten Baum-Hash
  (`tools/harness/working-tree-hash.sh`), der Stempel deckt den Ist-Baum
  (HEAD `5d5fa3eb`, Arbeitsbaum sauber) — der Nachzieh-Commit `5d5fa3eb`
  (id-unlinked-Behebung) ist gates-geprüft.
- **[`ADR-0110`](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md) Fitness-Zeile (reale Server-Instanz je Fläche)** — der
  Integrationstest trägt zwei Tests gegen die reale Instanz (positiver
  Empfang + `UNAUTHENTICATED`-Öffnung ohne Token), unverändert zur
  Haupt-Review-Prüfung.
- **`sdks/python/pgchangefeed/tests/test_grpc_client.py` (übrige Bindungen)**
  — die sechs Bestands-Tests (RPC-Pfad, Bearer-Metadata, filterloser
  Request, unmapped Yield, Unauthenticated-Durchreichung,
  Deskriptor-Bindung) sind durch die Tupel-Erweiterung nicht geschwächt.
- **`docs/reviews/verifikation-slice-sdk-python-grpc-client-flaeche.md`
  (neu im Diff)** — Verifier-Artefakt, gegen seinen eigenen Berichtsrahmen
  gelesen: Gegenstand/Commit-Liste stimmt mit `git log`, V-1…V-3 sind
  messbar und real gelöst (V-1: 31/8 gemessen; V-2: 1.40 eindeutig;
  V-3: Nachzug-Zeile real), die Zitate tragen Lauf-/Log-Anker, keine
  Hostpfade.

---

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 3 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** „Beleg trägt seinen Satz nicht"
(F-1, F-2 — F-1 zählt zur Herkunft
`BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht`, sechstes gezähltes
Vorkommen dieser Kette in diesem Slice-Verbund) · „Arbeit überholt
stehenden Träger" (F-3, F-4, F-7) · „Zahl im Träger ohne Ursprung — oder
gegen die Messung driftend" (F-5 — drittes Auftreten in dieser
Korrektur-Kette, in `5d5fa3eb` gelöst) · Übergabe-Artefakt-Zuordnung (F-6)

## Verdikt

**Merge-blockierend:** ja — **eine kleine Fixrunde am Implementer ist
nötig (1 MEDIUM: F-1)**. F-1 ist ein Fixrunden-Regression der
Bindungs-Breite in der Feldvollständigkeits-Assertion: die drei
String-Identitätsfelder brauchen je die Nichtleer-Bindung zurück (die
vormals korrekt bestand) — die Typ-Assertion bleibt daneben bestehen.
F-2…F-4 (LOW) und F-5…F-7 (INFO) gehen in dieselbe Rückgabe bzw. bleiben
als Carry-Over benannt; der DoD-Checkbox-Nachzug
(„Review durchgeführt, Report unter `docs/reviews/` liegt vor") läuft
regulär bei Schritt 21 des Implementer-Workflows, weil eine Fixrunde
existiert (`.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug ohne
Fixrunde greift nicht).

**Aufgelöst und geprüft:** F-1 (Zahl: 31/8, gemessen identisch), F-2
(typbewusste Prüfung für bytes/int, Rest in F-1 dieses Laufs), F-3
(Handbuch-Satzform entfernt, Version+Historie im selben Diff, V-2 real
gelöst), F-4 (Timeout je Aussage an der Eingabeseite gebunden), F-6 des
Haupt-Reviews (Pfad korrekt), V-1/V-2/V-3 (je eigene Messung oben).

**Nicht geprüft (Grenze dieses Laufs):** der reale
`make test-sdk-python-integration`-Lauf und der Pack-Lauf (DoD-Substanz —
Verifier, deren Fixrunden-Logs hier ausgewertet wurden); Mutationsproben
sind als statische Bindungsprüfung geführt (je Aussage ist die rot
machende Mutation benennbar, keine ausgeführte Mutation).