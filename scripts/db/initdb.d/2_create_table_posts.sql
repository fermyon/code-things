CREATE TABLE posts (
  "id"        SERIAL,
  "author_id" varchar(64) references profiles(id) NOT NULL,
  "content"   text NOT NULL,
  "type"      varchar(64) NOT NULL,
  "data"      text NOT NULL,
  PRIMARY KEY(id)
);