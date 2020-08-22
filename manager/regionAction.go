package manager

import "github.com/Selsynn/cbepbackend/communication"

type RegionAction struct {
	Description string
	Callback    func(communication.ActionToManager) *communication.ActionFromManager
}

func(ra *RegionAction) ThisIsARegionAction(){

}
