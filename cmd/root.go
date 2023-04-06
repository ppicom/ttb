/*
Copyright © 2023 Pere Picó Muntaner

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ppicom/ttb/internal/gen"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ttb <text>",
	Short: "Make text into a background",
	Long: `ttb allows you to create an image that contains the text you
	indicated.
	
	Use at your discretion.`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		text := strings.Join(args, " ")
		w, err := cmd.Flags().GetInt("width")
		if err != nil {
			return err
		}

		h, err := cmd.Flags().GetInt("height")
		if err != nil {
			return err
		}

		return generateAction(os.Stdout, text, w, h)
	},
}

func generateAction(out io.Writer, text string, width, height int) error {
	imgFname, err := gen.TextToImage(text, &gen.Config{Width: width, Height: height})
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(out, imgFname)
	return err
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

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ttb.yaml)")

	rootCmd.Flags().IntP("width", "x", 400, "width of the picture")
	rootCmd.Flags().IntP("height", "y", 400, "height of the picture")
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

		// Search config in home directory with name ".ttb" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".ttb")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
