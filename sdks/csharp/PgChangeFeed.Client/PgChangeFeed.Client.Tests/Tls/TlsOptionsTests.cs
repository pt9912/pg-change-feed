using Xunit;

namespace PgChangeFeed.Client.Tests.Tls;

/// <summary>
/// The trust anchor option is checked when the options are created: a path
/// that cannot be read and a file without a PEM certificate are a
/// configuration error of the caller, not a failure of the first request.
/// </summary>
public sealed class TlsOptionsTests : IDisposable
{
    private static readonly Uri Address = new("https://feed.example.invalid:8443");

    private readonly TestCertificates _certificates = new();

    public void Dispose() => _certificates.Dispose();

    [Fact]
    public void WithoutAnchor_TheAnchorIsNull()
    {
        var options = new PgChangeFeedClientOptions(Address, "token");

        Assert.Null(options.TrustAnchorFile);
    }

    [Fact]
    public void WithAnchor_TheOptionCarriesThePath()
    {
        var path = _certificates.WritePem(_certificates.Create(["feed.example.invalid"]));

        var options = new PgChangeFeedClientOptions(Address, "token", path);

        Assert.Equal(path, options.TrustAnchorFile);
    }

    [Fact]
    public void UnreadablePath_ThrowsWhenTheOptionsAreCreated()
    {
        var error = Assert.Throws<ArgumentException>(
            () => new PgChangeFeedClientOptions(Address, "token", _certificates.MissingPath()));

        Assert.Contains("missing.pem", error.Message);
    }

    [Fact]
    public void FileWithoutPemCertificate_ThrowsWhenTheOptionsAreCreated()
    {
        var path = _certificates.WriteText("this is not a certificate\n");

        var error = Assert.Throws<ArgumentException>(() => new PgChangeFeedClientOptions(Address, "token", path));

        Assert.Contains("no PEM certificate", error.Message);
    }

    [Fact]
    public void EmptyFile_ThrowsWhenTheOptionsAreCreated()
    {
        var path = _certificates.WriteText(string.Empty);

        Assert.Throws<ArgumentException>(() => new PgChangeFeedClientOptions(Address, "token", path));
    }
}
