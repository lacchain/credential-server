package lib

import (
	"fmt"
	calog "github.com/lacchain/credential-server/common/log"
	"github.com/cloudflare/cfssl/log"
	"net"
	"net/http"
	_ "net/http/pprof" // import to support profiling
	"os"
	"path/filepath"
	"strconv"

	"github.com/pkg/errors"

	"github.com/lacchain/credential-server/restapi"
	"github.com/lacchain/credential-server/util"
	"github.com/lacchain/credential-server/restapi/operations"
	loads "github.com/go-openapi/loads"
	flags "github.com/jessevdk/go-flags"
	bl "github.com/lacchain/credential-server/blockchain"
)

const (
	defaultClientAuth                   = "noclientcert"
	credentialProviderServerProfilePort = "CREDENTIAL_PROVIDER_SERVER_PROFILE_PORT"
	apiPathPrefix                       = "/v1/"
)

// Server is the credential provider server
type Server struct {
	// The home directory for the server
	HomeDir string
	// BlockingStart if true makes the Start function blocking;
	// It is non-blocking by default.
	BlockingStart bool
	// The server's configuration
	Config *restapi.ServerConfig
}

// Init initializes a Credential Provider Server
func (s *Server) Init(renew bool) (err error) {
	err = s.init(renew)
	return err
}

// init initializses the server
func (s *Server) init(renew bool) (err error) {
	serverVersion := "0.0.9"
	err = calog.SetLogLevel(s.Config.LogLevel, s.Config.Debug)
	if err != nil {
		return err
	}
	log.Infof("Server Version: %s", serverVersion)

	// Initialize the config
	err = s.initConfig()
	if err != nil {
		return err
	}

	// Successful initialization
	return nil
}

// Start the credential provider server
func (s *Server) Start() (err error) {
	log.Infof("configggg:",s.Config.Secret)
	log.Infof("keystore:",s.Config.Bccsp.KeyStore)
	s.getKey(s.Config.Bccsp.KeyStore,s.Config.Secret)

	log.Infof("Starting server in home directory: %s", s.HomeDir)

	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		log.Error(err)
	}

	api := operations.NewCredentialProviderAPI(swaggerSpec)
	server := restapi.NewServer(api)
	defer server.Shutdown()

	parser := flags.NewParser(server, flags.Default)
	parser.ShortDescription = "Credential Mother"
	parser.LongDescription = "This is a provider credential server that validates, signs, generates credential to identify persons, institutions and objects"

	server.ConfigureFlags()
	for _, optsGroup := range api.CommandLineOptionsGroups {
		_, err := parser.AddGroup(optsGroup.ShortDescription, optsGroup.LongDescription, optsGroup.Options)
		if err != nil {
			log.Error(err)
		}
	}

	if _, err := parser.Parse(); err != nil {
		code := 1
		if fe, ok := err.(*flags.Error); ok {
			if fe.Type == flags.ErrHelp {
				code = 0
			}
		}
		os.Exit(code)
	}

	server.ConfigureAPI(s.Config)

	if err := server.Serve(); err != nil {
		log.Error(err)
	}

	return nil
}

// initConfig initializes the configuration for the server
func (s *Server) initConfig() (err error) {
	// Home directory is current working directory by default
	if s.HomeDir == "" {
		s.HomeDir, err = os.Getwd()
		if err != nil {
			return errors.Wrap(err, "Failed to get server's home directory")
		}
	}
	// Make home directory absolute, if not already
	absoluteHomeDir, err := filepath.Abs(s.HomeDir)
	if err != nil {
		return fmt.Errorf("Failed to make server's home directory path absolute: %s", err)
	}
	s.HomeDir = absoluteHomeDir
	// Create config if not set
	if s.Config == nil {
		s.Config = new(restapi.ServerConfig)
	}

	err = s.generateKeyStore()
	if err != nil {
		return fmt.Errorf("Failed to create KeyStore: %s", err)
	}

	//get keystore filename
	var files []string

    root := "./keystore"
    err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
        files = append(files, path)
        return nil
    })
    if err != nil {
        panic(err)
    }

	util.ReplaceTextFile("./credential-provider-server-config.yaml","keystore/UTC--2019",files[1])

	err = s.deployRepositoryContract()
	if err != nil{
		return fmt.Errorf("Failed deploy repository smart contract: %s", err)
	}

	return nil
}

func (s *Server) generateKeyStore() (error){
	fmt.Println("secret:",s.Config.Secret)
	_, address, err := util.CreateKeyStore("./keystore",s.Config.Secret)
	if err != nil{
		fmt.Println("Error:",err)
	}

	s.Config.Address = address

	return nil
}

func (s *Server) deployRepositoryContract() (error){
	client := new(bl.Client)
	err := client.Connect(s.Config.Node)
	if err != nil {
		log.Error(err)
	}
	defer client.Close()
    contractAddress, err := client.DeployRepositoryContract(s.Config.PrivateKey) 
    if err != nil{
		log.Error(err)
		return err
	}

	repositoryConfig := new(restapi.RepositoryConfig) 
	repositoryConfig.Type = "SMARTCONTRACT"
	repositoryConfig.Address = contractAddress
	s.Config.Repository = *repositoryConfig

	util.ReplaceTextFile("./credential-provider-server-config.yaml","REPO_ADDRESS",contractAddress)

	return nil
}

func (s *Server) getKey(file, secret string) (error) {
	key,err:=util.GetKey(file, secret)
	if err !=nil{
		fmt.Println("Error getting Key from Keystore:",err)
	}

	s.Config.PrivateKey = util.FromECDSA(key.PrivateKey)
	s.Config.Address = key.Address.Hex()

	return nil
}

// Starting listening and serving
func (s *Server) listenAndServe() (err error) {

	c := s.Config

	// Set default listening address and port
	if c.Host == "" {
		c.Host = restapi.DefaultServerAddr
	}
	if c.Port == 0 {
		c.Port = restapi.DefaultServerPort
	}
	//	addr := net.JoinHostPort(c.Address, strconv.Itoa(c.Port))
	var addrStr string

	log.Infof("Listening on %s", addrStr)

	err = s.checkAndEnableProfiling()
	if err != nil {
		//s.closeListener()
		return errors.WithMessage(err, "TCP listen for profiling failed")
	}

	return nil
}

// checkAndEnableProfiling checks for FABRIC_CA_SERVER_PROFILE_PORT env variable
// if it is set, starts listening for profiling requests at the port specified
// by the environment variable
func (s *Server) checkAndEnableProfiling() error {
	// Start listening for profile requests
	pport := os.Getenv(credentialProviderServerProfilePort)
	if pport != "" {
		iport, err := strconv.Atoi(pport)
		if err != nil || iport < 0 {
			log.Warningf("Profile port specified by the %s environment variable is not a valid port, not enabling profiling",
				credentialProviderServerProfilePort)
		} else {
			addr := net.JoinHostPort(s.Config.Host, pport)
			listener, err1 := net.Listen("tcp", addr)
			log.Infof("Profiling enabled; listening for profile requests on port %s", pport)
			if err1 != nil {
				return err1
			}
			go func() {
				log.Debugf("Profiling enabled; waiting for profile requests on port %s", pport)
				err := http.Serve(listener, nil)
				log.Errorf("Stopped serving for profiling requests on port %s: %s", pport, err)
			}()
		}
	}
	return nil
}
