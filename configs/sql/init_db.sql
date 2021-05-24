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

ALTER USER postgres WITH PASSWORD 'ysnpkoyapassword';

create table if not exists users
(
    id        serial primary key,
    email     varchar(64)           default NULL,
    telephone varchar(12) unique    default NULL,
    password  text                  default NULL,
    name      varchar(64)  not null,
    surname   varchar(64)  not null,
    sex       varchar(12)           default 'notstated',
    birthdate date                  default NULL,
    reg_date  timestamp,
    latitude  float                 DEFAULT 55.753808,
    longitude float                 DEFAULT 37.620017,
    radius    int                   DEFAULT 0,
    address   varchar(128)          DEFAULT 'Москва',
    avatar    varchar(512) NOT NULL DEFAULT ''
);

create table if not exists users_oauth
(
    user_id    int unique  not null,
    oauth_type varchar(20) not null,
    oauth_id   float       not null,

    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
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
    close       boolean               DEFAULT false,
    buyer_id    int                   DEFAULT null,
    buyer_left_review boolean         DEFAULT null,
    seller_left_review boolean        DEFAULT null,

    FOREIGN KEY (owner_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES category (id) ON DELETE NO ACTION
);

CREATE INDEX idx_product ON product USING GIST
    (Geography(ST_SetSRID(ST_POINT(longitude, latitude), 4326)));


CREATE OR REPLACE FUNCTION check_achievement() RETURNS TRIGGER AS
$check_achievement$
DECLARE
    P_COUNT INTEGER;
BEGIN
    SELECT COUNT(*) from product where owner_id = NEW.owner_id INTO P_COUNT;
    IF (P_COUNT = 1) THEN
        INSERT INTO user_achievement (user_id, a_id) VALUES (NEW.owner_id, 1);
    end if;

    IF (P_COUNT = 10) THEN
        INSERT INTO user_achievement (user_id, a_id) VALUES (NEW.owner_id, 2);
    end if;

    IF (P_COUNT = 100) THEN
        INSERT INTO user_achievement (user_id, a_id) VALUES (NEW.owner_id, 3);
    end if;
    
    RETURN NEW; -- возвращаемое значение для триггера AFTER игнорируется
END;
$check_achievement$ LANGUAGE plpgsql;

CREATE TRIGGER check_achievement
    AFTER INSERT
    ON product
    FOR EACH ROW
EXECUTE PROCEDURE check_achievement();


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

create table if not exists chats
(
    id               serial primary key,
    creation_time    timestamp,

    last_msg_id      integer            default 0,    -- trigger
    last_msg_content text               default '',   -- trigger
    last_msg_time    timestamp not null default NOW() -- trigger
);



create table if not exists user_chats
(
    user_id          integer not null,
    partner_id       integer not null,
    product_id       integer not null,
    chat_id          integer not null,
    last_read_msg_id int default 0,

    new_messages     int default 0, -- trigger

    primary key (user_id, partner_id, product_id),
    foreign key (user_id) references users (id) on delete cascade,
    foreign key (partner_id) references users (id) on delete cascade,
    foreign key (product_id) references product (id) on delete cascade,
    foreign key (chat_id) references chats (id) on delete cascade
);


create table if not exists messages
(
    id            serial primary key,
    content       text      not null,
    creation_time timestamp not null default NOW(),
    chat_id       integer   not null,
    user_id       integer   not null,

    foreign key (user_id) references users (id) on delete cascade,
    foreign key (chat_id) references chats (id) on delete cascade
);

create table if not exists reviews
(
    id               serial primary key,
    creation_time    timestamp
);

create table if not exists user_reviews
(
  review_id        int      not null,
  content          text     not null,
  rating           float    not null,
  reviewer_id      int      not null,
  product_id       int      not null,
  target_id        int      not null ,
  type             varchar(12) not null,

  primary key (reviewer_id, product_id),
  foreign key (reviewer_id) references users (id) on delete cascade,
  foreign key (product_id) references product (id) on delete cascade,
  foreign key (review_id) references reviews(id) on delete cascade,
  foreign key (target_id) references users(id) on delete cascade
);

CREATE OR REPLACE FUNCTION msg_change() RETURNS TRIGGER AS
$msg_change$
BEGIN
    IF (TG_OP = 'INSERT') THEN
        UPDATE user_chats
        SET new_messages = new_messages + 1
        where user_chats.chat_id = NEW.chat_id
          AND user_chats.partner_id = NEW.user_id;

        UPDATE user_chats
        SET last_read_msg_id = NEW.id
        where user_chats.chat_id = NEW.chat_id
          AND user_chats.user_id = NEW.user_id;

        UPDATE chats
        SET last_msg_id      = NEW.id,
            last_msg_content = NEW.content,
            last_msg_time    = NEW.creation_time
        WHERE chats.id = NEW.chat_id;

    END IF;
    RETURN NULL; -- возвращаемое значение для триггера AFTER игнорируется
END;
$msg_change$ LANGUAGE plpgsql;


CREATE TRIGGER upd_msgs
    AFTER INSERT
    ON messages
    FOR EACH ROW
EXECUTE PROCEDURE msg_change();

INSERT INTO category (title)
VALUES ('Транспорт'),
       ('Недвижмость'),
       ('Работа'),
       ('Услуги'),
       ('Личные вещи'),
       ('Для дома и дачи'),
       ('Бытовая электрика'),
       ('Хобби и отдых'),
       ('Животные');



create table if not exists achievement
(
    id             serial primary key,
    title          text      not null,
    description    text      not null,
    link_pic       text      not null
);

INSERT INTO achievement (title, description, link_pic)
VALUES ('Транспорт', 'Транспорт', 'https://achievement-images.teamtreehouse.com/badge_javascript-array-iteration-methods_stage01.png'),
       ('Недвижмость','Транспорт', 'https://achievement-images.teamtreehouse.com/badge_javascript-array-iteration-methods_stage01.png'),
       ('Работа','Транспорт', 'https://achievement-images.teamtreehouse.com/badge_javascript-array-iteration-methods_stage01.png'),
       ('Услуги', 'Транспорт', 'https://achievement-images.teamtreehouse.com/badge_javascript-array-iteration-methods_stage01.png'),
       ('Личные вещи', 'Транспорт', 'https://achievement-images.teamtreehouse.com/badge_javascript-array-iteration-methods_stage01.png'),
       ('Для дома и дачи', 'Транспорт', 'https://achievement-images.teamtreehouse.com/badge_javascript-array-iteration-methods_stage01.png'),
       ('Бытовая электрика', 'Транспорт', 'https://achievement-images.teamtreehouse.com/badge_javascript-array-iteration-methods_stage01.png'),
       ('Хобби и отдых', 'Транспорт', 'https://achievement-images.teamtreehouse.com/badge_javascript-array-iteration-methods_stage01.png'),
       ('Животные', 'Транспорт', 'https://achievement-images.teamtreehouse.com/badge_javascript-array-iteration-methods_stage01.png');

create table if not exists user_achievement
(
    id             serial primary key,
    user_id        integer   not null,
    date           timestamp not null default NOW(),
    a_id           integer   not null,

    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE NO ACTION,
    FOREIGN KEY (a_id) REFERENCES achievement (id) ON DELETE NO ACTION
);


-- INSERT INTO users (email, telephone, password, name, surname, sex)
-- VALUES ('asd', '123', '123', '123', '123', 'M');

-- INSERT INTO product (name, amount, description, category_id, owner_id, address, longitude, latitude)
-- VALUES ('iPhone 10', 1000, 'hello', 1, 1, 'Москва', 37.620017, 55.753808),
--        ('iPhone 11', 1200, 'hello', 2, 1, 'Москва', 37.620017, 55.753808),
--        ('iPhone 12', 1300, 'hello', 3, 1, 'Москва', 37.620017, 55.753808),
--        ('iPhone 13', 1400, 'hello', 4, 1, 'Москва', 37.620017, 55.753808),
--        ('iPhone 14', 1500, 'hello', 5, 1, 'Москва', 37.620017, 55.753808),
--        ('iPhone 15', 1600, 'hello', 6, 1, 'Москва', 37.620017, 55.753808),
--        ('iPhone 16', 1700, 'hello', 7, 1, 'Москва', 37.620017, 55.753808),
--        ('iPhone 17', 1800, 'hello', 8, 1, 'Москва', 37.620017, 55.753808),
--        ('iPhone 18', 1900, 'hello', 8, 1, 'Москва', 37.620017, 55.753808),
--        ('iPhone 19', 2100, 'hello', 1, 1, 'Москва', 37.620017, 55.753808),
--        ('iPhone 20', 2400, 'hello', 2, 1, 'Москва', 37.620017, 55.753808);
--
-- INSERT INTO product_images (product_id, img_link)
-- VALUES (1, '/static/product/2e5659cd-72ac-43d8-8494-52bbc7a885fd.webp'),
--        (2, '/static/product/3af8506c-0608-498e-b2b7-7f7a445aa6df.webp'),
--        (3, '/static/product/4abcea38-00ad-4365-85af-1c144085ebd2.webp'),
--        (4, '/static/product/6d835ba7-1ecc-478d-8832-a64b3c58124c.webp'),
--        (5, '/static/product/8b644046-55b7-40a2-beab-308a964630ab.jpg'),
--        (6, '/static/product/697bade2-a4cb-49fc-bad3-c2205554b92a.jpeg'),
--        (7, '/static/product/936de281-1bdb-46e5-a404-3bf2a3fdbaac.webp'),
--        (8, '/static/product/ba1b1a47-97d3-4efb-aed2-f574fc28970f.webp'),
--        (9, '/static/product/dfc9f3d6-60cd-480f-97d1-c31c52dca48b.webp'),
--        (10, '/static/product/f75694ab-42a6-42f7-8f24-cd933cd4da2e.webp'),
--        (11, '/static/product/8776ad39-e754-4f29-8d19-640b1543fbfe.jpg');