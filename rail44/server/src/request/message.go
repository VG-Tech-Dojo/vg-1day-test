package request

type (
	Message struct {
		Body string `json:"body"`
		Username string `json:"user_name"`
	}

        UpdateMessage struct {
                Body string `json:"body"`
        }
)
