-- Run once with psql -v ON_ERROR_STOP=1.
-- Create a dedicated user's records and their events atomically.
BEGIN;

INSERT INTO user_change_counters (user_id, last_sequence)
VALUES ('ownership-test', 2)
ON CONFLICT (user_id)
DO UPDATE SET last_sequence = user_change_counters.last_sequence + 2;

INSERT INTO workouts (id, user_id, name, started_at)
VALUES (
    'f0000000-0000-4000-8000-000000000001',
    'ownership-test',
    'Ownership test',
    now()
);

INSERT INTO workout_sets
    (id, user_id, workout_id, exercise_name, reps, weight_kg)
VALUES (
    'f0000000-0000-4000-8000-000000000002',
    'ownership-test',
    'f0000000-0000-4000-8000-000000000001',
    'Ownership test',
    8,
    20
);

INSERT INTO change_log
    (user_id, sequence, entity_type, entity_id, operation, revision, payload)
SELECT w.user_id, c.last_sequence - 1, 'workout', w.id,
       'create', w.revision, to_jsonb(w)
FROM workouts AS w
JOIN user_change_counters AS c ON c.user_id = w.user_id
WHERE w.id = 'f0000000-0000-4000-8000-000000000001'
UNION ALL
SELECT s.user_id, c.last_sequence, 'workout_set', s.id,
       'create', s.revision, to_jsonb(s)
FROM workout_sets AS s
JOIN user_change_counters AS c ON c.user_id = s.user_id
WHERE s.id = 'f0000000-0000-4000-8000-000000000002';

COMMIT;

SELECT id AS set_id, workout_id, user_id, reps, weight_kg, revision
FROM workout_sets
WHERE id = 'f0000000-0000-4000-8000-000000000002';
