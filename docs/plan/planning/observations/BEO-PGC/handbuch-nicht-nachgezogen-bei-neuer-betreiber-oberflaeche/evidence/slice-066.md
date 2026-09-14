**Vorgang:** slice-066
**Fund:** Der Slice führte `cdc.exclude_column(...)`/`cdc.include_column(...)`
als administrative SQL-Funktionen ein — das Handbuch dokumentiert die
Geschwister `cdc.enable_table`/`cdc.disable_table` in §4 ausführlich, die
Spaltenfunktionen kommen im ganzen Dokument nicht vor (kein Treffer für
`cdc.exclude_column`, geprüft 2026-09-14). Die Versionshistorie wurde in
diesem Slice ebenfalls nicht angefasst.
