--
-- Microservice: Product Dictionary Service
--
-- Setup schema
-- The default COLLATE for utm8mb4 is utf8mb4_0900_ai_ci. No need to be explicit.
CREATE SCHEMA IF NOT EXISTS products CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;

-- Switch to using products schema
USE products;

CREATE TABLE categories (
  category_id      BINARY(16) NOT NULL PRIMARY KEY,
  category_id_text VARCHAR(36) generated always AS
   (insert(
      insert(
        insert(
          insert(hex(category_id),9,0,'-'),
          14,0,'-'),
        19,0,'-'),
      24,0,'-')
   ) virtual,
  name            varchar(64) NOT NULL,
  created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at      DATETIME,
  UNIQUE KEY uq__categories__name (name)
)
ENGINE=InnoDB
DEFAULT CHARSET=utf8mb4
COMMENT='Categories that a product may fall under';


CREATE TABLE products (
  product_id            BINARY(16) NOT NULL PRIMARY KEY,
  product_id_text       varchar(36) generated always AS
   (insert(
      insert(
        insert(
          insert(hex(product_id),9,0,'-'),
          14,0,'-'),
        19,0,'-'),
      24,0,'-')
   ) virtual,
  category_id           BINARY(16) NOT NULL,
  name                  varchar(64) NOT NULL,
  created_at            DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at            DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at            DATETIME,
  UNIQUE KEY uq__products__name (name),
  INDEX (product_id, category_id),
  CONSTRAINT fk__products__category_id  FOREIGN KEY (category_id) REFERENCES categories (category_id) ON DELETE CASCADE
)
ENGINE=InnoDB
DEFAULT CHARSET=utf8mb4
COMMENT='A service or product that caring markets for its b2b customers';
