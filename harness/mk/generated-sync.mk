# harness/mk/generated-sync.mk — Sync-Gate des generierten Protobuf-/gRPC-Codes
# (ADR-0084 Festlegung 1, Folgepflicht des Generators aus ADR-0060). Der
# gepinnte Generator (Dockerfile-Stufe `proto-export`, geteilt mit `make
# proto-generate` seit slice-104) erzeugt den Code zur Build-Zeit und gibt
# ihn als `tar`-Stream aus; tools/harness/generated-sync.sh extrahiert
# host-seitig in ein Temp-Verzeichnis (kein Bind-Mount, kein
# `--user`-Workaround, seit slice-generated-sync-tar-export) und vergleicht
# das Ergebnis mit dem committeten Erzeugnis, faerbt bei Abweichung mit Datei
# und Zeile rot. Der Arbeitsbaum bleibt dabei unberuehrt — das Ziel schreibt
# ihn nicht; darin unterscheidet es sich von `make proto-generate`, das
# in-place erzeugt und deshalb kein Pruef-Schritt ist.
.PHONY: generated-sync

GENERATED_SYNC_IMAGE ?= pg-change-feed:proto-sync

# GENERATED_SYNC_SOURCE_DIR ist ein Override: leer heisst, das Skript leitet
# das Quellverzeichnis selbst ab (`proto`) — es ist seit dem Stufen-Wechsel
# nur noch ein Existenz-/Berichts-Check, kein Generierungs-Parameter mehr
# (die Stufe `proto-export` traegt Quelle und Modulpfad fest).
generated-sync: ## Sync-Gate: committeter Protobuf-Code gegen die Ausgabe des gepinnten Generators (ADR-0084)
	@GENERATED_SYNC_IMAGE=$(GENERATED_SYNC_IMAGE) \
	  GENERATED_SYNC_SOURCE_DIR=$(GENERATED_SYNC_SOURCE_DIR) \
	  bash tools/harness/generated-sync.sh

GATE_CHECKS += generated-sync
