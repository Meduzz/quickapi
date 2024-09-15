package rpc

import (
	"fmt"

	"github.com/Meduzz/quickapi/model"
	"github.com/Meduzz/rpc"
	"gorm.io/gorm"
)

func For(db *gorm.DB, srv *rpc.RPC, prefix, queue string, entity model.Entity) {
	topic := entity.Name()

	if prefix != "" && topic != "" {
		topic = fmt.Sprintf("%s.%s", prefix, topic)
	}

	router := newRouter(db, entity)

	srv.HandleRPC(topicify(topic, "create"), "", router.Create)
	srv.HandleRPC(topicify(topic, "read"), "", router.Read)
	srv.HandleRPC(topicify(topic, "update"), "", router.Update)
	srv.HandleRPC(topicify(topic, "delete"), "", router.Delete)
	srv.HandleRPC(topicify(topic, "search"), "", router.Search)
	srv.HandleRPC(topicify(topic, "patch"), "", router.Patch)
}

func topicify(prefix, action string) string {
	if prefix != "" {
		return fmt.Sprintf("%s.%s", prefix, action)
	}

	return action
}

func createScopes(provided map[string]map[string]string, filters []*model.NamedFilter) []model.Hook {
	if len(filters) == 0 {
		return nil
	}

	scopes := []model.Hook{}

	for _, filter := range filters {
		data, ok := provided[filter.Name]

		if ok {
			scopes = append(scopes, filter.Scope(data))
		}
	}

	return scopes
}
