package detection

var LOCKS = map[string]string{
	"bun.lockb":           "bun",
	"bun.lock":            "bun",
	"deno.lock":           "deno",
	"package-lock.json":   "npm",
	"pnpm-workspace.yaml": "pnpm",
	"pnpm-lock.yaml":      "pnpm",
	"yarn.lock":           "yarn",
	"npm-shrinkwrap.json": "npm",
}

var INSTALL_METADATA = map[string]string{
	"node_modules/.deno/":             "deno",
	"node_modules/.pnpm/":             "pnpm",
	"node_modules/.yarn-state.yml":    "yarn",
	"node_modules/.yarn_integrity":    "yarn",
	"node_modules/.package-lock.json": "npm",
	".pnp.cjs":                        "yarn",
	".pnp.js":                         "yarn",
	"bun.lock":                        "bun",
	"bun.lockb":                       "bun",
}

var INSTALL_PACKAGE = map[string]string{
	"bun":        "https://bun.sh",
	"deno":       "https://deno.com",
	"pnpm":       "https://pnpm.io/installation",
	"pnpm@6":     "https://pnpm.io/6.x/installation",
	"yarn":       "https://classic.yarnpkg.com/en/docs/install",
	"yarn@berry": "https://yarnpkg.com/getting-started/install",
	"npm":        "https://docs.npmjs.com/cli/configuring-npm/install",
}

type AgentName string

const AgentNameNpm = "npm"
const AgentNameYarn = "yarn"
const AgentNamePnpm = "pnpm"
const AgentNameBun = "bun"
const AgentNameDeno = "deno"

type Agent string

const AgentNpm = "npm"
const AgentYarn = "yarn"
const AgentYarnBerry = "yarn@berry"
const AgentPnpm = "pnpm"
const AgentPnpm6 = "pnpm@6"
const AgentBun = "bun"
const AgentDeno = "deno"

var Agents = []string{
	AgentNpm,
	AgentYarn,
	AgentYarnBerry,
	AgentPnpm,
	AgentPnpm6,
	AgentBun,
	AgentDeno,
}

func FromStringToAgentName(s string) AgentName {
	switch s {
	case "npm":
		return AgentNameNpm
	case "yarn":
		return AgentNameYarn
	case "pnpm":
		return AgentNamePnpm
	case "bun":
		return AgentNameBun
	case "deno":
		return AgentNameDeno
	default:
		panic("invalid agent name")
	}
}
