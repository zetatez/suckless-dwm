package psl

import (
	"fmt"
	"net"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"

	"assistant/pkg/xlog"

	"github.com/spf13/viper"
)

var (
	onceConfig sync.Once
	config     *Config
)

func GetConfig() *Config { return config }

func InitConfig() error {
	var err error
	onceConfig.Do(func() {
		config, err = loadConfig()
	})
	return err
}

func loadConfig() (*Config, error) {
	home, _ := os.UserHomeDir()
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.SetEnvPrefix("APP")
	v.AutomaticEnv()
	v.AddConfigPath(".")
	v.AddConfigPath(path.Join(home, ".config/assistant"))

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read in config: %v", err)
	}
	var c Config
	if err := v.Unmarshal(&c); err != nil {
		return nil, fmt.Errorf("failed unmarshal config: %v", err)
	}
	c.resolveEnv()
	c.applyDefaults()
	if err := c.expandPaths(); err != nil {
		return nil, fmt.Errorf("expand paths: %w", err)
	}
	if err := c.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}
	return &c, nil
}

type Config struct {
	App         AppConfig         `mapstructure:"app"`
	Auth        AuthConfig        `mapstructure:"auth"`
	Log         xlog.LogConfig    `mapstructure:"log"`
	LLM         LLMConfig         `mapstructure:"llm"`
	Svc         SvcConfig         `mapstructure:"svc"`
	Channels    ChannelsConfig    `mapstructure:"channels"`
	Background  BackgroundConfig  `mapstructure:"background"`
	FileBrowser FileBrowserConfig `mapstructure:"filebrowser"`
}

type FileBrowserConfig struct {
	Root   string   `mapstructure:"root"`
	Allow  []string `mapstructure:"allow"`
	Deny   []string `mapstructure:"deny"`
	Public []string `mapstructure:"public"`
}

type AppConfig struct {
	Name      string `mapstructure:"name"`
	Host      string `mapstructure:"host"`
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
	SnipDir                string `mapstructure:"dir_snip"`
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
	Name      string `mapstructure:"name"`
	Command   string `mapstructure:"command"`
	Precursor string `mapstructure:"precursor"`
}

func (c *Config) applyDefaults() {
	if c.App.Host == "" {
		c.App.Host = "127.0.0.1"
	}
	if c.App.Port == 0 {
		c.App.Port = 4321
	}
	if c.App.Interface == "" {
		c.App.Interface = detectDefaultInterface()
	}
	if c.LLM.Timeout == 0 {
		c.LLM.Timeout = 60
	}
	if c.LLM.MaxTokens == 0 {
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
	if c.Svc.SnipDir == "" {
		c.Svc.SnipDir = "~/git/obsidian/.snippets"
	}
	if c.FileBrowser.Root == "" {
		home, _ := os.UserHomeDir()
		c.FileBrowser.Root = home
	}
	if len(c.FileBrowser.Deny) == 0 {
		c.FileBrowser.Deny = []string{".ssh", ".gnupg", ".config/assistant"}
	}
	if len(c.Background.Procs) == 0 {
		home, _ := os.UserHomeDir()
		c.Background.Procs = []BackgroundProc{
			{Name: "dwmblocks", Command: "dwmblocks"},
			{Name: "picom", Command: "picom --config " + home + "/.config/picom/picom.conf"},
			{Name: "dunst", Command: "dunst"},
		}
	}
}

func (c *Config) expandPaths() error {
	home, _ := os.UserHomeDir()
	for _, p := range []*string{
		&c.Log.Filename,
		&c.Svc.DirWallpaper,
		&c.Svc.WorkingLogbookDir,
		&c.Svc.SSHSecretFile,
		&c.Svc.KeyboardBrightnessPath,
		&c.Svc.SnipDir,
		&c.FileBrowser.Root,
	} {
		*p = expandHomePath(*p, home)
	}
	return nil
}

func expandHomePath(p, home string) string {
	if p == "" {
		return p
	}
	if p == "~" {
		return home
	}
	if strings.HasPrefix(p, "~/") {
		return path.Join(home, p[2:])
	}
	return p
}

func (c *Config) Validate() error {
	if c.App.Name == "" {
		return fmt.Errorf("app.name is required")
	}
	if c.App.Host != "" && net.ParseIP(c.App.Host) == nil && c.App.Host != "0.0.0.0" {
		return fmt.Errorf("app.host must be a valid IP address: %s", c.App.Host)
	}
	if c.LLM.Provider == "" {
		return fmt.Errorf("llm.provider is required")
	}
	return nil
}

func (c *Config) resolveEnv() {
	envPH := regexp.MustCompile(`\$\{(\w+)\}`)
	expand := func(p *string) { *p = expandEnvPH(*p, envPH) }
	expand(&c.LLM.APIKey)
	expand(&c.Channels.Feishu.AppID)
	expand(&c.Channels.Feishu.AppSecret)
	expand(&c.Channels.Feishu.ChatID)
}

func expandEnvPH(s string, envPH *regexp.Regexp) string {
	if s == "" {
		return s
	}
	return envPH.ReplaceAllStringFunc(s, func(m string) string {
		groups := envPH.FindStringSubmatch(m)
		if len(groups) < 2 {
			return m
		}
		return os.Getenv(groups[1])
	})
}

func detectDefaultInterface() string {
	data, err := os.ReadFile("/proc/net/route")
	if err != nil {
		return ""
	}
	lines := strings.Split(string(data), "\n")
	bestIface := ""
	bestMetric := uint32(^uint32(0))
	for _, line := range lines[1:] {
		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}
		if destHex := fields[1]; destHex != "00000000" || len(destHex) != 8 {
			continue
		}
		var metric uint32
		if len(fields) >= 7 {
			fmt.Sscanf(fields[6], "%d", &metric)
		}
		if metric < bestMetric {
			bestMetric = metric
			bestIface = fields[0]
		}
	}
	return bestIface
}

func (c *LLMConfig) GetAPIKey() string { return c.APIKey }
