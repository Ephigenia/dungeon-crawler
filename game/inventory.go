package game

// InventorySlot holds one item type and its stack count.
type InventorySlot struct {
	Item  *Item
	Count int
}

// Inventory holds the items a player is carrying and enforces weight/slot limits.
type Inventory struct {
	Items     []*InventorySlot
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
	for _, slot := range inv.Items {
		w += slot.Item.Weight * float64(slot.Count)
	}
	return w
}

// CanAdd reports whether item can be added without exceeding limits.
// For stackable items an existing partial stack counts as available space.
// Items that go in the backpack slot (which expand capacity) bypass the weight check.
func (inv *Inventory) CanAdd(item *Item) bool {
	// Stackable: check for an existing slot with room.
	if item.MaxStack > 1 {
		for _, slot := range inv.Items {
			if slot.Item == item && slot.Count < item.MaxStack {
				return true
			}
		}
	}
	// Need a new slot.
	if len(inv.Items) >= inv.MaxItems {
		return false
	}
	if item.FitsSlot(SlotBackpack) && item.StatMods.InvWeight > 0 {
		return true
	}
	return inv.CurrentWeight()+item.Weight <= inv.MaxWeight
}

// Add adds the item to the inventory. Returns false if limits would be exceeded.
// Stackable items are merged into an existing partial slot when possible.
func (inv *Inventory) Add(item *Item) bool {
	// Try to merge into an existing stack first.
	if item.MaxStack > 1 {
		for _, slot := range inv.Items {
			if slot.Item == item && slot.Count < item.MaxStack {
				slot.Count++
				return true
			}
		}
	}
	if !inv.CanAdd(item) {
		return false
	}
	inv.Items = append(inv.Items, &InventorySlot{Item: item, Count: 1})
	return true
}

// Remove removes the entire slot at the given index.
func (inv *Inventory) Remove(idx int) {
	if idx < 0 || idx >= len(inv.Items) {
		return
	}
	inv.Items = append(inv.Items[:idx], inv.Items[idx+1:]...)
}

// Consume removes one item from the stack at idx.
// When the stack reaches zero the slot is removed entirely.
func (inv *Inventory) Consume(idx int) {
	if idx < 0 || idx >= len(inv.Items) {
		return
	}
	inv.Items[idx].Count--
	if inv.Items[idx].Count <= 0 {
		inv.Remove(idx)
	}
}

// levelUp increases carrying capacity by 5%.
func (inv *Inventory) levelUp() {
	inv.MaxWeight *= 1.05
	inv.MaxItems = inv.MaxItems * 105 / 100
}
