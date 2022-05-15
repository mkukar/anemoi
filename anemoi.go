// main code
// My first go program!
// CB: Michael Kukar 2022
package main

import (
	"fmt"
	"sort"

	"golang.org/x/exp/maps"

	"github.com/mkukar/anemoi/nowdata"
	"github.com/mkukar/anemoi/parsedata"
)

// precipStr is precipitation or snowfall, since we use the same function
func summarizePrecipitationData(precipitationData map[int]interface{}, precipStr string) {
	annualPrecip := parsedata.GetAnnualSumPerYear(precipitationData)
	aveYearlyPrecip := parsedata.Average(annualPrecip)
	if aveYearlyPrecip == 0 { // no precipitation in this area
		return
	}
	m, _ := parsedata.LeastSquareFit(annualPrecip)
	if m < 0 {
		fmt.Println("- Yearly", precipStr, "is decreasing at a rate of", fmt.Sprintf("%.2f", -1*m), "inches per year.")
	}
	// checks historical data if precipitation is below average
	yearsPrecipBelowAverage := 0
	yearsPrecipBelowAverageInLastDecade := 0
	lastDecadeCounter := 10
	continueCountingYears := true
	for i := len(annualPrecip) - 1; i >= 0; i-- {
		if annualPrecip[i] < aveYearlyPrecip {
			if continueCountingYears {
				yearsPrecipBelowAverage += 1
			}
			if lastDecadeCounter > 0 {
				yearsPrecipBelowAverageInLastDecade += 1
			}
		} else if lastDecadeCounter == 0 {
			continueCountingYears = false
			break
		} else {
			continueCountingYears = false
		}
		lastDecadeCounter -= 1
	}

	if yearsPrecipBelowAverageInLastDecade > 1 {
		fmt.Println("- Total", precipStr, "has been below average for", yearsPrecipBelowAverageInLastDecade, "of the last 10 years")
	}
	if yearsPrecipBelowAverage > 1 {
		fmt.Println("- For the past", yearsPrecipBelowAverage, "years straight,", precipStr, "has been below the historical average of", fmt.Sprintf("%.2f", aveYearlyPrecip), "inches")
	}
}

func summarizeMinTempData(minTempData map[int]interface{}) {
	minYearlyTemp := parsedata.GetAnnualMinPerYear(minTempData)
	historicalMin, err := parsedata.Min(minYearlyTemp)
	if err != nil {
		return // data is bad, so just skip summary
	}
	degreesAboveHistoricalMinLastDecade := make([]float64, 10)
	lastDecadeCounter := 0
	for i := len(minYearlyTemp) - 1; i >= len(minYearlyTemp)-10; i-- {
		degreesAboveHistoricalMinLastDecade[lastDecadeCounter] = minYearlyTemp[i] - historicalMin
		lastDecadeCounter += 1
	}
	averageDegreesAboveMinLastDecade := parsedata.Average(degreesAboveHistoricalMinLastDecade)
	if averageDegreesAboveMinLastDecade > 0 {
		fmt.Println("- In the past decade, the average yearly minimum has been", averageDegreesAboveMinLastDecade, "degrees above the historical minimum temperature of", historicalMin, "degrees")
	}
	m, _ := parsedata.LeastSquareFit(minYearlyTemp)
	if m > 0 {
		fmt.Println("- The yearly minimum temperature has been increasing annually by", fmt.Sprintf("%.2f", m), "degrees")
	}
}

func summarizeMaxTempData(maxTempData map[int]interface{}) {
	maxYearlyTemp := parsedata.GetAnnualMaxPerYear(maxTempData)
	m, _ := parsedata.LeastSquareFit(maxYearlyTemp)
	if m > 0 {
		fmt.Println("- The yearly maximum temperature has been increasing annually by", fmt.Sprintf("%.2f", m), "degrees")
	}
}

func summarizeAverageTempData(aveTempData map[int]interface{}) {
	averageYearlyTemp := parsedata.GetAnnualAveragePerYear(aveTempData)
	m, _ := parsedata.LeastSquareFit(averageYearlyTemp)
	if m > 0 {
		fmt.Println("- The yearly average temperature has been increasing annually by", fmt.Sprintf("%.2f", m), "degrees")
	}
}

func main() {
	fmt.Println("Hello, I am enemoi. Let us see how climate change has affected your local area. First, a few questions...")

	// Interviews user
	regions := nowdata.GetFunctionMap()
	var selectedRegionIndx int
	regionSlice := maps.Keys(regions)
	sort.Strings(regionSlice)
	fmt.Println("Please select the number corresponding to the closest area to where you live:")
	for i, v := range regionSlice {
		fmt.Println(i, v)
	}
	fmt.Scan(&selectedRegionIndx)

	fmt.Println("Now please select the closest matching station. Please note that some may have less data than others:")
	stations := nowdata.GetStationList(regions[regionSlice[selectedRegionIndx]])
	stationSlice := maps.Keys(stations)
	sort.Strings(stationSlice)
	var selectedStation int
	for i, v := range stationSlice {
		fmt.Println(i, v)
	}
	fmt.Scan(&selectedStation)

	fmt.Println("Thank you. I am now looking back through time at station", stationSlice[selectedStation], "in the region of", regionSlice[selectedRegionIndx])
	fmt.Println("Please wait a few moments...")
	// gathers all the datasets
	precipitationData := nowdata.PostStationData(stations[stationSlice[selectedStation]], "pcpn", "sum")
	snowfallData := nowdata.PostStationData(stations[stationSlice[selectedStation]], "snow", "sum")
	aveMonthlyTempData := nowdata.PostStationData(stations[stationSlice[selectedStation]], "avgt", "mean")
	lowMonthlyTempData := nowdata.PostStationData(stations[stationSlice[selectedStation]], "mint", "min")
	highMonthlyTempData := nowdata.PostStationData(stations[stationSlice[selectedStation]], "maxt", "max")

	fmt.Println("Summarizing how climate change has affected your local area...")
	summarizePrecipitationData(precipitationData, "precipitation")
	summarizePrecipitationData(snowfallData, "snowfall")
	summarizeMinTempData(lowMonthlyTempData)
	summarizeMaxTempData(highMonthlyTempData)
	summarizeAverageTempData(aveMonthlyTempData)
	fmt.Println("Looking to see how you can help? See actions you can take here https://www.un.org/actnow")
}
