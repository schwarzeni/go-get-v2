package model

type VedioApi struct {
	Durl SingleVedioApi
}

type SingleVedioApi struct {
	Length int    `json:"length"`
	Order  int    `json:"order"`
	Size   int    `json:"size"`
	Url    string `json:"url"`
}
