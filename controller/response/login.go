package response

type Login struct {
	Id       *string `json:"id"`
	JwtToken *string `json:"jwt_token"`
}
