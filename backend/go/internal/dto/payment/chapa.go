package payment

type ChapaInitializeRequest struct {
	Amount        string            `json:"amount"`
	Currency      string            `json:"currency"`
	Email         string            `json:"email"`
	FirstName     string            `json:"first_name"`
	LastName      string            `json:"last_name"`
	// TxRef: a unique ID generated for the payment
	// Format: recipe-{recipeID}-user-{userID}-{timestamp}
	TxRef         string            `json:"tx_ref"`
	// CallbackURL: A backend endpoint Chapa calls (POST) after payment completes
	CallbackURL   string            `json:"callback_url"`
	// ReturnURL: Where the user's browser is redirected after they finish on the Chapa checkout page
	ReturnURL     string            `json:"return_url"`
	Customization map[string]string `json:"customization,omitempty"`
}

type ChapaInitializeResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
	Data    struct {
		CheckoutURL string `json:"checkout_url"`
	} `json:"data"`
}

type ChapaVerifyResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
	Data    struct {
		Amount   float64 `json:"amount"`
		Currency string  `json:"currency"`
		Status   string  `json:"status"`
		TxRef    string  `json:"tx_ref"`
	} `json:"data"`
}

// Request payload from our frontend
type InitializePaymentRequest struct {
	RecipeID int `json:"recipeId"`
}

type InitializePaymentResponse struct {
	CheckoutURL string `json:"checkoutUrl"`
}
