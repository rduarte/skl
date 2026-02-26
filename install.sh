#!/bin/bash
#
# SKL Installer
# Uso: curl -sSfL https://raw.githubusercontent.com/rduarte/skl/main/install.sh | bash
#

set -euo pipefail

REPO="rduarte/skl"
INSTALL_DIR="${HOME}/.local/bin"
BINARY_NAME="skl"

info()  { echo -e "\033[1;34m→\033[0m $*"; }
ok()    { echo -e "\033[1;32m✅\033[0m $*"; }
err()   { echo -e "\033[1;31m✗\033[0m $*" >&2; exit 1; }

# Detect architecture
ARCH=$(uname -m)
case "$ARCH" in
    x86_64)  ARCH="amd64" ;;
    aarch64) ARCH="arm64" ;;
    *)       err "Arquitetura não suportada: $ARCH" ;;
esac

OS="linux"

# Get latest release tag from GitHub API
info "Buscando última versão..."
LATEST=$(curl -sSf "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | cut -d'"' -f4)

if [ -z "$LATEST" ]; then
    err "Não foi possível obter a última versão. Verifique: https://github.com/${REPO}/releases"
fi

info "Versão: ${LATEST}"

# Download binary
ASSET_NAME="${BINARY_NAME}-${OS}-${ARCH}"
DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${LATEST}/${ASSET_NAME}"

info "Baixando ${DOWNLOAD_URL}..."
TMP_FILE=$(mktemp)
trap "rm -f ${TMP_FILE}" EXIT

if ! curl -sSfL -o "${TMP_FILE}" "${DOWNLOAD_URL}"; then
    err "Falha ao baixar o binário. Verifique: https://github.com/${REPO}/releases"
fi

# Install
mkdir -p "${INSTALL_DIR}"
mv "${TMP_FILE}" "${INSTALL_DIR}/${BINARY_NAME}"
chmod +x "${INSTALL_DIR}/${BINARY_NAME}"

ok "skl ${LATEST} instalado em ${INSTALL_DIR}/${BINARY_NAME}"

# Check if in PATH
if ! echo "$PATH" | tr ':' '\n' | grep -q "^${INSTALL_DIR}$"; then
    echo ""
    echo "⚠  ${INSTALL_DIR} não está no seu PATH."
    echo "   Adicione ao seu ~/.bashrc:"
    echo ""
    echo "   export PATH=\"\$HOME/.local/bin:\$PATH\""
    echo ""
fi
