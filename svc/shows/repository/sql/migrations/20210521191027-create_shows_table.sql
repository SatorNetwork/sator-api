-- +migrate Up
-- +migrate StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE
OR REPLACE FUNCTION shows_update_updated_at_column() RETURNS TRIGGER AS $$
BEGIN NEW .updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';
-- +migrate StatementEnd
CREATE TABLE IF NOT EXISTS shows (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR NOT NULL,
    cover VARCHAR NOT NULL,
    has_new_episode BOOLEAN NOT NULL,
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);
INSERT INTO shows (id, title, cover, has_new_episode, updated_at, created_at) VALUES
('0899a736-e40d-4ba1-b17e-94857e0c8ff0', 'Friends', 'https://hips.hearstapps.com/hmg-prod.s3.amazonaws.com/images/friends-tv-show-1542126105.jpg?crop=1.00xw:0.899xh;0,0&resize=1200:*', 'f', '2021-05-26 11:21:02.626882', '2021-05-26 11:19:22.78509'),
('0e72abac-ecef-47cc-87bd-07c09d8266b0', 'Friends', 'https://hips.hearstapps.com/hmg-prod.s3.amazonaws.com/images/friends-tv-show-1542126105.jpg?crop=1.00xw:0.899xh;0,0&resize=1200:*', 'f', '2021-05-26 11:21:02.626882', '2021-05-26 11:19:22.78509'),
('129ce6cd-bc3d-4823-a8d5-99b6f060b20b', 'Friends', 'https://hips.hearstapps.com/hmg-prod.s3.amazonaws.com/images/friends-tv-show-1542126105.jpg?crop=1.00xw:0.899xh;0,0&resize=1200:*', 'f', '2021-05-26 11:21:02.626882', '2021-04-25 19:59:14.129461'),
('1729013b-d0e4-4e60-838e-f75ee4973481', 'Peaky Blinders', 'https://public-media.si-cdn.com/filer/7b/ba/7bba298e-7e2e-44f0-adb9-b47dfdc1e240/p05m69vt.jpg', 'f', '2021-05-26 11:21:02.626882', '2021-05-26 11:19:22.78509'),
('17cca2de-cd6e-4bdc-affc-a98461cfdba7', 'Silicon Valley', 'https://miro.medium.com/max/3840/1*F-F_pytxhOLs64a18Fq6cw.jpeg', 't', '2021-05-26 11:21:02.626882', '2021-05-26 11:19:22.78509'),
('1c0cc678-6652-43c6-8a8e-3506bf16e2cb', 'Silicon Valley', 'https://miro.medium.com/max/3840/1*F-F_pytxhOLs64a18Fq6cw.jpeg', 't', '2021-05-26 11:21:02.626882', '2021-03-25 20:01:16.288009'),
('242da7cf-5beb-465c-a5ef-8dea1be86aaa', 'Friends', 'https://hips.hearstapps.com/hmg-prod.s3.amazonaws.com/images/friends-tv-show-1542126105.jpg?crop=1.00xw:0.899xh;0,0&resize=1200:*', 'f', '2021-05-26 11:21:02.626882', '2021-05-26 11:19:22.78509'),
('32e5b495-85b3-414c-86ae-1a54012424b2', 'Friends', 'https://hips.hearstapps.com/hmg-prod.s3.amazonaws.com/images/friends-tv-show-1542126105.jpg?crop=1.00xw:0.899xh;0,0&resize=1200:*', 'f', '2021-05-26 11:21:02.626882', '2021-05-26 11:19:22.78509'),
('35ca113b-297d-4b82-8df1-20e5aa8bbc23', 'Friends', 'https://hips.hearstapps.com/hmg-prod.s3.amazonaws.com/images/friends-tv-show-1542126105.jpg?crop=1.00xw:0.899xh;0,0&resize=1200:*', 'f', '2021-05-26 11:21:02.626882', '2021-05-26 11:19:22.78509'),
('374bc0bf-7b53-4942-9d5b-7065b55d6990', 'Friends', 'https://hips.hearstapps.com/hmg-prod.s3.amazonaws.com/images/friends-tv-show-1542126105.jpg?crop=1.00xw:0.899xh;0,0&resize=1200:*', 'f', '2021-05-26 11:21:02.626882', '2021-05-26 11:19:22.78509'),
('3e0ac50f-bda6-4815-9c43-2992c05576f6', 'Peaky Blinders', 'https://public-media.si-cdn.com/filer/7b/ba/7bba298e-7e2e-44f0-adb9-b47dfdc1e240/p05m69vt.jpg', 'f', '2021-05-26 11:21:02.626882', '2021-05-26 11:19:22.78509'),
('457daea5-9c93-44a1-813b-8d2ce01d4ecc', 'Friends', 'https://hips.hearstapps.com/hmg-prod.s3.amazonaws.com/images/friends-tv-show-1542126105.jpg?crop=1.00xw:0.899xh;0,0&resize=1200:*', 'f', '2021-05-26 11:21:02.626882', '2021-05-26 11:19:22.78509'),
('46a116bc-93c2-4003-ad2c-87c1b2493fd2', 'Friends', 'https://hips.hearstapps.com/hmg-prod.s3.amazonaws.com/images/friends-tv-show-1542126105.jpg?crop=1.00xw:0.899xh;0,0&resize=1200:*', 'f', '2021-05-26 11:21:02.626882', '2021-05-26 11:19:22.78509'),
('475dc35b-05ed-467b-945b-39500a87e86f', 'Silicon Valley', 'https://miro.medium.com/max/3840/1*F-F_pytxhOLs64a18Fq6cw.jpeg', 'f', '2021-05-26 11:21:02.626882', '2021-05-26 11:19:22.78509'),
('4f30f97b-08d0-4d5f-a8ff-fa6bb8fd9f14', 'Peaky Blinders', 'https://public-media.si-cdn.com/filer/7b/ba/7bba298e-7e2e-44f0-adb9-b47dfdc1e240/p05m69vt.jpg', 'f', '2021-05-26 11:21:02.626882', '2021-05-26 11:19:22.78509'),
('4fd74c51-75e1-4887-9e79-b51703b00164', 'Silicon Valley', 'https://miro.medium.com/max/3840/1*F-F_pytxhOLs64a18Fq6cw.jpeg', 'f', '2021-05-26 11:21:02.626882', '2021-05-26 11:19:22.78509'),
('6387f386-7e1a-48b1-90ba-73105f078fa3', 'Peaky Blinders', 'https://public-media.si-cdn.com/filer/7b/ba/7bba298e-7e2e-44f0-adb9-b47dfdc1e240/p05m69vt.jpg', 'f', '2021-05-26 11:21:02.626882', '2021-05-26 11:19:22.78509'),
('64d98bd6-44f6-4844-a3c8-9eb2545e0092', 'Friends', 'https://hips.hearstapps.com/hmg-prod.s3.amazonaws.com/images/friends-tv-show-1542126105.jpg?crop=1.00xw:0.899xh;0,0&resize=1200:*', 'f', '2021-05-26 11:21:02.626882', '2021-05-26 11:19:22.78509'),
('7aa40a54-d0b3-4d62-9ce5-090ae254016c', 'Peaky Blinders', 'https://public-media.si-cdn.com/filer/7b/ba/7bba298e-7e2e-44f0-adb9-b47dfdc1e240/p05m69vt.jpg', 'f', '2021-05-26 11:21:02.626882', '2021-05-26 11:19:22.78509'),
('80714ea2-496b-4de2-ad35-8b534194b931', 'Friends', 'https://hips.hearstapps.com/hmg-prod.s3.amazonaws.com/images/friends-tv-show-1542126105.jpg?crop=1.00xw:0.899xh;0,0&resize=1200:*', 'f', '2021-05-26 11:21:02.626882', '2021-05-26 11:19:22.78509'),
('81694a55-c4f7-42c4-a4db-2c526c103e55', 'Peaky Blinders', 'https://public-media.si-cdn.com/filer/7b/ba/7bba298e-7e2e-44f0-adb9-b47dfdc1e240/p05m69vt.jpg', 'f', '2021-05-26 11:21:02.626882', '2021-05-26 11:19:22.78509'),
('87fcb73b-f7c9-4a81-8cd7-6e0ce4a8cb4b', 'Silicon Valley', 'https://miro.medium.com/max/3840/1*F-F_pytxhOLs64a18Fq6cw.jpeg', 'f', '2021-05-26 11:21:02.626882', '2021-05-26 11:19:22.78509'),
('9388ff34-bc3d-4bc2-9a90-af05a4668111', 'Friends', 'https://hips.hearstapps.com/hmg-prod.s3.amazonaws.com/images/friends-tv-show-1542126105.jpg?crop=1.00xw:0.899xh;0,0&resize=1200:*', 'f', '2021-05-26 11:21:02.626882', '2021-05-26 11:19:22.78509'),
('93b21999-d81a-4568-8723-79c3ad5342c2', 'Peaky Blinders', 'https://public-media.si-cdn.com/filer/7b/ba/7bba298e-7e2e-44f0-adb9-b47dfdc1e240/p05m69vt.jpg', 'f', '2021-05-26 11:21:02.626882', '2021-05-26 11:19:22.78509'),
('93fae417-2ea3-46c8-ad9e-74f21bd60e7e', 'Peaky Blinders', 'https://public-media.si-cdn.com/filer/7b/ba/7bba298e-7e2e-44f0-adb9-b47dfdc1e240/p05m69vt.jpg', 'f', '2021-05-26 11:21:02.626882', '2021-05-26 11:19:22.78509'),
('9bbdea43-04b9-4a60-a6fa-7930c2190883', 'Friends', 'https://hips.hearstapps.com/hmg-prod.s3.amazonaws.com/images/friends-tv-show-1542126105.jpg?crop=1.00xw:0.899xh;0,0&resize=1200:*', 'f', '2021-05-26 11:21:02.626882', '2021-05-25 19:59:14.129461'),
('9d4d2a79-5761-4577-9f34-ad03a8e9c082', 'Peaky Blinders', 'https://public-media.si-cdn.com/filer/7b/ba/7bba298e-7e2e-44f0-adb9-b47dfdc1e240/p05m69vt.jpg', 'f', '2021-05-26 11:21:02.626882', '2021-05-26 11:19:22.78509'),
('a0879308-9555-4986-b69e-2e7479e06e98', 'Friends', 'https://hips.hearstapps.com/hmg-prod.s3.amazonaws.com/images/friends-tv-show-1542126105.jpg?crop=1.00xw:0.899xh;0,0&resize=1200:*', 'f', '2021-05-26 11:21:02.626882', '2021-05-26 11:19:22.78509'),
('a3c96a3c-3b30-40ad-9b4b-809b94e561ff', 'Silicon Valley', 'https://miro.medium.com/max/3840/1*F-F_pytxhOLs64a18Fq6cw.jpeg', 't', '2021-05-26 11:21:02.626882', '2021-05-26 11:19:22.78509'),
('a55608bb-47f7-4d0d-96b6-3bd5989edb71', 'Peaky Blinders', 'https://public-media.si-cdn.com/filer/7b/ba/7bba298e-7e2e-44f0-adb9-b47dfdc1e240/p05m69vt.jpg', 'f', '2021-05-26 11:21:02.626882', '2021-05-26 11:19:22.78509'),
('a7cb36d5-7479-4dcc-9fa6-1e7ec6e1b5c6', 'Silicon Valley', 'https://miro.medium.com/max/3840/1*F-F_pytxhOLs64a18Fq6cw.jpeg', 't', '2021-05-26 11:21:02.626882', '2021-05-26 11:19:22.78509'),
('a84a9c45-c8e7-40ba-a627-b4998b8ae8ce', 'Peaky Blinders', 'https://public-media.si-cdn.com/filer/7b/ba/7bba298e-7e2e-44f0-adb9-b47dfdc1e240/p05m69vt.jpg', 'f', '2021-05-26 11:21:02.626882', '2021-05-26 11:19:22.78509'),
('a8bddfdf-518f-429a-886f-b9ea33d47100', 'Silicon Valley', 'https://miro.medium.com/max/3840/1*F-F_pytxhOLs64a18Fq6cw.jpeg', 'f', '2021-05-26 11:21:02.626882', '2021-05-26 11:19:22.78509'),
('b042c280-cf07-40cc-bcb6-7489ee2bbdb6', 'Silicon Valley', 'https://miro.medium.com/max/3840/1*F-F_pytxhOLs64a18Fq6cw.jpeg', 't', '2021-05-26 11:21:02.626882', '2021-05-26 11:19:22.78509'),
('b28f21c5-bbd0-4fea-9c2a-613f42fde8c9', 'Friends', 'https://hips.hearstapps.com/hmg-prod.s3.amazonaws.com/images/friends-tv-show-1542126105.jpg?crop=1.00xw:0.899xh;0,0&resize=1200:*', 'f', '2021-05-26 11:21:02.626882', '2021-05-26 11:19:22.78509'),
('c0aa5955-23fa-447f-b735-5fff1ae5351f', 'Peaky Blinders', 'https://public-media.si-cdn.com/filer/7b/ba/7bba298e-7e2e-44f0-adb9-b47dfdc1e240/p05m69vt.jpg', 'f', '2021-05-26 11:21:02.626882', '2021-05-26 11:19:22.78509'),
('c82b93d2-6743-46d2-91f5-601bbee36451', 'Peaky Blinders', 'https://public-media.si-cdn.com/filer/7b/ba/7bba298e-7e2e-44f0-adb9-b47dfdc1e240/p05m69vt.jpg', 'f', '2021-05-26 11:21:02.626882', '2021-05-26 11:19:22.78509'),
('ccc6af01-5d04-4b57-8415-4afee168f2e9', 'Peaky Blinders', 'https://public-media.si-cdn.com/filer/7b/ba/7bba298e-7e2e-44f0-adb9-b47dfdc1e240/p05m69vt.jpg', 't', '2021-05-26 11:21:02.626882', '2021-05-21 20:03:29.284069'),
('d4b7550a-0ed7-49ef-95e2-14e23348e2a0', 'Silicon Valley', 'https://miro.medium.com/max/3840/1*F-F_pytxhOLs64a18Fq6cw.jpeg', 't', '2021-05-26 11:21:02.626882', '2021-05-26 11:19:22.78509'),
('d662040b-9241-481d-b3df-a2e4a65d1697', 'Silicon Valley', 'https://miro.medium.com/max/3840/1*F-F_pytxhOLs64a18Fq6cw.jpeg', 't', '2021-05-26 11:21:02.626882', '2021-05-26 11:19:22.78509'),
('e331073e-48d2-47ee-8b67-0ec3978a3861', 'Silicon Valley', 'https://miro.medium.com/max/3840/1*F-F_pytxhOLs64a18Fq6cw.jpeg', 'f', '2021-05-26 11:21:02.626882', '2021-05-26 11:19:22.78509'),
('e3979844-45fc-4d0a-9813-5204002eeb9a', 'Silicon Valley', 'https://miro.medium.com/max/3840/1*F-F_pytxhOLs64a18Fq6cw.jpeg', 'f', '2021-05-26 11:21:02.626882', '2021-05-26 11:19:22.78509'),
('e6fb2362-0b41-467f-b341-2df64a8445c4', 'Peaky Blinders', 'https://public-media.si-cdn.com/filer/7b/ba/7bba298e-7e2e-44f0-adb9-b47dfdc1e240/p05m69vt.jpg', 'f', '2021-05-26 11:21:02.626882', '2021-05-26 11:19:22.78509'),
('e7764cc1-7e0d-4320-b34c-c2c6c49ea90c', 'Friends', 'https://hips.hearstapps.com/hmg-prod.s3.amazonaws.com/images/friends-tv-show-1542126105.jpg?crop=1.00xw:0.899xh;0,0&resize=1200:*', 'f', '2021-05-26 11:21:02.626882', '2021-05-26 11:19:22.78509'),
('e8568ca6-c933-4436-a6c9-d3098ec49f45', 'Silicon Valley', 'https://miro.medium.com/max/3840/1*F-F_pytxhOLs64a18Fq6cw.jpeg', 't', '2021-05-26 11:21:02.626882', '2021-05-24 20:01:16.288009'),
('eca09dc8-d889-4754-adbe-f21fd17463cf', 'Peaky Blinders', 'https://public-media.si-cdn.com/filer/7b/ba/7bba298e-7e2e-44f0-adb9-b47dfdc1e240/p05m69vt.jpg', 't', '2021-05-26 11:21:02.626882', '2021-02-25 20:03:29.284069'),
('f5de5f5c-5250-4c81-b9a3-a4eb4c6a0058', 'Silicon Valley', 'https://miro.medium.com/max/3840/1*F-F_pytxhOLs64a18Fq6cw.jpeg', 'f', '2021-05-26 11:21:02.626882', '2021-05-26 11:19:22.78509'),
('fdc53e1d-f15a-4028-a216-27dd05578768', 'Silicon Valley', 'https://miro.medium.com/max/3840/1*F-F_pytxhOLs64a18Fq6cw.jpeg', 't', '2021-05-26 11:21:02.626882', '2021-05-26 11:19:22.78509');
CREATE INDEX ordering_shows_list ON shows USING BTREE (updated_at, created_at);
CREATE TRIGGER update_shows_modtime BEFORE
UPDATE ON shows FOR EACH ROW EXECUTE PROCEDURE shows_update_updated_at_column();
-- +migrate Down
DROP TRIGGER IF EXISTS update_shows_modtime ON shows;
DROP TABLE IF EXISTS shows;
DROP FUNCTION IF EXISTS shows_update_updated_at_column();