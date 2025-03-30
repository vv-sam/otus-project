package configuration

const (
	CONFIGURATION_TYPE_MINECRAFT = "minecraft"

	GAMEMODE_SURVIVAL  = "survival"
	GAMEMODE_CREATIVE  = "creative"
	GAMEMODE_ADVENTURE = "adventure"
	GAMEMODE_SPECTATOR = "spectator"
)

type Minecraft struct {
	Seed         string `json:"seed"`
	RconPort     uint16 `json:"rcon_port"`
	Gamemode     string `json:"gamemode"`
	ServerName   string `json:"server_name"`
	MaxPlayers   uint   `json:"max_players"`
	ViewDistance int    `json:"view_distance"`

	baseConfig
}
