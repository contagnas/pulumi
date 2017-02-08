// Copyright 2016 Marapongo, Inc. All rights reserved.

package binder

import (
	"github.com/marapongo/mu/pkg/compiler/core"
	"github.com/marapongo/mu/pkg/compiler/metadata"
	"github.com/marapongo/mu/pkg/compiler/symbols"
	"github.com/marapongo/mu/pkg/diag"
	"github.com/marapongo/mu/pkg/pack"
	"github.com/marapongo/mu/pkg/workspace"
)

// Binder annotates an existing parse tree with semantic information.
type Binder interface {
	core.Phase

	// Ctx represents the contextual information resulting from binding.
	Ctx() *Context

	// BindPackages takes a package AST, resolves all dependencies and tokens inside of it, and returns a fully bound
	// package symbol that can be used for semantic operations (like interpretation and evaluation).
	BindPackage(pkg *pack.Package) *symbols.Package
}

// New allocates a fresh binder object with the given workspace, context, and metadata reader.
func New(w workspace.W, ctx *core.Context, reader metadata.Reader) Binder {
	// Create a new binder with a fresh binding context.
	return &binder{
		w:      w,
		ctx:    NewContextFrom(ctx),
		reader: reader,
	}
}

type binder struct {
	w      workspace.W     // a workspace in which this compilation is happening.
	ctx    *Context        // a binding context shared with future phases of compilation.
	reader metadata.Reader // a metadata reader (in case we encounter package references).
}

func (b *binder) Ctx() *Context   { return b.ctx }
func (b *binder) Diag() diag.Sink { return b.ctx.Diag }
