-- +goose Up

CREATE TABLE IF NOT EXISTS calendar.event
(
    id uuid NOT NULL,
    title character varying COLLATE pg_catalog."default" NOT NULL,
    datetime timestamp without time zone NOT NULL,
    duration interval NOT NULL,
    userid bigint NOT NULL,
    CONSTRAINT event_pkey PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS calendar.event
    OWNER to postgres;

-- +goose Down
DROP TABLE IF EXISTS calendar.event;
