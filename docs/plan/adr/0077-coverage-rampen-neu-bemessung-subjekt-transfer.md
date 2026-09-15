# ADR-0077: Coverage-Rampen — Neu-Bemessung bei Subjekt-Transfer

**Status:** Accepted — **kein** Supersedes. Diese ADR führt den
Re-Evaluierungs-Trigger (b) der
[`ADR-0071`](0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
aus (die Naht ist gezogen; „die Rampe ist neu zu bemessen") und ergänzt den
dort beschriebenen Mechanismus um seine **abfließende** Hälfte. Der Text der
[`ADR-0071`](0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
bleibt unverändert (`AGENTS.md` §3.5); ihre vier Festlegungen zum
Messgegenstand und zum Subjekt-Namen werden hiermit **bestätigt** und nicht
wiederholt.

**Datum:** 2026-09-15

**Autor:** pt9912 (Architect-Rolle; anderer Kontext als der Implementer-Lauf,
der die Naht gezogen und die Zahlen gemessen hat — `slice-081`, Branch <!-- d-check:status-provenance -->
`slice-081-executor-naht`) <!-- d-check:status-provenance -->

**Bezug:**
[`ADR-0071`](0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
(deren Trigger (b) eintritt),
[`ADR-0054`](0054-coverage-gate-und-benchmark-infrastruktur.md) (§(a) — die
Eskalationsklausel, hier um die Senkungs-Hälfte ergänzt),
[`ADR-0030`](0030-testpyramide.md) (Testpyramide — Unit-Tier vs. DB-gestützter
Tier), `AGENTS.md` §3.5 (Accepted-ADRs sind immutable) · §3.6
(Schwellen-Änderung nur per ADR) · §3.7 (Ist-Zustand),
`harness/sensors/coverage-gate.md`,
`harness/sensors/db-adapter-coverage.md`, `harness/mk/coverage.mk`
(`THRESHOLD`), `tools/harness/db-coverage.sh` (`DB_COVERAGE_THRESHOLD`),
`Dockerfile` (Stufe `coverage`), `.github/workflows/e2e.yml`
(nicht-blockierender Träger der DB-gestützten Messung),
`docs/plan/planning/in-progress/slice-081-executor-naht.md` <!-- d-check:status-provenance -->
(der auslösende Vorgang),
`docs/plan/planning/observations/BEO-PGC/endstufe-unter-eigenem-messgegenstand-unerreichbar/observation.md`
(die verwandte, bereits aufgelöste Klasse)

**Schärft:** — (Prozess-/Tooling-ADR ohne Spec-Stratum, wie
[`ADR-0054`](0054-coverage-gate-und-benchmark-infrastruktur.md) und
[`ADR-0071`](0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md))

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

**(1) Die zwei Messungen partitionieren denselben Baum.** `ADR-0071`
Punkt 1–4 teilt `./internal/... ./cmd/...` in **zwei** Gegenstände: die
**netzlos prüfbare Fläche** (`make coverage-gate`, Träger `make gates`) und
die drei Pakete, deren Testlauf einen externen PostgreSQL voraussetzt —
`postgresstorage` ohne das Unterpaket `mapper`, `postgresack`,
`replication/receive` —, die ihre eigene, subjekt-qualifizierte
**DB-Adapter-Coverage** tragen (`make test-store`/`make test-replication`,
Träger der nicht-blockierende `.github/workflows/e2e.yml`). Die beiden
Gegenstände **überlappen nicht** und sind zusammen der ganze Baum: die Summe
ihrer Statements ist konstant.

**(2) Eine Naht ist per Konstruktion ein Transfer zwischen den zwei
Gegenständen.** `ADR-0071` Punkt 5 nennt die schmale Executor-Schnittstelle
(Option D) zulässig und **design**-begründet — schmale Abhängigkeit statt
konkreter `*pgxpool.Pool`, Fehlerklassifikation und Zeilen-Übersetzung als
reine Funktionen —, und ihr Re-Evaluierungs-Trigger (b) sagt voraus, dass mit
dem Ziehen der Naht netzlos prüfbare Logik aus den ausgenommenen Paketen in
die prüfbare Fläche wandert und die Rampe neu zu bemessen ist.
`slice-081` hat die Naht gezogen. <!-- d-check:status-provenance -->

**(3) Die realen Messungen dieses Zugs** (Implementer-Lauf auf dem Branch;
vom Architect nachgelesen, nicht aus zweiter Hand — die einzige `FAIL`-Zeile
des gesamten `make test-replication`-Logs ist das Coverage-Verdikt, **kein**
Testfehler):

| Größe | vor der Naht | nach der Naht |
|---|---|---|
| Unit-Gegenstand (Statements) | 1679 | 1817 (**+138**) |
| Unit-Ist (von der Stufe gedruckt) | 69,70 % | 71,30 % |
| DB-Gegenstand (Statements) | 788 | 650 (**−138**) |
| DB-Ist | 75,25 % (593/788) | 73,38 % (477/650) |
| DB-Einstieg der Rampe | 75 % | 75 % (unverändert) |

**(4) Die Arithmetik schließt sich, und sie ist der Kern.** Die 138
Statements, die den DB-Gegenstand verlassen, waren zu **116 gedeckt** (84 %).
Über einem Gegenstand, dessen Ist 75,25 % ist, liegen 84 % deutlich über dem
Durchschnitt: der Nenner schrumpft um 138, der Zähler nur um 116 — der Anteil
**fällt**, obwohl kein Verhalten sich ändert. Der Gegenstands-Schnitt aus
`ADR-0071` Punkt 1 hat denselben Effekt in **einer** Richtung behandelt (aus
dem Unit-Nenner entfernen); die **gespiegelte** Richtung — ein Schnitt, der
einen Gegenstand **entleert** — hat die Rampe nicht.

**(5) Die Dränage ist systematisch, nicht einmalig; und die Rampe kennt die
Richtung nicht.** Der Befund folgt aus der Partition selbst: **jeder** weitere
Transfer zwischen den Gegenständen wiederholt ihn. `slice-081` §4 nennt zwei <!-- d-check:status-provenance -->
Folge-Vorgänge je Paket (`postgresack`, `receive`), deren Naht jeweils weitere
netzlos prüfbare Logik aus dem DB-Gegenstand zieht.
Und der bootstrap-aware Mechanismus aus `ADR-0054` §(a) ist ein **Hochschalt**-
Mechanismus: Einstieg = gemessener Ist-Stand, abgerundet auf die nächste volle
5-%-Stufe, danach **aufwärts** bis zur Endstufe. Eine **abwärts** bewegte
Neu-Bemessung hat er nicht — die Dränage trifft eine Rampe, die nur in eine
Richtung läuft.

**(6) Die verwandte Klasse ist bereits registriert.**
`BEO-PGC/endstufe-unter-eigenem-messgegenstand-unerreichbar` (1×, direkt
aufgelöst durch `ADR-0071`) führt den Fall „eine Zielzahl, die ihr eigener
Messgegenstand strukturell nicht erreichen kann". Der Fall hier ist
**verwandt, aber verschiedenen Grades**: die Endstufe 80 % bleibt erreichbar
(73,38 % *kann* durch mehr Arbeit steigen); **stale** ist allein die
**Kalibrierung** — der Einstieg 75 % stand auf einem Ist-Stand, den die Naht
fortgenommen hat. Die Antwort ist deshalb nicht erneut eine
Gegenstands-Entscheidung, sondern eine Neu-Bemessung.

## Entscheidung

Wir wählen: **Bei einem Transfer zwischen den Messgegenständen werden die
Rampen neu bemessen — die abfließende nach unten auf ihren neu gemessenen
Ist-Stand, die aufnehmende nach oben nach dem bestehenden Mechanismus. Die
Endstufen (je 80 %) bleiben unverändert.** Vier Festlegungen:

1. **Der Träger der Neu-Bemessung ist der Transfer-Nachweis, nicht der
   Wunsch.** Neu bemessen wird nur, wenn Code den Gegenstand **gewechselt**
   hat — nicht, wenn die Quote gefallen ist. Der Nachweis ist die
   **Nenner-Paarung**: die Summe der Statements beider Gegenstände bleibt
   konstant, während sich die Verteilung verschiebt (der aufnehmende
   Gegenstand wächst um k, der abfließende schrumpft um k). Genau diese zwei
   Zahlen führen die Messungen selbst (`db-coverage.sh` druckt „gedeckt X von
   Y Statements"; die `coverage`-Stufe trägt ihre Statement-Zahl). **Fällt die
   Quote bei unverändertem Nenner, ist es keine Dränage, sondern eine
   Regression — dann steht die Schwelle.** Diese Bedingung ist es, die die
   Neu-Bemessung von einem Freibrief scheidet.

2. **Die Endstufen bleiben; nur der Einstieg bewegt sich.** Beide Endstufen
   stehen bei **80 %**. Neu bemessen wird der **Einstieg**: die DB-Adapter-
   Coverage von 75 % auf **70 %** (neu gemessener Ist-Stand 73,38 %, nächste
   volle 5-%-Stufe darunter) und die Unit-Coverage von 65 % auf **70 %** (neu
   gemessener Ist-Stand 71,30 %, nächste volle 5-%-Stufe; das ist der
   bestehende Hochschalt-Trigger, hier durch denselben Transfer ausgelöst).
   Das ist **keine** Endstufen-Senkung — `AGENTS.md` §3.6 bleibt unverletzt;
   der Einstieg hängt am real gemessenen Ist-Stand, und genau das ist der
   bootstrap-aware Mechanismus beider Rampen.

3. **Die Regel ist dauerhaft; ihre Anwendung ist keine neue Entscheidung.**
   Der Mechanismus wird **hier einmal** entschieden (diese ADR ist der
   §3.6-Träger für die Senkungs-Hälfte, die `ADR-0054` §(a) nicht hatte). Ein
   späterer Transfer, der ihn anwendet **und den Nenner-Nachweis führt**, ist
   die **Anwendung** der Regel, nicht eine neue Schwellen-Entscheidung. Eine
   Anwendung **ohne** den Nachweis ist keine.

4. **Der Gegenstand wird nicht neu definiert.** Der DB-Gegenstand bleibt an
   seinen Paketen verankert; ein Transfer verschiebt ihn, und die Zahl sagt
   die Wahrheit über den **verbliebenen** Code (siehe unten).

### Warum die Schwellen-Antwort trägt — und keine Gegenstands-Neudefinition

Die Frage ist berechtigt: eine Quote, die auf einen reinen Struktur-Umbau
ohne Verhaltens-Änderung fällt, sieht nach Goodhart aus. Vier Gründe, warum
die Neu-Bemessung trotzdem die richtige und die Neudefinition die falsche
Antwort ist:

- **Der gefallene Wert ist wahr.** 73,38 % ist die gegen echten PostgreSQL
  gemessene Deckung des **verbliebenen** DB-Adapter-Codes. Die Naht hat das
  Verhalten nicht verschlechtert, aber der verbliebene, DB-nahe Code ist
  tatsächlich weniger gedeckt als der Durchschnitt, auf dem der alte Einstieg
  stand. Die Zahl **beschreibt** den Rest; sie versteckt ihn nicht.
- **Eine transfer-invariante Neudefinition ist nicht ehrlich zu haben.** Jeder
  Trick, der die Zahl unbewegt lässt, ist einer von drei Fehlern: (a) er zählt
  netzlos geprüften Naht-Code in eine Zahl, die „gegen echten PostgreSQL
  gemessen" behauptet — das ist die Subjekt-Lüge, die `ADR-0071` Punkt 4
  schließt; (b) er lässt die zwei Gegenstände **überlappen** — er reißt die
  saubere Partition wieder auf, die Punkt 4 gegen die Doppelquelle verteidigt;
  oder (c) er speichert einen Basiswert — eine zweite Quelle für denselben
  Zustand, die Modul 6 verwirft. **Es gibt keine vierte, die invariant *und*
  wahr *und* überlappungsfrei ist**, weil der Fall (84 % gegen 75 %) intrinsisch
  ist: er kehrt mit dem Vorzeichen des Transfers wieder, egal wie man den
  Gegenstand zieht.
- **Die Neu-Bemessung ist billiger und die Fortsetzung des vorhandenen
  Verfahrens.** Sie ist eine Regel plus zwei Konstanten in bereits vorhandenen
  Trägern; ihr Kern — „Einstieg = real gemessener Ist-Stand" — steht schon im
  Mechanismus. Es fehlte nur die zweite Richtung, nicht ein neuer Gegenstand.
- **Sie ist sicherer, als sie aussieht.** Der Nenner-Nachweis (Festlegung 1)
  ist genau das, was eine Regression **nicht** erzeugen kann: wer Tests
  löscht, senkt den Zähler bei **unverändertem** Nenner — dann greift die
  Regel nicht.

### Was diese ADR nicht entscheidet

- **Nicht das Design der zwei Folge-Nähte** (`postgresack`, `receive`). Sie
  sind je ein eigener Vorgang (`ADR-0071` Punkt 5; `slice-081` §4); ihre <!-- d-check:status-provenance -->
  Schnittstelle ist eine andere (`*pgconn.PgConn` statt `*pgxpool.Pool`) und
  wird von dieser ADR nicht vorgezeichnet.
- **Nicht, ob die Naht hätte gezogen werden dürfen** — das ist mit `ADR-0071`
  Punkt 5 und `slice-081` entschieden. <!-- d-check:status-provenance -->
- **Nicht die Endstufen** (beide bleiben 80 %) und **nicht die Definition des
  DB-Gegenstands** (Festlegung 4).
- **Nicht der Lifecycle von `slice-081`.** Ob er nach dieser Entscheidung <!-- d-check:status-provenance -->
  abschließbar ist oder zurückgeführt gehört, ist Planner-Arbeit (Modul 5);
  die Implikation steht unter §Konsequenzen.
- **Kein Carveout.** Die Neu-Bemessung löst den roten Gate-Status auf; es
  bleibt kein roter Lauf, der eine `CO-<NNN>` bräuchte.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — Einstiege neu bemessen (**DB 75 → 70**, **Unit 65 → 70**) **plus** die Nenner-gebundene Neu-Bemessungs-Regel dauerhaft (**gewählt**) | der rote Lauf wird grün, ohne ein Verhalten zu ändern; die Endstufen bleiben (§3.6 unverletzt); die Regel ist **wiederholbar** — sie trifft `postgresack`/`receive` mit; der Nachweis bindet sie an einen Transfer, nicht an einen Wunsch | eine **zweite** Bewegungsrichtung im Mechanismus, die missbraucht werden könnte — gebunden an den Nenner-Nachweis, aber die Trennung „Transfer vs. Transfer-Regression" ist mechanisch nicht vollständig (benannte Grenze, §Konsequenzen); zwei Konstanten in zwei Trägern sind nachzuziehen |
| B — Schwellen **bleiben**; `slice-081` wird ein Carveout (`CO-<NNN>`) <!-- d-check:status-provenance --> | kein Eingriff am Mechanismus; die Naht landet schnell | lässt das **systematische** Problem stehen: die Dränage wiederholt sich bei jeder weiteren Naht, und die Rampe bleibt auf 75 % über einem Ist von 73,38 % — **jeder** spätere Lauf desselben Verfahrens wäre rot; der Carveout befreit _diesen_ Slice, nicht die Klasse |
| C — kein Schwellen-Schritt; **Neudefinition des DB-Gegenstands**, sodass ein Transfer die Zahl nicht bewegt | die Zahl wäre stabil | nicht ehrlich zu haben (siehe §Entscheidung „Warum die Schwellen-Antwort trägt"): Subjekt-Lüge, Überlappung oder gespeicherter Basiswert — jede Variante verletzt `ADR-0071` Punkt 4 oder Modul 6; und der Effekt (84 % gegen 75 %) kehrt unter jeder Gegenstands-Ziehung wieder |
| D — **anderer Träger**: Gate-Typ wird eine No-Regression-Ratsche auf einem gespeicherten Basiswert | fängt die reine Regression scharf | der gespeicherte Basiswert ist eine **zweite Quelle** für denselben Zustand (Modul 6) und müsste bei **jedem** Transfer mitgezogen werden — die Transfer-Frage bleibt, nur um eine Pflege-Last reicher; ersetzt den vorhandenen Mechanismus statt ihn zu vervollständigen |
| E — **eine** Gesamtzahl über beide Gegenstände (Träger wie Option C aus `ADR-0071`) | eine einzige Zahl | öffnet genau den verworfenen Pfad: PostgreSQL im `make gates`-Lauf (Option C aus `ADR-0071`) oder eine Zahl, deren Verfahren den Gegenstand nicht ausführt; `ADR-0071` hat ihn bereits entschieden |
| F — nichts tun | kein Aufwand | die Dränage bleibt, die Rampe bleibt rot über einem erreichbaren Ist, und Trigger (b) bleibt uneingelöst — `ADR-0071` verlangt für das Ziehen der Naht ausdrücklich eine Folge-ADR |
| G — `ADR-0071` **Trigger (b) per Teil-Supersede umschreiben** (die abfließende Seite ergänzen) | die Trigger-Formulierung wäre dann vollständig | Trigger (b) ist **nicht falsch** — er hat den Fall vorhergesehen und korrekt gefeuert; er ist nur **still** über die zweite Seite. `AGENTS.md` §3.5 verlangt ein `Supersedes` für eine **Korrektur** eines unwahren Satzes, nicht für die Ausführung eines wahren Triggers; ein Supersede wäre die Behauptung, dort habe etwas Falsches gestanden |

**Fazit:** A. Der Schnitt aus `ADR-0071` Punkt 1 hat eine Dränage erzeugt,
für die der Mechanismus aus `ADR-0054` §(a) keine Richtung hat; A schließt die
Lücke **im vorhandenen Verfahren**, statt es durch einen neuen Gegenstand (C,
D) oder eine neue Gesamtzahl (E) zu ersetzen, und statt die Klasse mit einem
Einzel-Carveout (B) stehen zu lassen.

## Konsequenzen

- Positiv: Der rote `make test-replication`-Lauf wird aufgelöst, **ohne** ein
  Verhalten zu ändern und **ohne** eine Endstufe zu senken — §3.6 ist über
  diese ADR als Träger erfüllt.
- Positiv: Der Mechanismus ist ab hier in **beide** Richtungen geschlossen;
  die Dränage ist nicht länger eine Lücke, sondern ein Nachweis.
- Positiv: die zwei Folge-Nähte (`postgresack`, `receive`) erben eine
  **wiederholbare** Regel statt eines neuen Einzel-Falls.
- Negativ mit Grenze: **Die Nenner-Paarung trennt einen Transfer nicht
  vollständig von einer gleichzeitigen Regression.** Träte beides zusammen
  auf, fiele der neu bemessene Einstieg tiefer als ein sauberer Transfer ihn
  setzte. Die Wächter dieses Rest-Falls sind die **Beleg-Disziplin** (die
  Naht ist ohne Verhaltens-Change zu haben — `slice-081` §2, belegt über <!-- d-check:status-provenance -->
  `make test-store`/`make test-replication` grün) und das **Review**, nicht
  dieser Sensor. Benannt, nicht verschwiegen.
- Negativ: Mit jeder weiteren Naht schrumpft der DB-Gegenstand. Eine 80-%-Rampe
  über einem Rest-Gegenstand misst einen Rest; das ist kein Defekt, aber es ist
  ein anderer Gegenstand als heute — deshalb der Re-Evaluierungs-Trigger (b)
  unten.
- Folgepflicht (Implementer-Zug, Tooling/Doku, **kein** Produkt-Code):
  `DB_COVERAGE_THRESHOLD` in `tools/harness/db-coverage.sh` von 75 auf 70;
  `THRESHOLD` in `harness/mk/coverage.mk` von 65 auf 70;
  `harness/sensors/db-adapter-coverage.md` und
  `harness/sensors/coverage-gate.md` (§Kalibrierungs-Bindung),
  `harness/README.md` §Sensors und `AGENTS.md` §4 ziehen die genannten
  Einstiegswerte nach. **Grün-Beleg:** `make test-replication` auf dem Branch
  endet Exit 0 (DB-Adapter-Coverage 73,38 % ≥ 70 %), `make coverage-gate`
  bleibt grün (71,30 % ≥ 70 %).
- Implikation für `slice-081` (Planner-Zug, **kein** Beschluss hier): die <!-- d-check:status-provenance -->
  DoD-Zeile „die reale Verdrahtung geht unverändert durch `make test-store`
  und `make test-replication` (Exit 0)" wird mit dem Landen der Neu-Bemessung
  grün; der Blocker war eine **unentschiedene** Architect-Frage, und die ist
  mit dieser ADR entschieden. Ob die Neu-Bemessung im Slice selbst oder in
  einem unmittelbaren Folge-Slice landet, ist der Slice-Schnitt der Planner-
  Rolle.
- Folgepflicht (Planner-Zug): die zwei Folge-Nähte (`postgresack`, `receive`)
  tragen je ihren Nenner-Nachweis und wenden diese Regel an, statt sie neu zu
  verhandeln.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| `tools/harness/db-coverage.sh` + `Dockerfile` (Stufe `coverage`) | Eine **Neu-Bemessung** eines Rampen-Einstiegs setzt die **Nenner-Paarung** voraus: die Statement-Summe beider Gegenstände konstant, die Verteilung verschoben. Eine Quote, die bei **unverändertem** Nenner fällt, ist eine Regression und wird **nicht** neu bemessen. Die tragenden Zahlen sind die gedruckten Zeilen beider Messungen | `make test-store` / `make test-replication` · `make coverage-gate` |
| (Disziplin, kein Sensor) | Der Rest-Fall „Transfer **und** Regression im selben Zug" ist von der Nenner-Paarung allein **nicht** fassbar; Wächter sind die Beleg-Disziplin der Naht und das Review | — |

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

Vier benannte Trigger, sonst permanent:

**(a)** Ein **weiterer Transfer** tritt ein (`postgresack`, `receive` oder ein
späterer Umbau) — dann wendet er **diese** Regel an (Neu-Bemessung mit
Nenner-Nachweis), ohne dass eine neue Schwellen-ADR fällig wird. **Kein**
Transfer: der Einstieg bleibt stehen, eine gesunkene Quote bei unverändertem
Nenner ist eine Regression und rotet.

**(b)** Der DB-Gegenstand fällt unter die **Hälfte** seiner
Kalibrierungs-Baseline (788 Statements, `ADR-0071` Punkt 3) — dann misst die
Endstufe 80 % einen Rest-Gegenstand, dessen Aussagekraft eine andere ist, und
die Endstufe ist neu zu beurteilen (Folge-ADR). Beobachtbar an der gedruckten
Statement-Zahl der DB-Messung.

**(c)** Ein Paket **kommt hinzu oder fällt weg** — nicht durch Transfer,
sondern durch Aufnahme/Austritt —, sodass sich die Gegenstands-Liste ändert:
`ADR-0071` Trigger (a) greift für die Liste; diese ADR greift für die
Bemessung, **sofern** die Nenner-Paarung einen Transfer zeigt, sonst steht die
Schwelle.

**(d)** Die **Partition** selbst ändert ihre Form (ein dritter Gegenstand
entsteht, etwa eine weitere dienst-gestützte Testebene) — dann ist die Zwei-
Gegenstände-Ordnung aus `ADR-0071` Punkt 4 neu zu beurteilen.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-15 | Accepted — Anlass: `slice-081` zieht die Naht (`ADR-0071` Punkt 5); der entleerte DB-Gegenstand fällt 75,25 % → 73,38 % auf einem verhaltens-erhaltenden Umbau (138 Statements transferiert, davon 116 gedeckt), `make test-replication` Exit 2, einzige `FAIL`-Zeile das Coverage-Verdikt. `ADR-0071` Trigger (b) tritt ein. Entscheidet die Nenner-gebundene Neu-Bemessungs-Regel und die zwei Einstiege (DB 75 → 70, Unit 65 → 70); bestätigt `ADR-0071` ohne Supersede | `docs/plan/planning/in-progress/slice-081-executor-naht.md` <!-- d-check:status-provenance --> |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0077` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
