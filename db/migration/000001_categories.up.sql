CREATE TABLE "categories" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "type" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "products" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "price" int NOT NULL,
  "category_id" int,
  "value" int NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

COMMENT ON COLUMN "products"."price" IS 'CHECK (price > 0)';

COMMENT ON COLUMN "products"."value" IS 'CHECK (price >= 0)';

ALTER TABLE "products" ADD FOREIGN KEY ("category_id") REFERENCES "categories" ("id");
