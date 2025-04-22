--
-- PostgreSQL database dump
--

-- Dumped from database version 14.17 (Ubuntu 14.17-0ubuntu0.22.04.1)
-- Dumped by pg_dump version 14.17 (Ubuntu 14.17-0ubuntu0.22.04.1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Data for Name: person; Type: TABLE DATA; Schema: public; Owner: colet
--

COPY public.person (id, slug, first, last, description, birthdate, profile_img, search_vector, insert_timestamp) FROM stdin;
1	kyle-mooney-1	Kyle	Mooney	\N	1984-09-03	kyle-mooney-1.jpg	'kyle':1 'mooney':2	2025-01-31 14:31:12.72663
2	tim-gilbert-4	Tim	Gilbert	this is the description	1983-05-13	tim-gilbert-4.jpg	'descript':6 'gilbert':2 'tim':1	2025-01-31 14:31:12.72663
3	james-hartnett-5	James	Hartnett	\N	\N	james-hartnett-5.jpg	'hartnett':2 'jame':1	2025-01-31 14:31:12.72663
4	test-alpha-4	Test	Alpha	\N	\N	james-hartnett-5.jpg	'alpha':2 'test':1	2025-01-31 14:31:12.72663
5	test-beta-5	Test	Beta	this is the description	1983-05-13	tim-gilbert-4.jpg	'beta':2 'descript':6 'test':1	2025-01-31 14:31:12.72663
6	test-charlie-6	Test	Charlie	\N	1984-09-03	kyle-mooney-1.jpg	'charli':2 'test':1	2025-01-31 14:31:12.72663
7	test-delta-6	Test	Delta	\N	1984-09-03	kyle-mooney-1.jpg	'delta':2 'test':1	2025-01-31 14:31:12.72663
\.


--
-- Data for Name: character; Type: TABLE DATA; Schema: public; Owner: colet
--

COPY public."character" (id, slug, name, description, img_name, insert_timestamp, search_vector, person_id) FROM stdin;
1	david-s-pumpkins-1	David S. Pumpkins	\N	\N	2025-01-31 14:31:12.727757	'david':1 'pumpkin':3	\N
2	dave-2	Dave	\N	\N	2025-01-31 14:31:12.727757	'dave':1	\N
3	test-char	Test Character	\N	\N	2025-01-31 14:31:12.727757	'charact':2 'test':1	\N
4	test-char-1	Test Character #1	\N	default-img.jpg	2025-01-31 14:31:12.727757	'1':3 'charact':2 'test':1	\N
5	test-char-2	Test Character #2	\N	default-img.jpg	2025-01-31 14:31:12.727757	'2':3 'charact':2 'test':1	\N
6	test-char-3	Test Character #3	\N	default-img.jpg	2025-01-31 14:31:12.727757	'3':3 'charact':2 'test':1	\N
7	test-char-4	Test Character #4	\N	default-img.jpg	2025-01-31 14:31:12.727757	'4':3 'charact':2 'test':1	\N
8	test-char-5	Test Character #5	\N	default-img.jpg	2025-01-31 14:31:12.727757	'5':3 'charact':2 'test':1	\N
9	test-char-6	Test Character #6	\N	default-img.jpg	2025-01-31 14:31:12.727757	'6':3 'charact':2 'test':1	\N
10	test-char-7	Test Character #7	\N	default-img.jpg	2025-01-31 14:31:12.727757	'7':3 'charact':2 'test':1	\N
11	test-char-8	Test Character #8	\N	default-img.jpg	2025-01-31 14:31:12.727757	'8':3 'charact':2 'test':1	\N
12	test-char-9	Test Character #9	\N	default-img.jpg	2025-01-31 14:31:12.727757	'9':3 'charact':2 'test':1	\N
13	test-char-10	Test Character #10	\N	default-img.jpg	2025-01-31 14:31:12.727757	'10':3 'charact':2 'test':1	\N
14	test-char-11	Test Character #11	\N	default-img.jpg	2025-01-31 14:31:12.727757	'11':3 'charact':2 'test':1	\N
15	test-char-12	Test Character #12	\N	default-img.jpg	2025-01-31 14:31:12.727757	'12':3 'charact':2 'test':1	\N
16	test-char-13	Test Character #13	\N	default-img.jpg	2025-01-31 14:31:12.727757	'13':3 'charact':2 'test':1	\N
17	test-char-14	Test Character #14	\N	default-img.jpg	2025-01-31 14:31:12.727757	'14':3 'charact':2 'test':1	\N
18	test-char-15	Test Character #15	\N	default-img.jpg	2025-01-31 14:31:12.727757	'15':3 'charact':2 'test':1	\N
19	test-char-16	Test Character #16	\N	default-img.jpg	2025-01-31 14:31:12.727757	'16':3 'charact':2 'test':1	\N
20	test-char-17	Test Character #17	\N	default-img.jpg	2025-01-31 14:31:12.727757	'17':3 'charact':2 'test':1	\N
21	test-char-18	Test Character #18	\N	default-img.jpg	2025-01-31 14:31:12.727757	'18':3 'charact':2 'test':1	\N
22	test-char-19	Test Character #19	\N	default-img.jpg	2025-01-31 14:31:12.727757	'19':3 'charact':2 'test':1	\N
23	test-char-20	Test Character #20	\N	default-img.jpg	2025-01-31 14:31:12.727757	'20':3 'charact':2 'test':1	\N
24	test-char-21	Test Character #21	\N	default-img.jpg	2025-01-31 14:31:12.727757	'21':3 'charact':2 'test':1	\N
25	test-char-22	Test Character #22	\N	default-img.jpg	2025-01-31 14:31:12.727757	'22':3 'charact':2 'test':1	\N
26	test-char-23	Test Character #23	\N	default-img.jpg	2025-01-31 14:31:12.727757	'23':3 'charact':2 'test':1	\N
27	test-char-24	Test Character #24	\N	default-img.jpg	2025-01-31 14:31:12.727757	'24':3 'charact':2 'test':1	\N
28	test-char-25	Test Character #25	\N	default-img.jpg	2025-01-31 14:31:12.727757	'25':3 'charact':2 'test':1	\N
29	test-char-26	Test Character #26	\N	default-img.jpg	2025-01-31 14:31:12.727757	'26':3 'charact':2 'test':1	\N
30	test-char-27	Test Character #27	\N	default-img.jpg	2025-01-31 14:31:12.727757	'27':3 'charact':2 'test':1	\N
31	test-char-28	Test Character #28	\N	default-img.jpg	2025-01-31 14:31:12.727757	'28':3 'charact':2 'test':1	\N
32	test-char-29	Test Character #29	\N	default-img.jpg	2025-01-31 14:31:12.727757	'29':3 'charact':2 'test':1	\N
33	test-char-30	Test Character #30	\N	default-img.jpg	2025-01-31 14:31:12.727757	'30':3 'charact':2 'test':1	\N
34	test-char-31	Test Character #31	\N	default-img.jpg	2025-01-31 14:31:12.727757	'31':3 'charact':2 'test':1	\N
35	test-char-32	Test Character #32	\N	default-img.jpg	2025-01-31 14:31:12.727757	'32':3 'charact':2 'test':1	\N
36	test-char-33	Test Character #33	\N	default-img.jpg	2025-01-31 14:31:12.727757	'33':3 'charact':2 'test':1	\N
37	test-char-34	Test Character #34	\N	default-img.jpg	2025-01-31 14:31:12.727757	'34':3 'charact':2 'test':1	\N
38	test-char-35	Test Character #35	\N	default-img.jpg	2025-01-31 14:31:12.727757	'35':3 'charact':2 'test':1	\N
39	test-char-36	Test Character #36	\N	default-img.jpg	2025-01-31 14:31:12.727757	'36':3 'charact':2 'test':1	\N
40	test-char-37	Test Character #37	\N	default-img.jpg	2025-01-31 14:31:12.727757	'37':3 'charact':2 'test':1	\N
41	test-char-38	Test Character #38	\N	default-img.jpg	2025-01-31 14:31:12.727757	'38':3 'charact':2 'test':1	\N
42	test-char-39	Test Character #39	\N	default-img.jpg	2025-01-31 14:31:12.727757	'39':3 'charact':2 'test':1	\N
43	test-char-40	Test Character #40	\N	default-img.jpg	2025-01-31 14:31:12.727757	'40':3 'charact':2 'test':1	\N
44	test-char-41	Test Character #41	\N	default-img.jpg	2025-01-31 14:31:12.727757	'41':3 'charact':2 'test':1	\N
45	test-char-42	Test Character #42	\N	default-img.jpg	2025-01-31 14:31:12.727757	'42':3 'charact':2 'test':1	\N
46	test-char-43	Test Character #43	\N	default-img.jpg	2025-01-31 14:31:12.727757	'43':3 'charact':2 'test':1	\N
47	test-char-44	Test Character #44	\N	default-img.jpg	2025-01-31 14:31:12.727757	'44':3 'charact':2 'test':1	\N
48	test-char-45	Test Character #45	\N	default-img.jpg	2025-01-31 14:31:12.727757	'45':3 'charact':2 'test':1	\N
49	test-char-46	Test Character #46	\N	default-img.jpg	2025-01-31 14:31:12.727757	'46':3 'charact':2 'test':1	\N
50	test-char-47	Test Character #47	\N	default-img.jpg	2025-01-31 14:31:12.727757	'47':3 'charact':2 'test':1	\N
51	test-char-48	Test Character #48	\N	default-img.jpg	2025-01-31 14:31:12.727757	'48':3 'charact':2 'test':1	\N
52	test-char-49	Test Character #49	\N	default-img.jpg	2025-01-31 14:31:12.727757	'49':3 'charact':2 'test':1	\N
53	test-char-50	Test Character #50	\N	default-img.jpg	2025-01-31 14:31:12.727757	'50':3 'charact':2 'test':1	\N
\.


--
-- Data for Name: show; Type: TABLE DATA; Schema: public; Owner: colet
--

COPY public.show (id, name, profile_img, slug) FROM stdin;
1	The Test Show	the-test-show.jpg	the-test-show
2	There's a new show on the block	fMUuYTB7_Wss1uRoe7_8cw.jpg	theres-a-new-show-on-the-block
4	ANOTHER SHOW YO12	9cdc0d9c-b1c0-455c-9c6d-75a674cebf93.jpg	this-is-a-test-XUQI0m5
\.


--
-- Data for Name: season; Type: TABLE DATA; Schema: public; Owner: colet
--

COPY public.season (id, show_id, season_number) FROM stdin;
1	1	1
2	1	2
3	1	3
4	1	4
5	4	1
6	4	2
7	4	3
8	4	4
9	4	5
10	4	6
\.


--
-- Data for Name: episode; Type: TABLE DATA; Schema: public; Owner: colet
--

COPY public.episode (id, season_id, episode_number, air_date, title, thumbnail_name) FROM stdin;
1	1	1	1975-01-01	\N	\N
2	5	1	2025-04-06	\N	8D7xTXqm9X2zV2DnTnxq5t.jpg
16	5	3	2025-04-02	adsfs	a2f14480-a86a-4153-8b3a-b34d2fc0d4b3.jpg
17	5	4	2025-04-16	asdf	d2123f56-68f9-4e1b-b391-1f09ed2b554a.jpg
21	5	6	2025-04-17	asdf	faa1c11d-9dcf-4ab6-8de0-dbc1210cabfa.jpg
23	5	7	2025-04-16	asdfasdf	9046d60e-55d3-4df0-a528-64fc14f61fa8.jpg
\.


--
-- Data for Name: video; Type: TABLE DATA; Schema: public; Owner: colet
--

COPY public.video (id, title, video_url, youtube_id, slug, thumbnail_name, description, upload_date, pg_rating, episode_id, part_number, sketch_number, search_vector, insert_timestamp) FROM stdin;
10	Test Video #10	localhost:4001	\N	test-video-10	missing-thumbnail.jpg	\N	\N	PG	\N	\N	\N	'10':3 'test':1 'video':2	\N
11	Test Video #11	localhost:4001	\N	test-video-11	missing-thumbnail.jpg	\N	\N	PG	\N	\N	\N	'11':3 'test':1 'video':2	\N
12	Test Video #12	localhost:4001	\N	test-video-12	missing-thumbnail.jpg	\N	\N	PG	\N	\N	\N	'12':3 'test':1 'video':2	\N
13	Test Video #13	localhost:4001	\N	test-video-13	missing-thumbnail.jpg	\N	\N	PG	\N	\N	\N	'13':3 'test':1 'video':2	\N
14	Test Video #14	localhost:4001	\N	test-video-14	missing-thumbnail.jpg	\N	\N	PG	\N	\N	\N	'14':3 'test':1 'video':2	\N
15	Test Video #15	localhost:4001	\N	test-video-15	missing-thumbnail.jpg	\N	\N	PG	\N	\N	\N	'15':3 'test':1 'video':2	\N
16	Test Video #16	localhost:4001	\N	test-video-16	missing-thumbnail.jpg	\N	\N	PG	\N	\N	\N	'16':3 'test':1 'video':2	\N
17	Test Video #17	localhost:4001	\N	test-video-17	missing-thumbnail.jpg	\N	\N	PG	\N	\N	\N	'17':3 'test':1 'video':2	\N
18	Test Video #18	localhost:4001	\N	test-video-18	missing-thumbnail.jpg	\N	\N	PG	\N	\N	\N	'18':3 'test':1 'video':2	\N
19	Test Video #19	localhost:4001	\N	test-video-19	missing-thumbnail.jpg	\N	\N	PG	\N	\N	\N	'19':3 'test':1 'video':2	\N
20	Test Video #20	localhost:4001	\N	test-video-20	missing-thumbnail.jpg	\N	\N	PG	\N	\N	\N	'20':3 'test':1 'video':2	\N
21	Test Video #21	localhost:4001	\N	test-video-21	missing-thumbnail.jpg	\N	\N	PG	\N	\N	\N	'21':3 'test':1 'video':2	\N
22	Test Video #22	localhost:4001	\N	test-video-22	missing-thumbnail.jpg	\N	\N	PG	\N	\N	\N	'22':3 'test':1 'video':2	\N
23	Test Video #23	localhost:4001	\N	test-video-23	missing-thumbnail.jpg	\N	\N	PG	\N	\N	\N	'23':3 'test':1 'video':2	\N
24	Test Video #24	localhost:4001	\N	test-video-24	missing-thumbnail.jpg	\N	\N	PG	\N	\N	\N	'24':3 'test':1 'video':2	\N
25	Test Video #25	localhost:4001	\N	test-video-25	missing-thumbnail.jpg	\N	\N	PG	\N	\N	\N	'25':3 'test':1 'video':2	\N
26	Test Video #26	localhost:4001	\N	test-video-26	missing-thumbnail.jpg	\N	\N	PG	\N	\N	\N	'26':3 'test':1 'video':2	\N
27	Test Video #27	localhost:4001	\N	test-video-27	missing-thumbnail.jpg	\N	\N	PG	\N	\N	\N	'27':3 'test':1 'video':2	\N
28	Test Video #28	localhost:4001	\N	test-video-28	missing-thumbnail.jpg	\N	\N	PG	\N	\N	\N	'28':3 'test':1 'video':2	\N
29	Test Video #29	localhost:4001	\N	test-video-29	missing-thumbnail.jpg	\N	\N	PG	\N	\N	\N	'29':3 'test':1 'video':2	\N
30	Test Video #30	localhost:4001	\N	test-video-30	missing-thumbnail.jpg	\N	\N	PG	\N	\N	\N	'30':3 'test':1 'video':2	\N
31	Test Video #31	localhost:4001	\N	test-video-31	missing-thumbnail.jpg	\N	\N	PG	\N	\N	\N	'31':3 'test':1 'video':2	\N
32	Test Video #32	localhost:4001	\N	test-video-32	missing-thumbnail.jpg	\N	\N	PG	\N	\N	\N	'32':3 'test':1 'video':2	\N
33	Test Video #33	localhost:4001	\N	test-video-33	missing-thumbnail.jpg	\N	\N	PG	\N	\N	\N	'33':3 'test':1 'video':2	\N
34	Test Video #34	localhost:4001	\N	test-video-34	missing-thumbnail.jpg	\N	\N	PG	\N	\N	\N	'34':3 'test':1 'video':2	\N
35	Test Video #35	localhost:4001	\N	test-video-35	missing-thumbnail.jpg	\N	\N	PG	\N	\N	\N	'35':3 'test':1 'video':2	\N
36	Test Video #36	localhost:4001	\N	test-video-36	missing-thumbnail.jpg	\N	\N	PG	\N	\N	\N	'36':3 'test':1 'video':2	\N
37	Test Video #37	localhost:4001	\N	test-video-37	missing-thumbnail.jpg	\N	\N	PG	\N	\N	\N	'37':3 'test':1 'video':2	\N
38	Test Video #38	localhost:4001	\N	test-video-38	missing-thumbnail.jpg	\N	\N	PG	\N	\N	\N	'38':3 'test':1 'video':2	\N
39	Test Video #39	localhost:4001	\N	test-video-39	missing-thumbnail.jpg	\N	\N	PG	\N	\N	\N	'39':3 'test':1 'video':2	\N
40	Test Video #40	localhost:4001	\N	test-video-40	missing-thumbnail.jpg	\N	\N	PG	\N	\N	\N	'40':3 'test':1 'video':2	\N
41	Test Video #41	localhost:4001	\N	test-video-41	missing-thumbnail.jpg	\N	\N	PG	\N	\N	\N	'41':3 'test':1 'video':2	\N
42	Test Video #42	localhost:4001	\N	test-video-42	missing-thumbnail.jpg	\N	\N	PG	\N	\N	\N	'42':3 'test':1 'video':2	\N
43	Test Video #43	localhost:4001	\N	test-video-43	missing-thumbnail.jpg	\N	\N	PG	\N	\N	\N	'43':3 'test':1 'video':2	\N
44	Test Video #44	localhost:4001	\N	test-video-44	missing-thumbnail.jpg	\N	\N	PG	\N	\N	\N	'44':3 'test':1 'video':2	\N
45	Test Video #45	localhost:4001	\N	test-video-45	missing-thumbnail.jpg	\N	\N	PG	\N	\N	\N	'45':3 'test':1 'video':2	\N
46	Test Video #46	localhost:4001	\N	test-video-46	missing-thumbnail.jpg	\N	\N	PG	\N	\N	\N	'46':3 'test':1 'video':2	\N
47	Test Video #47	localhost:4001	\N	test-video-47	missing-thumbnail.jpg	\N	\N	PG	\N	\N	\N	'47':3 'test':1 'video':2	\N
48	Test Video #48	localhost:4001	\N	test-video-48	missing-thumbnail.jpg	\N	\N	PG	\N	\N	\N	'48':3 'test':1 'video':2	\N
49	Test Video #49	localhost:4001	\N	test-video-49	missing-thumbnail.jpg	\N	\N	PG	\N	\N	\N	'49':3 'test':1 'video':2	\N
50	Test Video #50	localhost:4001	\N	test-video-50	missing-thumbnail.jpg	\N	\N	PG	\N	\N	\N	'50':3 'test':1 'video':2	\N
1	Test Video #1	localhost:4001	6aTqXkZHnQE	test-video-1	test-video-1.jpg	\N	2008-09-08	PG	2	\N	\N	'1':3 'test':1 'video':2	\N
2	Test Video #2	localhost:4001	\N	test-video-2	test-video-2.jpg	\N	2008-09-09	PG	2	\N	\N	'2':3 'test':1 'video':2	\N
3	Test Video #3 | A Long Title to Clamp for Those who Enjoy It	localhost:4001	\N	test-video-3	test-video-3.jpg	\N	2008-09-10	PG	2	\N	\N	'3':3 'clamp':8 'enjoy':12 'long':5 'test':1 'titl':6 'video':2	\N
5	Test Video #5	localhost:4001	\N	test-video-5	test-video-5.jpg	\N	2008-09-12	PG	2	\N	\N	'5':3 'test':1 'video':2	\N
6	Test Video #6	localhost:4001	\N	test-video-6	missing-thumbnail.jpg	\N	2008-09-13	PG	2	\N	\N	'6':3 'test':1 'video':2	\N
8	Test Video #8	localhost:4001	\N	test-video-8	missing-thumbnail.jpg	\N	\N	PG	2	\N	\N	'8':3 'test':1 'video':2	\N
9	Test Video #9	localhost:4001	\N	test-video-9	missing-thumbnail.jpg	\N	\N	PG	2	\N	\N	'9':3 'test':1 'video':2	\N
7	Test Video #7	localhost:4001	\N	test-video-7	missing-thumbnail.jpg	\N	2008-09-14	PG	2	\N	\N	'7':3 'test':1 'video':2	\N
4	Test Video #4	localhost:4001	\N	test-video-4	test-video-4.jpg	\N	2008-09-11	PG	2	\N	\N	'4':3 'test':1 'video':2	\N
\.


--
-- Data for Name: cast_members; Type: TABLE DATA; Schema: public; Owner: colet
--

COPY public.cast_members (id, video_id, person_id, character_name, character_id, "position", img_name, insert_timestamp, role) FROM stdin;
1	1	1	Kyle	\N	2	kyle-mooney-1.jpg	2025-01-31 14:31:12.730178	\N
2	1	2	David Pumpkins	1	1	nathan-fielder-2.jpg	2025-01-31 14:31:12.730178	\N
7	1	2	Tim	\N	\N	IsQUx4k2x-xX7n7z1p8SUl.jpg	2025-02-22 13:14:27.274317	\N
14	1	3	James	\N	\N	rx5fXPCwCpENbvliMwXItm.jpg	2025-02-22 13:58:10.535543	\N
15	1	3	James	2	\N	Q40ZKN0DMMku5A45waUtUL.jpg	2025-02-22 14:00:16.17135	\N
16	2	1	Kyle	\N	\N	jEe6UcyOr9TTqLQEnG6Z_U.jpg	2025-03-13 15:14:36.720926	\N
17	3	1	Kyle	\N	\N	JjXvjabAND5Hhl94cQJWfj.jpg	2025-03-13 15:16:08.37746	\N
18	4	1	Kyle	\N	\N	k4RWVHzLeak2mcsd_f2bvi.jpg	2025-03-13 15:19:34.598128	cast
\.


--
-- Data for Name: categories; Type: TABLE DATA; Schema: public; Owner: colet
--

COPY public.categories (id, created_at, name, slug) FROM stdin;
9	2025-02-22 23:59:26-08	Movies	movies
10	2025-02-23 14:48:41-08	Action	action
11	2025-02-23 14:49:31-08	Romance	romance
\.


--
-- Data for Name: creator; Type: TABLE DATA; Schema: public; Owner: colet
--

COPY public.creator (id, name, slug, page_url, description, profile_img, date_established, search_vector, insert_timestamp) FROM stdin;
1	nathanfielder	nathanfielder-1	https://www.youtube.com/@nathanfielder	\N	nathanfielder-1.jpg	2006-10-16	'nathanfield':1	2025-01-31 14:31:12.72841
2	A Long Ass Creator Name that May Certainly cause Problems	long-ass-name-2	localhost:4000	\N	missing-profile.jpg	2024-12-31	'ass':3 'caus':9 'certain':8 'creator':4 'long':2 'may':7 'name':5 'problem':10	2025-01-31 14:31:12.72841
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: colet
--

COPY public.users (id, created_at, username, email, password_hash, activated, role) FROM stdin;
1	2025-01-31 14:32:02-08	admin	admin@admin.com	\\x24326124313224746d637467354f556949737469506134562e643837756356785636544e313849507a335158637241373644453341564d4642443043	t	admin
2	2025-01-31 14:33:30-08	curator	curator@curator.com	\\x24326124313224754a7142786e626d622f4a617a54764537474f64652e78765276387534594c704857656a69325557765377494b59574471726c3632	t	editor
3	2025-01-31 14:34:55-08	viewer	viewer@viewer.com	\\x24326124313224695364626d6c316c323137586e763963507164635665767a795853764f436362785a643251796c6c50714f42517839484762673032	t	viewer
\.


--
-- Data for Name: likes; Type: TABLE DATA; Schema: public; Owner: colet
--

COPY public.likes (created_at, user_id, video_id) FROM stdin;
2025-02-17 22:52:39-08	1	1
2025-03-02 12:46:32-08	1	2
\.


--
-- Data for Name: sessions; Type: TABLE DATA; Schema: public; Owner: colet
--

COPY public.sessions (token, data, expiry) FROM stdin;
6MGUSzIHCjUuxRF5dVW_MZybHWbDF3TYfNA5No7fFHo	\\x257f030102ff800001020108446561646c696e6501ff8200010656616c75657301ff8400000010ff810501010454696d6501ff8200000027ff83040101176d61705b737472696e675d696e74657266616365207b7d01ff8400010c0110000032ff80010f010000000edf8e87871d3e2849ffff01011361757468656e7469636174656455736572494403696e740402000200	2025-04-13 21:35:19.490612-07
\.


--
-- Data for Name: tags; Type: TABLE DATA; Schema: public; Owner: colet
--

COPY public.tags (id, created_at, name, slug, category_id) FROM stdin;
2	2025-02-26 22:04:22-08	Transformers	action-transformers	10
1	2025-02-23 14:45:47-08	Pulp Fiction	movies-pulp-fiction	9
\.


--
-- Data for Name: video_creator_rel; Type: TABLE DATA; Schema: public; Owner: colet
--

COPY public.video_creator_rel (creator_id, video_id, "position", insert_timestamp) FROM stdin;
1	2	\N	2025-01-31 14:31:12.729426
2	3	\N	2025-01-31 14:31:12.729426
1	4	\N	2025-01-31 14:31:12.729426
2	5	\N	2025-01-31 14:31:12.729426
2	1	\N	2025-01-31 14:31:12.729426
\.


--
-- Data for Name: video_tags; Type: TABLE DATA; Schema: public; Owner: colet
--

COPY public.video_tags (video_id, tag_id) FROM stdin;
1	1
1	2
\.


--
-- Name: cast_members_id_seq; Type: SEQUENCE SET; Schema: public; Owner: colet
--

SELECT pg_catalog.setval('public.cast_members_id_seq', 1, false);


--
-- Name: categories_id_seq; Type: SEQUENCE SET; Schema: public; Owner: colet
--

SELECT pg_catalog.setval('public.categories_id_seq', 11, true);


--
-- Name: character_id_seq; Type: SEQUENCE SET; Schema: public; Owner: colet
--

SELECT pg_catalog.setval('public.character_id_seq', 2, false);


--
-- Name: creator_id_seq; Type: SEQUENCE SET; Schema: public; Owner: colet
--

SELECT pg_catalog.setval('public.creator_id_seq', 13, true);


--
-- Name: episode_id_seq; Type: SEQUENCE SET; Schema: public; Owner: colet
--

SELECT pg_catalog.setval('public.episode_id_seq', 23, true);


--
-- Name: person_id_seq; Type: SEQUENCE SET; Schema: public; Owner: colet
--

SELECT pg_catalog.setval('public.person_id_seq', 4, true);


--
-- Name: season_id_seq; Type: SEQUENCE SET; Schema: public; Owner: colet
--

SELECT pg_catalog.setval('public.season_id_seq', 10, true);


--
-- Name: show_id_seq; Type: SEQUENCE SET; Schema: public; Owner: colet
--

SELECT pg_catalog.setval('public.show_id_seq', 4, true);


--
-- Name: tags_id_seq; Type: SEQUENCE SET; Schema: public; Owner: colet
--

SELECT pg_catalog.setval('public.tags_id_seq', 2, true);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: colet
--

SELECT pg_catalog.setval('public.users_id_seq', 3, true);


--
-- Name: video_id_seq; Type: SEQUENCE SET; Schema: public; Owner: colet
--

SELECT pg_catalog.setval('public.video_id_seq', 51, true);


--
-- PostgreSQL database dump complete
--

