package entity

type User struct {
	Name , Password , Email , Phone string
}

func (m_User User) init(t_Name,t_PW,t_Email,t_Phone string)  {
	m_User.Name = t_Name
	m_User.Password = t_PW
	m_User.Email = t_Email
	m_User.Phone = t_Phone
}