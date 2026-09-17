using Xunit;

namespace CdcExamples.Http.Tests;

/// <summary>
/// Prüft den Aufbau der HTTP-Abfrage (<c>LH-FA-SST-006</c>) — netzlos, reine
/// Funktion. Form-Vorbild: <c>examples/http-client/tables_test.go</c>.
/// </summary>
public class TablesUrlBuilderTests
{
    [Fact]
    public void BuildCarriesBothRequiredFilters()
    {
        var got = TablesUrlBuilder.Build("feed:8080", "quelle-1", "pub_quelle_1");
        Assert.Equal("http://feed:8080/tables?publication=pub_quelle_1&source=quelle-1", got);
    }

    [Fact]
    public void BuildKeepsHostBoundary()
    {
        var got = TablesUrlBuilder.Build("localhost:9090", "src-e2e", "cdc_pub");
        Assert.Equal("http://localhost:9090/tables?publication=cdc_pub&source=src-e2e", got);
    }

    [Fact]
    public void BuildEscapesReservedCharacters()
    {
        var got = TablesUrlBuilder.Build("feed:8080", "quelle & test", "pub/1");
        Assert.Equal("http://feed:8080/tables?publication=pub%2F1&source=quelle%20%26%20test", got);
    }
}
