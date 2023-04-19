CREATE TABLE posts (
	id serial4 NOT NULL,
	author_id varchar(64) NOT NULL,
	"content" text NOT NULL,
	"type" varchar(64) NOT NULL,
	"data" text NOT NULL,
	visibility varchar(32) NOT NULL DEFAULT 'public',
	PRIMARY KEY (id),
	FOREIGN KEY (author_id) REFERENCES profiles(id)
);