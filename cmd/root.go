/*
Copyright © 2020 Graham Burgess <stormmore@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"

	"github.com/epip-io/go-birthday-calendar/pkg/conf"
	"github.com/epip-io/go-birthday-calendar/pkg/models"
)

var cfg models.Config

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-birthday-calendar",
	Short: "This is a rest API service to remind people that their birthday is coming up or wish them 'Happy birthday' on their birthday",
	Long: `This is a rest API service to remind people that their birthday is coming up or wish them 'Happy birthday' on their birthday.

Environment Variables:

It is possible to configure the service by using environment variables. To do so, export a variable that matches the desired flag option that meets the following rules:

- convert to upper case
- replace dashes (-) with underscores (_)
- prefix with BIRTHDAY_

Examples

- port -> BIRTHDAY_PORT
- tls-port -> BIRTHDAY_TLS_PORT

Configuration Reload:

If a configuration file is used, the service will monitor for changes to the file and will update the running configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		var wait time.Duration
		var s *http.Server

		logger.WithFields(log.Fields{
			"db-name":   cfg.DB.Name,
			"db-engine": cfg.DB.Engine,
			"db-user":   cfg.DB.User,
			"db-conn":   cfg.DB.Conn,
		}).Debug("configuring database connection")
		db, err := conf.ConfigureDatabase(&cfg.DB, logger)
		if err != nil {
			logger.Fatalf("unable to configure database connection: %+v", err)
		}
		defer db.Close()

		logger.Info("configuring router")
		r := conf.ConfigureRouter(&cfg, db, logger)

		if cfg.TLS.Enabled {
			s = &http.Server{
				Addr:         fmt.Sprintf(":%d", cfg.TLS.Port),
				WriteTimeout: time.Second * 15,
				ReadTimeout:  time.Second * 15,
				IdleTimeout:  time.Second * 60,
				Handler:      r,
			}

		} else {
			s = &http.Server{
				Addr:         fmt.Sprintf(":%d", cfg.Port),
				WriteTimeout: time.Second * 15,
				ReadTimeout:  time.Second * 15,
				IdleTimeout:  time.Second * 60,
				Handler:      r,
			}

			go func() {
				logger.WithFields(log.Fields{
					"port": cfg.Port,
					"path": cfg.Path,
				}).Info("starting server")
				if err := s.ListenAndServe(); err != nil {
					logger.Fatalf("server stopped: %s: %+v", fmt.Sprintf(":%d", cfg.Port), err)
				}
			}()
		}

		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c

		ctx, cancal := context.WithTimeout(context.Background(), wait)
		defer cancal()

		s.Shutdown(ctx)

		logger.Info("service shut down")
	},
}

// rootFlags represents the base command flags
var rootFlags = []models.Flag{
	{
		Type:    "int",
		Name:    "port",
		Default: 8080,
		Desc:    "port to serve HTTP requests",
	},
	{
		Type:    "string",
		Name:    "path",
		Default: "/",
		Desc:    "path to serve from, e.g. /hello",
	},
	{
		Type:    "bool",
		Name:    "redirect",
		Default: false, // Uncomment the following line if your bare application
		// has an action associated with it:
		Desc: "redirect HTTP to HTTPS",
	},
	{
		Type:    "bool",
		Name:    "tls-enable",
		Default: false,
		Desc:    "enable TLS",
	},
	{
		Type:    "int",
		Name:    "tls-port",
		Default: 8443,
		Desc:    "port to serve TLS/HTTPS requests",
	},
	{
		Type:    "string",
		Name:    "tls-cert",
		Default: "server.crt",
		Desc:    "path to server certificate file, include intermediary certificates",
	},
	{
		Type:    "string",
		Name:    "db-engine",
		Default: "sqlite3",
		Desc:    "database engine to use, supported engines: mysql, mssql, postgres, sqlite3",
	},
	{
		Type:    "string",
		Name:    "db-user",
		Default: "",
		Desc:    "database username, if required",
	},
	{
		Type:    "string",
		Name:    "db-pass",
		Default: "",
		Desc:    "database password, if required",
	},
	{
		Type:    "string",
		Name:    "db-host",
		Default: "localhost",
		Desc:    "database name",
	},
	{
		Type:    "int",
		Name:    "db-port",
		Default: 0,
		Desc:    "database port, default depends on db-engine",
	},
	{
		Type:    "string",
		Name:    "db-name",
		Default: "file:inmemdb1?mode=memory&cache=shared",
		Desc:    "database name",
	},
	{
		Type:    "string",
		Name:    "db-conn",
		Default: "",
		Desc:    "custom database connection string",
	},
	{
		Type:    "string",
		Name:    "log-level",
		Default: "INFO",
		Desc:    "log level",
	},
	{
		Type:    "string",
		Name:    "log-file",
		Default: "-",
		Desc:    "log file to us, defaults to stdin",
	},
}

var logger *log.Entry

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.go-birthday-calendar.yaml)")

	// Set up all the CLI flags for the base command
	for _, f := range rootFlags {
		switch f.Type {
		case "string":
			if f.Short == "" {
				rootCmd.Flags().String(f.Name, f.Default.(string), f.Desc)
			} else {
				rootCmd.Flags().StringP(f.Name, f.Short, f.Default.(string), f.Desc)
			}
		case "int":
			if f.Short == "" {
				rootCmd.Flags().Int(f.Name, f.Default.(int), f.Desc)
			} else {
				rootCmd.Flags().IntP(f.Name, f.Short, f.Default.(int), f.Desc)
			}
		case "bool":
			if f.Short == "" {
				rootCmd.Flags().Bool(f.Name, f.Default.(bool), f.Desc)
			} else {
				rootCmd.Flags().BoolP(f.Name, f.Short, f.Default.(bool), f.Desc)
			}
		}
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Bind base command flags to configuration
	for _, f := range rootFlags {
		viper.BindPFlag(strings.ReplaceAll(f.Name, "-", "."), rootCmd.Flags().Lookup(f.Name))
	}

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".go-birthday-calendar" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".go-birthday-calendar")
	}

	viper.SetEnvPrefix("BIRTHDAY")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		viper.WatchConfig()
		viper.OnConfigChange(func(e fsnotify.Event) {
			if err := viper.ReadInConfig(); err == nil {
				viper.Unmarshal(&cfg)
			}
		})
	}

	viper.Unmarshal(&cfg)

	var err error
	logger, err = conf.ConfigureLogger(&cfg.Log)
	if err != nil {
		fmt.Println("failed to configure logger:", err)
		os.Exit(1)
	}

	if viper.ConfigFileUsed() != "" {
		logger.WithFields(log.Fields{
			"config": viper.ConfigFileUsed(),
		}).Debug("using configuration file")
	}
	logger.WithFields(log.Fields{
		"log-level": cfg.Log.Level,
		"log-file":  cfg.Log.File,
		"db-name":   cfg.DB.Name,
		"db-engine": cfg.DB.Engine,
		"db-user":   cfg.DB.User,
		"db-conn":   cfg.DB.Conn,
		"port":      cfg.Port,
		"redirect":  cfg.Redirect,
		"tls_port":  cfg.TLS.Port,
		"tls_cert":  cfg.TLS.Cert,
		"tls_key":   cfg.TLS.Key,
	}).Debug("starting config")
}
