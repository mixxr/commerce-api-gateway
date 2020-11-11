package mockauth

type MockAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (o *MockAuth) IsValid() bool {
	return true
}

func (o *MockAuth) GetUsername() string {
	return o.Username
}
