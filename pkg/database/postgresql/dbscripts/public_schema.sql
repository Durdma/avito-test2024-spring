SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET client_min_messages = warning;
SET row_security = off;
CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA pg_catalog;
SET search_path = public, pg_catalog;
SET default_tablespace = '';

create table features (
    id bigserial not null,
    primary key (id)
);

create table banners (
    id bigserial not null,
    fk_feature_id bigserial,
    title varchar(250) not null,
    text varchar(1000) not null,
    is_active bool not null default true,
    created_at date not null,
    updated_at date not null,
    primary key (id),
    constraint fk_feature
    foreign key (fk_feature_id) references features(id)
        on delete restrict on update restrict
);

create table tags (
    id bigserial NOT NULL,
    primary key (id)
);

create table banners_tags (
    fk_banner_id bigserial not null,
    fk_tag_id bigserial not null,
    primary key (fk_banner_id, fk_tag_id),
    foreign key (fk_banner_id) references banners(id)
        on delete cascade on update restrict,
    foreign key (fk_tag_id) references tags(id)
        on delete cascade on update restrict
);
