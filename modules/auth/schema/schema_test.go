package schema

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/modules/crud"

	"github.com/spaceuptech/space-cloud/config"
)

var query = `
type Tweet {
	id: ID! @id
	createdAt: DateTime! @createdAt
	text: String
	owner: [String!]! @relation(field: INLINE)
	location: location!
	age : Integer!	
  }
  
  type User {
	id: ID! @id
	createdAt: DateTime! @createdAt
	updatedAt: DateTime! @updatedAt
	handle: String! @unique
	name: String
	tweets: [tweet!]!
  }
  
  type Location {
	id: ID! @id
	latitude: Float!
	longitude: Float!
	person : sharad!
  }

  type Sharad {
	  name : String!
	  sirName : String!
	  age : Integer!
	  isMale : Boolean!
	  dob : DateTime!
  }
`
var ParseData = config.Crud{
	"mongo": &config.CrudStub{
		Collections: map[string]*config.TableRule{
			"tweet": &config.TableRule{
				Schema: query,
			},
			"user": &config.TableRule{
				Schema: query,
			},
			"location": &config.TableRule{
				Schema: query,
			},
		},
	},
}

func TestParseSchema(t *testing.T) {
	temp := crud.Module{}
	s := Init(&temp)

	t.Run("Schema Parser", func(t *testing.T) {
		err := s.ParseSchema(ParseData)
		if err != nil {
			t.Fatal(err)
		}
		// uncomment the below statements to see the reuslt
		b, err := json.MarshalIndent(s.schemaDoc, "", "  ")
		if err != nil {
			fmt.Println("error:", err)
		}
		fmt.Print(string(b))
		t.Log("Logging Test Output :: ", s.schemaDoc)
	})
}

func TestValidateSchema(t *testing.T) {

	var arr []interface{}
	str := []int{1, 2, 3}
	for _, v := range str {
		arr = append(arr, v)
	}

	req := model.CreateRequest{
		Document: []map[string]interface{}{
			{
				"id":        "dfdsairfa",
				"createdAt": 986413662654,
				"text":      "Hello World!",
				"location": map[string]interface{}{
					"id":        "locatoinid",
					"latitude":  5.5,
					"longitude": 312.3,
					"person": map[string]interface{}{
						"name":    "sharad",
						"sirName": "Regoti",
						"age":     19,
						"isMale":  true,
						"dob":     "1999-10-19T11:45:26.371Z",
					},
				},
				"owner": arr,
			},
		},
	}

	tdd := []struct {
		dbName, coll, description string
		value                     model.CreateRequest
	}{{
		dbName:      "mongo",
		coll:        "tweet",
		description: "checking User defined type",
		value:       req,
	}}
	temp := crud.Module{}
	s := Init(&temp)
	err := s.ParseSchema(ParseData)
	if err != nil {
		t.Fatal(err)
	}

	for _, val := range tdd {
		t.Run(val.description, func(t *testing.T) {
			err := s.ValidateCreateOperation(val.dbName, val.coll, &val.value)
			if err != nil {
				t.Fatal(err)
			}
		})
	}

}
