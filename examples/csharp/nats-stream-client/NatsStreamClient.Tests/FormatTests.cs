using System.Text.Json;

using Xunit;

namespace CdcExamples.NatsStream.Tests;

/// <summary>
/// Prüft <see cref="Format.FormatChange"/>: Tabelle, Operation und
/// <c>change_id</c> stehen in der Ausgabezeile, das neue Row Image wird
/// unverändert angehängt — inklusive des <c>DELETE</c>-Falls, in dem das
/// Row Image auf dem Draht das JSON-Literal <c>null</c> trägt (dieselbe
/// Übersetzung wie <c>natsstream.rowImage</c> auf der Erzeugerseite). Form-
/// Vorbild: <c>examples/csharp/grpc-client/GrpcClient.Tests/FormatTests.cs</c>.
/// </summary>
public class FormatTests
{
    private static JsonElement Json(string raw) => JsonDocument.Parse(raw).RootElement;

    [Fact]
    public void FormatChange_CarriesIdentityAndPayload()
    {
        var change = new StreamMessage(
            ChangeId: "c-1",
            TransactionId: "t-1",
            SourceTableId: "st-1",
            Sequence: 1,
            Operation: "INSERT",
            OldImage: Json("null"),
            NewImage: Json("{\"id\":1}"),
            SchemaVersion: "v1",
            Schema: "public",
            Table: "orders");

        var got = Format.FormatChange(change);

        Assert.Equal("nats-stream-client: change_id=c-1 table=public.orders operation=INSERT new_image={\"id\":1}", got);
    }

    [Fact]
    public void FormatChange_HandlesMissingNewImage()
    {
        var change = new StreamMessage(
            ChangeId: "c-2",
            TransactionId: "t-1",
            SourceTableId: "st-1",
            Sequence: 2,
            Operation: "DELETE",
            OldImage: Json("{\"id\":1}"),
            NewImage: Json("null"),
            SchemaVersion: "v1",
            Schema: "public",
            Table: "orders");

        var got = Format.FormatChange(change);

        Assert.Equal("nats-stream-client: change_id=c-2 table=public.orders operation=DELETE new_image=null", got);
    }
}
