#!/bin/bash

air -c ./.air.toml 

npx tailwind \
  -c ./tailwind.config.js \
  -i ./ui/static/css/styles.css \
  -o ./ui/static/css/dist/styles.css \
  --watch 

npx esbuild ./ui/static/js/main.js \
  --bundle \
  --outfile=./ui/static/js/dist/main.js \
  --sourcemap \
  --watch=forever 
  # --minify FOR PRODUCTION

browser-sync start \
  --watch \
  --files 'cmd/**/*, internal/**/*, ui/**/*' \
  --reload-delay 1800 \
  --port 4001 \
  --proxy 'localhost:4000' \
  --middleware 'function(req, res, next) { \
    res.setHeader("Cache-Control", "no-cache, no-store, must-revalidate"); \
    return next(); \
  }' 
