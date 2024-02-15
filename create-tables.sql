-- CREATE USER 'authServiceGo'@'localhost' IDENTIFIED BY 'password';
-- GRANT ALL PRIVILEGES ON *.* TO 'authServiceGo'@'localhost' WITH GRANT OPTION;
-- FLUSH PRIVILEGES;
--
-- mysql -u authServiceGo -p users < ~/projects/auth-service-go/create-tables.sql
DROP TABLE IF EXISTS user;
CREATE TABLE user (
  email     VARCHAR(255) NOT NULL UNIQUE,
  password     VARCHAR(255) NOT NULL,
    id  int NOT NULL AUTO_INCREMENT UNIQUE,
  session     VARCHAR(255) NOT NULL,
  PRIMARY KEY (`email`)
);

-- INSERT INTO user
--   (email, password, session)
-- VALUES
--   ('testA@mail.com', 'testA', NULL),
--   ('testB@mail.com', 'testB', 'testB'),
--   ('testC@mail.com', 'testC', NULL);
