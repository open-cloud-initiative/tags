package controllers

import (
	"github.com/open-cloud-initiative/tags/internal/ports"

	"github.com/katallaxie/pkg/dbx"
	pb "github.com/open-cloud-initiative/specs/gen/go/tags/v1"
)

// TagsController is the controller for managing tags.
type TagsController struct {
	store dbx.Database[ports.ReadTx, ports.ReadWriteTx]
	pb.UnimplementedTagsServiceServer
}

// NewTagsController creates a new TagsController.
func NewTagsController(store dbx.Database[ports.ReadTx, ports.ReadWriteTx]) *TagsController {
	return &TagsController{
		store: store,
	}
}
