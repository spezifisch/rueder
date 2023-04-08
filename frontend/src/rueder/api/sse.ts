// SSE Message Types
// these need to be in sync with backend/internal/common/sse_types.go
export enum SSEMessageType {
    Raw = "raw", // special internal type, doesn't exist in backend
    FolderUpdate = "folder_update", // folder/feed list order or content was changed
}
