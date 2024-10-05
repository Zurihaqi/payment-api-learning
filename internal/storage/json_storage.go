package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"payment-api-learning/internal/models"
)

type Storage interface {
	GetCustomerByUsername(username string) (*models.Customer, error)
	GetCustomerByID(id string) (*models.Customer, error)
	LogActivity(userID, action, details string) error
}

type JSONStorage struct {
	customersFile string
	logsFile      string
	customers     []models.Customer
	logs          []models.LogEntry
}

func NewJSONStorage(customersFile, logsFile string) (*JSONStorage, error) {
	storage := &JSONStorage{
		customersFile: customersFile,
		logsFile:      logsFile,
	}

	if err := storage.loadCustomers(); err != nil {
		return nil, err
	}

	if err := storage.loadLogs(); err != nil {
		return nil, err
	}

	return storage, nil
}

func (s *JSONStorage) loadCustomers() error {
	file, err := os.Open(s.customersFile)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &s.customers)
}

func (s *JSONStorage) loadLogs() error {
	data, err := os.ReadFile(s.logsFile)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &s.logs)
}

func (s *JSONStorage) GetCustomerByUsername(username string) (*models.Customer, error) {
	for _, customer := range s.customers {
		if customer.Username == username {
			return &customer, nil
		}
	}
	return nil, fmt.Errorf("customer not found")
}

func (s *JSONStorage) GetCustomerByID(id string) (*models.Customer, error) {
	for _, customer := range s.customers {
		if customer.ID == id {
			return &customer, nil
		}
	}
	return nil, fmt.Errorf("customer not found")
}

func (s *JSONStorage) LogActivity(userID, action, details string) error {
	s.logs = append(s.logs, models.LogEntry{
		Timestamp: time.Now(),
		Action:    action,
		UserID:    userID,
		Details:   details,
	})

	logData, err := json.MarshalIndent(s.logs, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.logsFile, logData, 0644)
}
