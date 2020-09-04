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
VALUES ('126bea59-c9a7-44d0-bcd8-d710aad69676', 'bruno', '4321', 240.10);