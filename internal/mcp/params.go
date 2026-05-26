package mcp

// === Input Parameter Structs ===
// These structs define the JSON Schema for each tool's input.
// The MCP Go SDK auto-generates JSON Schema from struct tags.
// The "jsonschema" tag is used as the property description.
// Fields without "omitempty" in their json tag are automatically required.

// BragCreateParams defines input for the brag_create tool.
type BragCreateParams struct {
	UserID      int64    `json:"user_id" jsonschema:"ID of the brag owner"`
	Title       string   `json:"title" jsonschema:"Title of the achievement (min 5 chars)"`
	Description string   `json:"description" jsonschema:"Detailed description (min 20 chars)"`
	Category    string   `json:"category" jsonschema:"Achievement category: PROJECT, ACHIEVEMENT, SKILL, LEADERSHIP, or INNOVATION"`
	Tags        []string `json:"tags,omitempty" jsonschema:"Optional tag names to attach"`
}

// BragGetParams defines input for the brag_get tool.
type BragGetParams struct {
	ID int64 `json:"id" jsonschema:"Brag ID to retrieve"`
}

// BragListParams defines input for the brag_list tool.
type BragListParams struct {
	UserID int64 `json:"user_id" jsonschema:"User ID to list brags for"`
}

// BragSearchByTagsParams defines input for the brag_search_by_tags tool.
type BragSearchByTagsParams struct {
	UserID   int64    `json:"user_id" jsonschema:"User ID"`
	TagNames []string `json:"tag_names" jsonschema:"Tag names to search for (at least one required)"`
}

// BragSearchByCategoryParams defines input for the brag_search_by_category tool.
type BragSearchByCategoryParams struct {
	UserID   int64  `json:"user_id" jsonschema:"User ID"`
	Category string `json:"category" jsonschema:"Category to filter by: PROJECT, ACHIEVEMENT, SKILL, LEADERSHIP, or INNOVATION"`
}

// BragUpdateParams defines input for the brag_update tool.
type BragUpdateParams struct {
	ID          int64    `json:"id" jsonschema:"Brag ID to update"`
	Title       string   `json:"title,omitempty" jsonschema:"New title (min 5 chars)"`
	Description string   `json:"description,omitempty" jsonschema:"New description (min 20 chars)"`
	Category    string   `json:"category,omitempty" jsonschema:"New category: PROJECT, ACHIEVEMENT, SKILL, LEADERSHIP, or INNOVATION"`
	Tags        []string `json:"tags,omitempty" jsonschema:"New tags to replace existing"`
}

// BragDeleteParams defines input for the brag_delete tool.
type BragDeleteParams struct {
	ID int64 `json:"id" jsonschema:"Brag ID to delete"`
}

// TagCreateParams defines input for the tag_create tool.
type TagCreateParams struct {
	OwnerID int64  `json:"owner_id" jsonschema:"Owner user ID"`
	Name    string `json:"name" jsonschema:"Tag name (2-20 characters)"`
}

// TagListParams defines input for the tag_list tool.
type TagListParams struct {
	UserID int64 `json:"user_id" jsonschema:"User ID to list tags for"`
}

// TagAttachParams defines input for the tag_attach tool.
type TagAttachParams struct {
	BragID int64   `json:"brag_id" jsonschema:"Brag ID to attach tags to"`
	TagIDs []int64 `json:"tag_ids" jsonschema:"Tag IDs to attach (at least one required)"`
}

// TagDetachParams defines input for the tag_detach tool.
type TagDetachParams struct {
	BragID int64   `json:"brag_id" jsonschema:"Brag ID to detach tags from"`
	TagIDs []int64 `json:"tag_ids" jsonschema:"Tag IDs to detach (at least one required)"`
}

// TagDeleteParams defines input for the tag_delete tool.
type TagDeleteParams struct {
	ID int64 `json:"id" jsonschema:"Tag ID to delete"`
}

// TagGetOrCreateParams defines input for the tag_get_or_create tool.
type TagGetOrCreateParams struct {
	OwnerID int64  `json:"owner_id" jsonschema:"Owner user ID"`
	Name    string `json:"name" jsonschema:"Tag name (2-20 characters)"`
}

// DocGenerateParams defines input for the doc_generate tool.
type DocGenerateParams struct {
	UserID   int64  `json:"user_id" jsonschema:"User ID to generate document for"`
	Format   string `json:"format,omitempty" jsonschema:"Output format (only markdown supported)"`
	Template string `json:"template,omitempty" jsonschema:"Template to use (only default supported)"`
}

// UserListParams defines input for the user_list tool (no required params).
type UserListParams struct{}

// UserGetParams defines input for the user_get tool.
type UserGetParams struct {
	ID int64 `json:"id" jsonschema:"User ID"`
}

// UserGetByEmailParams defines input for the user_get_by_email tool.
type UserGetByEmailParams struct {
	Email string `json:"email" jsonschema:"User email address"`
}

// === Response Structs ===

// BragResponse represents a brag entry in tool responses.
type BragResponse struct {
	ID          int64    `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Category    string   `json:"category"`
	Tags        []string `json:"tags"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
}

// TagResponse represents a tag in tool responses.
type TagResponse struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	OwnerID   int64  `json:"owner_id"`
	CreatedAt string `json:"created_at"`
}

// UserResponse represents a user profile in tool responses.
type UserResponse struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	JobTitle  string `json:"job_title"`
	Company   string `json:"company"`
	Locale    string `json:"locale"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// DocResponse represents a generated document in tool responses.
type DocResponse struct {
	Content     string   `json:"content"`
	Title       string   `json:"title"`
	Author      string   `json:"author"`
	BragCount   int      `json:"brag_count"`
	Categories  []string `json:"categories"`
	Tags        []string `json:"tags"`
	GeneratedAt string   `json:"generated_at"`
}

// SuccessResponse represents a simple success/failure result.
type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
