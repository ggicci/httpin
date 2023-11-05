package internal

// func TestInterfaceEncoder_fallbackEncoder(t *testing.T) {
// 	loc := &Location{
// 		Latitude:  1.234,
// 		Longitude: 5.678,
// 	}

// 	actual, err := fallbackEncoder.Encode(reflect.ValueOf(loc))
// 	assert.NoError(t, err)
// 	assert.Equal(t, "1.234000,5.678000", actual)

// 	actual, err = fallbackEncoder.Encode(reflect.ValueOf(LocationFormValueMarshalerImpl(*loc)))
// 	assert.NoError(t, err)
// 	assert.Equal(t, "HttpinFormValue:1.234000,5.678000", actual)

// 	actual, err = fallbackEncoder.Encode(reflect.ValueOf(LocationTextMarshalerImpl(*loc)))
// 	assert.NoError(t, err)
// 	assert.Equal(t, "MarshalText:1.234000,5.678000", actual)
// }
