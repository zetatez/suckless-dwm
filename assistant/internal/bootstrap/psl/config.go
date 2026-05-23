package psl

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"

	"assistant/pkg/xlog"

	"github.com/spf13/viper"
)

var (
	config     *Config
	onceConfig sync.Once
)

func GetConfig() *Config { return config }

func InitConfig() error {
	var initErr error
	onceConfig.Do(func() {
		var err error
		config, err = LoadConfig()
		if err != nil {
			initErr = fmt.Errorf("load config failed: %w", err)
			return
		}
	})
	return initErr
}

type Config struct {
	App        AppConfig        `mapstructure:"app"`
	Auth       AuthConfig       `mapstructure:"auth"`
	Log        xlog.LogConfig   `mapstructure:"log"`
	LLM        LLMConfig        `mapstructure:"llm"`
	Svc        SvcConfig        `mapstructure:"svc"`
	Channels   ChannelsConfig   `mapstructure:"channels"`
	Background BackgroundConfig `mapstructure:"background"`
}

type AppConfig struct {
	Name      string `mapstructure:"name"`
	Port      int    `mapstructure:"port"`
	Interface string `mapstructure:"interface"`
}

type AuthConfig struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type LLMConfig struct {
	Provider    string            `mapstructure:"provider"`
	APIKey      string            `mapstructure:"api_key"`
	BaseURL     string            `mapstructure:"base_url"`
	Model       string            `mapstructure:"model"`
	Extra       map[string]string `mapstructure:"extra"`
	Timeout     int               `mapstructure:"timeout"`
	MaxTokens   int               `mapstructure:"max_tokens"`
	Temperature float32           `mapstructure:"temperature"`
}

type SvcConfig struct {
	ProxyServer            string `mapstructure:"vpn_proxy"`
	PrimaryMonitor         string `mapstructure:"default_monitor"`
	DirWallpaper           string `mapstructure:"dir_wallpaper"`
	WorkingLogbookDir      string `mapstructure:"dir_working_logbook"`
	KeyboardBrightnessPath string `mapstructure:"path_keyboard_brightness"`
	SSHSecretFile          string `mapstructure:"path_ssh_secret"`
	DefaultTerminal        string `mapstructure:"terminal_default"`
}

type ChannelsConfig struct {
	Feishu FeishuConfig `mapstructure:"feishu"`
}

type FeishuConfig struct {
	AppID     string `mapstructure:"app_id"`
	AppSecret string `mapstructure:"app_secret"`
	ChatID    string `mapstructure:"chat_id"`
}

type BackgroundConfig struct {
	Enabled bool             `mapstructure:"enabled"`
	Procs   []BackgroundProc `mapstructure:"procs"`
}

type BackgroundProc struct {
	Name    string `mapstructure:"name"`
	Command string `mapstructure:"command"`
}

func LoadConfig() (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.SetEnvPrefix("APP")
	v.AutomaticEnv()
	v.AddConfigPath(".")
	if home, err := os.UserHomeDir(); err == nil {
		v.AddConfigPath(path.Join(home, ".config/assistant"))
	}
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read in config: %v", err)
	}
	var cfg *Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed unmarshal config: %v", err)
	}
	cfg.resolveEnv()
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}
	return cfg, nil
}

func (c *Config) resolveEnv() {
	resolveEnv(&c.LLM.APIKey)
	resolveEnv(&c.Channels.Feishu.AppID)
	resolveEnv(&c.Channels.Feishu.AppSecret)
	resolveEnv(&c.Channels.Feishu.ChatID)
	if strings.HasPrefix(c.Log.Filename, "~/") {
		home, _ := os.UserHomeDir()
		c.Log.Filename = path.Join(home, c.Log.Filename[2:])
	}
}

func (c *Config) Validate() error {
	if c.App.Port <= 0 {
		c.App.Port = 4321
	}
	if c.App.Name == "" {
		return fmt.Errorf("app.name is required")
	}
	if c.App.Interface == "" {
		c.App.Interface = "wlan0"
	}
	if c.LLM.Provider == "" {
		return fmt.Errorf("llm.provider is required")
	}
	if c.LLM.Timeout <= 0 {
		c.LLM.Timeout = 60
	}
	if c.LLM.MaxTokens <= 0 {
		c.LLM.MaxTokens = 4096
	}
	if c.Svc.PrimaryMonitor == "" {
		c.Svc.PrimaryMonitor = "eDP-1"
	}
	if c.Svc.DirWallpaper == "" {
		c.Svc.DirWallpaper = "~/Pictures/wallpapers"
	}
	if c.Svc.WorkingLogbookDir == "" {
		c.Svc.WorkingLogbookDir = "~/git/working/logbook"
	}
	if c.Svc.KeyboardBrightnessPath == "" {
		c.Svc.KeyboardBrightnessPath = "/sys/class/leds/tpacpi::kbd_backlight/brightness"
	}
	if c.Svc.DefaultTerminal == "" {
		c.Svc.DefaultTerminal = "st"
	}
	if len(c.Background.Procs) == 0 {
		c.Background.Procs = []BackgroundProc{
			{Name: "dwmblocks", Command: "dwmblocks"},
			{Name: "picom", Command: "picom --config " + os.Getenv("HOME") + "/.config/picom/picom.conf"},
			{Name: "dunst", Command: "dunst"},
		}
	}
	return nil
}

var envPlaceholder = regexp.MustCompile(`\$\{(\w+)\}`)

func resolveEnv(val *string) {
	if val == nil || *val == "" {
		return
	}
	matches := envPlaceholder.FindStringSubmatch(*val)
	if len(matches) >= 2 {
		*val = os.Getenv(matches[1])
	}
}

func ResolveEnvPlaceholderStr(val string) string {
	if val == "" {
		return ""
	}
	matches := envPlaceholder.FindStringSubmatch(val)
	if len(matches) >= 2 {
		return os.Getenv(matches[1])
	}
	return val
}

func (c *LLMConfig) GetAPIKey() string {
	return ResolveEnvPlaceholderStr(c.APIKey)
}
