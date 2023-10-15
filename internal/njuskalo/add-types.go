package njuskalo

import "github.com/sycomancy/glasnik/internal/types"

type PropertyAdd struct {
	PriceEur string
	Location string
}

func (p *PropertyAdd) Process(e types.AdEntry) PropertyAdd {
	props := PropertyAdd{}
	description := e.Description
}
