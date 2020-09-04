CREATE TABLE users
(
    ID       UUID PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    balance  DOUBLE PRECISION
);

CREATE TABLE transactions
(
    ID             UUID PRIMARY KEY,
    source_user_id UUID REFERENCES users (ID)  NOT NULL,
    target_user_id UUID REFERENCES users (ID)  NOT NULL,
    amount         DOUBLE PRECISION            NOT NULL,
    created_at     TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now()
);

INSERT INTO users
VALUES ('256bea59-c9a7-44d0-bcd8-d710aad69676', 'breno', '1234', 10);

INSERT INTO users
VALUES ('c66af437-8536-4ac9-918c-5e73ef95578a', 'bruno', '4321', 100);

INSERT INTO users
VALUES ('9e321e7b-918b-4bef-9c85-81b1729b31d9', 'brono', 'abcd', 1000);

INSERT INTO users
VALUES ('007dcaec-6963-4d4c-a40d-9b5eda420f10', 'brano', 'abcdef', 10000);