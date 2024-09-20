package repository

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/abdelrahman146/digital-wallet/internal/model"
	"github.com/abdelrahman146/digital-wallet/internal/resource"
	"github.com/abdelrahman146/digital-wallet/pkg/api"
	"github.com/abdelrahman146/digital-wallet/pkg/logger"
	"gorm.io/gorm"
)

type TriggerRepo interface {
	// CreateTrigger creates a new trigger
	CreateTrigger(ctx context.Context, trigger *model.Trigger) error
	// UpdateTrigger updates an existing trigger
	UpdateTrigger(ctx context.Context, trigger *model.Trigger) error
	// DeleteTrigger deletes a trigger
	DeleteTrigger(ctx context.Context, trigger *model.Trigger) error
	// FetchTriggerByID retrieves a trigger by its ID
	FetchTriggerByID(ctx context.Context, id uint64) (*model.Trigger, error)
	// FetchTriggerBySlug retrieves a trigger by its slug
	FetchTriggerBySlug(ctx context.Context, slug string) (*model.Trigger, error)
	// FetchTriggers retrieves a paginated list of triggers
	FetchTriggers(ctx context.Context, page int, limit int) ([]model.Trigger, error)
	// CountTriggers retrieves the total number of triggers
	CountTriggers(ctx context.Context) (int64, error)
}

type triggerRepo struct {
	resources *resource.Resources
}

// NewTriggerRepo initializes the trigger repository
func NewTriggerRepo(resources *resource.Resources) TriggerRepo {
	return &triggerRepo{resources: resources}
}

// CreateTrigger creates a new trigger in the database
func (r *triggerRepo) CreateTrigger(ctx context.Context, trigger *model.Trigger) error {
	if err := r.resources.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(trigger).Error; err != nil {
			return err
		}
		topicDetail := sarama.TopicDetail{
			NumPartitions:     1,
			ReplicationFactor: 1,
		}
		if err := r.resources.Broker.CreateTopic(ctx, trigger.Slug, topicDetail); err != nil {
			return err
		}
		return nil
	}); err != nil {
		api.GetLogger(ctx).Error("Failed to create trigger", logger.Field("error", err), logger.Field("trigger", trigger))
		return err
	}

	return nil
}

// UpdateTrigger updates an existing trigger in the database
func (r *triggerRepo) UpdateTrigger(ctx context.Context, trigger *model.Trigger) error {
	// TODO: if slug is changed then update the topic name in broker
	if err := r.resources.DB.Save(trigger).Error; err != nil {
		api.GetLogger(ctx).Error("Failed to update trigger", logger.Field("error", err), logger.Field("trigger", trigger))
		return err
	}
	return nil
}

// DeleteTrigger deletes a trigger from the database
func (r *triggerRepo) DeleteTrigger(ctx context.Context, trigger *model.Trigger) error {
	if err := r.resources.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(trigger).Error; err != nil {
			return err
		}
		if err := r.resources.Broker.DeleteTopic(ctx, trigger.Slug); err != nil {
			return err
		}
		return nil
	}); err != nil {
		api.GetLogger(ctx).Error("Failed to delete trigger", logger.Field("error", err), logger.Field("trigger", trigger))
		return err
	}
	return nil
}

// FetchTriggerByID retrieves a trigger by its ID
func (r *triggerRepo) FetchTriggerByID(ctx context.Context, id uint64) (*model.Trigger, error) {
	var trigger model.Trigger
	err := r.resources.DB.Where("id = ?", id).First(&trigger).Error
	if err != nil {
		api.GetLogger(ctx).Error("Failed to retrieve trigger by ID", logger.Field("error", err), logger.Field("triggerId", id))
		return nil, err
	}
	return &trigger, nil
}

// FetchTriggerBySlug retrieves a trigger by its slug
func (r *triggerRepo) FetchTriggerBySlug(ctx context.Context, slug string) (*model.Trigger, error) {
	var trigger model.Trigger
	err := r.resources.DB.Where("slug = ?", slug).First(&trigger).Error
	if err != nil {
		api.GetLogger(ctx).Error("Failed to retrieve trigger by slug", logger.Field("error", err), logger.Field("triggerSlug", slug))
		return nil, err
	}
	return &trigger, nil
}

// FetchTriggers retrieves a paginated list of triggers
func (r *triggerRepo) FetchTriggers(ctx context.Context, page int, limit int) ([]model.Trigger, error) {
	var triggers []model.Trigger
	err := r.resources.DB.Order("created_at desc").Offset((page - 1) * limit).Limit(limit).Find(&triggers).Error
	if err != nil {
		api.GetLogger(ctx).Error("Failed to retrieve triggers", logger.Field("error", err))
		return nil, err
	}
	return triggers, nil
}

// CountTriggers retrieves the total number of triggers
func (r *triggerRepo) CountTriggers(ctx context.Context) (int64, error) {
	var total int64
	err := r.resources.DB.Model(&model.Trigger{}).Count(&total).Error
	if err != nil {
		api.GetLogger(ctx).Error("Failed to retrieve total triggers count", logger.Field("error", err))
		return 0, err
	}
	return total, nil
}
