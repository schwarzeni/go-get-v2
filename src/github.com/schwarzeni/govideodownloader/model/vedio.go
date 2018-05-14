package model

// http://bangumi.bilibili.com/player/
type Video struct {
	ApiUrl       string
	VedioSection []SingleVedioApi
	Cookie       string
	UserAgent    string
	WebPageUrl   string
}
