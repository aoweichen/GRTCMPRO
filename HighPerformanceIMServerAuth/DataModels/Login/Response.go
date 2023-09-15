package Login

type Response struct {
	ID              int64  `json:"id"`
	UID             string `json:"uid"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	Avatar          string `json:"avatar"`
	Token           string `json:"token"`
	ExpireTime      int64  `json:"expire_time"`
	TokenTimeToLive int64  `json:"token_time_to_live"`
}
