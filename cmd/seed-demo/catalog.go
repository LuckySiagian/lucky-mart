package main

// Demo catalogue. Prices are typical toko kelontong prices for October 2026,
// based on national food price data (PIHPS/Siskaperbapo) and Alfamart /
// Indomaret shelf prices; adjust them to the store's real prices.

type unitSpec struct {
	Name       string
	Conversion int // how many base units one of this unit holds
	Retail     float64
	Wholesale  float64 // B2B price at MinWholesale or more
	MinWhole   int
	Stock      int
}

type productSpec struct {
	SKU         string
	Name        string
	Description string
	MinStock    int
	Units       []unitSpec
}

type categorySpec struct {
	Name        string
	Description string
	Products    []productSpec
}

var catalog = []categorySpec{
	{"Beras", "Beras premium dan medium, kemasan 5 kg dan 25 kg", []productSpec{
		{"LM-BRS-001", "Beras Premium Pandan Wangi 5 kg", "Beras premium pulen dan wangi, kemasan karung 5 kg.", 10, []unitSpec{
			{"Karung 5 kg", 1, 89000, 86000, 5, 60}}},
		{"LM-BRS-002", "Beras Medium 5 kg", "Beras medium untuk kebutuhan sehari-hari, kemasan 5 kg.", 10, []unitSpec{
			{"Karung 5 kg", 1, 79000, 76500, 5, 70}}},
		{"LM-BRS-003", "Beras Premium 25 kg", "Beras premium karung besar, cocok untuk warung makan.", 3, []unitSpec{
			{"Karung 25 kg", 1, 425000, 412000, 2, 15}}},
	}},
	{"Minyak Goreng", "Minyak goreng kemasan dan curah", []productSpec{
		{"LM-MYK-001", "Minyak Goreng Bimoli 2 L", "Minyak goreng sawit kemasan pouch 2 liter.", 12, []unitSpec{
			{"Pouch 2 L", 1, 42500, 41000, 6, 48},
			{"Dus (6 pouch)", 6, 249000, 243000, 2, 8}}},
		{"LM-MYK-002", "Minyak Goreng Sunco 2 L", "Minyak goreng kemasan pouch 2 liter.", 12, []unitSpec{
			{"Pouch 2 L", 1, 46500, 45000, 6, 36}}},
		{"LM-MYK-003", "Minyak Goreng Fortune 2 L", "Minyak goreng kemasan pouch 2 liter.", 12, []unitSpec{
			{"Pouch 2 L", 1, 40500, 39000, 6, 40}}},
		{"LM-MYK-004", "MinyaKita 1 L", "Minyak goreng rakyat kemasan 1 liter.", 24, []unitSpec{
			{"Pouch 1 L", 1, 17500, 16800, 12, 96},
			{"Dus (12 pouch)", 12, 204000, 199000, 2, 6}}},
		{"LM-MYK-005", "Minyak Goreng Curah 1 kg", "Minyak goreng curah dikemas plastik 1 kg.", 10, []unitSpec{
			{"Plastik 1 kg", 1, 21500, 20800, 10, 40}}},
	}},
	{"Gula & Pemanis", "Gula pasir dan gula merah", []productSpec{
		{"LM-GLA-001", "Gula Pasir Gulaku 1 kg", "Gula kristal putih kemasan 1 kg.", 20, []unitSpec{
			{"Pak 1 kg", 1, 20000, 19300, 10, 120}}},
		{"LM-GLA-002", "Gula Pasir Curah 1 kg", "Gula pasir lokal dikemas plastik 1 kg.", 20, []unitSpec{
			{"Plastik 1 kg", 1, 19500, 18800, 10, 80}}},
		{"LM-GLA-003", "Gula Merah Kelapa 500 g", "Gula merah kelapa cetak.", 5, []unitSpec{
			{"Bungkus 500 g", 1, 12000, 11300, 10, 30}}},
	}},
	{"Telur", "Telur ayam ras segar", []productSpec{
		{"LM-TLR-001", "Telur Ayam Ras", "Telur ayam negeri segar, dijual per kg atau per tray isi 30 butir.", 10, []unitSpec{
			{"Kg", 1, 29500, 28500, 5, 60},
			{"Tray (30 butir)", 1, 56000, 54000, 3, 20}}},
	}},
	{"Tepung", "Tepung terigu, beras, dan tapioka", []productSpec{
		{"LM-TPG-001", "Tepung Terigu Segitiga Biru 1 kg", "Tepung terigu protein sedang serbaguna.", 10, []unitSpec{
			{"Pak 1 kg", 1, 14500, 13800, 10, 60}}},
		{"LM-TPG-002", "Tepung Terigu Cakra Kembar 1 kg", "Tepung terigu protein tinggi untuk roti dan mie.", 5, []unitSpec{
			{"Pak 1 kg", 1, 15500, 14800, 10, 30}}},
		{"LM-TPG-003", "Tepung Beras Rose Brand 500 g", "Tepung beras untuk kue dan gorengan.", 5, []unitSpec{
			{"Pak 500 g", 1, 8000, 7500, 10, 40}}},
		{"LM-TPG-004", "Tepung Tapioka Rose Brand 500 g", "Tepung tapioka (kanji).", 5, []unitSpec{
			{"Pak 500 g", 1, 7500, 7000, 10, 40}}},
	}},
	{"Mie Instan", "Mie instan satuan dan dus", []productSpec{
		{"LM-MIE-001", "Indomie Goreng 85 g", "Mi instan goreng rasa original.", 40, []unitSpec{
			{"Pcs", 1, 3500, 3200, 10, 300},
			{"Dus (40 pcs)", 40, 112000, 109000, 2, 15}}},
		{"LM-MIE-002", "Indomie Kuah Soto Mie 70 g", "Mi instan kuah rasa soto mie.", 40, []unitSpec{
			{"Pcs", 1, 3300, 3000, 10, 200},
			{"Dus (40 pcs)", 40, 108000, 105000, 2, 10}}},
		{"LM-MIE-003", "Mie Sedaap Goreng 90 g", "Mi instan goreng.", 40, []unitSpec{
			{"Pcs", 1, 3500, 3200, 10, 200},
			{"Dus (40 pcs)", 40, 110000, 107000, 2, 10}}},
	}},
	{"Bumbu & Bahan Dapur", "Bumbu, kecap, saus, dan bahan dapur segar", []productSpec{
		{"LM-BMB-001", "Garam Refina 250 g", "Garam meja beryodium.", 10, []unitSpec{
			{"Pak 250 g", 1, 3500, 3200, 12, 80}}},
		{"LM-BMB-002", "Kecap Manis Bango 520 ml", "Kecap manis kemasan isi ulang.", 10, []unitSpec{
			{"Pouch 520 ml", 1, 24500, 23500, 6, 40}}},
		{"LM-BMB-003", "Saus Sambal ABC 335 ml", "Saus sambal botol.", 6, []unitSpec{
			{"Botol 335 ml", 1, 13000, 12300, 6, 30}}},
		{"LM-BMB-004", "Masako Rasa Ayam 100 g", "Bumbu penyedap rasa ayam.", 10, []unitSpec{
			{"Bungkus 100 g", 1, 6500, 6100, 10, 50}}},
		{"LM-BMB-005", "Bawang Merah", "Bawang merah lokal segar.", 5, []unitSpec{
			{"Kg", 1, 42000, 40500, 5, 25},
			{"1/4 kg", 1, 11000, 11000, 1, 20}}},
		{"LM-BMB-006", "Bawang Putih Kating", "Bawang putih kating segar.", 5, []unitSpec{
			{"Kg", 1, 40000, 38500, 5, 20},
			{"1/4 kg", 1, 10500, 10500, 1, 20}}},
		{"LM-BMB-007", "Cabai Merah Keriting", "Cabai merah keriting segar. Harga mengikuti pasar.", 5, []unitSpec{
			{"Kg", 1, 82000, 80000, 3, 4},
			{"1/4 kg", 1, 21000, 21000, 1, 10}}},
	}},
	{"Kopi, Teh & Susu", "Kopi, teh, dan susu", []productSpec{
		{"LM-KTS-001", "Kopi Kapal Api Special Mix 25 g", "Kopi susu instan sachet.", 5, []unitSpec{
			{"Sachet", 1, 1700, 1550, 20, 200},
			{"Renceng (10 sachet)", 10, 15500, 15000, 5, 30}}},
		{"LM-KTS-002", "Kopi Bubuk Kapal Api Special 165 g", "Kopi bubuk murni.", 5, []unitSpec{
			{"Bungkus 165 g", 1, 15500, 14800, 6, 30}}},
		{"LM-KTS-003", "Teh Celup Sariwangi isi 25", "Teh celup asli.", 5, []unitSpec{
			{"Kotak", 1, 7500, 7100, 6, 40}}},
		{"LM-KTS-004", "Susu Kental Manis Frisian Flag 370 g", "Susu kental manis putih kaleng.", 6, []unitSpec{
			{"Kaleng 370 g", 1, 13500, 12900, 6, 36}}},
		{"LM-KTS-005", "Susu Bubuk Dancow FortiGro Instant 800 g", "Susu bubuk instan untuk keluarga.", 3, []unitSpec{
			{"Pak 800 g", 1, 112000, 108000, 3, 2}}},
	}},
	{"Minuman & Air", "Air mineral dan minuman kemasan", []productSpec{
		{"LM-MNM-001", "Aqua Air Mineral 600 ml", "Air mineral botol 600 ml.", 24, []unitSpec{
			{"Botol", 1, 4000, 3600, 12, 120},
			{"Dus (24 botol)", 24, 58000, 56000, 2, 15}}},
		{"LM-MNM-002", "Aqua Galon 19 L (Isi Ulang)", "Isi ulang galon Aqua; tukar galon kosong.", 5, []unitSpec{
			{"Galon", 1, 22000, 21000, 5, 30}}},
		{"LM-MNM-003", "Teh Pucuk Harum 350 ml", "Teh melati dalam botol.", 12, []unitSpec{
			{"Botol", 1, 4000, 3600, 12, 96}}},
	}},
	{"Gas", "Gas elpiji non-subsidi", []productSpec{
		{"LM-GAS-001", "Bright Gas 5,5 kg (Isi Ulang)", "Isi ulang tabung Bright Gas 5,5 kg; tukar tabung kosong.", 5, []unitSpec{
			{"Tabung", 1, 98000, 96000, 2, 3}}},
	}},
	{"Kebutuhan Rumah", "Sabun, deterjen, dan perlengkapan mandi", []productSpec{
		{"LM-RMH-001", "Sabun Mandi Lifebuoy 110 g", "Sabun batang antibakteri.", 10, []unitSpec{
			{"Batang", 1, 5000, 4600, 12, 60}}},
		{"LM-RMH-002", "Deterjen Rinso Anti Noda 770 g", "Deterjen bubuk.", 6, []unitSpec{
			{"Bungkus 770 g", 1, 25500, 24500, 6, 30}}},
		{"LM-RMH-003", "Sunlight Jeruk Nipis 755 ml", "Sabun cuci piring isi ulang.", 6, []unitSpec{
			{"Pouch 755 ml", 1, 17500, 16800, 6, 30}}},
		{"LM-RMH-004", "Pasta Gigi Pepsodent 190 g", "Pasta gigi pencegah gigi berlubang.", 6, []unitSpec{
			{"Tube 190 g", 1, 14500, 13800, 6, 30}}},
	}},
}
