# Slice slice-096: Konfigurationsdatei-Nachzug — Durchleitung und Feldmenge

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — Reaktives auf einen real gemessenen Defekt (Modul 6
§Wann Arbeit eine Welle braucht); `welle-20` ist geschlossen und davon
unberührt.

**Bezug:** [`ADR-0088`](../../adr/0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md)
(die Entscheidung: Durchleitung, Feldmenge, Zugangsdaten-Klasse) ·
[`ADR-0052`](../../adr/0052-optionale-yaml-konfigurationsdatei.md) (die Datei-Semantik;
`ADR-0088` löst genau seine Feldaufzählung ab) ·
[`ADR-0057`](../../adr/0057-http-grpc-api.md), [`ADR-0060`](../../adr/0060-grpc-streaming-mechanismus.md),
[`ADR-0061`](../../adr/0061-http-sse-zusaetzlich-zu-grpc.md) und
[`ADR-0055`](../../adr/0055-nats-change-notification-wecksignal.md) (die fünf
Variablen und die Oberflächen, die sie tragen) ·
`BEO-PGC/arbeit-ueberholt-stehenden-traeger` (**3×**) und
`BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht` (**5×**).

**Berührte Spec-Stellen:** `SPEC-016` (Feldform der Konfigurationsdatei) — die
Spec nennt diesen Slice nie (Referenz-Richtung SDP).

**Verantwortlich:** pt9912.

**Autor:** pt9912 (Planner). **Datum:** 2026-09-17.

---

## 1. Ziel und Abgrenzung

<!-- BEDIENHINWEIS: Ziel = ein Satz, Liefer-Fokus, kein "wir machen
aufraeumen". Abgrenzung = je Punkt eine Begruendung, nicht nur eine Nennung:
ein Ausschluss ohne Grund ist eine Behauptung. Keine Mindestzahl — ein echter
Ausschluss ist besser als vier erfundene. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Die Konfigurationsdatei (`CDC_CONFIG_FILE`) trägt wieder, was `SPEC-016`
zusagt. Drei Liefer-Punkte, und der erste ist ein **Defekt**, kein Nachzug:

**Der Defekt.** `ConfigFromEnv` liest fünf Variablen — `CDC_NATS_URL`,
`CDC_HTTP_ADDR`, `CDC_API_TOKEN_READER`, `CDC_API_TOKEN_ADMIN`, `CDC_GRPC_ADDR`;
`mergeConfig` (der Pfad bei gesetzter Datei, den `cmd/pg-change-feed/main.go`
für den Betriebslauf nimmt) setzt **keines** der fünf Felder. Folge: die
HTTP-API, der gRPC-Stream und das NATS-Wecksignal bleiben **still deaktiviert**,
beide Token-Klassen **leer** — und der Prozess startet **gesund**. Die
Haupt-Klausel von `SPEC-016` („additiv … Env schlägt Datei Feld für Feld") ist
für fünf Variablen heute **falsch**. Kein Test nennt die fünf: die Lücke ist
**unentdeckt, nicht entschieden**.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Die C#- und Kotlin-Beispiele** ([`ADR-0087`](../../adr/0087-beispiel-clients-csharp-kotlin.md)):
  eigene Slices, und sie setzen diesen voraus — sie brauchen die Adressen, die
  hier durchgereicht werden.
- **Der nicht-blockierende Workflow** über die zwei Sprachziele: ein anderer
  Vorgang (CI-Träger, `§3.10`-Risiko), kein Bestandteil der Datei-Semantik.
- **Der Go-gRPC-Client** und alles, was `gen/**` berührt: der fehlende
  Umzugs-Slice aus [`ADR-0076`](../../adr/0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
  §3 ist weiterhin ungeschnitten.
- **Die Adressen der anderen Seite** — die Beispiel-Clients lesen ihre Adresse
  und ihr Token aus **Env oder Flag** und **kein** `CDC_CONFIG_FILE`
  ([`ADR-0087`](../../adr/0087-beispiel-clients-csharp-kotlin.md)); sie sind nicht
  Gegenstand dieses Slice.
- **Ein Gate für die Datei-Semantik.** Das strikte Decoding und die
  Zugangsdaten-Prüfung sind Verhalten, das die Unit-Tests tragen; ein eigenes
  Gate wäre Zeremonie.

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

<!-- BEDIENHINWEIS: je Zeile ein pruefbares Kriterium. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] **LP1 — die fünf Variablen werden durchgereicht.** Unter gesetzter
      `CDC_CONFIG_FILE` trägt `Config` die Werte aus `CDC_NATS_URL`,
      `CDC_HTTP_ADDR`, `CDC_API_TOKEN_READER`, `CDC_API_TOKEN_ADMIN` und
      `CDC_GRPC_ADDR` — je aus ihrer Env-Herkunft, und **Env schlägt Datei
      Feld für Feld**. Je Variable ein Test **und** ein rot gesehenes
      Gegenbeispiel (Mutation am Durchreichen), plus der Nachweis, dass die
      Oberflächen unter geladener Datei **wirklich an** sind — **an den
      Prädikaten, die über den Start entscheiden**, nicht an einem laufenden
      Server: `Run` konstruiert vorher den Store und braucht eine erreichbare
      PostgreSQL-Instanz. Die Grenze steht im Testkommentar.
- [x] **LP2 — die Feldmenge und die Zugangsdaten-Klasse.** `http_addr` und
      `grpc_addr` sind Datei-Felder (`host:port`, nicht credential-tragend);
      `nats_url` und die zwei Tokens sind **env-exklusiv** und enden in
      `ErrConfiguration` mit einer **eigenen, den Grund nennenden Zeile** —
      nicht als „unbekannter Schlüssel". Je Klasse ein Test.
- [x] **LP3 — der Träger zieht nach.** `docs/user/benutzerhandbuch.md` §5.2
      führt die zwei neuen Felder und die env-exklusiven Namen mit ihrer
      Begründung, **und die Versionshistorie bekommt ihre Zeile** (verkörperte
      Regel aus `BEO-PGC/handbuch-versionshistorie-uebersprungen`).
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

<!-- BEDIENHINWEIS: Datei- oder Komponenten-Ebene reicht; der
Implementer-Agent erweitert die Liste in seinem ersten Lauf. -->

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/bootstrap/config_file.go` (`fileConfig`, `mergeConfig`) | update | Zwei Felder, die Durchleitung der fünf, die Zugangsdaten-Klasse samt Fehlerzeile. |
| `internal/bootstrap/config_file_internal_test.go` | Test neu/update | Je Feld und je Vorrangsrichtung; je Zugangsdaten-Klasse ein Ablehnungsfall. |
| `docs/user/benutzerhandbuch.md` §5.2 | update | Feldmenge, env-exklusive Namen mit Begründung, Versionshistorie. |
| `harness/sensors/coverage-gate.md` | **nicht** | Kein Messgegenstand berührt. |
| `spec/pflichtenheft.md` | **nicht** | `SPEC-016` ist mit [`ADR-0088`](../../adr/0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md) bereits nachgezogen — der Code folgt der Spec (GF-Modus). |

## 4. Trigger

<!-- BEDIENHINWEIS: Beispiele — "Wenn Welle X done." / "Wenn Carveout CO-NN
aufgeloest." -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): [`ADR-0088`](../../adr/0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md)
ist `Accepted` und `SPEC-016` trägt die neue Feldmenge. Ohne Rückfrage
feststellbar.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): wenn die Durchleitung
  **mehr als eine Stelle** im Kontrollfluss braucht (z. B. wenn die
  Datei-Herkunft sich nicht feldweise mit der Env-Herkunft mischen lässt, ohne
  `ConfigFromEnv` umzubauen). Dann ist der Schnitt falsch: dieser Slice **behebt**
  die Durchleitung, er baut die Konfiguration nicht um.
- `in-progress` → `open` (blockiert — Carveout?): wenn die Verhaltensänderung
  **mehr als die Datei-Semantik** trifft — etwa wenn ein Deployment sie
  voraussetzt, dass die Oberflächen unter Datei **aus** bleiben. Dann gehört das
  in eine Entscheidung (Rücknahme oder Übergangsform), nicht in diesen Slice.

## 5. Closure-Trigger

<!-- BEDIENHINWEIS: z.B. "DoD vollstaendig + PR gemerged + Closure-Notiz
geschrieben." -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

**Zwei beobachtbare Kriterien:** LP1–LP3 sind real belegt — die fünf Variablen
sind unter geladener Datei **wirksam** nachgewiesen (nicht nur „im Struct
gesetzt"), die zwei Felder und die zwei Zugangsdaten-Klassen sind je durch Test
und Mutation belegt, der Träger ist nachgezogen — **und** `make gates` ist grün.

**Lerneintrag:** geschärfte Regel, neuer Sensor oder benannte Spec-Lücke in §7.
Der naheliegende Kandidat: die Lücke war **unentdeckt, nicht entschieden** — ob
daraus eine Regel wird („ein Feld, das die Spec zusagt, hat einen Test"), oder ob
die bestehende §3.12-Instanz B trägt, entscheidet der Lauf.

## 6. Risiken und offene Punkte

<!-- BEDIENHINWEIS: Was koennte schief gehen? Welche Carveouts entstehen
ggf.? Die drei Ausgaenge stehen als Form in der Zeile darunter. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Die Durchleitung ist eine Verhaltensänderung.** Deployments mit gesetzter
  `CDC_CONFIG_FILE` betreiben die Oberflächen heute **still aus**; nach diesem
  Slice sind sie **an**. — **Ausgang: eingetreten — und begrenzt, gemessen.**
  Die Änderung ist real und beabsichtigt. **Kein Lauf dieses Repos ist
  betroffen:** `git grep CDC_CONFIG_FILE` findet Treffer nur in `docs/**` und
  `internal/bootstrap/**` — nicht in `compose.yaml`, `tools/`, `test/` oder
  `.github/`. Betroffen sind **fremde** Deployments; die Wirkung ist in
  `docs/user/benutzerhandbuch.md` §5.2 beschrieben (Version 1.18).
- **Die Zugangsdaten-Klasse greift zu weit** — ein legitim credential-freies Feld
  wird abgewiesen, weil seine Form (URL) auch Zugangsdaten tragen *könnte. —
  **Ausgang: entfallen — gemessen.** Die Klasse hat **sechs** Schlüssel und
  greift genau die gemeinten: drei DSN, zwei Token, `nats_url`. `http_addr` und
  `grpc_addr` sind **Datei-Felder** und werden **nicht** abgewiesen; die
  Gegenprobe (`TestConfigFromFileUnbekannterSchluesselOhneZugangsdatenGrund`)
  zeigt, dass ein gewöhnlicher Tippfehler-Schlüssel **ohne** die
  Klassen-Begründung abgewiesen wird.
- **Ein Träger wird überholt, den dieser Slice nicht anfasst** — er bewegt, welche
  Variablen unter einer Datei wirksam sind; welche Dokumente das beschreiben, weiß
  der Diff nicht (`BEO-PGC/arbeit-ueberholt-stehenden-traeger`, **3×**, seit
  `welle-20` eine Hard Rule). — **Ausgang: eingetreten — und behoben, an vier
  Stellen.** Der §3.13-Lauf hat den **Symbolnamen** gefunden, den die Umbenennung
  ungültig machte (`ADR-0088` §Bezug, gemeldet und per Zitat-Korrektur gezogen);
  **drei weitere** Stellen fand erst der Review — zwei **Zeilen-Lokatoren** in
  `ADR-0088` §Kontext und ein vierter Zitat-Anker, den der Architect beim
  Nachziehen entdeckte. **Und die Klasse traf die Träger-Aussagen selbst:**
  `ADR-0088` nannte die Klasse „fünf" statt sechs, `ADR-0089`s Ersatztext war für
  die Feldmenge zu weit — beide als Folge-ADR (`0091`, `0092`). **Die Lehre steht
  in §7.**
- **Die Lücke ist größer als die fünf** — die Messung des Vorgänger-Zugs nennt
  genau fünf; ob weitere Felder `mergeConfig` nicht erreichen, ist **nicht**
  gemessen. — **Ausgang: entfallen — gemessen.** Der Implementer hat die Zahl
  **syntaktisch** nachgezählt (`Config`-Struct gegen die `cfg.X =`-Zuweisungen in
  `mergeConfig`), am Parent **und** am Diff-Stand: **15** Felder, am Parent **10**
  erreicht, jetzt **15 von 15**. Die Differenz ist **genau** die genannten fünf.

## 7. Closure-Notiz

<!-- BEDIENHINWEIS — keine Norm; faellt beim Kopieren weg (README.md
§Verwendung, Schritt 5) und darf deshalb nichts Tragendes halten. Reihenfolge:
diese Sektion vor dem `git mv` nach done/ fuellen — einzige Ausnahme ist das
letzte DoD-Item in §2 (die Paarungen suchen in `done/`, also nach dem `git mv`).
Im Repo ohne Wellen-Betrieb braucht die Closure dadurch drei Commits: Inhalt,
`git mv`, Haekchen — das folgt aus der Hard Rule, es widerspricht ihr nicht. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<NNN>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks).

- **Was hat funktioniert:** Die **Zählung vor dem Bauen** — der Implementer hat die
  Zahl des Vorgänger-Zugs nicht übernommen, sondern **syntaktisch nachgezählt**
  (`Config`-Struct gegen die `cfg.X =`-Zuweisungen) und kam auf **15/10/15**: die
  fünf, genau benannt, an beiden Ständen gemessen. Zweitens: **jede Zusage mit
  ihrer Grenze** — die zwei Kommentar-Hälften (was gebunden ist / was nicht) sind
  je durch eine Mutation belegt, und die Prüfung hat sie für **alle drei**
  Zweige bestätigt. Drittens: **die Kette** — der §3.13-Lauf fand den
  Symbolnamen, der Review zwei Lokatoren und einen vierten Anker, und der
  Architect-Zug fand die zwei **falschen Zahlen in den ADRs selbst**. Vier Runden,
  jeder Fund an einer anderen Stelle, **kein** Fund im Mechanismus.
- **Was ging anders als geplant:** **Der Nachzug war mehr als ein Nachzug.** Erwartet
  war eine Feldtabelle, die drei Variablen nicht führt; gefunden wurde ein
  **stiller Defekt** — `mergeConfig` ließ fünf Variablen fallen, und ein
  Deployment mit Konfigurationsdatei betrieb HTTP-API, gRPC-Stream und
  NATS-Wecksignal **ohne Fehler, ohne Log, ohne Test** aus. Zweitens: **die
  ADRs dieses Tages rattern** — `0087` wurde von `0090`, `0088` von `0091`,
  `0089` von `0092` abgelöst, **alle am selben Tag**, jede wegen einer Zahl oder
  einer Reichweite, die nicht stimmte. Und die **Korrektur ist selbst eine ADR**,
  die irren kann: Das Behebungs-Werkzeug trägt dasselbe Risiko wie sein
  Gegenstand. Drittens: **der §3.13-Lauf war unvollständig** (drei gebrochene
  Anker, einer gemeldet) — der Suchlauf findet **Symbolnamen**; Lokatoren sind
  Zahlen und fallen durch. Und sein Ergebnis steht in **einem Bericht**, der
  keinen Repo-Träger hat.
- **Steering-Loop-Eintrag:** **eine geschärfte Regel — an ihrem Gegenstand
  gemessen.** `AGENTS.md` §3.13 ist seit der `welle-20`-Closure verkörpert; dieser
  Slice ist ihr **erster Fall** und hat **zwei Grenzen** gezeigt: (1) der
  Suchlauf greift **Symbolnamen**, nicht **Zahlen** — Zeilen-Lokatoren und
  Nummern in Sätzen fallen durch, obwohl sie dieselbe Klasse sind; (2) das
  Ergebnis gehört in einen **Träger im Repo**, nicht in einen Handoff-Bericht —
  sonst ist der Lauf selbst nicht nachlesbar. Beides ist als Beobachtung
  eingetragen (`BEO-PGC/regel-weiter-als-ihr-sensor` → **3×**, Schwelle erreicht)
  und wird **nicht** in diesem Zug verkörpert: eine Hard Rule zu ändern ist ein
  Architect-Zug. Auslöser: `BEO-PGC/regel-weiter-als-ihr-sensor` (`slice-096`).
- **Beobachtungs-Register (`../observations/`):** **drei Belege** ergänzt —
  `zahl-in-traeger-driftet-gegen-die-messung` → **8×** (die zwei falschen Zahlen
  in `ADR-0088`/`ADR-0089`, **ein** Vorgang); `arbeit-ueberholt-stehenden-traeger`
  → **4×** (der §3.13-Lauf und die vier Ankerstellen); `regel-weiter-als-ihr-sensor`
  → **3×** (**Schwelle erreicht** — den Ausgang weist der Lese-Schritt der
  nächsten Wellen-Closure zu, Modul 6). **Kein Zähler wird gesetzt.**
- **Folge-Slices:** **keine Datei in `open/`.** Die sieben Slices der Client-Matrix
  sind von [`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md)
  **empfohlen** und noch nicht geschnitten — sie hier zu nennen wäre „genannt
  ohne angelegt", dieselbe Klasse wie ein halluziniertes Gate. Ihre Adresse ist
  die ADR; `slice-097` (der Umzug) und `slice-095` liegen in `next/`.
- **Risiken aus §6:** vier, je ein Ausgang — R1 *eingetreten und begrenzt*
  (kein Lauf dieses Repos setzt `CDC_CONFIG_FILE`), R2 *entfallen* (die Klasse
  greift genau die sechs), R3 *eingetreten und behoben* (vier Ankerstellen, zwei
  falsche ADR-Zahlen), R4 *entfallen* (15/10/15 nachgezählt).
- **Drei Paarungen:** Dieses Repo führt Wellen-Betrieb, aber es ist **keine Welle
  offen** (`welle-20` ist geschlossen, `open/` und der flache Planning-Pfad
  tragen keine Welle-Datei) — die Prüfung fällt der **nächsten** Wellen-Closure
  zu. Vorab geprüft: die drei berührten Register-Adressen existieren als
  Verzeichnis und tragen ein nicht leeres `evidence/`.

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

**Vorgelagert — Sub-Area-Wahl prüfen:** Berührt sind `internal/bootstrap`
(Composition Root) und die zwei Doku-Träger — die repo-weite Default-Sub-Area
`*`/`PGC` aus der Modus-Deklaration in
[`harness/conventions.md`](../../../../harness/conventions.md).

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen. Drei
Treffer: `arbeit-ueberholt-stehenden-traeger` (**3×**, verkörpert in
`AGENTS.md` §3.13 — dieser Slice ist sein **erster Fall nach der Verkörperung**:
die Env-Variablen der Wellen 15/16/19 haben eine Eigenschaft bewegt, die
Feldmenge der Datei ist stehen geblieben); `beleg-befehl-traegt-seinen-satz-nicht`
(**5×**, verkörpert — betrifft die Zusage von `SPEC-016`, die für fünf Variablen
heute falsch ist); `dod-begruendung-unzutreffende-tatsachenbehauptung` (**4×**,
verkörpert — betrifft LP1: keine Zahl im DoD-Kriterium, die nicht gemessen ist).

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas sind **GF**. Der
Block pro Sub-Area entfällt; der **Abschnitt** bleibt.

<!-- Block für jede berührte Sub-Area duplizieren. Format identisch
mit dem im Baseline-Regelwerk §Ziel-Form: Sub-Area-Modus-Begründung
abgedruckten Block. -->

### Sub-Area: <Name>

- **Modus:** GF | BF | Hybrid
- **Konventionen-Dichte:** <Beleg aus `harness/conventions.md`,
  Adaptions-Block oder Code>
- **Phase-Reife:** Phase 0–5 <Begründung gegen die Phase × Modus-Matrix>
- **Evidenz-/Diskrepanz-Risiko:** <bei BF/Hybrid: was kann die
  Inventur sichtbar machen? bei GF: meist niedrig>
- **Reconciliation-Aufwand:** <Slice-Schätzung;
  Graduation-/Folge-Slice-Trigger>
