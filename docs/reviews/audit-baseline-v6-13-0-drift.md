# Audit: Baseline-Drift v6.9.0 → v6.13.0

**Auftrag:** `slice-harness-baseline-v6-13-0` §2 DoD 3 — Diff über beide
Bäume, je Modul Fund/Nichtfund, die verkörperte Form gegen die Fundstellen.
**Rolle:** Architect. **Datum:** 2026-09-29. **Bezug:**
[ADR-0045](../plan/adr/0045-commit-traceability-standing-gate.md).

Zahlen in diesem Bericht sind **gemessen** (Kommando/Ausgabe genannt) oder als
**erwartet/hergeleitet** gekennzeichnet.

---

## 1. Bezug und Verifikation (DoD 1)

| Schritt | Ergebnis |
|---|---|
| Download | `lab-regelwerk.zip` des Releases `v6.13.0` Docker-gekapselt über das gepinnte Toolchain-Image (`golang:1.27-alpine`, Digest `sha256:cf6fca66…3125`) mit `wget --network host` in ein Wegwerf-Verzeichnis außerhalb des Repos — Ausgabe `DOWNLOAD_OK` |
| Asset-Prüfsumme (gemessen) | `b5151e77807e2affebb25cfab9be24b88cf43afc42c075db0b982a2ff1052b96` |
| Entpacken | im Container; Struktur `regelwerk/` + `templates/` parallel — **das Asset trägt kein `SHA256SUMS`**; die Datei ist Bootstrap-eigenes Erzeugnis |
| `SHA256SUMS` erzeugt | nach der Methode des v6.9.0-Bootstraps (`47d86b26`: „eigene SHA256SUMS erzeugt, gleiche Methode wie tools/harness/baseline-verify.sh"): `find . ! -type d ! -path ./SHA256SUMS \| sed 's\|^\./\|\|' \| LC_ALL=C sort \| xargs sha256sum > SHA256SUMS` — 54 Einträge |
| Integritätsprüfung (gemessen) | `sha256sum -c SHA256SUMS` im Container: **54/54 OK, Exit 0**; zweite Prüfung nach `cp -r` in `.harness/baseline/v6.13.0/`: **54/54 OK, Exit 0** |
| Commit | `d443ee39` — 55 Dateien, 8829 Zeilen |

Der Rückführungs-Fall §4 („`SHA256SUMS` löst gegen das Asset nicht auf →
`open`") ist **nicht** eingetreten.

## 2. Diff-Umfang (gemessen)

- **Datei-Menge identisch:** 55 Pfade in beiden Bäumen (`regelwerk/` 26,
  `templates/` 28, `SHA256SUMS`). **0 neue, 0 gelöschte, 0 umbenannte
  Module.**
- **Inhaltlich unterschiedlich:** 35 Dateien — 26 `regelwerk/`-Dateien
  (alle, je „Quelle"-URL-Bump plus die Funde unten), 8 Templates, 1×
  `SHA256SUMS` (per Definition). 20 Templates sind byte-gleich — darunter
  **beide Reviewer-Skills** (`reviewer.template.md`,
  `closure-note-reviewer.template.md`, Modul-10/11-Form unverändert).
  Ursprung der Zählung: Nachmessen gegen den committeten Stand `d443ee39`
  (Blob-Vergleich über `git archive` + `diff -rq`, byte-gleich zur
  `cmp`-Messung des Reviews, F-1) — die Zählung der ersten Berichtsfassung
  (36/9/19) war falsch.
- **Differenz-Text:** 768 Zeilen (`diff -ru` über `regelwerk/`), 285 Zeilen
  (über `templates/`) — gegen den committeten Stand nachgemessen, unverändert.
- **Struktur erhalten:** über alle `regelwerk/`-Dateien ist genau **ein**
  `##`/`###`-Header neu (`### Nachzug ist keine Überschreibung`, Modul 4),
  **0** entfernt, **0** umbenannt; alle `<a id="…">`-Anker identisch
  (gemessen per Header- und Anker-Diff). Alle in der verkörperten Form
  zitierten Abschnitte lösen in v6.13.0 weiter auf (14 Abschnitts-Zitate
  gegen v6.13.0 geprüft, gemessen).
- `regelwerk/README.md`: Stand-Zeile `Kurs-Welle 137 · 2026-09-16` →
  `Kurs-Welle 153 · 2026-09-28`.

## 3. Fundstellen je Modul

| Modul | Fund | Art |
|---|---|---|
| `grundlagen-durchsetzungsschicht.md` | Workflow-Skelett (Slash-Commands): aus dem Register graduierte Regeln tragen Auflösungs-Trigger oder *permanent*; ist die Regel mechanisierbar, ist das Scharfschalten des Gates selbst der Trigger — Prosa-Entfernung ist DoD-Punkt **desselben** Slice | neu (Absatz) |
| `grundlagen-harness-dateien.md` | Schnitt-Prinzip (Kurzform im Index, Volltext in eigener Datei unter `harness/rules/<name>.md`) auf jede Guide-Datei — `AGENTS.md` eingeschlossen | neu (Absatz) |
| `grundlagen-klassifikation.md` | zwei neue Entropy-Klassen: **AGENTS.md-Wildwuchs** (in den Trigger-Audit aufgenommen) und **Guide-Datei-Wildwuchs** (Summenklasse; Gegenmittel: Zeilen-Obergrenze + Index-plus-Datei-Schnitt) | neu (Listeneinträge) |
| `grundlagen-referenz-richtung.md` | Spec-Historie nennt **keine** ADR und keinen Slice, sondern beim Vertrag den externen CR | verschärft |
| `grundlagen-source-precedence.md` | **neue ID-Reihe `<PREFIX>-RB-<NN>`** (Randbedingungen) im Vertrags-Stratum — ID-Tabelle und Mermaid-Diagramm | neu |
| `modul-00/01/03/04/11/12/13` | „Regeln gegen typische Fehlannahmen" in Zitat-Form umgestellt (cosmetic, Inhalt identisch); Modul 3 zusätzlich **vier neue** Regeln: Randbedingung ≠ Out-of-Scope, Usability/Internationalisierung in der Spec, Systemgrenze im Lastenheft, Prompts ersetzen keine Specs; Modul 3 + „Given/When/Then macht noch kein Lastenheft" (User Story ist Slice-Klasse); Modul 1 + „Validierung hat hier bewusst keine Station" | neu/umgestellt |
| `modul-02` | Bootstrap-Schritt 4: Lastenheft-Outline um `LH-RB-*` | erweitert |
| `modul-04` | **(a)** „Eine Gate-Erweiterung ist nicht automatisch ein ADR-Anlass" — Aufnahme eines bereits existierenden, unabhängig lauffähigen Wächters in `make gates` verweist auf dessen ADR, braucht keine eigene; neue Fehlerklasse/neuer Scope: weiterhin ADR. **(b)** neuer Abschnitt „**Nachzug ist keine Überschreibung**": Referenz-/Pfad-Nachzug und Template-Feld-Nachzug an `Accepted`-ADRs sind keine inhaltliche Überschreibung | neu |
| `modul-05` | `Verantwortlich:` trägt bei Parallel-Läufen Person **und** Zweig; **WIP-Limit pro Lauf, nicht pro Rolle**; Slice-Größe „mehrere" → „mehr als zwei" Schichten | verschärft |
| `modul-06` | Trigger-Audit um die vierte Klasse **Hard Rule** erweitert (Carveout · bootstrap-aware Gate · ADR · Hard Rule); „Verkörpert heißt nicht zwangsläufig automatisiert" (4. Auftreten einer Fehlerklasse → mechanischer Sensor oder begründete Absage); neuer Lese-Schritt für in `open/`/`next/` liegengebliebene Slices; „Erst gruppieren, dann entscheiden" | neu/erweitert |
| `modul-07` | Carveout mit Schwelle nennt die ADR, die die Schwelle setzt; ADR-Pfad korrigiert (`docs/plan/adr/<NNNN>-*.md`) | neu/korrigiert |
| `modul-08` | Übergabe-Tabelle: Trigger-Audit um Hard-Rule-Zweig; Rolleninhaber des WIP-Limits in dritter, engerer Lesart **pro Lauf** | erweitert |
| `modul-09` | „Gates dürfen nicht ohne ADR gelockert werden" geschärft: befristete Ausnahme für einen Teil = Carveout mit Trigger und Folge-Slice, keine Senkung; Hard-Rule-Liste trägt dasselbe Schnitt-Prinzip | geschärft |
| `modul-10/14/15/16` | nur „Quelle"-URL-Bump | kosmetisch |
| `templates/AGENTS.template.md` | §5 als Index-Tabelle (Regel-Datei-Spalte, `harness/rules/<name>.md`); §3.6 um Carveout-Formulierung; ID-Default um `<PREFIX>-RB-<NN>` | erweitert |
| `templates/.d-check.yml` | `ids`-ADR-Muster um optionales Bereichssegment + `link-policy: always`; matrix um `welle`/`carveout`/`roadmap`-Klassen und Regeln; neuer `file.max-lines`-Ratchet-Block (Guide-Datei-Wildwuchs) | erweitert |
| `templates/spec/lastenheft.template.md` | §4 „Nichtfunktionale Anforderungen und Randbedingungen" mit `LH-RB-`-Format und Systemgrenze-Leitstand | erweitert |
| `templates/` (übrige 5) | `NNNN-titel`/`slice.template`: RB-Bezug; `carveout.template`: Schwellen-ADR-Vermerk; `conventions.template`: URL + RB im MR-000-Default; `gate.template`: lebende Artefakte linken die Sensor-Datei direkt, einfrierende nennen `make <target>` als Token | erweitert |

## 4. Verkörperte Form gegen die Fundstellen

### 4.1 Pflicht-Nachzüge (forciert — durch die v6.9.0-Entfernung und das `versions`-Modul)

| # | Stelle | Änderung | Zwang |
|---|---|---|---|
| 1 | `AGENTS.md` §1 (URL-Zeile) | Download-URL `v6.9.0` → `v6.13.0` | DoD |
| 2 | `harness/conventions.md` §Baseline (Stand-Zeile, Adoptions-Zeile) | Stand `v6.13.0`, Datum der Adoption 2026-09-29, Release-URL, adoptierter Stand `Kurs-Welle 153 · 2026-09-28` | DoD |
| 3 | `harness/conventions.md` MR-000 (ID-Schema-Aufzählung) | um `<PREFIX>-RB-*` ergänzen — MR-000 behauptet „keine inhaltlichen Adaptionen ggü. Baseline-Default"; der Default ist um RB gewachsen, die Aufzählung muss ihn wiedergeben | Konsistenz |
| 4 | `harness/conventions/MR-001…MR-004` (je 1 Linkzeile) | `../../.harness/baseline/v6.9.0/…` → `v6.13.0` | **gate-fatal** — `harness/conventions/**` liegt unter `scan.roots ["."]`, `tracked` löst die Linkziele, `versions` meldet abweichende Pins |
| 5 | `.claude/agents/architect.md` (Vorlage-Zeile), `.claude/agents/reviewer.md` (Baseline-Zeile) | Pfad `v6.9.0` → `v6.13.0` | toter Pointer nach der Entfernung; `versions`-Pin (`.claude/**` ist gescannt) |
| 6 | `.harness/skills/closure-note-reviewer.md` (Template-Pfad) | `v6.9.0` → `v6.13.0` | toter Pointer nach der Entfernung (nicht gescannt — `.harness/**` ausgenommen) |
| 7 | `harness/sensors/baseline-verify.md` (Bindung-Abschluss) | „der adoptierte Stand `v6.9.0`" → `v6.13.0` | überholte Aussage über den eigenen Referenten (kein Pin-Muster, kein Gate-Zwang) |

Alle sieben sind **Ein-Zeilen-Mechanik** (Versions-/Pfad-Bump), kein
Regel-Nachzug.

### 4.2 Das `versions`-Modul nach dem Konventions-Update

`current-from: harness/conventions.md#baseline`, `pin-pattern: '\.harness/baseline/(v\d+\.\d+\.\d+)/'`,
Exemptions: `harness/conventions/done/**`, `docs/reviews/**`. Sobald §Baseline
auf `v6.13.0` steht, ist jeder restliche `v6.9.0`-Pin ein Befund
(**erwartet** — aus der Modulkonfiguration hergeleitet, nicht gemessen). Gemessene
Pin-Treffer außerhalb der Baseline (13 Zeilen in 12 Dateien, beide Stände
identisch):

- forciert nachzuziehen: MR-001…MR-004, `architect.md`, `reviewer.md` (4.1)
- [ADR-0095](../plan/adr/0095-review-klasse-exempt-status-check.md) (Accepted, unveränderbar):
  die Zeile zitiert die Baseline-Vorlage `.d-check.yml` als Messinstanz. Der
  Referent ist unverändert geprüft (der `status`-Block der v6.13.0-Vorlage
  trägt `matrix.status` weiterhin nur klassen-übergreifend) — die Korrektur
  der Version im Pfad ist **Zitat-Korrektur** nach
  [ADR-0073](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md):
  in-place, Commit nennt
  [ADR-0073](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md),
  §Geschichte-Zeile in der ADR.
- `docs/plan/planning/observations/BEO-PGC/…` (3 Zeilen in Evidence-/State-Dateien
  zu `slice-105`): Lauf-Belege mit derselben Record-Einfrierungs-Doktrin wie
  `docs/reviews/**` — aber **nicht exempt**. Entscheidung beim Planner:
  entweder Zitat-Korrektur (ungewöhnlich bei Records) oder Exemption für
  `docs/plan/planning/observations/**` — die eine Gate-Änderung ist und
  nach AGENTS §3.6 einen ADR-Träger braucht. Empfehlung: als offenen Punkt
  in die Closure tragen, nicht still entscheiden.
- Slice-Plan (1 Rest-Pin: die „entfernt"-Zeile der §3-Tabelle): dieselbe
  generische Pfadform wie in `c8e2904f` bereits angewendet.
- `baseline-verify.md`-Zeile: Prosa ohne Pin-Muster (gemessen), nur 4.1.7.

**Lesart zur
[ADR-0095](../plan/adr/0095-review-klasse-exempt-status-check.md)-Spanne
(F-4, Architect-Verdikt).** Der Pin sitzt in §Verglichene Alternativen — nach
AGENTS §3.5 eine der Sektionen, die „unberührbar" bleiben („ihre Änderung ist
eine neue ADR mit `Supersedes`, nie eine Zitat-Korrektur"); nach der
[ADR-0073](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md)-Kurzform
(„Gerüst ja, Aussage nie") ist eine reine Adress-Korrektur zulässig. Meine
Lesart: die Sektionen-Liste schützt die **normativen Aussagen** dieser
Sektionen; eine Änderung, die ausschließlich die Adresse im Zitatgerüst bewegt
und den Referenten unverändert lässt (hier gemessen), fällt unter die
[ADR-0073](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md)-Ausnahme,
die AGENTS §3.5 selbst eröffnet. Nach dieser Lesart ist die Ausführung
(Commit `00d96eb7` mit
[ADR-0073](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md)-Kennung,
`5bb4eabc` mit der §Geschichte-Zeile) in der Sache konform; der Wortlaut von
AGENTS §3.5 („nie") trägt freilich auch die engere Lesart, nach der der Weg nur
über eine Folge-ADR ging. Weil beide Lesarten vertretbar sind und die Klasse
**wiederkehrt** (jede künftige Versions-Angabe in einer Accepted-ADR, die in
einer geschützten Sektion sitzt), braucht die Klärung einen **eigenen Träger**:
Eintrag ins Beobachtungs-Register (`BEO-PGC`), vom Planner bei der Closure zu
öffnen; bis zur Klärung gilt für neue Fälle die engere Lesart (Folge-ADR statt
Zitat-Korrektur).
[ADR-0095](../plan/adr/0095-review-klasse-exempt-status-check.md) selbst bleibt
unberührt.

### 4.3 Inhaltliche Funde — Nachzug oder bewusstes Nicht-Nachziehen

| Fund | Entscheidung | Grund |
|---|---|---|
| `<PREFIX>-RB-<NN>` | Nachzug **nur** in MR-000 (4.1.3); keine Einführung im Lastenheft | Das Lastenheft trägt FA/QA-Reihen; RB ist eine Verfügung, kein Zwang — eine künftige Randbedingung erhält die Reihe über MR-000 und das v6.13.0-Template |
| Modul 4: Gate-Erweiterung ≠ ADR-Anlass | **nicht** nachziehen — als Fund an die offene Beobachtung `BEO-PGC/gate-scope-erweiterung-ohne-adr-traeger` gemeldet | Die Baseline nimmt zur Klasse Stellung („Im Zweifel: ADR"); die Auflösung der Beobachtung ist eigener Vorgang, nicht Teil dieses Slices (§1 „ausdrücklich nicht") |
| Modul 4: „Nachzug ist keine Überschreibung" | **nicht** nachziehen | Deckt sich mit [ADR-0073](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md); keine Divergenz |
| Modul 9: befristete Ausnahme = Carveout, keine Senkung | **nicht** nachziehen | Deckt sich mit AGENTS §3.6 und der Modul-7-Mechanik |
| Guide-Datei-Schnitt (`harness/rules/<name>.md`), `file.max-lines`-Ratchet | **bewusst nicht** übernommen | Die Form ist Wahl und „lohnt erst, sobald die Liste tatsächlich wächst" — AGENTS §3 trägt die Volltextform heute; ein neuer Träger für einen nicht eintretenden Fall wäre Pflichterfüllung ohne Gegenstand |
| Modul 5/8: WIP pro Lauf | **nicht** nachziehen | Gilt über die Baseline direkt; keine verkörperte Regel widerspricht |
| Modul 6: Hard-Rule-Trigger-Audit, Welle-Closure-Erweiterungen | **nicht** nachziehen | Träger ist das vendored Modul; kein Eigen-Text des Repos zitiert die Drei-Klassen-Form |
| Referenz-Richtung: Provenance-Regel (externer CR statt ADR/Slice in der Historie) | **nicht** nachziehen — bereits verkörpert | `spec/lastenheft.md` §7 trägt die verschärfte Form wörtlich („Keine ADR-, Slice-, Carveout- oder Welle-Verweise in dieser Tabelle") |
| `gate.template`: lebende/einfrierende Verweisform | **nicht** nachziehen | Ziel-Form für künftige Sensor-Dateien; bestehende Sensor-Dateien behalten ihre Form |
| `templates/.d-check.yml` (link-policy always, matrix-Erweiterungen) | **nicht** übernehmen | Gate-Vertragsänderung — laut Slice §1 ausdrücklich ausgeschlossen; ggf. Folge-Slice, Aufnahmen nach AGENTS §3.6 ADR-pflichtig |

## 5. Bewertung der §4-Rückführung

Der Slice-Plan führt „zu groß" bei „Nachzügen an mehr als den beiden Dateien
der §3-Tabelle". Das Audit zeigt Nachzüge an **neun** weiteren Dateien
(4.1.3–4.1.7 plus [ADR-0095](../plan/adr/0095-review-klasse-exempt-status-check.md)).
**Empfehlung: Slice halten.** Alle neun sind
mechanische Ein-Zeilen-Bumps derselben Art wie der geplante `AGENTS.md`-§1-Edit,
forciert durch die bereits geplante Entfernung — ohne sie bleibt
`make gates` nach dem Konventions-Update rot (`versions`/`tracked`); inhaltlich
braucht nichts davon einen Nachzug (4.3). Die Rückführung nach `next` wäre
Prozesskosten ohne neuen Entscheidungsinhalt. Die Entscheidung liegt beim
Planner.

## 6. Suchlauf (AGENTS §3.13)

```suchlauf
git grep -o 'v6\.9\.0' <Stand> -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
git grep -o 'Kurs-Welle 137' <Stand> -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
git grep -oE '\.harness/baseline/v[0-9]+\.[0-9]+\.[0-9]+/' <Stand> -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
```

| Muster | Stand `c8e2904f` (Parent) | Stand `d443ee39` (Arbeitsbaum) |
|---|---|---|
| `v6\.9\.0` | 21 Zeilen | 21 Zeilen |
| `Kurs-Welle 137` | 1 Zeile | 1 Zeile |
| Pin-Pfad `v<MAJOR>.<MINOR>.<PATCH>/` | 13 Zeilen (12 Dateien) | 13 Zeilen (12 Dateien) |

Nicht gefunden: jede weitere `Kurs-Welle`-Nennung des alten Stands. Die 21
bzw. 13 Treffer verteilen sich exakt auf die in §4.1/4.2 gelisteten Stellen;
der Bericht selbst (diese Datei) liegt in der ausgeschlossenen Fläche.

## 7. Verbleibende Risiken

1. **`make baseline-verify` war solange rot, bis `v6.9.0/` entfernt ist** —
   das Skript verlangt genau ein `<tag>`-Verzeichnis; der Bundled-Commit
   ließ beide bestehen (Auftrag: Entfernung ist Implementer-Zug). Nach
   `88cea828` ist die Entfernung vollzogen.
2. **`versions`-Befunde der Observations-Evidence** (4.2) sind bis zur
   Closure-Entscheidung offen — andernfalls bleibt `make gates` nach dem
   Konventions-Update rot, auch wenn 4.1 abgearbeitet ist.
3. Die **Spanne zur Zitat-Korrektur an
   [ADR-0095](../plan/adr/0095-review-klasse-exempt-status-check.md)**
   (AGENTS §3.5 Sektionen-Liste vs.
   [ADR-0073](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md)-Kurzform)
   ist in §4.2 mit Lesart dokumentiert; die Klärung braucht einen eigenen
   Träger (Beobachtungs-Eintrag), bis dahin gilt für neue Fälle die engere
   Lesart.
4. **Nicht gemessen, sondern hergeleitet:** die `versions`-Befunds-Erwartung
   in 4.2 folgt aus der Modulkonfiguration; der erste reale Gate-Lauf nach
   dem Konventions-Update ist der Beleg.
