using Xunit;

namespace CdcExamples.Nats.Tests;

/// <summary>
/// Prüft die Ableitung des Wecksignal-Subjekts aus den drei Bestandteilen
/// (<c>SPEC-017</c>) — netzlos, reine Funktion. Form-Vorbild:
/// <c>examples/nats-client/subject_test.go</c>.
/// </summary>
public class SubjectTests
{
    [Fact]
    public void BuildDerivesTableGranularSubject()
    {
        var got = Subject.Build("quelle-1", "public", "orders");
        Assert.Equal("cdc.changes.quelle-1.public.orders", got);
    }

    [Fact]
    public void BuildKeepsDistinctTables()
    {
        var orders = Subject.Build("quelle-1", "public", "orders");
        var invoices = Subject.Build("quelle-1", "public", "invoices");
        Assert.NotEqual(orders, invoices);
    }
}
