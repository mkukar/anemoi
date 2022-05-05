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

func summarizePrecipitationData(precipitationData map[int]interface{}) {
	annualRainfall := parsedata.GetAnnualSumPerYear(precipitationData)
	fmt.Println(annualRainfall)
	m, b := parsedata.LeastSquareFit(annualRainfall)
	fmt.Println(m, b)
	fmt.Println("PRECIPITATION TREND: ", m)
}

func summarizeSnowfallData(snowfallData map[int]interface{}) {
	annualSnowfall := parsedata.GetAnnualSumPerYear(snowfallData)
	fmt.Println(annualSnowfall)
	m, b := parsedata.LeastSquareFit(annualSnowfall)
	fmt.Println(m, b)
	fmt.Println("SNOWFALL TREND: ", m)
}

func summarizeMinTempData(minTempData map[int]interface{}) {
	minYearlyTemp := parsedata.GetAnnualMinPerYear(minTempData)
	fmt.Println(minYearlyTemp)
	m, b := parsedata.LeastSquareFit(minYearlyTemp)
	fmt.Println("MIN YEARLY TEMP: ", m, b)
}

func summarizeMaxTempData(maxTempData map[int]interface{}) {
	maxYearlyTemp := parsedata.GetAnnualMaxPerYear(maxTempData)
	fmt.Println(maxYearlyTemp)
	m, b := parsedata.LeastSquareFit(maxYearlyTemp)
	fmt.Println("MAX YEARLY TEMP: ", m, b)
}

func summarizeAverageTempData(aveTempData map[int]interface{}) {
	averageYearlyTemp := parsedata.GetAnnualAveragePerYear(aveTempData)
	fmt.Println(averageYearlyTemp)
	m, b := parsedata.LeastSquareFit(averageYearlyTemp)
	fmt.Println("AVERAGE YEARLY TEMP: ", m, b)
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
	summarizePrecipitationData(precipitationData)
	summarizeSnowfallData(snowfallData)
	summarizeMinTempData(lowMonthlyTempData)
	summarizeMaxTempData(highMonthlyTempData)
	summarizeAverageTempData(aveMonthlyTempData)
	fmt.Println("Looking to help? See actions you can take here https://www.un.org/actnow")
}
