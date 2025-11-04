-- Insert clothing sizes for subcategory ID 1
INSERT INTO clothing_size (id, id_clothing_category_sub, clothes_size_name, clothes_size_notes, clothes_size_status, created_at, updated_at)
VALUES
    (1, 1, 'XS', 'Extra Small', 1, datetime('now'), datetime('now')),
    (2, 1, 'S', 'Small', 1, datetime('now'), datetime('now')),
    (3, 1, 'M', 'Medium', 1, datetime('now'), datetime('now')),
    (4, 1, 'L', 'Large', 1, datetime('now'), datetime('now')),
    (5, 1, 'XL', 'Extra Large', 1, datetime('now'), datetime('now')),
    (6, 1, 'XXL', 'Double Extra Large', 1, datetime('now'), datetime('now'));