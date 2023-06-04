\c postgres;

CREATE TABLE transactions (
                              id SERIAL,
                              transaction_hash CHAR(66) NOT NULL UNIQUE PRIMARY KEY,
                              transaction_status SMALLINT NOT NULL,
                              block_hash CHAR(66) NOT NULL,
                              block_number INTEGER NOT NULL,
                              sender CHAR(42) NOT NULL ,
                              recipient CHAR(42),
                              contract_address VARCHAR(42),
                              logs_count INTEGER NOT NULL,
                              input TEXT NOT NULL,
                              value BIGINT NOT NULL
);