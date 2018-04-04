--Table To store client info
create table client(id serial primary key,client_no varchar(20),name varchar(256));

--Table to store ad campaign info
create table ad_campaign(id serial primary key,client_id int,name varchar(256),targeting_type_pos boolean,foreign key(client_id) references client(id));

--Table to store State info
create table state(id serial primary key,state_code varchar(20),name varchar(128),country_name varchar(128));

--Separate table to store postitive and negative ad campaign.Note: It can be done using one table but then will have to add a new boolean to distinguish between the two
create table ad_pos(ad_campaign_id int,state_ids int[],foreign key (ad_campaign_id) references ad_campaign(id));
create table ad_neg(ad_campaign_id int,state_ids int[],foreign key (ad_campaign_id) references ad_campaign(id));

--Table to store the images
create table images(id serial primary key,ad_campaign_id int,image_location varchar(512),image bytea,foreign key (ad_campaign_id) references ad_campaign(id));


--Dump
insert into client(client_no,name)values('001','Idea');
insert into client(client_no,name)values('002','Uninor');

insert into ad_campaign(client_id,name,targeting_type_pos)values(1,'Goa,Maharashtra & Rajasthan',True);
insert into ad_campaign(client_id,name,targeting_type_pos)values(2,'Not Jammu & Kashmir and Himachal Pradesh',False);

insert into state(state_code,name,country_name)values('MH','Maharashtra','India');
insert into state(state_code,name,country_name)values('RJ','Rajasthan','India');
insert into state(state_code,name,country_name)values('GJ','Gujarat','India');
insert into state(state_code,name,country_name)values('Hp','Himachal Pradesh','India');
insert into state(state_code,name,country_name)values('Jk','Jammu & Kashmir','India');
insert into state(state_code,name,country_name)values('GA','Goa','India');
insert into state(state_code,name,country_name)values('Ka','Karnataka','India');

insert into ad_pos(ad_campaign_id,state_ids)values(1,array[1,2,6]);
insert into ad_neg(ad_campaign_id,state_ids)values(2,array[4,5]);

insert into images(ad_campaign_id,image_location)values(1,'http://telecomtalk.info/wp-content/uploads/2016/03/idea-my-app-4G.png');
insert into images(ad_campaign_id,image_location)values(2,'http://telecomtalk.info/wp-content/uploads/2015/09/uninor-india.png');