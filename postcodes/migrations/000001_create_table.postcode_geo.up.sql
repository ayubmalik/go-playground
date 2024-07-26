create table postcode_geo
(
    id       integer       not null,
    postcode varchar(8)    not null,
    lat      numeric(9, 6) not null,
    lng      numeric(9, 6) not null,

    constraint pk_postcode_geo primary key (id)
);

create index idx_postcode_geo_lat_lng on postcode_geo (lat, lng);