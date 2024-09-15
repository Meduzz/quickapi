package rpc

import (
	"github.com/Meduzz/helper/http/herror"
	"github.com/Meduzz/quickapi/model"
	"github.com/Meduzz/quickapi/storage"
	"github.com/Meduzz/rpc"
	"github.com/Meduzz/rpc/messages"
	"gorm.io/gorm"
)

type (
	router struct {
		entity model.Entity
		storer storage.Storer
	}
)

func newRouter(db *gorm.DB, entity model.Entity) *router {
	storer := storage.NewStorer(db, entity)
	return &router{entity, storer}
}

func (r *router) Create(ctx *rpc.RpcContext) {
	entity := r.entity.Create()
	err := ctx.Bind(entity)

	if err != nil {
		ctx.ReplyBuilder(func(mb *messages.MsgBuilder) {
			mb.WithProblem(400, err.Error())
		})
	}

	entity, err = r.storer.Create(entity)

	if err != nil {
		code := herror.CodeFromError(err)
		ctx.ReplyBuilder(func(mb *messages.MsgBuilder) {
			mb.WithProblem(code, err.Error())
		})
	}

	ctx.ReplyBuilder(func(mb *messages.MsgBuilder) {
		mb.WithBody(entity)
		mb.WithCode(200)
	})
}

func (r *router) Read(ctx *rpc.RpcContext) {
	req := &ReadRequest{}
	err := ctx.Bind(req)

	if err != nil {
		ctx.ReplyBuilder(func(mb *messages.MsgBuilder) {
			mb.WithProblem(400, err.Error())
		})
	}

	data, err := r.storer.Read(req.ID)

	if err != nil {
		ctx.ReplyBuilder(func(mb *messages.MsgBuilder) {
			code := herror.CodeFromError(err)
			mb.WithProblem(code, err.Error())
		})
	}

	ctx.ReplyBuilder(func(mb *messages.MsgBuilder) {
		mb.WithBody(data)
		mb.WithCode(200)
	})
}

func (r *router) Update(ctx *rpc.RpcContext) {
	data := r.entity.Create()
	err := ctx.Bind(data)

	if err != nil {
		ctx.ReplyBuilder(func(mb *messages.MsgBuilder) {
			mb.WithProblem(400, err.Error())
		})
	}

	data, err = r.storer.Update(data)

	if err != nil {
		ctx.ReplyBuilder(func(mb *messages.MsgBuilder) {
			code := herror.CodeFromError(err)
			mb.WithProblem(code, err.Error())
		})
	}

	ctx.ReplyBuilder(func(mb *messages.MsgBuilder) {
		mb.WithBody(data)
		mb.WithCode(200)
	})
}

func (r *router) Delete(ctx *rpc.RpcContext) {
	req := &DeleteRequest{}
	err := ctx.Bind(req)

	if err != nil {
		ctx.ReplyBuilder(func(mb *messages.MsgBuilder) {
			mb.WithProblem(400, err.Error())
		})
	}

	err = r.storer.Delete(req.ID)

	if err != nil {
		ctx.ReplyBuilder(func(mb *messages.MsgBuilder) {
			code := herror.CodeFromError(err)
			mb.WithProblem(code, err.Error())
		})
	}

	ctx.ReplyBuilder(func(mb *messages.MsgBuilder) {
		mb.WithCode(200)
	})
}

func (r *router) Search(ctx *rpc.RpcContext) {
	req := &SearchRequest{}
	err := ctx.Bind(req)

	if err != nil {
		ctx.ReplyBuilder(func(mb *messages.MsgBuilder) {
			mb.WithProblem(400, err.Error())
		})
	}

	hooks := createScopes(req.Scopes, r.entity.Filters())

	data, err := r.storer.Search(req.Skip, req.Take, req.Where, hooks...)

	if err != nil {
		ctx.ReplyBuilder(func(mb *messages.MsgBuilder) {
			code := herror.CodeFromError(err)
			mb.WithProblem(code, err.Error())
		})
	}

	ctx.ReplyBuilder(func(mb *messages.MsgBuilder) {
		mb.WithBody(data)
		mb.WithCode(200)
	})
}

func (r *router) Patch(ctx *rpc.RpcContext) {
	req := &PatchRequest{}
	err := ctx.Bind(req)

	if err != nil {
		ctx.ReplyBuilder(func(mb *messages.MsgBuilder) {
			mb.WithProblem(400, err.Error())
		})
	}

	data, err := r.storer.Patch(req.ID, req.Data)

	if err != nil {
		ctx.ReplyBuilder(func(mb *messages.MsgBuilder) {
			code := herror.CodeFromError(err)
			mb.WithProblem(code, err.Error())
		})
	}

	ctx.ReplyBuilder(func(mb *messages.MsgBuilder) {
		mb.WithBody(data)
		mb.WithCode(200)
	})
}
