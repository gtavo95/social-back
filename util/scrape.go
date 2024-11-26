package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"social/model"
)

func encode_url(baseURL string) string {
	// Parse the URL
	u, err := url.Parse(baseURL)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
	}

	// Print the full URL with query parameters
	fmt.Println("Full URL:", u.String())

	return u.String()

}

func Scrape_url(baseURL string) model.ScrapeResult {

	// Create client
	client := &http.Client{}

	encoded_url := encode_url(baseURL)

	fmt.Println("URL CODIFICADA:", encoded_url)

	// Create request
	req, err := http.NewRequest("GET", "https://app.scrapingbee.com/api/v1/?api_key=CLDW8SV4VVM4AFW82355BQ9PV8AZSR3G33ZQ2BTQFO8QIAPZEVJRKWDOLNGNRLJBCMEVLZG8Z32HC4LR&url="+encoded_url+"&json_response=true", nil)

	if err != nil {
		panic(err)
	}

	parseFormErr := req.ParseForm()
	if parseFormErr != nil {
		fmt.Println(parseFormErr)
	}

	// Fetch Request
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Failure : ", err)
	}

	// Read Response Body
	respBody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	// Unmarshal the JSON response
	var responseData map[string]interface{}

	if err := json.Unmarshal(respBody, &responseData); err != nil {
		fmt.Println("Error parsing JSON:", err)

	}

	var description, logo string

	if metadata, ok := responseData["metadata"].(map[string]interface{}); ok {
		// Check if "opengraph" exists and is a non-empty list
		if opengraph, ok := metadata["opengraph"].([]interface{}); ok && len(opengraph) > 0 {
			// Get the first element as a map
			if firstOpengraph, ok := opengraph[0].(map[string]interface{}); ok {
				// Check for "og:description" and save it if exists
				if desc, ok := firstOpengraph["og:description"].(string); ok && description == "" {
					description = desc
				}
				// If "og:description" is not found, try "og:title"
				if description == "" {
					if title, ok := firstOpengraph["og:title"].(string); ok {
						description = title
					}
				}
				// Check for "og:image" and save it as logo if exists
				if img, ok := firstOpengraph["og:image"].(string); ok && logo == "" {
					logo = img
				}
			}
		}

		// Check if "dublincore" exists and is a non-empty list
		if dublincore, ok := metadata["dublincore"].([]interface{}); ok && len(dublincore) > 0 {
			// Loop through each element in "dublincore"
			for _, dc := range dublincore {
				// Check if each element is a map
				if dcMap, ok := dc.(map[string]interface{}); ok {
					// Check for "elements" key, which should be a list
					if elements, ok := dcMap["elements"].([]interface{}); ok && len(elements) > 0 {
						for _, element := range elements {
							// Each element should be a map, so assert it as a map
							if elemMap, ok := element.(map[string]interface{}); ok {
								// Look for "name" == "description" and retrieve "content"
								if name, ok := elemMap["name"].(string); ok && name == "description" && description == "" {
									if content, ok := elemMap["content"].(string); ok {
										description = content
										break // Exit loop once description is found
									}
								}
							}
						}
					}
				}
			}
		}

		// Check if "json-ld" exists and is a non-empty list
		if jsonLd, ok := metadata["json-ld"].([]interface{}); ok && len(jsonLd) > 0 {
			// Access the first element as a map
			if firstJsonLd, ok := jsonLd[0].(map[string]interface{}); ok {
				// Try to find "image" key for logo
				if img, ok := firstJsonLd["image"].(string); ok && logo == "" {
					logo = img
				}
				// If description is still empty, try to find publisher->name
				if description == "" {
					if publisher, ok := firstJsonLd["publisher"].(map[string]interface{}); ok {
						if name, ok := publisher["name"].(string); ok {
							description = name
						}
					}
				}
			}
		}
	}

	/*
		if metadata, ok := responseData["metadata"].(map[string]interface{}); ok {

			// Check if "opengraph" exists and is a non-empty list
			if opengraph, ok := metadata["opengraph"].([]interface{}); ok && len(opengraph) > 0 {
				// Get the first element as a map
				if firstOpengraph, ok := opengraph[0].(map[string]interface{}); ok {
					// Check for "og:description" and save it if exists
					if desc, ok := firstOpengraph["og:description"].(string); ok {
						description = desc
					}
					// Check for "og:image" and save it as logo if exists
					if img, ok := firstOpengraph["og:image"].(string); ok {
						logo = img
					}
				}
			} else if dublincore, ok := metadata["dublincore"].([]interface{}); ok && len(dublincore) > 0 {
				// Loop through each element in "dublincore"
				for _, dc := range dublincore {
					// Check if each element is a map
					if dcMap, ok := dc.(map[string]interface{}); ok {
						// Check for "elements" key, which should be a list
						if elements, ok := dcMap["elements"].([]interface{}); ok && len(elements) > 0 {
							for _, element := range elements {
								// Each element should be a map, so assert it as a map
								if elemMap, ok := element.(map[string]interface{}); ok {
									// Look for "name" == "description" and retrieve "content"
									if name, ok := elemMap["name"].(string); ok && name == "description" {
										if content, ok := elemMap["content"].(string); ok {
											description = content
											break // Exit loop once description is found
										}
									}
								}
							}
						}
					}
				}
			}
		}

	*/

	// response := map[string]interface{}{
	// 	"description": description,
	// 	"logo":        logo,
	// }

	result := model.ScrapeResult{
		Description: description,
		Logo:        logo,
	}

	// Send the response as JSON
	return result

}
