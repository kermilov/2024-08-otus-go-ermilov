-- +goose Up

CREATE UNIQUE INDEX event_datetime_duration
    ON calendar.event USING btree
    (datetime ASC NULLS LAST, duration ASC NULLS LAST)
    WITH (deduplicate_items=False)
;

-- +goose Down

DROP INDEX IF EXISTS event_datetime_duration;