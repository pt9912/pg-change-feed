# MR-002 — Slice-/Welle-Kennungen sind Namen, nicht Nummern (ab slice-105 exklusive)

Regeln dieser Datei: Pflichtfelder sind Datum, Geltungsbereich,
**Ersetzt-Baseline-Regel**, Adaption, Begründung und Auflösungs-Trigger.

- **Datum:** 2026-09-17
- **Geltungsbereich:** `harness/conventions.md` §Aktive Adaptionen; alle nach
  `slice-105` neu angelegten Slice- und Welle-Plan-Dateien dieses Repos.
  `slice-001`–`slice-105` sind von dieser Adaption nicht betroffen und bleiben
  numerisch (siehe Adaption unten).
- **Ersetzt-Baseline-Regel:** [`grundlagen-source-precedence.md`
  §Vergabe: woher die nächste Kennung
  kommt](../../.harness/baseline/v6.9.0/regelwerk/grundlagen-source-precedence.md#vergabe-woher-die-nächste-kennung-kommt)
  — konkret der Satz „Welle- und Slice-Kennungen sind Namen, nicht Nummern —
  unabhängig von der Schreiberzahl.“
- **Adaption:** Bestehende `slice-001`–`slice-104` (und `slice-105`, der
  Slice, der diese Adaption selbst einführt — die letzte nummerierte Kennung
  dieses Repos) bleiben **unangetastet** nummeriert; sie werden nicht
  umbenannt. **Ab dem nächsten neu angelegten Slice nach `slice-105`** gilt
  die Baseline-Regel unverändert und vollständig: Der Name trägt das Präfix
  eines vorhandenen Ankers (`LH-*`, `ADR-*`, `CO-*`), wenn einer existiert,
  sonst einen freien Slug — wörtlich die Baseline-Formulierung, keine
  abweichende Auslegung. Die Baseline definiert dazu: „Schreiber ist, was
  committet: ein Mensch, ein Agent, ein Automat.“ Die eigentliche Abweichung
  gegenüber der Baseline ist also nicht die künftige Namensform (die wird
  unverändert übernommen), sondern der **Bestandsschutz** für die 105
  bestehenden nummerierten Kennungen — ein bewusster, dauerhafter
  Altbestands-Schnitt.
- **Begründung:** Dieses Repo läuft — unabhängig davon, dass ein einzelner
  Mensch es lenkt — strukturell als Mehr-Schreiber-Repo im Sinn der
  Baseline-Definition: mehrere, teils parallel gestartete Claude-Agenten in
  den Rollen Architect, Implementer, Reviewer, Verifier und Planner committen
  jeweils eigenständig (z. B. die parallele Architect+Verifier-Sequenz bei
  `slice-101`). Die bisherige fortlaufende Nummernvergabe (`slice-<NNN>`)
  wurde nur durch manuelle Serialisierung kollisionsfrei gehalten — WIP-Limit
  1 pro Implementer-Rolleninhaber und ein `find`-Check auf die nächste freie
  Nummer vor jeder Neuanlage. Das Nummernschema selbst schützt nicht vor
  Kollision, wenn zwei Schreiber gleichzeitig eine neue Kennung ziehen; ein
  Name mit Anker-Präfix (`LH-*`/`ADR-*`/`CO-*`) oder freiem Slug macht eine
  Doppelvergabe an der Namensform selbst unwahrscheinlicher und an der
  `git`-Konfliktstelle sichtbar, statt sich auf durchgehaltene Disziplin zu
  verlassen. Eine Umbenennung der bestehenden 105 Kennungen wäre gegenüber
  dem erreichbaren Nutzen unverhältnismäßig disruptiv (siehe `slice-105` §1
  Abgrenzung: Querverweis-Bestand in ADRs, Beobachtungs-Register,
  Closure-Notizen, Review-/Verify-Reports, Commit-Messages) — daher der
  Bestandsschutz statt einer rückwirkenden Migration.
- **Auflösungs-Trigger:** permanent — die Entscheidung ist eine dauerhafte
  Konvention, kein Provisorium. Kein Re-Evaluierungs-Ereignis identifiziert;
  eine Auflösung wäre nur denkbar, wenn dieses Repo bewusst auf reinen
  Ein-Schreiber-Betrieb (ein einzelner committender Akteur, keine parallel
  agierenden Agentenrollen) umgestellt würde.
