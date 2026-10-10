# Architect-Verdikt — `matrix-inactive` meldet nur einen Link zwischen Dateien einer Klasse (`SPEC-040`, `slice-spec-festlegungen-doku-gates`, F-3)

**Datum:** 2026-10-10 · **Stand:** `7f139836` (Baum sauber) · **Rolle:** Architect (Modul 8),
Konflikt-Pfad, anderer Kontext als Implementer- und Reviewer-Lauf des Slice ·
**Rolleninhaber:** pt9912 (Architect-Agent im Auftrag) · **Anlass:**
[`review-slice-spec-festlegungen-doku-gates`](review-slice-spec-festlegungen-doku-gates.md)
F-3 (HIGH); der Implementer hat angehalten (Slice-Plan §6, Frage zu F-3; §3 Messung M32),
weil `ADR-0095` §Kontext eine Prüfung der Tokens nennt, die das Werkzeug nicht leistet.

**Bezug:** [`ADR-0095`](../plan/adr/0095-review-klasse-exempt-status-check.md)
§Kontext, §Entscheidung, §Konsequenzen, Fitness Function ·
[`ADR-0163`](../plan/adr/0163-schaerft-spec-040-provenance-marker-und-exempt-paths.md)
Entscheidung 1 (Tabelle, Zeile „Status-Prüfung“) und 3 ·
[`SPEC-040`](../../spec/pflichtenheft.md#7-festlegungen-der-harness-werkzeuge) ·
[`AGENTS.md`](../../AGENTS.md) §3.5, §3.12 · d-check v0.82.0 (Digest aus `d-check.mk`).

---

## 0. Ergebnis in einem Blick

- **Verdikt 1 aus Modul 8: die Entscheidung gilt, die Festlegung hat falsch behauptet.**
  Die Zeile `matrix-inactive` in `SPEC-040` geht auf das Werkzeug: ein **Link** aus
  einer Datei einer `matrix`-Klasse auf eine Datei einer `matrix`-Klasse, deren Status
  in `status.forbidden` steht. Ein Token ist kein Treffer; eine Datei ohne Klasse ist
  weder Quelle noch Ziel der Prüfung.
- **Keine Folge-ADR, keine neue ADR, keine Zitat-Korrektur an `ADR-0095`.** Der
  Wortlaut von §Entscheidung deckt die Arbeit (§1). Der Satz über die „ausgehenden
  `ADR-\d{4}`-Token“ steht in §Kontext, beschreibt den Mechanismus ungenau und bleibt
  als Geschichte stehen. Eine Zitat-Korrektur kommt nicht in Frage, weil sich der
  Referent ändern würde (`AGENTS.md` §3.5).
- **Vertrag `harness/sensors/docs-check.md`:** ein Satz in Grenze 12 (§3.2). Ein Leser
  von `AGENTS.md` oder `harness/README.md` würde sonst annehmen, ein Link dort auf eine
  abgelöste ADR fiele im Gate auf.
- **Akzeptiertes Negativ, kein Folge-Slice, kein Register-Eintrag:** Ein Token auf eine
  abgelöste ADR und ein Link aus einer Datei ohne Klasse bleiben grün (§2).

## 1. Begründung

1. **§Entscheidung von `ADR-0095` sagt keine Token-Prüfung zu.** Sie sagt nur, wovon
   `docs/reviews/*.md` ausgenommen wird („ausschließlich von der
   `matrix.status`-Prüfung, nicht von `matrix.rules`“). Ob diese Prüfung Links oder
   Tokens liest, legt sie nicht fest. Wie weit die Ausnahme reicht, hat
   `ADR-0163` Entscheidung 3 schon an das Werkzeug angepasst. §Konsequenzen sagt
   „verweist“, ohne eine Form zu nennen. Die Fitness Function (58 `matrix-inactive`
   ohne die Ausnahme) bleibt wahr, und die 57 Befunde aus M32 (dritter Lauf) hängen
   alle an Links (*übernommen*, Messung des Implementers).
2. **Ein Satz in §Kontext ist keine Entscheidung.** Er erklärt, warum 58 Befunde
   entstanden, und nennt dafür einen Mechanismus, der am gepinnten Werkzeug nicht
   zutrifft (Zeile T1 unten). Für den Befund selbst und für die Wahl von Option C
   trägt er nichts: Die Befunde entstehen, sobald `docs/reviews/*.md` einer Klasse
   angehört, ob über Link oder Token. Eine Folge-ADR, die nur diesen Satz ablöst,
   hätte keine Wirkung. Das ist dieselbe Lage wie bei den zwei Sätzen in
   `ADR-0163` Entscheidung 4.
3. **Kein Widerspruch zu `ADR-0163` Entscheidung 1.** Für die Status-Prüfung nennt die
   Herkunftstabelle `ADR-0095`. Da sie Link oder Token nicht festlegt, gilt der
   Grundsatz der Tabelle, dass die Festlegung das Werkzeug wiedergibt. Der Satz
   „weicht das Werkzeug ab, ist das Werkzeug falsch“ setzt eine Festlegung voraus,
   die selbst richtig ist. Die heutige Zeile war breiter als das Werkzeug, seit
   sie im ungemergten Slice steht. Die Engführung lockert kein Gate (`AGENTS.md`
   §3.6): Das Werkzeug meldet vorher wie nachher dasselbe.

## 2. Akzeptiertes Negativ

Die Prüfung meldet ein Token einer abgelösten ADR nicht und prüft keine Datei ohne
Klasse. Die Lücke bekommt keinen Sensor und keinen Folge-Slice:

- Im Bestand stehen Tokens abgelöster ADRs als legitime Herkunft, in
  Supersession-Ketten, in Belegen und Plänen. Grobmaß am Stand `7f139836`: 5 ADRs
  tragen den Status `Superseded`/`Deprecated` (0018, 0036, 0038, 0039, 0151). Ihre
  Kennungen stehen in 122 Zeilen in 27 Dateien unter `docs/plan/adr/` sowie in
  `slice-*.md` und `observations/` (*gemessen*, Befehl unten; das Maß schließt
  Lineage-Zeilen ein, die erlaubt wären). Eine Token-Prüfung würde vor allem
  Herkunft melden und bräuchte einen Marker-Ausweg, wie ihn `ADR-0163`
  Entscheidung 2 für `matrix-forbidden` hat.
- Der Link ist der Teil, der beim Lesen bricht. Ihn prüft das Gate in allen sechs
  Klassen (Zeilen L1, L4 bis L7; `review` per M32).
- Für `AGENTS.md`, `harness/**` und `README.md` bleibt das Review der Wächter, wie bei
  jeder Prosa-Aussage dort.

```text
git grep -l -iE '^\*\*Status:\*\* *(Superseded|Deprecated)' -- 'docs/plan/adr/[0-9]*.md'
git grep -nE '\b(ADR-0018|ADR-0036|ADR-0038|ADR-0039|ADR-0151)\b' -- 'docs/plan/adr/[0-9]*.md' 'docs/plan/planning/**/slice-*.md' 'docs/plan/planning/observations/**/*.md' | wc -l
git grep -lE '\b(ADR-0018|ADR-0036|ADR-0038|ADR-0039|ADR-0151)\b' -- 'docs/plan/adr/[0-9]*.md' 'docs/plan/planning/**/slice-*.md' 'docs/plan/planning/observations/**/*.md' | wc -l
```

Re-Evaluierung, falls ein Review einen Link aus einer Datei ohne Klasse auf eine
abgelöste ADR als echten Fehler findet: Dann kommt eine Klasse für diese Dateien in
Frage, per ADR, weil `.d-check.yml` sich ändert.

## 3. Auftrag an den Implementer — Wortlaut zum Übernehmen

### 3.1 `spec/pflichtenheft.md`, §7, `SPEC-040`

**Zeile `matrix-inactive` der Tabelle „Befund je Modul“**, Spalte „gilt als Treffer“:

> ein Link aus einer Datei einer Klasse zeigt auf eine Datei einer Klasse, deren
> Status in `status.forbidden` steht (`superseded`, `deprecated`)

**Absatz „Zu `SPEC-040` — Randformen der Verweis-Module“:** Der Satz „`matrix`
meldet auch eine Kennung der Zielklasse in Inline-Code, in einem Fence nicht.“ wird
ersetzt durch:

> Für `matrix-forbidden` meldet `matrix` auch eine Kennung der Zielklasse in
> Inline-Code, in einem Fence nicht. `matrix-inactive` meldet nur einen Link; eine
> Kennung, deren Datei einen verbotenen Status trägt, ist kein Treffer. Eine Datei,
> die keiner Klasse angehört (etwa `AGENTS.md` oder `harness/README.md`), prüft das
> Modul weder als Quelle noch als Ziel.

### 3.2 `harness/sensors/docs-check.md`, §Grenze, Punkt 12

Überschrift des Punkts: **`matrix` — zwei angenommene Lücken, der Fence und der
Status nur am Link.** Nach dem Satz „Ein Token in einem Fence ist für `matrix` kein
Treffer.“ steht neu:

> `matrix-inactive` meldet nur einen Link zwischen zwei Dateien, die einer Klasse
> angehören. Eine Kennung einer abgelösten ADR im Text und ein Link aus `AGENTS.md`,
> `README.md` oder `harness/**` auf eine abgelöste ADR bleiben grün; der Wächter ist
> das Review.

Der Satz nach „Die Wirkung steht in“ bleibt; ergänzt wird nach „Entscheidung 2 und 3“:
„, die Lesung von `matrix-inactive` im
[Architect-Verdikt](../../docs/reviews/architect-verdict-matrix-inactive-nur-link.md)“.

### 3.3 Slice-Plan

§6, Frage zu F-3: Ausgang *entfallen*, Grund: dieses Verdikt. §3: Gegenprobe zur
Status-Prüfung mit Ausgang S und Verweis auf dieses Verdikt. Die Messungen unten
tragen die Zeile; M32 bleibt der Beleg des Implementers. §8 (Historie) der Spec:
keine neue Zeile, die bestehende Zeile vom selben Datum trägt die Änderung (Review
F-11).

## 4. Messung

Gemessen von diesem Architect-Zug an einem Klon von `7f139836` im Scratchpad, d-check
v0.82.0 über `docker run --rm --network none -v <Klon>:/repo:ro <Digest>` (Digest aus
`d-check.mk`). Basislauf: `1832 Datei(en) geprüft, 0 Befund(e)`, Exit 0. Jede
Mutation ist **eine** Zeile, angehängt an **eine** Datei; danach
`git checkout -- . && git clean -fd`. Instanz ist der d-check-Lauf über den Klon,
kein Tabellentest. Ziel ist, wo nicht anders genannt, `ADR-0038` (Status
`Superseded by ADR-0039`).

| # | Mutation | Ergebnis |
|---|---|---|
| T1 | Token `ADR-0038` in Inline-Code in `ADR-0094` (Klasse `adr`) | `0 Befund(e)`, Exit 0 |
| T2 | dasselbe im Slice-Plan in `in-progress/` (Klasse `slice`) | `0 Befund(e)`, Exit 0 |
| L1 | Link auf `ADR-0038` in `ADR-0094` | 1 × `matrix-inactive`, Exit 1 |
| L2 | Link in `harness/README.md` (keine Klasse) | `0 Befund(e)`, Exit 0 |
| L3 | Link in `AGENTS.md` (keine Klasse) | `0 Befund(e)`, Exit 0 |
| L4 | Link im Slice-Plan in `in-progress/` | 1 × `matrix-inactive`, Exit 1 |
| L5 | Link in einer `observation.md` (`BEO-PGC/adapter-fehler-ausgang`) | 1 × `matrix-inactive`, Exit 1 |
| L6 | Link in `done/welle-12.md` (Klasse `welle`) | 1 × `matrix-inactive`, Exit 1 |
| L7 | Link in `spec/architecture.md` (Klasse `spec`) | `matrix-forbidden` und `matrix-inactive`, Exit 1 |
| Z1 | neue, in den Index aufgenommene Datei `docs/plan/probe-ohne-klasse.md` (keine Klasse) mit `**Status:** Superseded`, Link darauf aus `ADR-0094` | `0 Befund(e)`, Exit 0 |
| Z2 | `**Status:** Superseded` in `done/welle-12.md` (Klasse `welle`) eingefügt, Link darauf aus `ADR-0094` | 3 × `matrix-inactive`, Exit 1: aus `ADR-0094` und aus den zwei Slice-Plänen in `done/`, die schon auf die Datei verweisen |

**Reichweite.** Als Quelle sind fünf Klassen gemessen (`adr`, `slice`, `observation`,
`welle`, `spec`), die sechste (`review`) ist M32 des Implementers (*übernommen*). Als
Ziel sind `adr` und `welle` gemessen, für `spec`, `slice`, `review` und `observation`
ist das Verhalten *hergeleitet*. Dass ein Token kein `matrix-inactive` auslöst, ist
an `adr` und `slice` gemessen (T1, T2), an `observation` und `review` durch M32
(*übernommen*), an `welle` und `spec` *hergeleitet*. Gesehen ist, dass `Superseded`
den Eintrag `superseded` trifft und `Superseded by […]` ebenso (Z2, L1). Welche Regel
dahintersteht (Großschreibung, erstes Wort oder Präfix), ist *hergeleitet*; die
Festlegung sagt dazu nichts.
