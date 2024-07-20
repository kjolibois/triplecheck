package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models/schema"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("p7s5rowbx5o7ksg")
		if err != nil {
			return err
		}

		// add
		new_thought_process := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "4hbe5bag",
			"name": "thought_process",
			"type": "text",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), new_thought_process); err != nil {
			return err
		}
		collection.Schema.AddField(new_thought_process)

		// add
		new_asker_user_id := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "eqxddtyh",
			"name": "asker_user_id",
			"type": "relation",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"collectionId": "_pb_users_auth_",
				"cascadeDelete": false,
				"minSelect": null,
				"maxSelect": 1,
				"displayFields": null
			}
		}`), new_asker_user_id); err != nil {
			return err
		}
		collection.Schema.AddField(new_asker_user_id)

		// add
		new_reviewers := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "xaamtl3o",
			"name": "reviewers",
			"type": "relation",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"collectionId": "_pb_users_auth_",
				"cascadeDelete": false,
				"minSelect": null,
				"maxSelect": null,
				"displayFields": null
			}
		}`), new_reviewers); err != nil {
			return err
		}
		collection.Schema.AddField(new_reviewers)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("p7s5rowbx5o7ksg")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("4hbe5bag")

		// remove
		collection.Schema.RemoveField("eqxddtyh")

		// remove
		collection.Schema.RemoveField("xaamtl3o")

		return dao.SaveCollection(collection)
	})
}
