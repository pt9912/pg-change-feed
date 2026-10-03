# Review-Report: slice-sdk-meldungscodes-in-fehlertypen — 2026-10-03

**Review-Art:** Code — geprüft gegen Plan, [`ADR-0145`](../plan/adr/0145-sdk-meldungscodes-eigenschaft-der-fehlertypen.md),
[`ADR-0144`](../plan/adr/0144-meldungscodes-nutzerseitige-kennungen.md), das
[Architect-Verdikt](architect-verdict-sdk-meldungscodes-in-fehlertypen.md) und `AGENTS.md` Hard Rules
(Modul 10). Kein DoD-Abgleich (Verifier).

**Gegenstand:** Slice `sdk-meldungscodes-in-fehlertypen` (wellenlos), Diff-Range `85f0506c..a1947e8a`:
Implementierung `20ff9fe0` (C#), `5d8fe1af` (Kotlin), `7efea1fd` (Python), `a1947e8a`
(Runner-Phase je Sprache, Abdeckungs-Träger, Handbuch 1.93, `harness/README.md`, `harness/mk/sdk.mk`,
Plan-Nachzug); 38 geänderte Dateien. Plan:
[`slice-sdk-meldungscodes-in-fehlertypen`](../plan/planning/in-progress/slice-sdk-meldungscodes-in-fehlertypen.md).

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09“, seither um weitere
HIGH-/MEDIUM-Klassen ergänzt. **Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-03.

**Ablage:** Der Reviewer-Lauf hat diesen Report selbst geschrieben (Write-Werkzeug). Die Mutationsprobe lief
an einer Kopie von `sdks/csharp/` im Scratchpad (die Datei per Write ersetzt, kein `sed -i`, keine Umleitung auf
eine Repo-Datei; Bau per `docker build --no-cache --build-context proto=proto --target build` gegen die Kopie,
Image danach entfernt). Die Gate-Läufe liefen gegen den echten Arbeitsbaum, ungefiltert, Exit-Code direkt
gelesen (`AGENTS.md` §3.9). Keine Aktion wurde von der Berechtigungsschicht verweigert; ein Lauf des Edit-Werkzeugs
war in dieser Sitzung nicht verfügbar (keine Berechtigungs-Verweigerung eines Repo-Aufrufs, nur ein fehlendes
Werkzeug; der Plan-Haken unten ist deshalb per `git apply` gesetzt).

**Eingangs-Kontext:**

- Slice-Plan (oben), insbesondere §3 (Eingabetabellen, Mutationsbeleg, Belege des Implementers) und §6 (Risiken)
- [`ADR-0145`](../plan/adr/0145-sdk-meldungscodes-eigenschaft-der-fehlertypen.md) (Festlegungen 1 bis 6),
  [`ADR-0144`](../plan/adr/0144-meldungscodes-nutzerseitige-kennungen.md) (Wire-Form),
  [`ADR-0134`](../plan/adr/0134-sdk-public-doc-check-gate-make-gates.md) (keine interne Kennung unter `sdks/`),
  [`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md) (Herkunft von Aussagen)
- [`LH-FA-SST-009`](../../spec/lastenheft.md)
- `AGENTS.md` §3.1, §3.7, §3.9, §3.12, §3.13

---

## Findings

### F-1 — Basis-Konstruktor der C#-HTTP-Hierarchie: `null`-Literal wird mehrdeutig

- `kategorie`: LOW
- `quelle`: [`ADR-0145`](../plan/adr/0145-sdk-meldungscodes-eigenschaft-der-fehlertypen.md) Festlegung 4 („jede 0.5.x-Signatur bleibt binär und quellseitig erhalten“), `AGENTS.md` §3.12 Instanz B
- `pfad`: `sdks/csharp/PgChangeFeed.Client/Http/PgChangeFeedException.cs:29` und `:37` (protected `(int, string, Exception)` neben protected `(int, string, string?)`)
- `befund`: Ein Aufruf `base(status, text, null)` in einer außerhalb des Packages abgeleiteten Klasse bindet nach der C#-Überladungsauflösung an beide Konstruktoren und ist damit mehrdeutig (hergeleitet, nicht kompiliert); die ADR-Aussage „quellseitig erhalten“ ist an dieser Stelle breiter als ihr Beleg. Der Plan benennt die Mehrdeutigkeit für den eigenen Code (benannte Argumente), ein Test des Falls fehlt; die öffentlichen Blatt-Konstruktoren sind nicht betroffen.
- `verifizierbar`: nein — kein Test oder Gate kompiliert eine fremde Unterklasse gegen das Package.
- `klasse`: Aussage breiter als ihre Messung

### F-2 — Randfall `error` mit falschem JSON-Typ: die drei Fehlerbau-Parser lesen `code` verschieden

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.12; Plan §6 (Risiko „Drei Sprachen, drei Parser, drei Lesarten“); Beobachtung [`drei-sprachen-kopie-divergiert-am-randfall`](../plan/planning/observations/BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall/state.md)
- `pfad`: `sdks/csharp/PgChangeFeed.Client/Http/ErrorBody.cs:22-33`, `sdks/kotlin/pgchangefeed-kotlin/src/main/kotlin/io/github/pt9912/pgchangefeed/http/ErrorBody.kt:20-30`, `sdks/python/pgchangefeed/src/pgchangefeed/http_client.py:286-304`
- `befund`: Bei `{"error":5,"code":"PCF-E8051"}` (hergeleitet aus dem Code, nicht gefahren) liefert C# die Rohform als Text und `null` als Code (`System.Text.Json` wirft beim Lesen von `Error`, der `catch` verwirft den Code), Python `"5"` und den Code, Kotlin (Gson koerziert die Zahl) `"5"` und den Code. Die Eingabetabelle hat keine Zeile mit falsch getyptem `error` oder mit `code` ohne `error`; der Server sendet beides nicht.
- `verifizierbar`: ja — ein Tabellenfall je Sprache würde die Divergenz zeigen.
- `klasse`: Drei-Sprachen-Kopie divergiert am Randfall

### F-3 — Realserver-Fall pinnt den Code nicht auf den berichteten Wert

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.12 Instanz B; Reviewer-Skill „Beleg trägt seinen Satz nicht“
- `pfad`: `sdks/csharp/PgChangeFeed.Client.Integration/ErrorCodeRealserverTests.cs:64-65`, `sdks/python/pgchangefeed/integration/test_error_code_realserver.py:108-109`, Kotlin-Spiegel, Runner-Muster `PCF-E[0-9]{4}` in den drei Runnern
- `befund`: Der Plan (§3) und die Belege-Tabelle nennen `PCF-E8025` als am Server gemessen; die Tests prüfen nur Präfix `PCF-E` (bzw. vier Ziffern) und Gleichheit von HTTP- und gRPC-Wert. Die gedruckte Zeile des Laufs trägt den Wert, der Test hält ihn nicht; ein anderer `PCF-E`-Code am Server würde grün bleiben. Die Gleichheit beider Wege und das Nicht-`null` sind gebunden.
- `verifizierbar`: ja — `make test-sdk-csharp-integration` mit einer Erwartung auf den genauen Code.
- `klasse`: Beleg trägt seinen Satz nicht (schwach)

### F-4 — Diagnose-`error_code`-Tests prüfen nur den Durchreichungsweg

- `kategorie`: INFO
- `quelle`: Plan §2 (Punkt „Diagnose-`error_code`“), [`ADR-0145`](../plan/adr/0145-sdk-meldungscodes-eigenschaft-der-fehlertypen.md) Festlegung 6
- `pfad`: `sdks/kotlin/pgchangefeed-kotlin/src/test/kotlin/io/github/pt9912/pgchangefeed/MessageCodeTest.kt` (Test `diagnose error code arrives unchanged`) und die beiden Spiegel
- `befund`: Der Fake gibt das Protobuf-Objekt unverändert zurück; der Test wird nur rot, wenn der Client künftig einen Mapper um `Diagnose` legt, der das Feld verliert. Das entspricht dem Plan („kein Mapper-Code“); ein Beleg am Wire fehlt für `error_code` (kein Realserver-Fall der Diagnose in dieser Phase).
- `verifizierbar`: nein.
- `klasse`: Zusage nur an der Verdrahtung gebunden

### F-5 — Handbuch beschreibt Eigenschaft, die in keinem veröffentlichten Package liegt

- `kategorie`: INFO
- `quelle`: `.harness/skills/nutzerdoku-schreiben.md` (Ist-Zustand); Plan §1 („Release ist Folgeschritt“)
- `pfad`: `docs/user/benutzerhandbuch.md:1502` und `:1940-1946`
- `befund`: Das Handbuch sagt im Präsens, die SDK-Packages tragen den Code; die auf den Registries liegenden Packages (0.5.0) tragen ihn bis zum Release 0.6.0 nicht. Der Plan benennt den Release als Folgeschritt des Hauptlaufs, ohne ihn an das Handbuch zu binden; Versionsnummern in `csproj`, `build.gradle.kts` und `pyproject.toml` sind unverändert (gemessen, 0 Treffer im Diff).
- `verifizierbar`: nein.
- `klasse`: Träger vor dem Release beschreibt den Release-Zustand

### F-6 — Kotlin-Klassenpfad: Plan-Messung gegenstandslos

- `kategorie`: INFO
- `quelle`: Plan §3 („nicht gemessen: ob `com.google.rpc.Status` … im Kotlin-Klassenpfad liegt“)
- `pfad`: `sdks/kotlin/pgchangefeed-kotlin/src/main/kotlin/io/github/pt9912/pgchangefeed/grpc/StatusDetail.kt:5-8`, `sdks/kotlin/pgchangefeed-kotlin/build.gradle.kts:102`
- `befund`: Der Parser liest mit `CodedInputStream` und `com.google.protobuf.Any` aus `protobuf-java`, einer direkten `implementation`-Abhängigkeit (`build.gradle.kts:102`, unverändert); `com.google.rpc.*` wird nicht benötigt, die offene Messung ändert nichts am Ergebnis. Der Realserver-Lauf des Implementers (Kotlin, `NOT_FOUND` mit Code) trägt die Lesung am Wire (übernommen aus dem Bericht, nicht nachgefahren).
- `verifizierbar`: ja — `make test-sdk-kotlin-integration`.
- `klasse`: —

### F-7 — Runner-Phase `none`: die Identitäts-Zählung ist fest verdrahtet

- `kategorie`: INFO
- `quelle`: Maintainability; Plan §3 (Abweichung „eine Phase statt zwei“)
- `pfad`: `tools/harness/run-sdk-csharp-integration-tests.sh:301-305` und die beiden Spiegel
- `befund`: Für `sql_kind none` setzt der Runner `captured=1`; die nachfolgende Prüfung „genau eine Zeile“ ist damit für diese Phase konstruiert grün. Die Abweichung vom Plan-Vorschlag (eine Fehlercode-Phase mit beiden Wegen und Gegenprobe statt zwei Phasen) ist schlüssig begründet: der `READY`/`RECEIVED`/`REJECTED`-Mechanismus bleibt unverändert, die Zahl „vierzehn“ ist in den drei Hilfetexten, den drei Kopfkommentaren und den drei Zeilen von `harness/README.md` nachgezogen (4+1+4+4+1).
- `verifizierbar`: nein.
- `klasse`: —

## Negativbefunde

- geprüft, ohne Befund: `sdks/csharp/PgChangeFeed.Client/` (Eigenschaft an beiden Basen, Blatt-Typen erben; `StatusDetail.cs` — domain, `reason`, nicht lesbar gleich leer, erster passender `Any`, `type_url`-Suffix; `ErrorBody.cs`; alte Konstruktoren bleiben; keine neue Abhängigkeit — `Directory.Packages.props` und `.csproj` unverändert)
- geprüft, ohne Befund: `sdks/kotlin/pgchangefeed-kotlin/src/main` (`@JvmOverloads` erhält die alten JVM-Signaturen `(int,String)`/`(…,cause)`; `StatusDetail.kt`; `ErrorBody.kt` liest `code` nur als JSON-String; `JsonParser` wirft für `Bad Gateway` eine `JsonParseException`, der Zweig fängt sie)
- geprüft, ohne Befund: `sdks/python/pgchangefeed/src` (keyword-only `message_code`, `code` bleibt `grpc.StatusCode`; `_status_detail.py` — Varint-Leser, abgeschnittene Felder, ungültiges UTF-8 und Gruppen-Wire-Typen enden in `None`; `_extract_error`)
- geprüft, ohne Befund: Testquellen der drei Sprachen — je zwei Tabellen mit 12 Zeilen entsprechend Plan §3, das Statusdetail wird mit der Protobuf-Laufzeit gebaut, nicht mit dem Parser; Test des alten Konstruktors je Sprache
- geprüft, ohne Befund: Mutationsprobe (eigene, vom Implementer nicht berichtete Stelle): in der Kopie von `sdks/csharp/PgChangeFeed.Client/Http/ErrorBody.cs` die Typprüfung `ValueKind == String` durch `ToString()` jedes nicht fehlenden Werts ersetzt, Bau der Test-Stufe ohne Cache — `Failed!  - Failed:     2, Passed:   197, Skipped:     0, Total:   199`, rot genau an der Zeile 5 der HTTP-Tabelle (HTTP- und SSE-Fall). Die Mutationstabelle des Plans (M1 bis M4, Feldnummer) ist übernommen, nicht nachgefahren
- geprüft, ohne Befund: Kennungsfreiheit — `make sdk-public-doc-check` Exit 0 (`keine interne Kennung unter sdks`), `make handbuch-public-doc-check` Exit 0; ein `git diff`-Filter auf `ADR-|LH-|SPEC-|slice-|welle` über `sdks/` und `tools/` trifft nur Zeilen der Runner (Abdeckungs-Text, Ausgabezeilen, wie im Bestand), nichts in SDK-Quellen, Docstrings, KDoc, XML-Doku oder Fehlertexten
- geprüft, ohne Befund: Kommentarregeln `AGENTS.md` §3.7 — `make kommentar-kennungen DIFF=283d6175` Exit 0 ohne Kandidat; die neuen Kommentare der SDK-Dateien und Runner tragen keine Kennungskette, keine Chronik und keinen Konjunktiv über eine verworfene Alternative
- geprüft, ohne Befund: Handbuch 1.93 und Änderungshistorie (Betreibersicht, ohne Kennung), READMEs (Englisch, nur `PCF-…`-Beispiele, Stream-Grenze genannt); der Satz „die Beispiele und SDKs dieses Projekts werten es nicht aus“ ist berichtigt, die Beispiele werten das Detail tatsächlich nicht aus (Suchlauf-Zeile `examples` 0)
- geprüft, ohne Befund: Suchlauf `AGENTS.md` §3.13 — `make suchlauf-nachmessen PLAN=docs/plan/planning/in-progress/slice-sdk-meldungscodes-in-fehlertypen.md` Exit 0 (`21 Zeilen stimmen`); ein eigener `git grep -n -i dreizehn` über den Baum außer Records trifft nur Feldzahlen von `ChangeRecord`/`GET /changes` (andere Eigenschaft), keinen Träger der Phasenzahl
- geprüft, ohne Befund: `make gates` Exit 0 (ungefiltert in eine Log-Datei, Exit-Code gesondert gelesen); `make docs-check` Exit 0 (`1617 Datei(en) geprüft, 0 Befund(e)`); `make fmt-check` Exit 0
- geprüft, ohne Befund: Docker-only (`AGENTS.md` §3.1) — der Diff enthält kein Host-Werkzeug, das eine Repo-Datei umschreibt, keine Umleitung auf eine Repo-Datei, keinen Paketmanager; Gate-Läufe laut Plan-Tabelle mit gesichertem Exit (§3.9)
- geprüft, ohne Befund: `docs/user/sdk-e2e-abdeckung.md` — je Sprache eine Zeile in den Markern, Träger-Spalten stimmen mit Testklasse und Runner überein; Abhängigkeiten und Package-Versionen (`<Version>`, `build.gradle.kts`, `pyproject.toml`) im Diff unverändert

Nicht nachgefahren (übernommen aus dem Bericht des Implementers): die drei `make test-sdk-*-integration`-Läufe, `make sdk-pack-*` und `make examples-*` (Docker mit Netz, `:dev`-Image).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 3 |
| INFO | 4 |

**Finding-Klassen dieses Laufs:** Aussage breiter als ihre Messung · Drei-Sprachen-Kopie divergiert am Randfall · Beleg trägt seinen Satz nicht (schwach) · Zusage nur an der Verdrahtung gebunden · Träger vor dem Release beschreibt den Release-Zustand

## Verdikt

**Merge-blockierend:** nein — kein HIGH, kein MEDIUM. Die drei LOW-Findings (F-1 bis F-3) und die INFO-Findings gehen ohne Rückgabe-Pfeil an den Planner (Closure §7, ggf. Folge-Notiz); eine Fixrunde am Implementer ist nicht nötig. Die DoD-Zeile „Review durchgeführt“ im Plan ist mit diesem Report auf `[x]` gezogen (Reviewer-Skill §DoD-Checkbox-Nachzug ohne Fixrunde).

**Übergabe:** Dieser Report ist ein Lauf-Beleg; er ersetzt keine Verifikation — DoD-/Spec-Konformität
(darunter der Mutationsbeleg der DoD und die Läufe der Docker-Ziele) prüft der Verifier separat.
