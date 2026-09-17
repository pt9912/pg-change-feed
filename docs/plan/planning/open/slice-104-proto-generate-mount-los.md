# Slice slice-104: Mount-loser Umbau von `make proto-generate`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — es gibt keine Closure-Bedingung, die über die eigene
DoD dieses Slice hinausgeht (Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht).

**Bezug:** [`ADR-0060`](../../adr/0060-grpc-streaming-mechanismus.md)
(Folgepflicht „eine neue Build-Stufe im `Dockerfile` bzw. ein neues
`make`-Ziel für die Code-Generierung" — legt den **Generator und seine
Existenz** fest, nicht seinen Mount-Mechanismus; der Bind-Mount/`--user`-Weg
war eine Umsetzungsentscheidung von `slice-069`, keine `§Entscheidung`
dieser ADR) · [`ADR-0084`](../../adr/0084-sync-gate-fuer-generierte-artefakte.md)
(Kontext beschreibt `make proto-generate`s heutigen In-Place-Bind-Mount als
Begründung, warum `make generated-sync` **eigenständig** in ein
Temp-Verzeichnis erzeugt — dieser Slice ändert **nichts** an `ADR-0084`s
Entscheidung oder Fitness Function für `generated-sync`; er berührt nur
lebende Träger, die den heutigen `proto-generate`-Mechanismus im Präsens
beschreiben, siehe §3).

**Berührte Spec-Stellen:** — (geprüft: `grep -n "proto-generate\|protoc\|Bind-Mount" spec/lastenheft.md spec/pflichtenheft.md` ist leer; reiner Werkzeug-/Toolchain-Mechanismus ohne Vertrags- oder Draht-Bezug).

**Verantwortlich:** — (bis zur Priorisierung `open` → `next`).

**Autor:** pt9912 (Planner). **Datum:** 2026-09-17.

---

## 1. Ziel und Abgrenzung

<!-- BEDIENHINWEIS: Ziel = ein Satz, Liefer-Fokus, kein "wir machen
aufraeumen". Abgrenzung = je Punkt eine Begruendung, nicht nur eine Nennung:
ein Ausschluss ohne Grund ist eine Behauptung. Keine Mindestzahl — ein echter
Ausschluss ist besser als vier erfundene. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** `make proto-generate` (`Makefile:102-109`) verliert den
schreibbaren Bind-Mount des Arbeitsbaums (`-v "$(CURDIR)":/src`) und den
`--user "$(PROTO_RUN_USER)"`-Workaround gegen root-owned Dateien. An seine
Stelle tritt eine **mount-lose** Form: eine Dockerfile-Stufe kopiert die
`.proto`-Quelle in sich hinein und erzeugt den Go-/gRPC-Code **während des
Builds** (`RUN protoc …`, Ausgabe liegt im Image-Layer); ein
Kommando/Entrypoint dieser Stufe gibt das Erzeugnis über stdout aus
(`tar -cf - -C <verzeichnis> .`), und `make proto-generate` extrahiert es
**host-seitig** (`docker run --rm … $(PROTO_IMAGE) | tar -x -C .`). Die
extrahierten Dateien gehören dadurch automatisch dem aufrufenden Nutzer —
kein `docker run -v` mehr, kein `--user`-Trick.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **`make generated-sync` (`tools/harness/generated-sync.sh`) bleibt
  unverändert.** *Bestand bleibt bewusst stehen.* Der Sensor schreibt
  bereits **nicht** in den Arbeitsbaum: Die Quelle hängt als `:ro`-Mount,
  die Ausgabe geht in ein `mktemp -d`-Verzeichnis, das ein `trap`
  wegräumt. Das Problem, das dieser Slice löst — root-owned Dateien **im
  committeten Baum** —, tritt dort strukturell nicht auf, weil nichts
  Erzeugtes den Baum je erreicht. Ein Umbau brächte keinen Mehrwert und
  risikiert stattdessen Kollateralschaden an einem funktionierenden,
  unabhängigen Gate (`ADR-0084`) ohne Gegenleistung.
- **Der Inhalt der `.proto`-Quelle.** Der Draht-Vertrag ist
  [`ADR-0060`](../../adr/0060-grpc-streaming-mechanismus.md)s Gegenstand;
  dieser Slice ändert keine Nachricht, kein Feld — nur, **wie** der Code
  aus ihr entsteht.
- **`make schema-rollout`/`D_MIGRATE_RUN_USER`** (`Makefile:132-180`), das
  demselben `--user`-Muster folgt. *Anderer Vorgang.* Der Auftrag benennt
  ausschließlich `proto-generate`; ein Rollout mit Pflicht-Report und
  Rollback-Artefakt ([`ADR-0043`](../../adr/0043-schemamigrationen-mit-d-migrate.md))
  ist ein eigener Vertrag mit eigener Beweislast, kein Beifang dieses Zugs.
- **Ein neues Gate.** Der Wächter für die Erzeugnis-Treue bleibt
  `make generated-sync` ([`ADR-0084`](../../adr/0084-sync-gate-fuer-generierte-artefakte.md)).
  Dieser Slice ändert einen Mechanismus, keinen Prüf-Vertrag — ein
  zusätzliches Gate wäre Zeremonie.

**Keine Mindestzahl.** Ein Slice mit *einem* echten Ausschluss ist besser als
einer mit vier erfundenen; die vier Klassen sind ein Suchraster, keine
Ausfüll-Liste.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

<!-- BEDIENHINWEIS: je Zeile ein pruefbares Kriterium. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **LP1 — die neue Dockerfile-Stufe erzeugt zur Build-Zeit und gibt über
      stdout aus.** Eine neue Stufe (Name führt der umsetzende Zug; baut auf
      der bestehenden `proto`-Stufe auf, die `protoc`/die beiden
      `protoc-gen-*`-Plugins gepinnt trägt) kopiert die `.proto`-Quelle
      hinein, ruft `protoc` als `RUN`-Schritt auf (Ausgabe liegt im
      Image-Layer, kein Bind-Mount, kein Netz nötig — dieselbe
      `--network none`-Eigenschaft, die der heutige Lauf schon hat, gilt
      jetzt für den Build-Schritt selbst), und trägt ein
      Kommando/Entrypoint, das das Erzeugnis via `tar -cf - -C <verzeichnis> .`
      über stdout ausgibt.
- [ ] **LP2 — `make proto-generate` extrahiert host-seitig, kein
      `docker run -v` mehr.** Das Rezept baut die neue Stufe und ruft sie
      ohne Bind-Mount auf (`docker run --rm --network none $(PROTO_IMAGE) |
      tar -x -C .`); die Extraktions-Pipe sichert ihren Exit-Code gegen das
      §3.9-Muster ab (`AGENTS.md` §3.9 — ein `docker run`-Fehlschlag darf
      nicht hinter einem erfolgreichen, aber leeren `tar -x` verschwinden;
      Umsetzung z. B. `PIPESTATUS`/`pipefail` im Rezept). `PROTO_RUN_USER`
      und der `--user`-Aufruf entfallen, sofern nach dem Umbau ungenutzt
      (recherchiert: `grep -rn "PROTO_RUN_USER" --include="*.md"
      --include="Makefile" --include="*.mk" .` findet die Variable nur in
      der Makefile-Definition/-Verwendung selbst sowie in `done/` und
      `docs/reviews/` — beides unveränderliche Aufzeichnungen, kein lebender
      Träger; der Implementer-Suchlauf prüft das eigenständig erneut,
      §3.13). `make generated-sync` bleibt grün (unverändertes,
      byte-gleiches Erzeugnis).
- [ ] **LP3 — Träger-Nachzug: die Präsens-Beschreibungen des heutigen
      Bind-Mount-Mechanismus.** Gefunden (Suchlauf `grep -rn
      "proto-generate" --include="*.md" --include="*.sh" --include="Makefile"
      --include="Dockerfile" .`, `done/`/`docs/reviews/` ausgenommen): `harness/README.md`
      §Sensors (Zeile zu `make proto-generate`, nennt „läuft als
      Aufrufer-uid …") · `AGENTS.md` §4 (dieselbe Zeile, kürzer) ·
      `harness/sensors/generated-sync.md` (zwei Stellen: „das unterscheidet
      es von `make proto-generate`, das in-place in den Bind-Mount
      erzeugt …") · `tools/harness/generated-sync.sh` (Kommentarzeilen, die
      denselben Satz tragen) · der Kommentarblock über der `proto`-Stufe im
      `Dockerfile`. Alle fünf beschreiben den heutigen Mechanismus im
      Präsens und werden mit diesem Umbau falsch, wenn sie nicht
      mitgezogen werden (§3.13 — Arbeit, die eine beschriebene Eigenschaft
      bewegt, zieht ihre Träger nach; `BEO-PGC/arbeit-ueberholt-stehenden-traeger`
      steht bereits bei 5× und ist in `AGENTS.md` §3.13 verkörpert).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update für den geänderten Mechanismus siehe LP3 — kein weiterer
      öffentlicher Vertrag berührt.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. **Entfällt** — Repo ist Greenfield (`harness/conventions.md` Modus-Deklaration `*`/`PGC` = GF), keine `docs/plan/planning/reconciliation.md` vorhanden.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

<!-- BEDIENHINWEIS: Datei- oder Komponenten-Ebene reicht; der
Implementer-Agent erweitert die Liste in seinem ersten Lauf. -->

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `Dockerfile` | update | Neue Stufe (Name führt der umsetzende Zug), die auf `proto` aufbaut: `COPY proto/ ./proto`, `RUN protoc …` zur Build-Zeit, Kommando/Entrypoint mit `tar -cf - -C <verzeichnis> .`; Kommentarblock über der/den betroffenen Stufe(n) beschreibt den neuen Mechanismus statt des alten Bind-Mounts (LP1, LP3). |
| `Makefile` | update | `proto-generate`-Rezept baut die neue Stufe und extrahiert per `docker run … \| tar -x -C .`; `PROTO_RUN_USER` entfernt, falls ungenutzt; Exit-Code-Absicherung der Pipe (§3.9-Muster, siehe LP2). |
| `harness/README.md` | update | §Sensors-Zeile zu `make proto-generate` beschreibt den neuen Mechanismus statt „läuft als Aufrufer-uid …" (LP3). |
| `AGENTS.md` | update | §4-Zeile zu `make proto-generate`, falls sie den Mechanismus nennt (aktuell knapper als `harness/README.md`, prüfen statt annehmen). |
| `harness/sensors/generated-sync.md` | update | Zwei Stellen nennen `make proto-generate`s heutigen In-Place-Bind-Mount im Präsens als Abgrenzung zum eigenen Temp-Verzeichnis-Mechanismus — nach dem Umbau falsch (LP3). |
| `tools/harness/generated-sync.sh` | update | Kommentarzeilen (Skript-Kopf), die denselben Satz tragen wie oben (LP3, §3.7 — ein Kommentar beschreibt, was da ist). |
| `harness/image-hash.txt` | prüfen, ggf. update | `Dockerfile` ist eine Build-Kontext-Datei ([`ADR-0044`](../../adr/0044-image-beleg-semantik.md) Punkt 3) — `make image` läuft vor der Closure. Erwartet **unverändert**, weil die neue Stufe nicht Teil der `runtime`-Zielkette ist (BuildKit baut nur die für das Ziel-Stage nötigen Stufen) — das ist eine Erwartung, kein Messwert (§3.12); real zu prüfen, nicht anzunehmen. |

## 4. Trigger

<!-- BEDIENHINWEIS: Beispiele — "Wenn Welle X done." / "Wenn Carveout CO-NN
aufgeloest." -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `in-progress/` trägt aktuell keinen
Slice (Ruhe-Marker gesetzt, WIP-Limit frei) und dieser Slice hat keine
Abhängigkeit zu einer laufenden Welle oder einem anderen Slice. Ohne
Rückfrage feststellbar über `ls docs/plan/planning/in-progress/`.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): wenn der Umbau
  zusätzlich `make schema-rollout`/`D_MIGRATE_RUN_USER` oder
  `make generated-sync` mitziehen müsste, um konsistent zu bleiben (z. B.
  weil ein gemeinsamer Helfer entsteht, den beide Ziele teilen sollen).
  Dann ist der Schnitt falsch: dieser Slice bewegt **einen** Mechanismus,
  nicht ein Muster über mehrere Ziele hinweg.
- `in-progress` → `open` (blockiert — Carveout?): wenn das
  `tar`-Stdout-Streaming zwischen Container und Host nicht sauber
  funktioniert (z. B. Buffering-/Encoding-Probleme der konkreten
  Docker-Engine-Version) oder wenn `make generated-sync` durch den
  Mechanismus-Wechsel rot wird und der Grund nicht im selben Zug behebbar
  ist (z. B. eine unerwartete Divergenz zwischen build-zeitlich erzeugten
  und laufzeitlich erzeugten `.pb.go`-Dateien).

## 5. Closure-Trigger

<!-- BEDIENHINWEIS: z.B. "DoD vollstaendig + PR gemerged + Closure-Notiz
geschrieben." -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

**Zwei beobachtbare Kriterien:** (1) `make proto-generate` läuft ohne
`docker run -v` und ohne `--user`, die extrahierten `.pb.go`-Dateien gehören
real dem aufrufenden Nutzer (gemessen: `ls -l gen/cdc/stream/v1/*.go` nach
einem Lauf als normaler Nutzer, kein `sudo`/`chown` nötig). (2)
`make generated-sync` bleibt nach dem Umbau grün (byte-gleiches Erzeugnis,
[`ADR-0084`](../../adr/0084-sync-gate-fuer-generierte-artefakte.md) unberührt)
— **und** `make gates` insgesamt grün.

**Lerneintrag:** geschärfte Regel, neuer Sensor oder benannte Spec-Lücke in
§7. Naheliegende Kandidaten: (a) eine geschärfte Fassung von `AGENTS.md`
§3.1/§3.9 für das „Erzeugen-im-Build-plus-stdout-Extraktion"-Muster, falls
es sich als wiederverwendbar für `schema-rollout` erweist (dann eher
*geplant* im Register als *verkörpert*, siehe §1 Ausschluss); (b) eine
benannte Beobachtung, falls die Extraktions-Pipe tatsächlich einen
maskierten Exit-Code zeigt (§6 Risiko 2) — welcher Kandidat trägt,
entscheidet der Lauf.

## 6. Risiken und offene Punkte

<!-- BEDIENHINWEIS: Was koennte schief gehen? Welche Carveouts entstehen
ggf.? Die drei Ausgaenge stehen als Form in der Zeile darunter. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Datei-Attribute/Zeilenenden-Divergenz durch `tar`-Extraktion.** Die neue
  Form (Erzeugung im Image, Übertragung als `tar`-Stream, Extraktion via
  Host-`tar`) könnte andere Datei-Attribute (Zeitstempel, Zeilenenden bei
  einem `tar`, der auf einem anderen Betriebssystem-Layer läuft) erzeugen
  als das heutige direkte In-Place-Schreiben in den Bind-Mount —
  `make generated-sync` vergleicht **byte-genau** (`cmp`). Muss real
  geprüft werden, bevor der Slice schließt (ein Lauf von
  `make proto-generate`, gefolgt von `make generated-sync`). —
  **Ausgang:** <bei Closure einzutragen>
- **Die Extraktions-Pipe kann den Exit-Code von `docker run` maskieren.**
  `docker run … | tar -x -C .` folgt genau dem in `AGENTS.md` §3.9
  beschriebenen Muster: Der Gesamt-Exit-Code der Pipe ist der von `tar`
  (letztes Glied), nicht der von `docker run`. Ein gescheiterter
  Stufen-Build oder ein abgebrochener `docker run` könnte hinter einem
  erfolgreichen, aber leeren oder unvollständigen `tar -x` verschwinden —
  eine verwandte Klasse zur irreführenden Docker-Fehlermeldung aus
  `slice-102`/`slice-103` (dort: „pull access denied" statt fehlendem
  Bau-Kontext; hier: möglich maskierter Exit-Code statt eines sichtbar
  roten Laufs). Muss durch `PIPESTATUS`/`pipefail` oder einen zweistufigen
  Aufruf (Stream erst in eine Datei, Exit-Code prüfen, dann extrahieren)
  abgesichert und real mit einem absichtlich fehlschlagenden Lauf getestet
  werden. — **Ausgang:** <bei Closure einzutragen>
- **Ein unvollständiger Träger-Suchlauf lässt eine Präsens-Beschreibung
  stehen.** Fünf Träger sind in LP3 benannt; ob der Implementer-Suchlauf
  einen sechsten findet oder einen der fünf übersieht, ist offen —
  dieselbe Klasse wie `BEO-PGC/arbeit-ueberholt-stehenden-traeger` (bereits
  5× belegt, in `AGENTS.md` §3.13 verkörpert). — **Ausgang:** <bei Closure
  einzutragen>

## 7. Closure-Notiz

<!-- BEDIENHINWEIS — keine Norm; faellt beim Kopieren weg (README.md
§Verwendung, Schritt 5) und darf deshalb nichts Tragendes halten. Reihenfolge:
diese Sektion vor dem `git mv` nach done/ fuellen — einzige Ausnahme ist das
letzte DoD-Item in §2 (die Paarungen suchen in `done/`, also nach dem `git mv`).
Im Repo ohne Wellen-Betrieb braucht die Closure dadurch drei Commits: Inhalt,
`git mv`, Haekchen — das folgt aus der Hard Rule, es widerspricht ihr nicht. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<NNN>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks).

- **Was hat funktioniert:** <bei Closure auszufüllen>
- **Was ging anders als geplant:** <bei Closure auszufüllen>
- **Steering-Loop-Eintrag:** <bei Closure auszufüllen — oder „keiner, der
  Normalfall"; siehe §5 Kandidaten>
- **Beobachtungs-Register (`../observations/`):** <bei Closure auszufüllen>
- **Folge-Slices:** <bei Closure auszufüllen — voraussichtlich keine>
- **Risiken aus §6:** <jedes mit genau einem Ausgang — siehe §6>
- **Drei Paarungen:** <Repo ohne Wellen-Betrieb — nach dem `git mv` nach
  `done/` zu prüfen>

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Der Abschnitt selbst entfällt nie.** Die zwei vorgelagerten Prüfungen laufen
in **jedem** Slice-Plan — sie hängen weder am Modus noch am Slice-Typ. Bedingt
ist allein der Modus-Begründungsblock am Ende; deshalb nennt der Titel beide
Hälften.

**Vorgelagert — Sub-Area-Wahl prüfen:** Berührt sind ausschließlich
`Makefile`, `Dockerfile` und Doku-Träger unter `harness/` — die repo-weite
Default-Sub-Area `*`/`PGC` aus der Modus-Deklaration in
[`harness/conventions.md`](../../../../harness/conventions.md). Keine
feinere Sub-Area ist im Repo deklariert, die diese Pfade eigens führt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen
(`grep -rli "proto-generate\|bind-mount\|tar.*stdout\|mount-los\|PROTO_RUN_USER"
docs/plan/planning/observations/`) — **keine Treffer**. Das einzige
existierende Verzeichnis (`BEO-PGC/generierte-artefakte-ohne-sync-sensor/`)
betrifft die Sync-Gate-Klasse (`ADR-0084`), nicht den Mount-Mechanismus von
`proto-generate` selbst; sein Zähler ist bereits verkörpert und wird von
diesem Slice nicht berührt. Keine Treffer sind ebenfalls eine Antwort.

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas sind **GF**. Der
Block pro Sub-Area entfällt; der **Abschnitt** bleibt.

<!-- Block für jede berührte Sub-Area duplizieren. Format identisch
mit dem im Baseline-Regelwerk §Ziel-Form: Sub-Area-Modus-Begründung
abgedruckten Block. -->

### Sub-Area: `*`/`PGC` (Default)

- **Modus:** GF
- **Konventionen-Dichte:** hoch — `AGENTS.md` §3.1 (Docker-only) und §3.9
  (Pipe/Exit-Code-Disziplin) verankern beide berührten Regeln bereits als
  Hard Rules; das Muster „gepinnte Toolchain-Stufe + `make`-Ziel" ist an
  drei Stellen etabliert (`proto-generate`, `generated-sync`,
  `schema-rollout`).
- **Phase-Reife:** Phase 5 — die Sub-Area trägt bereits mehrere
  ADR-gebundene Sensoren und Toolchain-Stufen; keine Erstanlage.
- **Evidenz-/Diskrepanz-Risiko:** niedrig — GF, Doc führt; keine
  Register-Treffer zu diesem Mechanismus (siehe oben).
- **Reconciliation-Aufwand:** entfällt (GF, kein Brownfield-Bestand).
