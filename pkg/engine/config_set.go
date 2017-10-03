// Copyright 2017, Pulumi Corporation.  All rights reserved.

package engine

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi/pkg/tokens"
	"github.com/pulumi/pulumi/pkg/util/contract"
)

func (eng *Engine) SetConfig(envName tokens.QName, key tokens.ModuleMember, value string) error {
	contract.Require(envName != tokens.QName(""), "envName")

	info, err := eng.planContextFromEnvironment(envName, "")
	if err != nil {
		return err
	}

	config := info.Target.Config
	if config == nil {
		config = make(map[tokens.ModuleMember]string)
		info.Target.Config = config
	}

	config[key] = value

	if err = eng.Environment.SaveEnvironment(info.Target, info.Snapshot); err != nil {
		return errors.Wrap(err, "could not save configuration value")
	}

	return nil
}

// ReplaceConfig sets the config for an environment to match `newConfig` and then saves
// the environment. Note that config values that were present in the old environment but are
// not present in `newConfig` will be removed from the environment
func (eng *Engine) ReplaceConfig(envName tokens.QName, newConfig map[tokens.ModuleMember]string) error {
	info, err := eng.planContextFromEnvironment(envName, "")
	if err != nil {
		return err
	}

	config := make(map[tokens.ModuleMember]string)
	for key, v := range newConfig {
		config[key] = v
	}

	info.Target.Config = config

	if err = eng.Environment.SaveEnvironment(info.Target, info.Snapshot); err != nil {
		return errors.Wrap(err, "could not save configuration value")
	}

	return nil
}
