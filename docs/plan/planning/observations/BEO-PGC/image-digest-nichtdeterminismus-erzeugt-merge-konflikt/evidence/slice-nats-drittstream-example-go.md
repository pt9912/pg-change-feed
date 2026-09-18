**Vorgang:** slice-nats-drittstream-example-go
**Fund:** Ein `make image`-Lauf während der Implementierung committete einen
neuen Digest in `harness/image-hash.txt`, begründet mit einer falschen
Build-Kontext-Kausalität (Commit `7ea44409`, „examples/-Aenderungen liegen
im Build-Kontext"). Real widerlegt: `.dockerignore` schließt `examples/**`
vom Root-`Dockerfile`-Build-Kontext aus. Der Digest-Wechsel war reines
Builder-Rauschen zu unverändertem Binary-Inhalt (Review-Fund F-1,
zurückgenommen per `b67a1a63`). Der Auftraggeber benannte im Anschluss das
allgemeinere, wiederkehrende Ärgernis: ein solcher Diff stört beim Mergen
paralleler Branches, unabhängig davon, ob er berechtigt ist — Anlass dieser
Beobachtung.
