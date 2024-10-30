package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"social/model"
	"social/util"
	"strconv"

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

func SocialPostText(c *fiber.Ctx) error {

	instructions := c.FormValue("instructions")
	length := c.FormValue("length")

	var systemInstructions model.SystemInstructions
	err := json.Unmarshal([]byte(instructions), &systemInstructions)
	if err != nil {
		panic(err)
	}

	gem := util.Gem{}
	gem.Init()
	gem.SetModel()
	sysInstr := gem.CreateSystemStruction(systemInstructions.Params)
	gem.SetSystemInstructions(sysInstr)

	defer gem.Client.Close()

	ctx := context.Background()
	if length != "" {
		lengthInt, err := strconv.Atoi(length)
		if err != nil {
			log.Println("Error converting length to integer:", err)
			return c.Status(fiber.StatusBadRequest).JSON("Invalid length")
		}

		var genPars []genai.Part

		for i := 0; i < lengthInt; i++ {
			file, err := c.FormFile("file" + strconv.Itoa(i+1))

			if err != nil {
				panic(err)
			}

			f, err := file.Open()
			if err != nil {
				panic(err)
			}
			uri := gem.UploadToGemini(ctx, f, file.Filename, "application/pdf")
			genPars = append(genPars, genai.FileData{URI: uri})

		}
		gem.SetSession(genPars)
	}

	// generate parts
	parts := gem.SendRequest(ctx, systemInstructions.Prompt)

	log.Println(parts)

	// Dark Forest Part
	bf := util.BlackForest{}
	bf.Init()
	bf.SetPrompt(partsToStrings(parts))

	var samples []string

	for j := 0; j < 3; j++ {
		id := bf.Request()
		sample := bf.Poll(id)
		samples = append(samples, sample)

	}
	// log.Println(id)
	// sample := bf.Poll(id)

	result := map[string]interface{}{
		"result":  parts,
		"samples": samples,
		"status":  200,
	}

	return c.Status(fiber.StatusOK).JSON(result)

}
