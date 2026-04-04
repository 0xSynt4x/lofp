# Legends of Future Past

### **[Play it now — free forever — at lofp.metavert.io](https://lofp.metavert.io)**

---

<p align="center">
  <img src="frontend/src/assets/hero.png" alt="Legends of Future Past" width="300">
</p>

**Legends of Future Past** was the first commercial text-based [MUD](https://en.wikipedia.org/wiki/Multi-user_dungeon) to make the transition from a proprietary network provider (CompuServe) to the Internet. Designed by [Jon Radoff](https://metavert.io) and Angela Bull, it launched in 1992 and ran continuously until December 31, 1999.

The game was originally offered for $6.00 per hour via CompuServe, and later at reduced rates on the Internet. It was notable for having paid Game Masters who conducted live online events — a concept that would become standard in the MMO industry years later.

*Computer Gaming World* awarded *Legends of Future Past* a Special Award for Artistic Excellence in their 1993 Online Game of the Year competition, stating they were "overwhelmed by the creative power of storytelling and fertile liveliness." *Computer Game Review* gave it the Golden Triad Award. CGW wrote that it was "a rich, dynamic and lovingly supervised world of the imagination... Like most of these games, this one is extremely addicting — perhaps even more so."

*Legends* introduced one of the first (if not *the* first) crafting systems in an online game. Players could harvest ores, herbs, and skins, then use them to craft weapons, armor, and enchanted items. The game was skill-based with no class archetypes and no level caps — some dedicated players attained levels in the hundreds.

The game is credited with spawning a number of other online games and introducing top talent into the MMORPG industry. Many GameMasters and developers went on to become founders or product managers at top online games including SOE's *Star Wars Galaxies*, Worlds Apart Productions, and Dejobaan Games.

## Resurrection with AI

In 2026, Jon Radoff — the original creator — resurrected *Legends of Future Past* using [Claude Code](https://claude.ai/claude-code), Anthropic's agentic coding tool.

**The original game engine source code is lost.** But several former gamemasters had preserved copies of all 333 game script files, the GM manual, player documentation, scripting guides, a 1996 gameplay session capture, and an alchemy recipe guide. From these artifacts alone, Claude Code was able to reconstruct the entire game:

- **Parsed and interpreted** the custom scripting language that defined the game world — 2,273 rooms, 1,990 items, 297 monsters
- **Reverse-engineered combat mechanics** by analyzing session capture logs, cross-referencing monster stats with weapon damage, and reading the original GM documentation
- **Reconstructed the combat system** with faithful output formatting (`[ToHit: X, Roll: Y] Hit!`), damage severity tiers, weapon clash mechanics, and the original XP/build-point progression table
- **Implemented all major systems**: 60+ spells, 30+ psionic disciplines, 36 skills, a complete crafting/alchemy system, and monster AI with 7 strategy tiers
- **Matched original behavior** by comparing output against the 1996 session capture of actual gameplay

The entire reconstruction — from an empty repository to a fully playable multiplayer game — was accomplished through iterative collaboration between the original creator and Claude Code.

## Why This Matters

Classic games are [disappearing at an alarming rate](https://medium.com/mr-plan-publication/classic-games-disappearing-what-it-means-for-gamings-future-2a885dc3febc). The preservation crisis is especially acute for online games that required servers — when the servers shut down, the game is gone forever. Unlike single-player games that can be emulated, online games need their server code reconstructed to live again.

*Legends of Future Past* is now **free to play forever**, released under the MIT License, so it can live on in posterity. The server at [lofp.metavert.io](https://lofp.metavert.io) is running today, and anyone can clone this repository to host their own instance.

With AI and agentic coding, we now have a powerful new tool for digital preservation. Games like *Legends of Future Past* can be resurrected from their data files, documentation, and captured gameplay — even when the original source code is lost. This project demonstrates that approach, and we hope it inspires the preservation of other classic online games for future generations and computer archaeologists.

The problem of online game preservation is real and getting worse. With online games requiring servers to function, shutting down the servers means losing the game entirely. Today, with AI-powered tools like Claude Code, we can begin to reverse-engineer and preserve these experiences from whatever artifacts survive — script files, documentation, player captures, and community memories. This is how we keep gaming history alive.

## The World of Andor

The game takes place on **Andor**, set in the "Shattered Realms" — a world featuring a blend of fantasy and ancient technology. Most of the action revolves around the city of **Fayd**, which serves as the hub of activity for adventures, intrigue, and roleplaying events.

**Eight races** inhabit Andor:
- **Aelfen** — an elflike species
- **Drakin** — dragon-men with player-created languages and cultural institutions
- **Ephemerals** — a wraithlike species that cannot be harmed unless they choose to manifest
- **Highlander** — a stout, mountain-dwelling people
- **Humans** — the only people who can utilize ancient technology
- **Murg** — a proud warrior race
- **Mechanoids** — artificial beings
- **Wolflings** — a race of shapechangers

**Key features:**
- **Combat** with the original to-hit/damage system, weapon crits, slayer weapons, elemental immunities, and monster special attacks
- **Magic** across five schools: Conjuration, Enchantment, Necromancy, General, and Druidic
- **Psionics** with Mind over Matter and Mind over Mind disciplines
- **Crafting**: mining, smelting, forging, weaving, dyeing, and alchemy with 32 potion recipes
- **36 skills** from Edged Weapons to Lockpicking to Sagecraft, with build point costs and prerequisites
- **Real-time multiplayer** with cross-server coordination
- **Over 2,000 rooms** to explore across cities, forests, caverns, volcanic islands, and astral planes

## Architecture

| Component | Technology | Path |
|-----------|-----------|------|
| Backend | Go + gorilla/mux + MongoDB | `engine/` |
| Frontend | React 19 + TypeScript + Vite + Tailwind 4 | `frontend/` |
| Scripts | Original 1992-1999 game data (333 .SCR files) | `original/scripts/` |
| Documentation | GM Manual, player docs, session captures | `original/` |

## Running Your Own Server

```bash
# Prerequisites: Go 1.25+, Node.js 22+, MongoDB

# Clone the repository
git clone https://github.com/jonradoff/lofp.git
cd lofp

# Set up environment variables
cp .env.example .env
# Edit .env with your MONGODB_URI, JWT_SECRET, GOOGLE_CLIENT_ID

# Start both frontend and backend
./start.sh
# Frontend: http://localhost:4992
# Backend: http://localhost:4993
```

## Credits

**Original Game (1992-1999)** — Copyright (c) 1992-1999 Inner Circle Software / NovaLink USA Corp

| Role | Name |
|------|------|
| Created & Programmed by | Jon Radoff |
| Additional Programming | Ichiro Lambe |
| Co-Producer | Angela Bull |
| Legends Manager | Gary Whitten |
| World Building | Gary Whitten, David Goodman, Tony Spataro, Stacy Jannis, Kevin Jepson, Daniel Brainerd, Michael Hjerppe |
| Documentation | Gary Whitten |
| Quality Assurance | David Goodman, Stacy Jannis |

**2026 Resurrection** — Reimplemented from original script files and documentation by [Jon Radoff](https://metavert.io) using [Claude Code](https://claude.ai/claude-code).

Special thanks to **David Goodman** for supplying much of the original materials used to reconstruct the game.

## License

Released under the [MIT License](LICENSE) — Copyright (c) 2026 Metavert LLC.

## Links

- **Play now**: [lofp.metavert.io](https://lofp.metavert.io)
- **Version Notes**: [lofp.metavert.io/version-notes](https://lofp.metavert.io/version-notes)
- **Wikipedia**: [Legends of Future Past](https://en.wikipedia.org/wiki/Legends_of_Future_Past)
- **Online World Timeline**: [raphkoster.com](https://www.raphkoster.com/games/the-online-world-timeline/)
- **Game Preservation**: [Classic Games Disappearing](https://medium.com/mr-plan-publication/classic-games-disappearing-what-it-means-for-gamings-future-2a885dc3febc)
