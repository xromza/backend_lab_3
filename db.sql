CREATE TABLE applications (
  application_id int unsigned NOT NULL AUTO_INCREMENT,
  surname varchar(128) NOT NULL DEFAULT '',
  name varchar(128) NOT NULL DEFAULT '',
  midname varchar(128) NOT NULL DEFAULT '',
  phone_number varchar(128) NOT NULL DEFAULT '',
  email varchar(128) NOT NULL DEFAULT '',
  gender tinyint unsigned NOT NULL DEFAULT 0,
  bio text NOT NULL,
  birth_date DATE NOT NULL,
  PRIMARY KEY (application_id)
) Engine=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE languages (
language_id int unsigned NOT NULL AUTO_INCREMENT,
name varchar(64) NOT NULL,
PRIMARY KEY (language_id)
) Engine=InnoDB;

CREATE TABLE application_language (
application_id int unsigned NOT NULL,
language_id int unsigned NOT NULL,
PRIMARY KEY (application_id, language_id),
FOREIGN KEY (application_id) REFERENCES applications(application_id) ON DELETE CASCADE,
FOREIGN KEY (language_id) REFERENCES languages(language_id) ON DELETE CASCADE
) Engine=InnoDB;

INSERT INTO languages (name) VALUES 
('Pascal'), ('C'),('C++'), 
('JavaScript'),('PHP'), ('Python'),
('Java'), ('Haskell'),('Clojure'), 
('Prolog'),('Scala'), ('Go');
