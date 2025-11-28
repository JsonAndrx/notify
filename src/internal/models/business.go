package models

type Business struct {
	PK        string `dynamodbav:"PK"`
	SK        string `dynamodbav:"SK"`
	Name      string `dynamodbav:"name"`
	Email     string `dynamodbav:"email"`
	Phone     string `dynamodbav:"phone"`
	PlanID    string `dynamodbav:"planId"`
	APIKey    string `dynamodbav:"apiKey"`
	CreatedAt string `dynamodbav:"createdAt"`
	UpdatedAt string `dynamodbav:"updatedAt,omitempty"`
}
