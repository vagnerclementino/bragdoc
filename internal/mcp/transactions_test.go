package mcp

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vagnerclementino/bragdoc/internal/domain"
)

func callTool(t *testing.T, ctx context.Context, client *mcp.ClientSession, name string, arguments any) *mcp.CallToolResult {
	t.Helper()
	result, err := client.CallTool(ctx, &mcp.CallToolParams{Name: name, Arguments: arguments})
	require.NoError(t, err)
	return result
}

func createTestBrag(t *testing.T, ctx context.Context, client *mcp.ClientSession, userID int64, tags ...string) BragResponse {
	t.Helper()
	result := callTool(t, ctx, client, "brag_create", BragCreateParams{
		UserID: userID, Title: "Original achievement", Description: "Original description with enough detail to pass validation", Category: "PROJECT", Tags: tags,
	})
	require.False(t, result.IsError, "%s", extractText(result))
	var brag BragResponse
	require.NoError(t, unmarshalResult(result, &brag))
	return brag
}

func createTestUser(t *testing.T, ctx context.Context, is *integrationServer) int64 {
	t.Helper()
	user, err := is.userService.Create(ctx, &domain.User{Name: "Transaction Test", Email: "transaction@example.com"})
	require.NoError(t, err)
	return user.ID
}

func assertRowCounts(t *testing.T, ctx context.Context, is *integrationServer, brags, tags, links int) {
	t.Helper()
	for table, expected := range map[string]int{"brags": brags, "tags": tags, "brag_tags": links} {
		var count int
		require.NoError(t, is.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM "+table).Scan(&count))
		assert.Equal(t, expected, count, "unexpected rows in %s", table)
	}
}

func TestIntegration_CreateWithTagsRollsBack(t *testing.T) {
	for _, failure := range []string{"validation", "attachment"} {
		t.Run(failure, func(t *testing.T) {
			is := newIntegrationServer(t)
			ctx, client := newIntegrationClient(t, is.server)
			userID := createTestUser(t, ctx, is)
			tags := []string{"first", "x"}
			code := ErrCodeValidation
			if failure == "attachment" {
				// Fail after the brag, both tags, and the first association are written.
				_, err := is.db.ExecContext(ctx, `CREATE TRIGGER reject_attachment BEFORE INSERT ON brag_tags
WHEN (SELECT name FROM tags WHERE id = NEW.tag_id) = 'blocked'
BEGIN SELECT RAISE(ABORT, 'injected attachment failure'); END`)
				require.NoError(t, err)
				tags = []string{"first", "blocked"}
				code = ErrCodeInternal
			}
			result := callTool(t, ctx, client, "brag_create", BragCreateParams{
				UserID: userID, Title: "Atomic achievement", Description: "A complete achievement that must roll back when tagging fails", Category: "PROJECT", Tags: tags,
			})
			requireErrorEnvelope(t, result, code)
			assertRowCounts(t, ctx, is, 0, 0, 0)
			// A corrected retry produces exactly one complete achievement.
			brag := createTestBrag(t, ctx, client, userID, "first", "second")
			assert.ElementsMatch(t, []string{"first", "second"}, brag.Tags)
			assertRowCounts(t, ctx, is, 1, 2, 2)
		})
	}
}

func TestIntegration_UpdateWithTagsRollsBack(t *testing.T) {
	is := newIntegrationServer(t)
	ctx, client := newIntegrationClient(t, is.server)
	userID := createTestUser(t, ctx, is)
	original := createTestBrag(t, ctx, client, userID, "original")
	result := callTool(t, ctx, client, "brag_update", BragUpdateParams{
		ID: original.ID, Title: "Changed achievement", Description: "Changed description that should also be rolled back", Category: "SKILL", Tags: []string{"new-tag", "x"},
	})
	requireErrorEnvelope(t, result, ErrCodeValidation)
	got := callTool(t, ctx, client, "brag_get", BragGetParams{ID: original.ID})
	require.False(t, got.IsError)
	var after BragResponse
	require.NoError(t, unmarshalResult(got, &after))
	assert.Equal(t, original, after)
	assertRowCounts(t, ctx, is, 1, 1, 1)
}

func TestIntegration_UpdateTagReplacement(t *testing.T) {
	is := newIntegrationServer(t)
	ctx, client := newIntegrationClient(t, is.server)
	userID := createTestUser(t, ctx, is)
	original := createTestBrag(t, ctx, client, userID, "original")
	for _, test := range []struct {
		name      string
		arguments map[string]any
		tags      []string
	}{
		{"omitted tags are retained", map[string]any{"id": original.ID, "title": "Updated achievement"}, []string{"original"}},
		{"provided tags replace associations", map[string]any{"id": original.ID, "tags": []string{"replacement"}}, []string{"replacement"}},
		{"empty tags clear associations", map[string]any{"id": original.ID, "tags": []string{}}, []string{}},
	} {
		t.Run(test.name, func(t *testing.T) {
			result := callTool(t, ctx, client, "brag_update", test.arguments)
			require.False(t, result.IsError, "%s", extractText(result))
			var brag BragResponse
			require.NoError(t, unmarshalResult(result, &brag))
			assert.ElementsMatch(t, test.tags, brag.Tags)
			got := callTool(t, ctx, client, "brag_get", BragGetParams{ID: original.ID})
			assert.JSONEq(t, extractText(result), extractText(got))
		})
	}
	assertRowCounts(t, ctx, is, 1, 2, 0)
}

func TestIntegration_TagBatchWritesRollBack(t *testing.T) {
	const attach = "attach"
	for _, operation := range []string{attach, "detach", "delete"} {
		t.Run(operation, func(t *testing.T) {
			is := newIntegrationServer(t)
			ctx, client := newIntegrationClient(t, is.server)
			userID := createTestUser(t, ctx, is)
			first, err := is.tagService.Create(ctx, &domain.Tag{OwnerID: userID, Name: "first"})
			require.NoError(t, err)
			second, err := is.tagService.Create(ctx, &domain.Tag{OwnerID: userID, Name: "second"})
			require.NoError(t, err)
			brag := createTestBrag(t, ctx, client, userID)
			if operation != attach {
				result := callTool(t, ctx, client, "tag_attach", TagAttachParams{BragID: brag.ID, TagIDs: []int64{first.ID, second.ID}})
				require.False(t, result.IsError)
			}
			var trigger string
			var arguments any
			switch operation {
			case attach:
				trigger = fmt.Sprintf(`CREATE TRIGGER reject_batch BEFORE INSERT ON brag_tags WHEN NEW.tag_id = %d BEGIN SELECT RAISE(ABORT, 'injected failure'); END`, second.ID)
				arguments = TagAttachParams{BragID: brag.ID, TagIDs: []int64{first.ID, second.ID}}
			case "detach":
				trigger = fmt.Sprintf(`CREATE TRIGGER reject_batch BEFORE DELETE ON brag_tags WHEN OLD.tag_id = %d BEGIN SELECT RAISE(ABORT, 'injected failure'); END`, second.ID)
				arguments = TagDetachParams{BragID: brag.ID, TagIDs: []int64{first.ID, second.ID}}
			case "delete":
				trigger = fmt.Sprintf(`CREATE TRIGGER reject_batch BEFORE DELETE ON tags WHEN OLD.id = %d BEGIN SELECT RAISE(ABORT, 'injected failure'); END`, first.ID)
				arguments = TagDeleteParams{ID: first.ID}
			}
			_, err = is.db.ExecContext(ctx, trigger)
			require.NoError(t, err)
			result := callTool(t, ctx, client, "tag_"+operation, arguments)
			requireErrorEnvelope(t, result, ErrCodeInternal)
			if operation == attach {
				assertRowCounts(t, ctx, is, 1, 2, 0)
			} else {
				assertRowCounts(t, ctx, is, 1, 2, 2)
			}
		})
	}
}

func TestIntegration_ConcurrentWritesKeepTransactionsSeparate(t *testing.T) {
	is := newIntegrationServer(t)
	ctx, client := newIntegrationClient(t, is.server)
	userID := createTestUser(t, ctx, is)
	var results [8]*mcp.CallToolResult
	var failures [8]error
	var workers sync.WaitGroup
	for i := range results {
		workers.Add(1)
		go func() {
			defer workers.Done()
			tags := []string{"shared", fmt.Sprintf("tag-%d", i)}
			if i%2 == 0 {
				tags = append(tags, "x") // Fail after earlier writes in this request.
			}
			results[i], failures[i] = client.CallTool(ctx, &mcp.CallToolParams{Name: "brag_create", Arguments: BragCreateParams{
				UserID: userID, Title: fmt.Sprintf("Concurrent achievement %d", i), Description: "A concurrently recorded achievement with multiple tags", Category: "PROJECT", Tags: tags,
			}})
		}()
	}
	workers.Wait()
	for i, result := range results {
		require.NoError(t, failures[i])
		if i%2 == 0 {
			requireErrorEnvelope(t, result, ErrCodeValidation)
		} else {
			require.False(t, result.IsError, "%s", extractText(result))
		}
	}
	assertRowCounts(t, ctx, is, 4, 5, 8)
}
