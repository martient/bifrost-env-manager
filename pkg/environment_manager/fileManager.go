package environmentmanager

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func GenerateEnvFile(configJSON []byte, newEnvFilePath string) {
	rand.Seed(time.Now().UnixNano())

	var config Config
	var jsonRead map[string]interface{}
	err := json.Unmarshal([]byte(configJSON), &config)
	json.Unmarshal([]byte(configJSON), &jsonRead)
	if err != nil {
		fmt.Println("Error parsing config JSON:", err)
		return
	}
	generateEnvVariables(&config)

	err = generateStaticVariables(&config, jsonRead)
	if err != nil {
		fmt.Println("Error parsing config JSON:", err)
		return
	}
	err = generateRandomValueVariables(&config, jsonRead)
	if err != nil {
		fmt.Println("Error parsing config JSON:", err)
		return
	}
	err = generateCustomValueVariables(&config, jsonRead)
	if err != nil {
		fmt.Println("Error parsing config JSON:", err)
		return
	}

	var fileName string = ""
	if config.Filename != "" {
		fileName = config.Filename
	} else {
		fileName = ".env"
	}
	var outputFilePath string = newEnvFilePath + fileName

	err = writeVariablesToFile(&config, outputFilePath)
	if err != nil {
		fmt.Println("Error writing to .env file:", err)
		return
	}

	fmt.Println(".env file generated successfully!")
}

func generateEnvVariables(config *Config) {
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		config.EnvVariables = append(config.EnvVariables, Variable{Key: pair[0], Value: pair[1]})
	}
}

func generateStaticVariables(config *Config, jsonRead map[string]interface{}) error {
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
	length := intOrDefault(settings["length"], 16)
	availableCharacters := stringOrDefault(settings["available_character"], "")
	asSpecialCharacter := boolOrDefault(settings["as_special_character"], true)
	asUpperCase := boolOrDefault(settings["as_upper_case"], true)
	asLowerCase := boolOrDefault(settings["as_lower_case"], true)
	asDigit := boolOrDefault(settings["as_digit"], true)

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
	for _, variable := range jsonRead["custom_value_variables"].([]interface{}) {
		if v, ok := variable.(map[string]interface{}); ok {
			var key string = v["key"].(string)
			var line string = v["line"].(string)
			var CustomValues []Variable

			for _, variable := range v["values"].([]interface{}) {
				if v, ok := variable.(map[string]interface{}); ok {
					for key, value := range v {
						var subValueOfCustomValueVariable Variable = Variable{Key: key, Value: value.(string)}
						CustomValues = append(CustomValues, subValueOfCustomValueVariable)
					}
				}
			}

			var value string = replacePlaceholders(line, CustomValues)
			value = replacePlaceholders(value, config.StaticVariables)
			value = replacePlaceholders(value, config.RandomValueVariables)
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
	var existingVars []EnvFileVariable

	if _, err := os.Stat(outputFilePath); err == nil {
		envContent, err := os.ReadFile(outputFilePath)
		if err != nil {
			return err
		}

		lines := strings.Split(string(envContent), "\n")
		for _, line := range lines {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				existingVars = append(existingVars, EnvFileVariable{Key: parts[0], Value: parts[1], Exist: true})
			}
		}
	}

	var lines []string

	for _, variable := range config.StaticVariables {
		lines = append(lines, variable.Key+"="+variable.Value)
	}

	for _, variable := range config.RandomValueVariables {
		exisingVariable := searchExistingVariable(variable.Key, existingVars)
		if !exisingVariable.Exist {
			lines = append(lines, variable.Key+"="+variable.Value)
		} else {
			lines = append(lines, exisingVariable.Key+"="+exisingVariable.Value)

		}
	}

	for _, variable := range config.CustomValueVariables {
		exisingVariable := searchExistingVariable(variable.Key, existingVars)

		if !exisingVariable.Exist {
			lines = append(lines, variable.Key+"="+variable.Value)
		} else {
			lines = append(lines, exisingVariable.Key+"="+exisingVariable.Value)
		}
	}

	content := strings.Join(lines, "\n")

	err := os.WriteFile(outputFilePath, []byte(content), 0644)
	if err != nil {
		return err
	}

	return nil
}

func searchExistingVariable(key string, variables []EnvFileVariable) EnvFileVariable {
	for _, searchVar := range variables {
		if key == searchVar.Key {
			return searchVar
		}
	}
	return EnvFileVariable{Key: key, Exist: false}
}
