package characters

import "github.com/Selsynn/cbepbackend/business/item"

type Adventurer struct {
	Name          string
	Relationships map[*NPC]int
	ID            string
	Items         map[item.ID]*item.Resources
}

func (a *Adventurer) DeltaItems(res []*item.Resources) {
	for _, element := range res {
		if a.Items[element.Name] == nil {
			a.Items[element.Name] = &item.Resources{
				Item: item.Item{
					Name: element.Name,
				},
				Qty: 0,
			}
		}
		a.Items[element.Name].Qty += element.Qty
	}
}
