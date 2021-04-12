-- create database ysnpkoyaDB
-- 	with owner postgres
-- 	encoding 'utf8'
-- 	LC_COLLATE = 'ru_RU.UTF-8'
--     LC_CTYPE = 'ru_RU.UTF-8'
--     TABLESPACE = pg_default
--     TEMPLATE template0
-- 	;

CREATE EXTENSION postgis;
CREATE EXTENSION postgis_topology;


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
    latitude  float                       DEFAULT 55.753808,
    longitude float                       DEFAULT 37.620017,
    radius    int                         DEFAULT 0,
    address   varchar(128)                DEFAULT 'Москва',
    avatar    varchar(128)       NOT NULL DEFAULT ''
);

create table if not exists category
(
    id    serial primary key,
    title varchar(128) unique not null
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

    FOREIGN KEY (product_id) REFERENCES product (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS user_favorite
(
    user_id    int NOT NULL,
    product_id int NOT NULL,

    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES product (id) ON DELETE CASCADE
);

INSERT INTO users (email, telephone, password, name, surname, sex)
VALUES ('asd', '123', '123', '123', '123', 'M');


INSERT INTO category (title)
VALUES ('Транспорт'),
       ('Недвижмость'),
       ('Хобби и отдых'),
       ('Работа'),
       ('Для дома и дачи'),
       ('Бытовая электрика'),
       ('Личные вещи'),
       ('Животные');


INSERT INTO product (name, amount, description, category_id, owner_id, address, longitude, latitude) VALUES
    ('iPhone 10', 1000, 'hello', 1, 1, 'Москва', 37.620017, 55.753808),
    ('iPhone 11', 1200, 'hello', 2, 1, 'Москва', 37.620017, 55.753808),
    ('iPhone 12', 1300, 'hello', 3, 1, 'Москва', 37.620017, 55.753808),
    ('iPhone 13', 1400, 'hello', 4, 1, 'Москва', 37.620017, 55.753808),
    ('iPhone 14', 1500, 'hello', 5, 1, 'Москва', 37.620017, 55.753808),
    ('iPhone 15', 1600, 'hello', 6, 1, 'Москва', 37.620017, 55.753808),
    ('iPhone 16', 1700, 'hello', 7, 1, 'Москва', 37.620017, 55.753808),
    ('iPhone 17', 1800, 'hello', 8, 1, 'Москва', 37.620017, 55.753808),
    ('iPhone 18', 1900, 'hello', 8, 1, 'Москва', 37.620017, 55.753808),
    ('iPhone 19', 2100, 'hello', 1, 1, 'Москва', 37.620017, 55.753808),
    ('iPhone 20', 2400, 'hello', 2, 1, 'Москва', 37.620017, 55.753808);




INSERT INTO product_images (product_id, img_link)
VALUES (1, '/static/product/2e5659cd-72ac-43d8-8494-52bbc7a885fd.webp'),
       (2, '/static/product/3af8506c-0608-498e-b2b7-7f7a445aa6df.webp'),
       (3, '/static/product/4abcea38-00ad-4365-85af-1c144085ebd2.webp'),
       (4, '/static/product/6d835ba7-1ecc-478d-8832-a64b3c58124c.webp'),
       (5, '/static/product/8b644046-55b7-40a2-beab-308a964630ab.jpg'),
       (6, '/static/product/697bade2-a4cb-49fc-bad3-c2205554b92a.jpeg'),
       (7, '/static/product/936de281-1bdb-46e5-a404-3bf2a3fdbaac.webp'),
       (8, '/static/product/ba1b1a47-97d3-4efb-aed2-f574fc28970f.webp'),
       (9, '/static/product/dfc9f3d6-60cd-480f-97d1-c31c52dca48b.webp'),
       (10, '/static/product/f75694ab-42a6-42f7-8f24-cd933cd4da2e.webp'),
       (11, '/static/product/8776ad39-e754-4f29-8d19-640b1543fbfe.jpg');
