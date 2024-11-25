package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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
	// length := c.FormValue("length")

	var systemInstructions model.SystemInstructions
	err := json.Unmarshal([]byte(instructions), &systemInstructions)
	if err != nil {
		panic(err)
	}

	scrapeResult := util.Scrape_url(systemInstructions.Params.Url)
	log.Println("scrapeResult", scrapeResult)

	gem := util.Gem{}
	gem.Init()
	gem.SetModel()
	sysInstr := gem.CreateSystemStruction(systemInstructions.Params, scrapeResult.Description)
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

	result := map[string]interface{}{
		"result":  parts,
		"samples": samples,
		"status":  200,
	}

	return c.Status(fiber.StatusOK).JSON(result)

}
