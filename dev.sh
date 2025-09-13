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

browser-sync start \
  --watch \
  --files 'cmd/**/*, internal/**/*, ui/**/*' \
  --reload-delay 1800 \
  --port 4001 \
  --proxy '127.0.0.1:8080'
