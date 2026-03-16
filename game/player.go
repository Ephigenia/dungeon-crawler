package game

// Player holds the player's state and stats.
type Player struct {
	X, Y         int
	HP           int
	MaxHP        int
	Attack       int
	Defense      int
	Level        int
	EXP          int
	NextLevelEXP int
	Inventory    *Inventory
}

func newPlayer(x, y int) *Player {
	return &Player{
		X:            x,
		Y:            y,
		HP:           30,
		MaxHP:        30,
		Attack:       5,
		Defense:      2,
		Level:        1,
		EXP:          0,
		NextLevelEXP: 100,
		Inventory:    newInventory(),
	}
}

// levelUp increases the player's level and improves stats.
func (p *Player) levelUp() {
	p.Level++
	p.MaxHP = p.MaxHP * 110 / 100
	p.HP = p.MaxHP
	p.NextLevelEXP = p.NextLevelEXP * 125 / 100
	p.Inventory.levelUp()
}

// AddEXP adds exp points and calls levelUp each time the threshold is reached.
func (p *Player) AddEXP(amount int) {
	p.EXP += amount
	for p.EXP >= p.NextLevelEXP {
		p.EXP -= p.NextLevelEXP
		p.levelUp()
	}
}

// IsAlive returns true if the player has HP remaining.
func (p *Player) IsAlive() bool {
	return p.HP > 0
}

// TakeDamage reduces HP by the incoming attack minus defense, minimum 1.
func (p *Player) TakeDamage(attack int) {
	dmg := attack - p.Defense
	if dmg < 1 {
		dmg = 1
	}
	p.HP -= dmg
	if p.HP < 0 {
		p.HP = 0
	}
}
