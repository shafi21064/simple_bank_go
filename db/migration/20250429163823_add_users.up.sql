CREATE TABLE "users" (
  "id" bigserial,
  "user_name" varchar PRIMARY KEY,
  "hassed_password" varchar NOT NULL,
  "full_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "password_changed_at" timestamp NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "accounts" ADD FOREIGN KEY ("owner_name") REFERENCES "users" ("user_name");

-- CREATE UNIQUE INDEX ON "accounts" ("owner_name", "currency");

ALTER TABLE "accounts" ADD CONSTRAINT "owner_currency_key" UNIQUE ("owner_name", "currency");