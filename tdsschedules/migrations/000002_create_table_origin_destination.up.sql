CREATE TABLE origin_destinations
(
    origin      UUID REFERENCES stop_summary (id),
    destination UUID REFERENCES stop_summary (id)
);

CREATE UNIQUE INDEX origin_destinations_idx ON origin_destinations (origin, destination);
