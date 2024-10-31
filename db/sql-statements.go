package db

const (
	CREATE_TABLE = `
  create table services (id integer primary key, username text, password text, service text);
  delete from services;`

	INSERT_SERVICE = `
  insert into services(id, username, password, service) values(NULL,?,?,?)
  `

	DELETE_SERVICE = `
  delete from services where id = ?
  `

	SELECT_ALL_SERVICES = `
  select * from services
  `

	SELECT_SERVICE_BY_NAME = `
  select username, password, service from services where service = ?
  `
)
