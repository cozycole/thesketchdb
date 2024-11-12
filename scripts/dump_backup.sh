printf -v date '%(%Y-%m-%d)T' -1
pg_dump --data-only sketchdb > backups/sketchdb_backup$date.sql
