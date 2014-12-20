Pagination Over Martini
=======================

Simple service to create pagination GET REST API quickly.
Uses 'render' martini contrib package.

Example:
--------
using GORM

~~~ go
import "github.com/shlomomatichin/go-martini-paginition"

func ListOrganizationsView(p *pagination.Pagination, db gorm.DB) {
	var organizations []Organization
    var count unit64
	db.Offset(p.Offset).Limit(p.PerPage).Find(&organizations).Count(&count)
    p.SetTotal(count)
	for _, organization := range organizations {
		organizationJSON := map[string]interface{}{
			"Id":   organization.Id,
			"Name": organization.Name,
			"Center": map[string]interface{}{
				"Longitude": organization.CenterLongitude,
				"Latitude":  organization.CenterLatitude,
			},
		}
        p.Append(organizationJSON);
	}
}

func RegisterOrganizationRESTViews(m *martini.ClassicMartini) {
	m.Get("/api/v1/organization", RoleAllowed(), pagination.Service, ListOrganizationsView)
}
~~~
