// Package config loads runtime configuration from environment variables.
package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

const envPrefix = "GPUVM"

// Config holds all runtime configuration for the service.
type Config struct {
	Server    ServerConfig
	Tumblebug TumblebugConfig
	Template  TemplateConfig
}

// TemplateConfig selects the Infra template used as the provisioning baseline.
//
// The template ships two GPU NodeGroups (NVIDIA L4 and L40S). Provisioning both
// roughly doubles the hourly cost, so only the NodeGroup matching SpecID is kept.
type TemplateConfig struct {
	Namespace    string
	ID           string
	SpecID       string
	NodeUserName string
}

// ServerConfig holds HTTP server settings.
type ServerConfig struct {
	Port int
}

// TumblebugConfig holds CB-Tumblebug connection settings.
//
// Creating a GPU Infra blocks until provisioning finishes, so Timeout must stay
// generous: cancelling mid-flight aborts provisioning and can leave billable
// orphan resources behind.
type TumblebugConfig struct {
	BaseURL          string
	Username         string
	Password         string
	Timeout          time.Duration
	CredentialHolder string
}

// Load reads configuration from environment variables and validates it.
func Load() (*Config, error) {
	v := viper.New()
	v.SetEnvPrefix(envPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	v.SetDefault("server.port", 8080)
	v.SetDefault("tumblebug.timeout", 20*time.Minute)
	v.SetDefault("tumblebug.credential_holder", "")
	v.SetDefault("template.namespace", "system")
	v.SetDefault("template.id", "infra-aws-gpu-simple")
	v.SetDefault("template.spec_id", "aws+us-west-2+g6.8xlarge")
	v.SetDefault("template.node_user_name", "cb-user")

	// AutomaticEnv only resolves keys that viper already knows about.
	for _, key := range []string{
		"server.port",
		"tumblebug.base_url",
		"tumblebug.username",
		"tumblebug.password",
		"tumblebug.timeout",
		"tumblebug.credential_holder",
		"template.namespace",
		"template.id",
		"template.spec_id",
		"template.node_user_name",
	} {
		if err := v.BindEnv(key); err != nil {
			return nil, fmt.Errorf("failed to bind env for %q: %w", key, err)
		}
	}

	cfg := &Config{
		Server: ServerConfig{
			Port: v.GetInt("server.port"),
		},
		Tumblebug: TumblebugConfig{
			BaseURL:          strings.TrimRight(v.GetString("tumblebug.base_url"), "/"),
			Username:         v.GetString("tumblebug.username"),
			Password:         v.GetString("tumblebug.password"),
			Timeout:          v.GetDuration("tumblebug.timeout"),
			CredentialHolder: v.GetString("tumblebug.credential_holder"),
		},
		Template: TemplateConfig{
			Namespace:    v.GetString("template.namespace"),
			ID:           v.GetString("template.id"),
			SpecID:       v.GetString("template.spec_id"),
			NodeUserName: v.GetString("template.node_user_name"),
		},
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) validate() error {
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}
	if c.Tumblebug.BaseURL == "" {
		return fmt.Errorf("%s_TUMBLEBUG_BASE_URL is required", envPrefix)
	}
	if c.Tumblebug.Timeout <= 0 {
		return fmt.Errorf("invalid tumblebug timeout: %s", c.Tumblebug.Timeout)
	}
	if c.Template.Namespace == "" || c.Template.ID == "" || c.Template.SpecID == "" {
		return fmt.Errorf("template namespace, id and spec id are all required")
	}
	return nil
}
