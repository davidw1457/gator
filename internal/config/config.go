package config

import (
    "encoding/json"
    "fmt"
    "os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
    DbUrl string `json:"db_url"`
    CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {
    jsonPath, err := getJsonPath()
    if err != nil {
        return Config{}, fmt.Errorf("config.Read: %w", err)
    }

    jsonData, err := os.ReadFile(jsonPath)
    if err != nil {
        return Config{}, fmt.Errorf("config.Read: %w", err)
    }

    config := Config{}
    err = json.Unmarshal(jsonData, &config)
    if err != nil {
        return Config{}, fmt.Errorf("config.Read: %w", err)
    }
    return config, nil
}

func (c *Config) SetUser(userName string) error {
    jsonPath, err := getJsonPath()
    if err != nil {
        return fmt.Errorf("config.SetUser: %w", err)
    }

    c.CurrentUserName = userName

    jsonData, err := json.Marshal(c)
    if err != nil {
        return fmt.Errorf("config.SetUser: %w", err)
    }

    err = os.WriteFile(jsonPath, jsonData, 0660)
    if err != nil {
        return fmt.Errorf("config.SetUser: %w", err)
    }

    return nil
}

func getJsonPath() (string, error) {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return "", fmt.Errorf("config.getJsonPath: %w", err)
    }

    jsonPath := homeDir + string(os.PathSeparator) + configFileName

    return jsonPath, nil
}
