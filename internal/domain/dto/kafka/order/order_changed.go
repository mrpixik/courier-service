package order

const (
	StatusCreated   = "created"
	StatusCancelled = "cancelled"
	StatusCompleted = "completed"
)

type Event struct {
	OrderID   string `json:"order_id"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

type ProcessedEvent struct {
	OrderId   string `json:"order_id"`
	Status    string `json:"status"`
	CourierId int    `json:"courier_id"`
}
