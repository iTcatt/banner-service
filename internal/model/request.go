package model

type GetUserBannerParams struct {
	TagID           int
	FeatureID       int
	UseLastRevision bool
	IsAdmin         bool
}

type GetFilteredBannersParams struct {
	TagID     int
	FeatureID int
	Limit     int
	Offset    int
}

type BannerParams struct {
	TagIDs    []int       `json:"tag_ids"`
	FeatureID int         `json:"feature_id"`
	Content   interface{} `json:"content"`
	IsActive  bool        `json:"is_active"`
}
