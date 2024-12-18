package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"social/model"
	"social/util"

	"github.com/gofiber/fiber/v2"
	"github.com/google/generative-ai-go/genai"
)

func partsToStrings(parts []genai.Part) string {
	var result string
	for _, part := range parts {
		instruction := fmt.Sprintf("%s", part)
		result = result + " " + instruction
	}
	return result
}

// Define the struct to represent each promotion
type Promotion struct {
	Caption string `json:"caption"`
}

func CaptionStruct(data []genai.Part) []Promotion {
	// Variable to hold the parsed data

	log.Println("data", data)

	var promotions []Promotion

	var captions string
	for _, part := range data {
		captions = fmt.Sprintf("%s", part)
	}

	// // Parse the JSON
	err := json.Unmarshal([]byte(captions), &promotions)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		panic(err)
	}

	return promotions
}

func SocialPostText(c *fiber.Ctx) error {

	instructions := c.FormValue("instructions")

	var systemInstructions model.SystemInstructions
	err := json.Unmarshal([]byte(instructions), &systemInstructions)
	if err != nil {
		panic(err)
	}

	scrapeResult := util.Scrape_url(systemInstructions.Params.Url)

	//Nuevo codigo
	description := scrapeResult.Description
	logo := scrapeResult.Logo
	target_url := systemInstructions.Params.Url

	if description == "" && logo == "" {

		parsedURL, err := url.Parse(target_url)
		if err != nil {
			log.Println("Error parseando URL", target_url)
		}

		// Rebuild the URL with only the scheme and host
		target_url = fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host)
		log.Println("New URL", target_url)

		scrapeResult = util.Scrape_url(target_url)

		description = scrapeResult.Description
		logo = scrapeResult.Logo

	}

	gem := util.Gem{}
	gem.Init()
	gem.SetModel()
	log.Println("scrapeResult", scrapeResult)
	sysInstr := gem.CreateSystemStruction(systemInstructions.Params, systemInstructions.Meeting, description)
	gem.SetSystemInstructions(sysInstr)

	defer gem.Client.Close()

	ctx := context.Background()
	gem.SetSessionSimple()

	log.Println("start request")
	// generate parts
	parts := gem.SendRequest(ctx, systemInstructions.Prompt)
	promotions := CaptionStruct(parts)

	// Dark Forest Part
	bf := util.BlackForest{}
	bf.Init()
	//
	var samples []string

	log.Println("sysInstr", sysInstr)
	for _, promotion := range promotions {
		bf.SetPrompt(promotion.Caption)
		id := bf.Request()
		sample := bf.Poll(id)
		samples = append(samples, sample)
	}

	log.Println("parts", parts)

	result := map[string]interface{}{
		"result":  parts,
		"samples": samples,
		"status":  200,
	}

	return c.Status(fiber.StatusOK).JSON(result)

}
