package utils

import "database/sql"

func GetStringValue(ns sql.NullString) string {
    if ns.Valid {
        return ns.String
    }
    return ""
}

func GetNullString(s string) sql.NullString {
    if s == "" {
        return sql.NullString{Valid: false}
    }
    return sql.NullString{String: s, Valid: true}
}