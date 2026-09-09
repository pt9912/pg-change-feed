# `make a-check` — prüft die Hexagon-Schichten-Edges gegen den Go-Baum

## Vertrag

Wird dieses Target rot, verletzt der Go-Baum eine der deklarierten
Schichten-Regeln in `.a-check.yml` — die Maschinenform der
Schichten-Constraints aus [`spec/architecture.md` §2](../../spec/architecture.md):
Abhängigkeiten zeigen nach innen (`wrong-direction`, `core-impurity`,
`app-impurity`), Adapter sind lateral getrennt (`lateral-adapter`), Technik
leakt nicht in Ports (`tech-leak`, `port-impurity`), Composition Root außer
`bootstrap`/`cmd` verdrahtet nicht (`construct-leak`).

## Grenze — was das Grün nicht abdeckt

1. **Abdeckung folgt dem Baum** — Layer-Globs matchen nur vorhandene Dateien;
   eine Datei ohne Schicht fällt unter den Abdeckungs-Hinweis (kein
   Exit-Wechsel, aber auf stderr lesbar). Heilbar durch Glob-Nachzug in
   `.a-check.yml` je Slice.
2. **Auflösungs-Hinweis ist kein Grün** — löst kein Symbol auf eine Schicht
   auf, meldet der Lauf einen Hinweis; „alles grün“ sagt dann nichts über
   einen Prüfbereich. Beobachten, nicht durchwinken.
3. **Nur Pfad-/Import-Ebene, nur Go** — die `time`-Import-Regel der
   [`ADR-0040`](../../docs/plan/adr)-Fitness (Domain/Use-Cases importieren `time` nicht) ist nicht
   pfadgetrieben ausdrückbar; sie bleibt Review-Prüfpflicht bis zum
   Re-Evaluierungs-Trigger von [`ADR-0041`](../../docs/plan/adr/README.md).
   Permanent bis zu ergänzendem Tooling.
4. **Kein Runtime-Urteil** — der Lauf liest Quellen; Verhalten (Reihenfolge,
   Persistenz) prüfen Tests, nicht a-check. Permanent.

**Wie groß der Ausschnitt ist, sagt das Kommando, nicht diese Datei:**
`docker run --rm --network none -v "$PWD:/src:ro" <A_CHECK_IMAGE> /src`
über die `layers`-Globs der `.a-check.yml`; die Zeile „gesamt: 0 Befund(e)“
sagt etwas über diesen Ausschnitt, nicht über das Repo.

## Ausgabe und Ausgänge

| Exit | Bedeutung |
|---|---|
| 0 | keine Verstöße im Ausschnitt |
| 1 | mindestens ein Befund (`pfad:zeile: regel: meldung`) |
| 2 | Nutzungs-/Konfigurationsfehler — `.a-check.yml` strikt dekodiert, unbekannte Schlüssel brechen ab |

Hinweise (Abdeckung, Grenze, Auflösung) sind kein Befund und wechseln den
Exit nicht — sie stehen auf stderr und sind zu lesen, nicht zu ignorieren.

## Sperren

- Kein gültiges Image-Digest (`A_CHECK_IMAGE`-Platzhalter) — das Fragment
  bricht ab, bevor es scannt; Pin-Hebung = bewusster Commit.

## Bindung

[`ADR-0041`](../../docs/plan/adr/README.md) (Maschinenform der §2-Constraints,
supersedes [`ADR-0036`](../../docs/plan/adr)) · `.a-check.yml` (deklarativer Stand) · Image-Digest in
`a-check.mk` (Pin-Hebung = bewusster Commit).
