CREATE TABLE stop
(
    uuid         UUID PRIMARY KEY,
    station_name VARCHAR(50) UNIQUE,
    city         INTEGER,
    state        CHAR(2) REFERENCES state(abbreviation)
);
