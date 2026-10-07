# Harness-Konventionen

---

## Purpose

Diese Datei deklariert die *repo-lokalen* Strukturregeln dieses Repos
gegenüber der adoptierten Harnesskonvention (Baseline). Sie ist der
Default-Ort für:

- **Adaptionen** ggü. der Baseline (mit Begründung und Auflösungs-Trigger).
- **ID-Schema-Deklaration** — welches Präfix-Schema dieses Repo nutzt.
  Der Baseline-Default wird als Teil der `MR-000`-Aussage festgehalten;
  ein abweichendes Präfix oder Schema ist ein eigener `MR`-Eintrag.
- **Zusatzklassen-Deklarationen** für repo-spezifische
  Bindung-Klassen in der Sensors-Tabelle, die über die vier kanonischen
  hinausgehen (ADR, Carveout, Schwelle, Reproduzierbarkeit).
- **Modus-Deklarationen** pro Sub-Area (Greenfield / Brownfield /
  Hybrid) inklusive Konvergenz-Auftrag bei BF.

Bei Konflikt zwischen dieser Datei und einer kanonischen Quelle gilt die
kanonische Quelle (Source Precedence). Diese Datei ist konformitäts-
bringend für *Form*-Fragen, nicht autoritativ über Inhalt.

## Baseline

<!--
Welche Harnesskonvention wird adoptiert? Stand und Datum festhalten,
damit spätere Adaptionen einen Bezugspunkt haben.
-->

- **Konvention:** AI-Harness-Kurs (Baseline-Regelwerk, als Release-Asset
  vendored unter `.harness/baseline/`)
- **Stand:** v6.16.0
- **Datum der Adoption:** 2026-10-06

<!--
Der Stand ist eine VERSION, kein Datum: Er ist der Bezugspunkt, gegen den ein
Versions-Sensor die Baseline-Pins in den Adaptions-Eintraegen prueft (Muster in
`.d-check.yml`). Steht hier ein Datum, findet der Sensor keine Version und
bricht fail-closed ab — der Gate laeuft dann gar nicht mehr. Das Datum der
Adoption steht in der Zeile darunter.
-->

## Adoptierte Konventions-Quellen

<!--
Pointer auf die Quellen der Baseline. KEINE Wiederholung des Inhalts —
nur Verweise.

Das Agenten-Regelwerk ist die Quelle, die ein Code-Agent statt des
vollen Lehrmaterials liest (operatives Regelwerk ohne Didaktik). Es
ist derivativ — bei Konflikt gilt das Lehrmaterial.
-->

- **Extern (Lehrmaterial):** https://github.com/pt9912/ai-harness-course
  (das Kurs-Repo, dessen Release-Asset unten vendored ist; derselbe Ort, gegen
  den `make pin-stale-baseline` den adoptierten Stand prüft)
- **Vendored Baseline (Regelwerk + Templates):** aus dem self-contained
  Release-Asset
  https://github.com/pt9912/ai-harness-course/releases/download/v6.16.0/lab-regelwerk.zip
  nach `.harness/baseline/<tag>/{regelwerk,templates}/` entpackt (netzlos,
  `SHA256SUMS`) — adoptierter Stand: Kurs-Welle 159 · 2026-10-06 (Stand-Zeile
  in `regelwerk/README.md`; Wellen-Register: CHANGELOG.md im Kurs-Repo); für
  harte Reproduzierbarkeit das Asset eines Tags ziehen statt `latest`.
- **In-Repo (verkörperte Form):** `AGENTS.md`, `harness/README.md`, diese
  Datei, `docs/plan/planning/README.md`,
  `docs/plan/planning/in-progress/roadmap.md`, `docs/plan/adr/README.md` und
  `docs/plan/carveouts/README.md` (die Liste der abgeleiteten Dokumente im
  Bump-Ablauf von `harness/targets/pin-stale.md`) — die vendored
  `.harness/baseline/<tag>/templates/` sind die Referenz-Form („Ziel-Form" des
  Regelwerks); diese Dateien sind daraus kopiert und ausgefüllt.

## Adaptions-Block

Regeln dieser Sektion: Diese Datei trägt den **Index**, nicht die Einträge.
Jede Adaption ist eine eigene Datei unter `harness/conventions/`, kopiert aus der
gleichnamigen Eintrags-Vorlage `MR-NNN-titel.template.md` der vendored Baseline;
ist ihr Auflösungs-Trigger eingetreten, wandert sie per `git mv` nach `conventions/done/` — in einem eigenen Commit, der nur diesen Umzug trägt (dieser Index und der Nachfolge-Eintrag stehen in anderen Commits); die Message nennt die aufgelöste Kennung und [`ADR-0162`](../docs/plan/adr/0162-schaerft-spec-038-039-und-umzug-aufgeloester-mr-eintraege.md), weil `make doc-immutable` um diesen Commit in Teil-Ranges läuft (Entscheidung 3). Der Zustand ist die Verzeichnis-Position, kein
Status-Feld. Der Grund für den Schnitt: Was hier steht, liest **jeder**
Agentenlauf — aufgelöste Adaptionen gehören nicht in diesen Pfad
(Baseline-Regelwerk `grundlagen-harness-dateien.md`
§harness/conventions.md als Konventionsspeicher).

### MR-000 — Baseline-Aussage

Bleibt hier: Sie ist keine Adaption, sondern die Adoptions-Erklärung, und
sie gilt für jeden Lauf.

- **Datum:** 2026-09-09
- **Geltungsbereich:** gesamtes Repo
- **Ersetzt-Baseline-Regel:** — *(keine; dieser Eintrag ist die
  Adoptions-Erklärung, keine Adaption)*
- **Adaption:** *keine inhaltlichen Adaptionen ggü. Baseline-Default
  für Verzeichniskonvention, Lifecycle-Regeln, Carveout-Disziplin,
  ID-Schema (`<PREFIX>-FA-*`, `<PREFIX>-QA-*`, `<PREFIX>-RB-*`, `SPEC-<NNN>`, `ARC-<NNN>`,
  `ADR-<NNNN>`, `CO-<NNN>`, `slice-<Kennung>`, `MR-<NNN>`, `BEO-<KUERZEL>/<slug>`, `RC-<NNN>` — nur das
  Vertrags-Präfix wird repo-weit festgelegt, z. B. `LH`; `SPEC-*` und
  `ARC-*` kodieren das Stratum und sind fest, siehe Baseline-Regelwerk
  `grundlagen-source-precedence.md` §ID-Schema als Klammer;
  bei mehreren gleichzeitig schreibenden Entwicklern für die Artefakte mit
  je eigener Datei (`ADR-*`, `CO-*`, `slice-*`, `welle-*`) zusätzlich das
  Bereichssegment und damit den Zählraum je Sub-Area festlegen —
  `SPEC-*`/`ARC-*` bleiben davon ausgenommen und zählen fortlaufend je
  Datei, siehe Baseline-Regelwerk `grundlagen-source-precedence.md` §Vergabe).*
- **Begründung:** Initial-Setzung. Spätere Adaptionen werden als
  `MR-<NNN>` nachgetragen.
- **Auflösungs-Trigger:** permanent.

### Aktive Adaptionen

<!-- Eine Zeile je Datei in harness/conventions/. Geltungsbereich und
     Ersetzt-Baseline-Regel stehen hier, damit ein Agent ohne Öffnen
     entscheiden kann, ob der Eintrag ihn betrifft.

     Das <a id="mr-<NNN>"> in der MR-Zelle ist die Adresse, unter der andere
     Dateien diese Adaption referenzieren: conventions.md#mr-<NNN>. Es steht
     hier und nicht in der Eintrags-Datei, weil die Datei bei Auflösung nach
     conventions/done/ wandert und ein Pfad-Link dabei bricht — die Zeile
     wechselt nur die Tabelle, der Anker reist mit. Der Anker trägt die
     Kennung, nie den Titel: Titel werden umformuliert, Kennungen nicht.

     Kam dieses Repo von der Inline-Form (### MR-NNN — Titel), trägt die
     Zeile den alten Überschriften-Slug als ZWEITEN Anker daneben, sonst
     rotten die bereits veröffentlichten Verweise. -->

| MR | Titel | Geltungsbereich | Ersetzt-Baseline-Regel |
|---|---|---|---|
| MR-002 <a id="mr-002"></a> | [Slice-/Welle-Kennungen sind Namen, nicht Nummern (ab slice-105 exklusive)](conventions/MR-002-slice-welle-kennungen-sind-namen.md) | `harness/conventions.md` §Aktive Adaptionen; alle nach `slice-105` neu angelegten Slice-/Welle-Plan-Dateien | `grundlagen-source-precedence.md` §Vergabe (Bestandsschutz für `slice-001`–`slice-105`) |
| MR-003 <a id="mr-003"></a> | [Der PreToolUse-Guard blockt in-place Textwerkzeuge und Host-Interpreter auf Repo-Pfaden](conventions/MR-003-guard-inplace-textwerkzeug.md) | `.claude/hooks/pretooluse-command-guard.sh`, `tools/harness/mask-quotes.awk`, `tools/harness/run-command-guard-tests.sh`, `AGENTS.md` §3.1 „Durchsetzung“ | `grundlagen-durchsetzungsschicht.md` §Grenzen — ehrlich benannt |
| MR-004 <a id="mr-004"></a> | [Der PreToolUse-Guard sperrt Host-`python` und `python3` am Kopf eines Kommando-Segments unbedingt](conventions/MR-004-guard-host-python-am-kopf.md) | `tools/harness/blocked/python`, `.claude/hooks/pretooluse-command-guard.sh`, `tools/harness/run-command-guard-tests.sh`, `AGENTS.md` §3.1 „Durchsetzung“ | `modul-13-quality-gates.md` §Guard-Härtung |
| MR-005 <a id="mr-005"></a> | [Baseline-Pins in ADRs und MR-Einträgen bleiben auf dem Stand ihrer Abfassung](conventions/MR-005-baseline-pins-in-adrs-und-mr-eintraegen-eingefroren.md) | `.d-check.yml` `versions:`/`vcs:`, Pins in `docs/plan/adr/[0-9]*.md` und `harness/conventions/MR-[0-9]*.md`, `harness/targets/pin-stale.md` §Bump-Ablauf | `grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln („Zwei Rot-Quellen, ein Prinzip“) |
| MR-006 <a id="mr-006"></a> | [Technik-Dokument heißt Pflichtenheft](conventions/MR-006-technik-dokument-heisst-pflichtenheft.md) | `spec/pflichtenheft.md` (vormals `spezifikation.md`, Abschnitte 1–8), `harness/README.md`, `AGENTS.md`, `spec/lastenheft.md` | `grundlagen-referenz-richtung.md` §Spec-Straten (Rang-2-Datei-Name) |

### Aufgelöste Adaptionen

<!-- Eine Zeile je Datei in harness/conventions/done/ — nur ID und
     Nachfolger, damit die Kette auffindbar bleibt, ohne gelesen zu werden.

     Der Anker der Zeile zieht aus der Tabelle oben mit um; er ist der Grund,
     warum ein Verweis auf eine aufgelöste Adaption nicht bricht. -->

| MR | aufgelöst durch |
|---|---|
| [MR-001](conventions/done/MR-001-technik-dokument-heisst-pflichtenheft.md) <a id="mr-001"></a> | [MR-006](#mr-006) |

## Zusatzklassen-Deklaration für Sensors-Bindung

<!--
Die vier kanonischen Bindung-Klassen der Sensors-Tabelle in
`harness/README.md` (ADR, Carveout, Schwelle, Reproduzierbarkeit) sind
ohne Deklaration legitim.

Repos können weitere Klassen einführen — z. B. Anforderungs-Bindung
(`LH-...`), Compliance-Bindung (Regulatorik-Artikel), Modell-Version-
Bindung (für KI-Evals). Diese müssen hier deklariert werden, sonst sind
sie für Reviewer nicht von Tippfehlern unterscheidbar.

Eine nicht-deklarierte Zusatzklasse in der Sensors-Tabelle ist eine
stille Setzung und damit Harness-Lüge in derselben Klasse wie ein
halluziniertes Gate (Modul 13).
-->

| Klasse | Form | Bedeutung | Beispiel |
|---|---|---|---|
| LH-Bindung | `LH-(FA\|QA)-<BER>-<NNN>` als Link auf `spec/lastenheft.md` | die Zeile trägt oder belegt eine Anforderung des Lastenhefts | `LH-QA-POR-003` bei `make test-integration` |
| Vertragsdatei | Link auf `harness/sensors/<target>.md` bzw. `harness/targets/<gruppe>.md` | Vertrag, Grenze und Ausgänge stehen dort; die Index-Zeile trägt einen Satz | `harness/sensors/fmt-check.md` bei `make fmt-check` |
| Herkunfts-Anker | `seit slice-<Kennung>` / `seit welle-<Kennung>`, ergänzt um `erweitert seit …` | der Vorgang, der die Zeile eingeführt oder erweitert hat | `seit slice-harness-fmt-check` bei `make fmt-check` |
| Beobachtungs-Bindung | `BEO-PGC/<slug>` | der Register-Eintrag, aus dem das Werkzeug hervorging | `BEO-PGC/formatierungs-drift-ohne-gate` bei `make fmt-check` |
| MR-Bindung | `MR-<NNN>` als Link auf `harness/conventions/` | die Adaption, die die Regel des Werkzeugs trägt | `MR-003` bei `make test-command-guard` |
| Hard-Rule-Bindung | Link auf `AGENTS.md` mit Abschnitt `§<n>` | die harte Regel, deren Form das Werkzeug liest | `AGENTS.md` §3.7 bei `make kommentar-kennungen` |

## Modus-Deklaration pro Sub-Area

Die **Kürzel**-Spalte tragen nur Repos, deren Kennungen ein Bereichssegment
führen (`ADR-<KUERZEL>-NNNN`, `slice-<KUERZEL>-NNN`); wer ohne Segment zählt,
streicht sie. Regel: Baseline-Regelwerk `grundlagen-harness-dateien.md`
§Konventionsspeicher.

<!--
Pro Modul / Verzeichnis / Sub-Area: Modus festlegen.
- Greenfield (GF): Doc führt, Code folgt. Steady-State.
- Brownfield (BF): Code führt, Doc folgt. Übergangsmodus mit
  Konvergenz-Auftrag zu GF. Graduation-Bedingung benennen.
- Hybrid: gemischt pro Sub-Sub-Area.
- Permanent-BF (selten): nur für Code, der absehbar entfernt wird;
  mit Begründung und Folge-Slice analog zu permanentem Carveout.

Eine Sub-Area in BF *ohne* Graduation-Plan ist eine Harness-Lüge:
"permanente Ausnahme als temporär getarnt" (Modul 7 Analogie).

Kuerzel: kurz, GROSS, ohne Leerzeichen — und nach der ersten vergebenen
Kennung unveraenderlich. Die Bedingung, wann die Spalte ueberhaupt gefuehrt
wird, steht oben im Fliesstext, weil dieser Kommentar beim Kopieren wegfaellt.
-->

| Sub-Area (Pfad / Modul) | Kürzel | Modus | Begründung | Graduation-Bedingung / Folge-Slice |
|---|---|---|---|---|
| `*` (Default für gesamtes Repo) | `PGC` | Greenfield | Repo startet ohne Code-Bestand; die drei Spec-Straten (Lastenheft 0.3.0, Pflichtenheft, Architektur-Sicht) sind committet, bevor der erste Code entsteht — Doc führt | n/a (GF) |

## Glossar (optional)

<!--
Repo-spezifische Begriffe, die in den Kernbegriffen des
Baseline-Regelwerks nicht stehen. Nur ergänzen, nicht wiederholen.
-->

— keine Einträge: Die Begriffe des Produkts führt
[`spec/lastenheft.md`](../spec/lastenheft.md) §6 Glossar, die des Harness das
Baseline-Regelwerk; dieser Abschnitt wiederholt sie nicht.
