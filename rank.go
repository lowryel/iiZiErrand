package rank

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "net/url"
    "os"
	"github.com/eugene/iizi_errand/pkg/models"
)



func RankRunner(taskAttributes []string, runnerAttributes []string) {
	// ranking := []string{}
	matchCount := 0
	for _, taskAttr := range taskAttributes {
		for _, runnerAttr := range runnerAttributes {
			if runnerAttr == taskAttr {
				log.Printf("there's a match %v", runnerAttr)
				// ranking = append(ranking, runnerAttr)
				matchCount++
			}
		}
	}
	ranked := matchCount >= len(taskAttributes)/2
	if ranked {
		log.Printf("there is a %d chance", matchCount*100/len(taskAttributes))
	}else {
		log.Printf("low matching %d", matchCount*100/len(taskAttributes))
	}
}






func GetLocation() (*models.Location, error) {
    apiKey := os.Getenv("GEO_API_KEY")
    if apiKey == "" {
        return nil, fmt.Errorf("GEO_API_KEY is not set")
    }

    // Encode the API key
    encodedKey := url.QueryEscape(apiKey)
    url := fmt.Sprintf("https://api.ipgeolocation.io/ipgeo?apiKey=%s", encodedKey)

    // Create a new request
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        log.Printf("Error creating request: %v", err)
        return nil, err
    }

    // Set User-Agent header
    req.Header.Set("User-Agent", "YourAppName/1.0")

    // Send the request
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Printf("Error sending request: %v", err)
        return nil, err
    }
    defer resp.Body.Close()

    // Check the response status code
    if resp.StatusCode != http.StatusOK {
        bodyBytes, _ := ioutil.ReadAll(resp.Body)
        log.Printf("Unexpected status code: %d. Body: %s", resp.StatusCode, string(bodyBytes))
        return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
    }

    // Read the response body
    bodyBytes, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Printf("Error reading response body: %v", err)
        return nil, err
    }

    // Print the response body for debugging
    log.Printf("Response body: %s", string(bodyBytes))

    // Parse the response body
    var location models.Location
    if err := json.Unmarshal(bodyBytes, &location); err != nil {
        log.Printf("Error parsing JSON: %v", err)
        return nil, err
    }

    return &location, nil
}






