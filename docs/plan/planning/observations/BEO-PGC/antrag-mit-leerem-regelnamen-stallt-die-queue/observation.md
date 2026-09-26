# BEO-PGC/antrag-mit-leerem-regelnamen-stallt-die-queue

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft den Antragsweg der SQL-Administration, keine eigene Sub-Area im Sinn der Modus-Deklaration).

Die Beobachtung: Die administrativen SQL-Funktionen prüfen ihre Eingabe nicht (keine Domänenlogik in SQL); eine Zeile mit NULL oder leerem Regelnamen und mit SQL-NULL als Regelform entsteht als `pending`-Antrag. Der Antrags-Konstruktor lehnt sie beim **Lesen** der Queue ab, `ReadPendingRequests` gibt den Fehler zurück, `processAdministrationRequests` protokolliert ihn und liest im nächsten Durchlauf dieselbe Zeile erneut: kein Antrag der Queue, auch keiner einer anderen Art, wird verarbeitet, bis die Zeile entfernt oder ihr Status per `UPDATE` geändert ist. Die Spec sagt für diese Fälle `failed` mit einem Fehlertext; die Zeile soll gelesen und verarbeitet statt abgelehnt werden. Dieselbe Klasse besteht für `cdc.exclude_column(…, NULL)` und die Spaltenarten, deren Konstruktor eine leere Spalte beim Lesen ablehnt.

**Abgrenzung.** `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` betrifft Tests; hier ist die Ablehnung an der Lese-Stelle das Verhalten des Produktivpfads, die Stelle der Prüfung (Lesen statt Verarbeiten) die Ursache.

**Warum das zählt:** Ein einzelner ungültiger Antrag einer vertrauten Rolle (`cdc_admin`) hält die gesamte Administrations-Queue an; die Abhilfe verlangt SQL-Zugriff auf die Antrags-Tabelle, und die Ursache steht nur im Log.

Deklaration: `slice-transformationen-antragsweg-schema` (Review F-3).
