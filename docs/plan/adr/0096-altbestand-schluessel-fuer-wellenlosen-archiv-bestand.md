# ADR-0096: `altbestand` als dedizierter Schlüssel für den wellenlosen Archiv-Bestand

**Status:** Accepted

**Datum:** 2026-09-18

**Autor:** Architect (pt9912)

**Bezug:** `AGENTS.md` §3.5 (der hier gesetzte Schlüssel ist nach Annahme
praktisch ortsfest — ein zweiter Schlüssel für denselben Zweck wäre eine
neue Entscheidung, keine Korrektur), §3.6 (Sinngemäß: eine dauerhaft
wirkende, ins Repo schreibende Entscheidung verdient dieselbe Vorsicht wie
eine Gate-Lockerung, auch wenn `archive-welle` selbst kein Gate ist) ·
Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle braucht ·
`slice-archive-altbestand-adr` (Auftrag dieser ADR)

**Schärft:** — (Prozess-ADR ohne Spec-Stratum)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Dieses Repo hat noch nie archiviert. Ein Dry-Run des vendorten Werkzeugs
(`.harness/state/bin/ai-harness-init archive-welle --vorschau`, real am
Stand 2026-09-18 ausgeführt) zeigt: **jeder** erste schreibende Lauf —
gleich welche Welle-Kennung er trägt — sammelt zusätzlich zu seinen
eigenen Mitgliedern **jeden** wellenlosen `done/`-Slice ein, weil
`internal/archive/scan.go`s Untergrenzen-Wächter erst greift, sobald
irgendein `done/*/archiv.zip` existiert. Real gemessen: 41 wellenlose
Slices liegen flach in `docs/plan/planning/done/` (`grep` über
`**Welle:**` = „ohne Welle"). Ein Lauf gegen den bereits geschlossenen
`welle-d-check` (2 Mitglieder) würde denselben Bestand mitreißen und danach
43 statt 2 archivierte Vorgänge behaupten.

Real gemessen tragen sowohl der Schlüssel `altbestand` als auch
`welle-d-check` heute zwei Sperren, die der schreibende Lauf nähme:
`[untergrenze]` (keine Untergrenze existiert) und `[haenger]` (mehrere
`Accepted`-ADRs bzw. Records verweisen noch auf Review-Reports, die der
Lauf ersatzlos löschen würde — die zugehörige Zitat-Korrektur-Runde hat die
Zahl der Hänger-Kanten bereits von 217 auf 208 gesenkt, die Sperre selbst
bleibt aktiv). Für `altbestand` zusätzlich zwei weitere Sperren:
`[ergebnisnotiz]` (`docs/plan/planning/done/altbestand-results.md` fehlt)
und `[kein-plan]` (kein Welle-Plan `altbestand*` in `done/`) — die hier
vendorte Werkzeug-Version hebt diese beiden für einen Schlüssel ohne
Welle-Form (noch) nicht auf, anders als die neuere Fassung im
Schwester-Repo `ai-harness-init` (dort `AltbestandSchluessel`-Sonderfall in
`internal/archive/vorschau.go`, eingeführt über dessen eigenes `ADR-0041`).

Das Schwester-Repo hat exakt dieselbe Frage bereits für sich entschieden:
ein einzelnes, dediziertes Sammel-Archiv ohne Welle-Zugehörigkeit,
Schlüssel `altbestand`, statt Vermischung mit einer echten, inhaltlich
definierten Welle.

## Entscheidung

Wir wählen **einen dedizierten Schlüssel `altbestand`** für den
wellenlosen Archiv-Bestand dieses Repos — dieselbe Kennung wie im
Schwester-Repo, aber eine eigene Entscheidung für dieses Repo, nicht
übernommen. Der Schlüssel trägt **ausschließlich** Slices ohne eigene
Wellen-Zugehörigkeit, die vor der ersten realen Archivierung entstanden
sind; er ist nach Annahme ortsfest (`AGENTS.md` §3.5-Analogie).

`welle-d-check` bleibt davon unberührt: seine spätere eigene Archivierung
trägt weiterhin nur seine 2 Mitglieder. Beide Schlüssel sind unabhängige
Archiv-Läufe desselben Werkzeugs; keiner löst den anderen ab, keine
`Supersedes`-Beziehung. Welcher der beiden Läufe zuerst die Untergrenze
setzt, entscheidet der Vollzugs-Slice — die Zuordnungsfrage selbst hängt
nicht an der Reihenfolge.

Da die hier vendorte Werkzeug-Version die Sperren `[ergebnisnotiz]` und
`[kein-plan]` für den Schlüssel `altbestand` nicht aufhebt, schreibt der
Vollzugs-Slice vor dem realen Lauf zwei kurze, von Hand verfasste
Dokumente direkt nach `docs/plan/planning/done/` — einen knappen
Welle-Plan `altbestand.md` und seine Ergebnisnotiz
`altbestand-results.md`, aus den vorhandenen Welle-Vorlagen kopiert und
ausgefüllt. Das ist ein einmaliger, kleiner, klar umrissener Schritt —
keine neue Prüfregel, kein neues Werkzeug, keine Änderung an
`ai-harness-init`.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — `welle-d-check` als Vehikel | kein zusätzliches Dokument nötig (Plan und Ergebnisnotiz liegen bereits in `done/`); sofort lauffähig | behauptet danach dauerhaft eine Zugehörigkeit (43 statt 2 archivierte Vorgänge), die zu Ziel und Inhalt jener Welle nichts sagt — der Stub ist danach nicht mehr korrigierbar (kein Zitat-Korrektur-Fall, sondern eine Tatsachenbehauptung über den Archiv-Inhalt); dasselbe Schwester-Repo hat diese Option für dieselbe Frage bereits verworfen |
| B — nichts tun (Altbestand bleibt flach liegen) | kein Aufwand jetzt | `[untergrenze]` feuert für **jeden** ersten Lauf, gleich welcher Kennung — verwirft nicht nur die Archivierung des Altbestands, sondern jede künftige Archivierung dieses Repos, dauerhaft |
| **C — dedizierter Schlüssel `altbestand` — gewählt** | sauberer Geltungsbereich (kein Archiv behauptet eine Zugehörigkeit, die es nicht gibt); setzt die Untergrenze für alle künftigen Läufe; deckt sich mit der bereits gelösten Präzedenzfrage des Schwester-Repos, ohne sie blind zu kopieren | verlangt zwei von Hand verfasste Dokumente, weil die vendorte Werkzeug-Version die Sperren-Aufhebung für diesen Schlüssel noch nicht trägt — einmaliger, kleiner, klar begrenzter Mehraufwand |

**Fazit:** C. A löst zwar sofort, aber auf Kosten einer dauerhaften
Falschbehauptung, die dieselbe Fehlerklasse ist, die das Schwester-Repo für
sich bereits ausgeschlossen hat. B ist keine neutrale Option — es blockiert
jede künftige Archivierung. Der Mehraufwand von C ist klein, einmalig und
bereits im Vollzugs-Slice als Schritt vorgesehen; er erzeugt keine neue
Prüfregel und keine neue Bürokratie.

## Konsequenzen

- Positiv: Der Altbestand bekommt einen eigenen, wahrheitsgemäßen
  Archiv-Anker; `welle-d-check` bleibt unverfälscht bei seinen 2
  Mitgliedern.
- Positiv: Die erste erfolgreiche Archivierung (gleich welcher Schlüssel
  zuerst läuft) setzt die Untergrenze — danach greift die laufende
  Untergrenzen-Regel ohne weitere Zuordnungsfrage.
- Negativ: Zwei zusätzliche, von Hand verfasste Dokumente vor dem ersten
  realen Lauf gegen `altbestand` — bereits als Schritt im Vollzugs-Slice
  vorgesehen, kein neuer Slice nötig.
- Folgepflicht: `slice-archive-altbestand-vollzug` schreibt
  `docs/plan/planning/done/altbestand.md` und
  `docs/plan/planning/done/altbestand-results.md`, bevor er den realen,
  schreibenden Archiv-Lauf startet.

**Der `structure`-Modul-Blindfleck (Closure-Notiz-Regel) — akzeptiertes
Negativ, keine Folgepflicht.** `.d-check.yml`s `structure`-Regel prüft
`docs/plan/planning/done/slice-*.md` über einen flachen Glob (kein `**`)
auf eine vollständige `## 7. Closure-Notiz`. Ein durch diese ADR
archivierter Slice liegt danach unter `done/altbestand/slice-*.md` und
verlässt diesen Scan-Bereich — der an seiner Stelle stehende Stub
(`archiv-stub-slice.template.md`) trägt aber per Form **keine**
Abschnittsüberschriften mehr und behauptet damit auch keine
Closure-Notiz. Die Prüfung geht nicht verloren, weil der geprüfte
Gegenstand nicht mehr existiert, nicht weil eine Lücke im Sensor entstünde:
Ein Stub, der keine Closure-Notiz mehr behauptet, hat auch keine mehr zu
erfüllen. Kein neuer Beobachtungs-Eintrag, kein Folge-Slice.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| — | Die Schlüssel-Wahl selbst ist ein Urteil, kein Sensor prüft sie. Maschinell beobachtbar ist nur das Ergebnis eines Laufs: `archive-welle --vorschau altbestand` meldet nach den zwei Dokumenten und der `[haenger]`-Bereinigung „Sperren: keine" | — (kein Gate; `archive-welle` ist Werkzeug, `harness/README.md` §Sensors) |

## Re-Evaluierungs-Trigger

Beobachtbare Trigger: **(a)** die hier vendorte Werkzeug-Version wird auf
eine Fassung angehoben, die für den Schlüssel `altbestand` die Sperren
`[ergebnisnotiz]`/`[kein-plan]` selbst aufhebt — dann entfällt die Pflicht
zu den zwei Hand-Dokumenten, die Schlüssel-Wahl selbst bleibt unberührt;
**(b)** ein zweiter, inhaltlich eigenständiger wellenloser Bestand entsteht
nach der ersten Archivierung — dann entscheidet eine neue ADR, ob er
denselben Schlüssel erneut trägt oder einen eigenen bekommt (kein
automatischer Nachschlag). Andernfalls permanent.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-18 | Accepted — Anlass: `slice-archive-altbestand-adr`; dedizierter Schlüssel `altbestand` statt `welle-d-check`-Vehikel oder Nichtstun; `structure`-Modul-Blindfleck als akzeptiertes Negativ eingestuft | `slice-archive-altbestand-adr` |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0096` (Baseline-Regelwerk `modul-04-adrs.md` §Hard Rule für
Accepted-ADRs).
