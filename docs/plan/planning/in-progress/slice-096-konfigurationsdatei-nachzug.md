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

- [ ] **LP1 — die fünf Variablen werden durchgereicht.** Unter gesetzter
      `CDC_CONFIG_FILE` trägt `Config` die Werte aus `CDC_NATS_URL`,
      `CDC_HTTP_ADDR`, `CDC_API_TOKEN_READER`, `CDC_API_TOKEN_ADMIN` und
      `CDC_GRPC_ADDR` — je aus ihrer Env-Herkunft, und **Env schlägt Datei
      Feld für Feld**. Je Variable ein Test **und** ein rot gesehenes
      Gegenbeispiel (Mutation am Durchreichen), plus der Nachweis, dass die
      Oberflächen unter geladener Datei **wirklich an** sind.
- [ ] **LP2 — die Feldmenge und die Zugangsdaten-Klasse.** `http_addr` und
      `grpc_addr` sind Datei-Felder (`host:port`, nicht credential-tragend);
      `nats_url` und die zwei Tokens sind **env-exklusiv** und enden in
      `ErrConfiguration` mit einer **eigenen, den Grund nennenden Zeile** —
      nicht als „unbekannter Schlüssel". Je Klasse ein Test.
- [ ] **LP3 — der Träger zieht nach.** `docs/user/benutzerhandbuch.md` §5.2
      führt die zwei neuen Felder und die env-exklusiven Namen mit ihrer
      Begründung, **und die Versionshistorie bekommt ihre Zeile** (verkörperte
      Regel aus `BEO-PGC/handbuch-versionshistorie-uebersprungen`).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update für <Schnittstelle X> falls öffentlicher Vertrag berührt.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

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
  Slice sind sie **an**. — **Ausgang:** <…>
- **Die Zugangsdaten-Klasse greift zu weit** — ein legitim credential-freies Feld
  wird abgewiesen, weil seine Form (URL) auch Zugangsdaten tragen *könnte. —
  **Ausgang:** <…>
- **Ein Träger wird überholt, den dieser Slice nicht anfasst** — er bewegt, welche
  Variablen unter einer Datei wirksam sind; welche Dokumente das beschreiben, weiß
  der Diff nicht (`BEO-PGC/arbeit-ueberholt-stehenden-traeger`, **3×**, seit
  `welle-20` eine Hard Rule). — **Ausgang:** <…>
- **Die Lücke ist größer als die fünf** — die Messung des Vorgänger-Zugs nennt
  genau fünf; ob weitere Felder `mergeConfig` nicht erreichen, ist **nicht**
  gemessen. — **Ausgang:** <…>

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

- **Was hat funktioniert:** <…>
- **Was ging anders als geplant:** <…>
- **Steering-Loop-Eintrag:** <Guide oder Sensor> <geschärft/ergänzt>: <was genau>
  — liegt in `<AGENTS.md §X | Makefile:<target> | .harness/skills/…>`.
  Auslöser: `BEO-<NNN>` (<slice-NNN>, <slice-MMM>, <slice-KKK> — 3×).
  *(Wurde mit diesem Slice nichts verkörpert — der Normalfall —, entfällt die
  Teil-Zeile `— liegt in …` ersatzlos. Der Eintrag ist dann gezählt, nicht
  verkörpert.)*
- **Beobachtungs-Register (`../observations/`):** <`BEO-<KUERZEL>/<slug>/` neu angelegt, Beleg `evidence/slice-NNN.md` | `evidence/slice-NNN.md` in `BEO-<KUERZEL>/<slug>/` ergaenzt — Zaehler steht damit bei <N>x | keine Beobachtung angefallen>
- **Folge-Slices:** <slice-NNN (<Titel>) — ist eine Datei in `open/`>
- **Risiken aus §6:** <jedes mit genau einem Ausgang — siehe §6>
- **Drei Paarungen:** <nur im Repo ohne Wellen-Betrieb — Anker · Folge-Slice · Register, Ergebnis>

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
