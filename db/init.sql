CREATE TABLE IF NOT EXISTS instruments (
    id SERIAL PRIMARY KEY,
    name VARCHAR(64),
    description VARCHAR(200),
    price INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL
);

-- CREATE TABLE IF NOT EXISTS instruments (
--     id TEXT PRIMARY KEY DEFAULT gen_random_uuid(), -- O VARCHAR(255)
--     name VARCHAR(255) NOT NULL,
--     description TEXT,
--     price INT NOT NULL,
--     created_at TIMESTAMP DEFAULT NOW(),
--     updated_at TIMESTAMP NOT NULL
-- );

-- INSERT INTO instruments (id, name, description, price, created_at, updated_at)
-- VALUES 
-- ('a1b2c3d4-e5f6-7890-abcd-1234567890ab', 'Guitarra eléctrica', 'Guitarra Fender Stratocaster de seis cuerdas', 1200, NOW(), NOW()),
-- ('b2c3d4e5-f6a1-8901-bcde-2345678901bc', 'Batería acústica', 'Set completo de batería Pearl con platillos', 2300, NOW(), NOW()),
-- ('c3d4e5f6-a1b2-9012-cdef-3456789012cd', 'Teclado digital', 'Yamaha con 88 teclas contrapesadas', 850, NOW(), NOW()),
-- ('d4e5f6a1-b2c3-0123-def0-4567890123de', 'Violín', 'Violín acústico hecho a mano con arco y estuche', 600, NOW(), NOW()),
-- ('e5f6a1b2-c3d4-1234-ef01-5678901234ef', 'Saxofón alto', 'Saxofón profesional con boquilla y correa', 1500, NOW(), NOW());


INSERT INTO instruments (name, description, price, created_at, updated_at)
VALUES 
('Guitarra eléctrica', 'Guitarra Fender Stratocaster de seis cuerdas', 1200, NOW(), NOW()),
('Batería acústica', 'Set completo de batería Pearl con platillos', 2300, NOW(), NOW()),
('Teclado digital', 'Yamaha con 88 teclas contrapesadas', 850, NOW(), NOW()),
('Violín', 'Violín acústico hecho a mano con arco y estuche', 600, NOW(), NOW()),
('Saxofón alto', 'Saxofón profesional con boquilla y correa', 1500, NOW(), NOW());
