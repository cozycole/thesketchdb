printf -v date '%(%Y-%m-%d)T' -1
pg_dump --data-only test_sketch_data > backups/sketchdb_backup$date.sql
