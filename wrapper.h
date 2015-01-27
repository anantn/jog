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

// Get data at a given path. errno will be set if there's an error.
// If path is NULL, it will return the value at the current path.
int          GetInt(void* value, Path* path);
int          GetBool(void* value, Path* path); // 0 for FALSE, 1 for TRUE.
double       GetDouble(void* value, Path* path);
unsigned int GetUInt(void* value, Path* path);

// Don't free the return the value. NULL will be returned if there's an error.
const char* GetString(void* value, Path* path);
void*       GetObject(void* value, Path* path);

// Caller must free the array of pointers.
void**      GetArray(void* value, Path* path, size_t* length);

#ifdef __cplusplus
}
#endif
