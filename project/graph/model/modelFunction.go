// This file is used to define functions for model
// to satisfy interface requirement

package model

func (m *User) GetMM() *MongoModel {
	return m.Model
}
func (m *User) SetMM(mm *MongoModel) {
	m.Model = mm
}

func (m *Todo) GetMM() *MongoModel {
	return m.Model
}
func (m *Todo) SetMM(mm *MongoModel) {
	m.Model = mm
}
