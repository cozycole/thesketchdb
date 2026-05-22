ALTER TABLE sketch_creator_rel
DROP CONSTRAINT sketch_creator_rel_sketch_id_fkey;

ALTER TABLE sketch_creator_rel
ADD CONSTRAINT sketch_creator_rel_sketch_id_fkey
FOREIGN KEY (sketch_id)
REFERENCES sketch(id)
ON DELETE CASCADE;

ALTER TABLE sketch_creator_rel
DROP CONSTRAINT sketch_creator_rel_creator_id_fkey;

ALTER TABLE sketch_creator_rel
ADD CONSTRAINT sketch_creator_rel_creator_id_fkey
FOREIGN KEY (creator_id)
REFERENCES creator(id)
ON DELETE CASCADE;

ALTER TABLE transcription_lines
DROP CONSTRAINT transcription_lines_sketch_id_fkey;

ALTER TABLE transcription_lines
ADD CONSTRAINT transcription_lines_sketch_id_fkey
FOREIGN KEY (sketch_id)
REFERENCES sketch(id)
ON DELETE CASCADE;

ALTER TABLE sketch_video
DROP CONSTRAINT sketch_video_sketch_id_fkey;

ALTER TABLE sketch_video
ADD CONSTRAINT sketch_video_sketch_id_fkey
FOREIGN KEY (sketch_id)
REFERENCES sketch(id)
ON DELETE CASCADE;

ALTER TABLE cast_auto_screenshots
DROP CONSTRAINT cast_auto_screenshots_sketch_id_fkey;

ALTER TABLE cast_auto_screenshots
ADD CONSTRAINT cast_auto_screenshots_sketch_id_fkey
FOREIGN KEY (sketch_id)
REFERENCES sketch(id)
ON DELETE CASCADE;
