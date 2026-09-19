# Beleg: slice-release-multi-arch-image

Bei der lokalen Verifikation des Multi-Platform-Buildx-Builds
(`--platform linux/amd64,linux/arm64 --push`) fehlte `linux/amd64`
zunächst als natives Bau-Ziel auf dem Colima-Setup. Reihenfolge und
Befund:

1. `docker run --privileged --rm tonistiigi/binfmt --install all` —
   registrierte real alle Emulatoren (u. a. `rosetta`), `linux/amd64`
   blieb aber weiterhin außerhalb der vom Buildx-Builder gemeldeten
   Plattform-Liste.
2. `docker run --rm --privileged multiarch/qemu-user-static --reset -p
   yes` — als naheliegender zweiter Versuch. Beschädigte real den
   gesamten lokalen Docker-Runtime: jeder folgende Container-Start
   scheiterte mit `exec format error`, auch für triviale, native Images
   (`hello-world`, `docker run --rm hello-world`) und sogar für den
   Versuch, das Problem per weiterem `docker run` zu beheben.
3. Abhilfe: `colima restart` — einziger real wirksamer Weg, keine
   Handarbeit am laufenden System nötig. Danach lief `linux/amd64` real
   über Rosetta (`docker run --platform linux/amd64 alpine uname -m` →
   `x86_64`), und der eigentliche Multi-Platform-Push-Test gelang
   vollständig — obwohl `docker buildx inspect --bootstrap` `linux/amd64`
   bis zuletzt nicht in der Platform-Liste des Builders zeigte.

Erster Beleg dieser Beobachtung — unter der 3×-Schärfungsschwelle, kein
neuer Sensor oder verkörperte Regel ausgelöst.
