package cli

import (
	"os"
	"path/filepath"

	"cn.qfei/contract-cli/internal/build"
	"cn.qfei/contract-cli/internal/skillsync"
)

type skillCLIIdentity struct {
	Version    string `json:"version"`
	Executable string `json:"executable"`
	Error      string `json:"error,omitempty"`
}

type skillJSONMetadata struct {
	Directory   string `json:"directory"`
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
}

type skillRegistrationLocation struct {
	Target  string   `json:"target"`
	Version string   `json:"version"`
	Status  string   `json:"status"`
	Managed []string `json:"managed_skills"`
	Error   string   `json:"error,omitempty"`
}

type skillRegistrationSnapshot struct {
	Path      string                      `json:"path"`
	Confirmed bool                        `json:"confirmed"`
	Locations []skillRegistrationLocation `json:"locations"`
	ReadError string                      `json:"read_error,omitempty"`
}

type skillRecoveryDetails struct {
	Status      string `json:"status"`
	Description string `json:"description"`
}

type skillFailureDetails struct {
	CLI          skillCLIIdentity          `json:"cli"`
	Target       string                    `json:"target"`
	Registration skillRegistrationSnapshot `json:"registration"`
	Recovery     skillRecoveryDetails      `json:"recovery"`
	Guidance     string                    `json:"guidance"`
}

type skillCommandFailure struct {
	cause   error
	details skillFailureDetails
}

func (failure *skillCommandFailure) Error() string              { return failure.cause.Error() }
func (failure *skillCommandFailure) Unwrap() error              { return failure.cause }
func (failure *skillCommandFailure) commandFailureDetails() any { return failure.details }

func currentSkillCLIIdentity() skillCLIIdentity {
	identity := skillCLIIdentity{Version: build.Current().Version}
	executable, err := os.Executable()
	if err != nil {
		identity.Error = err.Error()
		return identity
	}
	identity.Executable = executable
	return identity
}

func readSkillRegistrationSnapshot(service *skillsync.Service, confirmed bool) skillRegistrationSnapshot {
	snapshot := skillRegistrationSnapshot{Path: filepath.Join(service.StateDir, "skill-locations.json"), Confirmed: confirmed, Locations: []skillRegistrationLocation{}}
	locations, err := service.RegisteredLocations()
	if err != nil {
		snapshot.ReadError = err.Error()
		snapshot.Locations = nil
		return snapshot
	}
	for _, location := range locations {
		snapshot.Locations = append(snapshot.Locations, skillRegistrationLocation{
			Target: location.Target, Version: location.Version, Status: location.Status, Managed: append([]string{}, location.Managed...), Error: location.Error,
		})
	}
	return snapshot
}

func (a *App) renderSkillsJSONList(skills []skillMetadata) error {
	metadata := make([]skillJSONMetadata, 0, len(skills))
	version := build.Current().Version
	for _, skill := range skills {
		metadata = append(metadata, skillJSONMetadata{Directory: skill.Dir, Name: skill.Name, Version: version, Description: skill.Description})
	}
	return a.renderJSONCommandResult(struct {
		CLI    skillCLIIdentity    `json:"cli"`
		Skills []skillJSONMetadata `json:"skills"`
	}{CLI: currentSkillCLIIdentity(), Skills: metadata})
}

func (a *App) renderSkillsJSONInstall(service *skillsync.Service, target string, result skillsync.InstallResult) error {
	return a.renderJSONCommandResult(struct {
		CLI                  skillCLIIdentity          `json:"cli"`
		Target               string                    `json:"target"`
		Installed            []string                  `json:"installed"`
		Skipped              []string                  `json:"skipped"`
		Registration         skillRegistrationSnapshot `json:"registration"`
		ClientReloadRequired bool                      `json:"client_reload_required"`
	}{
		CLI: currentSkillCLIIdentity(), Target: target, Installed: append([]string{}, result.Installed...), Skipped: append([]string{}, result.Skipped...),
		Registration: readSkillRegistrationSnapshot(service, result.Confirmed), ClientReloadRequired: true,
	})
}

func skillJSONFailure(err error, service *skillsync.Service, target string) error {
	recovery := skillsync.RecoveryStatus(err)
	return &skillCommandFailure{cause: err, details: skillFailureDetails{
		CLI: currentSkillCLIIdentity(), Target: target, Registration: readSkillRegistrationSnapshot(service, false),
		Recovery: skillRecoveryDetails{Status: recovery, Description: recoveryDescription(recovery)}, Guidance: skillFailureGuidance(target, recovery),
	}}
}
