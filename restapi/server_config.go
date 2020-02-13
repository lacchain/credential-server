package restapi

const (
	// DefaultServerPort is the default listening port for the fabric-ca server
	DefaultServerPort = 8000

	// DefaultServerAddr is the default listening address for the fabric-ca server
	DefaultServerAddr = "0.0.0.0"
)

// ServerConfig is the credential provider server's config
// The tags are recognized by the RegisterFlags function in credentialprovider/util/flag.go
// and are as follows:
// "def" - the default value of the field;
// "opt" - the optional one character short name to use on the command line;
// "help" - the help message to display on the command line;
// "skip" - to skip the field.
type ServerConfig struct {
	// Keystore Secret
	Secret string `def:"Password." opt:"x" help:"Keystore Secret"`
	// Listening port for the server
	Port int `def:"8000" opt:"p" help:"Listening port of credential-provider-server"`
	// Certificate to listen on TLS
	TLSCertificate string `def:"server.crt" opt:"s" help:"Certificate to listen on TLS"`
	// Certificate to listen on TLS
	TLSKey string `def:"server.key" opt:"t" help:"Certificate Private Key"`
	// Bind address for the server
	Host string `def:"0.0.0.0" help:"Listening address of credential-provider-server"`
	// Enables debug logging
	Debug bool `def:"false" opt:"d" help:"Enable debug level logging" hide:"true"`
	// Sets the logging level on the server
	LogLevel string `help:"Set logging level (info, warning, debug, error, fatal, critical)"`
	// Node Blockchain to connect
	Node string `help:"Node URL Blockchain to connect through RPC o IPC"`
	// Issuer of Credentials
	Issuer string `help:"DID issuer that sign credentials"`
	// Private Key
	PrivateKey string `help:"KeyStore location to save private key"`
	// Ethereum Address
	Address string `help:"Ethereum address of private Key"`
	// Repository in blockchain
	Repository RepositoryConfig `help:"Set repository address in blockchain"`
	// Proof to verifiable credential
	Proof ProofConfig `help:"Set method verification of verifiable credential"`
	// Blockchain Crypto Service Provider
	Bccsp BCCSP `help:"Set hash algorithm and path keystore"` 
}

// RepositoryConfig ...
type RepositoryConfig struct {
	Type    string `help:"repository type that saves credentials hash"`
	Address string `help:"smart contract address"`
}

// ProofConfig ...
type ProofConfig struct {
	Type         string `help:"proof type that can be smart contract, p-256 or secp256k1"`
	Verification string `help:"contract address or public key to verify digital signature"`
}

// BCCSP ...
type BCCSP struct {
	Hash string `help:"version hash to use"`
	Security string `help:"Algorithm hashing to use, SHA1, SHA2, SHA3"`
	KeyStore string `help:"Path to file used to save keystore"`
}
