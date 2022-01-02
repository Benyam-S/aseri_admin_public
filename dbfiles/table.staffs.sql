-- drop table staffs;

create table staffs(
    id varchar(50) primary key,
	first_name varchar(50),
	last_name varchar(50),
	phone_number varchar(13),
	email varchar(50) unique,
	profile_pic varchar(255),
	role varchar(10),
	created_at timestamp,
	updated_at timestamp
);

