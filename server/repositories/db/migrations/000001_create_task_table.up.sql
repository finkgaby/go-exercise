create table if not exists dataentries(
    "id" varchar,
    "title" varchar,
    "content" varchar,
    "views" int,
    "timestamp" timestamp default current_timestamp
    );