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
    avatar    varchar(512) NOT NULL DEFAULT '',
    score     int                   DEFAULT 0,
    reviews   int                   DEFAULT 0
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


CREATE OR REPLACE FUNCTION check_achievement_sold() RETURNS TRIGGER AS
$check_achievement_sold$
DECLARE
    S_COUNT INTEGER;
BEGIN
    SELECT COUNT(*) from product where owner_id = NEW.owner_id and product.close = true INTO S_COUNT;
    IF (S_COUNT = 1) THEN
        INSERT INTO user_achievement (user_id, a_id) VALUES (NEW.owner_id, 4);
    end if;

    IF (S_COUNT = 10) THEN
        INSERT INTO user_achievement (user_id, a_id) VALUES (NEW.owner_id, 5);
    end if;

    IF (S_COUNT = 100) THEN
        INSERT INTO user_achievement (user_id, a_id) VALUES (NEW.owner_id, 6);
    end if;
    
    RETURN NEW; -- возвращаемое значение для триггера AFTER игнорируется
END;
$check_achievement_sold$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION check_achievement_fav() RETURNS TRIGGER AS
$check_achievement_fav$
DECLARE
    F_COUNT INTEGER;
    Ach_Count INTEGER;
BEGIN

    SELECT COUNT(*) from user_favorite where user_id = NEW.user_id INTO F_COUNT;
    Select count(*) from user_achievement where user_achievement.user_id = New.user_id and user_achievement.a_id = 7 INTO Ach_Count;
    IF (F_COUNT = 10 and Ach_Count = 0) THEN
        INSERT INTO user_achievement (user_id, a_id) VALUES (NEW.user_id, 7);
    end if;

    RETURN NEW; -- возвращаемое значение для триггера AFTER игнорируется
END;
$check_achievement_fav$ LANGUAGE plpgsql;


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


CREATE OR REPLACE FUNCTION check_achievement_review() RETURNS TRIGGER AS
$check_achievement_review$
DECLARE
    R_COUNT INTEGER;
BEGIN

    SELECT COUNT(*) from user_reviews where reviewer_id = NEW.reviewer_id INTO R_COUNT;
    IF (R_COUNT = 1) THEN
        INSERT INTO user_achievement (user_id, a_id) VALUES (NEW.reviewer_id, 8);
    end if;

    IF (R_COUNT = 10) THEN
        INSERT INTO user_achievement (user_id, a_id) VALUES (NEW.reviewer_id, 9);
    end if;

    SELECT COUNT(rating) from user_reviews where target_id = NEW.target_id INTO R_COUNT;
    IF (R_COUNT = 10) THEN  
        INSERT INTO user_achievement (user_id, a_id) VALUES (NEW.target_id, 10);
    end if;
    RETURN NEW; 
END;
$check_achievement_review$ LANGUAGE plpgsql;


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
VALUES ('Новичок', 'Первое объявление', '/img/svg/ach1.svg'),
       ('Опытный','Десять объявлений', '/img/svg/ach3.svg'),
       ('Шарящий','Сто объявлений', '/img/svg/ach3.svg'),
       ('Малый бизнесмен', 'Первая продажа', '/img/svg/ach4.svg'),
       ('Бизнесмен', 'Десять продаж', '/img/svg/ach5.svg'),
       ('Директор по продажам', 'Сто продаж', '/img/svg/ach6.svg'),
       ('Любимка', 'Добавил 10 вещей в избранное', '/img/svg/ach7.svg'),
       ('Рецензент', 'Первый оставленный отзыв', '/img/svg/ach8.svg'),
       ('Главный рецензент', 'Десять оставленных отзывов', '/img/svg/ach9.svg'),
       ('Честный', '10 полученных 5-звездочных отзывов', '/img/svg/ach10.svg');

create table if not exists user_achievement
(
    user_id        integer   not null,
    date           timestamp not null default NOW(),
    a_id           integer   not null,

    primary key (user_id, a_id),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE NO ACTION,
    FOREIGN KEY (a_id) REFERENCES achievement (id) ON DELETE NO ACTION
);

CREATE TRIGGER check_achievement
    AFTER INSERT
    ON product
    FOR EACH ROW
EXECUTE PROCEDURE check_achievement();

CREATE TRIGGER check_achievement_sold
    AFTER UPDATE of close
    ON product
    FOR EACH ROW
EXECUTE PROCEDURE check_achievement_sold();

CREATE TRIGGER check_achievement_fav
    AFTER INSERT 
    ON user_favorite
    FOR EACH ROW
EXECUTE PROCEDURE check_achievement_fav();

CREATE TRIGGER check_achievement_review
    AFTER INSERT 
    ON user_reviews
    FOR EACH ROW
EXECUTE PROCEDURE check_achievement_review();
