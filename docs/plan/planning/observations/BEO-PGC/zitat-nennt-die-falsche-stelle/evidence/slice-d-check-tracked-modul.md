# Beleg: slice-d-check-tracked-modul

Vorgang: `slice-d-check-tracked-modul` — `d-check`-Modul `tracked` aktiviert
(Welle `welle-d-check`).

Fund (Review F-2, HIGH): `harness/sensors/docs-check.md` §Bindung zitierte für
die Aussage „keine eigene ADR nötig" die Dokumente `ADR-0072`/`ADR-0075`
(Präzedenzmuster `hostpaths`) als Beleg. Beide zitierten ADRs sind aber genau
die Dokumente, die **weil** die `hostpaths`-Aktivierung eine ADR brauchte,
geschrieben wurden — als Beleg für „keine ADR nötig" tragen sie die
gegenteilige Aussage. Anders als beim Erstauftreten (`slice-090`, Übernahme aus
einer Bericht-Kopfzeile) und beim zweiten (`slice-102`, Verwechslung zweier
benachbarter Slices) ist der Ursprung hier eine **falsche Analogie-Achse**: der
Implementer verglich „ist eine Bündel-Aufnahme" statt „löst die Aktivierung
Folgepflichten aus, die eine ADR bräuchten" (Architect-Verdikt §5).

Gefunden hat es der Reviewer (F-2), unabhängig bestätigt und mit dem
tragenden Präzedenzfall (`structure`-Modul, Commit `f9e5a3c`) korrigiert vom
nachträglichen Architect-Zug (der Architect-Verdikt zur ADR-Frage von
`slice-d-check-tracked-modul`).
Behoben in der Fixrunde (Commit `5800a83`), Delta-Review bestätigt den
korrigierten Wortlaut.

**Dritter Vorgang — Schwelle erreicht.** Anders als bei den ersten beiden
Fällen (Zahlen-/Kennungs-Zitat) ist der zitierte Gegenstand hier eine
**Entscheidungslage** (zwei ADRs als Präzedenzmuster) statt eines
Abschnitts-/Slice-Verweises — dieselbe Klasse auf einem dritten Gegenstandstyp.
Da dieser Slice `welle-d-check` angehört, liegt der Lese-Schritt (Ausgang
zuweisen) bei der Welle-Closure (Modul 6/8), nicht bei dieser Slice-Closure.

Quelle: Review zu `slice-d-check-tracked-modul` F-2 ·
der Architect-Verdikt zur ADR-Frage von `slice-d-check-tracked-modul` ·
die Fixrunde des Review-Berichts zu `slice-d-check-tracked-modul`.
