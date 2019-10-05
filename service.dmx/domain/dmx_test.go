package domain

//func TestUnmarshalFixture(t *testing.T) {
//	data := []byte(`{
//		"id": "test_fixture",
//		"name": "Test Fixture",
//		"type": "DMX",
//		"kind": "lamp",
//		"attributes": {
//			"fixture_type": "mega_par_profile",
//			"offset": 8
//		},
//		"room_id": "test_room",
//		"controller_name": "service.dmx"
//	}`)
//
//	var fix *Fixture
//	err := json.Unmarshal(data, &fix)
//	assert.NilError(t, err)
//
//	assert.Equal(t, "test_fixture", fix.ID)
//	assert.Equal(t, "Test Fixture", fix.Name)
//	assert.Equal(t, "DMX", fix.Type)
//	assert.Equal(t, "lamp", fix.Kind)
//	assert.Equal(t, "mega_par_profile", fix.Attributes.FixtureType)
//	assert.Equal(t, 8, fix.Attributes.Offset)
//	assert.Equal(t, "test_room", fix.RoomID)
//	assert.Equal(t, "service.dmx", fix.ControllerName)
//}
