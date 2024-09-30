-- Active: 1727140703105@@vividly-climbing-mallard.data-1.use1.tembo.io@5432@armazem_DB@public
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

INSERT INTO books (title, synopsis, price, amount) VALUES ('Lucas', 'Venda muito', 5.90, 4);
INSERT INTO authors (name,nationality) VALUES ('Paulo Freire', 'Brasileiro'),('Roberto Marinho','Portugues');
INSERT INTO genres (name) VALUES ('Drama');
UPDATE books SET amount = amount - 4 WHERE id = 'daeab76d';

SELECT * FROM books;
SELECT * FROM inventory_logs;
SELECT * FROM authors;
UPDATE books SET author_id = 2 WHERE id = 'c7a0972e';

ALTER TABLE books ALTER COLUMN author_id SET ON DELETE RESTRICT ON UPDATE CASCADE;

CREATE TABLE books(  
    id CHAR(8) PRIMARY KEY DEFAULT substring(md5(random()::text) FROM 1 FOR 8),
    title VARCHAR(200) NOT NULL,
    synopsis VARCHAR(500) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    amount INT DEFAULT 0,
    author_id INT REFERENCES authors (id) ON DELETE RESTRICT/NO ACTION
);

CREATE TABLE authors(
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  nationality VARCHAR(100)
);

CREATE TABLE genres(
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL
);

CREATE TABLE genre_book (
  fk_book_id CHAR(8) REFERENCES books(id),
  fk_genre_id INT REFERENCES genres(id),
  PRIMARY KEY (fk_book_id,fk_genre_id)
);

CREATE TABLE inventory_logs (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  quantity INT,
  update_type VARCHAR(10),
  book_id CHAR(8) REFERENCES books(id) 
);

-- Procedure para atualizar os logs
CREATE OR REPLACE FUNCTION log_inventory_changes()
RETURNS TRIGGER AS $$
BEGIN
    -- Determinar o tipo de atualização
    IF TG_OP = 'INSERT' THEN
        -- Se é um INSERT, é sempre uma entrada
        INSERT INTO inventory_logs (created_at, quantity, update_type, book_id)
        VALUES (CURRENT_TIMESTAMP, NEW.amount, 'entrada', NEW.id);
    ELSIF TG_OP = 'UPDATE' THEN
        IF NEW.amount > OLD.amount THEN
            INSERT INTO inventory_logs (created_at, quantity, update_type, book_id)
            VALUES (CURRENT_TIMESTAMP, NEW.amount - OLD.amount, 'entrada', NEW.id);
        ELSIF NEW.amount < OLD.amount THEN
            INSERT INTO inventory_logs (created_at, quantity, update_type, book_id)
            VALUES (CURRENT_TIMESTAMP, OLD.amount - NEW.amount, 'saída', NEW.id);
        END IF;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

 
CREATE TRIGGER track_inventory_changes
AFTER INSERT OR UPDATE ON books
FOR EACH ROW
EXECUTE FUNCTION log_inventory_changes();



CREATE OR REPLACE PROCEDURE add_book_with_genre(
  book_title VARCHAR,
  book_synopsis VARCHAR,
  book_price DECIMAL,
  book_amount INT,
  genre_id INT 
)
LANGUAGE plpgsql
AS $$
DECLARE
  new_book_id CHAR(8);
BEGIN 

  new_book_id := substring(md5(random()::text)FROM 1 for 8)

  INSERT INTO book (book_id, book_title, book_price, book_synopsis,)
  VALUES (new_book_id, book_title, book_price, book_amount, genre_id)
  RETURNING book_id INTO new_book_id;

  INSERT INTO genre_book (fk_book_id, fk_genre_id)
  VALUES (new_book_id, genre_id);
END;
$$;
