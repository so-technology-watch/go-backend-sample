package dao

// Get NewTaskDAOMock
func GetDAO() TaskDAO {
	return NewTaskDAOMock()
}
