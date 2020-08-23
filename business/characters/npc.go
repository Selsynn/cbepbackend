package characters

import (
	"fmt"
	"sync"
	"time"

	"github.com/Selsynn/cbepbackend/business/item"
	"github.com/pkg/errors"
)

var ActionNotAllowedErr = errors.New("Action is not allowed")

type Profession string

const (
	Merchant      Profession = "Merchant"
	WeaponCrafter Profession = "WeaponCrafter"
)

//durationPerUnit is the duration to do 1 Unit of work (1sec)
const durationPerUnit = time.Second

type (
	RelationLevel int
	NPCName       string
	NPC           struct {
		Name           NPCName
		Specialty      Profession
		RelQuest       map[int]Quest
		workQueue      []WorkItem
		workDone       []WorkItem
		nextWorkItem   *time.Timer
		lastWorkedTime time.Time
		workMutex      sync.Mutex
		Knowledge      []WorkItem
	}

	Quest struct{}

	WorkItem struct {
		CreatedAt    time.Time
		WorkCost     float64
		Player       *Adventurer
		Item         *item.Item
		WorkProgress float64
		Cost         []item.Resources
	}
)

func (n *NPC) AddWork(i *item.Item, cost float64, a *Adventurer) error {
	n.workMutex.Lock()
	defer n.workMutex.Unlock()
	t := time.Now()
	units := float64(t.Sub(n.lastWorkedTime) / durationPerUnit)
	n.lastWorkedTime = t

	work := WorkItem{
		CreatedAt:    t,
		Item:         i,
		Player:       a,
		WorkCost:     cost,
		WorkProgress: 0.0,
	}

	ref := &work

	for _, w := range n.workQueue {
		w.WorkProgress += units
		if (ref.WorkCost - ref.WorkProgress) > (w.WorkCost - w.WorkProgress) {
			ref = &w
		}
		if a.ID == ref.Player.ID {
			//warning there should not have any work already in progress with that npc
			//return ActionNotAllowedErr
			return errors.Wrapf(ActionNotAllowedErr, "Player %s cannot have new item with NPC %s as there is already the item %s ", a.ID, n.Name, w.Item.Name)
		}
	}

	n.workQueue = append(n.workQueue, work)

	nextWorkItemTime := durationPerUnit * time.Duration(ref.WorkCost-ref.WorkProgress) / time.Duration(len(n.workQueue))

	fmt.Printf("Next work item will be ready in %f s \n", nextWorkItemTime.Seconds())

	n.nextWorkItem = time.NewTimer(nextWorkItemTime)

	go func() {
		select {
		case <-n.nextWorkItem.C:
			n.RefreshWork()
		}
	}()

	return nil
}

func (n *NPC) RefreshWork() {
	n.updateTimeWork()
}

func (n *NPC) GetAdventurerWIPItem(a *Adventurer) *WorkItem {
	for _, w := range n.workQueue {
		if w.Player.ID == a.ID {
			return &w
		}
	}
	return nil
}

func (n *NPC) GetAdventurerDoneItems(a *Adventurer) []*item.Item {
	res := make([]*item.Item, 1)
	for _, i := range n.workDone {
		if i.Player.ID == a.ID {
			res = append(res, i.Item)
		}
	}
	return res
}

func (n *NPC) CancelWIPItem(a *Adventurer, i *item.Item) error {
	n.workMutex.Lock()
	defer n.workMutex.Unlock()

	res := make([]WorkItem, len(n.workQueue))
	for _, w := range n.workQueue {
		if w.Player.ID == a.ID && w.Item.Name == i.Name {
			continue
		}
		if w.Player.ID == a.ID {
			return errors.Wrapf(ActionNotAllowedErr, "Player has an wip item %s, different from expected wip item %s", w.Item.Name, i.Name)
		}

		res = append(res, w)
	}

	n.workQueue = res

	n.updateTimeWork()

	return nil
}

func (n *NPC) RecoltDoneItems(a *Adventurer) []item.Item {
	n.workMutex.Lock()
	defer n.workMutex.Unlock()

	res := make([]WorkItem, len(n.workDone))
	list := make([]item.Item, 5)
	for _, w := range n.workDone {
		if w.Player.ID == a.ID {
			list = append(list, *w.Item)
			continue
		}

		res = append(res, w)
	}

	n.workDone = res

	return list
}

func (n *NPC) updateTimeWork() {
	n.workMutex.Lock()
	defer n.workMutex.Unlock()

	t := time.Now()
	units := float64(t.Sub(n.lastWorkedTime) / durationPerUnit)
	n.lastWorkedTime = t

	newWork := make([]WorkItem, len(n.workQueue))

	var ref *WorkItem

	for _, w := range n.workQueue {
		w.WorkProgress += units
		if w.WorkProgress < w.WorkCost {
			newWork = append(newWork, w)
		} else {
			n.workDone = append(n.workDone, w)
		}
		if ref == nil || (ref.WorkCost-ref.WorkProgress) > (w.WorkCost-w.WorkProgress) {
			ref = &w
		}
	}
	n.workQueue = newWork

	n.nextWorkItem = time.NewTimer(durationPerUnit * time.Duration(ref.WorkCost-ref.WorkProgress))

	go func() {
		select {
		case <-n.nextWorkItem.C:
			n.RefreshWork()
		}
	}()

}

func InitNPCForNewCity() map[Profession]*NPC {
	return map[Profession]*NPC{
		Merchant: {
			lastWorkedTime: time.Now(),
			workDone:       make([]WorkItem, 0),
			workQueue:      []WorkItem{},
			Name:           "Gripsou",
			RelQuest:       map[int]Quest{},
			Specialty:      Merchant,
		},
		WeaponCrafter: {
			Name:           "Kreator",
			lastWorkedTime: time.Now(),
			workDone:       make([]WorkItem, 0),
			workQueue:      []WorkItem{},
			RelQuest:       map[int]Quest{},
			Specialty:      WeaponCrafter,
			Knowledge: []WorkItem{
				{
					Item: &item.Item{
						Name: item.Bow,
					},
					Cost: []item.Resources{
						{
							Item: item.Item{
								Name: item.Wood,
							},
							Qty: 10,
						},
					},
					WorkCost: 1,
				},
			},
		},
	}
}
