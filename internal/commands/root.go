package commands

import (
	"fmt"
	"os"

	"github.com/dkrichards86/gotodo/internal/gotodo"
	"github.com/mitchellh/go-homedir"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var todoList string

var rootCmd = &cobra.Command{
	Use:          "gotodo",
	Short:        "A CLI client to manage your todos",
	SilenceUsage: true,
}

// Execute runs the specified command.
func Execute() {
	rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gotodo.yaml)")
	rootCmd.PersistentFlags().StringVar(&todoList, "bucket", "", "todo bucket to use")

	viper.BindPFlag("bucket", rootCmd.PersistentFlags().Lookup("bucket"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		viper.AddConfigPath(home)
		viper.SetConfigName(".gotodo")
	}

	viper.SetDefault("bucket", "Todos")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func getManager() *gotodo.TodoManager {
	bucket := viper.Get("bucket").(string)
	storage := &gotodo.BoltStorage{Bucket: []byte(bucket)}
	return &gotodo.TodoManager{Storage: storage}
}

func drawTable(header []string, data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	table.SetAutoWrapText(false)
	table.SetBorder(false)
	table.SetAutoFormatHeaders(true)
	table.AppendBulk(data)
	table.Render()
}
