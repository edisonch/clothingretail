-- clothing_category contains all the available category for rental
-- id contains the id for category
-- clothes_cat_name contains the name of the category limit to 32 characters
-- clothes_notes contains the notes for the category limit to 256 characters
-- clothes_cat_status contains the status of the category: 1 = active, 2 = inactive
-- created_at contains the date and time when the category is created
-- updated_at contains the date and time when the category is updated
create table if not exists clothing_category (
    id integer primary key,
    clothes_cat_name text not null,
    clothes_notes text,
    clothes_cat_status integer not null default 1,
    created_at datetime not null,
    updated_at datetime not null
);

-- clothing_category_sub contains all the available subcategory for rental within the category
-- id contains the id for subcategory
-- id_clothing_category contains the id for the category
-- clothes_cat_name_sub contains the name of the subcategory limit to 32 characters
-- clothes_cat_location_sub contains the location of the subcategory limit to 64 characters
-- clothes_picture_1 contains the picture of the subcategory using base64 encoding
-- clothes_picture_2 contains the picture of the subcategory using base64 encoding
-- clothes_picture_3 contains the picture of the subcategory using base64 encoding
-- clothes_picture_4 contains the picture of the subcategory using base64 encoding
-- clothes_picture_5 contains the picture of the subcategory using base64 encoding
-- clothes_cat_status_sub contains the status of the subcategory: 1 = active, 2 = inactive
-- created_at contains the date and time when the subcategory is created
-- updated_at contains the date and time when the subcategory is updated
create table if not exists clothing_category_sub (
    id integer primary key,
    id_clothing_category integer not null REFERENCES clothing_category(id),
    clothes_cat_name_sub text not null,
    clothes_cat_location_sub text not null,
    clothes_picture_1 TEXT, -- base64 encoded image
    clothes_picture_2 TEXT, -- base64 encoded image
    clothes_picture_3 TEXT, -- base64 encoded image
    clothes_picture_4 TEXT, -- base64 encoded image
    clothes_picture_5 TEXT, -- base64 encoded image
    clothes_cat_status_sub integer not null default 1,
    created_at datetime not null,
    updated_at datetime not null
);

-- clothing_size contains all the available size for rental within the subcategory
-- id contains the id for size
-- id_clothing_category_sub contains the id for the category_sub
-- clothes_size_name contains the name of the size limit to 8 characters
-- clothes_size_notes contains the notes for the size limit to 256 characters
-- clothes_size_status contains the status of the size: 1 = active, 2 = inactive
-- created_at contains the date and time when the size is created
-- updated_at contains the date and time when the size is updated
create table if not exists clothing_size (
    id integer primary key,
    id_clothing_category_sub integer not null REFERENCES clothing_category_sub(id),
    clothes_size_name text not null,
    clothes_size_notes text,
    clothes_size_status integer not null default 1,
    created_at datetime not null,
    updated_at datetime not null
);

-- clothing_inventory_movement contains all the available inventory movement for rental within the subcategory
-- id contains the id for inventory movement
-- id_clothing_category_sub contains the id for the category_sub
-- id_clothing_size contains the id for the size
-- clothes_movement_action contains the action of the inventory movement:
-- 1 = BUY, 2 = SELL , 3 = RENT, 4 = RETURN, 5 = CANCEL, 6 = NOT RETURN, 7 = LOSS
-- clothes_qty_in contains the quantity of the inventory movement in
-- clothes_qty_out contains the quantity of the inventory movement out
-- clothes_qty_total contains the total quantity of the inventory movement
-- clothes_cat_status_sub contains the status of the inventory movement: 1 = active, 2 = inactive
-- created_at contains the date and time when the inventory movement is created
-- updated_at contains the date and time when the inventory movement is updated
create table if not exists clothing_inventory_movement (
    id integer primary key,
    id_clothing_category integer not null REFERENCES clothing_category_sub(id),
    id_clothing_size integer not null REFERENCES clothing_size(id),
    clothes_movement_action integer not null default 1,
    clothes_qty_in integer not null,
    clothes_qty_out integer not null,
    clothes_qty_total integer not null,
    clothes_cat_status_sub integer not null default 1,
    created_at datetime not null,
    updated_at datetime not null
);

-- clothing_customer contains all the available customer for rental
-- id contains the id for customer
-- cust_name contains the name of the customer limit to 64 characters
-- cust_address contains the address of the customer limit to 256 characters
-- cust_city contains the city of the customer limit to 64 characters
-- cust_phone contains the phone of the customer limit to 16 characters
-- cust_email contains the email of the customer limit to 128 characters
-- cust_notes contains the notes for the customer limit to 256 characters
-- cust_status contains the status of the customer: 1 = active, 2 = inactive
-- created_at contains the date and time when the customer is created
create table if not exists clothing_customer (
    id integer primary key,
    cust_name text not null,
    cust_address text not null,
    cust_city text not null,
    cust_phone text not null,
    cust_email text not null,
    cust_notes text,
    cust_status integer not null default 1,
    created_at datetime not null,
    updated_at datetime not null
);

-- clothing_rental contains all the available rental for customer
-- id contains the id for rental
-- id_clothing_category_sub contains the id for the category_sub
-- id_clothing_size contains the id for the size
-- id_clothing_customer contains the id for the customer
-- clothes_qty_rent contains the quantity of the rental
-- clothes_qty_return contains the quantity of the return
-- clothes_rent_status contains the status of the rental: 1 = rent, 2 = return, 3 = cancel, 4 = not_return, 5 = loss
-- clothes_rent_date_begin contains the date and time when the rental is begin
-- clothes_rent_date_end contains the date and time when the rental is end
-- clothes_rent_date_actual_pickup contains the date and time when the pickup is actual pickup
-- clothes_rent_date_actual_return contains the date and time when the rental is actual return
-- created_at contains the date and time when the rental is created
-- updated_at contains the date and time when the rental is updated
create table if not exists clothing_rental (
    id integer primary key,
    id_clothing_category_sub integer not null REFERENCES clothing_category_sub(id),
    id_clothing_size integer not null REFERENCES clothing_size(id),
    id_clothing_customer integer not null REFERENCES clothing_customer(id),
    clothes_qty_rent integer not null default 1,
    clothes_qty_return integer not null default 0,
    clothes_rent_date_begin datetime not null,
    clothes_rent_date_end datetime not null,
    clothes_rent_date_actual_pickup datetime not null,
    clothes_rent_date_actual_return datetime not null,
    clothes_rent_status integer not null default 1,
    created_at datetime not null,
    updated_at datetime not null
);

-- clothing_users contains all the available users for rental
-- id contains the id for users
-- username contains the username of the users limit to 32 characters
-- pin contains the pin of the users limit to 6 digits
-- user_status contains the status of the users: 1 = active, 2 = inactive, 3 = suspended
-- created_at contains the date and time when the users is created
-- updated_at contains the date and time when the users is updated
create table if not exists clothing_users (
    id integer primary key,
    username text not null,
    pin integer not null default 123456,
    user_status integer not null default 1,
    created_at datetime not null,
    updated_at datetime not null
);