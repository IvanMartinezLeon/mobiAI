package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/IvanMartinezLeon/mobiAI/cli/internal/check"
	"github.com/IvanMartinezLeon/mobiAI/cli/internal/config"
	"github.com/IvanMartinezLeon/mobiAI/cli/internal/pi"
	"github.com/spf13/cobra"
)

var installForce bool
var installNoNpm bool

func init() {
	installCmd.Flags().BoolVarP(&installForce, "force", "f", false, "Re-instalar aunque ya exista")
	installCmd.Flags().BoolVar(&installNoNpm, "no-npm", false, "Saltar instalación npm (solo copiar archivos)")
	rootCmd.AddCommand(installCmd)
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Instalar/configurar Pi con personalización MOBI AI",
	Long: `Instala el Pi Coding Agent globalmente y copia la personalización
MOBI AI (tema, extensión, APPEND_SYSTEM.md y configuración).

Flags:
  --force     Re-instala npm y sobreescribe archivos existentes
  --no-npm    Solo copia los archivos de personalización (si Pi ya está instalado)
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(`
  ███╗   ███╗ ██████╗ ██████╗ ██╗     █████╗ ██╗
  ████╗ ████║██╔═══██╗██╔══██╗██║    ██╔══██╗██║
  ██╔████╔██║██║   ██║██████╔╝██║    ███████║██║
  ██║╚██╔╝██║██║   ██║██╔══██╗██║    ██╔══██║██║
  ██║ ╚═╝ ██║╚██████╔╝██████╔╝██║    ██║  ██║██║
  ╚═╝     ╚═╝ ╚═════╝ ╚═════╝ ╚═╝    ╚═╝  ╚═╝╚═╝
  === Instalación MOBI AI ===`)

		repoRoot := detectRepoRoot()

		if !installNoNpm {
			if err := installPiAndPackages(); err != nil {
				return err
			}
		} else {
			fmt.Println("  → Saltando instalación npm (--no-npm)")
		}

		fmt.Println("\n📦 Copiando archivos de personalización...")
		ensureDirs := []string{pi.GetThemesDir(), pi.GetExtensionsDir()}
		for _, d := range ensureDirs {
			if err := pi.EnsureDirs(d); err != nil {
				return err
			}
		}

		type fileCopy struct {
			src string
			dst string
			label string
		}

		copies := []fileCopy{
			{filepath.Join(repoRoot, ".pi", "themes", "mobi-theme.json"), filepath.Join(pi.GetThemesDir(), "mobi-theme.json"), "Theme"},
			{filepath.Join(repoRoot, ".pi", "extensions", "mobi-header.ts"), filepath.Join(pi.GetExtensionsDir(), "mobi-header.ts"), "Extensión header"},
			{filepath.Join(repoRoot, ".pi", "APPEND_SYSTEM.md"), pi.GetAppendSystemPath(), "APPEND_SYSTEM.md"},
		}

		for _, c := range copies {
			if err := pi.CopyFile(c.src, c.dst); err != nil {
				return fmt.Errorf("error copiando %s: %s", c.label, err.Error())
			}
			fmt.Printf("  ✓ %s\n", c.label)
		}

		if err := mergeSettings(repoRoot); err != nil {
			return err
		}

		if user := pi.SudoUser(); user != "" {
			pi.FixPermissions(pi.HomeDir() + "/.pi")
		}

		fmt.Println("\n✅ Instalación completada correctamente")
		fmt.Println("\nComandos útiles:")
		fmt.Println("  mobi doctor    — Verificar que todo está correcto")
		fmt.Println("  mobi status    — Mostrar estado de la configuración")
		return nil
	},
}

func detectRepoRoot() string {
	cwd, _ := os.Getwd()
	if pi.FileExists(filepath.Join(cwd, ".pi")) {
		return cwd
	}
	if pi.FileExists(filepath.Join(cwd, "..", ".pi")) {
		return filepath.Dir(cwd)
	}
	return cwd
}

func installPiAndPackages() error {
	fmt.Println("\n🔍 Verificando requisitos...")
	results := check.RunEssential()
	for _, r := range results {
		if !r.Passed {
			fmt.Println(r.String())
			return fmt.Errorf("%s", r.Suggestion)
		}
		fmt.Println(r.String())
	}

	packages := []string{
		"npm:pi-subagents",
		"npm:pi-mcp-adapter",
		"npm:context-mode",
	}

	fmt.Println("\n📦 Instalando paquetes...")
	if err := pi.InstallPi(); err != nil {
		return err
	}
	if err := pi.InstallPackages(packages); err != nil {
		return err
	}
	return nil
}

func mergeSettings(repoRoot string) error {
	localPath := filepath.Join(repoRoot, ".pi", "settings.json")
	globalPath := pi.GetSettingsPath()

	var global *config.Settings
	if pi.FileExists(globalPath) {
		data, err := os.ReadFile(globalPath)
		if err == nil {
			global, _ = config.Parse(data)
		}
	}

	if pi.FileExists(localPath) {
		data, err := os.ReadFile(localPath)
		if err != nil {
			return fmt.Errorf("error leyendo settings locales: %s", err.Error())
		}
		local, err := config.Parse(data)
		if err != nil {
			return fmt.Errorf("error parseando settings locales: %s", err.Error())
		}
		merged := config.Merge(global, local)
		out, err := merged.Marshal()
		if err != nil {
			return fmt.Errorf("error serializando settings: %s", err.Error())
		}
		if err := os.WriteFile(globalPath, out, 0644); err != nil {
			return fmt.Errorf("error escribiendo settings globales: %s", err.Error())
		}
		fmt.Println("  ✓ Configuración (settings.json) actualizada")
	} else if global == nil {
		s := &config.Settings{Theme: "mobi-theme"}
		qt := true
		s.QuietStartup = &qt
		out, _ := s.Marshal()
		os.WriteFile(globalPath, out, 0644)
		fmt.Println("  ✓ Configuración (settings.json) creada")
	} else {
		fmt.Println("  ✓ Configuración (settings.json) ya existe, sin cambios")
	}
	return nil
}
