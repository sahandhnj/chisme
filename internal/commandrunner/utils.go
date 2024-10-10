package commandrunner

import "fmt"

// applyCommandRootElevation applies the root elevation to the command by adding sudo
func applyCommandRootElevation(command *string, askPassPath string) {
	switch {
	case askPassPath != "":
		*command = fmt.Sprintf("SUDO_ASKPASS=%s sudo -A %s", askPassPath, *command)
	default:
		*command = fmt.Sprintf("sudo -S %s", *command)
	}
}
