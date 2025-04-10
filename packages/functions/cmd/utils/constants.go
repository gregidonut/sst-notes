package utils

import "os"

var DYNAMODB_TABLE_NAME string = os.Getenv("NOTES_TABLE_NAME")
