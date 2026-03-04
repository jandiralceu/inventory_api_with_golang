-- Clean existing data before seeding
TRUNCATE TABLE public.inventory_transactions, public.inventory, public.products, public.suppliers, public.categories, public.users, public.warehouses, public.roles CASCADE;

--
-- PostgreSQL database dump
--

-- Dumped from database version 17.0 (Debian 17.0-1.pgdg120+1)
-- Dumped by pg_dump version 17.0 (Debian 17.0-1.pgdg120+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Data for Name: categories; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.categories (id, name, slug, description, parent_id, is_active, created_at, updated_at) FROM stdin;
d24c15fc-23ca-492b-841d-4e270437a975	Electronics	electronics	Electronic devices and accessories, including computers, smartphones, and peripherals	\N	t	2026-03-04 00:08:36.265622+00	2026-03-04 00:08:36.265622+00
33bdd5f9-ea38-4dc1-874a-b95575b24db8	Office Supplies	office-supplies	Everyday office essentials such as pens, paper, and organizational tools	\N	t	2026-03-04 00:22:23.442222+00	2026-03-04 00:22:23.442222+00
083ad58b-00e1-4392-a5c7-5980e2df4094	Furniture	furniture	Office and warehouse furniture including desks, chairs, and shelving units	\N	t	2026-03-04 00:22:42.978076+00	2026-03-04 00:22:42.978076+00
e510e60d-2e80-479e-bfa1-91b5c193bc6d	Tools & Hardware	tools-hardware	Hand tools, power tools, and hardware supplies for maintenance and operations	\N	t	2026-03-04 00:23:04.509988+00	2026-03-04 00:23:04.509988+00
ccf4dd63-12c8-4f66-b6ff-6e7ddc7a44ce	Cleaning & Hygiene	cleaning-hygiene	Cleaning products, sanitizers, and hygiene supplies for facility maintenance	\N	t	2026-03-04 00:23:26.837366+00	2026-03-04 00:23:26.837366+00
00ceeb46-aa42-4834-b61b-0730470a70c8	Packaging Materials	packaging-materials	Boxes, tape, bubble wrap, and other materials used for shipping and storage	\N	t	2026-03-04 00:23:58.573184+00	2026-03-04 00:23:58.573184+00
51691084-ded1-4a6c-afa2-fa8c7deb34e2	Safety & PPE	safety-ppe	Personal protective equipment including helmets, gloves, and safety vests	\N	t	2026-03-04 00:24:46.446741+00	2026-03-04 00:24:46.446741+00
417c41c0-427b-4275-bcfd-422b423df869	Food & Beverages	food-beverages	Non-perishable food items, beverages, and kitchen supplies	\N	t	2026-03-04 00:25:23.349282+00	2026-03-04 00:25:23.349282+00
46d6195b-36ef-4322-afa6-ab007d660878	Medical Supplies	medical-supplies	First aid kits, medications, and medical equipment for workplace use	\N	t	2026-03-04 00:25:45.899614+00	2026-03-04 00:25:45.899614+00
d912155a-faea-4bcf-a6dc-5003a11039cd	Raw Materials	raw-materials	Basic materials used in manufacturing, such as metals, plastics, and textiles	\N	t	2026-03-04 00:26:07.537337+00	2026-03-04 00:26:07.537337+00
\.


--
-- Data for Name: suppliers; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.suppliers (id, name, slug, description, tax_id, email, phone, address, contact_person, is_active, created_at, updated_at) FROM stdin;
6ce6beb4-c2b7-4e70-9919-c0824d4df4de	TechNova Solutions	technova-solutions	Leading supplier of enterprise electronics and IT infrastructure components.	EIN-12-3456789	contact@technovasolutions.com	+1-408-555-0101	{"city": "San Jose", "state": "CA", "number": "1500", "street": "Innovation Blvd", "country": "US", "zip_code": "95110"}	Zinedine Zidane	t	2026-03-04 01:05:45.892516+00	2026-03-04 01:05:45.892516+00
618e9e9d-5beb-4b41-a54f-10616f87bbdd	GlobalPack Industries	globalpack-industries	Manufacturer and distributor of industrial packaging materials and supplies.	EIN-34-5678901	sales@globalpack.com	+1-312-555-0202	{"city": "Chicago", "state": "IL", "number": "820", "street": "Commerce Drive", "country": "US", "zip_code": "60601"}	Viola Davis	t	2026-03-04 01:05:58.017064+00	2026-03-04 01:05:58.017064+00
23a18f3a-740a-4852-9cd7-7ac888fdc6ed	CleanCore Supply Co.	cleancore-supply-co	Wholesale supplier of professional cleaning and hygiene products.	EIN-56-7890123	orders@cleancoresupply.com	+1-305-555-0303	{"city": "Miami", "state": "FL", "number": "275", "street": "Hygiene Park Ave", "country": "US", "zip_code": "33101"}	Kyrie Irving	t	2026-03-04 01:06:34.171846+00	2026-03-04 01:06:34.171846+00
71ca31a5-e988-4278-b128-296c4bf82abd	SafeGuard PPE Ltd.	safeguard-ppe-ltd	Specialized supplier of personal protective equipment and workplace safety products.	EIN-78-9012345	info@safeguardppe.com	+1-713-555-0404	{"city": "Houston", "state": "TX", "number": "3300", "street": "Safety Lane", "country": "US", "zip_code": "77001"}	Serena Williams	t	2026-03-04 01:06:52.496905+00	2026-03-04 01:06:52.496905+00
22f8706a-3ee7-40d1-a138-c32eb97b4e30	FurniPro Workspace	furnipro-workspace	Commercial furniture supplier specializing in ergonomic office and warehouse solutions.	EIN-90-1234567	workspace@furnipro.com	+1-206-555-0505	{"city": "Seattle", "state": "WA", "number": "14", "street": "Design Quarter", "country": "US", "zip_code": "98101"}	Chadwick Boseman	t	2026-03-04 01:07:03.107276+00	2026-03-04 01:07:03.107276+00
60063305-6837-4c45-9f57-0d2799caa40a	MedFirst Supplies	medfirst-supplies	Trusted provider of first aid kits, medical devices, and workplace healthcare products.	EIN-23-4567890	support@medfirstsupplies.com	+1-617-555-0606	{"city": "Boston", "state": "MA", "number": "88", "street": "Health Center Rd", "country": "US", "zip_code": "02101"}	Lewis Hamilton	t	2026-03-04 01:07:18.728188+00	2026-03-04 01:07:18.728188+00
d98d41a1-2b67-4b9a-a3f6-9f1967e50a8c	RawEdge Materials	rawedge-materials	Bulk supplier of raw industrial materials including metals, plastics, and composites.	EIN-45-6789012	bulk@rawedgematerials.com	+1-412-555-0707	{"city": "Pittsburgh", "state": "PA", "number": "500", "street": "Industrial Parkway", "country": "US", "zip_code": "15201"}	Michael B. Jordan	t	2026-03-04 01:07:30.841726+00	2026-03-04 01:07:30.841726+00
f4b824b5-dd74-4735-92c8-12ab417c3c67	ToolMaster Hardware	toolmaster-hardware	Full-range supplier of hand tools, power tools, and hardware for maintenance operations.	EIN-67-8901234	sales@toolmasterhw.com	+1-602-555-0808	{"city": "Phoenix", "state": "AZ", "number": "940", "street": "Workshop Blvd", "country": "US", "zip_code": "85001"}	Naomi Osaka	t	2026-03-04 01:07:45.220871+00	2026-03-04 01:07:45.220871+00
ec9afadd-d1cb-413b-9531-c17a1b9ef529	NutriStore Wholesale	nutristore-wholesale	Wholesale distributor of non-perishable food items, beverages, and kitchen consumables.	EIN-89-0123456	wholesale@nutristore.com	+1-503-555-0909	{"city": "Portland", "state": "OR", "number": "77", "street": "Food Commerce Ave", "country": "US", "zip_code": "97201"}	Samuel L. Jackson	t	2026-03-04 01:08:07.505708+00	2026-03-04 01:08:07.505708+00
8a39b152-1531-440a-8133-816573637464	OfficeEssentials Corp.	officeessentials-corp	One-stop supplier for office supplies, stationery, and organizational tools.	EIN-01-2345678	orders@officeessentials.com	+1-404-555-1010	{"city": "Atlanta", "state": "GA", "number": "200", "street": "Business Center Dr", "country": "US", "zip_code": "30301"}	Simone Biles	t	2026-03-04 01:08:21.156413+00	2026-03-04 01:08:21.156413+00
\.


--
-- Data for Name: products; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.products (id, sku, slug, name, description, price, cost_price, category_id, supplier_id, is_active, reorder_level, reorder_quantity, weight_kg, images, metadata, created_at, updated_at) FROM stdin;
5bd4f74f-37e9-441d-8aef-af72ab943000	LAP-MBP-14-M3	macbook-pro-14-m3	MacBook Pro 14" M3	Apple M3 chip, 16GB RAM, 512GB SSD.	1999.00	1600.00	d24c15fc-23ca-492b-841d-4e270437a975	6ce6beb4-c2b7-4e70-9919-c0824d4df4de	t	5	10	1.550	\N	\N	2026-03-04 02:04:24.824145+00	2026-03-04 02:04:24.824145+00
64fd4b21-e470-4349-a94e-d44441667096	FUR-CHR-ERG-01	ergonomic-office-chair	Ergonomic Office Chair	High-back mesh chair with lumbar support.	349.99	150.00	083ad58b-00e1-4392-a5c7-5980e2df4094	22f8706a-3ee7-40d1-a138-c32eb97b4e30	t	10	20	18.000	\N	\N	2026-03-04 02:04:46.437006+00	2026-03-04 02:04:46.437006+00
f60af1ad-fc94-4f98-b616-4bbf284f1eab	SAF-HLM-IND-01	industrial-safety-helmet	Industrial Safety Helmet		45.00	15.00	51691084-ded1-4a6c-afa2-fa8c7deb34e2	71ca31a5-e988-4278-b128-296c4bf82abd	t	50	100	0.450	\N	\N	2026-03-04 02:04:58.545242+00	2026-03-04 02:04:58.545242+00
bc8c3c75-2fd0-4754-8d08-30a9f7b5283d	ELC-MON-27-4K	4k-computer-monitor-27	4K Computer Monitor 27"		499.00	320.00	d24c15fc-23ca-492b-841d-4e270437a975	6ce6beb4-c2b7-4e70-9919-c0824d4df4de	t	8	15	\N	\N	\N	2026-03-04 02:05:12.963693+00	2026-03-04 02:05:12.963693+00
868079c2-c462-4385-bd22-ee601b927b5c	MED-FAK-PRM	first-aid-kit-premium	First Aid Kit Premium		89.90	40.00	46d6195b-36ef-4322-afa6-ab007d660878	60063305-6837-4c45-9f57-0d2799caa40a	t	20	50	\N	\N	\N	2026-03-04 02:05:25.504673+00	2026-03-04 02:05:25.504673+00
62b1166a-4068-432c-bc75-0fef66e1015b	TOL-DRL-CRD-18V	cordless-power-drill	Cordless Power Drill		149.00	85.00	e510e60d-2e80-479e-bfa1-91b5c193bc6d	f4b824b5-dd74-4735-92c8-12ab417c3c67	t	15	30	\N	\N	\N	2026-03-04 02:05:38.595697+00	2026-03-04 02:05:38.595697+00
94e7ad88-4094-48eb-895b-0d37ef952baf	FUR-DSK-STD-ADJ	standing-desk-frame	Standing Desk Frame		275.00	180.00	083ad58b-00e1-4392-a5c7-5980e2df4094	22f8706a-3ee7-40d1-a138-c32eb97b4e30	t	5	10	\N	\N	\N	2026-03-04 02:05:50.174283+00	2026-03-04 02:05:50.174283+00
86ca0aaf-d1d7-4002-b084-d7bbc11fcec2	CLN-MSC-GAL	multi-surface-cleaner-gallon	Multi-Surface Cleaner (Gallon)		19.50	7.00	ccf4dd63-12c8-4f66-b6ff-6e7ddc7a44ce	23a18f3a-740a-4852-9cd7-7ac888fdc6ed	t	30	0	\N	\N	\N	2026-03-04 02:06:02.555148+00	2026-03-04 02:06:02.555148+00
8f9fd4c0-6535-4ccb-8202-144cc44283c2	RAW-ALU-SHT-10PR	aluminum-raw-sheets-10pk	Aluminum Raw Sheets (10pk)		450.00	280.00	d912155a-faea-4bcf-a6dc-5003a11039cd	d98d41a1-2b67-4b9a-a3f6-9f1967e50a8c	t	10	0	\N	\N	\N	2026-03-04 02:06:14.290673+00	2026-03-04 02:06:14.290673+00
d9a23a6e-4ba6-45ca-83ca-ac66d7397e91	ELC-KBD-MECH-RGB	mechanical-keyboard-rgb	Mechanical Keyboard RGB		129.00	60.00	d24c15fc-23ca-492b-841d-4e270437a975	6ce6beb4-c2b7-4e70-9919-c0824d4df4de	t	12	0	\N	\N	\N	2026-03-04 02:06:25.520939+00	2026-03-04 02:06:25.520939+00
1fad1bab-60ac-4874-aeb5-85eb542e62e4	OFF-STP-PRO	professional-stapler	Professional Stapler		25.00	10.00	33bdd5f9-ea38-4dc1-874a-b95575b24db8	8a39b152-1531-440a-8133-816573637464	t	25	0	\N	\N	\N	2026-03-04 02:06:37.373559+00	2026-03-04 02:06:37.373559+00
7060dbac-3adb-4a34-9fd5-b18f4cc7148d	PKG-BBW-ROLL-L	shipping-bubble-wrap-roll	Shipping Bubble Wrap (Roll)		35.00	12.00	00ceeb46-aa42-4834-b61b-0730470a70c8	618e9e9d-5beb-4b41-a54f-10616f87bbdd	t	40	0	\N	\N	\N	2026-03-04 02:06:52.04533+00	2026-03-04 02:06:52.04533+00
34d59bf2-03d6-486e-a8a5-8ce8ad75ee39	SAF-GGL-AFG	safety-goggles-anti-fog	Safety Goggles (Anti-fog)		18.00	5.00	51691084-ded1-4a6c-afa2-fa8c7deb34e2	71ca31a5-e988-4278-b128-296c4bf82abd	t	60	0	\N	\N	\N	2026-03-04 02:07:05.818277+00	2026-03-04 02:07:05.818277+00
8c015c9d-0a80-4211-b9ee-692b0b0b7c19	OFF-WBD-120-90	whiteboard-120x90cm	Whiteboard 120x90cm		85.00	45.00	33bdd5f9-ea38-4dc1-874a-b95575b24db8	8a39b152-1531-440a-8133-816573637464	t	5	0	\N	\N	\N	2026-03-04 02:07:17.048017+00	2026-03-04 02:07:17.048017+00
8913a04d-8b1c-4322-bc48-1d48ee64b7f7	ELC-MOU-WRL-G	wireless-gaming-mouse	Wireless Gaming Mouse		79.90	35.00	d24c15fc-23ca-492b-841d-4e270437a975	6ce6beb4-c2b7-4e70-9919-c0824d4df4de	t	20	0	\N	\N	\N	2026-03-04 02:07:33.472044+00	2026-03-04 02:07:33.472044+00
\.


--
-- Data for Name: warehouses; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.warehouses (id, name, slug, code, description, address, is_active, manager_name, phone, email, notes, created_at, updated_at) FROM stdin;
59e2ec0b-3797-4fcd-9d17-ede2e6627b99	North Coast Distribution Center	north-coast-distribution-center	WH-NC-001	Primary logistics hub for northern regional operations.	{"city": "New York", "state": "NY", "number": "500", "street": "Harbor Way", "country": "US", "zip_code": "10001"}	t	Morgan Freeman	+1-212-555-1111	north.coast@inventory.com		2026-03-04 01:16:54.992372+00	2026-03-04 01:16:54.992372+00
078821c6-d938-4618-ba2e-ab88f46d6da8	Silicon Valley Storage	silicon-valley-storage	WH-SV-002	High-security facility for electronic components and tech hardware.	{"city": "Palo Alto", "state": "CA", "number": "101", "street": "Tech Parkway", "country": "US", "zip_code": "94301"}	t	Lupita Nyong'o	+1-650-555-2222	sv.storage@inventory.com		2026-03-04 01:17:22.445461+00	2026-03-04 01:17:22.445461+00
f539d56f-53d4-4ea5-b5b6-476514108b55	Midwest Logistics Hub	midwest-logistics-hub	WH-MW-003	Centralized warehouse for nationwide shipping coordination.	{"city": "Chicago", "state": "IL", "number": "2500", "street": "Central Ave", "country": "US", "zip_code": "60607"}	t	Shaquille O'Neal	+1-312-555-3333	midwest.hub@inventory.com		2026-03-04 01:17:33.137799+00	2026-03-04 01:17:33.137799+00
e7658181-9571-49f4-a02a-2755a4cb1324	Southern Gate Fulfillment	southern-gate-fulfillment	WH-SG-004	Fast-paced fulfillment center for southern retail distribution.	{"city": "Atlanta", "state": "GA", "number": "88", "street": "Peach Tree St", "country": "US", "zip_code": "30303"}	t	Halle Berry	+1-404-555-4444	southern.gate@inventory.com		2026-03-04 01:17:43.500045+00	2026-03-04 01:17:43.500045+00
501c5669-48ae-44d9-a60a-a39e97a0cf8b	Pacific Rim Warehouse	pacific-rim-warehouse	WH-PR-005	Strategic port-side facility for international imports and exports.	{"city": "Los Angeles", "state": "CA", "number": "1200", "street": "Ocean View Dr", "country": "US", "zip_code": "90001"}	t	Idris Elba	+1-213-555-5555	pacific.rim@inventory.com		2026-03-04 01:18:43.092013+00	2026-03-04 01:18:43.092013+00
5f60421d-3cab-4a5f-b29d-a953f566e59a	Mountain View Depot	mountain-view-depot	WH-MV-006	Storage for tools, heavy machinery, and industrial raw materials.	{"city": "Denver", "state": "CO", "number": "45", "street": "Peak Road", "country": "US", "zip_code": "80201"}	t	Kobe Bryant	+1-303-555-6666	mountain.depot@inventory.com		2026-03-04 01:18:53.831829+00	2026-03-04 01:18:53.831829+00
c5cafcff-bffc-4773-ac98-5615e7f61866	Sunshine State Annex	sunshine-state-annex	WH-SS-007	Temperature-controlled facility for medical and food supplies.	{"city": "Miami", "state": "FL", "number": "300", "street": "Palm Blvd", "country": "US", "zip_code": "33101"}	t	Zendaya	+1-305-555-7777	sunshine.annex@inventory.com		2026-03-04 01:19:04.495711+00	2026-03-04 01:19:04.495711+00
cf3ccb80-f1df-4665-999e-18858016b1ac	New England Cold Storage	new-england-cold-storage	WH-NE-008	Specialized refrigeration unit for perishable goods and chemicals.	{"city": "Boston", "state": "MA", "number": "10", "street": "Beacon St", "country": "US", "zip_code": "02108"}	t	Forest Whitaker	+1-617-555-8888	ne.cold@inventory.com		2026-03-04 01:19:24.12041+00	2026-03-04 01:19:24.12041+00
74a99e7c-0df8-4bce-bf4d-9f6f5790166e	Lone Star Reserve	lone-star-reserve	WH-LS-009	Large-scale overflow warehouse for seasonal inventory peaks.	{"city": "Austin", "state": "TX", "number": "777", "street": "Austin Way", "country": "US", "zip_code": "73301"}	t	Usain Bolt	+1-512-555-9999	lonestar.res@inventory.com		2026-03-04 01:19:44.677836+00	2026-03-04 01:19:44.677836+00
b6c1bd6d-b7cd-44f9-9407-8edac8d40701	Liberty Bell Terminal	liberty-bell-terminal	WH-LB-010	East coast terminal for final-mile delivery operations.	{"city": "Philadelphia", "state": "PA", "number": "1776", "street": "Freedom Blvd", "country": "US", "zip_code": "19101"}	t	Angela Bassett	+1-215-555-1010	liberty.bell@inventory.com		2026-03-04 01:20:04.730452+00	2026-03-04 01:20:04.730452+00
\.


--
-- Data for Name: inventory; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.inventory (id, product_id, warehouse_id, quantity, reserved_quantity, location_code, min_quantity, max_quantity, version, last_counted_at, metadata, created_at, updated_at) FROM stdin;
9f684b87-e72a-465a-81fe-3d6c223d1c55	8913a04d-8b1c-4322-bc48-1d48ee64b7f7	078821c6-d938-4618-ba2e-ab88f46d6da8	50	0	SEC-A1-L4	5	100	1	\N	\N	2026-03-04 02:21:08.117917+00	2026-03-04 02:21:08.117917+00
2937d633-92d6-4f93-8b6c-33ab6b4526ef	94e7ad88-4094-48eb-895b-0d37ef952baf	f539d56f-53d4-4ea5-b5b6-476514108b55	120	0	BULK-04	20	200	1	\N	\N	2026-03-04 02:21:40.582711+00	2026-03-04 02:21:40.582711+00
2106c55c-d40e-437f-ba7f-92edd48f9e3a	34d59bf2-03d6-486e-a8a5-8ce8ad75ee39	59e2ec0b-3797-4fcd-9d17-ede2e6627b99	300	0	SAF-01	50	500	1	\N	\N	2026-03-04 02:22:00.092277+00	2026-03-04 02:22:00.092277+00
34cbc8e7-09b0-43ec-b651-ecfff7cf94c1	62b1166a-4068-432c-bc75-0fef66e1015b	5f60421d-3cab-4a5f-b29d-a953f566e59a	45	0	TOOL-B2	10	100	1	\N	\N	2026-03-04 02:22:14.892435+00	2026-03-04 02:22:14.892435+00
16e77907-701f-471c-8a00-b568f2c22013	d9a23a6e-4ba6-45ca-83ca-ac66d7397e91	078821c6-d938-4618-ba2e-ab88f46d6da8	30	0	ELC-C3	5	50	1	\N	\N	2026-03-04 02:22:27.531038+00	2026-03-04 02:22:27.531038+00
00e3fb0c-7d81-48f2-8bce-44913e92cd20	8f9fd4c0-6535-4ccb-8202-144cc44283c2	501c5669-48ae-44d9-a60a-a39e97a0cf8b	15	0	EXT-YARD-01	2	20	1	\N	\N	2026-03-04 02:22:38.056118+00	2026-03-04 02:22:38.056118+00
79490c89-4e03-4ec8-924f-0af9133b48f7	86ca0aaf-d1d7-4002-b084-d7bbc11fcec2	c5cafcff-bffc-4773-ac98-5615e7f61866	250	0	HAZ-M1	20	500	1	\N	\N	2026-03-04 02:22:51.379762+00	2026-03-04 02:22:51.379762+00
7726df26-8173-474c-b768-2f4a0166fd18	7060dbac-3adb-4a34-9fd5-b18f4cc7148d	e7658181-9571-49f4-a02a-2755a4cb1324	80	0	PKG-09	10	150	1	\N	\N	2026-03-04 02:23:08.680552+00	2026-03-04 02:23:08.680552+00
fcc35b71-d54d-48eb-ad32-2c8eaa90feda	1fad1bab-60ac-4874-aeb5-85eb542e62e4	b6c1bd6d-b7cd-44f9-9407-8edac8d40701	60	0	STN-C4	10	100	1	\N	\N	2026-03-04 02:23:20.81154+00	2026-03-04 02:23:20.81154+00
b0246cd0-7729-4617-b257-2eb3bd47f55a	8c015c9d-0a80-4211-b9ee-692b0b0b7c19	f539d56f-53d4-4ea5-b5b6-476514108b55	15	0	LRG-02	2	25	1	\N	\N	2026-03-04 02:23:31.438134+00	2026-03-04 02:23:31.438134+00
5e1b898c-fce7-41fe-a343-914f8db905f6	8913a04d-8b1c-4322-bc48-1d48ee64b7f7	59e2ec0b-3797-4fcd-9d17-ede2e6627b99	95	0	PICK-FAST-01	20	200	1	\N	\N	2026-03-04 02:23:43.508345+00	2026-03-04 02:23:43.508345+00
5684b141-43dc-497e-8bbd-8a12190f9d24	d9a23a6e-4ba6-45ca-83ca-ac66d7397e91	74a99e7c-0df8-4bce-bf4d-9f6f5790166e	200	0	OVER-B1	50	500	1	\N	\N	2026-03-04 02:23:56.165518+00	2026-03-04 02:23:56.165518+00
3e95c380-2748-4e0e-9704-14b75c857f0f	8c015c9d-0a80-4211-b9ee-692b0b0b7c19	cf3ccb80-f1df-4665-999e-18858016b1ac	40	0	COLD-05	10	100	1	\N	\N	2026-03-04 02:24:41.636788+00	2026-03-04 02:24:41.636788+00
a9c2879f-1c64-4370-970c-1072c917964b	94e7ad88-4094-48eb-895b-0d37ef952baf	e7658181-9571-49f4-a02a-2755a4cb1324	15	0	BULK-S2	5	30	1	\N	\N	2026-03-04 02:24:52.011946+00	2026-03-04 02:24:52.011946+00
2e9df62e-c2b2-4da7-8191-d5af30368fcc	8f9fd4c0-6535-4ccb-8202-144cc44283c2	5f60421d-3cab-4a5f-b29d-a953f566e59a	8	0	YARD-B	2	10	1	\N	\N	2026-03-04 02:25:04.027051+00	2026-03-04 02:25:04.027051+00
\.


--
-- Data for Name: roles; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.roles (id, name, description, created_at) FROM stdin;
a464fff6-4b82-45d6-ac18-2d08aeda4107	admin	Full access: manages users, products, and system configurations	2026-03-03 18:17:43.881514+00
9277e314-02da-4f0a-ad37-7af252175a38	manager	Manages products, suppliers, and views reports	2026-03-03 18:17:43.89104+00
8349d164-324a-4280-86cd-5fecbaf4c51a	operator	Daily operations: records inventory movements and checks stock	2026-03-03 18:17:43.893896+00
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.users (id, name, email, password_hash, role_id, created_at, updated_at) FROM stdin;
2dbdf352-7008-48a6-9637-62149a56126f	Jandir Alceu	jandiralceu@gmail.com	$argon2id$v=19$m=65536,t=1,p=4$wrfKGIDX9Uhqzm2qUT94kg$quyvfIeO9iR8xvU+B0nLO1PMFiwTU5WA9ivo31xpMSQ	a464fff6-4b82-45d6-ac18-2d08aeda4107	2026-03-03 18:19:22.057947+00	2026-03-03 18:19:22.057947+00
73b63226-4f0e-4cc4-bfaa-310928d89528	Michae Jordan	michaeljordan@email.com	$argon2id$v=19$m=65536,t=1,p=4$XXEB/SdNGEuXZBHd9xGKCA$KI8xsCq4yyPO80xUVhRoRcyjEgASvHa+aKwGyoNuSaw	a464fff6-4b82-45d6-ac18-2d08aeda4107	2026-03-03 23:55:50.334727+00	2026-03-03 23:55:50.334727+00
6bd4aefe-c013-43bb-9f0a-565026581a22	Dennis Rodman	dennisrodman@email.com	$argon2id$v=19$m=65536,t=1,p=4$XptntHT22zlj5adekcFspg$SfZH4f4/cfpn6WWOXLOjshRygB2xfF2Sa4Sz84wko7Y	a464fff6-4b82-45d6-ac18-2d08aeda4107	2026-03-03 23:56:12.398703+00	2026-03-03 23:56:12.398703+00
1e57c2e3-325e-4734-9cf0-e28555cb0f6c	Russell Westbrook	russellwestbrook@email.com	$argon2id$v=19$m=65536,t=1,p=4$AUiSvZptIemKKyTif9D6Iw$keha/ywmh7qc6dpAkqK4UlI3BpaBtuPDMRjupFT6Dbs	9277e314-02da-4f0a-ad37-7af252175a38	2026-03-03 23:57:18.302889+00	2026-03-03 23:57:18.302889+00
c29bc2d1-88e6-43a4-9781-3764624bf9cc	Lebron James	lebronjames@email.com	$argon2id$v=19$m=65536,t=1,p=4$anSkGKDm+fvgGEejrRlwaA$bkL+ALtsFWfY1o2/92zwDGqtUYKQcJ0NEL7gxqggDS0	9277e314-02da-4f0a-ad37-7af252175a38	2026-03-03 23:57:40.191921+00	2026-03-03 23:57:40.191921+00
6357c4c3-de7e-466b-9e8c-7cd207b7a47b	Denzel Washington	dw@email.com	$argon2id$v=19$m=65536,t=1,p=4$JFUlGi9PeU2mPKoiOAEoKg$HTwBLMApLTrKMw8J2GueeJmHKtaXulA75KmAlQMynG8	9277e314-02da-4f0a-ad37-7af252175a38	2026-03-04 00:00:54.069621+00	2026-03-04 00:00:54.069621+00
8b6d5b64-5207-4c98-a2c7-76067197f8c1	Will Smith	wsmith@email.com	$argon2id$v=19$m=65536,t=1,p=4$EIF1IyHqL2D7LZFWxnSnKA$wj1S9dD2cSLeM7pP2PGR/I01CZfG280Fa1hOD2iMCgs	9277e314-02da-4f0a-ad37-7af252175a38	2026-03-04 00:02:25.996711+00	2026-03-04 00:02:25.996711+00
34842db8-5d68-4944-ae73-b69d93e4cb1d	Nina Simone	ninasimone@email.com	$argon2id$v=19$m=65536,t=1,p=4$Lvt/zAFBI2ww5dMMuxnSWA$2brdqlFSoZ3a7qfNhKF69lZg6blLkpe9E/7IGXOeZF4	9277e314-02da-4f0a-ad37-7af252175a38	2026-03-04 00:02:47.703193+00	2026-03-04 00:02:47.703193+00
7d5a87d8-a895-4695-bb7b-5e2288e1f5c5	Jermaine Cole	jcole@email.com	$argon2id$v=19$m=65536,t=1,p=4$a9/4x9rDXBEZSGJGoeTYrg$llxCH76x5Rfa+MJKOyQllpS3ZDieiaTT9+D9mZTXUmU	9277e314-02da-4f0a-ad37-7af252175a38	2026-03-04 00:04:04.564207+00	2026-03-04 00:04:04.564207+00
\.


--
-- Data for Name: inventory_transactions; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.inventory_transactions (id, inventory_id, product_id, warehouse_id, user_id, quantity_change, quantity_balance, transaction_type, reference_id, reason, created_at) FROM stdin;
2e56d45c-7d03-46c7-8036-71115232cadd	9f684b87-e72a-465a-81fe-3d6c223d1c55	8913a04d-8b1c-4322-bc48-1d48ee64b7f7	078821c6-d938-4618-ba2e-ab88f46d6da8	\N	50	100	OPENING		Initial stock setup	2026-03-04 02:21:08.145497+00
ec5c520a-e785-42ca-9173-832aa8f09732	2937d633-92d6-4f93-8b6c-33ab6b4526ef	94e7ad88-4094-48eb-895b-0d37ef952baf	f539d56f-53d4-4ea5-b5b6-476514108b55	\N	120	240	OPENING		Initial stock setup	2026-03-04 02:21:40.585879+00
c4a8982f-f814-4b35-a0b1-1bc671c032bf	2106c55c-d40e-437f-ba7f-92edd48f9e3a	34d59bf2-03d6-486e-a8a5-8ce8ad75ee39	59e2ec0b-3797-4fcd-9d17-ede2e6627b99	\N	300	600	OPENING		Initial stock setup	2026-03-04 02:22:00.094016+00
19ed8103-146c-4b11-a96e-5d77407fd946	34cbc8e7-09b0-43ec-b651-ecfff7cf94c1	62b1166a-4068-432c-bc75-0fef66e1015b	5f60421d-3cab-4a5f-b29d-a953f566e59a	\N	45	90	OPENING		Initial stock setup	2026-03-04 02:22:14.894449+00
520fbf16-0448-48ef-9065-7d71daac7ead	16e77907-701f-471c-8a00-b568f2c22013	d9a23a6e-4ba6-45ca-83ca-ac66d7397e91	078821c6-d938-4618-ba2e-ab88f46d6da8	\N	30	60	OPENING		Initial stock setup	2026-03-04 02:22:27.533511+00
50ec34fd-4787-424d-a0c3-03f3dfa77a7c	00e3fb0c-7d81-48f2-8bce-44913e92cd20	8f9fd4c0-6535-4ccb-8202-144cc44283c2	501c5669-48ae-44d9-a60a-a39e97a0cf8b	\N	15	30	OPENING		Initial stock setup	2026-03-04 02:22:38.058867+00
27acc64f-26dd-47d1-9c9a-f5a135d68872	79490c89-4e03-4ec8-924f-0af9133b48f7	86ca0aaf-d1d7-4002-b084-d7bbc11fcec2	c5cafcff-bffc-4773-ac98-5615e7f61866	\N	250	500	OPENING		Initial stock setup	2026-03-04 02:22:51.381895+00
41e6efb4-cbd1-4c84-8d9c-07b99f4a14b1	7726df26-8173-474c-b768-2f4a0166fd18	7060dbac-3adb-4a34-9fd5-b18f4cc7148d	e7658181-9571-49f4-a02a-2755a4cb1324	\N	80	160	OPENING		Initial stock setup	2026-03-04 02:23:08.683921+00
1206d73a-29b9-4163-ba84-f2477c118f70	fcc35b71-d54d-48eb-ad32-2c8eaa90feda	1fad1bab-60ac-4874-aeb5-85eb542e62e4	b6c1bd6d-b7cd-44f9-9407-8edac8d40701	\N	60	120	OPENING		Initial stock setup	2026-03-04 02:23:20.81699+00
00208bf7-6a45-40de-bfcd-fb920489c154	b0246cd0-7729-4617-b257-2eb3bd47f55a	8c015c9d-0a80-4211-b9ee-692b0b0b7c19	f539d56f-53d4-4ea5-b5b6-476514108b55	\N	15	30	OPENING		Initial stock setup	2026-03-04 02:23:31.440177+00
1c8765c9-b12f-4b7b-9734-80e212b86c82	5e1b898c-fce7-41fe-a343-914f8db905f6	8913a04d-8b1c-4322-bc48-1d48ee64b7f7	59e2ec0b-3797-4fcd-9d17-ede2e6627b99	\N	95	190	OPENING		Initial stock setup	2026-03-04 02:23:43.510267+00
83eb33e1-9601-479e-9e8e-1f0f1f8d907f	5684b141-43dc-497e-8bbd-8a12190f9d24	d9a23a6e-4ba6-45ca-83ca-ac66d7397e91	74a99e7c-0df8-4bce-bf4d-9f6f5790166e	\N	200	400	OPENING		Initial stock setup	2026-03-04 02:23:56.171921+00
6e203bb1-3817-4867-a06d-75c243741c21	3e95c380-2748-4e0e-9704-14b75c857f0f	8c015c9d-0a80-4211-b9ee-692b0b0b7c19	cf3ccb80-f1df-4665-999e-18858016b1ac	\N	40	80	OPENING		Initial stock setup	2026-03-04 02:24:41.6402+00
9aa8cd75-67c9-4f58-ae0c-9f46d3192623	a9c2879f-1c64-4370-970c-1072c917964b	94e7ad88-4094-48eb-895b-0d37ef952baf	e7658181-9571-49f4-a02a-2755a4cb1324	\N	15	30	OPENING		Initial stock setup	2026-03-04 02:24:52.014829+00
836333a6-baa7-41e4-8a0e-291561e25b1e	2e9df62e-c2b2-4da7-8191-d5af30368fcc	8f9fd4c0-6535-4ccb-8202-144cc44283c2	5f60421d-3cab-4a5f-b29d-a953f566e59a	\N	8	16	OPENING		Initial stock setup	2026-03-04 02:25:04.034159+00
\.

--
-- PostgreSQL database dump complete
--

