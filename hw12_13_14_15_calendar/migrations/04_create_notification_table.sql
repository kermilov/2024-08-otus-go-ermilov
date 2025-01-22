-- +goose Up

CREATE TABLE IF NOT EXISTS calendar.notification
(
    id uuid NOT NULL,
    title character varying COLLATE pg_catalog."default" NOT NULL,
    datetime timestamp without time zone NOT NULL,
    userid bigint NOT NULL,
    CONSTRAINT notification_pkey PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS calendar.notification
    OWNER to postgres;

-- +goose Down
DROP TABLE IF EXISTS calendar.notification;
