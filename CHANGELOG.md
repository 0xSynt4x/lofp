# Changelog

## v0.9 — 2026-04-03

### Game Engine
- Script parser loads 2273 rooms, 1990 items, 297 monsters from original .SCR files
- Fixed description parser: `*DESCRIPTION_START ITEM READ/EXAM` no longer overwrites room descriptions
- Fixed item deduplication: later script definitions correctly override earlier ones (ANTI.SCR placeholders replaced)
- Fixed script block parser: missing ENDIF no longer consumes subsequent item definitions
- Seasonal scripts (ASCRIPT/WSCRIPT/SSCRIPT/PSCRIPT) skipped; base definitions used
- Formatted text preservation: descriptions with leading whitespace or blank lines retain formatting (poems, maps, etc.)

### Movement & Portals
- Portal traversal via GO command with VAL2 destination
- PORTAL_CLIMB, PORTAL_UP, PORTAL_DOWN, PORTAL_OVER, PORTAL_THROUGH support
- CLIMB verb for climbable portals
- Closed/locked portals block passage
- Directional look (LOOK N, LOOK NORTH, etc.) shows adjacent room description

### Script Execution
- IFENTRY blocks execute on room entry (movement, login, character creation)
- IFVAR conditions: INTNUM, ITEMBIT, ITEMVAL, ITEMADJ, LEV, RAC, SKILL, MISTFORM
- IFITEM conditions: check item open/closed/locked/unlocked state
- IFPREVERB/IFVERB execution on item interactions
- Actions: ECHO (PLAYER/ALL/OTHERS), EQUAL, ADD, SUB, NEWITEM, GMMSG, CLEARVERB, MOVE, SHOWROOM
- Actions: LOCK, UNLOCK, OPEN, CLOSE, REMOVEITEM, SETITEMVAL
- Nested conditional blocks supported
- Script text placeholders: %N, %n, %a, %h, %e, %o

### Commands
- Movement: N/S/E/W/NE/NW/SE/SW/UP/DOWN/OUT, GO, CLIMB
- Looking: LOOK, EXAMINE, INSPECT, LOOK IN/ON/UNDER/BEHIND, directional LOOK
- Items: GET, DROP, INVENTORY, WIELD, UNWIELD, WEAR, REMOVE
- Containers: OPEN, CLOSE
- Interaction: PULL, PUSH, TURN, RUB, TAP, TOUCH, SEARCH, DIG, RECALL
- Communication: SAY ('), WHISPER, YELL, RECITE, EMOTE
- Commerce: BUY, SELL (1482 store items across 20+ shops)
- Social: 60+ emotes (SMILE, BOW, KICK, etc.) with targeted second-person messages
- Info: STATUS, HEALTH, WEALTH, SKILLS, WHO, HELP, ASSIST
- Position: SIT, STAND, KNEEL, LAY
- Roleplay: ACT, EMOTE, RECITE
- READ items with room-scoped descriptions

### Multiplayer
- WebSocket-based real-time gameplay
- Cross-machine coordination via MongoDB Change Streams
- Player presence synced across multiple Fly.io machines
- Room state changes (item state, drops, pickups, script mutations) coordinated
- Emote targeting: second-person messages to target, third-person to room
- Room broadcasts, global broadcasts, GM broadcasts, whispers
- WebSocket reconnection with exponential backoff

### Characters & Auth
- Google OAuth login with 30-day JWT sessions
- Character creation with 8 races, stat rolling, name validation
- Character name uniqueness enforced (no logging into others' characters)
- Session persistence across server restarts
- Starting gear granted via IFENTRY scripts at room 201

### Commerce
- BUY command with store items (STOREITEM parsed from scripts)
- Efficient currency deduction (copper first, then silver, then gold with change)
- SELL command in rooms with BUY_ARMOR/BUY_SKINS/BUY_JEWELRY modifiers
- Price display as gold/silver/copper

### Admin Interface
- Rooms tab: searchable index with detail view, clickable exits, enriched items
- Items tab: searchable index with full properties, type/slot expansion
- Monsters tab: searchable index with combat stats and properties
- Players tab: character detail, GM toggle, account reassignment
- Users tab: account detail, admin toggle, character list
- Logs tab: filterable game event log with hyperlinked players/users
- Exact number match sorting in search results

### Logging
- MongoDB-backed game log with 90-day TTL
- Events: user login/logout, character game enter/exit, character creation, GM grant/revoke
- Log entries include user name, email, account ID
- Admin logs UI with event/player filtering

### Security
- Admin auth required for GM toggle, game world API endpoints
- WebSocket origin validation against frontend URL
- 64KB WebSocket message size limit
- Character name validation (alpha + ' + -, max 20 chars)
- Race/gender input validation
- Command input truncated to 500 characters

### Infrastructure
- Go + gorilla/mux backend, React 19 + TypeScript + Vite + Tailwind 4 frontend
- MongoDB Atlas for persistence
- Fly.io production deployment (2 machines, ord region)
- Custom domain: lofp.metavert.io
- Multi-stage Docker build (14MB production image)
- Separate dev/prod configs
