CREATE TABLE IF NOT EXISTS instruments (
    id SERIAL PRIMARY KEY,
    name VARCHAR(64),
    description VARCHAR(200),
    price INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL
);

INSERT INTO instruments (name, description, price, created_at, updated_at)
VALUES 
('Guitarra eléctrica', 'Guitarra Fender Stratocaster de seis cuerdas', 1200, NOW(), NOW()),
('Batería acústica', 'Set completo de batería Pearl con platillos', 2300, NOW(), NOW()),
('Teclado digital', 'Yamaha con 88 teclas contrapesadas', 850, NOW(), NOW()),
('Violín', 'Violín acústico hecho a mano con arco y estuche', 600, NOW(), NOW()),
('Saxofón alto', 'Saxofón profesional con boquilla y correa', 1500, NOW(), NOW());
