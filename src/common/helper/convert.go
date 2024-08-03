package helper

func ConvertToStringPointer(val string) *string {
	if val == "" {
		return nil
	}

	return &val
}
