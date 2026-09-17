# Review-Report: slice-101 — 2026-09-17

**Review-Art:** Code — Code-Review gegen Plan + Konventionen (Modul 10
§Drei Review-Arten): geprüft wird der Diff gegen den Slice-Plan, `ADR-0090`
und Modul 5/8 des Baseline-Regelwerks, nicht gegen die DoD (das ist
Verifier-Aufgabe, Modul 11).

**Gegenstand:** `git diff 03337c9..HEAD` — Commits `89f3ba5`, `7656c36`,
`79dbd5d` (Slice `slice-101`, Plan
`docs/plan/planning/in-progress/slice-101-nats-client-csharp-kotlin.md`).

**Skill:** `.harness/skills/reviewer.md` @ Accepted, Schärfung 2026-09-09 ·
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-17.

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-101-nats-client-csharp-kotlin.md`
  (§1 Abgrenzung, §2 DoD, §4 Trigger/Rückführungen, §6 Risiken)
- `ADR-0090` (Festlegung 1 „eine gepinnte öffentliche Client-Bibliothek je
  Sprache", Festlegung 4 „Was die Quellen benutzen dürfen")
- `ADR-0055`/`ADR-0056` (Wecksignal, Subjekt-Schema), `ADR-0079`
  (zweiseitiger Ablauf, Form-Vorbild), `ADR-0076` §1 (Startform/Env-Var-
  Konvention)
- Baseline-Regelwerk `modul-05-planning-harness.md` §Trigger je
  Lifecycle-Übergang, §Lifecycle als State Machine; `modul-08-agentenrollen.md`
  §Rollen-Regeln, §Konflikt-Pfad als Rollen-Sequenz
- `AGENTS.md` §3.12 (Herkunft von Aussagen), §3.13 (überholte Träger)
- Form-Vorbild `examples/nats-client` (Go)
- Externe Verifikation (netzlos nicht möglich, mit Netzzugriff durchgeführt):
  Maven Central (`jnats-2.26.3.pom`/`.module`, `bcprov-lts8on-2.73.12.1.pom`),
  NuGet (`nats.net`-Nuspec/-Index), `bouncycastle.org/licence.html`,
  OSV.dev (`org.bouncycastle:bcprov-lts8on`)

---

## Findings

### F-1 — Vorab benannter Rückführungs-Trigger (§4) traf real ein und wurde vom Implementer selbst weg-bewertet statt formal ausgelöst

- `kategorie`: HIGH
- `quelle`: Baseline-Regelwerk `modul-05-planning-harness.md` §Trigger je
  Lifecycle-Übergang / §Drei Übergänge tragen eine Pflicht; `modul-08-
  agentenrollen.md` §Rollen-Regeln, §Konflikt-Pfad als Rollen-Sequenz
- `pfad`: `docs/plan/planning/in-progress/slice-101-nats-client-csharp-
  kotlin.md:161-165` (§4, Rückführungs-Trigger 1) gegen
  `examples/kotlin/nats-client/build.gradle.kts:1-24` (Kommentar) und Commit
  `7656c36`
- `befund`: §4 des Slice-Plans benennt vorab, wörtlich: „`in-progress` →
  `next`… wenn eine der beiden gepinnten Client-Bibliotheken transitiv eine
  zweite, gepinnte Abhängigkeit erzwingt, die eine eigene Bewertung braucht
  (Lizenz, Sicherheits-Historie) — **dann ist die Größenannahme… falsch**".
  Der Implementer hat real gemessen (von mir unabhängig nachvollzogen: die
  `.pom`/`.module` von `jnats:2.26.3` listen `org.bouncycastle:bcprov-
  lts8on:2.73.12.1` mit `<scope>compile</scope>`, kein `optional`-Flag —
  unconditional transitive Abhängigkeit), dass genau diese Bedingung
  eingetreten ist. Statt den beschriebenen Übergang `in-progress → next`
  auszulösen (Rückgabe an Planner/Architect zur Neu-Zerlegung bzw.
  Trigger-Neubewertung, Modul 8 §Konflikt-Pfad), hat der Implementer die
  geforderte „eigene Bewertung" selbst im Pin-Kommentar des Build-Manifests
  vorgenommen und die Rückführung für nicht nötig erklärt. Der Wortlaut des
  Triggers ist eine deterministische Folge („X tritt ein → Größenannahme ist
  falsch"), keine an eine Ermessens-Klausel gebundene Bedingung — Modul 5
  trennt ausdrücklich *vorab benannte Bedingung* von *im Nachhinein
  nachgetragenem Grund* („was tatsächlich eintrat"), nicht von einer dritten
  Option „Bedingung eingetreten, aber wegbewertet". Modul 8 reserviert
  strukturanaloge Entscheidungen (Abweichung von einer vorab getroffenen
  Festlegung) für Architect/Planner — der Implementer „darf höchstens
  Folge-ADR vorschlagen, niemals stillschweigend… widersprechen"; hier ist es
  kein ADR, sondern ein Plan-Trigger, aber dieselbe Rollen-Grenze gilt: Der
  Implementer ist nicht die Rolle, die eine vorab als Rückführungs-Grund
  deklarierte Bedingung nachträglich für irrelevant erklären darf, so
  plausibel die Bewertung im Ergebnis auch ist.

  **Zur Bewertung selbst — inhaltlich nachgeprüft:** Die Lizenz-Hälfte
  trägt: „Bouncy Castle Licence" ist wörtlich eine MIT-Lizenz-Variante
  („Permission is hereby granted, free of charge…", real von
  `bouncycastle.org/licence.html` gelesen) — die Kommentar-Aussage
  „permissiv lizenziert" ist zutreffend und keine Übertreibung. Die
  Sicherheits-Historie-Hälfte des Triggers ist dagegen **nicht** mit
  vergleichbarer Sorgfalt dokumentiert: Kommentar und Commit-Message
  nennen ausschließlich Reputation („seit Jahrzehnten breit eingesetzt"),
  keinen CVE-/Advisory-Check. Eigene Nachprüfung (OSV.dev,
  `org.bouncycastle:bcprov-lts8on`) zeigt zwei bekannte Advisories
  (`GHSA-4h8f-2wvx-gg5w`/CVE-2024-34447, behoben ≤ 2.73.6;
  `GHSA-mx76-r943-rf8g`/CVE-2026-8149, behoben ≤ 2.73.11) — die gepinnte
  Version `2.73.12.1` liegt hinter beiden Fix-Ständen, eine gezielte Abfrage
  auf exakt diese Version liefert `{}` (keine offene Vulnerability). Die
  Bewertung trägt im Ergebnis, aber der Trigger-Text verlangt ausdrücklich
  **beide** Hälften („Lizenz, Sicherheits-Historie") als Teil der „eigenen
  Bewertung", und nur eine davon ist im Repo belegt — die andere wurde im
  Kommentar durch eine unbelegte Pauschalaussage ersetzt (Berührung mit
  `AGENTS.md` §3.12 „Beleg trägt seinen Satz").
- `verifizierbar`: teilweise — die transitive Abhängigkeit und ihre Lizenz
  sind mit `curl`/`grep` gegen Maven Central und die Bouncy-Castle-Website
  ohne Rückfrage nachvollziehbar (siehe oben); ob der Rollen-Übergang
  formal hätte laufen müssen, ist ein Modul-5/8-Urteil, kein Gate-Lauf
- `klasse`: Rückführungs-Trigger vom Implementer selbst weg-bewertet statt
  formal ausgelöst

### F-2 — Sicherheits-Historie-Hälfte der geforderten Bewertung nicht mit Beleg geführt

- `kategorie`: MEDIUM
- `quelle`: `AGENTS.md` §3.12 (Herkunft von Aussagen in Trägern)
- `pfad`: `examples/kotlin/nats-client/build.gradle.kts:1-24`, Commit
  `7656c36` (Commit-Message)
- `befund`: Die Trigger-Bewertung (siehe F-1) nennt für die Lizenz-Hälfte
  eine nachprüfbare Tatsachenbehauptung, für die Sicherheits-Historie-Hälfte
  dagegen nur eine unbelegte Reputationsaussage („breit eingesetzt… seit
  Jahrzehnten") ohne CVE-/Advisory-Abfrage oder deren Ergebnis. Getrennt von
  F-1 aufgeführt, weil dieser Befund auch dann bestehen bliebe, wenn der
  Rollen-Übergang aus F-1 nicht gefordert wäre — die inhaltliche Bewertung
  selbst ist unvollständig belegt, unabhängig davon, wer sie treffen durfte.
- `verifizierbar`: ja — Ergänzung eines nachgeprüften Advisory-/CVE-Belegs
  (z. B. OSV-Abfrage-Ergebnis) im selben Kommentar wäre der Fix
- `klasse`: Aussage ohne Ursprungsbeleg (Sicherheits-Historie-Hälfte)

## Negativbefunde

- geprüft, ohne Befund: `examples/csharp/nats-client/**` (zweiseitiger
  Ablauf, Subjekt-/URL-Aufbau, CLI-Parsing, Tests) — Form-Vorbild-treu
- geprüft, ohne Befund: `examples/kotlin/nats-client/**` (zweiseitiger
  Ablauf, Subjekt-/URL-Aufbau, CLI-Parsing, Tests) — Form-Vorbild-treu
- geprüft, ohne Befund: `NATS.Net:3.2.0`/`io.nats:jnats:2.26.3` sind real
  die neuesten stabilen (nicht-Prerelease) Versionen (eigene Abfrage gegen
  NuGet-Index und Maven-`maven-metadata.xml`)
- geprüft, ohne Befund: `examples/csharp/Dockerfile`,
  `examples/kotlin/Dockerfile`, `examples/kotlin/settings.gradle.kts`,
  `harness/mk/examples.mk` — drittes Programm/Modul samt `runtime-nats`-Stufe
  konsistent zur bestehenden `runtime-sse`-Form
- geprüft, ohne Befund: `harness/README.md` — „zwei"→„drei" Images-Korrektur
  präzise für beide Sprachziele
- geprüft, ohne Befund: `docs/user/benutzerhandbuch.md` §4 „Zugriff über das
  NATS-Wecksignal" — drei-Zeilen-Block (Go/C#/Kotlin) und
  Versionshistorie-Zeile 1.23 fortlaufend und inhaltlich zutreffend
- geprüft, ohne Befund: Out-of-Scope-Disziplin — keine Zustandsmaschine
  (Reconnect/Dedup/Rückstand) in beiden Clients, `.a-check.yml` unverändert,
  `examples/nats-client` (Go) unverändert
- geprüft, ohne Befund: Commit-Traceability aller drei Commits (`LH-FA-
  SST-007`, `ADR-0090`/`-0055`/`-0056`/`-0079` je Betreff, kein `SPEC-*`/
  `ARC-*` im Betreff)
- geprüft, ohne Befund: eigener §3.13-Suchlauf — `grep` über `*.md` nach
  `examples-csharp`/`examples-kotlin` findet keine weitere Fundstelle
  außerhalb der bereits im Diff gezogenen (`harness/README.md`); ADRs sind
  immutabel und `done/`-Slices sind historisch, beide zu Recht unangetastet

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 1 |
| LOW | 0 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** Rückführungs-Trigger vom Implementer selbst
weg-bewertet statt formal ausgelöst · Aussage ohne Ursprungsbeleg
(Sicherheits-Historie-Hälfte)

## Verdikt

**Merge-blockierend:** ja — F-1 ist HIGH mit Rollen-Widerspruch (Implementer
hat eine vorab als Rückführungs-Grund deklarierte, real eingetretene
Bedingung selbst für nicht bindend erklärt). Nach Modul 8 §Konflikt-Pfad
als Rollen-Sequenz läuft das **nicht** als stille Fixrunde am Implementer,
sondern als Sequenz mit Übergabe-Artefakt: Reviewer → Architect (ggf. über
Planner) mit der Frage, ob (a) der Trigger-Wortlaut in §4 tatsächlich einen
zwingenden `in-progress→next`-Übergang meint (dann: Rückführung nachholen,
Bewertung als Grund im Nachhinein dokumentieren) oder (b) die Implementer-
Bewertung als legitime Konkretisierung akzeptiert und der Trigger-Wortlaut
in einer Folge-Klärung geschärft wird, damit dieselbe Ambiguität nicht ein
drittes Mal auftritt. Beide Verdikte sind mit den hier real geprüften
Fakten vereinbar (Lizenz ist tatsächlich permissiv, keine offene CVE für die
gepinnte Version) — das Verdikt selbst ist damit keine Formsache, sondern
eine Rollen-Frage, die dieser Report nicht anstelle des Architects trifft.

F-2 bleibt unabhängig vom Ausgang von F-1 bestehen und wäre mit einer
Fixrunde am Implementer lösbar (Beleg nachtragen), sobald F-1 aufgelöst ist.

**Übergabe:** F-1 geht über den Konflikt-Pfad (Modul 8) an Architect/Planner,
nicht direkt an den Implementer zurück. F-2 geht regulär an den Implementer.
Beide Finding-Klassen gehen zusätzlich in die Slice-Closure §7 und von dort
in den Zähler (Beobachtungs-Register). Dieser Report selbst ist ein
Lauf-Beleg und wird über Läufe hinweg nicht wieder gelesen. Der Report
ersetzt keine Verifikation — DoD-/Spec-Konformität (inkl. `make gates`
grün) prüft der Verifier separat (Modul 11).
