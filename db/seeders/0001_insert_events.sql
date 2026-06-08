-- +goose Up
INSERT INTO public.events
(id, "name", slug, "location", "date", description, is_private, salt, created_at, deleted_at)
VALUES(1, 'Trah Projowinoto 2509', 'trah-projo-2509', 'Yogyakarta', '2025-09-15', 'Arisan Trah Keluarga Besar Projowinoto Jogja-Weleri-Semarang September 2025', true, '464f49d0e2f8c12359748ef94a775075', '2026-05-20 03:55:56.527', NULL);
INSERT INTO public.events
(id, "name", slug, "location", "date", description, is_private, salt, created_at, deleted_at)
VALUES(2, 'Desktop Setup August 2024', 'setup-aug-24', 'Yogyakarta', '2025-09-15', 'Desktop Setup for OchaXT Setup Reaction Video August 2024', false, '18b5b993ce45ee4bae818fe345692ad3', '2026-05-20 05:23:16.799', NULL);


-- +goose Down
DELETE FROM public.events
