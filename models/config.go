package models

// 视频源配置
type VideoSource struct {
	API    string `json:"api" yaml:"api"`
	Name   string `json:"name" yaml:"name"`
	Detail string `json:"detail,omitempty" yaml:"detail,omitempty"`
	Adult  bool   `json:"adult,omitempty" yaml:"adult,omitempty"`
}
