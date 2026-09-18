# ADR-0071: Coverage-Gate — Messgegenstand ist die netzlos prüfbare Fläche — Supersedes ADR-0054 (nur drei Klauseln)

**Status:** Accepted — Supersedes [`ADR-0054`](0054-coverage-gate-und-benchmark-infrastruktur.md)
in genau **drei** Klauseln aus deren §(a): dem **Scope-Bullet** (die
Definition des Messgegenstands), dem **Schwelle-Bullet** (der Gegenstand der
Endstufe 80 %) und der **Fitness-Function-Zeile**. Alles Übrige aus
`ADR-0054` — die vierte Docker-Stage `coverage`, das Gate-Skript
`tools/coverage-gate.sh`, das Suppression-Vollverbot (`AGENTS.md` §3.2), die
**Eskalationsklausel als Mechanismus** (bootstrap-aware Rampe), die
Abgrenzung `test/integration/` und der gesamte §(b) Benchmark-Teil — bleibt
unverändert bestehen und wird hier nicht wiederholt.

**Datum:** 2026-09-15

**Autor:** pt9912 (Architect-Rolle; anderer Kontext als der Planner-Lauf, der
die 80-%-Frage aufbereitet und die Zahlen gemessen hat, und als die
Implementer-/Reviewer-/Verifier-Läufe von `slice-049`/`slice-076`) <!-- d-check:status-provenance -->

**Bezug:** [`ADR-0054`](0054-coverage-gate-und-benchmark-infrastruktur.md)
(in drei Klauseln korrigiert), [`ADR-0030`](0030-testpyramide.md)
(Testpyramide — Unit-Tier vs. Integrations-Tier), `AGENTS.md` §3.1
(Docker-only), §3.6 (Schwellen-Änderung nur per ADR), §3.7 (Ist-Zustand),
`harness/sensors/coverage-gate.md` (die Sensor-Definition mit dem
Widerspruch), `harness/mk/coverage.mk` (`THRESHOLD`), `Makefile`
(`test-store`/`test-replication`/`test-notify`/`test-integration`),
`.github/workflows/e2e.yml` (nicht-blockierender Träger der DB-gestützten
Tier-Messung), der Architect-Verdikt dieses Zugs (Coverage-Gate-Messgegenstand),
`docs/plan/planning/in-progress/slice-076-coverage-gate-reifestufe-40.md` <!-- d-check:status-provenance -->
(dessen Kalibrierung dieser Schnitt neu bemisst).

**Schärft:** — (Prozess-/Tooling-ADR ohne Spec-Stratum, wie `ADR-0054`)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

`ADR-0054` §(a) legt in **einem** Absatz zweierlei fest: die Endstufe
**80 %** („Ziel-Endstufe 80 % (Nutzer-Entscheidung)") **und** den
Messumfang (`-coverpkg` über `./internal/... ./cmd/...`). Derselbe Absatz
sagt über diesen Umfang, die DB-/Replication-Adapter-Tests „skippen aber
ohne gesetzte `CDC_*_TEST_DSN`-Variable … der Coverage-Lauf im netzlosen
Docker-Build zählt sie nicht gegen die Schwelle, ohne den Build zu
brechen".

Das ist ein **Selbstwiderspruch**: Ein Messverfahren, das eine Code-Menge
strukturell nie ausführt, kann eine prozentuale Endstufe über genau diese
Menge nicht erreichen. Dieser Zug hat ihn nachgemessen (gepinntes
Toolchain-Image `golang:1.27-alpine@sha256:cf6fca66…`, eigener Bau der
`coverage`-Stage, eigene Deduplizierung über die Block-Position):

| Größe | Wert |
|---|---|
| Statements gesamt (`internal/…`+`cmd/…`) | **2467** |
| davon netzlos gedeckt | **1215** (49,25 %; das Gate meldet 49,30 %, der Planner 49,33 % — die Spanne ist die dokumentierte Lauf-zu-Lauf-Schwankung) |
| Statements in den drei DB-gestützten Paketen (`postgresstorage` 610, `replication/receive` 155, `postgresack` 23) | **788** |
| Decke des heutigen Verfahrens (100 % von allem außer den drei Paketen) | **1679 / 2467 = 68,06 %** |

**80 % > 68,06 %.** Die Endstufe ist unter dem eigenen Umfang unerreichbar;
die Kalibrierungs-Bindung („bis 80 % erreicht ist") führt eine Zahl, die das
Verfahren nicht liefern kann. Der Fall ist strukturell derselbe wie
`ADR-0069`/`ADR-0070` für `ADR-0062` Punkt 3: eine `Accepted`-ADR trägt
einen Satz, der gegen die reale Messung unwahr ist, und die Korrektur ist
eine Folge-ADR mit `Supersedes` (`AGENTS.md` §3.5).

Zwei Nebenbefunde dieses Zugs (nachgeprüft, nicht übernommen): SQL ist
bereits ausgelagert (`postgresstorage/queries` ist reiner SQL-Text ohne
ausführbare Statements, `postgresstorage/mapper` ist mit 80 % gedeckt); was
die DB-Adapter unprüfbar macht, ist allein die **Naht zur Ausführung** — sie
hängen am konkreten Typ (`store.go: pool *pgxpool.Pool`), es gibt kein
Interface. Der Methodenrumpf (Query absetzen, Rows iterieren, `defer Close`,
Scan über den Mapper, Fehler klassifizieren) ist deshalb nur mit lebender
PostgreSQL erreichbar.

Der Anlass ist eine Nutzerfrage: „Soll dieses Gate messen, was netzlos
prüfbar ist, oder den ganzen Baum?" Der Planner hat geantwortet, 80 % sei
„unter der eigenen Definition" unerreichbar. Dieser Zug bestätigt die
Aussage und entscheidet, welche Definition gilt.

## Entscheidung

Wir wählen: **Der Messgegenstand des Coverage-Gates ist die netzlos
prüfbare Fläche** — und die davon ausgenommene DB-gestützte Fläche bekommt
eine **eigene, subjekt-qualifizierte Messung**.

1. **Messgegenstand.** Das Gate misst `./internal/... ./cmd/...` **ohne**
   die Pakete, deren Tests ohne externen Dienst real überspringen —
   namentlich `internal/adapters/driven/postgresstorage` (ohne das
   Unterpaket `mapper`), `internal/adapters/driven/postgresack` und
   `internal/adapters/driving/replication/receive`. Die tragende Regel ist
   die Eigenschaft, nicht die Liste: **Ein Paket, dessen Testlauf einen
   externen Dienst voraussetzt, ist nicht Gegenstand dieses Gates.** Der
   Nenner ist damit 1679 Statements, die Decke 100 %.

2. **Die Endstufe 80 % bleibt** — sie gilt ab hier über der netzlos
   prüfbaren Fläche. Dieselbe gedeckte Menge (1215 Statements) über den
   kleineren Nenner ergibt **69,74 %** als Ist-Stand; die Rampe wird auf
   diesen Nenner neu kalibriert (Einstieg = real gemessener Ist-Stand,
   abgerundet auf die nächste volle 5-%-Stufe, nach dem unveränderten
   Mechanismus aus `ADR-0054` §(a)). Das ist **keine** Schwellen-Senkung
   (§3.6): die Endstufe bleibt 80 %, der Gegenstand wird präzise, der Wert
   steigt real von 49,3 % auf 69,7 % — nicht weil weniger verlangt, sondern
   weil der Nenner den Verfahrensbereich beschreibt.

3. **Die DB-gestützte Ebene bekommt ihre eigene Messung.** Die drei
   ausgenommenen Pakete werden dort gemessen, wo ihre Tests ohnehin gegen
   einen echten PostgreSQL laufen — `make test-store` und
   `make test-replication` —, mit eigenem `-coverprofile` und eigener
   Schwelle (**DB-Adapter-Coverage**). Diese Messung ist **nicht** Teil von
   `make gates`: das Gate läuft bei jedem Commit, ein PostgreSQL-Container
   dort wäre der von Option C verworfene Preis. Ihr Träger ist der
   **nicht-blockierende** `.github/workflows/e2e.yml`-Workflow, in dem
   `make test-integration` bereits läuft. Der Schwellen-Wert wird beim ersten
   Lauf real gemessen und nach demselben bootstrap-aware-Muster kalibriert
   wie der Unit-Wert; die Eskalationsklausel aus `ADR-0054` §(a) gilt dafür
   unverändert weiter.

4. **„Die Coverage" des Repos ist die Gate-getragene Zahl.** Wenn dieses
   Repo ohne Subjekt-Zusatz von „der Coverage" spricht, meint es die
   **Unit-Tier-Zahl** (`make coverage-gate`, in `make gates`) über der
   netzlos prüfbaren Fläche — die einzige mit einer blockierenden Schwelle.
   Die DB-Adapter-Coverage trägt ihren Subjekt-Namen **immer** mit und wird
   nie „die Coverage" genannt. Zwei Zahlen, zwei Gegenstände, ein Name je
   Gegenstand — keine Doppelquelle für denselben Zustand (Modul 6).

5. **Die Naht (Option D) ist zulässig, aber nicht Teil dieser Entscheidung.**
   Eine schmale Executor-Schnittstelle zieht netzlos prüfbare Logik
   (Fehlerklassifikation, Row-Übersetzung) aus den DB-Adaptern in die
   prüfbare Fläche und hebt die Decke ohne Messänderung. Ihre Rechtfertigung
   ist **Design** — schmale Abhängigkeit statt konkreter `*pgxpool.Pool`,
   Fehlerklassifikation als reine Funktion —, **nicht** die Zahl; dass der
   Wert danach steigt, ist Folge, nicht Zweck. Sie gehört als eigener,
   design-begründeter Vorgang in `next/`, nicht als Beigabe zum Schnitt aus
   Punkt 1.

### Was dieses Gate ist: kein Proxy für System-Korrektheit

Das Coverage-Gate ist eine **Messung der netzlos prüfbaren Fläche**, kein
Proxy für System-Korrektheit. Die Antwort auf die Ausgangsfrage folgt aus
dem Verfahren selbst: Ein netzloser Docker-Lauf kann nur beurteilen, was
netzlos läuft. Die DB-gestützte Korrektheit hat ihren eigenen
Beleg-Träger (die Integrations-Ebene) und ab Punkt 3 ihre eigene Zahl; sie
ist **nicht** durch eine hohe Unit-Coverage-Zahl vertreten. Wer das Gate als
System-Korrektheits-Proxy läse, müsste Option C wählen und bei jedem Commit
einen PostgreSQL-Container bezahlen.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — Umfang behalten, Endstufe neu auf die Decke (~68 %) schneiden | ein ADR, kein Code; die Zahl bleibt ein Anteil **des ganzen** Produkts | das Gate misst dauerhaft eine Fläche, die sein Verfahren nicht ausführt; 32 % Leergewicht im Nenner dämpfen jedes Signal (eine neu gedeckte Use-Case-Zeile bewegt die Zahl um einen Bruchteil dessen, was sie bei bereinigtem Nenner bewegte); die Decke wandert mit jedem neuen DB-Adapter-Code, das „Ziel" wäre dauerhaft eine bewegte Zahl |
| B — Umfang auf die netzlos prüfbare Fläche schneiden; Endstufe 80 % behalten (**gewählt, mit F**) | die Aussage „80 %" bekommt einen präzisen, erreichbaren Gegenstand; das Signal ist scharf (der Nenner trägt kein Leergewicht); §3.6-konform (Endstufe unverändert, Wert steigt real auf 69,7 %); der Nutzer-Wunsch „80 %" wird erreichbar | die DB-Adapter verlieren ihre (scheinbare) Zugehörigkeit zur Unit-Zahl und brauchen eine eigene Zusage (→ F); die Rampe ist neu zu kalibrieren |
| C — DB-Tests in den Coverage-Lauf holen (`make gates` mit PostgreSQL) | 80 % wird erreichbar, ohne Produktivcode zu ändern; eine einzige Zahl bleibt | `make gates` läuft **bei jedem Commit** und bekäme einen PostgreSQL-Container — Laufzeit, Ressourcen, eine neue Fehlerquelle im Gate, das gerade schnell und netzlos ist; widerspricht dem netzlosen Geltungsbereich von `ADR-0054` §(a) |
| D — Naht einziehen (Executor-Schnittstelle) + Logik pur ziehen | hebt die Decke **ohne** Messänderung; bessert das Design (schmale Abhängigkeit, Fehlerklassifikation als reine Funktion) | echte Refaktorierung von acht Adapterdateien — Produktivcode, den der Coverage-ADR nie im Blick hatte; als **eigener** Zweck „die Zahl" wäre sie Goodhart (die Zahl als Ziel verdrängt das Design als Grund) — deshalb nicht Teil dieser Entscheidung, sondern ihr eigener, design-begründeter Vorgang |
| E — nichts tun | kein Aufwand | die Endstufe bleibt eine aspirierte Zahl, die Kalibrierungs-Bindung („bis 80 % erreicht ist") bleibt unwahr, und der Widerspruch in einer `Accepted`-ADR bleibt stehen — `AGENTS.md` §3.5 verlangt für seine Auflösung eine Folge-ADR, nicht ein Liegenbleiben |
| F — eigene, zweite Schwelle für die DB-gestützte Testebene (**gewählt, mit B**) | jede Ebene bekommt ihre eigene, erreichbare Zusage (Unit: prüfbare Fläche; DB-Adapter: gegen echte PG); `make gates` bleibt schnell (Träger ist `e2e.yml`, nicht-blockierend); schließt die von B gerissene Lücke — die DB-Adapter bleiben nicht ohne eigene Zusage | eine zweite Messung und eine zweite Schwelle statt einer; ohne die Aussage aus §Entscheidung Punkt 4 wären zwei Zahlen eine Doppelquelle für „die Coverage" — diese Aussage ist deshalb Teil der Entscheidung, nicht Beiwerk |

**Fazit:** B und F **zusammen**. B allein ließe die DB-Adapter ohne eigene
Zusage; F allein ließe den Unit-Nenner verschmutzt. Die Kombination schneidet
den Unit-Gegenstand scharf **und** gibt der ausgenommenen Fläche ihre eigene
messbare Bindung, ohne `make gates` zu belasten.

## Konsequenzen

- Positiv: Die 80-%-Endstufe ist ab hier erreichbar und der Satz „bis 80 %
  erreicht ist" in der Kalibrierungs-Bindung wird wahr; §3.6 ist unverletzt
  (keine Senkung — die Endstufe bleibt, der Gegenstand wird präzise).
- Positiv: Das Unit-Signal wird scharf — der Nenner trägt kein Leergewicht
  mehr, jede neu gedeckte Zeile bewegt die Zahl um ihren realen Anteil.
- Positiv: Die DB-Adapter verlieren ihre scheinhafte Zugehörigkeit nicht
  ohne Ersatz: ihre Messung existiert (Punkt 3), subjekt-qualifiziert und
  ohne den netzlosen Lauf zu belasten.
- Negativ mit Grenze: Es gibt ab hier **zwei** Coverage-Zahlen. Die
  Doppelquelle ist durch §Entscheidung Punkt 4 geschlossen — „die Coverage"
  ohne Zusatz ist die Unit-Zahl; die zweite trägt immer ihr Subjekt. Diese
  Namens-Disziplin ist Prosa, kein Sensor; sie kann driften, wenn ein
  künftiger Lauf die Kurzform „die Coverage" auf die falsche Zahl bezieht.
- Negativ: Die Rampe (`THRESHOLD` in `harness/mk/coverage.mk`) ist gegen
  `slice-076`s Kalibrierung (Einstieg 40 %, auf dem alten Nenner) neu zu <!-- d-check:status-provenance -->
  bemessen; zwischen Schnitt und Neukalibrierung ist die alte Stufe real
  trivial grün. Das ist ein Übergangszustand, kein Defekt.
- Folgepflicht (Implementer-Zug, Tooling/Doku, **kein** Produkt-Code):
  `-coverpkg`/Profil-Filter ohne die drei Pakete, Neukalibrierung der Rampe,
  `harness/sensors/coverage-gate.md` (§Vertrag, §Kalibrierungs-Bindung,
  §Grenze), `harness/README.md` §Sensors-Zeile, `AGENTS.md` §4-Zeile; Grün-
  und Rot-Beleg. Größe: ein Slice, kein Produkt-Code.
- Folgepflicht (Implementer-Zug, Tooling): DB-Adapter-Coverage-Messung in
  `test-store`/`test-replication` mit `-coverprofile`, eigenes Target,
  Träger `.github/workflows/e2e.yml`, Erstkalibrierung. Eigener, wellenlos
  lieferbarer Slice.
- Folgepflicht (Planner-Zug): die Wave zum Erreichen der 80 % auf der
  netzlos prüfbaren Fläche schneiden (Schneide-Vorschlag: §Fitness Function
  und Verdikt §Folgearbeit).
- Offen (kein Folge-Slice hier): Option D bleibt als design-begründeter
  Vorgang in `next/` — nicht Teil dieser Entscheidung.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| `docker build --target coverage` (gepinntes Toolchain-Image `golang:1.27-alpine@sha256:cf6fca66…`) + `go tool cover -func` | **Messgegenstand:** keine Block-Position im Profil liegt in `postgresstorage` (ohne `mapper`), `postgresack` oder `replication/receive`; **Gesamt-Coverage** der verbleibenden Fläche ≥ `THRESHOLD` (Endstufe 80 %) | `make coverage-gate` (in `make gates`) |
| `make test-store` / `make test-replication` mit `-coverprofile` (gepinnte Toolchain-/PG-Digests aus `Makefile`) | **DB-Adapter-Coverage** ≥ ihre eigene Schwelle (Erstkalibrierung durch den Träger-Slice) | neues Target (geplant), getragen von `.github/workflows/e2e.yml` |

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

Beobachtbare Trigger: **(a)** ein Paket kommt hinzu oder fällt weg, dessen
Testlauf einen externen Dienst voraussetzt (die Regel aus §Entscheidung
Punkt 1 greift dann ohne Textänderung, aber die namentliche Liste ist
nachzuziehen); **(b)** die Naht (Option D) wird gezogen — dann wandert
netzlos prüfbare Logik aus den ausgenommenen Paketen in die prüfbare Fläche,
der Nenner wächst und die Rampe ist neu zu bemessen (Folge-ADR); **(c)** die
Integrations-Ebene erhält einen anderen Träger als `e2e.yml` (etwa blockierend)
— dann ist die Stellung der DB-Adapter-Coverage neu zu beurteilen. Andernfalls
permanent.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-15 | Accepted — Anlass: Nutzerfrage „netzlos prüfbare Fläche oder ganzer Baum?"; der Selbstwiderspruch der Endstufe 80 % (Decke 68,06 %) real nachgemessen; korrigiert `ADR-0054` §(a) Scope-, Schwelle- und Fitness-Function-Klausel | der Architect-Verdikt dieses Zugs (Coverage-Gate-Messgegenstand), eigene Messung (`coverage`-Stage, gepinntes Toolchain-Image) |
| 2026-09-18 | Zitat-Korrektur — `docs/reviews/**`-Pfade durch Kennung ersetzt (`ADR-0073`) | PENDING_COMMIT |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0071` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
