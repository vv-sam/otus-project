package configuration

const (
	CONFIGURATION_TYPE_FACTORIO = "factorio"
)

// Конфигурация для запуска сервера Factorio
type Factorio struct {
	Server ServerSetting  `json:"server_settings"`
	Map    MapSettings    `json:"map_settings"`
	MapGen MapGenSettings `json:"map_gen_settings"`

	baseConfig
}

type ServerSetting struct {
	Name                                 string   `json:"name"`
	Description                          string   `json:"description"`
	Tags                                 []string `json:"tags"`
	MaxPlayers                           int      `json:"max_players"`
	Username                             string   `json:"username"`
	password                             string
	Token                                string `json:"token"`
	gamePassword                         string
	RequireUserVerification              bool   `json:"require_user_verification"`
	MaxUploadInKilobytesPerSecond        int    `json:"max_upload_in_kilobytes_per_second"`
	MaxUploadSlots                       int    `json:"max_upload_slots"`
	MinimumLatencyInTicks                int    `json:"minimum_latency_in_ticks"`
	MaxHeartbeatsPerSecond               int    `json:"max_heartbeats_per_second"`
	IgnorePlayerLimitForReturningPlayers bool   `json:"ignore_player_limit_for_returning_players"`
	AllowCommands                        string `json:"allow_commands"`
	AutosaveInterval                     int    `json:"autosave_interval"`
	AutosaveSlots                        int    `json:"autosave_slots"`
	AfkAutokickInterval                  int    `json:"afk_autokick_interval"`
	AutoPause                            bool   `json:"auto_pause"`
	AutoPauseWhenPlayersConnect          bool   `json:"auto_pause_when_players_connect"`
	OnlyAdminsCanPauseTheGame            bool   `json:"only_admins_can_pause_the_game"`
	AutosaveOnlyOnServer                 bool   `json:"autosave_only_on_server"`
	NonBlockingSaving                    bool   `json:"non_blocking_saving"`
	MinimumSegmentSize                   int    `json:"minimum_segment_size"`
	MinimumSegmentSizePeerCount          int    `json:"minimum_segment_size_peer_count"`
	MaximumSegmentSize                   int    `json:"maximum_segment_size"`
	MaximumSegmentSizePeerCount          int    `json:"maximum_segment_size_peer_count"`
}

func (s *Factorio) SetPassword(pass string) {
	s.Server.password = pass
}

func (s *Factorio) SetGamePassword(pass string) {
	s.Server.gamePassword = pass
}

type MapSettings struct {
	DifficultySettings difficulty `json:"difficulty_settings"`
	Pollution          pollution  `json:"pollution"`
	EnemyEvolution     evolution  `json:"enemy_evolution"`
	EnemyExpansion     expansion  `json:"enemy_expansion"`
	UnitGroup          unitGroup  `json:"unit_group"`
}

type difficulty struct {
	TechnologyPriceMultiplier float32 `json:"technology_price_multiplier"`
	SpoilTimeModifier         float32 `json:"spoil_time_modifier"`
}

type pollution struct {
	Enabled                                 bool    `json:"enabled"`
	DiffusionRatio                          float64 `json:"diffusion_ratio"`
	MinToDiffuse                            int     `json:"min_to_diffuse"`
	Ageing                                  int     `json:"ageing"`
	ExpectedMaxPerChunk                     int     `json:"expected_max_per_chunk"`
	MinToShowPerChunk                       int     `json:"min_to_show_per_chunk"`
	MinPollutionToDamageTrees               int     `json:"min_pollution_to_damage_trees"`
	PollutionWithMaxForestDamage            int     `json:"pollution_with_max_forest_damage"`
	PollutionPerTreeDamage                  int     `json:"pollution_per_tree_damage"`
	PollutionRestoredPerTreeDamage          int     `json:"pollution_restored_per_tree_damage"`
	MaxPollutionToRestoreTrees              int     `json:"max_pollution_to_restore_trees"`
	EnemyAttackPollutionConsumptionModifier int     `json:"enemy_attack_pollution_consumption_modifier"`
}

type evolution struct {
	Enabled         bool    `json:"enabled"`
	TimeFactor      float64 `json:"time_factor"`
	DestroyFactor   float64 `json:"destroy_factor"`
	PollutionFactor float64 `json:"pollution_factor"`
}

type expansion struct {
	Enabled                          bool    `json:"enabled"`
	MaxExpansionDistance             int     `json:"max_expansion_distance"`
	FriendlyBaseInfluenceRadius      int     `json:"friendly_base_influence_radius"`
	EnemyBuildingInfluenceRadius     int     `json:"enemy_building_influence_radius"`
	BuildingCoefficient              float64 `json:"building_coefficient"`
	OtherBaseCoefficient             float64 `json:"other_base_coefficient"`
	NeighbouringChunkCoefficient     float64 `json:"neighbouring_chunk_coefficient"`
	NeighbouringBaseChunkCoefficient float64 `json:"neighbouring_base_chunk_coefficient"`
	MaxCollidingTilesCoefficient     float64 `json:"max_colliding_tiles_coefficient"`
	SettlerGroupMinSize              int     `json:"settler_group_min_size"`
	SettlerGroupMaxSize              int     `json:"settler_group_max_size"`
	MinExpansionCooldown             int     `json:"min_expansion_cooldown"`
	MaxExpansionCooldown             int     `json:"max_expansion_cooldown"`
}

type unitGroup struct {
	MinGroupGatheringTime          int     `json:"min_group_gathering_time"`
	MaxGroupGatheringTime          int     `json:"max_group_gathering_time"`
	MaxWaitTimeForLateMembers      int     `json:"max_wait_time_for_late_members"`
	MaxGroupRadius                 float64 `json:"max_group_radius"`
	MinGroupRadius                 float64 `json:"min_group_radius"`
	MaxMemberSpeedupWhenBehind     float64 `json:"max_member_speedup_when_behind"`
	MaxMemberSlowdownWhenAhead     float64 `json:"max_member_slowdown_when_ahead"`
	MaxGroupSlowdownFactor         float64 `json:"max_group_slowdown_factor"`
	MaxGroupMemberFallbackFactor   int     `json:"max_group_member_fallback_factor"`
	MemberDisownDistance           int     `json:"member_disown_distance"`
	TickToleranceWhenMemberArrives int     `json:"tick_tolerance_when_member_arrives"`
	MaxGatheringUnitGroups         int     `json:"max_gathering_unit_groups"`
	MaxUnitGroupSize               int     `json:"max_unit_group_size"`
}

type MapGenSettings struct {
	Width             int                 `json:"width"`
	Height            int                 `json:"height"`
	StartingArea      int                 `json:"starting_area"`
	PeacefulMode      bool                `json:"peaceful_mode"`
	AutoplaceControls map[string]resource `json:"autoplace_controls"`
	CliffSettings     cliff               `json:"cliff_settings"`
	Seed              int                 `json:"seed"`
}

type resource struct {
	Frequency float32 `json:"frequency"`
	Size      float32 `json:"size"`
	Richness  float32 `json:"richness"`
}

type cliff struct {
	Name                   string `json:"name"`
	CliffElevation0        int    `json:"cliff_elevation_0"`
	CliffElevationInterval int    `json:"cliff_elevation_interval"`
	Richness               int    `json:"richness"`
}
