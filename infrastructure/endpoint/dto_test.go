package endpoint

import (
	"testing"

	"github.com/google/uuid"
)

func TestPaginateEndpointQueryTagIDs(t *testing.T) {
	tag1 := uuid.MustParse("81a1fa69-12ba-48ac-84c1-feb71bcae113")
	tag2 := uuid.MustParse("1d88761f-ab43-468c-94e5-299b4adcd0eb")

	query := PaginateEndpointQuery{
		Tags: "81a1fa69-12ba-48ac-84c1-feb71bcae113,1d88761f-ab43-468c-94e5-299b4adcd0eb",
	}

	tagIDs := query.TagIDs()
	if len(tagIDs) != 2 {
		t.Fatalf("expected 2 tag ids, got %d", len(tagIDs))
	}

	if tagIDs[0] != tag1 || tagIDs[1] != tag2 {
		t.Fatalf("unexpected tag ids: %v", tagIDs)
	}
}
