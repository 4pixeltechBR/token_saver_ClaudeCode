#!/bin/sh
set -eu
package_root=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
binary="$package_root/bin/token-saver"
manifest="$package_root/SHA256SUMS.txt"
if [ ! -f "$binary" ] || [ ! -f "$manifest" ]; then
  printf '%s\n' 'Pacote incompleto. Extraia o ZIP inteiro.' >&2
  exit 1
fi
expected=$(awk '$2 == "bin/token-saver" {print $1}' "$manifest")
if command -v sha256sum >/dev/null 2>&1; then
  actual=$(sha256sum "$binary" | awk '{print $1}')
elif command -v shasum >/dev/null 2>&1; then
  actual=$(shasum -a 256 "$binary" | awk '{print $1}')
else
  printf '%s\n' 'Nao foi possivel verificar o pacote: sha256sum ou shasum necessario.' >&2
  exit 1
fi
if [ "${#expected}" -ne 64 ] || [ "$expected" != "$actual" ]; then
  printf '%s\n' 'O pacote esta corrompido. Baixe novamente a versao publicada.' >&2
  exit 1
fi
chmod u+x "$binary"
if ! command -v claude >/dev/null 2>&1; then
  printf '%s\n' 'Claude Code nao encontrado no PATH. Confira seu acesso antes de usar a skill.'
fi
exec "$binary" install "$@"
