package pokeAPI

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ziadehab433/pokedexcli/internal/pokecache"
)

type pokeData struct {
	Count    int
	Next     string
	Previous any
	Results  []struct {
		Name string
		URL  string
	}
}

type locationArea struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

var cache pokecache.Cache = *pokecache.NewCache(50 * time.Second)

func GetLocations(url string) (pokeData, error) {
	body, exist := cache.Get(url)
	if !exist {
		res, err := http.Get(url)
		if err != nil {
			return pokeData{}, err
		}

		newBody, err := io.ReadAll(res.Body)
		if err != nil {
			return pokeData{}, err
		}

		cache.Add(url, newBody)

		pokedata := pokeData{}
		err = json.Unmarshal(newBody, &pokedata)
		if err != nil {
			return pokeData{}, err
		}

		return pokedata, nil
	}

	pokedata := pokeData{}
	err := json.Unmarshal(body, &pokedata)
	if err != nil {
		fmt.Println(body)
		return pokeData{}, err
	}

	return pokedata, nil
}

func ExploreLocation(url string) (locationArea, error) {
	body, exist := cache.Get(url)
	if !exist {
		res, err := http.Get(url)
		if err != nil {
			return locationArea{}, err
		}

		newBody, err := io.ReadAll(res.Body)
		if err != nil {
			return locationArea{}, err
		}

		cache.Add(url, newBody)

		locationData := locationArea{}
		err = json.Unmarshal(newBody, &locationData)
		if err != nil {
			return locationArea{}, errors.New("invalid location")
		}

		return locationData, nil
	}

	locationData := locationArea{}
	err := json.Unmarshal(body, &locationData)
	if err != nil {
		return locationArea{}, errors.New("invalid location")
	}

	return locationData, nil
}
