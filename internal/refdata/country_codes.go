package refdata

// CountryCodes returns country, coalition, and organization codes used in
// ownerProducer, releasableTo, displayOnlyTo, and FGI source fields.
// This is a curated subset of ISO 3166 trigraphs plus IC-standard coalition
// and organization designators commonly used in ISM markings.
func CountryCodes() []CountryCode {
	return []CountryCode{
		// Countries (ISO 3166-1 alpha-3, curated for ISM use)
		{Code: "USA", Name: "United States of America", Type: "country"},
		{Code: "GBR", Name: "United Kingdom", Type: "country"},
		{Code: "CAN", Name: "Canada", Type: "country"},
		{Code: "AUS", Name: "Australia", Type: "country"},
		{Code: "NZL", Name: "New Zealand", Type: "country"},
		{Code: "FRA", Name: "France", Type: "country"},
		{Code: "DEU", Name: "Germany", Type: "country"},
		{Code: "ITA", Name: "Italy", Type: "country"},
		{Code: "JPN", Name: "Japan", Type: "country"},
		{Code: "KOR", Name: "Republic of Korea", Type: "country"},
		{Code: "ISR", Name: "Israel", Type: "country"},
		{Code: "NOR", Name: "Norway", Type: "country"},
		{Code: "DNK", Name: "Denmark", Type: "country"},
		{Code: "NLD", Name: "Netherlands", Type: "country"},
		{Code: "BEL", Name: "Belgium", Type: "country"},
		{Code: "ESP", Name: "Spain", Type: "country"},
		{Code: "PRT", Name: "Portugal", Type: "country"},
		{Code: "TUR", Name: "Turkey", Type: "country"},
		{Code: "POL", Name: "Poland", Type: "country"},
		{Code: "SWE", Name: "Sweden", Type: "country"},

		// Coalitions
		{Code: "FVEY", Name: "Five Eyes", Type: "coalition"},
		{Code: "ACGU", Name: "AUSCANNZUKUS (Four Eyes + AUS)", Type: "coalition"},
		{Code: "TEYE", Name: "Three Eyes (USA, GBR, CAN)", Type: "coalition"},
		{Code: "CMFC", Name: "Combined Maritime Forces Coalition", Type: "coalition"},
		{Code: "CMFP", Name: "Coalition Maritime Force Pacific", Type: "coalition"},
		{Code: "CPMT", Name: "Cooperative Maritime Partner Teaming", Type: "coalition"},
		{Code: "ISAF", Name: "International Security Assistance Force", Type: "coalition"},
		{Code: "KFOR", Name: "Kosovo Force", Type: "coalition"},
		{Code: "RSMA", Name: "Resolute Support Mission Afghanistan", Type: "coalition"},

		// Organizations
		{Code: "NATO", Name: "North Atlantic Treaty Organization", Type: "organization"},
		{Code: "BWCS", Name: "Biological Weapons Convention States", Type: "organization"},
	}
}
