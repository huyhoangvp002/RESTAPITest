CREATE TABLE categories (
  id bigserial PRIMARY KEY,
  name varchar NOT NULL,
  type varchar NOT NULL,
  account_id int,
  created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE products (
  id bigserial PRIMARY KEY,
  name varchar NOT NULL,
  price int NOT NULL,
  discount_price int NOT NULL,
  value int NOT NULL,
  account_id int,
  category_id int,
  created_at timestamptz NOT NULL DEFAULT now(),
  CHECK (price > 0),
  CHECK (discount_price > 0 AND discount_price <= price),
  CHECK (value >= 0)
);

CREATE TABLE discounts (
  id bigserial PRIMARY KEY,
  discount_value int NOT NULL,
  account_id int,
  product_id int,
  created_at timestamptz NOT NULL DEFAULT now(),
  CHECK (discount_value > 0 AND discount_value < 100)
);

CREATE TABLE accounts (
  id bigserial PRIMARY KEY,
  username varchar UNIQUE NOT NULL,
  hash_password varchar NOT NULL,
  role varchar NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE account_info (
  id bigserial PRIMARY KEY,
  name varchar NOT NULL,
  email varchar UNIQUE NOT NULL,
  phone_number varchar UNIQUE NOT NULL,
  address varchar NOT NULL,
  account_id int UNIQUE,
  created_at timestamptz NOT NULL DEFAULT now(),
  update_at timestamptz NOT NULL DEFAULT now(),
  CHECK (phone_number ~ '^\+?[0-9]{10,15}$')
);

CREATE TABLE cart (
  id bigserial PRIMARY KEY,
  value int NOT NULL CHECK (value > 0),
  account_id int NOT NULL,
  product_id int NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (account_id, product_id)
);


-- Foreign keys with CASCADE
ALTER TABLE categories ADD FOREIGN KEY (account_id) REFERENCES accounts (id) ON DELETE CASCADE;

ALTER TABLE products ADD FOREIGN KEY (account_id) REFERENCES accounts (id) ON DELETE CASCADE;
ALTER TABLE products ADD FOREIGN KEY (category_id) REFERENCES categories (id) ON DELETE CASCADE;

ALTER TABLE discounts ADD FOREIGN KEY (account_id) REFERENCES accounts (id) ON DELETE CASCADE;
ALTER TABLE discounts ADD FOREIGN KEY (product_id) REFERENCES products (id) ON DELETE CASCADE;

ALTER TABLE account_info ADD FOREIGN KEY (account_id) REFERENCES accounts (id) ON DELETE CASCADE;

ALTER TABLE cart ADD FOREIGN KEY (account_id) REFERENCES accounts (id) ON DELETE CASCADE;
ALTER TABLE cart ADD FOREIGN KEY (product_id) REFERENCES products (id) ON DELETE CASCADE;

INSERT INTO accounts (username, hash_password, role) VALUES
('user1', 'hashedpassword1', 'customer'),
('user2', 'hashedpassword2', 'customer');

INSERT INTO account_info (name, email, phone_number, address, account_id) VALUES
('Alice Smith', 'alice@example.com', '0912345678', '123 Main St', 1),
('Bob Brown', 'bob@example.com', '0987654321', '456 Elm St', 2);

INSERT INTO categories (name, type, account_id) VALUES
('Electronics', 'Gadgets', 1),
('Books', 'Media', 1),
('Clothing', 'Apparel', 2),
('Home', 'Furniture', 2),
('Toys', 'Kids', 1);


INSERT INTO products (name, price, discount_price, value, account_id, category_id) VALUES
('Phone', 500, 450, 100, 1, 1),
('Laptop', 1000, 900, 150, 1, 1),
('Novel', 20, 18, 50, 2, 2),
('T-Shirt', 30, 25, 60, 2, 3),
('Sofa', 400, 350, 80, 1, 4),
('Toy Car', 50, 45, 40, 2, 5),
('Headphones', 100, 85, 70, 1, 1),
('Cookbook', 35, 30, 40, 2, 2),
('Jeans', 60, 55, 90, 2, 3),
('Armchair', 200, 180, 75, 1, 4);

INSERT INTO discounts (discount_value, account_id, product_id) VALUES
(10, 1, 1),
(15, 2, 2),
(5, 1, 3),
(20, 2, 4),
(25, 1, 5),
(30, 2, 6),
(12, 1, 7),
(8, 2, 8),
(18, 1, 9),
(22, 2, 10);

