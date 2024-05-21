// This file is used to define functions for model
// to satisfy interface requirement

package model

func getIdFromModel(m *MongoModel) string {
	if m == nil || m.ID == nil {
		return ""
	}
	return *m.ID
}

func (m *User) GetId() string {
	return getIdFromModel(m.Model)
}

func (m *Todo) GetId() string {
	return getIdFromModel(m.Model)
}
