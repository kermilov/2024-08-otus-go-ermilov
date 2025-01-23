-- +goose Up

ALTER TABLE IF EXISTS calendar.event ADD COLUMN notification_duration interval;
ALTER TABLE IF EXISTS calendar.event ADD COLUMN is_send_notification boolean DEFAULT False;

CREATE INDEX event_datetime_notification_duration
    ON calendar.event USING btree
    (is_send_notification ASC NULLS LAST, datetime ASC NULLS LAST, notification_duration ASC NULLS LAST)
    WITH (deduplicate_items=False)
;

-- +goose Down

ALTER TABLE IF EXISTS calendar.event DROP COLUMN notification_duration;
ALTER TABLE IF EXISTS calendar.event DROP COLUMN is_send_notification;

DROP INDEX IF EXISTS event_datetime_notification_duration;
