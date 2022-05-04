// main code
// My first go program!
// CB: Michael Kukar 2022
package main

import (
	"fmt"
	"sort"

	"golang.org/x/exp/maps"

	"github.com/mkukar/anemoi/nowdata"
)

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
	// gathers all the datasets
	precipitationData := nowdata.PostStationData(stations[stationSlice[selectedStation]], "pcpn", "sum")
	snowfallData := nowdata.PostStationData(stations[stationSlice[selectedStation]], "snow", "sum")
	aveMonthlyTempData := nowdata.PostStationData(stations[stationSlice[selectedStation]], "avgt", "mean")
	lowMonthlyTempData := nowdata.PostStationData(stations[stationSlice[selectedStation]], "mint", "min")
	highMonthlyTempData := nowdata.PostStationData(stations[stationSlice[selectedStation]], "maxt", "max")

	fmt.Println(precipitationData)
	fmt.Println(snowfallData)
	fmt.Println(aveMonthlyTempData)
	fmt.Println(lowMonthlyTempData)
	fmt.Println(highMonthlyTempData)
	// TODO - parse datasets in an intelligent way

	// TODO - report data to user on global warming (interesting things only!)

}
