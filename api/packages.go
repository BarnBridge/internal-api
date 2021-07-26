package api

import (
	"github.com/barnbridge/internal-api/governance"
	"github.com/barnbridge/internal-api/smartexposure"
)

func (a *API) registerPackages() {
	a.packages = append(a.packages, governance.New(a.db))
	a.packages = append(a.packages, smartexposure.New(a.db))
}
