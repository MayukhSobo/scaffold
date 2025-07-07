-- +goose Up
-- +goose StatementBegin
INSERT INTO users (
    username, email, password_hash, first_name, last_name, avatar_url, bio, phone_number,
    address_street, address_city, address_state, address_postal_code, address_country,
    status, role, email_verified_at
) VALUES
(
    'johndoe', 'johndoe@example.com', '$2a$10$fS/C.DD.aQ7.SQjFwGvBzeN.9w.E5.e.N8.E2.L3.d.A5.g.Y5.y6',
    'John', 'Doe', 'https://example.com/avatars/johndoe.png', 'A software engineer from California.', '+1-202-555-0104',
    '123 Main St', 'Anytown', 'CA', '12345', 'USA',
    'active', 'user', NOW()
),
(
    'janedoe', 'janedoe@example.com', '$2a$10$fS/C.DD.aQ7.SQjFwGvBzeN.9w.E5.e.N8.E2.L3.d.A5.g.Y5.y6',
    'Jane', 'Doe', 'https://example.com/avatars/janedoe.png', 'A product manager from New York.', '+1-202-555-0168',
    '456 Oak Ave', 'Somecity', 'NY', '54321', 'USA',
    'active', 'admin', NOW()
),
(
    'peterjones', 'peterjones@example.com', '$2a$10$fS/C.DD.aQ7.SQjFwGvBzeN.9w.E5.e.N8.E2.L3.d.A5.g.Y5.y6',
    'Peter', 'Jones', 'https://example.com/avatars/peterjones.png', 'A data scientist from Texas.', '+1-202-555-0182',
    '789 Pine Ln', 'Anotherplace', 'TX', '67890', 'USA',
    'active', 'user', NOW()
),
(
    'marywhite', 'marywhite@example.com', '$2a$10$fS/C.DD.aQ7.SQjFwGvBzeN.9w.E5.e.N8.E2.L3.d.A5.g.Y5.y6',
    'Mary', 'White', 'https://example.com/avatars/marywhite.png', 'A UX designer from Washington.', '+1-202-555-0155',
    '101 Maple Dr', 'Newville', 'WA', '13579', 'USA',
    'suspended', 'user', NOW()
),
(
    'brianblack', 'brianblack@example.com', '$2a$10$fS/C.DD.aQ7.SQjFwGvBzeN.9w.E5.e.N8.E2.L3.d.A5.g.Y5.y6',
    'Brian', 'Black', 'https://example.com/avatars/brianblack.png', 'A DevOps engineer from Florida.', '+1-202-555-0199',
    '212 Birch Rd', 'Oldtown', 'FL', '24680', 'USA',
    'pending_verification', 'user', NULL
),
(
    'susanbrown', 'susanbrown@example.com', '$2a$10$fS/C.DD.aQ7.SQjFwGvBzeN.9w.E5.e.N8.E2.L3.d.A5.g.Y5.y6',
    'Susan', 'Brown', 'https://example.com/avatars/susanbrown.png', 'A project manager from Illinois.', '+1-202-555-0121',
    '333 Cedar Ct', 'Yourtown', 'IL', '97531', 'USA',
    'active', 'moderator', NOW()
),
(
    'kevingreen', 'kevingreen@example.com', '$2a$10$fS/C.DD.aQ7.SQjFwGvBzeN.9w.E5.e.N8.E2.L3.d.A5.g.Y5.y6',
    'Kevin', 'Green', 'https://example.com/avatars/kevingreen.png', 'A frontend developer from Arizona.', '+1-202-555-0134',
    '444 Spruce Pkwy', 'Hometown', 'AZ', '86420', 'USA',
    'inactive', 'user', NOW()
),
(
    'lisagray', 'lisagray@example.com', '$2a$10$fS/C.DD.aQ7.SQjFwGvBzeN.9w.E5.e.N8.E2.L3.d.A5.g.Y5.y6',
    'Lisa', 'Gray', 'https://example.com/avatars/lisagray.png', 'A full-stack developer from Colorado.', '+1-202-555-0143',
    '555 Willow Way', 'Lakeside', 'CO', '19283', 'USA',
    'active', 'user', NOW()
),
(
    'michaelyellow', 'michaelyellow@example.com', '$2a$10$fS/C.DD.aQ7.SQjFwGvBzeN.9w.E5.e.N8.E2.L3.d.A5.g.Y5.y6',
    'Michael', 'Yellow', 'https://example.com/avatars/michaelyellow.png', 'A backend developer from Oregon.', '+1-202-555-0177',
    '666 Redwood Ave', 'Mountainview', 'OR', '47568', 'USA',
    'pending_verification', 'user', NULL
),
(
    'karenpurple', 'karenpurple@example.com', '$2a$10$fS/C.DD.aQ7.SQjFwGvBzeN.9w.E5.e.N8.E2.L3.d.A5.g.Y5.y6',
    'Karen', 'Purple', 'https://example.com/avatars/karenpurple.png', 'A QA engineer from Nevada.', '+1-202-555-0110',
    '777 Sequoia St', 'Riverdale', 'NV', '30495', 'USA',
    'active', 'moderator', NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE users;
-- +goose StatementEnd
