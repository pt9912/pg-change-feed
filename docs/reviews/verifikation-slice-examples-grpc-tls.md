# Review-Report: Verifikation slice-examples-grpc-tls — 2026-10-05

**Review-Art:** Verifikation (Modul 11) — DoD- und Entscheidungs-Konformität, Plan gegen Code; Eingang ist der Review-Report `review-slice-examples-grpc-tls`

**Gegenstand:** Slice `slice-examples-grpc-tls`, Commits 71c6a886 (Runner-Pull), 5c614d7f (Slice), cac8e38b (Fixrunde)

**Skill:** `.claude/agents/verifier.md` @ 675246dd
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
- Review-Report `review-slice-examples-grpc-tls`
- [`ADR-0150`](../plan/adr/0150-tls-und-mehrfach-token.md)
- [`LH-FA-SST-011`](../../spec/lastenheft.md)

Der Report-Text stammt vom Verifier-Lauf (Rollenvorgabe: keine Report-Dateien); der Planner hat ihn hier angelegt.

**Eigene Sensor-Läufe:**

| Sensor | Ergebnis |
|---|---|
| `make gates` | Exit 0; docs-check 1686 Dateien, 0 Befunde; a-check 0 Befunde; übrige Gates grün |
| `make suchlauf-nachmessen PLAN=…` | Exit 0, 9 Zeilen stimmen (gemessen) |
| `make handbuch-public-doc-check` | Exit 0 |
| `make doc-commits`/`make doc-immutable` mit `RANGE` | Exit 0, 0 Befunde |
| `make examples-csharp` | Exit 0; frischer Lauf ohne Cache: GrpcClient.Tests 61 bestanden, 0 Fehler |
| `make examples-kotlin` | Exit 0; frischer Lauf, `:grpc-client:test` lief, Build erfolgreich |
| `make test` | Exit 0, Paket `examples/grpc-client` ok |

---

## Findings

Jedes Finding folgt dem §Output-Schema des Reviewer-Skills.

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| V-1 | MEDIUM | Das Handbuch sagt im Abschnitt zu TLS: gemessen, ohne Angabe ende jedes Beispiel mit Ausgang 1. Der Plan trägt: Ausgang 1 ist nur für Go gefahren; für C# und Kotlin aus dem Code hergeleitet. Verlangt: Aussage auf die Fehlermeldung zum nicht aufgebauten Kanal (für alle drei gemessen) begrenzen und Ausgang 1 nur für Go als gemessen führen, oder für C# und Kotlin nachmessen. | nicht erhoben | `docs/user/benutzerhandbuch.md`, Abschnitt zu TLS · Kurzzitat nicht erhoben | nicht erhoben | nicht erhoben |
| V-2 | LOW | Mutationen für Go und Kotlin (Optionswert, falsche Datei) sind vom Implementer übernommen; die Namensprüfung dort ist hergeleitet. | nicht erhoben | Plan `slice-examples-grpc-tls`, Mutationsbelege | nicht erhoben | nicht erhoben |
| V-3 | INFO | Offene DoD-Zeilen wie vor der Closure erwartet: Häkchen der Gate-Zeile, Review-Häkchen, Closure-Notiz, Register, Risiko-Ausgänge, Paarungen. | nicht erhoben | Plan `slice-examples-grpc-tls`, DoD | nicht erhoben | nicht erhoben |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Review F-1 (Namensprüfung): je Sprache ein Test mit Anker gleich Serverzertifikat und SAN `other.example` bei Verbindung an `127.0.0.1`, alle drei grün | behoben; die C#-Mutation (Namens-Bit entfernt, genau dieser Test rot) ist vom Implementer übernommen, nicht neu gefahren; Go und Kotlin delegieren den Namen an die Bibliothek, dort ist keine Mutation gefahren (hergeleitet) |
| Review F-5 (Kotlin, DER-Anker): Test `derAnchorIsAcceptedAndReachesTlsServer` exportiert das Zertifikat als DER und erreicht den TLS-Server | Widerlegung hält; ein Pfad, auf dem der Kanalbau nach bestandener Vorprüfung wirft, ist weder belegt noch ausgeschlossen (keine weitere Probe gefahren) |
| Review F-4 | als Grenze benannt (Handbuch, `examples/README.md`), hergeleitet, nicht gefahren |
| Review F-2, F-3, F-7 | behoben |
| Zertifikat im Repo | selbst gemessen: keines |

## Nicht neu gefahren (übernommen)

Realserver-Läufe je Sprache (TLS-Container, Fehlerfälle ohne Anker, Ausgang 2 bei defekter Datei): Implementer-Bericht, im Plan §3.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 1 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** nicht erhoben

## Verdikt

**Merge-blockierend:** nein — konform mit einem MEDIUM-Befund (V-1) und zwei Rest-Punkten. Gates und Tests sind eigenhändig grün gefahren; Review F-5 ist belegt widerlegt, F-1 ist behoben und gebunden. Die Begründung, V-1 nicht merge-blockierend zu führen, nennt der alte Report nicht (nicht erhoben).

**Übergabe:** V-1 an den Implementer bzw. Planner; Gate-Zeilen, Closure-Notiz und Register nach der Closure.
