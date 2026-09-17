using Xunit;

namespace CdcExamples.Nats.Tests;

/// <summary>
/// Prüft den Aufbau der HTTP-Abfrage des zweiseitigen Ablaufs
/// (<c>LH-FA-SST-006</c>, <c>ADR-0079</c>) — netzlos, reine Funktion.
/// Form-Vorbild: <c>examples/nats-client/subject_test.go</c>
/// (<c>TestChangesURL...</c>).
/// </summary>
public class ChangesUrlBuilderTests
{
    [Fact]
    public void BuildCarriesAllThreeFiltersOnReadPath()
    {
        var got = ChangesUrlBuilder.Build("feed:8080", "quelle-1", "public", "orders");
        Assert.Equal("http://feed:8080/changes?schema=public&source=quelle-1&table=orders", got);
    }

    [Fact]
    public void BuildKeepsHostBoundary()
    {
        var got = ChangesUrlBuilder.Build("localhost:9090", "src-e2e", "public", "feed_e2e");
        Assert.Equal("http://localhost:9090/changes?schema=public&source=src-e2e&table=feed_e2e", got);
    }

    [Fact]
    public void BuildEscapesReservedCharacters()
    {
        var got = ChangesUrlBuilder.Build("feed:8080", "quelle & test", "sch/1", "tab 1");
        Assert.Equal("http://feed:8080/changes?schema=sch%2F1&source=quelle%20%26%20test&table=tab%201", got);
    }
}
