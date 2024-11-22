CREATE TABLE origin_destinations
(
    origin      UUID REFERENCES stop (id),
    destination UUID REFERENCES stop (id)
);

CREATE UNIQUE INDEX origin_destinations_origin_idx ON origin_destinations (origin);
