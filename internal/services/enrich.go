package services

import (
	"effective_mobile_test/internal/models"
	"encoding/json"
	"net/http"
)

type AgeResponse struct {
	Age int `json:"age"`
}

type GenderResponse struct {
	Gender string `json:"gender"`
}

type NationalityResponse struct {
	Country []struct {
		CountryID string `json:"country_id"`
	} `json:"country"`
}

func EnrichPersonWithAge(person *models.Person) error {
	resp, err := http.Get("https://api.agify.io/?name=" + person.Name)
	if err != nil {
		return err
	}

	var ageResponse AgeResponse
	err = json.NewDecoder(resp.Body).Decode(&ageResponse)
	if err != nil {
		return err
	}

	person.Age = ageResponse.Age
	err = resp.Body.Close()
	if err != nil {
		return err
	}
	return nil
}

func EnrichPersonWithGender(person *models.Person) error {
	resp, err := http.Get("https://api.genderize.io/?name=" + person.Name)
	if err != nil {
		return err
	}

	var genderResponse GenderResponse
	err = json.NewDecoder(resp.Body).Decode(&genderResponse)
	if err != nil {
		return err
	}

	person.Gender = genderResponse.Gender
	err = resp.Body.Close()
	if err != nil {
		return err
	}
	return nil
}

func EnrichPersonWithNationality(person *models.Person) error {
	resp, err := http.Get("https://api.nationalize.io/?name=" + person.Name)
	if err != nil {
		return err
	}

	var nationalityResponse NationalityResponse
	err = json.NewDecoder(resp.Body).Decode(&nationalityResponse)
	if err != nil {
		return err
	}

	if len(nationalityResponse.Country) > 0 {
		person.Nationality = nationalityResponse.Country[0].CountryID
	} else {
		person.Nationality = ""
	}

	err = resp.Body.Close()
	if err != nil {
		return err
	}
	return nil
}
