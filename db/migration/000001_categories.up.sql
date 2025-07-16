-- Table `accounts` lưu thông tin đăng nhập và phân quyền
CREATE TABLE accounts (
  id BIGSERIAL PRIMARY KEY,
  username VARCHAR NOT NULL UNIQUE,
  hash_password VARCHAR NOT NULL,
  role VARCHAR NOT NULL CHECK (role IN ('admin', 'seller', 'buyer', 'guest')),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Table `account_info` lưu thông tin hồ sơ (1-1 với accounts)
CREATE TABLE account_info (
  id BIGSERIAL PRIMARY KEY,
  account_id BIGINT NOT NULL UNIQUE REFERENCES accounts(id) ON DELETE CASCADE,
  name VARCHAR NOT NULL,
  email VARCHAR NOT NULL UNIQUE,
  phone_number VARCHAR NOT NULL UNIQUE,
  address TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Table `categories` phân loại sản phẩm, do 1 seller tạo
CREATE TABLE categories (
  id BIGSERIAL PRIMARY KEY,
  account_id BIGINT NOT NULL REFERENCES accounts(id) ON DELETE SET NULL,
  name VARCHAR NOT NULL,
  type VARCHAR NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Table `products` lưu sản phẩm do seller đăng bán
CREATE TABLE products (
  id BIGSERIAL PRIMARY KEY,
  account_id BIGINT NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
  category_id BIGINT NOT NULL REFERENCES categories(id) ON DELETE SET NULL,
  name VARCHAR NOT NULL,
  price INT NOT NULL CHECK (price > 0),
  discount_price INT NOT NULL CHECK (discount_price > 0 AND discount_price <= price),
  stock_quantity INT NOT NULL CHECK (stock_quantity >= 0),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Table `discounts` lưu khuyến mãi cho sản phẩm hoặc toàn tài khoản
CREATE TABLE discounts (
  id BIGSERIAL PRIMARY KEY,
  account_id BIGINT REFERENCES accounts(id) ON DELETE SET NULL,
  product_id BIGINT REFERENCES products(id) ON DELETE CASCADE,
  discount_value INT NOT NULL CHECK (discount_value > 0 AND discount_value < 100),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Table `cart_items` lưu giỏ hàng tạm cho buyer
CREATE TABLE cart_items (
  id BIGSERIAL PRIMARY KEY,
  account_id BIGINT NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
  product_id BIGINT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
  quantity INT NOT NULL DEFAULT 1 CHECK (quantity > 0),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (account_id, product_id)
);

-- Table `orders` lưu đơn hàng giữa buyer và seller
CREATE TABLE orders (
  id BIGSERIAL PRIMARY KEY,
  buyer_id BIGINT NOT NULL REFERENCES accounts(id) ON DELETE RESTRICT,
  seller_id BIGINT NOT NULL REFERENCES accounts(id) ON DELETE RESTRICT,
  total_price BIGINT NOT NULL CHECK (total_price >= 0),
  cod BOOLEAN NOT NULL DEFAULT TRUE,
  status VARCHAR NOT NULL CHECK (status IN ('pending','confirmed','shipped','delivered','cancelled')),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Table `order_items` chi tiết sản phẩm trong đơn hàng
CREATE TABLE order_items (
  id BIGSERIAL PRIMARY KEY,
  order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
  product_id BIGINT NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
  quantity INT NOT NULL CHECK (quantity > 0),
  price_each BIGINT NOT NULL CHECK (price_each >= 0)
);

-- Optional: table `shipments` lưu thông tin vận chuyển
CREATE TABLE shipments (
  id BIGSERIAL PRIMARY KEY,
  order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
  shipment_code VARCHAR NOT NULL UNIQUE,
  fee BIGINT NOT NULL CHECK (fee >= 0),
  status VARCHAR NOT NULL CHECK (status IN ('created','picked','in_transit','delivered','failed')),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);


-- Indexes tối ưu truy vấn
CREATE INDEX idx_products_account ON products(account_id);
CREATE INDEX idx_products_category ON products(category_id);
CREATE INDEX idx_orders_buyer ON orders(buyer_id);
CREATE INDEX idx_orders_seller ON orders(seller_id);
CREATE INDEX idx_order_items_order ON order_items(order_id);
CREATE INDEX idx_cart_items_account ON cart_items(account_id);

INSERT INTO accounts (username, hash_password, role, created_at) VALUES
  ('user1', 'hashedpass1', 'buyer', NOW()),
  ('user2', 'hashedpass2', 'seller', NOW()),
  ('user3', 'hashedpass3', 'admin', NOW());

-- Insert 3 account_info
INSERT INTO account_info (account_id, name, email, phone_number, address, created_at, updated_at) VALUES
  (1, 'User One', 'user1@example.com', '123456789', 'Address 1', NOW(), NOW()),
  (2, 'User Two', 'user2@example.com', '987654321', 'Address 2', NOW(), NOW()),
  (3, 'User Three', 'user3@example.com', '555555555', 'Address 3', NOW(), NOW());

-- Insert 2 categories
INSERT INTO categories (account_id, name, type, created_at) VALUES
  (2, 'Electronics', 'Type A', NOW()),
  (2, 'Clothing', 'Type B', NOW());

-- Insert 10 products
INSERT INTO products (account_id, category_id, name, price, discount_price, stock_quantity, created_at) VALUES
  (2, 1, 'Product 1', 1000, 900, 50, NOW()),
  (2, 1, 'Product 2', 1500, 1400, 30, NOW()),
  (2, 1, 'Product 3', 2000, 1800, 20, NOW()),
  (2, 2, 'Product 4', 500, 450, 100, NOW()),
  (2, 2, 'Product 5', 800, 750, 80, NOW()),
  (2, 1, 'Product 6', 1200, 1100, 60, NOW()),
  (2, 2, 'Product 7', 700, 650, 90, NOW()),
  (2, 1, 'Product 8', 1300, 1200, 40, NOW()),
  (2, 2, 'Product 9', 900, 850, 70, NOW()),
  (2, 1, 'Product 10', 1600, 1500, 25, NOW());

-- Insert 10 discounts
INSERT INTO discounts (account_id, product_id, discount_value, created_at) VALUES
  (2, 1, 10, NOW()),
  (2, 2, 15, NOW()),
  (2, 3, 20, NOW()),
  (2, 4, 5, NOW()),
  (2, 5, 8, NOW()),
  (2, 6, 12, NOW()),
  (2, 7, 18, NOW()),
  (2, 8, 14, NOW()),
  (2, 9, 9, NOW()),
  (2, 10, 7, NOW());