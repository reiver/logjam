package models

type ConfigModel struct {
	SrcListenAddr      string `json:"src"`
	GoldGorillaSVCAddr string `json:"ggSVCAddr"`
}
