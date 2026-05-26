# MOBI AI — Personalización para Pi Coding Agent

![Go](https://img.shields.io/badge/CLI-Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)
![GitHub Release](https://img.shields.io/github/v/release/IvanMartinezLeon/mobiAI)

Ecosistema de personalización para el [Pi Coding Agent](https://github.com/earendil-works/pi-coding-agent) con tema oscuro, banner ASCII, detección automática de framework, métricas de sesión y más.

## 🚀 Instalación Sin Clonar (Recomendado)

No necesitas clonar este repositorio. Usa cualquiera de estos métodos para instalar inmediatamente:

### Opción 1: Binario Standalone (macOS/Linux)
```bash
curl -fsSL https://raw.githubusercontent.com/IvanMartinezLeon/mobiAI/main/cli/install.sh | sh
```

### Opción 2: Con Go Instalado
```bash
go install github.com/IvanMartinezLeon/mobiAI/cli@latest
```

### Opción 3: Desarrollo local
```bash
git clone https://github.com/IvanMartinezLeon/mobiAI.git
cd mobiAI/cli
go build -o mobi .
./mobi install
```

## 🖥️ CLI `mobi`

Una sola herramienta para instalar, diagnosticar y mantener la personalización.

### Comandos Disponibles

| Comando | Descripción |
|---------|-------------|
| `mobi install` | Instala Pi globalmente + copia toda la personalización |
| `mobi install --no-npm` | Solo copia archivos (si Pi ya está instalado) |
| `mobi doctor` | Diagnostica Node.js, Pi, theme, extensión, configuración |
| `mobi doctor --json` | Salida JSON para scripting |
| `mobi status` | Muestra versión de Pi, theme, extensión, settings |
| `mobi update` | Sincroniza `.pi/` local → `~/.pi/agent/` |

### Ejemplo de Uso

```bash
# Instalación completa (recomendado para primeros usuarios)
mobi install

# Verificar que todo está correcto
mobi doctor

# Ver estado actual
mobi status
```

## 📁 Estructura del Proyecto

```
mobiAI/
├── cli/                  ← CLI en Go (mobi)
├── packages/
│   └── mobi-cli/         ← Paquete npm (wrapper del binario Go)
├── .pi/                  ← Fuente de personalización
│   ├── themes/
│   │   └── mobi-theme.json
│   ├── extensions/
│   │   └── mobi-header.ts
│   ├── settings.json
│   └── APPEND_SYSTEM.md
└── legacy/               ← Scripts legacy (fallback)
    ├── install_mac.sh
    ├── install_linux.sh
    └── install_windows.ps1
```

## 🎨 Componentes de Personalización

### `mobi-theme.json`
Tema oscuro con paleta de 51 tokens, acentos en cian, bordes integrados y jerarquía visual limpia.

### `mobi-header.ts`
Extensión TypeScript que reemplaza el header de Pi con el logo ASCII de MOBI AI y añade un footer con:
- Framework del proyecto (detección automática)
- Rama activa de git
- Modelo en uso
- % de contexto utilizado
- Tokens de subida/bajada acumulados

### `APPEND_SYSTEM.md`
Instrucciones inyectadas en el prompt del agente para detectar el framework al iniciar la conversación.

### `settings.json`
`quietStartup: true` oculta las secciones de inicio; `theme: mobi-theme` como tema por defecto.

## 💻 Desarrollo

```bash
cd cli
go build -o mobi .    # Compilar la CLI
go run . doctor       # Probar sin compilar
```

### Lanzamiento de Nuevas Versiones

Etiquetar con `cli-v*` para lanzar release automático via GitHub Actions:

```bash
git tag cli-v0.1.0
git push origin cli-v0.1.0
```

## 🔧 Legacy

Los scripts de instalación originales (`install_mac.sh`, `install_linux.sh`, `install_windows.ps1`) están en `legacy/` como fallback. La CLI `mobi` es el método recomendado para todos los usuarios.

---
<sup>MOBI AI está inspirado y creado para mejorar la experiencia con [Pi Coding Agent](https://github.com/earendil-works/pi-coding-agent).</sup>