CREATE TABLE currencies (
    code text PRIMARY KEY
);

CREATE TABLE api_keys (
    project_id text PRIMARY KEY,
    hashed_secret text NOT NULL
);

CREATE TABLE exchange_rates (
    from_currency text NOT NULL REFERENCES currencies,
    to_currency text NOT NULL REFERENCES currencies,
    rate decimal NOT NULL,
    rateAt timestamp NOT NULL,
    PRIMARY KEY(from_currency, to_currency)
);
