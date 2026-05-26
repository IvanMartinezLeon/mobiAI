package pi

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func homeDir() string {
	h, _ := os.UserHomeDir()
	return h
}

func GetPiAgentDir() string {
	return filepath.Join(homeDir(), ".pi", "agent")
}

func GetThemesDir() string {
	return filepath.Join(GetPiAgentDir(), "themes")
}

func GetExtensionsDir() string {
	return filepath.Join(GetPiAgentDir(), "extensions")
}

func GetSettingsPath() string {
	return filepath.Join(GetPiAgentDir(), "settings.json")
}

func GetAppendSystemPath() string {
	return filepath.Join(GetPiAgentDir(), "APPEND_SYSTEM.md")
}

func runCmd(name string, args ...string) (string, error) {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command(name, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Env = os.Environ()
	err := cmd.Run()
	out := strings.TrimSpace(stdout.String())
	if out == "" {
		out = strings.TrimSpace(stderr.String())
	}
	if err != nil {
		return out, fmt.Errorf("%s: %s", err.Error(), strings.TrimSpace(stderr.String()))
	}
	return out, nil
}

func CheckNodeJS() (string, error) {
	return runCmd("node", "--version")
}

func CheckNpm() (string, error) {
	return runCmd("npm", "--version")
}

func CheckPiInstalled() (bool, string) {
	if ver, err := runCmd("pi", "--version"); err == nil && ver != "" {
		return true, ver
	}
	if ver, err := runCmd("npx", "-y", "@earendil-works/pi-coding-agent", "--version"); err == nil && ver != "" {
		return true, ver
	}
	piPath := findPiBinary()
	if piPath != "" {
		if ver, err := runCmd(piPath, "--version"); err == nil && ver != "" {
			return true, ver
		}
	}
	return false, ""
}

func findPiBinary() string {
	home, _ := os.UserHomeDir()
	candidates := []string{
		filepath.Join(home, "npx"),
		"/opt/homebrew/bin/pi",
		"/usr/local/bin/pi",
		"/usr/bin/pi",
		filepath.Join(home, ".npm-global", "bin", "pi"),
		filepath.Join(home, ".nvm", "current", "bin", "pi"),
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}

func InstallPi() error {
	fmt.Println("  → Instalando @earendil-works/pi-coding-agent...")
	out, err := runCmd("npm", "install", "-g", "--ignore-scripts", "@earendil-works/pi-coding-agent")
	if err != nil {
		return fmt.Errorf("error al instalar Pi: %s\n%s", err.Error(), out)
	}
	fmt.Println("  ✓ Pi Coding Agent instalado globalmente")
	return nil
}

func InstallPackage(pkg string) error {
	fmt.Printf("  → Instalando %s...\n", pkg)
	_, err := runCmd("pi", "install", pkg)
	if err != nil {
		return fmt.Errorf("error al instalar %s: %s", pkg, err.Error())
	}
	fmt.Printf("  ✓ %s instalado\n", pkg)
	return nil
}

func InstallPackages(packages []string) error {
	for _, pkg := range packages {
		if err := InstallPackage(pkg); err != nil {
			return err
		}
	}
	return nil
}

func EnsureDirs(dirs ...string) error {
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			return fmt.Errorf("error al crear directorio %s: %s", d, err.Error())
		}
	}
	return nil
}

func CopyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("error al leer %s: %s", src, err.Error())
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return fmt.Errorf("error al crear directorio %s: %s", filepath.Dir(dst), err.Error())
	}
	if err := os.WriteFile(dst, data, 0644); err != nil {
		return fmt.Errorf("error al escribir %s: %s", dst, err.Error())
	}
	return nil
}

func HomeDir() string {
	return homeDir()
}

func SudoUser() string {
	user := os.Getenv("SUDO_USER")
	if user != "" && user != "root" {
		return user
	}
	return ""
}

func FixPermissions(path string) error {
	user := SudoUser()
	if user == "" {
		return nil
	}
	_, err := runCmd("chown", "-R", user, path)
	return err
}

func IsAdmin() bool {
	if runtime.GOOS == "windows" {
		_, err := exec.LookPath("net")
		if err != nil {
			return false
		}
		out, err := runCmd("net", "session")
		return err == nil && !strings.Contains(out, "Access is denied")
	}
	return os.Geteuid() == 0
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func ReadFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
