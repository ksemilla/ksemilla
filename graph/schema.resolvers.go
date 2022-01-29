package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"ksemilla/database"
	"ksemilla/graph/generated"
	"ksemilla/graph/model"
	"ksemilla/middlewares"
	"reflect"
)

func (r *mutationResolver) CreateInvoice(ctx context.Context, input model.NewInvoice) (*model.Invoice, error) {
	return db.Save(&input), nil
}

func (r *mutationResolver) UpdateInvoice(ctx context.Context, input model.InvoiceInput) (*model.Invoice, error) {
	fmt.Println(input)
	return db.UpdateInvoice(&input), nil
}

func (r *mutationResolver) Login(ctx context.Context, input model.Login) (string, error) {
	return db.Login(&input)
}

func (r *mutationResolver) VerifyToken(ctx context.Context, input model.VerifyToken) (string, error) {
	fmt.Println(ctx.Value("user"), ctx.Value(middlewares.GetUserCtx()))
	return db.VerifyToken(&input)
}

func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (*model.User, error) {
	return db.CreateUser(&input), nil
}

func (r *mutationResolver) FindUserByID(ctx context.Context, input string) (*model.User, error) {
	return db.FindOneUser(input)
}

func (r *mutationResolver) UpdateUser(ctx context.Context, input model.UpdateUser) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Invoices(ctx context.Context) ([]*model.Invoice, error) {
	user := ctx.Value(middlewares.GetUserCtx())
	fmt.Println("INVOICES USER", user, reflect.TypeOf(user))
	// testUser := user.(model.User)
	// fmt.Println("testing", testUser.Role)
	// return nil, errors.New("NOT AALOWED")
	return db.All(), nil
}

func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	return db.AllUsers(), nil
}

func (r *queryResolver) InvoiceFilter(ctx context.Context, dateCreated string) ([]*model.Invoice, error) {
	res := db.PaginatedInvoice(dateCreated)
	return res, nil
}

func (r *queryResolver) PaginatedInvoices(ctx context.Context, page int) (*model.PaginatedInvoicesReturn, error) {
	invoices, total := db.InvoicesPaginated(int64(page))
	return &model.PaginatedInvoicesReturn{
		Data:  invoices,
		Total: int(total),
	}, nil
}

func (r *queryResolver) GetInvoice(ctx context.Context, id string) (*model.Invoice, error) {
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
var db = database.Connect()
