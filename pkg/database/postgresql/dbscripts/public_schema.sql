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

create table if not exists features (
    id bigserial not null,
    primary key (id)
);

create table if not exists banners (
    id bigserial not null,
    fk_feature_id int,
    content jsonb not null,
    is_active bool not null default true,
    created_at timestamp not null,
    updated_at timestamp not null,
    primary key (id),
    constraint fk_feature
    foreign key (fk_feature_id) references features(id)
        on delete restrict on update restrict
);

create table if not exists tags (
    id bigserial NOT NULL,
    primary key (id)
);

create table if not exists banners_tags (
    fk_banner_id int not null,
    fk_tag_id int not null,
    fk_feature_id int not null,
    primary key (fk_banner_id, fk_tag_id),
    foreign key (fk_banner_id) references banners(id)
        on delete cascade on update restrict,
    foreign key (fk_tag_id) references tags(id)
        on delete cascade on update restrict,
    foreign key (fk_feature_id) references features(id)
        on delete cascade on update restrict,
    constraint unique_banner_tag_feature unique (fk_feature_id, fk_tag_id)
);

create table if not exists users (
    id bigserial not null,
    fk_tag_id int,
    is_admin bool not null default false,
    primary key (id),
    constraint fk_tag
    foreign key (fk_tag_id) references tags(id)
        on delete restrict on update restrict
);