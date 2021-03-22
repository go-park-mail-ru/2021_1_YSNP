create database ysnpkoyaDB
	with owner postgres
	encoding 'utf8'
	LC_COLLATE = 'ru_RU.UTF-8'
    LC_CTYPE = 'ru_RU.UTF-8'
    TABLESPACE = pg_default
    TEMPLATE template0
	;
GRANT ALL PRIVILEGES ON database ysnpkoyaDB TO postgres;