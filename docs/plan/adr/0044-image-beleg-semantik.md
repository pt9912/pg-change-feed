# ADR-0044: Image-Beleg-Semantik (Digest ist lauf-gebunden)

**Status:** Accepted

**Datum:** 2026-09-09

**Autor:** pt9912

**Bezug:** [`ADR-0042`](0042-transport-typen-am-port.md) (Bindungs-Anker der
`make image`-Werkzeuge-Zeile); Belege:
[`docs/reviews/verify-slice-004.md`](../../../docs/reviews/verify-slice-004.md)
(F-2-Schiedsspruch),
[`docs/reviews/review-slice-005.md`](../../../docs/reviews/review-slice-005.md)
(F-7)

**Schärft:** — *(Prozess-ADR ohne Spec-Stratum; die Semantik des Belegs
tragen kein `SPEC-*`/`ARC-*`-Anker, sondern die Verkabelung — der
Dockerfile-Kopf und die `make image`-Werkzeuge-Zeile in
[`harness/README.md`](../../../harness/README.md), beide im selben Zug wie
diese ADR ersetzt. Wer diese ADR ändert, zieht die beiden Träger nach.)*

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

`harness/image-hash.txt` trägt seit dem ersten `make image`-Lauf
den Image-Digest, und die Regel aus welle-1 — ein Zug, der Build-Kontext-
Dateien ändert, erneuert den Beleg vor seiner Closure — hängt an seiner
Semantik. Commit `96c47af` (Planner-Deklaration, kein ADR) formulierte diese
Semantik als „genau dann"-Äquivalenz: Der Digest ändere sich **genau dann**,
wenn sich das exportierte Image ändert; deps-Layer-Änderungen ohne
Binary-Import ließen ihn unverändert (Dockerfile-Kopf
`Dockerfile:5-10`, Werkzeuge-Zeile in `harness/README.md`).

[`review-slice-005 F-7`](../../../docs/reviews/review-slice-005.md) widerlegt
die Formel am Verhalten des Ranges: Das Binary — mit exakt den
Dockerfile-Flags (`-trimpath -ldflags="-s -w"`, `CGO_ENABLED=0`) gebaut — ist
an beiden Range-Grenzen **bit-identisch** (sha256 `43c3aec0…`), und der
Digest wechselte trotzdem (`447eab36…` in `423cc5a` → `9ac4a9fb…` in
`cf3c57a`). Der F-2-Schiedsspruch des Verifiers
([`verify-slice-004.md`](../../../docs/reviews/verify-slice-004.md)) und die
Implementer-Probe zeigen dieselbe Richtung: Ein Re-Build am Aufzeichnungs-
Commit in einer anderen Umgebung lieferte `01b46e36…` statt des eingetragenen
`447eab36…` — der Digest ist **builder- und lauf-gebunden**; Provenance-
Metadaten und Builder-Umgebung variieren ihn bei byte-identischem Inhalt.
Als Inhalts-Adressierung über Läufe und Umgebungen hinweg ist der Digest
untauglich; die „genau dann"-Formel trifft den Mechanismus nicht.

Der Streit war im Implementer- und Review-Kontext unentschieden — die
Semantik-Korrektur ist eine Entscheidung (Zuständigkeit Architect, Modul 8),
und `96c47af` war eine Deklaration ohne Entscheidungsdokument. Diese ADR
holt die Entscheidung nach.

## Entscheidung

Wir wählen: **der Digest ist Lauf-Beleg, der Binary-Inhalt ist die
Entscheidungs-Größe.**

1. **`harness/image-hash.txt` ist der Digest des letzten offiziellen
   `make image`-Laufs** — ein Lauf-Beleg, kein Inhalts-Fingerabdruck. Er ist
   über Umgebungen und Läufe nicht inhaltlich reproduzierbar; die
   „genau dann"-Äquivalenz aus `96c47af` gilt nicht.
2. **Inhalts-Streits werden über den Binary-Hash entschieden** — der sha256
   des aus dem exportierten Image extrahierten Binaries (byte-identisch =
   inhaltlich identisch, unabhängig vom Digest).
3. **Die welle-1-Regel bleibt:** Ein Zug, der Build-Kontext-Dateien ändert,
   läuft `make image` vor seiner Closure; der Digest-Commit entfällt bei
   unverändertem Digest. Ein Digest-Vergleich über Umgebungen oder Läufe ist
   **kein Staleness-Beweis** — der Inhalt (Binary) ist die Entscheidungs-
   Größe; der unveränderte Digest innerhalb desselben Builder-Laufs bleibt
   ein gültiger Befund.
4. **Der Binary-Extraktions-Mechanismus ist der belegbare Weg:** Export des
   gebauten Images als Container, sha256 des extrahierten Binaries —
   dokumentiert als Muster in
   [`verify-slice-004.md`](../../../docs/reviews/verify-slice-004.md)
   (F-2-Schiedsspruch, Builds A/B/C mit `docker buildx build --load
   --metadata-file`, `/tmp`-Worktree) und
   [`review-slice-005.md`](../../../docs/reviews/review-slice-005.md)
   (F-7, Probe an beiden Range-Grenzen).

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — Digest-Wechsel als Frische-Marke | einfach: keine Extraktion; ein automatischer Staleness-Sensor könnte frischer Build gegen eingetragenen Digest vergleichen | **widerlegt** (F-7, F-2-Schiedsspruch): Binary bit-identisch, Digest wechselte (`447eab36…` → `9ac4a9fb…`); Re-Build am Aufzeichnungs-Commit lieferte `01b46e36…` statt `447eab36…` — der Vergleich über Umgebungen/Läufe entscheidet Staleness weder positiv noch negativ |
| B — Binary-Hash statt Digest in `image-hash.txt` | ein echter Inhalts-Fingerabdruck am HEAD; Vergleich über Umgebungen möglich | Wechsel des Datei-Vertrags: Der Digest ist die **Registry-Adresse** des exportierten Images (Archiv- und Push-Form, Kette zu `make image-stale`); der Binary-Hash ist die Inhalts-Größe, nicht die Adresse — der Datei-Vertrag würde beide Größen vermischen, und die bestehenden Modul-14-Beleg-Verweise zögen um |
| **C — Digest als Lauf-Beleg behalten, Inhalt über Binary-Extraktion prüfen (gewählt)** | der Datei-Vertrag bleibt (Digest = Registry-Adresse der Archiv-Form); die Inhalts-Frage bekommt ihre eigene Größe (Binary-sha256); keine Migration des Belegs, keine zweite Datei | der Lauf-Beleg ist über Umgebungen nicht inhaltlich vergleichbar; Inhalts-Streits verlangen den Extraktions-Schritt, der bisher nur als Proben-Muster in Review-/Verify-Läufen existiert (kein Make-Target) |

## Konsequenzen

- Positiv: Die Beleg-Semantik deckt sich mit dem beobachteten Verhalten;
  Staleness-Streitigkeiten entscheiden sich am Inhalt (Binary) statt an
  einem builder-gebundenen Digest; der Datei-Vertrag von
  `harness/image-hash.txt` bleibt stabil, bestehende Beleg-Verweise ziehen
  nicht um.
- Negativ: Der Digest allein beantwortet keine Inhaltsfrage mehr; ein
  Inhalts-Vergleich verlangt den Extraktions-Schritt (Container-Export +
  sha256), der noch kein Make-Target ist und bisher handgeführt als Probe
  lief.
- Folgepflicht: **Verkabelung** — der Dockerfile-Kopf und die
  `make image`-Werkzeuge-Zeile in `harness/README.md` tragen die korrigierte
  Semantik; die `96c47af`-Formel („ändert sich genau dann …") ist dort
  ersetzt. Soll ein Lauf-Zweig Inhalts-Vergleiche über Umgebungen führen,
  braucht es ein eigenes Binary-Beleg-Target (siehe
  Re-Evaluierungs-Trigger).

## Fitness Function (falls maschinell prüfbar)

Keine maschinell prüfbare Regel — der Beleg ist ein Werkzeug-Ausgang ohne
Gate (`make image` bleibt kein Gate, AGENTS §3.6 greift nicht); die
Semantik ist textlich verkörpert (Dockerfile-Kopf, Werkzeuge-Zeile), nicht
maschinell bewacht.

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

Beobachtbarer Trigger: **ein Lauf-Zweig braucht einen Inhalts-Vergleich über
Umgebungen oder Läufe** — etwa ein automatisierter Staleness- oder
Reproduzierbarkeits-Sensor, der zwei Builds am Inhalt statt am Lauf-Beleg
vergleichen muss. Dann wird der Binary-Hash als Beleg-Größe gehoben:
`harness/image-hash.txt` wird durch einen Binary-Beleg (sha256 des
extrahierten Binaries) ergänzt — als Folge-ADR mit `supersedes` (diese ADR
ist nach Accepted immutable). Sonst `permanent`.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-09 | Accepted — Anlass: [`review-slice-005 F-7`](../../../docs/reviews/review-slice-005.md) (Binary an beiden Range-Grenzen bit-identisch `43c3aec0…`, Digest wechselte `447eab36…` → `9ac4a9fb…`); korrigiert die `96c47af`-Deklaration („genau dann"-Äquivalenz), gestützt auf den [`F-2-Schiedsspruch`](../../../docs/reviews/verify-slice-004.md) | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
