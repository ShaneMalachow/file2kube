/*
Copyright © 2021 Shane Malachow

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
	"encoding/base64"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string
var filename string
var namespace string
var name string

var templateTxt = `---
apiVersion: v1
kind: Secret
metadata:
  name: {{ .SecretName }}
  namespace: {{ .Namespace }}
type: Opaque
data:
  {{ range $key, $val := .Files }}{{$key}}: {{$val}}{{end}}
`

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "file2kube",
	Short: "file2kube creates Kubernetes secrets from files",
	Long:  `file2kube is a utility that creates Kubernetes secrets from plaintext files, by doing the boilerplate for you and base64 encoding the file data.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: convert,
}

type TemplateData struct {
	Name      string
	Namespace string
	Files     map[string]string
}

func convert(cmd *cobra.Command, args []string) {
	data := TemplateData{
		Name:      name,
		Namespace: namespace,
		Files:     make(map[string]string),
	}
	for _, arg := range args {
		file, err := ioutil.ReadFile(arg)
		if err != nil {
			fmt.Printf("Couldn't read %s: %s\n", arg, err.Error())
		}
		enc := base64.StdEncoding.EncodeToString([]byte(file))
		// fmt.Printf("Encoded: %s\n", enc)
		data.Files[arg] = enc
	}

	t := template.Must(template.New("secret-tmpl").Parse(templateTxt))
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	err = t.Execute(f, data)
	if err != nil {
		fmt.Printf("Error printing template: '%s'\n", err)
	}
}

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

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.file2secret.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().StringVarP(&name, "name", "n", "new-file", "The name for the Secret object")
	rootCmd.PersistentFlags().StringVar(&namespace, "namespace", "default", "The Kubernetes Namespace to create the Secret in")
	rootCmd.PersistentFlags().StringVarP(&filename, "filename", "f", name+".yaml", "The name of the file")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
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

		// Search config in home directory with name ".file2secret" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".file2secret")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
