namespace CdcExamples.Grpc;

/// <summary>
/// Config trägt die Laufzeit-Eingabe des Beispiels: die Horch-Adresse des
/// gRPC-Servers, die zwei Token-Klassen und die Felder aller zwölf
/// Fähigkeiten — je Verb prüft <see cref="Validator.Validate"/> nur die
/// tatsächlich nötigen Felder. <c>Schema</c>/<c>Table</c> dienen doppelt: als
/// optionaler Stream-Filter (<c>ADR-0133</c>) und als Tabellen-Identität der
/// Verwaltungs-RPCs; <c>Target</c> (<c>--target</c>) wählt das Zustellziel einer
/// Change im Stream und in <c>read-changes</c>, leer ist kein Filter;
/// <c>CaFile</c> (<c>--ca-file</c>, <c>CDC_TLS_CA_FILE</c>) nennt das
/// Zertifikat (PEM) als Vertrauensanker — gesetzt spricht der Kanal TLS, leer
/// Klartext (<see cref="ChannelFactory"/>). Form-Vorbild: <c>examples/grpc-client</c> (Go),
/// <c>config</c> in <c>main.go</c>.
/// </summary>
public sealed record Config(
    string Addr,
    string Token,
    string AdminToken,
    string Verb,
    string Schema,
    string Table,
    string ConsumerId,
    string Name,
    ulong Offset,
    string TableId,
    string SchemaVersionId,
    long Version,
    string Source,
    string Publication,
    ulong From,
    ulong To,
    long Limit,
    long MinAgeNanos,
    string Target = "",
    string CaFile = "");
