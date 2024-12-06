CREATE TABLE origin_destination
(
    origin      UUID REFERENCES stop_summary (id),
    destination UUID REFERENCES stop_summary (id)
);

CREATE UNIQUE INDEX origin_destination_idx ON origin_destination (origin, destination);
