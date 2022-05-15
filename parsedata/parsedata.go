// parses retrieved data into usable formats
// All input format is map[int]interface{}, mapping the year to an array of monthly values or T for missing
// Example:
// { 1850 : [1.00, 2.00, T, 3.00, ... T ]}
// CB: Michael Kukar - 2022

package parsedata

import (
	"fmt"
	"sort"
	"strconv"

	"golang.org/x/exp/maps"
)

func GetAnnualSumPerYear(nowdata map[int]interface{}) []float64 {
	yearlyData := make([]float64, len(nowdata))
	years := getSortedYears(nowdata)
	for i, year := range years {
		floatData := getInterfaceAsFloatSliceRemovingTEntries(nowdata[year])
		yearSum := Sum(floatData)
		yearlyData[i] = yearSum
	}
	return yearlyData
}

func GetAnnualMinPerYear(nowdata map[int]interface{}) []float64 {
	yearlyData := make([]float64, len(nowdata))
	years := getSortedYears(nowdata)
	for i, year := range years {
		floatData := getInterfaceAsFloatSliceRemovingTEntries(nowdata[year])
		yearMin, err := Min(floatData)
		if err == nil {
			yearlyData[i] = yearMin
		} else {
			yearlyData[i] = 0.0
		}
	}
	return yearlyData
}

func GetAnnualMaxPerYear(nowdata map[int]interface{}) []float64 {
	yearlyData := make([]float64, len(nowdata))
	years := getSortedYears(nowdata)
	for i, year := range years {
		floatData := getInterfaceAsFloatSliceRemovingTEntries(nowdata[year])
		yearMax, err := Max(floatData)
		if err == nil {
			yearlyData[i] = yearMax
		} else {
			yearlyData[i] = 0.0
		}
	}
	return yearlyData
}

func GetAnnualAveragePerYear(nowdata map[int]interface{}) []float64 {
	yearlyData := make([]float64, len(nowdata))
	years := getSortedYears(nowdata)
	for i, year := range years {
		floatData := getInterfaceAsFloatSliceRemovingTEntries(nowdata[year])
		yearAve := Average(floatData)
		yearlyData[i] = yearAve
	}
	return yearlyData
}

func Sum(arr []float64) float64 {
	result := 0.0
	for _, v := range arr {
		result += v
	}
	return result
}

func Min(arr []float64) (float64, error) {
	if len(arr) == 0 {
		return 0.0, fmt.Errorf("Min function must have an input array length greater than 0")
	}
	min := arr[0]
	for _, v := range arr {
		if v < min {
			min = v
		}
	}
	return min, nil
}

func Max(arr []float64) (float64, error) {
	if len(arr) == 0 {
		return 0.0, fmt.Errorf("Max function must have input array of length greater than 0")
	}
	max := arr[0]
	for _, v := range arr {
		if v > max {
			max = v
		}
	}
	return max, nil
}

func Average(arr []float64) float64 {
	sum := Sum(arr)
	return sum / float64(len(arr))
}

// returns m,b of y=mx + b
func LeastSquareFit(arr []float64) (float64, float64) {
	n := float64(len(arr))
	sumX, sumY, sumXY, sumXX := 0.0, 0.0, 0.0, 0.0
	for x, y := range arr {
		floatX := float64(x)
		sumX += floatX
		sumY += y
		sumXY += floatX * y
		sumXX += floatX * floatX
	}
	divisor := (n*sumXX - sumX*sumX)
	m := (n*sumXY - sumX*sumY) / divisor
	b := (sumXX*sumY - sumXY*sumX) / divisor
	return m, b
}

// interface is of type []string, but is cast from interface -> []interface -> []string
func getInterfaceAsFloatSliceRemovingTEntries(data interface{}) []float64 {
	//strData := data.([]string)
	var floatData []float64
	for _, val := range data.([]interface{}) {
		strVal := val.(string)
		floatVal, err := strconv.ParseFloat(strVal, 64)
		if err == nil {
			floatData = append(floatData, floatVal)
		}
	}
	return floatData
}

func getSortedYears(nowdata map[int]interface{}) []int {
	years := maps.Keys(nowdata)
	sort.Ints(years)
	return years
}
