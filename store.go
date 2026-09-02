package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type LogEntry struct {
	Time    time.Time `json:"time"`
	Message string    `json:"message"`
	Type    string    `json:"type"`
	Pouches int       `json:"pouches,omitempty"`
}

type Purchase struct {
	Price    float64   `json:"price"`
	Pouches  int       `json:"pouches"`
	BoughtAt time.Time `json:"bought_at"`
	PerPouch float64   `json:"per_pouch"`
}

type Config struct {
	QuitDate        time.Time `json:"quit_date"`
	DosasPerDay     float64   `json:"dosas_per_day"`
	DosaPrice       float64   `json:"dosa_price"`
	PouchesPerDosa  int       `json:"pouches_per_dosa"`
	Goal            float64   `json:"goal"`
	GoalDescription string    `json:"goal_description"`
}

type Data struct {
	Config    Config     `json:"config"`
	Logs      []LogEntry `json:"logs"`
	Relapses  []LogEntry `json:"relapses"`
	Purchases []Purchase `json:"purchases"`
}

func dataPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".snusfria")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return filepath.Join(dir, "data.json"), nil
}

func LoadData() (*Data, error) {
	path, err := dataPath()
	if err != nil {
		return &Data{}, err
	}

	f, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Data{}, nil
		}
		return &Data{}, err
	}

	var data Data
	if err := json.Unmarshal(f, &data); err != nil {
		return &Data{}, fmt.Errorf("failed to parse data: %w", err)
	}
	return &data, nil
}

func SaveData(data *Data) error {
	path, err := dataPath()
	if err != nil {
		return err
	}

	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, b, 0644)
}
