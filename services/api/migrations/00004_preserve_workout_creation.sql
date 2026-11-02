-- +goose Up
ALTER TABLE workouts
    ADD COLUMN original_name TEXT;
ALTER TABLE workouts
    ADD COLUMN original_started_at TIMESTAMPTZ;
ALTER TABLE workouts
    ADD CONSTRAINT original_name CHECK ((original_name IS NULL) = (original_started_at IS NULL));
WITH original_creates AS (SELECT c.user_id, c.entity_id, c.payload ->> 'name' AS original_name, CASE WHEN pg_input_is_valid(c.payload ->> 'started_at', 'timestamptz') THEN (c.payload ->> 'started_at')::TIMESTAMPTZ ELSE NULL END AS original_started_at FROM change_log AS c JOIN workouts AS w ON w.id = c.entity_id AND w.user_id = c.user_id WHERE NOT EXISTS (SELECT 1 FROM change_log AS other WHERE other.user_id = c.user_id AND other.entity_id = c.entity_id AND other.entity_type = 'workout' AND other.operation = 'create' AND other.sequence <> c.sequence) AND c.revision = 1 AND (c.payload -> 'revision') = '1'::JSONB AND (c.payload ->> 'id') = c.entity_id::TEXT AND (c.payload ->> 'user_id') = c.user_id AND c.entity_type = 'workout' AND c.operation = 'create' AND jsonb_typeof(c.payload -> 'name') = 'string' AND jsonb_typeof(c.payload -> 'started_at')  = 'string' AND (c.payload ->> 'started_at') ~ '^[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}([.][0-9]{1,6})?(Z|[+-][0-9]{2}:[0-9]{2})$') UPDATE workouts AS w
SET original_name = src.original_name,
    original_started_at = src.original_started_at
FROM original_creates AS src
WHERE w.id = src.entity_id
  AND w.user_id = src.user_id
  AND src.original_started_at IS NOT NULL
  AND src.original_name <> ''
  AND w.original_name IS NULL
  AND w.original_started_at IS NULL;
-- +goose Down
ALTER TABLE workouts DROP CONSTRAINT original_name;
ALTER TABLE workouts
    DROP COLUMN original_name;
ALTER TABLE workouts
    DROP COLUMN original_started_at;
