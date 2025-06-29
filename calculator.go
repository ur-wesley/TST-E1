package main

import (
	"errors"
	"fmt"
	"time"
)

type Logger interface {
	Log(message string)
	LogError(message string, err error)
}

type ConsoleLogger struct{}

func (l *ConsoleLogger) Log(message string) {
	fmt.Printf("[info] %s\n", message)
}

func (l *ConsoleLogger) LogError(message string, err error) {
	fmt.Printf("[fehler] %s: %v\n", message, err)
}

type SilentLogger struct{}

func (l *SilentLogger) Log(message string) {
}

func (l *SilentLogger) LogError(message string, err error) {
}

func FormatCostEuro(cost float64) string {
	return fmt.Sprintf("%.2f €", cost)
}

func CalculateElectricityCost(tdpWatts, runtimeHours, pricePerKWh float64, logger Logger) (float64, error) {
	if logger != nil {
		logger.Log(fmt.Sprintf("berechne kosten für %.0fw gerät mit %.1f stunden laufzeit", tdpWatts, runtimeHours))
	}

	if tdpWatts < 0 {
		err := errors.New("tdp darf nicht negativ sein")
		if logger != nil {
			logger.LogError("ungültiger tdp wert", err)
		}
		return 0, err
	}
	if runtimeHours < 0 {
		err := errors.New("laufzeit darf nicht negativ sein")
		if logger != nil {
			logger.LogError("ungültige laufzeit", err)
		}
		return 0, err
	}
	if pricePerKWh <= 0 {
		err := errors.New("preis pro kwh muss positiv sein")
		if logger != nil {
			logger.LogError("ungültiger preis", err)
		}
		return 0, err
	}

	tdpKW := tdpWatts / 1000
	energyKWh := tdpKW * runtimeHours
	totalCost := energyKWh * pricePerKWh

	if logger != nil {
		logger.Log(fmt.Sprintf("berechnete kosten: %.2f € (%.3f kwh @ %.2f €/kwh)", totalCost, energyKWh, pricePerKWh))
	}
	return totalCost, nil
}

func CalculateMonthlyCost(tdpWatts, pricePerKWh float64, logger Logger) (float64, error) {
	hoursPerMonth := 24.0 * 30.0
	return CalculateElectricityCost(tdpWatts, hoursPerMonth, pricePerKWh, logger)
}

func ValidateDevice(tdpWatts float64, deviceName string) error {
	if tdpWatts < 0 {
		return errors.New("tdp darf nicht negativ sein")
	}
	if tdpWatts > 1000 {
		return fmt.Errorf("tdp von %.0fw scheint ungewöhnlich hoch für gerät '%s' - bitte überprüfen", tdpWatts, deviceName)
	}
	if deviceName == "" {
		return errors.New("gerätename darf nicht leer sein")
	}
	return nil
}

func GetCurrentElectricityPrice() (float64, error) {
	time.Sleep(100 * time.Millisecond)
	
	if time.Now().UnixNano()%13 == 0 {
		return 0, errors.New("netzwerk timeout: aktuelle strompreise konnten nicht von der api abgerufen werden")
	}
	
	return 0.30, nil
}
