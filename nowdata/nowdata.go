// API for NOWData (NOAA Online Weather Data)
// reverse-engineered from this website https://www.weather.gov/wrh/Climate?wfo=sgx
// CB: Michael Kukar - 2022

package nowdata

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

func GetFunctionMap() map[string]string {
	// key = human readable name, value is the 3-letter abbrev
	// TODO - finish filling this out using this URL https://www.weather.gov/wrh/climate
	regionMap := map[string]string{
		"Los Angeles":   "lox",
		"San Diego":     "sgx",
		"San Francisco": "mtr",
		"Hanford":       "hnx",
		"Las Vegas":     "vef",
		"Sacramento":    "sto",
		"Eureka":        "eka",
		"Medford":       "mef",
		"Reno":          "rev",
	}
	return regionMap
}

// get returns list of stations in format [stationid, string description, something else?]
// parses into dict of key is description, val is station id
func GetStationList(regionKey string) map[string]string {
	stationListUrl := "http://nowdata.rcc-acis.org/%s/station_list.txt"
	response, err := http.Get(fmt.Sprintf(stationListUrl, regionKey))
	if err != nil {
		panic(err)
	}
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	var responseArr [][]string
	json.Unmarshal(responseData, &responseArr)
	stationMap := make(map[string]string)
	for i := range responseArr {
		stationMap[responseArr[i][1]] = responseArr[i][0]
	}
	return stationMap
}

// post station data
// name types:
// pcpn Precipitation
// return is map of year to the interface (array of year data, with T denoting no data found e.g. { 1850 : [1,2,T...12] })
func PostStationData(stationId string, dataName string, reductionStrategy string) map[int]interface{} {
	fmt.Println(stationId)
	stationDataUrl := "http://data.rcc-acis.org/StnData"
	// monthly, from beginning of dataset to present, group by year, reduce is how it summarizes the monthly data (e.g. sum)
	paramsStr := `{"elems":[{"interval":"mly","duration":1,"name":"%s","reduce":{"reduce":"%s"},"maxmissing":"1","prec":3,"groupby":["year",1,12]}],"sid":"%s","sDate":"por","eDate":"por"}`
	var params = []byte(fmt.Sprintf(paramsStr, dataName, reductionStrategy, stationId))
	req, _ := http.NewRequest("POST", stationDataUrl, bytes.NewBuffer(params))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		panic(response.Status)
	}
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	var responseMap map[string]interface{}
	err = json.Unmarshal(responseData, &responseMap)
	if err != nil {
		panic(err)
	}
	resultData := make(map[int]interface{})
	for k, v := range responseMap {
		if k == "data" {
			var data []interface{} = v.([]interface{})
			for _, entry := range data {
				entryYear, _ := strconv.Atoi(entry.([]interface{})[0].(string))
				var entryMonthList []interface{} = entry.([]interface{})[1].([]interface{})
				resultData[entryYear] = entryMonthList
			}
		}
	}
	return resultData

}
