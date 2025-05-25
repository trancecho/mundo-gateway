package point_dto

type ChangePointReq struct {
	Token      string `json:"token"`
	UserID     string `json:"user_id"`
	Point      int64  `json:"point"`
	Experience int64  `json:"experience"`
	Reason     string `json:"reason"`
}

type SignReq struct {
	UserId string `json:"user_id"`
	Token  string `json:"token"`
}

type PointsInfoReq struct {
	UserId string `json:"user_id"`
	Token  string `json:"token"`
}

type AdminStatsReq struct {
	Token  string `json:"token"`
	UserId string `json:"user_id"`
}
