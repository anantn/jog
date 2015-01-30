#ifdef __cplusplus
extern "C" {
#endif

// A path is an array of keys + length.
typedef struct Path {
	char** keys;
	size_t length;
} Path;

// If the return value is NULL, an error message will be stored in *error.
// The caller must free the error message if present.
void* NewDocument(char* string, char** error);
void  DeleteDocument(void* value);

// Return the child value at given path. If the path is NULL, the provided
// value is returned as-is.
void*        Get(void* value, Path* path);

// Get data at a given path. errno will be set if there's an error.
// If path is NULL, it will return the value at the current path.
int          GetInt(void* value, Path* path);
bool         GetBool(void* value, Path* path);
double       GetDouble(void* value, Path* path);
unsigned int GetUInt(void* value, Path* path);

// Don't free the return the value. NULL will be returned if there's an error.
const char*  GetString(void* value, Path* path);

// Caller must free the array of returned pointers.
void**       GetArray(void* value, Path* path, size_t* length);

// Caller must free the array of returned pointers, the keys array and each member of the array.
void**       GetObject(void* value, Path* path, size_t* length, char*** keys);

// The type of a value is either "BOOL", "NULL", "ARRAY",
// "STRING", "OBJECT", or "NUMBER".
const char*  Type(void* value, Path* path);

#ifdef __cplusplus
}
#endif
