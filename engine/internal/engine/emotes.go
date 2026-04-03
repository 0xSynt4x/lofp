package engine

import (
	"fmt"
	"strings"
)

// emoteEntry defines the self/room messages for an emote.
// Placeholders: %N = actor first name, %P = his/her, %O = him/her, %E = he/she
type emoteEntry struct {
	Self       string // what the actor sees (no target)
	Room       string // what the room sees (no target)
	SelfTarget string // what the actor sees (with target) — %T = target name
	RoomTarget string // what the room sees (with target) — %T = target name
}

var emoteTable = map[string]emoteEntry{
	"SMILE":   {Self: "You smile.", Room: "%N smiles.", SelfTarget: "You smile at %T.", RoomTarget: "%N smiles at %T."},
	"BOW":     {Self: "You bow.", Room: "%N bows.", SelfTarget: "You bow to %T.", RoomTarget: "%N bows to %T."},
	"CURTSEY": {Self: "You curtsey.", Room: "%N curtseys.", SelfTarget: "You curtsey to %T.", RoomTarget: "%N curtseys to %T."},
	"WAVE":    {Self: "You wave.", Room: "%N waves.", SelfTarget: "You wave to %T.", RoomTarget: "%N waves to %T."},
	"NOD":     {Self: "You nod.", Room: "%N nods.", SelfTarget: "You nod to %T.", RoomTarget: "%N nods to %T."},
	"LAUGH":   {Self: "You laugh.", Room: "%N laughs.", SelfTarget: "You laugh at %T.", RoomTarget: "%N laughs at %T."},
	"CHUCKLE": {Self: "You chuckle.", Room: "%N chuckles.", SelfTarget: "You chuckle at %T.", RoomTarget: "%N chuckles at %T."},
	"GRIN":    {Self: "You grin.", Room: "%N grins.", SelfTarget: "You grin at %T.", RoomTarget: "%N grins at %T."},
	"FROWN":   {Self: "You frown.", Room: "%N frowns.", SelfTarget: "You frown at %T.", RoomTarget: "%N frowns at %T."},
	"SIGH":    {Self: "You sigh.", Room: "%N sighs.", SelfTarget: "You sigh at %T.", RoomTarget: "%N sighs at %T."},
	"SHRUG":   {Self: "You shrug.", Room: "%N shrugs.", SelfTarget: "You shrug at %T.", RoomTarget: "%N shrugs at %T."},
	"WINK":    {Self: "You wink.", Room: "%N winks.", SelfTarget: "You wink at %T.", RoomTarget: "%N winks at %T."},
	"CRY":     {Self: "You cry.", Room: "%N cries.", SelfTarget: "You cry on %T's shoulder.", RoomTarget: "%N cries on %T's shoulder."},
	"DANCE":   {Self: "You dance.", Room: "%N dances.", SelfTarget: "You dance with %T.", RoomTarget: "%N dances with %T."},
	"HUG":     {Self: "You hug yourself.", Room: "%N hugs themselves.", SelfTarget: "You hug %T.", RoomTarget: "%N hugs %T."},
	"KISS":    {Self: "You blow a kiss.", Room: "%N blows a kiss.", SelfTarget: "You kiss %T.", RoomTarget: "%N kisses %T."},
	"POKE":    {Self: "You poke yourself.", Room: "%N pokes themselves.", SelfTarget: "You poke %T.", RoomTarget: "%N pokes %T."},
	"TICKLE":  {Self: "You tickle yourself.", Room: "%N tickles themselves.", SelfTarget: "You tickle %T.", RoomTarget: "%N tickles %T."},
	"SLAP":    {Self: "You slap yourself.", Room: "%N slaps themselves.", SelfTarget: "You slap %T.", RoomTarget: "%N slaps %T."},
	"HOWL":    {Self: "You howl.", Room: "%N howls.", SelfTarget: "You howl at %T.", RoomTarget: "%N howls at %T."},
	"SING":    {Self: "You sing.", Room: "%N sings.", SelfTarget: "You sing to %T.", RoomTarget: "%N sings to %T."},
	"PACE":    {Self: "You pace back and forth.", Room: "%N paces back and forth."},
	"FIDGET":  {Self: "You fidget.", Room: "%N fidgets."},
	"SHIVER":  {Self: "You shiver.", Room: "%N shivers."},
	"SNORT":   {Self: "You snort.", Room: "%N snorts.", SelfTarget: "You snort at %T.", RoomTarget: "%N snorts at %T."},
	"GROAN":   {Self: "You groan.", Room: "%N groans."},
	"MUMBLE":  {Self: "You mumble something.", Room: "%N mumbles something."},
	"BABBLE":  {Self: "You babble.", Room: "%N babbles."},
	"BEAM":    {Self: "You beam.", Room: "%N beams.", SelfTarget: "You beam at %T.", RoomTarget: "%N beams at %T."},
	"SWOON":   {Self: "You swoon.", Room: "%N swoons."},
	"TOAST":   {Self: "You raise your glass in a toast.", Room: "%N raises a toast.", SelfTarget: "You raise a toast to %T.", RoomTarget: "%N raises a toast to %T."},
	"SHUDDER": {Self: "You shudder.", Room: "%N shudders."},
	"POINT":   {Self: "You point.", Room: "%N points.", SelfTarget: "You point at %T.", RoomTarget: "%N points at %T."},
	"KICK":    {Self: "You kick at the ground.", Room: "%N kicks at the ground.", SelfTarget: "You kick %T.", RoomTarget: "%N kicks %T."},
	"KNOCK":   {Self: "You knock.", Room: "%N knocks.", SelfTarget: "You knock on %T.", RoomTarget: "%N knocks on %T."},
	"TOUCH":   {Self: "You touch yourself.", Room: "%N touches themselves.", SelfTarget: "You touch %T.", RoomTarget: "%N touches %T."},
	"RUB":     {Self: "You rub your hands together.", Room: "%N rubs %P hands together.", SelfTarget: "You rub %T.", RoomTarget: "%N rubs %T."},
	"PET":     {Self: "You pet yourself.", Room: "%N pets themselves.", SelfTarget: "You pet %T.", RoomTarget: "%N pets %T."},
	"PUNCH":   {Self: "You punch the air.", Room: "%N punches the air.", SelfTarget: "You punch %T.", RoomTarget: "%N punches %T."},
	"SPIT":    {Self: "You spit.", Room: "%N spits.", SelfTarget: "You spit at %T.", RoomTarget: "%N spits at %T."},
	"GAZE":    {Self: "You gaze about.", Room: "%N gazes about.", SelfTarget: "You gaze at %T.", RoomTarget: "%N gazes at %T."},
	"GLARE":   {Self: "You glare.", Room: "%N glares.", SelfTarget: "You glare at %T.", RoomTarget: "%N glares at %T."},
	"SCOWL":   {Self: "You scowl.", Room: "%N scowls.", SelfTarget: "You scowl at %T.", RoomTarget: "%N scowls at %T."},
	"COMFORT": {Self: "You comfort yourself.", Room: "%N comforts themselves.", SelfTarget: "You comfort %T.", RoomTarget: "%N comforts %T."},
	"RECITE":  {Self: "You recite.", Room: "%N recites."},
	"YAWN":    {Self: "You yawn.", Room: "%N yawns.", SelfTarget: "You yawn at %T.", RoomTarget: "%N yawns at %T."},

	// New emotes from chat log analysis
	"BLINK":   {Self: "You blink.", Room: "%N blinks.", SelfTarget: "You blink at %T.", RoomTarget: "%N blinks at %T."},
	"BLUSH":   {Self: "You blush.", Room: "%N blushes."},
	"CRINGE":  {Self: "You cringe.", Room: "%N cringes."},
	"CUDDLE":  {Self: "You cuddle up.", Room: "%N cuddles up.", SelfTarget: "You cuddle up to %T.", RoomTarget: "%N cuddles up to %T."},
	"COUGH":   {Self: "You cough.", Room: "%N coughs."},
	"FURROW":  {Self: "You furrow your brow.", Room: "%N furrows %P brow."},
	"GASP":    {Self: "You gasp.", Room: "%N gasps."},
	"GIGGLE":  {Self: "You giggle.", Room: "%N giggles.", SelfTarget: "You giggle at %T.", RoomTarget: "%N giggles at %T."},
	"GRIMACE": {Self: "You grimace.", Room: "%N grimaces."},
	"GROWL":   {Self: "You growl.", Room: "%N growls.", SelfTarget: "You growl at %T.", RoomTarget: "%N growls at %T."},
	"GULP":    {Self: "You gulp.", Room: "%N gulps."},
	"JUMP":    {Self: "You jump up and down.", Room: "%N jumps up and down."},
	"LEAN":    {Self: "You lean back.", Room: "%N leans back.", SelfTarget: "You lean on %T.", RoomTarget: "%N leans on %T."},
	"NUZZLE":  {Self: "You nuzzle.", Room: "%N nuzzles.", SelfTarget: "You nuzzle %T affectionately.", RoomTarget: "%N nuzzles %T affectionately."},
	"PANT":    {Self: "You pant.", Room: "%N pants."},
	"PONDER":  {Self: "You ponder for a moment.", Room: "%N ponders."},
	"POUT":    {Self: "You pout.", Room: "%N pouts."},
	"ROLL":    {Self: "You roll your eyes.", Room: "%N rolls %P eyes.", SelfTarget: "You roll your eyes at %T.", RoomTarget: "%N rolls %P eyes at %T."},
	"SCREAM":  {Self: "You scream!", Room: "%N screams!"},
	"SMIRK":   {Self: "You smirk.", Room: "%N smirks.", SelfTarget: "You smirk at %T.", RoomTarget: "%N smirks at %T."},
	"SNICKER": {Self: "You snicker.", Room: "%N snickers."},
	"SALUTE":  {Self: "You salute.", Room: "%N salutes.", SelfTarget: "You salute %T.", RoomTarget: "%N salutes %T."},
	"STRETCH": {Self: "You stretch your arms lazily.", Room: "%N stretches %P arms lazily."},
	"TAP":     {Self: "You tap your foot.", Room: "%N taps %P foot.", SelfTarget: "You tap %T on the shoulder.", RoomTarget: "%N taps %T on the shoulder."},
	"TWIRL":   {Self: "You twirl around.", Room: "%N twirls around."},
	"WINCE":   {Self: "You wince.", Room: "%N winces."},
	"WHISTLE": {Self: "You whistle innocently.", Room: "%N whistles innocently."},
	"MUTTER":  {Self: "You mutter something under your breath.", Room: "%N mutters something you can't quite make out."},
	"CARESS":  {Self: "You caress yourself.", Room: "%N caresses themselves.", SelfTarget: "You caress %T.", RoomTarget: "%N caresses %T."},
	"NUDGE":   {Self: "You nudge.", Room: "%N nudges.", SelfTarget: "You nudge %T.", RoomTarget: "%N nudges %T."},
	"ARCH":    {Self: "You arch an eyebrow.", Room: "%N arches %P eyebrow.", SelfTarget: "You arch an eyebrow at %T.", RoomTarget: "%N arches %P eyebrow at %T."},
	"RAISE":   {Self: "You raise an eyebrow.", Room: "%N raises an eyebrow.", SelfTarget: "You raise an eyebrow towards %T.", RoomTarget: "%N raises an eyebrow towards %T."},
	"HEAD":    {Self: "You shake your head.", Room: "%N shakes %P head back and forth.", SelfTarget: "You shake your head at %T.", RoomTarget: "%N shakes %P head at %T."},
	"SCRATCH": {Self: "You scratch your head.", Room: "%N scratches %P head."},
	"CLAP":    {Self: "You clap.", Room: "%N claps.", SelfTarget: "You clap for %T.", RoomTarget: "%N claps for %T."},
	"SNIFF":   {Self: "You sniff.", Room: "%N sniffs.", SelfTarget: "You sniff %T.", RoomTarget: "%N sniffs %T."},
	"LISTEN":  {Self: "You listen carefully.", Room: "%N listens carefully.", SelfTarget: "You listen carefully to %T.", RoomTarget: "%N listens carefully to %T."},
}

// expandEmote replaces %N, %P, %O, %E, %T placeholders in emote strings.
func expandEmote(template string, actor *Player, targetName string) string {
	result := template
	for i := 0; i < len(result); i++ {
		if result[i] == '%' && i+1 < len(result) {
			var replacement string
			switch result[i+1] {
			case 'N':
				replacement = actor.FirstName
			case 'P':
				replacement = actor.Possessive()
			case 'O':
				replacement = actor.Objective()
			case 'E':
				replacement = actor.Pronoun()
			case 'T':
				replacement = targetName
			default:
				continue
			}
			result = result[:i] + replacement + result[i+2:]
			i += len(replacement) - 1
		}
	}
	return result
}

// processEmote handles emote commands using the emote table.
func (e *GameEngine) processEmote(player *Player, verb string, args []string) *CommandResult {
	entry, ok := emoteTable[verb]
	if !ok {
		// Fallback generic
		v := strings.ToLower(verb)
		return &CommandResult{
			Messages:      []string{fmt.Sprintf("You %s.", v)},
			RoomBroadcast: []string{fmt.Sprintf("%s %ss.", player.FirstName, v)},
		}
	}

	if len(args) > 0 {
		targetName := strings.ToLower(strings.Join(args, " "))

		// Check for "me"/"myself"
		if targetName == "me" || targetName == "myself" || targetName == "self" {
			selfMsg := expandEmote(entry.Self, player, player.FirstName)
			roomMsg := expandEmote(entry.Room, player, player.FirstName)
			return &CommandResult{Messages: []string{selfMsg}, RoomBroadcast: []string{roomMsg}}
		}

		// If targeted emote templates exist
		if entry.SelfTarget != "" && entry.RoomTarget != "" {
			// Try to resolve as a player in the room
			found := e.findPlayerInRoom(player, targetName)
			if found != nil {
				displayTarget := found.FirstName
				selfMsg := expandEmote(entry.SelfTarget, player, displayTarget)
				roomMsg := expandEmote(entry.RoomTarget, player, displayTarget)
				targetMsg := expandEmote(entry.RoomTarget, player, "you")
				return &CommandResult{
					Messages:      []string{selfMsg},
					RoomBroadcast: []string{roomMsg},
					TargetName:    found.FirstName,
					TargetMsg:     []string{targetMsg},
				}
			}

			// Try to resolve as a room item
			room := e.rooms[player.RoomNumber]
			if room != nil {
				for _, ri := range room.Items {
					itemDef := e.items[ri.Archetype]
					if itemDef == nil {
						continue
					}
					name := e.getItemNounName(itemDef)
					if matchesTarget(name, targetName, e.getAdjName(ri.Adj1)) {
						displayTarget := e.formatItemName(itemDef, ri.Adj1, ri.Adj2, ri.Adj3)
						selfMsg := expandEmote(entry.SelfTarget, player, displayTarget)
						roomMsg := expandEmote(entry.RoomTarget, player, displayTarget)
						return &CommandResult{Messages: []string{selfMsg}, RoomBroadcast: []string{roomMsg}}
					}
				}
			}

			// Also check player's items (inventory + worn + wielded)
			allItems := make([]InventoryItem, 0)
			allItems = append(allItems, player.Inventory...)
			allItems = append(allItems, player.Worn...)
			if player.Wielded != nil {
				allItems = append(allItems, *player.Wielded)
			}
			for _, ii := range allItems {
				itemDef := e.items[ii.Archetype]
				if itemDef == nil {
					continue
				}
				name := e.getItemNounName(itemDef)
				if matchesTarget(name, targetName, e.getAdjName(ii.Adj1)) || matchesTarget(name, targetName, e.getAdjName(ii.Adj3)) {
					displayTarget := e.formatItemName(itemDef, ii.Adj1, ii.Adj2, ii.Adj3)
					selfMsg := expandEmote(entry.SelfTarget, player, displayTarget)
					roomMsg := expandEmote(entry.RoomTarget, player, displayTarget)
					return &CommandResult{Messages: []string{selfMsg}, RoomBroadcast: []string{roomMsg}}
				}
			}

			// Nothing matched
			return &CommandResult{Messages: []string{fmt.Sprintf("You don't see '%s' here.", targetName)}}
		}
	}

	selfMsg := expandEmote(entry.Self, player, "")
	roomMsg := expandEmote(entry.Room, player, "")
	return &CommandResult{Messages: []string{selfMsg}, RoomBroadcast: []string{roomMsg}}
}
