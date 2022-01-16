CREATE TABLE IF NOT EXISTS items
(
    id integer,
    prices real[],
    sizes text[] COLLATE pg_catalog."default",
    colors text[] COLLATE pg_catalog."default",
    imagelinks text[] COLLATE pg_catalog."default",
    saleprices real[],
    infodate text[] COLLATE pg_catalog."default",
    count integer,
    category text COLLATE pg_catalog."default" DEFAULT NULL,
    items_id_key integer
);

CREATE TABLE IF NOT EXISTS categories
(
  name text,
  link text
);
