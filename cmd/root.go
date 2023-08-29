package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/karuppiah7890/go-jsonschema-generator"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func walk(key string, value string, jsonValues map[string]any) map[string]any {
	for _, subkey := range strings.Split(key, ".") {
		jsonValues = jsonValues["properties"].(map[string]any)[subkey].(map[string]any)
	}
	return jsonValues
}

type extractSource struct {
	option     string
	key        string
	value      string
	jsonValues map[string]any
}

func (s extractSource) extractEnums() error {
	if s.option == "@schemaEnum" {
		child := walk(s.key, s.value, s.jsonValues)
		var enums []string
		for _, v := range strings.Split(s.value, ",") {
			enums = append(enums, v)
		}

		child["enum"] = enums
	}
	return nil
}

func (s extractSource) extractRegex() error {
	if s.option == "@schemaRegex" {
		child := walk(s.key, s.value, s.jsonValues)
		child["pattern"] = s.value
	}
	return nil
}

func (s extractSource) extractMinimum() error {
	if s.option == "@schemaMinimum" {
		child := walk(s.key, s.value, s.jsonValues)

		v, err := strconv.ParseInt(s.value, 10, 32)
		if err != nil {
			return fmt.Errorf("'%s' must return an integer : %v", s.key, err)
		}
		child["minimum"] = v
	}
	return nil
}

func (s extractSource) extractMaximum() error {
	if s.option == "@schemaMaximum" {
		child := walk(s.key, s.value, s.jsonValues)

		v, err := strconv.ParseInt(s.value, 10, 32)
		if err != nil {
			return fmt.Errorf("'%s' must return an integer : %v", s.key, err)
		}
		child["maximum"] = v
	}
	return nil
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:           "helm schema-gen <values-yaml-file>",
	SilenceUsage:  true,
	SilenceErrors: true,
	Short:         "Helm plugin to generate json schema for values yaml",
	Long: `Helm plugin to generate json schema for values yaml

Examples:
  $ helm schema-gen values.yaml    # generate schema json
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("pass one values yaml file")
		}
		if len(args) != 1 {
			return fmt.Errorf("schema can be generated only for one values yaml at once")
		}

		valuesFilePath := args[0]
		values := make(map[string]interface{})
		valuesFileData, err := os.ReadFile(valuesFilePath)

		if err != nil {
			return fmt.Errorf("error when reading file '%s': %v", valuesFilePath, err)
		}
		err = yaml.Unmarshal(valuesFileData, &values)

		// Generate the JSON from YAML
		s := &jsonschema.Document{}
		s.ReadDeep(values)

		// Reload the JSON as map for editing
		jsonValues := make(map[string]any)
		err = json.Unmarshal([]byte(s.String()), &jsonValues)
		if err != nil {
			return fmt.Errorf("error when reading file '%s': %v", valuesFilePath, err)
		}

		for _, s := range strings.Split(string(valuesFileData), "\n") {
			params := strings.Split(s, " ")
			// Extract double hash comment only, with some actual command, will be ignored otherwise
			if strings.HasPrefix(s, "##") && len(params) > 3 {
				vars := extractSource{
					option:     params[1],
					key:        params[2],
					value:      params[3],
					jsonValues: jsonValues,
				}
				err1 := vars.extractEnums()
				err2 := vars.extractRegex()
				err3 := vars.extractMinimum()
				err4 := vars.extractMaximum()

				if err := errors.Join(err1, err2, err3, err4); err != nil {
					return err
				}
			}
		}

		j, err := json.MarshalIndent(jsonValues, "", "  ")
		if err != nil {
			return fmt.Errorf("error when decoding json: %v", err)
		}
		fmt.Println(string(j))
		return nil
	},
}

// Execute executes the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
