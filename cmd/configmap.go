/*
Copyright © 2021 Shane Malachow <shane.malachow@sciencelogic.com>

*/
package cmd

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var configTemplate = `---
apiVersion: v1
kind: ConfigMap
metadata:
  name: {{ .Name }}
  namespace: {{ .Namespace }}
data:
  {{ range $key, $val := .Files }}{{$key}}: |
    {{$val}}
	{{end}}
`

// configmapCmd represents the configmap command
var configmapCmd = &cobra.Command{
	Use:   "configmap",
	Short: "Converts files into a single Kubernetes ConfigMap",
	Long:  `This command converts files into a single Kubernetes ConfigMap with the key being the file name. Useful for things like config files.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("At least one argument required.")
			return
		}
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
			indentedFile := strings.ReplaceAll(string(file), "\n", "\n    ")
			data.Files[arg] = indentedFile
		}

		t := template.Must(template.New("tmpl").Parse(configTemplate))
		f, err := os.Create(filename)
		if err != nil {
			panic(err)
		}
		err = t.Execute(f, data)
		if err != nil {
			fmt.Printf("Error printing template: '%s'\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(configmapCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configmapCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configmapCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
