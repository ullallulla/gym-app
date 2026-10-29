-- +goose Up
ALTER TABLE workouts
    ADD CONSTRAINT workouts_user_id_id_key UNIQUE (user_id, id);
CREATE TABLE workout_sets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id TEXT NOT NULL,
    workout_id UUID NOT NULL,
    exercise_name TEXT NOT NULL CHECK (length(btrim(exercise_name)) > 0),
    reps INTEGER NOT NULL CHECK (reps > 0),
    weight_kg DOUBLE PRECISION NOT NULL
        CHECK (weight_kg >= 0 AND weight_kg < 'Infinity'::DOUBLE PRECISION),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    revision BIGINT NOT NULL DEFAULT 1 CHECK (revision > 0),
    CONSTRAINT workout_sets_workout_owner_fk
        FOREIGN KEY (user_id, workout_id)
        REFERENCES workouts (user_id, id)
);
CREATE INDEX workout_sets_workout_order_idx
    ON workout_sets (user_id, workout_id, created_at, id);

-- +goose Down
DROP TABLE workout_sets;
ALTER TABLE workouts DROP CONSTRAINT workouts_user_id_id_key;
