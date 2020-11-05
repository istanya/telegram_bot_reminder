CREATE TABLE "events" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "user_name" varchar NOT NULL,
  "dt_reminder" timestamptz NOT NULL,
  "bot_message_id" int NOT NULL,
  "message" varchar NOT NULL,
  "state" varchar NOT NULL,
  "dt_created" timestamptz NOT NULL DEFAULT (now())
);