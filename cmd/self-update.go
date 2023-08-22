package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/blang/semver"
	"github.com/martient/bifrost-env-manager/pkg/utils"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
)

func doConfirmAndSelfUpdate() {
	latest, found, err := selfupdate.DetectLatest("martient/bifrost-env-manager")
	if err != nil {
		utils.LogError("Error occurred while detecting version:", err, "Updater")
		return
	}

	v := semver.MustParse(BEMversion)
	if !found || latest.Version.LTE(v) {
		utils.LogInfo("Current version is the latest", "", "Updater")
		return
	}

	if b, _ := rootCmd.Flags().GetBool("yes"); !b {
		fmt.Print("Do you want to update to", latest.Version, "? (y/n): ")
		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil || (input != "y\n" && input != "n\n") {
			utils.LogError("Invalid input %s", err, "Updater")
			return
		}
		if input == "n\n" {
			return
		}
	}

	exe, err := os.Executable()
	if err != nil {
		utils.LogError("Could not locate executable path: %s", err, "Updater")
		return
	}
	if err := selfupdate.UpdateTo(latest.AssetURL, exe); err != nil {
		utils.LogError("Error occurred while updating binary: %s", err, "Updater")
		return
	}
	utils.LogInfo("Successfully updated to version %s", latest.Version, "Updater")
}
