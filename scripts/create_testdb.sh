psql -d test_sketch_data -f sql/schema.sql
psql -d test_sketch_data -f sql/triggers.sql
psql -d test_sketch_data -f sql/testdata/test1.sql