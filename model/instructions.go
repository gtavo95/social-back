package model

type Params struct {
	Tone     string `json:"tone"`
	Words    int16  `json:"words"`
	Hashtags bool   `json:"hashtags"`
	Emojis   bool   `json:"emojis"`
	Network  string `json:"network"`
	Context  bool   `json:"context"`
	Post     string `json:"posts"`
}

type SystemInstructions struct {
	Prompt string `json:"prompt"`
	Params Params `json:"params"`
}
