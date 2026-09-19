-- Original demonstration catalog. Product names, descriptions, and images are
-- created for this project and are not copied from a third-party retailer.
INSERT INTO ecommerce.categories (id, name)
VALUES
    ('10000000-0000-0000-0000-000000000001', 'Computers'),
    ('10000000-0000-0000-0000-000000000002', 'Components'),
    ('10000000-0000-0000-0000-000000000003', 'Graphics Cards'),
    ('10000000-0000-0000-0000-000000000004', 'Processors'),
    ('10000000-0000-0000-0000-000000000005', 'Motherboards'),
    ('10000000-0000-0000-0000-000000000006', 'Memory'),
    ('10000000-0000-0000-0000-000000000007', 'Storage'),
    ('10000000-0000-0000-0000-000000000008', 'Power Supplies'),
    ('10000000-0000-0000-0000-000000000009', 'Cooling'),
    ('10000000-0000-0000-0000-000000000010', 'Cases'),
    ('10000000-0000-0000-0000-000000000011', 'Peripherals'),
    ('10000000-0000-0000-0000-000000000012', 'Monitors'),
    ('10000000-0000-0000-0000-000000000013', 'Keyboards'),
    ('10000000-0000-0000-0000-000000000014', 'Mice'),
    ('10000000-0000-0000-0000-000000000015', 'Audio'),
    ('10000000-0000-0000-0000-000000000016', 'Networking'),
    ('10000000-0000-0000-0000-000000000017', 'Cables and Accessories'),
    ('10000000-0000-0000-0000-000000000018', 'Gaming')
ON CONFLICT (id) DO NOTHING;

INSERT INTO ecommerce.category_children (parent_id, child_id)
VALUES
    ('10000000-0000-0000-0000-000000000001', '10000000-0000-0000-0000-000000000002'),
    ('10000000-0000-0000-0000-000000000001', '10000000-0000-0000-0000-000000000011'),
    ('10000000-0000-0000-0000-000000000001', '10000000-0000-0000-0000-000000000016'),
    ('10000000-0000-0000-0000-000000000001', '10000000-0000-0000-0000-000000000018'),
    ('10000000-0000-0000-0000-000000000002', '10000000-0000-0000-0000-000000000003'),
    ('10000000-0000-0000-0000-000000000002', '10000000-0000-0000-0000-000000000004'),
    ('10000000-0000-0000-0000-000000000002', '10000000-0000-0000-0000-000000000005'),
    ('10000000-0000-0000-0000-000000000002', '10000000-0000-0000-0000-000000000006'),
    ('10000000-0000-0000-0000-000000000002', '10000000-0000-0000-0000-000000000007'),
    ('10000000-0000-0000-0000-000000000002', '10000000-0000-0000-0000-000000000008'),
    ('10000000-0000-0000-0000-000000000002', '10000000-0000-0000-0000-000000000009'),
    ('10000000-0000-0000-0000-000000000002', '10000000-0000-0000-0000-000000000010'),
    ('10000000-0000-0000-0000-000000000011', '10000000-0000-0000-0000-000000000012'),
    ('10000000-0000-0000-0000-000000000011', '10000000-0000-0000-0000-000000000013'),
    ('10000000-0000-0000-0000-000000000011', '10000000-0000-0000-0000-000000000014'),
    ('10000000-0000-0000-0000-000000000011', '10000000-0000-0000-0000-000000000015'),
    ('10000000-0000-0000-0000-000000000011', '10000000-0000-0000-0000-000000000017')
ON CONFLICT (parent_id, child_id) DO NOTHING;

INSERT INTO ecommerce.products
    (id, name, description, brand, code, price, stock, image_url, category_id)
VALUES
    ('20000000-0000-0000-0000-000000000001', 'Asterion RX 7800 XT', 'High-performance graphics card for smooth 1440p gaming with 16 GB of fast memory.', 'Asterion', 'AST-GPU-7800XT', 529.99, 12, 'https://placehold.co/600x400/png?text=Graphics+Card', '10000000-0000-0000-0000-000000000003'),
    ('20000000-0000-0000-0000-000000000002', 'Northstar RTX 4060 Twin Fan', 'Efficient graphics card for compact gaming systems and creative workloads.', 'Northstar', 'NOR-GPU-4060TF', 319.99, 18, 'https://placehold.co/600x400/png?text=Graphics+Card', '10000000-0000-0000-0000-000000000003'),
    ('20000000-0000-0000-0000-000000000003', 'Cobalt Core 7 8700X', 'Eight-core desktop processor designed for gaming, development, and everyday productivity.', 'Cobalt', 'COB-CPU-8700X', 349.99, 14, 'https://placehold.co/600x400/png?text=Processor', '10000000-0000-0000-0000-000000000004'),
    ('20000000-0000-0000-0000-000000000004', 'Cobalt Core 5 7600', 'Balanced six-core processor with excellent performance for a mid-range desktop build.', 'Cobalt', 'COB-CPU-7600', 199.99, 22, 'https://placehold.co/600x400/png?text=Processor', '10000000-0000-0000-0000-000000000004'),
    ('20000000-0000-0000-0000-000000000005', 'Forge AM5 Creator Board', 'ATX motherboard with modern connectivity, high-speed storage support, and four memory slots.', 'Forge', 'FOR-MB-AM5-C', 189.99, 9, 'https://placehold.co/600x400/png?text=Motherboard', '10000000-0000-0000-0000-000000000005'),
    ('20000000-0000-0000-0000-000000000006', 'Forge B760 Workbench Board', 'Reliable Intel-compatible motherboard with Wi-Fi, PCIe expansion, and USB-C connectivity.', 'Forge', 'FOR-MB-B760-W', 169.99, 11, 'https://placehold.co/600x400/png?text=Motherboard', '10000000-0000-0000-0000-000000000005'),
    ('20000000-0000-0000-0000-000000000007', 'Orbit DDR5 32GB Kit', 'Two-module 32 GB DDR5 kit tuned for responsive multitasking and modern desktop systems.', 'Orbit', 'ORB-MEM-DDR5-32', 94.99, 35, 'https://placehold.co/600x400/png?text=Memory', '10000000-0000-0000-0000-000000000006'),
    ('20000000-0000-0000-0000-000000000008', 'Orbit DDR4 16GB Kit', 'Affordable dual-channel memory kit for dependable home and office computers.', 'Orbit', 'ORB-MEM-DDR4-16', 42.99, 40, 'https://placehold.co/600x400/png?text=Memory', '10000000-0000-0000-0000-000000000006'),
    ('20000000-0000-0000-0000-000000000009', 'Vector NVMe 1TB', 'Fast PCIe solid-state drive for operating systems, applications, and game libraries.', 'Vector', 'VEC-SSD-NVME-1T', 69.99, 31, 'https://placehold.co/600x400/png?text=NVMe+SSD', '10000000-0000-0000-0000-000000000007'),
    ('20000000-0000-0000-0000-000000000010', 'Vector NVMe 2TB Pro', 'High-capacity NVMe drive with sustained performance for creative projects and large libraries.', 'Vector', 'VEC-SSD-NVME-2T', 139.99, 24, 'https://placehold.co/600x400/png?text=NVMe+SSD', '10000000-0000-0000-0000-000000000007'),
    ('20000000-0000-0000-0000-000000000011', 'Harbor 2TB Desktop Drive', 'Quiet mechanical storage for backups, media archives, and secondary capacity.', 'Harbor', 'HAR-HDD-2T', 59.99, 27, 'https://placehold.co/600x400/png?text=Hard+Drive', '10000000-0000-0000-0000-000000000007'),
    ('20000000-0000-0000-0000-000000000012', 'Anchor 650W Bronze PSU', 'Efficient 650 watt power supply with modular cabling for mainstream gaming builds.', 'Anchor', 'ANC-PSU-650B', 74.99, 16, 'https://placehold.co/600x400/png?text=Power+Supply', '10000000-0000-0000-0000-000000000008'),
    ('20000000-0000-0000-0000-000000000013', 'Anchor 850W Gold PSU', 'Fully modular 850 watt power supply for high-performance graphics cards and quiet operation.', 'Anchor', 'ANC-PSU-850G', 129.99, 10, 'https://placehold.co/600x400/png?text=Power+Supply', '10000000-0000-0000-0000-000000000008'),
    ('20000000-0000-0000-0000-000000000014', 'Breeze Tower Cooler', 'Compact tower air cooler with a quiet fan and broad socket compatibility.', 'Breeze', 'BRZ-CPU-COOL-T', 34.99, 21, 'https://placehold.co/600x400/png?text=CPU+Cooler', '10000000-0000-0000-0000-000000000009'),
    ('20000000-0000-0000-0000-000000000015', 'Breeze 240 Liquid Cooler', 'All-in-one liquid cooler for sustained processor workloads and clean system builds.', 'Breeze', 'BRZ-AIO-240', 89.99, 13, 'https://placehold.co/600x400/png?text=Liquid+Cooler', '10000000-0000-0000-0000-000000000009'),
    ('20000000-0000-0000-0000-000000000016', 'Frame ATX Airflow Case', 'Roomy ATX case with a mesh front, three included fans, and tool-free drive bays.', 'Frame', 'FRM-CASE-ATX-A', 84.99, 17, 'https://placehold.co/600x400/png?text=PC+Case', '10000000-0000-0000-0000-000000000010'),
    ('20000000-0000-0000-0000-000000000017', 'Frame Mini Mesh Case', 'Small-form-factor case for compact systems with full-size graphics card support.', 'Frame', 'FRM-CASE-MINI', 69.99, 8, 'https://placehold.co/600x400/png?text=PC+Case', '10000000-0000-0000-0000-000000000010'),
    ('20000000-0000-0000-0000-000000000018', 'Vista 27 QHD Monitor', '27-inch QHD display with a fast refresh rate and adjustable stand for productive gaming.', 'Vista', 'VIS-MON-27QHD', 249.99, 15, 'https://placehold.co/600x400/png?text=Monitor', '10000000-0000-0000-0000-000000000012'),
    ('20000000-0000-0000-0000-000000000019', 'Vista 24 Office Monitor', 'Comfortable 24-inch Full HD monitor with low-blue-light and height-adjustable viewing modes.', 'Vista', 'VIS-MON-24FHD', 119.99, 26, 'https://placehold.co/600x400/png?text=Monitor', '10000000-0000-0000-0000-000000000012'),
    ('20000000-0000-0000-0000-000000000020', 'Keyline Mechanical 75', 'Compact mechanical keyboard with hot-swappable switches and programmable layers.', 'Keyline', 'KEY-MECH-75', 89.99, 20, 'https://placehold.co/600x400/png?text=Keyboard', '10000000-0000-0000-0000-000000000013'),
    ('20000000-0000-0000-0000-000000000021', 'Keyline Quiet Office', 'Low-profile wireless keyboard with quiet keys and a rechargeable battery.', 'Keyline', 'KEY-OFFICE-W', 39.99, 33, 'https://placehold.co/600x400/png?text=Keyboard', '10000000-0000-0000-0000-000000000013'),
    ('20000000-0000-0000-0000-000000000022', 'Glide Wireless Mouse', 'Ergonomic wireless mouse with adjustable sensitivity and silent primary buttons.', 'Glide', 'GLD-MOUSE-W', 29.99, 42, 'https://placehold.co/600x400/png?text=Mouse', '10000000-0000-0000-0000-000000000014'),
    ('20000000-0000-0000-0000-000000000023', 'Glide Precision Mouse', 'Lightweight wired mouse with a flexible cable and configurable side buttons.', 'Glide', 'GLD-MOUSE-P', 44.99, 28, 'https://placehold.co/600x400/png?text=Mouse', '10000000-0000-0000-0000-000000000014'),
    ('20000000-0000-0000-0000-000000000024', 'Echo USB Headset', 'Clear voice headset with a flexible microphone for calls, classes, and casual gaming.', 'Echo', 'ECH-HEAD-USB', 49.99, 23, 'https://placehold.co/600x400/png?text=Headset', '10000000-0000-0000-0000-000000000015'),
    ('20000000-0000-0000-0000-000000000025', 'Echo Desktop Speakers', 'Compact stereo speakers powered over USB with a dedicated volume control.', 'Echo', 'ECH-SPK-DESK', 34.99, 19, 'https://placehold.co/600x400/png?text=Speakers', '10000000-0000-0000-0000-000000000015'),
    ('20000000-0000-0000-0000-000000000026', 'Beacon Wi-Fi 6 Router', 'Dual-band router with mesh support, guest networking, and parental controls.', 'Beacon', 'BCN-ROUTER-6', 99.99, 14, 'https://placehold.co/600x400/png?text=Wi-Fi+Router', '10000000-0000-0000-0000-000000000016'),
    ('20000000-0000-0000-0000-000000000027', 'Beacon 8-Port Switch', 'Quiet unmanaged gigabit switch for home offices, studios, and small labs.', 'Beacon', 'BCN-SWITCH-8', 24.99, 37, 'https://placehold.co/600x400/png?text=Network+Switch', '10000000-0000-0000-0000-000000000016'),
    ('20000000-0000-0000-0000-000000000028', 'LinkPro USB-C Dock', 'Multi-display USB-C dock with network, card reader, and charging connectivity.', 'LinkPro', 'LKP-DOCK-USBC', 109.99, 12, 'https://placehold.co/600x400/png?text=USB-C+Dock', '10000000-0000-0000-0000-000000000017'),
    ('20000000-0000-0000-0000-000000000029', 'LinkPro Display Cable', 'Durable high-bandwidth DisplayPort cable for high-resolution monitors.', 'LinkPro', 'LKP-CABLE-DP', 14.99, 55, 'https://placehold.co/600x400/png?text=Display+Cable', '10000000-0000-0000-0000-000000000017'),
    ('20000000-0000-0000-0000-000000000030', 'Forge Play Desktop', 'Balanced gaming desktop combining a modern processor, dedicated graphics, and fast NVMe storage.', 'Forge', 'FOR-PC-PLAY-01', 1249.99, 6, 'https://placehold.co/600x400/png?text=Gaming+Desktop', '10000000-0000-0000-0000-000000000018'),
    ('20000000-0000-0000-0000-000000000031', 'Forge Creator Desktop', 'Quiet workstation desktop with expanded memory and storage for editing and development.', 'Forge', 'FOR-PC-CREATE-01', 1799.99, 4, 'https://placehold.co/600x400/png?text=Creator+Desktop', '10000000-0000-0000-0000-000000000018')
ON CONFLICT (id) DO NOTHING;
