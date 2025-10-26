package fields

type UpdateFieldRequest struct {
	ImageUrl     string  `json:"image_url" form:"image_url" validate:"omitempty,url"`
	Name         string  `json:"name" form:"name" validate:"omitempty,min=3,max=100"`
	Type         string  `json:"type" form:"type" validate:"omitempty,min=3,max=100"`
	Description  string  `json:"desc" form:"desc"`
	Location     string  `json:"location" form:"location"`
	LocationLat  float64 `json:"location_lat" form:"location_lat"`
	LocationLon  float64 `json:"location_lon" form:"location_lon"`
	PricePerHour int64   `json:"price_per_hour" form:"price_per_hour" validate:"omitempty,min=1"`
}

type CreateFieldRequest struct {
	ImageUrl     string  `json:"image_url" form:"image_url" validate:"omitempty,url"`
	Name         string  `json:"name" form:"name" validate:"required,min=3,max=100"`
	Type         string  `json:"type" form:"type" validate:"required,min=3,max=100"`
	Description  string  `json:"desc" form:"desc"`
	Location     string  `json:"location" form:"location"`
	LocationLat  float64 `json:"location_lat" form:"location_lat"`
	LocationLon  float64 `json:"location_lon" form:"location_lon"`
	PricePerHour int64   `json:"price_per_hour" form:"price_per_hour" validate:"required,min=1"`
}
