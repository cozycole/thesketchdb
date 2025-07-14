printf -v date '%(%Y-%m-%d)T' -1
pg_dump --data-only thesketchdb > backups/sketchdb_backup$date.sql
