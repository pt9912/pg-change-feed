# Slice api-token-mehrfach-konfiguration: mehrere gleichzeitig gültige API-Token je Klasse für HTTP und gRPC

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle — es gibt keine Closure-Bedingung jenseits der DoD
dieses Slice (gemessen: `ls docs/plan/planning/*.md` am Planungsstand nennt nur
`README.md`; die Roadmap führt unter *Offene Wellen* keine Welle), siehe
Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle braucht
(Modul 6).

**Abhängigkeit:** keine. Dieser Slice ist der **erste** der drei
Umsetzungs-Slices zur Spec 0.15.0; er geht
[`slice-tls-http-grpc-server`](../open/slice-tls-http-grpc-server.md) und
[`slice-otlp-metrik-export`](../open/slice-otlp-metrik-export.md) voraus, weil alle drei
dieselben Dateien der Konfiguration (`internal/bootstrap/wiring.go`,
`internal/bootstrap/config_file.go`), die Code-Tabelle der Meldungscodes und
den Konfigurationsteil des Benutzerhandbuchs berühren und die Zugangsdaten-Klasse
der Konfigurationsdatei hier den ersten Schritt macht (§1, §3).

**Bezug:** [`LH-FA-SST-012`](../../../../spec/lastenheft.md) (Scope),
[`LH-FA-SST-006`](../../../../spec/lastenheft.md) und
[`LH-FA-SST-008`](../../../../spec/lastenheft.md) (die zwei Schnittstellen, deren
Token-Prüfung sich ändert; ihr Verhalten mit einem Token bleibt unverändert),
[`LH-QA-SEC-002`](../../../../spec/lastenheft.md),
[`ADR-0150`](../../adr/0150-tls-und-mehrfach-token.md) (Festlegung 2 bis 4),
[`ADR-0152`](../../adr/0152-zugangsdaten-klasse-elf-schluessel.md)
(Zugangsdaten-Klasse, Folgepflicht, hier die ersten zwei der vier Schlüssel),
[`ADR-0144`](../../adr/0144-meldungscodes-nutzerseitige-kennungen.md)
(Meldungscodes der Klasse `configuration`),
[`ADR-0057`](../../adr/0057-http-grpc-api.md) und
[`ADR-0060`](../../adr/0060-grpc-streaming-mechanismus.md) (Bestand der zwei
Prüfungen), [`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md).

**Berührte Spec-Stellen:** [`SPEC-035`](../../../../spec/pflichtenheft.md) ·
[`SPEC-016`](../../../../spec/pflichtenheft.md) (die zwei Listen-Schlüssel in der
Klasse, die zwei Variablen env-exklusiv) ·
[`SPEC-018`](../../../../spec/pflichtenheft.md) ·
[`SPEC-031`](../../../../spec/pflichtenheft.md) ·
[`SPEC-008`](../../../../spec/pflichtenheft.md) (Klasse `configuration`) ·
[`ARC-005`](../../../../spec/architecture.md) (Driving Adapter).
Der Verweis zeigt **aufwärts**: Die Spec nennt diesen Slice nie
(Baseline-Regelwerk `grundlagen-referenz-richtung.md`
§Referenz-Richtung (SDP), `grundlagen-source-precedence.md` §ID-Schema als Klammer).

**Verantwortlich:** —
<!-- BEDIENHINWEIS: Verantwortlich hält die Arbeit — der Rolleninhaber der
Implementer-Rolle, gesetzt beim Übergang open→next (Baseline-Regelwerk
modul-05-planning-harness.md §Lifecycle als State Machine). Der Autor schrieb
den Plan; zwei Felder, zwei Fragen. Kein Statuswert: der Zustand bleibt das
Verzeichnis. Kein Sensor prüft das Feld — es ist Deklaration. -->

**Autor:** Planner-Agent, direkt beauftragt (kein Architect: die Entscheidung
liegt mit [`ADR-0150`](../../adr/0150-tls-und-mehrfach-token.md) vor; die eine
Schichten-Frage dieses Slice — Ort des gemeinsamen Klassifikators — steht in §3
als Architect-Frage mit Vorschlag). **Datum:** 2026-10-04.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Anlass:** Je Klasse (lesend, administrativ) gilt heute genau ein Token
(`CDC_API_TOKEN_READER`, `CDC_API_TOKEN_ADMIN`); ein Wechsel hat ein Zeitfenster,
in dem Clients abgewiesen werden
([`LH-FA-SST-012`](../../../../spec/lastenheft.md)). Die Zuordnung eines Tokens zu
seiner Klasse steht **zweimal** im Baum, je einmal in
`internal/adapters/driving/http/middleware.go` und
`internal/adapters/driving/grpc/interceptor.go` (`classifyToken`, Vergleich mit
`==`); beide Kommentare nennen die zweite Fassung und verlangen die gemeinsame
Änderung ([`ADR-0150`](../../adr/0150-tls-und-mehrfach-token.md) Kontext, dort am
2026-10-04 gelesen; hier gemessen am Parent: `classifyToken` hat 27
Trefferzeilen, §3 Suchlauf Zeile 1).

**Ziel:** Je Token-Klasse der HTTP- und der gRPC-Schnittstelle sind mehrere
Token gleichzeitig gültig — über die env-exklusiven Variablen
`CDC_API_TOKENS_READER` und `CDC_API_TOKENS_ADMIN` neben den unveränderten
Singularen —, beide Schnittstellen entscheiden über **einen** gemeinsamen,
zeitkonstant vergleichenden Klassifikator, und der Wechsel (zwei Neustarts) ist
am laufenden Container belegt und im Benutzerhandbuch beschrieben.

**Schnitt der Zugangsdaten-Klasse (Entscheidung dieses Plans).** Die Klasse
der Konfigurationsdatei ([`ADR-0152`](../../adr/0152-zugangsdaten-klasse-elf-schluessel.md)
Festlegung 1) wächst um vier Schlüssel: `api_tokens_reader`, `api_tokens_admin`
(dieser Slice) und `otlp_endpoint`, `otlp_headers`
([`slice-otlp-metrik-export`](../open/slice-otlp-metrik-export.md)). Jeder Slice trägt
**genau seine zwei** Schlüssel ein: dieser Slice hebt die Klasse von **sieben auf
neun**, der OTLP-Slice von neun auf elf. Begründung: (1) das Handbuch beschreibt
den Ist-Zustand — `otlp_endpoint` in einer Verbotsliste, bevor es die Funktion
gibt, wäre eine Aussage über etwas, das nicht da ist; (2) jeder Slice ist einzeln
lieferbar und zeigt an seinem Endstand Code, Test und Handbuch mit derselben Zahl
(von Hand nachzuzählen,
[`ADR-0152`](../../adr/0152-zugangsdaten-klasse-elf-schluessel.md)
Fitness Function); (3) ein gemeinsamer erster Schritt als vierter Slice wäre
Verwaltung ohne Lieferwert. Die Halbänderungs-Gefahr (ein Slice ändert Liste und
Test, der andere nur eines von beiden) fängt ein Test, der **beide Listen
gegeneinander** prüft (§2 Liefer-Punkt 1). Zwischen der Closure dieses Slice und
der des OTLP-Slice nennt `SPEC-016` elf, der Code neun: das ist der
Umsetzungsrückstand, den
[`ADR-0152`](../../adr/0152-zugangsdaten-klasse-elf-schluessel.md) §Kontext schon
benennt, kein Widerspruch der Spec in sich; `tls_cert_file` und `tls_key_file`
gehören nicht zur Klasse.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **TLS** — [`slice-tls-http-grpc-server`](../open/slice-tls-http-grpc-server.md): andere
  Anforderung ([`LH-FA-SST-011`](../../../../spec/lastenheft.md)), andere
  Eingriffsstellen (Listener statt Token-Prüfung); ein gemeinsamer Slice machte
  den Wechsel-Beleg vom Transport abhängig.
- **Die Schlüssel `otlp_endpoint` und `otlp_headers` in der Zugangsdaten-Klasse**
  — [`slice-otlp-metrik-export`](../open/slice-otlp-metrik-export.md) (Begründung oben,
  Punkt 1); die Adresse nimmt die Sendung an: der Plan dort nennt beide
  Schlüssel in §2.
- **Neuladen der Token zur Laufzeit, Ablauffristen, Rotation, Hash-Form** —
  Out-of-Scope von [`LH-FA-SST-012`](../../../../spec/lastenheft.md); der Wechsel
  bleibt Konfigurationswechsel mit Neustart.
- **Das Token des NATS-Vollinhalts-Zustellwegs** (`CDC_NATS_STREAM_TOKEN`) — der
  NATS-Server vergibt und wechselt es
  ([`LH-FA-SST-012`](../../../../spec/lastenheft.md) Out-of-Scope).
- **Änderung der Antwortformen und Statuscodes** (`401`, `403`, gRPC
  `Unauthenticated`) — [`SPEC-018`](../../../../spec/pflichtenheft.md) bleibt
  unverändert; ein Client mit einem gültigen Token merkt nichts.
- **SDK-Packages und Beispielprogramme unter `examples/`** — sie senden ein
  Token je Aufruf; der Server-Vertrag ändert sich für sie nicht (zu belegen durch
  den Suchlauf in §3 Zeile 5: `CDC_API_TOKEN_*` bleibt in ihnen unverändert, kein
  Träger dort nennt „ein Token je Klasse“ als Servergrenze — Erwartung, am
  Arbeitsstand nachzumessen).

**Keine Mindestzahl.** Ein Slice mit *einem* echten Ausschluss ist besser als
einer mit vier erfundenen; die vier Klassen sind ein Suchraster, keine
Ausfüll-Liste. Suchreihenfolge: Was übernimmt ein **Folge-Slice** (mit
Kennung — und die Kennung muss den Punkt auch annehmen)? Was bleibt als
**Bestand** bewusst stehen (mit Begründung)? Was wäre ein **anderer Vorgang**?
Welche **Schicht** rührt der Slice nicht an?

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

Alle Beleg-Angaben dieser Liste sind **Zusagen** („zu belegen durch …“): der
Planungsstand hat keinen davon gefahren
([`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md)
Instanz B). Die gedruckte Zeile bzw. der Exit-Code steht im Bericht des
Implementers und in §7.

- [ ] **Klassifikator, Konfiguration, Zugangsdaten-Klasse (Liefer-Punkt 1,
      [`LH-FA-SST-012`](../../../../spec/lastenheft.md) Happy Path, Boundary,
      Negative).**
      (a) *Ein* Klassifikator: das gemeinsame Paket (§3) trägt Rollen-Typ,
      Konstruktion aus den zwei Token-Mengen und die Zuordnung; die HTTP-Middleware
      (`withToken`) und der gRPC-Interceptor (Stream-Interceptor, Administration-
      Rollentabelle `administrationRPCRoles`) rufen sie, und die
      zwei Fassungen samt der Kommentare „zweite, wortgleiche Fassung“ sind
      entfallen (zu belegen durch die Suchlauf-Zeilen 1 und 2 am `diff`-Stand:
      kein `classifyToken` und keine wortgleiche Fassung mehr in `adapters/`).
      (b) Verhalten, je Fall ein Tabellentest-Eintrag mit Meldungstext
      (`make test`, gedruckte `ok`-Zeilen der Pakete): beide Token einer Klasse
      werden gleich bedient; ein entferntes Token ist `roleNone`; **derselbe Wert in
      beiden Klassen ergibt die administrative Klasse**; ein leerer Wert matcht in
      keiner Klasse (auch nicht gegen eine leere, ungesetzte Klasse); ein Token mit
      Komma im Singular bleibt ein Token (`a,b` als Singular wird nicht zerschnitten);
      Vereinigung Singular plus Liste. Dieselben Fälle sind je einmal durch den
      HTTP-Server-Test (`401`/`403`/`200`) und den gRPC-Interceptor-Test
      (`Unauthenticated`/durchgelassen) gefahren, nicht nur im Paket.
      (c) Konfiguration: `CDC_API_TOKENS_READER`/`CDC_API_TOKENS_ADMIN` werden in
      **beiden** Zugriffswegen gelesen (`ConfigFromEnv` und `mergeConfig`/
      `ConfigFromEnvAndFile`), auch unter geladener Datei (env-exklusiv); ein
      leeres Element, ein Element mit Leerraum (Leerzeichen, Tabulator) und eine
      Liste aus nur Trennzeichen enden mit `ErrConfiguration` und einem **neuen**
      Meldungscode (nächste freie Nummer ab `PCF-E2008`, §3 Zuweisungsregel; Eintrag
      in `internal/domain/messagecode/codes.go`, Katalog-Zeile im Handbuch,
      `make meldungscodes-check` Exit 0). Eine **leere Zeichenkette** der
      Plural-Variable gilt wie „ungesetzt“ (die Umgebung unterscheidet beides nicht;
      Festlegung dieses Plans, Spec-Lücke-Kandidat §6). Jeder dieser
      Negativfälle prüft den Meldungscode **und** den Klartext des betroffenen
      Elements und trägt eine Gegenprobe ohne die Verletzung (`a,b` ist gültig),
      damit der Test an seiner Eingabe hängt
      ([`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`](../observations/BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe/observation.md)).
      (d) Zugangsdaten-Klasse **neun**: `forbiddenFileCredentialKeys` trägt
      zusätzlich `api_tokens_reader` und `api_tokens_admin`;
      `TestConfigFromFileLehntZugangsdatenAb` iteriert über die neun (Liste **und**
      Kommentar „sieben“ geändert); ein **neuer** Test hält die Liste des Codes und
      die des Tests in beiden Richtungen gleich (Mengengleichheit, gegen die
      Halbänderung durch den Folge-Slice); eine Datei mit `api_tokens_reader:`
      endet mit `ErrConfiguration` und der Zeile „Zugangsdaten bleiben
      env-var-exklusiv“.
      (e) **Mutationsprobe (Zusage, an einer Kopie im Scratchpad, nicht an der
      Datei im Arbeitsbaum — [`AGENTS.md`](../../../../AGENTS.md) §3.1; Instanz: der
      Go-Test des jeweiligen Pakets im Toolchain-Container mit der Kopie als
      Mount):** vier Stellen einzeln mutiert, je die rote Farbe im Bericht —
      Prüfreihenfolge Reader vor Admin im Klassifikator (Fall „derselbe Wert“ rot),
      Entfernen der Leer-Token-Wache (Fall „leerer Wert“ rot), Entfernen der
      Leerraum-Prüfung des Listen-Parsers (Fall „Element mit Leerraum“ rot),
      Streichen von `api_tokens_admin` aus `forbiddenFileCredentialKeys`
      (Mengengleichheit und Klassen-Fall rot). Die Verallgemeinerung auf „alle
      Stellen“ ist **hergeleitet**, nicht erprobt; die **Zeitkonstanz** der
      Prüfung ist als *Erwartung* geführt (die Umsetzung vergleicht gegen alle
      Token ohne Abbruch beim ersten Treffer; ein Test misst keine Laufzeit, die
      Lesung des Reviewers trägt sie, §6).
- [ ] **Realserver-Beleg des Wechselablaufs (Liefer-Punkt 2,
      [`LH-FA-SST-012`](../../../../spec/lastenheft.md) Happy Path und Negative am
      laufenden Prozess).** Eine neue Phase von
      `tools/harness/run-integration-tests.sh` (`make test-integration`) fährt den
      Wechsel am Feed-Container über eine Compose-Override-Datei im Temp-Verzeichnis
      des Runners (Muster: die Phasen zur Leerlauf-Bestätigung, `compose.yaml`
      bleibt unverändert, der Container wird am Phasen-Ende ohne Override
      wiederhergestellt): (1) altes Token als Singular, neues in der Liste — beide
      werden über HTTP (`GET /tables` mit `reader`, ein administrativer Endpunkt
      mit `admin`) **und** gRPC (Wegwerf-Client) bedient; (2) Neustart ohne das
      alte Token — das alte endet mit `401` bzw. `Unauthenticated`, das neue wird
      bedient; (3) derselbe Wert in beiden Klassen erreicht einen administrativen
      Endpunkt; (4) eine Liste mit leerem Element lässt den Container nicht starten
      (Beendigung mit der Fehlerklasse `configuration` und dem neuen Meldungscode in
      der Log-Zeile, per `docker inspect`/`docker logs` gelesen — nicht aus dem
      Unit-Test übernommen,
      [`BEO-PGC/adapter-unittest-verdeckt-bootstrap-luecke`](../observations/BEO-PGC/adapter-unittest-verdeckt-bootstrap-luecke/observation.md)).
      Die Phase trägt eine `abdeckung_declare`-Zeile für
      [`LH-FA-SST-012`](../../../../spec/lastenheft.md);
      `docs/user/e2e-abdeckung.md` ist vom Runner neu geschrieben (Erzeugnis,
      committet); `make doc-trace` druckt die Zeile mit der Waisen-Zahl, und
      `LH-FA-SST-012` ist darin **keine** Waise (gedruckte Zeile im Bericht).
      Läuft die Phase auf dem Stack nicht (Runner-Struktur), bleibt der
      Befund im Bericht, keine stille Auslassung.
- [ ] **Benutzerhandbuch (Liefer-Punkt 3).** `docs/user/benutzerhandbuch.md`:
      Umgebungsvariablen-Tabelle (§5) um die zwei Plural-Variablen (Form, Vereinigung
      mit dem Singular, Fehlerfälle, Vorrang `admin`), der Absatz
      „Authentifizierung“ der HTTP-/JSON-API (kein Satz mehr, der „ein Token je
      Klasse“ als Grenze nennt), ein Abschnitt **Token-Wechsel in zwei Neustarts**
      (neues Token ergänzen, Clients umstellen, altes entfernen; kein Neuladen zur
      Laufzeit), der Klassen-Satz der Konfigurationsdatei (§5, „Zugangsdaten
      bleiben env-var-exklusiv“) mit den **neun** Schlüsseln und die
      Katalog-Zeile des neuen Meldungscodes; Änderungshistorie-Zeile mit der
      nächsten freien Nummer (`BEO-PGC/handbuch-versionshistorie-uebersprungen`).
      Kennungsfrei: `make handbuch-public-doc-check` Exit 0.

Gate- und Lauf-Pflichten (zählen nicht zu den Liefer-Punkten):

- [ ] `make gates` grün (Exit direkt ausgewertet, am Endstand;
      [`AGENTS.md`](../../../../AGENTS.md) §3.9), `make docs-check` Exit 0
      (Kennungen in diesem Plan verlinkt), `make test`, `make test-store`
      (die Verdrahtungs-Tests von `internal/bootstrap`) und
      `make test-integration` Exit 0, `make image` Exit 0 (der Zug ändert
      Build-Kontext-Dateien).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update: [`harness/README.md`](../../../../harness/README.md) §Sensors
      (Zeile `make test-integration`: ein Satz zum neuen Rundlauf);
      gemeldete Träger fremder Dateien mit der Closure nachgezogen (§3 Suchlauf,
      [`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von der Slice-Closure selbst, weil die Roadmap unter *Offene Wellen* keine Welle führt (gemessen: `docs/plan/planning/` trägt keine flache Welle-Datei) und „die nächste Welle-Closure“ damit keine Adresse ist.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/application/port/apiauth/apiauth.go` (Arbeitsname) | neu | der gemeinsame Klassifikator: Rollen-Typ, Konstruktion aus zwei Token-Mengen, Zuordnung `Token → Rolle`; Vergleich zeitkonstant gegen alle Token beider Klassen (`crypto/subtle`) ohne Abbruch beim ersten Treffer; leerer Wert nie gültig, `admin` gewinnt |
| `internal/application/port/apiauth/apiauth_test.go` | neu | Tabellentest je Fall aus §2 (Happy: zwei Token; Boundary: Singular allein, derselbe Wert in beiden Klassen, Komma im Singular; Negative: entfernt, leer, unbekannt) — nach [`LH-FA-SST-012`](../../../../spec/lastenheft.md) |
| `internal/adapters/driving/http/middleware.go`, `server.go` | update | `withToken` ruft den gemeinsamen Klassifikator; lokale `role`/`classifyToken` und der Kommentar zur zweiten Fassung entfallen; `Config` trägt je Klasse die Menge |
| `internal/adapters/driving/grpc/interceptor.go`, `server.go` | update | Stream-Interceptor und `administrationRPCRoles` auf denselben Rollen-Typ; lokale Fassung und Kommentar entfallen |
| `internal/adapters/driving/http/*_test.go`, `internal/adapters/driving/grpc/*_test.go` | update | mechanische Migration der `Config`-Felder (am Parent 99 Trefferzeilen in 17 Dateien, §3 Suchlauf Zeile 1 zählt die Symbole) plus je ein Fall „zweites Token“ und „entferntes Token“ |
| `internal/bootstrap/wiring.go` | update | `Config`-Felder, Konstanten der zwei Plural-Variablen, `ConfigFromEnv`, Übergabe an beide Server; Kommentare, die „ein Token je Klasse“ nennen |
| `internal/bootstrap/config_file.go` | update | `mergeConfig` (zweiter Zugriffsweg), `forbiddenFileCredentialKeys` (7 → 9) samt Kommentar |
| `internal/bootstrap/config_file_internal_test.go`, `internal/bootstrap/wiring_test.go` | update | `TestConfigFromFileLehntZugangsdatenAb` auf neun, Mengengleichheit-Test, Parser-Fälle (leeres Element, Leerraum, nur Trennzeichen, leere Zeichenkette) für beide Zugriffswege |
| `internal/domain/messagecode/codes.go` (und der Test des Pakets, falls er die Menge zählt) | update | neuer(r) Code(s) der Klasse `configuration`, Zuweisungsregel in diesem Abschnitt |
| `tools/harness/run-integration-tests.sh` | update | neue Phase „Token-Wechsel“, `abdeckung_declare` für [`LH-FA-SST-012`](../../../../spec/lastenheft.md); jede neue `func TestE2E*` stünde im `-run` (Vollständigkeits-Test `test/integration/runner_vollstaendigkeit_test.go`) |
| `docs/user/e2e-abdeckung.md` | Erzeugnis | vom Runner geschrieben, committet |
| `docs/user/benutzerhandbuch.md` | update | §2 Liefer-Punkt 3 |
| `harness/README.md` | update | Zeile `make test-integration` |

- **Architect-Frage — Ort des gemeinsamen Klassifikators (Schichten-Entscheidung,
  Bestätigung vor dem Start, §4).** `.a-check.yml` führt `adapters` als eine
  Schicht ohne Kante von Adapter zu Adapter (gelesen am Parent); die zwei
  Fassungen von `classifyToken` stehen deshalb getrennt. **Vorschlag dieses
  Plans: ein Paket unter `internal/application/port/`** (Arbeitsname `apiauth`),
  weil (1) die Schicht `ports` (`internal/application/port/**`) bereits die
  Kante `adapters → ports` trägt, `.a-check.yml` also unverändert bleibt; (2) das
  Paket `internal/application/port/apifault` genau dieses Muster im Code
  vorgibt (gemessen: sein Kopf nennt „Beide Driving-Adapter fragen dieselbe
  Stelle“ und ordnet Fehler für HTTP und gRPC einheitlich zu); (3) die
  Token-Zuordnung eine Zugriffsrichtlinie der Driving-Seite ist, keine
  Domänenregel ([`ADR-0002`](../../adr/0002-abhaengigkeitsrichtung.md): die Domain
  kennt keine Treiber). **Offen:** keine ADR benennt das Muster `apifault`
  (gemessen: `git grep -l apifault -- docs/plan/adr` ohne Treffer) — die
  Platzierung hat Präzedenz im Code, keine Entscheidung; ein Architect
  bestätigt sie (oder wählt `internal/domain/…`) **vor** dem Start. Eine
  Kante in `.a-check.yml` wird in keinem Fall ergänzt (`adapters → adapters`
  bleibt verboten); `make a-check` Exit 0 ist DoD.
- **Zuweisungsregel der Meldungscodes (gilt für diesen und die zwei
  Folge-Slices).** Der neue Code der Klasse `configuration` ist die **nächste freie
  Nummer ab `PCF-E2008`** zum Zeitpunkt der Arbeit; der Implementer misst den
  Stand von `internal/domain/messagecode/codes.go` an seinem Arbeitsstand
  (`git grep -n 'PCF-E20' -- internal/domain/messagecode`); es gibt keinen
  Vorab-Block je Slice. Die drei Slices laufen in der Reihenfolge dieses Plans
  nacheinander, die Nummern kollidieren deshalb nicht
  ([`ADR-0144`](../../adr/0144-meldungscodes-nutzerseitige-kennungen.md): ein Code
  wird nie neu belegt). Ein Code für „leeres Element“ und „Leerraum im Element“
  gemeinsam ist zulässig (ein Fehler der Token-Liste, ein Klartext je Fall);
  die Zahl der Codes legt der Implementer fest.
- **Parser-Form.** Trennzeichen ist das Komma; jedes Element bleibt unverändert
  (kein `TrimSpace`: Leerraum ist ein Fehler, kein Rauschen). Der Singular wird
  nicht zerlegt. Die Plural-Variable wird in **beiden** Zugriffswegen
  (`ConfigFromEnv`, `mergeConfig`) von **einer** Funktion geparst — die zwei
  Wege dürfen nicht auseinanderlaufen
  ([`BEO-PGC/lese-doppelquelle`](../observations/BEO-PGC/lese-doppelquelle/observation.md)).
- **Vergleich.** Zeitkonstant ist eine *Erwartung* der Spec an die Umsetzung
  ([`SPEC-035`](../../../../spec/pflichtenheft.md)); `subtle.ConstantTimeCompare`
  liefert bei verschiedener Länge sofort 0 und verdeckt die Länge nicht. Ob
  beide Seiten vor dem Vergleich gehasht werden (gleiche Länge), entscheidet der
  Implementer und nennt die Wahl und die benannte Grenze im Kommentar der
  Funktion; eine Messung der Laufzeit gehört nicht in den Test.
- **Reihenfolge:** (1) Startmessung: Suchlauf unten am Arbeitsstand, `git grep`
  der zwei Fassungen; (2) Paket und Tabellentest; (3) Adapter und deren
  Tests migrieren, `make a-check`; (4) Konfiguration, Codes, Zugangsdaten-Klasse,
  Mengengleichheit-Test; (5) Mutationsläufe an Kopien; (6) Runner-Phase,
  `make image`, `make test-integration`; (7) Handbuch, README, `make gates`.
- **Suchlauf (§3.13 der Regeln, [`AGENTS.md`](../../../../AGENTS.md)).** Bewegte
  Eigenschaft: die **Token-Klassifikation** (Symbol `classifyToken`, Zählwort
  „zwei Token-Klassen“/„beiden Token-Klassen“, Beschreibung „zweite, wortgleiche
  Fassung“) und die **Zugangsdaten-Klasse** (Zählwort „sieben“, die Aufzählung
  `nats_stream_token`). Parent ist der Stand `a53f4e75…` (`git rev-parse HEAD`
  am Planungsstand, vor dem Plan-Commit; nie `HEAD`). Suchraum: ganzer Baum ohne
  `docs/reviews`, `done/`, `observations/`, `.harness/baseline`, die offenen
  Pläne und `docs/plan/adr` (Accepted ADRs halten ihren Stand,
  [`AGENTS.md`](../../../../AGENTS.md) §3.5). **Parent-Zeilen gemessen; die
  `diff`-Zeilen tragen den Planungsstand** (Arbeitsbaum gleich Parent, die
  Plan-Dateien sind ausgeschlossen) — der Implementer setzt sie nach der Arbeit
  auf den Endstand und liest jede Trefferzeile. **Erwartung am Endstand**
  (hergeleitet, nicht gemessen): Zeilen 1 und 2 nennen in `adapters/` keinen
  Treffer mehr; Zeile 4 hat „sieben“ nicht mehr; Zeile 3 verlangt von jedem
  Träger, der „zwei/beiden Token-Klassen“ im Sinn von „ein Token je Klasse“
  liest, eine Nachzug-Entscheidung (die Kommentare in `compose.yaml` und den
  Beispiel-Clients bleiben richtig: sie nennen die zwei Klassen, nicht die Zahl der
  Token).

```suchlauf
a53f4e75ca51bdaa4319aeabba96f89c48b6448f 27 -n -E 'classifyToken' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!docs/plan/planning/open' ':!.harness/baseline' ':!docs/plan/adr'
a53f4e75ca51bdaa4319aeabba96f89c48b6448f 4 -n -E 'zweite, wortgleiche Fassung|eigene Fassung derselben Zuordnung' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!docs/plan/planning/open' ':!.harness/baseline' ':!docs/plan/adr'
a53f4e75ca51bdaa4319aeabba96f89c48b6448f 17 -n -E 'beiden Token-Klassen|zwei Token-Klassen|zwei Rechtsklassen|beiden konfigurierten Klassen' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!docs/plan/planning/open' ':!.harness/baseline' ':!docs/plan/adr'
a53f4e75ca51bdaa4319aeabba96f89c48b6448f 9 -n -E 'sieben Schlüssel|sieben Zugangsdaten|nats_stream_token' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!docs/plan/planning/open' ':!.harness/baseline' ':!docs/plan/adr'
a53f4e75ca51bdaa4319aeabba96f89c48b6448f 131 -n -E 'CDC_API_TOKEN_(READER|ADMIN)' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!docs/plan/planning/open' ':!.harness/baseline' ':!docs/plan/adr'
diff 27 -n -E 'classifyToken' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!docs/plan/planning/open' ':!.harness/baseline' ':!docs/plan/adr'
diff 4 -n -E 'zweite, wortgleiche Fassung|eigene Fassung derselben Zuordnung' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!docs/plan/planning/open' ':!.harness/baseline' ':!docs/plan/adr'
diff 17 -n -E 'beiden Token-Klassen|zwei Token-Klassen|zwei Rechtsklassen|beiden konfigurierten Klassen' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!docs/plan/planning/open' ':!.harness/baseline' ':!docs/plan/adr'
diff 9 -n -E 'sieben Schlüssel|sieben Zugangsdaten|nats_stream_token' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!docs/plan/planning/open' ':!.harness/baseline' ':!docs/plan/adr'
diff 131 -n -E 'CDC_API_TOKEN_(READER|ADMIN)' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!docs/plan/planning/open' ':!.harness/baseline' ':!docs/plan/adr'
```

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `in-progress/` trägt keinen anderen Slice
(WIP-Limit 1), **und** ein Architect hat die Platzierung des gemeinsamen
Klassifikators (§3, Architect-Frage) bestätigt oder geändert — als committetes
Artefakt (Verdikt unter `docs/reviews/` oder Nachtrag in §3), nicht nur im
Gespräch
([`BEO-PGC/start-trigger-ohne-uebergabe-artefakt`](../observations/BEO-PGC/start-trigger-ohne-uebergabe-artefakt/observation.md)),
**und** der Implementer hat die Startmessung gefahren: `make suchlauf-nachmessen
PLAN=` mit diesem Plan endet am Arbeitsstand mit Exit 0 (die `diff`-Zeilen
tragen den Planungsstand; eine Abweichung seit dem Plan-Commit ist ein Befund
für den Planner, bevor Code entsteht).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): der gemeinsame
  Klassifikator verlangt eine neue Kante in `.a-check.yml` oder ein Paket der
  Domain gegen die Bestätigung (dann zuerst der Architect-Zug), oder der
  Liefer-Umfang wächst über die drei Punkte aus §2 (etwa eine nötige Änderung
  der SDK-Packages).
- `in-progress` → `open` (blockiert — Carveout?): die Bestätigung des
  Architect bleibt aus, die Spec 0.15.0 wird im Wortlaut von
  [`SPEC-035`](../../../../spec/pflichtenheft.md) geändert, oder der
  Compose-Stack ist am Messhost nicht startbar; kein Carveout, weil dann kein
  Gate rot ist.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

Die DoD aus §2 ist vollständig, `make gates` endet mit Exit 0 am Stand der
Closure, `make test-integration` und `make a-check` enden mit Exit 0,
`make suchlauf-nachmessen PLAN=` mit diesem Plan endet am Endstand mit Exit 0
(die `diff`-Zeilen auf den Endstand gesetzt, jede Trefferzeile gelesen), und die
Closure-Notiz in §7 trägt einen Lerneintrag (geschärfte Regel, neuer Sensor oder
benannte Spec-Lücke; Kandidaten: der Mengengleichheit-Test der
Zugangsdaten-Klasse als Sensor gegen die Halbänderung, die Festlegung zur leeren
Plural-Variable als benannte Spec-Lücke). Ein Gate, das am Stand der Closure rot
ist, geht nur mit dokumentiertem Carveout nach `done/`.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

Ausgangsform je Risiko: eingetreten (CO-NNN oder Folge-Slice) · entfallen
(Grund) · weiter offen (BEO-Eintrag im Register). Alle Ausgänge sind bis zur
Closure **offen** (Platzhalter `Ausgang: offen bis Closure`).

- **Die Zeitkonstanz der Prüfung ist nicht belegt.**
  [`ADR-0150`](../../adr/0150-tls-und-mehrfach-token.md) und
  [`SPEC-035`](../../../../spec/pflichtenheft.md) führen sie als Erwartung; kein
  Test misst Laufzeit, und `subtle.ConstantTimeCompare` verdeckt die Länge nicht.
  Der Reviewer liest den Vergleich (kein früher Abbruch), der Kommentar der
  Funktion nennt die Grenze. — **Ausgang:** offen bis Closure.
- **Halbänderung der Zugangsdaten-Klasse zwischen den zwei Slices** (Code und
  Test stehen bei neun, `SPEC-016` bei elf, bis der OTLP-Slice schließt): der
  Mengengleichheit-Test (§2) fängt eine Abweichung von Code und Test, nicht
  die von Handbuch und Code — die Paarung bleibt von Hand
  ([`ADR-0089`](../../adr/0089-feldmengen-paarung-kein-sensor-review-waechter.md)).
  — **Ausgang:** offen bis Closure.
- **Spec-Lücke zur Plural-Variable:** [`SPEC-035`](../../../../spec/pflichtenheft.md)
  sagt nicht, ob eine leere Zeichenkette als „ungesetzt“ oder als
  leeres Element gilt, und nicht, wie ein Token mit Komma in der Liste auszudrücken ist
  (es gibt keinen Weg: das Komma ist Trenner). Der Plan legt „leer = ungesetzt“
  fest und nennt die Grenze im Handbuch; eine Spec-Klarstellung wäre ein
  Zug des Planners/Architect, kein Zug dieses Slice. — **Ausgang:** offen bis
  Closure (Kandidat: benannte Spec-Lücke im Lerneintrag).
- **Zwei Konfigurationswege lesen verschieden:** `ConfigFromEnv` und
  `mergeConfig` bilden die Konfiguration getrennt; liest nur einer die
  Plural-Variable, bleibt das im Unit-Test eines Weges unsichtbar
  ([`BEO-PGC/adapter-unittest-verdeckt-bootstrap-luecke`](../observations/BEO-PGC/adapter-unittest-verdeckt-bootstrap-luecke/observation.md),
  fünf Evidenz-Dateien am Planungsstand gezählt). Gegenmaßnahme: Test je Weg und
  die Runner-Phase (der Container liest über `ConfigFromEnvAndFile`). —
  **Ausgang:** offen bis Closure.
- **Die neue Runner-Phase stört die Folge-Phasen** (Container mit Override
  neu erzeugt, Zustand muss am Phasen-Ende ohne Override wiederhergestellt sein;
  `run-integration-tests.sh` hat über 5000 Zeilen, Reihenfolge der Phasen
  ist Teil des Belegs). — **Ausgang:** offen bis Closure.
- **Das Handbuch zieht nicht mit** (Betreiber-Oberfläche wächst um zwei
  Variablen):
  [`BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`](../observations/BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche/observation.md)
  (drei Evidenz-Dateien am Planungsstand, also eine Lücke und keine Notiz) —
  deshalb ein eigener Liefer-Punkt in §2, nicht ein Anhang. — **Ausgang:**
  offen bis Closure.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<NNN>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks). Ging der Gegenstand an einen anderen Slice oder entfiel er, trägt
diese Sektion die Zeile `Gegenstand:` mit Kennung oder Grund und jedes Risiko
aus §6 seinen Ausgang; die Liefer-Punkte der DoD bleiben leer
(`modul-05-planning-harness.md` §Ein Slice, dessen Gegenstand ein anderer
übernimmt).

*Der Plan füllt diese Sektion nicht; sie wird bei der Closure vor dem
`git mv` nach `done/` geschrieben (§2: die Paarungs-Zeile nach dem `git mv`).*

- **Was hat funktioniert:** (bei der Closure zu füllen)
- **Was ging anders als geplant:** (bei der Closure zu füllen)
- **Steering-Loop-Eintrag:** (bei der Closure zu füllen; Kandidaten §5)
- **Beobachtungs-Register (`../observations/`):** (bei der Closure zu füllen;
  zu lesen sind die in §6 genannten Einträge)
- **Folge-Slices:** [`slice-tls-http-grpc-server`](../open/slice-tls-http-grpc-server.md)
  (TLS, nach diesem Slice) und
  [`slice-otlp-metrik-export`](../open/slice-otlp-metrik-export.md) (hebt die
  Zugangsdaten-Klasse von neun auf elf) — beides Dateien in `open/`
- **Risiken aus §6:** (bei der Closure zu füllen, je genau ein Ausgang)
- **Drei Paarungen:** (bei der Closure zu füllen: Anker · Folge-Slice · Register)

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Der Abschnitt selbst entfällt nie.** Die zwei vorgelagerten Prüfungen laufen
in **jedem** Slice-Plan — sie hängen weder am Modus noch am Slice-Typ. Bedingt
ist allein der Modus-Begründungsblock am Ende; deshalb nennt der Titel beide
Hälften.

**Vorgelagert — Sub-Area-Wahl prüfen:** [`harness/conventions.md`](../../../../harness/conventions.md)
§Modus-Deklaration führt **eine** Sub-Area (`*`, Kürzel `PGC`, Greenfield); die
berührten Pfade (`internal/…`, `tools/harness/`, `docs/user/`, `harness/`)
liegen alle in ihr, es entsteht keine neue Sub-Area — die Schwelle
(≥ 2 von 3 Achsen) wird nicht angewandt, weil kein Pfad ausdifferenziert werden muss.

**Vorgelagert — offene Beobachtungen sichten:** Register
(`docs/plan/planning/observations/BEO-PGC/`) am Planungsstand durchgegangen;
Treffer für die berührte Sub-Area, Zähler = Zahl der `evidence/`-Dateien:
`adapter-unittest-verdeckt-bootstrap-luecke` (5, erreicht 3× — Gegenmaßnahme
im Plan: Runner-Phase und Container-Start, §2/§6),
`handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche` (3, eigener
Liefer-Punkt), `handbuch-versionshistorie-uebersprungen` (3, Historien-Zeile im
Liefer-Punkt 3), `zwei-quellen-drift-handbuch-gegen-pflichtenheft` (3, Handbuch
und Pflichtenheft werden für die Token-Form gegeneinander gelesen),
`lese-doppelquelle` (3, ein Parser für beide Konfigurationswege, §3),
`negativtest-ohne-bindung-an-seine-eingabe` (25, Gegenprobe je Negativfall,
§2), `kommentar-herkunft-als-kette` (4, höchstens eine Kennung je Kommentar im
neuen Code, [`AGENTS.md`](../../../../AGENTS.md) §3.7),
`start-trigger-ohne-uebergabe-artefakt` (1, Start-Trigger §4). Die Einträge mit
mindestens drei Dateien sind **mit** diesem Slice berührt und tragen deshalb
hier ihre Gegenmaßnahme.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

**Modus:** alle berührten Sub-Areas GF (`*`/`PGC`, Greenfield).
