# ADR-0062: Lokaler `commit-msg`-Hook als nicht-durchsetzende Vorab-Meldung — Supersedes ADR-0045 (nur die Hook-Klausel)

**Status:** Accepted — Supersedes [`ADR-0045`](0045-commit-traceability-standing-gate.md)
(nur die Klausel „Kein commit-msg-Hook" in deren Abschnitt „Entscheidung";
das Standing-Gate selbst — Fenster-Semantik `HEAD~5..HEAD`, Zwei-Werkzeug-
Aufteilung d-check-Modul `commits` + `tools/harness/commit-traceability.sh`,
`GATE_CHECKS`-Bindung — bleibt unverändert bestehen und wird hier nicht
wiederholt)

**Datum:** 2026-09-14

**Autor:** pt9912 (Architect-Rolle, Modul 8 — unabhängiger Architect-Zug,
gegenprüfend zum ersten Architect-Verdikt zur Commit-Traceability
(kein Vorab-Hook), der als Rollenspiel im Planner-Kontext lief und diesen
Konflikt nicht gegen den vollständigen Text von `ADR-0045` geprüft hatte)

**Bezug:** [`ADR-0045`](0045-commit-traceability-standing-gate.md)
(korrigierte Klausel), `AGENTS.md` §5 / `harness/README.md`
§Traceability rules (die getragene Regel selbst, unverändert),
`docs/plan/planning/open/slice-073-commit-msg-git-hook.md` (Umsetzungs-Slice, <!-- d-check:status-provenance -->
wellenlos), `docs/plan/planning/observations/BEO-PGC/commit-traceability-kein-vorab-hook`
(3×-Beobachtung, Auslöser), `.d-check.yml` (`exempt-pattern: '^(Merge |Revert
)'` — Merge-/Revert-Ausnahme der d-check-Positiv-Hälfte, die der Hook
spiegeln muss).

**Schärft:** — (Prozess-ADR ohne Spec-Stratum, wie `ADR-0045`)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

[`ADR-0045`](0045-commit-traceability-standing-gate.md) (Accepted,
2026-09-10) hat für die Commit-Traceability-Regel ein **Standing-Gate**
gewählt (`make commit-traceability`, im `make gates`-Bündel, Fenster
`HEAD~5..HEAD`) und dabei in ihrem eigenen Abschnitt „Entscheidung"
ausdrücklich festgehalten: „**Kein commit-msg-Hook** — der Hook fängt nur
die Pending-Message des Committers, aber die Regel gilt für Commits jedes
Urhebers (Planner-, Zweitschreiber-, Allowlist-Commits), und der
Gate-Nachweis (record-gates) sieht einen Hook nicht." Dieselbe Ablehnung
erscheint als Option B in ihrer Alternativen-Tabelle: „Hook-Installation
ist client-seitig nicht erzwingbar (kein Wächter); Commits außerhalb des
Harness bleiben ungedeckt; der Gate-Nachweis-Stempel sieht den Hook nicht."

Das Beobachtungs-Register (`BEO-PGC/commit-traceability-kein-vorab-hook`)
erreichte mit `slice-059` real 3× (Verstöße, die erst nach dem `git push` <!-- d-check:status-provenance -->
durch `make gates` auffielen, teils nicht mehr per `git commit --amend`
korrigierbar) und schlug als Lösung genau das vor, was `ADR-0045` bereits
abgelehnt hatte: einen lokalen `commit-msg`-Hook. Ein erster
Architect-Zug (Architect-Verdikt zur Commit-Traceability, kein Vorab-Hook,
2026-09-14) prüfte diesen Vorschlag, zitierte `ADR-0045` in seinem
Bezugs-Feld als „bindend", stellte aber nicht fest, dass `ADR-0045`s
eigener Entscheidungstext einen Hook explizit ausschließt — er kam zum
Schluss „keine ADR nötig, da keine Architektur-Entscheidung im Sinn von
Modul 4 getroffen wird (reine Tooling-Ergänzung, DevEx)" und plante
`slice-073` mit dem Bezugs-Feld „keine Änderung an `ADR-0045` nötig". Das <!-- d-check:status-provenance -->
ist die Lücke: Dieser Zug lief im selben Kontextfenster wie die
vorausgehende Planner-Arbeit an der `welle-16`-Closure (transparent
benanntes Rollenspiel, kein unabhängiger Architect-Lauf) und hat die
bereits getroffene Entscheidung nicht gegengelesen — exakt der blinde
Fleck, den Modul 8s Kontext-Trennung verhindern soll.

Ein unabhängiger Architect-Zug (dieser hier) hat `ADR-0045` vollständig
gelesen und den Widerspruch festgestellt: `slice-073` plant die Anlage <!-- d-check:status-provenance -->
von `.githooks/commit-msg` — genau das Artefakt, dessen Anlage `ADR-0045`s
Entscheidungstext ausdrücklich verneint („Kein commit-msg-Hook"). Nach
`AGENTS.md` §3.5 / Modul 8 §Rollen-Regeln ist eine Accepted-ADR
immutable; ein Implementer (oder ein Folge-Slice, das sich stillschweigend
über die Klausel hinwegsetzt) darf ihr nicht stillschweigend widersprechen
— das wäre „Drift, kein pragmatisches Implementieren". Der korrekte Pfad
ist Verdikt 2 aus Modul 8 §Konflikt-Pfad: Folge-ADR mit `Supersedes`.

**Warum die Klausel dennoch korrigiert werden soll, statt `slice-073` <!-- d-check:status-provenance -->
fallenzulassen:** `ADR-0045`s eigene Ablehnungsbegründung zielt auf einen
Hook **als alleinige oder ersetzende Durchsetzung** — „Hook-Installation
ist client-seitig nicht erzwingbar (kein Wächter)". Genau dieses Risiko
benennt `slice-073`s Design bereits explizit und macht es zur Nicht-Lösung: <!-- d-check:status-provenance -->
Der Hook ersetzt das Standing-Gate nicht, bleibt Opt-in, und `make
commit-traceability`/das d-check-Modul `commits` bleiben laut
`slice-073` §1 die „kanonische, in `make gates` gebundene Prüfung". Das <!-- d-check:status-provenance -->
zweite Argument aus `ADR-0045` — „der Gate-Nachweis (record-gates) sieht
einen Hook nicht" — trägt unverändert und wird hier **bestätigt, nicht
widerlegt**: Der Hook bleibt bewusst außerhalb des Gate-Nachweises: eine
schnellere, lokale Vorab-Meldung für den committierenden Menschen/Agenten,
kein Ersatz für den Nachweis-Pfad. Die ursprüngliche Ablehnung war korrekt
für die Frage „ersetzt ein Hook das Standing-Gate" (Antwort: nein) —
und wird hier nicht revidiert. Sie beantwortete aber nie die andere Frage
„darf ein Hook **zusätzlich**, nicht-durchsetzend, neben dem Standing-Gate
existieren" — diese ADR beantwortet genau diese zweite, bislang offene
Frage.

## Entscheidung

Wir korrigieren `ADR-0045`s Entscheidungs-Klausel „Kein commit-msg-Hook"
auf: **Ein lokaler, opt-in `commit-msg`-Hook ist zulässig, sofern er
strukturell nicht-durchsetzend bleibt** — konkret:

1. Der Hook ersetzt **nicht** das Standing-Gate (`make
   commit-traceability` bleibt die einzige im Gate-Nachweis (`record-gates`)
   sichtbare, kanonische Prüfung).
2. Aktivierung bleibt lokale, dokumentierte Nutzer-Entscheidung
   (`git config core.hooksPath`) — kein versionierter, automatisch
   wirksamer Mechanismus.
3. Der Hook ist bash-only (kein Docker-Aufruf, keine Latenz-Regression
   für `git commit`) und spiegelt exakt die zwei bestehenden Regeln
   (positiv: ≥ 1 `LH-*`/`ADR-*`-Kennung im Betreff; negativ: keine
   `SPEC-*`/`ARC-*`-Struktur-ID im Betreff) **inklusive** der
   Merge-/Revert-Ausnahme, die die d-check-Positiv-Hälfte bereits trägt
   (`.d-check.yml` `exempt-pattern: '^(Merge |Revert )'`) — ohne diese
   Ausnahme würde der Hook jeden Merge-Commit fälschlich zurückweisen, den
   das bestehende Gate zulässt, und die beiden Prüfungen liefen bei
   Merge-Commits strukturell auseinander.
4. `--no-verify` bleibt ein gültiger, unveränderter Umgehungsweg — der
   Hook ist kein Sicherheitsmechanismus.

Alle übrigen Entscheidungen aus `ADR-0045` (Standing-Gate im
`make gates`-Bündel, Fenster `HEAD~5..HEAD`, Zwei-Werkzeug-Aufteilung
zwischen d-check-Modul `commits` und `tools/harness/commit-traceability.sh`,
`RANGE`-Override, `exempt-pattern`) bleiben unverändert bestehen und
werden hier nicht wiederholt — verbindlich bleibt `ADR-0045` selbst für
alles außer dieser einen Klausel.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun; `ADR-0045`s Ablehnung bleibt bestehen, `slice-073` wird verworfen oder gestrichen | keine neue ADR nötig; kein Duplikations-Risiko zwischen Hook und Standing-Gate | die belegte 3×-Beobachtung (`BEO-PGC/commit-traceability-kein-vorab-hook`) bleibt ungelöst — Verstöße bleiben weiterhin erst nach `git push` sichtbar, teils nicht mehr per `--amend` korrigierbar (Review zu `slice-041`); die Register-Schwelle fordert einen Ausgang, „entfallen" träfe nur mit Begründung zu, die hier nicht vorliegt (die Beobachtung ist real, nicht obsolet) | <!-- d-check:status-provenance -->
| B — `ADR-0045` vollständig neu fassen (ganze ADR Superseded, nicht nur eine Klausel) | ein einziges, vollständiges Nachfolge-Dokument statt zweier ADRs, die zusammengelesen werden müssen | unverhältnismäßig — 95 % von `ADR-0045`s Inhalt (Standing-Gate, Fenster, Werkzeug-Teilung) bleiben unverändert richtig; eine volle Neufassung dupliziert stabilen Text und ist genau das Muster, von dem `ADR-0048` (partielle Korrektur von `ADR-0047`, nur ein Grant-Text) bereits abweicht |
| C — Hook ruft d-check/Docker direkt auf (keine Duplikat-Implementierung) | eine einzige Implementierung der Regeln, kein Drift-Risiko zwischen Hook und Standing-Gate | jeder lokale `git commit` bräuchte einen vollen Containerstart (mehrere Sekunden) — zieht die Docker-Abhängigkeit in einen Pfad, der bislang bewusst ohne sie auskam; bereits im vorausgehenden Architect-Zug „geprüft und verworfen" |
| **D — Klausel-Korrektur per Folge-ADR: Hook zulässig, aber strukturell nicht-durchsetzend (gewählt)** | löst die belegte 3×-Beobachtung (schnellere Vorab-Meldung); respektiert `ADR-0045`s Kernargument vollständig (Hook ersetzt nie das Standing-Gate, Gate-Nachweis bleibt hook-blind); folgt dem im Repo bereits etablierten Muster einer engen Klausel-Korrektur (`ADR-0048` Supersedes `ADR-0047`) | zwei ADRs (`0045`+`0062`) müssen zusammengelesen werden, um die volle Hook-Politik zu verstehen — mildert die Immutabilitäts-Regel bewusst nicht ab, sondern befolgt sie (keine stille Bearbeitung von `ADR-0045`) |

## Konsequenzen

- Positiv: Die belegte 3×-Beobachtung
  (`BEO-PGC/commit-traceability-kein-vorab-hook`) bekommt einen
  ADR-gedeckten Umsetzungspfad — `slice-073` darf `.githooks/commit-msg` <!-- d-check:status-provenance -->
  anlegen, ohne einer Accepted-ADR stillschweigend zu widersprechen.
- Positiv: `ADR-0045`s zentrales Argument (Hook ersetzt nie das
  Standing-Gate) bleibt vollständig erhalten und wird hier verschärft
  festgeschrieben (Punkt 1–4 der Entscheidung), nicht aufgeweicht.
- Negativ: Zwei Regel-Implementierungen derselben zwei Prüfungen
  (Hook + Standing-Gate) — ein Drift-Risiko, das `slice-073` bereits als <!-- d-check:status-provenance -->
  benanntes Risiko in seinem §6 trägt; diese ADR löst es nicht auf,
  sondern macht es zur bewussten, dokumentierten Abwägung statt einer
  stillen.
- Folgepflicht: `slice-073`s Kopf-Feld „Bezug" und Plan-Zeile zu <!-- d-check:status-provenance -->
  `ADR-0045` werden auf diese ADR nachgezogen (Korrektur in diesem Zug);
  `slice-073`s Plan-Tabelle (§3) bekommt eine Zeile für die <!-- d-check:status-provenance -->
  Merge-/Revert-Ausnahme im Hook-Skript, `slice-073`s DoD bekommt einen <!-- d-check:status-provenance -->
  eigenen Rückweisungs-/Erfolgsbeleg für einen Merge-Commit-Versuch.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| `.githooks/commit-msg` (Umsetzung in `slice-073`) | spiegelt exakt die zwei Regeln aus `tools/harness/commit-traceability.sh` + d-check-Modul `commits`, inklusive `exempt-pattern: '^(Merge |Revert )'`; kein Aufruf in `make gates` (nicht Teil des Gate-Nachweises) | kein Gate — lokaler Hook, siehe `slice-073` DoD für reale Rückweisungs-/Erfolgsbelege | <!-- d-check:status-provenance -->

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

Beobachtbare Trigger: **(a)** dieselben Trigger wie `ADR-0045` §Re-
Evaluierungs-Trigger (a)/(b) — trifft einer von beiden ein, wird auch diese
Klausel im selben Folge-ADR-Zug mitgeprüft; **(b)** der Hook und das
Standing-Gate weichen real auseinander (eine der beiden Seiten wird
geschärft, die andere nicht nachgezogen — das in Punkt „Konsequenzen"
benannte Drift-Risiko tritt ein) — dann Vereinheitlichung oder Rückbau des
Hooks als Folge-ADR prüfen. Andernfalls permanent.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-14 | Accepted — Anlass: unabhängige Architect-Gegenprüfung des ersten Architect-Verdikts zur Commit-Traceability (kein Vorab-Hook) deckte auf, dass `slice-073` einer bereits in `ADR-0045` getroffenen Entscheidung ohne Folge-ADR widersprochen hätte; Verdikt 2 aus Modul 8 §Konflikt-Pfad (Folge-ADR statt stiller Lockerung) | der Architect-Verdikt dieses Zugs (Gegenprüfung), `slice-073` (`open/`) | <!-- d-check:status-provenance -->
| 2026-09-18 | Zitat-Korrektur — `docs/reviews/**`-Pfade durch Kennung ersetzt (`ADR-0073`) | PENDING_COMMIT |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0062` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
