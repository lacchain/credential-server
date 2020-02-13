package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/lacchain/credential-server/lib"
	"github.com/lacchain/credential-server/restapi"
	"github.com/lacchain/credential-server/util"
	"github.com/cloudflare/cfssl/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	version = "version"
)

// ServerCmd encapsulates cobra command that provides command line interface
// for the Credential Provider server and the configuration used by the Credential Provider server
type ServerCmd struct {
	// name of the credential-provider-server command (init, start, version)
	name string
	// rootCmd is the cobra command
	rootCmd *cobra.Command
	// My viper instance
	myViper *viper.Viper
	// blockingStart indicates whether to block after starting the server or not
	blockingStart bool
	// cfgFileName is the name of the configuration file
	cfgFileName string
	// homeDirectory is the location of the server's home directory
	homeDirectory string
	// serverCfg is the server's configuration
	cfg *restapi.ServerConfig
}

// NewCommand returns new ServerCmd ready for running
func NewCommand(name string, blockingStart bool) *ServerCmd {
	s := &ServerCmd{
		name:          name,
		blockingStart: blockingStart,
		myViper:       viper.New(),
	}
	s.init()
	return s
}

// Execute runs this ServerCmd
func (s *ServerCmd) Execute() error {
	return s.rootCmd.Execute()
}

// init initializes the ServerCmd instance
// It intializes the cobra root and sub commands and
// registers command flgs with viper
func (s *ServerCmd) init() {
	// root command
	rootCmd := &cobra.Command{
		Use:   cmdName,
		Short: longName,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			err := s.configInit()
			if err != nil {
				return err
			}
			cmd.SilenceUsage = true
			util.CmdRunBegin(s.myViper)
			return nil
		},
	}
	s.rootCmd = rootCmd

	// initCmd represents the server init command
	initCmd := &cobra.Command{
		Use:   "init",
		Short: fmt.Sprintf("Initialize the %s", shortName),
		Long:  "Generate the key material needed by the server if it doesn't already exist",
	}
	initCmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			return errors.Errorf(extraArgsError, args, initCmd.UsageString())
		}
		err := s.getServer().Init(false)
		if err != nil {
			util.Fatal("Initialization failure: %s", err)
		}
		log.Info("Initialization was successful")
		return nil
	}
	s.rootCmd.AddCommand(initCmd)

	// startCmd represents the server start command
	startCmd := &cobra.Command{
		Use:   "start",
		Short: fmt.Sprintf("Start the %s", shortName),
	}

	startCmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			return errors.Errorf(extraArgsError, args, startCmd.UsageString())
		}
		err := s.getServer().Start()
		if err != nil {
			return err
		}
		return nil
	}
	s.rootCmd.AddCommand(startCmd)

	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Prints Credential Provider Server version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Print("0.0.9")
		},
	}
	s.rootCmd.AddCommand(versionCmd)
	s.registerFlags()
}

// registerFlags registers command flags with viper
func (s *ServerCmd) registerFlags() {
	// Get the default config file path
	cfg := util.GetDefaultConfigFile(cmdName)

	// All env variables must be prefixed
	s.myViper.SetEnvPrefix(envVarPrefix)
	s.myViper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Set specific global flags used by all commands
	pflags := s.rootCmd.PersistentFlags()
	pflags.StringVarP(&s.cfgFileName, "config", "c", "", "Configuration file")
	pflags.MarkHidden("config")
	// Don't want to use the default parameter for StringVarP. Need to be able to identify if home directory was explicitly set
	pflags.StringVarP(&s.homeDirectory, "home", "H", "", fmt.Sprintf("Server's home directory (default \"%s\")", filepath.Dir(cfg)))

	// Register flags for all tagged and exported fields in the config
	s.cfg = &restapi.ServerConfig{}

	err := util.RegisterFlags(s.myViper, pflags, s.cfg, nil)
	if err != nil {
		panic(err)
	}
}

// Configuration file is not required for some commands like version
func (s *ServerCmd) configRequired() bool {
	return s.name != version
}

// getServer returns a lib.Server for the init and start commands
func (s *ServerCmd) getServer() *lib.Server {
	return &lib.Server{
		HomeDir:       s.homeDirectory,
		Config:        s.cfg,
		BlockingStart: s.blockingStart,
	}
}
