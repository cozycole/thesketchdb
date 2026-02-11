# --- main site assets ---
FROM node:18-alpine AS assets
WORKDIR /app

COPY package*.json .
COPY tailwind.config.js .
RUN npm install

COPY ./ui ./ui
RUN mkdir -p dist/css dist/js dist/img dist/fonts
RUN npx esbuild ./ui/static/js/main.js --bundle --minify --outfile=./dist/js/main.js
RUN JS_HASH=$(sha256sum ./dist/js/main.js | cut -c1-8) && mv ./dist/js/main.js ./dist/js/main.${JS_HASH}.js

RUN npx tailwindcss -i ./ui/static/css/styles.css -o ./dist/css/styles.css --minify --verbose
RUN CSS_HASH=$(sha256sum ./dist/css/styles.css | cut -c1-8) && mv ./dist/css/styles.css ./dist/css/styles.${CSS_HASH}.css
RUN cp -r ./ui/static/img ./dist 2>/dev/null 
RUN cp -r ./ui/static/fonts ./dist 2>/dev/null 

RUN JS_FILE=$(basename $(ls ./dist/js/main.*.js)) && \
    CSS_FILE=$(basename $(ls ./dist/css/styles.*.css)) && \
    echo "{ \"css\": \"$CSS_FILE\", \"js\": \"$JS_FILE\" }" > ./dist/manifest.json

# --- CMS build ---
FROM node:20-alpine AS cms_build
WORKDIR /cms

COPY ./cms/package*.json ./
RUN npm install

COPY ./cms ./

RUN npm run build


# --- Go build ---
FROM golang:1.24 AS build
WORKDIR /app
COPY . .
COPY --from=assets /app/dist ./dist
COPY --from=cms_build /cms/dist ./cmsdist
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ./bin/web-app ./cmd/web


# --- runtime ---
FROM alpine:3.17
RUN apk --no-cache add ca-certificates
WORKDIR /app

COPY --from=build /app/bin/web-app ./web-app
COPY --from=build /app/dist /tmp/dist
COPY --from=build /app/cmsdist ./cmsdist

EXPOSE 8080
ENTRYPOINT sh -c 'cp -r /tmp/dist/* /app/dist/ && exec /app/web-app --addr 0.0.0.0:8080'
