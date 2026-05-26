package cmd

import (
	"fmt"
	"strings"

	"github.com/IvanMartinezLeon/mobiAI/cli/internal/config"
	"github.com/IvanMartinezLeon/mobiAI/cli/internal/pi"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(statusCmd)
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Mostrar estado actual de la configuración",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("\n📊 Estado MOBI AI")
		fmt.Println("=================")

		piVer := "—"
		if ok, ver := pi.CheckPiInstalled(); ok {
			piVer = ver
		}
		fmt.Printf("  Pi Coding Agent: %s\n", piVer)

		themePath := pi.GetThemesDir() + "/mobi-theme.json"
		if pi.FileExists(themePath) {
			fmt.Printf("  Theme: ✓ mobi-theme\n")
		} else {
			fmt.Printf("  Theme: ✗ no instalado\n")
		}

		extPath := pi.GetExtensionsDir() + "/mobi-header.ts"
		if pi.FileExists(extPath) {
			content, _ := pi.ReadFile(extPath)
			lines := strings.Split(content, "\n")
			hasFooter := strings.Contains(content, "setFooter")
			hasTokens := strings.Contains(content, "totalInput")
			extras := []string{}
			if hasFooter {
				extras = append(extras, "footer")
			}
			if hasTokens {
				extras = append(extras, "tokens")
			}
			detail := ""
			if len(extras) > 0 {
				detail = " (" + strings.Join(extras, ", ") + ")"
			}
			fmt.Printf("  Extensión: ✓ mobi-header.ts (%d líneas%s)\n", len(lines), detail)
		} else {
			fmt.Printf("  Extensión: ✗ no instalada\n")
		}

		appendPath := pi.GetAppendSystemPath()
		if pi.FileExists(appendPath) {
			fmt.Printf("  APPEND_SYSTEM.md: ✓ presente\n")
		} else {
			fmt.Printf("  APPEND_SYSTEM.md: ✗ no instalado\n")
		}

		settingsPath := pi.GetSettingsPath()
		if pi.FileExists(settingsPath) {
			data, _ := pi.ReadFile(settingsPath)
			s, err := config.Parse([]byte(data))
			if err == nil {
				quietStr := "false"
				if s.QuietStartup != nil && *s.QuietStartup {
					quietStr = "true"
				}
				fmt.Printf("  Configuración:\n")
				fmt.Printf("    - theme: %s\n", s.Theme)
				fmt.Printf("    - quietStartup: %s\n", quietStr)
				if s.DefaultProvider != "" {
					fmt.Printf("    - defaultProvider: %s\n", s.DefaultProvider)
				}
				if s.DefaultModel != "" {
					fmt.Printf("    - defaultModel: %s\n", s.DefaultModel)
				}
			}
		} else {
			fmt.Printf("  Configuración: ✗ no existe\n")
		}

		return nil
	},
}
