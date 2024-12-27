CREATE TABLE origin_destination
(
    origin      CHAR(36) REFERENCES stop_summary (id),
    destination CHAR(36) REFERENCES stop_summary (id)
);

CREATE UNIQUE INDEX origin_destination_idx ON origin_destination (origin, destination);
