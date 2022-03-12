package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"
	"ksemilla/database"
	"ksemilla/graph/generated"
	"ksemilla/graph/model"
	"ksemilla/middlewares"
)

func (r *mutationResolver) CreateInvoice(ctx context.Context, input model.NewInvoice) (*model.Invoice, error) {
	err := GetUserPermission(ctx, []string{OWNER})
	if err != nil {
		return nil, err
	}
	return db.Save(&input), nil
}

func (r *mutationResolver) UpdateInvoice(ctx context.Context, input model.InvoiceInput) (*model.Invoice, error) {
	err := GetUserPermission(ctx, []string{OWNER})
	if err != nil {
		return nil, err
	}
	return db.UpdateInvoice(&input), nil
}

func (r *mutationResolver) Login(ctx context.Context, input model.Login) (*model.LoginReturn, error) {
	fmt.Println("xxx")
	return db.Login(&input)
}

func (r *mutationResolver) VerifyToken(ctx context.Context, input model.VerifyToken) (*model.User, error) {
	return db.VerifyToken(&input)
}

func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (*model.User, error) {
	err := GetUserPermission(ctx, []string{OWNER})
	if err != nil {
		return nil, err
	}
	return db.CreateUser(&input), nil
}

func (r *mutationResolver) FindUserByID(ctx context.Context, input string) (*model.User, error) {
	err := GetUserPermission(ctx, []string{OWNER})
	if err != nil {
		return nil, err
	}
	return db.FindOneUser(input)
}

func (r *mutationResolver) UpdateUser(ctx context.Context, input model.UpdateUser) (*model.User, error) {
	err := GetUserPermission(ctx, []string{OWNER})
	if err != nil {
		return nil, err
	}
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Invoices(ctx context.Context, page int) (*model.PaginatedInvoicesReturn, error) {
	err := GetUserPermission(ctx, []string{OWNER, ACCT})
	if err != nil {
		return nil, err
	}
	invoices, total := db.InvoicesPaginated(int64(page))
	return &model.PaginatedInvoicesReturn{
		Data:  invoices,
		Total: int(total),
	}, nil
}

func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	err := GetUserPermission(ctx, []string{OWNER})
	if err != nil {
		return nil, err
	}
	return db.AllUsers(), nil
}

func (r *queryResolver) GetInvoice(ctx context.Context, id string) (*model.Invoice, error) {
	err := GetUserPermission(ctx, []string{OWNER, ACCT})
	if err != nil {
		return nil, err
	}
	return db.GetInvoice(id), nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.

// func (r *queryResolver) InvoiceFilter(ctx context.Context, dateCreated string) ([]*model.Invoice, error) {
// 	err := GetUserPermission(ctx, []string{OWNER})
// 	if err != nil {
// 		return nil, err
// 	}
// 	res := db.PaginatedInvoice(dateCreated)
// 	return res, nil
// }

var OWNER = "owner"
var ACCT = "acct"

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
func GetUserPermission(ctx context.Context, role []string) error {
	user := ctx.Value(middlewares.GetUserCtx())

	if user == nil {
		return errors.New("auth failed")
	}

	testUser := user.(model.User)
	fmt.Println("testing", testUser.Role)

	if !contains(role, testUser.Role) {
		return errors.New("permission denied")
	}

	return nil
}

var db = database.Connect()
