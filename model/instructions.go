package model

type Params struct {
	Tone     string `json:"tone"`
	Words    int16  `json:"words"`
	Hashtags bool   `json:"hashtags"`
	Emojis   bool   `json:"emojis"`
	Network  string `json:"network"`
	Context  bool   `json:"context"`
	Post     string `json:"posts"`
	Url      string `json:"url"`
}

type Meeting struct {
	Link      string `json:"link"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

type SystemInstructions struct {
	Prompt  string  `json:"prompt"`
	Params  Params  `json:"params"`
	Meeting Meeting `json:"meeting"`
}

type ScrapeResult struct {
	Description string `json:"description"`
	Logo        string `json:"logo"`
}
