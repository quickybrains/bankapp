create table if not exists transactions (
  id uuid primary key,
  from_account_id uuid not null,
  to_account_id uuid not null,
  amount double precision not null check (amount > 0),
  descr varchar(100)
);

create index if not exists from_account_id_idx on transactions using hash (from_account_id);
create index if not exists to_account_id_idx on transactions using hash (to_account_id);