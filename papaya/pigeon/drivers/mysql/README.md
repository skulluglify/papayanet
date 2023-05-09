# MySQL

- var use ``` `` ```
- string use `""`

# list all users
## MySQL

```sql
SELECT `user`, `host`, `plugin` FROM `mysql`.`user`;
CREATE USER `user`@`localhost` IDENTIFIED WITH 'caching_sha2_password' BY 'User@1234' WITH MAX_USER_CONNECTIONS 2;
CREATE DATABASE IF NOT EXISTS `main` CHARACTER SET 'utf8mb4';
GRANT ALL ON `main`.* TO `user`@`localhost`;
FLUSH PRIVILEGES;
```

## mysql cli

```sh
mysql -h localhost -P 3306 -u user -p -D main
```