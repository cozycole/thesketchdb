--
-- PostgreSQL database dump
--

-- Dumped from database version 14.15 (Ubuntu 14.15-0ubuntu0.22.04.1)
-- Dumped by pg_dump version 14.15 (Ubuntu 14.15-0ubuntu0.22.04.1)

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
-- Data for Name: creator; Type: TABLE DATA; Schema: public; Owner: colet
--

COPY public.creator (id, name, slug, page_url, description, profile_img, date_established, search_vector, insert_timestamp) FROM stdin;
1	nathanfielder	nathanfielder-1	https://www.youtube.com/@nathanfielder	\N	nathanfielder-1.jpg	2006-10-16	'nathanfield':1	2025-01-31 14:31:12.72841
2	A Long Ass Creator Name that May Certainly cause Problems	long-ass-name-2	localhost:4000	\N	default-img.jpg	2024-12-31	'ass':3 'caus':9 'certain':8 'creator':4 'long':2 'may':7 'name':5 'problem':10	2025-01-31 14:31:12.72841
\.


--
-- Data for Name: sessions; Type: TABLE DATA; Schema: public; Owner: colet
--

COPY public.sessions (token, data, expiry) FROM stdin;
-aASf1wg96eFzycvtsIZzOaxab4mkcunx1bDBGyUALk	\\x257f030102ff800001020108446561646c696e6501ff8200010656616c75657301ff8400000010ff810501010454696d6501ff8200000027ff83040101176d61705b737472696e675d696e74657266616365207b7d01ff8400010c0110000016ff80010f010000000edf2fef2200568f24ffff010000	2025-02-01 02:32:02.005672-08
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
-- Data for Name: video; Type: TABLE DATA; Schema: public; Owner: colet
--

COPY public.video (id, title, video_url, youtube_id, slug, thumbnail_name, description, upload_date, pg_rating, search_vector, insert_timestamp) FROM stdin;
1	Test Video #1	localhost:4001	6aTqXkZHnQE	test-video-1	test-video-1.jpg	\N	2008-09-08	PG	'1':3 'test':1 'video':2	\N
2	Test Video #2	localhost:4001	\N	test-video-2	test-video-2.jpg	\N	2008-09-08	PG	'2':3 'test':1 'video':2	\N
3	Test Video #3 | A Long Title to Clamp for Those who Enjoy It	localhost:4001	\N	test-video-3	test-video-3.jpg	\N	2008-09-08	PG	'3':3 'clamp':8 'enjoy':12 'long':5 'test':1 'titl':6 'video':2	\N
4	Test Video #4	localhost:4001	\N	test-video-4	test-video-4.jpg	\N	2008-09-08	PG	'4':3 'test':1 'video':2	\N
5	Test Video #5	localhost:4001	\N	test-video-5	test-video-5.jpg	\N	2008-09-08	PG	'5':3 'test':1 'video':2	\N
6	Test Video #6	localhost:4001	\N	test-video-6	test-video-6.jpg	\N	2008-09-08	PG	'6':3 'test':1 'video':2	\N
7	Test Video #7	localhost:4001	\N	test-video-7	test-video-7.jpg	\N	2008-09-08	PG	'7':3 'test':1 'video':2	\N
8	Test Video #8	localhost:4001	\N	test-video-8	test-video-8.jpg	\N	\N	PG	'8':3 'test':1 'video':2	\N
9	Test Video #9	localhost:4001	\N	test-video-9	test-video-9.jpg	\N	\N	PG	'9':3 'test':1 'video':2	\N
10	Test Video #10	localhost:4001	\N	test-video-10	test-video-10.jpg	\N	\N	PG	'10':3 'test':1 'video':2	\N
11	Test Video #11	localhost:4001	\N	test-video-11	test-video-11.jpg	\N	\N	PG	'11':3 'test':1 'video':2	\N
12	Test Video #12	localhost:4001	\N	test-video-12	test-video-12.jpg	\N	\N	PG	'12':3 'test':1 'video':2	\N
13	Test Video #13	localhost:4001	\N	test-video-13	test-video-13.jpg	\N	\N	PG	'13':3 'test':1 'video':2	\N
14	Test Video #14	localhost:4001	\N	test-video-14	test-video-14.jpg	\N	\N	PG	'14':3 'test':1 'video':2	\N
15	Test Video #15	localhost:4001	\N	test-video-15	test-video-15.jpg	\N	\N	PG	'15':3 'test':1 'video':2	\N
16	Test Video #16	localhost:4001	\N	test-video-16	test-video-16.jpg	\N	\N	PG	'16':3 'test':1 'video':2	\N
17	Test Video #17	localhost:4001	\N	test-video-17	test-video-17.jpg	\N	\N	PG	'17':3 'test':1 'video':2	\N
18	Test Video #18	localhost:4001	\N	test-video-18	test-video-18.jpg	\N	\N	PG	'18':3 'test':1 'video':2	\N
19	Test Video #19	localhost:4001	\N	test-video-19	test-video-19.jpg	\N	\N	PG	'19':3 'test':1 'video':2	\N
20	Test Video #20	localhost:4001	\N	test-video-20	test-video-20.jpg	\N	\N	PG	'20':3 'test':1 'video':2	\N
21	Test Video #21	localhost:4001	\N	test-video-21	test-video-21.jpg	\N	\N	PG	'21':3 'test':1 'video':2	\N
22	Test Video #22	localhost:4001	\N	test-video-22	test-video-22.jpg	\N	\N	PG	'22':3 'test':1 'video':2	\N
23	Test Video #23	localhost:4001	\N	test-video-23	test-video-23.jpg	\N	\N	PG	'23':3 'test':1 'video':2	\N
24	Test Video #24	localhost:4001	\N	test-video-24	test-video-24.jpg	\N	\N	PG	'24':3 'test':1 'video':2	\N
25	Test Video #25	localhost:4001	\N	test-video-25	test-video-25.jpg	\N	\N	PG	'25':3 'test':1 'video':2	\N
26	Test Video #26	localhost:4001	\N	test-video-26	test-video-26.jpg	\N	\N	PG	'26':3 'test':1 'video':2	\N
27	Test Video #27	localhost:4001	\N	test-video-27	test-video-27.jpg	\N	\N	PG	'27':3 'test':1 'video':2	\N
28	Test Video #28	localhost:4001	\N	test-video-28	test-video-28.jpg	\N	\N	PG	'28':3 'test':1 'video':2	\N
29	Test Video #29	localhost:4001	\N	test-video-29	test-video-29.jpg	\N	\N	PG	'29':3 'test':1 'video':2	\N
30	Test Video #30	localhost:4001	\N	test-video-30	test-video-30.jpg	\N	\N	PG	'30':3 'test':1 'video':2	\N
31	Test Video #31	localhost:4001	\N	test-video-31	test-video-31.jpg	\N	\N	PG	'31':3 'test':1 'video':2	\N
32	Test Video #32	localhost:4001	\N	test-video-32	test-video-32.jpg	\N	\N	PG	'32':3 'test':1 'video':2	\N
33	Test Video #33	localhost:4001	\N	test-video-33	test-video-33.jpg	\N	\N	PG	'33':3 'test':1 'video':2	\N
34	Test Video #34	localhost:4001	\N	test-video-34	test-video-34.jpg	\N	\N	PG	'34':3 'test':1 'video':2	\N
35	Test Video #35	localhost:4001	\N	test-video-35	test-video-35.jpg	\N	\N	PG	'35':3 'test':1 'video':2	\N
36	Test Video #36	localhost:4001	\N	test-video-36	test-video-36.jpg	\N	\N	PG	'36':3 'test':1 'video':2	\N
37	Test Video #37	localhost:4001	\N	test-video-37	test-video-37.jpg	\N	\N	PG	'37':3 'test':1 'video':2	\N
38	Test Video #38	localhost:4001	\N	test-video-38	test-video-38.jpg	\N	\N	PG	'38':3 'test':1 'video':2	\N
39	Test Video #39	localhost:4001	\N	test-video-39	test-video-39.jpg	\N	\N	PG	'39':3 'test':1 'video':2	\N
40	Test Video #40	localhost:4001	\N	test-video-40	test-video-40.jpg	\N	\N	PG	'40':3 'test':1 'video':2	\N
41	Test Video #41	localhost:4001	\N	test-video-41	test-video-41.jpg	\N	\N	PG	'41':3 'test':1 'video':2	\N
42	Test Video #42	localhost:4001	\N	test-video-42	test-video-42.jpg	\N	\N	PG	'42':3 'test':1 'video':2	\N
43	Test Video #43	localhost:4001	\N	test-video-43	test-video-43.jpg	\N	\N	PG	'43':3 'test':1 'video':2	\N
44	Test Video #44	localhost:4001	\N	test-video-44	test-video-44.jpg	\N	\N	PG	'44':3 'test':1 'video':2	\N
45	Test Video #45	localhost:4001	\N	test-video-45	test-video-45.jpg	\N	\N	PG	'45':3 'test':1 'video':2	\N
46	Test Video #46	localhost:4001	\N	test-video-46	test-video-46.jpg	\N	\N	PG	'46':3 'test':1 'video':2	\N
47	Test Video #47	localhost:4001	\N	test-video-47	test-video-47.jpg	\N	\N	PG	'47':3 'test':1 'video':2	\N
48	Test Video #48	localhost:4001	\N	test-video-48	test-video-48.jpg	\N	\N	PG	'48':3 'test':1 'video':2	\N
49	Test Video #49	localhost:4001	\N	test-video-49	test-video-49.jpg	\N	\N	PG	'49':3 'test':1 'video':2	\N
50	Test Video #50	localhost:4001	\N	test-video-50	test-video-50.jpg	\N	\N	PG	'50':3 'test':1 'video':2	\N
\.


--
-- Data for Name: video_creator_rel; Type: TABLE DATA; Schema: public; Owner: colet
--

COPY public.video_creator_rel (creator_id, video_id, "position", insert_timestamp) FROM stdin;
1	1	\N	2025-01-31 14:31:12.729426
1	2	\N	2025-01-31 14:31:12.729426
2	3	\N	2025-01-31 14:31:12.729426
1	4	\N	2025-01-31 14:31:12.729426
2	5	\N	2025-01-31 14:31:12.729426
\.


--
-- Data for Name: cast; Type: TABLE DATA; Schema: public; Owner: colet
--

COPY public.cast (id, video_id, person_id, character_id, character_name, "position", img_name, insert_timestamp) FROM stdin;
1	1	1	\N	2	\N	2025-01-31 14:31:12.730178
2	1	2	1	1	\N	2025-01-31 14:31:12.730178
\.


--
-- Name: character_id_seq; Type: SEQUENCE SET; Schema: public; Owner: colet
--

SELECT pg_catalog.setval('public.character_id_seq', 2, false);


--
-- Name: creator_id_seq; Type: SEQUENCE SET; Schema: public; Owner: colet
--

SELECT pg_catalog.setval('public.creator_id_seq', 13, true);


--
-- Name: person_id_seq; Type: SEQUENCE SET; Schema: public; Owner: colet
--

SELECT pg_catalog.setval('public.person_id_seq', 4, true);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: colet
--

SELECT pg_catalog.setval('public.users_id_seq', 3, true);


--
-- Name: video_id_seq; Type: SEQUENCE SET; Schema: public; Owner: colet
--

SELECT pg_catalog.setval('public.video_id_seq', 51, true);


--
-- Name: cast_id_seq; Type: SEQUENCE SET; Schema: public; Owner: colet
--

SELECT pg_catalog.setval('public.cast_id_seq', 2, true);


--
-- PostgreSQL database dump complete
--

