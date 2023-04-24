# PostgreSQL

- var use `""`
- string use `''`

## list all tables

```sql
SELECT "table_name", "table_type" 
FROM "information_schema"."tables" 
WHERE "table_schema" = 'public' 
ORDER BY "table_name" ASC;
```

## list all columns

```sql
SELECT "column_name", "data_type", "character_maximum_length"
FROM "information_schema"."columns" 
WHERE "table_schema" = 'public'
ORDER BY "column_name" ASC;
```
