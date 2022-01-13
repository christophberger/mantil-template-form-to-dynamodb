package form

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/mantil-io/mantil.go"
)

const (
	TableKey     = "email"
	TableSortKey = "name"
)

type Form struct {
	Name              string   `json:"What is your name?"`
	CanYouAttend      string   `json:"Can you attend?"`
	Count             string   `json:"How many of you are attending?"`
	Items             []string `json:"What will you be bringing?"`
	Restrictions      string   `json:"Do you have any allergies or dietary restrictions?"`
	Email             string   `json:"What is your email address?"`
	table             *dynamodb.Client
	tableName         string
	tableResourceName *string
}

func New() *Form {
	tableName := os.Getenv("TABLE_NAME")
	if tableName == "" {
		tableName = "MantilPartyTable"
	}
	table, err := mantil.DynamodbTable(tableName, TableKey, TableSortKey)
	if err != nil {
		log.Fatalf("Cannot create table: %s", err)
	}
	return &Form{
		table:             table,
		tableResourceName: aws.String(mantil.Resource(tableName).Name),
	}
}

type DefaultRequest struct{}

func (f *Form) Default(ctx context.Context, req *DefaultRequest) error {
	log.Println("Default called")
	return nil
}

type SaveResponse struct {
	msg string
}

// Save receives the Party form response and saves it to a DynamoDB table.
func (f *Form) Save(ctx context.Context, req *Form) (*SaveResponse, error) {
	log.Printf("Save: req is '%+v'", req)
	items := &types.AttributeValueMemberSS{Value: req.Items}
	if items == nil || len(items.Value) == 0 {
		// A string set may not be empty
		items = &types.AttributeValueMemberSS{Value: []string{""}}
	}
	_, err := f.table.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: f.tableResourceName,
		Item: map[string]types.AttributeValue{
			"name":         &types.AttributeValueMemberS{Value: req.Name},
			"canattend":    &types.AttributeValueMemberS{Value: req.CanYouAttend},
			"count":        &types.AttributeValueMemberS{Value: req.Count},
			"items":        items,
			"restrictions": &types.AttributeValueMemberS{Value: req.Restrictions},
			"email":        &types.AttributeValueMemberS{Value: req.Email},
		},
	})
	if err != nil {
		log.Printf("Cannot save form: %s", err)
		return &SaveResponse{msg: "Cannot save form"}, err
	}
	return &SaveResponse{msg: fmt.Sprintf("%s saved", req.Name)}, nil
}

// List returns a list of all parties in the DynamoDB table. Useful for quickly checking the table contents via "mantil invoke form/list".
func (f *Form) List(ctx context.Context, req *DefaultRequest) (*[]Form, error) {
	log.Println("List called")
	out, err := f.table.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: f.tableResourceName,
	})
	if err != nil {
		log.Printf("Cannot scan table: %s", err)
	}
	log.Println(len(out.Items), out.Items)
	var forms []Form
	for _, item := range out.Items {
		forms = append(forms, Form{
			Name:         (item["name"]).(*types.AttributeValueMemberS).Value,
			CanYouAttend: (item["canattend"]).(*types.AttributeValueMemberS).Value,
			Count:        (item["count"]).(*types.AttributeValueMemberS).Value,
			Items:        (item["items"]).(*types.AttributeValueMemberSS).Value,
			Restrictions: (item["restrictions"]).(*types.AttributeValueMemberS).Value,
			Email:        (item["email"]).(*types.AttributeValueMemberS).Value,
		})
	}
	return &forms, nil
}
