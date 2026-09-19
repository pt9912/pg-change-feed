using System.Text.Json.Serialization;

namespace PgChangeFeed.Client.Http.Models;

/// <summary>
/// <c>RegisterConsumer</c> request — <c>POST /consumers</c> (SPEC-018), both
/// fields mandatory.
/// </summary>
public sealed record RegisterConsumerRequest(
    [property: JsonPropertyName("consumer_id")] string ConsumerId,
    [property: JsonPropertyName("name")] string Name);

/// <summary>
/// <c>RegisterConsumer</c> response (<c>201</c>) — <c>already_registered</c>
/// carries idempotency forward, there is no separate status code for it
/// (SPEC-018).
/// </summary>
public sealed record RegisterConsumerResponse(
    [property: JsonPropertyName("consumer_id")] string ConsumerId,
    [property: JsonPropertyName("name")] string Name,
    [property: JsonPropertyName("already_registered")] bool AlreadyRegistered);

/// <summary>
/// <c>AcknowledgeConsumer</c> request — <c>POST /consumers/acknowledge</c>
/// (SPEC-018), all three fields mandatory. A position older than the
/// current one, or one from a different source, ends <c>400</c>.
/// </summary>
public sealed record AcknowledgeConsumerRequest(
    [property: JsonPropertyName("consumer_id")] string ConsumerId,
    [property: JsonPropertyName("source_id")] string SourceId,
    [property: JsonPropertyName("offset")] ulong Offset);

/// <summary>
/// <c>AcknowledgeConsumer</c> response (<c>200</c>) — the position now in
/// effect after the acknowledgement.
/// </summary>
public sealed record AcknowledgeConsumerResponse(
    [property: JsonPropertyName("consumer_id")] string ConsumerId,
    [property: JsonPropertyName("source_id")] string SourceId,
    [property: JsonPropertyName("offset")] ulong Offset);

/// <summary>
/// <c>GetConsumerPosition</c> response (<c>200</c>) —
/// <c>acknowledged = false</c> reads the defined starting position without
/// any prior acknowledgement (SPEC-018).
/// </summary>
public sealed record ConsumerPositionResponse(
    [property: JsonPropertyName("consumer_id")] string ConsumerId,
    [property: JsonPropertyName("source_id")] string SourceId,
    [property: JsonPropertyName("offset")] ulong Offset,
    [property: JsonPropertyName("acknowledged")] bool Acknowledged);

/// <summary>
/// <c>RemoveConsumer</c> request — <c>POST /consumers/remove</c> (SPEC-018).
/// </summary>
public sealed record RemoveConsumerRequest(
    [property: JsonPropertyName("consumer_id")] string ConsumerId);

/// <summary>
/// <c>RemoveConsumer</c> response (<c>200</c>) — a never-registered consumer
/// reports <c>removed = false</c>, never <c>404</c> (idempotency instead of
/// an error against an unknown resource, SPEC-018).
/// </summary>
public sealed record RemoveConsumerResponse(
    [property: JsonPropertyName("consumer_id")] string ConsumerId,
    [property: JsonPropertyName("removed")] bool Removed);
