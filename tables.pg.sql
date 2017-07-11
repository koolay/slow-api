create table slow_sql.slow_sql
(
	id int not null auto_increment
		primary key,
	user varchar(200) null,
	host varchar(100) null,
	query_time double null,
	lock_time double null,
	rows_sent int null,
	rows_examined int null,
	`sql` text null,
	`when` timestamp default CURRENT_TIMESTAMP not null,
	created_on datetime default CURRENT_TIMESTAMP null
);
