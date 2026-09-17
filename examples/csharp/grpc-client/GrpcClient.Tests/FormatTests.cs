using Cdc.Stream.V1;
using Google.Protobuf;
using Xunit;

namespace CdcExamples.Grpc.Tests;

/// <summary>
/// Prüft <see cref="Format.FormatChange"/>: Tabelle, Operation und
/// <c>change_id</c> stehen in der Ausgabezeile, das neue Row Image wird
/// unverändert angehängt — inklusive des leeren Falls bei <c>DELETE</c>.
/// Form-Vorbild: <c>examples/grpc-client/format_test.go</c>.
/// </summary>
public class FormatTests
{
    [Fact]
    public void FormatChange_CarriesIdentityAndPayload()
    {
        var change = new Change
        {
            ChangeId = "c-1",
            Schema = "public",
            Table = "orders",
            Operation = "INSERT",
            NewImage = ByteString.CopyFromUtf8("{\"id\":1}"),
        };

        var got = Format.FormatChange(change);

        Assert.Equal("grpc-client: change_id=c-1 table=public.orders operation=INSERT new_image={\"id\":1}", got);
    }

    [Fact]
    public void FormatChange_HandlesEmptyNewImage()
    {
        var change = new Change
        {
            ChangeId = "c-2",
            Schema = "public",
            Table = "orders",
            Operation = "DELETE",
        };

        var got = Format.FormatChange(change);

        Assert.Equal("grpc-client: change_id=c-2 table=public.orders operation=DELETE new_image=", got);
    }
}
