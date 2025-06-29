package main

import (
	"testing"
)

type MockLogger struct {
	messages []string
	errors   []string
}

func (m *MockLogger) Log(message string) {
	m.messages = append(m.messages, message)
}

func (m *MockLogger) LogError(message string, err error) {
	m.errors = append(m.errors, message+": "+err.Error())
}

func TestCalculateElectricityCost(t *testing.T) {
	testCases := []struct {
		name         string
		tdpWatts     float64
		runtimeHours float64
		pricePerKWh  float64
		expectedCost float64
	}{
		{"Raspberry Pi 4 - 24h", 15, 24, 0.25, 0.09},
		{"Gaming PC - 5h", 300, 5, 0.25, 0.375},
		{"Home Server - 24h", 150, 24, 0.25, 0.9},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cost, err := CalculateElectricityCost(tc.tdpWatts, tc.runtimeHours, tc.pricePerKWh, nil)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			epsilon := 0.001
			if cost < tc.expectedCost-epsilon || cost > tc.expectedCost+epsilon {
				t.Errorf("Expected cost %.3f, got %.3f", tc.expectedCost, cost)
			}
		})
	}
}

func TestExceptionHandling(t *testing.T) {
	testCases := []struct {
		name        string
		tdpWatts    float64
		hours       float64
		pricePerKWh float64
		expectError bool
		expectedMsg string
	}{
		{
			name:        "Negative TDP Exception",
			tdpWatts:    -100,
			hours:       5,
			pricePerKWh: 0.25,
			expectError: true,
			expectedMsg: "tdp darf nicht negativ sein",
		},
		{
			name:        "Negative Hours Exception",
			tdpWatts:    100,
			hours:       -5,
			pricePerKWh: 0.25,
			expectError: true,
			expectedMsg: "laufzeit darf nicht negativ sein",
		},
		{
			name:        "Invalid Price Exception",
			tdpWatts:    100,
			hours:       5,
			pricePerKWh: -0.25,
			expectError: true,
			expectedMsg: "preis pro kwh muss positiv sein",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cost, err := CalculateElectricityCost(tc.tdpWatts, tc.hours, tc.pricePerKWh, nil)

			if tc.expectError {
				if err == nil {
					t.Fatalf("Expected error for test case: %s", tc.name)
				}
				if err.Error() != tc.expectedMsg {
					t.Errorf("Expected error message '%s', got '%s'", tc.expectedMsg, err.Error())
				}
				if cost != 0 {
					t.Errorf("Expected cost to be 0 when error occurs, got %.3f", cost)
				}
			} else {
				if err != nil {
					t.Fatalf("Expected no error for test case: %s, got %v", tc.name, err)
				}
			}
		})
	}
}

func TestDependencyInjection(t *testing.T) {
	mockLogger := &MockLogger{}

	cost, err := CalculateElectricityCost(100, 5, 0.20, mockLogger)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedCost := 0.1
	epsilon := 0.001
	if cost < expectedCost-epsilon || cost > expectedCost+epsilon {
		t.Errorf("Expected cost %.3f, got %.3f", expectedCost, cost)
	}

	if len(mockLogger.messages) == 0 {
		t.Error("Logger dependency not called through injection")
	}

	expectedMessages := 2
	if len(mockLogger.messages) != expectedMessages {
		t.Errorf("Expected %d log messages, got %d", expectedMessages, len(mockLogger.messages))
	}
}
