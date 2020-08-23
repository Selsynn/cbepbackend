package item

import "fmt"

type ID string

const (
	Bow     ID = "bow"
	Wood    ID = "wood"
	Leather ID = "leather"
)

type Item struct {
	Name ID
}

type Resources struct {
	Item
	Qty int
}

func (r *Resources) Describe() string {
	if r.Qty == 0 {
		return ""
	}
	return fmt.Sprintf("(%d *%s*)", r.Qty, r.Name)
}
