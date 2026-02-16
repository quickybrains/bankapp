create table if not exists transactions (
  id uuid primary key,
  from_account_id uuid not null references accounts(id),
  to_account_id uuid not null references accounts(id),
  amount double precision not null check (amount > 0),
  descr varchar(100),
  created_at timestamp not null,
  processed_at timestamp not null
);

create index if not exists from_account_id_idx on transactions using hash (from_account_id);
create index if not exists to_account_id_idx on transactions using hash (to_account_id);