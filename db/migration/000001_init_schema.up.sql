CREATE TABLE "accounts" (
	  "id" bigserial PRIMARY KEY,
	  "owner" varchar NOT NULL,
	  "balance" bigint NOT NULL,
	  "currency" varchar NOT NULL,
	  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "entries" (
	  "id" bigserial PRIMARY KEY,
	  "account_id" bigint NOT NULL,
	  "amount" bigint NOT NULL,
	  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "transfers" (
	  "id" bigserial PRIMARY KEY,
	  "src_account_id" bigint NOT NULL,
	  "dest_account_id" bigint NOT NULL,
	  "amount" bigint NOT NULL,
	  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE INDEX ON "accounts" ("owner");

CREATE INDEX ON "entries" ("account_id");

CREATE INDEX ON "transfers" ("src_account_id");

CREATE INDEX ON "transfers" ("dest_account_id");

CREATE INDEX ON "transfers" ("src_account_id", "dest_account_id");

COMMENT ON COLUMN "transfers"."amount" IS 'must be positive';

ALTER TABLE "entries" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("src_account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("dest_account_id") REFERENCES "accounts" ("id");
