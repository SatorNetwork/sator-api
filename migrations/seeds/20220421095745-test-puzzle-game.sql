INSERT INTO public.shows (id, title, cover, has_new_episode, updated_at, created_at, category, description, realms_title, realms_subtitle, watch, archived) VALUES
('6d3eb699-b725-49d5-ac62-84667b462b88', 'Test', 'test', true, NULL, '2022-04-21 03:55:22.706128', 'test', 'test', NULL, NULL, NULL, false);

INSERT INTO public.episodes (id, show_id, season_id, episode_number, cover, title, description, release_date, updated_at, created_at, challenge_id, verification_challenge_id, hint_text, watch, archived) VALUES
('455acf92-1b11-4861-af6b-7351f57d060d', '6d3eb699-b725-49d5-ac62-84667b462b88', null, 1, 'Test', 'Test', 'Test', null, null, '2022-04-21 04:02:28.990313', null, null, 'Test', 'Test', false);

INSERT INTO public.puzzle_games (id, episode_id, prize_pool, parts_x, parts_y, updated_at, created_at) VALUES
('7801d5d3-2d2c-4f85-9190-3fa82527f2af', '455acf92-1b11-4861-af6b-7351f57d060d', 1, 4, 4, null, '2022-04-21 04:07:13.093272');

INSERT INTO public.files (id, file_name, file_path, file_url, created_at) VALUES
('b7684f9f-91d6-4683-8489-f16612f28aa5', 'test', 'test', 'test', '2022-04-21 04:51:41.063916'),
('3320cf8c-f6d9-42f1-b8d6-cea843910c43', 'test', 'test', 'test', '2022-04-21 04:51:41.063916'),
('e260ec9c-13a7-42b7-a945-40ef2dbc0303', 'test', 'test', 'test', '2022-04-21 04:51:41.063916'),
('9fe63234-90ba-4aa6-9fb2-62ca738f597d', 'test', 'test', 'test', '2022-04-21 04:51:41.063916'),
('ac2acc14-ce99-4c28-9596-5a26d574e965', 'test', 'test', 'test', '2022-04-21 04:51:41.063916'),
('998baee6-74d4-4898-8a94-32a4e867749e', 'test', 'test', 'test', '2022-04-21 04:51:41.063916'),
('49062388-3ea3-4833-b23a-7d3e6081739d', 'test', 'test', 'test', '2022-04-21 04:51:41.063916'),
('a851b737-3c87-4bc4-8d04-39bcc3c3f34a', 'test', 'test', 'test', '2022-04-21 04:51:41.063916'),
('c8195e1d-83d1-4701-8ebc-0b237040101b', 'test', 'test', 'test', '2022-04-21 04:51:41.063916'),
('194ad696-d7d9-4b97-90ab-f52f514027e5', 'test', 'test', 'test', '2022-04-21 04:51:41.063916'),
('e3693fe0-c157-4742-9781-a96bccc35c35', 'test', 'test', 'test', '2022-04-21 04:51:41.063916'),
('681b7f1d-870a-4367-9b60-e3df60955ead', 'test', 'test', 'test', '2022-04-21 04:51:41.063916'),
('29e7a54b-c86c-4021-ad17-3c27d7cd6b9c', 'test', 'test', 'test', '2022-04-21 04:51:41.063916'),
('311cfc0e-f1bb-4c54-baac-bb82e9964c6e', 'test', 'test', 'test', '2022-04-21 04:51:41.063916'),
('08e18925-bc61-40db-90fe-d60f363c5487', 'test', 'test', 'test', '2022-04-21 04:51:41.063916'),
('a706620b-a179-41f8-9217-830327dea1f5', 'test', 'test', 'test', '2022-04-21 04:51:41.063916');

INSERT INTO puzzle_games_to_images(file_id, puzzle_game_id) VALUES
('b7684f9f-91d6-4683-8489-f16612f28aa5', '7801d5d3-2d2c-4f85-9190-3fa82527f2af'),
('3320cf8c-f6d9-42f1-b8d6-cea843910c43', '7801d5d3-2d2c-4f85-9190-3fa82527f2af'),
('e260ec9c-13a7-42b7-a945-40ef2dbc0303', '7801d5d3-2d2c-4f85-9190-3fa82527f2af'),
('9fe63234-90ba-4aa6-9fb2-62ca738f597d', '7801d5d3-2d2c-4f85-9190-3fa82527f2af'),
('ac2acc14-ce99-4c28-9596-5a26d574e965', '7801d5d3-2d2c-4f85-9190-3fa82527f2af'),
('998baee6-74d4-4898-8a94-32a4e867749e', '7801d5d3-2d2c-4f85-9190-3fa82527f2af'),
('49062388-3ea3-4833-b23a-7d3e6081739d', '7801d5d3-2d2c-4f85-9190-3fa82527f2af'),
('a851b737-3c87-4bc4-8d04-39bcc3c3f34a', '7801d5d3-2d2c-4f85-9190-3fa82527f2af'),
('c8195e1d-83d1-4701-8ebc-0b237040101b', '7801d5d3-2d2c-4f85-9190-3fa82527f2af'),
('194ad696-d7d9-4b97-90ab-f52f514027e5', '7801d5d3-2d2c-4f85-9190-3fa82527f2af'),
('e3693fe0-c157-4742-9781-a96bccc35c35', '7801d5d3-2d2c-4f85-9190-3fa82527f2af'),
('681b7f1d-870a-4367-9b60-e3df60955ead', '7801d5d3-2d2c-4f85-9190-3fa82527f2af'),
('29e7a54b-c86c-4021-ad17-3c27d7cd6b9c', '7801d5d3-2d2c-4f85-9190-3fa82527f2af'),
('311cfc0e-f1bb-4c54-baac-bb82e9964c6e', '7801d5d3-2d2c-4f85-9190-3fa82527f2af'),
('08e18925-bc61-40db-90fe-d60f363c5487', '7801d5d3-2d2c-4f85-9190-3fa82527f2af'),
('a706620b-a179-41f8-9217-830327dea1f5', '7801d5d3-2d2c-4f85-9190-3fa82527f2af');

INSERT INTO public.puzzle_game_unlock_options (id, steps, amount, disabled, updated_at, created_at, locked) VALUES
('899fcefd-9b4b-4c67-905f-b1e580fcaf78', 32, 0, false, null, '2022-04-21 05:27:55.628799', false);