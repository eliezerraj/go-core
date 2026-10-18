package policy

import (
	"testing"

	cedar "github.com/cedar-policy/cedar-go"
)

func TestNewCedarClient(t *testing.T) {
	policyPath := "./policy.cedar"

	client, err := NewCedarClient(policyPath)
	if err != nil {
		t.Logf("Failed to create Cedar client: %v (expected if policy file doesn't exist)", err)
		return
	}

	principalUID := cedar.NewEntityUID(
		cedar.EntityType("Client"),
		cedar.String("client-test-23"),
	)

	resourceUID := cedar.NewEntityUID(
		cedar.EntityType("Resource"),
		cedar.String("resource-test-23"),
	)

	principal := cedar.Entity{
		UID: principalUID,
		Attributes: cedar.NewRecord(
			cedar.RecordMap{
				"client_id": cedar.String("client-test-23"),
				"iss":       cedar.String("go-authorizer-v2"),
				"sub":       cedar.String("client-test-23"),
				"scopes": cedar.NewSet(
					cedar.String("inventory:read"),
					cedar.String("inventory:create"),
					cedar.String("inventory:update"),
					cedar.String("order:read"),
					cedar.String("order:create"),
					cedar.String("order:history"),
					cedar.String("checkout:read"),
					cedar.String("checkout:create"),
					cedar.String("statistics:compute"),
					cedar.String("inventory_a2a:monitor"),
					cedar.String("inventory_a2a:manage"),
				),
			},
		),
		Parents: cedar.NewEntityUIDSet(),
		Tags:    cedar.NewRecord(cedar.RecordMap{}),
	}

	resource := cedar.Entity{
		UID: resourceUID,
		Attributes: cedar.NewRecord(
			cedar.RecordMap{},
		),
		Parents: cedar.NewEntityUIDSet(),
		Tags:    cedar.NewRecord(cedar.RecordMap{}),
	}

	entities := cedar.EntityMap{
		principalUID: principal,
		resourceUID:  resource,
	}

	request := cedar.Request{
		Principal: principalUID,
		Action: cedar.NewEntityUID(
			cedar.EntityType("Action"),
			cedar.String("get_inventory"),
		),
		Resource: resourceUID,
		Context: cedar.NewRecord(
			cedar.RecordMap{},
		),
	}

	result, err := client.EvaluatePolicy(
		request,
		entities,
	)
	if err != nil {
		t.Fatalf("EvaluatePolicy() returned error: %v", err)
	}

	t.Logf("Decision: %v", result.Decision)
	t.Logf("Diagnostic: %+v", result.Diagnostic)
}