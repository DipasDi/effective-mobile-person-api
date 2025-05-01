create table if not exists people (
   id          serial primary key,
   name        varchar(100) not null,
   surname     varchar(100) not null,
   patronymic  varchar(100),
   age         integer,
   gender      varchar(10),
   nationality varchar(100)
);