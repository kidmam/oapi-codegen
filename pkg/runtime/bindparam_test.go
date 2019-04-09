package runtime

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitParameter(t *testing.T) {
	// Please see the parameter serialization docs to understand these test
	// cases

	expectedPrimitive := []string{"5"}
	expectedArray := []string{"3", "4", "5"}
	expectedObject := []string{"role", "admin", "firstName", "Alex"}
	expectedExplodedObject := []string{"role=admin", "firstName=Alex"}

	var result []string
	var err error
	//  ------------------------ simple style ---------------------------------
	result, err = splitParameter("simple",
		false,
		false,
		"id",
		"5")
	assert.NoError(t, err)
	assert.EqualValues(t, expectedPrimitive, result)

	result, err = splitParameter("simple",
		false,
		false,
		"id",
		"3,4,5")
	assert.NoError(t, err)
	assert.EqualValues(t, expectedArray, result)

	result, err = splitParameter("simple",
		false,
		true,
		"id",
		"role,admin,firstName,Alex")
	assert.NoError(t, err)
	assert.EqualValues(t, expectedObject, result)

	result, err = splitParameter("simple",
		true,
		false,
		"id",
		"5")
	assert.NoError(t, err)
	assert.EqualValues(t, expectedPrimitive, result)

	result, err = splitParameter("simple",
		true,
		false,
		"id",
		"3,4,5")
	assert.NoError(t, err)
	assert.EqualValues(t, expectedArray, result)

	result, err = splitParameter("simple",
		true,
		true,
		"id",
		"role=admin,firstName=Alex")
	assert.NoError(t, err)
	assert.EqualValues(t, expectedExplodedObject, result)

	//  ------------------------ label style ---------------------------------
	result, err = splitParameter("label",
		false,
		false,
		"id",
		".5")
	assert.NoError(t, err)
	assert.EqualValues(t, expectedPrimitive, result)

	result, err = splitParameter("label",
		false,
		false,
		"id",
		".3,4,5")
	assert.NoError(t, err)
	assert.EqualValues(t, expectedArray, result)

	result, err = splitParameter("label",
		false,
		true,
		"id",
		".role,admin,firstName,Alex")
	assert.NoError(t, err)
	assert.EqualValues(t, expectedObject, result)

	result, err = splitParameter("label",
		true,
		false,
		"id",
		".5")
	assert.NoError(t, err)
	assert.EqualValues(t, expectedPrimitive, result)

	result, err = splitParameter("label",
		true,
		false,
		"id",
		".3.4.5")
	assert.NoError(t, err)
	assert.EqualValues(t, expectedArray, result)

	result, err = splitParameter("label",
		true,
		true,
		"id",
		".role=admin.firstName=Alex")
	assert.NoError(t, err)
	assert.EqualValues(t, expectedExplodedObject, result)

	//  ------------------------ matrix style ---------------------------------
	result, err = splitParameter("matrix",
		false,
		false,
		"id",
		";id=5")
	assert.NoError(t, err)
	assert.EqualValues(t, expectedPrimitive, result)

	result, err = splitParameter("matrix",
		false,
		false,
		"id",
		";id=3,4,5")
	assert.NoError(t, err)
	assert.EqualValues(t, expectedArray, result)

	result, err = splitParameter("matrix",
		false,
		true,
		"id",
		";id=role,admin,firstName,Alex")
	assert.NoError(t, err)
	assert.EqualValues(t, expectedObject, result)

	result, err = splitParameter("matrix",
		true,
		false,
		"id",
		";id=5")
	assert.NoError(t, err)
	assert.EqualValues(t, expectedPrimitive, result)

	result, err = splitParameter("matrix",
		true,
		false,
		"id",
		";id=3;id=4;id=5")
	assert.NoError(t, err)
	assert.EqualValues(t, expectedArray, result)

	result, err = splitParameter("matrix",
		true,
		true,
		"id",
		";role=admin;firstName=Alex")
	assert.NoError(t, err)
	assert.EqualValues(t, expectedExplodedObject, result)

	// ------------------------ form style ---------------------------------
	result, err = splitParameter("form",
		false,
		false,
		"id",
		"id=5")
	assert.NoError(t, err)
	assert.EqualValues(t, expectedPrimitive, result)

	result, err = splitParameter("form",
		false,
		false,
		"id",
		"id=3,4,5")
	assert.NoError(t, err)
	assert.EqualValues(t, expectedArray, result)

	result, err = splitParameter("form",
		false,
		true,
		"id",
		"id=role,admin,firstName,Alex")
	assert.NoError(t, err)
	assert.EqualValues(t, expectedObject, result)

	result, err = splitParameter("form",
		true,
		false,
		"id",
		"id=5")
	assert.NoError(t, err)
	assert.EqualValues(t, expectedPrimitive, result)

	result, err = splitParameter("form",
		true,
		false,
		"id",
		"id=3&id=4&id=5")
	assert.NoError(t, err)
	assert.EqualValues(t, expectedArray, result)

	result, err = splitParameter("form",
		true,
		true,
		"id",
		"role=admin&firstName=Alex")
	assert.NoError(t, err)
	assert.EqualValues(t, expectedExplodedObject, result)
}
