package models

// 视频源配置
type VideoSource struct {
	API    string `mapstructure:"api" json:"api"`
	Name   string `mapstructure:"name" json:"name"`
	Detail string `mapstructure:"detail,omitempty" json:"detail,omitempty"`
	Adult  bool   `mapstructure:"adult" json:"adult,omitempty"`
}
