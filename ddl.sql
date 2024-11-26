-- ===========================
-- 1. Author and Genre Tables
-- ===========================

-- Author Table
CREATE TABLE IF NOT EXISTS author
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Genre Table
CREATE TABLE IF NOT EXISTS genre
(
    id   SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL
);

-- ===========================
-- 2. Book Tables
-- ===========================

-- Book Table
CREATE TABLE IF NOT EXISTS book
(
    id           SERIAL PRIMARY KEY,
    title        VARCHAR(200) UNIQUE NOT NULL,
    synopsis     TEXT                NOT NULL,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    fk_author_id INTEGER references author (id) ON DELETE RESTRICT
);

CREATE TYPE book_stock_status AS ENUM ('available', 'borrowed', 'missing');

CREATE TABLE IF NOT EXISTS book_stock
(
    id         SERIAL PRIMARY KEY,
    status     book_stock_status DEFAULT 'available',
    code       INTEGER NOT NULL UNIQUE,
    created_at TIMESTAMP         DEFAULT CURRENT_TIMESTAMP,
    fk_book_id INTEGER REFERENCES book (id) ON DELETE CASCADE
);

CREATE OR REPLACE FUNCTION prevent_book_stock_delete_if_borrowed()
    RETURNS TRIGGER AS
$$
BEGIN
    IF OLD.status = 'borrowed' THEN
        RAISE EXCEPTION 'Cannot delete book_stock unless the status is ''available'' or ''missing''.';
    END IF;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER before_stock_delete_check_status
    BEFORE DELETE
    on book_stock
    FOR EACH ROW
EXECUTE FUNCTION prevent_book_stock_delete_if_borrowed();

-- Function to prevent deletion if active reservations exceed available stock
CREATE OR REPLACE FUNCTION prevent_book_stock_delete_if_reserved_exceeds_available()
    RETURNS TRIGGER AS
$$
DECLARE
    stock_and_reservations_count RECORD;
BEGIN
    -- Single query to get both total stock count and active reservations count
    SELECT COUNT(bs.id) AS total_stock_count,
           COUNT(r.id)  AS active_reservations_count
    INTO stock_and_reservations_count
    FROM book_stock bs
             LEFT JOIN reservation r
                       ON r.fk_book_id = bs.fk_book_id
                           AND r.status = 'pending'
                           AND r.expires_at > CURRENT_TIMESTAMP
    WHERE bs.id = OLD.id AND bs.status = 'available';

    -- Prevent deletion if active reservations >= available stock
    IF stock_and_reservations_count.active_reservations_count >= stock_and_reservations_count.total_stock_count THEN
        RAISE EXCEPTION 'Cannot delete book_stock as active reservations exceed or equal available stock.';
    END IF;

    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

-- Trigger to check reservations before deleting book_stock
CREATE OR REPLACE TRIGGER before_book_stock_delete_check_reservations_and_stock
    BEFORE DELETE
    ON book_stock
    FOR EACH ROW
EXECUTE FUNCTION prevent_book_stock_delete_if_reserved_exceeds_available();


CREATE OR REPLACE FUNCTION update_book_timestamp()
    RETURNS TRIGGER AS
$$
BEGIN
    -- Set the `updated_at` field to the current timestamp
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER book_update_timestamp
    BEFORE UPDATE
    ON book
    FOR EACH ROW
EXECUTE FUNCTION update_book_timestamp();

-- Book-Genre Relationship Table
CREATE TABLE IF NOT EXISTS book_genre
(
    fk_book_id  INTEGER NOT NULL REFERENCES book (id) ON DELETE CASCADE,
    fk_genre_id INTEGER NOT NULL REFERENCES genre (id) ON DELETE CASCADE,
    PRIMARY KEY (fk_book_id, fk_genre_id)
);

-- ===========================
-- 3. User Management Tables
-- ===========================

-- Account Role Table
CREATE TABLE IF NOT EXISTS account_role
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(30) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- User Account Table
CREATE TABLE IF NOT EXISTS user_account
(
    id              SERIAL PRIMARY KEY,
    name            VARCHAR(150)        NOT NULL,
    cpf             CHAR(11) UNIQUE     NOT NULL,
    phone           VARCHAR(20) UNIQUE  NOT NULL,
    email           VARCHAR(150) UNIQUE NOT NULL,
    password_hash   VARCHAR(255)        NOT NULL,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_active       BOOLEAN   DEFAULT TRUE,
    fk_account_role INTEGER REFERENCES account_role (id) ON DELETE RESTRICT
);

-- Trigger to update `updated_at` on user account updates
CREATE OR REPLACE FUNCTION update_user_account_timestamp()
    RETURNS TRIGGER AS
$$
BEGIN
    -- Set the `updated_at` field to the current timestamp
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER user_account_update_timestamp
    BEFORE UPDATE
    ON user_account
    FOR EACH ROW
EXECUTE FUNCTION update_user_account_timestamp();

-- ===========================
-- 4. Reservation and Loan Tables
-- ===========================

-- Reservation Status Enum Type
CREATE TYPE reservation_status AS ENUM ('pending', 'cancelled', 'expired', 'collected', 'finished');

-- Reservation Table
CREATE TABLE IF NOT EXISTS reservation
(
    id            SERIAL PRIMARY KEY,
    reserved_at   TIMESTAMP          DEFAULT CURRENT_TIMESTAMP,
    expires_at    TIMESTAMP GENERATED ALWAYS AS ( reserved_at + INTERVAL '24 hours') STORED,
    borrowed_days INTEGER NOT NULL CHECK ( borrowed_days <= 90 ),
    status        reservation_status DEFAULT 'pending',
    fk_user_id    INTEGER REFERENCES user_account (id) ON DELETE CASCADE,
    fk_admin_id   INTEGER REFERENCES user_account (id) ON DELETE SET NULL,
    fk_book_id    INTEGER REFERENCES book (id) ON DELETE CASCADE
);

-- Trigger to prevent deletion of active reservations
CREATE OR REPLACE FUNCTION prevent_reservation_delete_if_active()
    RETURNS TRIGGER AS
$$
BEGIN
    IF OLD.status IN ('pending', 'collected') AND OLD.expires_at > CURRENT_TIMESTAMP THEN
        RAISE EXCEPTION 'Cannot delete reservation unless the status is ''expired'',''cancelled'' or ''finished''.';
    END IF;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER before_reservation_delete_check_status
    BEFORE DELETE
    ON reservation
    FOR EACH ROW
EXECUTE FUNCTION prevent_reservation_delete_if_active();

-- ===========================
-- 5. Loan Tables
-- ===========================

-- Loan Status Enum Type
CREATE TYPE loan_status AS ENUM ('borrowed', 'returned');

-- Loan Table
CREATE TABLE IF NOT EXISTS loan
(
    id                SERIAL PRIMARY KEY,
    loaned_at         TIMESTAMP   DEFAULT CURRENT_TIMESTAMP,
    return_by         TIMESTAMP NOT NULL,
    returned_at       TIMESTAMP CHECK (returned_at IS NULL OR returned_at > loaned_at),
    status            loan_status DEFAULT 'borrowed',
    fk_admin_id       INTEGER   REFERENCES user_account (id) ON DELETE SET NULL,
    fk_book_stock_id  INTEGER REFERENCES book_stock (id) ON DELETE CASCADE,
    fk_reservation_id INTEGER REFERENCES reservation (id) ON DELETE CASCADE
);

-- Trigger to increment book stock and update reservation status on loan return
CREATE OR REPLACE FUNCTION increment_book_on_loan_return_and_finish_reservation()
    RETURNS TRIGGER AS
$$
BEGIN
    -- Check if the loan status is changed to 'returned'
    IF NEW.status = 'returned' AND OLD.status != 'returned' THEN
        UPDATE book_stock
        SET status = 'available'
        WHERE id = NEW.fk_book_stock_id;

        -- Update the reservation status to 'finished'
        UPDATE reservation
        SET status = 'finished'
        WHERE id = NEW.fk_reservation_id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER loan_status_change
    AFTER UPDATE
    ON loan
    FOR EACH ROW
    WHEN (NEW.status IS DISTINCT FROM OLD.status)
EXECUTE FUNCTION increment_book_on_loan_return_and_finish_reservation();

-- Trigger to prevent loan deletion when its status is 'returned'
CREATE OR REPLACE FUNCTION prevent_loan_delete_if_active()
    RETURNS TRIGGER AS
$$
BEGIN
    IF OLD.status != 'returned' THEN
        RAISE EXCEPTION 'Cannot delete loan unless the status is ''returned''.';

    END IF;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER before_loan_delete_check_status
    BEFORE DELETE
    ON loan
    FOR EACH ROW
EXECUTE FUNCTION prevent_loan_delete_if_active()
