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

create table if not exists category
(
    id        serial        primary key,
    title     varchar(128)  unique not null
);

CREATE TABLE IF NOT EXISTS product
(
    id          serial PRIMARY KEY,
    name        varchar(128) NOT NULL,
    date        date         not null default '1970-1-1',
    amount      int          not null,
    description text         NOT NULL,
    category_id int          not null,
    owner_id    int          not null,
    address     varchar(128),
    longitude   float,
    latitude    float,
    likes       int                   DEFAULT 0, -- триггер на каждый лайк/дизлайк
    views       int                   DEFAULT 0,
    tariff      int                   DEFAULT 0,

    FOREIGN KEY (owner_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES category (id) ON DELETE NO ACTION
);

CREATE TABLE IF NOT EXISTS product_images
(
    product_id int                 NOT NULL,
    img_link   varchar(128) unique NOT NULL,

FOREIGN KEY (product_id) REFERENCES product(id) ON DELETE CASCADE
);


INSERT INTO users (email, telephone, password, name, surname, sex) VALUES ('asd', '123', '123', '123', '123', 'M');

INSERT INTO category (title) VALUES 
    ('Транспорт'),
    ('Недвижмость'),
    ('Хобби и отдых'),
    ('Работа'),
    ('Для дома и дачи'),
    ('Бытовая электрика'),
    ('Личные вещи'),
    ('Животные');

INSERT INTO product (name, amount, description, category_id, owner_id, address, longitude, latitude) VALUES
    ('iPhone 10', 1000, 'hello', 1, 1, 'Москва', 55.753808, 37.620017),
    ('iPhone 11', 1200, 'hello', 2, 1, 'Москва', 55.753808, 37.620017),
    ('iPhone 12', 1300, 'hello', 3, 1, 'Москва', 55.753808, 37.620017),
    ('iPhone 13', 1400, 'hello', 4, 1, 'Москва', 55.753808, 37.620017),
    ('iPhone 14', 1500, 'hello', 5, 1, 'Москва', 55.753808, 37.620017),
    ('iPhone 15', 1600, 'hello', 6, 1, 'Москва', 55.753808, 37.620017),
    ('iPhone 16', 1700, 'hello', 7, 1, 'Москва', 55.753808, 37.620017),
    ('iPhone 17', 1800, 'hello', 8, 1, 'Москва', 55.753808, 37.620017),
    ('iPhone 18', 1900, 'hello', 8, 1, 'Москва', 55.753808, 37.620017),
    ('iPhone 19', 2100, 'hello', 1, 1, 'Москва', 55.753808, 37.620017),
    ('iPhone 20', 2400, 'hello', 2, 1, 'Москва', 55.753808, 37.620017);


    INSERT INTO product_images (product_id, img_link) VALUES 
    (1, 'asd2'),
        (2, 'as3d'),
            (3, 'as4d'),
                (4, 'as5d'),
                    (5, 'a6sd'),
                        (6, 'as7d'),
                            (7, 'as8d'),
                                (8, 'as9d'),
                    (9, 'a52sd'),
                        (10, 'a32sd'),
                            (11, 'as43d');


