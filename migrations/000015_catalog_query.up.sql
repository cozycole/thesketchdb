-- Run once before switching to sketch_listing_refactor.go.
-- Assumes sketch.id is integer. If it is bigint, change sketch_id and the
-- refresh function's argument to bigint.

CREATE TABLE IF NOT EXISTS sketch_search_documents (
    sketch_id integer PRIMARY KEY REFERENCES sketch(id) ON DELETE CASCADE,
    document tsvector NOT NULL
);

CREATE INDEX IF NOT EXISTS sketch_search_documents_document_gin
    ON sketch_search_documents USING GIN (document);

-- Rebuild the text-search document for ONE sketch.
-- Call SELECT refresh_sketch_search($1) after changing a sketch's title,
-- creator assignments, cast, tags, episode/show assignment, or the names
-- of referenced people/characters/creators/tags/shows.
-- Ideally call it within the same transaction as each CMS edit.
CREATE OR REPLACE FUNCTION refresh_sketch_search(p_sketch_id integer)
RETURNS void
LANGUAGE sql
AS $$
    INSERT INTO sketch_search_documents (sketch_id, document)
    SELECT
        v.id,
        setweight(to_tsvector('english', COALESCE(v.title, '')), 'A')
        || setweight(to_tsvector('english', concat_ws(' ',
                creators.words,
                cast_words.words,
                sh.name,
                sh.aliases
            )), 'B')
        || setweight(to_tsvector('english', COALESCE(tags.words, '')), 'C')
    FROM sketch v
    LEFT JOIN episode e ON e.id = v.episode_id
    LEFT JOIN season se ON se.id = e.season_id
    LEFT JOIN show sh ON sh.id = se.show_id
    LEFT JOIN LATERAL (
        SELECT string_agg(concat_ws(' ', c.name, c.alias), ' ') AS words
        FROM sketch_creator_rel scr
        JOIN creator c ON c.id = scr.creator_id
        WHERE scr.sketch_id = v.id
    ) creators ON TRUE
    LEFT JOIN LATERAL (
        SELECT string_agg(concat_ws(' ',
                   p.first, p.last, cm.character_name, ch.name
               ), ' ') AS words
        FROM cast_members cm
        LEFT JOIN person p ON p.id = cm.person_id
        LEFT JOIN character ch ON ch.id = cm.character_id
        WHERE cm.sketch_id = v.id
    ) cast_words ON TRUE
    LEFT JOIN LATERAL (
        SELECT string_agg(t.name, ' ') AS words
        FROM sketch_tags st
        JOIN tags t ON t.id = st.tag_id
        WHERE st.sketch_id = v.id
    ) tags ON TRUE
    WHERE v.id = p_sketch_id
    ON CONFLICT (sketch_id) DO UPDATE
        SET document = EXCLUDED.document;
$$;

-- Initial backfill. Re-run this after bulk imports or to rebuild all documents.
SELECT refresh_sketch_search(id) FROM sketch;

-- Supporting indexes for listing filters and thumbnail lookup.
-- Verify which indexes already exist before adding them in production.
CREATE INDEX IF NOT EXISTS sketch_creator_rel_creator_sketch_idx
    ON sketch_creator_rel (creator_id, sketch_id);
CREATE INDEX IF NOT EXISTS sketch_creator_rel_sketch_creator_idx
    ON sketch_creator_rel (sketch_id, creator_id);
CREATE INDEX IF NOT EXISTS cast_members_character_sketch_idx
    ON cast_members (character_id, sketch_id);
CREATE INDEX IF NOT EXISTS cast_members_person_sketch_idx
    ON cast_members (person_id, sketch_id);
CREATE INDEX IF NOT EXISTS cast_members_sketch_position_idx
    ON cast_members (sketch_id, position, id);
CREATE INDEX IF NOT EXISTS sketch_tags_tag_sketch_idx
    ON sketch_tags (tag_id, sketch_id);
CREATE INDEX IF NOT EXISTS sketch_tags_sketch_tag_idx
    ON sketch_tags (sketch_id, tag_id);

-- Consider adding sort indexes only if EXPLAIN shows sorting dominates.
CREATE INDEX IF NOT EXISTS sketch_popular_sort_idx
    ON sketch (popularity_score DESC, upload_date DESC, id DESC);
CREATE INDEX IF NOT EXISTS sketch_newest_sort_idx
    ON sketch (upload_date DESC, title ASC, id DESC);
