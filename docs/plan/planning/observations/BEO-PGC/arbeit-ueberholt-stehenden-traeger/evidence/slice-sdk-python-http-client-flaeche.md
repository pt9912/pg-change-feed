# Beleg: slice-sdk-python-http-client-flaeche

Vorgang: `slice-sdk-python-http-client-flaeche` — liefert
`PgChangeFeedHttpClient` (zehn öffentliche Methoden, `SPEC-018`/`SPEC-022`)
als real existierende HTTP-API-Client-Fläche des Python-SDK-Packages.

Fund: Der Implementer führte denselben §3.13-Suchlauf durch, den das
C#-Geschwister-Vorkommen bereits gelehrt hatte
(`evidence/slice-sdk-csharp-http-client-flaeche.md`, 17. Beleg), und fand
proaktiv zwei Stellen (`README.md`, `__init__.py`) noch vor dem ersten
Review — ein realer Lernerfolg gegenüber dem C#-Zyklus, wo derselbe
README-Fund erst der Reviewer machte. Der Suchlauf selbst blieb dennoch in
zwei voneinander unabhängigen Dimensionen zugleich unvollständig: (1) sein
Datei-Glob (`sdks/python/README.md
sdks/python/pgchangefeed/src/pgchangefeed/*.py`) schloss `pyproject.toml`
aus, und (2) sein Wortmuster (`follow-up|added by`) deckte die tatsächlich
verwendete **deutsche** Formulierung „folgt erst" nicht ab.
`pyproject.toml:22-24` behauptete deshalb weiterhin „die HTTP-Client-Flaeche
selbst folgt erst mit slice-sdk-python-http-client-flaeche" — exakt der
Slice, der sie real liefert. Gefunden hat den dritten Treffer erst der
**Reviewer** (F-1, HIGH, merge-blockierend), über einen bewusst breiteren
`grep`, der über den dokumentierten Implementer-Suchlauf hinausging.

Behoben in der Fixrunde (`d0da688b`): der Kommentar trägt jetzt eine
indikative Ist-Zustands-Beschreibung ohne Slice-Namen. Fixrunden-Nachprüfung
und Verifikation bestätigen unabhängig voneinander — beide mit einem noch
breiteren Wortmuster, u. a. „folgt mit"/„steht noch aus"/„is added"/„will be
added" —, dass keine vierte Fundstelle verblieben ist.

**Warum dieser Beleg trotz proaktivem, aus dem Vorgänger-Fund abgeleitetem
Implementer-Suchlauf zählt:** Er zeigt, dass ein aus einem vorigen Vorkommen
kopierter Suchlauf zwei unabhängige Lücken zugleich tragen kann — eine
Datei-Glob-Lücke und eine Sprach-/Wortmuster-Lücke —, ohne dass eine der
beiden allein den Fund erklärt. Beide Achsen (welche Dateien, welche
Formulierungen) müssen unabhängig breit genug gewählt werden; ein aus dem
letzten Fund übernommenes Muster garantiert keine Vollständigkeit am
nächsten, strukturell ähnlichen Vorgang.

Quelle: `docs/reviews/review-slice-sdk-python-http-client-flaeche.md` <!-- d-check:status-provenance -->
F-1 und §Fixrunden-Nachprüfung F-1 · `docs/reviews/verifikation-slice-sdk-python-http-client-flaeche.md` <!-- d-check:status-provenance -->
§1.5 · Commit `d0da688b`.
