#include <errno.h>
#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>

#include "rapid.h"
#include "writer.h"
#include "document.h"
#include "stringbuffer.h"

using namespace rapidjson;

void* NewDocument(char* string, char** error) {
    Document* doc = new Document();
    if (!string) {
        delete doc;
        return NULL;
    }

    if (doc->ParseInsitu(string).HasParseError()) {
        char* msg = new char[100];
        size_t offset = doc->GetErrorOffset();
        ParseErrorCode code = doc->GetParseError();
        switch (code) {
            case kParseErrorNone:
                sprintf(msg, "[%lu] %s", offset, "No error.");
                break;
            case kParseErrorDocumentEmpty:
                sprintf(msg, "[%lu] %s", offset, "The document is empty.");
                break;
            case kParseErrorDocumentRootNotSingular:
                sprintf(msg, "[%lu] %s", offset, "The document root must not follow by other values.");
                break;
            case kParseErrorValueInvalid:
                sprintf(msg, "[%lu] %s", offset, "Invalid value.");
                break;
            case kParseErrorObjectMissName:
                sprintf(msg, "[%lu] %s", offset, "Missing a name for object member.");
                break;
            case kParseErrorObjectMissColon:
                sprintf(msg, "[%lu] %s", offset, "Missing a colon after a name of object member.");
                break;
            case kParseErrorObjectMissCommaOrCurlyBracket:
                sprintf(msg, "[%lu] %s", offset, "Missing a comma or '}' after an object member.");
                break;
            case kParseErrorArrayMissCommaOrSquareBracket:
                sprintf(msg, "[%lu] %s", offset, "Missing a comma or ']' after an array element.");
                break;
            case kParseErrorStringUnicodeEscapeInvalidHex:
                sprintf(msg, "[%lu] %s", offset, "Incorrect hex digit after \\u escape in string.");
                break;
            case kParseErrorStringUnicodeSurrogateInvalid:
                sprintf(msg, "[%lu] %s", offset, "The surrogate pair in string is invalid.");
                break;
            case kParseErrorStringEscapeInvalid:
                sprintf(msg, "[%lu] %s", offset, "Invalid escape character in string.");
                break;
            case kParseErrorStringMissQuotationMark:
                sprintf(msg, "[%lu] %s", offset, "Missing a closing quotation mark in string.");
                break;
            case kParseErrorStringInvalidEncoding:
                sprintf(msg, "[%lu] %s", offset, "Invalid encoding in string.");
                break;
            case kParseErrorNumberTooBig:
                sprintf(msg, "[%lu] %s", offset, "Number too big to be stored in double.");
                break;
            case kParseErrorNumberMissFraction:
                sprintf(msg, "[%lu] %s", offset, "Miss fraction part in number.");
                break;
            case kParseErrorNumberMissExponent:
                sprintf(msg, "[%lu] %s", offset, "Miss exponent in number.");
                break;
            case kParseErrorTermination:
                sprintf(msg, "[%lu] %s", offset, "Parsing was terminated.");
                break;
            case kParseErrorUnspecificSyntaxError:
                sprintf(msg, "[%lu] %s", offset, "Unspecific syntax error.");
                break;
            default:
                sprintf(msg, "[%lu] %s", offset, "Unrecognized error code.");
                break;
        }

        *error = msg;
        delete doc;
        return NULL;
    }

    return doc;
}

void DeleteDocument(void *value) {
    Document *d = static_cast<Document*>(value);
    delete d;
}

void* Get(void* value, Path* path) {
    if (!value) {
        return NULL;
    }
    if (!path) {
        return static_cast<Value*>(value);
    }

    int i;
    Value* val = static_cast<Value*>(value);
    for (i = 0; i < path->length; i++) {
        char* key = path->keys[i];
        Value::MemberIterator itr = val->FindMember(key);
        if (itr == val->MemberEnd()) {
            return NULL;
        }
        val = &(itr->value);
    }
    return val;
}

int GetInt(void* value, Path* path) {
    Value* val = (Value*) Get(value, path);
    if (!val || !val->IsInt()) {
        errno = EINVAL;
        return -1;
    }
    return val->GetInt();
}

bool GetBool(void* value, Path* path) {
    Value* val = (Value*) Get(value, path);
    if (!val || !val->IsBool()) {
        errno = EINVAL;
        return -1;
    }
    return val->GetBool();
}

double GetDouble(void* value, Path* path) {
    Value* val = (Value*) Get(value, path);
    if (!val || !val->IsDouble()) {
        errno = EINVAL;
        return -1;
    }
    return val->GetDouble();
}

unsigned int GetUInt(void* value, Path* path) {
    Value* val = (Value*) Get(value, path);
    if (!val || !val->IsUint()) {
        errno = EINVAL;
        return -1;
    }
    return val->GetUint();
}

const char* GetString(void* value, Path* path) {
    Value* val = (Value*) Get(value, path);
    if (!val || !val->IsString()) {
        return NULL;
    }
    return val->GetString();
}

void** GetArray(void* value, Path* path, size_t* length) {
    Value* val = (Value*) Get(value, path);
    if (!val || !val->IsArray()) {
        return NULL;
    }

    SizeType len = val->Size();
    void** array = (void**) malloc(sizeof(void *) * len);
    Value::ValueIterator itr = val->Begin();
    for (int i = 0; itr != val->End(); i++, itr++) {
        array[i] = (void*) itr;
    }

    *length = len;
    return array;
}

void** GetObject(void* value, Path* path, size_t* length, char*** keys) {
    Value* val = (Value*) Get(value, path);
    if (!val || !val->IsObject()) {
        return NULL;
    }

    SizeType len = val->MemberCount();
    char** names = (char**) malloc(sizeof(char *) * len);
    void** members = (void**) malloc(sizeof(void *) * len);

    Value::MemberIterator itr = val->MemberBegin();
    for (int i = 0; itr != val->MemberEnd(); i++, itr++) {
        names[i] = (char*) itr->name.GetString();
        members[i] = (void*) &(itr->value);
    }

    *length = len;
    *keys = names;
    return members;
}

const char* Type(void *value, Path* path) {
    Value* val = (Value*) Get(value, path);
    if (!val) {
        return NULL;
    }
    switch (val->GetType()) {
        case kTrueType:
        case kFalseType:
            return "BOOL";
        case kNullType:
            return "NULL";
        case kArrayType:
            return "ARRAY";
        case kStringType:
            return "STRING";
        case kObjectType:
            return "OBJECT";
        case kNumberType:
            return "NUMBER";
    }
    return "UNKNOWN";
}

char* Stringify(void *value, Path* path) {
    Value* val = (Value*) Get(value, path);
    if (!val) {
        return NULL;
    }
    StringBuffer buffer;
    Writer<StringBuffer> writer(buffer);
    val->Accept(writer);

    const char* str = buffer.GetString();
    char* retstr = (char*) malloc(sizeof(char) * (strlen(str) + 1));
    strcpy(retstr, str);
    return retstr;
}
