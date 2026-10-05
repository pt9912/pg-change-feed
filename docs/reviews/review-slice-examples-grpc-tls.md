# Review-Report: slice-examples-grpc-tls — 2026-10-05

**Review-Art:** Code — gegen Plan und Hard Rules

**Gegenstand:** Diff 1ed7a246..5c614d7f (Commits 71c6a886 Runner-Pull, 5c614d7f Slice)

**Skill:** `.harness/skills/reviewer.md` @ 675246dd
**Modell:** Sonnet 5.5 · **Datum:** 2026-10-05

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

- Slice-Plan `slice-examples-grpc-tls`
- [`ADR-0150`](../plan/adr/0150-tls-und-mehrfach-token.md)
- [`LH-FA-SST-011`](../../spec/lastenheft.md)
- [`AGENTS.md`](../../AGENTS.md) §3.1, §3.7, §3.12, §3.13

Der Report-Text stammt vom Reviewer-Lauf; der Reviewer darf keine Report-Dateien anlegen, die Datei hat der Planner aus der Rückgabe angelegt. Die Mutationsproben hat der Reviewer nicht neu gefahren, sondern am Quelltext nachvollzogen.

---

## Findings

Jedes Finding folgt dem §Output-Schema des Reviewer-Skills.

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Alle Tests verbinden an `127.0.0.1` mit einem Zertifikat, dessen SAN `127.0.0.1` ist; kein Test verbindet an einen Namen außerhalb des SAN. In C# ist die Namensprüfung eigener Code (Callback, Bit `RemoteCertificateNameMismatch`); entfernt man das Bit aus der Maske, bleibt jeder Test grün (hergeleitet). Go und Kotlin delegieren den Namen an die Bibliothek, auch dort deckt ihn kein Test. Handbuch und Plan sagen den Namensvergleich zu. | nicht erhoben | `examples/csharp/grpc-client/ChannelFactory.cs` · `Validate`; Tests `ChannelFactoryTests.cs`, `examples/grpc-client/tls_test.go`, `ChannelFactoryTest.kt` | ja — Mutation an einer Kopie | Namensprüfung der TLS-Beispiele von keinem Test gebunden |
| F-2 | LOW | Der Handbuch-Kopf `Stand:` ist gegenüber der Historienzeile 1.99 nicht fortgeschrieben. | nicht erhoben | `docs/user/benutzerhandbuch.md` · Kopfzeile `Stand:` | nicht erhoben | Handbuch-Kopf nicht fortgeschrieben |
| F-3 | LOW | Chronik-Wortlaut in der Betreibersicht; Ist-Zustand: „ohne Angabe verbindet es im Klartext“. | nicht erhoben | `docs/user/benutzerhandbuch.md` · „im Klartext wie bisher“ und Historienzeile 1.99 | nicht erhoben | Chronik-Wortlaut in der Betreibersicht |
| F-4 | LOW | Die Kette wird mit `CustomRootTrust` ohne `ExtraStore` gebaut; nennt der Anker eine CA statt des Blattes und sendet der Server eine Zwischen-CA, scheitert der Aufbau (hergeleitet). Der Fehlerfall ist sicher (Ablehnung). | nicht erhoben | `examples/csharp/grpc-client/ChannelFactory.cs` · `Validate` | nicht erhoben | Zwischenzertifikate des Servers ignoriert |
| F-5 | INFO | Eine DER-Datei besteht die Vorprüfung und wirft danach im `build()` ohne Fang (hergeleitet). | nicht erhoben | `ChannelFactory.kt` · Kurzzitat nicht erhoben | nicht erhoben | nicht erhoben |
| F-6 | INFO | Mutationsbelege plausibel. C#: zwei Mutationen in einem Lauf, Zuordnung Test zu Mutation offen; Namensprüfung offen (F-1). | nicht erhoben | Plan `slice-examples-grpc-tls`, Mutationsbelege | nicht erhoben | nicht erhoben |
| F-7 | INFO | Exit 1 ohne Anker ist für C# und Kotlin nur im Code belegt, nicht in der berichteten Messung. Die Begründung für „keine Runner-Phase“ fehlt im Plan. | nicht erhoben | Plan `slice-examples-grpc-tls`, Realserver-Beleg · Kurzzitat nicht erhoben | nicht erhoben | nicht erhoben |
| F-8 | INFO | Die DoD-Zeile „Doku-Update“ steht vorzeitig auf erledigt; der Teil „gemeldete Träger fremder Dateien mit der Closure nachgezogen“ ist offen. | nicht erhoben | Plan `slice-examples-grpc-tls`, DoD-Zeile „Doku-Update“ | nicht erhoben | nicht erhoben |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Suchlauf (9 Zeilen stimmen) | geprüft, ohne Befund |
| Kommentar-Regeln | geprüft, ohne Befund |
| TLS-Sicherheit in Go (kein Überspringen der Prüfung) | geprüft, ohne Befund |
| TLS-Sicherheit in C# (kein `return true`) | geprüft, ohne Befund |
| TLS-Sicherheit in Kotlin | geprüft, ohne Befund |
| Tests ohne gespeicherte Zertifikate | geprüft, ohne Befund |
| `harness/README.md`-Zeilen | geprüft, ohne Befund |
| Runner-Commit (`docker image inspect` vor `docker pull`) | geprüft, ohne Befund |
| `make fmt-check` | geprüft, ohne Befund |
| Traceability | geprüft, ohne Befund |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 3 |
| INFO | 4 |

**Finding-Klassen dieses Laufs:** Namensprüfung der TLS-Beispiele von keinem Test gebunden · Handbuch-Kopf nicht fortgeschrieben · Chronik-Wortlaut in der Betreibersicht · Zwischenzertifikate des Servers ignoriert

## Verdikt

**Merge-blockierend:** ja — 0 HIGH, 1 MEDIUM; Fixrunde am Implementer: F-1 (Namensfehler-Test in allen drei Sprachen), F-2, F-3; F-4 und F-5 als benannte Grenze oder Fix nach Aufwand; F-7 Begründung im Plan nachtragen.

**Übergabe:** Findings gehen an den Implementer; die **Finding-Klassen** gehen zusätzlich in die Slice-Closure §7 und von dort in den Zähler. Dieser Report ist ein **Lauf-Beleg** und ersetzt keine Verifikation — DoD-/Spec-Konformität prüft der Verifier separat (Modul 11).
