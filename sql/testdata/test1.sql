--
-- PostgreSQL database dump
--

-- Dumped from database version 14.13 (Ubuntu 14.13-0ubuntu0.22.04.1)
-- Dumped by pg_dump version 14.13 (Ubuntu 14.13-0ubuntu0.22.04.1)

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

COPY public.person (id, slug, first, last, birthdate, profile_img, description) FROM stdin;
1	kyle-mooney-1	Kyle	Mooney	1984-09-03	kyle-mooney-1.jpg	\N
2	tim-gilbert-4	Tim	Gilbert	1983-05-13	tim-gilbert-4.jpg	this is the description
3	james-hartnett-5	James	Hartnett	\N	james-hartnett-5.jpg	\N
4	test-alpha-4	Test	Alpha	\N	james-hartnett-5.jpg	\N
5	test-beta-5	Test	Beta	1983-05-13	tim-gilbert-4.jpg	this is the description
6	test-charlie-6	Test	Charlie	1984-09-03	kyle-mooney-1.jpg	\N
7	test-delta-6	Test	Delta	1984-09-03	kyle-mooney-1.jpg	\N
\.


--
-- Data for Name: character; Type: TABLE DATA; Schema: public; Owner: colet
--
COPY public."character" (id, name, description, img_name, person_id, slug) FROM stdin;
1	David S. Pumpkins	\N	\N	\N	david-s-pumpkins-1
2	Dave	\N	\N	\N	dave-2
3	Test Character	\N	\N	\N	test-char
4	Test Character #1	\N	default-img.jpg	\N	test-char-1
5	Test Character #2	\N	default-img.jpg	\N	test-char-2
6	Test Character #3	\N	default-img.jpg	\N	test-char-3
7	Test Character #4	\N	default-img.jpg	\N	test-char-4
8	Test Character #5	\N	default-img.jpg	\N	test-char-5
9	Test Character #6	\N	default-img.jpg	\N	test-char-6
10	Test Character #7	\N	default-img.jpg	\N	test-char-7
11	Test Character #8	\N	default-img.jpg	\N	test-char-8
12	Test Character #9	\N	default-img.jpg	\N	test-char-9
13	Test Character #10	\N	default-img.jpg	\N	test-char-10
14	Test Character #11	\N	default-img.jpg	\N	test-char-11
15	Test Character #12	\N	default-img.jpg	\N	test-char-12
16	Test Character #13	\N	default-img.jpg	\N	test-char-13
17	Test Character #14	\N	default-img.jpg	\N	test-char-14
18	Test Character #15	\N	default-img.jpg	\N	test-char-15
19	Test Character #16	\N	default-img.jpg	\N	test-char-16
20	Test Character #17	\N	default-img.jpg	\N	test-char-17
21	Test Character #18	\N	default-img.jpg	\N	test-char-18
22	Test Character #19	\N	default-img.jpg	\N	test-char-19
23	Test Character #20	\N	default-img.jpg	\N	test-char-20
24	Test Character #21	\N	default-img.jpg	\N	test-char-21
25	Test Character #22	\N	default-img.jpg	\N	test-char-22
26	Test Character #23	\N	default-img.jpg	\N	test-char-23
27	Test Character #24	\N	default-img.jpg	\N	test-char-24
28	Test Character #25	\N	default-img.jpg	\N	test-char-25
29	Test Character #26	\N	default-img.jpg	\N	test-char-26
30	Test Character #27	\N	default-img.jpg	\N	test-char-27
31	Test Character #28	\N	default-img.jpg	\N	test-char-28
32	Test Character #29	\N	default-img.jpg	\N	test-char-29
33	Test Character #30	\N	default-img.jpg	\N	test-char-30
34	Test Character #31	\N	default-img.jpg	\N	test-char-31
35	Test Character #32	\N	default-img.jpg	\N	test-char-32
36	Test Character #33	\N	default-img.jpg	\N	test-char-33
37	Test Character #34	\N	default-img.jpg	\N	test-char-34
38	Test Character #35	\N	default-img.jpg	\N	test-char-35
39	Test Character #36	\N	default-img.jpg	\N	test-char-36
40	Test Character #37	\N	default-img.jpg	\N	test-char-37
41	Test Character #38	\N	default-img.jpg	\N	test-char-38
42	Test Character #39	\N	default-img.jpg	\N	test-char-39
43	Test Character #40	\N	default-img.jpg	\N	test-char-40
44	Test Character #41	\N	default-img.jpg	\N	test-char-41
45	Test Character #42	\N	default-img.jpg	\N	test-char-42
46	Test Character #43	\N	default-img.jpg	\N	test-char-43
47	Test Character #44	\N	default-img.jpg	\N	test-char-44
48	Test Character #45	\N	default-img.jpg	\N	test-char-45
49	Test Character #46	\N	default-img.jpg	\N	test-char-46
50	Test Character #47	\N	default-img.jpg	\N	test-char-47
51	Test Character #48	\N	default-img.jpg	\N	test-char-48
52	Test Character #49	\N	default-img.jpg	\N	test-char-49
53	Test Character #50	\N	default-img.jpg	\N	test-char-50
\.


--
-- Data for Name: creator; Type: TABLE DATA; Schema: public; Owner: colet
--

COPY public.creator (id, name, profile_img, page_url, date_established, slug, description) FROM stdin;
1	nathanfielder	nathanfielder-1.jpg	https://www.youtube.com/@nathanfielder	2006-10-16	nathanfielder-1	\N
2	A Long Ass Creator Name that May Certainly cause Problems	default-img.jpg	localhost:4000	2024-12-31	long-ass-name-2	\N
\.



--
-- Data for Name: video; Type: TABLE DATA; Schema: public; Owner: colet
--

COPY public.video (id, title, video_url, thumbnail_name, upload_date, pg_rating,  insert_timestamp, slug) FROM stdin;
1	Test Video #1	localhost:4001	test-video-1.jpg	2008-09-08	PG	\N	test-video-1
2	Test Video #2	localhost:4001	test-video-2.jpg	2008-09-08	PG	\N	test-video-2
3	Test Video #3 | A Long Title to Clamp for Those who Enjoy It	localhost:4001	test-video-3.jpg	2008-09-08	PG	\N	test-video-3
4	Test Video #4	localhost:4001	test-video-4.jpg	2008-09-08	PG	\N	test-video-4
5	Test Video #5	localhost:4001	test-video-5.jpg	2008-09-08	PG	\N	test-video-5
6	Test Video #6	localhost:4001	test-video-6.jpg	2008-09-08	PG	\N	test-video-6
7	Test Video #7	localhost:4001	test-video-7.jpg	2008-09-08	PG	\N	test-video-7
8	Test Video #8	localhost:4001	test-video-8.jpg	\N	PG	\N	test-video-8
9	Test Video #9	localhost:4001	test-video-9.jpg	\N	PG	\N	test-video-9
10	Test Video #10	localhost:4001	test-video-10.jpg	\N	PG	\N	test-video-10
11	Test Video #11	localhost:4001	test-video-11.jpg	\N	PG	\N	test-video-11
12	Test Video #12	localhost:4001	test-video-12.jpg	\N	PG	\N	test-video-12
13	Test Video #13	localhost:4001	test-video-13.jpg	\N	PG	\N	test-video-13
14	Test Video #14	localhost:4001	test-video-14.jpg	\N	PG	\N	test-video-14
15	Test Video #15	localhost:4001	test-video-15.jpg	\N	PG	\N	test-video-15
16	Test Video #16	localhost:4001	test-video-16.jpg	\N	PG	\N	test-video-16
17	Test Video #17	localhost:4001	test-video-17.jpg	\N	PG	\N	test-video-17
18	Test Video #18	localhost:4001	test-video-18.jpg	\N	PG	\N	test-video-18
19	Test Video #19	localhost:4001	test-video-19.jpg	\N	PG	\N	test-video-19
20	Test Video #20	localhost:4001	test-video-20.jpg	\N	PG	\N	test-video-20
21	Test Video #21	localhost:4001	test-video-21.jpg	\N	PG	\N	test-video-21
22	Test Video #22	localhost:4001	test-video-22.jpg	\N	PG	\N	test-video-22
23	Test Video #23	localhost:4001	test-video-23.jpg	\N	PG	\N	test-video-23
24	Test Video #24	localhost:4001	test-video-24.jpg	\N	PG	\N	test-video-24
25	Test Video #25	localhost:4001	test-video-25.jpg	\N	PG	\N	test-video-25
26	Test Video #26	localhost:4001	test-video-26.jpg	\N	PG	\N	test-video-26
27	Test Video #27	localhost:4001	test-video-27.jpg	\N	PG	\N	test-video-27
28	Test Video #28	localhost:4001	test-video-28.jpg	\N	PG	\N	test-video-28
29	Test Video #29	localhost:4001	test-video-29.jpg	\N	PG	\N	test-video-29
30	Test Video #30	localhost:4001	test-video-30.jpg	\N	PG	\N	test-video-30
31	Test Video #31	localhost:4001	test-video-31.jpg	\N	PG	\N	test-video-31
32	Test Video #32	localhost:4001	test-video-32.jpg	\N	PG	\N	test-video-32
33	Test Video #33	localhost:4001	test-video-33.jpg	\N	PG	\N	test-video-33
34	Test Video #34	localhost:4001	test-video-34.jpg	\N	PG	\N	test-video-34
35	Test Video #35	localhost:4001	test-video-35.jpg	\N	PG	\N	test-video-35
36	Test Video #36	localhost:4001	test-video-36.jpg	\N	PG	\N	test-video-36
37	Test Video #37	localhost:4001	test-video-37.jpg	\N	PG	\N	test-video-37
38	Test Video #38	localhost:4001	test-video-38.jpg	\N	PG	\N	test-video-38
39	Test Video #39	localhost:4001	test-video-39.jpg	\N	PG	\N	test-video-39
40	Test Video #40	localhost:4001	test-video-40.jpg	\N	PG	\N	test-video-40
41	Test Video #41	localhost:4001	test-video-41.jpg	\N	PG	\N	test-video-41
42	Test Video #42	localhost:4001	test-video-42.jpg	\N	PG	\N	test-video-42
43	Test Video #43	localhost:4001	test-video-43.jpg	\N	PG	\N	test-video-43
44	Test Video #44	localhost:4001	test-video-44.jpg	\N	PG	\N	test-video-44
45	Test Video #45	localhost:4001	test-video-45.jpg	\N	PG	\N	test-video-45
46	Test Video #46	localhost:4001	test-video-46.jpg	\N	PG	\N	test-video-46
47	Test Video #47	localhost:4001	test-video-47.jpg	\N	PG	\N	test-video-47
48	Test Video #48	localhost:4001	test-video-48.jpg	\N	PG	\N	test-video-48
49	Test Video #49	localhost:4001	test-video-49.jpg	\N	PG	\N	test-video-49
50	Test Video #50	localhost:4001	test-video-50.jpg	\N	PG	\N	test-video-50
\.


--
-- Data for Name: video_creator_rel; Type: TABLE DATA; Schema: public; Owner: colet
--

COPY public.video_creator_rel (creator_id, video_id) FROM stdin;
1	1
1	2
2	3
1	4
2	5
\.



-- Data for Name: video_person_rel; Type: TABLE DATA; Schema: public; Owner: colet
--

COPY public.video_person_rel (id, person_id, video_id, character_id) FROM stdin;
1	1	1	\N
2	2	1	\N
\.


--
-- Name: person_id_seq; Type: SEQUENCE SET; Schema: public; Owner: colet
--

SELECT pg_catalog.setval('public.person_id_seq', 4, true);


--
-- Name: character_id_seq; Type: SEQUENCE SET; Schema: public; Owner: colet
--

SELECT pg_catalog.setval('public.character_id_seq', 2, false);


--
-- Name: creator_id_seq; Type: SEQUENCE SET; Schema: public; Owner: colet
--

SELECT pg_catalog.setval('public.creator_id_seq', 13, true);


--
-- Name: video_id_seq; Type: SEQUENCE SET; Schema: public; Owner: colet
--

SELECT pg_catalog.setval('public.video_id_seq', 1, true);


--
-- Name: video_person_rel_id_seq; Type: SEQUENCE SET; Schema: public; Owner: colet
--

SELECT pg_catalog.setval('public.video_person_rel_id_seq', 2, true);
