package conf

import (
	"os"
	"tv/models"

	"github.com/goccy/go-yaml"
)

var VideoSources map[string]models.VideoSource

func LoadVideoSources(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, &VideoSources)
	if err != nil {
		return err
	}
	return nil
}

// 根据 key 获取视频源
func GetVideoSource(key string) (models.VideoSource, bool) {
	source, exists := VideoSources[key]
	return source, exists
}

// 获取所有视频源
func GetAllVideoSources() map[string]models.VideoSource {
	return VideoSources
}

// todo 获取所有激活的视频源（排除成人内容）
func GetActiveVideoSources() map[string]models.VideoSource {
	sources := make(map[string]models.VideoSource)
	for key, source := range VideoSources {
		if !source.Adult {
			sources[key] = source
		}
	}
	return sources
}
