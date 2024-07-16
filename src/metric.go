package src

import (
	"errors"
	"strconv"
)

var (
	Gauge      int8 = -1
	GaugeVec   int8 = -2
	Counter    int8 = -3
	CounterVec int8 = -4
)

type Metric struct {
	Name        string
	Help        string
	MetricType  int8
	FindKeyword string
	WithLabels  []WithLabels
}

type WithLabels struct {
	Key         string
	Value       string
	FindKeyword string
}

func FindValue(result map[string]string, findKeyword string) (string, float64, error) {
	stringValue, ok := result[findKeyword]
	if !ok {
		return "", 0, errors.New("Key '" + findKeyword + "' not found in the map")
	}

	// Convert string to float64
	tempValue, err := strconv.ParseFloat(stringValue, 64)
	if err != nil {
		return stringValue, 0, err
	}

	return stringValue, tempValue, nil
}
