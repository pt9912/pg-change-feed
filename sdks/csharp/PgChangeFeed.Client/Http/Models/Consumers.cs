using System.Text.Json.Serialization;

namespace PgChangeFeed.Client.Http.Models;

/// <summary>
/// The consumer to register: a unique <c>consumer_id</c> and a display
/// <c>name</c>, both mandatory.
/// </summary>
public sealed record RegisterConsumerRequest(
    [property: JsonPropertyName("consumer_id")] string ConsumerId,
    [property: JsonPropertyName("name")] string Name);

/// <summary>
/// The registered consumer. <c>already_registered</c> is true when it existed
/// before; there is no separate status code for that.
/// </summary>
public sealed record RegisterConsumerResponse(
    [property: JsonPropertyName("consumer_id")] string ConsumerId,
    [property: JsonPropertyName("name")] string Name,
    [property: JsonPropertyName("already_registered")] bool AlreadyRegistered);

/// <summary>
/// The position to store: <c>offset</c> is the <c>commit_position</c> of the
/// last change the consumer has processed for <c>source_id</c>; all three
/// fields are mandatory. A position older than the current one, or one from a
/// different source, ends <c>400</c>.
/// </summary>
public sealed record AcknowledgeConsumerRequest(
    [property: JsonPropertyName("consumer_id")] string ConsumerId,
    [property: JsonPropertyName("source_id")] string SourceId,
    [property: JsonPropertyName("offset")] ulong Offset);

/// <summary>
/// The position now stored for the consumer after the acknowledgement.
/// </summary>
public sealed record AcknowledgeConsumerResponse(
    [property: JsonPropertyName("consumer_id")] string ConsumerId,
    [property: JsonPropertyName("source_id")] string SourceId,
    [property: JsonPropertyName("offset")] ulong Offset);

/// <summary>
/// The stored position of a consumer. <c>acknowledged = false</c> reads the
/// defined starting position of a consumer that never acknowledged.
/// </summary>
public sealed record ConsumerPositionResponse(
    [property: JsonPropertyName("consumer_id")] string ConsumerId,
    [property: JsonPropertyName("source_id")] string SourceId,
    [property: JsonPropertyName("offset")] ulong Offset,
    [property: JsonPropertyName("acknowledged")] bool Acknowledged);

/// <summary>
/// The consumer to remove.
/// </summary>
public sealed record RemoveConsumerRequest(
    [property: JsonPropertyName("consumer_id")] string ConsumerId);

/// <summary>
/// A consumer that was never registered reports <c>removed = false</c>, never
/// <c>404</c>.
/// </summary>
public sealed record RemoveConsumerResponse(
    [property: JsonPropertyName("consumer_id")] string ConsumerId,
    [property: JsonPropertyName("removed")] bool Removed);
