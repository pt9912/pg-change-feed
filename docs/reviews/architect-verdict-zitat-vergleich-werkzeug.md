# Architect-Verdikt — Werkzeug für die Referent-Messung (`slice-zitat-vergleich-werkzeug`): normativer Träger, F-1 bis F-3, F-5, F-6, Form

**Datum:** 2026-10-06 · **Stand:** `2d21f9ab` (Baum sauber) · **Rolle:** Architect (Modul 8),
Zug vor dem Code von `slice-zitat-vergleich-werkzeug`, anderer Kontext als die Architect-Läufe
von [`ADR-0158`](../plan/adr/0158-zitat-korrektur-vergleichseinheit-je-verweisform.md) und
[`ADR-0159`](../plan/adr/0159-zitat-korrektur-html-id-mr-pins-und-gehaertete-befehlsform.md)
und als der Reviewer-Lauf · **Rolleninhaber:** pt9912 (Architect-Agent im Auftrag) ·
**Anlass:** Re-Evaluierungs-Trigger (a) von `ADR-0159` ist eingetreten
([`review-slice-zitat-korrektur-vergleichseinheit-fixrunde`](review-slice-zitat-korrektur-vergleichseinheit-fixrunde.md)
F-1 bis F-3); Vorgabe des Auftraggebers: `ADR-0159` ist die letzte Folge-ADR der Kette
`ADR-0157` → `ADR-0158` → `ADR-0159`.

**Bezug:** [`ADR-0159`](../plan/adr/0159-zitat-korrektur-html-id-mr-pins-und-gehaertete-befehlsform.md)
(Entscheidung 1, 3, 4, Re-Evaluierungs-Trigger (a), Option E),
[`ADR-0158`](../plan/adr/0158-zitat-korrektur-vergleichseinheit-je-verweisform.md) (Entscheidung 1),
[`ADR-0157`](../plan/adr/0157-zitat-korrektur-reichweite-nach-aussage-und-mr-pins.md) (Entscheidung 4),
[`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md),
[`AGENTS.md`](../../AGENTS.md) §3.1, §3.5, §3.12.

---

## 0. Ergebnis in einem Blick

- **(1) Normativer Träger — keine neue ADR, dieses Verdikt genügt.** `ADR-0159` bleibt Träger
  der **Semantik** (Einheiten, Normalisierung, die vier Zusagen von Entscheidung 4). Das
  Werkzeug wird Träger der **Messung**. Der `bash`-Block in Entscheidung 4 bleibt als
  historische Fassung stehen.
- **(2) F-1, F-2, F-3 — Auslegung, keine Semantik-Änderung.** F-2 und F-3 sind
  Erkennungsfehler. F-1 ist der Grenzfall; er wird nach dem erkennbaren Zweck von
  Entscheidung 1 ausgelegt (§3). Wer das als Änderung liest, hat die Ausweichform in §3.1.
- **(3) F-5:** Ein Tag-Paar wird nur angenommen, wenn `.harness/baseline/<alt-tag>` am alten
  und `.harness/baseline/<neu-tag>` am neuen Stand liegt. **F-6:** Das Skript setzt
  `LC_ALL=C.UTF-8` selbst und prüft vor der Messung, ob `awk` Multibyte kann. Beides ist
  fail-closed.
- **(4) Form:** `tools/harness/zitat-vergleich.sh`, `make zitat-vergleich ARGS=…`,
  `make test-zitat-vergleich`, Vertrag `harness/targets/zitat-vergleich.md`. GNU Awk wird
  **nicht** nach Name verlangt, sondern über eine Fähigkeitsprobe (§5).
- **Datei-Schleife aus `ADR-0157` Entscheidung 4:** Sie bleibt draußen, als akzeptiertes Negativ.
  Ihr Exit wird nicht gefärbt, und der Verifier liest ihre Ausgabe ohnehin. Dass sie draußen
  bleibt, verschiebt nichts und lässt keine Spur verschwinden.

---

## 1. Messungen dieses Zugs

Klon des Repos im Scratchpad des Zugs (`architect-werkzeug/`), Stand `2d21f9ab`, GNU Awk 5.2.1,
GNU bash, Linux. Die Befehlsform ist per `awk` wörtlich aus dem `bash`-Block von `ADR-0159`
Entscheidung 4 gezogen (56 Zeilen, `bash -n` Exit 0) und per `source` geladen („ADR-Form“).
Der **Prototyp** ist eine Kopie davon mit den drei Änderungen aus §3 (nur im Scratchpad, nicht
committet). Die Probe-Commits `b82c0d8c` und `3fee3c7c` gibt es nur im Klon. Sie enthalten eine
Datei `probe/f.md` mit vier Fällen: gestapelte `id`s vor einem Heading, `<a id="attr"
class="k"></a>`, `<a id="selbst"/>` sowie ein Fence aus vier Backticks mit einem inneren Fence.
Der zweite Commit ändert in jedem Abschnittskörper ein Wort.

| Nr | Was, womit | Ergebnis (gedruckt, gekürzt) |
|---|---|---|
| M1 | ADR-Form, `vergleich b82c0d8c probe/f.md '<ref>' 3fee3c7c probe/f.md '<ref>'` | `#oben` `cmp 0`, `#unten` `cmp 1`, `#attr` `cmp 0`, `#selbst` `cmp 0`, `#ungerade` `cmp 0` — F-1, F-2 und F-3 reproduziert |
| M2 | Prototyp, dieselben Aufrufe | `#oben` `cmp 1`, `#unten` `cmp 1`, `#attr` `einheit: leere Einheit` / `keine Einheit, Exit 2`, `#selbst` ebenso Exit 2, `#ungerade` `cmp 1`, `#ende` `cmp 0` (der Abschnitt hinter dem Fence beginnt am richtigen Heading) |
| M3 | Gleichstand ADR-Form gegen Prototyp: `einheit HEAD <datei> '#<anker>'` für jeden Heading-Slug und jede `<a id="…">` in allen getrackten `.md` unter `.harness/baseline/v6.14.1/regelwerk/`, `docs/plan/adr/`, `harness/` und `AGENTS.md`; Vergleich von Ausgabe und Exit je Anker | `anker=2103 abweichend=0 leer_alt=6 leer_neu=6` |
| M4 | `git grep` über getrackte `.md` am Stand `2d21f9ab` | Fences aus vier oder mehr Zeichen: 0. Um 1 bis 3 Leerzeichen eingerückte Fences: 46 Zeilen (in M3 enthalten, Gleichstand). `<a id="…"` mit weiterem Attribut oder selbstschließend: 4 Treffer, alle in Inline-Code (Slice-Plan und Re-Review) |
| M5 | `git ls-tree --name-only <stand> .harness/baseline/` | `5d8855d9~1` und `5d8855d9`: `v6.14.0 v6.14.1`; `11a5bac5`: `v6.14.0`; `625ddbef` und `2d21f9ab`: `v6.14.1` |
| M6 | ADR-Form, `vergleich HEAD <modul-13> '#guard-härtung-wächter-reifen-in-wellen-modul-13'` gegen sich selbst unter drei Locales; `printf 'Ä' \| LC_ALL=<L> awk '{print tolower($0), length($0)}'` | `LC_ALL=C` Exit 2, `C.UTF-8` Exit 0, `de_DE.UTF-8` Exit 0. Unter `C`: `tolower` = `Ä`, Länge 2. Unter `C.UTF-8`: `ä`, Länge 1 |
| M7 | Wegwerf-Makefile, `bash -c '…' x $(ARGS)`, drei Aufrufe | `ARGS="a p.md '#x' b p.md '#x'"` ergibt 6 Argumente. `ARGS="a p.md #x b p.md #x"` ergibt **2** Argumente, weil die Shell ab `#` einen Kommentar liest. `ARGS="a p.md '' b p.md ''"` ergibt 6 Argumente, darunter zwei leere |

**Übernommen** (nicht nachgemessen): Die Befunde F-1 bis F-6 und ihre Proben N3 bis N7 und L
stammen aus dem Re-Review. Dass `make` den Exit eines Rezepts ≠ 0 als Exit 2 meldet, steht
so im Bestand (`harness/README.md`, Zeile `make kommentar-kennungen`).

---

## 2. (1) Normativer Träger

**Verdikt:** Kein Folge-ADR. Ein Verdikt-Artefakt genügt. **Begründung:** `ADR-0159` sieht
den Weg selbst vor. Re-Evaluierungs-Trigger (a) lautet: „dann Option E als eigener Slice, statt
einer weiteren Folge-ADR zur Befehlsform“. Der Satz in Entscheidung 4 („bleibt eine
Verfahrensregel in dieser ADR, kein Werkzeug“) beschreibt den Stand bis zu diesem Trigger.
Mit dem Trigger hat sich die Entscheidung nicht geändert, sie läuft nur auf dem Weg, den sie
selbst vorsieht.

Rangfolge, sobald das Werkzeug gemergt ist:

| Frage | Träger |
|---|---|
| Was ist die Einheit, was wird normalisiert, was ist „roh“, wann endet der Lauf mit 2 | `ADR-0159` Entscheidung 1 bis 3, 5 und die vier Zusagen von Entscheidung 4 (mit den Teilen von `ADR-0158`/`ADR-0157`, die sie in Kraft lässt), ausgelegt durch §3 dieses Verdikts |
| Womit wird gemessen, was gilt als Beleg | das Werkzeug (`make zitat-vergleich`) und seine gedruckte Zeile |
| Der `bash`-Block in Entscheidung 4 | historische Fassung; die ADR ist unveränderlich ([`AGENTS.md`](../../AGENTS.md) §3.5), es gibt keine Zitat-Korrektur und keine Zeile in §Geschichte |

- **Werkzeug und Block laufen auseinander:** In den fünf Punkten aus §3 und §4 (F-1, F-2,
  F-3, F-5, F-6) ist das gewollt. Der Vertrag nennt diese Abweichungen. An jeder anderen
  Stelle ist das ein Fehler des Werkzeugs. Er wird im Skript und im Tabellentest behoben,
  ohne ADR.
- **Werkzeug und Semantik laufen auseinander:** Dann ist das Werkzeug falsch. Die Semantik
  ändert nur eine Folge-ADR mit `Supersedes ADR-0159`.
- **Weitere Korrekturen an der Messform:** Sie gehen ins Skript und in den Tabellentest.
  Trigger (a) ist damit verbraucht. Die Trigger (b) und (c) gelten unverändert.
- **Risiko 1 aus §6 des Plans („Zwei Träger“):** Ausgang *entfallen*. Die Rangfolge oben
  löst die Frage, welcher Träger gilt. Den Gleichstand außerhalb der fünf Punkte hat M3
  für die Einheit gemessen. Der Implementer fährt M3 am fertigen Skript nach und trägt die
  gedruckte Zeile in den Plan ein.

---

## 3. (2) Re-Review F-1 bis F-3

| Finding | Einordnung | Regel für das Werkzeug | Beleg |
|---|---|---|---|
| **F-3** Fence aus vier Zeichen | Erkennungsfehler. Die Semantik „Headings in Code-Fences zählen nicht“ (`ADR-0158` Entscheidung 1) steht fest, die ADR-Form erkennt den Fence falsch | Ein Fence öffnet mit 0 bis 3 Leerzeichen Einzug und mindestens drei gleichen Zeichen (`` ` `` oder `~`); Zeichen und Länge werden gemerkt. Er schließt nur mit demselben Zeichen in mindestens dieser Länge, danach höchstens Leerraum (CommonMark). | M2 `#ungerade` `cmp 1`, `#ende` `cmp 0`; M3 Gleichstand, auch an den 46 eingerückten Fences (M4) |
| **F-2** `id` mit Attribut, selbstschließend | Erkennungsfehler gegen den Wortlaut von Entscheidung 1: „Gelesen wird nur `<a id="…">`; `name=` und eine `id` an einem anderen Element ergeben eine leere Einheit (Exit 2)“. `<a id="x" class="k">` und `<a id="x"/>` haben diese Form nicht, also fallen sie in denselben Zweig wie `name=` | Gelesen wird nur das öffnende Tag `<a id="X">` in genau dieser Form. Steht `<a id="X"` außerhalb von Fence und Inline-Code in einer anderen Form, endet die Einheit mit Exit 2 und mit einer Zeile, die die Form nennt. Das gilt auch dann, wenn ein Heading-Slug `X` trifft. | M2 `#attr` und `#selbst` Exit 2; M4: 0 Fundstellen außerhalb von Inline-Code |
| **F-1** gestapelte `id`s | **Auslegung** (Grenzfall, Begründung unten) | Eine Zeile, die nach dem Entfernen aller `<a id="…"></a>` nur Leerraum trägt, ist eine **Zeile ohne Inhalt**. „Nächste nicht leere Zeile“ in den Zeilen 2 und 3 der Tabelle von Entscheidung 1 heißt: die nächste Zeile **mit Inhalt**. Gestapelte `id`s vor einem Heading adressieren also alle den Abschnittskörper dieses Headings. Vor einem Absatz adressieren sie den Block ab diesem Absatz. | M2 `#oben` `cmp 1`, `#unten` `cmp 1`; M3 Gleichstand |

**Warum F-1 eine Auslegung ist und keine Änderung:**

- (a) `ADR-0159` definiert „ohne Inhalt“ selbst. §Kontext sagt: „Die Zeile `<a id="…"></a>`
  vor einem Heading hat keinen Inhalt“. Die eigene Befehlsform prüft genau das: `gsub` der
  `id`-Tags, dann Prüfung auf Leerraum. Die Auslegung wendet diese Definition auch auf die
  Folgezeile an.
- (b) Wörtlich auf den gestapelten Fall angewandt, ergibt die Tabelle eine Einheit ohne
  Inhalt. Das ist genau der Defekt, gegen den Entscheidung 1 geschrieben ist (erstes
  Review F-1). Es widerspricht außerdem der Konsequenz „Ein Verweis auf eine HTML-`id` vor
  einem Heading misst den Abschnitt“ und dem Grundsatz „fail-closed gewählt“.
- (c) Trigger (a) nennt „`id`-Stellung“ ausdrücklich als Gegenstand von Option E statt einer
  Folge-ADR.
- (d) Es gibt 0 Fundstellen. Die Fitness Function von `ADR-0159` führt nicht gemessene
  Stellungen als *hergeleitet*. Kein gemessenes Ergebnis ändert sich (M3).

### 3.1 Ausweichform, falls F-1 als Änderung gelesen wird

Gestapelte `id`s ohne Inhalt enden mit Exit 2 („Stellung nicht gelesen“). Das ist fail-closed
und kostet heute nichts, weil es 0 Fundstellen gibt. Die Semantik bliebe dann offen, bis eine
Folge-ADR sie festlegt. Ich empfehle diese Form **nicht**: Die Auslegung trägt (a) bis (d)
und liefert am selben Fall eine echte Messung statt eines „nicht messbar“.

---

## 4. (3) F-5 und F-6

- **F-5 Tag-Paar — Durchsetzung von Entscheidung 3 („das Tag-Paar des Bumps“), keine Änderung.**
  - Ein siebtes Argument wird nur angenommen, wenn alle drei Bedingungen gelten:
    - Es hat die Form `v<X.Y.Z>:v<X.Y.Z>`.
    - Beide Tags sind verschieden.
    - `git ls-tree` findet `.harness/baseline/<alt-tag>` am alten Stand **und**
      `.harness/baseline/<neu-tag>` am neuen Stand.
  - Sonst endet der Lauf mit Exit 2 und einer gedruckten Zeile.
  - Verlangt wird **nicht**, dass der neue Tag am alten Stand fehlt. Am Pin-Commit
    `5d8855d9~1` liegen schon beide Tags (M5); diese Bedingung würde den realen Fall abweisen.
  - Die geforderten Bäume liegen an allen real gemessenen Paaren vor (M5): `5d8855d9`,
    `11a5bac5` gegen `625ddbef`. Ein fremdes Paar wie `v0.79.0:v0.80.0` hat keinen solchen Baum
    und wird abgewiesen. Das ist *hergeleitet* aus M5; der Tabellentest belegt es.
  - Der Tabellentest legt im Wegwerf-Repo die Verzeichnisse `.harness/baseline/<tag>/` an.
- **F-6 Locale — fail-closed und reproduzierbar.**
  - Das Skript setzt `export LC_ALL=C.UTF-8` selbst und hängt damit nicht von der Locale des
    Aufrufers ab.
  - Vor der ersten Messung prüft es, ob `awk` Multibyte kann:
    `printf 'Ä' | awk '{print tolower($0), length($0)}'` muss `ä 1` ergeben. Sonst endet der
    Lauf mit Exit 2 und der Zeile „awk ohne Multibyte/UTF-8“. Damit sind eine fehlende Locale
    und ein `awk` ohne Multibyte-Unterstützung abgedeckt (M6).
  - Der Tabellentest ruft das Skript unter `LC_ALL=C` auf und erwartet am Umlaut-Slug grün.
    Mit einem Stub-`awk` im `PATH`, der die Probe nicht besteht, erwartet er Exit 2.
  - Damit hat Risiko 2 aus §6 des Plans seinen Beleg.

---

## 5. (4) Form des Werkzeugs — Vorgabe an den Implementer

- **Dateien:**
  - `tools/harness/zitat-vergleich.sh`: enthält die Funktionen `einheit`, `tagnorm` und
    `vergleich`; der Aufruf ist `vergleich "$@"`.
  - `tools/harness/run-zitat-vergleich-tests.sh`: der Tabellentest im Wegwerf-Repo, im Muster
    von `run-suchlauf-nachmessen-tests.sh`. Der Prüfling lässt sich per `PROG=<Datei>`
    übersteuern, für Mutationsläufe an Kopien.
  - `harness/targets/zitat-vergleich.md`: der Vertrag.
- **Make-Ziele** stehen im `Makefile` neben `suchlauf-nachmessen`, nicht in `GATE_CHECKS`:
  - `make zitat-vergleich ARGS="<alt-stand> <alt-pfad> '<alt-ref>' <neu-stand> <neu-pfad> '<neu-ref>' [<alt-tag>:<neu-tag>]"`.
    Ohne `ARGS` bricht das Ziel mit `$(error …)` ab.
  - `make test-zitat-vergleich`.
- **Argumente fail-closed:**
  - Genau 6 oder 7 Argumente, sonst Exit 2 mit Gebrauchszeile. M7 zeigt, dass ein
    ungequotetes `#…` in `ARGS` die Argumente stumm auf 2 kürzt.
  - Die Referenz ist leer, beginnt mit `#` oder hat die Form `L<a>-<b>`. Alles andere endet
    mit Exit 2, wie in der ADR-Form.
  - Der Vertrag zeigt die Quotierung von `'#anker'` und `''` an einem Beispiel.
- **Exit und Farbe:**
  - Das Skript endet mit 0, 1 oder 2 und druckt die Zeile vor dem Ende, wie in der ADR-Form.
  - Über `make` kommt jeder Exit ≠ 0 als Exit 2 an. Der Vertrag sagt deshalb: Die Farbe steht
    in der gedruckten Zeile (`cmp 0`, `cmp 1`, `… Exit 2`). Ein Beleg zitiert diese Zeile,
    nicht den Exit von `make`.
- **Host-Werkzeuge** nach [`AGENTS.md`](../../AGENTS.md) §3.1 (Klasse ohne Installation):
  - `bash`, `git`, `awk`, `sed`, `cmp`.
  - `awk` muss Multibyte unterstützen; das prüft die Fähigkeitsprobe aus §4. Gemessen ist das
    an GNU Awk 5.2.1.
  - **GNU Awk wird nicht nach Name verlangt**, und das Skript benutzt keine Funktion, die
    nur gawk kennt (`gensub`, `PROCINFO`, `asorti` und ähnliche). Kein anderes Skript im
    Repo setzt gawk voraus; `tools/harness/extract-command.awk` erklärt sich ausdrücklich als
    POSIX-awk.
  - Die Probe ist der fail-closed Wächter für die eine nötige Fähigkeit. Diese Grenze nennt
    der Vertrag.
- **Regeln aus §3/§4** setzt das Skript um, und der Vertrag nennt sie als Abweichungen vom
  historischen Block: F-1 (Zeile ohne Inhalt), F-2 (nur `<a id="X">`, sonst Exit 2), F-3
  (CommonMark-Fence), F-5 (Tag-Paar an den Baseline-Bäumen der Stände), F-6 (Locale und
  Fähigkeitsprobe).
  - Ein Randfall, den der Prototyp nicht trägt: Ein Fence direkt nach einer `id`-Zeile ohne
    Inhalt (Zustand „vor“) beginnt den Block, und seine Öffnungszeile gehört zur Einheit.
  - Die Zusagen von Entscheidung 4 (Quotierung, Exit 2 auf jeder Seite, Zeile vor `return`,
    roh mit Wächter-Zeichen) gelten unverändert.
- **Tabellentest — zusätzlich zu den Fällen in DoD-Punkt 2:**
  - Falsche Zahl der Argumente.
  - Ein fremdes Tag-Paar, das formal gültig ist, aber am Stand keinen Baseline-Baum hat.
  - Ein Tag-Paar mit gleichen Tags.
  - Ein Stub-`awk` ohne Multibyte.
  - Eine `id` in Attributform **und** ein gleichnamiger Heading-Slug in derselben Datei:
    Exit 2.
  - Gestapelte `id`s vor einem Absatz: Einheit ist der Block.
  - Ein Fence mit `~~~` innerhalb eines Fence mit `` ``` ``.
- **Mutationen (§3.12):** Je neuer Regel eine Mutation am Skript, die ihren Fall rot färbt.
  Je Mutation nennt der Plan Stelle, Instanz (Tabellentest) und gesehene Farbe. Was darüber
  hinaus verallgemeinert wird, steht als *hergeleitet*.
- **Gleichstand:** M3 am fertigen Skript nachfahren, über dieselben Dateien wie oben, mit der
  ADR-Form als Gegenseite. Die gedruckte Zeile steht im Plan. Eine Abweichung außerhalb der
  fünf Punkte ist ein Fehler des Werkzeugs.
- **Träger (DoD-Punkt 3):** wie im Plan. `AGENTS.md` §3.5 nennt das Ziel und dieses Verdikt
  als Rang-Zeiger für die Auslegung. `ADR-0159` und `docs/plan/adr/README.md` bleiben
  unberührt.
- **Nicht im Slice:** die Datei-Schleife aus `ADR-0157` Entscheidung 4 (§0), ein Gate
  ([`AGENTS.md`](../../AGENTS.md) §3.6) und eine Folge-ADR.
