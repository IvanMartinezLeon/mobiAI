package check

import (
	"fmt"
	"strings"

	"github.com/anomalyco/mobiAI/cli/internal/pi"
)

type Result struct {
	Name       string
	Passed     bool
	Detail     string
	Suggestion string
}

func (r *Result) String() string {
	status := "✓"
	if !r.Passed {
		status = "✗"
	}
	line := fmt.Sprintf("  %s %s", status, r.Name)
	if r.Detail != "" {
		line += fmt.Sprintf(" (%s)", r.Detail)
	}
	if !r.Passed && r.Suggestion != "" {
		line += fmt.Sprintf("\n     → %s", r.Suggestion)
	}
	return line
}

func RunAll() []Result {
	return []Result{
		checkNodeJS(),
		checkNpm(),
		checkPiInstalled(),
		checkTheme(),
		checkExtension(),
		checkAppendSystem(),
		checkSettings(),
	}
}

func RunEssential() []Result {
	return []Result{
		checkNodeJS(),
		checkPiInstalled(),
	}
}

func checkNodeJS() Result {
	r := Result{Name: "Node.js instalado"}
	ver, err := pi.CheckNodeJS()
	if err != nil {
		r.Detail = err.Error()
		r.Suggestion = "Instala Node.js desde https://nodejs.org/ o con brew install node"
		return r
	}
	r.Passed = true
	r.Detail = ver
	return r
}

func checkNpm() Result {
	r := Result{Name: "npm disponible"}
	ver, err := pi.CheckNpm()
	if err != nil {
		r.Detail = err.Error()
		r.Suggestion = "npm debería venir con Node.js. Reinstala Node.js."
		return r
	}
	r.Passed = true
	r.Detail = ver
	return r
}

func checkPiInstalled() Result {
	r := Result{Name: "Pi Coding Agent instalado"}
	ok, ver := pi.CheckPiInstalled()
	if !ok {
		r.Suggestion = "Ejecuta mobi install para instalarlo"
		return r
	}
	r.Passed = true
	r.Detail = ver
	return r
}

func checkTheme() Result {
	r := Result{Name: "Theme MOBI AI instalado"}
	themePath := pi.GetThemesDir() + "/mobi-theme.json"
	if !pi.FileExists(themePath) {
		r.Suggestion = "Ejecuta mobi install para copiar el theme"
		return r
	}
	r.Passed = true
	r.Detail = themePath
	return r
}

func checkExtension() Result {
	r := Result{Name: "Extensión mobi-header instalada"}
	extPath := pi.GetExtensionsDir() + "/mobi-header.ts"
	if !pi.FileExists(extPath) {
		r.Suggestion = "Ejecuta mobi install para copiar la extensión"
		return r
	}
	r.Passed = true
	r.Detail = extPath
	return r
}

func checkAppendSystem() Result {
	r := Result{Name: "APPEND_SYSTEM.md instalado"}
	path := pi.GetAppendSystemPath()
	if !pi.FileExists(path) {
		r.Suggestion = "Ejecuta mobi install para copiar APPEND_SYSTEM.md"
		return r
	}
	r.Passed = true
	r.Detail = path
	return r
}

func checkSettings() Result {
	r := Result{Name: "Configuración global correcta"}
	path := pi.GetSettingsPath()
	if !pi.FileExists(path) {
		r.Suggestion = "Ejecuta mobi install para crear la configuración"
		return r
	}
	data, err := pi.ReadFile(path)
	if err != nil {
		r.Detail = err.Error()
		r.Suggestion = "Revisa el archivo de configuración"
		return r
	}
	if !strings.Contains(data, `"theme": "mobi-theme"`) {
		r.Suggestion = "El tema no está configurado. Ejecuta mobi install."
		return r
	}
	r.Passed = true
	r.Detail = "theme: mobi-theme, quietStartup: " + fmt.Sprintf("%v", strings.Contains(data, `"quietStartup": true`))
	return r
}

func Summary(results []Result) string {
	var passed, failed int
	for _, r := range results {
		if r.Passed {
			passed++
		} else {
			failed++
		}
	}
	return fmt.Sprintf("\nResultado: %d correctos, %d fallos", passed, failed)
}
