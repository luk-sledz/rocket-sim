#!/usr/bin/env bash
# Start serwera wizualizacji + otwarcie przeglądarki.
# Użycie: ./run-web.sh   (albo: bash run-web.sh)
set -e
cd "$(dirname "$0")"

URL="http://localhost:8080"
PORT=8080

# Stary proces na :8080 = stary kod (np. zamrożona telemetria przy łapaniu).
if command -v fuser >/dev/null 2>&1; then
  fuser -k "${PORT}/tcp" >/dev/null 2>&1 || true
elif command -v lsof >/dev/null 2>&1; then
  pid=$(lsof -t -iTCP:"$PORT" -sTCP:LISTEN 2>/dev/null || true)
  if [ -n "$pid" ]; then kill $pid 2>/dev/null || true; fi
fi
sleep 0.2

# Otwórz przeglądarkę po krótkiej chwili (serwer musi zdążyć wstać).
(sleep 1 && xdg-open "$URL") >/dev/null 2>&1 &

echo "Serwer: $URL  (Ctrl+C = stop)"
exec go run ./cmd/web
