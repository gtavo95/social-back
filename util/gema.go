package util

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"social/model"
	"strings"

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

	apiKey := os.Getenv("GEMINI_API_KEY")

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

}

func (g *Gem) CreateSystemStruction(params model.Params, identity string) string {

	withHashtags := "no"
	if params.Hashtags {
		withHashtags = ""
	}

	withEmojis := "no"
	if params.Emojis {
		withEmojis = ""
	}

	withContext := ""
	if params.Context {
		withContext = "To do this, use the identity of the company and the context of the brand" + identity
	}

	instruction := fmt.Sprintf(`You are a creative assistant guiding a marketer to create content about %s. For %s.
		Task Details:
		Generate markdown instructions with exactly %d words.
		Include %s hashtags (at the end of the caption).
		Add %s relevant emojis spread throughout the content.
		Ensure the tone is %s.
		Provide %s options of captions.
		Mandatory Details:
		Include the URL: http://localhost:3000/honey/90140547.
		Add the meeting date and time: Use the current time + 30 minutes, in this format: MonthName/Day HH:MM.
		Use formatting for clarity: Ensure captions have line breaks and spaces to enhance readability.
		Generate captions in JSON format as follows: 
		[[ 
		  { \"caption\": \" Title in markdown\nA brief description in *markdown* format.\" },
		  { \"caption\": \" Another title\nAdditional details in **markdown**.\" } 
		]]. Ensure all text have proper formats, add bolds and enters where is necessary.`,
		withContext, params.Network, params.Words, withHashtags, withEmojis, params.Tone, params.Post)

	return instruction
}

func (g *Gem) SetSystemInstructions(instruction string) {
	g.Model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(instruction)},
	}
}

// func (g *Gem) UploadToGemini(ctx context.Context, file multipart.File, name, mimeType string) string {
//
// 	options := genai.UploadFileOptions{
// 		DisplayName: name,
// 		MIMEType:    mimeType,
// 	}
// 	fileData, err := g.Client.UploadFile(ctx, "", file, &options)
// 	if err != nil {
// 		log.Fatalf("Error uploading file: %v", err)
// 	}
//
// 	return fileData.URI
// }

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

func (g *Gem) SetSessionSimple() {
	session := g.Model.StartChat()
	g.Session = session
	g.Session.History = []*genai.Content{}

}

func (g *Gem) UploadImageFromURL(imgURL string) genai.Blob {
	response, err := http.Get(imgURL)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, response.Body)
	if err != nil {
		panic(err)
	}

	parts := strings.Split(imgURL, ".")
	imgType := parts[len(parts)-1]

	genaiImgData := genai.ImageData(imgType, buf.Bytes())
	return genaiImgData

}

func (g *Gem) SendRequest(ctx context.Context, prompt string) []genai.Part {
	log.Println(prompt)
	resp, err := g.Session.SendMessage(ctx, genai.Text(prompt))
	if err != nil {
		log.Println("Error sending message", err.Error())
		panic(err)
	}
	// for _, part := range resp.Candidates[0].Content.Parts {
	// 	fmt.Printf("%v\n", part)
	// }
	return resp.Candidates[0].Content.Parts
}
