package cli

import (
	"fmt"

	"cn.qfei/contract-cli/internal/build"
)

func (a *App) runVersion(args []string) (err error) {
	a.logger.Info("version command started")
	defer func() {
		if err != nil {
			a.logger.Error("version command failed", "error_type", fmt.Sprintf("%T", err))
		}
	}()
	parsed, err := parseArgs(args, map[string]struct{}{"--output": {}}, nil)
	if err != nil {
		return err
	}
	if len(parsed.positionals) != 0 {
		return fmt.Errorf("version does not accept positional arguments")
	}
	jsonOutput, err := parseJSONOutputFormat(parsed.String("--output"))
	if err != nil {
		return err
	}
	if jsonOutput {
		current := build.Current()
		return a.renderJSONCommandResult(map[string]string{"name": current.Name, "version": current.Version,
			"commit": current.Commit, "built_at": current.Date, "feature_baseline": build.FeatureBaseline()})
	}
	_, err = fmt.Fprintln(a.stdout, build.Current().String())
	return err
}
