namespace CdcExamples.Http;

/// <summary>
/// Config trägt die Laufzeit-Eingabe des Beispiels: die Horch-Adresse der
/// HTTP-API, die zwei Token-Klassen (<c>CDC_API_TOKEN_READER</c>/
/// <c>CDC_API_TOKEN_ADMIN</c>) und die Felder aller zehn Fähigkeiten — je
/// Verb prüft <see cref="Validator.Validate"/> nur die tatsächlich nötigen
/// Felder. Form-Vorbild: <c>examples/http-client</c> (Go), <c>config</c> in
/// <c>main.go</c>.
/// </summary>
public sealed record Config(
    string Addr,
    string Token,
    string AdminToken,
    string Verb,
    string Source,
    string Publication,
    string ConsumerId,
    string Name,
    ulong Offset,
    string Schema,
    string Table,
    string TableId,
    string SchemaVersionId,
    long Version,
    string From,
    string To,
    string Limit,
    long MinAgeNanos,
    string Target = "");
