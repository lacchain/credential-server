package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	calog "github.com/lacchain/credential-server/common/log"
	"github.com/lacchain/credential-server/lib"
	"github.com/lacchain/credential-server/util"
	"github.com/cloudflare/cfssl/log"
	"github.com/pkg/errors"
)

const (
	longName     = "Digital Identity Verifiable Credential Provider Server"
	shortName    = "credential provider server"
	cmdName      = "credential-provider-server"
	envVarPrefix = "CREDENTIAL_PROVIDER_SERVER"
	homeEnvVar   = "CREDENTIAL_PROVIDER_SERVER_HOME"
	caNameReqMsg = "ca.name property is required but is missing from the configuration file"
)

//Version ...
const Version = "0.0.9"

const (
	defaultCfgTemplate = `#############################################################################
#   This is a configuration file for the provider-credential-server command.
#
#   COMMAND LINE ARGUMENTS AND ENVIRONMENT VARIABLES
#   ------------------------------------------------
#   Each configuration element can be overridden via command line
#   arguments or environment variables.  The precedence for determining
#   the value of each element is as follows:
#   1) command line argument
#      Examples:
#      a) --port 443
#         To set the listening port
#   2) environment variable
#      Examples:
#      a) CREDENTIAL_PROVIDER_SERVER_PORT=443
#         To set the listening port
#      b) CREDENTIAL_PROVIDER_SERVER_KEYFILE="../mykey.pem"
#         To set the "keyfile" element in the "ca" section below;
#         note the '_' separator character.
#   3) configuration file
#   4) default value (if there is one)
#      All default values are shown beside each element below.
#
#   FILE NAME ELEMENTS
#   ------------------
#   The value of all fields whose name ends with "file" or "files" are
#   name or names of other files.
#   For example, see "tls.certfile" and "tls.clientauth.certfiles".
#   The value of each of these fields can be a simple filename, a
#   relative path, or an absolute path.  If the value is not an
#   absolute path, it is interpretted as being relative to the location
#   of this configuration file.
#
#############################################################################

# Version of config file
version: <<<VERSION>>>

# Server's listening port (default: 8000)
port: 8000

# Enables debug logging (default: false)
debug: false

# Node Ethereum Blockchain
node: http://ropsten.com

# Issuer 
issuer: did:ita:12345

#############################################################################
#  Repository section
#  Supported types are: "smartContract".
#############################################################################
repository:
  type: smartContract
  address: REPO_ADDRESS

#############################################################################
#  Proof section
#  Supported types are: "SmartContract, P-256, Secp256k1".
#############################################################################
proof:
  type: SmartContract
  verification: 0x0f6487d9640f4230e09d0c2c0ef8b2bef6592573

#############################################################################
# BCCSP (BlockChain Crypto Service Provider) section is used to select which
# crypto library implementation to use
#############################################################################
bccsp:
    hash: SHA2
    security: 256
    # The directory used for the software file-based keystore
    keyStore: ./keystore/UTC--2019
`	
)

var (
	extraArgsError = "Unrecognized arguments found: %v\n\n%s"
)

// Initialize config
func (s *ServerCmd) configInit() (err error) {
	if !s.configRequired() {
		return nil
	}

	s.cfgFileName, s.homeDirectory, err = util.ValidateAndReturnAbsConf(s.cfgFileName, s.homeDirectory, cmdName)
	if err != nil {
		return err
	}

	s.myViper.AutomaticEnv() // read in environment variables that match
	logLevel := s.myViper.GetString("loglevel")
	debug := s.myViper.GetBool("debug")
	calog.SetLogLevel(logLevel, debug)

	log.Debugf("Home directory: %s", s.homeDirectory)

	// If the config file doesn't exist, create a default one
	if !util.FileExists(s.cfgFileName) {
		err = s.createDefaultConfigFile()
		if err != nil {
			return errors.WithMessage(err, "Failed to create default configuration file")
		}
		log.Infof("Created default configuration file at %s", s.cfgFileName)
	} else {
		log.Infof("Configuration file location: %s", s.cfgFileName)
	}

	// Read the config
	err = lib.UnmarshalConfig(s.cfg, s.myViper, s.cfgFileName, true)

	if err != nil {
		return err
	}

	return nil
}

func (s *ServerCmd) createDefaultConfigFile() error {
	var myhost string
	var err error
	myhost, err = os.Hostname()
	if err != nil {
		return err
	}

	// Do string subtitution to get the default config
	cfg := strings.Replace(defaultCfgTemplate, "<<<VERSION>>>", Version, 1)
	cfg = strings.Replace(cfg, "<<<MYHOST>>>", myhost, 1)

	// Now write the file
	cfgDir := filepath.Dir(s.cfgFileName)
	err = os.MkdirAll(cfgDir, 0755)
	if err != nil {
		return err
	}

	// Now write the file
	return ioutil.WriteFile(s.cfgFileName, []byte(cfg), 0644)
}
