package modules

import (
	"github.com/Sirupsen/logrus"
	"github.com/etcinit/phabulous/app/interfaces"
	"github.com/etcinit/phabulous/app/modules/core"
	"github.com/etcinit/phabulous/app/modules/dev"
	"github.com/etcinit/phabulous/app/modules/extension"
	"github.com/jacobstr/confer"
)

// ModuleFactory provides facilities for building sets of modules to be used
// in the chat server.
type ModuleFactory struct {
	config *confer.Config
	logger *logrus.Logger
}

// NewModuleFactory constructs a new instance of a ModuleFactory.
func NewModuleFactory(
	config *confer.Config,
	logger *logrus.Logger,
) *ModuleFactory {
	return &ModuleFactory{
		config: config,
		logger: logger,
	}
}

// Make loads module configuration and initializes them.
func (f *ModuleFactory) Make() []interfaces.Module {
	modules := []interfaces.Module{}
	moduleMap := map[string]interfaces.Module{}
	modulesToLoad := f.config.GetStringSlice("server.modules")

	// We go over the configured slice of modules to load, attempt to find a
	// module with that name, and initialize it.
	for _, moduleName := range modulesToLoad {
		// We won't load the same module more than once.
		if _, ok := moduleMap[moduleName]; ok {
			continue
		}

		if moduleName == "core" {
			moduleMap["core"] = &core.Module{}
		} else if moduleName == "dev" {
			moduleMap["dev"] = &dev.Module{}
		} else if moduleName == "extension" {
			moduleMap["extension"] = &extension.Module{}
		} else {
			f.logger.Panicf(
				"Unable to load modules. Unknown module '%s'.",
				moduleName,
			)
		}
	}

	// We extract all the values of the map into a slice.
	for _, module := range moduleMap {
		modules = append(modules, module)
	}

	return modules
}
