package engine

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jonradoff/lofp/internal/gameworld"
)

// ScriptContext holds state for script execution within a single trigger.
type ScriptContext struct {
	Player   *Player
	Room     *gameworld.Room
	Engine   *GameEngine
	Messages []string // ECHO PLAYER messages to send to the player
	RoomMsgs []string // ECHO ALL / ECHO OTHERS messages for the room
	GMMsgs   []string // GMMSG messages for gamemasters
	Blocked  bool     // CLEARVERB: block the triggering action
	MoveTo   int      // MOVE: destination room (0 = no move)

	// Item interaction context (set when running IFPREVERB/IFVERB on a room item)
	ItemRef *gameworld.RoomItem // the room item being interacted with
	ItemDef *gameworld.ItemDef  // its archetype definition
}

// RunEntryScripts executes all IFENTRY script blocks for a room.
func (e *GameEngine) RunEntryScripts(player *Player, room *gameworld.Room) *ScriptContext {
	sc := &ScriptContext{
		Player: player,
		Room:   room,
		Engine: e,
	}
	for _, block := range room.Scripts {
		if block.Type == "IFENTRY" {
			sc.execBlock(block)
		}
	}
	return sc
}

// RunPreverbScripts executes IFPREVERB blocks for a specific verb and item ref.
// Returns the script context. Check sc.Blocked to see if the action should be cancelled.
func (e *GameEngine) RunPreverbScripts(player *Player, room *gameworld.Room, verb string, ri *gameworld.RoomItem, def *gameworld.ItemDef) *ScriptContext {
	sc := &ScriptContext{
		Player:  player,
		Room:    room,
		Engine:  e,
		ItemRef: ri,
		ItemDef: def,
	}
	refStr := fmt.Sprintf("%d", ri.Ref)
	verb = strings.ToUpper(verb)

	// Check room-level scripts
	for _, block := range room.Scripts {
		if block.Type == "IFPREVERB" && len(block.Args) >= 2 {
			if strings.ToUpper(block.Args[0]) == verb && block.Args[1] == refStr {
				sc.execBlock(block)
			}
		}
	}

	// Check item-level scripts (on the archetype definition)
	for _, block := range def.Scripts {
		if block.Type == "IFPREVERB" && len(block.Args) >= 1 {
			if strings.ToUpper(block.Args[0]) == verb {
				// Item scripts use -1 as self-reference
				if len(block.Args) < 2 || block.Args[1] == "-1" {
					sc.execBlock(block)
				}
			}
		}
	}

	return sc
}

// RunVerbScripts executes IFVERB blocks for a specific verb and item.
func (e *GameEngine) RunVerbScripts(player *Player, room *gameworld.Room, verb string, ri *gameworld.RoomItem, def *gameworld.ItemDef) *ScriptContext {
	sc := &ScriptContext{
		Player:  player,
		Room:    room,
		Engine:  e,
		ItemRef: ri,
		ItemDef: def,
	}
	verb = strings.ToUpper(verb)

	// Check item-level scripts
	for _, block := range def.Scripts {
		if block.Type == "IFVERB" && len(block.Args) >= 1 {
			if strings.ToUpper(block.Args[0]) == verb {
				if len(block.Args) < 2 || block.Args[1] == "-1" {
					sc.execBlock(block)
				}
			}
		}
	}

	return sc
}

// execBlock executes a script block if its condition is met.
func (sc *ScriptContext) execBlock(block gameworld.ScriptBlock) {
	switch block.Type {
	case "IFENTRY":
		sc.execChildren(block)

	case "IFPREVERB", "IFVERB":
		// Condition already matched by caller; execute body
		sc.execChildren(block)

	case "IFVAR":
		if sc.evalIfVar(block.Args) {
			sc.execChildren(block)
		}

	case "IFITEM":
		if sc.evalIfItem(block.Args) {
			sc.execChildren(block)
		}
	}
}

// execChildren runs the actions and nested blocks within a script block.
func (sc *ScriptContext) execChildren(block gameworld.ScriptBlock) {
	for _, action := range block.Actions {
		sc.execAction(action)
	}
	for _, child := range block.Children {
		sc.execBlock(child)
	}
}

// execAction executes a single script action.
func (sc *ScriptContext) execAction(action gameworld.ScriptAction) {
	switch action.Command {
	case "ECHO":
		sc.doEcho(action.Args)
	case "EQUAL":
		sc.doEqual(action.Args)
	case "NEWITEM":
		sc.doNewItem(action.Args)
	case "GMMSG":
		sc.doGMMsg(action.Args)
	case "CLEARVERB":
		sc.Blocked = true
	case "MOVE":
		sc.doMove(action.Args)
	case "SHOWROOM":
		sc.doShowRoom(action.Args)
	case "PLREVENT", "CONTPLREVENT":
		// Timing/event delay — not yet implemented; ignore silently
	case "AFFECT":
		// Room synchronization — not yet implemented; ignore silently
	case "ADD":
		sc.doAdd(action.Args)
	case "SUB":
		sc.doSub(action.Args)
	case "SETITEMVAL":
		sc.doSetItemVal(action.Args)
	case "REMOVEITEM":
		sc.doRemoveItem(action.Args)
	case "LOCK":
		sc.doItemState(action.Args, "LOCKED")
	case "UNLOCK":
		sc.doItemState(action.Args, "UNLOCKED")
	case "OPEN":
		sc.doItemState(action.Args, "OPEN")
	case "CLOSE":
		sc.doItemState(action.Args, "CLOSED")
	}
}

// doEcho handles ECHO PLAYER, ECHO ALL, ECHO OTHERS.
func (sc *ScriptContext) doEcho(args []string) {
	if len(args) < 2 {
		return
	}
	target := strings.ToUpper(args[0])
	text := strings.Join(args[1:], " ")
	text = sc.expandScriptText(text)

	switch target {
	case "PLAYER":
		sc.Messages = append(sc.Messages, text)
	case "ALL":
		sc.Messages = append(sc.Messages, text)
		sc.RoomMsgs = append(sc.RoomMsgs, text)
	case "OTHERS":
		sc.RoomMsgs = append(sc.RoomMsgs, text)
	}
}

// doEqual handles EQUAL INTNUMn value — sets a variable on the player.
func (sc *ScriptContext) doEqual(args []string) {
	if len(args) < 2 {
		return
	}
	varName := strings.ToUpper(args[0])
	val, err := strconv.Atoi(args[1])
	if err != nil {
		return
	}
	sc.setVar(varName, val)
}

// doAdd handles ADD INTNUMn value — increments a variable.
func (sc *ScriptContext) doAdd(args []string) {
	if len(args) < 2 {
		return
	}
	varName := strings.ToUpper(args[0])
	val, err := strconv.Atoi(args[1])
	if err != nil {
		return
	}
	sc.setVar(varName, sc.getVar(varName)+val)
}

// doSub handles SUB INTNUMn value — decrements a variable.
func (sc *ScriptContext) doSub(args []string) {
	if len(args) < 2 {
		return
	}
	varName := strings.ToUpper(args[0])
	val, err := strconv.Atoi(args[1])
	if err != nil {
		return
	}
	sc.setVar(varName, sc.getVar(varName)-val)
}

// doNewItem handles NEWITEM ref archetype [ADJ1=n] [ADJ2=n] [VAL1=n] ...
// ref -1 means add to player inventory.
func (sc *ScriptContext) doNewItem(args []string) {
	if len(args) < 2 {
		return
	}
	ref, err := strconv.Atoi(args[0])
	if err != nil {
		return
	}
	archetype, err := strconv.Atoi(args[1])
	if err != nil {
		return
	}

	item := InventoryItem{Archetype: archetype}
	for _, arg := range args[2:] {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.ToUpper(parts[0])
		val, err := strconv.Atoi(parts[1])
		if err != nil {
			continue
		}
		switch key {
		case "ADJ1":
			item.Adj1 = val
		case "ADJ2":
			item.Adj2 = val
		case "ADJ3":
			item.Adj3 = val
		case "VAL1":
			item.Val1 = val
		case "VAL2":
			item.Val2 = val
		case "VAL3":
			item.Val3 = val
		case "VAL4":
			item.Val4 = val
		case "VAL5":
			item.Val5 = val
		}
	}

	if ref == -1 {
		sc.Player.Inventory = append(sc.Player.Inventory, item)
	}
}

// doGMMsg broadcasts a message to all online GMs.
func (sc *ScriptContext) doGMMsg(args []string) {
	if len(args) == 0 {
		return
	}
	text := strings.Join(args, " ")
	text = sc.expandScriptText(text)
	sc.GMMsgs = append(sc.GMMsgs, fmt.Sprintf("[GM] %s", text))
}

// doMove handles MOVE <room> or MOVE ITEMVAL2, etc.
func (sc *ScriptContext) doMove(args []string) {
	if len(args) == 0 {
		return
	}
	dest := sc.resolveNumericArg(args[0])
	if dest > 0 {
		sc.MoveTo = dest
	}
}

// doShowRoom handles SHOWROOM <room> or SHOWROOM ITEMVAL2, etc.
func (sc *ScriptContext) doShowRoom(args []string) {
	if len(args) == 0 {
		return
	}
	roomNum := sc.resolveNumericArg(args[0])
	if roomNum > 0 {
		if room := sc.Engine.rooms[roomNum]; room != nil {
			sc.Messages = append(sc.Messages, fmt.Sprintf("[%s]", room.Name))
			if room.Description != "" {
				sc.Messages = append(sc.Messages, descriptionToMessages(room.Description)...)
			}
		}
	}
}

// doSetItemVal handles SETITEMVAL ref valIndex value.
func (sc *ScriptContext) doSetItemVal(args []string) {
	// Not yet fully implemented; needs room item mutation
}

// doRemoveItem handles REMOVEITEM ref — removes item from player or room.
func (sc *ScriptContext) doRemoveItem(args []string) {
	if len(args) == 0 {
		return
	}
	ref, err := strconv.Atoi(args[0])
	if err != nil {
		return
	}
	if ref == -1 && sc.ItemRef != nil {
		// Remove current item from inventory (by archetype match)
		for i, ii := range sc.Player.Inventory {
			if ii.Archetype == sc.ItemRef.Archetype {
				sc.Player.Inventory = append(sc.Player.Inventory[:i], sc.Player.Inventory[i+1:]...)
				break
			}
		}
	}
}

// doItemState sets the state of a room item (LOCK, UNLOCK, OPEN, CLOSE).
func (sc *ScriptContext) doItemState(args []string, state string) {
	if len(args) == 0 {
		return
	}
	ref, err := strconv.Atoi(args[0])
	if err != nil {
		return
	}
	for i := range sc.Room.Items {
		if sc.Room.Items[i].Ref == ref && !sc.Room.Items[i].IsPut {
			sc.Room.Items[i].State = state
			sc.Engine.notifyRoomChange(RoomChange{RoomNumber: sc.Room.Number, Type: "item_state", ItemRef: ref, NewState: state})
			break
		}
	}
}

// evalIfVar evaluates IFVAR conditions like "INTNUM6 = 2" or "ITEMBIT5 = 0".
func (sc *ScriptContext) evalIfVar(args []string) bool {
	if len(args) < 3 {
		return false
	}
	varName := strings.ToUpper(args[0])
	op := args[1]
	expected, err := strconv.Atoi(args[2])
	if err != nil {
		return false
	}

	actual := sc.getVar(varName)

	switch op {
	case "=":
		return actual == expected
	case "!":
		return actual != expected
	case ">":
		return actual > expected
	case "<":
		return actual < expected
	case ">=":
		return actual >= expected
	case "<=":
		return actual <= expected
	}
	return false
}

// evalIfItem evaluates IFITEM conditions like "IFITEM 0 OPEN" or "IFITEM -1 CLOSED".
func (sc *ScriptContext) evalIfItem(args []string) bool {
	if len(args) < 2 {
		return false
	}
	ref, err := strconv.Atoi(args[0])
	if err != nil {
		return false
	}
	expectedState := strings.ToUpper(args[1])

	var ri *gameworld.RoomItem
	if ref == -1 && sc.ItemRef != nil {
		ri = sc.ItemRef
	} else {
		for i := range sc.Room.Items {
			if sc.Room.Items[i].Ref == ref && !sc.Room.Items[i].IsPut {
				ri = &sc.Room.Items[i]
				break
			}
		}
	}
	if ri == nil {
		return false
	}

	state := strings.ToUpper(ri.State)
	switch expectedState {
	case "OPEN":
		return state == "OPEN"
	case "CLOSED":
		return state == "CLOSED" || state == ""
	case "LOCKED":
		return state == "LOCKED"
	case "UNLOCKED":
		return state == "UNLOCKED" || state == "OPEN"
	}
	return false
}

// getVar retrieves a variable value for the player or current item.
func (sc *ScriptContext) getVar(name string) int {
	if strings.HasPrefix(name, "INTNUM") {
		idx, err := strconv.Atoi(name[6:])
		if err != nil {
			return 0
		}
		if sc.Player.IntNums == nil {
			return 0
		}
		return sc.Player.IntNums[idx]
	}
	if strings.HasPrefix(name, "ITEMBIT") {
		idx, err := strconv.Atoi(name[7:])
		if err != nil || sc.ItemRef == nil {
			return 0
		}
		if sc.ItemRef.Val4&(1<<idx) != 0 {
			return 1
		}
		return 0
	}
	if strings.HasPrefix(name, "ITEMVAL") {
		idx, err := strconv.Atoi(name[7:])
		if err != nil || sc.ItemRef == nil {
			return 0
		}
		switch idx {
		case 1:
			return sc.ItemRef.Val1
		case 2:
			return sc.ItemRef.Val2
		case 3:
			return sc.ItemRef.Val3
		case 4:
			return sc.ItemRef.Val4
		case 5:
			return sc.ItemRef.Val5
		}
		return 0
	}
	if strings.HasPrefix(name, "ITEMADJ") {
		idx, err := strconv.Atoi(name[7:])
		if err != nil || sc.ItemRef == nil {
			return 0
		}
		switch idx {
		case 1:
			return sc.ItemRef.Adj1
		case 2:
			return sc.ItemRef.Adj2
		case 3:
			return sc.ItemRef.Adj3
		}
		return 0
	}
	// SKILL variables — stored in player Skills map
	if strings.HasPrefix(name, "SKILL") {
		idx, err := strconv.Atoi(name[5:])
		if err != nil {
			return 0
		}
		if sc.Player.Skills == nil {
			return 0
		}
		return sc.Player.Skills[idx]
	}
	switch name {
	case "LEV":
		return sc.Player.Level
	case "RAC":
		return sc.Player.Race
	case "MISTFORM":
		// Ephemeral race (8) has innate mist form; otherwise check IntNum flag
		if sc.Player.Race == 8 {
			return 1
		}
		return 0
	}
	return 0
}

// setVar sets a variable value on the player or current item.
func (sc *ScriptContext) setVar(name string, val int) {
	if strings.HasPrefix(name, "INTNUM") {
		idx, err := strconv.Atoi(name[6:])
		if err != nil {
			return
		}
		if sc.Player.IntNums == nil {
			sc.Player.IntNums = make(map[int]int)
		}
		sc.Player.IntNums[idx] = val
		return
	}
	if strings.HasPrefix(name, "ITEMVAL") {
		idx, err := strconv.Atoi(name[7:])
		if err != nil || sc.ItemRef == nil {
			return
		}
		switch idx {
		case 1:
			sc.ItemRef.Val1 = val
		case 2:
			sc.ItemRef.Val2 = val
		case 3:
			sc.ItemRef.Val3 = val
		case 4:
			sc.ItemRef.Val4 = val
		case 5:
			sc.ItemRef.Val5 = val
		}
		// Sync full item snapshot to other machines
		itemCopy := *sc.ItemRef
		sc.Engine.notifyRoomChange(RoomChange{
			RoomNumber: sc.Room.Number, Type: "item_update",
			ItemRef: sc.ItemRef.Ref, Item: &itemCopy,
		})
		return
	}
}


// resolveNumericArg resolves a script argument that can be a literal number
// or a variable reference like ITEMVAL2.
func (sc *ScriptContext) resolveNumericArg(arg string) int {
	upper := strings.ToUpper(arg)
	if strings.HasPrefix(upper, "ITEMVAL") {
		return sc.getVar(upper)
	}
	val, err := strconv.Atoi(arg)
	if err != nil {
		return 0
	}
	return val
}

// expandScriptText replaces script placeholders in text.
func (sc *ScriptContext) expandScriptText(text string) string {
	text = strings.ReplaceAll(text, "%N", sc.Player.FirstName)
	text = strings.ReplaceAll(text, "%n", sc.Player.FirstName)
	if sc.ItemRef != nil && sc.ItemDef != nil {
		itemName := "it"
		if name, ok := sc.Engine.nouns[sc.ItemDef.NameID]; ok {
			itemName = name
		}
		text = strings.ReplaceAll(text, "%a", itemName)
	}
	// Gender-based pronouns
	if sc.Player.Gender == 0 {
		text = strings.ReplaceAll(text, "%h", "his")
		text = strings.ReplaceAll(text, "%e", "he")
		text = strings.ReplaceAll(text, "%o", "him")
	} else {
		text = strings.ReplaceAll(text, "%h", "her")
		text = strings.ReplaceAll(text, "%e", "she")
		text = strings.ReplaceAll(text, "%o", "her")
	}
	return text
}
