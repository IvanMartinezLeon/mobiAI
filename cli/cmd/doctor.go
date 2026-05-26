package cmd

import (
	"fmt"

	"github.com/anomalyco/mobiAI/cli/internal/check"
	"github.com/spf13/cobra"
)

var doctorVerbose bool
var doctorJSON bool

func init() {
	doctorCmd.Flags().BoolVarP(&doctorVerbose, "verbose", "v", false, "Mostrar detalle completo")
	doctorCmd.Flags().BoolVar(&doctorJSON, "json", false, "Salida en formato JSON")
	rootCmd.AddCommand(doctorCmd)
}

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Diagnosticar el estado de la instalación",
	Long: `Ejecuta una serie de comprobaciones para verificar que
la personalización MOBI AI está correctamente instalada y configurada.

Flags:
  --verbose    Muestra información detallada de cada comprobación
  --json       Salida en formato JSON para scripting
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("\n🔍 Diagnóstico MOBI AI")
		fmt.Println("=====================")

		results := check.RunAll()

		if doctorJSON {
			fmt.Println("{")
			for i, r := range results {
				comma := ","
				if i == len(results)-1 {
					comma = ""
				}
				fmt.Printf("  %q: {\"passed\": %v, \"detail\": %q, \"suggestion\": %q}%s\n",
					r.Name, r.Passed, r.Detail, r.Suggestion, comma)
			}
			fmt.Println("}")
			return nil
		}

		fmt.Println()
		for _, r := range results {
			fmt.Println(r.String())
		}
		fmt.Println(check.Summary(results))
		return nil
	},
}
