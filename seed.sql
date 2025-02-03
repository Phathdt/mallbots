-- Seed Categories
INSERT INTO categories (name, created_at, updated_at) VALUES
    ('Smartphones', NOW(), NOW()),
    ('Laptops', NOW(), NOW()),
    ('Tablets', NOW(), NOW()),
    ('Accessories', NOW(), NOW()),
    ('Smart Home', NOW(), NOW()),
    ('Audio', NOW(), NOW()),
    ('Gaming', NOW(), NOW()),
    ('Cameras', NOW(), NOW());

-- Seed Products
INSERT INTO products (name, description, price, category_id, created_at, updated_at) VALUES
    -- Smartphones
    ('iPhone 15 Pro', 'Latest Apple iPhone with A17 Pro chip', 999.99, 1, NOW(), NOW()),
    ('Samsung Galaxy S24', 'Flagship Android phone with AI features', 899.99, 1, NOW(), NOW()),
    ('Google Pixel 8', 'Pure Android experience with amazing camera', 799.99, 1, NOW(), NOW()),

    -- Laptops
    ('MacBook Pro 14"', 'Professional laptop with M3 Pro chip', 1999.99, 2, NOW(), NOW()),
    ('Dell XPS 15', 'Premium Windows laptop with OLED display', 1799.99, 2, NOW(), NOW()),
    ('Lenovo ThinkPad X1', 'Business laptop with great keyboard', 1599.99, 2, NOW(), NOW()),

    -- Tablets
    ('iPad Pro 12.9"', 'Powerful tablet for professionals', 1099.99, 3, NOW(), NOW()),
    ('Samsung Galaxy Tab S9', 'Premium Android tablet with S Pen', 849.99, 3, NOW(), NOW()),

    -- Accessories
    ('AirPods Pro', 'Wireless earbuds with noise cancellation', 249.99, 4, NOW(), NOW()),
    ('Apple Watch Series 9', 'Advanced health and fitness tracking', 399.99, 4, NOW(), NOW()),
    ('Samsung Galaxy Watch 6', 'Elegant smartwatch with health features', 299.99, 4, NOW(), NOW()),

    -- Smart Home
    ('Amazon Echo Show', 'Smart display with Alexa', 129.99, 5, NOW(), NOW()),
    ('Google Nest Hub', 'Smart home controller with display', 99.99, 5, NOW(), NOW()),
    ('Philips Hue Starter Kit', 'Smart lighting system', 199.99, 5, NOW(), NOW()),

    -- Audio
    ('Sony WH-1000XM5', 'Premium noise-cancelling headphones', 399.99, 6, NOW(), NOW()),
    ('Bose QuietComfort', 'Comfortable noise-cancelling headphones', 379.99, 6, NOW(), NOW()),
    ('JBL Flip 6', 'Portable Bluetooth speaker', 129.99, 6, NOW(), NOW()),

    -- Gaming
    ('PS5', 'Next-gen gaming console', 499.99, 7, NOW(), NOW()),
    ('Xbox Series X', 'Powerful gaming console', 499.99, 7, NOW(), NOW()),
    ('Nintendo Switch OLED', 'Hybrid gaming console', 349.99, 7, NOW(), NOW()),

    -- Cameras
    ('Sony A7 IV', 'Full-frame mirrorless camera', 2499.99, 8, NOW(), NOW()),
    ('Canon EOS R6', 'Professional mirrorless camera', 2299.99, 8, NOW(), NOW()),
    ('DJI Air 3', 'Premium consumer drone with 4K camera', 999.99, 8, NOW(), NOW());

-- Example of how to verify the seed
-- SELECT c.name as category, COUNT(p.id) as product_count
-- FROM categories c
-- LEFT JOIN products p ON c.id = p.category_id
-- GROUP BY c.name
-- ORDER BY c.name;
