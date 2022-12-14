package data

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"regexp"
	"sort"

	"github.com/techmet/round-the-world/models"
	"github.com/techmet/round-the-world/utils"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const citiesURL = "https://s3.us-west-2.amazonaws.com/secure.notion-static.com/4be05480-e7fc-4b41-b642-fb26dcaa4c39/cities.json?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Content-Sha256=UNSIGNED-PAYLOAD&X-Amz-Credential=AKIAT73L2G45EIPT3X45%2F20221013%2Fus-west-2%2Fs3%2Faws4_request&X-Amz-Date=20221013T164617Z&X-Amz-Expires=86400&X-Amz-Signature=7c5ec6581ef355024f42a3cd1df3dfe1509e1af0be78d7948b24a67f5b87bf37&X-Amz-SignedHeaders=host&response-content-disposition=filename%20%3D%22cities.json%22&x-id=GetObject"

var cityMap models.CityMap = models.CityMap{}
var continentCitiesMap models.ContinentCityMap = models.ContinentCityMap{}
var totalContinents []string = []string{}

func GetCitiesFromURL() models.CityMap {
	if len(cityMap) == 0 {
		req, err := http.NewRequest(http.MethodGet, citiesURL, nil)
		if err != nil {
			fmt.Printf("client: could not create request: %s\n", err)
			os.Exit(1)
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Printf("client: error making http request: %s\n", err)
			os.Exit(1)
		}

		fmt.Printf("client: status code: %d\n", res.StatusCode)

		body, err := httputil.DumpResponse(res, true)

		if err != nil {
			log.Fatal(err)
		}

		// wrap the cached response
		r := bufio.NewReader(bytes.NewReader(body))
		res, err = http.ReadResponse(r, nil)

		if err != nil {
			log.Fatal(err)
		}

		byteValue, err := io.ReadAll(res.Body)

		if err != nil {
			fmt.Printf("client: error reading http response: %s\n", err)
			os.Exit(1)
		}

		var result models.CityMap
		json.Unmarshal([]byte(byteValue), &result)
		cityMap = result
		return result
	} else {
		return cityMap
	}
}

func GetCities() models.CityMap {
	if len(cityMap) == 0 {
		byteValue, err := os.ReadFile("data/cities.json")
		if err != nil {
			fmt.Print("Error in reading file: ", err)
			os.Exit(1)
		}
		var result models.CityMap
		json.Unmarshal(byteValue, &result)
		cityMap = result
		return result
	} else {
		return cityMap
	}
}

func GetCitiesForSearch() []models.SearchCity {
	citiesMap := GetCities()
	cities := make([]models.SearchCity, 0, len(citiesMap))
	for _, value := range citiesMap {
		cities = append(cities, models.SearchCity{
			ID:   value.ID,
			Name: value.Name,
		})
	}
	return cities
}

func GetContinentCityMap() models.ContinentCityMap {
	if len(continentCitiesMap) == 0 {
		citiesMap := GetCities()
		continentCityMap := models.ContinentCityMap{}
		for _, value := range citiesMap {
			if val, ok := continentCityMap[value.ContID]; ok {
				continentCityMap[value.ContID] = append(val, value)
			} else {
				cities := make([]*models.City, 0)
				cities = append(cities, value)
				continentCityMap[value.ContID] = cities
			}
		}
		continentCitiesMap = continentCityMap
	}
	return continentCitiesMap
}

func GetContinents(contientMap models.ContinentCityMap) []string {
	if len(totalContinents) == 0 {
		continents := make([]string, 0, len(contientMap))
		for k := range contientMap {
			continents = append(continents, k)
		}
		totalContinents = continents
	}
	return totalContinents
}

func PopulateNeighboursForEachCity() {
	continentMap := GetContinentCityMap()
	continents := GetContinents(continentMap)
	for continent, cities := range continentMap {
		otherContinents := getOtherContinents(continent, continents)
		for _, city := range cities {
			city.NeighbouringCities = createNeighbours(*city, otherContinents)
		}
	}
}

func GetRoundTrip(cityId string) models.TripDetails {
	currentCity := cityMap[cityId]
	continentMap := GetContinentCityMap()
	continents := GetContinents(continentMap)
	otherContinents := getOtherContinents(currentCity.ContID, continents)

	path := []string{
		getPathTextFromCity(*currentCity),
	}
	return getShortestPath(*currentCity, currentCity.ID, otherContinents, models.TripDetails{
		ID:            currentCity.ID,
		Path:          path,
		TotalDistance: 0,
		Coordinates: []models.LocationDetails{{
			Lat: currentCity.Location.Lat,
			Lon: currentCity.Location.Lon,
		}},
	})

}

func getShortestPath(originalCity models.City, currentCityID string, remainContinents []string, tripDetails models.TripDetails) models.TripDetails {
	cityMap := GetCities()
	currentCity := *cityMap[currentCityID]

	if len(remainContinents) == 0 {
		tripDetails.Path = append(tripDetails.Path, originalCity.ID+" (Back to "+formatName(originalCity.Name)+")")
		tripDetails.TotalDistance = tripDetails.TotalDistance + utils.GetDistanceFromLatLonInKm(currentCity, originalCity)
		return tripDetails
	}

	continentMap := currentCity.NeighbouringCities

	var nextCity *models.NeighbouringCity

	for _, remainContinent := range remainContinents {
		if nextCity == nil || nextCity.Distance > continentMap[remainContinent][0].Distance {
			nextCity = (continentMap[remainContinent])[0]
		}
	}

	nextContinent := nextCity.ContID
	tripDetails.Path = append(tripDetails.Path, getPathText(*nextCity))
	tripDetails.Coordinates = append(tripDetails.Coordinates, models.LocationDetails{
		Lat: nextCity.Location.Lat,
		Lon: nextCity.Location.Lon,
	})
	tripDetails.TotalDistance = tripDetails.TotalDistance + nextCity.Distance

	updatedContinents := make([]string, 0, len(remainContinents)-1)

	for _, continent := range remainContinents {
		if continent != nextContinent {
			updatedContinents = append(updatedContinents, continent)
		}
	}
	return getShortestPath(originalCity, nextCity.ID, updatedContinents, tripDetails)
}

func getOtherContinents(currentContinent string, continents []string) []string {
	result := []string{}
	for _, continent := range continents {
		if currentContinent != continent {
			result = append(result, continent)
		}
	}
	return result
}

func createNeighbours(city models.City, otherContinents []string) models.NeighbouringContinentCityMap {
	otherContinentCityMap := models.NeighbouringContinentCityMap{}
	for _, continent := range otherContinents {
		neighbouringCities := make([]*models.NeighbouringCity, 0)
		for _, nextCity := range continentCitiesMap[continent] {
			distance := utils.GetDistanceFromLatLonInKm(city, *nextCity)
			neighbouringCities = append(neighbouringCities, &models.NeighbouringCity{
				ID:          nextCity.ID,
				Name:        nextCity.Name,
				CountryName: nextCity.CountryName,
				ContID:      nextCity.ContID,
				Distance:    distance,
				Location:    nextCity.Location,
			})
		}
		sort.Slice(neighbouringCities, func(i, j int) bool {
			return neighbouringCities[i].Distance < neighbouringCities[j].Distance
		})
		otherContinentCityMap[continent] = neighbouringCities
	}
	return otherContinentCityMap
}

func getPathText(city models.NeighbouringCity) string {
	return city.ID + " (" + formatName(city.Name) + ", " + formatName(city.ContID) + ")"
}

func getPathTextFromCity(city models.City) string {
	return city.ID + " (" + formatName(city.Name) + ", " + formatName(city.ContID) + ")"
}

func formatName(name string) string {
	re, err := regexp.Compile(`[^\w]`)
	if err != nil {
		log.Fatal(err)
	}
	return cases.Title(language.English).String(re.ReplaceAllString(name, " "))
}
