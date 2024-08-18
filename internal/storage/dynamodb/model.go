package dynamodb

type User struct {
	UserID string `dynamodbav:"UserID"`
	Name   string `dynamodbav:"Name"`
	Age    int    `dynamodbav:"Age"`
}
