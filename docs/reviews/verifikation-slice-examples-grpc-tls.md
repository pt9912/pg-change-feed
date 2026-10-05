# Verifikation slice-examples-grpc-tls (Modul 11)

**Gegenstand:** [Slice-Plan](../plan/planning/in-progress/slice-examples-grpc-tls.md), Commits 71c6a886 (Runner-Pull), 5c614d7f (Slice), cac8e38b (Fixrunde), gegen [ADR-0150](../plan/adr/0150-tls-und-mehrfach-token.md) und [LH-FA-SST-011](../../spec/lastenheft.md); Eingang ist der [Review-Report](review-slice-examples-grpc-tls.md).

Der Report-Text stammt vom Verifier-Lauf (Rollenvorgabe: keine Report-Dateien); der Planner hat ihn hier angelegt.

## Verdikt

Konform mit einem MEDIUM-Befund (V-1) und zwei Rest-Punkten. Gates und Tests sind eigenhändig grün gefahren; F-5 ist belegt widerlegt, F-1 ist behoben und gebunden.

## 1. Eigene Sensor-Läufe

| Sensor | Ergebnis |
|---|---|
| `make gates` | Exit 0; docs-check 1686 Dateien, 0 Befunde; a-check 0 Befunde; übrige Gates grün |
| `make suchlauf-nachmessen PLAN=…` | Exit 0, 9 Zeilen stimmen (gemessen) |
| `make handbuch-public-doc-check` | Exit 0 |
| `make doc-commits`/`make doc-immutable` mit `RANGE` | Exit 0, 0 Befunde |
| `make examples-csharp` | Exit 0; frischer Lauf ohne Cache: GrpcClient.Tests 61 bestanden, 0 Fehler |
| `make examples-kotlin` | Exit 0; frischer Lauf, `:grpc-client:test` lief, Build erfolgreich |
| `make test` | Exit 0, Paket `examples/grpc-client` ok |

## 2. Review-Findings

- **F-1 (Namensprüfung):** behoben. Je Sprache ein Test mit Anker gleich Serverzertifikat und SAN `other.example` bei Verbindung an `127.0.0.1`; alle drei laufen grün. Die C#-Mutation (Namens-Bit entfernt, genau dieser Test rot) ist vom Implementer übernommen, nicht neu gefahren; Go und Kotlin delegieren den Namen an die Bibliothek, dort ist keine Mutation gefahren (hergeleitet).
- **F-5 (Kotlin, DER-Anker):** Widerlegung hält. Der Test `derAnchorIsAcceptedAndReachesTlsServer` exportiert das Zertifikat als DER und erreicht den TLS-Server. Ein Pfad, auf dem der Kanalbau nach bestandener Vorprüfung wirft, ist weder belegt noch ausgeschlossen (keine weitere Probe gefahren).
- **F-4:** als Grenze benannt (Handbuch, `examples/README.md`), hergeleitet, nicht gefahren.
- **F-2, F-3, F-7:** behoben.

## 3. Befunde der Verifikation

### V-1 MEDIUM — Handbuch gibt eine nur für Go gemessene Aussage als gemessen für alle drei aus
- Das Handbuch sagt im Abschnitt zu TLS: gemessen, ohne Angabe ende jedes Beispiel mit Ausgang 1. Der Plan trägt: Ausgang 1 ist nur für Go gefahren; für C# und Kotlin aus dem Code hergeleitet.
- Verlangt: Aussage auf die Fehlermeldung zum nicht aufgebauten Kanal (für alle drei gemessen) begrenzen und Ausgang 1 nur für Go als gemessen führen, oder für C# und Kotlin nachmessen.

### V-2 LOW — Mutationsprobe je Sprache
Mutationen für Go und Kotlin (Optionswert, falsche Datei) sind vom Implementer übernommen; die Namensprüfung dort ist hergeleitet.

### V-3 INFO — Offene DoD-Zeilen
Offen wie vor der Closure erwartet: Häkchen der Gate-Zeile, Review-Häkchen, Closure-Notiz, Register, Risiko-Ausgänge, Paarungen.

## 4. Nicht neu gefahren (übernommen)

Realserver-Läufe je Sprache (TLS-Container, Fehlerfälle ohne Anker, Ausgang 2 bei defekter Datei): Implementer-Bericht, im Plan §3. Selbst gemessen: kein Zertifikat im Repo.
