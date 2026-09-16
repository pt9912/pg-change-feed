# harness/mk/generated-sync.mk — Sync-Gate des generierten Protobuf-/gRPC-Codes
# (ADR-0084 Festlegung 1, Folgepflicht des Generators aus ADR-0060). Der
# gepinnte Generator (Dockerfile-Stufe `proto`) laeuft gegen die committeten
# .proto-Quellen und schreibt in ein Temp-Verzeichnis;
# tools/harness/generated-sync.sh vergleicht das Ergebnis mit dem committeten
# Erzeugnis und faerbt bei Abweichung mit Datei und Zeile rot. Der Baum haengt
# dabei als `:ro`-Bind-Mount im Generator-Container, das Ziel schreibt den
# Arbeitsbaum also nicht — darin unterscheidet es sich von `make
# proto-generate`, das in-place erzeugt und deshalb kein Pruef-Schritt ist.
.PHONY: generated-sync

GENERATED_SYNC_IMAGE ?= pg-change-feed:proto-sync

# GENERATED_SYNC_SOURCE_DIR/MODULE sind Overrides: leer heisst, das Skript
# leitet sie selbst ab (Quellverzeichnis `proto`, Modulpfad aus go.mod).
generated-sync: ## Sync-Gate: committeter Protobuf-Code gegen die Ausgabe des gepinnten Generators (ADR-0084)
	@GENERATED_SYNC_IMAGE=$(GENERATED_SYNC_IMAGE) \
	  GENERATED_SYNC_SOURCE_DIR=$(GENERATED_SYNC_SOURCE_DIR) \
	  GENERATED_SYNC_MODULE=$(GENERATED_SYNC_MODULE) \
	  bash tools/harness/generated-sync.sh

GATE_CHECKS += generated-sync
