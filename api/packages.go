package api

import (
	"github.com/barnbridge/internal-api/governance"
	"github.com/barnbridge/internal-api/smartyield"
	"github.com/barnbridge/internal-api/yieldfarming"
)

func (a *API) registerPackages() {
	a.packages = append(a.packages, governance.New(a.db))
	a.packages = append(a.packages, yieldfarming.New(a.db))
	a.packages = append(a.packages, smartyield.New(a.db))
}
