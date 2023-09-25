/*
Validator APIs for SLSA Provenance v1 protos.
*/
package v1

import (
	"errors"
	"fmt"
)

// all of the following errors apply to SLSA Build L1 and above
var (
	ErrBuilderRequired         = errors.New("RunDetails.Builder required")
	ErrBuilderIdRequired       = errors.New("Builder.Id required")
	ErrBuildDefinitionRequired = errors.New("BuildeDefinition required")
	ErrBuildTypeRequired       = errors.New("BuildDefinition.BuildType required")
	ErrExternalParamsRequired  = errors.New("BuildDefinition.ExternalParameters required")
	ErrRunDetailsRequired      = errors.New("RunDetails required")
)

func (b *Builder) Validate() error {
	// the id field is required for SLSA Build L1
	if b.GetId() == "" {
		return ErrBuilderIdRequired
	}

	// check that all builderDependencies are valid RDs
	builderDeps := b.GetBuilderDependencies()
	if len(builderDeps) > 0 {
		for i, rd := range builderDeps {
			if err := rd.Validate(); err != nil {
				return fmt.Errorf("Invalid Builder.BuilderDependencies[%d]: %w", i, err)
			}
		}
	}

	return nil
}

func (b *BuildDefinition) Validate() error {
	// the buildType field is required for SLSA Build L1
	if b.GetBuildType() == "" {
		return ErrBuildTypeRequired
	}

	// the externalParameters field is required for SLSA Build L1
	if b.GetExternalParameters() == nil {
		return ErrExternalParamsRequired
	}

	// check that all resolvedDependencies are valid RDs
	resolvedDeps := b.GetResolvedDependencies()
	if len(resolvedDeps) > 0 {
		for i, rd := range resolvedDeps {
			if err := rd.Validate(); err != nil {
				return fmt.Errorf("Invalid BuildDefinition.ResolvedDependencies[%d]: %w", i, err)
			}
		}
	}

	return nil
}

func (r *RunDetails) Validate() error {
	// the builder field is required for SLSA Build L1
	builder := r.GetBuilder()
	if builder == nil {
		return ErrBuilderRequired
	}

	// check the Builder
	if err := builder.Validate(); err != nil {
		return err
	}

	// check that all byproducts are valid RDs
	byproducts := r.GetByproducts()
	if len(byproducts) > 0 {
		for i, rd := range byproducts {
			if err := rd.Validate(); err != nil {
				return fmt.Errorf("Invalid RunDetails.Byproducts[%d]: %w", i, err)
			}
		}
	}

	return nil
}

func (p *Provenance) Validate() error {
	// the buildDefinition field is required for SLSA Build L1
	buildDef := p.GetBuildDefinition()
	if buildDef == nil {
		return ErrBuildDefinitionRequired
	}

	// check the BuildDefinition
	if err := buildDef.Validate(); err != nil {
		return err
	}

	// the runDetails field is required for SLSA Build L1
	runDetails := p.GetRunDetails()
	if runDetails == nil {
		return ErrRunDetailsRequired
	}

	// check the RunDetails
	if err := runDetails.Validate(); err != nil {
		return err
	}

	return nil
}
