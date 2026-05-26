# Instalación de @earendil-works/pi-coding-agent

Este repositorio contiene scripts automatizados para instalar el paquete global `@earendil-works/pi-coding-agent` de forma segura utilizando la opción `--ignore-scripts`.

La instalación manual equivale a ejecutar:
```bash
npm install -g --ignore-scripts @earendil-works/pi-coding-agent
```

## Requisitos Previos

Asegúrate de tener instalados **Node.js** y **npm** en tu sistema antes de ejecutar cualquiera de los scripts.

---

## Instrucciones de Instalación por Plataforma

### 🍎 macOS
1. Abre tu terminal.
2. Dirígete a este directorio y ejecuta:
   ```bash
   ./install_mac.sh
   ```
*Nota: Si la instalación global requiere privilegios de administrador, el script te solicitará ejecutarlo con `sudo`:*
```bash
sudo ./install_mac.sh
```

---

### 🐧 Linux
1. Abre tu terminal.
2. Dirígete a este directorio y ejecuta:
   ```bash
   ./install_linux.sh
   ```
*Nota: Si tu instalación de Node.js es a nivel de sistema (system-wide), es probable que necesites ejecutar el script como root:*
```bash
sudo ./install_linux.sh
```

---

### 🪟 Windows
1. Abre **PowerShell** (se recomienda ejecutarlo como Administrador si tu instalación global de npm lo requiere).
2. Dirígete a este directorio y ejecuta:
   ```powershell
   .\install_windows.ps1
   ```
*Nota: Si PowerShell bloquea la ejecución de scripts por las políticas de seguridad por defecto, puedes ejecutarlo con la política omitida para este proceso:*
```powershell
Set-ExecutionPolicy Bypass -Scope Process; .\install_windows.ps1
```

---

## Personalización de Tema e Interfaz (MOBI AI)

Este proyecto incluye una configuración local en el directorio `.pi/` para dotar a la terminal de Pi con un diseño personalizado y un banner ASCII de **MOBI AI** en el inicio de la sesión.

### Componentes de Personalización

1. **`.pi/themes/mobi-theme.json`**:
   * Define una paleta de colores oscura, moderna y premium (basada en el esquema oficial de 51 tokens del Pi Coding Agent).
   * Utiliza acentos en color cian, textos limpios y bordes oscuros integrados.

2. **`.pi/extensions/mobi-header.ts`**:
   * Una extensión de TypeScript auto-ejecutable que reemplaza el encabezado inicial de Pi con un logo estilizado en ASCII de **MOBI AI** empleando el color de acento del tema.
   * Añade una línea de subtítulo personalizada y registra un comando opcional `/builtin-header` por si deseas volver al banner nativo.

3. **`.pi/settings.json`**:
   * Configura localmente el agente para cargar por defecto el tema `mobi-theme`.

### Cómo se aplica

Los scripts de instalación automatizan este proceso de la siguiente manera:
1. **Respaldan** los archivos de personalización (`mobi-theme.json` y `mobi-header.ts`) en una ubicación temporal.
2. **Eliminan** la carpeta local `.pi/` del directorio de trabajo para evitar la advertencia de colisión de temas de Pi.
3. **Instalan** el Pi Coding Agent globalmente.
4. **Restauran** los archivos de personalización en la configuración global de Pi (`~/.pi/agent/` o equivalentemente en Windows) para que la personalización esté disponible en cualquier proyecto sin generar advertencias de conflicto.

Si no utilizas los instaladores automatizados, puedes realizar el proceso de forma manual copiando los contenidos a tu carpeta de configuración global de Pi:
* **Tema global**: Copia `mobi-theme.json` a `~/.pi/agent/themes/`
* **Extensión global**: Copia `mobi-header.ts` a `~/.pi/agent/extensions/`
* **Ajustes globales**: Añade `"theme": "mobi-theme"` a tu `~/.pi/agent/settings.json`


