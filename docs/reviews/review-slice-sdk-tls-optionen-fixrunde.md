# Review-Report: slice-sdk-tls-optionen (Fixrunde) — 2026-10-05

**Review-Art:** Code — gegen Plan, Entscheidungen und Hard Rules (Folgelauf zu `review-slice-sdk-tls-optionen`)

**Gegenstand:** Diff bd71fe7b..4e4a2a7e, die Commits fb4e8d70 (C#-Fix, Tests in drei Sprachen, READMEs) und 4e4a2a7e (Plan-Belege)

**Skill:** `.harness/skills/reviewer.md` (Stand Baseline v6.14.0)
**Modell:** Sonnet 5.5 · **Datum:** 2026-10-05

> **Zitier-Form.** Dieser Report friert ein; er zitiert Kennungen, nicht
> Adressen (`slice-<Kennung>`, `make <target>`, Baseline-Stellen als Tag +
> Pfad in Inline-Code).

**Eingangs-Kontext:**

- Slice-Plan `slice-sdk-tls-optionen`, Vorlauf-Report `review-slice-sdk-tls-optionen` (F-1 bis F-7)
- `LH-FA-SST-013`, `SPEC-037`, `ADR-0150`, `ADR-0145`
- `AGENTS.md` §3.1, §3.7, §3.12, §3.13

**Eigene Proben dieses Laufs:** `make suchlauf-nachmessen PLAN=<Plan>` (6 Zeilen stimmen, darunter `diff 38` zu „plaintext“), `make kommentar-kennungen DIFF=bd71fe7b PATHS=sdks` (kein Kandidat, Exit 0; Form, nicht Wahrheit), `make sdk-public-doc-check` (Exit 0). Gelesen und nachverfolgt: `TlsTransport.cs` im Ganzen (Reihenfolge der Zweige), die drei Test-Dateien, Fixture `TlsFixture.kt`, `tls.py`. **Nicht wiederholt, übernommen** aus dem Plan: die Unit-Läufe (254/252 Fälle), die drei Mutationsläufe, die Realserver-Läufe und der rote Lauf von `make test-sdk-kompat`.

---

## Vorlauf-Findings: Stand

| Vorlauf-ID | Stand | Beleg |
|---|---|---|
| F-1 (Prozess, Umleitung) | unverändert, nicht Gegenstand dieser Runde | Beobachtung geht an den Planner (Vorlauf-Verdikt); kein Diff-Anteil |
| F-2 (C# Aussteller-Zertifikat als Anker) | behoben | siehe Prüfpunkt 1 und 2 |
| F-3 (Systemanker-Ausschluss) | zu einem Teil gebunden, Rest benannt (N-2) | siehe Prüfpunkt 3 |
| F-4 (Kompat-Messung) | offen, nicht gemessen; Werkzeug-Drift (N-3) | siehe Prüfpunkt 5 |
| F-5, F-6, F-7 | F-6 für Kotlin-Gegenlauf geschlossen, Rest unverändert INFO | Plan-Beleg `SDK_TLS_CA_FILE=/tls/other.pem` rot |

## Findings (neu in dieser Runde)

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| N-1 | LOW | Die Systemanker-Zusage („mit Anker gelten die Anker des Betriebssystems nicht zusätzlich“) ist nur in Python HTTP und SSE an der Eingabeseite gebunden (`SSL_CERT_FILE` auf das Server-Zertifikat, Kontrolle ohne Anker verbindet, fremder Anker scheitert; Mutation `load_default_certs()` rot). In C#, Kotlin und Python gRPC bleibt sie ungebunden: eine Umstellung auf `TrustMode.System` bzw. den Standard-Speicher bleibt grün. Der Plan benennt die Grenze und ihren Grund (prozessweiter Wurzelspeicher), die Aussage über den Code ist als *hergeleitet* gekennzeichnet; das ist die ehrliche Form, kein Scheintest. Einstufung vom Vorlauf-MEDIUM auf LOW wegen der neuen Evidenz (benannte, begründete Grenze, eine Fläche gebunden), nicht wegen eines Widerspruchs. | `SPEC-037`; `AGENTS.md` §3.12 | `docs/plan/planning/in-progress/slice-sdk-tls-optionen.md` · `Nicht gebunden bleiben:` | ja — ein Kindprozess-Test je Fall würde die Mutation färben | Systemanker-Ausschluss in drei Flächen ungebunden, benannt |
| N-2 | LOW | Die Upgrading-Einträge von C# und Kotlin tragen die Überschrift `0.7.0`, die Package-Versionen stehen auf `0.6.1`, und die Tags `…-v0.6.1` existieren. Der Text beschreibt den Inhalt, den die nächste Minor-Version tragen wird (laut `SPEC-037`, im Plan so benannt); bis zum Anheben der Version steht in einem gebauten `0.6.1`-Paket eine README, die eine Version nennt, die es nicht gibt. Inhaltlich sind die Einträge wahr (C#: Konstruktor, Typwechsel `Internal` → `UnexpectedStatus`, `PgChangeFeedGrpcException`-Fänger unberührt — verifiziert, beide Typen leiten von `PgChangeFeedGrpcException` ab; Kotlin: `https` spricht jetzt TLS), englisch, ohne interne Kennung, in der Form der Vorgänger-Einträge; `make sdk-public-doc-check` Exit 0. Das Anheben der Versionen vor dem Veröffentlichen ist die Auflösung (Vorlauf F-5). | `AGENTS.md` §3.12 (Aussage und Stand); `SPEC-037` | `sdks/csharp/README.md` · `- **0.7.0** — \`PgChangeFeedClientOptions\` gained the optional third argument` | nein — Release-Zug außerhalb des Diffs | Versionsstand vs. Inhalt |
| N-3 | LOW | `make test-sdk-kompat` im Standardmodus endet rot: `tools/harness/sdk-kompat/csharp/Dockerfile` löst `-p:PgcfVersion=0.6.0` auf, `sdks/csharp/dist` trägt seit dem Release-Commit c520f3cd `0.6.1` (Plan: `NU1603 … 0.6.1 was resolved instead`; **übernommen**, nicht gefahren, das Ziel braucht Netz). Die Ursache liegt vor diesem Diff, der Plan sagt das und führt die Messung als „nicht gemessen“; er benennt aber keine Adresse für die Reparatur. Failure-Szenario: der Release-Zug fährt das Ziel, um die Verhaltensänderungen (Typwechsel in C#, `https` → TLS in Kotlin) zu messen, und bekommt Rot an der Auflösung statt an der Sache; die Version 0.7.0 bricht denselben Mechanismus erneut. Das Ziel liest die Fehlertypen, nicht dieses Mapping, deshalb wäre ein grüner Lauf ohnehin keine Messung der Änderung. Der Träger liegt in fremden Dateien (Werkzeug), also ist es eine Meldung mit Frist (`AGENTS.md` §3.13) und keine stille Mitänderung. | `AGENTS.md` §3.13; Skill „Aufschub-Adresse“ | `tools/harness/sdk-kompat/csharp/Dockerfile` · `-p:PgcfVersion=0.6.0` | ja — `make test-sdk-kompat` | Werkzeug-Drift nach Versionssprung |
| N-4 | INFO | Im Gleichheitszweig von `TlsTransport.cs` ist die Untergrenze `NotBefore <= now` durch keine Mutation gebunden (der gemeldete Ablauf-Fall deckt `now <= NotAfter`); ein Test mit einem noch nicht gültigen Anker-Blatt fehlt. Die Mutation, die `return true` im Zweig setzt, färbt nur den Ablauf-Fall. | Maintainability | `sdks/csharp/PgChangeFeed.Client/TlsTransport.cs` · `return leaf.NotBefore <= now && now <= leaf.NotAfter;` | ja — Fall mit `notBefore` in der Zukunft | Mutationsbeleg einseitig |

## Prüfpunkte der Anfrage

1. **C# `TlsTransport.cs`, Gleichheitszweig.** Reihenfolge nachgefahren: `certificate is null` → Rückweisung; dann Name/Verfügbarkeit aus `SslPolicyErrors` (`RemoteCertificateNameMismatch | RemoteCertificateNotAvailable` → false); erst danach `IsAnchor`. Die Gleichheit kann die Namensprüfung also nicht umgehen; der Test `ServerNameNotInCertificate_FailsTheConnection` hält das an der Eingabeseite (Anker ist das Blatt selbst, Adresse nennt die IP; würde die Namensprüfung hinter den Zweig rücken, wäre er rot). Maskiert werden bewusst nur die Chain-Fehler (`RemoteCertificateChainErrors`), weil die Kette neu gebaut wird; im Gleichheitszweig entfällt der Bau, was bei byte-gleichem Anker die Pinning-Semantik ist (der Handshake beweist den Besitz des Schlüssels). Gültigkeit: `X509Certificate2.NotBefore/NotAfter` sind Ortszeit, verglichen mit `DateTime.Now` — konsistent, kein Zeitzonenfehler. Revocation ist in beiden Zweigen aus (`NoCheck`), unverändert. Mehrere Anker: die Schleife prüft jeden, Test `AnchorFileWithSeveralCertificates_…` deckt es. Leere Anker-Datei: der Konstruktor wirft beim Erzeugen („contains no PEM certificate“), erreicht den Handler nicht. Ein Zertifikat ohne gültigen Zeitraum-Auswertung gibt es nicht; N-4 nennt die offene Untergrenze. Ohne Befund bis auf N-4.
   Mutationsbelege (3, **übernommen**): Eingabeseite ist das präsentierte Zertifikat; Zahlen konsistent (254 = 246 + 8 neue Fälle, 4 Flächen je Fall); jede Mutation färbt die zum Zweig passenden Tests, die dritte („jeder Anker gleich“) die Negativfälle. Plausibel und an der Eingabe gebunden.
2. **Neue Tests in drei Sprachen.** C# (Erfolg: ausgestelltes Blatt als alleiniger Anker; Fehler: Blatt eines anderen Servers mit gleichem Namen), Kotlin (`leafPem` als dritter Konstruktor-Parameter mit Standard `pem`, nur für die ausgestellte Variante auf `signed` gesetzt — bestehende Fälle unberührt) und Python (`server_cert_path` ist in der Fixture-Klasse der ausgestellten Variante definiert) sind symmetrisch aufgebaut; die Fehlerfälle benutzen ein anderes Blatt gleichen Namens, scheitern also an der Kette und nicht am Namen. Ohne Befund.
3. **F-3.** Siehe N-1. Der Python-Test ist valide: der Kontext ohne Anker (`httpx` mit Standard-Verifikation) liest `SSL_CERT_FILE`, die Kontrolle (`use` ohne Anker) wird im Test selbst gefahren, danach scheitert dieselbe Verbindung mit fremdem Anker; `SSL_CERT_DIR` ist entfernt. Die Mutation ist eingabeseitig (zusätzlicher Standard-Speicher).
4. **README-Einträge.** Siehe N-2.
5. **`make test-sdk-kompat`.** Siehe N-3.
6. **Suchlauf und Kommentare.** `make suchlauf-nachmessen` stimmt (der geänderte Sollwert 38 ist gelesen im Plan begründet: die 38. Zeile ist der Kotlin-Upgrading-Eintrag); `make kommentar-kennungen DIFF=…` ohne Kandidat; die neuen Kommentare (`TlsTransport.cs`, Python-Test) tragen Zustand, keine Chronik und keine Kennung.

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `sdks/csharp/.../TlsTransport.cs` | geprüft, ohne Befund außer N-4 |
| C#-, Kotlin-, Python-Tests und Kotlin-Fixture | geprüft, ohne Befund |
| `sdks/*/README.md` (Upgrading) | geprüft, Befund N-2 (Versionsstand), sonst ohne Befund |
| `docs/plan/planning/in-progress/slice-sdk-tls-optionen.md` | geprüft, ohne Befund außer N-1/N-3 — Nenner und Läufe als gemessen/übernommen/hergeleitet erkennbar, Suchlauf-Block nachgemessen |
| Prozess dieser Runde (Umleitungen, Host-Werkzeuge) | aus dem Diff nicht ablesbar, kein Hinweis; Vorlauf-F-1 unverändert beim Planner |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 3 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Systemanker-Ausschluss in drei Flächen ungebunden, benannt · Versionsstand vs. Inhalt · Werkzeug-Drift nach Versionssprung · Mutationsbeleg einseitig

## Verdikt

**Merge-blockierend:** nein. Der Vorlauf-Fund F-2 ist behoben (Fix, drei Mutationen, Tests in drei Sprachen); F-3 ist zur benannten, begründeten Grenze geworden (N-1); der HIGH-Prozessfund F-1 war ein Prozessverstoß ohne Diff-Anteil und liegt beim Planner, er verlangt keine Rückgabe an den Implementer. Keine weitere Fixrunde. Die DoD-Zeile „Review durchgeführt“ ist im Plan nachgezogen.

**Übergabe:** N-3 an den Planner (Werkzeug-Pflege `tools/harness/sdk-kompat/` vor dem Release-Zug, mit Adresse und Frist = Closure des Slice oder Folge-Slice); N-2 an den Release-Zug (Versionen anheben, vor jeder Veröffentlichung); N-1 und N-4 sind Hinweise ohne Aktion. Der Report ersetzt keine Verifikation; DoD-/Spec-Konformität prüft der Verifier.
