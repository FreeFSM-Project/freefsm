package services

import (
	"context"
	"fmt"

	"github.com/MartialM1nd/freefsm/internal/ent"
	"github.com/MartialM1nd/freefsm/internal/ent/tag"
	"github.com/MartialM1nd/freefsm/internal/ent/taglink"
)

type TagLinkService struct {
	client *ent.Client
}

func NewTagLinkService(client *ent.Client) *TagLinkService {
	return &TagLinkService{client: client}
}

func (s *TagLinkService) Attach(ctx context.Context, tagID int64, objectType string, objectID int64) (*ent.TagLink, error) {
	// Check if link already exists
	exists, err := s.client.TagLink.Query().
		Where(taglink.TagIDEQ(tagID), taglink.ObjectTypeEQ(objectType), taglink.ObjectIDEQ(objectID)).
		Exist(ctx)
	if err != nil {
		return nil, fmt.Errorf("check tag link: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("tag already attached")
	}

	l, err := s.client.TagLink.Create().
		SetTagID(tagID).
		SetObjectType(objectType).
		SetObjectID(objectID).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("attach tag: %w", err)
	}
	return l, nil
}

func (s *TagLinkService) Detach(ctx context.Context, tagID int64, objectType string, objectID int64) error {
	_, err := s.client.TagLink.Delete().
		Where(taglink.TagIDEQ(tagID), taglink.ObjectTypeEQ(objectType), taglink.ObjectIDEQ(objectID)).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("detach tag: %w", err)
	}
	return nil
}

func (s *TagLinkService) ListForObject(ctx context.Context, objectType string, objectID int64) ([]*ent.Tag, error) {
	links, err := s.client.TagLink.Query().
		Where(taglink.ObjectTypeEQ(objectType), taglink.ObjectIDEQ(objectID)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("list tags for object: %w", err)
	}

	if len(links) == 0 {
		return nil, nil
	}

	tagIDs := make([]int64, len(links))
	for i, l := range links {
		tagIDs[i] = l.TagID
	}

	tags, err := s.client.Tag.Query().
		Where(tag.IDIn(tagIDs...)).
		Order(ent.Asc(tag.FieldName)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetch tags: %w", err)
	}
	return tags, nil
}

func (s *TagLinkService) ListObjectsWithTag(ctx context.Context, tagID int64, objectType string) ([]*ent.TagLink, error) {
	return s.client.TagLink.Query().
		Where(taglink.TagIDEQ(tagID), taglink.ObjectTypeEQ(objectType)).
		All(ctx)
}
