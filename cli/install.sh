#!/bin/sh
set -e

echo "=== Instalando CLI mobi ==="

if ! command -v curl >/dev/null 2>&1; then
  echo "Error: curl no está instalado"
  exit 1
fi

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
  x86_64)  ARCH="amd64" ;;
  aarch64) ARCH="arm64" ;;
  arm64)   ARCH="arm64" ;;
  *)
    echo "Arquitectura no soportada: $ARCH"
    exit 1
    ;;
esac

case "$OS" in
  darwin|linux) ;;
  *)
    echo "Sistema operativo no soportado: $OS"
    exit 1
    ;;
esac

INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

echo "  → Detectado: $OS/$ARCH"
echo "  → Instalando en: $INSTALL_DIR/mobi"

LATEST=$(curl -s https://api.github.com/repos/IvanMartinezLeon/mobiAI/releases/latest 2>/dev/null | grep '"tag_name"' | cut -d'"' -f4)
if [ -z "$LATEST" ]; then
  echo "  ⚠ No se pudo obtener la última versión, usando 'latest'"
  LATEST="latest"
fi

URL="https://github.com/IvanMartinezLeon/mobiAI/releases/download/$LATEST/mobi_${OS}_${ARCH}.tar.gz"

TMP_FILE=$(mktemp)
trap "rm -f $TMP_FILE" EXIT

echo "  → Descargando $URL..."
curl -fsSL "$URL" -o "$TMP_FILE"

tar -xzf "$TMP_FILE" -C /tmp mobi 2>/dev/null || {
  echo "  → Descargando binario directo..."
  BIN_URL="https://github.com/IvanMartinezLeon/mobiAI/releases/download/$LATEST/mobi_${OS}_${ARCH}"
  curl -fsSL "$BIN_URL" -o /tmp/mobi
  chmod +x /tmp/mobi
}

if [ ! -f /tmp/mobi ]; then
  echo "Error: no se pudo descargar el binario"
  exit 1
fi

chmod +x /tmp/mobi

if [ -w "$INSTALL_DIR" ]; then
  mv /tmp/mobi "$INSTALL_DIR/mobi"
else
  echo "  → No hay permisos de escritura en $INSTALL_DIR, usando sudo..."
  sudo mv /tmp/mobi "$INSTALL_DIR/mobi"
fi

echo ""
echo "✅ CLI mobi instalada correctamente"
echo ""
echo "Próximos pasos:"
echo "  mobi install    — Instalar/configurar Pi con personalización MOBI AI"
echo "  mobi doctor     — Diagnosticar el estado"
echo "  mobi status     — Mostrar estado actual"
