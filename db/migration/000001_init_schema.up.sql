CREATE TABLE "accounts" (
  "id" BIGSERIAL PRIMARY KEY,
  "owner" VARCHAR NOT NULL,
  "balance" BIGINT NOT NULL,
  "created_at" TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE "entries" (
  "id" BIGSERIAL PRIMARY KEY,
  "account_id" BIGINT NOT NULL,
  "amount" BIGINT NOT NULL,
  "created_at" TIMESTAMP NOT NULL DEFAULT NOW(),
  FOREIGN KEY ("account_id") REFERENCES "accounts" ("id") ON DELETE CASCADE
);

CREATE TABLE "transfers" (
  "id" BIGSERIAL PRIMARY KEY,
  "from_account_id" BIGINT NOT NULL,
  "to_account_id" BIGINT NOT NULL,
  "amount" BIGINT NOT NULL,
  "created_at" TIMESTAMP NOT NULL DEFAULT NOW(),
  CHECK ("from_account_id" <> "to_account_id"),
  FOREIGN KEY ("from_account_id") REFERENCES "accounts" ("id") ON DELETE CASCADE,
  FOREIGN KEY ("to_account_id") REFERENCES "accounts" ("id") ON DELETE CASCADE
);


CREATE INDEX ON "accounts" ("owner");

CREATE INDEX ON "entries" ("account_id");

CREATE INDEX ON "transfers" ("from_account_id");

CREATE INDEX ON "transfers" ("to_account_id");

CREATE INDEX ON "transfers" ("from_account_id", "to_account_id");