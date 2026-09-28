using Xunit;

namespace CdcExamples.Grpc.Tests;

/// <summary>
/// Prüft <see cref="Validator.Validate"/> — netzlos, reine Funktion. Deckt
/// für jedes der zwölf Verben die Rechtsklassen-Bindung (das jeweils falsche
/// Token allein genügt nicht) und die dokumentierten Pflichtfeld-/Default-
/// Ausnahmen. Form-Vorbild: <c>examples/grpc-client/main_test.go</c>.
/// </summary>
public class ValidatorTests
{
    private static Config BaseConfig(string verb) => new(
        Addr: "", Token: "", AdminToken: "", Verb: verb,
        Schema: "", Table: "", ConsumerId: "", Name: "", Offset: 0,
        TableId: "", SchemaVersionId: "", Version: 1,
        Source: "", Publication: "", From: 0, To: 0, Limit: 0, MinAgeNanos: 0);

    [Fact]
    public void ValidateRejectsUnknownVerb()
    {
        var cfg = BaseConfig("unbekannt") with { Addr = "feed:9090" };

        var err = Validator.Validate(cfg);

        Assert.NotNull(err);
        Assert.Contains("unbekanntes --verb", err);
    }

    [Fact]
    public void ValidateRequiresAddrBeforeVerbFields()
    {
        var err = Validator.Validate(BaseConfig("list-tables"));

        Assert.NotNull(err);
        Assert.Contains("CDC_GRPC_ADDR", err);
    }

    [Fact]
    public void ValidateStreamDefaultAcceptsReaderTokenWithoutFilter()
    {
        var cfg = BaseConfig("stream") with { Addr = "feed:9090", Token = "reader-token" };

        Assert.Null(Validator.Validate(cfg));
    }

    [Fact]
    public void ValidateStreamRequiresReaderTokenNotAdminToken()
    {
        var cfg = BaseConfig("stream") with { Addr = "feed:9090", AdminToken = "admin-token" };

        var err = Validator.Validate(cfg);

        Assert.NotNull(err);
        Assert.Contains("CDC_API_TOKEN_READER", err);
    }

    [Fact]
    public void ValidateEnableTableAllowsDefaultTableId()
    {
        var cfg = BaseConfig("enable-table") with
        {
            Addr = "feed:9090", AdminToken = "admin-token",
            Source = "quelle-1", Schema = "public", Table = "orders", Publication = "pub_quelle_1",
        };

        Assert.Null(Validator.Validate(cfg));
    }

    [Fact]
    public void ValidateRunRetentionAllowsZeroMinAge()
    {
        var cfg = BaseConfig("run-retention") with { Addr = "feed:9090", AdminToken = "admin-token", Source = "quelle-1" };

        Assert.Null(Validator.Validate(cfg));
    }

    [Fact]
    public void ValidateAcknowledgeConsumerRequiresAdminTokenNotReaderToken()
    {
        var cfg = BaseConfig("acknowledge-consumer") with
        {
            Addr = "feed:9090", Token = "reader-token", ConsumerId = "consumer-1", Source = "quelle-1",
        };

        var err = Validator.Validate(cfg);

        Assert.NotNull(err);
        Assert.Contains("Admin-Token", err);
    }

    [Fact]
    public void ValidateRegisterConsumerRequiresAdminTokenNotReaderToken()
    {
        var cfg = BaseConfig("register-consumer") with
        {
            Addr = "feed:9090", Token = "reader-token", ConsumerId = "consumer-1", Name = "Consumer",
        };

        var err = Validator.Validate(cfg);

        Assert.NotNull(err);
        Assert.Contains("Admin-Token", err);
    }

    [Fact]
    public void ValidateRemoveConsumerRequiresAdminTokenNotReaderToken()
    {
        var cfg = BaseConfig("remove-consumer") with { Addr = "feed:9090", Token = "reader-token", ConsumerId = "consumer-1" };

        var err = Validator.Validate(cfg);

        Assert.NotNull(err);
        Assert.Contains("Admin-Token", err);
    }

    [Fact]
    public void ValidateGetConsumerPositionRequiresReaderTokenNotAdminToken()
    {
        var cfg = BaseConfig("get-consumer-position") with { Addr = "feed:9090", AdminToken = "admin-token", ConsumerId = "consumer-1" };

        var err = Validator.Validate(cfg);

        Assert.NotNull(err);
        Assert.Contains("CDC_API_TOKEN_READER", err);
    }

    [Fact]
    public void ValidateEnableTableRequiresAdminTokenNotReaderToken()
    {
        var cfg = BaseConfig("enable-table") with
        {
            Addr = "feed:9090", Token = "reader-token",
            Source = "quelle-1", Schema = "public", Table = "orders", Publication = "pub_quelle_1",
        };

        var err = Validator.Validate(cfg);

        Assert.NotNull(err);
        Assert.Contains("Admin-Token", err);
    }

    [Fact]
    public void ValidateDisableTableRequiresAdminTokenNotReaderToken()
    {
        var cfg = BaseConfig("disable-table") with
        {
            Addr = "feed:9090", Token = "reader-token",
            Source = "quelle-1", Schema = "public", Table = "orders", Publication = "pub_quelle_1",
        };

        var err = Validator.Validate(cfg);

        Assert.NotNull(err);
        Assert.Contains("Admin-Token", err);
    }

    [Fact]
    public void ValidateGetTableStatusRequiresReaderTokenNotAdminToken()
    {
        var cfg = BaseConfig("get-table-status") with
        {
            Addr = "feed:9090", AdminToken = "admin-token",
            Source = "quelle-1", Schema = "public", Table = "orders", Publication = "pub_quelle_1",
        };

        var err = Validator.Validate(cfg);

        Assert.NotNull(err);
        Assert.Contains("CDC_API_TOKEN_READER", err);
    }

    [Fact]
    public void ValidateListTablesRequiresReaderTokenNotAdminToken()
    {
        var cfg = BaseConfig("list-tables") with { Addr = "feed:9090", AdminToken = "admin-token", Source = "quelle-1", Publication = "pub_quelle_1" };

        var err = Validator.Validate(cfg);

        Assert.NotNull(err);
        Assert.Contains("CDC_API_TOKEN_READER", err);
    }

    [Fact]
    public void ValidateRunRetentionRequiresAdminTokenNotReaderToken()
    {
        var cfg = BaseConfig("run-retention") with { Addr = "feed:9090", Token = "reader-token", Source = "quelle-1" };

        var err = Validator.Validate(cfg);

        Assert.NotNull(err);
        Assert.Contains("Admin-Token", err);
    }

    [Fact]
    public void ValidateDiagnoseRequiresReaderTokenNotAdminToken()
    {
        var cfg = BaseConfig("diagnose") with { Addr = "feed:9090", AdminToken = "admin-token", Source = "quelle-1" };

        var err = Validator.Validate(cfg);

        Assert.NotNull(err);
        Assert.Contains("CDC_API_TOKEN_READER", err);
    }

    [Fact]
    public void ValidateReadChangesRequiresReaderTokenNotAdminToken()
    {
        var cfg = BaseConfig("read-changes") with { Addr = "feed:9090", AdminToken = "admin-token", Source = "quelle-1" };

        var err = Validator.Validate(cfg);

        Assert.NotNull(err);
        Assert.Contains("CDC_API_TOKEN_READER", err);
    }

    [Fact]
    public void ValidateDiagnoseRequiresSource()
    {
        var cfg = BaseConfig("diagnose") with { Addr = "feed:9090", Token = "reader-token" };

        var err = Validator.Validate(cfg);

        Assert.NotNull(err);
        Assert.Contains("-source ist Pflicht", err);
    }
}
