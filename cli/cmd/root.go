package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mobi",
	Short: "MobiAI CLI — Instalación y gestión del Pi Coding Agent",
	Long: `MobiAI CLI unifica la instalación, diagnóstico y actualización
de la personalización MOBI AI para el Pi Coding Agent.

Ejemplo:
  mobi install    Instalar/configurar Pi con personalización MOBI AI
  mobi doctor     Diagnosticar el estado de la instalación
  mobi status     Mostrar estado actual de la configuración
  mobi update     Sincronizar personalización desde el repositorio local
`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
