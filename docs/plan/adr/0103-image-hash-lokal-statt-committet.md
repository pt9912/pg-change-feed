# ADR-0103: `harness/image-hash.txt` lokal statt committet (Supersedes ADR-0044, teilweise)

**Status:** Accepted — Supersedes [`ADR-0044`](0044-image-beleg-semantik.md)
(nur deren Entscheidungspunkt 3, die welle-1-Regel „Ein Zug, der
Build-Kontext-Dateien ändert, läuft `make image` vor seiner Closure; der
Digest-Commit entfällt bei unverändertem Digest"; die Entscheidungspunkte
1, 2 und 4 dieser ADR — Digest als Lauf-Beleg-Semantik, Inhalts-Streit über
den Binary-sha256, der Extraktions-Mechanismus als belegbarer Weg — bleiben
unverändert bestehen und werden hier nicht wiederholt)

**Datum:** 2026-09-18

**Autor:** pt9912

**Bezug:** [`ADR-0044`](0044-image-beleg-semantik.md) (die teilweise
supersedete ADR) ·
`docs/plan/planning/observations/BEO-PGC/image-digest-nichtdeterminismus-erzeugt-merge-konflikt/`
(Anlass) · Baseline-Regelwerk `modul-14-docker-harness.md` §Zwei Formen des
Reproduzierbarkeits-Ankers (Archiv- vs. Rezept-Form) · [`ADR-0084`](0084-sync-gate-fuer-generierte-artefakte.md)
Entscheidungspunkt 4 (bleibt unverändert richtig, siehe §Konsequenzen) ·
`.gitignore` · `harness/README.md` §Werkzeuge (`make image`-Zeile) ·
`Makefile` (`image`-Target)

**Schärft:** — (Prozess-ADR ohne Spec-Stratum, wie [`ADR-0044`](0044-image-beleg-semantik.md))

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

**(1) Der Anlass.** `harness/image-hash.txt` trägt seit `welle-1` den
Image-Digest des letzten `make image`-Laufs als committetes Artefakt.
[`ADR-0044`](0044-image-beleg-semantik.md) hat bereits real belegt, dass
dieser Digest **builder- und lauf-gebunden** ist, nicht inhaltsadressierend
— ein bit-identisches Binary erzeugte an zwei Range-Grenzen zwei
verschiedene Digests. Genau diese Eigenschaft erzeugt reales
Merge-Rauschen: Jeder Zug, der `make image` erneut laufen lässt, ändert
potenziell den committeten Wert dieser einen Zeile, unabhängig davon, ob
sich der tatsächliche Build-Inhalt geändert hat. Zwei parallele Branches,
die beide `make image` laufen lassen, kollidieren beim Merge auf dieser
Zeile, ohne dass einer von beiden inhaltlich etwas Konfliktträchtiges
beigetragen hätte — real aufgetreten und dokumentiert in
`BEO-PGC/image-digest-nichtdeterminismus-erzeugt-merge-konflikt`
(`slice-nats-drittstream-example-go`: ein fälschlich committeter
Digest-Wechsel, weil ein Build-Kontext-Pfad fälschlich als betroffen
angenommen wurde, im selben Zug per Review-Fund zurückgenommen).

**(2) Warum der Baseline-Grund fürs Committen hier nicht trägt — die
Archiv-Bedingung.** Baseline-Regelwerk `modul-14-docker-harness.md` §Zwei
Formen des Reproduzierbarkeits-Ankers unterscheidet zwei Formen: die
**Archiv-Form** (der Digest des gebauten Images — Bedingung: das Image
wird aufbewahrt, gepusht oder nachweislich vorgehalten) und die
**Rezept-Form** (Commit der Quellen plus gepinnte Digests der
Eingangs-Images — braucht kein Lager, kann den Digest des *gebauten*
Images aber nicht als Griff nutzen). `make image` baut in diesem Repo
ausschließlich lokal (`docker buildx build --load`); kein `docker push`
existiert irgendwo im Repo (real geprüft: `grep -rn "docker push"` ohne
Treffer außerhalb der vendored Baseline). Die Archiv-Bedingung ist damit
nie erfüllt — ein committeter Digest adressiert ein Image, das der lokale
Docker-Cache jederzeit verwerfen kann (`docker system prune` o. ä.); dieses
Repo betreibt faktisch die Rezept-Form, ohne sie so zu benennen.

**(3) Warum der zweite Baseline-Grund ebenfalls nicht trägt — Modul 12.**
Derselbe Baseline-Abschnitt nennt einen zweiten Zweck: Ein Replay-Lauf
gegen ein altes Golden Set (Baseline-Regelwerk
`modul-12-replay-evaluierung.md`) braucht den *Image-Hash von damals* —
dafür muss die Datei über Commits hinweg nachschlagbar, also committet
sein. Dieses Repo führt kein Replay-Feature, kein Golden Set und keinen
`image_hash`-Slot eines Replay-Manifests (real geprüft: kein Treffer für
„replay"/„Golden Set"/„image_hash" im gesamten Repo außerhalb der
vendored Baseline). Der einzige Grund, der eine committete Historie
tatsächlich bräuchte, existiert hier nicht als Anwendungsfall.

**(4) Was vom committeten Artefakt übrig bleibt.** Ohne (2) und (3) trägt
nur noch [`ADR-0044`](0044-image-beleg-semantik.md)s eigener Zweck: „Digest
ist der Lauf-Beleg des letzten offiziellen `make image`-Laufs." Dafür
genügt eine lokal geschriebene, nicht getrackte Datei ebenso gut — kein
Gate, kein Sensor und kein Replay-Prozess liest diesen Wert je aus der
`git`-Historie zurück; `make image` schreibt ihn bei jedem Lauf neu.

## Entscheidung

Wir wählen: **`harness/image-hash.txt` bleibt ein von `make image`
geschriebenes Werkzeug-Ausgangsartefakt, wird aber nicht mehr committet.**

1. **`.gitignore` trägt `harness/image-hash.txt`.** Die Datei entsteht bei
   jedem `make image`-Lauf lokal neu (unverändertes Rezept), verlässt den
   Arbeitsbaum aber nie in Richtung `git add`/Commit.
2. **Einmalige Entfernung aus dem Index** (`git rm --cached
   harness/image-hash.txt`) — der aktuell committete Stand wird aus der
   Git-Nachverfolgung genommen; die Datei bleibt lokal auf der Platte
   bestehen, bis der nächste `make image`-Lauf sie überschreibt.
3. **Die welle-1-Regel aus [`ADR-0044`](0044-image-beleg-semantik.md)
   Entscheidungspunkt 3 entfällt ersatzlos.** Es gibt keinen
   „Digest-Commit" mehr, weder als Regelfall noch als entfallenden
   Sonderfall bei unverändertem Digest — die Datei wird grundsätzlich
   nicht committet.
4. **`make image` läuft unverändert vor jeder Closure eines Zugs, der
   Build-Kontext-Dateien ändert** ([`ADR-0085`](0085-build-kontext-ausnahme-test-only-zweck.md)
   §Referenz) — als realer Bau-Beleg innerhalb des Laufs, nur ohne den
   anschließenden Commit-Schritt.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun; der Digest bleibt committet | kein Eingriff | das real belegte Merge-Rauschen (`BEO-PGC/image-digest-nichtdeterminismus-erzeugt-merge-konflikt`) bleibt; der Beleg zeigt laut §Kontext (2) auf ein Image, das nie aufbewahrt wird — er ist inhaltlich bereits heute wertlos, das Rauschen erkauft nichts |
| B — Datei-Inhalt auf `sha256` des extrahierten Binaries umstellen (entspricht [`ADR-0044`](0044-image-beleg-semantik.md)s bereits verglichener, verworfener Option B) | deterministisch bei unverändertem `cmd/`/`internal/`/`gen/`/`proto/`-Inhalt — kein Diff für Züge, die diese Pfade nicht berühren | löst nur die Nichtdeterminismus-Hälfte, nicht die Grundfrage; bräuchte einen neuen Extraktionsschritt (Container-Export) in jedem betroffenen Lauf; die Datei bliebe committet und würde bei echten Code-Änderungen weiterhin einen (jetzt berechtigten) Merge-relevanten Diff erzeugen, den zwei parallele Branches ebenso treffen können |
| **C — lokal halten, nicht mehr committen (gewählt)** | beseitigt das Merge-Rauschen vollständig — keine committete Zeile mehr, die kollidieren könnte; der einzige tragende Zweck (Lauf-Beleg) braucht keine Persistenz über Commits hinweg; kleinster Eingriff (`.gitignore` + einmaliges `git rm --cached`) | die welle-1-Praxis „Digest committen" entfällt; ein späterer Leser findet den Digest eines vergangenen Laufs nicht mehr in der `git`-Historie (siehe Re-Evaluierungs-Trigger) |

**Fazit:** C. A lässt ein bereits wertloses Rauschen unangetastet; B löst
nur die halbe Frage und behält das Konfliktpotential.

## Konsequenzen

- Positiv: Kein Merge-Konflikt mehr auf `harness/image-hash.txt` — die
  Datei ist nach dieser ADR gar nicht mehr im Baum, den zwei Branches
  gemeinsam sehen.
- Positiv: kleinster Eingriff, keine neue Werkzeugkette, keine Änderung
  am `make image`-Rezept selbst.
- Negativ: Der historische Image-Hash eines vergangenen Commits ist nicht
  mehr aus der `git`-Historie ablesbar. Bewusst in Kauf genommen, weil
  dieser Anwendungsfall aktuell nirgends existiert (kein Push, kein
  Replay-Feature, §Kontext (2)/(3)).
- Hinweis: [`ADR-0084`](0084-sync-gate-fuer-generierte-artefakte.md)
  Entscheidungspunkt 4 („`harness/image-hash.txt` ist kein
  Sync-Gegenstand") bleibt unverändert richtig — eine nicht mehr getrackte
  Datei ist a fortiori kein Sync-Gegenstand; kein Supersede nötig.
- Folgepflicht (Implementer-Zug): `.gitignore`, `git rm --cached
  harness/image-hash.txt`, `harness/README.md` §Werkzeuge (`make
  image`-Zeile), `Makefile`-Kommentar der `image:`-Zeile.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| — | Kein Sensor: `make image` bleibt kein Gate ([`AGENTS.md`](../../../AGENTS.md) §3.6 greift nicht, wie schon [`ADR-0044`](0044-image-beleg-semantik.md)); ob die Datei committet ist, prüft niemand maschinell — es ist die `.gitignore`-Zeile selbst, die es verhindert | — |

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

Beobachtbarer Trigger: **ein Lauf-Zweig braucht einen historischen,
über Commits hinweg nachschlagbaren Image-Beleg** — etwa ein
`docker push`-Workflow (Archiv-Form wird real erfüllt) oder ein
Replay-Feature nach Modul 12 (ein `image_hash`-Slot eines
Replay-Manifests). Dann wird der Beleg als Folge-ADR mit `supersedes`
wieder gehoben — entweder als committeter Digest (Archiv-Form, sofern
gepusht wird) oder als committeter Binary-Hash (Option B dieser ADR bzw.
[`ADR-0044`](0044-image-beleg-semantik.md)s Option B, sofern kein Push
existiert). Sonst `permanent`.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-18 | Accepted — Anlass: `BEO-PGC/image-digest-nichtdeterminismus-erzeugt-merge-konflikt` (real aufgetreten bei `slice-nats-drittstream-example-go`); im Dialog mit dem Auftraggeber geprüft, dass weder die Archiv-Bedingung (kein `docker push`) noch die Replay-Bedingung (kein Modul-12-Feature in diesem Repo) für ein Committen sprechen — Entscheidung für die lokale, nicht getrackte Form | `BEO-PGC/image-digest-nichtdeterminismus-erzeugt-merge-konflikt/observation.md` |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0103` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
