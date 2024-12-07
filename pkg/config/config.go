package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strconv"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Configuration struct {
	Unsupported              map[string]string `key:"UNSUPPORTED"`
	Environment              string            `yaml:"environment"`
	FromEnvBool              bool              `key:"BOOL"`
	FromEnvInt               int8              `key:"INT"`
	FromEnvUint              uint16            `key:"UINT"`
	FromEnvFloat             float32           `key:"FLOAT"`
	FromEnv                  string            `key:"ENV_ARG"`
	FromSettingsFile         string            `yaml:"fromSettings"`
	FromCommandLineArguments string            `key:"CL_ARG"`
	Sub                      *SubConfiguration `yaml:"sub" key:"SUB"`
}

type SubConfiguration struct {
	FromEnvBool              bool    `key:"BOOL"`
	FromEnvInt               int8    `key:"INT"`
	FromEnvUint              uint16  `key:"UINT"`
	FromEnvFloat             float32 `key:"FLOAT"`
	FromEnv                  string  `key:"ENV_ARG"`
	FromCommandLineArguments string  `key:"CL_ARG"`
	FromSettingsFile         string  `yaml:"fromSettings"`
}

func Load(configFile string) (*Configuration, error) {
	cfg := &Configuration{}

	loadDefaults(cfg)
	err := loadFromSettingsFile(configFile, cfg)
	if err != nil {
		return nil, err
	}
	err = loadFromEnvironmentVariables(cfg)
	if err != nil {
		return nil, err
	}
	loadFromCommandLineArguments(cfg)

	return cfg, nil
}

func loadDefaults(cfg *Configuration) {
	cfg.Environment = "Development"
}

func loadFromSettingsFile(configFile string, cfg *Configuration) error {
	if configFile == "" {
		return nil
	}

	ext := filepath.Ext(configFile)

	switch {
	case ext == ".yaml" || ext == ".yml":
		yamlFile, err := os.Open(configFile)
		if err != nil {
			return err
		}
		defer yamlFile.Close()

		bytes, err := io.ReadAll(yamlFile)
		if err != nil {
			return err
		}

		err = yaml.Unmarshal(bytes, cfg)
		if err != nil {
			return err
		}
	case ext == ".json":
		jsonFile, err := os.Open(configFile)
		if err != nil {
			return err
		}
		defer jsonFile.Close()

		bytes, err := io.ReadAll(jsonFile)
		if err != nil {
			return err
		}

		err = json.Unmarshal(bytes, cfg)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported settings file extension: %s", ext)
	}

	return nil
}

func loadFromEnvironmentVariables(cfg *Configuration) error {
	if cfg.Environment == "Development" {
		err := godotenv.Load(os.Getenv("DOTENV_FILE"))
		if err != nil {
			return err
		}
	}

	return parseEnv("", cfg)
}

func parseEnv(prefix string, to interface{}) error {
	if prefix != "" {
		prefix += "__"
	}

	v := reflect.ValueOf(to)
	if v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}

	vType := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := vType.Field(i)
		tag := fieldType.Tag.Get("key")

		if tag == "" || tag == "-" {
			continue
		}

		tagWithPrefix := prefix + tag

		if field.Kind() == reflect.Ptr && !field.IsNil() && field.Type().Elem().Kind() == reflect.Struct {
			err := parseEnv(tagWithPrefix, field.Interface())
			if err != nil {
				return err
			}
			continue
		}

		envValue := os.Getenv(tagWithPrefix)
		if envValue == "" {
			continue
		}

		switch field.Kind() {
		case reflect.String:
			field.SetString(envValue)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			intValue, err := strconv.ParseInt(envValue, 10, 64)
			if err != nil {
				return err
			}
			field.SetInt(intValue)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			uintValue, err := strconv.ParseUint(envValue, 10, 64)
			if err != nil {
				return err
			}
			field.SetUint(uintValue)
		case reflect.Float32, reflect.Float64:
			floatValue, err := strconv.ParseFloat(envValue, 64)
			if err != nil {
				return err
			}
			field.SetFloat(floatValue)
		case reflect.Bool:
			boolValue, err := strconv.ParseBool(envValue)
			if err != nil {
				return err
			}
			field.SetBool(boolValue)
		default:
			return fmt.Errorf("unsupported type for field '%s': %s", tag, field.Kind())
		}
	}

	return nil
}

func loadFromCommandLineArguments(cfg *Configuration) {

}
