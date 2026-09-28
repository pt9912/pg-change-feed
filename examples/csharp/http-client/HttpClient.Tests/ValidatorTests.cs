using Xunit;

namespace CdcExamples.Http.Tests;

/// <summary>
/// Prüft <see cref="Validator.Validate"/> — netzlos, reine Funktion.
/// Form-Vorbild: <c>examples/http-client/main_test.go</c>.
/// </summary>
public class ValidatorTests
{
    private static Config BaseConfig(string verb) => new(
        Addr: "", Token: "", AdminToken: "", Verb: verb,
        Source: "", Publication: "", ConsumerId: "", Name: "", Offset: 0,
        Schema: "", Table: "", TableId: "", SchemaVersionId: "", Version: 1,
        From: "", To: "", Limit: "", MinAgeNanos: 0);

    [Fact]
    public void ValidateRejectsUnknownVerb()
    {
        var cfg = BaseConfig("unbekannt") with { Addr = "feed:8080" };

        var err = Validator.Validate(cfg);

        Assert.NotNull(err);
        Assert.Contains("unbekanntes --verb", err);
    }

    [Fact]
    public void ValidateRequiresAddrBeforeVerbFields()
    {
        var err = Validator.Validate(BaseConfig("tables"));

        Assert.NotNull(err);
        Assert.Contains("CDC_HTTP_ADDR", err);
    }

    [Fact]
    public void ValidateEnableTableAllowsDefaultTableId()
    {
        var cfg = BaseConfig("enable-table") with
        {
            Addr = "feed:8080", AdminToken = "admin-token",
            Source = "quelle-1", Schema = "public", Table = "orders", Publication = "pub_quelle_1",
        };

        Assert.Null(Validator.Validate(cfg));
    }

    [Fact]
    public void ValidateRetentionRunAllowsZeroMinAge()
    {
        var cfg = BaseConfig("retention-run") with { Addr = "feed:8080", AdminToken = "admin-token", Source = "quelle-1" };

        Assert.Null(Validator.Validate(cfg));
    }

    [Fact]
    public void ValidateAcknowledgeRequiresAdminTokenNotReaderToken()
    {
        var cfg = BaseConfig("acknowledge") with
        {
            Addr = "feed:8080", Token = "reader-token", ConsumerId = "consumer-1", Source = "quelle-1",
        };

        var err = Validator.Validate(cfg);

        Assert.NotNull(err);
        Assert.Contains("Admin-Token", err);
    }

    [Fact]
    public void ValidateChangesRequiresOnlySource()
    {
        var cfg = BaseConfig("changes") with { Addr = "feed:8080", Token = "reader-token", Source = "quelle-1" };

        Assert.Null(Validator.Validate(cfg));
    }

    [Fact]
    public void ValidateConsumerPositionRequiresReaderTokenNotAdminToken()
    {
        var cfg = BaseConfig("consumer-position") with { Addr = "feed:8080", AdminToken = "admin-token", ConsumerId = "c-1" };

        var err = Validator.Validate(cfg);

        Assert.NotNull(err);
        Assert.Contains("CDC_API_TOKEN_READER", err);
    }
}
