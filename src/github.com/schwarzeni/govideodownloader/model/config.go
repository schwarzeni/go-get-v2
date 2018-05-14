package model

type Config struct {
	Data   []VedioInfoConfig `json:"data"`
	Cookie string            `json:"cookie"`
}

type VedioInfoConfig struct {
	ApiUrl     string `json:"apiUrl"`
	WebpageUrl string `json:"webpageUrl"`
	SavePath   string `json:"savePath"`
}
