package conf

import (
	"time"
	"tv/models"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// Config 配置结构体
type Config struct {
	App struct {
		APIVersion string `mapstructure:"api_version"`
		Password   string `mapstructure:"password"`
		Port       string `mapstructure:"port"`
	} `mapstructure:"app"`

	Cache struct {
		Search time.Duration `mapstructure:"search"`
		ID     time.Duration `mapstructure:"id"`
		Hot    time.Duration `mapstructure:"hot"`
	} `mapstructure:"cache"`

	Sources map[string]models.VideoSource `mapstructure:"sources"`
}

var Cfg Config

func InitConfig(configPath string) error {

	viper.SetConfigFile(configPath) // 指定配置文件路径
	viper.SetConfigType("yaml")     // 配置文件类型

	// 默认值
	viper.SetDefault("app.mode", "debug")
	viper.SetDefault("app.api_version", "v1")

	if err := viper.ReadInConfig(); err != nil {
		log.Err(err).Msg("读取配置文件失败")
		return err
	}

	// 解析到结构体
	if err := viper.Unmarshal(&Cfg); err != nil {
		log.Err(err).Msg("解析配置失败")
		return err
	}

	names := make([]string, 0, len(Cfg.Sources))
	for _, src := range Cfg.Sources {
		names = append(names, src.Name)
	}

	log.Info().Strs("数据源", names).Msg("已成功加载数据源")

	return nil
}

// 根据 key 获取视频源
func (cfg Config) GetVideoSource(key string) (models.VideoSource, bool) {
	source, exists := cfg.Sources[key]
	return source, exists
}

// 获取所有视频源
func (cfg Config) GetAllVideoSources() map[string]models.VideoSource {
	return cfg.Sources
}

// 获取所有激活的视频源（排除成人内容）
func (cfg Config) GetActiveVideoSources() map[string]models.VideoSource {
	sources := make(map[string]models.VideoSource)
	for key, source := range cfg.Sources {
		if !source.Adult {
			sources[key] = source
		}
	}
	return sources
}
