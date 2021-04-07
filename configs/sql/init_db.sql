-- create database ysnpkoyaDB
-- 	with owner postgres
-- 	encoding 'utf8'
-- 	LC_COLLATE = 'ru_RU.UTF-8'
--     LC_CTYPE = 'ru_RU.UTF-8'
--     TABLESPACE = pg_default
--     TEMPLATE template0
-- 	;
GRANT ALL PRIVILEGES ON database ysnpkoyaDB TO postgres;

create table if not exists users
(
    id        serial primary key,
    email     varchar(64)        not null,
    telephone varchar(12) unique not null,
    password  text               not null,
    name      varchar(64)        not null,
    surname   varchar(64)        not null,
    sex       varchar(12)        not null,
    birthdate date,
    reg_date  timestamp,
    avatar    varchar(128)       NOT NULL DEFAULT ''
);


CREATE TABLE IF NOT EXISTS product
(
    id          serial PRIMARY KEY,
    name        varchar(128) NOT NULL,
    date        date         not null default '1970-1-1',
    amount      int          not null,
    description text         NOT NULL,
    category    varchar(64)  not null,
    owner_id    int          not null,
    address varchar(128),
    longitude varchar(64),
    latitude varchar(64),
    likes       int                   DEFAULT 0, -- триггер на каждый лайк/дизлайк
    views       int                   DEFAULT 0,
    tariff      int                   DEFAULT 0,

    FOREIGN KEY (owner_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS product_images
(
    product_id int                 NOT NULL,
    img_link   varchar(128) unique NOT NULL,

    FOREIGN KEY (product_id) REFERENCES product (id) ON DELETE CASCADE
);