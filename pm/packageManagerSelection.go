package pm

import (
	"github.com/samtv12345/gnpm/pm/impl/npm"
	"github.com/samtv12345/gnpm/pm/impl/pnpm"
	"github.com/samtv12345/gnpm/pm/impl/yarnClassic"
	"github.com/samtv12345/gnpm/pm/interfaces"
	"go.uber.org/zap"
)

func GetPackageManagerSelection(pm string, logger *zap.SugaredLogger) interfaces.IPackageManager {
	if pm == "npm" {
		return npm.Npm{
			Logger: logger,
		}
	} else if pm == "pnpm" {
		return pnpm.Pnpm{
			Logger: logger,
		}
	} else if pm == "yarn@classic" {
		return yarnClassic.Yarn{
			Logger: logger,
		}
	} else if pm == "yarn@berry" {
		panic("Not implemented yet")

	} else if pm == "bun" {
		panic("Not implemented yet")
	} else {
		panic("Not implemented yet")
	}
}
