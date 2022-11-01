package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/youngshawn/go-project-demo/course/routes"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "course",
	Short: "A Course Management system",
	Long:  `A Course Management system`,
	Run: func(cmd *cobra.Command, args []string) {
		router := gin.Default()
		router.SetTrustedProxies([]string{"127.0.0.1"})

		routes.InstallRoutes(router)

		log.Fatal(router.Run(":3000"))
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.course.yaml)")

	//rootCmd.Flags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.course.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".course.yaml".
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".course")
	}

	//pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	//pflag.Parse()
	//viper.BindPFlags(pflag.CommandLine)

	viper.SetEnvPrefix("COURSE")
	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
