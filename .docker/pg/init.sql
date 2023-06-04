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

CREATE TABLE users (
                       username VARCHAR(50) NOT NULL PRIMARY KEY,
                       password VARCHAR(100) NOT NULL
);

INSERT INTO users (username, password) VALUES ('alice', '$2a$10$K01poobswi.hbfUiJG9go.kQRgpZ6kWRlJFJrjH1/iiqJJJEwlTv6');
INSERT INTO users (username, password) VALUES ('bob', '$2a$10$J2uuNk6KC4l/qia/XcxDBOJkTM7EqHusIukrQY4KfET617cy1rPpW');
INSERT INTO users (username, password) VALUES ('carol', '$2a$10$/K5q5phzlps4IA84EgsMB.tGVrA7UDszkndI0.d80lMLEc6hrnEr6');
INSERT INTO users (username, password) VALUES ('dave', '$2a$10$h/Nbf8Rezyv3etqwBq7cSeg/IBN.YDJ0lP54fKQ95XqLbPx79W7bK');


CREATE TABLE transaction_history (
                                     id SERIAL PRIMARY KEY,
                                     username VARCHAR(50) NOT NULL,
                                     transaction_hash CHAR(66) NOT NULL
);
CREATE INDEX idx_username ON transaction_history (username);