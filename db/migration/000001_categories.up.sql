CREATE TABLE "categories" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "type" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "products" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "price" int CHECK (price > 0) NOT NULL,
  "discount_price" int CHECK (0 < discount_price) NOT NULL,
  "category_id" int,
  "value" int CHECK (value > -1) NOT NULL,
  "customers_id" int,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "discounts" (
  "id" bigserial PRIMARY KEY,
  "discount_value" int CHECK (discount_value > -1 AND discount_value < 101) NOT NULL,
  "product_id" int,
  "created_at" timestamptz NOT NULL DEFAULT 'now()',
  "customer_id" int
);

CREATE TABLE "accounts" (
  "id" bigserial PRIMARY KEY,
  "username" varchar UNIQUE NOT NULL,
  "hash_password" varchar NOT NULL,
  "role" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "customers" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "account_id" int UNIQUE,
  "email" varchar UNIQUE NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

ALTER TABLE "products" ADD FOREIGN KEY ("category_id") REFERENCES "categories" ("id");

ALTER TABLE "products" ADD FOREIGN KEY ("customers_id") REFERENCES "customers" ("id");

ALTER TABLE "discounts" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id");

ALTER TABLE "discounts" ADD FOREIGN KEY ("customer_id") REFERENCES "customers" ("id");

ALTER TABLE "customers" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

INSERT INTO accounts (username, hash_password, role, created_at) VALUES
('newuser1', 'hashedpassword1', 'customer', NOW()),
('newdealer1', 'hashedpassword2', 'user', NOW()),
('newadmin1', 'hashedpassword3', 'admin', NOW());

INSERT INTO customers (id, name, account_id, email, created_at) VALUES
(1, 'Customer One', 1, 'customer1@example.com', NOW()),
(2, 'Customer Two', 2, 'customer2@example.com', NOW());

INSERT INTO categories (id, name, type, created_at) VALUES
(1, 'Electronics', 'Gadget', NOW()),
(2, 'Books', 'Media', NOW());

INSERT INTO products (name, price, discount_price, category_id, value, customers_id, created_at) VALUES
('Product A', 100, 90, 1, 10, 1, NOW()),
('Product B', 200, 180, 2, 15, 2, NOW()),
('Product C', 150, 135, 1, 12, 1, NOW()),
('Product D', 120, 100, 2, 8, 2, NOW()),
('Product E', 300, 270, 1, 20, 1, NOW()),
('Product F', 400, 350, 2, 25, 2, NOW()),
('Product G', 50, 45, 1, 5, 1, NOW()),
('Product H', 250, 220, 2, 18, 2, NOW()),
('Product I', 180, 160, 1, 14, 1, NOW()),
('Product J', 220, 200, 2, 16, 2, NOW());

INSERT INTO discounts (discount_value, product_id, customer_id, created_at) VALUES
(10, 1, 1, NOW()),
(15, 2, 2, NOW()),
(5, 3, 1, NOW()),
(20, 4, 2, NOW()),
(25, 5, 1, NOW()),
(30, 6, 2, NOW()),
(8, 7, 1, NOW()),
(12, 8, 2, NOW()),
(18, 9, 1, NOW()),
(22, 10, 2, NOW());
