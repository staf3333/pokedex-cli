package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/staf3333/pokedexcli/internal/pokecache"
)

// Struct for the JSON returned from the PokeAPI
// apparently the strings next to each field in the struct provide metadata about how
// the fields of the struct should be handled
// In context of JSON parsing and serialization, they define mapping between JSON keys and struct fields
// Particularly useful when the JSON field names don't match the Go struct field names exactly
// use upper case name if needed to be used across multiple packages
type PokemonLocationResponse struct {
	Count    int        `json:"count"`
	Next     string     `json:"next"`
	Previous *string    `json:"previous"`
	Results  []Location `json:"results"`
}

// Create separate structs for locations themselves
// use capital names for struct field so the `encoding/json` package can access them
// need to be able to marshall and unmarshal
type Location struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// implementing Cache:
// usual approach: check cache for requested resource. If found in cache, return it to client
// else, proceed down whatever logic you have to get the data!

func GetData(url string, cache *pokecache.Cache) PokemonLocationResponse{
	var body []byte
	// attempt to get data from the Cache first, if not found in the cache, get from the API
	body, ok := cache.Get(url)
	if !ok {		
		body = getFromPokeAPI(url)
		cache.Add(url, body)
	} else {
		fmt.Println("Found data in cache!")
	}
	// apiResponse is what we want to cache for a given url!
	apiResponse := PokemonLocationResponse{}
	err := json.Unmarshal(body, &apiResponse)
	if err != nil {			
		fmt.Println(err)
	}
	for _, location := range apiResponse.Results {
		fmt.Println(location.Name)
	}
	return apiResponse
}

func getFromPokeAPI(url string) ([]byte) {
	// base url for PokeAPI: https://pokeapi.co/api/v2/{endpoint}/
	// url for locations: https://pokeapi.co/api/v2/location/
	// list by default contains 20 resources
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("Error in api call")
	}
	// res contains req but use io.ReadAll to make code simpler
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if res.StatusCode > 299 {
		fmt.Printf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		fmt.Println("Error reading response body")
	}
	return body
}