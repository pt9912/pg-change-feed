using System.Text.Json;
using PgChangeFeed.Client.Nats.Models;
using Xunit;

namespace PgChangeFeed.Client.Tests.Nats;

/// <summary>
/// Message-schema completeness for <see cref="Change"/> — field for field,
/// analogous to
/// <c>PgChangeFeed.Client.Tests.Sse.ChangeMessageSchemaTests</c> on the SSE
/// client, here from the JSON payload a NATS stream message carries.
/// </summary>
public class ChangeMessageSchemaTests
{
    [Fact]
    public void AllTenFieldsRoundTrip()
    {
        const string json = """
            {
              "change_id": "change-1",
              "transaction_id": "tx-1",
              "source_table_id": "table-1",
              "sequence": 2,
              "operation": "UPDATE",
              "old_image": {"id": 1},
              "new_image": {"id": 1, "bestellstatus": "bezahlt"},
              "schema_version": "table-1-v1",
              "schema": "public",
              "table": "orders"
            }
            """;

        var change = JsonSerializer.Deserialize<Change>(json);

        Assert.NotNull(change);
        Assert.Equal("change-1", change!.ChangeId);
        Assert.Equal("tx-1", change.TransactionId);
        Assert.Equal("table-1", change.SourceTableId);
        Assert.Equal(2, change.Sequence);
        Assert.Equal("UPDATE", change.Operation);
        Assert.Equal("""{"id": 1}""", change.OldImage!.Value.GetRawText());
        Assert.Equal("""{"id": 1, "bestellstatus": "bezahlt"}""", change.NewImage!.Value.GetRawText());
        Assert.Equal("table-1-v1", change.SchemaVersion);
        Assert.Equal("public", change.Schema);
        Assert.Equal("orders", change.Table);
    }

    /// <summary>
    /// A defensive guard against silent schema drift: the stream message has
    /// exactly ten fields (the same ten as the SSE event). If a change added
    /// or removed a property without updating this test, the field-by-field
    /// assertion above could stay green while the schema no longer matches
    /// the stream — this count assertion is what would actually go red.
    /// </summary>
    [Fact]
    public void ChangeHasExactlyTenProperties()
    {
        Assert.Equal(10, typeof(Change).GetProperties().Length);
    }
}
