package town

import "github.com/Selsynn/cbepbackend/business/command"

type (
	RegionID    string
	RegionLevel int
	Region      struct {
		Name    RegionID
		Level   RegionLevel
		Command command.ID
	}
)

const (
	Forest          RegionID = "Forest"
	EnchantedForest RegionID = "Enchanted Forest"
)

func (t *Town) NewRegionLevel(rid RegionID, cid command.ID, newLevel RegionLevel) {
	var found *Region
	for _, r := range t.Regions {
		if r.Name == rid && r.Level >= newLevel {
			return
		}
		if r.Name == rid {
			found = r
		}
	}
	if found == nil {
		found = &Region{
			Name:    rid,
			Command: cid,
		}
	}
	found.Level = newLevel
}
