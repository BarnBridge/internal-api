package api

import "github.com/barnbridge/internal-api/governance"

func (a *API) registerPackages() {
	a.packages = append(a.packages, governance.New(a.db))
}
