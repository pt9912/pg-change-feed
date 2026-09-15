# Slice slice-073: Lokaler `commit-msg`-Git-Hook für Commit-Traceability

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — die Closure-Bedingung ist die eigene DoD dieses
Slice, kein *Mehr* jenseits ihrer (Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht).

**Bezug:** [`ADR-0062`](../../adr/0062-lokaler-commit-msg-hook-ergaenzt-standing-gate.md)
(bindend — korrigiert `ADR-0045`s Klausel „Kein commit-msg-Hook" und
erlaubt diesen Slice erst; ohne `ADR-0062` widerspräche dieser Slice
`ADR-0045` stillschweigend, siehe
`docs/reviews/architect-verdict-commit-traceability-kein-vorab-hook-gegengeprueft.md`),
[`ADR-0045`](../../adr/0045-commit-traceability-standing-gate.md)
(weiterhin bindend für alles außer der von `ADR-0062` korrigierten
Klausel — Standing-Gate, Fenster-Semantik, Werkzeug-Aufteilung); die
getragene Regel bleibt `AGENTS.md` §5/`harness/README.md` §Traceability
rules, dieser Slice ergänzt eine schnellere, lokale Vorab-Meldung
derselben Regel, kein neuer Vertrag.

**Berührte Spec-Stellen:** — (kein Vertrags-/Technik-Bezug, reine
Entwickler-Werkzeugkette).

**Verantwortlich:** pt9912.

**Autor:** pt9912 (Planner-Lauf, Architect-Verdikt vorausgehend; Plan
korrigiert durch einen unabhängigen, gegenprüfenden Architect-Zug,
2026-09-14 — siehe
`docs/reviews/architect-verdict-commit-traceability-kein-vorab-hook-gegengeprueft.md`).
**Datum:** 2026-09-14.

---

## 1. Ziel und Abgrenzung

**Ziel:** Ein lokaler, opt-in `commit-msg`-Git-Hook (`.githooks/commit-msg`)
weist einen Commit bereits vor seinem Abschluss zurück, wenn sein Betreff
keine `LH-*`-/`ADR-*`-Kennung trägt oder eine `SPEC-*`-/`ARC-*`-Struktur-ID
enthält — dieselben zwei Regeln, die `make commit-traceability`
(`tools/harness/commit-traceability.sh` + d-check-Modul `commits`) heute
erst nachträglich meldet — **inklusive** der Merge-/Revert-Ausnahme, die
die d-check-Positiv-Hälfte bereits trägt (`.d-check.yml`
`exempt-pattern: '^(Merge |Revert )'`, siehe DoD und Plan-Tabelle unten).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Ersatz oder Änderung des bestehenden Standing-Gates (`ADR-0045`)** —
  `make commit-traceability` und das d-check-Modul `commits` bleiben die
  kanonische, in `make gates` gebundene Prüfung; der Hook ist eine
  zusätzliche, schnellere Vorab-Meldung, keine neue Quelle der Wahrheit
  (`ADR-0062` Entscheidung Punkt 1; ursprünglich Architect-Verdikt
  `docs/reviews/architect-verdict-commit-traceability-kein-vorab-hook.md`,
  gegengeprüft und in der ADR-Frage korrigiert durch
  `docs/reviews/architect-verdict-commit-traceability-kein-vorab-hook-gegengeprueft.md`).
- **Automatische, erzwungene Aktivierung** (z. B. ein Mechanismus, der
  `core.hooksPath` ungefragt setzt) — Git kennt keinen versionierten Weg,
  lokale Hook-Konfiguration ohne Zutun des Nutzers zu aktivieren, der nicht
  selbst ein zweites Vertrauensproblem wäre; Aktivierung bleibt ein
  dokumentierter, einmaliger Opt-in-Schritt (`ADR-0062` Entscheidung
  Punkt 2).
- **Docker-Aufruf innerhalb des Hooks** — ein `docker run`-Aufruf bei jedem
  lokalen `git commit` würde eine spürbare Latenz-Regression einführen, die
  heute nicht besteht (`git commit` selbst braucht bislang kein Docker);
  der Hook bleibt bash-only (`ADR-0062` Entscheidung Punkt 3).
- **CI-/Server-seitige Durchsetzung** — ein lokaler Hook lässt sich mit
  `--no-verify` umgehen und ersetzt keine serverseitige Kontrolle; das
  bestehende `make commit-traceability` in `make gates` bleibt die
  durchsetzende Instanz, dieser Slice fügt nur ein schnelleres
  Frühwarnsignal hinzu (`ADR-0062` Entscheidung Punkt 4).

**Keine Mindestzahl.** Ein Slice mit *einem* echten Ausschluss ist besser als
einer mit vier erfundenen; die vier Klassen sind ein Suchraster, keine
Ausfüll-Liste.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

- [x] `.githooks/commit-msg` weist real einen Commit-Versuch mit
      `SPEC-*`/`ARC-*`-Struktur-ID im Betreff zurück (nicht-Null-Exit vor
      Commit-Abschluss) — real ausgeführt, nicht nur am Skripttext behauptet.
- [x] `.githooks/commit-msg` weist real einen Commit-Versuch ganz ohne
      `LH-*`-/`ADR-*`-Kennung im Betreff zurück — real ausgeführt.
- [x] Ein regulärer, konformer Commit-Versuch läuft real ungehindert durch
      den Hook — real ausgeführt (kein Fehlalarm auf einem gültigen
      Betreff).
- [x] Ein Merge-Commit-Versuch (Betreff nach dem Muster `Merge …`, ohne
      `LH-*`/`ADR-*`-Kennung) läuft real ungehindert durch den Hook —
      real ausgeführt; der Hook spiegelt hier dieselbe Ausnahme wie die
      d-check-Positiv-Hälfte (`.d-check.yml` `exempt-pattern:
      '^(Merge |Revert )'`), sonst wiese er jeden regulären
      `git merge`-Commit fälschlich zurück, den das bestehende Gate
      zulässt (`ADR-0062` Entscheidung Punkt 3).
- [x] Aktivierung dokumentiert (`harness/README.md` oder
      `AGENTS.md`-Onboarding-Hinweis): einmaliger, expliziter
      `git config core.hooksPath .githooks`-Schritt, keine automatische
      Aktivierung.
- [x] `make gates` grün (der Hook selbst ist kein Gate-Ziel und wird von
      `make gates` nicht aufgerufen — er läuft ausschließlich lokal vor
      `git commit`).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Doku-Update: `harness/README.md` (Onboarding-Hinweis) trägt den neuen
      Opt-in-Schritt.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben,
      **falls dieser Slice einen Inventur-Fund auflöst** — entfällt: Repo
      ist GF (`harness/conventions.md` Modus-Deklaration `PGC`), keine
      `reconciliation.md` vorhanden.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen
      `evidence/`; keine Beobachtung angefallen ist ebenfalls eine Antwort
      und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen /
      weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen —
      im Repo **ohne** Wellen-Betrieb hier geprüft (dieser Slice trägt keine
      Welle).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `.githooks/commit-msg` | neu | bash-only Hook, spiegelt die zwei Regeln aus `tools/harness/commit-traceability.sh`/d-check-Modul `commits` als Regex-Abgleich auf `$1`, **inklusive** der Merge-/Revert-Ausnahme (`^(Merge |Revert )`) der d-check-Positiv-Hälfte — ohne sie liefe der Hook bei jedem Merge-Commit strukturell anders als das Standing-Gate. **Plan-Nachzug (Implementer-Lauf):** Die positive Hälfte prüft die ganze Message-Datei, nicht nur den Betreff (das d-check-Modul `commits` liest die gesamte Message — am gepinnten Image gemessen); die Grenz-Hälfte prüft den Betreff; die Merge-/Revert-Ausnahme waivet nur die positive Hälfte. So weist der Hook keinen Commit zurück, den das Standing-Gate zulässt |
| `harness/README.md` | update | Onboarding-Hinweis: optionaler `core.hooksPath`-Aktivierungsschritt |
| `docs/plan/adr/0045-commit-traceability-standing-gate.md` | keine Änderung | `Accepted`-ADR bleibt unverändert (`AGENTS.md` §3.5); die zuvor kollidierende Klausel „Kein commit-msg-Hook" ist bereits durch `ADR-0062` (Supersedes, nur diese Klausel) korrigiert — kein weiterer Eingriff in diesem Slice |
| `docs/plan/adr/0062-lokaler-commit-msg-hook-ergaenzt-standing-gate.md` | bereits vorhanden | von der unabhängigen Architect-Gegenprüfung vorab geschrieben, nicht Teil der Implementer-Arbeit dieses Slice |

## 4. Trigger

**Start** (`next` → `in-progress`): priorisiert, `Verantwortlich:` gesetzt,
WIP-Limit (1 je Implementer) frei.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Zeigt sich, dass
  das Regex-Muster die beiden bestehenden Prüfungen nicht deckungsgleich
  abbilden kann (z. B. weil das d-check-Modul `commits` zusätzliche,
  nicht rein betreff-basierte Bedingungen trägt, die ein Hook ohne
  Docker-Zugriff nicht real prüfen kann), gehört die Klärung dieser
  Deckungsgleichheit zurück zur Zerlegung, bevor der Hook geschrieben wird.
- `in-progress` → `open` (blockiert — Carveout?): Kein bekannter Blocker —
  `ADR-0062` und das vorausgehende Architect-Verdikt treffen die
  Design-Fragen bereits vorab (bash-only, Opt-in-Aktivierung, kein
  Docker-Aufruf, Merge-/Revert-Ausnahme).

## 5. Closure-Trigger

DoD vollständig **und** `make gates` grün **und** Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

- Der Hook dupliziert die beiden Regeln aus `tools/harness/commit-traceability.sh`
  und dem d-check-Modul `commits` als zweite, unabhängige Implementierung —
  driftet eine der beiden Seiten (z. B. eine künftig geschärfte Regel im
  d-check-Modul), ohne dass die andere nachgezogen wird, meldet der Hook
  fälschlich grün oder rot. `ADR-0062` benennt dies als bewusst
  dokumentiertes, nicht aufgelöstes Restrisiko (Re-Evaluierungs-Trigger
  (b)). — **Ausgang:** <bei Closure zu füllen>
- Der Hook ist per `--no-verify` umgehbar und aktiviert sich nicht
  automatisch bei einem neuen Klon/einer neuen Session (`core.hooksPath`
  ist lokale Konfiguration) — ein Nutzer, der den Aktivierungsschritt nie
  ausführt, bleibt ohne Vorab-Meldung und verlässt sich weiterhin
  ausschließlich auf `make gates`. — **Ausgang:** <bei Closure zu füllen>

## 7. Closure-Notiz

<…>

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Einzige berührte Sub-Area ist die
Repo-weite Default-Sub-Area `*`/`PGC` (`harness/conventions.md` führt keine
eigene Sub-Area für Git-/Harness-Tooling außerhalb der Produktionsschichten).

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen
(`docs/plan/planning/observations/BEO-PGC/`), insbesondere gegen Tooling/
Traceability:

- `BEO-PGC/commit-traceability-kein-vorab-hook` — Auslöser dieses Slice
  selbst (3×, Ausgang *geplant* → dieser Slice, siehe Architect-Verdikt
  `docs/reviews/architect-verdict-commit-traceability-kein-vorab-hook.md`
  und dessen unabhängige Gegenprüfung
  `docs/reviews/architect-verdict-commit-traceability-kein-vorab-hook-gegengeprueft.md`).
  Kein weiterer Treffer nötig — die Beobachtung *ist* der Auftrag.
- `BEO-PGC/coverage-stage-dockerignore-blockiert-tooling` — **offen**, 1×.
  Kein Treffer: Dieser Slice führt keine neue Docker-Build-Stage ein, der
  Hook läuft ausschließlich lokal per bash, kein Container-Build.
- Übrige Einträge (Planning-Harness-Prozess, Rollen-/Adapter-spezifische
  Beobachtungen): keine Treffer.

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas GF (nur
`*`/`PGC`).

Alle berührten Sub-Areas GF (nur `*`/`PGC`).
