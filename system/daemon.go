package system

import (
	"io/ioutil"
	"os/exec"
	"regexp"
	"strings"
)

// SystemctlEnable call systemctl enable on thing.
func SystemctlEnable(thing string) error {
	if out, err := exec.Command("systemctl", "enable", thing).CombinedOutput(); err != nil {
		return ErrorLog("%v - Failed to call systemctl enable on %s - %s", err, thing, string(out))
	}
	return nil
}

// SystemctlDisable call systemctl disable on thing.
func SystemctlDisable(thing string) error {
	if out, err := exec.Command("systemctl", "disable", thing).CombinedOutput(); err != nil {
		return ErrorLog("%v - Failed to call systemctl disable on %s - %s", err, thing, string(out))
	}
	return nil
}

// SystemctlRestart call systemctl restart on thing.
func SystemctlRestart(thing string) error {
	if IsSystemRunning() {
		if out, err := exec.Command("systemctl", "restart", thing).CombinedOutput(); err != nil {
			return ErrorLog("%v - Failed to call systemctl restart on %s - %s", err, thing, string(out))
		}
	}
	return nil
}

// SystemctlStart call systemctl start on thing.
func SystemctlStart(thing string) error {
	if IsSystemRunning() {
		if out, err := exec.Command("systemctl", "start", thing).CombinedOutput(); err != nil {
			return ErrorLog("%v - Failed to call systemctl start on %s - %s", err, thing, string(out))
		}
	}
	return nil
}

// SystemctlStop call systemctl stop on thing.
func SystemctlStop(thing string) error {
	if IsSystemRunning() {
		if out, err := exec.Command("systemctl", "stop", thing).CombinedOutput(); err != nil {
			return ErrorLog("%v - Failed to call systemctl stop on %s - %s", err, thing, string(out))
		}
	}
	return nil
}

// SystemctlEnableStart call systemctl enable and then systemctl start on thing.
func SystemctlEnableStart(thing string) error {
	if err := SystemctlEnable(thing); err != nil {
		return err
	}
	err := SystemctlStart(thing)
	return err
}

// SystemctlDisableStop call systemctl disable and then systemctl stop on thing.
// Panic on error.
func SystemctlDisableStop(thing string) error {
	if err := SystemctlDisable(thing); err != nil {
		return err
	}
	err := SystemctlStop(thing)
	return err
}

// SystemctlIsRunning return true only if systemctl suggests that the thing is
// running.
func SystemctlIsRunning(thing string) bool {
	if _, err := exec.Command("systemctl", "is-active", thing).CombinedOutput(); err == nil {
		return true
	}
	return false
}

// IsSystemRunning returns true, if 'is-system-running' reports 'running'
// or 'starting'. In all other cases it returns false, which means: do not
// call 'start' or 'restart' to prevent 'Transaction is destructive' messages
func IsSystemRunning() bool {
	match := false
	out, err := exec.Command("/usr/bin/systemctl", "is-system-running").CombinedOutput()
	DebugLog("IsSystemRunning - /usr/bin/systemctl is-system-running : '%+v %s'", err, string(out))
	for _, line := range strings.Split(string(out), "\n") {
		if strings.TrimSpace(line) == "starting" || strings.TrimSpace(line) == "running" || strings.TrimSpace(line) == "degraded" {
			DebugLog("IsSystemRunning - system is degraded/starting/running, match true")
			match = true
			break
		}
	}
	return match
}

// WriteTunedAdmProfile write new profile to tuned, used instead of sometimes
// unreliable 'tuned-adm' command
func WriteTunedAdmProfile(profileName string) error {
	err := ioutil.WriteFile("/etc/tuned/active_profile", []byte(profileName), 0644)
	if err != nil {
		return ErrorLog("Failed to write tuned profile '%s' to '%s': %v", profileName, "/etc/tuned/active_profile", err)
	}
	return nil
}

// GetTunedProfile returns the currently active tuned profile by reading the
// file /etc/tuned/active_profile
// may be unreliable in newer tuned versions, so better use 'tuned-adm active'
// Return empty string if it cannot be determined.
func GetTunedProfile() string {
	content, err := ioutil.ReadFile("/etc/tuned/active_profile")
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(content))
}

// TunedAdmOff calls tuned-adm to switch off the active profile.
func TunedAdmOff() error {
	if out, err := exec.Command("tuned-adm", "off").CombinedOutput(); err != nil {
		return ErrorLog("Failed to call tuned-adm to switch off the active profile - %v %s", err, string(out))
	}
	return nil
}

// TunedAdmProfile calls tuned-adm to switch to the specified profile.
// newer versions of tuned seems to be reliable with this command and they
// changed the behaviour/handling of the file /etc/tuned/active_profile
func TunedAdmProfile(profileName string) error {
	if out, err := exec.Command("tuned-adm", "profile", profileName).CombinedOutput(); err != nil {
		return ErrorLog("Failed to call tuned-adm to active profile %s - %v %s", profileName, err, string(out))
	}
	return nil
}

// GetTunedAdmProfile return the currently active tuned profile.
// Return empty string if it cannot be determined.
func GetTunedAdmProfile() string {
	out, err := exec.Command("tuned-adm", "active").CombinedOutput()
	if err != nil {
		_ = ErrorLog("Failed to call tuned-adm to get the active profile - %v %s", err, string(out))
		return ""
	}
	re := regexp.MustCompile(`Current active profile: ([\w-]+)`)
	matches := re.FindStringSubmatch(string(out))
	if len(matches) == 0 {
		return ""
	}
	return matches[1]
}
