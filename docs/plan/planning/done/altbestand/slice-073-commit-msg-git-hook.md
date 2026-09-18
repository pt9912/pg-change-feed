# Slice slice-073: Lokaler `commit-msg`-Git-Hook für Commit-Traceability

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — die Closure-Bedingung ist die eigene DoD dieses
Slice, kein *Mehr* jenseits ihrer (Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht).

**Bezug:** [`ADR-0062`](../../adr/0062-lokaler-commit-msg-hook-ergaenzt-standing-gate.md)
(bindend in ihren Punkten 1, 2 und 4 — das Standing-Gate bleibt die
durchsetzende Instanz, die Aktivierung bleibt lokaler Opt-in, `--no-verify`
bleibt ein gültiger Umgehungsweg; sie korrigiert `ADR-0045`s Klausel „Kein
commit-msg-Hook" und erlaubt diesen Slice erst; ihr Punkt 3 ist durch
`ADR-0069` ersetzt),
[`ADR-0069`](../../adr/0069-commit-msg-hook-einseitige-zusage.md)
(Teil-Supersede von `ADR-0062` Punkt 3 — einseitige Zusage des Hooks:
message-weite positive Hälfte, Hook-Logik unverändert; die Klassengrenze ist
nachgezogen durch
[`ADR-0070`](../../adr/0070-supersede-reichweite-und-klassengrenze.md) —
drei benannte Divergenz-Klassen (a)/(b)/(c) sind **gemessen, nicht
erschöpfend**; eine vierte ist benannt: führende Leerzeile vor einem
`Merge …`-Betreff unter `--cleanup=verbatim`),
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
weist einen Commit bereits vor seinem Abschluss zurück, wenn die Message
keine `LH-*`-/`ADR-*`-Kennung trägt (**message-weit** — die ganze rohe
Message-Datei, nicht nur der Betreff) oder wenn der Betreff eine
`SPEC-*`-/`ARC-*`-Struktur-ID enthält — dieselben zwei Regeln, die
`make commit-traceability` (`tools/harness/commit-traceability.sh` +
d-check-Modul `commits`) heute erst nachträglich meldet — **inklusive** der
Merge-/Revert-Ausnahme, die die d-check-Positiv-Hälfte bereits trägt
(`.d-check.yml` `exempt-pattern: '^(Merge |Revert )'`, siehe DoD und
Plan-Tabelle unten).

Die Zusage des Hooks ist **einseitig**
([`ADR-0069`](../../adr/0069-commit-msg-hook-einseitige-zusage.md)): er weist
keinen Commit zurück, den das Standing-Gate zulässt; **vollständig** ist er
nicht — drei benannte Divergenz-Klassen (a)/(b)/(c) zum Modul bleiben offen
(Kennung nur in einer `#`-Kommentarzeile · Kennung nur hinter der
scissors-Zeile/im Verbose-Diff · Struktur-ID auf der Fortsetzungszeile des
ersten Absatzes). Die drei Klassen sind **gemessen, nicht erschöpfend**
([`ADR-0070`](../../adr/0070-supersede-reichweite-und-klassengrenze.md));
eine vierte, ebenfalls laxere Klasse ist benannt — eine führende Leerzeile vor
einem `Merge …`-Betreff, nur unter `--cleanup=verbatim` erreichbar: der Hook
lässt durch, das Modul färbt rot.

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
      Commit-Abschluss) — real ausgeführt, nicht nur am Skripttext behauptet;
      der Nachweis führt einen **isolierenden** Input (Kennung **und**
      Struktur-ID, `docs: update SPEC-012 table (ADR-0045)`), den allein die
      Grenz-Hälfte zurückweist — die positive Hälfte passiert ihn.
- [x] `.githooks/commit-msg` weist real einen Commit-Versuch ganz ohne
      `LH-*`-/`ADR-*`-Kennung in der Message zurück (message-weit) — real
      ausgeführt.
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
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8);
      Erstbefund [`docs/reviews/review-slice-073.md`](../../../reviews/review-slice-073.md),
      Bestätigungslauf nach der Fixrunde
      [`docs/reviews/review-slice-073-fixrunde.md`](../../../reviews/review-slice-073-fixrunde.md).
- [x] Doku-Update: `harness/README.md` (Onboarding-Hinweis) trägt den neuen
      Opt-in-Schritt.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag (§7).
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben,
      **falls dieser Slice einen Inventur-Fund auflöst** — entfällt: Repo
      ist GF (`harness/conventions.md` Modus-Deklaration `PGC`), keine
      `reconciliation.md` vorhanden.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen
      `evidence/`; keine Beobachtung angefallen ist ebenfalls eine Antwort
      und wird in §7 notiert. Zwei neue Einträge
      (`BEO-PGC/spiegelung-ist-approximation`,
      `BEO-PGC/vorab-bedingung-nach-umsetzung-geprueft`) und ein überführter
      Ausgang (`BEO-PGC/commit-traceability-kein-vorab-hook`: `geplant` →
      `verkörpert`).
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen /
      weiter offen) — beide Ausgänge stehen in §6.
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen —
      im Repo **ohne** Wellen-Betrieb hier geprüft (dieser Slice trägt keine
      Welle); geprüft **nach** dem `git mv` nach `done/`.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `.githooks/commit-msg` | neu | bash-only Hook, Regex-Abgleich auf `$1`; die getragene Zusage, die Prüf-Weiten der beiden Hälften (positiv message-weit, Grenze betreff-scoped) und die drei bewusst offenen Divergenz-Klassen trägt [`ADR-0069`](../../adr/0069-commit-msg-hook-einseitige-zusage.md); die Klassen sind **gemessen, nicht erschöpfend** — eine vierte ist benannt (führende Leerzeile vor `Merge …` unter `--cleanup=verbatim`), [`ADR-0070`](../../adr/0070-supersede-reichweite-und-klassengrenze.md) |
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
  (b)). — **Ausgang:** *weiter offen → Beobachtungs-Register*: die
  Prämisse des Risikos (zwei deckungsgleiche Implementierungen, die
  auseinanderdriften können) ist durch die Messung widerlegt — die beiden
  waren nie deckungsgleich; was bleibt, ist die gemessene Laxheit des Hooks
  in den benannten Klassen, erreichbar und nicht geschlossen. Eingetragen als
  `BEO-PGC/spiegelung-ist-approximation`, Wächter sind
  [`ADR-0070`](../../adr/0070-supersede-reichweite-und-klassengrenze.md)
  §Re-Evaluierungs-Trigger (a) und (c). Der Ausgang *eingetreten* trägt hier
  nicht: er verlangt Carveout oder Folge-Slice, und die Auflösung ist eine
  Fest-Schreibung, keine Ausnahme auf Zeit — der Architect hat die
  nachziehende Option als eigene entschieden verworfen.
- Der Hook ist per `--no-verify` umgehbar und aktiviert sich nicht
  automatisch bei einem neuen Klon/einer neuen Session (`core.hooksPath`
  ist lokale Konfiguration) — ein Nutzer, der den Aktivierungsschritt nie
  ausführt, bleibt ohne Vorab-Meldung und verlässt sich weiterhin
  ausschließlich auf `make gates`. — **Ausgang:** *entfallen — gestrichen mit
  Begründung*: der beschriebene Zustand ist die entschiedene Eigenschaft des
  Hooks (lokaler Opt-in, kein Ersatz des Gates —
  [`ADR-0062`](../../adr/0062-lokaler-commit-msg-hook-ergaenzt-standing-gate.md)
  Punkte 1 und 2), kein Risiko: die Zusage „der Hook ersetzt das
  Standing-Gate nicht" deckt genau diesen Fall, es entsteht keine falsche
  Sicherheit. `--no-verify` ist derselbe Weg und ebenso entschieden.

## 7. Closure-Notiz

**Geliefert:** der lokale, nicht-durchsetzende `commit-msg`-Hook
(`.githooks/commit-msg`), seine Aktivierungs-Doku (`harness/README.md`
§Traceability rules) und die Nachzüge aus dem Konflikt-Pfad. Der Hook meldet
vor dem Commit-Abschluss und ersetzt das Standing-Gate nicht; seine Zusage ist
**einseitig** (er weist keinen Commit zurück, den das Gate zulässt) und
**nicht vollständig** — die Klassen, in denen er blind ist, sind benannt und am
gepinnten Image gemessen
([`ADR-0069`](../../adr/0069-commit-msg-hook-einseitige-zusage.md),
[`ADR-0070`](../../adr/0070-supersede-reichweite-und-klassengrenze.md)).

**Closure-Kriterien** (§5): DoD vollständig, `make gates` grün, diese Notiz.
Verifikation:
[`docs/reviews/verify-slice-073.md`](../../../reviews/verify-slice-073.md) —
DoD- und Entscheidungs-Konformität bestätigt, kein blockierender Befund;
davor der Erstbefund
[`review-slice-073.md`](../../../reviews/review-slice-073.md) und der
Bestätigungslauf
[`review-slice-073-fixrunde.md`](../../../reviews/review-slice-073-fixrunde.md).

**Was funktionierte:** die Rollentrennung hat den Fall getragen. Der
Implementer hat die Abweichung von der Entscheidung begründet statt still
vorgenommen, der Reviewer hat sie nicht herabgestuft, sondern als HIGH mit
Rollen-Widerspruch in den Konflikt-Pfad gegeben, und der Architect hat den
Träger der Zusage von einer Plan-Zelle in eine Entscheidung verschoben. Der
Bestätigungslauf hat die Fixrunde als **frisches Artefakt** geprüft statt als
Bestätigung des eigenen Befunds — der Nachweis, dass die Hook-Logik unverändert
ist, läuft dort über den `sha256` der Nicht-Kommentar-Zeilen.

**Was anders lief — §4-Grund:** §4 nannte die Kongruenz-Frage vorab und
verlangte ihre Klärung „bevor der Hook geschrieben wird". Die Bedingung trat
ein (der Hook bildet die beiden Prüfungen nicht deckungsgleich ab); der Slice
folgte ihr aber nicht durch `in-progress` → `next`, weil das Artefakt zu diesem
Zeitpunkt bereits geschrieben war und die Deckungsgleichheit strukturell
unerreichbar ist — der Hook läuft vor `git`s Cleanup. Die Antwort war eine
Entscheidung, kein Neuschnitt; der Preis der Reihenfolge war der volle
Konflikt-Pfad (zwei HIGH, zwei Review-Runden, zwei Folge-ADRs).

**Steering-Loop-Einträge:**

- *geschärfte Regel* — verkörpert in
  [`ADR-0069`](../../adr/0069-commit-msg-hook-einseitige-zusage.md): wer eine
  bereits gate-getragene Regel außerhalb des Gates nachbaut, beansprucht keine
  Fidelität, sondern formuliert seine Zusage einseitig und benennt die
  Klassen, in denen er blind ist; der §Warum-die-Approximation-Absatz trägt
  den Grund allgemein. Herkunfts-Anker `seit slice-073`.
- *benannte Spec-Lücke* — §4 des Slice-Plans verlangt, die
  Rückführungs-Bedingung vorab zu benennen, aber keine Stelle verlangt, sie
  **beim Schreiben der Umsetzung** auszuwerten; gelesen wird sie erst bei der
  Closure. Beobachtung: `BEO-PGC/vorab-bedingung-nach-umsetzung-geprueft`.
- *kein neuer Sensor* — die Divergenz ist am gepinnten Image messbar, aber der
  Hook ist kein Gate-Ziel; ein Sensor müsste ihn aufrufen, was `ADR-0062`
  Punkt 1 ausschließt.

**Beobachtungs-Register:** zwei neue Einträge (s. o.) und der überführte
Ausgang `BEO-PGC/commit-traceability-kein-vorab-hook` (`geplant` →
`verkörpert`, derselbe Anker). Kein Eintrag hat mit diesem Slice die
3×-Schwelle erreicht. Der Register-Stand ist **vor** dem `git mv` geschrieben;
die Lage des Belegs prüft die Register-Paarung danach.

**Paarungen (nach dem `git mv`):** *Anker* — der Herkunfts-Anker der
verkörperten Regel steht an einem veränderlichen Träger, `harness/README.md`
§Traceability rules (`· seit slice-073`); die Accepted-ADR
[`ADR-0069`](../../adr/0069-commit-msg-hook-einseitige-zusage.md) selbst kann
ihn nicht nachtragen (Immutabilität, `AGENTS.md` §3.5) und ist deshalb
Zielort, nicht Ankerträger. *Folge-Slice* — keiner genannt. *Register* — beide
genannten Einträge existieren; die zweite Hälfte der Prüfung („jedes
Verzeichnis trägt mindestens einen Beleg") fand **zwei Vorbestände** mit leerem
`evidence/` trotz Ausgang `verkörpert`: `BEO-PGC/verwaltung-keine-sql-administration`
und `BEO-PGC/retention-keine-loeschausfuehrung`. Beide Belege sind aus dem
jeweils eigenen `state.md` rekonstruiert, das den auflösenden Vorgang selbst
nennt, und nachgetragen (`evidence/slice-036.md`, `evidence/slice-043.md`).
**Benannte Lücke:** für diese zweite Hälfte der Register-Paarung führt das Repo
keinen Sensor — sie ist nur durch die Closure selbst geprüft, ein
`verkörpert`-Eintrag ohne Beleg bleibt für jedes Gate unsichtbar.

**Folge-Slices:** keine. Die vom Architect verworfene Option — den Hook
nachziehen — wäre der einzige Folge-Schnitt gewesen; sie ist als Entscheidung
ausgeschieden ([`ADR-0069`](../../adr/0069-commit-msg-hook-einseitige-zusage.md)
§Verglichene Alternativen, Option B).

**Archivierung:** nicht ausgeführt — das Repo führt kein Archivierungswerkzeug
(kein `*-archiv.zip` unter `done/`); nach Modul 6 bleibt es ohne sie konform,
und diese Feststellung ersetzt den Handlauf.

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
