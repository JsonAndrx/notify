package models

type Usage struct {
	PK                string `dynamodbav:"PK"`
	SK                string `dynamodbav:"SK"`
	BusinessID        string `dynamodbav:"businessId"`
	PlanID            string `dynamodbav:"planId"`
	NotificationCount int    `dynamodbav:"notificationCount"`
	PeriodStart       string `dynamodbav:"periodStart"`
	PeriodEnd         string `dynamodbav:"periodEnd"`
	CreatedAt         string `dynamodbav:"createdAt"`
	UpdatedAt         string `dynamodbav:"updatedAt,omitempty"`
}
