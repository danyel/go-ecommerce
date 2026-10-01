-- +goose Up
SET search_path TO ecommerce, public;

INSERT INTO categories (id, name, slug) VALUES
                                            ('00000000-0000-0000-0000-000000000001', 'Computers', 'computers'),
                                            ('00000000-0000-0000-0000-000000000002', 'Laptops', 'laptops'),
                                            ('00000000-0000-0000-0000-000000000003', 'Accessories', 'accessories')
    ON CONFLICT (id) DO NOTHING;
UPDATE categories SET parent_id = '00000000-0000-0000-0000-000000000001'
WHERE id IN ('00000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000003');

INSERT INTO products (id, category_id, brand, name, description, code, price, stock, image_url, metadata) VALUES
                                                                                                              ('10000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000002', 'Northstar', 'Northstar 14', 'A dependable 14 inch laptop for everyday work.', 'DEMO-LAPTOP-14', 899.00, 25, 'https://images.example.invalid/northstar-14.jpg', '{"featured":true}'),
                                                                                                              ('10000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000002', 'Northstar', 'Northstar Pro 16', 'A high performance 16 inch laptop.', 'DEMO-LAPTOP-16', 1499.00, 12, 'https://images.example.invalid/northstar-pro-16.jpg', '{"featured":true}'),
                                                                                                              ('10000000-0000-0000-0000-000000000003', '00000000-0000-0000-0000-000000000003', 'Northstar', 'USB-C Dock', 'A compact dock for modern workspaces.', 'DEMO-DOCK-01', 129.00, 40, 'https://images.example.invalid/usb-c-dock.jpg', '{"featured":false}')
    ON CONFLICT (id) DO NOTHING;

INSERT INTO cms (code, value, language) VALUES
                                            ('demo.catalog.title', 'Catalog', 'en'),
                                            ('demo.catalog.add_to_basket', 'Add to basket', 'en'),
                                            ('demo.basket.title', 'Your basket', 'en')
    ON CONFLICT (code, language) DO UPDATE SET value = EXCLUDED.value;

-- +goose Down
DELETE FROM cms WHERE code IN ('demo.catalog.title', 'demo.catalog.add_to_basket', 'demo.basket.title') AND language = 'en';
DELETE FROM products WHERE id IN ('10000000-0000-0000-0000-000000000001', '10000000-0000-0000-0000-000000000002', '10000000-0000-0000-0000-000000000003');
DELETE FROM categories WHERE id IN ('00000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000003', '00000000-0000-0000-0000-000000000001');
