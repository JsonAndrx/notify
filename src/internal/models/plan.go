package models

type Plan struct {
	PK                string  `dynamodbav:"PK"`
	SK                string  `dynamodbav:"SK"`
	Name              string  `dynamodbav:"name"`
	NotificationLimit int     `dynamodbav:"notificationLimit"`
	PeriodDays        int     `dynamodbav:"periodDays"`
	Price             float64 `dynamodbav:"price"`
	Description       string  `dynamodbav:"description"`
	Active            bool    `dynamodbav:"active"`
	CreatedAt         string  `dynamodbav:"createdAt"`
}
