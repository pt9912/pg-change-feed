# Verifikationsbericht: slice-104 — 2026-09-17

**Rolle:** Verifier (Modul 8/11) — „Bauen wir es richtig?" gegen den
DoD-Vertrag (`slice-104` §2, LP1–LP3) und die im Slice referenzierten
Entscheidungen [`ADR-0060`](../plan/adr/0060-grpc-streaming-mechanismus.md)
(Folgepflicht „eine neue Build-Stufe im `Dockerfile` bzw. ein neues
`make`-Ziel für die Code-Generierung" — Generator und seine Existenz, nicht
sein Mount-Mechanismus) und [`ADR-0084`](../plan/adr/0084-sync-gate-fuer-generierte-artefakte.md)
(Kontext-Bezug, unberührt) sowie die Hard Rules `AGENTS.md` §3.1
(Docker-only), §3.9 (Exit-Code-Disziplin), §3.7/§3.13 (Kommentar-Klassen,
Träger-Nachzug). **Nicht** gegen den Diff als solchen (Reviewer-Aufgabe,
[`review-slice-104.md`](review-slice-104.md)) und **nicht** gegen realen
Bedarf (Validator, nicht ausgelöst).

**Frischer Kontext.** Diese Sitzung hat den Slice-Plan vollständig gelesen
(§1–§8), die beiden Implementierungs-Commits (`338cfe0`, `59d5b53`) samt
vollem Diff (`Dockerfile`, `Makefile`, `.dockerignore`, `.gitignore`, sowie
die fünf nachgezogenen Träger `AGENTS.md`, `harness/README.md`,
`harness/sensors/generated-sync.md`, `tools/harness/generated-sync.sh`), den
Review-Report (`c7fa290`), `ADR-0060` vollständig (§Entscheidung,
§Folgepflichten) und das Beobachtungs-Register
`BEO-PGC/coverage-stage-dockerignore-blockiert-tooling/`
(`observation.md`, `state.md`, beide Evidence-Dateien) gelesen. Jede Aussage
dieses Berichts stammt aus einem hier selbst gefahrenen Lauf oder einer hier
selbst gelesenen/abgefragten Quelle — Review-Zahlen und -Befunde waren
Kontext, nicht übernommen.

**Gegenstand.** `338cfe0` (LP1+LP2, mount-loser Umbau) → `59d5b53` (LP3,
Träger-Nachzug) → `c7fa290` (Review-Report, 0 HIGH, 1 MEDIUM, 1 LOW). Der
Slice liegt in `in-progress/`; der `git mv` nach `done/` ist **nicht**
erfolgt — erwartungsgemäß. `git status` war vor dieser Sitzung sauber; alle
in dieser Sitzung erzeugten Artefakte (`.proto-generate.tar`,
Test-`Makefile`-Fragment) liegen ausschließlich im Scratchpad bzw. wurden
per `.gitignore` ausgeschlossen und rückstandslos entfernt.

---

## 1. Eigene Messungen dieses Laufs

Jeder Lauf ungepiped, Exit-Code direkt aus einem eigenen, abgeschlossenen
Schritt gelesen (`AGENTS.md` §3.9).

| # | Lauf | Exit | Ergebnis |
|---|---|---|---|
| 1 | `make proto-generate` (realer, unveränderter Nutzer, kein `sudo`) | **0** | `docker build --target proto-export` (gecacht) → `docker run --rm --network none … > .proto-generate.tar` → `tar -xf … -C .` → `rm -f`; extrahierte Dateien (`gen/cdc/stream/v1/{changestream.pb.go,changestream_grpc.pb.go}`) gehören real `db:db` (aufrufender Nutzer) |
| 2 | Inhalts-/Eigentümer-Vergleich (Kopie vor dem Lauf vs. danach) | — | `diff` beider Dateien: **byte-identisch**; `ls -l`: Eigentümer unverändert `db:db`; `git status --short`: leer (kein Arbeitsbaum-Drift) |
| 3 | Reale Fehlschlag-Simulation: `docker run --rm --network none <nicht-existentes-image> > datei` | **125** | `docker run` scheitert real mit „pull access denied"; die Datei bleibt leer (0 Byte) |
| 4 | Isoliertes `make`-Testrezept mit identischem Rezeptmuster (`> Datei` gefolgt von `tar -xf`) gegen dasselbe nicht-existente Image | **2** (make), Rezeptzeile mit **125** | `make: *** [test.mk:4: t] Fehler 125` — die Rezeptzeile mit `tar -xf` läuft **nachweislich nicht**, der `echo`-Marker danach erscheint nicht im Output |
| 5 | `make generated-sync` (nach Lauf 1, unverändertes Erzeugnis) | **0** | „byte-gleich der Ausgabe des gepinnten Generators" — bestätigt eigenständig |
| 6 | `make gates` (auf `HEAD` = `c7fa290`) | **2** | **rot** — `docs-check` bricht mit `docs/reviews/review-slice-104.md:34 ADR-0060 id-unlinked` ab; `baseline-verify` und `coverage-gate` liefen zuvor grün, `a-check`/`commit-traceability`/`generated-sync` liefen wegen des Abbruchs nicht mehr in diesem Aufruf |
| 7 | `make a-check` (separat) | **0** | `gesamt: 0 Befund(e)` |
| 8 | `make commit-traceability` (separat, `HEAD~5..HEAD`) | **0** | 5 Commits, 0 Befunde, Betreffs ohne Struktur-ID |
| 9 | `grep -n -A10 "^proto-generate:" Makefile` | — | reine Umleitung (`>` bei `docker run`), kein Pipe-Zeichen, kein `-`-Präfix vor `tar -xf`/`rm -f` |
| 10 | `grep -in "bind.mount\|docker run.*-v\|mount" docs/plan/adr/0060-grpc-streaming-mechanismus.md` | — | leer — bestätigt Review-F-1 unabhängig: `ADR-0060` trägt keine Bind-Mount-Aussage |
| 11 | `grep -rn "PROTO_RUN_USER" …` (lebende Träger, `done/`/`docs/reviews/` ausgenommen) | — | keine Treffer außerhalb des Slice-Plans selbst |
| 12 | Eigener, breiterer Träger-Suchlauf (`grep -rn "proto-generate"`, `grep -in "bind-mount"`) über alle lebenden `.md`/`.sh`/`Makefile`/`Dockerfile` | — | alle verbleibenden Bind-Mount-Erwähnungen betreffen entweder Historie (`vormals …`, Dockerfile:41), einen anderen Mechanismus (`generated-sync`s eigener `:ro`-Mount, `schema-rollout`s `D_MIGRATE_RUN_USER`) oder sind bereits korrekt auf den neuen Mechanismus umgestellt — kein sechster Fund |

**Befund 6 ist der zentrale Fund dieser Verifikation** und wird unten (§4)
ausführlich behandelt.

---

## 2. LP1 — neue Dockerfile-Stufe, Build-Zeit-Erzeugung, stdout-Ausgabe

**Erfüllt, gemessen.** `Dockerfile` trägt eine neue Stufe `proto-export`
(`FROM proto AS proto-export`), die die `.proto`-Quelle per `COPY proto/
proto/` in den Image-Layer holt (kein Bind-Mount), `protoc` als
`RUN`-Schritt zur Build-Zeit ausführt (Ausgabe unter `/out` im Layer), und
per `ENTRYPOINT ["tar", "-cf", "-", "-C", "/out", "."]` das Erzeugnis über
stdout ausgibt. Eigener Build (Lauf 1) bestätigt: `--network none` gilt
bereits für den Build-Schritt selbst (kein Netz nötig, dieselbe Eigenschaft
wie der bisherige Lauf). Der Kommentarblock über der Stufe beschreibt den
Mechanismus korrekt im Präsens.

## 3. LP2 — host-seitige Extraktion, kein `docker run -v`, `PROTO_RUN_USER` entfernt

**Erfüllt, gemessen — mit einer begrüßenswerten Abweichung vom im Plan
skizzierten Standard-Pfad.** Das Rezept (`grep`-Beleg, Lauf 9) verwendet
**keine** Pipe (`docker run … | tar -x`), sondern eine reine Ausgabe-
Umleitung in eine Zwischendatei (`> $(PROTO_GENERATE_TARBALL)`), gefolgt von
einem separaten `tar -xf`-Aufruf. Das ist **eine** der beiden vom Plan §6
Risiko 2 explizit genannten Optionen („ein zweistufiger Aufruf: Stream erst
in eine Datei, Exit-Code prüfen, dann extrahieren") — und robuster als die
andere (`PIPESTATUS`/`pipefail`), weil sie unabhängig von der Shell des
Make-Rezepts funktioniert (Makefile-Kommentar begründet dies selbst korrekt:
`/bin/sh` auf diesem Host ist `dash`, das kein `set -o pipefail` trägt).
Eigene Fehlschlag-Proben (Lauf 3, Lauf 4) bestätigen real: Ein
`docker run`-Fehlschlag stoppt die Make-Rezeptzeile über ihren eigenen
Exit-Code (GNU-Make-Default, kein Pipe-Glied), bevor `tar -xf` überhaupt
läuft. `PROTO_RUN_USER` ist aus dem `Makefile` entfernt (bestätigt: keine
Fundstelle mehr in lebenden Trägern, Lauf 11); die extrahierten Dateien
gehören real dem aufrufenden Nutzer (Lauf 1/2). `make generated-sync` bleibt
grün (Lauf 5).

## 4. LP3 — Fünf Träger nachgezogen, plus ein selbst gefundener, aktuell roter sechster

**Die fünf im Plan benannten Träger sind korrekt nachgezogen** (`AGENTS.md`
§4, `harness/README.md` §Sensors, `harness/sensors/generated-sync.md`,
`tools/harness/generated-sync.sh`, Dockerfile-Kommentar) — eigener,
breiterer Suchlauf (Lauf 12) findet keinen weiteren lebenden Träger, der den
alten Mechanismus im Präsens beschreibt. `tools/harness/generated-sync.sh`s
Diff ändert ausschließlich Kommentarzeilen (Skript-Logik byte-identisch,
eigener `diff` bestätigt).

**Ein sechster, vom Slice nicht vorgesehener Träger ist jedoch selbst
betroffen — und aktuell rot.** `docs/reviews/review-slice-104.md:34` trägt
eine nackte `ADR-0060`-Kennung ohne Link auf ihre Definition, außerhalb von
Inline-Code (der Satz zitiert eine Dockerfile-Kommentarzeile im
Fließtext: `„…siehe ADR-0060"` — die schließende Anführung folgt direkt auf
die nackte Kennung; alle sechs anderen `ADR-0060`-Erwähnungen im selben
Report stehen korrekt in Backticks und werden von `docs-check`s
`ids`-Modul nicht beanstandet). Eigener `make gates`-Lauf (Lauf 6) bestätigt
dies real: `docs-check` bricht mit exakt diesem einen Befund ab, `make
gates` liefert Exit 2. Das ist **kein Effekt des Umbaus selbst** (LP1/LP2
sind funktional korrekt und real geprüft grün) — der Fund liegt im
Review-Report-Text, nicht im produktiven Diff.

Dieser Fund entspricht **strukturell** der bereits verkörperten
Beobachtungsklasse `BEO-PGC/report-nackte-id-ohne-link`
(5×, `seit slice-063` in `AGENTS.md` §3.9 verkörpert — dort allerdings unter
einer inzwischen auf die *Sequenzierungs*-Fehlerklasse geschärften Regel,
„Prüfung und Folgehandlung sind zwei Schritte"). Der hier vorliegende Fund
ist mechanisch derselbe **Basis**-Fehler, den der Registereintrag
ursprünglich benennt (nackte Kennung im Fließtext eines Reports) — ob er
als 6. Beleg in dasselbe Verzeichnis gehört oder eine neue, engere Klasse
braucht, ist eine Register-Urteilsfrage für die Closure (Modul 6: „Mensch
urteilt"); ich benenne ihn hier, ohne die Eintragung selbst vorzunehmen
(Eintragung ist Planner-Arbeit bei der Slice-Closure, nicht Verifier-Arbeit).

**Konsequenz für die DoD-Zeile „`make gates` grün":** Diese Zeile ist
**zum jetzigen Zeitpunkt nicht erfüllt** und darf nicht angehakt werden,
bis der Fund in `docs/reviews/review-slice-104.md:34` behoben ist (ein
Markdown-Link auf `docs/plan/adr/0060-grpc-streaming-mechanismus.md` oder
Backticks analog den sechs anderen Stellen — das ist eine
Zitat-Korrektur an einem noch nicht `done/`-abgelegten, nicht-immutablen
Bericht, keine ADR-Änderung).

---

## 5. Plan-vs-Code-Diff

**Was der Diff nicht enthält, obwohl der Plan es nennt:** nichts Fehlendes
gefunden — `harness/image-hash.txt` ist unverändert (`git diff` leer), wie
im Plan §3 als Erwartung (nicht Messwert) benannt; eigener `make image`
ergänzend nicht erneut gelaufen, da Reviewer bereits real geprüft hat
(„denselben Digest wie `harness/image-hash.txt`") und der Dockerfile-Diff
(neue Stufe hängt nicht an der `runtime`-Kette, Zeilen 19/35/50/76/102/108
unverändert in ihrer `FROM`-Kette) das strukturell stützt.

**Was der Diff enthält, obwohl der Plan es nicht als eigene Zeile nennt:**
`.gitignore` (neue Zeile für `.proto-generate.tar`, transientes
Zwischenprodukt — konsistent mit dem im Plan §6 Risiko 2 skizzierten
zweistufigen Aufruf, kein stiller Umfangs-Zuwachs) und
`docs/reviews/review-slice-104.md` (Modul-8-Übergabe, erwartet).

Kein Schicht- oder Umfangs-Bruch gegenüber §1 (Abgrenzung): `generated-sync`,
`.proto`-Inhalt, `schema-rollout`/`D_MIGRATE_RUN_USER` und ein neues Gate
bleiben unangetastet — eigene Diffs bestätigen dies (Lauf 5,
`tools/harness/generated-sync.sh`-Logik-Diff = 0 Zeilen außerhalb von
Kommentaren, kein `schema-rollout`-Treffer im Diff).

---

## 6. F-2 — eigenes Urteil zur Beobachtungsklassen-Frage

Der Reviewer hält die `!proto/`-`.dockerignore`-Ausnahme **nicht** für eine
dritte Instanz von `BEO-PGC/coverage-stage-dockerignore-blockiert-tooling`.
Eigene, unabhängige Lektüre von `observation.md`/`state.md`/beiden
Evidence-Dateien bestätigt dieses Urteil:

- Der Registereintrag ist über **Pfad und Text** eng auf „eine neue
  Docker-Stage, die ein **Skript unter `tools/`** braucht" gezogen — beide
  Evidence-Dateien (`slice-049`, `slice-093`) betreffen `tools/coverage-gate.sh`
  bzw. dasselbe Muster; `state.md`s eigene Handlungsanweisung für künftige
  Wiederholung nennt wörtlich „eine weitere Docker-Stage, die ein Skript
  unter `tools/` braucht". Der Titel selbst (`…-tooling`) ist die
  Pfad-Kennung, und Modul 6 ist hier ausdrücklich strikt: „Ein Prosa-Name
  taugt nicht — er darf umformuliert werden, ein Pfad nicht."
- Der `!proto/`-Fund betrifft ein **Quellverzeichnis für eine neue
  Erzeugungsstufe** (`proto/`), kein Skript unter `tools/`; die zweite,
  bei beiden bisherigen Belegen mitlaufende Facette (Alpine-Basis ohne
  `bash`) tritt hier gar nicht auf.
- Der zugrunde liegende **Mechanismus** ist zwar real derselbe
  (`.dockerignore`s Allow-Listen-Default-Deny bricht bei jeder neuen
  Docker-Stufe, die einen bisher ungelisteten Pfad braucht) — eigene
  Gegenprobe nicht nötig, da bereits vom Reviewer real durchgeführt
  (`COPY proto/ proto/` ohne `!proto/`-Zeile scheitert mit „not found") und
  von mir durch den erfolgreichen Lauf 1 (mit der Zeile) indirekt
  bestätigt.

**Eigenes Urteil:** Ich teile das Reviewer-Urteil. Der Fund ist derselbe
**Mechanismus**, aber **nicht** dieselbe **benannte Klasse** — die
Registereintrag-Identität ist über den Pfad `.../blockiert-tooling` und den
`tools/`-Skript-Bezug beider Evidence-Dateien enger gezogen als der neue
Fund trägt. Eine Verbuchung als 3. Instanz würde den Zähler unter einem zu
engen Titel erhöhen und die daraus abgeleitete Regel („bei einem Skript
unter `tools/` vorab prüfen") träfe den neuen Fall nicht mehr präzise ab.
Die für die Closure richtige Antwort ist entweder ein **neuer**
`BEO-PGC/`-Eintrag für den generalisierten Mechanismus (z. B.
„docker-stage-dockerignore-default-deny" oder vergleichbar präzise
benannt) oder eine bewusste, im bestehenden Eintrag selbst dokumentierte
Erweiterung von Titel und Text — beides ist Planner-Arbeit bei der Closure,
keine Verifier-Entscheidung. Ich lege mich **nicht** auf eine der beiden
Formen fest; das Urteil, welche träfe, bleibt beim Planner.

---

## 7. §6-Risiken — nicht schädlich eingetreten

| Risiko | Ausgang (eigenes Urteil) | Beleg |
|---|---|---|
| 1. Datei-Attribute/Zeilenenden-Divergenz durch `tar`-Extraktion | **entfallen** — nicht eingetreten | Lauf 2 (byte-identischer `diff`), Lauf 5 (`generated-sync` grün) |
| 2. Extraktions-Pipe maskiert Exit-Code | **entfallen** — durch Architekturentscheidung strukturell vermieden (kein Pipe-Glied, reine Umleitung + separater Schritt) | Lauf 3, Lauf 4 (reale Fehlschlag-Proben, Exit-Code sichtbar propagiert) |
| 3. Unvollständiger Träger-Suchlauf | **entfallen** — eigener, breiterer Suchlauf findet keinen sechsten *im Plan erwarteten* Träger | Lauf 12 |

Keines der drei geplanten Risiken ist schädlich eingetreten. Der in §4
benannte sechste Träger-Fund lag **außerhalb** der drei geplanten Risiken
(ein Report-Text, kein Präsens-Beschreibungs-Träger im Sinn von LP3) und
ist deshalb kein Ausgang eines der drei — er ist ein eigenständiger,
unerwarteter Fund dieser Verifikation und gehört als solcher in die
Closure-Notiz, nicht als vierter Unterpunkt unter §6.

---

## 8. Verdikt

**Noch nicht DoD-konform — ein selbst gefundener, aktuell roter Blocker.**
LP1 und LP2 sind gemessen vollständig erfüllt (eigener Build, eigener Lauf,
eigene Fehlschlag-Simulation, eigener Byte-Vergleich). LP3 ist für die fünf
im Plan benannten Träger vollständig erfüllt. **Aber:** `make gates` ist zum
Verifikationszeitpunkt **rot** (Exit 2) — `docs-check` findet eine nackte,
unverlinkte `ADR-0060`-Kennung in `docs/reviews/review-slice-104.md:34`,
außerhalb von Inline-Code. Die DoD-Zeile „`make gates` grün" darf deshalb
nicht angehakt werden, und der Slice darf **nicht** nach `done/` wandern,
bis dieser Fund behoben ist (Zitat-Korrektur am Review-Report: Backticks
oder Markdown-Link, analog den sechs anderen `ADR-0060`-Stellen im selben
Dokument).

**F-1 (Reviewer, LOW):** eigenständig bestätigt — `ADR-0060` trägt
nachweislich keine Bind-Mount-Aussage (Lauf 10).

**F-2 (Reviewer, MEDIUM, Beobachtungsklassen-Frage):** eigenes Urteil
**teilt** das Reviewer-Urteil — der `!proto/`-Fund ist derselbe Mechanismus,
aber nicht dieselbe benannte Registerklasse wie
`BEO-PGC/coverage-stage-dockerignore-blockiert-tooling` (§6 oben). Ein
neuer Registereintrag oder eine dokumentierte Erweiterung des bestehenden
ist bei Closure fällig; welche Form, entscheidet der Planner.

**ADR-Konformität:** `ADR-0060` verlangt ausschließlich Existenz und
Generator, nicht den Mount-Mechanismus (eigener `grep`, Lauf 10 und
zusätzliche Suche nach „Folgepflicht"/„Build-Stufe" in der ADR); kein neues
ADR nötig — bestätigt.

**Für die Closure noch offen (Planner-/Implementer-Arbeit):**

1. **Blocker beheben:** `docs/reviews/review-slice-104.md:34` — nackte
   `ADR-0060`-Kennung verlinken oder in Backticks setzen; danach `make
   gates` real erneut prüfen, bevor die DoD-Zeile angehakt wird.
2. §2-Checkboxen (LP1–LP3, `make gates`) erst nach Punkt 1 anhaken.
3. §6 — die drei Risiken mit den in §7 dieses Berichts vorgeschlagenen
   Ausgängen (alle *entfallen*) eintragen; der sechste Träger-Fund gehört
   zusätzlich in die Closure-Notiz, nicht in §6 selbst.
4. §7 Closure-Notiz inkl. Steering-Loop-Lerneintrag — Kandidat: der
   zweistufige Extraktions-Aufruf (Umleitung statt Pipe) als
   verallgemeinerbares Muster für „Erzeugen-im-Build-plus-Extraktion"
   (Plan §5 (a) bereits als Kandidat benannt).
5. F-2 / §6 dieses Berichts: Registerentscheidung treffen (neuer Eintrag
   vs. Erweiterung des bestehenden).
6. Der sechste Träger-Fund (§4) gehört als Beleg-Kandidat für
   `BEO-PGC/report-nackte-id-ohne-link/` in die Closure-Prüfung — Urteil,
   ob 6. Instanz derselben Klasse oder eigenständig, liegt beim Planner.
7. Der `git mv` nach `done/` folgt erst nach 1–6.

**Nicht Gegenstand dieser Verifikation:** Validierung gegen realen Bedarf
(Validator, hier nicht ausgelöst).
