package environmentmanager

import (
	"encoding/json"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/martient/bifrost-env-manager/pkg/utils"
)

func GenerateEnvFile(configJSON []byte, newEnvFilePath string, readOnlyEnvFilesPath string) int {
	rand.Seed(time.Now().UnixNano())

	var config Config
	var jsonRead map[string]interface{}
	var fileName string = ""

	err := json.Unmarshal([]byte(configJSON), &config)
	json.Unmarshal([]byte(configJSON), &jsonRead)
	if err != nil {
		utils.LogError("Error parsing config JSON: %s", err, "Environment manager")
		return 1
	}

	if config.Filename != "" {
		fileName = config.Filename
	} else {
		fileName = ".env"
	}
	var outputFilePath string = newEnvFilePath + fileName

	generateEnvVariables(&config)

	generateExistingVariables(&config, outputFilePath)

	generateReadOnlyVariables(&config, readOnlyEnvFilesPath)

	err = generateStaticVariables(&config, jsonRead)
	if err != nil {
		utils.LogError("Error parsing config JSON: %s", err, "Environment manager")
		return 1
	}
	err = generateRandomValueVariables(&config, jsonRead)
	if err != nil {
		utils.LogError("Error parsing config JSON: %s", err, "Environment manager")
		return 1
	}
	err = generateCustomValueVariables(&config, jsonRead)
	if err != nil {
		utils.LogError("Error parsing config JSON: %s", err, "Environment manager")
		return 1
	}

	err = writeVariablesToFile(&config, outputFilePath)
	if err != nil {
		utils.LogError("Error writing to .env file:", err, "Environment manager")
		return 1
	}

	utils.LogInfo("%s file generated successfully!\n", outputFilePath, "Environment manager")
	return 0
}

func generateEnvVariables(config *Config) {
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		config.EnvVariables = append(config.EnvVariables, Variable{Key: pair[0], Value: pair[1]})
	}
}

func readEnvFile(filePath string) ([]string, error) {
	_, err := os.Stat(filePath)

	if err == nil {
		envContent, err := os.ReadFile(filePath)
		if err != nil {
			return nil, err
		}

		lines := strings.Split(string(envContent), "\n")
		return lines, nil
	} else {
		return nil, err
	}
}

func generateExistingVariables(config *Config, outputFilePath string) {
	lines, err := readEnvFile(outputFilePath)

	if err != nil {
		utils.LogDebug("Error parsing existing env files: %s", err, "Environment manager")
	}

	for _, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			config.ExistingVariables = append(config.ExistingVariables, Variable{Key: parts[0], Value: parts[1]})
		}
	}
}

func generateReadOnlyVariables(config *Config, readOnlyEnvFilesPath string) {
	if readOnlyEnvFilesPath != "" {
		files := strings.Split(readOnlyEnvFilesPath, ";")

		for _, file := range files {
			lines, err := readEnvFile(file)

			if err != nil {
				utils.LogError("Error parsing existing env files: %s", err, "Environment manager")
				continue
			}

			for _, line := range lines {
				parts := strings.SplitN(line, "=", 2)
				if len(parts) == 2 {
					config.ReadOnlyVariables = append(config.ReadOnlyVariables, Variable{Key: parts[0], Value: parts[1]})
				}
			}
		}
	}
}

func generateStaticVariables(config *Config, jsonRead map[string]interface{}) error {
	_, isDefine := jsonRead["static_variables"].([]interface{})
	if !isDefine {
		return nil
	}
	for _, variable := range jsonRead["static_variables"].([]interface{}) {
		if v, ok := variable.(map[string]interface{}); ok {
			for key, value := range v {
				config.StaticVariables = append(config.StaticVariables, Variable{Key: key, Value: value.(string)})
			}
		}
	}
	return nil
}

func generateRandomValueVariables(config *Config, jsonRead map[string]interface{}) error {
	_, isDefine := jsonRead["random_value_variables"].([]interface{})
	if !isDefine {
		return nil
	}
	for _, variable := range jsonRead["random_value_variables"].([]interface{}) {
		if v, ok := variable.(map[string]interface{}); ok {
			var key string = v["key"].(string)
			var value string = generateRandomValue(v)
			var newRandomValueVariable Variable = Variable{Key: key, Value: value}
			config.RandomValueVariables = append(config.RandomValueVariables, newRandomValueVariable)
		}
	}
	return nil
}

func generateRandomValue(settings map[string]interface{}) string {
	length := utils.IntOrDefault(settings["length"], 16)
	availableCharacters := utils.StringOrDefault(settings["available_character"], "")
	asSpecialCharacter := utils.BoolOrDefault(settings["as_special_character"], true)
	asUpperCase := utils.BoolOrDefault(settings["as_upper_case"], true)
	asLowerCase := utils.BoolOrDefault(settings["as_lower_case"], true)
	asDigit := utils.BoolOrDefault(settings["as_digit"], true)

	availableChars := ""
	if availableCharacters != "" {
		availableChars = availableCharacters
	} else {
		if asUpperCase {
			availableChars += "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		}
		if asLowerCase {
			availableChars += "abcdefghijklmnopqrstuvwxyz"
		}
		if asDigit {
			availableChars += "0123456789"
		}
		if asSpecialCharacter {
			availableChars += "!#$%^&*()-_=+[]{}|;:,.<>/?"
		}
	}

	result := make([]byte, length)
	for i := range result {
		result[i] = availableChars[rand.Intn(len(availableChars))]
	}
	return string(result)
}

func generateCustomValueVariables(config *Config, jsonRead map[string]interface{}) error {
	_, isDefine := jsonRead["custom_value_variables"].([]interface{})
	if !isDefine {
		return nil
	}
	for _, variable := range jsonRead["custom_value_variables"].([]interface{}) {
		if v, ok := variable.(map[string]interface{}); ok {
			var key string = v["key"].(string)
			var line string = v["line"].(string)
			var CustomValues []Variable

			_, isDefine := v["values"].([]interface{})
			if isDefine {
				for _, variable := range v["values"].([]interface{}) {
					if v, ok := variable.(map[string]interface{}); ok {
						for key, value := range v {
							var subValueOfCustomValueVariable Variable = Variable{Key: key, Value: value.(string)}
							CustomValues = append(CustomValues, subValueOfCustomValueVariable)
						}
					}
				}
			}

			var value string = replacePlaceholders(line, CustomValues)
			value = replacePlaceholders(value, config.StaticVariables)
			value = replacePlaceholders(value, config.ExistingVariables)
			value = replacePlaceholders(value, config.RandomValueVariables)
			value = replacePlaceholders(value, config.ReadOnlyVariables)
			value = replacePlaceholders(value, config.EnvVariables)
			config.CustomValueVariables = append(config.CustomValueVariables, Variable{Key: key, Value: value})
		}
	}
	return nil
}

func replacePlaceholders(line string, values []Variable) string {
	for _, variable := range values {
		line = strings.ReplaceAll(line, "{{ "+variable.Key+" }}", variable.Value)
	}
	return line
}

func writeVariablesToFile(config *Config, outputFilePath string) error {
	var lines []string

	for _, variable := range config.StaticVariables {
		lines = append(lines, variable.Key+"="+variable.Value)
	}

	for _, variable := range config.RandomValueVariables {
		exisingVariable := searchExistingVariable(variable.Key, config)
		if exisingVariable == nil {
			utils.LogInfo("Existing avariable null %s", exisingVariable, "env mana")
			lines = append(lines, variable.Key+"="+variable.Value)
		} else {
			lines = append(lines, exisingVariable.Key+"="+exisingVariable.Value)
		}
	}

	for _, variable := range config.CustomValueVariables {
		lines = append(lines, variable.Key+"="+variable.Value)
	}

	content := strings.Join(lines, "\n")

	err := os.WriteFile(outputFilePath, []byte(content), 0644)
	if err != nil {
		return err
	}

	return nil
}

func searchExistingVariable(key string, config *Config) *Variable {
	for _, searchVar := range config.ExistingVariables {
		if key == searchVar.Key {
			return &searchVar
		}
	}
	return nil
}
