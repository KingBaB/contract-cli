package cli

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

func defaultSaveFileDialog(ctx context.Context, suggestedName string) (string, error) {
	suggestedName = strings.TrimSpace(suggestedName)
	if suggestedName == "" {
		suggestedName = "download"
	}

	switch runtime.GOOS {
	case "darwin":
		return runSaveDialogCommand(ctx, "osascript",
			"-e", `set targetFile to choose file name with prompt "Save contract file as:" default name `+appleScriptString(suggestedName),
			"-e", "POSIX path of targetFile",
		)
	case "windows":
		script := strings.Join([]string{
			"Add-Type -AssemblyName System.Windows.Forms",
			"$dialog = New-Object System.Windows.Forms.SaveFileDialog",
			"$dialog.FileName = " + powershellString(suggestedName),
			"if ($dialog.ShowDialog() -eq [System.Windows.Forms.DialogResult]::OK) { $dialog.FileName } else { exit 2 }",
		}, "; ")
		return runSaveDialogCommand(ctx, "powershell", "-NoProfile", "-Command", script)
	default:
		if path, err := exec.LookPath("zenity"); err == nil {
			return runSaveDialogCommand(ctx, path, "--file-selection", "--save", "--confirm-overwrite", "--filename", suggestedName)
		}
		if path, err := exec.LookPath("kdialog"); err == nil {
			return runSaveDialogCommand(ctx, path, "--getsavefilename", suggestedName)
		}
		return "", fmt.Errorf("no supported save dialog command found")
	}
}

func runSaveDialogCommand(ctx context.Context, name string, args ...string) (string, error) {
	output, err := exec.CommandContext(ctx, name, args...).Output()
	if err != nil {
		return "", err
	}
	path := strings.TrimSpace(string(output))
	if path == "" {
		return "", fmt.Errorf("save dialog returned empty path")
	}
	return path, nil
}

func appleScriptString(value string) string {
	return `"` + strings.ReplaceAll(value, `"`, `\"`) + `"`
}

func powershellString(value string) string {
	return `'` + strings.ReplaceAll(value, `'`, `''`) + `'`
}
