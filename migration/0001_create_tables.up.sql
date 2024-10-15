CREATE TABLE USERS (
                       id SERIAL PRIMARY KEY,
                       FIRST_NAME VARCHAR(50) NOT NULL,
                       LAST_NAME VARCHAR(50) NOT NULL
);

CREATE TABLE BOOKS (
                       id SERIAL PRIMARY KEY,
                       TITLE VARCHAR(50) NOT NULL,
                       QUANTITY INT NOT NULL CHECK (QUANTITY >= 0),
                       BORROWED_COUNT INT DEFAULT 0 CHECK (BORROWED_COUNT >= 0)
);

ALTER TABLE BOOKS
    ADD CONSTRAINT chk_borrowed_count_quantity CHECK (BORROWED_COUNT <= QUANTITY);

CREATE TABLE BORROW (
                        id SERIAL PRIMARY KEY,
                        USER_ID INT NOT NULL REFERENCES USERS(id),
                        BOOK_ID INT NOT NULL REFERENCES BOOKS(id),
                        BORROWED_AT TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        RETURNED_AT TIMESTAMP
);