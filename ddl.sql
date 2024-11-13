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
    amount       INTEGER   DEFAULT 0 CHECK ( amount >= 0 ),
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    fk_author_id INTEGER references author (id) ON DELETE RESTRICT 
);

-- Trigger to update `updated_at` on book updates
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
    fk_user_id    INTEGER REFERENCES user_account (id) ON DELETE RESTRICT, 
    fk_admin_id   INTEGER REFERENCES user_account (id) ON DELETE SET NULL,
    fk_book_id    INTEGER REFERENCES book (id) ON DELETE CASCADE
);

-- Trigger to decrement book stock on reservation creation
CREATE OR REPLACE FUNCTION decrement_book_on_reservation_create()
    RETURNS TRIGGER AS
$$
BEGIN
    -- Lock the row for the book to prevent race conditions
    SELECT amount
    FROM book
    WHERE id = NEW.fk_book_id
        FOR UPDATE;

    -- Attempt to decrement the book amount conditionally
    UPDATE book
    SET amount = amount - 1
    WHERE id = NEW.fk_book_id
      AND amount > 0
    RETURNING amount;

    -- If no row was updated (i.e., stock is 0), raise an exception
    IF NOT FOUND THEN
        RAISE EXCEPTION 'Book is out of stock, cannot reserve.';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER reservation_create
    AFTER INSERT
    ON reservation
    FOR EACH ROW
EXECUTE FUNCTION decrement_book_on_reservation_create();

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

-- Trigger to handle changes in reservation status
CREATE OR REPLACE FUNCTION create_loan_or_increment_book_on_reservation_status_change()
    RETURNS TRIGGER AS
$$
BEGIN
    -- Handle reservation status change to 'collected'
    IF NEW.status = 'collected' AND OLD.status != 'collected' THEN
        -- Insert a new loan record
        INSERT INTO loan (return_by, fk_reservation_id)
        VALUES (CURRENT_TIMESTAMP + NEW.borrowed_days * INTERVAL '1 day', NEW.id);

        -- Handle reservation status change to 'cancelled' or 'expired'
    ELSIF NEW.status IN ('cancelled', 'expired') AND OLD.status != NEW.status THEN
        -- Increase the book amount by 1
        UPDATE book
        SET amount = amount + 1
        WHERE id = NEW.fk_book_id;
    END IF;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER reservation_status_change
    AFTER UPDATE
    ON reservation
    FOR EACH ROW
    WHEN (NEW.status IS DISTINCT FROM OLD.status)
EXECUTE FUNCTION create_loan_or_increment_book_on_reservation_status_change();

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
    fk_reservation_id INTEGER REFERENCES reservation (id) ON DELETE CASCADE
);

-- Trigger to increment book stock and update reservation status on loan return
CREATE OR REPLACE FUNCTION increment_book_on_loan_return_and_finish_reservation()
    RETURNS TRIGGER AS
$$
BEGIN
    -- Check if the loan status is changed to 'returned'
    IF NEW.status = 'returned' AND OLD.status != 'returned' THEN
        -- Increment the book amount by 1
        UPDATE book
        SET amount = amount + 1
        WHERE id = (SELECT fk_book_id FROM reservation WHERE id = NEW.fk_reservation_id);

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
