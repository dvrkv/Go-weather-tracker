package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
)

//Data struct to pass API key to it
type APIkey struct {
	OpenWeatherMapApiKey string `json:"OpenWeatherMapApiKey"`
}

//Weather data. Example of API response https://openweathermap.org/current#current_JSON
type weatherData struct {
	City string `json:"name"`
	Main struct {
		Temperature float32 `json:"temp"`
		Feels_like  float32 `json:"feels_like"`
		Humidity    int     `json:"humidity"`
	} `json:"main"`
}

//Function to get API key from .API-key file
func loadAPIkey(filename string) (APIkey, error) {
	//The ReadFile reads a file and returns the contents
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return APIkey{}, err
	}

	var key APIkey
	//Parsing data from JSON and saving the result to the "key" variable
	err = json.Unmarshal(bytes, &key)
	if err != nil {
		fmt.Println("Ooops :(\n Ошибка с парсингом JSON!")
	}
	return key, nil
}

//Application health check
func test(w http.ResponseWriter, r *http.Request) {
	tpl, _ := template.ParseFiles("static/test.html")
	tpl.Execute(w, nil)
}

//The query function takes a city name and returns the required data
func query(city string) (weatherData, error) {
	apiConfig, _ := loadAPIkey(".API-key")

	//API call. Built-in request by city name https://openweathermap.org/current#name
	response, err := http.Get("https://api.openweathermap.org/data/2.5/weather?q=" + city + "&appid=" + apiConfig.OpenWeatherMapApiKey + "&lang=" + "en" + "&units=" + "metric")
	if err != nil {
		return weatherData{}, err
	}

	defer response.Body.Close()

	var data weatherData
	//Decode a stream of distinct JSON values
	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		return weatherData{}, err
	}
	fmt.Println(data)
	return data, err
}

func weatherDataHandler(w http.ResponseWriter, r *http.Request) {
	city := strings.SplitN(r.URL.Path, "/", 3)[2]
	data, _ := query(city)
	//Response to frontend
	if data.Main.Temperature < 0 {
		t, _ := template.ParseFiles("static/0 and less.html")
		t.Execute(w, data)
	} else if data.Main.Temperature >= 0 && data.Main.Temperature <= 10 {
		t, _ := template.ParseFiles("static/0-10.html")
		t.Execute(w, data)
	} else if data.Main.Temperature > 10 && data.Main.Temperature <= 20 {
		t, _ := template.ParseFiles("static/10-20.html")
		t.Execute(w, data)
	} else if data.Main.Temperature > 20 && data.Main.Temperature <= 30 {
		t, _ := template.ParseFiles("static/20-30.html")
		t.Execute(w, data)
	} else if data.Main.Temperature > 30 {
		t, _ := template.ParseFiles("static/30 and more.html")
		t.Execute(w, data)
	}
}

//Call handlers
func main() {
	//Responds to the "/test" command and calls the test function
	http.HandleFunc("/test", test)

	//Responds to the command "/weather/" and calls the weatherDataHandler function
	http.HandleFunc("/weather/", weatherDataHandler)

	/*
		This function listens on a TCP network address,
		accepts incoming HTTP connections,
		and creates a separate goroutine for each connection
	*/
	port := ":8080"
	fmt.Print("Server listen on port", port, "\n")
	http.ListenAndServe(port, nil)
}
