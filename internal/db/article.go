package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/caring/go-packages/pkg/errors"
	"github.com/google/uuid"

	"github.com/caring/articles/pb"
)



// articleService provides an API for interacting with the articles table
type articleService struct {
	db    *sql.DB
	stmts map[string]*sql.Stmt
}

// Article is a struct representation of a row in the articles table
type Article struct {
	ID  	uuid.UUID
	Name  string
}

// protoArticle is an interface that most proto article objects will satisfy
type protoArticle interface {
	GetName() string
}

// NewArticle is a convenience helper cast a proto article to it's DB layer struct
func NewArticle(ID string, proto protoArticle) (*Article, error) {
	mID, err := ParseUUID(ID)
	if err != nil {
		return nil, err
	}

	return &Article{
		ID:  	mID,
		Name: proto.GetName(),
	}, nil
}

// ToProto casts a db article into a proto response object
func (m *Article) ToProto() *pb.ArticleResponse {
	return &pb.ArticleResponse{
		Id:  				m.ID.String(),
		Name:       m.Name,
	}
}

// Get fetches a single article from the db
func (svc *articleService) Get(ctx context.Context, ID uuid.UUID) (*Article, error) {
	return svc.get(ctx, false, ID)
}

// GetTx fetches a single article from the db inside of a tx from ctx
func (svc *articleService) GetTx(ctx context.Context, ID uuid.UUID) (*Article, error) {
	return svc.get(ctx, true, ID)
}

// get fetches a single article from the db
func (svc *articleService) get(ctx context.Context, useTx bool, ID uuid.UUID) (*Article, error) {
	errMsg := func() string { return "Error executing get article - " + fmt.Sprint(ID) }

	var (
		stmt *sql.Stmt
		err  error
		tx   *sql.Tx
	)

	if useTx {

		if tx, err = FromCtx(ctx); err != nil {
			return nil, err
		}

		stmt = tx.Stmt(svc.stmts["get-article"])
	} else {
		stmt = svc.stmts["get-article"]
	}

	p := Article{}

	err = stmt.QueryRowContext(ctx, ID).
		Scan(&m.ArticleID, &m.Name)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Wrap(ErrNotFound, errMsg())
		}

		return nil, errors.Wrap(err, errMsg())
	}

	return &p, nil
}

// Create a new article
func (svc *articleService) Create(ctx context.Context, input *Article) error {
	return svc.create(ctx, false, input)
}

// CreateTx creates a new article withing a tx from ctx
func (svc *articleService) CreateTx(ctx context.Context, input *Article) error {
	return svc.create(ctx, true, input)
}

// create a new article. if useTx = true then it will attempt to create the article within a transaction
// from context.
func (svc *articleService) create(ctx context.Context, useTx bool, input *Article) error {
	errMsg := func() string { return "Error executing create article - " + fmt.Sprint(input) }

	var (
		stmt *sql.Stmt
		err  error
		tx   *sql.Tx
	)

	if useTx {

		if tx, err = FromCtx(ctx); err != nil {
			return err
		}

		stmt = tx.Stmt(svc.stmts["create-article"])
	} else {
		stmt = svc.stmts["create-article"]
	}

	result, err := stmt.ExecContext(ctx, input.ArticleID, input.Name)
	if err != nil {
		return errors.Wrap(err, errMsg())
	}

	rowCount, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, errMsg())
	}

	if rowCount == 0 {
		return errors.Wrap(ErrNotCreated, errMsg())
	}

	return nil
}

// Update updates a single article row in the DB
func (svc *articleService) Update(ctx context.Context, input *Article) error {
	return svc.update(ctx, false, input)
}

// UpdateTx updates a single article row in the DB within a tx from ctx
func (svc *articleService) UpdateTx(ctx context.Context, input *Article) error {
	return svc.update(ctx, true, input)
}

// update a article. if useTx = true then it will attempt to update the article within a transaction
// from context.
func (svc *articleService) update(ctx context.Context, useTx bool, input *Article) error {
	errMsg := func() string { return "Error executing update article - " + fmt.Sprint(input) }

	var (
		stmt *sql.Stmt
		err  error
		tx   *sql.Tx
	)

	if useTx {

		if tx, err = FromCtx(ctx); err != nil {
			return err
		}

		stmt = tx.Stmt(svc.stmts["update-article"])
	} else {
		stmt = svc.stmts["update-article"]
	}

	result, err := stmt.ExecContext(ctx, input.Name, input.ArticleID)
	if err != nil {
		return errors.Wrap(err, errMsg())
	}

	rowCount, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, errMsg())
	}

	if rowCount == 0 {
		return errors.Wrap(ErrNoRowsAffected, errMsg())
	}

	return nil
}

// Delete sets deleted_at for a single articles row
func (svc *articleService) Delete(ctx context.Context, ID uuid.UUID) error {
	return svc.delete(ctx, false, ID)
}

// DeleteTx sets deleted_at for a single articles row within a tx from ctx
func (svc *articleService) DeleteTx(ctx context.Context, ID uuid.UUID) error {
	return svc.delete(ctx, true, ID)
}

// delete a article by setting deleted at. if useTx = true then it will attempt to delete the article within a transaction
// from context.
func (svc *articleService) delete(ctx context.Context, useTx bool, ID uuid.UUID) error {
	errMsg := func() string { return "Error executing delete article - " + ID.String() }

	var (
		stmt *sql.Stmt
		err  error
		tx   *sql.Tx
	)

	if useTx {

		if tx, err = FromCtx(ctx); err != nil {
			return err
		}

		stmt = tx.Stmt(svc.stmts["delete-article"])
	} else {
		stmt = svc.stmts["delete-article"]
	}

	result, err := stmt.ExecContext(ctx, ID)
	if err != nil {
		return errors.Wrap(err, errMsg())
	}

	rowCount, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, errMsg())
	}

	if rowCount == 0 {
		return errors.Wrap(ErrNotFound, errMsg())
	}

	return nil
}

