package topic

import (
	"database/sql"
	db "main/infrastructure/database"
	"main/model"
)

type topicRepository struct {
	*db.Database
}

func NewTopicRepository(db *db.Database) TopicRepository {
	return &topicRepository{db}
}

func (r *topicRepository) Insert(topic string, id int64) error {
	err := r.Execute("usp_RepositoryTopic_Insert",
		sql.Named("Topic", topic),
		sql.Named("RepositoryId", id),
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *topicRepository) Delete(id int64) error {
	err := r.Execute("[dbo].[usp_RepositoryTopic_Delete]",
		sql.Named("RepositoryId", id))
	if err != nil {
		return err
	}
	return nil
}

func (r *topicRepository) SelectByOption(opt *model.FilterOptions) ([]model.Topic, error) {
	var topics []model.Topic
	rows, err := r.Query("[dbo].[usp_RepositoryTopic_Select_PopularTopic]",
		sql.Named("Offset", opt.Offset),
		sql.Named("RowCount", opt.Filter),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	mapRows, err := r.RowsToMap(rows)
	if err != nil {
		return nil, err
	}

	for _, v := range mapRows {
		var topic model.Topic
		topic.Topic = v["Topic"].(string)
		topic.Total = v["Total"].(int64)

		topics = append(topics, topic)
	}

	return topics, nil
}
