package game

// Inventory holds the items a player is carrying and enforces weight/slot limits.
type Inventory struct {
	Items     []*Item
	MaxWeight float64
	MaxItems  int
}

func newInventory() *Inventory {
	return &Inventory{
		MaxWeight: 20.0,
		MaxItems:  15,
	}
}

// CurrentWeight returns the total weight of all items in the inventory.
func (inv *Inventory) CurrentWeight() float64 {
	w := 0.0
	for _, item := range inv.Items {
		w += item.Weight
	}
	return w
}

// CanAdd reports whether item can be added without exceeding limits.
// Backpacks that increase weight capacity bypass the weight check.
func (inv *Inventory) CanAdd(item *Item) bool {
	if len(inv.Items) >= inv.MaxItems {
		return false
	}
	if item.Category == CategoryBackpack && item.StatMods.InvWeight > 0 {
		return true
	}
	return inv.CurrentWeight()+item.Weight <= inv.MaxWeight
}

// Add adds the item to the inventory. Returns false if limits would be exceeded.
func (inv *Inventory) Add(item *Item) bool {
	if !inv.CanAdd(item) {
		return false
	}
	inv.Items = append(inv.Items, item)
	return true
}

// Remove removes the item at the given index.
func (inv *Inventory) Remove(idx int) {
	if idx < 0 || idx >= len(inv.Items) {
		return
	}
	inv.Items = append(inv.Items[:idx], inv.Items[idx+1:]...)
}

// levelUp increases carrying capacity by 5%.
func (inv *Inventory) levelUp() {
	inv.MaxWeight *= 1.05
	inv.MaxItems = inv.MaxItems * 105 / 100
}
