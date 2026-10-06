# Review-Report: slice-zitat-korrektur-vergleichseinheit — 2026-10-06

**Review-Art:** Code — gegen Plan, Entscheidungen und Hard Rules (nicht gegen die DoD; die prüft der Verifier)

**Gegenstand:** Diff `b7d97cca..e0c82845` (3 Commits): `07923146`
([`ADR-0158`](../plan/adr/0158-zitat-korrektur-vergleichseinheit-je-verweisform.md)
samt ADR-Index, Architect), `fc6d104c` (`AGENTS.md` §3.5, Implementer),
`e0c82845` (Slice-Plan, Belege und Suchlauf, Implementer).

**Skill:** `.harness/skills/reviewer.md` @ e0c82845
**Modell:** claude-opus-5-5 · **Datum:** 2026-10-06

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis; die `<Platzhalter>` darin sind Formbeispiele)*. Dieser
> Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link
> (`v<X.Y.Z>` · `regelwerk/<datei>.md` §<Abschnitt>). Der vendored Baum trägt
> genau einen Tag; der Sprung löscht den alten, und ein Link darauf färbt beim
> nächsten Bump ein Artefakt rot, das niemand mehr anfassen darf. Ein `pfad`-Feld
> auf den **geprüften Gegenstand** ist davon nicht betroffen — es zitiert den
> Stand des Laufs und darf ihn festhalten (`v<X.Y.Z>` ·
> `regelwerk/grundlagen-harness-dateien.md` §harness/README.md als
> Einstiegspunkt — diese Zeile ist selbst ein Beispiel der Form).

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde — ohne
diese Liste ist der Lauf nicht reproduzierbar):

- Slice-Plan `slice-zitat-korrektur-vergleichseinheit` (Stand `e0c82845`)
- [`ADR-0158`](../plan/adr/0158-zitat-korrektur-vergleichseinheit-je-verweisform.md) (neu, Supersedes [`ADR-0157`](../plan/adr/0157-zitat-korrektur-reichweite-nach-aussage-und-mr-pins.md) teilweise)
- [`ADR-0157`](../plan/adr/0157-zitat-korrektur-reichweite-nach-aussage-und-mr-pins.md) (Entscheidung 1, 3, 4), [`ADR-0073`](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md) (Entscheidung 3, 4), [`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md)
- Vorheriges Review am selben Gegenstand: [`review-slice-harness-baseline-v6-14-1-fixrunde`](review-slice-harness-baseline-v6-14-1-fixrunde.md) F-1
- `AGENTS.md` §3.1, §3.5, §3.9, §3.12, §3.13
- `harness/conventions/MR-001-technik-dokument-heisst-pflichtenheft.md`, `MR-004-guard-host-python-am-kopf.md` (reale Verweise mit Anker); `harness/targets/pin-stale.md`, `.claude/agents/verifier.md`
- Keine `LH-*`-ID berührt (Harness-Regel)

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | HIGH | Für einen Anker als HTML-`id` ist die Einheit nur die Zeile mit der `id`; vor einem Heading ist das eine Markierungszeile ohne Inhalt, die an jedem Stand gleich bleibt. Gemessen: Abschnitt hinter `<a id="marke">` geändert, die Zeile der `id` `cmp 0` (Gegenprobe Abschnitt `cmp 1`); realer Fall `MR-004` → `modul-13-quality-gates.md#guard-haertung`, beim nächsten Bump besteht er unabhängig vom Inhalt von „Guard-Härtung“. Die Konsequenz „misst jetzt den Referenten jeder Form aus (a)“ ist breiter als die Messung (der HTML-`id`-Fall ist dort nur hergeleitet). | `AGENTS.md` §3.12 „Verfasser einer ADR“; [`ADR-0158`](../plan/adr/0158-zitat-korrektur-vergleichseinheit-je-verweisform.md) Entscheidung 1, §Konsequenzen | `docs/plan/adr/0158-zitat-korrektur-vergleichseinheit-je-verweisform.md` · „Anker als HTML-`id` außerhalb eines Headings \| die Zeile, die die `id` trägt“ | ja — Probe P4 unten (Befehlsform Entscheidung 6) | ADR-Aussage breiter als ihre Messung |
| F-2 | HIGH | Die Konsequenz „MR-Pins `5d8855d9` sind Versions-Segmente auf ganze Dateien … dort ist Einheit gleich Datei, und der … nach `ADR-0157` Entscheidung 4 normalisierte `cmp` ist der hier verlangte“ trägt nicht: `5d8855d9` bewegt vier Verweise **mit Anker**; nach Entscheidung 1 ist die Einheit dort der Abschnitt. Gemessen: zwei Abschnitte `cmp 0`, zwei enden mit Exit 2 (`#guard-haertung` ist eine HTML-`id`; `#spec-straten-mehr-als-ein-spec-dokument` löst in `grundlagen-source-precedence.md` nicht auf). Der `cmp` aus `ADR-0157` Entscheidung 4 vergleicht die MR-Datei selbst, nicht den Referenten. Dieselbe Aussage steht im Plan §1 („an ihnen trägt der Datei-`cmp`“) und begründet dort die Einstufung von `harness/targets/pin-stale.md` und `.claude/agents/verifier.md` als „nicht geändert“. | `AGENTS.md` §3.12; [`ADR-0157`](../plan/adr/0157-zitat-korrektur-reichweite-nach-aussage-und-mr-pins.md) Entscheidung 3, 4 | `docs/plan/adr/0158-zitat-korrektur-vergleichseinheit-je-verweisform.md` · „sind Versions-Segmente auf ganze Dateien und“ | ja — Probe P6 unten | Beleg trägt seinen Satz nicht |
| F-3 | MEDIUM | Die Befehlsform expandiert die Normalisierung ungequotet (`$n`); unter `shopt -s nullglob` oder `failglob` scheitert `sed`, beide Prozess-Substitutionen sind leer, und `vergleich … norm` druckt `cmp 0` mit Exit 0 — gemessen auch für zwei völlig verschiedene Einheiten (`probe/p.md#kind-a` gegen `ADR-0100`; ohne die Option `cmp 1`). Die Zusage „eine leere oder nicht lesbare Einheit endet mit 2“ und die Konsequenz „die Abweichung macht den Lauf rot, nicht grün“ gelten für diesen Pfad nicht. | [`ADR-0158`](../plan/adr/0158-zitat-korrektur-vergleichseinheit-je-verweisform.md) Entscheidung 6 | `docs/plan/adr/0158-zitat-korrektur-vergleichseinheit-je-verweisform.md` · `cmp <(printf '%s\n' "$x" \| $n)` | ja — Probe P5 unten | Messbefehl fail-open |
| F-4 | MEDIUM | Der neue Satz „Der Beleg einer Zitat-Korrektur nennt je Verweis Form, Einheit, beide Stände …“ steht unmittelbar nach „Records … bei ihnen ist die Commit-Kennung der Beleg“ und ohne den Geltungsbereich der Quelle (`ADR-0158` Entscheidung 5: „nach `ADR-0157` Entscheidung 1“; Entscheidung 4: „nicht messbar“ mit Grund). Eine Hostpfad-Korrektur in einem Record liest danach zwei widersprechende Belegformen. Daneben nennt §3.5 als Einheit „den Abschnitt hinter dem Anker“, die ADR für eine HTML-`id` die Zeile — zwei Lesarten derselben Regel. | `AGENTS.md` §3.5; [`ADR-0158`](../plan/adr/0158-zitat-korrektur-vergleichseinheit-je-verweisform.md) Entscheidung 1, 4, 5 | `AGENTS.md` · „Der Beleg einer Zitat-Korrektur“ | nein — Lese-Handlung | Nachzug widerspricht dem Nachbarn im selben Träger |
| F-5 | LOW | Unter `set -e` bricht ein roter `vergleich` am `cmp` ab, bevor die Zeile „vergleich …: … cmp 1“ gedruckt ist (gemessen: nur die `cmp`-Meldung, Exit 1); der Exit bleibt rot, die von Entscheidung 5 verlangte gedruckte Zeile fehlt für genau diesen Fall. | [`ADR-0158`](../plan/adr/0158-zitat-korrektur-vergleichseinheit-je-verweisform.md) Entscheidung 5, 6 | `docs/plan/adr/0158-zitat-korrektur-vergleichseinheit-je-verweisform.md` · `local ec=$?` | ja — Probe P7 unten | Befehlsform-Randverhalten unbenannt |
| F-6 | LOW | „Verglichen wird roh … Keine Normalisierung von Leerraum, Zeilenenden“: die Befehlsform schneidet über `$(…)` abschließende Leerzeilen ab — `x\n` gegen `x\n\n\n\n` `cmp 0`, der rohe Blob-`cmp` derselben Dateien `cmp 1`; ein Zeilen-Lokator auf eine Leerzeile endet als „leere Einheit“ mit 2. | [`ADR-0158`](../plan/adr/0158-zitat-korrektur-vergleichseinheit-je-verweisform.md) Entscheidung 2, 6 | `docs/plan/adr/0158-zitat-korrektur-vergleichseinheit-je-verweisform.md` · „Keine Normalisierung von Leerraum, Zeilenenden“ | ja — Probe P8 unten | Befehlsform-Randverhalten unbenannt |
| F-7 | LOW | Die Normalisierung ersetzt jedes `/v<X.Y.Z>/` im Text der Einheit, nicht nur den bewegten Tag; ein Wechsel eines fremden Werkzeug-Pins (`…/download/v0.79.0/` → `v0.80.0`) besteht normalisiert (`cmp 0`, roh `cmp 1`). Heute ohne Fundort: im Baum `.harness/baseline/v6.14.1` trägt kein Text ein anderes `/v<X.Y.Z>/` als den eigenen Tag (gemessen, 0 Treffer). Die Grenze ist in der ADR nicht benannt. | [`ADR-0158`](../plan/adr/0158-zitat-korrektur-vergleichseinheit-je-verweisform.md) Entscheidung 2 | `docs/plan/adr/0158-zitat-korrektur-vergleichseinheit-je-verweisform.md` · „auch auf Upstream-URLs im Text der Einheit“ | ja — Probe P5 unten (ohne Shell-Option) | Grenze der Normalisierung unbenannt |
| F-8 | INFO | Weil die Heading-Zeile nicht zur Einheit gehört, besteht ein Anker-Wechsel zwischen zwei Abschnitten mit gleichem Körper (`#pflicht` → `#verbot`, Körper je „Siehe Tabelle unten.“: `cmp 0`); die Konsequenz „Ein Anker-Wechsel auf einen anderen Abschnitt fällt“ steht ohne diese Bedingung. Im Baum (ohne Records) gibt es heute keinen solchen Abschnitt (gemessen, 0). | [`ADR-0158`](../plan/adr/0158-zitat-korrektur-vergleichseinheit-je-verweisform.md) §Konsequenzen | `docs/plan/adr/0158-zitat-korrektur-vergleichseinheit-je-verweisform.md` · „Ein Anker-Wechsel auf einen anderen Abschnitt fällt“ | ja — Probe P3 unten | ADR-Aussage breiter als ihre Messung |
| F-9 | INFO | Außerhalb des Diffs (Meldung an den Planner, Träger-Nachzug nach `AGENTS.md` §3.13): `MR-001` verweist auf `grundlagen-source-precedence.md#spec-straten-mehr-als-ein-spec-dokument`; das Heading „Spec-Straten: mehr als ein Spec-Dokument“ steht in `regelwerk/grundlagen-referenz-richtung.md`. `make docs-check` meldet es nicht; die Befehlsform aus Entscheidung 6 findet es (Exit 2). | [`ADR-0157`](../plan/adr/0157-zitat-korrektur-reichweite-nach-aussage-und-mr-pins.md) Entscheidung 3 | `harness/conventions/MR-001-technik-dokument-heisst-pflichtenheft.md` · „grundlagen-source-precedence.md#spec-straten-mehr-als-ein-spec-dokument“ | ja — Probe P6 unten | Verweis löst nicht auf |

## Proben

Alle an einem Klon im Scratchpad (Stand `e0c82845`, GNU Awk 5.2.1,
`LANG=de_DE.UTF-8`). Die Befehlsform ist per `awk` wörtlich aus dem
`bash`-Block von `ADR-0158` Entscheidung 6 gezogen (31 Zeilen, eingerückte
drei Leerzeichen entfernt) und per `source` geladen. Probe-Commits existieren
nur im Klon: `514405ba` (Probe-Datei `probe/p.md` v1) und `f27b4e03` (v2: ein
Zeichen im Körper „Kind A“, ein Wort im Abschnitt hinter `<a id="marke">`,
`v0.79.0`→`v0.80.0` und `v6.14.0`→`v6.14.1` im Abschnitt „Pins“). Gedruckte
Zeilen gekürzt auf die Kernaussage.

- **P0 Nachgefahren (Architect):** Anker-Wechsel `#teilfrage-2-…` → `#teilfrage-3-…` an `ADR-0100` `cmp 1`; gleicher Anker `cmp 0`; → `#entscheidung` `cmp 1`; unbekannter Anker `einheit: leere Einheit`, Exit 2. Versions-Segment `11a5bac5`/`b7d97cca`: `.d-check.yml` roh `cmp 0`, `modul-05…#ziel-form-slice` roh `cmp 0`, `grundlagen-begriffe.md` roh `cmp 1`, norm `cmp 0`. Stimmt mit der §Fitness Function überein.
- **P1 Heading in Code-Fence (vom Architect nicht gefahren):** `einheit 514405ba probe/p.md '#kind-a'` liefert den Körper einschließlich des Fence mit der Zeile `## kein Heading, nur Kommentar im Fence` und der Zeile danach — der Fence beendet den Abschnitt nicht. Hält.
- **P2 Elternabschnitt und Mutation:** `#eltern` umfasst beide Unterabschnitte. Mutation ein Zeichen im Körper von „Kind A“: gleicher Anker `vergleich roh: …#kind-a <-> …#kind-a cmp 1`; am Elternanker ebenfalls `cmp 1`. Rot, wie behauptet.
- **P3 Dubletten und gleiche Körper:** `#beispiel` und `#beispiel-1` lösen getrennt auf (Zähler hält). Zwei verschiedene Abschnitte mit gleichem Körper: `vergleich roh: HEAD:probe/t3.md#pflicht <-> HEAD:probe/t3.md#verbot cmp 0` (F-8). Zählung gleicher Körper über alle getrackten `.md` ohne `docs/reviews/` und `done/`: 0.
- **P4 HTML-`id`:** `vergleich 514405ba probe/p.md L30-30 f27b4e03 probe/p.md L30-30` → `cmp 0`, obwohl `vergleich … '#gezielt' … '#gezielt'` → `cmp 1` (F-1). `L27-27` (Leerzeile) → `leere Einheit`, Exit 2 (F-6). Reale `id`-Anker in `.harness/baseline/v6.14.1/regelwerk`: 5 Dateien; `git grep` der fünf `id`-Werte außerhalb der Baseline: ein Treffer, `MR-004` `#guard-haertung`.
- **P5 Normalisierung:** `#pins` v1 gegen v2 roh `cmp 1`, norm `cmp 0` (F-7). Mit `shopt -s nullglob` bzw. `failglob`: `vergleich norm: 514405ba:probe/p.md#kind-a <-> f27b4e03:docs/plan/adr/0100-….md cmp 0`, Exit 0; ohne Option Exit 1 (F-3).
- **P6 MR-Pins `5d8855d9`:** je Anker `vergleich 5d8855d9~1 …v6.14.0/… '#…' 5d8855d9 …v6.14.1/… '#…'`: `#vergabe-woher-die-nächste-kennung-kommt` `cmp 0`, `#grenzen--ehrlich-benannt` `cmp 0`, `#spec-straten-mehr-als-ein-spec-dokument` `leere Einheit` (Exit 2), `#guard-haertung` `leere Einheit` (Exit 2); `modul-13-quality-gates.md` als Datei roh `cmp 1`, norm `cmp 0` (F-2, F-9).
- **P7 `set -euo pipefail`:** Skript mit `set -euo pipefail`, `source` der Befehlsform, roter `vergleich`: Ausgabe nur `/dev/fd/63 /dev/fd/62 sind verschieden: …`, Exit 1, keine `vergleich`-Zeile (F-5). SIGPIPE: grüner `vergleich` am frühen Abschnitt `#1-einleitung` von `docs/user/benutzerhandbuch.md` (232049 Byte) unter `pipefail`, dreimal: Exit 0, Zeile gedruckt — kein Abbruch durch das frühe `exit` von `awk`.
- **P8 Rohheit:** `printf 'x\n'` gegen `printf 'x\n\n\n\n'` als Dateien: `vergleich … '' … ''` `cmp 0`; `cmp <(git show …t1.md) <(git show …t2.md)` Exit 1 (F-6).

**Suchlauf ([`AGENTS.md`](../../AGENTS.md) §3.13):** `make suchlauf-nachmessen PLAN=<Plan>` am Stand `e0c82845`: „8 Zeilen stimmen“, Exit 0. Fremde Träger: `git grep -iE 'Vergleichseinheit|zitat-korrektur-vergleichseinheit'` ohne `docs/reviews`, `done/`, `docs/plan/adr` und die Plan-Datei: Treffer in `AGENTS.md` (eigener Nachzug), in den zwei gemeldeten `state.md` und in `BEO-PGC/adr-aussage-breiter-als-ihre-messung/evidence/slice-harness-baseline-v6-14-1.md` (Record, unveränderlich) — die Meldung der zwei `state.md` ist vollständig. Die Einstufung von `harness/targets/pin-stale.md` und `.claude/agents/verifier.md` als „nicht geändert“ hängt an der Aussage aus F-2.

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `ADR-0158` Teil-Supersede (nur `ADR-0157` Entscheidung 1 (b)), Status-Zeile, Form wie bei `ADR-0157` gegenüber `ADR-0073` | geprüft, ohne Befund |
| `docs/plan/adr/README.md` Index-Zeile `ADR-0158`, Vermerk an der Zeile `ADR-0157`, Konventions-Satz | geprüft, ohne Befund — folgt dem Muster der Zeile `ADR-0073` |
| Lücke aus dem Re-Review F-1 (Anker-Wechsel auf dieselbe Datei; verschobener Lokator) | geprüft, ohne Befund — geschlossen, nachgemessen (P0) |
| Fence, Elternabschnitt, Unterabschnitte, Dubletten-Zähler, SIGPIPE | geprüft, ohne Befund (P1–P3, P7) |
| `ADR-0158` §Fitness Function: gemessene Zeilen gegen eigene Messung | geprüft, ohne Befund — die Trennung gemessen/hergeleitet steht in der Tabelle; die Abweichungen stehen in §Konsequenzen (F-1, F-2, F-8) |
| `AGENTS.md` §3.5 Kernsatz gegen `ADR-0158` Entscheidung 1/2 | geprüft — Befund F-4 |
| Suchlauf: Zahlen, Stände, Meldung fremder Träger | geprüft, ohne Befund — Einstufung zweier Träger hängt an F-2 |
| Commit-Messages `07923146`, `fc6d104c`, `e0c82845` (Traceability, `ADR-*`-Kennung) | geprüft, ohne Befund |
| Docker-only / Umleitung auf Repo-Dateien im Diff | geprüft, ohne Befund — reine Markdown-Änderungen |
| Produkt-Code, Spec-Straten | nicht berührt |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 2 |
| MEDIUM | 2 |
| LOW | 3 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** ADR-Aussage breiter als ihre Messung · Beleg trägt seinen Satz nicht · Messbefehl fail-open · Nachzug widerspricht dem Nachbarn im selben Träger · Befehlsform-Randverhalten unbenannt · Grenze der Normalisierung unbenannt · Verweis löst nicht auf

## Verdikt

**Merge-blockierend:** ja — F-1 und F-2 sind HIGH an einer `Accepted`-ADR im
Diff. Eine Berichtigung ist nach `AGENTS.md` §3.5 eine Folge-ADR; der Weg läuft
über den Architect (Konflikt-Pfad als Sequenz, Baseline-Regelwerk `v6.14.1` ·
`regelwerk/modul-08-agentenrollen.md` §Konflikt-Pfad als Rollen-Sequenz), nicht
über eine Herabstufung. F-3 betrifft dieselbe Befehlsform und gehört in
denselben Zug. F-4 ist ein Implementer-Nachzug in `AGENTS.md` §3.5, der dem
Architect-Verdikt folgt.

**DoD-Häkchen „Review durchgeführt“:** nicht nachgezogen — die Findings
verlangen eine Fixrunde (`.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug
ohne Fixrunde, Grenze).

**Übergabe:** Findings an den Architect (F-1, F-2, F-3, F-7, F-8) und den
Implementer (F-4); F-9 an den Planner als fremder Träger. Die Finding-Klassen
gehen in die Slice-Closure §7. Dieser Report ersetzt keine Verifikation.
