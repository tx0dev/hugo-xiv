package types

// Config represents the configuration structure
type Config struct {
	RateLimit     int           `json:"rateLimit"`
	XIVApiUrl     string        `json:"xivApiUrl"`
	Currencies    string        `json:"currencies"`
	Icons         string        `json:"icons"`
	ScssVariables string        `json:"scssVariables"`
	Sources       SourcesConfig `json:"sources"`
	// xivapi        string        `json:"_xivapiv2"`
}

type SourcesConfig struct {
	Maps       string `json:"maps"`
	Places     string `json:"places"`
	Aetherytes string `json:"aetherytes"`
}

// Currencies represents the currencies.json structure
type Currencies struct {
	Currency  map[string]CurrencyItem   `json:"currency"`
	Tomestone map[string]TomestoneItem  `json:"tomestone"`
	Scrip     map[string]ScripItem      `json:"scrip"`
	Society   map[string]BeastTribeItem `json:"society"`
}

type CurrencyItem struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type TomestoneItem struct {
	Name     string  `json:"name"`
	ID       string  `json:"id"`
	Added    string  `json:"added"`
	Obsolete *string `json:"obsolete"`
}

type ScripItem struct {
	Gatherer string  `json:"gatherer"`
	Crafter  string  `json:"crafter"`
	Added    string  `json:"added"`
	Obsolete *string `json:"obsolete"`
}

type BeastTribeItem struct {
	Name       string `json:"name"`
	SocietyID  string `json:"societyId"`
	CurrencyID string `json:"currencyId"`
}

// Icons represents the icons.json structure
type Icons map[string]map[string]string

// Map represents a map
type Map struct {
	ID             int    `json:"id"`
	Hierarchy      int    `json:"hierarchy"`
	PriorityUI     int    `json:"priority_ui"`
	Image          string `json:"image"`
	OffsetX        int    `json:"offset_x"`
	OffsetY        int    `json:"offset_y"`
	OffsetZ        *int   `json:"offset_z"` // Pointer to int to handle null
	MapMarkerRange int    `json:"map_marker_range"`
	PlaceNameID    int    `json:"placename_id"`
	PlaceNameSubID int    `json:"placename_sub_id"`
	RegionID       int    `json:"region_id"`
	ZoneID         int    `json:"zone_id"`
	SizeFactor     int    `json:"size_factor"`
	TerritoryID    int    `json:"territory_id"`
	Index          int    `json:"index"`
	Dungeon        bool   `json:"dungeon"`
	Housing        bool   `json:"housing"`

	Aetherytes []Aetheryte `json:"aetherytes"` // Add field for aetherytes within this map
	PlaceName  Place       `json:"name"`       // add the place name
	SubName    Place       `json:"sub_name"`   // Add the sub place name
}

// Place represents a place item
type Place struct {
	EN string `json:"en"`
	JA string `json:"ja"`
	DE string `json:"de"`
	FR string `json:"fr"`
	//ID string `json:"id"` // ID is a string in places.json, handle accordingly
}

// Aetheryte represents an aetheryte item
type Aetheryte struct {
	ID              int     `json:"id"`
	ZoneID          int     `json:"zoneid"`
	MapID           int     `json:"map"` // Renamed from `map` to `MapID` to avoid conflict if needed
	X               float64 `json:"x"`
	Y               float64 `json:"y"`
	Z               float64 `json:"z"`
	Type            int     `json:"type"`
	TerritoryID     int     `json:"territory"`
	NameID          int     `json:"nameid"`
	AethernetCoords struct {
		X int `json:"x"`
		Y int `json:"y"`
	} `json:"aethernetCoords"`
	Name Place `json:"name"` // Add the aetheryte place name
}
