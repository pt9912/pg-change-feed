using System.Net;
using System.Text.Json.Serialization;

namespace CdcExamples.Http;

/// <summary>
/// RegisterConsumerRequest/-Response spiegeln den JSON-Vertrag von
/// <c>POST /consumers</c> (<c>LH-FA-CON-001</c>): beide Request-Felder
/// Pflicht, <c>AlreadyRegistered</c> trägt die Idempotenz-Antwort ohne
/// eigenen Statuscode. Form-Vorbild: <c>examples/http-client/consumer.go</c>.
/// </summary>
public sealed record RegisterConsumerRequest(
    [property: JsonPropertyName("consumer_id")] string ConsumerId,
    [property: JsonPropertyName("name")] string Name);

public sealed record RegisterConsumerResponse(
    [property: JsonPropertyName("consumer_id")] string ConsumerId,
    [property: JsonPropertyName("name")] string Name,
    [property: JsonPropertyName("already_registered")] bool AlreadyRegistered);

/// <summary>
/// AcknowledgeConsumerRequest/-Response spiegeln
/// <c>POST /consumers/acknowledge</c> (<c>LH-FA-CON-004</c>).
/// </summary>
public sealed record AcknowledgeConsumerRequest(
    [property: JsonPropertyName("consumer_id")] string ConsumerId,
    [property: JsonPropertyName("source_id")] string SourceId,
    [property: JsonPropertyName("offset")] ulong Offset);

public sealed record AcknowledgeConsumerResponse(
    [property: JsonPropertyName("consumer_id")] string ConsumerId,
    [property: JsonPropertyName("source_id")] string SourceId,
    [property: JsonPropertyName("offset")] ulong Offset);

/// <summary>
/// ConsumerPositionResponse spiegelt <c>GET /consumers/position</c>
/// (<c>LH-FA-CON-005</c>): <c>Acknowledged</c> unterscheidet die definierte
/// Anfangsposition (kein Nachweis) von einer echten Bestätigung mit
/// Offset 0.
/// </summary>
public sealed record ConsumerPositionResponse(
    [property: JsonPropertyName("consumer_id")] string ConsumerId,
    [property: JsonPropertyName("source_id")] string SourceId,
    [property: JsonPropertyName("offset")] ulong Offset,
    [property: JsonPropertyName("acknowledged")] bool Acknowledged);

/// <summary>
/// RemoveConsumerRequest/-Response spiegeln <c>POST /consumers/remove</c>
/// (<c>LH-FA-CON-006</c>): <c>Removed</c> trägt den Idempotenz-Ausgang — ein
/// nie registrierter Consumer meldet <c>false</c>, keinen <c>404</c>.
/// </summary>
public sealed record RemoveConsumerRequest(
    [property: JsonPropertyName("consumer_id")] string ConsumerId);

public sealed record RemoveConsumerResponse(
    [property: JsonPropertyName("consumer_id")] string ConsumerId,
    [property: JsonPropertyName("removed")] bool Removed);

/// <summary>
/// ConsumerPositionUrlBuilder baut die Lese-Adresse von
/// <c>GET /consumers/position</c>: ihr einziges Pflichtfeld ist
/// <c>consumer_id</c>. Form-Vorbild: <c>ConsumerPositionURL</c> in
/// <c>examples/http-client/consumer.go</c>.
/// </summary>
public static class ConsumerPositionUrlBuilder
{
    public static string Build(string addr, string consumerId) =>
        $"http://{addr}/consumers/position?consumer_id={Uri.EscapeDataString(consumerId)}";
}

/// <summary>
/// ConsumerClient ruft die vier Consumer-Verwaltungs-Fähigkeiten der
/// HTTP-/JSON-API auf. Form-Vorbild: <c>examples/http-client/consumer.go</c>.
/// </summary>
public static class ConsumerClient
{
    public static Task<RegisterConsumerResponse> RegisterConsumerAsync(
        System.Net.Http.HttpClient httpClient, Config cfg, CancellationToken cancellationToken = default) =>
        RequestHelper.SendJsonAsync<RegisterConsumerResponse>(
            httpClient, HttpMethod.Post, $"http://{cfg.Addr}/consumers", cfg.AdminToken,
            new RegisterConsumerRequest(cfg.ConsumerId, cfg.Name), HttpStatusCode.Created, cancellationToken);

    public static Task<AcknowledgeConsumerResponse> AcknowledgeConsumerAsync(
        System.Net.Http.HttpClient httpClient, Config cfg, CancellationToken cancellationToken = default) =>
        RequestHelper.SendJsonAsync<AcknowledgeConsumerResponse>(
            httpClient, HttpMethod.Post, $"http://{cfg.Addr}/consumers/acknowledge", cfg.AdminToken,
            new AcknowledgeConsumerRequest(cfg.ConsumerId, cfg.Source, cfg.Offset), HttpStatusCode.OK, cancellationToken);

    public static Task<ConsumerPositionResponse> ConsumerPositionAsync(
        System.Net.Http.HttpClient httpClient, Config cfg, CancellationToken cancellationToken = default) =>
        RequestHelper.SendJsonAsync<ConsumerPositionResponse>(
            httpClient, HttpMethod.Get, ConsumerPositionUrlBuilder.Build(cfg.Addr, cfg.ConsumerId), cfg.Token,
            null, HttpStatusCode.OK, cancellationToken);

    public static Task<RemoveConsumerResponse> RemoveConsumerAsync(
        System.Net.Http.HttpClient httpClient, Config cfg, CancellationToken cancellationToken = default) =>
        RequestHelper.SendJsonAsync<RemoveConsumerResponse>(
            httpClient, HttpMethod.Post, $"http://{cfg.Addr}/consumers/remove", cfg.AdminToken,
            new RemoveConsumerRequest(cfg.ConsumerId), HttpStatusCode.OK, cancellationToken);
}
