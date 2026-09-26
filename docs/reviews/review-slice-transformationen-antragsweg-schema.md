# Review-Report: slice-transformationen-antragsweg-schema — 2026-09-26

**Review-Art:** Code — der Diff liefert die Antragsarten `set_transformation`/`remove_transformation` auf der
Antrags-Queue: die zwei nullable Spalten `rule_name`/`rule_spec` im neutralen Modell, die um zwei Werte
erweiterte `request_kind`-Menge, die zwei SQL-Funktionen samt `REVOKE … FROM PUBLIC`/`GRANT … TO cdc_admin`,
die Guard-Einträge (`knownForeignObjects`, neun Objekte), den Rollen-Test der Funktions-Rechte, den
Alt-Tag-Lauf des Guard-Skripts, die `.dockerignore`-Negation, die Domäne (Antragsarten, Konstruktor), das
Lesen der zwei Spalten im Store, den Fehlertext des `default`-Zweigs in `applyAdministrationRequest` und den
Architect-Zug [`ADR-0125`](../plan/adr/0125-transformationen-parametertyp-regelform-json.md) mit dem
`SPEC-019`-Nachzug; geprüft gegen Plan, ADRs, Spec-Stellen und `AGENTS.md` Hard Rules (Modul 10 §Drei
Review-Arten). Kein DoD-Abgleich — das ist Verifier-Aufgabe (Modul 11).

**Gegenstand:** Slice `slice-transformationen-antragsweg-schema` (Welle `welle-transformationen`), Diff-Range
`8d5d7d2e..f88463cc` (10 Commits, 29 Dateien, +1483/−448; Baum sauber, nicht gepusht).

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09“ (seither um weitere
HIGH-Klassen ergänzt, u. a. Kommentar-Zusage, Zahl-im-Träger mit `suchlauf`-Probe, Beleg-Satz,
Zusage-ohne-Eingabeseite, Nachzug-Nachbar, Neue-Betreiber-Oberfläche).
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-26.

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis; die `<Platzhalter>` darin sind Formbeispiele)*. Dieser
> Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link
> (`v<X.Y.Z>` · `regelwerk/<datei>.md` §<Abschnitt>). Der vendored Baum trägt
> genau einen Tag; der Sprung löscht den alten, und ein Link darauf färbt beim
> nächsten Bump ein Artefakt rot, das niemand mehr anfassen darf. Ein `pfad`-Feld
> auf den **geprüften Gegenstand** ist davon nicht betroffen — es zitiert den
> Stand des Laufs und darf ihn festhalten (`v<X.Y.Z>` ·
> `regelwerk/grundlagen-harness-dateien.md` §harness/README.md als
> Einstiegspunkt — diese Zeile ist selbst ein Beispiel der Form).

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde — ohne
diese Liste ist der Lauf nicht reproduzierbar):

- Slice-Plan `slice-transformationen-antragsweg-schema` (§1 Ziel und Abgrenzung, §2 DoD-Wortlaut als Bezug,
  §3 Plan mit Suchlauf-Feld, Mutationen und Belegen des Laufs, §6 Risiken); die Pläne der Adressen
  `slice-transformationen-betriebsdoku` und `slice-transformationen-antragsweg-usecase`
- [`ADR-0112`](../plan/adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) (Teilfrage 1,
  Folgepflicht 3), [`ADR-0125`](../plan/adr/0125-transformationen-parametertyp-regelform-json.md),
  [`ADR-0050`](../plan/adr/0050-sql-administration-antragsqueue-und-live-reload.md),
  [`ADR-0043`](../plan/adr/0043-schemamigrationen-mit-d-migrate.md),
  [`ADR-0047`](../plan/adr/0047-rollenspezifische-dsn-verdrahtung.md),
  [`ADR-0046`](../plan/adr/0046-sql-driving-adapter-lese-schreib-trennung.md),
  [`ADR-0114`](../plan/adr/0114-schema-rollout-vorlauf-view-signatur.md) (Entscheidung 7),
  [`ADR-0085`](../plan/adr/0085-build-kontext-ausnahme-test-only-zweck.md) (Festlegung 3),
  [`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md)
- [`SPEC-019`](../../spec/pflichtenheft.md) (Antrags-Datensatz, Transformations-Antragsarten, Fehlertext-Tabelle),
  [`SPEC-030`](../../spec/pflichtenheft.md), [`ARC-005`](../../spec/architecture.md)
- [`LH-FA-CFG-007`](../../spec/lastenheft.md), [`LH-FA-ADM-001`](../../spec/lastenheft.md),
  [`LH-QA-SEC-002`](../../spec/lastenheft.md), [`LH-QA-OPS-005`](../../spec/lastenheft.md)
- `AGENTS.md` (Hard Rules §3.1, §3.3, §3.5, §3.7, §3.9, §3.12, §3.13), `harness/conventions.md`
  (`MR-000`/`MR-001`/`MR-002`), `harness/targets/schema-rollout.md`
- Report-Gerüst: `docs/reviews/review-report.template.md`, Formvorbild
  `docs/reviews/review-slice-transformationen-kern-rename.md`

**Eigenständig durchgeführte Prüfungen** (gemessen, nicht aus dem Implementer-Bericht übernommen;
Exit-Codes ungepiped in Log-Dateien gesichert, gedruckte Zeilen zitiert; je ein schwerer Docker-Lauf zugleich,
`free -m` vor dem ersten Lauf 16,4 GB verfügbar):

- **Sensoren am Stand `f88463cc`:** `make test` (Race-Detektor) Exit 0, alle Pakete `ok`, das Paket
  `tools/schema/rolloutguard` und `internal/bootstrap` eingeschlossen; `make test-store` Exit 0, gedruckt
  „db-coverage: OK — DB-Adapter-Coverage 82.61% erfuellt Schwelle 80%“; `make a-check` Exit 0, gedruckt
  „gesamt: 0 Befund(e)“; `make coverage-gate` Exit 0, gedruckt „coverage-gate: OK — Coverage 83.90% erfüllt
  Schwelle 80%“; `make gates` (einmal) Exit 0, gedruckt „d-check: 1217 Datei(en) geprüft, 0 Befund(e)“,
  „commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID“,
  „generated-sync: OK“, „gesamt: 0 Befund(e)“; `make commit-traceability RANGE=8d5d7d2e..f88463cc` Exit 0,
  gedruckt „OK — 10 Commit(s) in "8d5d7d2e..f88463cc", Betreffs ohne Struktur-ID“.
- **Rollout, zweimal, gegen eine frische Wegwerf-PostgreSQL** (PostgreSQL-18-Digest aus `run-store-tests.sh`,
  eigenes Docker-Netz, `make schema-rollout` mit `SCHEMA_TARGET`/`SCHEMA_ROLLOUT_NETWORK`): Exit 0 und Exit 0.
  `tools/schema/down.sql` nach dem ersten Lauf byte-gleich dem committeten Stand, `tools/schema/plan.yaml`
  unterscheidet sich in genau einer Zeile, der `target`-Zeile (`postgres://cdc:***@rv-pg:5432/…` statt der
  Compose-Zeile); der `CREATE TABLE "administration_request"` des Reports trägt `"rule_name" TEXT`,
  `"rule_spec" JSONB`. Katalog danach: `set_transformation … p_rule_spec json`, `remove_transformation`,
  beide `prosecdef = t`, `proconfig = {"search_path=cdc, pg_temp"}`. Arbeitsbaum danach per
  `git checkout` der zwei Erzeugnisse zurückgesetzt.
- **`bash tools/harness/run-schema-rollout-guard-test.sh`** Exit 0, alle sechs Läufe; gedruckt: „Lauf 5 OK —
  Tag v0.2.0: Exit 0 (Rollout des Tags), Exit 0 (Arbeitsbaum, ohne Vorlauf), Exit 0 (Arbeitsbaum, zweiter Lauf);
  Zeile alttag-ch über cdc.changes lesbar; 19 Tabellen-/View-Rechte … EXECUTE auf 3 Funktionen … allein für
  cdc_admin (nicht PUBLIC), Spalten rule_name:text:YES,rule_spec:jsonb:YES (Alt-Zeile alttag-req: NULL),
  request_kind-Menge
  backfill,disable,enable,exclude_column,include_column,remove_transformation,set_transformation,
  cdc.set_transformation/cdc.remove_transformation unter cdc_admin schreiben pending-Anträge, cdc_reader:
  permission denied for function“ und die Schlusszeile „OK — alle Belege real erbracht (… Alt-Tag v0.2.0,
  Negativ-Abbruch: unbekannte Funktion bleibt bestehen, make-Exit 2/2 mit d-migrate-Exit 8)“. Skript gelesen:
  die drei Vorbedingungen prüfen das Delta dieses Slice (`fail` bei Abweichung, kein stilles Grün), die
  Rechte-Prüfung läuft je Funktion über `has_function_privilege` für die drei Rollen und über `aclexplode` für
  `PUBLIC`.
- **Die Begründung von `ADR-0125` nachgestellt:** die Funktionssignaturen in `nacharbeit-administration.sql`
  auf `jsonb` gestellt (Kopie, danach zurückgenommen), zwei `make schema-rollout` gegen eine frische
  Wegwerf-PostgreSQL: Exit 0, dann Exit 2 mit „make: *** [Makefile:278: schema-rollout] Fehler 5“; die Wache
  meldet davor „bekannte Fremdobjekt-Blocker (…) - --execute laeuft mit --allow-destructive“ — der Fehler
  liegt in `--execute`, nicht bei der Wache, wie `ADR-0125` es festhält. Unmutiert derselbe Lauf: Exit 0 und 0.
- **Randformen des Parameters `rule_spec`** (dieselbe Instanz, Superuser, `SELECT cdc.set_transformation('rv-src',
  'public','t','r', <Form>)`): `NULL` angenommen (Zeile mit `rule_spec` NULL); `'null'::json` angenommen (Text
  `null`); `'[1]'::json` angenommen; `'{oops'::json` „invalid input syntax for type json“, keine Zeile;
  `'{"a":"\u0000"}'::json` „ERROR: unsupported Unicode escape sequence“ („\u0000 cannot be converted to text“),
  keine Zeile; `'{"kind":"a","kind":"b", "z":1,  "a":2}'::json` angenommen, gelesen als
  `{"a": 2, "z": 1, "kind": "b"}`; `'{"kind":"x"}'::jsonb` „function cdc.set_transformation(unknown, unknown,
  unknown, unknown, jsonb) does not exist“. `set_transformation(…, NULL, '{}'::json)` und
  `remove_transformation(…, '')` schreiben je eine Zeile mit NULL bzw. leerem `rule_name`.
- **Suchlauf-Feld des Plans:** `make suchlauf-nachmessen PLAN=<Plan-Datei>` Exit 0, gedruckt
  „suchlauf-nachmessen: 24 Zeilen stimmen“. **Von Hand** mit `git grep` nachgefahren (Zahlen und Stände
  stimmten, die Befehle liefen ohne das Werkzeug): `exclude_column` am Parent `2f5ed3dc` 95, am Diff 96;
  „neun bekannten/neun Objekte“ am Diff 12, „sieben bekannten/sieben Objekte/sieben Fremdobjekte“ am Diff 0;
  `administration_request` in `internal tools` am Parent 137, am Diff 152; `administration_request` in
  `postgresstorage/schema.sql` 0/0; „geschlossene Menge“ in `wiring.go` 1/0; „sechs Läufe“ 6/6. Träger, die der
  Plan als „nicht gefunden“ führt, selbst gesucht: `postgresstorage/schema.sql` führt `administration_request`
  nicht (0 Treffer, Plan bestätigt); `docs/user` trägt weder `set_transformation` noch `remove_transformation`
  (bestätigt); `git grep -n -i 'create table.*administration_request'` trifft nur Plan und `plan.yaml`
  (bestätigt). Zusätzlich gesucht mit einer Wortform, die das Suchmuster des Plans nicht trägt („fünf
  Antrags-Funktionen“): **ein Treffer** (F-5).
- **Eingabeseiten-Mutationen** (Kopie der Datei über `sed` in Ausgabedatei und `cp`, kein `sed -i`; je einzeln
  gefahren, danach `git checkout` der Datei; Baum am Ende sauber):

  | Nr. | Mutation | Ergebnis |
  |---|---|---|
  | M1 | Konstruktor: `set_transformation` ohne Prüfung des Regelnamens (`ruleName == "" \|\|` gestrichen) | rot: `TestNewAdministrationRequestRejectsInvariantViolations` |
  | M2 | `translate.go`: `ruleName`/`ruleSpec` im Konstruktor-Aufruf vertauscht | rot: `TestReadPendingRequestsTranslatesRequests` |
  | M3 | `guard.go`: Eintrag `set_transformation(…)` gestrichen | rot: `TestDecideAllowsDestructiveWhenAllBlockersKnown`, `TestDecideViewSignatureWithKnownForeignObjects` |
  | M3b | `guard.go`: `in:json` → `in:jsonb` | rot: dieselben zwei und `TestDecideRefusesOtherSignatureSpelling` (drei, wie `ADR-0125` es nennt) |
  | M5 | `wiring.go`: `include_column` aus `processedAdministrationKinds` | rot: `TestApplyAdministrationRequestRejectsUnprocessedKind` |
  | M6 | `wiring.go`: `request.Kind` im Fehlertext durch ein Literal ersetzt | rot: `TestApplyAdministrationRequestRejectsUnprocessedKind`, `TestProcessAdministrationRequestsMarksTransformationRequestsFailed` |
  | R1 | SQL: `cdc.remove_transformation(…)` aus der `GRANT`-Liste | rot: `TestAdministrationDateiTraegtDieFunktionsRechte` („fehlt in einer `GRANT EXECUTE … TO cdc_admin`-Anweisung“) |
  | R2 | SQL: `cdc.set_transformation(…)` aus der `REVOKE`-Liste | rot: derselbe Test („fehlt in einer `REVOKE EXECUTE … FROM PUBLIC`-Anweisung“) |
  | R3 | SQL: zusätzliche Zeile `GRANT EXECUTE ON FUNCTION cdc.set_transformation(…) TO cdc_reader;` | rot: derselbe Test („cdc_reader trägt EXECUTE auf …“) |
  | R4 | SQL: neue Funktion `cdc.zz_neu()` ohne Rechte-Zeile | rot: derselbe Test |
  | R5 | SQL: `GRANT EXECUTE ON ALL FUNCTIONS IN SCHEMA cdc TO PUBLIC;` | **grün** (F-2) |
  | R6 | SQL: `GRANT ALL ON FUNCTION cdc.set_transformation(…) TO cdc_reader;` | **grün** (F-2) |
  | R7 | SQL: neue Funktion `public.zz_neu()` ohne Rechte-Zeile | **grün** (F-2) |
  | R8 | SQL: neue Funktion `"cdc"."zz_neu"()` ohne Rechte-Zeile | **grün** (F-2) |
  | R9 | SQL: `GRANT ALL PRIVILEGES ON ALL FUNCTIONS IN SCHEMA cdc TO cdc_reader;` | **grün** (F-2) |
  | S1 | SQL: `'remove_transformation'` aus dem CHECK (`make test-store`) | rot: `TestAdministrationRequestKindCheckCarriesExactlyTheSevenKinds`, `TestAdministrationRequestTransformationRequestsCarryRuleAndNotify` (SQLSTATE 23514) |
  | S2 | SQL: `remove_transformation` schreibt `NULL` statt `p_rule_name` (`make test-store`) | rot: `TestAdministrationRequestTransformationRequestsCarryRuleAndNotify` |
  | S3 | SQL: `remove_transformation` sendet auf dem Kanal `cdc_administratio` (`make test-store`) | rot: derselbe Test (`WaitForNotification` läuft in den Timeout, 5,04 s) |
  | F1 | Funktion `set_transformation` ohne `SET search_path` (auf der Wegwerf-Instanz eingespielt, Aufrufe wie Lauf 5 unter `cdc_admin`/`cdc_reader`) | **kein Test rot** — der Aufruf unter `cdc_admin` gelingt, unter `cdc_reader` „permission denied for function“ (F-6) |
  | F2 | Funktion `set_transformation` mit `SECURITY INVOKER` (dieselbe Instanz, Aufruf wie Lauf 5 unter `cdc_admin`) | rot in Lauf 5: „permission denied for table administration_request“ |

- **`.dockerignore`-Negation ([`ADR-0085`](../plan/adr/0085-build-kontext-ausnahme-test-only-zweck.md) Merkmal (iv)):**
  `make image` am Stand `f88463cc`: `harness/image-hash.txt` `sha256:b9cfd65527f004eead067ab4bbc7102b19f9236b1eda6e03ded2d74557f39fe4`;
  dieselbe Zeile `!tools/schema/nacharbeit-administration.sql` aus `.dockerignore` entfernt (Kopie), `make image`:
  derselbe Digest; Zeile wiederhergestellt, `make image`: derselbe Digest (Builder innerhalb der Sitzung
  derselbe, [`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md)). Merkmale (i)–(iii): eine Datei, gelesen von
  `TestAdministrationDateiTraegtDieFunktionsRechte`, der Leser steht im Kommentar. `make image` überschreibt die
  lokale, nicht committete Datei `harness/image-hash.txt` mit demselben Wert.
- **Kommentar-Probe nach `AGENTS.md` §3.7:** `git diff 8d5d7d2e..f88463cc` über `*.go`, `*.sh`, `*.sql`, `*.yaml`,
  `.dockerignore` (ohne `plan.yaml`), hinzugefügte Zeilen gegen `früher|bisher|wäre|würde|vorher|zuvor|seit
  slice|slice-[0-9a-z-]+|welle-|weiterhin|jetzt|neu `: zwei Treffer („weiterhin“ in einem Vorbestands-Kommentar
  des Guard-Skripts, „ein Login ohne diese Mitgliedschaft scheitert“ ohne Bezug zur Chronik); keine
  Slice-/Wellen-Chronik in Produktionspfaden, kein Konjunktiv über eine verworfene Alternative.
- **Umgebung:** Wegwerf-Container und -Netz (`rv-pg`, `rv-net`) nach den Läufen entfernt, kein `prune`;
  dangling-Volumes vor dem Lauf 34, nach dem Lauf 34; Arbeitsbaum am Ende sauber.
- **Nicht gefahren (Grenze):** `make test-replication` und `make test-integration` (der Diff berührt weder den
  Replication-Pfad noch den Compose-Rundlauf; beide rollen dieselbe Datei über `apply-rollout.sh` bzw. das
  Guard-Skript aus, dessen Läufe 1–6 gefahren sind).

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Die Festlegung 1 von `ADR-0125` und der Absatz „Transformations-Antragsarten“ in `SPEC-019` sagen, alles, was gültiges JSON ist — auch `NULL`, JSON-`null` und ein Wert ohne Objekt —, werde als Antrag angenommen. Gemessen wird ein gültiges JSON mit `\u0000` (`'{"a":"\u0000"}'::json`) beim Aufruf mit „unsupported Unicode escape sequence“ abgelehnt, ohne Zeile; `NULL` ist zudem kein JSON. Die Messung der ADR nennt sieben Formen (Tabelle in §Kontext, dort „die sieben Formen der Tabelle sind die gesamte Messung, eine Aussage über weitere Formen folgt nicht daraus“), der Satz der Festlegung gilt für alle; JSON-`null` steht in keiner der sieben Formen. | `AGENTS.md` §3.12 Instanz B (Verfasser einer ADR: Aussage über alle Werte einer Menge nennt die geprüfte Menge); `BEO-PGC/adr-aussage-breiter-als-ihre-messung`; [`SPEC-019`](../../spec/pflichtenheft.md) | `docs/plan/adr/0125-transformationen-parametertyp-regelform-json.md:121-124`; `spec/pflichtenheft.md:660-666` | ja — `SELECT cdc.set_transformation('s','public','t','r','{"a":"\u0000"}'::json)` gegen die Wegwerf-Instanz | ADR-Aussage breiter als ihre Messung |
| F-2 | MEDIUM | `TestAdministrationDateiTraegtDieFunktionsRechte` sagt in seinem Doc-Kommentar, „keine andere Rolle“ trage `EXECUTE` auf eine der Funktionen, und eine „achte Funktion ohne Rechte-Zeile“ färbe ihn rot. Gefahren: `GRANT ALL ON FUNCTION … TO cdc_reader` (R6), `GRANT EXECUTE ON ALL FUNCTIONS IN SCHEMA cdc TO PUBLIC` (R5), `GRANT ALL PRIVILEGES ON ALL FUNCTIONS … TO cdc_reader` (R9) und eine neue Funktion in `public.` (R7) oder mit quotiertem Namen (R8) lassen ihn grün — die regulären Ausdrücke lesen nur `GRANT EXECUTE ON FUNCTION` und `CREATE FUNCTION cdc.<a-z_>`. Die „Benannte Grenze“ im Kommentar nennt Instanz-Belege und dynamische Grants, nicht diese Formen. Für die zwei neuen Funktionen fängt sie der Realinstanz-Test (`…FunctionsRequireCdcAdminMembership`, `cdc_reader`/`cdc_capture` → 42501) und Lauf 5; eine spätere Funktion trägt diesen Schutz nicht. | Zusage ohne Bindung an ihre Eingabeseite (Reviewer-Skill, hier MEDIUM: die Hauptformen R1–R4 sind rot gebunden, die Lücke ist auf Grant-Schreibweisen begrenzt); `BEO-PGC/rollen-test-abdeckungsluecken` (offen, 2×); [`ADR-0047`](../plan/adr/0047-rollenspezifische-dsn-verdrahtung.md); [`LH-QA-SEC-002`](../../spec/lastenheft.md) | `internal/bootstrap/roles_rollout_file_internal_test.go:427-434`, `:466-500` | ja — Mutationen R5–R9 an einer Kopie, `make test` bleibt Exit 0 | Rollen-Test Abdeckungslücke (Grant-Schreibweise) |
| F-3 | MEDIUM | `cdc.set_transformation(…, NULL, '{}'::json)` (Regelname NULL), `cdc.set_transformation(…, 'r', NULL)` (Regelform NULL) und `cdc.remove_transformation(…, '')` (leerer Regelname) schreiben je eine `pending`-Zeile (gemessen), die der Konstruktor beim Lesen mit `ErrEmptyIdentifier` ablehnt (aus dem Quelltext hergeleitet; der Plan führt `ListPending: leere Kennung` als am Konstruktor-Pfad gemessen); `ReadPendingRequests` gibt den Fehler zurück, `processAdministrationRequests` protokolliert ihn und liest im nächsten Durchlauf dieselbe Zeile erneut — kein Antrag der Queue (auch `enable`/`disable`/`backfill`) wird verarbeitet, bis die Zeile entfernt oder ihr Status per `UPDATE` (Recht von `cdc_admin`) geändert ist. **Schwere bewertet:** nicht HIGH — nur die vertraute Rolle `cdc_admin` löst es aus, es gibt eine Abhilfe ohne Superuser, dieselbe Klasse besteht am Parent für `cdc.exclude_column(…, NULL)` (der Konstruktor lehnt dort `column == ""` mit demselben Sentinel ab, `git show 8d5d7d2e:internal/domain/model/administrationrequest.go`), und der Plan führt sie in §6 als offenes Risiko mit Adresse; MEDIUM, weil dieser Diff die Auslöser um zwei Funktionen und zwei Eingaben vermehrt und `SPEC-019` (Fehlertext-Tabelle) für genau diese Fälle `failed` zusagt, das bis zum Folge-Slice nicht eintritt. | Maintainability (unklare Fehlerbehandlung am Rand, Reviewer-Skill MEDIUM-Liste); [`SPEC-019`](../../spec/pflichtenheft.md); [`ADR-0050`](../plan/adr/0050-sql-administration-antragsqueue-und-live-reload.md) | `internal/adapters/driven/postgresstorage/sqlexec/translate.go:312-318`; `internal/bootstrap/wiring.go:1309-1312` | ja — Antrag mit leerem `rule_name` einfügen, `ListPending` lesen | Ungültiger Antrag stallt die Lesung der Queue (benannt, adressiert) |
| F-4 | MEDIUM | Die Adressen des Aufschubs tragen zwei Sachverhalte dieses Slice nicht: (a) `ADR-0125` Folgepflicht 2 und Plan §6 verlangen, dass der Betreiber-Text von `slice-transformationen-betriebsdoku` die Aufrufform der Regelform (Literal oder `::json`, ein `jsonb`-Wert wird abgelehnt) nennt — `git grep -n -i -e '::json' -e 'Aufrufform' -e jsonb` im Plan von `betriebsdoku` liefert 0 Treffer, der Plan bleibt laut §6 „unverändert“; (b) die Fälle `rule_spec` SQL-NULL und leerer/NULL-`rule_name` stallen die Queue (F-3); der Plan von `slice-transformationen-antragsweg-usecase` führt in seiner DoD den „leeren, fehlenden oder ungültigen Regelnamen“, nicht `rule_spec` SQL-NULL und nicht, dass die Lesung der Queue die Zeile liefern statt ablehnen muss. Die Meldung an den Planner steht im Bericht des Implementers, nicht in einem committeten Träger der Adresse. §1 des Plans sagt, `betriebsdoku` §2 nenne den aufgeschobenen Gegenstand „vollständig“; §6 räumt zwei fehlende Sachverhalte ein. | `AGENTS.md` §3.13 (Träger einer fremden Datei wird gemeldet); Reviewer-Skill HIGH „Neue Betreiber-Oberfläche ohne Handbuch-Zug“ (hier nicht ausgelöst: Aufschub mit Adresse steht) und MEDIUM „Wiederholung eines Musters“; `BEO-PGC/aufschub-adresse-nimmt-sendung-nicht-an` (offen, 2×; dies wäre das dritte Auftreten); [`ADR-0125`](../plan/adr/0125-transformationen-parametertyp-regelform-json.md) Folgepflicht 2 | `docs/plan/planning/open/slice-transformationen-betriebsdoku.md` §2; `docs/plan/planning/open/slice-transformationen-antragsweg-usecase.md` §2 (Zeilen 76–100); Plan `slice-transformationen-antragsweg-schema` §1 (Handbuch-Punkt), §6 (Punkt „Das Handbuch nennt die Funktionen nicht“) | ja — `git grep` in den zwei Plänen | Aufschub-Adresse nimmt die Sendung nicht an |
| F-5 | LOW | Der Kommentar an `TestAdministrationRequestColumnEndToEndAgainstPostgreSQL` sagt, der Schema-Stand trage „die fünf Antrags-Funktionen“ (`tools/schema/nacharbeit-administration.sql`); die Datei definiert sieben. Das Suchlauf-Feld des Plans führt „fünf Funktionen/fünf SQL-Funktionen“, nicht die Wortform „fünf Antrags-Funktionen“, und hält als Nichtgefundenes „kein weiterer Träger, der die Menge als geschlossene Liste von fünf führt“. | `AGENTS.md` §3.13; `BEO-PGC/arbeit-ueberholt-stehenden-traeger` (verkörpert); Reviewer-Skill „Beleg trägt seinen Satz nicht“ (das Suchmuster deckt die Wortform nicht) — LOW: Test-Kommentar, keine Zusage | `internal/bootstrap/administration_endtoend_test.go:37`; Plan §3 Suchlauf-Feld, Zeile „Aufzählungen der Antragsarten“ | ja — `git grep -n 'fünf Antrags-Funktionen'` | Träger vom Suchmuster nicht gefunden |
| F-6 | LOW | Die Zusage „`SECURITY DEFINER`, gepinnter `search_path`“ (Plan §1) ist an ihrer Eingabeseite nur zur Hälfte gebunden: `SECURITY INVOKER` färbt Lauf 5 rot (F2), das Entfernen von `SET search_path = cdc, pg_temp` färbt keinen Test (F1) — kein Test liest `proconfig` oder `prosecdef`, `git grep -n -e prosecdef -e proconfig -e pg_temp -- '*.go' '*.sh'` liefert 0 Treffer. Dieselbe Lücke besteht für die fünf Funktionen des Parents. | Zusage ohne Bindung an ihre Eingabeseite (Reviewer-Skill; LOW: kein Test behauptet die Bindung, die Eigenschaft steht nur im Zielsatz); Sicherheits-Härtung von `SECURITY DEFINER` | `tools/schema/nacharbeit-administration.sql:180`, `:199` | ja — Mutation F1 | Sicherheits-Härtung ohne Testbindung |
| F-7 | INFO | `jsonb` normalisiert die Regelform vor dem Lesen: `{"kind":"a","kind":"b","z":1,"a":2}` wird als `{"a": 2, "z": 1, "kind": "b"}` gelesen (letzter Wert eines doppelten Schlüssels gewinnt, Schlüssel nach Länge und Alphabet, Leerraum normalisiert). Der Store liefert `rule_spec::text` an den Use Case; die „strikte Dekodierung“ von `antragsweg-usecase` sieht doppelte Schlüssel nicht. Für die Reihenfolge trägt `SPEC-030` „nicht zugesagt“; über Duplikate sagt weder `ADR-0112` noch `SPEC-019` etwas. | Maintainability; [`ADR-0112`](../plan/adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) Teilfrage 1 (Spalte `jsonb`) | `internal/adapters/driven/postgresstorage/queries/queries.go:330-336` | ja — Randform an der Wegwerf-Instanz | Normalisierung durch die Spalte (benannt) |
| F-8 | INFO | `processedAdministrationKinds` (`enable/disable/exclude_column/include_column/backfill`) ist eine handgeführte Zeichenkette neben der `switch`-Anweisung, die sie beschreibt; ein Test bindet ihren Inhalt an den Fehlertext, keiner an die Fälle des `switch`. Der Slice, der `set_transformation` verarbeitet, färbt `TestApplyAdministrationRequestRejectsUnprocessedKind` rot (der Test erwartet einen Fehler für beide Arten) und zieht die Konstante damit nach; ein anderer neuer Fall ohne Konstante fiele nicht auf. | Maintainability | `internal/bootstrap/wiring.go:1335`, `:1372-1376` | nein | Zweite Quelle der verarbeiteten Menge (benannt) |
| F-9 | INFO | Der Commit `f88463cc` trägt den Typ `docs(plan)` und ändert neben dem Plan und `harness/targets/schema-rollout.md` je einen Kommentar in `tools/schema/nacharbeit-administration.sql` und `tools/schema/rolloutguard/guard.go` (Verweis auf `ADR-0125`); das Verhalten bleibt unverändert (`git show f88463cc -- tools/schema/` zeigt nur Kommentarzeilen). | Maintainability | Commit `f88463cc` | ja — `git show --stat f88463cc` | Commit-Typ deckt eine Kommentar-Änderung im Produktionsskript nicht |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `tools/schema/nacharbeit-administration.sql` (zwei Funktionen: Muster, Payload, Ereignis-Name, `search_path`, `SECURITY DEFINER`, `REVOKE`/`GRANT`, CHECK-Menge, `::jsonb`-Cast) | geprüft, ohne Befund über F-1, F-6 hinaus: die Funktionen gleichen den Nachbarn Zeile für Zeile (Payload `v_id`, Kanal `cdc_administration`, `RETURNS text`, `gen_random_uuid()`), `column_name` bleibt NULL (Test liest `column != nil`), `remove_transformation` schreibt `rule_spec` nicht; CHECK trägt genau sieben Werte (Store-Test bindet: S1 rot, achte Art `truncate` endet mit 23514); `NULL`/`'null'`/`'[1]'` werden wie in `ADR-0125` beschrieben angenommen (gemessen), ungültiges JSON scheitert ohne Zeile |
| `tools/schema/schema.yaml`, `tools/schema/plan.yaml`, `tools/schema/down.sql` | geprüft, ohne Befund: `rule_spec: { type: json }` erzeugt im Report `"rule_spec" JSONB`; `plan.yaml` und `down.sql` sind Betriebs-Erzeugnisse des ersten Rollouts (`harness/targets/schema-rollout.md` §Erzeugnisse) — nach eigenem Lauf unterscheidet sich `plan.yaml` nur in der `target`-Zeile, `down.sql` ist byte-gleich; der Test-Lauf stellt beide über `tools/schema/rollout-restore.sh` wieder her (`git status` nach `make test-store` und Guard-Skript leer) |
| `tools/schema/rolloutguard` (`knownForeignObjects`, neun Einträge; `guard_test.go`) | geprüft, ohne Befund: die Schreibweise `in:json` stimmt mit dem realen `--plan-only`-Report überein (Guard-Skript Lauf 2 nimmt den `--allow-destructive`-Pfad, Lauf 6 belegt den Abbruch bei unbekannter Funktion); M3, M3b rot; der neue Test bindet `in:jsonb` als abgelehnt; die Kommentar-Zählung „neun“ in `guard.go`, `guard_test.go`, `harness/targets/schema-rollout.md`, `harness/README.md`, `run-integration-tests.sh`, `run-schema-rollout-guard-test.sh` steht (Suchlauf 12/0) |
| `tools/harness/run-schema-rollout-guard-test.sh` (Lauf 5, Alt-Tag) | geprüft, ohne Befund: Vorbedingungen sind das Delta dieses Slice (Funktionen, Spalten, Antragsarten) und brechen mit `fail`; die Prüfung der Rechte läuft je Funktion für drei Rollen und `PUBLIC`; die `request_kind`-Menge wird als genau sieben Werte verglichen (Mutation im Lauf des Implementers rot; eigener Lauf grün); Aufruf unter `cdc_admin` schreibt `pending`, unter `cdc_reader` „permission denied for function“; Alt-Zeile trägt NULL in beiden Spalten; das Risiko der Kopplung an den jüngsten Tag steht mit Adresse in Plan §6 |
| `.dockerignore` | geprüft, ohne Befund: eine Datei, Leser im Kommentar, Image-Digest mit und ohne Zeile gleich (gemessen); die Fitness-Zeile 1 von `ADR-0085` („keine geänderte Datei unter `internal/**` oder `cmd/**`, die nicht auf `_test.go` endet“) ist die Regel der Coverage-Messungs-Slices und für diesen Feature-Slice nicht einschlägig |
| `internal/bootstrap/roles_rollout_file_internal_test.go` | geprüft, ohne Befund über F-2 hinaus: Umfang aus den `CREATE FUNCTION`-Zeilen abgeleitet, Kommas in Signaturen an den Klammern zerlegt, Kommentar-Zeilen entfernt, Leerraum und Groß-/Kleinschreibung normalisiert; R1–R4 rot; mehrzeilige Formen tragen die Ausdrücke (`(?is)`, `[^)]*` über Zeilen) |
| `internal/domain/model/administrationrequest.go` und Aufrufer | geprüft, ohne Befund: Parameterreihenfolge `column, ruleName, ruleSpec, kind`; Invarianten `rule_name` für beide Arten, `rule_spec` nur für `set_transformation` (M1 rot, Implementer-Mutationen an `remove` und `ruleSpec`); `git grep -n 'NewAdministrationRequest('` findet einen Produktions-Aufrufer (`sqlexec/translate.go:312`), zieht acht Argumente nach und ist von M2 gebunden; der Doc-Kommentar zählt sieben Arten |
| `internal/adapters/driven/postgresstorage` (`queries.go`, `translate.go`, Tests) | geprüft, ohne Befund über F-3, F-7 hinaus: `COALESCE(rule_name, '')`, `COALESCE(rule_spec::text, '')` — SQL-NULL wird der leere Wert, JSON-`null` der Text `null` (Store-Test bindet beide); Round-Trip liest Art, Regelname, Regelform und Spalte; Bereinigung skopiert über die eigenen Antrags-IDs bzw. `table_name` (`BEO-PGC/test-isolation-geteilter-zustand`); `make test-store` Exit 0, S1–S3 rot |
| `internal/bootstrap/wiring.go` (`default`-Zweig), `administration_internal_test.go` | geprüft, ohne Befund über F-8 hinaus: der Fehlertext nennt die Antragsart und die verarbeitete Menge (`enable`, `disable`, `exclude_column`, `include_column`, `backfill`); der erste Test bindet Kind und die fünf Namen einzeln (M5, M6 rot), der zweite geht über `processAdministrationRequests` bis zum `failed`-Vermerk statt `applied`; der Kommentar vor `applyAdministrationRequest` trägt Zustand im Indikativ |
| `docs/plan/adr/0125-transformationen-parametertyp-regelform-json.md`, `docs/plan/adr/README.md`, `spec/pflichtenheft.md` | geprüft, ohne Befund über F-1 hinaus: Kopf mit „Supersedes `ADR-0112` in genau zwei Stellen“ trifft die zwei Stellen (Festlegung von Teilfrage 1, Contra-Zelle der Option D — am Wortlaut von `ADR-0112` geprüft); Messungs-Tabelle mit Aufruf-Formen und die Rollout-Läufe eigenständig bestätigt (jsonb: Exit 0/2, Fehler 5; json: Exit 0/0); Option B als „übernommen … nicht wiederholt“ gekennzeichnet; Fitness-Function-Zeilen nennen Tests, die `rolloutguard`-Zeile stimmt mit M3b überein (drei Tests); Index-Zeile vorhanden, Form der Teil-Supersedes wie `ADR-0122`; die `SPEC-019`-Spalten-Zeilen („fünf übrigen“, „sechs übrigen“) stimmen mit sieben Arten; keine ADR- oder Slice-Kennung im Spec-Text (§3.4) |
| Suchlauf-Feld des Plans §3 (24 Zeilen, beide Stände) | geprüft, ohne Befund über F-5 hinaus: mit dem Werkzeug Exit 0 und neun Zeilen von Hand nachgefahren; die Nichtgefundenen-Aussagen zu `postgresstorage/schema.sql`, `docs/user` und zur Fixture-DDL bestätigt |
| Kommentare nach `AGENTS.md` §3.7 in allen geänderten Go-, Skript- und SQL-Dateien | geprüft, ohne Befund: keine Slice-/Wellen-Chronik in Produktionspfaden, keine verworfene Alternative im Konjunktiv (der Kommentar zum `json`-Parameter nennt die geltende Eigenschaft von d-migrate im Indikativ), Test-Kommentare nennen den Test als Subjekt |
| Hard Rules und Commit-Struktur | geprüft, ohne Befund über F-9 hinaus: Docker-only (alle Läufe über `make` oder `docker run`), keine Suppression, `make a-check` „gesamt: 0 Befund(e)“, Coverage 83.90 % gegen 80 %, die zwei Lifecycle-Moves sind reine Renames (`4c91dc1c`, `2f5ed3dc`: 0 Zeilen), Inhalt steht in eigenen Commits (`AGENTS.md` §3.3), alle zehn Betreffs nennen `LH-FA-CFG-007` und/oder `ADR-*`, keiner trägt `SPEC-*`/`ARC-*` im Betreff |
| Handbuch-Kandidatenlauf (`docs/user/benutzerhandbuch.md`) | geprüft, ohne Befund über F-4 hinaus: neue Betreiber-Oberfläche (zwei `cdc.*`-SQL-Funktionen) ohne Handbuch-Zug im Diff, aber mit benanntem Aufschub und Adresse `slice-transformationen-betriebsdoku` (existiert in `open/`, §2 nennt SQL-Funktionen, Rolle, Status, Fehlertexte, Wirkung, Abhilfe, Glossar, Änderungshistorie); keine Handbuch-Änderung, also keine Versionshistorie-Pflicht |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 4 |
| LOW | 2 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** ADR-Aussage breiter als ihre Messung · Rollen-Test Abdeckungslücke (Grant-Schreibweise) · Ungültiger Antrag stallt die Lesung der Queue (benannt, adressiert) · Aufschub-Adresse nimmt die Sendung nicht an · Träger vom Suchmuster nicht gefunden · Sicherheits-Härtung ohne Testbindung · Normalisierung durch die Spalte (benannt) · Zweite Quelle der verarbeiteten Menge (benannt) · Commit-Typ deckt eine Kommentar-Änderung im Produktionsskript nicht

Hinweis zum Zähler: `BEO-PGC/adr-aussage-breiter-als-ihre-messung` steht bei 5× (`AGENTS.md` §3.12,
verkörpert); F-1 wäre das sechste Auftreten, diesmal an einer ADR desselben Architect-Zugs, der die Regel
trägt. `BEO-PGC/rollen-test-abdeckungsluecken` (offen, 2×) und
`BEO-PGC/aufschub-adresse-nimmt-sendung-nicht-an` (offen, 2×) hätten mit F-2 bzw. F-4 je ein drittes
Auftreten.

## Verdikt

**Merge-blockierend:** ja — F-1 (MEDIUM, eine Aussage in `SPEC-019` und in der Festlegung 1 von `ADR-0125` reicht
über die Messung hinaus und ist an einer Randform falsch; die Spec-Aussage ist in place änderbar, die
`Accepted`-ADR nur über eine Folge-ADR, `AGENTS.md` §3.5) und F-2 (MEDIUM, der Zusage-Umfang des
Doc-Kommentars von `TestAdministrationDateiTraegtDieFunktionsRechte` reicht über die Formen hinaus, die der
Parser liest) gehören in eine Fixrunde. F-3 ist bewertet und **nicht HIGH**: eine benannte Grenze mit Adresse
(`slice-transformationen-antragsweg-usecase`), die an ihrer Adresse vollständig ankommen muss (F-4); F-3 und F-4
gehen an den Planner und hängen nicht an einer Änderung im Code dieses Slice. Der Rest des Diffs trägt: der
Antragsweg (Funktionen, CHECK, Spalten, Rechte) ist gegen 14 rot gesehene Mutationen an der Eingabeseite
gebunden, der Rollout ist über zwei Läufe und den Alt-Tag-Lauf real belegt, die Begründung von `ADR-0125` ist
nachgestellt und trägt.

**Übergabe:** Findings gehen an den Implementer (Fixrunde nötig, deshalb bleibt die DoD-Zeile „Review
durchgeführt“ im Plan offen und wird bei Schritt 21 des Implementer-Workflows nachgezogen); F-1 berührt
zusätzlich den Architect (Folge-ADR oder benannte Grenze zu `ADR-0125` Festlegung 1). Die **Finding-Klassen**
gehen zusätzlich in die Slice-Closure §7 und von dort in den Zähler. F-3 und F-4 gehen mit der Frist „Closure
dieses Slice“ an den Planner (Pläne von `betriebsdoku` und `antragsweg-usecase` nachziehen). Dieser Report
selbst ist ein **Lauf-Beleg** (Audit: dieser Diff, dieser Skill, dieses Modell, dieses Verdikt) — er wird über
Läufe hinweg nicht wieder gelesen, und muss es nicht. Der Report ersetzt keine Verifikation —
DoD-/Spec-Konformität prüft der Verifier separat (Modul 11; anderes Prüf-Artefakt, anderer Eingabe-Kontext).
