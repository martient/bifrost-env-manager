package environmentmanager

import (
	"encoding/json"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/martient/golang-utils/utils"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func GenerateEnvFile(configJSON []byte, newEnvFilePath string, readOnlyEnvFilesPath string) int {
	rand.Seed(time.Now().UnixNano())

	var config Config
	var jsonRead map[string]interface{}
	var fileName string = ""

	err := json.Unmarshal([]byte(configJSON), &config)
	json.Unmarshal([]byte(configJSON), &jsonRead)
	if err != nil {
		utils.LogError("Error parsing config JSON: %s", "Environment manager", err)
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
		utils.LogError("Error parsing config JSON: %s", "Environment manager", err)
		return 1
	}
	err = generateRandomValueVariables(&config, jsonRead)
	if err != nil {
		utils.LogError("Error parsing config JSON: %s", "Environment manager", err)
		return 1
	}
	err = generateCustomValueVariables(&config, jsonRead)
	if err != nil {
		utils.LogError("Error parsing config JSON: %s", "Environment manager", err)
		return 1
	}

	err = writeVariablesToFile(&config, outputFilePath)
	if err != nil {
		utils.LogError("Error writing to .env file:", "Environment manager", err)
		return 1
	}

	utils.LogInfo("%s file generated successfully!\n", "Environment manager", outputFilePath)
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
		utils.LogDebug("Error parsing existing env files: %s", "Environment manager", err)
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
				utils.LogError("Error parsing existing env files: %s", "Environment manager", err)
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

			var value string = replacePlaceholders(line, CustomValues, config)
			value = replacePlaceholders(value, config.StaticVariables, config)
			value = replacePlaceholders(value, config.ExistingVariables, config)
			value = replacePlaceholders(value, config.RandomValueVariables, config)
			value = replacePlaceholders(value, config.ReadOnlyVariables, config)
			value = replacePlaceholders(value, config.EnvVariables, config)
			config.CustomValueVariables = append(config.CustomValueVariables, Variable{Key: key, Value: value})
		}
	}
	return nil
}

func formatVariableWithFlag(value string, flag string) string {
	switch flag {
	case "UPPERCASE":
		return strings.ToUpper(value)
	case "LOWERCASE":
		return strings.ToLower(value)
	case "CAPITALIZE":
		return cases.Title(language.Und).String(value)
	case "POSTGRESQL_MODEL":
		return strings.ReplaceAll(value, "-", "_")
	default:
		return value
	}
}

func checkIfVariableExistAnyConfig(target string, config *Config) bool {
	if searchIfVariableExist(target, config.CustomValueVariables) != nil || searchIfVariableExist(target, config.StaticVariables) != nil || searchIfVariableExist(target, config.ExistingVariables) != nil || searchIfVariableExist(target, config.RandomValueVariables) != nil || searchIfVariableExist(target, config.ReadOnlyVariables) != nil || searchIfVariableExist(target, config.EnvVariables) != nil {
		return true
	}
	return false
}

func parseInsideCustomVarString(input string, variableToFind string, config *Config) (string, string) {
	re := regexp.MustCompile(`{{\s*(\w+)\s*\|\|\s*(\w+)\s*`)

	matches := re.FindAllStringSubmatch(input, -1)
	if checkIfVariableExistAnyConfig(matches[0][2], config) {
		return matches[0][2], "{{ " + matches[0][1] + " || " + matches[0][2]
	}
	return matches[0][1], "{{ " + matches[0][1] + " || " + matches[0][2]
}

func replacePlaceholders(line string, values []Variable, config *Config) string {
	for _, variable := range values {
		var flagsToReplace []string
		var keyChanged bool = false

		var target string = variable.Key
		if strings.Contains(line, "{{ "+variable.Key+" ||") {
			replace, targetForReplace := parseInsideCustomVarString(line, variable.Key, config)
			target = replace
			line = strings.ReplaceAll(line, targetForReplace, "{{ "+target)
			keyChanged = true
		}

		for _, flag := range FLAGS {
			if strings.Contains(line, "{{ "+variable.Key+" %"+flag+"% }}") {
				flagsToReplace = append(flagsToReplace, flag)
			}
		}
		for _, flagToUseDuringReplace := range flagsToReplace {
			line = strings.ReplaceAll(line, "{{ "+variable.Key+" %"+flagToUseDuringReplace+"% }}", formatVariableWithFlag(variable.Value, flagToUseDuringReplace))

		}
		if strings.Contains(line, "{{ "+variable.Key+" }}") {
			line = strings.ReplaceAll(line, "{{ "+variable.Key+" }}", variable.Value)
		} else if keyChanged {
			line = replacePlaceholders(line, config.CustomValueVariables, config)
			line = replacePlaceholders(line, config.StaticVariables, config)
			line = replacePlaceholders(line, config.ExistingVariables, config)
			line = replacePlaceholders(line, config.RandomValueVariables, config)
			line = replacePlaceholders(line, config.ReadOnlyVariables, config)
			line = replacePlaceholders(line, config.EnvVariables, config)
			keyChanged = false
		}
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
	return searchIfVariableExist(key, config.ExistingVariables)
}

func searchIfVariableExist(key string, variables []Variable) *Variable {
	for _, searchVar := range variables {
		if key == searchVar.Key {
			return &searchVar
		}
	}
	return nil
}
