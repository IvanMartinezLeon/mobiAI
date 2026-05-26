package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/IvanMartinezLeon/mobiAI/cli/internal/pi"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(updateCmd)
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Sincronizar personalización desde el repositorio local",
	Long: `Copia los archivos de personalización (.pi/) al directorio global
de Pi (~/.pi/agent/), actualizando theme, extensión, APPEND_SYSTEM.md
y haciendo merge de la configuración.

Ejecuta esto cuando hagas cambios locales en .pi/ y quieras aplicarlos.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("\n🔄 Actualizando personalización MOBI AI...")

		repoRoot := detectRepoRoot()

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
			srcData, err := os.ReadFile(c.src)
			if err != nil {
				fmt.Printf("  ⚠ No se pudo leer %s: %s\n", c.label, err.Error())
				continue
			}
			changed := true
			if pi.FileExists(c.dst) {
				dstData, _ := os.ReadFile(c.dst)
				changed = string(srcData) != string(dstData)
			}
			if err := pi.CopyFile(c.src, c.dst); err != nil {
				fmt.Printf("  ✗ Error copiando %s: %s\n", c.label, err.Error())
				continue
			}
			if changed {
				fmt.Printf("  ✓ %s actualizado\n", c.label)
			} else {
				fmt.Printf("  · %s sin cambios\n", c.label)
			}
		}

		if err := mergeSettings(repoRoot); err != nil {
			fmt.Printf("  ⚠ Error actualizando settings: %s\n", err.Error())
		}

		fmt.Println("\n✅ Personalización actualizada")
		return nil
	},
}
