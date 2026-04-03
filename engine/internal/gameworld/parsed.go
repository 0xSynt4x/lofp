package gameworld

// ParsedData holds all data loaded from script files, used to initialize the engine.
type ParsedData struct {
	Rooms        []Room
	Items        []ItemDef
	Monsters     []MonsterDef
	Nouns        []NounDef
	Adjectives   []AdjDef
	MonsterAdjs  []MonsterAdjDef
	Variables    []Variable
	Regions      []Region
	MonsterLists []MonsterList
	StartRoom    int
	BumpRoom     int
}
