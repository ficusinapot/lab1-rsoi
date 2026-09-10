package persons

type IDInput struct {
	ID int `path:"id" example:"1" doc:"Person ID"`
}

type Input struct {
	Body Request
}

type WithIDInput struct {
	ID   int `path:"id" example:"1" doc:"Person ID"`
	Body Request
}

type ListOutput struct {
	Body []Response
}

type Output struct {
	Body Response
}

type Request struct {
	Name    *string `json:"name,omitempty"`
	Age     *int    `json:"age,omitempty"`
	Address *string `json:"address,omitempty"`
	Work    *string `json:"work,omitempty"`
}

type Response struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Address string `json:"address"`
	Work    string `json:"work"`
}
