using Cdc.Stream.V1;
using Google.Protobuf;
using Xunit;

namespace PgChangeFeed.Client.Tests.Grpc;

/// <summary>
/// Message-schema completeness for the generated <see cref="Change"/> stub —
/// field for field, from the consumer side: every one of the ten fields
/// round-trips through the generated stub unchanged.
/// </summary>
public class ChangeMessageSchemaTests
{
    [Fact]
    public void AllTenFieldsRoundTrip()
    {
        var change = new Change
        {
            ChangeId = "change-1",
            TransactionId = "tx-1",
            SourceTableId = "table-1",
            Sequence = 2,
            Operation = "UPDATE",
            OldImage = ByteString.CopyFromUtf8("""{"id":1}"""),
            NewImage = ByteString.CopyFromUtf8("""{"id":1,"bestellstatus":"bezahlt"}"""),
            SchemaVersion = "table-1-v1",
            Schema = "public",
            Table = "orders",
        };

        Assert.Equal("change-1", change.ChangeId);
        Assert.Equal("tx-1", change.TransactionId);
        Assert.Equal("table-1", change.SourceTableId);
        Assert.Equal(2, change.Sequence);
        Assert.Equal("UPDATE", change.Operation);
        Assert.Equal("""{"id":1}""", change.OldImage.ToStringUtf8());
        Assert.Equal("""{"id":1,"bestellstatus":"bezahlt"}""", change.NewImage.ToStringUtf8());
        Assert.Equal("table-1-v1", change.SchemaVersion);
        Assert.Equal("public", change.Schema);
        Assert.Equal("orders", change.Table);
    }

    /// <summary>
    /// A defensive guard against silent schema drift: <c>Change</c> has
    /// exactly ten fields. If a <c>.proto</c> change added or removed a field
    /// without updating this test, the field-by-field assertion above could
    /// stay green while the schema no longer matches the stream — this count
    /// assertion is what would actually go red.
    /// </summary>
    [Fact]
    public void ChangeDescriptorHasExactlyTenFields()
    {
        Assert.Equal(10, Change.Descriptor.Fields.InFieldNumberOrder().Count);
    }
}
