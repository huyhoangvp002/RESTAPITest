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
  "discount_price" int NOT NULL,
  "category_id" int,
  "value" int NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "discounts" (
  "id" bigserial PRIMARY KEY,
  "discount_value" int NOT NULL,
  "product_id" int
);

CREATE TABLE "accounts" (
  "id" bigserial PRIMARY KEY,
  "username" varchar UNIQUE NOT NULL,
  "hash_password" varchar NOT NULL,
  "role" varchar NOT NULL
);

COMMENT ON COLUMN "products"."price" IS 'CHECK (price > 0)';

COMMENT ON COLUMN "products"."discount_price" IS 'CHECK (0 < discount_price <= price)';

COMMENT ON COLUMN "products"."value" IS 'CHECK (price >= 0)';

COMMENT ON COLUMN "discounts"."discount_value" IS '0<value<100';

ALTER TABLE "products" ADD FOREIGN KEY ("category_id") REFERENCES "categories" ("id");

ALTER TABLE "discounts" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id");
