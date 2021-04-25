package shipping

import (
	"math/rand"

	"github.com/hojulian/mdb-bench/shipping/cargo"
	"github.com/hojulian/mdb-bench/shipping/location"
)

var (
	locations = []location.UNLocode{
		location.SESTO,
		location.AUMEL,
		location.CNHKG,
		location.USNYC,
		location.USCHI,
		location.JNTKO,
		location.DEHAM,
		location.NLRTM,
		location.FIHEL,
	}

	trackingIDs = []cargo.TrackingID{
		"000003A6",
		"00002627",
		"0000D969",
		"0001890B",
		"0002405A",
		"00028D5C",
		"000522EC",
		"0005A65C",
		"00066356",
		"00076DA3",
		"00078394",
		"0007A3AC",
		"0008E5DB",
		"00099445",
		"0009A2B3",
		"0009D13E",
		"000A4959",
		"000AD24C",
		"000B4674",
		"000BA3D9",
		"000BCB7C",
		"000C0656",
		"000C6CD3",
		"000CE5A9",
		"000DDD9E",
		"000E272F",
		"000ED99E",
		"000F3E6D",
		"000F599D",
		"00118AFE",
		"00123EB9",
		"0012C012",
		"001344A6",
		"0013AA93",
		"0013B6B0",
		"00141C3E",
		"001424D4",
		"00142A31",
		"001459C2",
		"00148A5E",
		"00151D36",
		"00152484",
		"0016FC21",
		"00178AD2",
		"0017C91D",
		"0017E3B6",
		"0017E96F",
		"00185A17",
		"001C6C93",
		"001D5EC5",
		"001F0D40",
		"001F1583",
		"001F68C6",
		"001FC145",
		"001FD8EF",
		"00202211",
		"00208FE2",
		"0020FBE1",
		"002149A4",
		"00214B80",
		"0021B5D2",
		"0021BECC",
		"00237A78",
		"002380DF",
		"0023D67E",
		"0024DDD9",
		"00256450",
		"0025843E",
		"0025AB07",
		"002693F1",
		"0026C25C",
		"00275CBB",
		"0027834D",
		"00286388",
		"00289BD9",
		"0028EC71",
		"002A4DAA",
		"002C0EA7",
		"002C63B0",
		"002D5392",
		"002E50E1",
		"002F10D5",
		"002FBDD1",
		"00303EE3",
		"00311AC8",
		"0031773B",
		"00319F09",
		"0031B172",
		"0031BE71",
		"0032D889",
		"0033AD84",
		"0035CD41",
		"00369BD2",
		"0036A20D",
		"0036D987",
		"00375E2E",
		"0037F704",
		"003931D0",
		"00394586",
		"003945E9",
		"00398F6E",
		"003A647E",
		"003A67D4",
		"003AA8DC",
		"003B4EDA",
		"003BBCA0",
		"003BBD4E",
		"003CBC58",
		"003CFE97",
		"003D0146",
		"003D6A7D",
		"003E8CD4",
		"00405687",
		"0040ED24",
		"00419D6C",
		"0041CB46",
		"00424E60",
		"00428862",
		"00430623",
		"004494B3",
		"0045BC44",
		"0046ACD6",
		"00474200",
		"00479EF1",
		"0047EF83",
		"0047FBFA",
		"0048482D",
		"00484B0F",
		"0049D918",
		"004AF2FE",
		"004B2980",
		"004B9303",
		"004BDE74",
		"004C73ED",
		"004DE992",
		"004E4FD5",
		"004E5AD3",
		"004F23A9",
		"0050BD51",
		"00514EE2",
		"00539EC9",
		"0053CB78",
		"00540E35",
		"0054CAC4",
		"005638D6",
		"005737E2",
		"005749FC",
		"005934EB",
		"005A7D1E",
		"005BCC4E",
		"005BE552",
		"005DECD4",
		"005F4970",
		"005FCE40",
		"00617D88",
		"0061E484",
		"0062DD1C",
		"0063A12F",
		"0065027A",
		"0065C968",
		"00668709",
		"0066EADA",
		"006756A0",
		"00684170",
		"0068BD7D",
		"006A2A3E",
		"006A8FC3",
		"006AB35B",
		"006BB35C",
		"006C457E",
		"006CFBF9",
		"006D2477",
		"006DD718",
		"006E348C",
		"006E861E",
		"00723A02",
		"0073F480",
		"0074C475",
		"00758039",
		"00763B00",
		"00774B2A",
		"0078D332",
		"0078FC26",
		"007ADB98",
		"007BBF5A",
		"007CD8FE",
		"007CE478",
		"007D1D4C",
		"007DBF0C",
		"007EEFB3",
		"007EFEE2",
		"007F5183",
		"007FC94C",
		"007FE61D",
		"00824999",
		"0082C16A",
		"0082E45B",
		"0083100A",
		"0083227A",
		"0084054D",
	}
)

func randLocation() location.UNLocode {
	return locations[rand.Intn(len(locations))]
}

func randTrackingID() cargo.TrackingID {
	return trackingIDs[rand.Intn(len(trackingIDs))]
}
