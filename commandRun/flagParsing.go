package commandRun

import (
	"flag"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type FlagArguments struct {
	RuntimeVersion        *string
	PackageManagerVersion *string
}

func ParseFlags() FlagArguments {
	flag.String("runtimeVersion", "", "runtimeVersion flag")
	flag.String("packageManagerVersion", "", "runtimeVersion flag")

	pflag.CommandLine.ParseErrorsWhitelist.UnknownFlags = true
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	var parsedRuntimeVersion = viper.GetString("runtimeVersion")
	var parsedPackageManagerVersion = viper.GetString("packageManagerVersion")

	var flagArguments = FlagArguments{}

	if parsedRuntimeVersion != "" {
		flagArguments.RuntimeVersion = &parsedRuntimeVersion
	}
	if parsedPackageManagerVersion != "" {
		flagArguments.PackageManagerVersion = &parsedPackageManagerVersion
	}

	return flagArguments
}
