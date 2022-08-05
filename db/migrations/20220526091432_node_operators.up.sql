create table if not exists node_operators
(
    id                     bigserial PRIMARY KEY NOT NULL,
    index                  integer   not null,
    block_number           integer   not null,
    tx_hash                text      not null,
    pubkey                 text      not null,
    withdrawal_credentials text      not null,

    UNIQUE (index)
);
