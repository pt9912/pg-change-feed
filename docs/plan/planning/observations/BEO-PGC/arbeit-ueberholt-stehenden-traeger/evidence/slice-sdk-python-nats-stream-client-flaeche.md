# Beleg: slice-sdk-python-nats-stream-client-flaeche

Vorgang: `slice-sdk-python-nats-stream-client-flaeche` liefert die
öffentliche NATS-Vollinhalts-Client-Fläche (`SPEC-024`), hebt die
Package-Version auf `0.2.0` und bündelt den Träger-Nachzug für alle drei
neuen Flächen (`AGENTS.md` §3.13).

Fund: Der §3.13-Suchlauf des Implementers war an die schmale
Plan-Mustersatz-Grep-Form gebunden („deckt HTTP-API|pgchangefeed-0.1.0|no
gRPC/SSE/NATS surface exists yet" über `spec/`, `docs/`, `sdks/python/`)
und verfehlte alle fünf Fundstellen in deren tatsächlicher Schreibform;
`harness/` und die Wurzel-READMEs lagen außerhalb des Pfadraums. Die Kette
trägt ihre Fundstellen je selbst:

- **Haupt-Review F-1 (HIGH, vom Reviewer gefunden, Klasse doppelt
  benannt):** die `SPEC-027`-Vertragszeile in `spec/pflichtenheft.md` §6
  trug weiter „aktuell `0.1.0`" und die Ein-Vertrags-Liste, während
  derselbe Diff in derselben Datei (§1) auf die Vier-Wege-Form + `0.2.0`
  zog und `pyproject.toml` real `0.2.0` misst. Die Zahl war bei ihrer
  letzten Niederschrift wahr und wurde durch diese Arbeit falsch — nach
  der Klassen-Grenze („dort driftet ein Träger schon; hier driftet er
  durch diese Arbeit") zählt der Vorgang hier, die
  Zahl-Drift-Klasse (`BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`)
  trägt kein eigenes Vorkommen. Fixrunde 1 zog die Zeile auf „aktuell
  `0.2.0`" + Vier-Verträge-Liste.
- **Haupt-Review F-4 (MEDIUM, vom Reviewer gefunden):** fünf
  Träger-Stellen (`harness/README.md` sdk-pack- und
  test-sdk-Zeilen, Wurzel-`README.md`, `README.de.md`,
  `sdks/python/README.md` Intro) trugen die Ein-Flächen-/`0.1.0`-Formen;
  die Plan-deklarierte Grep-Form traf keine einzige in ihrer Schreibform.
  Fixrunde 1 zog drei der fünf Stellen.
- **Fixrunden-Re-Review R-1 (HIGH) / R-2 (MEDIUM):** das committete
  Suchlauf-Feld behauptete für genau die restlichen Stellen vollzogene
  Nachzüge — der Baum widerlegte beide Behauptungen (`README.de.md:27`
  unverändert, `harness/README.md`-Reject-Satz unverändert); das Feld ist
  selbst ein Träger, dessen Behandlungs-Angaben gemessen werden müssen.
  Fixrunde 2 zog die drei Reststellen und stellte die Feld-Zeilen auf den
  realen Ist-Stand.

Einordnung: dieselbe Unter-Klasse „gefunden vom Reviewer, nicht vom
Implementer-Suchlauf" wie bei `slice-095`/`slice-097`/
`slice-sdk-csharp-http-client-flaeche`/`slice-sdk-kotlin-nats-stream-client-flaeche`/`slice-sdk-python-sse-client-flaeche`
— neu an dieser Kette ist die dritte Stufe: der Suchlauf-Befund selbst
wurde zum Träger, dessen Behauptungen das Re-Review am Baum widerlegte
(R-1 trägt zugleich die Klasse
`BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht`). Geschärfte Lehre
(Closure-Notiz §7 des Slices): der Suchlauf folgt §3.13s Wortlaut (über
die Träger nach der bewegten Eigenschaft), nicht dem Plan-Mustersatz, und
das committete Suchlauf-Feld wird am Baum nachgemessen, nicht deklariert.
Ausgang bleibt **verkörpert** (`AGENTS.md` §3.13), kein neuer
Schwellen-Übertritt — die Schärfung ist eine Anwendungs-Schärfung der
verkörperten Regel; ob sie einen eigenen Satz im Träger trägt, prüft der
Lese-Schritt der Welle-Closure.

Quelle: Haupt-Review `slice-sdk-python-nats-stream-client-flaeche`
(F-1, F-4, F-5), Fixrunden-Report desselben Slices (R-1, R-2),
Verifikationsbericht desselben Slices (§2.5, V-2), Commits `0f8cc4f2`,
`b4d1d352`, `a3e9d64e`.