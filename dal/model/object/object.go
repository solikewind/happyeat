package object

import (
	"context"

	"github.com/solikewind/happyeat/dal/model/ent"
	entobj "github.com/solikewind/happyeat/dal/model/ent/object"
)

type Object struct {
	c *ent.Client
}

func NewObject(c *ent.Client) *Object {
	return &Object{c: c}
}

type CreateInput struct {
	Name        string
	Key         string
	URL         string
	ContentType string
	Suffix      string
	Size        int64
	Hash        string
}

func (o *Object) GetByID(ctx context.Context, id uint64) (*ent.Object, error) {
	return o.c.Object.Get(ctx, id)
}

func (o *Object) GetByHash(ctx context.Context, hash string) (*ent.Object, error) {
	return o.c.Object.Query().Where(entobj.HashEQ(hash)).First(ctx)
}

func (o *Object) Create(ctx context.Context, in CreateInput) (*ent.Object, error) {
	create := o.c.Object.Create().
		SetName(in.Name).
		SetKey(in.Key).
		SetURL(in.URL).
		SetSize(in.Size).
		SetHash(in.Hash)
	if in.ContentType != "" {
		create = create.SetContentType(in.ContentType)
	}
	if in.Suffix != "" {
		create = create.SetSuffix(in.Suffix)
	}
	return create.Save(ctx)
}

func (o *Object) DeleteByID(ctx context.Context, id uint64) error {
	return o.c.Object.DeleteOneID(id).Exec(ctx)
}
