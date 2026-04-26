package Query

const (
	GetStatesDropdownQuery      = `SELECT "stateid" AS id, "statename" AS name FROM "states" ORDER BY "statename" ASC;`
	GetSupplyTypesDropdownQuery = `SELECT "supplytypeid" AS id, "typename" AS name FROM "supplytypes" ORDER BY "typename" ASC;`
	GetCountriesDropdownQuery   = `SELECT "countryid" AS id, "countryname" AS name FROM "countries" ORDER BY "countryname" ASC;`
	GetRolesDropdownQuery       = `SELECT "roleid" AS id, "rolename" AS name FROM "roles" ORDER BY "rolename" ASC;`
)