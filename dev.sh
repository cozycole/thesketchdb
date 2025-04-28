#!/bin/bash

air -c ./.air.toml| sed 's/^/[AIR] /' & 

npx tailwind \
  -c ./tailwind.config.js \
  -i ./ui/styles.css \
  -o ./ui/static/css/styles.css \
  --watch=always | sed 's/^/[TAILWIND] /' &

npx esbuild ./ui/static/js/main.js \
  --bundle \
  --outfile=./ui/static/js/dist/main.js \
  --sourcemap \
  --watch=forever | sed 's/^/[ESBUILD] /' & 
  # --minify FOR PRODUCTION

browser-sync start \
  --watch \
  --files 'cmd/**/*, internal/**/*, ui/**/*' \
  --reload-delay 1000 \
  --port 4001 \
  --proxy 'localhost:4000' \
  --middleware 'function(req, res, next) { \
    res.setHeader("Cache-Control", "no-cache, no-store, must-revalidate"); \
    return next(); \
  }' | sed 's/^/[BROWSERSYNC] /'
