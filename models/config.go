package models

type ConfigModel struct {
	SrcListenAddr        string `json:"src"`
	AuxiliaryNodeSVCAddr string `json:"anSVCAddr"`
}
