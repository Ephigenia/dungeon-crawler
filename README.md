# Dungeon Crawler

A 2D dungeon crawler written in Go using [Ebitengine](https://ebitengine.org/). Combines hack-and-slash combat with roguelike dungeon generation and RPG inventory/equipment systems.

Levels are **procedurally generated** with a rooms-and-corridors algorithm: random rectangular rooms are placed and connected by L-shaped corridors.

## Run

```bash
make run
# or
go run .
# or
make build && ./dungeon-crawler
```

## Controls

### Movement
| Key | Action |
|-----|--------|
| `WASD` / Arrow keys | Move one tile |
| `R` | Generate a new dungeon |

### World
| Key | Action |
|-----|--------|
| `P` | Pick up item at current position |
| `Q` | Quit |

### Inventory (`I` to open/close)
| Key | Action |
|-----|--------|
| `Tab` | Switch focus between item grid and equipment slots |
| `WASD` / Arrow keys | Navigate |
| `U` / `Enter` | Use (consumable) or Equip/Unequip (equipment) |
| `X` | Destroy item (cannot destroy equipped items) |

## Gameplay

- **Combat** – bump into an enemy to attack; enemies retaliate immediately
- **Leveling** – earn EXP from hits and kills; leveling up improves HP, carry weight, and inventory slots
- **Pickups** – items are scattered across rooms; walk over them and press `P` to pick up
- **Equipment** – items stay in inventory when equipped (shown with a gold border); multiple slots supported (e.g. rings, weapons fit either hand)
- **Backpacks** – equipping a backpack increases maximum carry weight and inventory slots

## Items

| Category | Examples |
|----------|---------|
| Consumables | Health potions (small/medium/large), food (apple, bread, meat, pizza…) |
| Weapons | Iron sword |
| Armor | Basic, bronze, complex, gold |
| Shields | Wooden, metal, gold, bronze |
| Helmets | Coif, basic, full, horn, gold |
| Gloves | Leather, finger, leather-metal, metal |
| Shoes | Simple, leather, metal, gold |
| Accessories | Necklaces (skull, diamond, star, tooth), rings (gold, silver, diamond) |
| Backpacks | Small, medium, large |

## Code Structure

```
main.go                  – entry point, window setup, embedded assets
dungeon/                 – procedural dungeon generation (rooms + corridors)
game/
  game.go                – Game struct, constants, Layout()
  game_init.go           – New(), resetEntities(), Regenerate(), potionAt(), enemyAt()
  input.go               – Update(), combat resolution, inventory input handlers
  render.go              – Draw(), HUD, world rendering, shared draw helpers
  render_inventory.go    – inventory overlay and detail panel
  player.go              – Player struct, leveling, equip/unequip, stat modifiers
  enemy.go               – Enemy struct, damage, spawn table
  inventory.go           – Inventory weight/slot management
  equipment.go           – Equipment slots, StatModifiers, slot labels
  item.go                – Item struct, ItemCategory constants
  items.go               – All item definitions, AllItems, SpawnableItems
  item_images.go         – Sprite loading from embedded FS
  potion.go              – Map pickup entity (holds any Item)
```

## Ideas

- easy
  - add enemy types
  - enemy images
- multiple items on the floor, "p" opens interface to select which one to take, a (all)
- item spawn probabilties
- graphics: levels
- enemies: different types
- enemies: 
- items: modifier target: player, enemy
- weapons: range
- weapons: different attack modes?
- fight: step-based fight concept
- item rarities (influencing the stats and positively increasing stats)
- Drop items onto the map
- Level/stat requirements for equipping items
- Item highlighting in inventory when a compatible equipment slot is selected
- Enemy variants with unique abilities
- Animated sprites and asset abstraction layer
