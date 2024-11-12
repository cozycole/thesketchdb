dir_path="./sql/schema"
db_name="sketch_data"

for file in "$dir_path"/*; do
    if [ -f "$file" ]; then
        psql $db_name -f "$file"
    fi
done
