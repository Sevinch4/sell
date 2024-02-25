create table if not exists incomes(
    id uuid primary key ,
    branch_id uuid references branches(id),
    price int,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP DEFAULT NULL
);

create table if not exists income_products(
  id uuid primary key ,
  income_id uuid references incomes(id),
  product_id uuid references products(id),
  price int,
  count int,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  deleted_at TIMESTAMP DEFAULT NULL
);