package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/valtlfelipe/secret-editor/backend/types"
)

type PreferencesStorage struct {
	// mutex       sync.Mutex
	Preferences types.Preferences
}

func NewPreferences() *PreferencesStorage {
	return &PreferencesStorage{}
}

func getConfigPath() (path string, err error) {
	value, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	path = filepath.Join(value, "secret-editor")

	return path, nil
}

func getConfigFilePath() (filePath string, err error) {
	path, err := getConfigPath()
	if err != nil {
		return "", err
	}
	filePath = filepath.Join(path, "config.json")

	return filePath, nil
}

func writeFile(path string, data *types.Preferences) error {
	file, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = os.WriteFile(path, []byte(file), 0755)
	if err != nil {
		return err
	}

	return nil
}

func (p *PreferencesStorage) LoadPreferences() error {

	path, err := getConfigPath()
	if err != nil {
		return err
	}
	fullPath, _ := getConfigFilePath()

	// Check for file existance
	if _, err := os.Stat(fullPath); errors.Is(err, os.ErrNotExist) {

		// Create app's config directory
		err = os.Mkdir(path, 0700)
		if err != nil {
			return err
		}

		// Create file
		err := writeFile(fullPath, &p.Preferences)
		if err != nil {
			return err
		}

	} else if err != nil {
		return err
	}

	// Open file
	file, err := os.Open(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Read the file content
	content, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	// Parse the JSON content
	err = json.Unmarshal(content, &p.Preferences)
	if err != nil {
		return err
	}

	return nil
}

func (p *PreferencesStorage) SetPreference(key string, value string) error {
	parts := strings.Split(key, ".")
	if len(parts) > 0 {
		var reflectValue reflect.Value
		if reflect.TypeOf(&p.Preferences).Kind() == reflect.Ptr {
			reflectValue = reflect.ValueOf(&p.Preferences).Elem()
		} else {
			reflectValue = reflect.ValueOf(&p.Preferences)
		}
		for i, part := range parts {
			part = strings.ToUpper(part[:1]) + part[1:]
			reflectValue = reflectValue.FieldByName(part)
			if reflectValue.IsValid() {
				if i == len(parts)-1 {
					reflectValue.Set(reflect.ValueOf(value))

					// TODO: handle error
					err := p.SavePreference()
					if err != nil {
						log.Printf("### Error saving preference: %v", err)
					}
					return nil
				}
			} else {
				break
			}
		}
	}

	return fmt.Errorf("invalid key path(%s)", key)
}

func (p *PreferencesStorage) SavePreference() error {
	fullPath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	err = writeFile(fullPath, &p.Preferences)
	if err != nil {
		return err
	}

	return nil
}
