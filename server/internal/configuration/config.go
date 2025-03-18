package configuration

import (
	"fmt"
	"strings"
)

type Config struct {
	Discord        DiscordConfig `json:"discord"`
	changedHandler func(conf *Config) error
}

// Discord section
type DiscordConfig struct {
	// Roles having unlimited access to service. Optional if AdminProfileId is set
	adminRoles []string `json:"admin_roles"`

	// Id of user who has unlimited access to service. Optional if AdminRoles are set
	adminProfileId string `json:"admin_profile_id"`

	// Token for bot access
	BotToken string
}

func (conf *DiscordConfig) GetAdminRoles() []string {
	c := make([]string, len(conf.adminRoles))
	copy(c, conf.adminRoles)

	return c
}

func (conf *DiscordConfig) GetAdminProfileId() string {
	return conf.adminProfileId
}

// Из пальца высосанная валидация для демонстрации инкапсуляции unexported полей
func (conf *Config) SetAdminRoles(roles []string) error {
	for _, role := range roles {
		r := strings.TrimSpace(role)
		if r == "" {
			return fmt.Errorf("incorrect role name: %s", role)
		}
	}

	conf.Discord.adminRoles = make([]string, len(roles))
	copy(conf.Discord.adminRoles, roles)

	conf.changedHandler(conf)

	return nil
}

func (conf *Config) HandleConfigChanged(f func(conf *Config) error) {
	conf.changedHandler = f
}
