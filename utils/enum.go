package utils

import (
	"database/sql/driver"
	"encoding/binary"
	"fmt"
	"math/rand"
)

type Country int64

// Taken from https://en.wikipedia.org/wiki/List_of_ISO_3166_country_codes
const (
	COUNTRY_AF = iota
	COUNTRY_AX
	COUNTRY_AL
	COUNTRY_DZ
	COUNTRY_AS
	COUNTRY_AD
	COUNTRY_AO
	COUNTRY_AI
	COUNTRY_ATA
	COUNTRY_AG
	COUNTRY_AR
	COUNTRY_AM
	COUNTRY_AW
	COUNTRY_AUS
	COUNTRY_AT
	COUNTRY_AZ
	COUNTRY_BS
	COUNTRY_BH
	COUNTRY_BD
	COUNTRY_BB
	COUNTRY_BY
	COUNTRY_BE
	COUNTRY_BZ
	COUNTRY_BJ
	COUNTRY_BM
	COUNTRY_BT
	COUNTRY_BO
	COUNTRY_BQ
	COUNTRY_BA
	COUNTRY_BW
	COUNTRY_BV
	COUNTRY_BR
	COUNTRY_IO
	COUNTRY_BRN
	COUNTRY_BG
	COUNTRY_BF
	COUNTRY_BI
	COUNTRY_CV
	COUNTRY_KH
	COUNTRY_CM
	COUNTRY_CA
	COUNTRY_KY
	COUNTRY_CF
	COUNTRY_TD
	COUNTRY_CL
	COUNTRY_CN
	COUNTRY_CX
	COUNTRY_CC
	COUNTRY_CO
	COUNTRY_KM
	COUNTRY_CD
	COUNTRY_CG
	COUNTRY_CK
	COUNTRY_CR
	COUNTRY_CI
	COUNTRY_HR
	COUNTRY_CU
	COUNTRY_CW
	COUNTRY_CY
	COUNTRY_CZ
	COUNTRY_DK
	COUNTRY_DJ
	COUNTRY_DM
	COUNTRY_DO
	COUNTRY_EC
	COUNTRY_EG
	COUNTRY_SV
	COUNTRY_GQ
	COUNTRY_ER
	COUNTRY_EE
	COUNTRY_SZ
	COUNTRY_ET
	COUNTRY_FLK
	COUNTRY_FO
	COUNTRY_FJ
	COUNTRY_FI
	COUNTRY_FR
	COUNTRY_GF
	COUNTRY_PF
	COUNTRY_ATF
	COUNTRY_GA
	COUNTRY_GM
	COUNTRY_GE
	COUNTRY_DE
	COUNTRY_GH
	COUNTRY_GI
	COUNTRY_GR
	COUNTRY_GL
	COUNTRY_GD
	COUNTRY_GP
	COUNTRY_GU
	COUNTRY_GT
	COUNTRY_GG
	COUNTRY_GN
	COUNTRY_GW
	COUNTRY_GY
	COUNTRY_HT
	COUNTRY_HM
	COUNTRY_HN
	COUNTRY_HK
	COUNTRY_HU
	COUNTRY_IS
	COUNTRY_IN
	COUNTRY_ID
	COUNTRY_IR
	COUNTRY_IQ
	COUNTRY_IE
	COUNTRY_IM
	COUNTRY_IL
	COUNTRY_IT
	COUNTRY_JM
	COUNTRY_JP
	COUNTRY_JE
	COUNTRY_JO
	COUNTRY_KZ
	COUNTRY_KE
	COUNTRY_KI
	COUNTRY_KP
	COUNTRY_KR
	COUNTRY_KW
	COUNTRY_KG
	COUNTRY_LA
	COUNTRY_LV
	COUNTRY_LB
	COUNTRY_LS
	COUNTRY_LR
	COUNTRY_LY
	COUNTRY_LI
	COUNTRY_LT
	COUNTRY_LU
	COUNTRY_MAC
	COUNTRY_MG
	COUNTRY_MW
	COUNTRY_MY
	COUNTRY_MV
	COUNTRY_ML
	COUNTRY_MT
	COUNTRY_MH
	COUNTRY_MQ
	COUNTRY_MR
	COUNTRY_MU
	COUNTRY_YT
	COUNTRY_MX
	COUNTRY_FM
	COUNTRY_MD
	COUNTRY_MC
	COUNTRY_MN
	COUNTRY_ME
	COUNTRY_MS
	COUNTRY_MA
	COUNTRY_MZ
	COUNTRY_MMR
	COUNTRY_NA
	COUNTRY_NR
	COUNTRY_NP
	COUNTRY_NL
	COUNTRY_NC
	COUNTRY_NZ
	COUNTRY_NI
	COUNTRY_NE
	COUNTRY_NG
	COUNTRY_NU
	COUNTRY_NF
	COUNTRY_MKD
	COUNTRY_MP
	COUNTRY_NO
	COUNTRY_OM
	COUNTRY_PK
	COUNTRY_PW
	COUNTRY_PS
	COUNTRY_PA
	COUNTRY_PG
	COUNTRY_PY
	COUNTRY_PE
	COUNTRY_PH
	COUNTRY_PCN
	COUNTRY_PL
	COUNTRY_PT
	COUNTRY_PR
	COUNTRY_QA
	COUNTRY_RE
	COUNTRY_RO
	COUNTRY_RU
	COUNTRY_RW
	COUNTRY_BL
	COUNTRY_SH
	COUNTRY_KN
	COUNTRY_LC
	COUNTRY_MF
	COUNTRY_PM
	COUNTRY_VC
	COUNTRY_WS
	COUNTRY_SM
	COUNTRY_ST
	COUNTRY_SA
	COUNTRY_SN
	COUNTRY_RS
	COUNTRY_SC
	COUNTRY_SL
	COUNTRY_SG
	COUNTRY_SX
	COUNTRY_SK
	COUNTRY_SI
	COUNTRY_SB
	COUNTRY_SO
	COUNTRY_ZA
	COUNTRY_GS
	COUNTRY_SS
	COUNTRY_ES
	COUNTRY_LK
	COUNTRY_SD
	COUNTRY_SR
	COUNTRY_SJ
	COUNTRY_SE
	COUNTRY_CH
	COUNTRY_SY
	COUNTRY_TJ
	COUNTRY_TZ
	COUNTRY_TH
	COUNTRY_TL
	COUNTRY_TG
	COUNTRY_TK
	COUNTRY_TO
	COUNTRY_TT
	COUNTRY_TN
	COUNTRY_TR
	COUNTRY_TM
	COUNTRY_TC
	COUNTRY_TV
	COUNTRY_UG
	COUNTRY_UA
	COUNTRY_AE
	COUNTRY_GB
	COUNTRY_UMI
	COUNTRY_US
	COUNTRY_UY
	COUNTRY_UZ
	COUNTRY_VU
	COUNTRY_VE
	COUNTRY_VN
	COUNTRY_VGB
	COUNTRY_VI
	COUNTRY_WF
	COUNTRY_YE
	COUNTRY_ZM
	COUNTRY_ZW
	COUNTRY_PLACEHOLDER_LAST
)

var countryToString map[Country]string = map[Country]string{
	COUNTRY_AF:  "Afghanistan",
	COUNTRY_AX:  "Åland Islands",
	COUNTRY_AL:  "Albania",
	COUNTRY_DZ:  "Algeria",
	COUNTRY_AS:  "American Samoa",
	COUNTRY_AD:  "Andorra",
	COUNTRY_AO:  "Angola",
	COUNTRY_AI:  "Anguilla",
	COUNTRY_ATA: "Antarctica ",
	COUNTRY_AG:  "Antigua and Barbuda",
	COUNTRY_AR:  "Argentina",
	COUNTRY_AM:  "Armenia",
	COUNTRY_AW:  "Aruba",
	COUNTRY_AUS: "Australia ",
	COUNTRY_AT:  "Austria",
	COUNTRY_AZ:  "Azerbaijan",
	COUNTRY_BS:  "Bahamas (the)",
	COUNTRY_BH:  "Bahrain",
	COUNTRY_BD:  "Bangladesh",
	COUNTRY_BB:  "Barbados",
	COUNTRY_BY:  "Belarus",
	COUNTRY_BE:  "Belgium",
	COUNTRY_BZ:  "Belize",
	COUNTRY_BJ:  "Benin",
	COUNTRY_BM:  "Bermuda",
	COUNTRY_BT:  "Bhutan",
	COUNTRY_BO:  "Bolivia (Plurinational State of)",
	COUNTRY_BQ:  "Bonaire;Sint Eustatius;Saba",
	COUNTRY_BA:  "Bosnia and Herzegovina",
	COUNTRY_BW:  "Botswana",
	COUNTRY_BV:  "Bouvet Island",
	COUNTRY_BR:  "Brazil",
	COUNTRY_IO:  "British Indian Ocean Territory (the)",
	COUNTRY_BRN: "Brunei Darussalam ",
	COUNTRY_BG:  "Bulgaria",
	COUNTRY_BF:  "Burkina Faso",
	COUNTRY_BI:  "Burundi",
	COUNTRY_CV:  "Cabo Verde ",
	COUNTRY_KH:  "Cambodia",
	COUNTRY_CM:  "Cameroon",
	COUNTRY_CA:  "Canada",
	COUNTRY_KY:  "Cayman Islands (the)",
	COUNTRY_CF:  "Central African Republic (the)",
	COUNTRY_TD:  "Chad",
	COUNTRY_CL:  "Chile",
	COUNTRY_CN:  "China",
	COUNTRY_CX:  "Christmas Island",
	COUNTRY_CC:  "Cocos (Keeling) Islands (the)",
	COUNTRY_CO:  "Colombia",
	COUNTRY_KM:  "Comoros (the)",
	COUNTRY_CD:  "Congo (the Democratic Republic of the)",
	COUNTRY_CG:  "Congo (the) ",
	COUNTRY_CK:  "Cook Islands (the)",
	COUNTRY_CR:  "Costa Rica",
	COUNTRY_CI:  "Côte d'Ivoire ",
	COUNTRY_HR:  "Croatia",
	COUNTRY_CU:  "Cuba",
	COUNTRY_CW:  "Curaçao",
	COUNTRY_CY:  "Cyprus",
	COUNTRY_CZ:  "Czechia ",
	COUNTRY_DK:  "Denmark",
	COUNTRY_DJ:  "Djibouti",
	COUNTRY_DM:  "Dominica",
	COUNTRY_DO:  "Dominican Republic (the)",
	COUNTRY_EC:  "Ecuador",
	COUNTRY_EG:  "Egypt",
	COUNTRY_SV:  "El Salvador",
	COUNTRY_GQ:  "Equatorial Guinea",
	COUNTRY_ER:  "Eritrea",
	COUNTRY_EE:  "Estonia",
	COUNTRY_SZ:  "Eswatini ",
	COUNTRY_ET:  "Ethiopia",
	COUNTRY_FLK: "Falkland Islands (the)",
	COUNTRY_FO:  "Faroe Islands (the)",
	COUNTRY_FJ:  "Fiji",
	COUNTRY_FI:  "Finland",
	COUNTRY_FR:  "France",
	COUNTRY_GF:  "French Guiana",
	COUNTRY_PF:  "French Polynesia",
	COUNTRY_ATF: "French Southern Territories (the)",
	COUNTRY_GA:  "Gabon",
	COUNTRY_GM:  "Gambia (the)",
	COUNTRY_GE:  "Georgia",
	COUNTRY_DE:  "Germany",
	COUNTRY_GH:  "Ghana",
	COUNTRY_GI:  "Gibraltar",
	COUNTRY_GR:  "Greece",
	COUNTRY_GL:  "Greenland",
	COUNTRY_GD:  "Grenada",
	COUNTRY_GP:  "Guadeloupe",
	COUNTRY_GU:  "Guam",
	COUNTRY_GT:  "Guatemala",
	COUNTRY_GG:  "Guernsey",
	COUNTRY_GN:  "Guinea",
	COUNTRY_GW:  "Guinea-Bissau",
	COUNTRY_GY:  "Guyana",
	COUNTRY_HT:  "Haiti",
	COUNTRY_HM:  "Heard Island and McDonald Islands",
	COUNTRY_HN:  "Honduras",
	COUNTRY_HK:  "Hong Kong",
	COUNTRY_HU:  "Hungary",
	COUNTRY_IS:  "Iceland",
	COUNTRY_IN:  "India",
	COUNTRY_ID:  "Indonesia",
	COUNTRY_IR:  "Iran (Islamic Republic of)",
	COUNTRY_IQ:  "Iraq",
	COUNTRY_IE:  "Ireland",
	COUNTRY_IM:  "Isle of Man",
	COUNTRY_IL:  "Israel",
	COUNTRY_IT:  "Italy",
	COUNTRY_JM:  "Jamaica",
	COUNTRY_JP:  "Japan",
	COUNTRY_JE:  "Jersey",
	COUNTRY_JO:  "Jordan",
	COUNTRY_KZ:  "Kazakhstan",
	COUNTRY_KE:  "Kenya",
	COUNTRY_KI:  "Kiribati",
	COUNTRY_KP:  "Korea (the Democratic People's Republic of)",
	COUNTRY_KR:  "Korea (the Republic of)",
	COUNTRY_KW:  "Kuwait",
	COUNTRY_KG:  "Kyrgyzstan",
	COUNTRY_LA:  "Lao People's Democratic Republic (the)",
	COUNTRY_LV:  "Latvia",
	COUNTRY_LB:  "Lebanon",
	COUNTRY_LS:  "Lesotho",
	COUNTRY_LR:  "Liberia",
	COUNTRY_LY:  "Libya",
	COUNTRY_LI:  "Liechtenstein",
	COUNTRY_LT:  "Lithuania",
	COUNTRY_LU:  "Luxembourg",
	COUNTRY_MAC: "Macao ",
	COUNTRY_MG:  "Madagascar",
	COUNTRY_MW:  "Malawi",
	COUNTRY_MY:  "Malaysia",
	COUNTRY_MV:  "Maldives",
	COUNTRY_ML:  "Mali",
	COUNTRY_MT:  "Malta",
	COUNTRY_MH:  "Marshall Islands (the)",
	COUNTRY_MQ:  "Martinique",
	COUNTRY_MR:  "Mauritania",
	COUNTRY_MU:  "Mauritius",
	COUNTRY_YT:  "Mayotte",
	COUNTRY_MX:  "Mexico",
	COUNTRY_FM:  "Micronesia (Federated States of)",
	COUNTRY_MD:  "Moldova (the Republic of)",
	COUNTRY_MC:  "Monaco",
	COUNTRY_MN:  "Mongolia",
	COUNTRY_ME:  "Montenegro",
	COUNTRY_MS:  "Montserrat",
	COUNTRY_MA:  "Morocco",
	COUNTRY_MZ:  "Mozambique",
	COUNTRY_MMR: "Myanmar ",
	COUNTRY_NA:  "Namibia",
	COUNTRY_NR:  "Nauru",
	COUNTRY_NP:  "Nepal",
	COUNTRY_NL:  "Netherlands (Kingdom of the)",
	COUNTRY_NC:  "New Caledonia",
	COUNTRY_NZ:  "New Zealand",
	COUNTRY_NI:  "Nicaragua",
	COUNTRY_NE:  "Niger (the)",
	COUNTRY_NG:  "Nigeria",
	COUNTRY_NU:  "Niue",
	COUNTRY_NF:  "Norfolk Island",
	COUNTRY_MKD: "North Macedonia ",
	COUNTRY_MP:  "Northern Mariana Islands (the)",
	COUNTRY_NO:  "Norway",
	COUNTRY_OM:  "Oman",
	COUNTRY_PK:  "Pakistan",
	COUNTRY_PW:  "Palau",
	COUNTRY_PS:  "Palestine, State of",
	COUNTRY_PA:  "Panama",
	COUNTRY_PG:  "Papua New Guinea",
	COUNTRY_PY:  "Paraguay",
	COUNTRY_PE:  "Peru",
	COUNTRY_PH:  "Philippines (the)",
	COUNTRY_PCN: "Pitcairn ",
	COUNTRY_PL:  "Poland",
	COUNTRY_PT:  "Portugal",
	COUNTRY_PR:  "Puerto Rico",
	COUNTRY_QA:  "Qatar",
	COUNTRY_RE:  "Réunion",
	COUNTRY_RO:  "Romania",
	COUNTRY_RU:  "Russian Federation (the)",
	COUNTRY_RW:  "Rwanda",
	COUNTRY_BL:  "Saint Barthélemy",
	COUNTRY_SH:  "Saint Helena Ascension Island Tristan da Cunha",
	COUNTRY_KN:  "Saint Kitts and Nevis",
	COUNTRY_LC:  "Saint Lucia",
	COUNTRY_MF:  "Saint Martin (French part)",
	COUNTRY_PM:  "Saint Pierre and Miquelon",
	COUNTRY_VC:  "Saint Vincent and the Grenadines",
	COUNTRY_WS:  "Samoa",
	COUNTRY_SM:  "San Marino",
	COUNTRY_ST:  "Sao Tome and Principe",
	COUNTRY_SA:  "Saudi Arabia",
	COUNTRY_SN:  "Senegal",
	COUNTRY_RS:  "Serbia",
	COUNTRY_SC:  "Seychelles",
	COUNTRY_SL:  "Sierra Leone",
	COUNTRY_SG:  "Singapore",
	COUNTRY_SX:  "Sint Maarten (Dutch part)",
	COUNTRY_SK:  "Slovakia",
	COUNTRY_SI:  "Slovenia",
	COUNTRY_SB:  "Solomon Islands",
	COUNTRY_SO:  "Somalia",
	COUNTRY_ZA:  "South Africa",
	COUNTRY_GS:  "South Georgia and the South Sandwich Islands",
	COUNTRY_SS:  "South Sudan",
	COUNTRY_ES:  "Spain",
	COUNTRY_LK:  "Sri Lanka",
	COUNTRY_SD:  "Sudan (the)",
	COUNTRY_SR:  "Suriname",
	COUNTRY_SJ:  "Svalbard Jan Mayen",
	COUNTRY_SE:  "Sweden",
	COUNTRY_CH:  "Switzerland",
	COUNTRY_SY:  "Syrian Arab Republic (the)",
	COUNTRY_TJ:  "Tajikistan",
	COUNTRY_TZ:  "Tanzania, the United Republic of",
	COUNTRY_TH:  "Thailand",
	COUNTRY_TL:  "Timor-Leste",
	COUNTRY_TG:  "Togo",
	COUNTRY_TK:  "Tokelau",
	COUNTRY_TO:  "Tonga",
	COUNTRY_TT:  "Trinidad and Tobago",
	COUNTRY_TN:  "Tunisia",
	COUNTRY_TR:  "Türkiye",
	COUNTRY_TM:  "Turkmenistan",
	COUNTRY_TC:  "Turks and Caicos Islands (the)",
	COUNTRY_TV:  "Tuvalu",
	COUNTRY_UG:  "Uganda",
	COUNTRY_UA:  "Ukraine",
	COUNTRY_AE:  "United Arab Emirates (the)",
	COUNTRY_GB:  "United Kingdom of Great Britain and Northern Ireland (the)",
	COUNTRY_UMI: "United States Minor Outlying Islands (the)",
	COUNTRY_US:  "United States of America (the)",
	COUNTRY_UY:  "Uruguay",
	COUNTRY_UZ:  "Uzbekistan",
	COUNTRY_VU:  "Vanuatu",
	COUNTRY_VE:  "Venezuela (Bolivarian Republic of)",
	COUNTRY_VN:  "Viet Nam",
	COUNTRY_VGB: "Virgin Islands (British)",
	COUNTRY_VI:  "Virgin Islands (U.S.) ",
	COUNTRY_WF:  "Wallis and Futuna",
	COUNTRY_YE:  "Yemen",
	COUNTRY_ZM:  "Zambia",
	COUNTRY_ZW:  "Zimbabwe",
}

func (c Country) String() string {
	var str string
	var ok bool
	if str, ok = countryToString[c]; !ok {
		str = "Unknown"
	}
	return fmt.Sprintf("Country(%s, %d)", str, c)
}

func RandomCountry() Country {
	return Country(rand.Intn(COUNTRY_PLACEHOLDER_LAST))
}

type CountryBitset [4]uint64

// A set of countries. Repeated elements will be ignored.
type Countries []Country

func RandomCountries() Countries {
	number := rand.Intn(COUNTRY_PLACEHOLDER_LAST)
	countries := make(Countries, number)
	for i := range number {
		countries[i] = RandomCountry()
	}
	return countries
}

func (c *Countries) ToBitset() CountryBitset {
	var e CountryBitset

	for _, country := range *c {
		arrIndex := int(country / 64)
		bitIndex := int(country % 64)
		mask := uint64(1) << bitIndex
		e[arrIndex] = e[arrIndex] | mask
	}
	return e
}

func (c *CountryBitset) ToCountries() Countries {
	countries := Countries{}
	for i, elem := range c {
		bitIndex := 0
		for elem != 0 {
			if elem&1 == 1 {
				country := Country(bitIndex + i*64)
				countries = append(countries, country)
			}
			bitIndex++
			elem = elem >> 1
		}
	}
	return countries
}

func (c *Countries) Value() (driver.Value, error) {
	b := make([]byte, 32)
	bitset := c.ToBitset()
	binary.BigEndian.PutUint64(b[0:8], bitset[0])
	binary.BigEndian.PutUint64(b[8:16], bitset[1])
	binary.BigEndian.PutUint64(b[16:24], bitset[2])
	binary.BigEndian.PutUint64(b[24:32], bitset[3])
	return b, nil
}

func (e *Countries) Scan(src any) error {
	var bitset CountryBitset
	switch val := src.(type) {
	case []byte:
		bitset[0] = binary.BigEndian.Uint64(val[0:8])
		bitset[1] = binary.BigEndian.Uint64(val[8:16])
		bitset[2] = binary.BigEndian.Uint64(val[16:24])
		bitset[3] = binary.BigEndian.Uint64(val[24:32])
	default:
		return fmt.Errorf("Unknown type of src: %v", src)
	}
	*e = bitset.ToCountries()
	return nil
}
