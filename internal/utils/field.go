package utils

func GetStringField[T any](obj *T, getter func(*T) string) string {
	if obj != nil {
		return getter(obj)
	}
	return ""
}

func GetInt64Field[T any](obj *T, getter func(*T) int64) int64 {
	if obj != nil {
		return getter(obj)
	}
	return 0
}

func GetBoolField[T any](obj *T, getter func(*T) bool) bool {
	if obj != nil {
		return getter(obj)
	}
	return false
}

func MergeStringField(primary, fallback string) string {
	if primary != "" {
		return primary
	}
	return fallback
}

func MergeInt64Field(primary, fallback int64) int64 {
	if primary > 0 {
		return primary
	}
	return fallback
}

func MergeBoolField(primary, fallback bool) bool {
	return primary || fallback
}
