# Review-Report: slice-beispiele-compose-bootstrap — 2026-09-18

**Review-Art:** Code — Diff gegen Plan/`ADR-0098`/Konventionen (Modul 10 §Drei
Review-Arten), nicht gegen DoD (Verifier-Aufgabe).

**Gegenstand:** Commit `309d39a`, Slice `slice-beispiele-compose-bootstrap`,
Welle `welle-beispiele-start-ueber-make`.

**Skill:** `.harness/skills/reviewer.md`
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-18

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-beispiele-compose-bootstrap.md` (Plan)
- `docs/plan/adr/0098-beispiel-clients-start-ueber-make-dockerfile.md`
  (Festlegung 3 Umgebungsdatei-Kontrakt, Festlegung 4 festes Netzwerk)
- `docs/plan/planning/done/slice-beispiele-go-dockerfile-start.md` +
  dessen Review-Report (Stilreferenz)
- `AGENTS.md` §3.1, §3.5, §3.7, §3.9, §3.11-§3.13
- `docs/plan/planning/observations/README.md`
- Eigene Messung: `make example-demo-up` real ausgeführt (Bootstrap-Log
  bestätigt Reihenfolge Schema-Rollout → Quelle/Tabelle → Feed-Start →
  Slot-Poll → Demo-Zeile), `GET /changes?source=demo-source` über einen
  Wegwerf-`curl`-Container gegen `cdc-examples` (reale Antwort mit
  Demo-Zeile), zweiter `make example-demo-up`-Lauf (Exit 0, keine
  Duplikate, `count(*)=1`), Root-`compose.yaml` parallel hochgefahren
  (kein Netzwerk-/Container-Namen-Konflikt), `make example-demo-down`
  (vollständiger Abbau, `docker ps -a`/`docker network ls` geprüft),
  `make example-run-go SURFACE=http ...` gegen die Demo-Umgebung (reale
  Antwort), `make gates` (ungepiped Exit 0).

---

## Findings

Keine HIGH-, keine MEDIUM-Findings.

### F-1 — Kopfkommentar in `examples/.env` lässt die ADR-interne Klammer-Referenz weg

- `kategorie`: LOW
- `quelle`: `ADR-0098` Festlegung 3
- `pfad`: `examples/.env:1-2`
- `befund`: Der ADR-Text lautet wörtlich „… im isolierten Docker-Netzwerk
  `cdc-examples` **(Festlegung 4)** — niemals …"; die committete Datei
  trägt denselben Satz ohne den Klammer-Verweis „(Festlegung 4)".
  Inhaltlich (Netzwerk-Name, Isolationsaussage, Verwendungsverbot) ist der
  Kommentar vollständig und stimmt mit dem übrigen ADR-Zitat überein — der
  weggelassene Teil ist ein ADR-interner Querverweis, der in einer
  `.env`-Datei ohnehin ohne Kontext stünde.
- `verifizierbar`: ja (Diff gegen ADR-Zitat)
- `klasse`: „Zitat nicht wortgleich, Inhalt unverändert" — kein
  Fixrunden-Bedarf.

## Negativbefunde

- geprüft, ohne Befund: `examples/compose.yaml` — Netzwerkname
  `cdc-examples` exakt (kein Typo, gegen `harness/mk/examples.mk`
  gegengeprüft), kein `build:`-Block im `pg-change-feed`-Service,
  `../tools/schema/compose-init`-Wiederverwendung statt SQL-Duplikation,
  Digest-Pins identisch zur Wurzel-`compose.yaml`.
- geprüft, ohne Befund: `examples/bootstrap.sh` — reale Ausführungsreihenfolge
  (Tabelle leer → Feed/Slot-Start → danach erst Demo-Zeile) deckt sich mit
  dem behaupteten Fix des Snapshot-Ordnungsfehlers; Idempotenz real durch
  zweiten Lauf bestätigt (Exit 0, keine Duplikate).
- geprüft, ohne Befund: `examples/.env` — nur Demo-Werte (dieselbe Klasse
  wie Wurzel-`compose.yaml`), Compose-Servicenamen statt `localhost`, alle
  Variablennamen einzeln gegen `docs/user/benutzerhandbuch.md` §5
  bestätigt.
- geprüft, ohne Befund: `harness/mk/examples.mk` —
  `example-demo-up`/`example-demo-down` real funktionsfähig, `include
  harness/mk/*.mk` deckt die Datei ab.
- geprüft, ohne Befund: Doku-Nachzug `examples/README.md`
  §Demo-Umgebung, `harness/README.md` §Werkzeuge — inhaltlich korrekt,
  Stil passend zu `examples-csharp`/`examples-kotlin`/`example-run-go`,
  „kein Gate" korrekt markiert.
- geprüft, ohne Befund: reale Ende-zu-Ende-Kette (Hochfahren, `GET
  /changes`, `example-run-go`, Parallelbetrieb mit Root-`compose.yaml`,
  Idempotenz, Abbau) — alle Implementer-Behauptungen selbst nachgefahren,
  keine Abweichung gefunden.
- geprüft, ohne Befund: `AGENTS.md` §3.7 (Kommentarklassen), §3.11
  (host-lokale Pfade), §3.13 (Träger-Nachzug) — keine Verstöße.
- geprüft, ohne Befund: Traceability — Commit-Message nennt `ADR-0098` und
  `LH-QA-OPS-001`.
- geprüft, ohne Befund: `097c649` (Portabilitäts-Fix, unmittelbar davor
  gelandet) — real unabhängig von diesem Slice, keine Dateiüberschneidung.
- geprüft, ohne Befund: `make gates` — 6/6 Gates grün, Exit-Code ungepiped
  direkt gelesen; Zahlen im Closure-Bericht (701 Dateien/0 Befunde,
  83.40 %) decken sich mit der eigenen Messung.

## INFO

- **INFO-1** — Die Bootstrapping-Idempotenz-Wache (`to_regclass`-Existenz-
  Check in `examples/bootstrap.sh`) erkennt nur „migriert ja/nein", nicht
  „migriert auf welchem Schema-Stand". Bereits vom Implementer ehrlich als
  Grenze benannt; der reguläre Reset-Pfad (`make example-demo-down`
  entfernt Volumes) trägt das ausreichend. Kein Fixrunden-Bedarf.
- **INFO-2** — `BEO-PGC/schema-rollout-fremdobjekte` erreicht mit diesem
  Slice real die 3×-Schwelle (Zähler nachgezählt: `slice-016`,
  `slice-063`, `slice-beispiele-compose-bootstrap`). Korrekt als „Ausgang
  fällig bei der nächsten Welle-Closure" markiert, nicht selbst
  entschieden — der Lese-Schritt bei der Closure von
  `welle-beispiele-start-ueber-make` muss diese Beobachtung auflösen oder
  bewusst vertagen.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Zitat nicht wortgleich, Inhalt
unverändert · bekannte, offen benannte Grenze einer Workaround-Wache ·
3×-Schwelle real erreicht — Lese-Schritt fällig

## Verdikt

**Merge-blockierend:** nein — 0 HIGH, 0 MEDIUM. Keine Fixrunde nötig.

Die DoD-Checkbox „Review durchgeführt, Report unter `docs/reviews/` liegt
vor" ist mit diesem Report erfüllt. Der Lese-Schritt zu
`BEO-PGC/schema-rollout-fremdobjekte` (INFO-2) wird an die Welle-Closure
von `welle-beispiele-start-ueber-make` weitergegeben, nicht hier
vorweggenommen.

Dieser Report ersetzt keine Verifikation — DoD-/Spec-Konformität prüft der
Verifier separat (Modul 11), sofern dieser Slice eine eigene Verifikation
vorsieht.
