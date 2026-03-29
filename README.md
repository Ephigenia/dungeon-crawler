# Dungeon Crawler

A tile-based 2D dungeon crawler written in Go using [Ebitengine](https://ebitengine.org/). Combines bump-attack combat with procedural dungeon generation and a full RPG inventory/equipment system.

Dungeons are **procedurally generated** using a rooms-and-corridors algorithm: random rectangular rooms are placed and connected by L-shaped corridors. Every run produces a different layout.

## Run

```bash
make run
# or
go run .
# or
make build && ./dungeon-crawler
```

## Controls

### World
| Key | Action |
|-----|--------|
| `WASD` / Arrow keys | Move one tile |
| `P` | Pick up item at current position |
| `O` | Open chest/object on adjacent tile |
| `R` | Generate a new dungeon |
| `Q` | Quit |

### Inventory (`I` to open/close)
| Key | Action |
|-----|--------|
| `Tab` | Switch focus between item grid and equipment slots |
| `Arrow keys` / `WAS` | Navigate |
| `U` / `Enter` | Use consumable · Equip/Unequip equipment |
| `D` | Drop selected item onto the floor (one at a time) |
| `X` | Destroy selected item (one at a time, cannot destroy equipped items) |

## Player

### Starting Stats
| Stat | Value | Description |
|------|-------|-------------|
| HP | 30 | Current / maximum hit points |
| Attack | 5 | Base melee attack power |
| Defense | 2 | Damage reduction |
| Agility | 5 | Synergises with weapon speed to amplify damage |
| Level | 1 | Increases via EXP; improves stats and damage variance |
| EXP to next level | 100 | Gained from hits (+5) and kills (+20) |
| Carry weight | 20 kg | Total inventory weight limit |
| Inventory slots | 15 | Total item slot limit |

The HUD and inventory panel show stats as `base / effective`, where *effective* includes all equipment bonuses.

### Leveling Up
Each level-up applies the following stat growth and then restores HP to the new effective maximum:

| Stat | Growth |
|------|--------|
| MaxHP | +10% per level |
| Attack | +1 every level |
| Defense | +1 every 2 levels |
| Agility | +1 every 3 levels |
| Inventory | +5% slots and weight per level |

The EXP threshold follows an exponential curve: `100 × level^1.5` (e.g. 283 for L2→L3, 520 for L3→L4).

Equipment bonuses are **never** baked into the base stats. Each effective stat is computed on the fly as `base + sum(equipped item bonuses)`, so equip/unequip always yields the exact correct value.

## Combat

Combat is triggered by **bumping into an enemy** (moving onto its tile). The enemy retaliates immediately if it survives.

### Player → Enemy
```
weaponContrib   = weaponPower × (1 + weaponSpeed × agility / 100)
effectiveAttack = (baseAttack + weaponContrib) × (1 + (level − 1) × 0.05)
randomBonus     = random(0 … level × 2)
damage          = max(0, int(effectiveAttack) − enemyDefense + randomBonus)
```
- Weapon **speed × agility** synergy rewards fast weapons on agile characters
- Each level adds **+5%** to effective attack and widens the random bonus range
- Damage can be **0** if defense exceeds the effective attack

### Enemy → Player
```
factor      = enemyAttack / playerDefense
randomBonus = random(0 … enemyAttack / 2)
damage      = max(0, (enemyAttack − playerDefense) × factor + randomBonus)
```
- The attack/defense ratio acts as a multiplier — a large gap is amplified quadratically
- Random bonus scales with the enemy's raw attack strength
- If `playerDefense = 0`, damage equals raw `enemyAttack + randomBonus`

### Enemies
| Name | HP | Attack | Defense |
|------|----|--------|---------|
| Goblin | 8 | 3 | 0 |
| Skeleton | 10 | 4 | 1 |
| Orc | 15 | 5 | 2 |
| Troll | 22 | 7 | 3 |

Defeated enemies have a level-scaled chance to drop a random item.

## Equipment

Items are equipped from the inventory and **stay in the inventory slot** (shown with a gold border). Unequipping returns the stat bonus immediately.

### Equipment Slots
| Slot | Notes |
|------|-------|
| Head | Helmets |
| Body | Armor |
| Legs | Pants |
| Feet | Shoes |
| Necklace | Stat accessories |
| Left / Right Hand | Gloves |
| Left / Right Ring | Stat rings |
| Right Weapon | Weapons only (one weapon at a time) |
| Left Weapon | Shields |
| Backpack | Increases carry weight and inventory slots |

### Weapons
Each weapon has a **Power** (raw damage added to the attack formula) and a **Speed** (synergises with player agility). Only one weapon can be equipped at a time (right weapon slot).

| Weapon | Power | Speed |
|--------|-------|-------|
| Iron Sword | 3 | 5 |
| Broadsword | 5 | 3 |
| Golden Sword | 7 | 5 |
| Sword Jeweled | 8 | 5 |
| Mega Sword | 14 | 2 |
| Saber | 4 | 7 |
| Rapier (×2) | 3 / 6 | 9 |
| Axe | 5 | 4 |
| Hatchet | 4 | 7 |
| Knights Axe | 9 | 3 |
| Executioner's Axe | 12 | 2 |

## Items

### Consumables
Health potions are **stackable** (up to 5 per slot). All consumables restore HP on use.

| Item | Heal |
|------|------|
| Small Health Potion | 5 HP |
| Medium Health Potion | 10 HP |
| Large Health Potion | 20 HP |
| Food (apple, bread, grapes, egg, meat, mushroom, pizza) | 1–5 HP |

### Other Categories
| Category | Examples |
|----------|---------|
| Armor | Basic (+1 DEF), bronze (+3), complex (+2), gold (+5) |
| Shields | Wooden (+2 DEF), bronze (+2), metal (+4), gold (+4) |
| Helmets | Coif, basic (+1 DEF), full (+2), horn (+2), gold (+3) |
| Gloves | Finger, leather (+1 DEF), leather-metal (+2), metal (+3) |
| Shoes | Simple, leather, metal (+1 DEF), gold (+2) |
| Necklaces | Skull (+20 HP), diamond/star/tooth (+5 HP each) |
| Rings | Gold (+2 ATK), silver (+1 DEF), diamond (+3 ATK), diamond2 (+2 DEF) |
| Backpacks | Small (+10 slots, +5 kg), medium (+15 slots, +7 kg), large (+15 slots, +20 kg) |

## Map Objects

**Chests** (wooden or iron) spawn at room centres with a **25% chance per room**. They block movement until opened.

- Stand adjacent and press `O` to open
- A 20-frame opening animation plays
- On completion, **1–5 random items** are dropped at the player's position
- Objects have `PassableByPlayer` and `PassableByEnemy` flags (chests are impassable by default)

## Code Structure

```
main.go                  – entry point, window setup, embedded assets
dungeon/                 – procedural dungeon generation (rooms + corridors)
game/
  game.go                – Game struct, constants, Layout()
  game_init.go           – New(), resetEntities(), Regenerate(), helpers (potionAt, enemyAt, objectAt)
  input.go               – Update(), resolveCombat(), inventory input handlers
  render.go              – Draw(), HUD, world rendering, shared draw helpers
  render_inventory.go    – inventory overlay and item detail panel
  player.go              – Player struct, stats, leveling, equip/unequip, TakeDamage
  enemy.go               – Enemy struct, spawn table, calcDamage, calcPlayerDamage
  inventory.go           – InventorySlot, weight/slot management, stacking, Consume()
  equipment.go           – EquipmentSlot constants, StatModifiers, slot labels
  item.go                – Item struct (Power, Speed, MaxStack, OnUse…), ItemCategory
  items.go               – All item definitions, AllItems, SpawnableItems
  item_images.go         – Sprite loading from embedded FS
  potion.go              – Map pickup entity (wraps any Item)
  object.go              – Object struct (map-placed interactables: chests), ObjectKind/State
```

## Ideas

- [ ] general: non-grid movement, move by pixels and not grid positions
- [ ] buffs: timed effects on items & player stats (weapons, armor, equipables, books, scrolls, oils)
- [ ] bug: when l.hand/r.hand item is used another item which matches the slot both items are highlighted as if they are used
- [ ] general: general game concept
- [ ] graphics: high-resolution images
- [ ] graphics: pixel perfect font rendering
- [ ] debug: add debug menu/command line
- [ ] enemies: sensor based follow (https://www.lexaloffle.com/bbs/?tid=48889)
- [ ] enemies: attack types?
- [ ] enemies: spawn positions
- [ ] enemies: visualize difficulty?
- [ ] enemies: wayfinding
- [ ] fight: different attack modes?
- [ ] fight: range?
- [ ] fight: critical chance
- [ ] fight: armor + armor penetration
- [ ] fight: step-based fight concept
- [ ] items: durability for weapons
- [ ] items: modifier target: player, enemy
- [ ] items: rarities (influencing the stats and positively increasing stats)
- [ ] items: requirements for equipping items
- [ ] map: 47 floor tileset
- [ ] map: room & floor decorations, corners, shaped walls, traps, doors, secrets
- [ ] map: secrets (walk next to secrets to discover)
- [ ] map: generation algorithms
- [ ] player: stats vs. attributes (attributes long term development, stats short term), core attribute vs. attribute
- [ ] player: core attributes: strength, 
- [ ] QoC: abstraction for image loading, assets, spritemaps, tilemaps
- [ ] QoC: object, items reference instance of spreadsheet (name, index)
- [ ] QoC: abstraction for tilemaps
- [ ] QoC: abstraction for tilemaps with animations
- [ ] QoL: multiple items on the floor, "a" to get them all
- [ ] QoL: multiple items on the floor, "p" opens interface to select which one to take, a (all)
- [ ] tests: add tests for certain assumptions
- [ ] tests: attack increases with level
- [ ] tests: items cannot be equipped when stats not matching


#### Done

- [x] debug: include "r" which re-creates level
- [x] map: walkable and not walkable items, (shelf)
- [x] map: destructable objects, not walkable but attackable, have health-points and spawn items on destruction
- [x] map: shelfes
- [x] map: simple tileset floor
- [x] map: simple tileset walls
- [x] map: simple base map (floor/wall) graphics
- [x] enemies: walking speed
- [x] enemies: looking distance
- [x] enemies: idle+chase
- [x] debug: fps counter
- [x] fight: how to handle attacks from enemies (counter-attacks, blocks)