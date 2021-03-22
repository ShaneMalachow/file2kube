/*
Copyright © 2021 Shane Malachow <shane.malachow@sciencelogic.com>

*/
package cmd

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
)

var secretTemplate = `---
apiVersion: v1
kind: Secret
metadata:
  name: {{ .Name }}
  namespace: {{ .Namespace }}
type: Opaque
data:
  {{ range $key, $val := .Files }}{{$key}}: {{$val}}{{end}}
`

// secretCmd represents the secret command
var secretCmd = &cobra.Command{
	Use:   "secret",
	Short: "Converts files into Kubernetes secrets with base64 encoding",
	Long:  `This command converts files into a single Kubernetes Secret by base64 encoding all the data and including it with the key being the file name. Useful for things like certificates and keys.`,
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
			enc := base64.StdEncoding.EncodeToString([]byte(file))
			// fmt.Printf("Encoded: %s\n", enc)
			data.Files[arg] = enc
		}

		t := template.Must(template.New("tmpl").Parse(secretTemplate))
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
	rootCmd.AddCommand(secretCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// secretCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// secretCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
