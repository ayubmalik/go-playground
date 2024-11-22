CREATE TABLE stop_summary
(
    id           UUID PRIMARY KEY,
    station_name VARCHAR(50),
    station_code CHAR(10) UNIQUE,
    city_name    VARCHAR(50),
    state_code   CHAR(2)
);
