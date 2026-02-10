create table if not exists accounts (
  id uuid primary key,
  user_id uuid not null,
  fullname varchar(100),
  currency varchar(10),
  balance double precision not null
);

create index if not exists accounts_user_id_idx on accounts using hash (user_id);