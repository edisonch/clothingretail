insert into clothing_category values (1,'maxim mode','my first maxim mode',1,datetime('now'),datetime('now'));
insert into clothing_category_sub values (1,1,'maxim mode 3 yrs old','JKT HQ', null, null, null, null, null, 1, datetime('now'),datetime('now'));
INSERT INTO clothing_size (id, id_clothing_category_sub, clothes_size_name, clothes_size_notes, clothes_size_status, created_at, updated_at)
VALUES (1, 1, 'XS', 'Extra Small', 1, datetime('now'), datetime('now')),
    (2, 1, 'S', 'Small', 1, datetime('now'), datetime('now')),
    (3, 1, 'M', 'Medium', 1, datetime('now'), datetime('now')),
    (4, 1, 'L', 'Large', 1, datetime('now'), datetime('now')),
    (5, 1, 'XL', 'Extra Large', 1, datetime('now'), datetime('now')),
    (6, 1, 'XXL', 'Double Extra Large', 1, datetime('now'), datetime('now'));

INSERT INTO clothing_inventory_movement values (1,1,1,1,10,0,10,1,datetime('now'),datetime('now'));
INSERT INTO clothing_inventory_movement values (2,1,2,1,18,0,18,1,datetime('now'),datetime('now'));
INSERT INTO clothing_inventory_movement values (3,1,3,1,20,0,18,1,datetime('now'),datetime('now'));
INSERT INTO clothing_inventory_movement values (4,1,4,1,5,0,5,1,datetime('now'),datetime('now'));
INSERT INTO clothing_inventory_movement values (5,1,5,1,8,0,8,1,datetime('now'),datetime('now'));
INSERT INTO clothing_inventory_movement values (6,1,6,1,11,0,11,1,datetime('now'),datetime('now'));

