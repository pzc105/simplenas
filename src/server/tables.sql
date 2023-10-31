create database pnas;

use pnas;

create table user (
  id bigint not null auto_increment,
  name varchar(256) not null,
  email varchar(256) not null,
  passwd varchar(256) not null,
  auth varbinary(64) not null,
  directory_id bigint not null,
  primary key(id),
  unique key email (email)
);

/*
drop table video;
drop table category_items;
drop table category_type;
*/

create table video (
  id bigint not null auto_increment,
  file_name varchar(512) not null,
  hls_created tinyint default 0 not null,
  created_at datetime default current_timestamp not null,
  primary key(id),
  key file (file_name)
);

drop procedure if exists new_video;
delimiter //
create procedure new_video(in abs_file_name varchar(512), out new_video_id bigint)
begin
  insert into video (file_name) values(abs_file_name);
  select last_insert_id() into new_video_id;
  select new_video_id;
end//
delimiter ;


create table category_type (
  id int not null,
  type_name varchar(32) not null,
  primary key(id)
);
insert into category_type(id, type_name) values 
(1, "home"),
(2, "directory"),
(3, "video"),
(4, "other")
;

create table category_items (
  id bigint not null auto_increment,
  parent_id bigint not null,
  type_id int not null,
  name varchar(256) not null,
  creator bigint not null,
  auth varbinary(64) not null,
  resource_path varchar(256) not null,
  poster_path varchar(256) not null,
  introduce text not null,
  other text not null,
  created_at datetime default current_timestamp not null,
  updated_at timestamp default current_timestamp on update current_timestamp not null,
  primary key(id),
  key name (name),
  key resource(creator, type_id, resource_path),
  key parent(parent_id, name),
  fulltext fulltext_ni(name, introduce) with parser ngram
);

insert into category_items (id, parent_id, type_id, name, creator, auth, resource_path, poster_path, introduce) values 
  (1, 0, 1, "root", 1, "", "", "", "");
insert into category_items (id, parent_id, type_id, name, creator, auth, resource_path, poster_path, introduce) values 
  (2, 1, 2, "tmp", 1, "", "", "", "");
insert into category_items (id, parent_id, type_id, name, creator, auth, resource_path, poster_path, introduce) values 
  (3, 1, 1, "users", 1, "", "", "", "");
insert into category_items (id, parent_id, type_id, name, creator, auth, resource_path, poster_path, introduce) values 
  (4, 1, 2, "magnet-shares", 1, "", "", "", "");

insert into user(id, name, email, passwd, auth, directory_id) values
  (1, "admin", "admin@admin.cn", "202cb962ac59075b964b07152d234b70", "", 1);

drop procedure if exists new_category;
delimiter //
create procedure new_category(in type_id int,
                              in name varchar(256),
                              in creator bigint,
                              in auth varbinary(64),
                              in resource_path varchar(256),
                              in poster_path varchar(256),
                              in introduce text,
                              in other text,
                              in parent_id bigint,
                              out new_item_id bigint)
begin
  declare parent_count int default 0;
  start transaction;
  select count(*) into parent_count from pnas.category_items where id = parent_id for update;
  if parent_count = 1 then
    insert into pnas.category_items (parent_id, type_id, name, creator, auth, resource_path, poster_path, introduce, other) values 
      (parent_id, type_id, name, creator, auth, resource_path, poster_path, introduce, other);
    select last_insert_id() into new_item_id;
  else
    set new_item_id = -2;
  end if;
  select new_item_id;
  commit;
end//
delimiter ;

drop procedure if exists del_category;
delimiter //
create procedure del_category(in del_item_id bigint)
begin
  start transaction;
  delete from pnas.category_items where id = del_item_id;
  commit;
end//
delimiter ;

drop procedure if exists new_user;
delimiter //
create procedure new_user(in name varchar(256),
                          in email varchar(256),
                          in passwd varchar(256),
                          in auth varbinary(64),
                          in homeAuth varbinary(64),
                          out new_user_id bigint,
                          out new_home_id bigint)
begin
  start transaction;
  insert into pnas.user (name, email, passwd, auth, directory_id) values(name, email, passwd, auth, 0);
  select last_insert_id() into new_user_id;
  insert into pnas.category_items (parent_id, type_id, name, creator, auth, resource_path, poster_path, introduce) values 
      (3, 1, name, new_user_id, homeAuth, "", "", "");
  select last_insert_id() into new_home_id;
  update pnas.user set directory_id=new_home_id where id=new_user_id;
  select new_user_id, new_home_id;
  commit;
end//
delimiter ;


/*
drop table torrent;
drop table user_torrent;

*/

create table torrent (
  id bigint not null auto_increment,
  name varchar(128) not null,
  version int not null,
  info_hash varbinary(64) not null,
  state int default 0 not null comment 'unknown = 0,checking_files = 1, downloading_metadata = 2, downloading = 3, finished = 4, seeding = 5, checking_resume_data = 7',
  total_size bigint default 0  not null,
  piece_length int default 0  not null,
  num_pieces int default 0  not null,
  introduce text not null,
  resume_data longblob not null,
  created_at datetime default current_timestamp not null,
  updated_at timestamp default current_timestamp on update current_timestamp not null,

  primary key(id),
  key info_hash (info_hash, version),
  fulltext fulltext_i(introduce) with parser ngram
);

create table user_torrent (
  user_id bigint not null,
  torrent_id bigint not null,
  created_at datetime default current_timestamp not null,

  primary key(user_id, torrent_id)
);

drop procedure if exists new_torrent;
delimiter //
create procedure new_torrent(in version int,
                             in info_hash varbinary(64),
                             in user_id bigint, out torrent_id bigint)
begin
  declare ut_count int default 0;
  start transaction;
  select count(*) into ut_count from pnas.user_torrent u left join torrent t on t.id = u.torrent_id
    where u.user_id = user_id and t.info_hash = info_hash;
  if ut_count = 0 then
    insert into pnas.torrent (name, version, info_hash, introduce, resume_data) values ('', version, info_hash, '', '');
    select last_insert_id() into torrent_id;
    insert into pnas.user_torrent (user_id, torrent_id) values (user_id, torrent_id);
    select torrent_id;
  end if;
  commit;
end//
delimiter ;


create table magnet (
  id bigint not null auto_increment,
  version int not null,
  info_hash varbinary(64) not null,
  magnet_uri text not null,
  created_at datetime default current_timestamp not null,
  updated_at timestamp default current_timestamp on update current_timestamp not null,

  primary key(id),
  key info_hash (info_hash, version)
);