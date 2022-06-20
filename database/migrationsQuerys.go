package database

var Usuarios string = `CREATE TABLE IF NOT EXISTS usuarios ( 
	id serial, 
	username varchar(255) not null,
	lastname varchar(255) not null,
	password varchar(255) not null,
	active boolean,
	primary key (id)
	);`

var index_users string = `CREATE UNIQUE INDEX username_table_users ON usuarios (username);`	

var Friends string = `CREATE TABLE IF NOT EXISTS friends ( 
	id serial, 
	primary key (id),
	IDuser1 int not null,
	IDuser2 int not null,
	foreign key (IDuser1) references usuarios(id),
	foreign key (IDuser2) references usuarios(id)
	ON DELETE CASCADE
	);`
