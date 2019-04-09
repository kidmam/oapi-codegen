package runtime

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// This function binds a parameter as described in the Path Parameters
// section here to a Go object:
// https://swagger.io/docs/specification/serialization/
func BindStyledParameter(style string, explode bool, paramName string,
	value string, dest interface{}) error {

	if value == "" {
		return fmt.Errorf("parameter '%s' is empty, can't bind its value", paramName)
	}

	// Everything comes in by pointer, dereference it
	v := reflect.Indirect(reflect.ValueOf(dest))

	// This is the basic type of the destination object.
	t := v.Type()

	if t.Kind() == reflect.Struct {
		// We've got a destination object, we'll create a JSON representation
		// of the input value, and let the json library deal with the unmarshaling
		parts, err := splitParameter(style, explode, true, paramName, value)
		if err != nil {
			return fmt.Errorf("error splitting input '%s' into parts: %s", value, err)
		}

		var fields []string
		if explode {
			fields = make([]string, len(parts))
			for i, property := range parts {
				propertyParts := strings.Split(property, "=")
				if len(propertyParts) != 2 {
					return fmt.Errorf("parameter '%s' has invalid exploded format", paramName)
				}
				fields[i] = "\"" + propertyParts[0] + "\":\"" + propertyParts[1] + "\""
			}
		} else {
			if len(parts)%2 != 0 {
				return fmt.Errorf("parameter '%s' has invalid format, property/values need to be pairs", paramName)
			}
			fields = make([]string, len(parts)/2)
			for i := 0; i < len(parts); i += 2 {
				key := parts[i]
				value := parts[i+1]
				fields[i/2] = "\"" + key + "\":\"" + value + "\""
			}
		}
		jsonParam := "{" + strings.Join(fields, ",") + "}"
		err = json.Unmarshal([]byte(jsonParam), dest)
		if err != nil {
			return fmt.Errorf("error binding parameter %s fields: %s", paramName, err)
		}
		return nil
	}

	if t.Kind() == reflect.Slice {
		// Chop up the parameter into parts based on its style
		parts, err := splitParameter(style, explode, false, paramName, value)
		if err != nil {
			return fmt.Errorf("error splitting input '%s' into parts: %s", value, err)
		}

		// We've got a destination array, bind each object one by one.
		// This generates a slice of the correct element type and length to
		// hold all the parts.
		newArray := reflect.MakeSlice(t, len(parts), len(parts))
		for i, p := range parts {
			err = BindStringToObject(p, newArray.Index(i).Addr().Interface())
			if err != nil {
				return fmt.Errorf("error setting array element: %s", err)
			}
		}
		v.Set(newArray)
		return nil
	}

	// Try to bind the remaining types as a base type.
	return BindStringToObject(value, dest)
}

func splitParameter(style string, explode bool, object bool, paramName string, value string) ([]string, error) {
	switch style {
	case "simple":
		// In the simple case, we always split on comma
		parts := strings.Split(value, ",")
		return parts, nil
	case "label":
		// In the label case, it's more tricky. In the no explode case, we have
		// /users/.3,4,5 for arrays
		// /users/.role,admin,firstName,Alex for objects
		// in the explode case, we have:
		// /users/.3.4.5
		// /users/.role=admin.firstName=Alex
		if explode {
			// In the exploded case, split everything on periods.
			parts := strings.Split(value, ".")
			// The first part should be an empty string because we have a
			// leading period.
			if parts[0] != "" {
				return nil, fmt.Errorf("invalid format for label parameter '%s', should start with '.'", paramName)
			}
			return parts[1:], nil

		} else {
			// In the unexploded case, we strip off the leading period.
			if value[0] != '.' {
				return nil, fmt.Errorf("invalid format for label parameter '%s', should start with '.'", paramName)
			}
			// The rest is comma separated.
			return strings.Split(value[1:], ","), nil
		}

	case "matrix":
		if explode {
			// In the exploded case, we break everything up on semicolon
			parts := strings.Split(value, ";")
			// The first part should always be empty string, since we started
			// with ;something
			if parts[0] != "" {
				return nil, fmt.Errorf("invalid format for matrix parameter '%s', should start with ';'", paramName)
			}
			parts = parts[1:]
			// Now, if we have an object, we just have a list of x=y statements.
			// for a non-object, like an array, we have id=x, id=y. id=z, etc,
			// so we need to strip the prefix from each of them.
			if !object {
				prefix := paramName + "="
				for i := range parts {
					parts[i] = strings.TrimPrefix(parts[i], prefix)
				}
			}
			return parts, nil
		} else {
			// In the unexploded case, parameters will start with ;paramName=
			prefix := ";" + paramName + "="
			if !strings.HasPrefix(value, prefix) {
				return nil, fmt.Errorf("expected parameter '%s' to start with %s", paramName, prefix)
			}
			str := strings.TrimPrefix(value, prefix)
			return strings.Split(str, ","), nil
		}
	case "form":
		var parts []string
		if explode {
			parts = strings.Split(value, "&")
			if !object {
				prefix := paramName + "="
				for i := range parts {
					parts[i] = strings.TrimPrefix(parts[i], prefix)
				}
			}
			return parts, nil
		} else {
			parts = strings.Split(value, ",")
			prefix := paramName + "="
			for i := range parts {
				parts[i] = strings.TrimPrefix(parts[i], prefix)
			}
		}
		return parts, nil
	}

	return nil, fmt.Errorf("unhandled parameter style: %s", style)
}
