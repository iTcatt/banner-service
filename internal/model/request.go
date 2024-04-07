package model

type GetUserBannerParams struct {
	TagID           string
	FeatureID       string
	UseLastRevision bool
}

type GetBannerWithFiltersParams struct {
	TagID     string
	FeatureID string
	Limit     string
	Offset    string
}
