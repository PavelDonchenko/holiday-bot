package utils

import (
	"fmt"
	"math"
	"strings"
)

func GetUnitedCoordinate(longitude float64, latitude float64) string {
	return fmt.Sprintf("%f:%f", longitude, latitude)
}

func GetDividedCoordinate(coordinate string) []string {
	return strings.Split(coordinate, ":")
}

func Round(num float64, places int) float64 {
	shift := math.Pow(10, float64(places))
	return math.Round(num*shift) / shift
}
