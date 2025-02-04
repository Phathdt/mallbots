package errorx

import "errors"

var (
	ErrCannotGetUser     = errors.New("cannot get user")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrCreateUser        = errors.New("create user failed")
	ErrPasswordNotMatch  = errors.New("password not match")
	ErrGenToken          = errors.New("when gen token")
	ErrCannotLogin       = errors.New("cannot login")
)

var (
	// Order errors
	ErrCartEmpty               = errors.New("cart is empty")
	ErrOrderNotFound           = errors.New("order not found")
	ErrCannotCreateOrder       = errors.New("cannot create order")
	ErrCannotUpdateOrder       = errors.New("cannot update order")
	ErrInvalidOrderStatus      = errors.New("invalid order status")
	ErrInvalidStatusTransition = errors.New("invalid status transition")
	ErrUnauthorizedOrderAccess = errors.New("unauthorized access to order")

	// Payment errors
	ErrInvalidPaymentStatus           = errors.New("invalid payment status")
	ErrInvalidPaymentStatusTransition = errors.New("invalid payment status transition")
	ErrPaymentFailed                  = errors.New("payment failed")
	ErrPaymentAlreadyProcessed        = errors.New("payment already processed")

	// Shipping errors
	ErrInvalidShippingAddress    = errors.New("invalid shipping address")
	ErrInvalidShippingCountry    = errors.New("shipping not available in this country")
	ErrShippingCalculationFailed = errors.New("failed to calculate shipping cost")

	// Order Items errors
	ErrProductOutOfStock      = errors.New("product is out of stock")
	ErrInvalidProductQuantity = errors.New("invalid product quantity")
	ErrProductPriceChanged    = errors.New("product price has changed")
	ErrCannotCreateOrderItems = errors.New("cannot create order items")

	// Business Logic errors
	ErrOrderAlreadyCancelled        = errors.New("order is already cancelled")
	ErrOrderNotCancellable          = errors.New("order cannot be cancelled at this stage")
	ErrOrderNotRefundable           = errors.New("order is not eligible for refund")
	ErrMinimumOrderAmountNotMet     = errors.New("minimum order amount not met")
	ErrMaximumOrderQuantityExceeded = errors.New("maximum order quantity exceeded")
)
