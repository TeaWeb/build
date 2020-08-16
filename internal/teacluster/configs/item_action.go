package configs

type ItemActionString = string

const (
	ItemActionAdd    = "add"
	ItemActionChange = "change"
	ItemActionRemove = "remove"
)

type ItemAction struct {
	ItemId string
	Action ItemActionString
	Item   *Item
}

func NewItemAction(itemId string, action ItemActionString, item *Item) *ItemAction {
	return &ItemAction{
		ItemId: itemId,
		Action: action,
		Item:   item,
	}
}
