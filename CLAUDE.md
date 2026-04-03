# LoFP — Legends of Future Past

Resurrecting a 1990s MUD from original script files. The original game engine source code is lost; only the content scripts and documentation survive. We reverse-engineer a working game from those.

## Architecture

- **Backend**: Go + gorilla/mux + MongoDB at `engine/`
- **Frontend**: React 19 + TypeScript + Vite + Tailwind 4 at `frontend/`
- **Original Scripts**: `original/scripts/` (read-only reference, 333 .SCR files)
- Frontend dev server: port 4992, Backend server: port 4993
- Start both: `./start.sh`

## Build Check

```sh
cd engine && go build ./...
cd frontend && npx tsc --noEmit
```

## Multi-Machine Coordination

Production runs multiple Fly.io machines. ALL mutable world state must be coordinated via the MongoDB-backed hub (`engine/internal/hub/`):

- **Messages**: broadcasts, whispers, global announcements → published to `events` collection, delivered via Change Streams
- **Player presence**: WHO list, room occupancy → `presence` collection with TTL heartbeat
- **Room state changes**: item open/close/lock, item drops/pickups, script mutations (vals, itembits) → `room_state_change` events via hub
- **Any new mutable state** added to rooms, items, or global data MUST call `notifyRoomChange()` or publish through the hub, or players on different machines will see inconsistent worlds

## After Server Changes

After making changes to the backend (engine/), restart the Go server:
```sh
kill $(lsof -ti:4993) 2>/dev/null; sleep 1; cd engine && go run cmd/lofp/main.go &
```
Load .env first if needed: `source .env`

## Script Language

The game world is defined in a custom scripting language (documented in `original/gmscript.doc`):
- **Rooms**: `NUMBER`, `NAME`, `*DESCRIPTION_START/END`, `EXIT`, `ITEM`, terrain, lighting
- **Items**: `INUMBER`, `NAME` (noun ref), type, weight, volume, substance, worn slots
- **Monsters**: `MNUMBER`, body parts, stats, AI strategy, weapons, spells
- **Events**: `IFVERB/IFPREVERB/IFSAY/IFENTRY/IFVAR...ENDIF` conditional blocks
- **Variables**: Named variables + internal vars (stats, time, flags, item vals)
- Config file: `original/scripts/LEGENDS.CFG` lists all scripts to load in order

## Current State (Phase 1)

- Script parser loads 2843 rooms, 2898 items, 319 monsters, 1221 nouns, 1352 adjectives
- Game engine supports: movement, look, get/drop, inventory, wield/wear, open/close, status, emotes
- Players start at Room 201 (City Gate) — tutorial room 3950 needs script execution engine
- WebSocket-based real-time gameplay
- Admin panel for inspecting rooms/items/scripts

## Units

| Unit | Code | Path | Description |
|------|------|------|-------------|
| Engine | ENG | `engine` | Go backend: script parser, game engine, command interpreter, MongoDB persistence |
| Frontend | FRONT | `frontend` | React + Tailwind: player UI (text client) and admin interface |
| Scripts | SCR | `original` | Original game script files and documentation (read-only reference) |

## Future Phases

- Script execution engine (IFVERB, IFSAY, IFENTRY, ECHO, MOVE, SPELL, etc.)
- Combat system
- NPC/monster AI and spawning
- Spells and psionics
- Crafting, mining, foraging
- Tutorial room (3950) sequence
- Multiplayer support
