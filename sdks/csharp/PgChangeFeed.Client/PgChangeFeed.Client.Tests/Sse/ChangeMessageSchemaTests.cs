using System.Text.Json;
using PgChangeFeed.Client.Sse.Models;
using Xunit;

namespace PgChangeFeed.Client.Tests.Sse;

/// <summary>
/// Message-schema completeness for <see cref="Change"/> (SPEC-021) — field
/// for field, analogous to
/// <c>PgChangeFeed.Client.Tests.Grpc.ChangeMessageSchemaTests</c> on the
/// gRPC surface, here from the JSON wire form the SSE frame's <c>data:</c>
/// line carries.
/// </summary>
public class ChangeMessageSchemaTests
{
    [Fact]
    public void AllTenSpec021FieldsRoundTrip()
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
    /// A defensive guard against silent schema drift: SPEC-021 names exactly
    /// ten fields. If a future change added or removed a property without
    /// updating this test, the field-by-field assertion above could stay
    /// green while the schema no longer matches SPEC-021 — this count
    /// assertion is what would actually go red.
    /// </summary>
    [Fact]
    public void ChangeHasExactlyTenProperties()
    {
        Assert.Equal(10, typeof(Change).GetProperties().Length);
    }
}
