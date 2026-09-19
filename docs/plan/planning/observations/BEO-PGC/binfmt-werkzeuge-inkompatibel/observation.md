# BEO-PGC/binfmt-werkzeuge-inkompatibel

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die
lokale Docker-Entwicklungsumgebung, keine eigene Sub-Area im Sinn der
Modus-Deklaration).

Die Beobachtung: Zwei gängige, beide öffentlich empfohlene Docker-
Binfmt-Registrierungswerkzeuge (`tonistiigi/binfmt` und
`multiarch/qemu-user-static`) sind auf demselben Colima-Setup
**nicht** nacheinander ausführbar, ohne den lokalen Docker-Runtime real
zu beschädigen. Ein `docker run --privileged --rm tonistiigi/binfmt
--install all` registrierte zunächst alle Emulatoren erfolgreich
(inklusive `rosetta`), zeigte aber `linux/amd64` weiterhin nicht als vom
Buildx-Builder unterstützte Plattform. Der naheliegende zweite Versuch,
`docker run --rm --privileged multiarch/qemu-user-static --reset -p
yes`, überschrieb daraufhin die `binfmt_misc`-Registrierungen so, dass
**jeder** Container-Start scheiterte — auch für triviale, native Images
(`hello-world`) — mit `exec format error` auf den
`containerd-shim`-Prozess selbst, nicht nur auf das Zielimage. Ein
weiterer `docker run`-Aufruf zur Fehlerbehebung war dadurch selbst nicht
mehr möglich (derselbe Fehler). Einzige real wirksame Abhilfe: `colima
restart` — setzt den VM-Kernelzustand inklusive `binfmt_misc` sauber
zurück, keine Handarbeit im laufenden System nötig. Danach funktionierte
`linux/amd64` real über Rosetta, ohne dass der Buildx-Builder es in
seiner eigenen Plattform-Liste (`docker buildx inspect --bootstrap`)
anzeigte — die Laufzeit-Emulation über Rosetta ist von dieser Anzeige
unabhängig.

Deklaration: `release-multi-arch-image` §6, real aufgetreten bei der
lokalen Verifikation des Multi-Platform-Builds.
