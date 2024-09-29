air -c ./.air.toml & \
npx tailwind \
  -i './ui/styles.css' \
  -o './ui/static/css/styles.css' \
  --watch & \
browser-sync start \
  --files 'cmd/*, internal/*, ui/*' \
  --port 4001 \
  --proxy 'localhost:4000' \
  --middleware 'function(req, res, next) { \
    res.setHeader("Cache-Control", "no-cache, no-store, must-revalidate"); \
    return next(); \
  }'