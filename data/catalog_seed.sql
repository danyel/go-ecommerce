-- Optional catalog seed for databases that already ran the Goose baseline.
-- The same deterministic rows are installed by migrations/20260919000200_seed_catalog.sql.
SET search_path TO ecommerce, public;

INSERT INTO categories (id, name, slug)
VALUES
    ('00000000-0000-0000-0000-000000000001', 'Computers', 'computers'),
    ('00000000-0000-0000-0000-000000000002', 'Laptops', 'laptops'),
    ('00000000-0000-0000-0000-000000000003', 'Accessories', 'accessories')
ON CONFLICT (id) DO NOTHING;

UPDATE categories
SET parent_id = '00000000-0000-0000-0000-000000000001'
WHERE id IN (
    '00000000-0000-0000-0000-000000000002',
    '00000000-0000-0000-0000-000000000003'
);

INSERT INTO products (id, category_id, brand, name, description, code, price, stock, image_url, metadata)
VALUES
    ('10000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000002', 'Northstar', 'Northstar 14', 'A dependable 14 inch laptop for everyday work.', 'DEMO-LAPTOP-14', 899.00, 25, 'https://images.example.invalid/northstar-14.jpg', '{"featured":true}'),
    ('10000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000002', 'Northstar', 'Northstar Pro 16', 'A high performance 16 inch laptop.', 'DEMO-LAPTOP-16', 1499.00, 12, 'https://images.example.invalid/northstar-pro-16.jpg', '{"featured":true}'),
    ('10000000-0000-0000-0000-000000000003', '00000000-0000-0000-0000-000000000003', 'Northstar', 'USB-C Dock', 'A compact dock for modern workspaces.', 'DEMO-DOCK-01', 129.00, 40, 'https://images.example.invalid/usb-c-dock.jpg', '{"featured":false}')
ON CONFLICT (id) DO UPDATE SET stock = EXCLUDED.stock, price = EXCLUDED.price;
