# Review-Report: slice-examples-grpc-tls — 2026-10-05

**Review-Art:** Code (gegen Plan, Hard Rules)
**Gegenstand:** Diff 1ed7a246..5c614d7f (Commits 71c6a886 Runner-Pull, 5c614d7f Slice)
**Skill:** `.harness/skills/reviewer.md` · **Modell:** Sonnet 5.5

**Eingangs-Kontext:**
- Plan [`slice-examples-grpc-tls`](../plan/planning/in-progress/slice-examples-grpc-tls.md)
- [`ADR-0150`](../plan/adr/0150-tls-und-mehrfach-token.md)
- [`LH-FA-SST-011`](../../spec/lastenheft.md)
- [`AGENTS.md`](../../AGENTS.md) §3.1, §3.7, §3.12, §3.13

Der Report-Text stammt vom Reviewer-Lauf; der Reviewer darf keine Report-Dateien anlegen, die Datei hat der Planner aus der Rückgabe angelegt. Die Mutationsproben hat der Reviewer nicht neu gefahren, sondern am Quelltext nachvollzogen.

## Verdikt

0 HIGH, 1 MEDIUM, 3 LOW, 4 INFO. Fixrunde am Implementer für das MEDIUM.

## Findings

### F-1 MEDIUM — Namensprüfung der TLS-Beispiele ist von keinem Test gebunden
- Pfad: `examples/csharp/grpc-client/ChannelFactory.cs` (`Validate`), Tests `ChannelFactoryTests.cs`, `examples/grpc-client/tls_test.go`, `ChannelFactoryTest.kt`
- Befund: Alle Tests verbinden an `127.0.0.1` mit einem Zertifikat, dessen SAN `127.0.0.1` ist. Kein Test verbindet an einen Namen außerhalb des SAN. In C# ist die Namensprüfung eigener Code (Callback, Bit `RemoteCertificateNameMismatch`); entfernt man das Bit aus der Maske, bleibt jeder Test grün (hergeleitet). Go und Kotlin delegieren den Namen an die Bibliothek, auch dort deckt ihn kein Test. Handbuch und Plan sagen den Namensvergleich zu.
- Verifizierbar: ja (Mutation an einer Kopie).

### F-2 LOW — Handbuch-Kopf `Stand:` nicht fortgeschrieben
- Pfad: `docs/user/benutzerhandbuch.md`, Kopfzeile `Stand:` gegenüber der Historienzeile 1.99.

### F-3 LOW — Chronik-Wortlaut in der Betreibersicht
- Pfad: `docs/user/benutzerhandbuch.md`, „im Klartext wie bisher“ und Historienzeile 1.99. Ist-Zustand: „ohne Angabe verbindet es im Klartext“.

### F-4 LOW — C#-Prüfung ignoriert vom Server gesendete Zwischenzertifikate
- Pfad: `examples/csharp/grpc-client/ChannelFactory.cs` (`Validate`)
- Befund: Die Kette wird mit `CustomRootTrust` ohne `ExtraStore` gebaut; nennt der Anker eine CA statt des Blattes und sendet der Server eine Zwischen-CA, scheitert der Aufbau (hergeleitet). Der Fehlerfall ist sicher (Ablehnung).

### F-5 INFO — Kotlin: DER-Zertifikat und ungefangener Pfad
- Pfad: `ChannelFactory.kt`
- Befund: Eine DER-Datei besteht die Vorprüfung und wirft danach im `build()` ohne Fang (hergeleitet).

### F-6 INFO — Mutationsbelege plausibel
- C#: zwei Mutationen in einem Lauf, Zuordnung Test zu Mutation offen. Namensprüfung offen (F-1).

### F-7 INFO — Realserver-Beleg und Handbuch-Messung
- Exit 1 ohne Anker ist für C# und Kotlin nur im Code belegt, nicht in der berichteten Messung. Die Begründung für „keine Runner-Phase“ fehlt im Plan.

### F-8 INFO — DoD-Zeile „Doku-Update“ steht vorzeitig auf erledigt
- Der Teil „gemeldete Träger fremder Dateien mit der Closure nachgezogen“ ist offen.

## Geprüft, ohne Befund

Suchlauf (9 Zeilen stimmen), Kommentar-Regeln, TLS-Sicherheit in Go (kein Überspringen der Prüfung), C# (kein `return true`), Kotlin, Tests ohne gespeicherte Zertifikate, harness/README.md-Zeilen, Runner-Commit (`docker image inspect` vor `docker pull`), `make fmt-check`, Traceability.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 3 |
| INFO | 4 |

**Verdikt:** Fixrunde: F-1 (Namensfehler-Test in allen drei Sprachen), F-2, F-3; F-4 und F-5 als benannte Grenze oder Fix nach Aufwand; F-7 Begründung im Plan nachtragen.
