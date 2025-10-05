package commandRun

import (
	"flag"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type FlagArguments struct {
	RuntimeVersion        *string
	PackageManagerVersion *string
	Env                   bool
}

func ParseFlags() FlagArguments {
	flag.String("runtimeVersion", "", "runtimeVersion flag")
	flag.String("packageManagerVersion", "", "runtimeVersion flag")
	flag.Bool("gnpmEnv", false, "env flag for printing environment")

	pflag.CommandLine.ParseErrorsWhitelist.UnknownFlags = true
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	var flagArguments = FlagArguments{}

	var parsedRuntimeVersion = viper.GetString("runtimeVersion")
	var parsedPackageManagerVersion = viper.GetString("packageManagerVersion")
	flagArguments.Env = viper.GetBool("gnpmEnv")

	if parsedRuntimeVersion != "" {
		flagArguments.RuntimeVersion = &parsedRuntimeVersion
	}
	if parsedPackageManagerVersion != "" {
		flagArguments.PackageManagerVersion = &parsedPackageManagerVersion
	}

	return flagArguments
}

var knownFlags = []string{
	"--runtimeVersion", "--packageManagerVersion",
	"-runtimeVersion", "-packageManagerVersion",
	"--gnpmEnv", "-gnpmEnv",
}

func FilterArgs(args []string) []string {
	var filtered []string
	skipNext := false
	for i, arg := range args {
		if skipNext {
			skipNext = false
			continue
		}
		isFlag := false
		for _, flag := range knownFlags {
			if arg == flag || strings.HasPrefix(arg, flag+"=") {
				isFlag = true
				if !strings.Contains(arg, "=") && i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
					skipNext = true
				}
				break
			}
		}
		if !isFlag {
			filtered = append(filtered, arg)
		}
	}
	return filtered
}
