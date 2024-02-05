drop type if exists payment_type_enum ;
drop type if exists status_enum;
drop type if exists transaction_type_enum;
drop type if exists source_type_enum;
drop type if exists tariff_type_enum;
drop type if exists staff_type_enum;
drop type if exists repository_transaction_type_enum;

drop table if exists repository_transactions;
drop table if exists transactions;
drop table if exists staffs;
drop table if exists staff_tariffs;
drop table if exists baskets;
drop table if exists sales;
drop table if exists repositories;
drop table if exists branches;
drop table if exists products;
drop table if exists categories;