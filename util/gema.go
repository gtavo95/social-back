package util

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"social/model"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type Gem struct {
	Client  *genai.Client
	Model   *genai.GenerativeModel
	Session *genai.ChatSession
}

func (g *Gem) Init() {
	ctx := context.Background()

	// apiKey, ok := os.LookupEnv("GEMINI_API_KEY")
	apiKey := os.Getenv("GEMINI_API_KEY")
	log.Println("apikey", apiKey)

	// if !ok {
	// 	log.Fatalln("Environment variable GEMINI_API_KEY not set")
	// }

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}
	g.Client = client
}

func (g *Gem) SetModel() {
	model := g.Client.GenerativeModel("gemini-1.5-flash")
	g.Model = model

	g.Model.SetTemperature(1)
	g.Model.SetTopK(64)
	g.Model.SetTopP(0.95)
	g.Model.SetMaxOutputTokens(8192)
	g.Model.ResponseMIMEType = "application/json"
	g.Model.SafetySettings = []*genai.SafetySetting{
		{
			Category:  genai.HarmCategoryHarassment,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategoryHateSpeech,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategorySexuallyExplicit,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategoryDangerousContent,
			Threshold: genai.HarmBlockNone,
		},
	}
	// g.Model.ResponseMIMEType = "text/plain"

}

func (g *Gem) CreateSystemStruction(params model.Params) string {

	withHashtags := "no"
	if params.Hashtags {
		withHashtags = "yes"
	}

	withEmojis := "no"
	if params.Emojis {
		withEmojis = "yes"
	}

	withContext := ""
	if params.Context {
		withContext = "To do this, you have in I have upload files about the identity of the company and the context of the brand"
	}

	instruction := fmt.Sprintf("You are a creative assistant, who seeks to guide and help a marketer to create content on %s. Generate markdown instruction with exactly %d words, including %s hashtags and %s relevant emojis. Ensure the tone is %s. %s. Provide %s options. The answer in Json format [{caption: '', caption: '', ...}]",
		params.Network, params.Words, withHashtags, withEmojis, params.Tone, withContext, params.Post)

	log.Println(instruction)
	return instruction
}

func (g *Gem) SetSystemInstructions(instruction string) {
	g.Model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(instruction)},
	}
}

func (g *Gem) UploadToGemini(ctx context.Context, file multipart.File, name, mimeType string) string {

	options := genai.UploadFileOptions{
		DisplayName: name,
		MIMEType:    mimeType,
	}
	fileData, err := g.Client.UploadFile(ctx, "", file, &options)
	if err != nil {
		log.Fatalf("Error uploading file: %v", err)
	}

	return fileData.URI
}

func (g *Gem) SetSession(parts []genai.Part) {
	session := g.Model.StartChat()
	g.Session = session
	g.Session.History = []*genai.Content{
		{
			Role:  "user",
			Parts: parts,
		},
	}

}

func (g *Gem) SendRequest(ctx context.Context, prompt string) []genai.Part {
	log.Println(prompt)
	resp, err := g.Session.SendMessage(ctx, genai.Text(prompt))
	if err != nil {
		log.Println("Error sending message: %v", err.Error())
		panic(err)
	}
	for _, part := range resp.Candidates[0].Content.Parts {
		fmt.Printf("%v\n", part)
	}
	return resp.Candidates[0].Content.Parts
}
