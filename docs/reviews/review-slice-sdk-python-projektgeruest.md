# Review-Report: slice-sdk-python-projektgeruest — 2026-09-19

**Review-Art:** Code — geprüft gegen Plan
(`docs/plan/planning/in-progress/slice-sdk-python-projektgeruest.md`),
`ADR-0107` (Accepted) und `AGENTS.md` Hard Rules (Modul 10 §Drei
Review-Arten).

**Gegenstand:** Diff-Range `f3e10429..HEAD` (Welle-Eröffnung bis
Implementer-Commit), Slice `slice-sdk-python-projektgeruest`,
Welle `welle-sdk-python-lh-fa-sst-009`. Vier Commits: `1d922647`
(open→next, reiner Move), `5663406e` (Verantwortlich gesetzt, nur
Slice-Datei), `8cfb5fd1` (next→in-progress, reiner Move), `07a464eb`
(Inhalt: neuer Baum `sdks/python/`).

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09"
(vier repo-spezifische HIGH-Regeln, drei MEDIUM-Klassen).
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-19.

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-sdk-python-projektgeruest.md`
  (§1 Ziel/Abgrenzung, §2 DoD, §3 Plan inkl. Plan-Nachzug, §6 Risiken,
  §8 Sub-Area/Modus)
- `ADR-0107` (Accepted) — Festlegung 1/3/4, §Kontext „Was das ändert"
- `AGENTS.md` §3.1 (Docker-only), §3.3 (git mv + Inhalt = zwei Commits),
  §3.7 (Kommentar-Disziplin), §3.11 (kein host-lokaler Pfad), §3.13
  (Träger-Nachzug)
- `harness/conventions.md` (MR-000 ID-Schema, MR-002 Slice-Kennungen)
- `sdks/csharp/{Dockerfile,PgChangeFeed.Client/PgChangeFeedClientOptions.cs,README.md}`
  als bereits **geprüftes** Formvorbild (`docs/reviews/review-slice-sdk-csharp-projektgeruest.md`)
  — direkter Nebeneinander-Vergleich, nicht nur Erinnerung an das Muster

**Eigenständig durchgeführte Prüfungen (nicht nur Commit-Message
übernommen):**

- `docker build --no-cache -f sdks/python/Dockerfile -t
  pgcf-python-review-test sdks/python` real ausgeführt, Exit-Code ungepiped
  in eine Log-Datei umgeleitet und separat geprüft: `0`. `pytest`-Ausgabe:
  `3 passed in 0.01s` — bestätigt die Implementer-Behauptung „3/3 Tests
  grün" wortgleich. Test-Image danach mit `docker rmi` entfernt
  (`git status --short` davor/danach leer).
- `docker buildx imagetools inspect python:3.13-slim` real gegen die
  Registry ausgeführt: Index-Digest
  `sha256:8d9d0b8bcf6506481eae4907c18f5e3e7902e629f5f6d684f9e7c32e85e3ddf0`
  — identisch mit dem im Dockerfile-Kommentar genannten Digest. Der Pin
  trägt real noch das behauptete Tag.
- `make docs-check` real ausgeführt, Exit-Code direkt (ungepiped) geprüft:
  `0` (`d-check: 821 Datei(en) geprüft, 0 Befund(e)`).
- `git show 1d922647/8cfb5fd1 --stat`: beide zeigen ausschließlich den
  Rename `{open,next}⇒{next,in-progress}`, 0 Insertions/Deletions —
  reine Moves. `git show 5663406e --stat`: ändert ausschließlich die
  Slice-Datei (1 Zeile) — das 3-Commit-Muster (Move · Inhalt · Move) ist
  sauber getrennt, keine Move+Inhalt-Vermischung (`AGENTS.md` §3.3).
- `grep -rn "internal/\|cmd/pg-change-feed" sdks/python/` — kein Treffer
  (Exit 1).
- `git ls-files sdks/python` gelesen: `README.md` erscheint **genau
  einmal** im Git-Index (`sdks/python/README.md`); die zweite Kopie
  (`pgchangefeed/README.md`) entsteht ausschließlich im Docker-Bau-Layer
  über `COPY`, nicht im committeten Baum.
- `pyproject.toml` Feld für Feld gegen die DoD-Liste geprüft: `name`,
  `version = "0.1.0"`, `description`, `authors`, `license = "MIT"`
  (`LICENSE` im Repo-Root real gegengelesen — MIT, Copyright pt9912 2026),
  `[project.urls]` (Homepage/Repository auf dieses Repo), `readme`
  vorhanden; `dependencies = ["httpx>=0.27"]` ist die einzige
  Laufzeit-Fremdabhängigkeit (`grep -n dependencies` zeigt nur den einen
  Block plus `[project.optional-dependencies] test = ["pytest>=8"]` —
  eine Test-Extra, keine Laufzeitabhängigkeit).
- **README-Workaround-Vergleich (`sdks/csharp/Dockerfile` vs.
  `sdks/python/Dockerfile` nebeneinander gelesen):** Das C#-`.proto`-Muster
  kopiert eine Datei, die **außerhalb** des Bau-Kontexts liegt, über einen
  zusätzlichen, **benannten** Bau-Kontext (`--build-context proto=proto`)
  herein — ohne den Zusatzkontext bricht der Bau ab, kein stiller
  Fallback. Das Python-README-Muster kopiert dagegen eine Datei, die
  bereits **innerhalb** desselben Bau-Kontexts liegt (`sdks/python/`), nur
  von einer relativen Stelle an eine andere — ein einfaches `COPY
  README.md pgchangefeed/README.md` im selben, einzigen Kontext, kein
  zweiter benannter Kontext, keine Abbruch-an-fehlendem-Kontext-Semantik.
  Es ist dieselbe **Grundidee** („erwartete Datei per COPY an die vom
  Projekt verlangte relative Stelle bringen, ohne sie im committeten Baum
  zu duplizieren"), aber technisch das einfachere der beiden Muster, nicht
  identisch — der Plan-Nachzug benennt diesen Unterschied bereits korrekt
  („Nicht übertragbar: das .proto-Zusatzkontext-COPY-Muster … aber das
  Grundmuster … wurde real … wiederverwendet"). Kein Befund: Die Analogie
  ist ehrlich benannt, nicht als 1:1-Identität behauptet.
- `docker buildx imagetools inspect` bestätigt zusätzlich, dass
  `python:3.13-slim` als Manifest-Index-Digest gepinnt ist (alle
  Plattformen), wie im Plan-Nachzug §3 behauptet.
- Suchlauf nach host-lokalen absoluten Pfaden über `sdks/python/` und die
  Slice-Datei (`AGENTS.md` §3.11) — kein Treffer.
- `git diff f3e10429..HEAD --stat -- .a-check.yml harness/README.md spec/
  AGENTS.md harness/conventions.md docs/user/version.md
  docs/user/benutzerhandbuch.md docs/plan/adr/ .github/` — leer; bestätigt
  `ADR-0107` §6 „Was diese ADR nicht ändert" wortgleich für diesen Slice.
- `sdks/python/README.md` gegen relative Pfade geprüft (`grep -n
  '](\.\.' `/`'](/' `) — kein Treffer, ausschließlich absolute
  GitHub-Blob-URLs.
- Quervergleich `sdks/csharp/README.md` vs. `sdks/python/README.md`:
  Beide verwenden denselben etablierten, bereits geprüften Satzbaustein
  „NATS-vollinhalt(s) delivery remain … (`ADR-01XX` Festlegung 1)" — kein
  neuer Sprachbruch, sondern ein bereits akzeptiertes, bilinguales
  Zitierschema (`Festlegung N`, Feature-Eigenname `NATS-Vollinhalt`), kein
  Befund.
- `git status --short` nach allen eigenen Läufen leer.

---

## Findings

### F-1 — Slice-Chronik mit benannten Slice-IDs und Vorher/Nachher-Sprache im Produktionscode-Modul

- `kategorie`: HIGH
- `quelle`: Hard-Rule-Name — Skill-Regel „Slice-/Wellen-Chronik in
  Produktionscode-Kommentar" (`.harness/skills/reviewer.md` §Klassifikation
  HIGH)
- `pfad`: `sdks/python/pgchangefeed/src/pgchangefeed/__init__.py:1-6`
- `befund`: Der Modul-Docstring von `__init__.py` — Teil des shipped
  PyPI-Packages, nicht eines Dockerfile-Kopfkommentars — lautet: „Package
  skeleton (slice-sdk-python-projektgeruest, `ADR-0107` Festlegung 1/3/4).
  The HTTP API client surface itself … is added by the follow-up slice
  (slice-sdk-python-http-client-flaeche) — this release exposes only the
  shared connection configuration." Das Satzsubjekt ist der
  Produktionscode-Pfad selbst („Package skeleton … is added by …"), nicht
  ein Testfall — die in der Abgrenzung der Skill-Regel als zulässig
  genannte Testfall-Provenienz-Form greift hier nicht. Der Docstring
  benennt sowohl den aktuellen als auch einen konkreten künftigen
  Slice-Namen und begründet die Aussage damit explizit über
  Vorher/Nachher-Sprache („is added by the follow-up slice …") statt
  ausschließlich über `ADR-*`/`LH-*` oder den Herkunfts-Anker
  „· seit slice-<NNN>". Das direkte C#-Formvorbild
  (`sdks/csharp/PgChangeFeed.Client/PgChangeFeedClientOptions.cs:1-13`,
  bereits geprüft und akzeptiert) vermeidet exakt das: Es zitiert nur
  `ADR-0106 Festlegung 1` und spricht generisch von „the follow-up slices",
  ohne einen einzelnen Slice-Namen zu nennen. Der bereits akzeptierte
  Bare-Slice-Name im `sdks/csharp/Dockerfile`-Kopfkommentar ist kein
  Gegenbeleg: Dort steht kein Vorher/Nachher-Satz, sondern eine reine
  Herkunftsangabe im Bau-/Infrastruktur-Artefakt — hier dagegen ein
  Fließtext-Satz mit Zukunftsbezug in einer Bibliotheksdatei, die
  Konsumenten des Packages importieren.
- `verifizierbar`: nein — kein Gate liest Kommentar-Semantik (Skill-Text:
  „Kein Gate fängt das … repo-weiter Textmuster-Sensor geprüft und
  verworfen").
- `klasse`: Slice-Chronik in Produktionscode-Kommentar (fünftes benanntes
  Auftreten der Klasse, erste Instanz in einer Python-Sub-Area)

## Negativbefunde

- geprüft, ohne Befund: Docker-only-Disziplin (`AGENTS.md` §3.1) —
  `sdks/python/Dockerfile` trägt `pip install`/`pytest` ausschließlich im
  gepinnten Container; `.gitignore`-Kommentar benennt einen Host-Bau
  ausdrücklich nur als Schutz gegen ein Versehen.
- geprüft, ohne Befund: `git mv` + Inhaltsänderung als getrennte Commits
  (`AGENTS.md` §3.3) — siehe Prüfungsliste oben.
- geprüft, ohne Befund: Import-Grenze (`ADR-0107` §Entscheidung
  Festlegung 3, „verschärft") — kein Treffer für
  `internal/`/`cmd/pg-change-feed` unter `sdks/python/`; die einzige
  Fremdabhängigkeit ist `httpx`.
- geprüft, ohne Befund: `pyproject.toml`-Metadaten-Vollständigkeit gegen
  die DoD-Liste — jedes verlangte Feld einzeln geprüft, keine Lücke, keine
  weitere Laufzeit-Fremdabhängigkeit.
- geprüft, ohne Befund: doppelte README-Kopie im committeten Baum —
  `git ls-files sdks/python | grep -c README.md` = `1`; die Docker-Layer-
  Kopie entsteht nur im Bau, nicht im Git-Index.
- geprüft, ohne Befund: README-Workaround-Muster — technisch einfacher
  als, aber ehrlich als Analogie (nicht Identität) zum `.proto`-Muster in
  `sdks/csharp/Dockerfile` benannt; kein irreführender Vergleich.
- geprüft, ohne Befund: Basis-Image-Pin (`python:3.13-slim`) — real gegen
  die Registry verifiziert, Index-Digest stimmt mit dem Kommentar überein.
- geprüft, ohne Befund: `options.py`/`ClientOptions` — kein Vorgriff auf
  Endpunkt-Methoden, ausschließlich Adresse/Token, dieselbe minimale Form
  wie das C#-Vorbild; Konstruktions- und beide Negativtests
  (leere Adresse, leeres Token) vorhanden und binden real an die
  `__post_init__`-Validierung (nicht nur den Rückgabewert).
- geprüft, ohne Befund: Kommentar-Disziplin (`AGENTS.md` §3.7) im
  Dockerfile und in `pyproject.toml`/`options.py` — beschreiben
  durchgängig den geltenden Zustand oder eine Abgrenzung zu Folge-Slices,
  keine Konjunktiv-Begründung über eine verworfene Alternative, kein
  abwesender Text (Ausnahme: F-1, betrifft ausschließlich `__init__.py`).
- geprüft, ohne Befund: `README.md` dupliziert nicht die kanonische
  Draht-Doku, ist durchgängig Englisch, ausschließlich absolute
  GitHub-Blob-URLs; der bilinguale Baustein „NATS-vollinhalt(s) …
  Festlegung N" ist ein bereits etabliertes, akzeptiertes Zitierschema
  (identisch in `sdks/csharp/README.md`), kein neuer Sprachbruch.
- geprüft, ohne Befund: host-lokale absolute Pfade (`AGENTS.md` §3.11) —
  kein Treffer unter `sdks/python/` oder der Slice-Datei.
- geprüft, ohne Befund: Träger, die `ADR-0107` §6 ausdrücklich als
  „bleibt unberührt" benennt (`.a-check.yml`, `harness/README.md`,
  `spec/**`, `AGENTS.md`, `harness/conventions.md`,
  `docs/user/version.md`, `docs/user/benutzerhandbuch.md`,
  `docs/plan/adr/`, `.github/`) — tatsächlich in diesem Diff unangetastet.
- geprüft, ohne Befund: Scope-Treue gegen §1 „Ausdrücklich NICHT in
  diesem Slice" — kein `make sdk-pack-python`, kein Publish-Workflow, kein
  Umbau von `sdks/csharp/`, kein Pflichtenheft-Träger-Nachzug im Diff.
- geprüft, ohne Befund: Traceability — Commit-Betreff `07a464eb` nennt
  `LH-FA-SST-009` und `ADR-0107`, kein `SPEC-*`/`ARC-*` im Betreff;
  `harness/conventions.md` MR-002 (Slice-Kennungen sind Namen) —
  `slice-sdk-python-projektgeruest` ist bereits ein Name.
- geprüft, ohne Befund: reale Docker-Build- und Testausführung — 3/3
  Tests grün, Build ohne Fehler (einzige Meldungen: der übliche
  `pip`-Root-User-Hinweis und ein plattform-bedingter Buildx-Hinweis,
  identisch zum bereits bestehenden `sdks/csharp/Dockerfile`-Verhalten —
  kein Befund, kein neues Muster).
- geprüft, ohne Befund: `make gates`-Nachweis — nicht separat erneut
  gefahren (Review beschränkt sich auf `make docs-check` als engsten
  einschlägigen Sensor für einen reinen neuen-Sprach-Baum-Diff, plus den
  realen Docker-Build als Gegenprobe zur Implementer-Behauptung „3/3 Tests
  grün"); `make docs-check` real grün, siehe oben.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** Slice-Chronik in Produktionscode-Kommentar.

## Verdikt

**Merge-blockierend:** ja — 1 HIGH. Der Modul-Docstring in `__init__.py`
ist eine Fixrunde beim Implementer wert: die konkreten Slice-Namen und die
Vorher/Nachher-Formulierung durch eine ADR-*/LH-*-Referenz bzw. einen
`· seit slice-<NNN>`-Anker ersetzen, analog dem bereits akzeptierten
Formvorbild `PgChangeFeedClientOptions.cs`.

**Übergabe:** F-1 geht als Rückkante Review → Implementer für eine
Fixrunde. Kein Architect-Eskalationspfad nötig — ein einzelnes HIGH ohne
Rollen-Widerspruch, isoliert im selben Kontext behebbar (Skill: der
Konflikt-Pfad über den Architect greift erst „ab HIGH mit
Rollen-Widerspruch — oder ab dem dritten gleichen Konflikttyp"; dies ist
das erste Auftreten dieser Klasse in einer Python-Sub-Area). Da eine
Fixrunde folgt, bleibt die DoD-Checkbox „Review durchgeführt, Report
unter `docs/reviews/` liegt vor" im Slice-Plan **offen** (Skill-Regel
„DoD-Checkbox-Nachzug ohne Fixrunde" greift nur bei 0 HIGH). Die
**Finding-Klasse** geht zusätzlich in die Slice-Closure §7 und von dort in
den Steering-Loop-Zähler. Dieser Report ersetzt keine Verifikation gegen
die DoD — das bleibt Verifier-Aufgabe (Modul 11).
